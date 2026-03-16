package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/qiwi1272/ethereal-go"
	"github.com/qiwi1272/ethereal-go/rest"
)

func getTestOrder() ethereal.Order {
	return ethereal.Order{
		Sender:     "0xdeadbeef00000000000000000000000000000000",
		Subaccount: "0x123456789abcde00000000000000000000000000000000000000000000000000",
		Quantity:   "1",
		Price:      "3000",
		ReduceOnly: false,
		Side:       ethereal.BUY,
		EngineType: ethereal.PERPETUAL,
		OnchainID:  2, // later -> ProductId
		Nonce:      "1764897077655477722",
		SignedAt:   int64(1764897077),
	}
}

func reverseHex(s string) ([]byte, error) {
	clean := strings.TrimPrefix(s, "0x")
	return hex.DecodeString(clean)
}
func TestOrders(t *testing.T) {
	orderType := abi.TypedData{
		Types: abi.Types{"TradeOrder": []abi.Type{
			{Name: "sender", Type: "address"},
			{Name: "subaccount", Type: "bytes32"},
			{Name: "quantity", Type: "uint128"},
			{Name: "price", Type: "uint128"},
			{Name: "reduceOnly", Type: "bool"},
			{Name: "side", Type: "uint8"},
			{Name: "engineType", Type: "uint8"},
			{Name: "productId", Type: "uint32"},
			{Name: "nonce", Type: "uint64"},
			{Name: "signedAt", Type: "uint64"},
		}},
	}
	order := getTestOrder()

	message, err := order.ToMessage()
	if err != nil {
		panic(err)
	}
	// We do a pretty print of the message to visually inspect it during test runs.
	// We convert it to to a json string for better readability.
	messageJSON, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Order Message JSON:\n%s\n", string(messageJSON))

	SenderBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][0].Type, message["sender"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(SenderBytes) != "000000000000000000000000deadbeef00000000000000000000000000000000" {
		panic("SenderBytes")
	}

	SubaccountBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][1].Type, message["subaccount"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(SubaccountBytes) != "123456789abcde00000000000000000000000000000000000000000000000000" {
		panic("SubaccountBytes")
	}

	QuantityBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][2].Type, message["quantity"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(QuantityBytes) != "000000000000000000000000000000000000000000000000000000003b9aca00" {
		panic("QuantityBytes")
	}

	PriceBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][3].Type, message["price"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(PriceBytes) != "000000000000000000000000000000000000000000000000000002ba7def3000" {
		panic("PriceBytes")
	}

	ReduceOnlyBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][4].Type, message["reduceOnly"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(ReduceOnlyBytes) != "0000000000000000000000000000000000000000000000000000000000000000" {
		panic("ReduceOnlyBytes")
	}

	SideBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][5].Type, message["side"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(SideBytes) != "0000000000000000000000000000000000000000000000000000000000000000" {
		panic("SideBytes")
	}

	EngineTypeBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][6].Type, message["engineType"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(EngineTypeBytes) != "0000000000000000000000000000000000000000000000000000000000000000" {
		panic("EngineTypeBytes")
	}

	OnchainIDBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][7].Type, message["productId"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(OnchainIDBytes) != "0000000000000000000000000000000000000000000000000000000000000002" {
		panic("OnchainIDBytes")
	}

	NonceBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][8].Type, message["nonce"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(NonceBytes) != "000000000000000000000000000000000000000000000000187e2c8a92c79dda" {
		panic("NonceBytes")
	}

	SignedAtBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][9].Type, message["signedAt"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(SignedAtBytes) != "0000000000000000000000000000000000000000000000000000000069323135" {
		panic("SignedAtBytes")
	}
}

type TestSigner struct {
	addr  string
	pk    *ecdsa.PrivateKey
	types *abi.TypedData
}

func getSigner(pk string) (*TestSigner, error) {
	if ecdsa, err := crypto.HexToECDSA(pk); err == nil {
		return &TestSigner{
			addr: crypto.PubkeyToAddress(ecdsa.PublicKey).Hex(),
			pk:   ecdsa,
		}, nil
	} else {
		return nil, err
	}
}

func (s *TestSigner) GetPk() *ecdsa.PrivateKey {
	return s.pk
}

func (s *TestSigner) GetTypes() *abi.TypedData {
	return s.types
}

// We test signing the order here as well, to ensure that the message encoding is correct.
// If the encoding is wrong, the signature will also be wrong.
func TestOrderSigning(t *testing.T) {

	order := getTestOrder()
	cxt := context.Background()
	pk := "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a" // private key for 0xdeadbeef...

	signer, err := getSigner(pk)
	client, err := rest.NewClient(cxt, pk, rest.Testnet)
	account := client.GetSubaccount()

	signer.types = client.GetTypes()

	fmt.Println("Expected Signer address: ", account.Account)

	domainHashString, err := client.InitDomain(cxt)
	if err != nil {
		panic(err)
	}
	fmt.Println("Domain Hash:", domainHashString)

	expectedDomainHash := "baf501bc2614cf7092d082742580b04c176be1815f46e407eab1bc37ba543c05"
	if domainHashString != expectedDomainHash {
		panic("Domain hash does not match expected value")
	}

	msg, err := order.ToMessage()
	if err != nil {
		panic("Unable to convert order to message: " + err.Error())
	}
	messageHash, err := client.GetTypes().HashStruct("TradeOrder", msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Message Hash:", common.Bytes2Hex(messageHash))
	signature, err := ethereal.Sign(&order, "TradeOrder", signer)
	if err != nil {
		panic(err)
	}

	domainBytes, err := reverseHex(domainHashString)
	if err != nil {
		panic(err)
	}
	fullHash := ethereal.MakeFullHash(domainBytes, messageHash)
	fmt.Println("Full Hash:", common.Bytes2Hex(fullHash))

	expectedSignature := "0x82aed7486e9855459f58537e413760597e689d3ba7b859f56b6edc730e044fff2888ccf92cd282a8299d8d6a76f8bf0aa93d97f30340c4bb0d27b626aca62f211b"
	if signature != expectedSignature {
		panic("Signature does not match expected value")
	}
	fmt.Println("Order Signature:", signature)

	// We extract the exact payload

	payload := ethereal.SignedMessage[*ethereal.Order]{
		Data:      &order,
		Signature: signature,
	}
	payloadJson, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("Signed Order Payload JSON:\n", string(payloadJson))

}
