package main

import (
	"fmt"

	"github.com/subhammahanty235/bizzgenie/internal/redisclient"
)

type BizzGenie struct {
	internalRedis *redisclient.NewInternalRedisClient
	externalRedis *redisclient.NewExternalRedisClient
}

func main() {
	fmt.Println("Bizz Genie is Running and Ready to process")
	//handshake mechanism --> times based

	//health checkup --> timer based on regular intervals

	//redis instance health update --> while start
}
