package snark

import (
	"github.com/hyle-org/hyle/x/zktx"
)

type HyleCircomProof struct {
	Proof       interface{}     `json:"proof"`
	PublicInput interface{}     `json:"publicInput"`
	HyleOutput  zktx.HyleOutput `json:"hyleOutput"`
}
