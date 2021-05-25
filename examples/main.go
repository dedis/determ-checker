package main

import (
	"fmt"

	"github.com/dedis/deter-checker/examples/sigver"
	"go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/sign/bls"
	"go.dedis.ch/kyber/v3/util/random"
)

func test_bls() error {
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

func main() {
	err := test_bls()
	if err != nil {
		fmt.Println("Signature verification failed:", err)
	}
}
