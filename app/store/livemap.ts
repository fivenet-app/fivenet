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

export const useLivemapStore = defineStore(
    'livemap',
    () => {
        const { $grpc } = useNuxtApp();

        // State
        const error = ref<RpcError | undefined>(undefined);
        const abort = ref<AbortController | undefined>(undefined);
        const reconnecting = ref<boolean>(false);
        const reconnectBackoffTime = ref<number>(0);

        const location = ref<Coordinate>({ x: 0, y: 0 });
        const showLocationMarker = ref<boolean>(false);
        const zoom = ref<number>(2);

        const initiated = ref<boolean>(false);
        const userOnDuty = ref<boolean>(false);

        const jobsMarkers = ref<Job[]>([]);
        const jobsUsers = ref<Job[]>([]);

        const markersMarkers = ref<Map<number, MarkerMarker>>(new Map());
        const markersUsers = ref<Map<number, UserMarker>>(new Map());

        const selectedMarker = ref<UserMarker | undefined>(undefined);

        // Actions
        const startStream = async (): Promise<void> => {
            if (abort.value !== undefined) {
                return;
            }

            logger.debug('Starting Stream');

            // Access settings
            const settingsStore = useSettingsStore();
            const { livemap } = storeToRefs(settingsStore);

            abort.value = new AbortController();
            error.value = undefined;
            reconnecting.value = false;

            try {
                const call = $grpc.livemapper.livemapper.stream({}, { abort: abort.value.signal });

                // For partial user updates
                const foundUsers: number[] = [];

                for await (const resp of call.responses) {
                    error.value = undefined;

                    if (!resp || !resp.data) {
                        continue;
                    }

                    if (resp.userOnDuty !== undefined) {
                        userOnDuty.value = resp.userOnDuty;
                    }

                    logger.debug('Received change - Kind:', resp.data.oneofKind, resp.data);

                    if (resp.data.oneofKind === 'jobs') {
                        jobsMarkers.value = resp.data.jobs.markers;
                        jobsUsers.value = resp.data.jobs.users;
                    } else if (resp.data.oneofKind === 'markers') {
                        const foundMarkers: number[] = [];
                        resp.data.markers.markers.forEach((v) => {
                            foundMarkers.push(v.info!.id);
                            addOrUpdateMarkerMarker(v);
                        });

                        // Remove markers not found in the latest state
                        let removedMarkers = 0;
                        markersMarkers.value.forEach((_, id) => {
                            if (!foundMarkers.includes(id)) {
                                markersMarkers.value.delete(id);
                                removedMarkers++;
                            }
                        });
                        foundMarkers.length = 0;
                        logger.debug(`Removed ${removedMarkers} old marker markers`);
                    } else if (resp.data.oneofKind === 'users') {
                        resp.data.users.users.forEach((v) => {
                            foundUsers.push(v.info!.id);
                            addOrUpdateUserMarker(v);

                            // If auto-centering on selected marker
                            if (livemap.value.centerSelectedMarker && v.info!.id === selectedMarker.value?.info?.id) {
                                selectedMarker.value = v;
                            }
                        });

                        if (resp.data.users.part <= 0) {
                            // Remove user markers not found in the latest state
                            let removedMarkers = 0;
                            markersUsers.value.forEach((_, id) => {
                                if (!foundUsers.includes(id)) {
                                    markersUsers.value.delete(id);

                                    if (id === selectedMarker.value?.info?.id) {
                                        selectedMarker.value = undefined;
                                    }
                                    removedMarkers++;
                                }
                            });
                            foundUsers.length = 0;
                            logger.debug(`Removed ${removedMarkers} old user markers`);
                        }

                        initiated.value = true;
                    } else {
                        logger.warn('Unknown data received - Kind: ' + resp.data.oneofKind);
                    }
                }
            } catch (e) {
                const err = e as RpcError;

                // Only restart if not CANCELLED or ABORTED
                if (err.code !== 'CANCELLED' && err.code !== 'ABORTED') {
                    logger.error('Stream failed', err.code, err.message, err.cause);

                    // If we haven't manually aborted, attempt restart
                    if (!abort.value?.signal.aborted) {
                        restartStream();
                    } else {
                        error.value = err;
                    }
                } else {
                    error.value = undefined;

                    // Restart only if not manually aborted
                    if (!abort.value?.signal.aborted) {
                        await restartStream();
                    }
                }
            }

            logger.debug('Stream ended');
        };

        const stopStream = async (): Promise<void> => {
            if (!abort.value) {
                return;
            }

            abort.value.abort();
            logger.debug('Stopping Stream');
            abort.value = undefined;
        };

        const restartStream = async (): Promise<void> => {
            if (!abort.value || abort.value.signal.aborted) {
                return;
            }

            reconnecting.value = true;

            // Reset back off time if it exceeds max
            if (reconnectBackoffTime.value > maxBackOffTime) {
                reconnectBackoffTime.value = initialReconnectBackoffTime;
            } else {
                reconnectBackoffTime.value += initialReconnectBackoffTime;
            }

            logger.debug('Restart back off time in', reconnectBackoffTime.value, 'seconds');
            await stopStream();

            setTimeout(async () => {
                if (reconnecting.value) {
                    startStream();
                }
            }, reconnectBackoffTime.value * 1000);
        };

        const addOrUpdateMarkerMarker = (marker: MarkerMarker): void => {
            const m = markersMarkers.value.get(marker.info!.id);
            if (!m) {
                markersMarkers.value.set(marker.info!.id, marker);
            } else {
                updateMarkerInfo(m.info!, marker.info!);

                if (m.type !== marker.type) {
                    m.type = marker.type;
                }
                m.creatorId = marker.creatorId;
                if (marker.creator !== undefined) {
                    updateUserInfo(m.creator!, marker.creator);
                }
                m.data = marker.data;
                if (m.expiresAt !== marker.expiresAt) {
                    m.expiresAt = marker.expiresAt;
                }
            }
        };

        const addOrUpdateUserMarker = (marker: UserMarker): void => {
            const m = markersUsers.value.get(marker.info!.id);
            if (!m) {
                markersUsers.value.set(marker.info!.id, marker);
            } else {
                updateMarkerInfo(m.info!, marker.info!);

                if (m.userId !== marker.userId) {
                    m.userId = marker.userId;
                    updateUserInfo(m.user!, marker.user!);
                }
                if (m.unitId !== marker.unitId) {
                    m.unitId = marker.unitId;
                    m.unit = marker.unit;
                }
            }
        };

        const updateMarkerInfo = (dest: MarkerInfo, src: MarkerInfo): void => {
            if (dest.updatedAt !== src.updatedAt) {
                dest.updatedAt = src.updatedAt;
            }
            if (dest.job !== src.job) {
                dest.job = src.job;
            }
            if (dest.jobLabel !== src.jobLabel) {
                dest.jobLabel = src.jobLabel;
            }
            if (dest.name !== src.name) {
                dest.name = src.name;
            }
            if (dest.description !== src.description) {
                dest.description = src.description;
            }
            if (dest.x !== src.x) {
                dest.x = src.x;
            }
            if (dest.y !== src.y) {
                dest.y = src.y;
            }
            if (dest.postal !== src.postal) {
                dest.postal = src.postal;
            }
            if (dest.color !== src.color) {
                dest.color = src.color;
            }
            if (dest.icon !== src.icon) {
                dest.icon = src.icon;
            }
        };

        const updateUserInfo = (dest: UserShort, src: UserShort): void => {
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
        };

        const deleteMarkerMarker = (id: number): void => {
            markersMarkers.value.delete(id);
        };

        const goto = async (loc: Coordinate): Promise<void> => {
            location.value = loc;

            // Set in-game waypoint via NUI
            return setWaypoint(loc.x, loc.y);
        };

        return {
            // State
            error,
            abort,
            reconnecting,
            reconnectBackoffTime,
            location,
            showLocationMarker,
            zoom,
            initiated,
            userOnDuty,
            jobsMarkers,
            jobsUsers,
            markersMarkers,
            markersUsers,
            selectedMarker,

            // Actions
            startStream,
            stopStream,
            restartStream,
            addOrUpdateMarkerMarker,
            addOrUpdateUserMarker,
            updateMarkerInfo,
            updateUserInfo,
            deleteMarkerMarker,
            goto,
        };
    },
    {
        persist: false,
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useLivemapStore, import.meta.hot));
}
