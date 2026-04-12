// Chase / follow-the-market (composition pattern).
//
// There is no "chase" or pegged-order API in ethereal-rest. This example shows
// cancel-then-replace driven by a reference price that you must obtain yourself
// (exchange feed, indexer, etc.). The stub below uses a fake moving reference.
//
// Races: if the market moves between cancel and create, you can end up with no
// working order or an unexpected fill. Use ClientOrderId on Order when the venue
// supports idempotent client keys.
package main

import (
	"context"
	"log"
	"os"
	"time"

	rest "github.com/roundinternetmoney/ethereal-rest"
)

// referencePriceStub stands in for an external oracle or ticker. Replace with
// your own data source; nothing in client.go provides live quotes for this loop.
func referencePriceStub(iter int) float64 {
	return 1000.0 + float64(iter)*0.5
}

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

	iterations := 2
	qty := 0.01

	var active *rest.OrderCreated
	for i := 0; i < iterations; i++ {
		ref := referencePriceStub(i)
		// Place or replace at reference (offset for a passive resting quote in demo).
		targetPx := ref - 50.0

		if active != nil {
			cancelReq := rest.NewCancelOrderFromCreated(active)
			if _, err := client.CancelOrder(ctx, cancelReq); err != nil {
				log.Fatalf("iteration %d: cancel: %v", i, err)
			}
			active = nil
		}

		placed, err := client.CreateOrder(ctx, p.NewOrder(rest.ORDER_LIMIT, qty, targetPx, false, rest.BUY, rest.TIF_GTD))
		if err != nil {
			log.Fatalf("iteration %d: create: %v", i, err)
		}
		active = &placed
		log.Printf("iteration %d: ref=%.4f limit=%.4f order=%+v", i, ref, targetPx, placed)
		time.Sleep(500 * time.Millisecond)
	}

	if active != nil {
		if _, err := client.CancelOrder(ctx, rest.NewCancelOrderFromCreated(active)); err != nil {
			log.Fatalf("final cancel: %v", err)
		}
	}
}
