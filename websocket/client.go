package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/qiwi1272/ethereal-go/pb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
)

type Environment string

const (
	Testnet Environment = "wss://ws2.etherealtest.net/v1/stream"
	Mainnet Environment = "wss://ws2.ethereal.trade/v1/stream"
)

type Client struct {
	Con                *websocket.Conn
	conMu              *sync.Mutex
	env                Environment
	bookHandler        func(*pb.L2Book)      // non-array
	priceHandler       func(*pb.MarketPrice) // non-array
	tradeFillHandler   func(*pb.TradeFillEvent)
	liquidationHandler func(*pb.SubaccountLiquidationEvent)
	orderUpdateHandler func(*pb.OrderUpdateEvent)
	orderFillHandler   func(*pb.OrderFillEvent)
	transferHandler    func(*pb.Transfer) // non-array
	hbCancel           context.CancelCauseFunc
}

func NewClient(parent context.Context, env Environment) *Client {
	ctx, cancel := context.WithCancelCause(parent)
	c, _, err := websocket.Dial(ctx, string(env), nil)
	if err != nil {
		log.Fatal(err)
	}

	cl := &Client{
		Con:   c,
		conMu: &sync.Mutex{},
		env:   env,
	}

	cl.keepalive(ctx, cancel)
	cl.hbCancel = cancel

	return cl
}

type Intent string

const (
	sub   Intent = "subscribe"
	unsub Intent = "unsubscribe"
)

type SubIntent[T eventData] struct {
	I Intent    `json:"event"`
	D eventData `json:"data"`
}

func marshalSubscribe[T eventData, I Intent](data T) (b []byte, err error) {
	req := &SubIntent[T]{I: sub, D: data}
	return json.Marshal(req)
}

func marshalUnsubscribe[T eventData](data T) (b []byte, err error) {
	req := &SubIntent[T]{I: unsub, D: data}
	return json.Marshal(req)
}

func marshalToValueCallback[T proto.Message](data []byte, pb T, cb func(T)) (err error) {
	if err := protojson.Unmarshal(data, pb); err != nil {
		return err
	}
	cb(pb)
	return
}

type eventData interface{}

type SymbolEvent struct {
	eventData
	T string `json:"type"`
	S string `json:"symbol"`
}

type SubaccountEvent struct {
	eventData
	T string `json:"type"`
	S string `json:"subaccountId"`
}

func (c *Client) req(ctx context.Context, payload []byte) (err error) {
	return c.Con.Write(ctx, websocket.MessageBinary, payload)
}

func (c *Client) SubscribeBook(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SymbolEvent{
		T: "L2Book",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeBook(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SymbolEvent{
		T: "L2Book",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeMarketPrice(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SymbolEvent{
		T: "MarketPrice",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeMarketPrice(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SymbolEvent{
		T: "MarketPrice",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeFill(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SymbolEvent{
		T: "TradeFill",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeFill(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SymbolEvent{
		T: "TradeFill",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeLiquidation(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SubaccountEvent{
		T: "SubaccountLiquidation",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeLiquidation(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SubaccountEvent{
		T: "SubaccountLiquidation",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeOrderUpdate(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SubaccountEvent{
		T: "OrderUpdate",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeOrderUpdate(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SubaccountEvent{
		T: "OrderUpdate",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeOrderFill(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SubaccountEvent{
		T: "OrderFill",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeOrderFill(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SubaccountEvent{
		T: "OrderFill",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeTokenTransfer(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SubaccountEvent{
		T: "TokenTransfer",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeTokenTransfer(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SubaccountEvent{
		T: "TokenTransfer",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) OnBook(callback func(*pb.L2Book)) {
	c.bookHandler = callback
}

func (c *Client) OnPrice(callback func(*pb.MarketPrice)) {
	c.priceHandler = callback
}

func (c *Client) OnTradeFill(callback func(*pb.TradeFillEvent)) {
	c.tradeFillHandler = callback
}

func (c *Client) OnLiquidation(callback func(*pb.SubaccountLiquidationEvent)) {
	c.liquidationHandler = callback
}

func (c *Client) OnOrderUpdate(callback func(*pb.OrderUpdateEvent)) {
	c.orderUpdateHandler = callback
}

func (c *Client) OnOrderFill(callback func(*pb.OrderFillEvent)) {
	c.orderFillHandler = callback
}

func (c *Client) OnTransfer(callback func(*pb.Transfer)) {
	c.transferHandler = callback
}

type wssMsg struct {
	Event string          `json:"e"`
	Ts    int64           `json:"t"`
	Data  json.RawMessage `json:"data"`
}

func (c *Client) Listen(parent context.Context) error {
	ctx, cancel := context.WithCancelCause(parent)
	defer cancel(nil)
	defer c.Close()

	for {
		_, data, err := c.Con.Read(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return context.Cause(ctx)
			}
			cancel(err)
			return err
		}

		var msg wssMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			cancel(err)
			return err
		}

		switch msg.Event {
		case "L2Book":
			var diff pb.L2Book
			if err := marshalToValueCallback(msg.Data, &diff, c.bookHandler); err != nil {
				cancel(err)
				return err
			}

		case "MarketPrice":
			var mp pb.MarketPrice
			if err := marshalToValueCallback(msg.Data, &mp, c.priceHandler); err != nil {
				cancel(err)
				return err
			}

		case "SubaccountLiquidation":
			var lq pb.SubaccountLiquidationEvent
			if err := marshalToValueCallback(msg.Data, &lq, c.liquidationHandler); err != nil {
				fmt.Println(string(data))
				cancel(err)
				return err
			}

		case "OrderFill":
			var ou pb.OrderFillEvent
			if err := marshalToValueCallback(data, &ou, c.orderFillHandler); err != nil {
				cancel(err)
				return err
			}

		case "OrderUpdate":
			var ou pb.OrderUpdateEvent
			if err := marshalToValueCallback(data, &ou, c.orderUpdateHandler); err != nil {
				cancel(err)
				return err
			}

		case "TradeFill":
			var tf pb.TradeFillEvent
			if err := marshalToValueCallback(data, &tf, c.tradeFillHandler); err != nil {
				cancel(err)
				return err
			}

		case "TokenTransfer":
			var t pb.Transfer
			if err := marshalToValueCallback(data, &t, c.transferHandler); err != nil {
				fmt.Println(string(data))
				cancel(err)
				return err
			}

		default:
			fmt.Printf("unknown event, raw: %s\n", string(data))
		}
	}
}

func (c *Client) keepalive(ctx context.Context, cancel context.CancelCauseFunc) {
	go func() {
		t := time.NewTicker(20 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				// Ping will return error if connection is dead
				if err := c.Con.Ping(ctx); err != nil {
					cancel(err)
					return
				}
			}
		}
	}()
}

func (c *Client) Resubscribe(parent context.Context) error {
	c.Close()

	c.conMu.Lock()
	defer c.conMu.Unlock()

	ctx, cancel := context.WithCancelCause(parent)

	// replace con and restart listener with new context
	var err error
	c.Con, _, err = websocket.Dial(ctx, string(c.env), nil)
	if err != nil {
		cancel(err)
		return err
	}
	c.hbCancel = cancel
	c.keepalive(ctx, cancel)

	return nil
}

func (c *Client) Close() {
	c.conMu.Lock()
	defer c.conMu.Unlock()
	if c.hbCancel != nil {
		c.hbCancel(nil)
	}
	if c.Con != nil {
		c.Con.Close(websocket.StatusNormalClosure, "closing")
	}
}
