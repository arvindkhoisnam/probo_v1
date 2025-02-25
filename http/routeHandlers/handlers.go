package routeHandlers

import (
	"github.com/arvindkhoisnam/go_probo_http/redisManager"
	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context){
	c.JSON(200, gin.H{"message":"Healthy Server."})
}

func CreateMarket(c *gin.Context){
	symbol := c.Param("symbol")
	data:= &redisManager.Incoming{
		Event: "createMarket",
		Payload: redisManager.Data{
			Stock: symbol,
		},
	}
	outgoing := redisManager.PushToRedisAwait(data)
	c.JSON(outgoing.StatusCode,gin.H{"message":outgoing.Message})
}

func CreateUser(c *gin.Context){
	userId := c.Param("userId")
	data := &redisManager.Incoming{
		Event: "createUser",
		Payload: redisManager.Data{
			UserId: userId,
		},
	}
	outgoing := redisManager.PushToRedisAwait(data)
	c.JSON(outgoing.StatusCode,gin.H{"message": outgoing.Message})
}

func OnrampINR(c *gin.Context){
	type reqBody struct {
		UserId string `json:"userId"`
		Amount int 	  `json:"amount"`
	}
	body := &reqBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
		}
	data := &redisManager.Incoming{
		Event: "onramp",
		Payload: redisManager.Data{
			UserId: body.UserId,
			Amount: body.Amount,
		},
	}

	outgoing := redisManager.PushToRedisAwait(data)
	c.JSON(outgoing.StatusCode,gin.H{"message": outgoing.Message})
}

func GetInrBal(c *gin.Context){
	userId := c.Param("userId")
	data := &redisManager.Incoming{
		Event: "getInrBal",
		Payload: redisManager.Data{
			UserId: userId,
		},
	}

	outgoing := redisManager.PushToRedisAwait(data)
	c.JSON(outgoing.StatusCode,gin.H{"data": outgoing.Payload.INRBalance})
}

func Mint(c *gin.Context){
	type reqBody struct {
		UserId 	 string `json:"userId"`
		Stock    string `json:"stock"`
		Quantity int  	`json:"quantity"`
		Price 	 int 	`json:"price"`
	}

	body := &reqBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
		}
	data := &redisManager.Incoming{
		Event: "mint",
		Payload: redisManager.Data{
			UserId: body.UserId,
			Stock: body.Stock,
			Quantity: body.Quantity,
			Price: body.Price,
		},
	}

	outgoing := redisManager.PushToRedisAwait(data)
	c.JSON(outgoing.StatusCode,gin.H{"message":outgoing.Message})
}

func GetStockBal(c *gin.Context){
	userId := c.Param("userId")

	data := &redisManager.Incoming{
		Event: "getStockBal",
		Payload: redisManager.Data{
			UserId: userId,
		},
	}

	outgoing := redisManager.PushToRedisAwait(data)
	c.JSON(outgoing.StatusCode,gin.H{"data":outgoing.Payload.StockBalance})
}

func AllMarkets(c *gin.Context){
	data:= &redisManager.Incoming{
		Event: "allMarkets",
	}
	outgoing := redisManager.PushToRedisAwait(data)
	c.JSON(outgoing.StatusCode,gin.H{"data":outgoing.Payload.Markets})
}

func SellOrder(c *gin.Context){
	type reqBody struct {
		UserId 	 string `json:"userId"`
		Stock    string `json:"stock"`
		Type	 string `json:"type"`
		Quantity int  	`json:"quantity"`
		Price 	 int 	`json:"price"`
	}

	body := &reqBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
		}
	data := &redisManager.Incoming{
		Event: "placeOrder",
		Payload: redisManager.Data{
			UserId: body.UserId,
			Stock: body.Stock,
			StockType: body.Type,
			OrderType: "sell",
			Quantity: body.Quantity,
			Price: body.Price,
		},
	}
	outgoing :=  redisManager.PushToRedisAwait(data)
	c.JSON(outgoing.StatusCode,gin.H{"data":outgoing.Message})
}

func BuyOrder(c *gin.Context){
	type reqBody struct {
		UserId 	 string `json:"userId"`
		Stock    string `json:"stock"`
		Type	 string `json:"type"`
		Quantity int  	`json:"quantity"`
		Price 	 int 	`json:"price"`
	}
	body := &reqBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
		}
	data := &redisManager.Incoming{
		Event: "placeOrder",
		Payload: redisManager.Data{
			UserId: body.UserId,
			Stock: body.Stock,
			StockType: body.Type,
			OrderType: "buy",
			Quantity: body.Quantity,
			Price: body.Price,
		},
	}
	outgoing :=  redisManager.PushToRedisAwait(data)
	c.JSON(outgoing.StatusCode,gin.H{"message":outgoing.Message,"filled":outgoing.Payload.FilledOrders,"pending":outgoing.Payload.PendingOrders})
}

func GetSellOB(c *gin.Context){
	stock := c.Param("stock")
	data := &redisManager.Incoming{
		Event: "sellOrderbook",
		Payload: redisManager.Data{
			Stock: stock,
		},
	}
	outgoing := redisManager.PushToRedisAwait(data)
	c.JSON(outgoing.StatusCode,gin.H{"yes": outgoing.Payload.SellOBYes,"no":outgoing.Payload.SellOBNo})
}