package chain

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"encoding/base64"
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"

)

const (
	ErrorCode = 0
)

// 临时变量，因为使用了一个递归函数，返回值有点问题
var listRst []interface{}

// define block info
type BlockInfo struct {
	TransactionData      []TransactionData `json:"transction_list"`
	Number               uint64            `json:"number,omitempty"`
	PreviousHash         string            `json:"previous_hash,omitempty"`
	DataHash             string            `json:"data_hash,omitempty"`
	OriginData           map[string]interface{}            `json:"origin_data,omitempty"`
} 

// define transaction struct for return data
type TransactionData struct {
	Function string   `json:"function"`
	Key string        `json:"key"`
	Value []string      `json:"value"`
	Message string    `json:"message"`
	Status float64    `json:"status"`
	Timestamp string  `json:"timestamp"`
	TXId string       `json:"tx_id"`
	ChannelId string  `json:"channel_id"`
}

/**

// 返回数据
proposal_response_payload: {
	extension: {
		chaincode_id: {
			name: "mycc",
			path: "",
			version: "1.0"
		},
		response: {
			message: "",
			status: 200
		},
		results: "EhQKBGxzY2MSDAoKCgRteWNjEgIIAxIdCgRteWNjEhUaEwoIbmFtZWRkZGQaB2FhYWFhYWE="
	},
	proposal_hash: "IVidVKfFsYKIT9F8sjUuBHVqiCQLMEne5W29Xe8xhIc="
}

// 请求数据
chaincode_spec: {
	chaincode_id: {
		name: "mycc",
		path: "",
		version: ""
	},
	input: {
		args: [
			"d3JpdGU=",
			"bmFtZWRkZGQ=",
			"YWFhYWFhYQ=="
		]
	},
	timeout: 0,
	type: "GOLANG"
}

// header
channel_header: {
	channel_id: "mychannel",
	epoch: "0",
	extension: "EgYSBG15Y2M=",
	timestamp: "2018-10-18T09:50:24.094223902Z",
	tx_id: "947e69dcbacf80285a5fe96dcbf2b644749078ff538e9d50429f18c7a6b2f683",
	type: 3,
	version: 0
},
*/

// 解析区块数据
func decodeBlockJson(jsonStr []byte) ([]TransactionData, error) {
	jsonMap := make(map[string]interface{})
	json.Unmarshal(jsonStr, &jsonMap)
	listRst = nil
	findDataFromJson(jsonMap, "proposal_response_payload")
	var chaincodeProposalPayloads []interface{}  // find key proposal_response_payload
	chaincodeProposalPayloads = listRst
	listRst = nil
	findDataFromJson(jsonMap, "chaincode_spec")
	var chaincodeSpecs []interface{}
	chaincodeSpecs = listRst
	listRst = nil
	findDataFromJson(jsonMap, "channel_header")
	var channelHeaders []interface{}
	channelHeaders = listRst

	if len(chaincodeProposalPayloads) != len(chaincodeSpecs) && len(chaincodeSpecs) != len(channelHeaders) {
		return nil, errors.New("数据解析失败")
	}
	var index int
	var transactionDataList []TransactionData
	for index = 0; index < len(chaincodeProposalPayloads); index++ {
		message, status := decodeChaincodeProposalPayload(chaincodeProposalPayloads[index])
		function, key, value := decodeChaincodeSpec(chaincodeSpecs[index])
		channelId, timestamp, txId := decodeChannelHeader(channelHeaders[index])
		transactionDataList = append(transactionDataList, TransactionData{
			function,
			key,
			value,
			message,
			status,
			timestamp,
			txId,
			channelId,
		})
	}

	return transactionDataList, nil
}

// channel_id, timestamp, tx_id
func decodeChannelHeader(channelHeader interface{}) (string, string, string) {
	if channelHeader == nil {
		return "", "", ""
	}
	channelHeaderMap, ok := channelHeader.(map[string]interface{})
	if ok {
		return channelHeaderMap["channel_id"].(string), channelHeaderMap["timestamp"].(string), channelHeaderMap["tx_id"].(string)
	}
	return "", "", ""
}

