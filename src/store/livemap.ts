import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { StoreDefinition, defineStore } from 'pinia';
import { Marker, UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import { Job } from '~~/gen/ts/resources/users/jobs';
import { LivemapperServiceClient } from '~~/gen/ts/services/livemapper/livemap.client';

// In seconds
const initialBackoffTime = 2;

export interface LivemapState {
    error: RpcError | undefined;
    abort: AbortController | undefined;
    restarting: boolean;
    restartBackoffTime: number;

    location: Coordinate | undefined;
    zoom: number;
    jobs: {
        users: Job[];
        markers: Job[];
    };
    markers: {
        users: UserMarker[];
        markers: Marker[];
    };
}

export const useLivemapStore = defineStore('livemap', {
    state: () =>
        ({
            error: undefined,
            abort: undefined,
            restarting: false,
            restartBackoffTime: 0,

            location: { x: 0, y: 0 },
            zoom: 2,
            jobs: {
                users: [] as Job[],
            },
            markers: {
                users: [] as UserMarker[],
            },
        }) as LivemapState,
    persist: false,
    actions: {
        async startStream(): Promise<void> {
            if (this.abort !== undefined) return;

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

                for await (let resp of call.responses) {
                    this.error = undefined;

                    if (resp === undefined) {
                        continue;
                    }

                    this.jobs.users = resp.jobsUsers;
                    this.markers.users = resp.users;
                    this.jobs.markers = resp.jobsMarkers;
                    this.markers.markers = resp.markers;
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
