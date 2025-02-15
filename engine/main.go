package main

import (
	"github.com/arvindkhoisnam/go_probo_engine/packages"
)

func main(){
	engine := packages.InitEngine()
	// engine.CreateMint("AAPL")
	// engine.CreateMint("TSLA")
	// engine.CreateMint("AAPL")
	// engine.CreateMint("NFLX")
	engine.CreateUser("user1")
	engine.CreateUser("user2")
	engine.OnrampINR("user1",1000)
	engine.OnrampINR("user3",1000)
	engine.OnrampINR("user2",1000)
	// fmt.Println(engine)
}