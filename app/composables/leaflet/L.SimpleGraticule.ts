import L from 'leaflet';
/**
 * File: L.SimpleGraticule.js
 * Desc: A graticule for Leaflet maps in the L.CRS.Simple coordinate system.
 * Auth: Andrew Blakey (ablakey@gmail.com)
 * License: BSD 2-Clause "Simplified" License - https://github.com/ablakey/Leaflet.SimpleGraticule/blob/master/LICENSE
 *
 * Code has been modified to use TypeScript and to fit into the current project structure.
 */
interface ZoomInterval {
    start: number;
    end: number;
    interval: number;
}

interface SimpleGraticuleOptions extends L.LayerOptions {
    interval?: number;
    showOriginLabel?: boolean;
    redraw?: string;
    hidden?: boolean;
    zoomIntervals?: ZoomInterval[];
}

export class SimpleGraticule extends L.LayerGroup {
    declare _map: L.Map;
    declare _bounds: L.LatLngBounds;
    override options: SimpleGraticuleOptions = {};
    lineStyle: L.PolylineOptions = {
        stroke: true,
        color: '#000000',
        opacity: 0.6,
        weight: 1,
        interactive: false,
    };

    constructor(options?: SimpleGraticuleOptions) {
        super();
        L.Util.setOptions(this, options);
    }

    override onAdd(map: L.Map): this {
        this._map = map;
        this.redraw();
        this._map.on('viewreset ' + this.options.redraw, this.redraw, this);
        this.eachLayer(map.addLayer, map);
        return this;
    }

    override onRemove(map: L.Map): this {
        map.off('viewreset ' + this.options.redraw, this.redraw, this);
        this.eachLayer(this.removeLayer, this);
        return this;
    }

    hide(): void {
        this.options.hidden = true;
        this.redraw();
    }

    show(): void {
        this.options.hidden = false;
        this.redraw();
    }

    redraw(): this {
        this._bounds = this._map.getBounds().pad(0.5);

        this.clearLayers();

        if (!this.options.hidden) {
            const currentZoom = this._map.getZoom();

            if (this.options.zoomIntervals) {
                for (let i = 0; i < this.options.zoomIntervals.length; i++) {
                    const zi = this.options.zoomIntervals[i]!;
                    if (currentZoom >= zi.start && currentZoom <= zi.end) {
                        this.options.interval = zi.interval;
                        break;
                    }
                }
            }

            this.constructLines(this.getMins(), this.getLineCounts());

            if (this.options.showOriginLabel) {
                this.addLayer(this.addOriginLabel());
            }
        }

        return this;
    }

    getLineCounts(): { x: number; y: number } {
        const interval = this.options.interval ?? 20;
        return {
            x: Math.ceil((this._bounds.getEast() - this._bounds.getWest()) / interval),
            y: Math.ceil((this._bounds.getNorth() - this._bounds.getSouth()) / interval),
        };
    }

    getMins(): { x: number; y: number } {
        const s = this.options.interval ?? 20;
        return {
            x: Math.floor(this._bounds.getWest() / s) * s,
            y: Math.floor(this._bounds.getSouth() / s) * s,
        };
    }

    constructLines(mins: { x: number; y: number }, counts: { x: number; y: number }): void {
        const lines: L.Polyline[] = [];
        const labels: L.Marker[] = [];

        // Horizontal lines
        for (let i = 0; i <= counts.x; i++) {
            const x = mins.x + i * (this.options.interval ?? 20);
            lines.push(this.buildXLine(x));
            labels.push(this.buildLabel('gridlabel-horiz', x));
        }

        // Vertical lines
        for (let j = 0; j <= counts.y; j++) {
            const y = mins.y + j * (this.options.interval ?? 20);
            lines.push(this.buildYLine(y));
            labels.push(this.buildLabel('gridlabel-vert', y));
        }

        lines.forEach((l) => this.addLayer(l));
        labels.forEach((l) => this.addLayer(l));
    }

    buildXLine(x: number): L.Polyline {
        const bottomLL = new L.LatLng(this._bounds.getSouth(), x);
        const topLL = new L.LatLng(this._bounds.getNorth(), x);
        return new L.Polyline([bottomLL, topLL], this.lineStyle);
    }

    buildYLine(y: number): L.Polyline {
        const leftLL = new L.LatLng(y, this._bounds.getWest());
        const rightLL = new L.LatLng(y, this._bounds.getEast());
        return new L.Polyline([leftLL, rightLL], this.lineStyle);
    }

    buildLabel(axis: 'gridlabel-horiz' | 'gridlabel-vert', val: number): L.Marker {
        const bounds = this._map.getBounds().pad(-0.003);
        let latLng: L.LatLng;
        if (axis === 'gridlabel-horiz') {
            latLng = new L.LatLng(bounds.getNorth(), val);
        } else {
            latLng = new L.LatLng(val, bounds.getWest());
        }

        return L.marker(latLng, {
            interactive: false,
            contextmenu: false,
            contextmenuItems: [],
            icon: L.divIcon({
                iconSize: [0, 0],
                className: 'leaflet-grid-label',
                html: `<div class="${axis}">${val}</div>`,
            }),
        });
    }

    addOriginLabel(): L.Marker {
        return L.marker([0, 0], {
            interactive: false,
            contextmenu: false,
            contextmenuItems: [],
            icon: L.divIcon({
                iconSize: [0, 0],
                className: 'leaflet-grid-label',
                html: '<div class="gridlabel-horiz">(0,0)</div>',
            }),
        });
    }
}

// Factory function for compatibility
export function simpleGraticule(options?: SimpleGraticuleOptions): SimpleGraticule {
    return new SimpleGraticule(options);
}
