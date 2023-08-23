import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { StoreDefinition, defineStore } from 'pinia';
import { UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import { Job } from '~~/gen/ts/resources/users/jobs';
import { LivemapperServiceClient } from '~~/gen/ts/services/livemapper/livemap.client';

export interface LivemapState {
    error: RpcError | undefined;
    abort: AbortController | undefined;
    location: Coordinate | undefined;
    zoom: number;
    jobs: {
        users: Job[];
    };
    markers: {
        users: UserMarker[];
    };
}

export const useLivemapStore = defineStore('livemap', {
    state: () =>
        ({
            error: undefined,
            abort: undefined,
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
            try {
                this.abort = new AbortController();

                const { $grpc } = useNuxtApp();
                const call = new LivemapperServiceClient($grpc.getTransport()).stream(
                    {},
                    {
                        abort: this.abort.signal,
                    },
                );

                for await (let resp of call.responses) {
                    this.error = undefined;

                    if (resp === undefined || !resp.jobsUsers) {
                        continue;
                    }

                    this.jobs.users = resp.jobsUsers;
                    this.markers.users = resp.users;
                }
            } catch (e) {
                this.error = e as RpcError;
                if (this.error)
                    console.error('Livemap: Data Stream Failed', this.error.code, this.error.message, this.error.cause);

                this.restartStream();
            }

            console.debug('Livemap: Data Stream Ended');
        },
        async stopStream(): Promise<void> {
            console.debug('Livemap: Stopping Data Stream');
            if (this.abort) this.abort.abort();
            this.abort = undefined;
            this.$reset();
        },
        async restartStream(): Promise<void> {
            console.debug('Centrum: Restarting Data Stream');
            await this.stopStream();

            setTimeout(async () => this.startStream(), 250);
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useLivemapStore as unknown as StoreDefinition, import.meta.hot));
}
