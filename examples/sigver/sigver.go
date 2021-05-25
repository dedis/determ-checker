package sigver

import (
	"errors"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
)

type hashablePoint interface {
	Hash([]byte) kyber.Point
}

func Verify(suite pairing.Suite, X kyber.Point, msg, sig []byte) error {
	hashable, ok := suite.G1().Point().(hashablePoint)
	if !ok {
		return errors.New("bls: point needs to implement hashablePoint")
	}
	HM := hashable.Hash(msg)
	left := suite.Pair(HM, X)
	s := suite.G1().Point()
	if err := s.UnmarshalBinary(sig); err != nil {
		return err
	}
	right := suite.Pair(s, suite.G2().Point().Base())
	if !left.Equal(right) {
		return errors.New("bls: invalid signature")
	}
	return nil
}
