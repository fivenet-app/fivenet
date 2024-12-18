import { defineStore } from 'pinia';
import type { Coordinate } from '~/composables/livemap';
import type { MarkerInfo, MarkerMarker, UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import type { Job } from '~~/gen/ts/resources/users/jobs';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import { useSettingsStore } from './settings';

const logger = useLogger('🗺️ Livemap');

// In seconds
const maxBackOffTime = 15;
const initialReconnectBackoffTime = 1.75;

export interface LivemapState {
    error: RpcError | undefined;
    abort: AbortController | undefined;
    reconnecting: boolean;
    reconnectBackoffTime: number;

    location: Coordinate | undefined;
    showLocationMarker: boolean;
    zoom: number;

    initiated: boolean;

    userOnDuty: boolean;
    jobsMarkers: Job[];
    jobsUsers: Job[];

    markersMarkers: Map<string, MarkerMarker>;
    markersUsers: Map<string, UserMarker>;

    selectedMarker: UserMarker | undefined;
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
            zoom: 2,

            initiated: false,

            userOnDuty: false,
            jobsMarkers: [],
            jobsUsers: [],

            markersMarkers: new Map<string, MarkerMarker>(),
            markersUsers: new Map<string, UserMarker>(),

            selectedMarker: undefined,
        }) as LivemapState,
    persist: false,
    actions: {
        async startStream(): Promise<void> {
            if (this.abort !== undefined) {
                return;
            }

            logger.debug('Starting Stream');

            const settingsStore = useSettingsStore();
            const { livemap } = storeToRefs(settingsStore);

            this.abort = new AbortController();
            this.error = undefined;
            this.reconnecting = false;

            try {
                const call = getGRPCLivemapperClient().stream(
                    {},
                    {
                        abort: this.abort.signal,
                    },
                );

                const foundUsers: string[] = [];

                for await (const resp of call.responses) {
                    this.error = undefined;

                    if (resp === undefined || !resp.data) {
                        continue;
                    }

                    if (resp.userOnDuty !== undefined) {
                        this.userOnDuty = resp.userOnDuty;
                    }

                    logger.debug('Received change - Kind:', resp.data.oneofKind, resp.data);

                    if (resp.data.oneofKind === 'jobs') {
                        this.jobsMarkers = resp.data.jobs.markers;
                        this.jobsUsers = resp.data.jobs.users;
                    } else if (resp.data.oneofKind === 'markers') {
                        const foundMarkers: string[] = [];
                        resp.data.markers.markers.forEach((v) => {
                            foundMarkers.push(v.info!.id);
                            this.addOrUpdateMarkerMarker(v);
                        });
                        // Remove marker markers not found in latest state
                        let removedMarkers = 0;
                        this.markersMarkers.forEach((_, id) => {
                            if (!foundMarkers.includes(id)) {
                                this.markersMarkers.delete(id);
                                removedMarkers++;
                            }
                        });
                        foundMarkers.length = 0;
                        logger.debug(`Removed ${removedMarkers} old marker markers`);
                    } else if (resp.data.oneofKind === 'users') {
                        resp.data.users.users.forEach((v) => {
                            foundUsers.push(v.info!.id);
                            this.addOrUpdateUserMarker(v);

                            if (livemap.value.centerSelectedMarker && v.info!.id === this.selectedMarker?.info?.id) {
                                this.selectedMarker = v;
                            }
                        });

                        if (resp.data.users.part <= 0) {
                            // Remove user markers not found in latest state
                            let removedMarkers = 0;
                            this.markersUsers.forEach((_, id) => {
                                if (!foundUsers.includes(id)) {
                                    this.markersUsers.delete(id);

                                    if (id === this.selectedMarker?.info?.id) {
                                        this.selectedMarker = undefined;
                                    }
                                    removedMarkers++;
                                }
                            });
                            foundUsers.length = 0;
                            logger.debug(`Removed ${removedMarkers} old user markers`);
                        }

                        this.initiated = true;
                    } else {
                        logger.warn('Unknown data received - Kind: ' + resp.data.oneofKind);
                    }
                }
            } catch (e) {
                const error = e as RpcError;
                // Only restart when not cancelled and abort is still valid
                if (error.code !== 'CANCELLED' && error.code !== 'ABORTED') {
                    logger.error('Stream failed', error.code, error.message, error.cause);

                    // Only set error if we don't need to restart
                    if (this.abort !== undefined && !this.abort?.signal.aborted) {
                        this.restartStream();
                    } else {
                        this.error = error;
                    }
                } else {
                    this.error = undefined;

                    // Only restart stream when not aborted
                    if (!this.abort?.signal.aborted) {
                        await this.restartStream();
                    }
                }
            }

            logger.debug('Stream ended');
        },

        async stopStream(): Promise<void> {
            if (this.abort === undefined) {
                return;
            }

            this.abort.abort();
            this.abort = undefined;
            logger.debug('Stopping Stream');
        },

        async restartStream(): Promise<void> {
            if (this.abort === undefined || this.abort.signal.aborted) {
                return;
            }

            this.reconnecting = true;

            // Reset back off time when over 10 seconds
            if (this.reconnectBackoffTime > maxBackOffTime) {
                this.reconnectBackoffTime = initialReconnectBackoffTime;
            } else {
                this.reconnectBackoffTime += initialReconnectBackoffTime;
            }

            logger.debug('Restart back off time in', this.reconnectBackoffTime, 'seconds');
            await this.stopStream();

            setTimeout(async () => {
                if (this.reconnecting) {
                    this.startStream();
                }
            }, this.reconnectBackoffTime * 1000);
        },

        addOrUpdateMarkerMarker(marker: MarkerMarker): void {
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
                if (m.expiresAt !== marker.expiresAt) {
                    m.expiresAt = marker.expiresAt;
                }
            }
        },
        addOrUpdateUserMarker(marker: UserMarker): void {
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
            if (dest!.updatedAt !== src.updatedAt) {
                dest.updatedAt = src.updatedAt;
            }
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
            if (dest.avatar !== src.avatar) {
                dest.avatar = src.avatar;
            }
        },

        deleteMarkerMarker(id: string): void {
            this.markersMarkers.delete(id);
        },

        async goto(loc: Coordinate): Promise<void> {
            this.location = loc;

            // Set in-game waypoint via NUI
            return setWaypoint(loc.x, loc.y);
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useLivemapStore, import.meta.hot));
}
