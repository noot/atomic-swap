package tests

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/noot/atomic-swap/common"
	"github.com/noot/atomic-swap/common/types"
	"github.com/noot/atomic-swap/monero"
	"github.com/noot/atomic-swap/rpcclient"
	"github.com/noot/atomic-swap/rpcclient/wsclient"

	"github.com/stretchr/testify/require"
)

const (
	testsEnv          = "TESTS"
	integrationMode   = "integration"
	generateBlocksEnv = "GENERATEBLOCKS"
	falseStr          = "false"

	defaultXMRTakerDaemonEndpoint   = "http://localhost:5001"
	defaultXMRTakerDaemonWSEndpoint = "ws://localhost:8081"
	defaultXMRMakerDaemonEndpoint   = "http://localhost:5002"
	defaultXMRMakerDaemonWSEndpoint = "ws://localhost:8082"
	defaultCharlieDaemonWSEndpoint  = "ws://localhost:8083"

	defaultDiscoverTimeout = 2 // 2 seconds

	xmrmakerProvideAmount = float64(1.0)
	exchangeRate          = float64(0.05)
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(0)
	}

	if os.Getenv(testsEnv) != integrationMode {
		os.Exit(0)
	}

	if os.Getenv(generateBlocksEnv) != falseStr {
		generateBlocks(512)
	}

	go generateBlocksAsync()

	os.Exit(m.Run())
}

func generateBlocks(num uint) {
	c := monero.NewClient(common.DefaultXMRMakerMoneroEndpoint)
	d := monero.NewDaemonClient(common.DefaultMoneroDaemonEndpoint)
	xmrmakerAddr, err := c.GetAddress(0)
	if err != nil {
		panic(err)
	}

	fmt.Println("> Generating blocks for test setup...")
	_ = d.GenerateBlocks(xmrmakerAddr.Address, num)
	err = c.Refresh()
	if err != nil {
		panic(err)
	}

	fmt.Println("> Completed generating blocks.")
}

func generateBlocksAsync() {
	c := monero.NewClient(common.DefaultXMRMakerMoneroEndpoint)
	d := monero.NewDaemonClient(common.DefaultMoneroDaemonEndpoint)
	xmrmakerAddr, err := c.GetAddress(0)
	if err != nil {
		panic(err)
	}

	// generate 1 block per second
	for {
		time.Sleep(time.Second)
		_ = d.GenerateBlocks(xmrmakerAddr.Address, 1)
		err = c.Refresh()
		if err != nil {
			panic(err)
		}
	}
}

func TestXMRTaker_Discover(t *testing.T) {
	bc := rpcclient.NewClient(defaultXMRMakerDaemonEndpoint)
	_, err := bc.MakeOffer(xmrmakerProvideAmount, xmrmakerProvideAmount, exchangeRate)
	require.NoError(t, err)

	c := rpcclient.NewClient(defaultXMRTakerDaemonEndpoint)
	providers, err := c.Discover(types.ProvidesXMR, defaultDiscoverTimeout)
	require.NoError(t, err)
	require.Equal(t, 1, len(providers))
	require.GreaterOrEqual(t, len(providers[0]), 2)
}

func TestXMRMaker_Discover(t *testing.T) {
	c := rpcclient.NewClient(defaultXMRMakerDaemonEndpoint)
	providers, err := c.Discover(types.ProvidesETH, defaultDiscoverTimeout)
	require.NoError(t, err)
	require.Equal(t, 0, len(providers))
}

func TestXMRTaker_Query(t *testing.T) {
	bc := rpcclient.NewClient(defaultXMRMakerDaemonEndpoint)
	_, err := bc.MakeOffer(xmrmakerProvideAmount, xmrmakerProvideAmount, exchangeRate)
	require.NoError(t, err)

	c := rpcclient.NewClient(defaultXMRTakerDaemonEndpoint)

	providers, err := c.Discover(types.ProvidesXMR, defaultDiscoverTimeout)
	require.NoError(t, err)
	require.Equal(t, 1, len(providers))
	require.GreaterOrEqual(t, len(providers[0]), 2)

	resp, err := c.Query(providers[0][0])
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(resp.Offers), 1)
	require.Equal(t, xmrmakerProvideAmount, resp.Offers[0].MinimumAmount)
	require.Equal(t, xmrmakerProvideAmount, resp.Offers[0].MaximumAmount)
	require.Equal(t, exchangeRate, float64(resp.Offers[0].ExchangeRate))
}

