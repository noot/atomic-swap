package wsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/noot/atomic-swap/common/rpctypes"
	"github.com/noot/atomic-swap/common/types"

	"github.com/gorilla/websocket"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("rpcclient")

// WsClient ...
type WsClient interface {
	Close()
	Discover(provides types.ProvidesCoin, searchTime uint64) ([][]string, error)
	Query(maddr string) (*rpctypes.QueryPeerResponse, error)
	SubscribeSwapStatus(id uint64) (<-chan types.Status, error)
	TakeOfferAndSubscribe(multiaddr, offerID string,
		providesAmount float64) (id uint64, ch <-chan types.Status, err error)
	MakeOfferAndSubscribe(min, max float64,
		exchangeRate types.ExchangeRate) (string, <-chan types.Status, error)
}

type wsClient struct {
	wmu  sync.Mutex
	rmu  sync.Mutex
	conn *websocket.Conn
}

// NewWsClient ...
func NewWsClient(ctx context.Context, endpoint string) (*wsClient, error) { ///nolint:revive
	conn, resp, err := websocket.DefaultDialer.DialContext(ctx, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial endpoint: %w", err)
	}

	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	return &wsClient{
		conn: conn,
	}, nil
}

func (c *wsClient) Close() {
	_ = c.conn.Close()
}

func (c *wsClient) writeJSON(msg *rpctypes.Request) error {
	c.wmu.Lock()
	defer c.wmu.Unlock()
	return c.conn.WriteJSON(msg)
}

