package chain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)


func(blockChain *BlockChain) GetBlockchainInfo() uint64 {/*
	err := common.InitConfig("core")
	if err != nil { // Handle errors reading the config file
		fmt.Printf("Cannot init configure, error=[%v]", err)
		os.Exit(1)
	}

	provider, err := kvledger.NewProvider()   // core/ledger/kvledger/kv_ledger_provider.go
	if err != nil {
		fmt.Printf("Cannot new provider, error=[%s]", err)
		os.Exit(1)
	}
	defer provider.Close()

	// Print channel list
	channels, err := provider.List()        // core/ledger/kvledger/kv_ledger_provider.go
	if err != nil {
		fmt.Printf("Cannot get channel list, error=[%v]\n", err)
		os.Exit(1)
	}
	fmt.Printf("channels=[%v]\n", channels)

	// Open a channel
	ledger, err := provider.Open(channels[0])   // core/ledger/kvledger/kv_ledger_provider.go
	if err != nil {
		fmt.Printf("Cannot open channel ledger, error=[%v]\n", err)
		os.Exit(1)
	}
	defer ledger.Close()

	// Return ledger as kvLedger is defined in core/ledger/kvledger/kv_ledger.go, following API:
	//  func (l *kvLedger) GetBlockchainInfo() (*common.BlockchainInfo, error)
	//  func (l *kvLedger) GetTransactionByID(txID string) (*peer.ProcessedTransaction, error)
	//  func (l *kvLedger) GetBlockByNumber(blockNumber uint64) (*common.Block, error)
	//  func (l *kvLedger) GetBlockByHash(blockHash []byte) (*common.Block, error)
	//  func (l *kvLedger) GetBlockByTxID(txID string) (*common.Block, error)
	//  func (l *kvLedger) GetBlocksIterator(startBlockNumber uint64) (commonledger.ResultsIterator, error)
	//  func (l *kvLedger) GetTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error)
	//  func (l *kvLedger) Close()

	// Get basic channel information
	chainInfo, err := ledger.GetBlockchainInfo() // (*common.BlockchainInfo, error)
	if err != nil {
		fmt.Printf("Cannot get block chain info, error=[%v]\n", err)
		os.Exit(1)
	}
	fmt.Printf("chainInfo: Height=[%d], CurrentBlockHash=[%s], PreviousBlockHash=[%s]\n",
		chainInfo.GetHeight(),
		base64.StdEncoding.EncodeToString(chainInfo.CurrentBlockHash),
		base64.StdEncoding.EncodeToString(chainInfo.PreviousBlockHash))

	// Retrieve blocks based on block number
	for i := uint64(0); i < chainInfo.GetHeight(); i++ {
		block, err := ledger.GetBlockByNumber(i) // (blockNumber uint64) (*common.Block, error)
		if err != nil {
			fmt.Printf("Cannot get block for %d, error=[%v]\n", i, err)
			os.Exit(1)
		}
		printBlock("Get", block)
	}


	// Retrieve blocks based on iterator
	itr, err := ledger.GetBlocksIterator(0) // (ResultsIterator, error)
	if err != nil {
		fmt.Printf("Cannot get iterator, error=[%v]\n", err)
		os.Exit(1)
	}
	defer itr.Close()

	queryResult, err := itr.Next()    // commonledger.QueryResult
	for i := uint64(0); err == nil; i++ {
		block := queryResult.(*cb.Block)
		printBlock("Iterator", block)
		if i >= chainInfo.GetHeight() - 1 {
			break
		}
		queryResult, err = itr.Next()    // commonledger.QueryResult
	}
*/
	return 12

}


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



/*
func (blockChain *BlockChain) QueryBlockFile() {
	fileName = "/opt/app/fabric/peer/data/ledgersData/chains/chains/mychannel/blockfile_000000"

	var err error
	if file, err = os.OpenFile(fileName, os.O_RDONLY, 0600); err != nil {
		fmt.Printf("ERROR: Cannot Open file: [%s], error=[%v]\n", fileName, err)
		return
	}
	defer file.Close()


	if fileInfo, err := file.Stat(); err != nil {
		fmt.Printf("ERROR: Cannot Stat file: [%s], error=[%v]\n", fileName, err)
		return
	} else {
		fileOffset = 0
		fileSize   = fileInfo.Size()
		fileReader = bufio.NewReader(file)
	}

	// Loop each block
	for {
		if blockBytes, err := nextBlockBytes(); err != nil {
			fmt.Printf("ERROR: Cannot read block file: [%s], error=[%v]\n", fileName, err)
			break
		} else if blockBytes == nil {
			// End of file
			break
		} else {
			if block, err := deserializeBlock(blockBytes); err != nil {
				fmt.Printf("ERROR: Cannot deserialize block from file: [%s], error=[%v]\n", fileName, err)
				break
			} else {
				handleBlock(block)
			}
		}
	}
}*/