package cpabse

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"

	"github.com/Nik-U/pbc"
)

func CP_Keygen(pk *CpabePm, mk *CpabeMsk, attrs string) *CpabeSk {
	attr := strings.Split(attrs, " ")
	sk := Keygen(pk, mk, attr)
	return sk
}

func CP_Enc(pm *CpabePm, policy string, msk *CpabeMsk, w string) ([]byte, *CpabeCph) {
	cph := EncKeyword(pm, policy, msk, w)
	bcph := CphToByte(cph)
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err1 := enc.Encode(bcph); err1 != nil {
		fmt.Println(err1)
	}
	//fmt.Printf("序列化后：%x\n", buffer.Bytes())
	c := []byte(buffer.Bytes())

	return c, cph
}

func CphDec(c string, p *pbc.Pairing) *CpabeCph {
	var c1 []uint8
	t := strings.Fields(c)
	for _, v := range t {
		s, _ := strconv.Atoi(v)
		temp := uint8(s)
		c1 = append(c1, temp)
	}

	bcph := new(ByteCph)
	var buffer bytes.Buffer
	buffer.Write(c1)
	F := gob.NewDecoder(&buffer)
	F.Decode(&bcph)
	cph := ByteToCph(bcph, p)
	return cph
}

func CP_TkEnc(sk *CpabeSk, w string, msk *CpabeMsk, pm *CpabePm) ([]byte, *Token) {
	tk := TokenGen(sk, w, msk, pm)
	btk := TkToByte(tk)
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err1 := enc.Encode(btk); err1 != nil {
		fmt.Println(err1)
	}
	t := []byte(buffer.Bytes())

	return t, tk
}

func TkDec(t string, p *pbc.Pairing) *Token {
	var t1 []uint8
	s := strings.Fields(t)
	for _, v := range s {
		k, _ := strconv.Atoi(v)
		temp := uint8(k)
		t1 = append(t1, temp)
	}

	btk := new(ByteTk)
	var buffer bytes.Buffer
	buffer.Write(t1)
	F := gob.NewDecoder(&buffer)
	F.Decode(&btk)
	tk := ByteToTk(btk, p)
	return tk
}