func TestSuccess(t *testing.T) {
	const testTimeout = time.Second * 60

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bwsc, err := wsclient.NewWsClient(ctx, defaultXMRMakerDaemonWSEndpoint)
	require.NoError(t, err)

	offerID, statusCh, err := bwsc.MakeOfferAndSubscribe(0.1, xmrmakerProvideAmount,
		types.ExchangeRate(exchangeRate))
	require.NoError(t, err)

	bc := rpcclient.NewClient(defaultXMRMakerDaemonEndpoint)
	offersBefore, err := bc.GetOffers()
	require.NoError(t, err)

	errCh := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			select {
			case status := <-statusCh:
				fmt.Println("> XMRMaker got status:", status)
				if status.IsOngoing() {
					continue
				}

				if status != types.CompletedSuccess {
					errCh <- fmt.Errorf("swap did not complete successfully: got %s", status)
				}

				return
			case <-time.After(testTimeout):
				errCh <- errors.New("make offer subscription timed out")
			}
		}
	}()

	c := rpcclient.NewClient(defaultXMRTakerDaemonEndpoint)
	wsc, err := wsclient.NewWsClient(ctx, defaultXMRTakerDaemonWSEndpoint)
	require.NoError(t, err)

	// TODO: implement discovery over websockets
	providers, err := c.Discover(types.ProvidesXMR, defaultDiscoverTimeout)
	require.NoError(t, err)
	require.Equal(t, 1, len(providers))
	require.GreaterOrEqual(t, len(providers[0]), 2)

	takerStatusCh, err := wsc.TakeOfferAndSubscribe(providers[0][0], offerID, 0.05)
	require.NoError(t, err)

	go func() {
		defer wg.Done()
		for status := range takerStatusCh {
			fmt.Println("> XMRTaker got status:", status)
			if status.IsOngoing() {
				continue
			}

			if status != types.CompletedSuccess {
				errCh <- fmt.Errorf("swap did not complete successfully: got %s", status)
			}

			return
		}
	}()

	wg.Wait()

	select {
	case err = <-errCh:
		require.NoError(t, err)
	default:
	}

	offersAfter, err := bc.GetOffers()
	require.NoError(t, err)
	require.Equal(t, 1, len(offersBefore)-len(offersAfter))
}

func TestRefund_XMRTakerCancels(t *testing.T) {
	if os.Getenv(generateBlocksEnv) != falseStr {
		generateBlocks(64)
	}

	const (
		testTimeout = time.Second * 60
		swapTimeout = 5 // 5s
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bwsc, err := wsclient.NewWsClient(ctx, defaultXMRMakerDaemonWSEndpoint)
	require.NoError(t, err)

	offerID, statusCh, err := bwsc.MakeOfferAndSubscribe(0.1, xmrmakerProvideAmount,
		types.ExchangeRate(exchangeRate))
	require.NoError(t, err)

	bc := rpcclient.NewClient(defaultXMRMakerDaemonEndpoint)
	offersBefore, err := bc.GetOffers()
	require.NoError(t, err)

	errCh := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			select {
			case status := <-statusCh:
				fmt.Println("> XMRMaker got status:", status)
				if status.IsOngoing() {
					continue
				}

				if status != types.CompletedRefund {
					errCh <- fmt.Errorf("swap did not refund successfully for XMRMaker: got %s", status)
				}

				return
			case <-time.After(testTimeout):
				errCh <- errors.New("make offer subscription timed out")
			}
		}
	}()

	c := rpcclient.NewClient(defaultXMRTakerDaemonEndpoint)
	wsc, err := wsclient.NewWsClient(ctx, defaultXMRTakerDaemonWSEndpoint)
	require.NoError(t, err)

	err = c.SetSwapTimeout(swapTimeout)
	require.NoError(t, err)

	providers, err := c.Discover(types.ProvidesXMR, defaultDiscoverTimeout)
	require.NoError(t, err)
	require.Equal(t, 1, len(providers))
	require.GreaterOrEqual(t, len(providers[0]), 2)

	takerStatusCh, err := wsc.TakeOfferAndSubscribe(providers[0][0], offerID, 0.05)
	require.NoError(t, err)

	go func() {
		defer wg.Done()
		for status := range takerStatusCh {
			fmt.Println("> XMRTaker got status:", status)
			if status != types.ETHLocked {
				continue
			}

			fmt.Println("> XMRTaker cancelled swap!")
			exitStatus, err := c.Cancel(offerID) //nolint:govet
			if err != nil {
				t.Log("XMRTaker got error", err)
				errCh <- err
				return
			}

			if exitStatus != types.CompletedRefund {
				errCh <- fmt.Errorf("did not refund successfully for XMRTaker: exit status was %s", exitStatus)
			}

			return
		}
	}()

	wg.Wait()

	select {
	case err = <-errCh:
		require.NoError(t, err)
	default:
	}

	offersAfter, err := bc.GetOffers()
	require.NoError(t, err)
	require.Equal(t, len(offersBefore), len(offersAfter))
}

