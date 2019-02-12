package kafka

import (
	. "light/http/controllers"
	"github.com/gin-gonic/gin"
	"light/service/kafka"
)

func Sync(c *gin.Context)  {
	defer func() {
		if err := recover(); err != nil {
			SendResponse(c, nil, "出错")

		}
	}()
	kafka.SyncProducer()
	SendResponse(c, nil, "success")
}
