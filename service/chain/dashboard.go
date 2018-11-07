package chain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (blockChain *BlockChain)QueryConfig() (BlockChainResponse, error) {
	ledgerClient, err := getLedgerClient()
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	channelCfg, err := ledgerClient.QueryConfig()
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	config := make(map[string]interface{})
	config["BlockNumber"] = channelCfg.BlockNumber()
	config["AnchorPeers"] = channelCfg.AnchorPeers()
	config["ID"] = channelCfg.ID()
	config["MSPs"] = channelCfg.MSPs()
	config["Orderers"] = channelCfg.Orderers()

	return BlockChainResponse{config, "", channel.Response{}}, nil
}

func (blockChain *BlockChain)QueryInfo() (BlockChainResponse, error) {
	ledgerClient, err := getLedgerClient()
	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	channelCfg, err := ledgerClient.QueryInfo()

	if err != nil {
		return BlockChainResponse{"", "", channel.Response{}}, err
	}
	config := make(map[string]interface{})
	config["channelCfg"] = channelCfg
	return BlockChainResponse{config, "", channel.Response{}}, nil
}