// TestRefund_XMRMakerCancels_untilAfterT1 tests the case where XMRTaker and XMRMaker
// both lock their funds, but XMRMaker goes offline
// until time t1 in the swap contract passes. This triggers XMRTaker to refund, which XMRMaker will then
// "come online" to see, and he will then refund also.
func TestRefund_XMRMakerCancels_untilAfterT1(t *testing.T) {
	testRefundXMRMakerCancels(t, 5, types.CompletedRefund)
	time.Sleep(time.Second * 5)
}

// TestRefund_XMRMakerCancels_afterIsReady tests the case where XMRTaker and XMRMaker both lock their funds,
// but XMRMaker goes offline until past isReady==true and t0, but comes online before t1.
//  When XMRMaker comes back online, he should claim the ETH, causing XMRTaker to also claim the XMR.
func TestRefund_XMRMakerCancels_afterIsReady(t *testing.T) {
	testRefundXMRMakerCancels(t, 30, types.CompletedSuccess)
}

func testRefundXMRMakerCancels(t *testing.T, swapTimeout uint64, expectedExitStatus types.Status) {
	if os.Getenv(generateBlocksEnv) != falseStr {
		generateBlocks(64)
	}

	const testTimeout = time.Second * 60

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bcli := rpcclient.NewClient(defaultXMRMakerDaemonEndpoint)
	bwsc, err := wsclient.NewWsClient(ctx, defaultXMRMakerDaemonWSEndpoint)
	require.NoError(t, err)

	offerID, statusCh, err := bwsc.MakeOfferAndSubscribe(0.1, xmrmakerProvideAmount,
		types.ExchangeRate(exchangeRate))
	require.NoError(t, err)

	offersBefore, err := bcli.GetOffers()
	require.NoError(t, err)

	errCh := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			select {
			case status := <-statusCh:
				fmt.Println("> XMRMaker got status:", status)
				if status != types.XMRLocked {
					continue
				}

				fmt.Println("> XMRMaker cancelled swap!")
				exitStatus, err := bcli.Cancel(offerID) //nolint:govet
				if err != nil {
					errCh <- err
					return
				}

				if exitStatus != expectedExitStatus {
					errCh <- fmt.Errorf("did not get expected exit status for XMRMaker: got %s, expected %s", exitStatus, expectedExitStatus) //nolint:lll
					return
				}

				fmt.Println("> XMRMaker refunded successfully")
				return
			case <-time.After(testTimeout):
				errCh <- errors.New("make offer subscription timed out")
			}
		}
	}()

	c := rpcclient.NewClient(defaultXMRTakerDaemonEndpoint)
	wsc, err := wsclient.NewWsClient(ctx, defaultXMRTakerDaemonWSEndpoint)
	require.NoError(t, err)

	err = c.SetSwapTimeout(swapTimeout)
	require.NoError(t, err)

	providers, err := c.Discover(types.ProvidesXMR, defaultDiscoverTimeout)
	require.NoError(t, err)
	require.Equal(t, 1, len(providers))
	require.GreaterOrEqual(t, len(providers[0]), 2)

	takerStatusCh, err := wsc.TakeOfferAndSubscribe(providers[0][0], offerID, 0.05)
	require.NoError(t, err)

	go func() {
		defer wg.Done()

		for status := range takerStatusCh {
			fmt.Println("> XMRTaker got status:", status)
			if status.IsOngoing() {
				continue
			}

			if status != expectedExitStatus {
				errCh <- fmt.Errorf("did not get expected exit status for XMRTaker: got %s, expected %s", status, expectedExitStatus) //nolint:lll
				return
			}

			fmt.Println("> XMRTaker refunded successfully")
			return
		}
	}()

	wg.Wait()

	select {
	case err = <-errCh:
		require.NoError(t, err)
	default:
	}

	offersAfter, err := bcli.GetOffers()
	require.NoError(t, err)
	if expectedExitStatus != types.CompletedSuccess {
		require.Equal(t, len(offersBefore), len(offersAfter))
	} else {
		require.Equal(t, 1, len(offersBefore)-len(offersAfter))
	}
}

