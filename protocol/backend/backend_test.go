package backend

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/noot/atomic-swap/common"
	"github.com/noot/atomic-swap/tests"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestWaitForReceipt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ec, err := ethclient.Dial(common.DefaultEthEndpoint)
	require.NoError(t, err)
	defer ec.Close()

	privKey, err := ethcrypto.HexToECDSA(tests.GetTakerTestKey(t))
	require.NoError(t, err)

	publicKey := privKey.Public().(*ecdsa.PublicKey)

	nonce, err := ec.PendingNonceAt(ctx, ethcrypto.PubkeyToAddress(*publicKey))
	require.NoError(t, err)

	to := ethcommon.Address{}
	txInner := &ethtypes.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    big.NewInt(99),
		Gas:      21000,
		GasPrice: big.NewInt(2000000000),
	}

	tx, err := ethtypes.SignNewTx(privKey,
		ethtypes.LatestSignerForChainID(big.NewInt(common.GanacheChainID)),
		txInner,
	)
	require.NoError(t, err)

	err = ec.SendTransaction(ctx, tx)
	require.NoError(t, err)

	b := &backend{
		ethClient: ec,
	}

	receipt, err := b.WaitForReceipt(ctx, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, tx.Hash(), receipt.TxHash)
}
