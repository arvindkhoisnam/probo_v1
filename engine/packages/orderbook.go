package packages

import (
	"fmt"
	"sort"

	"github.com/arvindkhoisnam/go_probo_engine/models"
	"github.com/arvindkhoisnam/go_probo_engine/redisManager"
	"golang.org/x/exp/maps"

	"math"
	"time"
)

type ORDERBOOK struct{
	StockSymbol  string
	CurrYesPrice int
	CurrNoPrice  int
	Buy          models.BuyOrderbook
	Sell         models.SellOrderbook
}

func CreateOrderbook(stock string) *ORDERBOOK {
	ob := ORDERBOOK {
		StockSymbol: stock,
		CurrYesPrice: 0,
		CurrNoPrice: 0,
		Buy: models.BuyOrderbook{
			Type: map[models.StockTypeEnum]models.StrikePrice{},
		},
		Sell: models.SellOrderbook{
			Type: map[models.StockTypeEnum]models.StrikePrice{
				models.Yes : {
					Strike: map[int]models.Orders{},
				},
				models.No : {
					Strike: map[int]models.Orders{},
				},
			},
		},
	}
	return &ob
}

func (ob *ORDERBOOK)PlaceSellOrder(redisChan,userId,stockType string, quantity,price int){
	var st models.StockTypeEnum
	if stockType == "yes" {
		st = models.Yes
	}else if stockType == "no" {
		st = models.No
	}
	if exists := ob.CheckBuyer(st,price); !exists {
		if _,typeExists := ob.Sell.Type[st]; !typeExists {
			ob.Sell.Type[st] = models.StrikePrice{
				Strike: map[int]models.Orders{},
			}
		}
		if _,strikeExists := ob.Sell.Type[st].Strike[price]; !strikeExists{
			ob.Sell.Type[st].Strike[price] = models.Orders{
				TotalOrders: 0,
				TimeStamp: map[int]models.User{},
			}
		}

		strikePrice := ob.Sell.Type[st].Strike[price]
		strikePrice.TotalOrders += quantity
		strikePrice.TimeStamp[int(time.Now().Unix())] = models.User{
			UserId: userId,
			Quantity: quantity,
		}

		ob.Sell.Type[st].Strike[price] = strikePrice

		outgoing := &redisManager.Outgoing{
			StatusCode: 200,
			Message: fmt.Sprintf("Order for %d %s stocks of %s placed at %d.",quantity,stockType,ob.StockSymbol,price),
		}
		redisManager.PubToRedis(redisChan, outgoing )
	}
}


func (ob *ORDERBOOK)CheckBuyer(st models.StockTypeEnum,price int) bool {
	_,exists := ob.Buy.Type[st].Strike[price]
	return exists
}
func (ob *ORDERBOOK)CheckSeller(st models.StockTypeEnum,price,quantity int) bool {
	_,availabe := ob.Sell.Type[st]
	if !availabe {
		return availabe
	}
	fillable := ob.CalcFillableBuyQty(price,quantity,st)
	return fillable > 0
}

func (ob *ORDERBOOK)PlaceBuyOrder(redisChan,userId,stockType string, quantity,price int, e *Engine){
	var st models.StockTypeEnum
	var mst models.StockEnum
	if stockType == "yes" {
		st = models.Yes
		mst = models.YesStock
	}else if stockType == "no" {
		st = models.No
		mst = models.NoStock
	}

	isAvailable := ob.CheckSeller(st,price,quantity)
	if !isAvailable {
		fmt.Println("No orders can be matched at the moment.")
		ob.ReverseOrder(redisChan,userId,price,quantity,st)
		outgoing := &redisManager.Outgoing{
			StatusCode: 200,
			Message: fmt.Sprintf("Reverse order for %d %s stocks of %s placed at %d.",quantity,stockType,ob.StockSymbol,price),
		}
		redisManager.PubToRedis(redisChan, outgoing )
		return
	}

	ob.matchOrder(redisChan,userId,price,quantity,st,mst,e)
}

