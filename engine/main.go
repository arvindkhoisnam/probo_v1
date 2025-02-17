package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/arvindkhoisnam/go_probo_engine/packages"
	"github.com/redis/go-redis/v9"
)

func main(){

	opt, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		panic(err)
	}
	
	client := redis.NewClient(opt)

	ctx := context.Background()
	engine := packages.InitEngine()
fmt.Println(engine)
	for {
		// Continuously pop messages from the Redis list
		res, err := client.RPop(ctx, "engine").Result()
		if err == redis.Nil {
			// No messages in queue, continue polling
			continue
		} else if err != nil {
			log.Println("Redis error:", err)
			continue
		}

		// Process the received message
		var parsed packages.Incoming
		json.Unmarshal([]byte(res),&parsed) 
		switch parsed.Event {
		case "onramp":
			fmt.Println(parsed.Payload.UserId)
			fmt.Println(parsed.Payload.Amount)
		}
	}
}