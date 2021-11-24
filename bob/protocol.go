package bob

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/noot/atomic-swap/common"
	"github.com/noot/atomic-swap/monero"
	"github.com/noot/atomic-swap/net"
	"github.com/noot/atomic-swap/swap-contract"

	logging "github.com/ipfs/go-log"
)

var (
	log = logging.Logger("bob")
)

// bob implements the functions that will be called by a user who owns XMR
// and wishes to swap for ETH.
type bob struct {
	ctx context.Context

	client                     monero.Client
	daemonClient               monero.DaemonClient
	walletFile, walletPassword string

	ethClient  *ethclient.Client
	ethPrivKey *ecdsa.PrivateKey
	auth       *bind.TransactOpts
	callOpts   *bind.CallOpts
	ethAddress ethcommon.Address

	net net.MessageSender

	swapMu    sync.Mutex
	swapState *swapState
}

// NewBob returns a new instance of Bob.
// It accepts an endpoint to a monero-wallet-rpc instance where account 0 contains Bob's XMR.
func NewBob(ctx context.Context, moneroEndpoint, moneroDaemonEndpoint, ethEndpoint, ethPrivKey, walletFile, walletPassword string) (*bob, error) {
	pk, err := crypto.HexToECDSA(ethPrivKey)
	if err != nil {
		return nil, err
	}

	ec, err := ethclient.Dial(ethEndpoint)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(1337)) // ganache chainID
	if err != nil {
		return nil, err
	}

	pub := pk.Public().(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*pub)

	// monero-wallet-rpc client
	walletClient := monero.NewClient(moneroEndpoint)

	// open Bob's XMR wallet
	if err = walletClient.OpenWallet(walletFile, walletPassword); err != nil {
		return nil, err
	}

	return &bob{
		ctx:            ctx,
		client:         walletClient,
		daemonClient:   monero.NewClient(moneroDaemonEndpoint),
		walletFile:     walletFile,
		walletPassword: walletPassword,
		ethClient:      ec,
		ethPrivKey:     pk,
		auth:           auth,
		callOpts: &bind.CallOpts{
			From:    addr,
			Context: ctx,
		},
		ethAddress: addr,
	}, nil
}

func (b *bob) SetMessageSender(n net.MessageSender) {
	b.net = n
}

func (b *bob) openWallet() error { //nolint
	return b.client.OpenWallet(b.walletFile, b.walletPassword)
}

// generateKeys generates Bob's spend and view keys (s_b, v_b)
// It returns Bob's public spend key and his private view key, so that Alice can see
// if the funds are locked.
func (s *swapState) generateKeys() (*monero.PublicKey, *monero.PrivateViewKey, error) {
	if s.privkeys != nil {
		return s.pubkeys.SpendKey(), s.privkeys.ViewKey(), nil
	}

	var err error
	s.privkeys, err = monero.GenerateKeys()
	if err != nil {
		return nil, nil, err
	}

	// TODO: configure basepath
	// TODO: write swap ID
	if err := common.WriteKeysToFile("/tmp/bob-xmr", s.privkeys); err != nil {
		return nil, nil, err
	}

	s.pubkeys = s.privkeys.PublicKeyPair()
	return s.pubkeys.SpendKey(), s.privkeys.ViewKey(), nil
}

// setAlicePublicKeys sets Alice's public spend and view keys
func (s *swapState) setAlicePublicKeys(sk *monero.PublicKeyPair) {
	s.alicePublicKeys = sk
}

// setContract sets the contract in which Alice has locked her ETH.
func (s *swapState) setContract(address ethcommon.Address) error {
	var err error
	s.contractAddr = address
	s.contract, err = swap.NewSwap(address, s.bob.ethClient)
	return err
}

// watchForReady watches for Alice to call Ready() on the swap contract, allowing
// Bob to call Claim().
func (s *swapState) watchForReady() (<-chan struct{}, error) { //nolint:unused
	watchOpts := &bind.WatchOpts{
		Context: s.ctx,
	}

	done := make(chan struct{})
	ch := make(chan *swap.SwapIsReady)
	defer close(done)

	// watch for Refund() event on chain, calculate unlock key as result
	sub, err := s.contract.WatchIsReady(watchOpts, ch)
	if err != nil {
		return nil, err
	}

	defer sub.Unsubscribe()

	go func() {
		for {
			select {
			case event := <-ch:
				if !event.B {
					continue
				}

				// contract is ready!!
				close(done)
			case <-s.ctx.Done():
				return
			}
		}
	}()

	return done, nil
}

