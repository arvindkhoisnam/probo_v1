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
	Markets []ORDERBOOK
	InrBalance INR_BALANCE
	StockBalance STOCK_BALANCE
}

var (
	once sync.Once
	engineInstance *Engine
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
	Event   string `json:"event"`
	Payload Data
}


func InitEngine() *Engine {
	once.Do(func ()  {
		engineInstance = &Engine{
			Markets: make([]ORDERBOOK, 0),
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

func (e *Engine)StartEngine(incoming Incoming){
	switch incoming.Event {
	case "createMarket" :
		e.CreateMarket(incoming.Payload.Stock)
	case "createUser" :
		e.CreateUser(incoming.Payload.UserId)
	case "onramp" :
		e.OnrampINR(incoming.Payload.UserId,incoming.Payload.Amount)
	case "getInrBal" :
		e.GetInrBal(incoming.Payload.UserId)
	case "getStockBal" :
		e.GetStockBal(incoming.Payload.UserId)
	case "allMarkets" :
		e.AllMarkets()
	case "mint" :
		e.Mint(incoming.Payload.UserId,incoming.Payload.Stock,incoming.Payload.Quantity,incoming.Payload.Price)
	case "placeOrder":
		e.PlaceOrder(incoming.Payload.UserId,incoming.Payload.Stock,incoming.Payload.StockType,incoming.Payload.OrderType,incoming.Payload.Quantity,incoming.Payload.Price)
	}

}

func (e *Engine)CreateMarket(stock string){
	if _,exists := e.checkMarket(stock);exists{
		fmt.Println("ALREADY EXISTS")
			return
	}
	orderbook := CreateOrderbook(stock)
	e.Markets = append(e.Markets, *orderbook)
	fmt.Println(e.Markets)
}

func (e *Engine)CreateUser(userId string){
	if exists := e.checkUser(userId);exists{
		fmt.Println("User already created.")
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
	if exists := e.checkUser(userId); !exists{
		fmt.Println("No user found. Please create user.")
		return
	}
	user := e.InrBalance.User[userId]
	user.Balance += amount
	e.InrBalance.User[userId] = user
	fmt.Println(user)
}

func (e *Engine)GetInrBal(userId string) {
	if exists := e.checkUser(userId);!exists{
		fmt.Println("User does not exist")
	}

	userBal := e.InrBalance.User[userId]
	fmt.Println(userBal)
}

func (e *Engine)GetStockBal(userId string) {
	if exists := e.checkUser(userId);!exists{
		fmt.Println("User does not exist")
	}

	userStocks := e.StockBalance.User[userId]
	fmt.Println(userStocks)
}
type market struct{
	stockSymbol string
	currYesPrice int
	currNoPrice int
}
func (e *Engine)AllMarkets()[]market{
	var allMarkets []market
	for _, val := range e.Markets{
		m := market{
			stockSymbol: val.StockSymbol,
			currYesPrice: val.CurrYesPrice,
			currNoPrice: val.CurrNoPrice,
		}
		allMarkets = append(allMarkets,m)
	}
	return allMarkets
}

func (e *Engine)Mint(userId, stock string, qty,price int){
	if _,exists := e.checkMarket(stock); !exists{
		fmt.Printf("Markets does not exist for %s stock.\n",stock)
		return
	}

	if exists := e.checkUser(userId); !exists{
		fmt.Println("No user found. Please create user.")
		return
	}

	if sufficient := e.checkInrBal(userId,qty*price); !sufficient {
		fmt.Println("Insufficient balance")
		return
	}

	userStocks := e.StockBalance.User[userId]
	userStocks.Symbol = map[string]stockType{
		stock: {
			Type: map[StockType]quantity{
				YesStock: {Available: qty, Locked: 0},
				NoStock:  {Available: qty, Locked: 0},
			},
		},
	}
	e.StockBalance.User[userId] = userStocks

	userBal := e.InrBalance.User[userId]
	userBal.Balance -= qty*price
	e.InrBalance.User[userId] = userBal
}


func (e *Engine)PlaceOrder(userId,stock,stockType,orderType string, quantity,price int){
	 market ,exists := e.checkMarket(stock)
	if !exists{
		fmt.Printf("Markets does not exist for %s stock.\n",stock)
		return
	}
	sufficientStocks := e.checkAndLockStock(userId,stock,stockType,quantity)

	if !sufficientStocks {
		fmt.Println("failed at checkAndLock")
		return
	}

	market.PlaceSellOrder(userId,stockType,quantity,price)

}

// Helper function
func (e *Engine) checkUser(userId string)  bool {
	_, exists := e.InrBalance.User[userId]
	return exists
}

func (e *Engine)checkMarket(stock string)(*ORDERBOOK,bool){
	for _, val := range e.Markets {
		if val.StockSymbol == stock{
			return &val,true
		}
	}
	return nil, false
}

func (e *Engine)checkInrBal(userId string, request int)bool{
	user := e.InrBalance.User[userId]
	if user.Balance >= request {
		return user.Balance >= request
	}
	return false
}

func (e *Engine)checkAndLockStock(userId,stock,stockType string, quantity int ) bool {
	var st StockType
	if stockType == "yes" {
		st = YesStock
	}else if stockType == "no" {
		st = NoStock
	} else {
		fmt.Println("stock type not available")
		return false
	}
	 userStocks,exists := e.StockBalance.User[userId].Symbol[stock].Type[st]; 
	if !exists {
		fmt.Printf("You do not have %s stocks",stock)
		return false
	}
	if userStocks.Available < quantity {
		fmt.Printf("You do not have enough %s %s stocks \n",stockType,stock)
		return false
	}
	userStocks.Available -= quantity
	userStocks.Locked += quantity

	e.StockBalance.User[userId].Symbol[stock].Type[st] = userStocks
	return true
}