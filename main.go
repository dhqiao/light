package main

import (
	"light/config"
	"github.com/spf13/viper"
	"light/http/route"
	"light/http/middleware"
	"github.com/gin-gonic/gin"
	"light/models"
	"light/library/envcheck"
)

func main()  {
        // update test branch oooooo
	
	// init configuration
	if err := config.Init(""); err != nil {
		panic(err)
	}
	// init db
	models.DB.Init()
	defer models.DB.Close()

	// Set gin mode.
	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()
	// Routes.
	router.Load(
		// Cores
		g,
		// Middlwares.
		middleware.Logging(),
		middleware.RequestId(),
		middleware.SetGlobalData(),
	)
	envcheck.ServerCheck(g)
}
