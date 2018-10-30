package router

import (
	"net/http"

	"light/http/controllers/health"
	"light/http/controllers/user"
	"light/http/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"light/http/controllers/chain"
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
		blockChain.GET("get", chain.Get)
		blockChain.GET("set", chain.Set)
		blockChain.GET("history", chain.History)
		blockChain.GET("queryBlockByTxID", chain.QueryBlockByTxID)
		blockChain.GET("queryBlock", chain.QueryBlock)
		blockChain.GET("queryTransaction", chain.QueryTransaction)
		blockChain.GET("queryBlockByHash", chain.QueryBlockByHash)
		blockChain.GET("getBlockchainInfo", chain.GetBlockchainInfo)
		blockChain.GET("queryConfig", chain.QueryConfig)
		blockChain.GET("queryInfo", chain.QueryInfo)
		blockChain.GET("getAllIdentities", chain.GetAllIdentities)
		blockChain.GET("del", chain.Del)
		blockChain.GET("getByRange", chain.GetByRange)

		blockChain.GET("queryResult", chain.QueryResult)
		blockChain.GET("test", chain.Test)


	}

	return g
}
