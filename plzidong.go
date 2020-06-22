package main

import (
	"bufio"
	"cpabse"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func set(pm *cpabse.CpabePm, msk *cpabse.CpabeMsk, ad string) {
	file, err := os.Open("plsc.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	num, _:= strconv.Atoi(ad)
	for scanner.Scan() {
		lineText := scanner.Text()
		if len(lineText) == 0{
			continue
		}
		pak := strings.Split(lineText, "||")

		policy := pak[0]
		keyword := pak[1]
		data := pak[2]

		c, _ := cpabse.CP_Enc(pm, policy, msk, keyword)
		c1 := strconv.Itoa(int(c[0]))
		for i := 1; i < len(c); i++ {
			temp := strconv.Itoa(int(c[i]))
			c1 += " "
			c1 += temp
		}

		ns := strconv.Itoa(num)
		comm := `peer chaincode invoke -n my -c '{"Args":["set","` + ns + `","` + c1 + `","` + string(data) + `"]}' -C myc`
		cmd := exec.Command("/bin/sh", "-c", comm)
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		num++
	}
	file.Close()
	fmt.Println("\033[32mUpload Complete\033[0m\n")
}

func main() {
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

	fmt.Println("Please input the first address")
	var ad string
	fmt.Scanln(&ad)
	set(pm, msk, ad)
}
