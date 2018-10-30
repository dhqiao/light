package chain

import (
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"fmt"
	"log"
	"errors"
)

func Validate() error {
	//读取配置文件(config.yaml)中的组织(member1.example.com)的用户(Admin)
	mspClient, err := mspclient.New(FabSDK.Context(),
		mspclient.WithOrg("member1.example.com"))
	if err != nil {
		log.Fatalf("create msp client fail: %s\n", err.Error())
		return errors.New("create msp client fail: " + err.Error())
	}

	adminIdentity, err := mspClient.GetSigningIdentity("Admin")
	if err != nil {
		log.Fatalf("get admin identify fail: %s\n", err.Error())
		return errors.New("get admin identify fail: " + err.Error())
	} else {
		fmt.Println("AdminIdentify is found:")
		fmt.Println(adminIdentity)
	}
	return nil
}