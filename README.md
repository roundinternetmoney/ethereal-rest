# Ethereal Go Client

Lightweight golang client for interacting with the Ethereal API.

## Features

- Experimental protobuf support.
- Order placement and cancellation for REST, Websocket, and Socket.IO
- EIP-712 data signing
- Batch execution support (concurrent, unordered, type-safe)
- Automatic nonce and timestamp handling
- Minimal dependencies

## Getting started

- Requires Go 1.25+.
- Install from GitHub: `go get github.com/roundinternetmoney/ethereal-wss`

## Example Usage

From the client directory:

- `make examples`
- `ETHEREAL_PK=0x0000 bin/example_account_balance`

For more complete usage examples (batching, cancel orders, etc.),
see the [examples/](./examples/) folder in this repository.

## Configuration Notes

- If no private key is passed to the rest client, an error will be returned.
- All signable request messages implement the `Signable` interface.
- Only one subaccount is currently supported; by default the first one discovered is used.

## Modifying the package
- This client depends on protobuf wrappers from [pkg.go.dev/roundinternet.money/pb-dex](https://pkg.go.dev/roundinternet.money/pb-dex)
- If you want to extend the `.proto` files directly, see the Buf module at [buf.build/round-internet-money/dex](https://buf.build/round-internet-money/dex)
- Otherwise, use or fork [github.com/roundinternetmoney/pb-dex](https://github.com/Round-Internet-Money/pb-dex)

Contributing
-------------
Contributions are welcome! Please open issues or pull requests as needed.

## Todo

- proto Client -> Server messaging.
- add missing rest methods