// function key value
func decodeChaincodeSpec(chaincodeSpec interface{}) (string, []string, string ) {
	var values []string

	if chaincodeSpec == nil {
		return "", values, ""
	}
	chaincodeSpecMap, ok := chaincodeSpec.(map[string]interface{})
	if ok {
		input, ok := chaincodeSpecMap["input"]
		if ok {
			inputMap, ok := input.(map[string]interface{})
			if ok {
				args, ok := inputMap["args"]
				if ok {
					argsList, ok := args.([]interface{})

					if ok {
						// del
						if len(argsList) == 2 {
							values = append(values, decodeString(argsList[1].(string)))
							return decodeString(argsList[0].(string)), values, ""
						}
						/*
						if len(argsList) == 3 {
							values = append(values, decodeString(argsList[1].(string)))
							return decodeString(argsList[0].(string)), values, decodeString(argsList[2].(string))
						}*/
						// if len(argsList) > 3 不是设置的key value的形式，而是多值输入
						if len(argsList) >= 3 {
							len := len(argsList)
							index := 1
							for ; index < len-1; index ++ {
								values = append(values, decodeString(argsList[index].(string)))
							}
							return decodeString(argsList[0].(string)), values, decodeString(argsList[2].(string))
						}
					}
				}
			}
		}
	}

	return "", values, ""
}

func decodeString(str string) string {
	stringR, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(stringR[:])
}

func encodeToString(byte []byte) string {
	return base64.StdEncoding.EncodeToString(byte)
}

// message, status
func decodeChaincodeProposalPayload(chaincodeProposalPayload interface{}) (string, float64) {
	if chaincodeProposalPayload == nil {
		return "", ErrorCode
	}
	payload, ok := chaincodeProposalPayload.(map[string]interface{})
	if ok {
		extension, ok := payload["extension"]

		if ok {
			extensionMap, ok := extension.(map[string]interface{})

			if ok {
				response, ok := extensionMap["response"]

				if ok {
					responseMap, ok := response.(map[string]interface{})

					if ok {
						return responseMap["message"].(string), (responseMap["status"]).(float64)
					}
				}
			}
		}
	}
	return "", ErrorCode
}

// 递归查找
func findDataFromJson(jsonData interface{}, findKey string) error {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("errrrrrrrrrrrrrrrrr", p)
		}
	}()

	if jsonData == nil {
		return nil
	}

	switch jsonData.(type) {
	case map[string]interface{}:
		jsonMap := changeDataToMap(jsonData)
		if jsonMap != nil {
			for key, value := range jsonMap {
				fmt.Println("************************* key: ", key)
				if key == findKey {
					fmt.Println(">>>>>>>>>>>>>>>>>>>>>>",key, value)
					listRst = append(listRst, value)
					return  nil
				} else {
					findDataFromJson(value, findKey)
				}
			}
		}
	case []interface{}:
		jsonSlice := changeDataToSlice(jsonData)
		if jsonSlice != nil {
			var index int
			for index = 0; index < len(jsonSlice); index++ {
				fmt.Println("************************* index: ", index)
				findDataFromJson(jsonSlice[index], findKey)
			}
		}
	default:
		return nil
	}
	return nil
}

func changeDataToMap(jsonData interface{}) map[string]interface{} {
	jsonMap := jsonData.(map[string]interface{})
	defer func() map[string]interface{} {
		return nil
	}()
	return jsonMap
}

func changeDataToSlice(jsonData interface{}) []interface{} {
	jsonSlice := jsonData.([]interface{})
	defer func() []interface{} {
		return nil
	}()
	return jsonSlice
}

// 解析区块数据
func ParseBlockInfo(block *common.Block) (BlockInfo, error ){
	var blockInfo BlockInfo
	b, err := proto.Marshal(block)
	if err != nil {
		return blockInfo, err
	}
	httpResp, err := http.Post("http://127.0.0.1:7059/protolator/decode/common.Block", "application/octet-stream", bytes.NewReader(b))
	resBody, err := ioutil.ReadAll(httpResp.Body)

	httpResp.Body.Close()
	transactionList, err := decodeBlockJson(resBody)
	if err != nil{
		fmt.Println("<<<<<<<<<<<<<<<<<<<<decode json error")
	}
	blockInfo.PreviousHash = encodeToString(block.Header.PreviousHash)
	blockInfo.DataHash = encodeToString(block.Header.DataHash)
	blockInfo.TransactionData = transactionList
	blockInfo.Number = block.Header.Number

	// 不需要的时候注释掉
	jsonMap := make(map[string]interface{})
	json.Unmarshal(resBody, &jsonMap)
	blockInfo.OriginData = jsonMap

	return blockInfo, nil
}