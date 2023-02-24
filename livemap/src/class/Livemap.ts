import { LatLng, Map, MapOptions } from 'leaflet';

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
		console.log(layer);
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
