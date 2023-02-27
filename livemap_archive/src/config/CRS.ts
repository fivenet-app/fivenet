import { extend, CRS, Transformation, Projection } from 'leaflet';

export const centerX = 117.3;
export const centerY = 172.8;
export const scaleX = 0.02072;
export const scaleY = 0.0205;

export const customCRS = extend({}, CRS.Simple, {
	projection: Projection.LonLat,
	scale: function (zoom: number) {
		return Math.pow(2, zoom);
	},
	zoom: function (sc: number) {
		return Math.log(sc) / 0.6931471805599453;
	},
	distance: function (pos1: L.LatLng, pos2: L.LatLng) {
		var x_difference = pos2.lng - pos1.lng;
		var y_difference = pos2.lat - pos1.lat;
		return Math.sqrt(x_difference * x_difference + y_difference * y_difference);
	},
	transformation: new Transformation(scaleX, centerX, -scaleY, centerY),
	infinite: true,
});
