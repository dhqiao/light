package chain

import (
	. "light/http/controllers"
	"github.com/gin-gonic/gin"
	"light/service/chain"
	"fmt"
	"bytes"
	"encoding/pem"
	"crypto/x509"
	"io/ioutil"
)

func Get(c *gin.Context) {
	key := c.Query("key")
	blockChain := &chain.BlockChain{}
	fmt.Println(">>>>>>>>>>>>>>>>>", key)
	rst, err := blockChain.Get(key)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func Set(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")
	blockChain := &chain.BlockChain{}
	fmt.Println("=-----------set value-----------", key, value)
	rst, err := blockChain.Set(key, value)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func History(c *gin.Context) {
	key := c.Query("key")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.History(key)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

// 删除键
func Del(c *gin.Context) {
	key := c.Query("key")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.Del(key)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

// 获取两个key之间的数据
func GetByRange(c *gin.Context) {
	startKey := c.Query("startKey")
	endKey := c.Query("endKey")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.GetByRange(startKey, endKey)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}


func QueryResult(c *gin.Context) {
	key := c.Query("key")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryResult(key)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func QueryBlock(c *gin.Context) {
	blockNum := c.Query("blockNum")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryBlock(blockNum)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func QueryTransaction(c *gin.Context) {
	txID := c.Query("txid")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryTransaction(txID)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}


func QueryBlockByTxID(c *gin.Context) {
	txID := c.Query("txid")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>", txID)
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryBlockByTxID(txID)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func Test(c *gin.Context) {
	txID := c.Query("txid")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>", txID)
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.Test(txID)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func QueryBlockByHash(c *gin.Context)  {
	hash := c.Query("hash")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryBlockByHash(hash)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

// 获取配置信息 只能到order节点和锚点peer
func QueryConfig(c *gin.Context)  {
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryConfig()
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

// 获取区块总数量，最后一个block，倒数第二个block hash， 和背书节点
func QueryInfo(c *gin.Context)  {
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryInfo()
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

// 获取所有的成员，我们没有成员管理
func GetAllIdentities(c *gin.Context)  {
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.GetAllIdentities()
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}



func GetPrivate(c *gin.Context) {
	collection := c.Query("collection")
	key := c.Query("key")
	blockChain := &chain.BlockChain{}
	fmt.Println(">>>>>>>>>>>>>>>>>", key)
	rst, err := blockChain.GetPrivate(collection, key)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func SetPrivate(c *gin.Context) {
	collection := c.Query("collection")
	key := c.Query("key")
	value := c.Query("value")
	blockChain := &chain.BlockChain{}
	fmt.Println("=-----------set value-----------", key, value)
	rst, err := blockChain.SetPrivate(collection, key, value)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}


// 获取两个key之间的数据
func GetPrivateByRange(c *gin.Context) {
	collection := c.Query("collection")
	startKey := c.Query("startKey")
	endKey := c.Query("endKey")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.GetPrivateByRange(collection, startKey, endKey)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func WriteHouse(c *gin.Context)  {
	id := c.Query("id")
	location := c.Query("location")
	owner := c.Query("owner")

	blockChain := &chain.BlockChain{}
	fmt.Println("=-----------set value-----------", owner, location, id)
	rst, err := blockChain.WriteHouse(id, location, owner)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

// 权限测试
func TestCertificate(c *gin.Context)  {
	key := c.Query("key")
	value := c.Query("value")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.TestCertificate(key, value)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func ReadFile(c *gin.Context)  {
	keyPath := "service/crypto-config/peerOrganizations/member1.example.com/users/Admin@member1.example.com/tls/client.key"

	//certPath := "service/crypto-config/peerOrganizations/member1.example.com/users/Admin@member1.example.com/tls/client.crt"
	//client1 := "/home/felix/fabric/fabric-pa/Admin@org1.example.com/tls/client.crt"
	path := c.Query("path")
	var creatorByte []byte
	var err error
	if path != "" {
		creatorByte, err = ioutil.ReadFile(path)
	} else {
		creatorByte, err = ioutil.ReadFile(keyPath)
	}


	//creatorByte,_:= stub.GetCreator()
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
	SendResponse(c, nil, cert)
}