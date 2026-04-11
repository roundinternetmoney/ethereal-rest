package etherealRest

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const USER_AGENT = "ethereal-go/1.0.0dev"

type Environment string

const (
	Testnet Environment = "https://api.etherealtest.net"
	Mainnet Environment = "https://api.ethereal.trade"
)

type Client struct {
	BaseURL string
	Http    *http.Client
	account *Signer
}

func (e *Client) GetSubaccount() *Subaccount {
	return e.account.Subaccount
}

func (e *Client) GetTypes() *abi.TypedData {
	return e.account.GetTypes()
}

func (e *Client) Do(ctx context.Context, method, path string, body any) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, e.BaseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.Http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	out := new(bytes.Buffer)
	_, err = out.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ethereal error %d: %s", resp.StatusCode, out.String())
	}
	return out.Bytes(), nil
}

func NewClient(ctx context.Context, pk string, env Environment) (*Client, error) {
	transport := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   2 * time.Second,
		ExpectContinueTimeout: 0,
		DisableCompression:    true,
		ForceAttemptHTTP2:     true,
	}

	client := &Client{
		BaseURL: string(env),
		Http: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}

	// load pk
	if pk == "" {
		return nil, errors.New("no private key provided; ETHEREAL_PK not set in environment")
	}

	// parse key, set address
	if len(pk) > 1 && pk[:2] == "0x" {
		pk = pk[2:]
	}
	if ecdsa, err := crypto.HexToECDSA(pk); err == nil {
		client.account = NewSigner(ecdsa)
	} else {
		return nil, err
	}
	// ethereal env setup
	var err error
	_, err = client.InitDomain(ctx)
	if err != nil {
		return nil, errors.Join(errors.New("unable to compute domain hash: "), err)
	}

	if err := client.InitSubaccount(ctx); err != nil {
		return nil, errors.Join(errors.New("failed to fetch subaccount: "), err)
	}

	return client, nil
}

// ---------- REST ----------

type Response[T any] struct {
	Data T `json:"data"`
}

// ---------- Setup ----------
func (e *Client) InitDomain(ctx context.Context) (string, error) {
	// init eip 712 data from rpc
	data, err := e.Do(ctx, "GET", "/v1/rpc/config", nil)
	if err != nil {
		return "", err
	}
	var resp struct {
		Domain   abi.TypedDataDomain `json:"domain"`
		SigTypes map[string]string   `json:"signatureTypes"`
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return "", err
	}

	// parse flattened type data
	parsedTypes := abi.Types{}
	for primaryType, schema := range resp.SigTypes {
		types, err := ParseTypeSchema(schema)
		if err != nil {
			return "", err
		}
		parsedTypes[primaryType] = types
	}
	// hardcode domain type
	parsedTypes["EIP712Domain"] = []abi.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	}

	types := &abi.TypedData{
		Types:  parsedTypes,
		Domain: resp.Domain,
	}

	e.account.SetTypes(types)

	domain, err := types.HashStruct("EIP712Domain", types.Domain.Map())
	if err != nil {
		panic("failed to compute domain hash: " + err.Error())
	}
	DomainHash = domain
	return hex.EncodeToString(domain), nil
}

func (e *Client) InitSubaccount(ctx context.Context) error {
	path := fmt.Sprintf("/v1/subaccount?sender=%s", e.account.Address)
	data, err := e.Do(ctx, "GET", path, nil)
	if err != nil {
		return err
	}
	var resp Response[[]Subaccount]
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	if len(resp.Data) == 0 {
		return errors.New("no subaccounts found")
	}
	e.account.Subaccount = &resp.Data[0] // NOTE: currently only one subaccount per client is supported

	return nil
}

// ---------- Methods ----------

type sendable interface {
	OrderCreated | []*OrderCancelled
}

type Sendable[T sendable] interface {
	Signable
	Build(SubaccountHolder)
	Send(context.Context, OrderClient, *Signer) (T, error)
}

func Send[T sendable](ctx context.Context, cl *Client, msg Sendable[T]) (resp T, err error) {
	return msg.Send(ctx, cl, cl.account)
}

func (e *Client) CreateOrder(ctx context.Context, msg Sendable[OrderCreated]) (resp OrderCreated, err error) {
	return msg.Send(ctx, e, e.account)
}

func (e *Client) CreateOrders(ctx context.Context, msg []*Order) (resp []*OrderCreated, err error) {
	order := NewOrderBatch(msg)
	return order.SendBatch(ctx, e, Create, e.account)
}

func (e *Client) CancelOrder(ctx context.Context, msg Sendable[[]*OrderCancelled]) (resp []*OrderCancelled, err error) {
	return msg.Send(ctx, e, e.account)
}

func (e *Client) CancelOrdersFromCreated(ctx context.Context, msg []*OrderCreated) (resp []*OrderCancelled, err error) {
	cancel := NewCancelOrderFromCreated(msg...)
	return cancel.Send(ctx, e, e.account)
}

func (e *Client) GetPosition(ctx context.Context) ([]*Position, error) {
	path := fmt.Sprintf("/v1/position?subaccountId=%s&open=%v", e.account.Subaccount.Id, true)
	data, err := e.Do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp Response[[]*Position]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (e *Client) GetAccountBalance(ctx context.Context) ([]*AccountBalance, error) {
	path := fmt.Sprintf("/v1/subaccount/balance?subaccountId=%s", e.account.Subaccount.Id)
	data, err := e.Do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	var resp Response[[]*AccountBalance]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (e *Client) GetProductMap(ctx context.Context) (map[string]Product, error) {
	data, err := e.Do(ctx, "GET", "/v1/product", nil)
	if err != nil {
		return nil, err
	}
	var resp Response[[]Product]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	products := make(map[string]Product)

	for _, p := range resp.Data {
		products[p.Ticker] = p
	}

	return products, nil
}
