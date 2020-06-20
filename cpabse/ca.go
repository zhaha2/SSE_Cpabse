package cpabse

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"

	"github.com/Nik-U/pbc"
)

var curveParams = "type a\n" +
	"q 87807107996633125224377819847540498158068831994142082" +
	"1102865339926647563088022295707862517942266222142315585" +
	"8769582317459277713367317481324925129998224791\n" +
	"h 12016012264891146079388821366740534204802954401251311" +
	"822919615131047207289359704531102844802183906537786776\n" +
	"r 730750818665451621361119245571504901405976559617\n" +
	"exp2 159\n" + "exp1 107\n" + "sign1 1\n" + "sign0 1\n"

func Psetup(pm *CpabePm) {
	params := new(pbc.Params)
	params, _ = pbc.NewParamsFromString(curveParams)
	pm.PairingDesc = curveParams
	pm.P = pbc.NewPairing(params)
}

func Setup(pm *CpabePm, msk *CpabeMsk) {
	var g_x *pbc.Element
	params := new(pbc.Params)
	params, _ = pbc.NewParamsFromString(curveParams)
	pm.PairingDesc = curveParams
	pm.P = pbc.NewPairing(params)
	pairing := pm.P

	pm.G = pairing.NewG1()
	msk.A = pairing.NewZr()
	msk.B = pairing.NewZr()
	msk.C = pairing.NewZr()
	msk.X = pairing.NewZr()
	msk.Y = pairing.NewZr()
	pm.G_a = pairing.NewG1()
	pm.G_b = pairing.NewG1()
	pm.G_c = pairing.NewG1()
	pm.G_hat_x = pairing.NewGT()

	msk.A.Rand()
	msk.B.Rand()
	msk.C.Rand()
	msk.X.Rand()
	msk.Y.Rand()
	pm.G.Rand()

	pm.G_a = pm.G.NewFieldElement().Set(pm.G)
	pm.G_a.PowZn(pm.G, msk.A)

	pm.G_b = pm.G.NewFieldElement().Set(pm.G)
	pm.G_b.PowZn(pm.G, msk.B)

	pm.G_c = pm.G.NewFieldElement().Set(pm.G)
	pm.G_c.PowZn(pm.G, msk.C)

	pm.G_y = pm.G.NewFieldElement().Set(pm.G)
	pm.G_y.PowZn(pm.G, msk.Y)

	g_x = pm.G.NewFieldElement().Set(pm.G)
	g_x.PowZn(g_x, msk.X)

	pm.G_hat_x.Pair(g_x, pm.G)
}

func Keygen(pm *CpabePm, msk *CpabeMsk, attrs []string) *CpabeSk {
	//attrs := strings.Split(attr, " ")

	sk := new(CpabeSk)
	var a, b, c, r, x, y, m, n, g_r *pbc.Element
	var pairing *pbc.Pairing

	/* initialize */
	pairing = pm.P
	a = msk.A
	b = msk.B
	c = msk.C
	x = msk.X
	y = msk.Y
	r = pairing.NewZr()
	m = pairing.NewG1()
	n = pairing.NewG1()
	sk.D = pairing.NewG1()
	g_r = pairing.NewG1()

	/* compute */
	/* compute sk.A*/
	r.Rand()
	sk.R = r.NewFieldElement().Set(r)

	m = a.NewFieldElement().Set(a)
	m.Mul(m, c)
	m.Sub(m, r)
	m.Div(m, b)
	sk.A = pm.G.NewFieldElement().Set(pm.G)
	sk.A.PowZn(sk.A, m)

	n = x.NewFieldElement().Set(x)
	n.Add(n, r)
	n.Div(m, y)
	sk.B = pm.G.NewFieldElement().Set(pm.G)
	sk.B.PowZn(sk.B, n)

	g_r = pm.G.NewFieldElement().Set(pm.G)
	g_r.PowZn(g_r, r)

	len := len(attrs)
	for i := 0; i < len; i++ {
		comp := new(CpabeSkComp)
		var h_rp, rp *pbc.Element

		comp.Attr = attrs[i]
		comp.A_j = pairing.NewG1()
		comp.B_j = pairing.NewG1()
		h_rp = pairing.NewG1()
		rp = pairing.NewZr()

		elementFromString(h_rp, comp.Attr) //hash1(at_j)
		rp.Rand()
		h_rp.PowZn(h_rp, rp)

		comp.A_j = g_r.NewFieldElement().Set(g_r)
		comp.A_j.Mul(comp.A_j, h_rp)
		comp.B_j = pm.G.NewFieldElement().Set(pm.G)
		comp.B_j.PowZn(comp.B_j, rp)

		sk.Comps = append(sk.Comps, comp)
	}
	return sk
}

