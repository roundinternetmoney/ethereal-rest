# Golang Rest Client for Ethereal API

[Go Reference](https://pkg.go.dev/github.com/roundinternetmoney/ethereal-rest)
[Go Report Card](https://goreportcard.com/report/github.com/roundinternetmoney/ethereal-rest)
[Release](https://github.com/roundinternetmoney/ethereal-rest/actions/workflows/release.yaml)

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
- Run a binary with your key set, for example: `ETHEREAL_PK=<hex> ./bin/example_account_balance`

All examples read **`ETHEREAL_PK`** (hex private key, with or without `0x` prefix). They target **testnet** by default (`rest.Testnet`).

| Example | Binary (`make examples`) | What it shows |
|--------|---------------------------|---------------|
| [examples/account_balance](./examples/account_balance/) | `bin/example_account_balance` | Read-only balances after client init |
| [examples/limit_single](./examples/limit_single/) | `bin/example_limit_single` | Single limit order, cancel, `Send` helper |
| [examples/limit_multiple](./examples/limit_multiple/) | `bin/example_limit_multiple` | `CreateOrders` batch + multi-ID cancel |
| [examples/positions](./examples/positions/) | `bin/example_positions` | Open positions (`GetPosition`) |
| [examples/cancel_replace](./examples/cancel_replace/) | `bin/example_cancel_replace` | Cancel then submit a new order (replace) |
| [examples/twap](./examples/twap/) | `bin/example_twap` | **Composition:** time-sliced orders (not a venue TWAP type) |
| [examples/chase](./examples/chase/) | `bin/example_chase` | **Composition:** cancel/replace loop with a stub reference price |

For more detail, see the [examples/](./examples/) folder and the file comments in `twap` and `chase`.

## Configuration Notes

- If no private key is passed to the rest client, an error will be returned.
- All signable request messages implement the `Signable` interface.
- Only one subaccount is currently supported; by default the first one discovered is used.

## Modifying the package

- This client depends on protobuf wrappers from [pkg.go.dev/roundinternet.money/protos](https://pkg.go.dev/roundinternet.money/protos)
- If you want to extend the `.proto` files directly, see the Buf module at [buf.build/round-internet-money/dex](https://buf.build/round-internet-money/dex)
- Otherwise, use or fork [github.com/roundinternetmoney/protos](github.com/roundinternetmoney/protos)

## Contributing

Contributions are welcome! Please open issues or pull requests as needed.

## Todo

- proto Client -> Server messaging.
- add missing rest methods