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
} //实现Init函数

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
<<<<<<< HEAD
		ret, err = chainquery(stub, args[0])
	}
	if r != nil {
=======
		result, err = chainquery(stub, args[0])
	}
	if err != nil {
>>>>>>> c013c505f79be6fe06cff2bcc6c61cc29642d753
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success(result)
}

func upload(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	Key := "Key" + args[0]
	var key, _ = stub.CreateCompositeKey(Key, []string{args[0], args[1]})
	if len(args) != 3 {
		return []byte(""), fmt.Errorf("Incorrect arguments. Expecting ID ,a key and a value")
	}

	err := stub.PutState(key, []byte(args[2]))
	if err != nil {
		return []byte(""), fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return []byte(args[0]), nil
}

func query(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	start := time.Now()
	var p *pbc.Pairing
	params := new(pbc.Params)
	params, _ = pbc.NewParamsFromString(curveParams)
	p = pbc.NewPairing(params)
	if len(args) != 1 {
		return []byte(""), fmt.Errorf("Incorrect arguments. Expecting a token")
	}

	tk := cpabse.TkDec(args[0], p)

	queryResultsIterator, err := stub.GetStateByPartialCompositeKey("Key", []string{})
	if err != nil {
		return []byte(""), fmt.Errorf("Incorrect")
	}
	defer queryResultsIterator.Close()

	var result []string

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
	end := time.Now()
	b := end.Sub(start)
	fmt.Println(b)
	return byteSlice, nil
}

func chainquery(stub shim.ChaincodeStubInterface, Key string) ([]byte, error) {
	start := time.Now()
	fmt.Println(Key)
	queryResultsIterator, err := stub.GetStateByPartialCompositeKey(Key, []string{})
	if err != nil {
		return []byte(""), fmt.Errorf("Incorrect")
	}
	defer queryResultsIterator.Close()

	var result []string

	for queryResultsIterator.HasNext() {

		responseRange, err := queryResultsIterator.Next()
		if err != nil {
			return []byte(""), fmt.Errorf("Incorrect")
		}
		fmt.Printf("\033[1;31;40m%s\033[0m\n", responseRange.Value)
		Value := strings.Split(string(responseRange.Value), "::")
		result = append(result, Value[0])
		fmt.Println(Value[0])
		fmt.Println()
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
	b := end.Sub(start)
	fmt.Println(b)
	return byteSlice, nil
}

func main() {
	err := shim.Start(new(MyChaincode))
	if err != nil {
		fmt.Printf("error start MyChaincode")
	}
}
