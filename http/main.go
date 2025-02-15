// package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/arvindkhoisnam/go_probo_http/packages"
// 	"github.com/gin-gonic/gin"
// )

// // Health Check Endpoint
// func healthEndpoint(c *gin.Context) {
// 	c.JSON(200, gin.H{"message":"Healthy server."})
// }

// // Create User Endpoint
// func createUser(instance *packages.BALANCE_INR) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userId := c.Param("userId") // Get userId from route param
// 		data, err := instance.Onramp(userId)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"data": data})
// 	}
// }

// // Get All Users' Balances
// func balancesInr(address *packages.BALANCE_INR) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		users,length := address.GetAllUsers()
// 		c.JSON(http.StatusOK, gin.H{"length":length,"data": users})
// 	}
// }

// // Setup Routes
// func Routes(app *gin.Engine, inrInstance *packages.BALANCE_INR) {
// 	api := app.Group("/api/v1")
// 	{
// 		api.GET("/health", healthEndpoint)
// 		api.GET("/balances/inr", balancesInr(inrInstance))
// 		api.POST("/user/create/:userId", createUser(inrInstance))
// 	}
// }

// func main() {
// 	app := gin.Default() // Gin comes with built-in Logger & Recovery middleware
// 	INR_BALANCE_INSTANCE := packages.Init()
// 	Routes(app, INR_BALANCE_INSTANCE)

// 	log.Println("Server is running on port 3000")
// 	log.Fatal(app.Run(":3000"))
// }

package main

import (
	"fmt"
	"log"

	"github.com/arvindkhoisnam/go_probo_http/packages"
	"github.com/gin-gonic/gin"
)

func healthHandler(c *gin.Context){
	c.JSON(200, gin.H{"message":"Healthy Server."})
}

func handleCreateUser(instance *packages.BALANCE_INR) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		data,err := instance.Onramp(userId)
		if err != nil {
			ctx.JSON(400,gin.H{"message":"user already exists"})
		}
		ctx.JSON(200,gin.H{"message":fmt.Sprintf("%s successfully created.",userId),"data": *data})
	}
}

func handleBalancesInr(instance *packages.BALANCE_INR) gin.HandlerFunc {
	return func(ctx *gin.Context){
		data,length := instance.GetAllUsers()
		ctx.JSON(200,gin.H{"data":data,"length": length})
	}
}
func routes(app *gin.Engine,instance *packages.BALANCE_INR){
	api := app.Group("/api/v1")
	{
		api.GET("/health",healthHandler)
		api.POST("/user/create/:userId",handleCreateUser(instance))	
		api.GET("/balances/inr",handleBalancesInr(instance))	
	}
}

func main() {
	app := gin.Default()

	INR_BALANANCE_INSTANCE := packages.Init()

	routes(app,INR_BALANANCE_INSTANCE)
	log.Println("LISTENING ON PORT 3000")
	log.Fatal(app.Run(":3000"))
}