.PHONY: clean
clean:
	rm -rf ./dist/ ./public/tiles/*/
	rm -rf gdal2tiles-leaflet

gdal2tiles-leaflet:
	git clone https://github.com/commenthol/gdal2tiles-leaflet.git

.PHONY: gen-tiles
gen-tiles: gdal2tiles-leaflet
	./gdal2tiles-leaflet/gdal2tiles.py -l -p raster -z 0-6 -w none ./maps/GTAV_ATLAS_8192x8192.jpg ./public/tiles/atlas
	./gdal2tiles-leaflet/gdal2tiles.py -l -p raster -z 0-6 -w none ./maps/GTAV_ROAD_8192x8192.jpg ./public/tiles/road
	./gdal2tiles-leaflet/gdal2tiles.py -l -p raster -z 0-6 -w none ./maps/GTAV_SATELITE_8192x8192.jpg ./public/tiles/satelite

.PHONY: optimize-tiles
optimize-tiles:
	find -iname '*.png' -print0 | xargs -n1 -P6 -0 optipng -strip all -clobber -fix -o9

.PHONY: tiles
tiles:
	$(MAKE) generate-tiles
	$(MAKE) optimize-tiles

.PHONY: watch
watch:
	yarn dev
