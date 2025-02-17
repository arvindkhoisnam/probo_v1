package routeHandlers

import (
	"fmt"

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
	redisManager.PushToRedisAwait(data)
	c.JSON(200,gin.H{"message":fmt.Sprintf("%s successfully created.",symbol)})
}

func CreateUser(c *gin.Context){
	userId := c.Param("userId")
	data := &redisManager.Incoming{
		Event: "createUser",
		Payload: redisManager.Data{
			UserId: userId,
		},
	}
	redisManager.PushToRedisAwait(data)
	c.JSON(200,gin.H{"message":fmt.Sprintf("%s successfully created.",userId)})
}

func OnrampINR(c *gin.Context){
	type reqBody struct {
		UserId string `json:"userId"`
		Amount int `json:"amount"`
	}
	body := &reqBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
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

	redisManager.PushToRedisAwait(data)
	c.JSON(200,gin.H{"message":fmt.Sprintf("Successfully onramped %d to %s",body.Amount,body.UserId)})
}