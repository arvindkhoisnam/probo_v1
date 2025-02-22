package packages

import (
	"fmt"
	"sort"

	"golang.org/x/exp/maps"

	"math"
	"time"
)

type StockTypeEnum int
const (
	Yes StockTypeEnum = iota 
	No
)

type User struct {
	UserId string
	Quantity int
}
type Orders struct {
	TotalOrders int
	TimeStamp map[int] User
}

type StrikePrice struct {
	Strike map[int] Orders
}

type BuyOrderbook struct {
	Type map[StockTypeEnum] StrikePrice
}
type SellOrderbook struct {
	Type map[StockTypeEnum] StrikePrice
}
type ORDERBOOK struct{
	StockSymbol  string
	CurrYesPrice int
	CurrNoPrice  int
	Buy          BuyOrderbook
	Sell         SellOrderbook
}

func CreateOrderbook(stock string) *ORDERBOOK {
	ob := ORDERBOOK {
		StockSymbol: stock,
		CurrYesPrice: 0,
		CurrNoPrice: 0,
		Buy: BuyOrderbook{
			Type: map[StockTypeEnum]StrikePrice{},
		},
		Sell: SellOrderbook{
			Type: map[StockTypeEnum]StrikePrice{
				Yes : {
					Strike: map[int]Orders{},
				},
				No : {
					Strike: map[int]Orders{},
				},
			},
		},
	}
	return &ob
}

func (ob *ORDERBOOK)PlaceSellOrder(userId,stockType string, quantity,price int){
	var st StockTypeEnum
	if stockType == "yes" {
		st = Yes
	}else if stockType == "no" {
		st = No
	}
	if exists := ob.CheckBuyer(st,price); !exists {
		if _,typeExists := ob.Sell.Type[st]; !typeExists {
			ob.Sell.Type[st] = StrikePrice{
				Strike: map[int]Orders{},
			}
		}
		if _,strikeExists := ob.Sell.Type[st].Strike[price]; !strikeExists{
			ob.Sell.Type[st].Strike[price] = Orders{
				TotalOrders: 0,
				TimeStamp: map[int]User{},
			}
		}

		strikePrice := ob.Sell.Type[st].Strike[price]
		strikePrice.TotalOrders += quantity
		strikePrice.TimeStamp[int(time.Now().Unix())] = User{
			UserId: userId,
			Quantity: quantity,
		}

		ob.Sell.Type[st].Strike[price] = strikePrice
	}
}


func (ob *ORDERBOOK)CheckBuyer(st StockTypeEnum,price int) bool {
	_,exists := ob.Buy.Type[st].Strike[price]
	return exists
}
func (ob *ORDERBOOK)CheckSeller(st StockTypeEnum,price,quantity int) bool {
	_,availabe := ob.Sell.Type[st]
	if !availabe {
		return availabe
	}
	fillable := ob.CalcFillableBuyQty(price,quantity,st)
	return fillable > 0
}

func (ob *ORDERBOOK)PlaceBuyOrder(userId,stockType string, quantity,price int, e *Engine){
	var st StockTypeEnum
	var mst StockEnum
	if stockType == "yes" {
		st = Yes
		mst = YesStock
	}else if stockType == "no" {
		st = No
		mst = NoStock
	}

	isAvailable := ob.CheckSeller(st,price,quantity)
	if !isAvailable {
		fmt.Println("No orders can be matched at the moment.")
		return
	}

	ob.matchOrder(userId,price,quantity,st,mst,e)
}

func (ob *ORDERBOOK)matchOrder(userId string ,price,qty int, st StockTypeEnum,mst StockEnum,e *Engine){
	var filledQty int
	pendingOrders := qty
	fillable := ob.CalcFillableBuyQty(price,qty,st)
	strikes := ob.Sell.Type[st]
	sortedStrikeKeys := maps.Keys(strikes.Strike)
	sort.Ints(sortedStrikeKeys)
	fmt.Println("fillable",fillable)
	for _,currStrike := range sortedStrikeKeys{
		order := strikes.Strike[currStrike]
		if  filledQty < fillable{
			filled := math.Min(float64(pendingOrders),float64(order.TotalOrders))
			order.TotalOrders -= int(filled)
			timestamps := &order.TimeStamp
			ob.manageStocksAndInr(userId,int(filled),currStrike,mst,e,timestamps)
			pendingOrders -= int(filled)
			filledQty += int(filled)
			ob.Sell.Type[st].Strike[currStrike] = order
		}
	}
	fmt.Println(ob.Sell.Type[st])
	fmt.Println("pending :",pendingOrders)
	fmt.Println("filled :",filledQty)
}



func(ob *ORDERBOOK)manageStocksAndInr(buyer string, toFill,strike int,mst StockEnum, e *Engine, ts *map[int]User){
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
			
			sellerStocks := e.StockBalance.User[seller].Symbol[ob.StockSymbol].Type[mst]
			sellerStocks.Locked -= int(toDeduct)
			e.StockBalance.User[seller].Symbol[ob.StockSymbol].Type[mst] = sellerStocks

			_,existing := e.StockBalance.User[buyer].Symbol[ob.StockSymbol]
			if !existing {
				e.StockBalance.User[buyer].Symbol[ob.StockSymbol] = StockType{
					Type: map[StockEnum]Quantity{
						YesStock: {
							Available: 0,
							Locked: 0,
						},
						NoStock :{
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
func (ob *ORDERBOOK)manageINR(buyer,seller string, qty,price int, e *Engine){
	buyerBal := e.InrBalance.User[buyer]
	sellerBal := e.InrBalance.User[seller]

	buyerBal.Locked -= qty * price
	sellerBal.Balance += qty * price

	e.InrBalance.User[buyer] = buyerBal
	e.InrBalance.User[seller] = sellerBal
}

func (ob *ORDERBOOK)CalcFillableBuyQty( price,qty int, st StockTypeEnum) int {
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

