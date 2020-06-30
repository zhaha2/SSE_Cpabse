package main

import (
	"bytes"
	"cpabse"
	"encoding/gob"
	"fmt"
	"os"
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
		result, err, _ = chainquery(stub, args[0])         //arg[0]为地址（Key1, Key2, ...）
	} else if fn == "batchcq" {
		fmt.Println("batchcq")
		var totalTime = int64(0)
		result, err, totalTime = batchcq(stub, args[0])         //arg[0]为地址串，以空格分开（Key1, Key2, ...）
		tt := strconv.FormatFloat(float64(totalTime) * 0.000001, 'f', 2,64)
		fmt.Println("Total time: " + tt +" ms")

		//将搜索总时间写入文件
		f, err := os.OpenFile("totalTime.txt", os.O_CREATE |os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			fmt.Println(err)
			return shim.Error(err.Error())
		}
		_, err = fmt.Fprint(f, tt + " ")
		if err != nil {
			fmt.Println(err)
			f.Close()
			return shim.Error(err.Error())
		}
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return shim.Error(err.Error())
		}
		fmt.Println("Write file complete!")
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
	start := time.Now()             //开始时间，一会计时用

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

func chainquery(stub shim.ChaincodeStubInterface, Key string) ([]byte, error, int64) {
	//开始时间，一会计时用
	start := time.Now()
	fmt.Printf("\033[33m%s\033[0m\n", "Key: " + Key)
	//用符合键中一个键进行查询，这里用地址Key查询
	queryResultsIterator, err := stub.GetStateByPartialCompositeKey(Key, []string{})
	if err != nil {
		return []byte(""), fmt.Errorf("Incorrect"), 0
	}
	defer queryResultsIterator.Close()

	var result []string

	//返回结果
	for queryResultsIterator.HasNext() {

		responseRange, err := queryResultsIterator.Next()
		if err != nil {
			return []byte(""), fmt.Errorf("Incorrect"), 0
		}
		fmt.Println(responseRange.Value)
		//分离数据和地址
		Value := strings.Split(string(responseRange.Value), "::")
		//输出结果（数据）
		result = append(result, Value[0])
		fmt.Printf("\033[34mResult: %s\033[0m\n", Value[0])
		fmt.Println()
		if len(Value) == 1{
			break
		}
		key_next := Value[1]

		//如果没有地址则为链尾，停止搜索，否则继续搜索文件链
		for {
			fmt.Printf("\033[33m%s\033[0m\n", "Key: " + key_next)
			//用复合键中一个键进行查询，这里用地址Key查询
			queryResultsIterator1, err := stub.GetStateByPartialCompositeKey(key_next, []string{})
			if err != nil {
				return []byte(""), fmt.Errorf("Incorrect"), 0
			}
			defer queryResultsIterator1.Close()

			responseRange, err = queryResultsIterator1.Next()
			if err != nil {
				return []byte(""), fmt.Errorf("Incorrect"), 0
			}
			fmt.Println(responseRange.Value)
			//分离数据和地址
			Value := strings.Split(string(responseRange.Value), "::")
			//输出结果（数据）
			result = append(result, Value[0])

			//找到链尾 搜索结束
			if len(Value) == 1{
				fmt.Printf("\033[34mResult: %s\033[0m\n", Value[0])
				fmt.Println()
				break
			}

			//否则，继续按地址搜索
			key_next = Value[1]
			fmt.Printf("\033[34mResult: %s\033[0m\n", Value[0])
			fmt.Println()
		}

		//if len(Value) != 2 {
		//      break
		//} else {
		//      var r []string
		//      res, _ := chainquery(stub, Value[1])
		//      r = append(r, strconv.Itoa(int(res[0])))
		//      for i := 1; i < len(res); i++ {
		//              r = append(r, strconv.Itoa(int(res[i])))
		//      }
		//      result = append(result, r...)
		//}
	}
	queryResultsIterator.Close()
	buf2 := &bytes.Buffer{}
	gob.NewEncoder(buf2).Encode(result)
	byteSlice := []byte(buf2.Bytes())

	end := time.Now()
	//输出总时间
	b := end.Sub(start)
	fmt.Printf("Chainquery time cost: %s \n", b)
	return byteSlice, nil, b.Nanoseconds()
}

//批量搜索，统计时间，这里搜索时间不计写入账本所需时间
func batchcq(stub shim.ChaincodeStubInterface, Key string) ([]byte, error, int64) {
	keyList := strings.Split(Key, " ")
	toatalTime := int64(0)
	for _, v := range keyList{
		_, err, timecost := chainquery(stub, "Key" + v)
		toatalTime += timecost
		if err != nil {
			return []byte(""), fmt.Errorf("Incorrect"), 0
		}
	}
	return []byte(""), nil, toatalTime
}

func main() {
	err := shim.Start(new(MyChaincode))
	if err != nil {
		fmt.Printf("error start MyChaincode")
	}
}
