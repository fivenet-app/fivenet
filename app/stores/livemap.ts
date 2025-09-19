import type { RpcError, ServerStreamingCall } from '@protobuf-ts/runtime-rpc';
import { defineStore } from 'pinia';
import { getLivemapLivemapClient } from '~~/gen/ts/clients';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import type { MarkerMarker } from '~~/gen/ts/resources/livemap/marker_marker';
import type { UserMarker } from '~~/gen/ts/resources/livemap/user_marker';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { StreamRequest } from '~~/gen/ts/services/centrum/centrum';
import type { StreamResponse } from '~~/gen/ts/services/livemap/livemap';
import type { Coordinate } from '~~/shared/types/types';
import { useSettingsStore } from './settings';

const logger = useLogger('🗺️ Livemap');

// In seconds
const maxBackOffTime = 15;
const initialReconnectBackoffTime = 1.75;

export const useLivemapStore = defineStore(
    'livemap',
    () => {
        // State
        const error = ref<RpcError | undefined>(undefined);
        const abort = ref<AbortController | undefined>(undefined);
        const stopping = ref<boolean>(false);
        const reconnectBackoffTime = ref<number>(0);

        const location = ref<Coordinate | undefined>();
        const showLocationMarker = ref<boolean>(false);
        const zoom = ref<number>(2);

        const initiated = ref<boolean>(false);
        const userOnDuty = ref<boolean>(false);

        const jobsMarkers = ref<Job[]>([]);
        const jobsUsers = ref<Job[]>([]);

        const markersMarkers = ref<Map<number, MarkerMarker>>(new Map());
        const markersUsers = ref<Map<number, UserMarker>>(new Map());
        const ownMarker = ref<UserMarker | undefined>();

        const selectedMarker = ref<UserMarker | undefined>(undefined);

        // Actions
        const cleanupMarkerMarkers = (): void => {
            const now = new Date();
            markersMarkers.value.forEach((m) => {
                if (!m.expiresAt) return;

                if (toDate(m.expiresAt).getTime() > now.getTime()) return;

                markersMarkers.value.delete(m.id);
            });
        };

        // Stream
        let currentStream: ServerStreamingCall<StreamRequest, StreamResponse> | undefined = undefined;

        const startStream = async (): Promise<void> => {
            if (abort.value !== undefined) return;
            stopping.value = false;

            logger.debug('Starting Stream');

            const { activeChar } = useAuth();
            // Access settings
            const settingsStore = useSettingsStore();
            const { livemap } = storeToRefs(settingsStore);

            abort.value = new AbortController();
            error.value = undefined;

            // Tracking marker and user markers between part responses
            const foundMarkers: number[] = [];

            cleanupMarkerMarkers();

            const livemapLivemapClient = await getLivemapLivemapClient();

            try {
                currentStream = livemapLivemapClient.stream({}, { abort: abort.value.signal });

                for await (const respRaw of currentStream.responses) {
                    // The gRPC stream may yield unknown, so cast to the expected type
                    const resp = respRaw as StreamResponse;
                    error.value = undefined;

                    if (!resp || !resp.data) {
                        continue;
                    }

                    if (resp.userOnDuty !== undefined) {
                        userOnDuty.value = resp.userOnDuty;
                    }

                    logger.debug('Received change - oneofKind:', resp.data.oneofKind, resp.data);

                    if (resp.data.oneofKind === 'jobs') {
                        jobsMarkers.value = resp.data.jobs.markers;
                        jobsUsers.value = resp.data.jobs.users;
                    } else if (resp.data.oneofKind === 'markers') {
                        resp.data.markers.updated.forEach((v: MarkerMarker) => {
                            // Only record found users for non-partial responses
                            if (resp.data.oneofKind === 'markers' && !resp.data.markers.partial) {
                                foundMarkers.push(v.id);
                            }

                            addOrUpdateMarkerMarker(v);
                        });

                        resp.data.markers.deleted.forEach((id: number) => markersMarkers.value.delete(id));

                        if (!resp.data.markers.partial) {
                            if (resp.data.markers.part <= 0) {
                                // Remove markers not found in the latest full state
                                let removedMarkers = 0;
                                markersMarkers.value.forEach((_, id) => {
                                    if (!foundMarkers.includes(id)) {
                                        markersMarkers.value.delete(id);
                                        removedMarkers++;
                                    }
                                });
                                foundMarkers.length = 0;
                                logger.debug(`Removed ${removedMarkers} old marker markers`);
                            }
                        }
                    } else if (resp.data.oneofKind === 'snapshot') {
                        // Handle snapshot response
                        const snapshot = resp.data.snapshot;

                        markersUsers.value.clear();
                        // Add all markers from the snapshot
                        snapshot.markers.forEach((marker: UserMarker) => {
                            addOrUpdateUserMarker(marker);
                        });

                        initiated.value = true;
                    } else if (resp.data.oneofKind === 'userDeletes') {
                        // Handle user deletion
                        for (const userDelete of resp.data.userDeletes.deletes) {
                            const userId = userDelete.id;

                            const um = markersUsers.value.get(userId);
                            if (!um || um.job !== userDelete.job) {
                                continue; // Skip if the user marker does not match the job
                            }

                            markersUsers.value.delete(userId);

                            // If the deleted user was selected, clear the selection
                            if (selectedMarker.value?.userId === userId) {
                                selectedMarker.value = undefined;
                            }
                            if (ownMarker.value?.userId === userId) {
                                ownMarker.value = undefined;
                            }

                            logger.debug('User marker deleted:', userId);
                        }
                    } else if (resp.data.oneofKind === 'userUpdates') {
                        for (const userUpdate of resp.data.userUpdates.updates) {
                            addOrUpdateUserMarker(userUpdate);
                            // If a marker is selected, update it
                            if (livemap.value.centerSelectedMarker && userUpdate.userId === selectedMarker.value?.userId) {
                                selectedMarker.value = userUpdate;
                            }
                            if (activeChar.value?.userId === userUpdate.userId) {
                                ownMarker.value = userUpdate;
                            }
                        }
                    } else {
                        logger.warn('Unknown data received - oneofKind:' + resp.data.oneofKind);
                    }

                    if (!userOnDuty.value) {
                        if (markersUsers.value.size > 0) markersUsers.value.clear();
                        logger.info('User is not on duty, clearing user markers');
                    }
                }
            } catch (e) {
                const err = e as RpcError;

                // Always clear the error first
                error.value = undefined;

                // If the stream was cancelled or aborted
                if (err.code === 'CANCELLED' || err.code === 'ABORTED') {
                    // Only restart if not manually aborted
                    if (!abort.value?.signal.aborted) {
                        await restartStream();
                    }
                    // Otherwise, do nothing (intentional stop)
                    return;
                }

                // For all other errors, log and attempt restart if not manually aborted
                logger.error('Stream failed', err.code, err.message, err.cause);

                if (!abort.value?.signal.aborted) {
                    await restartStream();
                } else {
                    error.value = err;
                }
            }

            logger.debug('Stream ended');
        };

        const stopStream = async (end?: boolean): Promise<void> => {
            if (end === true) stopping.value = true;

            if (abort.value) {
                abort.value.abort();
                logger.debug('Stopping Stream');
            }

            abort.value = undefined;
        };

        const restartStream = async (): Promise<void> => {
            if (!abort.value || abort.value.signal.aborted) return;

            // Reset back off time if it exceeds max
            if (reconnectBackoffTime.value > maxBackOffTime) {
                reconnectBackoffTime.value = initialReconnectBackoffTime;
            } else {
                reconnectBackoffTime.value += initialReconnectBackoffTime;
            }

            logger.debug('Restart back off time in', reconnectBackoffTime.value, 'seconds');
            await stopStream();

            useTimeoutFn(async () => {
                if (!stopping.value) {
                    startStream();
                }
            }, reconnectBackoffTime.value * 1000);
        };

        const addOrUpdateMarkerMarker = (marker: MarkerMarker): void => {
            const m = markersMarkers.value.get(marker.id);
            if (!m) {
                markersMarkers.value.set(marker.id, marker);
            } else {
                updateMarkerMarker(m, marker);
            }
        };

        const updateMarkerMarker = (dest: MarkerMarker, src: MarkerMarker): void => {
            if (dest.x !== src.x) {
                dest.x = src.x;
            }
            if (dest.y !== src.y) {
                dest.y = src.y;
            }
            if (dest.createdAt !== src.createdAt) {
                dest.createdAt = src.createdAt;
            }
            if (dest.updatedAt !== src.updatedAt) {
                dest.updatedAt = src.updatedAt;
            }
            if (dest.expiresAt !== src.expiresAt) {
                dest.expiresAt = src.expiresAt;
            }
            if (dest.deletedAt !== src.deletedAt) {
                dest.deletedAt = src.deletedAt;
            }
            if (dest.name !== src.name) {
                dest.name = src.name;
            }
            if (dest.description !== src.description) {
                dest.description = src.description;
            }
            if (dest.postal !== src.postal) {
                dest.postal = src.postal;
            }
            if (dest.color !== src.color) {
                dest.color = src.color;
            }
            if (dest.job !== src.job) {
                dest.job = src.job;
            }
            if (dest.jobLabel !== src.jobLabel) {
                dest.jobLabel = src.jobLabel;
            }
            if (dest.type !== src.type) {
                dest.type = src.type;
            }
            if (dest.data !== src.data) {
                dest.data = src.data;
            }

            dest.creatorId = src.creatorId;
            if (src.creator !== undefined) {
                updateUserInfo(dest.creator!, src.creator);
            }
        };

        const addOrUpdateUserMarker = (marker: UserMarker): void => {
            const m = markersUsers.value.get(marker.userId);
            if (!m) {
                markersUsers.value.set(marker.userId, marker);
            } else {
                updateUserMarker(m, marker);
            }
        };

        const updateUserMarker = (dest: UserMarker, src: UserMarker): void => {
            if (dest.x !== src.x) {
                dest.x = src.x;
            }
            if (dest.y !== src.y) {
                dest.y = src.y;
            }
            if (dest.updatedAt !== src.updatedAt) {
                dest.updatedAt = src.updatedAt;
            }
            if (dest.postal !== src.postal) {
                dest.postal = src.postal;
            }
            if (dest.color !== src.color) {
                dest.color = src.color;
            }
            if (dest.job !== src.job) {
                dest.job = src.job;
            }
            if (dest.jobLabel !== src.jobLabel) {
                dest.jobLabel = src.jobLabel;
            }

            if (dest.userId !== src.userId) {
                dest.userId = src.userId;
                updateUserInfo(dest.user!, src.user!);
            }
            if (dest.unitId !== src.unitId) {
                dest.unitId = src.unitId;
                dest.unit = src.unit;
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
            if (dest.profilePicture !== src.profilePicture) {
                dest.profilePicture = src.profilePicture;
            }
        };

        const deleteMarkerMarker = (id: number): void => {
            markersMarkers.value.delete(id);
        };

        const goto = async (loc: Coordinate, ingame = true): Promise<void> => {
            location.value = loc;

            if (ingame) {
                // Set in-game waypoint via NUI
                return setWaypoint(loc.x, loc.y);
            }
        };

        return {
            // State
            error,
            abort,
            stopping,
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
            ownMarker,
            selectedMarker,

            // Actions
            startStream,
            stopStream,
            restartStream,
            addOrUpdateMarkerMarker,
            addOrUpdateUserMarker,
            updateUserMarker,
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
