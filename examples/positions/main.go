package main

import (
	"context"
	"fmt"
	"log"
	"os"

	rest "github.com/roundinternetmoney/ethereal-rest"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pk := os.Getenv("ETHEREAL_PK")
	if pk == "" {
		log.Fatal("ETHEREAL_PK is required (hex private key, with or without 0x prefix)")
	}

	client, err := rest.NewClient(ctx, pk, rest.Testnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}

	positions, err := client.GetPosition(ctx)
	if err != nil {
		log.Fatalf("failed to get positions: %v", err)
	}

	if len(positions) == 0 {
		fmt.Println("no open positions")
		return
	}

	for _, p := range positions {
		fmt.Println(p)
	}
}
