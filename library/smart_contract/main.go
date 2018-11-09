package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
	"bytes"
	"encoding/pem"
)

func ToChaincodergs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

type Chaincode struct {
}

type House struct {
	ID string `json:"id"`
	Location string `json:"location"`
	Owner string `json:"owner"`
}

//{"Args":["attr", "name"]}'
func (t *Chaincode) attr(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("parametes's number is wrong")
	}
	fmt.Println("get attr: ", args[0])
	value, ok, err := cid.GetAttributeValue(stub, args[0])
	if err != nil {
		return shim.Error("get attr error: " + err.Error())
	}

	if ok == false {
		value = "not found"
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return shim.Error("json marshal error: " + err.Error())
	}
	return shim.Success(bytes)
}

//{"Args":["creator2"]}'
func (t *Chaincode) creator2(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var cinfo struct {
		ID   string
		ORG  string
		CERT *x509.Certificate
	}

	fmt.Println("creator2: ", args)

	id, err := cid.GetID(stub)
	if err != nil {
		return shim.Error("getid error: " + err.Error())
	}

	id_readable, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return shim.Error("base64 decode error: " + err.Error())
	}
	cinfo.ID = string(id_readable)

	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error("getmspid error: " + err.Error())
	}
	cinfo.ORG = mspid

	cert, err := cid.GetX509Certificate(stub)
	if err != nil {
		return shim.Error("getX509Cert error: " + err.Error())
	}
	cinfo.CERT = cert

	bytes, err := json.Marshal(cinfo)
	if err != nil {
		return shim.Error("json marshal error: " + err.Error())
	}
	return shim.Success(bytes)
}

//{"Args":["creator"]}'
func (t *Chaincode) creator(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("creator: ", args)
	bytes, err := stub.GetCreator()
	if err != nil {
		return shim.Error("get creator error: " + err.Error())
	}

	// TODO: status: 500, message: unmarshal creator error: invalid character '\x19'
	//looking for beginning of value, bytes:
	/*
		var creator msp.SerializedIdentity
		if err := json.Unmarshal(bytes, &creator); err != nil {
			return shim.Error("unmarshal creator error: " + err.Error())
		}
	*/
	return shim.Success(bytes)
}

//{"Args":["call","chaincode","method"...]}'
func (t *Chaincode) call(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("call: ", args)
	sub_args := args[1:]
	return stub.InvokeChaincode(args[0], ToChaincodergs(sub_args...), stub.GetChannelID())
}

//{"Args":["append","key", ...]}'
func (t *Chaincode) append(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	key := args[0]
	value := args[1:]
	var data []string

	bytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("query " + key + " fail: " + err.Error())
	}

	if bytes != nil {
		if err := json.Unmarshal(bytes, &data); err != nil {
			return shim.Error(err.Error())
		}
	}

	data = append(data, value...)
	new_bytes, err := json.Marshal(data)
	if err != nil {
	}

	if err := stub.PutState(key, new_bytes); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//{"Args":["query_chaincode","chaincode","key"]}'
func (t *Chaincode) query_chaincode(stub shim.ChaincodeStubInterface, chaincode, key string) pb.Response {
	fmt.Printf("query %s in %s\n", key, chaincode)
	return stub.InvokeChaincode(chaincode, ToChaincodergs("query", key), stub.GetChannelID())
}

//{"Args":["write_chaincode","chaincode","key","value"]}'
func (t *Chaincode) write_chaincode(stub shim.ChaincodeStubInterface, chaincode, key, value string) pb.Response {
	fmt.Printf("write %s to %s, value is %s\n", key, chaincode, value)
	return stub.InvokeChaincode(chaincode, ToChaincodergs("write", key, value), stub.GetChannelID())
}

//{"Args":["query","key"]}'
func (t *Chaincode) query(stub shim.ChaincodeStubInterface, key string) pb.Response {
	fmt.Printf("query %s\n", key)
	bytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("query fail " + err.Error())
	}
	return shim.Success(bytes)
}