// TestAbort_XMRTakerCancels tests the case where XMRTaker cancels the swap before any funds are locked.
// Both parties should abort the swap successfully.
func TestAbort_XMRTakerCancels(t *testing.T) {
	if os.Getenv(generateBlocksEnv) != falseStr {
		generateBlocks(64)
	}

	const testTimeout = time.Second * 60

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bwsc, err := wsclient.NewWsClient(ctx, defaultXMRMakerDaemonWSEndpoint)
	require.NoError(t, err)

	offerID, statusCh, err := bwsc.MakeOfferAndSubscribe(0.1, xmrmakerProvideAmount,
		types.ExchangeRate(exchangeRate))
	require.NoError(t, err)

	bc := rpcclient.NewClient(defaultXMRMakerDaemonEndpoint)
	offersBefore, err := bc.GetOffers()
	require.NoError(t, err)

	errCh := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			select {
			case status := <-statusCh:
				fmt.Println("> XMRMaker got status:", status)
				if status.IsOngoing() {
					continue
				}

				if status != types.CompletedAbort {
					errCh <- fmt.Errorf("swap did not exit successfully for XMRMaker: got %s", status)
				}

				return
			case <-time.After(testTimeout):
				errCh <- errors.New("make offer subscription timed out")
			}
		}
	}()

	c := rpcclient.NewClient(defaultXMRTakerDaemonEndpoint)
	wsc, err := wsclient.NewWsClient(ctx, defaultXMRTakerDaemonWSEndpoint)
	require.NoError(t, err)

	providers, err := c.Discover(types.ProvidesXMR, defaultDiscoverTimeout)
	require.NoError(t, err)
	require.Equal(t, 1, len(providers))
	require.GreaterOrEqual(t, len(providers[0]), 2)

	takerStatusCh, err := wsc.TakeOfferAndSubscribe(providers[0][0], offerID, 0.05)
	require.NoError(t, err)

	go func() {
		defer wg.Done()
		for status := range takerStatusCh {
			fmt.Println("> XMRTaker got status:", status)
			if status != types.ExpectingKeys {
				continue
			}

			fmt.Println("> XMRTaker cancelled swap!")
			exitStatus, err := c.Cancel(offerID) //nolint:govet
			if err != nil {
				errCh <- err
				return
			}

			if exitStatus != types.CompletedAbort {
				errCh <- fmt.Errorf("did not refund exit: exit status was %s", exitStatus)
			}

			return
		}
	}()

	wg.Wait()

	select {
	case err = <-errCh:
		require.NoError(t, err)
	default:
	}

	offersAfter, err := bc.GetOffers()
	require.NoError(t, err)
	require.Equal(t, len(offersBefore), len(offersAfter))
}

