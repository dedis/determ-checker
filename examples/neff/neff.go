package neff

import (
	"errors"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/proof"
	"go.dedis.ch/kyber/v3/shuffle"
)

func (ps *PairShuffle) Verify(g, h kyber.Point, X, Y, Xbar, Ybar []kyber.Point,
	ctx proof.VerifierContext) error {

	// Validate all vector lengths
	grp := ps.grp
	k := ps.k
	if len(X) != k || len(Y) != k || len(Xbar) != k || len(Ybar) != k {
		panic("mismatched vector lengths")
	}

	// P step 1
	p1 := &ps.p1
	if err := ctx.Get(p1); err != nil {
		return err
	}

	// V step 2
	v2 := &ps.v2
	if err := ctx.PubRand(v2); err != nil {
		return err
	}
	B := make([]kyber.Point, k)
	for i := 0; i < k; i++ {
		P := grp.Point().Mul(v2.Zrho[i], g)
		B[i] = P.Sub(P, p1.U[i])
	}

	// P step 3
	p3 := &ps.p3
	if err := ctx.Get(p3); err != nil {
		return err
	}

	// V step 4
	v4 := &ps.v4
	if err := ctx.PubRand(v4); err != nil {
		return err
	}

	// P step 5
	p5 := &ps.p5
	if err := ctx.Get(p5); err != nil {
		return err
	}

	// P,V step 6: simple k-shuffle
	if err := ps.pv6.Verify(g, p1.Gamma, ctx); err != nil {
		return err
	}

	// V step 7
	Phi1 := grp.Point().Null()
	Phi2 := grp.Point().Null()
	P := grp.Point() // scratch
	Q := grp.Point() // scratch
	for i := 0; i < k; i++ {
		Phi1 = Phi1.Add(Phi1, P.Mul(p5.Zsigma[i], Xbar[i])) // (31)
		Phi1 = Phi1.Sub(Phi1, P.Mul(v2.Zrho[i], X[i]))
		Phi2 = Phi2.Add(Phi2, P.Mul(p5.Zsigma[i], Ybar[i])) // (32)
		Phi2 = Phi2.Sub(Phi2, P.Mul(v2.Zrho[i], Y[i]))
		if !P.Mul(p5.Zsigma[i], p1.Gamma).Equal( // (33)
			Q.Add(p1.W[i], p3.D[i])) {
			return errors.New("invalid PairShuffleProof")
		}
	}

	if !P.Add(p1.Lambda1, Q.Mul(p5.Ztau, g)).Equal(Phi1) || // (34)
		!P.Add(p1.Lambda2, Q.Mul(p5.Ztau, h)).Equal(Phi2) { // (35)
		return errors.New("invalid PairShuffleProof")
	}

	return nil
}

func Verifier(group kyber.Group, g, h kyber.Point,
	X, Y, Xbar, Ybar []kyber.Point) proof.Verifier {

	ps := PairShuffle{}
	ps.Init(group, len(X))
	verifier := func(ctx proof.VerifierContext) error {
		return ps.Verify(g, h, X, Y, Xbar, Ybar, ctx)
	}
	return verifier
}

func VerifyShuffle(suite shuffle.Suite, g, h kyber.Point, X, Y, Xbar,
	Ybar []kyber.Point, prf []byte) error {

	verifier := Verifier(suite, nil, h, X, Y, Xbar, Ybar)
	// For now I could not bring this in because of the expected function
	// parameters in kyber.Proof. However, this is still OK for now because
	// the verifier function, where most of the crypto is done, is already
	// in this filer. See the function Verifier() above.
	err := proof.HashVerify(suite, "PairShuffle", verifier, prf)
	return err
}
