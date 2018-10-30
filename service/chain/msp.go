package chain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (blockChain *BlockChain) GetAllIdentities() (BlockChainResponse, error)  {
	//读取配置文件(config.yaml)中的组织(member1.example.com)的用户(Admin)
	mspClient, err := msp.New(FabSDK.Context(),
		msp.WithOrg("member1.example.com"))
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	identities, err := mspClient.GetAllIdentities()
	return BlockChainResponse{identities, "", channel.Response{}}, err
}