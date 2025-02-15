package packages

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
			Type: map[orderType]strikePrice{
				Yes: {},
				No:  {},
			},
		},
		Sell: sellOrderbook{
			Type: map[orderType]strikePrice{
				Yes :{},
				No:{},
			},
		},
	}
	return &ob
}