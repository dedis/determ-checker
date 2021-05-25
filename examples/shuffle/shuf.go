package shuffle

import (
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/xof/blake2xb"
)

var suite = edwards25519.NewBlakeSHA256Ed25519WithRand(blake2xb.New(nil))
