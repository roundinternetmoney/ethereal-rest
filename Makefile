.PHONY: rest examples test test-integration deps version tag release all fmt tidy vet build

# Default CI / release checks: vet, all packages tests, build.
# Integration (live Testnet): ETHEREAL_INTEGRATION=1 — see `make test-integration` (runs go test ./...).

LATEST_TAG := $(shell git describe --tags --abbrev=0 --match 'v[0-9]*' 2>/dev/null || echo v0.0.0)
VERSION := $(patsubst v%,%,$(LATEST_TAG))
PATCH_VERSION := $(shell echo $(VERSION) | awk -F. '{printf "%d.%d.%d", $$1, $$2, $$3+1}')
NEW_VERSION ?= $(PATCH_VERSION)

fmt:
	go fmt ./...

tidy:
	go mod tidy

vet:
	go vet ./...

build:
	go build ./...

rest: vet test
	go build ./...

# Optional: hits api.etherealtest.net (requires network; may need ETHEREAL_PK). Runs the full tree.
test:
	ETHEREAL_INTEGRATION=1 go test -v ./...

examples:
	mkdir -p bin
	go build -o bin/example_account_balance examples/account_balance/main.go
	go build -o bin/example_limit_multiple examples/limit_multiple/main.go
	go build -o bin/example_limit_single examples/limit_single/main.go
	go build -o bin/example_positions examples/positions/main.go
	go build -o bin/example_cancel_replace examples/cancel_replace/main.go
	go build -o bin/example_twap examples/twap/main.go
	go build -o bin/example_chase examples/chase/main.go

deps:
	go get -u ./...

version:
	@echo "Current version: $(VERSION)"
	@echo "Release version: $(NEW_VERSION)"

tag:
	git tag -a v$(NEW_VERSION) -m "Release v$(NEW_VERSION)"
	git push origin v$(NEW_VERSION)

release: rest tag
	gh release create v$(NEW_VERSION) \
		--title "v$(NEW_VERSION)" \
		--notes "Release v$(NEW_VERSION)"

all: fmt tidy rest
