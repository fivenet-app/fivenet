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
	$(GO) install github.com/bufbuild/buf/cmd/buf@v1.26.1
endif

protoc-gen-go:
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc-gen-go-grpc:
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc-gen-validate: build_dir
	if test ! -d $(BUILD_DIR)validate-$(VALIDATE_VERSION)/; then \
		git clone --branch $(VALIDATE_VERSION) https://github.com/bufbuild/protoc-gen-validate.git $(BUILD_DIR)validate-$(VALIDATE_VERSION); \
	else \
		git -C $(BUILD_DIR)validate-$(VALIDATE_VERSION)/ pull --all; \
		git -C $(BUILD_DIR)validate-$(VALIDATE_VERSION)/ checkout $(VALIDATE_VERSION); \
	fi

	cd $(BUILD_DIR) && ln -sfn validate-$(VALIDATE_VERSION)/ validate

protoc-gen-customizer:
	$(GO) build -o ./internal/cmd/protoc-gen-customizer ./internal/cmd/protoc-gen-customizer

protoc-gen-fronthelper:
	$(GO) build -o ./internal/cmd/protoc-gen-fronthelper ./internal/cmd/protoc-gen-fronthelper

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
		-tags=jsoniter \
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
gen-proto: protoc-gen-go protoc-gen-go-grpc protoc-gen-validate protoc-gen-customizer protoc-gen-fronthelper protoc-gen-doc
	mkdir -p ./gen/go/proto
	PATH="$$PATH:./internal/cmd/protoc-gen-customizer/" \
	$(PROTOC) \
		--proto_path=./$(BUILD_DIR)validate-$(VALIDATE_VERSION) \
		--proto_path=./proto \
		--go_out=./gen/go/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=./gen/go/proto \
		--go-grpc_opt=paths=source_relative \
		--validate_opt=paths=source_relative \
		--validate_out="lang=go:./gen/go/proto" \
		--customizer_opt=paths=source_relative \
		--customizer_out=./gen/go/proto \
		--doc_opt=./internal/protoc-gen-doc-md.tmpl,grpc-api.md \
		--doc_out=./gen \
		$(shell find proto/ -iname "*.proto")

	# Inject Go field tags into generated fields
	find ./gen/go/proto/ -iname "*.pb.go" \
		-exec protoc-go-inject-tag \
			-input={} \;

	mkdir -p ./gen/ts
	PATH="$$PATH:node_modules/@protobuf-ts/plugin/bin/:./internal/cmd/protoc-gen-fronthelper/" \
	$(PROTOC) \
		--proto_path=./$(BUILD_DIR)validate-$(VALIDATE_VERSION) \
		--proto_path=./proto \
		--ts_out=./gen/ts \
		--ts_opt=optimize_speed,long_type_number,force_server_none \
		--fronthelper_opt=paths=source_relative \
		--fronthelper_out=./gen/ts \
		$(shell find proto/ -iname "*.proto")

	# Fix ignore TS typecheck comment
	find ./gen/ts/ -type f -iname "*.ts" -print0 | \
		xargs -0 sed -i 's~// tslint:disable~// @ts-nocheck~g'

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
