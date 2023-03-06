.DEFAULT: run-server

.PHONY: clean
clean:
	rm -rf ./dist/ ./public/tiles/*/
	rm -rf gdal2tiles-leaflet

.PHONY: watch
watch:
	yarn dev

.PHONY: build
build:
	go build -o arpanet .

.PHONY: run-server
run-server:
	go run . server

.PHONY: gen
gen: gen-sql gen-proto

IGNORED_TABLES := $(shell paste -s -d, ./query/jet_ignored_tables.txt)

.PHONY: gen-sql
gen-sql:
	jet -source=mysql \
		-dsn="arpanet:changeme@tcp(localhost:3306)/arpanet" \
		-path=./query \
		-ignore-tables "$(IGNORED_TABLES)"

protoc-gen-validate:
	if test ! -d validate/; then \
		git clone https://github.com/bufbuild/protoc-gen-validate.git validate; \
	else \
		git -C validate/ pull --all; \
	fi

.PHONY: gen-proto
gen-proto: protoc-gen-validate
	protoc \
		--proto_path=./validate \
		--proto_path=./proto \
		--go_out=./proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=./proto \
		--go-grpc_opt=paths=source_relative \
		--validate_opt=paths=source_relative \
		--validate_out="lang=go:./proto" \
		$(shell find proto/ -iname "*.proto")

	PATH="$$PATH:node_modules/protoc-gen-js/bin/" \
	protoc \
		--proto_path=./validate \
		--proto_path=./proto \
		--js_out=import_style=commonjs,binary:./gen \
		--grpc-web_out=import_style=typescript,mode=grpcwebtext:./gen \
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
	find -iname '*.png' -print0 | xargs -n1 -P6 -0 optipng -strip all -clobber -fix -o9

.PHONY: tiles
tiles:
	$(MAKE) generate-tiles
	$(MAKE) optimize-tiles