// GetQueryResult
func (t *Chaincode) queryResult(stub shim.ChaincodeStubInterface, key string) pb.Response {
	fmt.Printf("queryResult %s\n", key)
	iter, err := stub.GetQueryResult(key)
	if err != nil {
		return shim.Error("queryResult fail " + err.Error())
	}
	defer iter.Close()
	if err != nil {
		return shim.Error("queryResult fail " + err.Error())
	}

	values := make(map[string]string)

	for iter.HasNext() {
		fmt.Println("next \n")
		if kv, err := iter.Next(); err == nil {
			fmt.Println("id: %s value: %s namespace: %s\n", kv.Key, kv.Value, kv.Namespace)
			values[kv.Key] = string(kv.Value)
		}
		if err != nil {
			return shim.Error("iterator queryResult fail: " + err.Error())
		}
	}
	bytes, err := json.Marshal(values)
	if err != nil {
		return shim.Error("json marshal fail: " + err.Error())
	}
	return shim.Success(bytes)
}

// {"Args":["del","key"]}'
func (t *Chaincode) del(stub shim.ChaincodeStubInterface, key string) pb.Response {
	fmt.Printf("del %s\n", key)
	err := stub.DelState(key)
	if err != nil {
		return shim.Error("del fail " + err.Error())
	}
	return shim.Success(nil)
}

func (t *Chaincode) getByRange(stub shim.ChaincodeStubInterface, startKey, endKey string) pb.Response {
	fmt.Printf("getByRange %s - %s\n", startKey, endKey)
	iter, err := stub.GetStateByRange(startKey, endKey)
	defer iter.Close()
	if err != nil {
		return shim.Error("getByRange fail " + err.Error())
	}

	values := make(map[string]string)

	for iter.HasNext() {
		fmt.Println("next \n")
		if kv, err := iter.Next(); err == nil {
			fmt.Println("id: %s value: %s namespace: %s\n", kv.Key, kv.Value, kv.Namespace)
			values[kv.Key] = string(kv.Value)
		}
		if err != nil {
			return shim.Error("iterator getByRange fail: " + err.Error())
		}
	}
	bytes, err := json.Marshal(values)
	if err != nil {
		return shim.Error("json marshal fail: " + err.Error())
	}

	return shim.Success(bytes)

}


//{"Args":["history","key"]}'
func (t *Chaincode) history(stub shim.ChaincodeStubInterface, key string) pb.Response {
	fmt.Printf("history %s\n", key)
	iter, err := stub.GetHistoryForKey(key)
	defer iter.Close()
	if err != nil {
		return shim.Error("query fail " + err.Error())
	}

	values := make(map[string]string)

	for iter.HasNext() {
		fmt.Printf("next\n")
		if kv, err := iter.Next(); err == nil {
			fmt.Printf("id: %s value: %s\n", kv.TxId, kv.Value)
			values[kv.TxId] = string(kv.Value)
		}
		if err != nil {
			return shim.Error("iterator history fail: " + err.Error())
		}
	}

	bytes, err := json.Marshal(values)
	if err != nil {
		return shim.Error("json marshal fail: " + err.Error())
	}

	return shim.Success(bytes)
}

//{"Args":["write","key","value"]}'
func (t *Chaincode) write(stub shim.ChaincodeStubInterface, key, value string) pb.Response {
	fmt.Printf("new felix  write %s, value is %s\n", key, value)
	if err := stub.PutState(key, []byte(value)); err != nil {
		return shim.Error("write fail " + err.Error())
	}
	return shim.Success(nil)
}

//{"Args":["write","key","value"]}'
func (t *Chaincode) writeHouse(stub shim.ChaincodeStubInterface, id, location, owner string) pb.Response {
	key := id
	fmt.Printf("writeHouse  write %s, value is %s\n", key, id)
	var house = House{id, location, owner}
	houseBytes, _ := json.Marshal(house)
	if err := stub.PutState(key, houseBytes); err != nil {
		return shim.Error("write fail " + err.Error())
	}
	return shim.Success(nil)
}

// 隐私数据 1.2 新特性
func (t *Chaincode) writePrivateData(stub shim.ChaincodeStubInterface, collection, key, value string) pb.Response {
	fmt.Println("write private data collection s%, key is s%, value is s% \n", collection,  key, value)
	if err := stub.PutPrivateData(collection, key, []byte(value)); err != nil {
		return shim.Error("write fail " + err.Error())
	}
	return shim.Success(nil)
}

// 隐私数据 1.2 新特性
func (t *Chaincode) getPrivateData(stub shim.ChaincodeStubInterface, collection, key string) pb.Response {
	fmt.Println("get private data collection s%, key is s% \n", collection,  key)
	bytes, err := stub.GetPrivateData(collection, key)
	if err != nil {
		return shim.Error("query fail " + err.Error())
	}
	return shim.Success(bytes)
}

