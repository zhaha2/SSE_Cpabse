package main

import (
	"bufio"
	"cpabse"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//IndexGen
func set(pm *cpabse.CpabePm, msk *cpabse.CpabeMsk, num int) {
	var policy, keyword, data string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter your policy:")
	scanner.Scan()
	policy = scanner.Text()
	fmt.Println("Please enter your keyword:")
	scanner.Scan()
	keyword = scanner.Text()
	fmt.Println("Please enter your data:")
	scanner.Scan()
	data = scanner.Text()

	//c为cph?
	c, _ := cpabse.CP_Enc(pm, policy, msk, keyword)
	fmt.Println(c)
	//c1为字符串形式的c
	c1 := strconv.Itoa(int(c[0]))
	for i := 1; i < len(c); i++ {
		temp := strconv.Itoa(int(c[i]))
		c1 += " "
		c1 += temp
	}
	//ns为main()中的n，计数器（地址）
	ns := strconv.Itoa(num)

	//让chiancode执行操作
	comm := `peer chaincode invoke -n my -c '{"Args":["set","` + ns + `","` + c1 + `","` + string(data) + `"]}' -C myc`
	cmd := exec.Command("/bin/sh", "-c", comm)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	fmt.Printf("\033[32m%s\033[0m", "Upload Complete!\n")
}

func query(pm *cpabse.CpabePm, msk *cpabse.CpabeMsk) {
	var attrs, keyword string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter you attrs:")
	scanner.Scan()
	attrs = scanner.Text()
	fmt.Println("Please enter your keyword:")
	scanner.Scan()
	keyword = scanner.Text()

	prv := cpabse.CP_Keygen(pm, msk, attrs)
	//tocken
	t, _ := cpabse.CP_TkEnc(prv, keyword, msk, pm)
	t1 := strconv.Itoa(int(t[0]))
	for i := 1; i < len(t); i++ {
		temp := strconv.Itoa(int(t[i]))
		t1 += " "
		t1 += temp
	}

	comm := `peer chaincode invoke -n my -c '{"Args":["query","` + t1 + `"]}' -C myc`
	cmd := exec.Command("/bin/sh", "-c", comm)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	fmt.Printf("\033[32m%s\033[0m", "Query Complete!\n")

}

func query1() {
	var Key string
	scanner := bufio.NewScanner(os.Stdin)
	//Key为地址，即上传时的num，这里只能用这个搜
	fmt.Println("Please enter you Key:")
	scanner.Scan()
	Key = scanner.Text()
	//Key字符转换为数字k
	k, _ := strconv.Atoi(Key)
	//这个循环为了模拟多次执行？，这里设置为不循环
	for i := k; i < k+1; i++ {
		ks := strconv.Itoa(i)
		//因为加密时地址为Key1，Key2...，这里也做成这个形式
		k := "Key" + ks

		comm := `peer chaincode invoke -n my -c '{"Args":["query","` + k + `"]}' -C myc`
		cmd := exec.Command("/bin/sh", "-c", comm)
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		fmt.Printf("\033[32m%s\033[0m", "Query Complete!\n")
	}

}


func chainquery() {
	var Key string
	scanner := bufio.NewScanner(os.Stdin)
	//Key为地址，即上传时的num，这里只能用这个搜
	fmt.Println("Please enter you Key:")
	scanner.Scan()
	Key = scanner.Text()
	//Key字符转换为数字k
	k, _ := strconv.Atoi(Key)
	//这个循环为了模拟多次执行？，这里设置为不循环
	for i := k; i < k+1; i++ {
		ks := strconv.Itoa(i)
		//因为加密时地址为Key1，Key2...，这里也做成这个形式
		k := "Key" + ks

		comm := `peer chaincode invoke -n my -c '{"Args":["chainquery","` + k + `"]}' -C myc`
		cmd := exec.Command("/bin/sh", "-c", comm)
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		fmt.Printf("\033[32m%s\033[0m", "Chainquery Complete!\n")
	}
}

//一次搜索多个关键字
func batchq() {

	//n为要返回多少结果
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter the number of files you want to get:")
	scanner.Scan()
	num := scanner.Text()

	//搜索30次
	for i := 0; i < 30; i++ {
		comm := `peer chaincode invoke -n my -c '{"Args":["batchq","` + num + `"]}' -C myc`
		cmd := exec.Command("/bin/sh", "-c", comm)
		cmd.Stdout = os.Stdout
		_ = cmd.Run()


		fmt.Printf("\033[32m%s\033[0m", "BatchChainquery Complete\n")
	}
}

//一次搜索多个关键字
func batchcq() {

	//将kwStartIndex读入本地保存
	kwStartIndex, err := ioutil.ReadFile("kwStartIndex.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	kwStartIndexList := strings.Split(string(kwStartIndex), " ")

	//要返回多少结果
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter the number of files you want to get:")
	scanner.Scan()
	n, err1 := strconv.Atoi(scanner.Text())
	if err1 != nil {
		fmt.Println("Input error!", err)
		return
	}

	//搜索30次
	for i := 0; i < 30; i++ {
		rand.Seed(int64(i))
		searchkey := ""
		//每次搜索n/5个关键字，共返回n个结果
		for j := 0; j < (n / 10); j++{
			searchkey = searchkey + kwStartIndexList[rand.Intn(1000)] + " "
		}
		comm := `peer chaincode invoke -n my -c '{"Args":["batchcq","` + string(searchkey) + `"]}' -C myc`
		cmd := exec.Command("/bin/sh", "-c", comm)
		cmd.Stdout = os.Stdout
		_ = cmd.Run()

		fmt.Printf("\033[32m%s\033[0m", "BatchChainquery Complete\n")
	}
}

func main() {
	//生成pm
	bpm := new(cpabse.BytePm)
	f, _ := os.Open("pm.txt")
	dec := gob.NewDecoder(f)
	_ = dec.Decode(&bpm)
	pm := new(cpabse.CpabePm)
	cpabse.Psetup(pm)
	cpabse.BpmToPm(pm, bpm)

	//生成msk
	bmsk := new(cpabse.ByteMsk)
	f1, _ := os.Open("msk.txt")
	dec1 := gob.NewDecoder(f1)
	_ = dec1.Decode(&bmsk)
	msk := new(cpabse.CpabeMsk)
	cpabse.BmskToMsk(msk, bmsk, pm)

	fmt.Println("Please choose the operation: upload/query/chainquery(U/Q/CQ/BQ/BCQ)")
	var c string
	fmt.Scanln(&c)
	n := 1	//n只有在 U 时才++，注意不要和plzidong的地址冲突
	if c == "U" || c =="u" {
		set(pm, msk, n)
		n++
	} else if c == "Q" || c =="q" {
		query1()
	} else if c == "CQ" || c =="cq" {	//文件链查询
		chainquery()
	} else if c == "BQ" || c =="bq" {
		batchq()
	}else if c == "BCQ" || c =="bcq" {	//文件链查询
		batchcq()
	} else {
		fmt.Println("error input")
		return
	}
}
