package packages

import (
	"fmt"
	"sync"

	"github.com/arvindkhoisnam/go_probo_engine/models"
	"github.com/arvindkhoisnam/go_probo_engine/redisManager"
)
type Engine struct {
	Markets []ORDERBOOK
	InrBalance models.INR_BALANCE
	StockBalance models.STOCK_BALANCE
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
	Event   	 string `json:"event"`
	RedisChannel string `json:"redisChannel"`
	Payload 	 Data
}


func InitEngine() *Engine {
	once.Do(func ()  {
		engineInstance = &Engine{
				Markets: make([]ORDERBOOK, 0),
				InrBalance: models.INR_BALANCE{
					User: make(map[string]models.UserBalance),
				},
				StockBalance: models.STOCK_BALANCE{
					User: make(map[string]models.StockSymbol),
				},
		}
	})
	return engineInstance
}

func (e *Engine)StartEngine(incoming Incoming){
	switch incoming.Event {
	case "createMarket" :
		e.CreateMarket(incoming.Payload.Stock,incoming.RedisChannel)
	case "createUser" :
		e.CreateUser(incoming.Payload.UserId,incoming.RedisChannel)
	case "onramp" :
		e.OnrampINR(incoming.Payload.UserId,incoming.RedisChannel,incoming.Payload.Amount)
	case "getInrBal" :
		e.GetInrBal(incoming.Payload.UserId,incoming.RedisChannel)
	case "getStockBal" :
		e.GetStockBal(incoming.Payload.UserId,incoming.RedisChannel)
	case "allMarkets" :
		e.AllMarkets(incoming.RedisChannel)
	case "mint" :
		e.Mint(incoming.Payload.UserId,incoming.Payload.Stock,incoming.RedisChannel,incoming.Payload.Quantity,incoming.Payload.Price)
	case "placeOrder":
		e.PlaceOrder(incoming.Payload.UserId,incoming.Payload.Stock,incoming.Payload.StockType,incoming.Payload.OrderType,incoming.RedisChannel,incoming.Payload.Quantity,incoming.Payload.Price)
	case "sellOrderbook":
		e.GetSellOB(incoming.Payload.Stock,incoming.RedisChannel)
	case "getDepth":
		e.GetDepth(incoming.RedisChannel,incoming.Payload.Stock)	
	}

}

func (e *Engine)CreateMarket(stock,redisChan string){
	if _,exists := e.checkMarket(stock);exists{
		outgoing := &redisManager.Outgoing{
			StatusCode: 400,
			Message:  fmt.Sprintf("%s market already exists.",stock),
		}
		redisManager.PubToRedis(redisChan,outgoing)
		return
	}
	orderbook := CreateOrderbook(stock)
	
	e.Markets = append(e.Markets, *orderbook)
	outgoing := &redisManager.Outgoing{
		StatusCode: 200,
		Message:  fmt.Sprintf("%s market successfully created.",stock),
	}
	redisManager.PubToRedis(redisChan,outgoing)
}

func (e *Engine)CreateUser(userId, redisChan string){
	if exists := e.checkUser(userId);exists{
		outgoing := &redisManager.Outgoing{
			StatusCode: 400,
			Message:  fmt.Sprintf("%s already exists.",userId),
		}
		redisManager.PubToRedis(redisChan,outgoing)
		return
	}
		 e.InrBalance.User[userId] = models.UserBalance{
			Balance: 0,
			Locked: 0,
		 }
		 e.StockBalance.User[userId] = models.StockSymbol{
			Symbol: map[string]models.StockType{
				
			},
		 }
		 outgoing := &redisManager.Outgoing{
			StatusCode: 200,
			Message:  fmt.Sprintf("%s successfully created.",userId),
		}
		redisManager.PubToRedis(redisChan,outgoing)
}

func (e *Engine)OnrampINR(userId,redisChan string, amount int){
	if exists := e.checkUser(userId); !exists{
		outgoing := &redisManager.Outgoing{
			StatusCode: 400,
			Message:  "No user found. Please create user.",
		}
		redisManager.PubToRedis(redisChan,outgoing)
		return
	}
	user := e.InrBalance.User[userId]
	user.Balance += amount
	e.InrBalance.User[userId] = user
	outgoing := &redisManager.Outgoing{
		StatusCode: 200,
		Message:  fmt.Sprintf("Onramped %d to %s successfully .",amount,userId),
	}
	redisManager.PubToRedis(redisChan,outgoing)
}

