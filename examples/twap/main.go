// TWAP-style execution (composition pattern).
//
// This is not a native TWAP or time-sliced order type in ethereal-rest: it is
// user-land scheduling that submits ordinary limit (or other) orders on a timer.
// Use tiny size on testnet; stop or back off on repeated API errors.
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

	slices := 3
	interval := 2 * time.Second
	qtyPerSlice := 0.01
	// Far-from-market prices for demo only.
	basePx := 1000.0

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var placed []rest.OrderCreated
	for i := 0; i < slices; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}

		px := basePx + float64(i)*0.1
		o, err := client.CreateOrder(ctx, p.NewOrder(rest.ORDER_LIMIT, qtyPerSlice, px, false, rest.BUY, rest.TIF_GTD))
		if err != nil {
			log.Fatalf("slice %d: create order: %v", i, err)
		}
		placed = append(placed, o)
		log.Printf("slice %d placed: %+v", i, o)
	}

	// Best-effort cleanup: one cancel request with all venue order IDs.
	created := make([]*rest.OrderCreated, len(placed))
	for i := range placed {
		created[i] = &placed[i]
	}
	if _, err := client.CancelOrdersFromCreated(ctx, created); err != nil {
		log.Fatalf("cancel TWAP legs: %v", err)
	}
}
