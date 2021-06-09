package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	client.Ping(context.Background()).Result()
	log.Println(client.String())
	client.Set(context.Background(), "foo", "bar", 0)
	log.Println(client.Get(context.Background(), "foo").Result())
}
