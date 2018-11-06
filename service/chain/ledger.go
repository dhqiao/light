package chain

import (
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/utils"
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

func GetLedgerClient() (*ledger.Client, error) {
	channelProvider := FabSDK.ChannelContext(ChannelID,
		fabsdk.WithUser(UserName),
		fabsdk.WithOrg(AimOrg))

	return  ledger.New(channelProvider)
}