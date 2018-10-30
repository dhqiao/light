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
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	"strconv"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/utils"
	"encoding/base64"
)

const (
	// block chain config
	ChainCodeID = "mycc"
	ChannelID = "mychannel" // channel id
	UserName = "Admin"
	AimOrg = "member2.example.com"

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

func (blockChain *BlockChain)QueryBlockByTxID(txID string) (BlockChainResponse, error) {
	channelProvider := FabSDK.ChannelContext(ChannelID,
		fabsdk.WithUser(UserName),
		fabsdk.WithOrg(AimOrg))

	client, err := ledger.New(channelProvider)

	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	// block, err := client.QueryBlockByHash(txID)
	var id fab.TransactionID
	id = fab.TransactionID(txID)
	fmt.Println("...........1...........", id, txID)
	block, err := client.QueryBlockByTxID(id)
	fmt.Println("...........block...........", block.Header)

	fmt.Println("...........block-len...........", len(block.Data.Data))
	fmt.Println("...........block-byte2Json...........", byte2Json(block.Data.Data[0]))
	envelope, err:= utils.ExtractEnvelope(block, 0)
	fmt.Println("..................envelope............", envelope)

	fmt.Println("...........block3...........", string(block.Metadata.Metadata[0]))

	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	return BlockChainResponse{block, "", channel.Response{}}, nil
}


func (blockChain *BlockChain)QueryBlockByHash(blockHash string) (BlockChainResponse, error) {
	channelProvider := FabSDK.ChannelContext(ChannelID,
		fabsdk.WithUser(UserName),
		fabsdk.WithOrg(AimOrg))

	client, err := ledger.New(channelProvider)

	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}

	blockHashByte, err := base64.StdEncoding.DecodeString(blockHash)
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	block, err := client.QueryBlockByHash(blockHashByte)
	//block, err := client.QueryBlockByTxID(id)
	fmt.Println("...........block...........", block, blockHash)

	fmt.Println("...........block1...........", string(block.Data.Data[0]), len(block.Data.Data))
	fmt.Println("...........block-len...........", len(block.Data.Data))
	fmt.Println("...........block-byte2Json...........", byte2Json(block.Data.Data[0]))



	fmt.Println("...........block3...........", string(block.Metadata.Metadata[0]))

	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	return BlockChainResponse{block, "", channel.Response{}}, nil
}

func (blockChain *BlockChain)QueryBlock(blockNum string) (BlockChainResponse, error) {
	ledgerClient, err := getLedgerClient()
	blockNo, err := strconv.ParseUint(blockNum, 0, 64)
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	block, err := ledgerClient.QueryBlock(blockNo)
	dataa := block.Metadata.Metadata
	fmt.Println("...........dataa[0]...........", dataa[0])
	fmt.Println("...........dataa[0].string...........", string(dataa[0]))
	fmt.Println("...........dataa[0][0].string...........", string(dataa[0][0]))



	rst := make(map[string]interface{})
	//txsFltr := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
//fmt.Println("txsFltr", txsFltr)
	for _, r := range block.Data.Data {
		tx, _ := getTxPayload1(r)
		if tx != nil {
			chdr, err := utils.UnmarshalChannelHeader(tx.Header.ChannelHeader)
			rst["chdr"] = chdr
			if err != nil {
				fmt.Print("Error extracting channel header\n")
			}
			//getChainCodeEvents1(r)

		}
	}












	rst["block"] = block


	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	return BlockChainResponse{rst, "", channel.Response{}}, nil
}

func (blockChain *BlockChain)QueryTransaction(txID string) (BlockChainResponse, error) {

	ledgerClient, err := getLedgerClient()
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	processedTransaction, err := ledgerClient.QueryTransaction(fab.TransactionID(txID))
	fmt.Println("..............QueryTransaction.................", processedTransaction)
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	mapRst := make(map[string]interface{})

	envelope, err := proto.Marshal(processedTransaction.TransactionEnvelope)
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	httpResp, err := http.Post("http://127.0.0.1:7059/protolator/decode/common.Envelope", "application/octet-stream", bytes.NewReader(envelope))

	resBody, err := ioutil.ReadAll(httpResp.Body)

	transactionList, err := decodeBlockJson(resBody)

	mapRst["transactionList"] = transactionList

	return BlockChainResponse{mapRst, "", channel.Response{}}, nil

}