func EncKeyword(pm *CpabePm, policy string, msk *CpabeMsk, w string) *CpabeCph {
	cph := new(CpabeCph)
	var r1, r2, x1, x2, r *pbc.Element

	/* initialize */
	pairing := pm.P
	r1 = pairing.NewZr()
	r2 = pairing.NewZr()
	r = pairing.NewZr()
	cph.W = pairing.NewG1()
	cph.W0 = pairing.NewG1()
	cph.W1 = pairing.NewG1()
	x1 = pairing.NewZr()
	x2 = pairing.NewG1()
	cph.P = parsePolicyPostfix(policy)

	/* compute */
	r1.Rand()
	r2.Rand()
	cph.W = pm.G.NewFieldElement().Set(pm.G)
	cph.W.PowZn(cph.W, msk.C)
	cph.W.PowZn(cph.W, r1)

	elementFromString(x1, w)
	cph.W0 = pm.G.NewFieldElement().Set(pm.G)
	cph.W0.PowZn(cph.W0, r1)
	cph.W0.PowZn(cph.W0, msk.B)
	cph.W0.PowZn(cph.W0, x1)
	x2 = pm.G.NewFieldElement().Set(pm.G)
	r = r1.NewFieldElement().Set(r1)
	r.Add(r, r2)
	x2.PowZn(x2, r)
	x2.PowZn(x2, msk.A)
	cph.W0.Mul(cph.W0, x2)

	cph.W1 = pm.G.NewFieldElement().Set(pm.G)
	cph.W1.PowZn(cph.W1, msk.B)
	cph.W1.PowZn(cph.W1, r2)

	fillPolicy(cph.P, pm, r2)

	return cph
}

func TokenGen(sk *CpabeSk, w string, msk *CpabeMsk, pm *CpabePm) *Token {
	tk := new(Token)
	var x1, s *pbc.Element

	pairing := pm.P
	x1 = pairing.NewZr()
	s = pairing.NewZr()
	s.Rand()
	elementFromString(x1, w)
	x1.Mul(x1, msk.B)
	x1.Add(x1, msk.A)
	tk.Tok1 = pm.G.NewFieldElement().Set(pm.G)
	tk.Tok1.PowZn(tk.Tok1, x1)
	tk.Tok1.PowZn(tk.Tok1, s)

	tk.Tok2 = pm.G.NewFieldElement().Set(pm.G)
	tk.Tok2.PowZn(tk.Tok2, msk.C)
	tk.Tok2.PowZn(tk.Tok2, s)

	tk.Tok3 = sk.A.NewFieldElement().Set(sk.A)
	tk.Tok3 = tk.Tok3.PowZn(tk.Tok3, s)

	len := len(sk.Comps)
	for i := 0; i < len; i++ {
		comp := new(TokenComp)
		comp.Attr = sk.Comps[i].Attr
		comp.A_j1 = sk.Comps[i].A_j.NewFieldElement().Set(sk.Comps[i].A_j)
		comp.A_j1.PowZn(comp.A_j1, s)
		comp.B_j1 = sk.Comps[i].B_j.NewFieldElement().Set(sk.Comps[i].B_j)
		comp.B_j1.PowZn(comp.B_j1, s)

		tk.Comps = append(tk.Comps, comp)
	}

	return tk
}

func Check(tk *Token, cph *CpabeCph, p *pbc.Pairing) bool {
	var t, s, s1, s2 *pbc.Element

	t = p.NewGT()
	s = p.NewGT()
	s1 = p.NewGT()
	s2 = p.NewGT()

	checkSatisfy(cph.P, tk)
	if !cph.P.Satisfiable {
		//fmt.Println("cannot decrypt, attributes in key do not satisfy policy")
		return false
	}

	pickSatisfyMinLeaves(cph.P, tk)
	decFlatten(t, cph.P, tk, p)
	s.Pair(cph.W0, tk.Tok2)
	s1.Pair(cph.W, tk.Tok1)
	s2.Pair(tk.Tok3, cph.W1)
	s1.Mul(s1, t)
	s1.Mul(s1, s2)
	if !s.Equals(s1) {
		return false
	} else {
		return true
	}
}

