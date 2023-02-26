.DEFAULT: run-server

.PHONY: clean
clean:
	rm -rf ./dist/

.PHONY: build
build:
	go build -o arpanet .

.PHONY: run-server
run-server:
	go run . server

.PHONY: gen-gorm
gen-gorm:
	go run ./gorm/main.go

.PHONY: gen-proto
gen-proto:
	protoc \
		--proto_path=./proto \
		--go_out=./gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=./gen \
		--go-grpc_opt=paths=source_relative \
		--validate_opt=paths=source_relative \
		--validate_out="lang=go:./gen" \
		$(shell find proto/ -iname "*.proto")

	PATH="$$PATH:node_modules/protoc-gen-js/bin/" protoc \
		--proto_path=./proto \
		--js_out=import_style=commonjs,binary:./gen \
		--grpc-web_out=import_style=typescript,mode=grpcwebtext:./gen \
		$(shell find proto/ -iname "*.proto")

localhost-certs:
	openssl req -x509 -out localhost.crt -keyout localhost.key \
		-newkey rsa:2048 -nodes -sha256 \
		-subj '/CN=localhost' -extensions EXT -config <( \
		printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
