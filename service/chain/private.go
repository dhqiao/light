package chain

/**
 * test private data for fabric 1.2 new feature
 */
// 数据获取
func (blockChain *BlockChain) GetPrivate(collection, key string) (BlockChainResponse, error) {
	var args [][]byte
	args = append(args, []byte(collection))
	args = append(args, []byte(key))
	return blockChain.Request(ChainCodeID, QueryPrivateData, args)
}

// 数据设置
func (blockChain *BlockChain) SetPrivate(collection, key, value string) (BlockChainResponse, error) {
	var args [][]byte
	args = append(args, []byte(collection))
	args = append(args, []byte(key))
	args = append(args, []byte(value))

	return blockChain.Request(ChainCodeID, WritePrivateData, args)
}

// 获取两个key之间的数据
func (blockChain *BlockChain) GetPrivateByRange(collection, startKey, endKey string) (BlockChainResponse, error) {
	var args [][]byte
	args = append(args, []byte(collection))
	args = append(args, []byte(startKey))
	args = append(args, []byte(endKey))
	return blockChain.Request(ChainCodeID, GetPrivateByRange, args)
}
