package main

import (
	"context"
	"log"
	"os"
	"time"

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

	products, err := client.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

	p := products["ETHUSD"]

	// Initial resting limit (far from market on purpose for a safe demo).
	placed, err := client.CreateOrder(ctx, p.NewOrder(rest.ORDER_LIMIT, 0.01, 1000.0, false, rest.BUY, rest.TIF_GTD))
	if err != nil {
		log.Fatalf("failed to place limit order: %v", err)
	}

	cancelReq := rest.NewCancelOrderFromCreated(&placed)
	cancelled, err := client.CancelOrder(ctx, cancelReq)
	if err != nil {
		log.Fatalf("failed to cancel limit order: %v", err)
	}
	for _, c := range cancelled {
		log.Printf("cancelled: %+v", c)
	}

	// Replace: new order with updated price (two server operations; not atomic).
	replaced, err := client.CreateOrder(ctx, p.NewOrder(rest.ORDER_LIMIT, 0.01, 1001.0, false, rest.BUY, rest.TIF_GTD))
	if err != nil {
		log.Fatalf("failed to place replacement order: %v", err)
	}
	log.Printf("replacement placed: %+v", replaced)

	cancelReq = rest.NewCancelOrderFromCreated(&replaced)
	if _, err := client.CancelOrder(ctx, cancelReq); err != nil {
		log.Fatalf("failed to cancel replacement: %v", err)
	}

	time.Sleep(time.Second)
}
