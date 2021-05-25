package neff

import (
	"bytes"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/proof"
	"go.dedis.ch/kyber/v3/shuffle"
)

type ega1 struct {
	Gamma            kyber.Point
	A, C, U, W       []kyber.Point
	Lambda1, Lambda2 kyber.Point
}

// V (Verifier) step 2: random challenge t
type ega2 struct {
	Zrho []kyber.Scalar
}

// P step 3: Theta vectors
type ega3 struct {
	D []kyber.Point
}

// V step 4: random challenge c
type ega4 struct {
	Zlambda kyber.Scalar
}

// P step 5: alpha vector
type ega5 struct {
	Zsigma []kyber.Scalar
	Ztau   kyber.Scalar
}

// P and V, step 5: simple k-shuffle proof
type ega6 struct {
	shuffle.SimpleShuffle
}

// PairShuffle creates a proof of the correctness of a shuffle
// of a series of ElGamal pairs.
//
// The caller must first invoke Init()
// to establish the cryptographic parameters for the shuffle:
// in particular, the relevant cryptographic Group,
// and the number of ElGamal pairs to be shuffled.
//
// The caller then may either perform its own shuffle,
// according to a permutation of the caller's choosing,
// and invoke Prove() to create a proof of its correctness;
// or alternatively the caller may simply invoke Shuffle()
// to pick a random permutation, compute the shuffle,
// and compute the correctness proof.
type PairShuffle struct {
	grp kyber.Group
	k   int
	p1  ega1
	v2  ega2
	p3  ega3
	v4  ega4
	p5  ega5
	pv6 shuffle.SimpleShuffle
}

func (ps *PairShuffle) Init(grp kyber.Group, k int) *PairShuffle {

	if k <= 1 {
		panic("can't shuffle permutation of size <= 1")
	}

	// Create a well-formed PairShuffleProof with arrays correctly sized.
	ps.grp = grp
	ps.k = k
	ps.p1.A = make([]kyber.Point, k)
	ps.p1.C = make([]kyber.Point, k)
	ps.p1.U = make([]kyber.Point, k)
	ps.p1.W = make([]kyber.Point, k)
	ps.v2.Zrho = make([]kyber.Scalar, k)
	ps.p3.D = make([]kyber.Point, k)
	ps.p5.Zsigma = make([]kyber.Scalar, k)
	ps.pv6.Init(grp, k)

	return ps
}

type hashVerifier struct {
	suite   proof.Suite
	proof   bytes.Buffer // Buffer with which to read the proof
	prbuf   []byte       // Byte-slice underlying proof buffer
	pubrand kyber.XOF
}
