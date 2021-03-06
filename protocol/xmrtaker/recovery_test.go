package xmrtaker

import (
	"path"
	"testing"
	"time"

	"github.com/noot/atomic-swap/common"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

func newTestRecoveryState(t *testing.T) *recoveryState {
	s := newTestInstance(t)
	s.SetSwapTimeout(time.Second * 10)
	akp, err := generateKeys()
	require.NoError(t, err)

	s.privkeys = akp.PrivateKeyPair
	s.pubkeys = akp.PublicKeyPair
	s.secp256k1Pub = akp.Secp256k1PublicKey
	s.dleqProof = akp.DLEqProof

	s.setXMRMakerKeys(s.pubkeys.SpendKey(), s.privkeys.ViewKey(), akp.Secp256k1PublicKey)
	s.xmrmakerAddress = s.EthAddress()

	_, err = s.lockETH(common.NewEtherAmount(1))
	require.NoError(t, err)

	basePath := path.Join(t.TempDir(), "test-infoFile")
	rs, err := NewRecoveryState(s, basePath, s.privkeys.SpendKey(), s.contractSwapID, s.contractSwap)
	require.NoError(t, err)
	return rs
}

func TestClaimOrRefund_Claim(t *testing.T) {
	// test case where XMRMaker has claimed the ether, so XMRTaker should be able to
	// claim the monero.
	rs := newTestRecoveryState(t)

	// call swap.Ready()
	err := rs.ss.ready()
	require.NoError(t, err)

	// call swap.Claim()
	sc := rs.ss.getSecret()
	txOpts, err := rs.ss.TxOpts()
	require.NoError(t, err)

	_, err = rs.ss.Contract().Claim(txOpts, rs.ss.contractSwap, sc)
	require.NoError(t, err)

	t.Log("XMRMaker claimed ETH...")

	// assert we can claim the monero
	res, err := rs.ClaimOrRefund()
	require.NoError(t, err)
	require.True(t, res.Claimed)
}

func TestClaimOrRefund_Refund_beforeT0(t *testing.T) {
	// test case where XMRMaker hasn't claimed the ether, and it's before
	// t0/IsReady, so XMRTaker should be able to refund.
	rs := newTestRecoveryState(t)

	// assert we can refund the ether
	res, err := rs.ClaimOrRefund()
	require.NoError(t, err)
	require.True(t, res.Refunded)
}

func TestClaimOrRefund_Refund_afterT1(t *testing.T) {
	// test case where XMRMaker hasn't claimed the ether, and it's after
	// t1, so XMRTaker should be able to refund.
	rs := newTestRecoveryState(t)

	rpcClient, err := rpc.Dial(common.DefaultEthEndpoint)
	require.NoError(t, err)

	var result string
	err = rpcClient.Call(&result, "evm_snapshot")
	require.NoError(t, err)

	err = rpcClient.Call(nil, "evm_increaseTime", rs.ss.SwapTimeout().Seconds()*2+360)
	require.NoError(t, err)

	defer func() {
		var ok bool
		err = rpcClient.Call(&ok, "evm_revert", result)
		require.NoError(t, err)
	}()

	// assert we can refund the ether
	res, err := rs.ClaimOrRefund()
	require.NoError(t, err)
	require.True(t, res.Refunded)
}
