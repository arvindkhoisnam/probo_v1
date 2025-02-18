package packages

import (
	"fmt"
)

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

type BuyOrderbook struct {
	Type map[OrderType] StrikePrice
}
type SellOrderbook struct {
	Type map[OrderType] StrikePrice
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
			Type: map[OrderType]StrikePrice{},
		},
		Sell: SellOrderbook{
			Type: map[OrderType]StrikePrice{
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
	var st OrderType
	if stockType == "yes" {
		st = Yes
	}else if stockType == "no" {
		st = No
	}
    // strike := ob.Sell.Type[st].Strike[price]
    strike := ob.Sell.Type[st].Strike[price]
	strike.TotalOrders += quantity
	if strike.Order == nil {
		strike.Order = map[string]int{}
	}
	strike.Order[userId] += quantity
	ob.Sell.Type[st].Strike[price] = strike
	fmt.Println(ob.Sell)
}
