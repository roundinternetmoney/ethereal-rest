# Golang Rest Client for Ethereal API

[![Go Reference](https://pkg.go.dev/badge/github.com/roundinternetmoney/ethereal-rest.svg)](https://pkg.go.dev/github.com/roundinternetmoney/ethereal-rest)
[![Go Report Card](https://goreportcard.com/badge/github.com/roundinternetmoney/ethereal-rest)](https://goreportcard.com/report/github.com/roundinternetmoney/ethereal-rest)
[![Release](https://github.com/roundinternetmoney/ethereal-rest/actions/workflows/Release.yml/badge.svg)](https://github.com/roundinternetmoney/ethereal-rest/actions/workflows/Release.yml)

Lightweight golang client for interacting with the Ethereal API.

## Features

- Server -> Client Protobuf support.
- Batch execution support (concurrent, ordered, type-safe)
- EIP-712 data signing
- Automatic nonce, timestamp, and subaccount handling for requests
- Minimal dependencies

## Getting started

- Requires Go 1.25+.
- Install from GitHub: `go get github.com/roundinternetmoney/ethereal-rest`
- Import path: `github.com/roundinternetmoney/ethereal-rest`

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
- This client depends on protobuf wrappers from [pkg.go.dev/roundinternet.money/protos](https://pkg.go.dev/roundinternet.money/protos)
- If you want to extend the `.proto` files directly, see the Buf module at [buf.build/round-internet-money/dex](https://buf.build/round-internet-money/dex)
- Otherwise, use or fork [github.com/roundinternetmoney/protos](github.com/roundinternetmoney/protos)

Contributing
-------------
Contributions are welcome! Please open issues or pull requests as needed.

## Todo

- proto Client -> Server messaging.
- add missing rest methods
