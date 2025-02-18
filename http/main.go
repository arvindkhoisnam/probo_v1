package main

import (
	"log"

	"github.com/arvindkhoisnam/go_probo_http/redisManager"
	"github.com/arvindkhoisnam/go_probo_http/routeHandlers"
	"github.com/gin-gonic/gin"
)




func routes(app *gin.Engine){
	api := app.Group("/api/v1")
	{
		api.GET("/health",routeHandlers.HealthHandler)
		api.POST("/symbol/create/:symbol",routeHandlers.CreateMarket)
		api.POST("/user/create/:userId",routeHandlers.CreateUser)
		api.POST("/onramp/inr",routeHandlers.OnrampINR)
		api.POST("/trade/mint",routeHandlers.Mint)
		api.GET("/balances/inr/:userId",routeHandlers.GetInrBal)
		api.GET("/balances/stock/:userId",routeHandlers.GetStockBal)
		api.GET("allmarkets",routeHandlers.AllMarkets)
	}
}


func main() {
	redisManager.InitRedis()
	app := gin.Default()
	routes(app)
	log.Println("LISTENING ON PORT 3000")
	log.Fatal(app.Run(":3000"))
}