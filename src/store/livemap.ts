import { RpcError } from '@protobuf-ts/runtime-rpc';
import { defineStore, type StoreDefinition } from 'pinia';
import { Marker, UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import { Job } from '~~/gen/ts/resources/users/jobs';
import { LivemapperServiceClient } from '~~/gen/ts/services/livemapper/livemap.client';
import { type Coordinate } from '~/composables/livemap';

// In seconds
const initialBackoffTime = 1.75;

export interface LivemapState {
    error: RpcError | undefined;
    abort: AbortController | undefined;
    restarting: boolean;
    restartBackoffTime: number;

    location: Coordinate | undefined;
    offsetLocationZoom: boolean;
    zoom: number;

    initiated: boolean;

    jobsMarkers: Job[];
    jobsUsers: Job[];

    markersMarkers: Marker[];
    markersUsers: UserMarker[];
}

export const useLivemapStore = defineStore('livemap', {
    state: () =>
        ({
            error: undefined,
            abort: undefined,
            restarting: false,
            restartBackoffTime: 0,

            location: { x: 0, y: 0 },
            offsetLocationZoom: false,
            zoom: 2,

            initiated: false,

            jobsMarkers: [] as Job[],
            jobsUsers: [] as Job[],

            markersMarkers: [] as Marker[],
            markersUsers: [] as UserMarker[],
        }) as LivemapState,
    persist: false,
    actions: {
        async startStream(): Promise<void> {
            if (this.abort !== undefined) {
                return;
            }

            console.debug('Livemap: Starting Data Stream');

            this.abort = new AbortController();
            this.error = undefined;
            this.restarting = false;
            const { $grpc } = useNuxtApp();

            try {
                const call = new LivemapperServiceClient($grpc.getTransport()).stream(
                    {},
                    {
                        abort: this.abort.signal,
                    },
                );

                for await (const resp of call.responses) {
                    this.error = undefined;

                    if (resp === undefined) {
                        continue;
                    }

                    if (!this.initiated) {
                        this.jobsMarkers = resp.jobsMarkers;
                        this.jobsUsers = resp.jobsUsers;

                        this.initiated = true;
                    }

                    // Sort markers by id
                    this.markersUsers = resp.users.sort((a, b) => (a.info?.id ?? '0').localeCompare(b.info?.id ?? '0'));
                    this.markersMarkers = resp.markers.sort((a, b) => (a.info?.id ?? '0').localeCompare(b.info?.id ?? '0'));
                }
            } catch (e) {
                const error = e as RpcError;
                if (error) {
                    // Only restart when not cancelled and abort is still valid
                    if (error.code !== 'CANCELLED' && error.code !== 'ABORTED') {
                        console.error('Livemap: Data Stream Failed', error.code, error.message, error.cause);

                        // Only set error if we don't need to restart
                        if (this.abort !== undefined && !this.abort?.signal.aborted) {
                            this.restartStream();
                        } else {
                            this.error = error;
                        }
                    } else {
                        this.error = undefined;
                    }
                }
            }

            console.debug('Livemap: Data Stream Ended');
        },
        async stopStream(): Promise<void> {
            if (this.abort !== undefined) this.abort.abort();
            this.abort = undefined;
            console.debug('Livemap: Stopping Data Stream');
        },
        async restartStream(): Promise<void> {
            this.restarting = true;

            // Reset back off time when over 10 seconds
            if (this.restartBackoffTime > 10) {
                this.restartBackoffTime = initialBackoffTime;
            } else {
                this.restartBackoffTime += initialBackoffTime;
            }

            console.debug('Livemap: Restart back off time in', this.restartBackoffTime, 'seconds');
            await this.stopStream();

            setTimeout(async () => {
                if (this.restarting) {
                    this.startStream();
                }
            }, this.restartBackoffTime * 1000);
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useLivemapStore as unknown as StoreDefinition, import.meta.hot));
}
