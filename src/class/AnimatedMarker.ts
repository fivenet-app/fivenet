import L from 'leaflet';

export interface AnimatedMarkerOptions extends L.MarkerOptions {
	distance?: number;
	interval?: number;
}

export class AnimatedMarker extends L.Marker {
	public distance: number;
	public interval: number;

	private _i: number = 0;
	private _latlngs: Array<L.LatLng> = [];

	constructor(latlng: L.LatLng, options: AnimatedMarkerOptions | undefined = undefined) {
		super(latlng, options);

		this.distance = options?.distance ? options.distance : 200;
		this.interval = options?.interval ? options.interval : 1000;

		return this;
	}

	private _chunk(latlngs: Array<L.LatLng>): Array<L.LatLng> {
		const len = latlngs.length;
		const chunkedLatLngs: Array<L.LatLng> = [];

		for (let i = 1; i < latlngs.length; i++) {
			let curr = latlngs[i - 1];
			let next = latlngs[i];
			let dist = curr.distanceTo(next);
			const factor = this.distance / dist;
			const dLat = factor * (next.lat - curr.lat);
			const dLng = factor * (next.lng - curr.lng);

			if (dist > this.distance) {
				while (dist > this.distance) {
					curr = new L.LatLng(curr.lat + dLat, curr.lng + dLng);
					dist = curr.distanceTo(next);
					chunkedLatLngs.push(curr);
				}
			} else {
				chunkedLatLngs.push(curr);
			}
		}

		chunkedLatLngs.push(latlngs[len - 1]);
		return chunkedLatLngs;
	}

	public setLine(latlng: L.LatLng): void {
		this._latlngs = this._chunk([latlng]);
		this._i = 0;
	}

	public animate(): void {
		const len = this._latlngs.length;
		let speed = this.interval;

		if (this._i < len && this._i > 0) {
			speed = (this._latlngs[this._i - 1].distanceTo(this._latlngs[this._i]) / this.distance) * this.interval;
		}

		if (L.DomUtil.TRANSITION) {
			if (this._shadow) {
				this._shadow.style.transition = 'all ' + speed + 'ms ease';
			}

			if (this.getIcon()) {
				let style = document.getElementById('AnimatedMarker');
				if (!style) {
					const head = document.getElementsByTagName('head')[0];
					head.innerHTML = head.innerHTML + '<style id="AnimatedMarker"></style>';
					style = document.getElementById('AnimatedMarker');
				}

				if (style) style.innerHTML = '.leaflet-marker-icon { transition: all ' + speed + 'ms ease }';
			}
		}

		this.setLatLng(this._latlngs[this._i]);
		this._i++;

		setTimeout(() => {
			if (this._i !== len) this.animate();
		}, speed);
	}

	moveTo(latlng: L.LatLng): this {
		this.setLine(latlng);
		this.animate();
		return this;
	}
}