func (blockChain *BlockChain) Test (txID string) (BlockChainResponse, error){

	channelProvider := FabSDK.ChannelContext(ChannelID,
		fabsdk.WithUser(UserName),
		fabsdk.WithOrg(AimOrg))

	client, err := ledger.New(channelProvider)

	block, err := client.QueryBlockByTxID(fab.TransactionID(txID))
	if err != nil {
		fmt.Printf("Ledger_QueryBlockByTxID return error: %v", err)
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	//header, _ := json.Marshal(block.Header)

	b, err := proto.Marshal(block)
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	httpResp, err := http.Post("http://127.0.0.1:7059/protolator/decode/common.Block", "application/octet-stream", bytes.NewReader(b))

	resBody, err := ioutil.ReadAll(httpResp.Body)

	httpResp.Body.Close()



	transactionList, err := decodeBlockJson(resBody)
	var blockInfo BlockInfo
	if err != nil{
		fmt.Println("<<<<<<<<<<<<<<<<<<<<decode json error")
	}
	blockInfo.PreviousHash = encodeToString(block.Header.PreviousHash)
	blockInfo.DataHash = encodeToString(block.Header.DataHash)
	blockInfo.TransactionData = transactionList
	blockInfo.Number = block.Header.Number

	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	return BlockChainResponse{blockInfo, "", channel.Response{}}, err
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

func byte2Json(data []byte) interface{} {
	var rst interface{}
	json.Unmarshal(data, &rst)
	fmt.Println("-----------------------", rst)
	if rst == nil {
		//rst = string(response.Payload[:])
	}
	return rst
}

func printMap(params map[string]*common.ConfigGroup) {
	fmt.Println("..............print map begin............")
	for key, value := range params {
		fmt.Println("key:", key)
		fmt.Println("value:",value)
	}
	fmt.Println("..............print map end............")
}

func getTxPayload1(tdata []byte) (*common.Payload, error) {
	if tdata == nil {
		return nil, errors.New("Cannot extract payload from nil transaction")
	}

	if env, err := utils.GetEnvelopeFromBlock(tdata); err != nil {
		return nil, fmt.Errorf("Error getting tx from block(%s)", err)
	} else if env != nil {
		// get the payload from the envelope
		payload, err := utils.GetPayload(env)
		if err != nil {
			return nil, fmt.Errorf("Could not extract payload from envelope, err %s", err)
		}
		return payload, nil
	}
	return nil, nil
}


// getChainCodeEvents parses block events for chaincode events associated with individual transactions
func getChainCodeEvents1(tdata []byte) {
	env, err := utils.GetEnvelopeFromBlock(tdata)
	fmt.Println("-----------------env: ", env, err)

			if err != nil {

			} else if env != nil {
				// get the payload from the envelope
				payload, err := utils.GetPayload(env)
				if err != nil {

				}

				//chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
				if err != nil {

				}

				//if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {
					tx, err := utils.GetTransaction(payload.Data)
					if err != nil {

					}
					chaincodeActionPayload, err := utils.GetChaincodeActionPayload(tx.Actions[0].Payload)
					fmt.Println("-------------------chaincodeActionPayload: ", chaincodeActionPayload.Action)
				fmt.Println("-------------------chaincodeActionPayload: ", chaincodeActionPayload.Action.ProposalResponsePayload)

				if err != nil {

					}
					propRespPayload, err := utils.GetProposalResponsePayload(chaincodeActionPayload.Action.ProposalResponsePayload)
					if err != nil {

					}
					caPayload, err := utils.GetChaincodeAction(propRespPayload.Extension)
					fmt.Println("---------------------caPayload: ", caPayload)
					if err != nil {

					}
					ccEvent, err := utils.GetChaincodeEvents(caPayload.Events)
					fmt.Println("----------------------ccEvent: ", ccEvent)
					if ccEvent != nil {

					}

			}

}
