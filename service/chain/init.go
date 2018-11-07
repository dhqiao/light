package chain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
)

var FabSDK *fabsdk.FabricSDK

func init()  {

	//读取配置文件，创建SDK
	configProvider := config.FromFile("config/conf/block_chain_config.yaml")
	var err error
	FabSDK, err = fabsdk.New(configProvider)
	if err != nil {
		log.Fatalf("create sdk fail: %s\n", err.Error())
	}

	// 配置验证
	err = Validate()
	if err != nil {
		log.Fatalf("create sdk fail: %s\n", err.Error())
	}

}