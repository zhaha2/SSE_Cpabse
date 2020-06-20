package cpabse

import "github.com/Nik-U/pbc"

type ByteTk struct {
	Tok1  []byte /*(g^a*g^bH2(w))^s*/
	Tok2  []byte
	Tok3  []byte
	Comps []*ByteTkComp
}

type ByteTkComp struct {
	Attr string
	A_j1 []byte
	B_j1 []byte
}

type ByteCph struct {
	W  []byte /* G_T */
	W0 []byte /* G_1 */
	W1 []byte
	S  []byte
	P  *BytePolicy
}

type BytePolicy struct {
	/* k=1 if leaf, otherwise threshould */
	K int
	/* attribute string if leaf, otherwise null */
	Attr string
	C    []byte /* G_1 only for leaves */ //W
	Cp   []byte /* G_1 only for leaves */ //D
	/* array of BswabePolicy and length is 0 for leaves */
	Children []*BytePolicy

	/* only used during encryption */
	Q *BytePolynomial

	/* only used during decription */
	Satisfiable bool
	Min_leaves  int
	Attri       int
	Satl        []int
}

type BytePolynomial struct {
	Deg int
	/* coefficients from [0] x^0 to [deg] x^deg */
	Coef [][]byte /* G_T (of length deg+1) */
}

type BytePm struct {
	BG       []byte
	BG_a     []byte
	BG_b     []byte
	BG_c     []byte
	BG_y     []byte
	BG_hat_x []byte
}

type ByteMsk struct {
	BA []byte
	BB []byte
	BC []byte
	BX []byte
	BY []byte
}

func MskToBmsk(msk *CpabeMsk) *ByteMsk {
	bmsk := new(ByteMsk)
	bmsk.BA = msk.A.Bytes()
	bmsk.BB = msk.A.Bytes()
	bmsk.BC = msk.A.Bytes()
	bmsk.BX = msk.A.Bytes()
	bmsk.BY = msk.A.Bytes()
	return bmsk
}

func BmskToMsk(msk *CpabeMsk, bmsk *ByteMsk, pm *CpabePm) {
	msk.A = pm.P.NewZr()
	msk.B = pm.P.NewZr()
	msk.C = pm.P.NewZr()
	msk.X = pm.P.NewZr()
	msk.Y = pm.P.NewZr()
	msk.A = msk.A.SetBytes(bmsk.BA)
	msk.B = msk.B.SetBytes(bmsk.BB)
	msk.C = msk.C.SetBytes(bmsk.BC)
	msk.X = msk.X.SetBytes(bmsk.BX)
	msk.Y = msk.Y.SetBytes(bmsk.BY)
}

func PmToBpm(pm *CpabePm) *BytePm {
	bpm := new(BytePm)
	bpm.BG = pm.G.Bytes()
	bpm.BG_a = pm.G_a.Bytes()
	bpm.BG_b = pm.G_b.Bytes()
	bpm.BG_c = pm.G_c.Bytes()
	bpm.BG_y = pm.G_y.Bytes()
	bpm.BG_hat_x = pm.G_hat_x.Bytes()
	return bpm
}

func BpmToPm(pm *CpabePm, bpm *BytePm) {
	pm.G = pm.P.NewG1()
	pm.G_a = pm.P.NewG1()
	pm.G_b = pm.P.NewG1()
	pm.G_c = pm.P.NewG1()
	pm.G_hat_x = pm.P.NewGT()
	pm.G_y = pm.P.NewG1()
	pm.G = pm.G.SetBytes(bpm.BG)
	pm.G_a = pm.G_a.SetBytes(bpm.BG_a)
	pm.G_b = pm.G_b.SetBytes(bpm.BG_b)
	pm.G_c = pm.G_c.SetBytes(bpm.BG_c)
	pm.G_hat_x = pm.G_hat_x.SetBytes(bpm.BG_hat_x)
	pm.G_y = pm.G_y.SetBytes(bpm.BG_y)
}

func TkToByte(tk *Token) *ByteTk {
	btk := new(ByteTk)
	btk.Tok1 = tk.Tok1.Bytes()
	btk.Tok2 = tk.Tok2.Bytes()
	btk.Tok3 = tk.Tok3.Bytes()
	for _, x := range tk.Comps {
		btk.Comps = append(btk.Comps, TkcToBTkc(x))
	}
	return btk
}