func (t *Chaincode) getPrivateByRange(stub shim.ChaincodeStubInterface, collection, startKey, endKey string) pb.Response {
	fmt.Printf("getCollectionByRange s% - %s - %s\n", collection, startKey, endKey)
	iter, err := stub.GetPrivateDataByRange(collection, startKey, endKey)
	defer iter.Close()
	if err != nil {
		return shim.Error("getByRange fail " + err.Error())
	}

	values := make(map[string]string)

	for iter.HasNext() {
		fmt.Println("next \n")
		if kv, err := iter.Next(); err == nil {
			fmt.Println("id: %s value: %s namespace: %s\n", kv.Key, kv.Value, kv.Namespace)
			values[kv.Key] = string(kv.Value)
		}
		if err != nil {
			return shim.Error("iterator getByRange fail: " + err.Error())
		}
	}
	bytes, err := json.Marshal(values)
	if err != nil {
		return shim.Error("json marshal fail: " + err.Error())
	}

	return shim.Success(bytes)
}


//{"Args":["init"]}
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Init Chaincode Chaincode")
	return shim.Success(nil)
}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	//返回调用者信息
	case "creator":
		return t.creator(stub, args)
	//返回调用者信息，方法2
	case "creator2":
		return t.creator2(stub, args)
	//调用改合约中的其它方法，用来演示复杂的调用
	case "call":
		return t.call(stub, args)
	//直接对key的内容进行append，用来演示这样操作的结果
	case "append":
		return t.append(stub, args)
	//读取当前用户的属性值
	case "attr":
		return t.attr(stub, args)
	//查询一个key的当前值
	case "query":
		if len(args) != 1 {
			return shim.Error("parametes's number is wrong")
		}
		return t.query(stub, args[0])
	//查询一个key的所有历史值
	case "history":
		if len(args) != 1 {
			return shim.Error("parametes's number is wrong")
		}
		return t.history(stub, args[0])
	//创建一个key，并写入key的值
	case "write": //写入
		if len(args) != 2 {
			return shim.Error("parametes's number is wrong")
		}
		return t.write(stub, args[0], args[1])
		//创建一个key，并写入key的值 House，处理比write复杂
	case "writeHouse": //写入
		if len(args) != 3 {
			return shim.Error("writeHouse parametes's number is wrong")
		}
		return t.writeHouse(stub, args[0], args[1], args[2])

		//创建一个key隐私数据，并写入key的值
	case "writePrivateData": //写入
		if len(args) != 3 {
			return shim.Error("parametes's number is wrong")
		}
		return t.writePrivateData(stub, args[0], args[1], args[2])
	// 查询隐私数据
	case "getPrivateData":
		if len(args) != 2 {
			return shim.Error("parametes's number is wrong")
		}
		return t.getPrivateData(stub, args[0], args[1])
	// 隐私数据
	case "getPrivateByRange":
		if len(args) != 3 {
			return shim.Error("parametes's number is wrong")
		}
		return t.getPrivateByRange(stub, args[0], args[1], args[2])

		//删除一个key
	case "del":
		if len(args) != 1 {
			return shim.Error("parametes's number is wrong")
		}
		return t.del(stub, args[0])

	// getByRange
	case "getByRange":
		if len(args) != 2 {
			return shim.Error("parametes's number is wrong")
		}
		return t.getByRange(stub, args[0], args[1])

	//queryResult
	case "queryResult":
		if len(args) != 1 {
			return shim.Error("parametes's number is wrong")
		}
		return t.queryResult(stub, args[0])
	// testCertificate
	case "testCertificate":
		return t.testCertificate(stub)


		//通过当前合约，到另一个合约中进行查询
	case "query_chaincode":
		if len(args) != 2 {
			return shim.Error("parametes's number is wrong")
		}
		return t.query_chaincode(stub, args[0], args[1])
	//通过当前合约，到另一个合约中进行写入
	case "write_chaincode":
		if len(args) != 3 {
			return shim.Error("parametes's number is wrong")
		}
		return t.write_chaincode(stub, args[0], args[1], args[2])
	default:
		return shim.Error("Invalid invoke function name.")
	}
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting Chaincode chaincode: %s", err)
	}
}



func (t *Chaincode) testCertificate(stub shim.ChaincodeStubInterface) pb.Response{
	creatorByte,_:= stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		fmt.Errorf("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Errorf("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Errorf("ParseCertificate failed")
	}
	uname:=cert.Subject.CommonName
	fmt.Println("Name:"+uname)
	fmt.Println("all-----:", cert)
	return shim.Success(bl.Bytes)
}