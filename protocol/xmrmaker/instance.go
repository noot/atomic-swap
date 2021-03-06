package xmrmaker

import (
	"sync"

	"github.com/noot/atomic-swap/common"
	"github.com/noot/atomic-swap/common/types"
	"github.com/noot/atomic-swap/protocol/backend"

	logging "github.com/ipfs/go-log"
)

var (
	log = logging.Logger("xmrmaker")
)

// Instance implements the functionality that will be needed by a user who owns XMR
// and wishes to swap for ETH.
type Instance struct {
	backend  backend.Backend
	basepath string

	walletFile, walletPassword string

	offerManager *offerManager

	swapMu     sync.Mutex
	swapStates map[types.Hash]*swapState
}

// Config contains the configuration values for a new XMRMaker instance.
type Config struct {
	Backend                    backend.Backend
	Basepath                   string
	WalletFile, WalletPassword string
}

// NewInstance returns a new *xmrmaker.Instance.
// It accepts an endpoint to a monero-wallet-rpc instance where account 0 contains XMRMaker's XMR.
func NewInstance(cfg *Config) (*Instance, error) {
	if cfg.WalletFile != "" {
		if err := cfg.Backend.OpenWallet(cfg.WalletFile, cfg.WalletPassword); err != nil {
			return nil, err
		}
	} else {
		log.Warn("monero wallet-file not set; must be set via RPC call personal_setMoneroWalletFile before making an offer")
	}

	return &Instance{
		backend:        cfg.Backend,
		basepath:       cfg.Basepath,
		walletFile:     cfg.WalletFile,
		walletPassword: cfg.WalletPassword,
		offerManager:   newOfferManager(cfg.Basepath),
		swapStates:     make(map[types.Hash]*swapState),
	}, nil
}

// SetMoneroWalletFile sets the Instance's current monero wallet file.
func (b *Instance) SetMoneroWalletFile(file, password string) error {
	_ = b.backend.CloseWallet()
	return b.backend.OpenWallet(file, password)
}

func (b *Instance) openWallet() error { //nolint
	return b.backend.OpenWallet(b.walletFile, b.walletPassword)
}

// GetOngoingSwapState ...
func (b *Instance) GetOngoingSwapState(id types.Hash) common.SwapState {
	return b.swapStates[id]
}