func TkcToBTkc(tkc *TokenComp) *ByteTkComp {
	btkc := new(ByteTkComp)
	btkc.Attr = tkc.Attr
	btkc.A_j1 = tkc.A_j1.Bytes()
	btkc.B_j1 = tkc.B_j1.Bytes()
	return btkc
}

func ByteToTk(btk *ByteTk, p *pbc.Pairing) *Token {
	tk := new(Token)
	tk.Tok1 = p.NewG1()
	tk.Tok2 = p.NewG1()
	tk.Tok3 = p.NewG1()
	tk.Tok1 = tk.Tok1.SetBytes(btk.Tok1)
	tk.Tok2 = tk.Tok2.SetBytes(btk.Tok2)
	tk.Tok3 = tk.Tok3.SetBytes(btk.Tok3)
	for _, x := range btk.Comps {
		tk.Comps = append(tk.Comps, BTkcToTkc(x, p))
	}
	return tk
}

func BTkcToTkc(btkc *ByteTkComp, p *pbc.Pairing) *TokenComp {
	tkc := new(TokenComp)
	tkc.Attr = btkc.Attr
	tkc.A_j1 = p.NewG2()
	tkc.B_j1 = p.NewG2()
	tkc.A_j1 = tkc.A_j1.SetBytes(btkc.A_j1)
	tkc.B_j1 = tkc.B_j1.SetBytes(btkc.B_j1)
	return tkc
}

func CphToByte(cph *CpabeCph) *ByteCph {
	bcph := new(ByteCph)
	bcph.W = cph.W.Bytes()
	bcph.W0 = cph.W0.Bytes()
	bcph.W1 = cph.W1.Bytes()
	//	bcph.S = cph.S.Bytes()
	bcph.P = CPToBP(cph.P)
	return bcph
}

func CPToBP(cp *CpabePolicy) *BytePolicy {
	bp := new(BytePolicy)
	bp.K = cp.K
	bp.Attr = cp.Attr
	if cp.C != nil && cp.Cp != nil {
		bp.C = cp.C.Bytes()
		bp.Cp = cp.Cp.Bytes()
	}
	for _, x := range cp.Children {
		bp.Children = append(bp.Children, CPToBP(x))
	}
	bp.Q = CPolyToBPoly(cp.Q)
	bp.Satisfiable = cp.Satisfiable
	bp.Min_leaves = cp.Min_leaves
	bp.Attri = cp.Attri
	bp.Satl = cp.Satl
	return bp
}

func CPolyToBPoly(cpo *CpabePolynomial) *BytePolynomial {
	bpo := new(BytePolynomial)
	bpo.Deg = cpo.Deg
	for _, x := range cpo.Coef {
		bpo.Coef = append(bpo.Coef, x.Bytes())
	}
	return bpo
}

func ByteToCph(bcph *ByteCph, p *pbc.Pairing) *CpabeCph {
	cph := new(CpabeCph)
	cph.W = p.NewG1()
	cph.W0 = p.NewG1()
	cph.W1 = p.NewG1()
	//cph.S = p.NewZr()
	cph.W = cph.W.SetBytes(bcph.W)
	cph.W0 = cph.W0.SetBytes(bcph.W0)
	cph.W1 = cph.W1.SetBytes(bcph.W1)
	//cph.S = cph.S.SetBytes(bcph.S)
	cph.P = BPToCP(bcph.P, p)
	return cph
}

func BPToCP(bp *BytePolicy, p *pbc.Pairing) *CpabePolicy {
	cp := new(CpabePolicy)
	cp.K = bp.K
	cp.Attr = bp.Attr
	if bp.C != nil && bp.Cp != nil {
		cp.C = p.NewG1()
		cp.Cp = p.NewG1()
		cp.C = cp.C.SetBytes(bp.C)
		cp.Cp = cp.Cp.SetBytes(bp.Cp)
	}
	for _, x := range bp.Children {
		cp.Children = append(cp.Children, BPToCP(x, p))
	}
	cp.Q = BPolyToCPoly(bp.Q, p)
	cp.Satisfiable = bp.Satisfiable
	cp.Min_leaves = bp.Min_leaves
	cp.Attri = bp.Attri
	cp.Satl = bp.Satl
	return cp
}

func BPolyToCPoly(bpo *BytePolynomial, p *pbc.Pairing) *CpabePolynomial {
	cpo := new(CpabePolynomial)
	cpo.Deg = bpo.Deg
	for _, x := range bpo.Coef {
		var temp *pbc.Element
		temp = p.NewZr()
		cpo.Coef = append(cpo.Coef, temp.SetBytes(x))
	}
	return cpo
}

