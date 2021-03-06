package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/noot/atomic-swap/common/rpctypes"
	"github.com/noot/atomic-swap/common/types"
	mcrypto "github.com/noot/atomic-swap/crypto/monero"
	"github.com/noot/atomic-swap/protocol/txsender"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
)

const (
	subscribeNewPeer    = "net_subscribeNewPeer"
	subscribeMakeOffer  = "net_makeOfferAndSubscribe"
	subscribeTakeOffer  = "net_takeOfferAndSubscribe"
	subscribeSwapStatus = "swap_subscribeStatus"
	subscribeSigner     = "signer_subscribe"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: checkOriginFunc,
}

func checkOriginFunc(r *http.Request) bool {
	return true
}

type wsServer struct {
	ctx     context.Context
	sm      SwapManager
	ns      *NetService
	backend ProtocolBackend
	signer  *txsender.ExternalSender
}

func newWsServer(ctx context.Context, sm SwapManager, ns *NetService, backend ProtocolBackend,
	signer *txsender.ExternalSender) *wsServer {
	s := &wsServer{
		ctx:     ctx,
		sm:      sm,
		ns:      ns,
		backend: backend,
		signer:  signer,
	}

	return s
}

// ServeHTTP ...
func (s *wsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnf("failed to update connection to websockets: %s", err)
		return
	}

	defer conn.Close() //nolint:errcheck

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Warnf("failed to read websockets message: %s", err)
			break
		}

		var req *rpctypes.Request
		err = json.Unmarshal(message, &req)
		if err != nil {
			_ = writeError(conn, err)
			continue
		}

		log.Debugf("received message over websockets: %s", message)
		err = s.handleRequest(conn, req)
		if err != nil {
			_ = writeError(conn, err)
		}
	}
}

func (s *wsServer) handleRequest(conn *websocket.Conn, req *rpctypes.Request) error {
	switch req.Method {
	case subscribeSigner:
		var params *rpctypes.SignerRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return fmt.Errorf("failed to unmarshal parameters: %w", err)
		}

		return s.handleSigner(s.ctx, conn, params.OfferID, params.EthAddress, params.XMRAddress)
	case subscribeNewPeer:
		return errUnimplemented
	case "net_discover":
		var params *rpctypes.DiscoverRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return fmt.Errorf("failed to unmarshal parameters: %w", err)
		}

		resp := new(rpctypes.DiscoverResponse)
		err := s.ns.Discover(nil, params, resp)
		if err != nil {
			return err
		}

		return writeResponse(conn, resp)
	case "net_queryPeer":
		var params *rpctypes.QueryPeerRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return fmt.Errorf("failed to unmarshal parameters: %w", err)
		}

		resp := new(rpctypes.QueryPeerResponse)
		err := s.ns.QueryPeer(nil, params, resp)
		if err != nil {
			return err
		}

		return writeResponse(conn, resp)
	case subscribeSwapStatus:
		var params *rpctypes.SubscribeSwapStatusRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return fmt.Errorf("failed to unmarshal parameters: %w", err)
		}

		return s.subscribeSwapStatus(s.ctx, conn, params.ID)
	case subscribeTakeOffer:
		var params *rpctypes.TakeOfferRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return fmt.Errorf("failed to unmarshal parameters: %w", err)
		}

		ch, infofile, err := s.ns.takeOffer(params.Multiaddr, params.OfferID, params.ProvidesAmount)
		if err != nil {
			return err
		}

		return s.subscribeTakeOffer(s.ctx, conn, ch, infofile)
	case subscribeMakeOffer:
		var params *rpctypes.MakeOfferRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return fmt.Errorf("failed to unmarshal parameters: %w", err)
		}

		offerID, offerExtra, err := s.ns.makeOffer(params)
		if err != nil {
			return err
		}

		s.ns.net.Advertise()
		return s.subscribeMakeOffer(s.ctx, conn, offerID, offerExtra)
	default:
		return errInvalidMethod
	}
}

