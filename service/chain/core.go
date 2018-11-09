package chain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"log"
	"errors"
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
)

const (
	// block chain config
	ChainCodeID = "mycc"
	ChannelID = "mychannel" // channel id
	UserName = "Admin"
	AimOrg = "member1.example.com"

	// next words are chain code has supplied functions' name
	Query = "query"
	Write = "write"
	History = "history"
	Del = "del"
	GetByRange = "getByRange" // startKey, endKey
	QueryResult = "queryResult"

	QueryPrivateData = "getPrivateData"
	WritePrivateData = "writePrivateData"
	GetPrivateByRange = "getPrivateByRange" // startKey, endKey

	// business
	WriteHouse = "writeHouse"

	// test
	TestCertificate = "testCertificate"

)

type BlockChain struct {}

type BlockChainResponse struct {
	Data           interface{}            `json:"data"`
	TransactionID  fab.TransactionID      `json:"transactionId"`
	Response       channel.Response       `json:"response"`
}

var response channel.Response
var blockChainResponse BlockChainResponse

// 数据获取
func (blockChain *BlockChain) Get(key string) (BlockChainResponse, error) {
	var args [][]byte
	args = append(args, []byte(key))
	return blockChain.Request(ChainCodeID, Query, args)
}

// 数据设置
func (blockChain *BlockChain) Set(key, value string) (BlockChainResponse, error) {
	var args [][]byte
	args = append(args, []byte(key))
	args = append(args, []byte(value))

	return blockChain.Request(ChainCodeID, Write, args)
}

// 查看历史变更
func (blockChain *BlockChain) History(key string) (BlockChainResponse, error) {
	var args [][]byte
	args = append(args, []byte(key))
	return blockChain.Request(ChainCodeID, History, args)
}


// 删除
func (blockChain *BlockChain) Del(key string) (BlockChainResponse, error) {
	var args [][]byte
	args = append(args, []byte(key))
	return blockChain.Request(ChainCodeID, Del, args)
}

// 获取两个key之间的数据
func (blockChain *BlockChain) GetByRange(startKey, endKey string) (BlockChainResponse, error) {
	var args [][]byte
	args = append(args, []byte(startKey))
	args = append(args, []byte(endKey))
	return blockChain.Request(ChainCodeID, GetByRange, args)
}


// 获取结果
func (blockChain *BlockChain) QueryResult(key string) (BlockChainResponse, error) {
	var args [][]byte
	args = append(args, []byte(key))
	return blockChain.Request(ChainCodeID, QueryResult, args)
}



// 调用智能合约
func (blockChain *BlockChain) Request(chainCodeId, fcn string, args [][]byte) (BlockChainResponse, error){
	channelClient, err := getChannelClient()
	if err != nil {
		return blockChainResponse, errors.New("create channel client fail: " + err.Error())
	}
	request := channel.Request{
		ChaincodeID: chainCodeId,
		Fcn:         fcn,
		Args:        args,
	}

	response, err = call(channelClient, request)
	if err != nil {
		log.Fatal("query fail: ", err.Error())
		return blockChainResponse, nil
	}
	log.Println("response is: ", response.Proposal.TxnID, response.ChaincodeStatus, response.TransactionID, response.TxValidationCode)
	data := make(map[string]interface{})
	var payload interface{}
	json.Unmarshal(response.Payload, &payload)
	if payload == nil {
		payload = string(response.Payload[:])
	}
	data["payload"] = payload
	blockChainResponse.Data = data
	blockChainResponse.TransactionID = response.TransactionID
	blockChainResponse.Response = response


	var data1 interface{}
	json.Unmarshal(response.Proposal.Payload, &data1)

	fmt.Println("................2..........", string(response.Proposal.Payload[:]), data1)

	var data2 interface{}
	json.Unmarshal(response.Responses[0].Payload, &data2)
	fmt.Println("................3..........", string(response.Responses[0].Payload[:]), data2)

	var data3 interface{}
	json.Unmarshal(response.Responses[0].Response.Payload, &data3)
	fmt.Println("................4..........", string(response.Responses[0].Response.Payload[:]), data3)


	return blockChainResponse, nil
}

// call function of blockchain
func call(channelClient *channel.Client,request channel.Request) (channel.Response, error) {
	switch request.Fcn {
	case Query:
		return channelClient.Query(request)
	case Write:
		return channelClient.Execute(request)
	case History:
		return channelClient.Query(request)
	case GetByRange:
		return channelClient.Query(request)
	case Del:
		return channelClient.Execute(request)
	case QueryResult:
		return channelClient.Query(request)
	case QueryPrivateData:
		return channelClient.Query(request)
	case WritePrivateData:
		return channelClient.Execute(request)
	case GetPrivateByRange:
		return channelClient.Query(request)
	case WriteHouse:
		return channelClient.Execute(request)
	case TestCertificate:
		return channelClient.Query(request)
	}
	return channel.Response{}, errors.New("call unknown func: " + request.Fcn)
}

func getChannelClient() (*channel.Client, error)  {
	//调用合约
	channelProvider := FabSDK.ChannelContext(ChannelID,
		fabsdk.WithUser(UserName),
		fabsdk.WithOrg(AimOrg))
	return channel.New(channelProvider)
}

func getLedgerClient() (*ledger.Client, error) {
	channelProvider := FabSDK.ChannelContext(ChannelID,
		fabsdk.WithUser(UserName),
		fabsdk.WithOrg(AimOrg))

	return ledger.New(channelProvider)
}