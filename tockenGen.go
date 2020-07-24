package main

import (
	"cpabse"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"time"
)

func tockenGen(pm *cpabse.CpabePm, msk *cpabse.CpabeMsk)(int64){
	var attrs, keyword string
	attrs = "baf"
	keyword = "lbw"

	start := time.Now()
	prv := cpabse.CP_Keygen(pm, msk, attrs)
	//tocken
	_, _ = cpabse.CP_TkEnc(prv, keyword, msk, pm)

	cost := time.Since(start)

	return cost.Nanoseconds()

}

func main()  {
	bpm := new(cpabse.BytePm)
	f, _ := os.Open("pm.txt")
	dec := gob.NewDecoder(f)
	_ = dec.Decode(&bpm)
	pm := new(cpabse.CpabePm)
	cpabse.Psetup(pm)
	cpabse.BpmToPm(pm, bpm)

	bmsk := new(cpabse.ByteMsk)
	f1, _ := os.Open("msk.txt")
	dec1 := gob.NewDecoder(f1)
	_ = dec1.Decode(&bmsk)
	msk := new(cpabse.CpabeMsk)
	cpabse.BmskToMsk(msk, bmsk, pm)


	sum := int64(0)
	fmt.Println("Please input the number of tockens:")
	var n string
	fmt.Scanln(&n)
	num, _ := strconv.Atoi(n)
	s := int64(0)

	for j:=0; j<30; j++ {
		sum = int64(0)
		for i := 0; i < num; i++ {
			sum += tockenGen(pm, msk)
		}
		s += sum
	}
	fmt.Print(float64(s)/30 * 0.000001)
	fmt.Println("ms")

	f.Close()
	f1.Close()

}

