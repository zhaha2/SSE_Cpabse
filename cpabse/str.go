package cpabse

import "github.com/Nik-U/pbc"

type CpabePm struct {
/*
 * A public key
 */
PairingDesc string       `json:"pairing_desc"` //a.param
P           *pbc.Pairing `json:"p"`            //e
G           *pbc.Element `json:"g"`            /* G_1 */
//r *pbc.Element `json:"r"`
G_a     *pbc.Element `json:"g_a"`
G_b     *pbc.Element `json:"g_b"`
G_c     *pbc.Element `json:"g_c"` /* G_T */ //e(g,g)^alpha
G_y     *pbc.Element `json:"g_y"`
G_hat_x *pbc.Element `json:"g_hat_x"` //e(g^x,g)

}

type CpabeMsk struct {
/*
 * A master secret key
 */
A *pbc.Element `json:"a"` /* Z_r */
B *pbc.Element `json:"b"` /* Z_r */
C *pbc.Element `json:"c"` /* Z_r */
X *pbc.Element `json:"x"`
Y *pbc.Element `json:"y"`
}

type CpabeSk struct {
/*
 * A private key
 */
B     *pbc.Element
A     *pbc.Element
Comps []*CpabeSkComp /* CpabeSkComp */
R     *pbc.Element

D *pbc.Element
}

type CpabeSkComp struct {
Attr string
A_j  *pbc.Element
B_j  *pbc.Element
}

type CpabePolicy struct {
/* k=1 if leaf, otherwise threshould */
K int
/* attribute string if leaf, otherwise null */
Attr string
C    *pbc.Element /* G_1 only for leaves */ //W
Cp   *pbc.Element /* G_1 only for leaves */ //D
/* array of BswabePolicy and length is 0 for leaves */
Children []*CpabePolicy

/* only used during encryption */
Q *CpabePolynomial

/* only used during decription */
Satisfiable bool
Min_leaves  int
Attri       int
Satl        []int
}

type CpabeCphKey struct {
/*
 * This class is defined for some classes who return both cph and key.
 */
Cph *CpabeCph
//Key *pbc.Element
ciphertext []byte
}

type CpabeCph struct {
/*
 * A ciphertext. Note that this library only handles encrypting a single
 * group element, so if you want to encrypt something bigger, you will have
 * to use that group element as a symmetric key for hybrid encryption (which
 * you do yourself).
 */
W  *pbc.Element /* G_T */
W0 *pbc.Element /* G_1 */
W1 *pbc.Element
//S  *pbc.Element
P *CpabePolicy

//cs *pbc.Element 		/* G_T */
//c *pbc.Element	 		/* G_1 */
}

type CpabePolynomial struct {
Deg int
/* coefficients from [0] x^0 to [deg] x^deg */
Coef []*pbc.Element /* G_T (of length deg+1) */
}

type Token struct {
Tok1  *pbc.Element /*(g^a*g^bH2(w))^s*/
Tok2  *pbc.Element
Tok3  *pbc.Element
Comps []*TokenComp
}

type TokenComp struct {
Attr string
A_j1 *pbc.Element
B_j1 *pbc.Element
}