func decFlatten(r *pbc.Element, cp *CpabePolicy, tk *Token, p *pbc.Pairing) {
	var one *pbc.Element
	one = p.NewZr()
	one.Set1()
	r.Set1()

	decNodeFlatten(r, one, cp, tk, p)
}

func decNodeFlatten(r *pbc.Element, exp *pbc.Element, cp *CpabePolicy, tk *Token, p *pbc.Pairing) {
	if cp.Children == nil || len(cp.Children) == 0 {
		decLeafFlatten(r, exp, cp, tk, p)
	} else {
		decInternalFlatten(r, exp, cp, tk, p)
	}
}

func decLeafFlatten(r *pbc.Element, exp *pbc.Element, cp *CpabePolicy, tk *Token, p *pbc.Pairing) {
	c := new(TokenComp)
	var s, t *pbc.Element

	c = tk.Comps[cp.Attri]
	s = p.NewGT()
	t = p.NewGT()
	s.Pair(c.A_j1, cp.C)  /* num_pairings++; */
	t.Pair(c.B_j1, cp.Cp) /* num_pairings++; */
	t.Invert(t)
	s.Mul(s, t)     /* num_muls++; */
	s.PowZn(s, exp) /* num_exps++; */
	r.Mul(r, s)     /* num_muls++; */
}

func decInternalFlatten(r *pbc.Element, exp *pbc.Element, cp *CpabePolicy, tk *Token, p *pbc.Pairing) {
	var i int
	var t, expnew *pbc.Element
	t = p.NewZr()
	expnew = p.NewZr()
	for i = 0; i < len(cp.Satl); i++ {
		lagrangeCoef(t, cp.Satl, cp.Satl[i])
		expnew = exp.NewFieldElement().Set(exp)
		expnew.Mul(expnew, t)
		decNodeFlatten(r, expnew, cp.Children[cp.Satl[i]-1], tk, p)
	}
}

func lagrangeCoef(r *pbc.Element, s []int, i int) {
	var j, k int
	var t *pbc.Element

	t = r.NewFieldElement().Set(r)

	r.Set1()
	for k = 0; k < len(s); k++ {
		j = s[k]
		if j == i {
			continue
		}
		t.SetInt32(int32(-j))
		r.Mul(r, t) /* num_muls++; */
		t.SetInt32(int32(i - j))
		t.Invert(t)
		r.Mul(r, t) /* num_muls++; */
	}
}

func checkSatisfy(p *CpabePolicy, tk *Token) { //check access tree
	var i, l int
	var tkAttr string

	p.Satisfiable = false
	if p.Children == nil || len(p.Children) == 0 {
		for i = 0; i < len(tk.Comps); i++ {
			tkAttr = tk.Comps[i].Attr
			if strings.Compare(tkAttr, p.Attr) == 0 {
				p.Satisfiable = true
				p.Attri = i
				break
			}
		}
	} else {
		for i = 0; i < len(p.Children); i++ {
			checkSatisfy(p.Children[i], tk)
		}

		l = 0
		for i = 0; i < len(p.Children); i++ {
			if p.Children[i].Satisfiable {
				l++
			}
		}

		if l >= p.K {
			p.Satisfiable = true
		}
	}
}

func pickSatisfyMinLeaves(p *CpabePolicy, tk *Token) {
	var i, k, l, c_i int
	var c []int

	if p.Children == nil || len(p.Children) == 0 {
		p.Min_leaves = 1
	} else {
		len := len(p.Children)
		for i = 0; i < len; i++ {
			if p.Children[i].Satisfiable {
				pickSatisfyMinLeaves(p.Children[i], tk)
			}
		}

		for i = 0; i < len; i++ {
			c = append(c, i)
		}

		//TODO 这里的排序需要进一步改写,min_leaves是从小到大排序的，用了很low的冒泡排序。。。
		for i := 0; i < len; i++ {
			for j := 0; j < len-i-1; j++ {
				if p.Children[c[j]].Min_leaves > p.Children[c[j+1]].Min_leaves {
					c[j], c[j+1] = c[j+1], c[j]
				}
			}
		}

		p.Min_leaves = 0
		l = 0

		for i = 0; i < len && l < p.K; i++ {
			c_i = c[i] /* c[i] */
			if p.Children[c_i].Satisfiable {
				l++
				p.Min_leaves += p.Children[c_i].Min_leaves
				k = c_i + 1
				p.Satl = append(p.Satl, k)
			}
		}
	}
}

