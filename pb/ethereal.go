package pb

import (
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _NO_INTENT_ERROR = errors.New("No marshal intent")

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

type Intent string

const (
	Sub   Intent = "subscribe"
	Unsub Intent = "unsubscribe"
)

type SubscriptionIntent[T eventData] struct {
	I Intent    `json:"event"`
	D eventData `json:"data"`
}

// intent abstraction   |   TODO: Subscription Intent protos
func sub(e eventData) ([]byte, error) {
	return json.Marshal(&SubscriptionIntent[eventData]{I: Sub, D: e})
}

func unsub(e eventData) ([]byte, error) {
	return json.Marshal(&SubscriptionIntent[eventData]{I: Unsub, D: e})
}

type Event[T proto.Message] interface {
	EventName() string
	EventStruct() (Event[T], error)
	MarshalIntent(to string, i Intent) ([]byte, error)
	UnmarshalToCallback(b json.RawMessage, cb func(T)) error
}

/////////////////
// ENUM EVENTS //
/////////////////

// server -> client lookup

func EventEnum(e string) EventType {
	switch e {
	case "L2Book":
		return EventType_L2_BOOK
	case "Ticker":
		return EventType_TICKER
	case "TradeFill":
		return EventType_TRADE_FILL
	case "SubaccountLiquidation":
		return EventType_SUBACCOUNT_LIQUIDATION
	case "PositionUpdate":
		return EventType_POSITION_UPDATE
	case "OrderUpdate":
		return EventType_ORDER_UPDATE
	case "OrderFill":
		return EventType_ORDER_FILL
	case "TokenTransfer":
		return EventType_TOKEN_TRANSFER
	default:
		return EventType_EVENT_UNSPECIFIED
	}
}

// client -> server lookup

func (e EventType) EventName() string {
	switch e {
	case EventType_L2_BOOK:
		return "L2Book"
	case EventType_TICKER:
		return "Ticker"
	case EventType_TRADE_FILL:
		return "TradeFill"
	case EventType_SUBACCOUNT_LIQUIDATION:
		return "SubaccountLiquidation"
	case EventType_POSITION_UPDATE:
		return "PositionUpdate"
	case EventType_ORDER_UPDATE:
		return "OrderUpdate"
	case EventType_ORDER_FILL:
		return "OrderFill"
	case EventType_TOKEN_TRANSFER:
		return "TokenTransfer"
	default:
		return ""
	}
}

func (e EventType) EventStruct() (Event[proto.Message], error) {
	switch e {
	case EventType_L2_BOOK:
		return new(L2Book), nil
	case EventType_TICKER:
		return new(Ticker), nil
	case EventType_TRADE_FILL:
		return new(TradeFill), nil
	case EventType_SUBACCOUNT_LIQUIDATION:
		return new(SubaccountLiquidation), nil
	case EventType_POSITION_UPDATE:
		return new(PositionUpdate), nil
	case EventType_ORDER_UPDATE:
		return new(OrderUpdate), nil
	case EventType_ORDER_FILL:
		return new(OrderFill), nil
	case EventType_TOKEN_TRANSFER:
		return new(TokenTransfer), nil
	default:
		return nil, _NO_INTENT_ERROR
	}
}

func (e EventType) MarshalIntent(to string, i Intent) ([]byte, error) {
	if s, err := e.EventStruct(); err == nil {
		return s.MarshalIntent(to, i)
	} else {
		fmt.Println(err)
	}

	return nil, _NO_INTENT_ERROR
}

func (e EventType) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) error {
	if s, err := e.EventStruct(); err == nil {
		return s.UnmarshalToCallback(b, cb)
	}
	return _NO_INTENT_ERROR
}

/////////////////////
// PROTOBUF EVENTS //
/////////////////////

func (*L2Book) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SymbolEvent{
		T: "L2Book",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*Ticker) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SymbolEvent{
		T: "Ticker",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*TradeFill) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SymbolEvent{
		T: "TradeFill",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*SubaccountLiquidation) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SubaccountEvent{
		T: "SubaccountLiquidation",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*PositionUpdate) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SubaccountEvent{
		T: "PositionUpdate",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*OrderUpdate) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SubaccountEvent{
		T: "OrderUpdate",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*OrderFill) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SubaccountEvent{
		T: "OrderFill",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*TokenTransfer) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SubaccountEvent{
		T: "TokenTransfer",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (l *L2Book) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, l); err == nil {
		cb(l)
	}
	return
}

func (m *Ticker) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, m); err == nil {
		cb(m)
	}
	return
}

func (s *TradeFill) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, s); err == nil {
		cb(s)
	}
	return
}

func (of *SubaccountLiquidation) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, of); err == nil {
		cb(of)
	}
	return
}

func (ou *PositionUpdate) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, ou); err == nil {
		cb(ou)
	}
	return
}

func (tf *OrderUpdate) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, tf); err == nil {
		cb(tf)
	}
	return
}

func (t *OrderFill) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, t); err == nil {
		cb(t)
	}
	return
}

func (t *TokenTransfer) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, t); err == nil {
		cb(t)
	}
	return
}

func (*L2Book) EventName() string {
	return "L2Book"
}
func (*Ticker) EventName() string {
	return "Ticker"
}
func (*TradeFill) EventName() string {
	return "TradeFill"
}
func (*SubaccountLiquidation) EventName() string {
	return "SubaccountLiquidation"
}
func (*PositionUpdate) EventName() string {
	return "PositionUpdate"
}
func (*OrderUpdate) EventName() string {
	return "OrderUpdate"
}
func (*OrderFill) EventName() string {
	return "OrderFill"
}
func (*TokenTransfer) EventName() string {
	return "TokenTransfer"
}

func (p *L2Book) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *Ticker) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *TradeFill) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *SubaccountLiquidation) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *PositionUpdate) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *OrderUpdate) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *OrderFill) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *TokenTransfer) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
