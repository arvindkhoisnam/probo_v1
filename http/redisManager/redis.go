package redisManager

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)


var (	
	client *redis.Client
	once sync.Once
)
type Data struct {
	UserId    string `json:"userId,omitempty"`
	Stock     string `json:"stock,omitempty"`
	StockType string `json:"stockType,omitempty"`
	OrderType string `json:"orderType,omitempty"`
	Quantity  int    `json:"quantity,omitempty"`
	Price     int    `json:"price,omitempty"`
	Amount    int    `json:"amount,omitempty"`
}
type Incoming struct {
	Event   string
	Payload Data
}
func InitRedis() *redis.Client{
	once.Do(func(){
		opt, err := redis.ParseURL("redis://localhost:6379")
		if err != nil {
			panic(err)
		}
		client = redis.NewClient(opt)
	})
	return client
}

func PushToRedisAwait(data *Incoming){
	stringified,err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	client.LPush(ctx,"engine",stringified)
}