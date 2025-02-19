package packages

import (
	"fmt"
)

type StockTypeEnum int
const (
	Yes StockTypeEnum = iota 
	No
)
type Orders struct {
	TotalOrders int
	Order map[string] int
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
				Order: map[string]int{},
			}
		}

		existingStrike := ob.Sell.Type[st].Strike[price]
		existingStrike.TotalOrders += quantity
		if existingUser, userExists := existingStrike.Order[userId]; userExists {
			existingStrike.Order[userId] = existingUser + quantity
		} else {
				existingStrike.Order[userId] = quantity
		}
		ob.Sell.Type[st].Strike[price] = existingStrike
		fmt.Println(ob)
	}
}


func (ob *ORDERBOOK)CheckBuyer(st StockTypeEnum,price int) bool {
	_,exists := ob.Buy.Type[st].Strike[price]
	return exists
}
// func (ob *ORDERBOOK)CheckSeller(st StockTypeEnum,price int) bool {
// 	_,exists := ob.Sell.Type[st].
// }

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
	sellOB := ob.Sell.Type[st].Strike[price]
	sellOB.TotalOrders -= quantity
	sellOB.Order["user1"] -= quantity
	
	buyer, buyerStocks := e.StockBalance.User[userId]
	seller, sellerStocks := e.StockBalance.User["user1"]

	if sellerStocks {
		stocks := seller.Symbol[ob.StockSymbol].Type[mst]
		stocks.Locked -= quantity
		e.StockBalance.User["user1"].Symbol[ob.StockSymbol].Type[mst] = stocks
	}
	if !buyerStocks {
		buyer = StockSymbol{
			Symbol: map[string]StockType{},
		}
	}
	stocks, stockExists := buyer.Symbol[ob.StockSymbol]
	if !stockExists {
		stocks = StockType{
			Type: map[StockEnum]Quantity{
				mst :{
					Available: quantity,
					Locked: 0,
				},
			},
		}
	} else {
		stocks.Type[mst] = Quantity{
			Available:  stocks.Type[mst].Available + quantity,
			Locked:  stocks.Type[mst].Locked,
		}
	}

	if buyer.Symbol == nil {
		buyer.Symbol = make(map[string]StockType)
	}

	buyer.Symbol[ob.StockSymbol] = stocks
	e.StockBalance.User[userId] = buyer

	buyerINR := e.InrBalance.User[userId]
	sellerINR := e.InrBalance.User["user1"]

	buyerINR.Locked -= quantity*price
	sellerINR.Balance += quantity*price

	e.InrBalance.User[userId] = buyerINR
	e.InrBalance.User["user1"] = sellerINR
}

func CheckFillableBuyQty(ob *ORDERBOOK, price,qty int){
	strikes := ob.Sell.Type[Yes]
	var temp int
	for strike, order := range strikes.Strike {
		if strike <= price && temp <= qty{
			temp += order.TotalOrders
		}
	}
}
