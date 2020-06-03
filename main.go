package main

import (
	"fmt"

	"github.com/lamhai1401/gord/redis"
)

func main() {
	client := redis.NewClient(&redis.Config{
		Addrs: []string{"127.0.0.1:6379"},
		DB:    0,
		Pwd:   "",
	})

	fmt.Print(client)

}