func (ob *ORDERBOOK)matchOrder(redisChan, userId string ,price,qty int, st models.StockTypeEnum,mst models.StockEnum,e *Engine){
	var filledQty int
	pendingQty := qty
	fillable := ob.CalcFillableBuyQty(price,qty,st)
	strikes := ob.Sell.Type[st]
	sortedStrikeKeys := maps.Keys(strikes.Strike)
	sort.Ints(sortedStrikeKeys)
	for _,currStrike := range sortedStrikeKeys{
		order := strikes.Strike[currStrike]
		if  filledQty < fillable{
			filled := math.Min(float64(pendingQty),float64(order.TotalOrders))
			order.TotalOrders -= int(filled)
			timestamps := &order.TimeStamp
			ob.manageStocksAndInr(userId,int(filled),currStrike,mst,e,timestamps)
			pendingQty -= int(filled)
			filledQty += int(filled)
			ob.Sell.Type[st].Strike[currStrike] = order
			if order.TotalOrders == 0 {
				delete(ob.Sell.Type[st].Strike,currStrike)
			}
			ob.updateLTP(&st,currStrike)
			ob.pubToWS()
		}
	}
	outgoing := &redisManager.Outgoing{
		StatusCode: 200,
		Message: "Order placed successfully.",
		Payload: redisManager.Data{
			FilledOrders: filledQty,
			PendingOrders: pendingQty,
		},
	}
	redisManager.PubToRedis(redisChan, outgoing )
}

func (ob *ORDERBOOK)updateLTP(st *models.StockTypeEnum,LTP int){
	if *st == models.Yes {
		ob.CurrYesPrice = LTP
		ob.CurrNoPrice = 10 - LTP
		}else{
			ob.CurrNoPrice =  LTP
			ob.CurrYesPrice = 10 - LTP
		}
}
func (ob *ORDERBOOK)pubToWS(){
	depth := ob.GetDepth()
	ticker := ob.GetTicker()
	outgoing2 := &redisManager.Outgoing{
		Payload: redisManager.Data{
			Depth: models.DepthType(depth),
			Ticker: models.TickerType(ticker),
		},
	}
	redisManager.PubToRedis(ob.StockSymbol, outgoing2 )
}
func(ob *ORDERBOOK)manageStocksAndInr(buyer string, toFill,strike int,mst models.StockEnum, e *Engine, ts *map[int]models.User){
	sortedTimestamps := maps.Keys(*ts)
	sort.Ints(sortedTimestamps)
	tempToFill := toFill
	for _,time := range sortedTimestamps {
		if tempToFill > 0{
			orders := (*ts)[time]
			seller := orders.UserId
			toDeduct := math.Min(float64(orders.Quantity),float64(tempToFill))
			orders.Quantity -= int(toDeduct)
			tempToFill -= int(toDeduct)

			(*ts)[time] = orders
			if orders.Quantity == 0 {
				delete(*ts,time)
			}
			if orders.ReverseOrder {
				ob.fillReverseOrder(seller,buyer,strike,int(toDeduct),mst,e)
			} else {
				sellerStocks := e.StockBalance.User[seller].Symbol[ob.StockSymbol].Type[mst]
				sellerStocks.Locked -= int(toDeduct)
				e.StockBalance.User[seller].Symbol[ob.StockSymbol].Type[mst] = sellerStocks
				
				_,existing := e.StockBalance.User[buyer].Symbol[ob.StockSymbol]
				if !existing {
					e.StockBalance.User[buyer].Symbol[ob.StockSymbol] = models.StockType{
						Type: map[models.StockEnum]models.Quantity{
							models.YesStock: {
								Available: 0,
								Locked: 0,
							},
							models.NoStock :{
								Available: 0,
								Locked: 0,
							},
						},
					}
				}
				buyerStocks := e.StockBalance.User[buyer].Symbol[ob.StockSymbol].Type[mst]
				buyerStocks.Available += int(toDeduct)
				e.StockBalance.User[buyer].Symbol[ob.StockSymbol].Type[mst] = buyerStocks
				ob.manageINR(buyer,seller,int(toDeduct),strike,e)
			}
		}
	}
} 
func (ob *ORDERBOOK)manageINR(buyer,seller string, qty,price int, e *Engine){
	buyerBal := e.InrBalance.User[buyer]
	sellerBal := e.InrBalance.User[seller]

	buyerBal.Locked -= qty * price
	sellerBal.Balance += qty * price

	e.InrBalance.User[buyer] = buyerBal
	e.InrBalance.User[seller] = sellerBal
}

