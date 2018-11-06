package chain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"encoding/base64"
	"strconv"
	"github.com/golang/protobuf/proto"

)

func (blockChain *BlockChain)QueryBlockByTxID(txID string) (BlockChainResponse, error) {

	client, err := GetLedgerClient()
	// if err...
	block, err := client.QueryBlockByTxID(fab.TransactionID(txID))
	if err != nil {
		fmt.Printf("Ledger_QueryBlockByTxID return error: %v", err)
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	//header, _ := json.Marshal(block.Header)

	blockInfo, err := ParseBlockInfo(block)
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	return BlockChainResponse{blockInfo, "", channel.Response{}}, err
}


func (blockChain *BlockChain)QueryBlockByHash(blockHash string) (BlockChainResponse, error) {
	client, err := getLedgerClient()

	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}

	blockHashByte, err := base64.StdEncoding.DecodeString(blockHash)
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	block, err := client.QueryBlockByHash(blockHashByte)
	// if err
	blockInfo, err := ParseBlockInfo(block)

	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	return BlockChainResponse{blockInfo, "", channel.Response{}}, nil
}

func (blockChain *BlockChain)QueryBlock(blockNum string) (BlockChainResponse, error) {
	ledgerClient, err := getLedgerClient()
	blockNo, err := strconv.ParseUint(blockNum, 0, 64)
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	block, err := ledgerClient.QueryBlock(blockNo)
	// if err
	blockInfo, err := ParseBlockInfo(block)
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	return BlockChainResponse{blockInfo, "", channel.Response{}}, nil
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
	blockInfo, err := ParseBlockInfo(block)

	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	return BlockChainResponse{blockInfo, "", channel.Response{}}, err
}

// 写死了
func GetLedgerClient() (*ledger.Client, error) {
	channelProvider := FabSDK.ChannelContext(ChannelID,
		fabsdk.WithUser(UserName),
		fabsdk.WithOrg(AimOrg))

	return  ledger.New(channelProvider)
}