// This test simulates the case where neither XMRTaker and XMRMaker have
// locked funds yet, and XMRMaker cancels the swap.
// The swap should abort on XMRMaker's side, but might abort *or* refund on XMRTaker's side, in case she ended up
// locking ETH before she was notified that XMRMaker disconnected.
func TestAbort_XMRMakerCancels(t *testing.T) {
	if os.Getenv(generateBlocksEnv) != falseStr {
		generateBlocks(64)
	}

	const testTimeout = time.Second * 60

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bcli := rpcclient.NewClient(defaultXMRMakerDaemonEndpoint)
	bwsc, err := wsclient.NewWsClient(ctx, defaultXMRMakerDaemonWSEndpoint)
	require.NoError(t, err)

	offerID, statusCh, err := bwsc.MakeOfferAndSubscribe(0.1, xmrmakerProvideAmount,
		types.ExchangeRate(exchangeRate))
	require.NoError(t, err)

	bc := rpcclient.NewClient(defaultXMRMakerDaemonEndpoint)
	offersBefore, err := bc.GetOffers()
	require.NoError(t, err)

	errCh := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			select {
			case status := <-statusCh:
				fmt.Println("> XMRMaker got status:", status)
				if status != types.KeysExchanged {
					continue
				}

				fmt.Println("> XMRMaker cancelled swap!")
				exitStatus, err := bcli.Cancel(offerID) //nolint:govet
				if err != nil {
					errCh <- err
					return
				}

				if exitStatus != types.CompletedAbort {
					errCh <- fmt.Errorf("did not abort successfully: exit status was %s", exitStatus)
					return
				}

				fmt.Println("> XMRMaker exited successfully")
				return
			case <-time.After(testTimeout):
				errCh <- errors.New("make offer subscription timed out")
			}
		}
	}()

	c := rpcclient.NewClient(defaultXMRTakerDaemonEndpoint)
	wsc, err := wsclient.NewWsClient(ctx, defaultXMRTakerDaemonWSEndpoint)
	require.NoError(t, err)

	providers, err := c.Discover(types.ProvidesXMR, defaultDiscoverTimeout)
	require.NoError(t, err)
	require.Equal(t, 1, len(providers))
	require.GreaterOrEqual(t, len(providers[0]), 2)

	takerStatusCh, err := wsc.TakeOfferAndSubscribe(providers[0][0], offerID, 0.05)
	require.NoError(t, err)

	go func() {
		defer wg.Done()

		for status := range takerStatusCh {
			fmt.Println("> XMRTaker got status:", status)
			if status.IsOngoing() {
				continue
			}

			if status != types.CompletedAbort && status != types.CompletedRefund {
				errCh <- fmt.Errorf("swap did not exit successfully: got %s", status)
				return
			}

			fmt.Println("> XMRTaker exited successfully")
			return
		}
	}()

	wg.Wait()

	select {
	case err = <-errCh:
		require.NoError(t, err)
	default:
	}

	offersAfter, err := bc.GetOffers()
	require.NoError(t, err)
	require.Equal(t, len(offersBefore), len(offersAfter))
}

// TestError_ShouldOnlyTakeOfferOnce tests the case where two takers try to take the same offer concurrently.
// Only one should succeed, the other should return an error or Abort status.
func TestError_ShouldOnlyTakeOfferOnce(t *testing.T) {
	const testTimeout = time.Second * 60

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bc := rpcclient.NewClient(defaultXMRMakerDaemonEndpoint)
	offerID, err := bc.MakeOffer(xmrmakerProvideAmount, xmrmakerProvideAmount, exchangeRate)
	require.NoError(t, err)

	ac := rpcclient.NewClient(defaultXMRTakerDaemonEndpoint)

	providers, err := ac.Discover(types.ProvidesXMR, defaultDiscoverTimeout)
	require.NoError(t, err)
	require.Equal(t, 1, len(providers))
	require.GreaterOrEqual(t, len(providers[0]), 2)

	errCh := make(chan error)

	go func() {
		wsc, err := wsclient.NewWsClient(ctx, defaultXMRTakerDaemonWSEndpoint)
		require.NoError(t, err)

		takerStatusCh, err := wsc.TakeOfferAndSubscribe(providers[0][0], offerID, 0.05)
		if err != nil {
			errCh <- err
			return
		}

		for status := range takerStatusCh {
			fmt.Println("> XMRTaker got status:", status)
			if status.IsOngoing() {
				continue
			}

			if status != types.CompletedSuccess {
				errCh <- fmt.Errorf("swap did not exit successfully: got %s", status)
				return
			}

			fmt.Println("> XMRTaker exited successfully")
			return
		}
	}()

	go func() {
		wsc, err := wsclient.NewWsClient(ctx, defaultCharlieDaemonWSEndpoint)
		require.NoError(t, err)

		takerStatusCh, err := wsc.TakeOfferAndSubscribe(providers[0][0], offerID, 0.05)
		if err != nil {
			errCh <- err
			return
		}

		for status := range takerStatusCh {
			fmt.Println("> XMRTaker got status:", status)
			if status.IsOngoing() {
				continue
			}

			if status != types.CompletedSuccess {
				errCh <- fmt.Errorf("swap did not exit successfully: got %s", status)
				return
			}

			fmt.Println("> XMRTaker exited successfully")
			return
		}
	}()

	select {
	case err := <-errCh:
		require.NotNil(t, err)
		t.Log("got expected error:", err)
	case <-time.After(testTimeout):
		t.Fatalf("did not get error from XMRTaker or Charlie")
	}

	select {
	case err := <-errCh:
		t.Fatalf("should only have one error! also got %s", err)
	default:
	}
}

