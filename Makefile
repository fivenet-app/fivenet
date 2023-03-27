VERSION := $(shell cat VERSION)

VALIDATE_VERSION ?= v0.10.1

.DEFAULT: run-server

.PHONY: clean
clean:
	rm -rf ./dist/ ./public/tiles/*/
	rm -rf gdal2tiles-leaflet

.PHONY: watch
watch:
	yarn dev

.PHONY: build-container
build-container:
	docker build -t galexrt/arpanet:latest .

.PHONY: release
release:
	docker tag galexrt/arpanet:latest galexrt/arpanet:$(VERSION)

.PHONY: build-go
build-go:
	go build -o arpanet .

.PHONY: build-yarn
build-yarn:
	rm -rf ./dist/
	VITE_BASE="/dist" \
		yarn build

.PHONY: run-server
run-server:
	mkdir -p ./dist/
	go run . server

.PHONY: gen
gen: gen-sql gen-proto

.PHONY: gen-sql
gen-sql:
	go run ./query/gen/

	# Remove schema/database name from the generated table code, so it uses the currently selected database
	find ./query/arpanet/table -type f -iname '*.go' -exec sed -i 's~("arpanet", ~("", ~g' {} \;

protoc-gen-validate:
	if test ! -d validate/; then \
		git clone --branch $(VALIDATE_VERSION) https://github.com/bufbuild/protoc-gen-validate.git validate; \
	else \
		git -C validate/ pull --all; \
		git -C validate/ checkout $(VALIDATE_VERSION); \
	fi

protoc-gen-customizer:
	go build -o ./cmd/protoc-gen-customizer ./cmd/protoc-gen-customizer

protoc-gen-customizerweb:
	go build -o ./cmd/protoc-gen-customizerweb ./cmd/protoc-gen-customizerweb

.PHONY: gen-proto
gen-proto: protoc-gen-validate protoc-gen-customizer protoc-gen-customizerweb
	PATH="$$PATH:cmd/protoc-gen-customizer/" \
	protoc \
		--proto_path=./validate \
		--proto_path=./proto \
		--go_out=./proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=./proto \
		--go-grpc_opt=paths=source_relative \
		--validate_opt=paths=source_relative \
		--validate_out="lang=go:./proto" \
		--customizer_opt=paths=source_relative \
		--customizer_out="./proto" \
		$(shell find proto/ -iname "*.proto")

	find proto/ -iname "*.pb.go" \
		-exec protoc-go-inject-tag \
			-input={} \;

	PATH="$$PATH:cmd/protoc-gen-customizerweb/:node_modules/protoc-gen-js/bin/" \
	protoc \
		--proto_path=./validate \
		--proto_path=./proto \
		--js_out=import_style=commonjs,binary:./gen \
		--grpc-web_out=import_style=typescript,mode=grpcwebtext:./gen \
		--customizerweb_opt=paths=source_relative \
		--customizerweb_out="./gen" \
		$(shell find proto/ -iname "*.proto")

	# Remove validate_pb imports from JS files
	find ./gen -type f \( -iname '*.js' -o -iname '*.ts' \) -exec sed -i '/validate_pb/d' {} +

	# Update local yarn package
	yarn upgrade '@arpanet/gen@file:./gen'

gdal2tiles-leaflet:
	if test ! -d gdal2tiles-leaflet/; then \
		git clone https://github.com/commenthol/gdal2tiles-leaflet.git gdal2tiles-leaflet; \
	else \
		git -C gdal2tiles-leaflet pull --all; \
	fi

.PHONY: gen-tiles
gen-tiles: gdal2tiles-leaflet
	./gdal2tiles-leaflet/gdal2tiles.py -l -p raster -z 0-6 -w none ./maps/GTAV_ATLAS_8192x8192.jpg ./public/tiles/atlas
	./gdal2tiles-leaflet/gdal2tiles.py -l -p raster -z 0-6 -w none ./maps/GTAV_POSTAL_8192x8192.jpg ./public/tiles/postal
	./gdal2tiles-leaflet/gdal2tiles.py -l -p raster -z 0-6 -w none ./maps/GTAV_ROAD_8192x8192.jpg ./public/tiles/road
	./gdal2tiles-leaflet/gdal2tiles.py -l -p raster -z 0-6 -w none ./maps/GTAV_SATELITE_8192x8192.jpg ./public/tiles/satelite

.PHONY: optimize-tiles
optimize-tiles:
	find ./public/tiles -iname '*.png' -print0 | xargs -n1 -P6 -0 optipng -strip all -clobber -fix -o9

.PHONY: tiles
tiles:
	$(MAKE) gen-tiles
	$(MAKE) optimize-tiles
