import { Map, MapOptions } from 'leaflet';

import { Hash } from './Hash';

export class Livemap extends Map {
	private element: string | HTMLElement;
	public hash: Hash | undefined;
	public hasLoaded: boolean = false;

	constructor(element: string | HTMLElement, options?: MapOptions | undefined) {
		super(element, options);
		this.element = element;
        
		this.on('load', () => (this.hasLoaded = true));
	}

	public addHash(): void {
		this.hash = new Hash(this, this.element);
	}

	public removeHash(): void {
		this.hash?.remove();
	}
}
