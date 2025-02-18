package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/arvindkhoisnam/go_probo_engine/packages"
	"github.com/arvindkhoisnam/go_probo_engine/redisManager"
	"github.com/redis/go-redis/v9"
)

func main(){
	redisClient := redisManager.InitRedis()
	engine := packages.InitEngine()
	ctx := context.Background()
	for {
		// Continuously pop messages from the Redis list
		res, err := redisClient.RPop(ctx, "engine").Result()
		if err == redis.Nil {
			// No messages in queue, continue polling
			continue
		} else if err != nil {
			log.Println("Redis error:", err)
			continue
		}
		var parsed packages.Incoming
		json.Unmarshal([]byte(res),&parsed) 
		// engine.StartEngine(parsed)
		engine.StartEngine(parsed)
	}
}