func (ob *ORDERBOOK)CalcFillableBuyQty( price,qty int, st models.StockTypeEnum) int {
	strikes := ob.Sell.Type[st]
	var temp int
	remainingQty := qty
	for strike, order := range strikes.Strike {
		if strike <= price && temp <= qty{
			if remainingQty >= order.TotalOrders {
				temp += order.TotalOrders
				remainingQty -= order.TotalOrders
			} else {
				temp += remainingQty
			}
		}
	}
	return temp
}

func (ob *ORDERBOOK)ReverseOrder(redisChan,userId string, price,quantity int, st models.StockTypeEnum){
	if st == models.Yes {
		st = models.No
	}else {
		st = models.Yes
	}

	newStrikePrice := 10 - price
	if _,typeExists := ob.Sell.Type[st]; !typeExists {
		ob.Sell.Type[st] =models.StrikePrice{
			Strike: map[int]models.Orders{},
		}
	}
	if _,strikeExists := ob.Sell.Type[st].Strike[newStrikePrice]; !strikeExists{
		ob.Sell.Type[st].Strike[newStrikePrice] = models.Orders{
			TotalOrders: 0,
			TimeStamp: map[int]models.User{},
		}
	}
	strikePrice := ob.Sell.Type[st].Strike[newStrikePrice]
	strikePrice.TotalOrders += quantity
	strikePrice.TimeStamp[int(time.Now().Unix())] = models.User{
		ReverseOrder: true,
		UserId: userId,
		Quantity: quantity,
	}

	ob.Sell.Type[st].Strike[newStrikePrice] = strikePrice
}

func (ob *ORDERBOOK)fillReverseOrder(seller,buyer string ,price,qty int, mst models.StockEnum,e *Engine){
	// buyer
	_,existing := e.StockBalance.User[buyer].Symbol[ob.StockSymbol]
	if !existing {
		e.StockBalance.User[buyer].Symbol[ob.StockSymbol] = models.StockType{
			Type: map[models.StockEnum]models.Quantity{
				models.YesStock: {
					Available: 0,
					Locked: 0,
				},
				models.NoStock :{
					Available: 0,
					Locked: 0,
				},
			},
		}
	}
	buyerStocks := e.StockBalance.User[buyer].Symbol[ob.StockSymbol].Type[mst]
	buyerStocks.Available += int(qty)
	e.StockBalance.User[buyer].Symbol[ob.StockSymbol].Type[mst] = buyerStocks

	buyerBal := e.InrBalance.User[buyer]
	buyerBal.Locked -= qty * price
	e.InrBalance.User[buyer] = buyerBal

	// seller
	if mst == models.YesStock {
		mst = models.NoStock
	}else {
		mst = models.YesStock
	}

	_,existingStocks := e.StockBalance.User[seller].Symbol[ob.StockSymbol]
	if !existingStocks {
		e.StockBalance.User[seller].Symbol[ob.StockSymbol] = models.StockType{
			Type: map[models.StockEnum]models.Quantity{
				models.YesStock: {
					Available: 0,
					Locked: 0,
				},
				models.NoStock :{
					Available: 0,
					Locked: 0,
				},
			},
		}
	}
	sellerStocks := e.StockBalance.User[seller].Symbol[ob.StockSymbol].Type[mst]
	sellerStocks.Available += int(qty)
	e.StockBalance.User[seller].Symbol[ob.StockSymbol].Type[mst] = sellerStocks

	sellerBal := e.InrBalance.User[seller]
	sellerBal.Locked -= qty * (10-price)
	e.InrBalance.User[seller] = sellerBal
}

func(ob *ORDERBOOK)GetTicker()models.TickerType{
	ticker := models.TickerType {
		Ticker: ob.StockSymbol,
		YesPrice: ob.CurrYesPrice,
		NoPrice: ob.CurrNoPrice,
	}

 	return ticker
}

func (ob *ORDERBOOK)GetDepth()models.DepthType{
	depth := models.DepthType{
		YesMarket: map[int]int{},
		NoMarket: map[int]int{},
	}

	for strike,orders := range ob.Sell.Type[models.Yes].Strike{
		depth.YesMarket[strike] = orders.TotalOrders
	}
	for strike,orders := range ob.Sell.Type[models.No].Strike{
		depth.NoMarket[strike] = orders.TotalOrders
	}
	return depth
}