func (c *wsClient) read() ([]byte, error) {
	c.rmu.Lock()
	defer c.rmu.Unlock()
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (c *wsClient) Discover(provides types.ProvidesCoin, searchTime uint64) ([][]string, error) {
	params := &rpctypes.DiscoverRequest{
		Provides:   provides,
		SearchTime: searchTime,
	}

	bz, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req := &rpctypes.Request{
		JSONRPC: rpctypes.DefaultJSONRPCVersion,
		Method:  "net_discover",
		Params:  bz,
		ID:      0,
	}

	if err = c.writeJSON(req); err != nil {
		return nil, err
	}

	message, err := c.read()
	if err != nil {
		return nil, fmt.Errorf("failed to read websockets message: %s", err)
	}

	var resp *rpctypes.Response
	err = json.Unmarshal(message, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("websocket server returned error: %w", resp.Error)
	}

	log.Debugf("received message over websockets: %s", message)
	var dresp *rpctypes.DiscoverResponse
	if err := json.Unmarshal(resp.Result, &dresp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal swap ID response: %s", err)
	}

	return dresp.Peers, nil
}

func (c *wsClient) Query(maddr string) (*rpctypes.QueryPeerResponse, error) {
	params := &rpctypes.QueryPeerRequest{
		Multiaddr: maddr,
	}

	bz, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req := &rpctypes.Request{
		JSONRPC: rpctypes.DefaultJSONRPCVersion,
		Method:  "net_queryPeer",
		Params:  bz,
		ID:      0,
	}

	if err = c.writeJSON(req); err != nil {
		return nil, err
	}

	// read ID from connection
	message, err := c.read()
	if err != nil {
		return nil, fmt.Errorf("failed to read websockets message: %s", err)
	}

	var resp *rpctypes.Response
	err = json.Unmarshal(message, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("websocket server returned error: %w", resp.Error)
	}

	log.Debugf("received message over websockets: %s", message)
	var dresp *rpctypes.QueryPeerResponse
	if err := json.Unmarshal(resp.Result, &dresp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal swap ID response: %s", err)
	}

	return dresp, nil
}

// SubscribeSwapStatus returns a channel that is written to each time the swap's status updates.
// If there is no swap with the given ID, it returns an error.
func (c *wsClient) SubscribeSwapStatus(id types.Hash) (<-chan types.Status, error) {
	params := &rpctypes.SubscribeSwapStatusRequest{
		ID: id,
	}

	bz, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req := &rpctypes.Request{
		JSONRPC: rpctypes.DefaultJSONRPCVersion,
		Method:  "swap_subscribeStatus",
		Params:  bz,
		ID:      0,
	}

	if err = c.writeJSON(req); err != nil {
		return nil, err
	}

	respCh := make(chan types.Status)

	go func() {
		defer close(respCh)

		for {
			message, err := c.read()
			if err != nil {
				log.Warnf("failed to read websockets message: %s", err)
				break
			}

			var resp *rpctypes.Response
			err = json.Unmarshal(message, &resp)
			if err != nil {
				log.Warnf("failed to unmarshal response: %s", err)
				break
			}

			if resp.Error != nil {
				log.Warnf("websocket server returned error: %s", resp.Error)
				break
			}

			log.Debugf("received message over websockets: %s", message)
			var status *rpctypes.SubscribeSwapStatusResponse
			if err := json.Unmarshal(resp.Result, &status); err != nil {
				log.Warnf("failed to unmarshal response: %s", err)
				break
			}

			s := types.NewStatus(status.Status)
			respCh <- s
			if !s.IsOngoing() {
				return
			}
		}
	}()

	return respCh, nil
}

func (c *wsClient) TakeOfferAndSubscribe(multiaddr, offerID string,
	providesAmount float64) (ch <-chan types.Status, err error) {
	params := &rpctypes.TakeOfferRequest{
		Multiaddr:      multiaddr,
		OfferID:        offerID,
		ProvidesAmount: providesAmount,
	}

	bz, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req := &rpctypes.Request{
		JSONRPC: rpctypes.DefaultJSONRPCVersion,
		Method:  "net_takeOfferAndSubscribe",
		Params:  bz,
		ID:      0,
	}

	if err = c.writeJSON(req); err != nil {
		return nil, err
	}

	// read ID from connection
	message, err := c.read()
	if err != nil {
		return nil, fmt.Errorf("failed to read websockets message: %s", err)
	}

	var resp *rpctypes.Response
	err = json.Unmarshal(message, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("websocket server returned error: %w", resp.Error)
	}

	log.Debugf("received message over websockets: %s", message)
	var idResp *rpctypes.TakeOfferResponse
	if err := json.Unmarshal(resp.Result, &idResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal swap ID response: %s", err)
	}

	respCh := make(chan types.Status)

	go func() {
		defer close(respCh)

		for {
			message, err := c.read()
			if err != nil {
				log.Warnf("failed to read websockets message: %s", err)
				break
			}

			var resp *rpctypes.Response
			err = json.Unmarshal(message, &resp)
			if err != nil {
				log.Warnf("failed to unmarshal response: %s", err)
				break
			}

			if resp.Error != nil {
				log.Warnf("websocket server returned error: %s", resp.Error)
				break
			}

			log.Debugf("received message over websockets: %s", message)
			var status *rpctypes.SubscribeSwapStatusResponse
			if err := json.Unmarshal(resp.Result, &status); err != nil {
				log.Warnf("failed to unmarshal swap status response: %s", err)
				break
			}

			s := types.NewStatus(status.Status)
			respCh <- s
			if !s.IsOngoing() {
				return
			}
		}
	}()

	return respCh, nil
}

func (c *wsClient) MakeOfferAndSubscribe(min, max float64,
	exchangeRate types.ExchangeRate) (string, <-chan types.Status, error) {
	params := &rpctypes.MakeOfferRequest{
		MinimumAmount: min,
		MaximumAmount: max,
		ExchangeRate:  exchangeRate,
	}

	bz, err := json.Marshal(params)
	if err != nil {
		return "", nil, err
	}

	req := &rpctypes.Request{
		JSONRPC: rpctypes.DefaultJSONRPCVersion,
		Method:  "net_makeOfferAndSubscribe", // TODO: use const
		Params:  bz,
		ID:      0,
	}

	if err = c.writeJSON(req); err != nil {
		return "", nil, err
	}

	// read ID from connection
	message, err := c.read()
	if err != nil {
		return "", nil, fmt.Errorf("failed to read websockets message: %s", err)
	}

	var resp *rpctypes.Response
	err = json.Unmarshal(message, &resp)
	if err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.Error != nil {
		return "", nil, fmt.Errorf("websocket server returned error: %w", resp.Error)
	}

	// read synchronous response (offer ID and infofile)
	var respData *rpctypes.MakeOfferResponse
	if err := json.Unmarshal(resp.Result, &respData); err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal response: %s", err)
	}

	respCh := make(chan types.Status)

	go func() {
		defer close(respCh)

		for {
			message, err := c.read()
			if err != nil {
				log.Warnf("failed to read websockets message: %s", err)
				break
			}

			var resp *rpctypes.Response
			err = json.Unmarshal(message, &resp)
			if err != nil {
				log.Warnf("failed to unmarshal response: %s", err)
				break
			}

			if resp.Error != nil {
				log.Warnf("websocket server returned error: %s", resp.Error)
				break
			}

			log.Debugf("received message over websockets: %s", message)
			var status *rpctypes.SubscribeSwapStatusResponse
			if err := json.Unmarshal(resp.Result, &status); err != nil {
				log.Warnf("failed to unmarshal response: %s", err)
				break
			}

			s := types.NewStatus(status.Status)
			respCh <- s
			if !s.IsOngoing() {
				return
			}
		}
	}()

	return respData.ID, respCh, nil
}
