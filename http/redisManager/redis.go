package redisManager

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)


var (	
	client *redis.Client
	subscriber *redis.Client
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
	Event   	 string
	RedisChannel string
	Payload 	 Data
}
func InitRedis() *redis.Client{
	once.Do(func(){
		opt, err := redis.ParseURL("redis://localhost:6379")
		if err != nil {
			panic(err)
		}
		client = redis.NewClient(opt)
		subscriber = redis.NewClient(opt)
	})
	return client
}
type Market struct{
	StockSymbol string
	CurrYesPrice int
	CurrNoPrice int
}
type UserBalance struct {
	Balance int
	Locked  int
}
type Quantity struct {
	Available int
	Locked int
}
type StockEnum int
const (
	YesStock StockEnum = iota
	NoStock
)
type StockType struct {
	Type map[StockEnum] Quantity
}
type StockSymbol struct {
	Symbol map[string] StockType
}
type OrderType int
const (
	Yes OrderType = iota
	No
)
type Orders struct {
	TotalOrders int
	Order map[string] int
}

type StrikePrice struct {
	Strike map[int] Orders
}
type Data2 struct {
	UserId string `json:"userId,omitempty"`
	BalanceInr int	  `json:"balanceInr,omitempty"`
	LockedInr int `json:"lockedInr,omitempty"`
	Markets []Market
	INRBalance UserBalance
	StockBalance StockSymbol
	SellOBYes StrikePrice
	SellOBNo StrikePrice

}
type Outgoing struct {
	StatusCode int
	Message string
	Payload Data2
}
func PushToRedisAwait(data *Incoming) *Outgoing{
	ctx := context.Background()
	randId := genRandId()
	data.RedisChannel = randId
	stringified,err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	sub := subscriber.Subscribe(ctx, randId)
	defer sub.Close() // Ensure cleanup
	client.LPush(ctx,"engine",stringified)
	// Wait for a response (blocking)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
	}
	
	outgoing := &Outgoing{}
	json.Unmarshal([]byte(msg.Payload),outgoing)
	fmt.Println(outgoing)
	return outgoing
}

func genRandId() string {
	b := make([]byte, 8) // 8 bytes
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	randomStr := base64.URLEncoding.EncodeToString(b)
	return randomStr
}