package models

type UserBalance struct {
	Balance int
	Locked  int
}
type INR_BALANCE struct {
	User  map[string]UserBalance
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

type STOCK_BALANCE struct {
	User map[string] StockSymbol
}

type Engine struct {
	Markets []ORDERBOOK
	InrBalance INR_BALANCE
	StockBalance STOCK_BALANCE
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

type Market struct{
	StockSymbol string
	CurrYesPrice int
	CurrNoPrice int
}
