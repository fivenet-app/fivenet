import { RpcError } from '@protobuf-ts/runtime-rpc';
import { defineStore, type StoreDefinition } from 'pinia';
import { type Coordinate } from '~/composables/livemap';
import { Marker, MarkerInfo, UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import { Job } from '~~/gen/ts/resources/users/jobs';
import { type UserShort } from '~~/gen/ts/resources/users/users';
import { LivemapperServiceClient } from '~~/gen/ts/services/livemapper/livemap.client';

// In seconds
const initialReconnectBackoffTime = 1.75;

export interface LivemapState {
    error: RpcError | undefined;
    abort: AbortController | undefined;
    reconnecting: boolean;
    reconnectBackoffTime: number;

    location: Coordinate | undefined;
    showLocationMarker: boolean;
    offsetLocationZoom: boolean;
    zoom: number;

    initiated: boolean;

    jobsMarkers: Job[];
    jobsUsers: Job[];

    markersMarkers: Map<string, Marker>;
    markersUsers: Map<string, UserMarker>;
}

export const useLivemapStore = defineStore('livemap', {
    state: () =>
        ({
            error: undefined,
            abort: undefined,
            reconnecting: false,
            reconnectBackoffTime: 0,

            location: { x: 0, y: 0 },
            showLocationMarker: false,
            offsetLocationZoom: false,
            zoom: 2,

            initiated: false,

            jobsMarkers: [],
            jobsUsers: [],

            markersMarkers: new Map<string, Marker>(),
            markersUsers: new Map<string, UserMarker>(),
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
            this.reconnecting = false;
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

                    if (resp === undefined || !resp.data) {
                        continue;
                    }

                    console.debug('Centrum: Received change - Kind:', resp.data.oneofKind, resp.data);

                    if (resp.data.oneofKind === 'jobs') {
                        this.jobsMarkers = resp.data.jobs.markers;
                        this.jobsUsers = resp.data.jobs.users;

                        this.initiated = true;
                    } else if (resp.data.oneofKind === 'markers') {
                        const foundMarkers: string[] = [];
                        resp.data.markers.markers.forEach((v) => {
                            foundMarkers.push(v.info!.id);
                            this.addOrpdateMarkerMarker(v);
                        });
                        // Remove marker markers not found in latest state
                        let removedMarkers = 0;
                        this.markersMarkers.forEach((_, id) => {
                            if (!foundMarkers.includes(id)) {
                                this.markersMarkers.delete(id);
                                removedMarkers++;
                            }
                        });
                        console.debug(`Livemap: Removed ${removedMarkers} old marker markers`);
                    } else if (resp.data.oneofKind === 'users') {
                        const foundUsers: string[] = [];
                        resp.data.users.users.forEach((v) => {
                            foundUsers.push(v.info!.id);
                            this.addOrpdateUserMarker(v);
                        });
                        // Remove user markers not found in latest state
                        let removedMarkers = 0;
                        this.markersUsers.forEach((_, id) => {
                            if (!foundUsers.includes(id)) {
                                this.markersUsers.delete(id);
                                removedMarkers++;
                            }
                        });
                        console.debug(`Livemap: Removed ${removedMarkers} old user markers`);
                    }
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
            this.reconnecting = true;

            // Reset back off time when over 10 seconds
            if (this.reconnectBackoffTime > 10) {
                this.reconnectBackoffTime = initialReconnectBackoffTime;
            } else {
                this.reconnectBackoffTime += initialReconnectBackoffTime;
            }

            console.debug('Livemap: Restart back off time in', this.reconnectBackoffTime, 'seconds');
            await this.stopStream();

            setTimeout(async () => {
                if (this.reconnecting) {
                    this.startStream();
                }
            }, this.reconnectBackoffTime * 1000);
        },

        addOrpdateMarkerMarker(marker: Marker): void {
            const m = this.markersMarkers.get(marker.info!.id);
            if (m === undefined) {
                this.markersMarkers.set(marker.info!.id, marker);
            } else {
                this.updateMarkerInfo(m.info!, marker.info!);

                if (m.type !== marker.type) {
                    m.type = marker.type;
                }
                m.creatorId = marker.creatorId;
                if (marker.creator !== undefined) {
                    this.updateUserInfo(m.creator!, marker.creator);
                }
                m.data = marker.data;
                m.expiresAt = marker.expiresAt;
            }
        },
        addOrpdateUserMarker(marker: UserMarker): void {
            const m = this.markersUsers.get(marker.info!.id);
            if (m === undefined) {
                this.markersUsers.set(marker.info!.id, marker);
            } else {
                this.updateMarkerInfo(m.info!, marker.info!);

                if (m.userId !== marker.userId) {
                    m.userId = marker.userId;
                    this.updateUserInfo(m.user!, marker.user!);
                }
                if (m.unitId !== marker.unitId) {
                    m.unitId = marker.unitId;
                    m.unit = marker.unit;
                }
            }
        },

        updateMarkerInfo(dest: MarkerInfo, src: MarkerInfo): void {
            dest!.updatedAt = src.updatedAt;
            if (dest!.job !== src!.job) {
                dest!.job = src.job;
            }
            if (dest!.jobLabel !== src.jobLabel) {
                dest!.jobLabel = src.jobLabel;
            }
            if (dest!.name !== src.name) {
                dest!.name = src.name;
            }
            if (dest!.description !== src.description) {
                dest!.description = src.description;
            }
            if (dest!.x !== src.x) {
                dest!.x = src.x;
            }
            if (dest!.y !== src.y) {
                dest!.y = src.y;
            }
            if (dest!.postal !== src.postal) {
                dest!.postal = src.postal;
            }
            if (dest!.color !== src.color) {
                dest!.color = src.color;
            }
            if (dest!.icon !== src.icon) {
                dest!.icon = src.icon;
            }
        },
        updateUserInfo(dest: UserShort, src: UserShort): void {
            if (dest.firstname !== src.firstname) {
                dest.firstname = src.firstname;
            }
            if (dest.lastname !== src.lastname) {
                dest.lastname = src.lastname;
            }
            if (dest.job !== src.job) {
                dest.job = src.job;
            }
            if (dest.jobLabel !== src.jobLabel) {
                dest.jobLabel = src.jobLabel;
            }
            if (dest.jobGrade !== src.jobGrade) {
                dest.jobGrade = src.jobGrade;
            }
            if (dest.jobGradeLabel !== src.jobGradeLabel) {
                dest.jobGradeLabel = src.jobGradeLabel;
            }
            if (dest.dateofbirth !== src.dateofbirth) {
                dest.dateofbirth = src.dateofbirth;
            }
            if (dest.phoneNumber !== src.phoneNumber) {
                dest.phoneNumber = src.phoneNumber;
            }
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useLivemapStore as unknown as StoreDefinition, import.meta.hot));
}