func (e *Engine)GetInrBal(userId,redisChan string) {
	if exists := e.checkUser(userId);!exists{
		outgoing := &redisManager.Outgoing{
			StatusCode: 400,
			Message:  "No user found. Please create user.",
		}
		redisManager.PubToRedis(redisChan,outgoing)
		return
	}
	userBal := e.InrBalance.User[userId]
	// temp := models.UserBalance{
	// 	Balance: userBal.Balance,
	// 	Locked: userBal.Locked,
	// }
	outgoing := &redisManager.Outgoing{
		StatusCode: 200,
		Payload: redisManager.Data{
			INRBalance: models.UserBalance(userBal),
		},
	}
	redisManager.PubToRedis(redisChan, outgoing)
}

func (e *Engine)GetStockBal(userId,redisChan string) {
	if exists := e.checkUser(userId);!exists{
		outgoing := &redisManager.Outgoing{
			StatusCode: 400,
			Message:  "No user found. Please create user.",
		}
		redisManager.PubToRedis(redisChan,outgoing)
		return
	}
	userStocks := e.StockBalance.User[userId]

	// Convert StockType map to models.StockType map
	convertedSymbol := make(map[string]models.StockType)
	for stockKey, stockVal := range userStocks.Symbol {
		convertedType := make(map[models.StockEnum]models.Quantity)

		// Convert StockEnum & Quantity
		for enumKey, qtyVal := range stockVal.Type {
			convertedType[models.StockEnum(enumKey)] = models.Quantity{
				Available: qtyVal.Available,
				Locked:    qtyVal.Locked,
			}
		}

		convertedSymbol[stockKey] = models.StockType{
			Type: convertedType,
		}
	}

	// Correctly populate temp with the converted data
	temp := models.StockSymbol{
		Symbol: convertedSymbol,
	}

	outgoing := &redisManager.Outgoing{
		StatusCode: 200,
		Payload: redisManager.Data{
			StockBalance:  temp,
		},
	}
	redisManager.PubToRedis(redisChan, outgoing)
}

func (e *Engine)AllMarkets(redisChan string){
	var allMarkets []models.Market
	for _, val := range e.Markets{
		m := models.Market{
			StockSymbol: val.StockSymbol,
			CurrYesPrice: val.CurrYesPrice,
			CurrNoPrice: val.CurrNoPrice,
		}
		allMarkets = append(allMarkets,m)
	}
		outgoing := &redisManager.Outgoing{
		StatusCode: 200,
		Payload: redisManager.Data{
			Markets: allMarkets,
		},
	}
	redisManager.PubToRedis(redisChan, outgoing)
}

func (e *Engine) Mint(userId, stock, redisChan string, qty, price int) {
	if _, exists := e.checkMarket(stock); !exists {
		fmt.Printf("Markets do not exist for %s stock.\n", stock)
		return
	}

	if exists := e.checkUser(userId); !exists {
		fmt.Println("No user found. Please create a user.")
		return
	}

	if sufficient := e.checkInrBal(userId, qty*price); !sufficient {
		fmt.Println("Insufficient balance")
		return
	}

	// Fetch user stocks
	userStocks, userExists := e.StockBalance.User[userId]
	if !userExists {
		// Initialize user stock balance if it does not exist
		userStocks = models.StockSymbol{Symbol: make(map[string]models.StockType)}
	}

	// Check if stock exists, if not initialize it
	stocks, stockExists := userStocks.Symbol[stock]
	if stockExists {
		// Update existing stock quantities
		stocks.Type[models.YesStock] =models.Quantity{
			Available: stocks.Type[models.YesStock].Available + qty,
			Locked:    stocks.Type[models.YesStock].Locked,
		}

		stocks.Type[models.NoStock] = models.Quantity{
			Available: stocks.Type[models.NoStock].Available + qty,
			Locked:    stocks.Type[models.NoStock].Locked,
		}
	} else {
		// Initialize new stock entry
		stocks = models.StockType{
			Type: map[models.StockEnum]models.Quantity{
				models.YesStock: {Available: qty, Locked: 0},
				models.NoStock:  {Available: qty, Locked: 0},
			},
		}
	}

	// Ensure userStocks.Symbol is initialized
	if userStocks.Symbol == nil {
		userStocks.Symbol = make(map[string]models.StockType)
	}

	// Assign updated/new stock back to user's stock balance
	userStocks.Symbol[stock] = stocks
	e.StockBalance.User[userId] = userStocks

	// Deduct balance
	userBal := e.InrBalance.User[userId]
	userBal.Balance -= qty * price
	e.InrBalance.User[userId] = userBal

	// Send success response
	outgoing := &redisManager.Outgoing{
		StatusCode: 200,
		Message:    fmt.Sprintf("%d yes and no stocks of %s have been minted to %s.", qty, stock, userId),
	}
	redisManager.PubToRedis(redisChan, outgoing)
}



