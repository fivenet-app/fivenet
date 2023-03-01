import { CRS, extend, Map, MapOptions, Projection, Transformation } from 'leaflet';

import { Hash } from './Hash';

export class Livemap extends Map {
	public hash: Hash | undefined;
	public hasLoaded: boolean = false;

	private element: HTMLElement;

	constructor(element: string | HTMLElement, options?: MapOptions | undefined) {
		super(element, options);
		this.element = typeof element === 'string' ? (document.getElementById(element) as HTMLElement) : element;

		this.on('load', () => (this.hasLoaded = true));
		this.on('baselayerchange', (context) => this.updateBackground(context.name));
	}

	public addHash(): void {
		this.hash = new Hash(this, this.element);
	}

	public removeHash(): void {
		this.hash?.remove();
	}

	public updateBackground(layer: string): void {
		switch (layer) {
			case 'Atlas':
				this.element.style.backgroundColor = '#0fa8d2';
				return;
			case 'Satelite':
				this.element.style.backgroundColor = '#143d6b';
				return;
			case 'Road':
				this.element.style.backgroundColor = '#1862ad';
				return;
		}
	}
}

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