// watchForRefund watches for the Refund event in the contract.
// This should be called before LockFunds.
// If a keypair is sent over this channel, the rest of the protocol should be aborted.
//
// If Alice chooses to refund and thus reveals s_a,
// the private spend and view keys that contain the previously locked monero
// ((s_a + s_b), (v_a + v_b)) are sent over the channel.
// Bob can then use these keys to move his funds if he wishes.
func (s *swapState) watchForRefund() (<-chan *monero.PrivateKeyPair, error) { //nolint:unused
	watchOpts := &bind.WatchOpts{
		Context: s.ctx,
	}

	out := make(chan *monero.PrivateKeyPair)
	ch := make(chan *swap.SwapRefunded)
	defer close(out)

	// watch for Refund() event on chain, calculate unlock key as result
	sub, err := s.contract.WatchRefunded(watchOpts, ch)
	if err != nil {
		return nil, err
	}

	defer sub.Unsubscribe()

	go func() {
		for {
			select {
			case refund := <-ch:
				if refund == nil {
					continue
				}

				// got Alice's secret
				sa := refund.S
				skA, err := monero.NewPrivateSpendKey(sa[:])
				if err != nil {
					log.Info("failed to convert Alice's secret into a key: %w", err)
					return
				}

				vkA, err := skA.View()
				if err != nil {
					log.Info("failed to get view key from Alice's secret spend key: %w", err)
					return
				}

				skAB := monero.SumPrivateSpendKeys(skA, s.privkeys.SpendKey())
				vkAB := monero.SumPrivateViewKeys(vkA, s.privkeys.ViewKey())
				kpAB := monero.NewPrivateKeyPair(skAB, vkAB)
				out <- kpAB
			case <-s.ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

// lockFunds locks Bob's funds in the monero account specified by public key
// (S_a + S_b), viewable with (V_a + V_b)
// It accepts the amount to lock as the input
// TODO: units
func (s *swapState) lockFunds(amount uint64) (monero.Address, error) {
	kp := monero.SumSpendAndViewKeys(s.alicePublicKeys, s.pubkeys)

	log.Debug("public spend keys for lock account: ", kp.SpendKey().Hex())
	log.Debug("public view keys for lock account: ", kp.ViewKey().Hex())
	log.Infof("going to lock XMR funds, amount=%d", amount)

	balance, err := s.bob.client.GetBalance(0)
	if err != nil {
		return "", err
	}

	log.Debug("XMR balance: ", balance.Balance)
	log.Debug("unlocked XMR balance: ", balance.UnlockedBalance)
	log.Debug("blocks to unlock: ", balance.BlocksToUnlock)

	address := kp.Address()
	if err := s.bob.client.Transfer(address, 0, uint(amount)); err != nil {
		return "", err
	}

	bobAddr, err := s.bob.client.GetAddress(0)
	if err != nil {
		return "", err
	}

	// TODO: this has gonna be removed before stagenet/mainnet, will need to add
	// waiting logs or such
	if err := s.bob.daemonClient.GenerateBlocks(bobAddr.Address, 1); err != nil {
		return "", err
	}

	if err := s.bob.client.Refresh(); err != nil {
		return "", err
	}

	log.Infof("successfully locked XMR funds: address=%s", address)
	return address, nil
}

// claimFunds redeems Bob's ETH funds by calling Claim() on the contract
func (s *swapState) claimFunds() (string, error) {
	pub := s.ethPrivKey.Public().(*ecdsa.PublicKey)
	addr := ethcrypto.PubkeyToAddress(*pub)

	balance, err := s.ethClient.BalanceAt(s.ctx, addr, nil)
	if err != nil {
		return "", err
	}

	log.Info("Bob's balance before claim: ", balance)

	// call swap.Swap.Claim() w/ b.privkeys.sk, revealing Bob's secret spend key
	secret := s.privkeys.SpendKeyBytes()
	var sc [32]byte
	copy(sc[:], common.Reverse(secret))

	tx, err := s.contract.Claim(s.bob.auth, sc)
	if err != nil {
		return "", err
	}

	log.Infof("success! Bob claimed funds, tx hash=%s", tx.Hash())

	receipt, err := s.bob.ethClient.TransactionReceipt(s.ctx, tx.Hash())
	if err != nil {
		return "", err
	}

	log.Infof("included in block number=%d gas used=%d s_a=%x",
		receipt.Logs[0].BlockNumber,
		receipt.CumulativeGasUsed,
		secret,
	)

	balance, err = s.bob.ethClient.BalanceAt(s.ctx, addr, nil)
	if err != nil {
		return "", err
	}

	log.Info("Bob's balance after claim: ", balance)
	s.success = true
	return tx.Hash().String(), nil
}
