package packages

import "fmt"

type orderType int
const (
	Yes orderType = iota 
	No
)
type orders struct {
	TotalOrders int
	Order map[string] int
}

type strikePrice struct {
	Strike map[int] orders
}

type buyOrderbook struct {
	Type map[orderType] strikePrice
}
type sellOrderbook struct {
	Type map[orderType] strikePrice
}
type ORDERBOOK struct{
	StockSymbol  string
	CurrYesPrice int
	CurrNoPrice  int
	Buy          buyOrderbook
	Sell         sellOrderbook
}

func CreateOrderbook(stock string) *ORDERBOOK {
	ob := ORDERBOOK {
		StockSymbol: stock,
		CurrYesPrice: 0,
		CurrNoPrice: 0,
		Buy: buyOrderbook{
			Type: map[orderType]strikePrice{},
		},
		Sell: sellOrderbook{
			Type: map[orderType]strikePrice{
				Yes : {
					Strike: map[int]orders{},
				},
				No : {
					Strike: map[int]orders{},
				},
			},
		},
	}
	return &ob
}

func (ob *ORDERBOOK)PlaceSellOrder(userId,stockType string, quantity,price int){
	var st orderType
	if stockType == "yes" {
		st = Yes
	}else if stockType == "no" {
		st = No
	}
    strike := ob.Sell.Type[st].Strike[price]
	strike.TotalOrders += quantity
	if strike.Order == nil {
		strike.Order = map[string]int{}
	}
	strike.Order[userId] += quantity
	ob.Sell.Type[st].Strike[price] = strike
	fmt.Println(ob.Sell)
}
func (ob *ORDERBOOK)PlaceBuyOrder(){}