func parsePolicyPostfix(s string) *CpabePolicy {
	var toks []string
	var tok string
	var stack []*CpabePolicy
	var root *CpabePolicy

	toks = strings.Split(s, " ")

	toks_cnt := len(toks)
	for index := 0; index < toks_cnt; index++ {
		var i, k, n int

		tok = toks[index]
		if !strings.Contains(tok, "of") {
			stack = append(stack, baseNode(1, tok))
		} else {
			var node *CpabePolicy

			/* parse k of n node */
			k_n := strings.Split(tok, "of")
			k, _ = strconv.Atoi(k_n[0])
			n, _ = strconv.Atoi(k_n[1])

			if k < 1 {
				fmt.Println("error parsing " + s + ": trivially satisfied operator " + tok)
				return nil
			} else if k > n {
				fmt.Println("error parsing " + s + ": unsatisfiable operator " + tok)
				return nil
			} else if n == 1 {
				fmt.Println("error parsing " + s + ": indentity operator " + tok)
				return nil
			} else if n > len(stack) {
				fmt.Println("error parsing " + s + ": stack underflow at " + tok)
				return nil
			}

			/* pop n things and fill in children */
			node = baseNode(k, "")
			node.Children = make([]*CpabePolicy, n)

			for i = n - 1; i >= 0; i-- {
				node.Children[i] = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}

			/* push result */
			stack = append(stack, node)
		}
	}

	if len(stack) > 1 {
		fmt.Println("error parsing " + s + ": extra node left on the stack")
		return nil
	} else if len(stack) < 1 {
		fmt.Println("error parsing " + s + ": empty policy")
		return nil
	}

	root = stack[0]
	return root
}

func randPoly(deg int, zeroVal *pbc.Element) *CpabePolynomial {
	var i int
	q := new(CpabePolynomial)
	q.Deg = deg
	q.Coef = make([]*pbc.Element, deg+1)

	for i = 0; i < deg+1; i++ {
		q.Coef[i] = zeroVal.NewFieldElement().Set(zeroVal)
	}

	q.Coef[0].Set(zeroVal)

	for i = 1; i < deg+1; i++ {
		q.Coef[i].Rand()
	}

	return q
}

func fillPolicy(p *CpabePolicy, pm *CpabePm, e *pbc.Element) {
	var i int
	var r, t, h *pbc.Element
	pairing := pm.P
	r = pairing.NewZr()
	t = pairing.NewZr()
	h = pairing.NewG1()
	p.Q = randPoly(p.K-1, e)

	if p.Children == nil || len(p.Children) == 0 {
		p.C = pairing.NewG1()
		p.Cp = pairing.NewG1()

		elementFromString(h, p.Attr)
		p.C = pm.G.NewFieldElement().Set(pm.G)
		p.C.PowZn(p.C, p.Q.Coef[0])
		p.Cp = h.NewFieldElement().Set(h)
		p.Cp.PowZn(p.Cp, p.Q.Coef[0])
	} else {
		for i = 0; i < len(p.Children); i++ {
			r.SetInt32(int32(i + 1))
			evalPoly(t, p.Q, r)
			fillPolicy(p.Children[i], pm, t)
		}
	}
}

func evalPoly(r *pbc.Element, q *CpabePolynomial, x *pbc.Element) {
	var i int
	var s, t *pbc.Element

	s = r.NewFieldElement().Set(r)
	t = r.NewFieldElement().Set(r)

	r.Set0()
	t.Set1()

	for i = 0; i < q.Deg+1; i++ {
		/* r += q->coef[i] * t */
		s = q.Coef[i].NewFieldElement().Set(q.Coef[i])
		s.Mul(s, t)
		r.Add(r, s)

		/* t *= x */
		t.Mul(t, x)
	}

}

func baseNode(k int, s string) *CpabePolicy {
	p := new(CpabePolicy)

	p.K = k
	if !(s == "") {
		p.Attr = s
	} else {
		p.Attr = ""
	}
	p.Q = nil

	return p
}

func elementFromString(h *pbc.Element, s string) {
	sha := sha1.Sum([]byte(s))
	digest := sha[:]
	h.SetFromHash(digest)
}

