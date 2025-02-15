package packages

import (
	"fmt"
	"sync"
)

type userBalance struct {
	Balance int
	Locked  int
}
type INR_BALANCE struct {
	User  map[string]userBalance
}

type quantity struct {
	Available int
	Locked int
}
type StockType int
const (
	YesStock StockType = iota 
	NoStock
)
type stockType struct {
	Type map[StockType] quantity
}
type stockSymbol struct {
	Symbol map[string] stockType
}

type STOCK_BALANCE struct {
	User map[string] stockSymbol
}

type Engine struct {
	Orderbook []ORDERBOOK
	InrBalance INR_BALANCE
	StockBalance STOCK_BALANCE
}

var (
	once sync.Once
	engineInstance *Engine
)

func InitEngine() *Engine {
	once.Do(func ()  {
		engineInstance = &Engine{
			Orderbook: make([]ORDERBOOK, 0),
			InrBalance: INR_BALANCE{
				User: make(map[string]userBalance),
			},
			StockBalance: STOCK_BALANCE{
				User: make(map[string]stockSymbol),
			},
		}
	})
	return engineInstance
}

func (e *Engine)CreateMint(stock string){
	for _, val := range e.Orderbook {
		if val.StockSymbol == stock{
			fmt.Println("ALREADY EXISTS")
			return
		}
	}
	
	orderbook := CreateOrderbook(stock)
	e.Orderbook = append(e.Orderbook, *orderbook)
}

func (e *Engine)CreateUser(userId string){
	_,alreadyExists := e.InrBalance.User[userId]
	if (alreadyExists) {
		fmt.Println("User already created.")
		return
	} 
		 e.InrBalance.User[userId] = userBalance{
			Balance: 0,
			Locked: 0,
		 }
		 e.StockBalance.User[userId] = stockSymbol{}
		 fmt.Println(e.InrBalance)
		fmt.Println(e.StockBalance) 
	
}

func (e *Engine)OnrampINR(userId string, amount int){
	user,userExists := e.InrBalance.User[userId]
	if (!userExists){
		fmt.Println("No user found. Please create user.")
		return
	}
	user.Balance += amount
	fmt.Println(user)
}