import { DomEvent, LatLng, Util } from 'leaflet';

import { Livemap } from './Livemap';

export interface parsedHash {
    center: LatLng;
    zoom: number;
}

export class Hash {
    private HAS_HASHCHANGE: boolean = 'onhashchange' in window;

    private map: Livemap | undefined;
    private app: HTMLElement | null;

    private lastHash: string = '';
    private changeDefer: number = 100;
    private changeTimeout: NodeJS.Timeout | undefined = undefined;
    private movingMap: boolean = false;
    private isListening: boolean = false;
    private hashChangeInterval: NodeJS.Timer | undefined = undefined;

    constructor(map: Livemap, element: HTMLElement) {
        this.onHashChange = Util.bind(this.onHashChange, this);
        this.app = element;
        this.map = map;

        this.onHashChange();
        if (!this.isListening) this.startListening();
    }

    public remove(): void {
        if (this.changeTimeout) clearTimeout(this.changeTimeout);
        if (this.isListening) this.stopListening();

        this.map = undefined;
    }

    private parseHash(hash: string): parsedHash | false {
        if (hash.indexOf('#') === 0)  hash = hash.substring(1);

        const args = hash.split('/');
        if (args.length == 3) {
            const zoom = parseInt(args[0], 10),
                lat = parseFloat(args[1]),
                lon = parseFloat(args[2]);
            if (isNaN(zoom) || isNaN(lat) || isNaN(lon)) {
                return false;
            } else {
                return {
                    center: new LatLng(lat, lon),
                    zoom: zoom,
                };
            }
        } else {
            return false;
        }
    }

    private formatHash(): string {
        if (!this.map) throw new Error('map is undefined');

        const center = this.map.getCenter();
        const zoom = this.map.getZoom();
        const precision = Math.max(0, Math.ceil(Math.log(zoom) / Math.LN2));

        const hash = '#' + [zoom, center.lat.toFixed(precision), center.lng.toFixed(precision)].join('/');
        return hash;
    }

    private onMapMove(): void {
        if (!this.map) return;

        if (this.movingMap || !this.map.hasLoaded) return;

        const hash = this.formatHash();
        if (this.lastHash != hash) {
            window.location.replace(hash);
            this.lastHash = hash;
        }
    }

    private update(): void {
        if (!this.map) return;

        const hash = window.location.hash;
        if (hash === this.lastHash) {
            return;
        }
        const parsed = this.parseHash(hash);
        if (parsed) {
            this.movingMap = true;

            this.map.setView(parsed.center, parsed.zoom);

            this.movingMap = false;
        } else {
            this.onMapMove();
        }
    }

    private onHashChange(): void {
        if (!this.changeTimeout) {
            const that = this;
            this.changeTimeout = setTimeout(() => {
                that.update();
                that.changeTimeout = undefined;
            }, this.changeDefer);
        }
    }

    private startListening(): void {
        if (!this.map) return;

        this.map.on('moveend', this.onMapMove, this);

        if (this.HAS_HASHCHANGE) {
            DomEvent.addListener(this.app as HTMLElement, 'hashchange', this.onHashChange);
        } else {
            clearInterval(this.hashChangeInterval);
            this.hashChangeInterval = setInterval(this.onHashChange, 50);
        }

        this.isListening = true;
    }

    private stopListening(): void {
        if (!this.map) return;

        this.map.off('moveend', this.onMapMove, this);

        if (this.HAS_HASHCHANGE) {
            DomEvent.removeListener(this.app as HTMLElement, 'hashchange', this.onHashChange);
        } else {
            clearInterval(this.hashChangeInterval);
        }

        this.isListening = false;
    }
}
