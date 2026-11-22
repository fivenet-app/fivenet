import type * as L from 'leaflet';
import { DivIcon, LatLng, LayerGroup, Marker, Polyline, Util } from 'leaflet';
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

export class SimpleGraticule extends LayerGroup {
    static defaultOptions: SimpleGraticuleOptions = {
        interval: 20,
        showOriginLabel: false,
        hidden: false,
    };

    override options: SimpleGraticuleOptions = { ...SimpleGraticule.defaultOptions };
    private redrawFrame: number | undefined;
    lineStyle: L.PolylineOptions = {
        stroke: true,
        color: '#000000',
        opacity: 0.6,
        weight: 1,
        interactive: false,
    };

    constructor(options?: SimpleGraticuleOptions) {
        super();
        Util.setOptions(this, options);
    }

    override onAdd(map: L.Map): this {
        super.onAdd(map);
        map.on(this.getRedrawEvents(), this.scheduleRedraw, this);
        this.redraw();
        return this;
    }

    override onRemove(map: L.Map): this {
        map.off(this.getRedrawEvents(), this.scheduleRedraw, this);
        this.cancelScheduledRedraw();
        super.onRemove(map);
        this.clearLayers();
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
        this.cancelScheduledRedraw();

        if (!this._map) return this;

        const bounds = this._map.getBounds().pad(0.5);

        this.clearLayers();

        if (this.options.hidden) return this;

        const interval = this.getInterval();
        this.constructLines(bounds, interval);

        if (this.options.showOriginLabel && bounds.contains([0, 0])) {
            this.addLayer(this.addOriginLabel());
        }

        return this;
    }

    private getRedrawEvents(): string {
        return ['viewreset', this.options.redraw].filter((event): event is string => !!event).join(' ');
    }

    private scheduleRedraw(): void {
        if (this.redrawFrame !== undefined) return;

        this.redrawFrame = requestAnimationFrame(() => {
            this.redrawFrame = undefined;
            this.redraw();
        });
    }

    private cancelScheduledRedraw(): void {
        if (this.redrawFrame === undefined) return;

        cancelAnimationFrame(this.redrawFrame);
        this.redrawFrame = undefined;
    }

    private getInterval(): number {
        const currentZoom = this._map.getZoom();
        const interval =
            this.options.zoomIntervals?.find(({ start, end }) => currentZoom >= start && currentZoom <= end)?.interval ??
            this.options.interval ??
            SimpleGraticule.defaultOptions.interval!;

        if (!Number.isFinite(interval) || interval <= 0) {
            throw new RangeError('SimpleGraticule interval must be a finite number greater than zero.');
        }

        return interval;
    }

    private getLineCounts(bounds: L.LatLngBounds, interval: number): { x: number; y: number } {
        return {
            x: Math.ceil((bounds.getEast() - bounds.getWest()) / interval),
            y: Math.ceil((bounds.getNorth() - bounds.getSouth()) / interval),
        };
    }

    private getMins(bounds: L.LatLngBounds, interval: number): { x: number; y: number } {
        return {
            x: Math.floor(bounds.getWest() / interval) * interval,
            y: Math.floor(bounds.getSouth() / interval) * interval,
        };
    }

    private constructLines(bounds: L.LatLngBounds, interval: number): void {
        const mins = this.getMins(bounds, interval);
        const counts = this.getLineCounts(bounds, interval);
        const labelBounds = this._map.getBounds().pad(-0.003);

        // Horizontal lines
        for (let i = 0; i <= counts.x; i++) {
            const x = mins.x + i * interval;
            this.addLayer(this.buildXLine(bounds, x));
            this.addLayer(this.buildLabel(labelBounds, 'gridlabel-horiz', x));
        }

        // Vertical lines
        for (let j = 0; j <= counts.y; j++) {
            const y = mins.y + j * interval;
            this.addLayer(this.buildYLine(bounds, y));
            this.addLayer(this.buildLabel(labelBounds, 'gridlabel-vert', y));
        }
    }

    private buildXLine(bounds: L.LatLngBounds, x: number): L.Polyline {
        const bottomLL = new LatLng(bounds.getSouth(), x);
        const topLL = new LatLng(bounds.getNorth(), x);
        return new Polyline([bottomLL, topLL], this.lineStyle);
    }

    private buildYLine(bounds: L.LatLngBounds, y: number): L.Polyline {
        const leftLL = new LatLng(y, bounds.getWest());
        const rightLL = new LatLng(y, bounds.getEast());
        return new Polyline([leftLL, rightLL], this.lineStyle);
    }

    private buildLabel(bounds: L.LatLngBounds, axis: 'gridlabel-horiz' | 'gridlabel-vert', val: number): Marker {
        let latLng: LatLng;
        if (axis === 'gridlabel-horiz') {
            latLng = new LatLng(bounds.getNorth(), val);
        } else {
            latLng = new LatLng(val, bounds.getWest());
        }

        return new Marker(latLng, {
            interactive: false,
            icon: new DivIcon({
                iconSize: [0, 0],
                className: 'leaflet-grid-label',
                html: `<div class="${axis}">${val}</div>`,
            }),
        });
    }

    private addOriginLabel(): Marker {
        return new Marker([0, 0], {
            interactive: false,
            icon: new DivIcon({
                iconSize: [0, 0],
                className: 'leaflet-grid-label',
                html: '<div class="gridlabel-horiz">(0,0)</div>',
            }),
        });
    }
}