func (s *wsServer) handleSigner(ctx context.Context, conn *websocket.Conn, offerIDStr, ethAddress,
	xmrAddr string) error {
	if s.signer == nil {
		return errSignerNotRequired
	}

	if err := mcrypto.ValidateAddress(xmrAddr); err != nil {
		return err
	}

	s.backend.SetEthAddress(ethcommon.HexToAddress(ethAddress))

	offerID, err := offerIDStringToHash(offerIDStr)
	if err != nil {
		return err
	}

	s.backend.SetXMRDepositAddress(mcrypto.Address(xmrAddr), offerID)

	s.signer.AddID(offerID)
	defer s.signer.DeleteID(offerID)

	txsOutCh, err := s.signer.OngoingCh(offerID)
	if err != nil {
		return err
	}

	txsInCh, err := s.signer.IncomingCh(offerID)
	if err != nil {
		return err
	}

	for {
		select {
		// TODO: check if conn closes or swap exited
		case <-time.After(time.Hour): // TODO: vary timeout based on env
			return fmt.Errorf("signer timed out")
		case <-ctx.Done():
			return nil
		case tx := <-txsOutCh:
			log.Debugf("outbound tx: %v", tx)
			resp := &rpctypes.SignerResponse{
				OfferID: offerIDStr,
				To:      tx.To.String(),
				Data:    tx.Data,
				Value:   tx.Value,
			}

			err := conn.WriteJSON(resp)
			if err != nil {
				return err
			}

			_, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			var params *rpctypes.SignerTxSigned
			if err := json.Unmarshal(message, &params); err != nil {
				return fmt.Errorf("failed to unmarshal parameters: %w", err)
			}

			if params.OfferID != offerIDStr {
				return fmt.Errorf("got unexpected offerID %s, expected %s", params.OfferID, offerID)
			}

			txsInCh <- ethcommon.HexToHash(params.TxHash)
		}
	}
}

func (s *wsServer) subscribeTakeOffer(ctx context.Context, conn *websocket.Conn,
	statusCh <-chan types.Status, infofile string) error {
	resp := &rpctypes.TakeOfferResponse{
		InfoFile: infofile,
	}

	if err := writeResponse(conn, resp); err != nil {
		return err
	}

	for {
		select {
		case status, ok := <-statusCh:
			if !ok {
				return nil
			}

			resp := &rpctypes.SubscribeSwapStatusResponse{
				Status: status.String(),
			}

			if err := writeResponse(conn, resp); err != nil {
				return err
			}

			if !status.IsOngoing() {
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *wsServer) subscribeMakeOffer(ctx context.Context, conn *websocket.Conn,
	offerID string, offerExtra *types.OfferExtra) error {
	resp := &rpctypes.MakeOfferResponse{
		ID:       offerID,
		InfoFile: offerExtra.InfoFile,
	}

	if err := writeResponse(conn, resp); err != nil {
		return err
	}

	for {
		select {
		case status, ok := <-offerExtra.StatusCh:
			if !ok {
				return nil
			}

			resp := &rpctypes.SubscribeSwapStatusResponse{
				Status: status.String(),
			}

			if err := writeResponse(conn, resp); err != nil {
				return err
			}

			if !status.IsOngoing() {
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// subscribeSwapStatus writes the swap's stage to the connection every time it updates.
// when the swap completes, it writes the final status then closes the connection.
// example: `{"jsonrpc":"2.0", "method":"swap_subscribeStatus", "params": {"id": 0}, "id": 0}`
func (s *wsServer) subscribeSwapStatus(ctx context.Context, conn *websocket.Conn, id types.Hash) error {
	info := s.sm.GetOngoingSwap(id)
	if info == nil {
		return s.writeSwapExitStatus(conn, id)
	}

	statusCh := info.StatusCh()
	for {
		select {
		case status, ok := <-statusCh:
			if !ok {
				return nil
			}

			resp := &rpctypes.SubscribeSwapStatusResponse{
				Status: status.String(),
			}

			if err := writeResponse(conn, resp); err != nil {
				return err
			}

			if !status.IsOngoing() {
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *wsServer) writeSwapExitStatus(conn *websocket.Conn, id types.Hash) error {
	info := s.sm.GetPastSwap(id)
	if info == nil {
		return errNoSwapWithID
	}

	resp := &rpctypes.SubscribeSwapStatusResponse{
		Status: info.Status().String(),
	}

	if err := writeResponse(conn, resp); err != nil {
		return err
	}

	return nil
}

func writeResponse(conn *websocket.Conn, result interface{}) error {
	bz, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp := &rpctypes.Response{
		Version: rpctypes.DefaultJSONRPCVersion,
		Result:  bz,
	}

	return conn.WriteJSON(resp)
}

func writeError(conn *websocket.Conn, err error) error {
	resp := &rpctypes.Response{
		Version: rpctypes.DefaultJSONRPCVersion,
		Error: &rpctypes.Error{
			Message: err.Error(),
		},
	}

	return conn.WriteJSON(resp)
}
