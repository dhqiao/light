package chain

import (
	. "light/http/controllers"
	"github.com/gin-gonic/gin"
	"light/service/chain"
	"fmt"
)

func Get(c *gin.Context) {
	key := c.Query("key")
	blockChain := &chain.BlockChain{}
	fmt.Println(">>>>>>>>>>>>>>>>>", key)
	rst, err := blockChain.Get(key)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func Set(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")
	blockChain := &chain.BlockChain{}
	fmt.Println("=-----------set value-----------", key, value)
	rst, err := blockChain.Set(key, value)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func History(c *gin.Context) {
	key := c.Query("key")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.History(key)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

// 删除键
func Del(c *gin.Context) {
	key := c.Query("key")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.Del(key)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

// 获取两个key之间的数据
func GetByRange(c *gin.Context) {
	startKey := c.Query("startKey")
	endKey := c.Query("endKey")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.GetByRange(startKey, endKey)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}


func QueryResult(c *gin.Context) {
	key := c.Query("key")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryResult(key)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func QueryBlock(c *gin.Context) {
	blockNum := c.Query("blockNum")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryBlock(blockNum)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func QueryTransaction(c *gin.Context) {
	txID := c.Query("txid")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryTransaction(txID)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}


func QueryBlockByTxID(c *gin.Context) {
	txID := c.Query("txid")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>", txID)
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryBlockByTxID(txID)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func Test(c *gin.Context) {
	txID := c.Query("txid")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>", txID)
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.Test(txID)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

func QueryBlockByHash(c *gin.Context)  {
	hash := c.Query("hash")
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryBlockByHash(hash)
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

// 没有使用
func GetBlockchainInfo(c *gin.Context)  {
	blockChain := &chain.BlockChain{}
	rst := blockChain.GetBlockchainInfo()
	SendResponse(c, nil, rst)
}

// 获取配置信息 只能到order节点和锚点peer
func QueryConfig(c *gin.Context)  {
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryConfig()
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

// 获取区块总数量，最后一个block，倒数第二个block hash， 和背书节点
func QueryInfo(c *gin.Context)  {
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.QueryInfo()
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

// 获取所有的成员，我们没有成员管理
func GetAllIdentities(c *gin.Context)  {
	blockChain := &chain.BlockChain{}
	rst, err := blockChain.GetAllIdentities()
	if err != nil {
		SendResponse(c, nil, err.Error())
	}
	SendResponse(c, nil, rst)
}

