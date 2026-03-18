.PHONY: rest examples test deps version tag release all

LATEST_TAG := $(shell git describe --tags --abbrev=0 --match 'v[0-9]*' 2>/dev/null || echo v0.0.0)
VERSION := $(patsubst v%,%,$(LATEST_TAG))
PATCH_VERSION := $(shell echo $(VERSION) | awk -F. '{printf "%d.%d.%d", $$1, $$2, $$3+1}')
NEW_VERSION ?= $(PATCH_VERSION)

rest:
	go vet ./...
	go test ./...
	go build ./...

examples: 
	mkdir -p bin
	go build -o bin/example_account_balance examples/account_balance/main.go
	go build -o bin/example_limit_multiple examples/limit_multiple/main.go
	go build -o bin/example_limit_single examples/limit_single/main.go

test:
	go test -v tests/order_signing/signing_test.go

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

all: fmt tidy vet test build
