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
//////////////////////////////////////////
type TickerType struct {
	Ticker string
	YesPrice int
	NoPrice int
}
type DepthType struct {
	YesMarket map[int]int
	NoMarket map[int]int
}
type StockTypeEnum int

const (
	Yes StockTypeEnum = iota 
	No
)

type User struct {
	ReverseOrder bool
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

type Market struct{
	StockSymbol string
	CurrYesPrice int
	CurrNoPrice int
}
