package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/qiwi1272/ethereal-go"
	"github.com/qiwi1272/ethereal-go/pb"
	"github.com/qiwi1272/ethereal-go/rest"
	ws "github.com/qiwi1272/ethereal-go/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// load all the symbols using a rest client
	rest, err := rest.NewClient(ctx, os.Getenv("ETHEREAL_PK"), rest.Testnet)
	if err != nil {
		panic(err)
	}
	sid := rest.Subaccount.Id

	var symbols map[string]ethereal.Product
	if symbols, err = rest.GetProductMap(ctx); err != nil {
		panic(err)
	}

	ws := ws.NewClient(ctx, ws.Testnet)
	defer ws.Close()

	diffCb := func(m proto.Message) {
		diff := m.(*pb.L2Book)
		fmt.Printf("diff: %v\n", diff)
	}

	cb := func(m proto.Message) {
		fmt.Println(m)
	}

	// many different subscription examples, all achieving the same result.
	// error channel for callbacks, ethereal.Subscription namespace, "resubscribe" intent

	// for all exchange symbols
	for symbolKey := range symbols {
		// callbacks here are overwritten every iteration.
		// subscribe to book with a protobuf enum, and callback to diffCB
		if err := ws.SubscribeWithCallback(ctx, pb.EventType_L2_BOOK, symbolKey, diffCb); err != nil {
			log.Fatal("EventType_EVENT_TYPE_L2_BOOK: ", err)
		}
		// subscribe to trade fill with a protobuf struct
		if bytes, err := new(pb.Ticker).MarshalIntent(symbolKey, pb.Sub); err == nil {
			ws.Req(ctx, bytes)
		} else {
			log.Fatal("OrderUpdate: ", err)
		}
		ws.OnEvent(new(pb.OrderUpdate), cb)
		// subscribe to trade fill with a protobuf struct
		if err := ws.Subscribe(ctx, &pb.TradeFill{}, symbolKey); err != nil {
			log.Fatal("TradeFill: ", err)
		}
		ws.OnEvent(&pb.TradeFill{}, cb)
	}

	// subscribe to a SubaccountLiquidation protobuf struct, with a callback
	if err := ws.SubscribeWithCallback(ctx, &pb.SubaccountLiquidation{}, sid, cb); err != nil {
		log.Fatal("SubaccountLiquidation: ", err)
	}
	// subscribe to a order fill protobuf enum, with a callback
	if err := ws.SubscribeWithCallback(ctx, pb.EventType_ORDER_FILL, sid, func(m proto.Message) {
		fill := m.(*pb.OrderFill)
		fmt.Println(fill)
	}); err != nil {
		log.Fatal("EventType_EVENT_TYPE_ORDER_FILL: ", err)
	}

	// underlying calls of subscribing to a protobuf struct
	if bytes, err := new(pb.OrderUpdate).MarshalIntent(sid, pb.Sub); err == nil {
		ws.Req(ctx, bytes)
	} else {
		log.Fatal("OrderUpdate: ", err)
	}
	ws.OnEvent(new(pb.OrderUpdate), cb)
	// underlying calls of subscribing to a protobuf enum
	if bytes, err := pb.EventType_TOKEN_TRANSFER.MarshalIntent(sid, pb.Sub); err == nil {
		ws.Req(ctx, bytes)
	} else {
		log.Fatal("EventType_EVENT_TYPE_TRANSFER: ", err)
	}
	ws.OnEvent(pb.EventType_TOKEN_TRANSFER, cb)

	errCh := make(chan error, 1)
	go func() {
		errCh <- ws.Listen(ctx) // blocking
	}()

	select {
	case <-ctx.Done():
		return
	case err := <-errCh:
		if err != nil && ctx.Err() == nil {
			log.Fatal(err)
		}
	}

}