func (e *Engine)PlaceOrder(userId,stock,stockType,orderType,redisChan string, quantity,price int){
	 market ,exists := e.checkMarket(stock)
	if !exists{
			outgoing := &redisManager.Outgoing{
			StatusCode: 200,
			Message: fmt.Sprintf("Market does not exist for %s stock.\n",stock),
		}
		redisManager.PubToRedis(redisChan, outgoing)
		return
	}

	if  orderType == "sell"{
		sufficientStocks := e.checkAndLockStock(userId,stock,stockType,quantity)
		if !sufficientStocks {
			outgoing := &redisManager.Outgoing{
				StatusCode: 200,
				Message: "Insufficient stocks",
			}
			redisManager.PubToRedis(redisChan, outgoing)
			return
		}
		market.PlaceSellOrder(redisChan,userId,stockType,quantity,price)
	} else {
		sufficeintBalance := e.checkAndLockInr(userId,quantity,price)
		if !sufficeintBalance {
			outgoing := &redisManager.Outgoing{
				StatusCode: 200,
				Message: "Insufficient inr balance.",
			}
			redisManager.PubToRedis(redisChan, outgoing)
			return
		}
		market.PlaceBuyOrder(redisChan,userId,stockType,quantity,price,e)
	}
}

func (e *Engine)GetSellOB(stock,redisChan string)  {
	market,exists := e.checkMarket(stock)
	if ! exists {
		fmt.Printf("Markets do not exist for %s stock.\n", stock)
		outgoing := &redisManager.Outgoing{
			StatusCode: 400,
			Message: fmt.Sprintf("Markets do not exist for %s stock.\n", stock),
		}
		redisManager.PubToRedis(redisChan, outgoing )
		return
	}

	yesOB := models.StrikePrice{
		Strike: map[int]models.Orders{},
	}
	noOB := models.StrikePrice{
		Strike: make(map[int]models.Orders),
	}

	for strike,order := range market.Sell.Type[models.Yes].Strike{
		yesOB.Strike[strike] = models.Orders{
			TotalOrders: order.TotalOrders,
			TimeStamp: map[int]models.User{},
		}
		for time, val := range order.TimeStamp {
			yesOB.Strike[strike].TimeStamp[time] = models.User{
				UserId: val.UserId,
				Quantity: val.Quantity,
				ReverseOrder: val.ReverseOrder,
			}
		}
	}
	for strike,order := range market.Sell.Type[models.No].Strike{
		noOB.Strike[strike] = models.Orders{
			TotalOrders: order.TotalOrders,
			TimeStamp: map[int]models.User{},
		}

		for time,val := range order.TimeStamp{
			noOB.Strike[strike].TimeStamp[time] = models.User{
				UserId: val.UserId,
				Quantity: val.Quantity,
				ReverseOrder: val.ReverseOrder,
			}
		}
	}
	outgoing := &redisManager.Outgoing{
		StatusCode: 200,
		Payload: redisManager.Data{
			SellOBYes: yesOB,
			SellOBNo: noOB,
		},
	}
	redisManager.PubToRedis(redisChan, outgoing )
}

func (e *Engine)GetDepth(redisChan,stock string){
	market,exists := e.checkMarket(stock)
	if !exists {
		outgoing := &redisManager.Outgoing{
		StatusCode: 400,
		Message: fmt.Sprintf("Market does not exist for %s stock.\n",stock),
		}
		redisManager.PubToRedis(redisChan,outgoing)
		return
	}
	depth := market.GetDepth()
	outgoing := &redisManager.Outgoing{
		StatusCode: 200,
		Payload: redisManager.Data{
			Depth: models.DepthType(depth),
		},
	}
	redisManager.PubToRedis(redisChan, outgoing )
}
// Helper function
func (e *Engine) checkUser(userId string)  bool {
	_, exists := e.InrBalance.User[userId]
	return exists
}

func (e *Engine)checkMarket(stock string)(*ORDERBOOK, bool){
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
	var st models.StockEnum
	if stockType == "yes" {
		st = models.YesStock
	}else if stockType == "no" {
		st = models.NoStock
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

func (e *Engine)checkAndLockInr(userId string, quantity,price int) bool {
	 exists := e.checkUser(userId)
	if !exists {
		fmt.Println("User does not exist")
		return false
	}

	availableBalance := e.InrBalance.User[userId]
	if availableBalance.Balance < quantity*price {
	fmt.Println("Insufficient balance")
	return false
	}
	availableBalance.Balance -= quantity*price
	availableBalance.Locked += quantity*price
	e.InrBalance.User[userId] = availableBalance
	return true
}

