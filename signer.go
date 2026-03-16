package ethereal // import "ethereal-dev"

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const USER_AGENT = "ethereal-go/1.0.0dev"

type Signer struct {
	Subaccount *Subaccount
	types      *abi.TypedData
	pk         *ecdsa.PrivateKey
	Address    string
}

func NewSigner(pk *ecdsa.PrivateKey) *Signer {
	return &Signer{
		pk:      pk,
		Address: crypto.PubkeyToAddress(pk.PublicKey).Hex(),
	}
}

func (r *Signer) GetPk() *ecdsa.PrivateKey {
	return r.pk
}

func (r *Signer) SetTypes(t *abi.TypedData) {
	r.types = t
}

func (r *Signer) GetTypes() *abi.TypedData {
	return r.types
}
