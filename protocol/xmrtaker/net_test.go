package xmrtaker

import (
	"testing"

	"github.com/noot/atomic-swap/common"
	"github.com/noot/atomic-swap/common/types"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func newTestXMRTaker(t *testing.T, ec *ethclient.Client) *Instance {
	b := newBackend(t, ec)
	cfg := &Config{
		Backend:  b,
		Basepath: "/tmp/xmrtaker",
	}

	xmrtaker, err := NewInstance(cfg)
	require.NoError(t, err)
	return xmrtaker
}

func TestXMRTaker_InitiateProtocol(t *testing.T) {
	ec, err := ethclient.Dial(common.DefaultEthEndpoint)
	require.NoError(t, err)
	defer ec.Close()

	a := newTestXMRTaker(t, ec)
	offer := &types.Offer{
		ExchangeRate: 1,
	}
	s, err := a.InitiateProtocol(3.33, offer)
	require.NoError(t, err)
	require.Equal(t, a.swapStates[offer.GetID()], s)
}
