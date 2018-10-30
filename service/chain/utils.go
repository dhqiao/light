package chain

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"encoding/base64"
)

const (
	ErrorCode = 0
)

type BlockInfo struct {
	TransactionData      []TransactionData `json:"transction_list"`
	Number               uint64            `json:"number,omitempty"`
	PreviousHash         string            `json:"previous_hash,omitempty"`
	DataHash             string            `json:"data_hash,omitempty"`
} 

type TransactionData struct {
	Function string   `json:"function"`
	Key string        `json:"key"`
	Value string      `json:"value"`
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

var listRst []interface{}

func findDataFromJsonByte(jsonStr []byte, findKey string) ([]interface{}, error) {
	jsonMap := make(map[string]interface{})
	json.Unmarshal(jsonStr, &jsonMap)
	findDataFromJson(jsonMap, findKey)
	return listRst, nil
}

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
func decodeChaincodeSpec(chaincodeSpec interface{}) (string, string, string ) {
	if chaincodeSpec == nil {
		return "", "", ""
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
							return decodeString(argsList[0].(string)), decodeString(argsList[1].(string)), ""
						}
						if len(argsList) == 3 {
							return decodeString(argsList[0].(string)), decodeString(argsList[1].(string)), decodeString(argsList[2].(string))
						}
					}
				}
			}
		}
	}

	return "", "", ""
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