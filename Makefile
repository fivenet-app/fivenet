VERSION := $(shell cat VERSION)

VALIDATE_VERSION ?= v1.0.2
BUILD_DIR := .build/

.DEFAULT: run-server

build_dir:
	mkdir -p $(BUILD_DIR)

.PHONY: clean
clean:
	rm -rf ./.nuxt/dist/ ./src/public/images/livemap/tiles/*/
	rm -rf gdal2tiles-leaflet

.PHONY: watch
watch:
	yarn dev

.PHONY: build-container
build-container:
	docker build -t galexrt/fivenet:latest .

.PHONY: release
release:
	docker tag galexrt/fivenet:latest galexrt/fivenet:$(VERSION)

.PHONY: build-go
build-go:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o fivenet .

.PHONY: build-yarn
build-yarn:
	rm -rf ./.nuxt/dist/
	yarn build

.PHONY: run-server
run-server:
	mkdir -p ./.nuxt/dist/
	go run . server

.PHONY: gen
gen: gen-sql gen-proto

.PHONY: gen-sql
gen-sql:
	go run ./query/gen/

	# Remove schema/database name from the generated table code, so it uses the currently selected database
	find ./query/fivenet/table -type f -iname '*.go' -exec sed -i 's~("fivenet", ~("", ~g' {} \;

protoc-gen-validate: build_dir
	if test ! -d $(BUILD_DIR)validate-$(VALIDATE_VERSION)/; then \
		git clone --branch $(VALIDATE_VERSION) https://github.com/bufbuild/protoc-gen-validate.git $(BUILD_DIR)validate-$(VALIDATE_VERSION); \
	else \
		git -C $(BUILD_DIR)validate-$(VALIDATE_VERSION)/ pull --all; \
		git -C $(BUILD_DIR)validate-$(VALIDATE_VERSION)/ checkout $(VALIDATE_VERSION); \
	fi

	cd $(BUILD_DIR) && ln -sfn validate-$(VALIDATE_VERSION)/ validate

protoc-gen-customizer:
	go build -o ./cmd/protoc-gen-customizer ./cmd/protoc-gen-customizer

.PHONY: gen-proto
gen-proto: protoc-gen-validate protoc-gen-customizer
	PATH="$$PATH:./cmd/protoc-gen-customizer/" \
	npx protoc \
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
		$(shell find proto/ -iname "*.proto")

	# Inject Go field tags into generated fields
	find ./gen/go/proto/ -iname "*.pb.go" \
		-exec protoc-go-inject-tag \
			-input={} \;

	npx protoc \
		--proto_path=./$(BUILD_DIR)validate-$(VALIDATE_VERSION) \
		--proto_path=./proto \
		--ts_out=./gen/ts \
		--ts_opt=optimize_code_size,long_type_bigint \
		$(shell find proto/ -iname "*.proto")

	node ./internal/scripts/proto-patch.js

	# Remove validate_pb imports from JS files
	find ./gen -type f \( -iname '*.js' -o -iname '*.ts' \) -exec sed -i '/validate_pb/d' {} +

.PHONY: gen-licenses
gen-licenses:
	yarn licenses generate-disclaimer > ./src/public/licenses/frontend.txt
	go-licenses report . --template internal/scripts/licenses-backend.txt.tpl > ./src/public/licenses/backend.txt

gdal2tiles-leaflet: build_dir
	if test ! -d $(BUILD_DIR)gdal2tiles-leaflet/; then \
		git clone https://github.com/commenthol/gdal2tiles-leaflet.git $(BUILD_DIR)gdal2tiles-leaflet; \
	else \
		git -C $(BUILD_DIR)gdal2tiles-leaflet pull --all; \
	fi

.PHONY: gen-tiles
gen-tiles: gdal2tiles-leaflet
	$(BUILD_DIR)gdal2tiles-leaflet/gdal2tiles.py -l -p raster -z 0-7 -w none ./maps/GTAV_ATLAS_8192x8192.jpg ./src/public/images/livemap/tiles/atlas
	$(BUILD_DIR)gdal2tiles-leaflet/gdal2tiles.py -l -p raster -z 0-7 -w none ./maps/GTAV_POSTAL_8192x8192.jpg ./src/public/images/livemap/tiles/postal
	$(BUILD_DIR)gdal2tiles-leaflet/gdal2tiles.py -l -p raster -z 0-7 -w none ./maps/GTAV_ROAD_8192x8192.jpg ./src/public/images/livemap/tiles/road
	$(BUILD_DIR)gdal2tiles-leaflet/gdal2tiles.py -l -p raster -z 0-7 -w none ./maps/GTAV_SATELITE_8192x8192.jpg ./src/public/images/livemap/tiles/satelite

.PHONY: optimize-tiles
optimize-tiles:
	find ./src/public/tiles -iname '*.png' -print0 | xargs -n1 -P6 -0 optipng -strip all -clobber -fix -o9

.PHONY: tiles
tiles:
	$(MAKE) gen-tiles
	$(MAKE) optimize-tiles

# ====================================================================================
# Makefile helper functions for helm-docs: https://github.com/norwoodj/helm-docs
#

HELM_DOCS_VERSION := v1.11.0
HELM_DOCS := helm-docs
HELM_DOCS_REPO := github.com/norwoodj/helm-docs/cmd/helm-docs

bin-$(HELM_DOCS): ## Installs helm-docs
	@GO111MODULE=on go install $(HELM_DOCS_REPO)@$(HELM_DOCS_VERSION)

helm-docs: bin-$(HELM_DOCS) ## Use helm-docs to generate documentation from helm charts
	$(HELM_DOCS) -c charts/fivenet \
		-o README.md \
		-t README.gotmpl.md \
		-t _templates.gotmpl
