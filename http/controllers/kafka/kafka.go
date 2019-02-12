package kafka

import (
	. "light/http/controllers"
	"github.com/gin-gonic/gin"
	"light/service/kafka"
	"sync"
)

func Sync(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			SendResponse(c, nil, "出错")

		}
	}()
	kafka.SyncProducer()
	SendResponse(c, nil, "success")
}

func Async(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			SendResponse(c, nil, "出错")
		}
	}()
	kafka.SaramaProducer()
	SendResponse(c, nil, "success")
}

func Consumer(c *gin.Context)  {
	groupId := c.Query("id")
	defer func() {
		if err := recover(); err != nil {
			SendResponse(c, nil, "出错")
			return
		}
	}()
	var wg = &sync.WaitGroup{}
	wg.Add(2)
	go kafka.ClusterConsumer(wg, groupId)
	go kafka.ClusterConsumer(wg, groupId)
	SendResponse(c, nil, "success")
}
