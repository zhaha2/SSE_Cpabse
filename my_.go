package main

import (
	"bytes"
	"cpabse"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Nik-U/pbc"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var curveParams = "type a\n" +
	"q 87807107996633125224377819847540498158068831994142082" +
	"1102865339926647563088022295707862517942266222142315585" +
	"8769582317459277713367317481324925129998224791\n" +
	"h 12016012264891146079388821366740534204802954401251311" +
	"822919615131047207289359704531102844802183906537786776\n" +
	"r 730750818665451621361119245571504901405976559617\n" +
	"exp2 159\n" + "exp1 107\n" + "sign1 1\n" + "sign0 1\n"

type MyChaincode struct {
}

//实现Init函数，初始化
func (t *MyChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("my init")
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

func (t *MyChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke")
	var result []byte
	var err error
	if fn == "set" {
		fmt.Println("set")
		result, err = upload(stub, args)
	} else if fn == "query" {
		fmt.Println("query")
		result, err = query(stub, args)
	} else if fn == "chainquery" {
		fmt.Println("chainquery")
		result, err = chainquery(stub, args[0])		//arg[0]为地址（Key1, Key2, ...）
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success(result)
}

func upload(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	Key := "Key" + args[0]
	//args[0]为id，即计数器n，args[1]为key，关键字，args[2]为值，即关键字对应的数据
	//此处的key为计数器n和kw的复合键，即最后字典中键值对的键
	var key, _ = stub.CreateCompositeKey(Key, []string{args[0], args[1]})
	if len(args) != 3 {
		return []byte(""), fmt.Errorf("Incorrect arguments. Expecting ID ,a key and a value")
	}

	err := stub.PutState(key, []byte(args[2]))
	if err != nil {
		return []byte(""), fmt.Errorf("Failed to set asset: %s", args[0])
	}
	fmt.Printf("\033[32m%s\033[0m\n", "id: " + args[0] + " data: " + args[2])
	fmt.Println()
	return []byte(args[0]), nil
}

func query(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	start := time.Now()		//开始时间，一会计时用

	var p *pbc.Pairing
	params := new(pbc.Params)
	params, _ = pbc.NewParamsFromString(curveParams)
	p = pbc.NewPairing(params)
	if len(args) != 1 {
		return []byte(""), fmt.Errorf("Incorrect arguments. Expecting a token")
	}

	tk := cpabse.TkDec(args[0], p)

	//用符合键中一个键进行查询，这里用kw查询
	queryResultsIterator, err := stub.GetStateByPartialCompositeKey("Key", []string{})
	if err != nil {
		return []byte(""), fmt.Errorf("Incorrect")
	}
	defer queryResultsIterator.Close()

	var result []string

	//返回结果
	for queryResultsIterator.HasNext() {

		responseRange, err := queryResultsIterator.Next()
		if err != nil {
			return []byte(""), fmt.Errorf("Incorrect")
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return []byte(""), fmt.Errorf("Incorrect")
		}

		cph := cpabse.CphDec(compositeKeyParts[1], p)
		if cpabse.Check(tk, cph, p) {
			result = append(result, string(responseRange.Value))
			fmt.Println(string(responseRange.Value))
		} else {
			continue
		}
	}
	buf2 := &bytes.Buffer{}
	gob.NewEncoder(buf2).Encode(result)
	byteSlice := []byte(buf2.Bytes())

	//结束时间
	end := time.Now()
	//输出总时间
	b := end.Sub(start)
	fmt.Println("Query time cost: %s", b)
	return byteSlice, nil
}

func chainquery(stub shim.ChaincodeStubInterface, Key string) ([]byte, error) {
	//开始时间，一会计时用
	start := time.Now()
	fmt.Printf("\033[33m%s\033[0m\n", "Key: " + Key)
	//用符合键中一个键进行查询，这里用地址Key查询
	queryResultsIterator, err := stub.GetStateByPartialCompositeKey(Key, []string{})
	if err != nil {
		return []byte(""), fmt.Errorf("Incorrect")
	}
	defer queryResultsIterator.Close()

	var result []string

	//返回结果
	for queryResultsIterator.HasNext() {

		responseRange, err := queryResultsIterator.Next()
		if err != nil {
			return []byte(""), fmt.Errorf("Incorrect")
		}
		fmt.Println(responseRange.Value)
		//分离数据和地址
		Value := strings.Split(string(responseRange.Value), "::")
		//输出结果（数据）
		result = append(result, Value[0])
		fmt.Printf("\033[34mResult: %s\033[0m\n", Value[0])
		fmt.Println()

		//如果没有地址则为链尾，停止搜索，否则继续搜索文件链
		if len(Value) != 2 {
			break
		} else {
			var r []string
			res, _ := chainquery(stub, Value[1])
			r = append(r, strconv.Itoa(int(res[0])))
			for i := 1; i < len(res); i++ {
				r = append(r, strconv.Itoa(int(res[i])))
			}
			result = append(result, r...)
		}
	}
	buf2 := &bytes.Buffer{}
	gob.NewEncoder(buf2).Encode(result)
	byteSlice := []byte(buf2.Bytes())

	end := time.Now()
	//输出总时间
	b := end.Sub(start)
	fmt.Printf("Chainquery time cost: %s \n", b)
	return byteSlice, nil
}

func main() {
	err := shim.Start(new(MyChaincode))
	if err != nil {
		fmt.Printf("error start MyChaincode")
	}
}
