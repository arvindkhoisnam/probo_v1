package redisManager

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/arvindkhoisnam/go_probo_engine/models"
	"github.com/redis/go-redis/v9"
)


var (
redisClient *redis.Client
Once sync.Once
)

type Data struct {
	UserId 	      string 			`json:"userId,omitempty"`
	Markets       []models.Market  `json:"markets"`
	INRBalance 	  models.UserBalance
	StockBalance  models.StockSymbol
	SellOBYes     models.StrikePrice
	SellOBNo      models.StrikePrice
	FilledOrders  int
	PendingOrders int
}

type Outgoing struct {
	StatusCode int
	Message    string
	Payload    Data
}

func InitRedis() *redis.Client {
	Once.Do(func ()  {	
		opt, err := redis.ParseURL("redis://localhost:6379")
		if err != nil {
			panic(err)
		}
		redisClient = redis.NewClient(opt)
	})
	return redisClient
}

func PubToRedis(channel string, data *Outgoing){
	fmt.Println(data)
	stringified,err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	redisClient.Publish(ctx,channel,stringified)
}