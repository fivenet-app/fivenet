VERSION := $(shell cat VERSION)

VALIDATE_VERSION ?= v1.0.2
BUILD_DIR := .build/

GO ?= go
PROTOC ?= protoc

.DEFAULT: run-server

# Read .env file and export the vars
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build_dir:
	mkdir -p $(BUILD_DIR)

# Build, Format, etc., Tools, Dependency checkouts

buf:
ifeq (, $(shell which buf))
	$(GO) install github.com/bufbuild/buf/cmd/buf@v1.55.1
endif

protoc-gen-go:
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest

protoc-gen-go-grpc:
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc-gen-doc:
	$(GO) install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest

go-licenses:
ifeq (, $(shell which go-licenses))
	@GO111MODULE=on $(GO) install github.com/google/go-licenses@latest
endif

# Actual targets

.PHONY: clean
clean:
	@pnpx nuxi cleanup
	rm -rf ./.nuxt/dist/

.PHONY: watch
watch:
	pnpm dev

.PHONY: build-container
build-container:
	docker build -t docker.io/fivenet-app/fivenet:latest \
		--build-arg NUXT_UI_PRO_LICENSE=$(NUXT_UI_PRO_LICENSE) \
		.

.PHONY: release
release:
	docker tag docker.io/fivenet-app/fivenet:latest docker.io/fivenet-app/fivenet:$(VERSION)

.PHONY: tests
tests: tests-go

.PHONY: tests-go
tests-go:
	$(GO) test -v ./...

.PHONY: build-go
build-go:
	CGO_ENABLED=0 $(GO) \
		build \
		-a \
		-installsuffix cgo \
		-ldflags "-X github.com/fivenet-app/fivenet/v2025/pkg/version.Version=$(shell git describe --tags --exclude='fivenet-*')" \
		-o fivenet \
		.

.PHONY: build-js
build-js:
	rm -rf ./.nuxt/dist/
	NODE_OPTIONS="--max-old-space-size=8192" pnpm build

.PHONY: run-server
run-server:
	mkdir -p ./.nuxt/dist/
	$(GO) run . server

.PHONY: gen
gen: gen-sql gen-proto

.PHONY: gen-sql
gen-sql:
	$(GO) run ./query/gen/

	# Remove schema/database name from the generated table code, so it uses the currently selected database
	find ./query/fivenet/table -type f -iname '*.go' -exec sed -i 's~("fivenet", ~("", ~g' {} \;

.PHONY: gen-proto
gen-proto: buf protoc-gen-go protoc-gen-go-grpc protoc-gen-doc
	mkdir -p ./gen/go/proto ./gen/ts

	buf generate
	buf generate --template buf.gen.tag.yaml

.PHONY: fmt
fmt:
	$(MAKE) fmt-proto gen-proto
	$(MAKE) fmt-js

.PHONY: fmt-proto
fmt-proto: buf
	buf format --write ./proto

.PHONY: fmt-js
fmt-js:
	pnpm prettier --write ./app

.PHONY: gen-licenses
gen-licenses: go-licenses
	$(MAKE) gen-licenses-js gen-licenses-go

gen-licenses-js:
	pnpm licenses list --prod --json | \
		pnpx @quantco/pnpm-licenses generate-disclaimer --json-input \
			--output-file=./public/licenses/frontend.txt

gen-licenses-go:
	go-licenses report ./... --ignore $$($(GO) list -m) --include_tests \
		--ignore $$($(GO) list std | awk 'NR > 1 { printf(",") } { printf("%s",$$0) } END { print "" }') \
		--template internal/scripts/go-licenses-backend.txt.tpl > ./public/licenses/backend.txt
