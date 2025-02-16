package main

import (
	"fmt"

	"github.com/arvindkhoisnam/go_probo_engine/packages"
)

func main(){
	engine := packages.InitEngine()
	engine.CreateMarket("AAPL")
	engine.CreateUser("user1")
	engine.OnrampINR("user1",10000)
	all := engine.AllMarkets()
	fmt.Println(all)
	fmt.Println("---------------------")
	engine.Mint("user1","AAPL",5,1000)
	engine.GetInrBal("user1")
	engine.GetStockBal("user1")
}