package chain

// 写入house info
func (blockChain *BlockChain) WriteHouse(id, location, owner string)  (BlockChainResponse, error) {
	var args [][]byte
	args = append(args, []byte(id))
	args = append(args, []byte(location))
	args = append(args, []byte(owner))
	return blockChain.Request(ChainCodeID, Write, args)
}