func TestSuccess_ConcurrentSwaps(t *testing.T) {
	const testTimeout = time.Second * 180
	const numConcurrentSwaps = 10

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type makerTest struct {
		offerID  string
		statusCh <-chan types.Status
		errCh    chan error
	}

	makerTests := make([]*makerTest, numConcurrentSwaps)

	for i := 0; i < numConcurrentSwaps; i++ {
		bwsc, err := wsclient.NewWsClient(ctx, defaultXMRMakerDaemonWSEndpoint)
		require.NoError(t, err)

		offerID, statusCh, err := bwsc.MakeOfferAndSubscribe(0.1, xmrmakerProvideAmount,
			types.ExchangeRate(exchangeRate))
		require.NoError(t, err)

		fmt.Println("maker made offer ", offerID)

		makerTests[i] = &makerTest{
			offerID:  offerID,
			statusCh: statusCh,
			errCh:    make(chan error, 2),
		}
	}

	bc := rpcclient.NewClient(defaultXMRMakerDaemonEndpoint)
	offersBefore, err := bc.GetOffers()
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(2 * numConcurrentSwaps)

	for _, tc := range makerTests {
		go func(tc *makerTest) {
			defer wg.Done()

			for {
				select {
				case status := <-tc.statusCh:
					fmt.Println("> XMRMaker got status:", status)
					if status.IsOngoing() {
						continue
					}

					if status != types.CompletedSuccess {
						tc.errCh <- fmt.Errorf("swap did not complete successfully: got %s", status)
					}

					return
				case <-time.After(testTimeout):
					tc.errCh <- errors.New("make offer subscription timed out")
				}
			}
		}(tc)
	}

	type takerTest struct {
		statusCh <-chan types.Status
		errCh    chan error
	}

	takerTests := make([]*takerTest, numConcurrentSwaps)

	for i := 0; i < numConcurrentSwaps; i++ {
		c := rpcclient.NewClient(defaultXMRTakerDaemonEndpoint)
		wsc, err := wsclient.NewWsClient(ctx, defaultXMRTakerDaemonWSEndpoint) //nolint:govet
		require.NoError(t, err)

		// TODO: implement discovery over websockets
		providers, err := c.Discover(types.ProvidesXMR, defaultDiscoverTimeout)
		require.NoError(t, err)
		require.Equal(t, 1, len(providers))
		require.GreaterOrEqual(t, len(providers[0]), 2)

		offerID := makerTests[i].offerID
		takerStatusCh, err := wsc.TakeOfferAndSubscribe(providers[0][0], offerID, 0.05)
		require.NoError(t, err)

		fmt.Println("taker took offer ", offerID)

		takerTests[i] = &takerTest{
			statusCh: takerStatusCh,
			errCh:    make(chan error, 2),
		}
	}

	for _, tc := range takerTests {
		go func(tc *takerTest) {
			defer wg.Done()
			for status := range tc.statusCh {
				fmt.Println("> XMRTaker got status:", status)
				if status.IsOngoing() {
					continue
				}

				if status != types.CompletedSuccess {
					tc.errCh <- fmt.Errorf("swap did not complete successfully: got %s", status)
				}

				return
			}
		}(tc)
	}

	wg.Wait()

	for _, tc := range makerTests {
		select {
		case err = <-tc.errCh:
			require.NoError(t, err)
		default:
		}
	}

	for _, tc := range takerTests {
		select {
		case err = <-tc.errCh:
			require.NoError(t, err)
		default:
		}
	}

	offersAfter, err := bc.GetOffers()
	require.NoError(t, err)
	require.Equal(t, numConcurrentSwaps, len(offersBefore)-len(offersAfter))
}
