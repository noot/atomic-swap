package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/noot/atomic-swap/common/types"
)

// Type represents the type of a network message
type Type byte

const (
	QueryResponseType Type = iota //nolint
	SendKeysType
	NotifyContractDeployedType // TODO: rename to NotifyETHLockType
	NotifyXMRLockType
	NotifyReadyType
	NotifyClaimedType
	NotifyRefundType
	NilType
)

func (t Type) String() string {
	switch t {
	case QueryResponseType:
		return "QueryResponse"
	case SendKeysType:
		return "SendKeysMessage"
	case NotifyContractDeployedType:
		return "NotifyContractDeployed"
	case NotifyXMRLockType:
		return "NotifyXMRLock"
	case NotifyClaimedType:
		return "NotifyClaimed"
	case NotifyRefundType:
		return "NotifyRefund"
	default:
		return "unknown"
	}
}

// Message must be implemented by all network messages
type Message interface {
	String() string
	Encode() ([]byte, error)
	Type() Type
}

// DecodeMessage decodes the given bytes into a Message
func DecodeMessage(b []byte) (Message, error) {
	if len(b) == 0 {
		return nil, errors.New("invalid message bytes")
	}

	switch Type(b[0]) {
	case QueryResponseType:
		var m *QueryResponse
		if err := json.Unmarshal(b[1:], &m); err != nil {
			return nil, err
		}
		return m, nil
	case SendKeysType:
		var m *SendKeysMessage
		if err := json.Unmarshal(b[1:], &m); err != nil {
			return nil, err
		}
		return m, nil
	case NotifyContractDeployedType:
		var m *NotifyContractDeployed
		if err := json.Unmarshal(b[1:], &m); err != nil {
			return nil, err
		}
		return m, nil
	case NotifyXMRLockType:
		var m *NotifyXMRLock
		if err := json.Unmarshal(b[1:], &m); err != nil {
			return nil, err
		}
		return m, nil
	case NotifyReadyType:
		var m *NotifyReady
		if err := json.Unmarshal(b[1:], &m); err != nil {
			return nil, err
		}
		return m, nil
	case NotifyClaimedType:
		var m *NotifyClaimed
		if err := json.Unmarshal(b[1:], &m); err != nil {
			return nil, err
		}
		return m, nil
	default:
		return nil, errors.New("invalid message type")
	}
}

// QueryResponse ...
type QueryResponse struct {
	Offers []*types.Offer
}

// String ...
func (m *QueryResponse) String() string {
	return fmt.Sprintf("QueryResponse Offers=%v",
		m.Offers,
	)
}

// Encode ...
func (m *QueryResponse) Encode() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append([]byte{byte(QueryResponseType)}, b...), nil
}

// Type ...
func (m *QueryResponse) Type() Type {
	return QueryResponseType
}

// The below messages are sawp protocol messages, exchanged after the swap has been agreed
// upon by both sides.

// SendKeysMessage is sent by both parties to each other to initiate the protocol
type SendKeysMessage struct {
	OfferID            string
	ProvidedAmount     float64
	PublicSpendKey     string
	PublicViewKey      string
	PrivateViewKey     string
	DLEqProof          string
	Secp256k1PublicKey string
	EthAddress         string
}

// String ...
func (m *SendKeysMessage) String() string {
	return fmt.Sprintf("SendKeysMessage OfferID=%s ProvidedAmount=%v PublicSpendKey=%s PublicViewKey=%s PrivateViewKey=%s DLEqProof=%s Secp256k1PublicKey=%s EthAddress=%s", //nolint:lll
		m.OfferID,
		m.ProvidedAmount,
		m.PublicSpendKey,
		m.PublicViewKey,
		m.PrivateViewKey,
		m.DLEqProof,
		m.Secp256k1PublicKey,
		m.EthAddress,
	)
}

// Encode ...
func (m *SendKeysMessage) Encode() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append([]byte{byte(SendKeysType)}, b...), nil
}

// Type ...
func (m *SendKeysMessage) Type() Type {
	return SendKeysType
}

// NotifyContractDeployed is sent by Alice to Bob after deploying the swap contract
// and locking her ether in it
type NotifyContractDeployed struct {
	Address        string
	ContractSwapID *big.Int
}

// String ...
func (m *NotifyContractDeployed) String() string {
	return fmt.Sprintf("NotifyContractDeployed Address=%s ContractSwapID=%d", m.Address, m.ContractSwapID)
}

// Encode ...
func (m *NotifyContractDeployed) Encode() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append([]byte{byte(NotifyContractDeployedType)}, b...), nil
}

// Type ...
func (m *NotifyContractDeployed) Type() Type {
	return NotifyContractDeployedType
}

// NotifyXMRLock is sent by Bob to Alice after locking his XMR.
type NotifyXMRLock struct {
	Address string
}

// String ...
func (m *NotifyXMRLock) String() string {
	return "NotifyXMRLock"
}

// Encode ...
func (m *NotifyXMRLock) Encode() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append([]byte{byte(NotifyXMRLockType)}, b...), nil
}

// Type ...
func (m *NotifyXMRLock) Type() Type {
	return NotifyXMRLockType
}

// NotifyReady is sent by Alice to Bob after calling Ready() on the contract.
type NotifyReady struct{}

// String ...
func (m *NotifyReady) String() string {
	return "NotifyReady"
}

// Encode ...
func (m *NotifyReady) Encode() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append([]byte{byte(NotifyReadyType)}, b...), nil
}

// Type ...
func (m *NotifyReady) Type() Type {
	return NotifyReadyType
}

// NotifyClaimed is sent by Bob to Alice after claiming his ETH.
type NotifyClaimed struct {
	TxHash string
}

// String ...
func (m *NotifyClaimed) String() string {
	return fmt.Sprintf("NotifyClaimed %s", m.TxHash)
}

// Encode ...
func (m *NotifyClaimed) Encode() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append([]byte{byte(NotifyClaimedType)}, b...), nil
}

// Type ...
func (m *NotifyClaimed) Type() Type {
	return NotifyClaimedType
}

// NotifyRefund is sent by Alice to Bob after calling Refund() on the contract.
type NotifyRefund struct {
	TxHash string
}

// String ...
func (m *NotifyRefund) String() string {
	return fmt.Sprintf("NotifyClaimed %s", m.TxHash)
}

// Encode ...
func (m *NotifyRefund) Encode() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append([]byte{byte(NotifyRefundType)}, b...), nil
}

// Type ...
func (m *NotifyRefund) Type() Type {
	return NotifyRefundType
}
