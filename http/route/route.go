package router

import (
	"net/http"

	"light/http/controllers/health"
	"light/http/controllers/user"
	"light/http/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"light/http/controllers/chain"
	"light/http/controllers/kafka"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	//g.Use(middleware.NoCache)
	//g.Use(middleware.Options)
	//g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// pprof router
	pprof.Register(g)

	// api for authentication functionalities
	g.POST("/login", user.Login)

	// The user handlers, requiring authentication
	u := g.Group("/v1/user")
	u.Use(middleware.AuthMiddleware())
	{
		u.POST("", user.Create)
		u.DELETE("/:id", user.Delete)
		u.PUT("/:id", user.Update)
		u.GET("", user.List)
		u.GET("/:username", user.Get)
	}

	// The health check handlers
	svcd := g.Group("/health")
	{
		svcd.GET("/check", health.HealthCheck)
		svcd.GET("/disk", health.DiskCheck)
		svcd.GET("/cpu", health.CPUCheck)
		svcd.GET("/ram", health.RAMCheck)
	}

	blockChain := g.Group("/block/chain")
	{
		// base function
		blockChain.GET("get", chain.Get)
		blockChain.GET("set", chain.Set)
		blockChain.POST("setpost", chain.SetPost)
		blockChain.GET("history", chain.History)

		// block info
		blockChain.GET("queryBlockByTxID", chain.QueryBlockByTxID)
		blockChain.GET("queryBlock", chain.QueryBlock)
		blockChain.GET("queryTransaction", chain.QueryTransaction)
		blockChain.GET("queryBlockByHash", chain.QueryBlockByHash)
		blockChain.GET("queryConfig", chain.QueryConfig)
		blockChain.GET("queryInfo", chain.QueryInfo)
		blockChain.GET("getAllIdentities", chain.GetAllIdentities)
		blockChain.GET("del", chain.Del)
		blockChain.GET("getByRange", chain.GetByRange)

		blockChain.GET("queryResult", chain.QueryResult)
		blockChain.GET("test", chain.Test)

		// private test
		blockChain.GET("getPrivate", chain.GetPrivate)
		blockChain.GET("setPrivate", chain.SetPrivate)
		blockChain.GET("getPrivateByRange", chain.GetPrivateByRange)

		// business
		blockChain.GET("writeHouse", chain.WriteHouse)

		// test
		blockChain.GET("testCertificate", chain.TestCertificate)
		blockChain.GET("readFile", chain.ReadFile)



	}

	kafkaController := g.Group("kafka")
	{
		kafkaController.GET("send", kafka.Sync)
	}

	return g
}
