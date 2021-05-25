package main

import (
	"fmt"

	"github.com/dedis/deter-checker/examples/neff"
	"github.com/dedis/deter-checker/examples/sigver"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/proof"
	"go.dedis.ch/kyber/v3/shuffle"
	"go.dedis.ch/kyber/v3/sign/bls"
	"go.dedis.ch/kyber/v3/util/random"
	"go.dedis.ch/kyber/v3/xof/blake2xb"
)

var k = 5
var N = 7

func test_bls() error {
	fmt.Println("Starting bls test...")
	msg := []byte("Static deter-checker!")
	suite := bn256.NewSuite()
	private, public := bls.NewKeyPair(suite, random.New())
	sig, err := bls.Sign(suite, private, msg)
	if err != nil {
		return err
	}
	err = sigver.Verify(suite, public, msg, sig)
	return err
}

func test_neff() error {
	fmt.Println("Starting neff test...")
	suite := edwards25519.NewBlakeSHA256Ed25519WithRand(blake2xb.New(nil))
	rand := suite.RandomStream()

	// Create a "server" private/public keypair
	h := suite.Scalar().Pick(rand)
	H := suite.Point().Mul(h, nil)

	// Create a set of ephemeral "client" keypairs. Public keys are going
	// to be shuffled.
	c := make([]kyber.Scalar, k)
	C := make([]kyber.Point, k)
	for i := 0; i < k; i++ {
		c[i] = suite.Scalar().Pick(rand)
		C[i] = suite.Point().Mul(c[i], nil)
	}

	// ElGamal-encrypt the keypairs with the "server" key
	X := make([]kyber.Point, k)
	Y := make([]kyber.Point, k)
	r := suite.Scalar() // temporary
	for i := 0; i < k; i++ {
		r.Pick(rand)
		X[i] = suite.Point().Mul(r, nil)
		Y[i] = suite.Point().Mul(r, H) // ElGamal blinding factor
		Y[i].Add(Y[i], C[i])           // Encrypted client public key
	}

	for i := 0; i < N; i++ {
		Xbar, Ybar, prover := shuffle.Shuffle(suite, nil, H, X, Y, rand)
		prf, err := proof.HashProve(suite, "PairShuffle", prover)
		if err != nil {
			err := fmt.Errorf("Shuffle proof failed: %v", err.Error())
			return err
		}
		err = neff.VerifyShuffle(suite, nil, H, X, Y, Xbar, Ybar, prf)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := test_bls()
	if err != nil {
		fmt.Println("Signature verification failed:", err)
	}
	err = test_neff()
	if err != nil {
		fmt.Println("Neff shuffle verification failed:", err)
	}
}
