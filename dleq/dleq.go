package dleq

import (
	"github.com/noot/atomic-swap/crypto/secp256k1"
)

// Interface ...
type Interface interface {
	Prove() (*Proof, error)
	Verify(*Proof) (*VerifyResult, error)
}

// Proof represents a DLEq proof
type Proof struct {
	secret [32]byte
	proof  []byte
}

// NewProofWithoutSecret returns a new Proof without a secret from the given proof slice
func NewProofWithoutSecret(p []byte) *Proof {
	return &Proof{
		proof: p,
	}
}

// NewProofWithSecret returns a new Proof with the given secret.
// Note that the returned proof actually lacks the `proof` field.
func NewProofWithSecret(s [32]byte) *Proof {
	return &Proof{
		secret: s,
	}
}

// Secret returns the proof's 32-byte secret
func (p *Proof) Secret() [32]byte {
	return p.secret
}

// Proof returns the encoded DLEq proof
func (p *Proof) Proof() []byte {
	return p.proof
}

// VerifyResult contains the public keys resulting from verifying a DLEq proof
type VerifyResult struct {
	ed25519Pub   [32]byte
	secp256k1Pub *secp256k1.PublicKey
}

// Secp256k1PublicKey returns the secp256k1 public key associated with the DLEq verification
func (r *VerifyResult) Secp256k1PublicKey() *secp256k1.PublicKey {
	return r.secp256k1Pub
}
