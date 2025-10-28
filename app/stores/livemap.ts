import type { RpcError, ServerStreamingCall } from '@protobuf-ts/runtime-rpc';
import { defineStore } from 'pinia';
import { getLivemapLivemapClient } from '~~/gen/ts/clients';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import type { MarkerMarker } from '~~/gen/ts/resources/livemap/marker_marker';
import type { UserMarker } from '~~/gen/ts/resources/livemap/user_marker';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { StreamRequest } from '~~/gen/ts/services/centrum/centrum';
import type {
    JobsList,
    MarkerMarkersUpdates,
    Snapshot,
    StreamResponse,
    UserDeletes,
    UserUpdates,
} from '~~/gen/ts/services/livemap/livemap';
import type { Coordinate } from '~~/shared/types/types';
import { useSettingsStore, type LivemapSettings } from './settings';

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
        const followMarker = ref<boolean>(false);

        // Actions

        // Marker Management
        /**
         * Remove expired marker markers from the map.
         */
        const cleanupMarkerMarkers = (): void => {
            const now = new Date();
            markersMarkers.value.forEach((m) => {
                if (!m.expiresAt) return;
                if (toDate(m.expiresAt).getTime() > now.getTime()) return;
                markersMarkers.value.delete(m.id);
            });
        };

        /**
         * Add or update a marker marker in the map.
         */
        const addOrUpdateMarkerMarker = (marker: MarkerMarker): void => {
            const m = markersMarkers.value.get(marker.id);
            if (!m) {
                markersMarkers.value.set(marker.id, marker);
            } else {
                updateMarkerMarker(m, marker);
            }
        };

        /**
         * Update an existing marker marker with new data.
         */
        const updateMarkerMarker = (dest: MarkerMarker, src: MarkerMarker): void => {
            if (dest.x !== src.x) dest.x = src.x;
            if (dest.y !== src.y) dest.y = src.y;
            if (dest.createdAt !== src.createdAt) dest.createdAt = src.createdAt;
            if (dest.updatedAt !== src.updatedAt) dest.updatedAt = src.updatedAt;
            if (dest.expiresAt !== src.expiresAt) dest.expiresAt = src.expiresAt;
            if (dest.deletedAt !== src.deletedAt) dest.deletedAt = src.deletedAt;
            if (dest.name !== src.name) dest.name = src.name;
            if (dest.description !== src.description) dest.description = src.description;
            if (dest.postal !== src.postal) dest.postal = src.postal;
            if (dest.color !== src.color) dest.color = src.color;
            if (dest.job !== src.job) dest.job = src.job;
            if (dest.jobLabel !== src.jobLabel) dest.jobLabel = src.jobLabel;
            if (dest.type !== src.type) dest.type = src.type;
            if (dest.data !== src.data) dest.data = src.data;
            dest.creatorId = src.creatorId;
            if (src.creator !== undefined) {
                updateUserInfo(dest.creator!, src.creator);
            }
        };

        /**
         * Delete a marker marker by id.
         */
        const deleteMarkerMarker = (id: number): void => {
            markersMarkers.value.delete(id);
        };

        // User Marker Management
        /**
         * Add or update a user marker in the map.
         */
        const addOrUpdateUserMarker = (marker: UserMarker): void => {
            const m = markersUsers.value.get(marker.userId);
            if (!m) {
                markersUsers.value.set(marker.userId, marker);
            } else {
                updateUserMarker(m, marker);
            }
        };

        /**
         * Update an existing user marker with new data.
         */
        const updateUserMarker = (dest: UserMarker, src: UserMarker): void => {
            if (dest.x !== src.x) dest.x = src.x;
            if (dest.y !== src.y) dest.y = src.y;
            if (dest.updatedAt !== src.updatedAt) dest.updatedAt = src.updatedAt;
            if (dest.postal !== src.postal) dest.postal = src.postal;
            if (dest.color !== src.color) dest.color = src.color;
            if (dest.job !== src.job) dest.job = src.job;
            if (dest.jobLabel !== src.jobLabel) dest.jobLabel = src.jobLabel;
            if (dest.userId !== src.userId) {
                dest.userId = src.userId;
                updateUserInfo(dest.user!, src.user!);
            }
            if (dest.unitId !== src.unitId) {
                dest.unitId = src.unitId;
                dest.unit = src.unit;
            }
        };

        /**
         * Update user info for a UserShort object.
         */
        const updateUserInfo = (dest: UserShort, src: UserShort): void => {
            if (dest.firstname !== src.firstname) dest.firstname = src.firstname;
            if (dest.lastname !== src.lastname) dest.lastname = src.lastname;
            if (dest.job !== src.job) dest.job = src.job;
            if (dest.jobLabel !== src.jobLabel) dest.jobLabel = src.jobLabel;
            if (dest.jobGrade !== src.jobGrade) dest.jobGrade = src.jobGrade;
            if (dest.jobGradeLabel !== src.jobGradeLabel) dest.jobGradeLabel = src.jobGradeLabel;
            if (dest.dateofbirth !== src.dateofbirth) dest.dateofbirth = src.dateofbirth;
            if (dest.phoneNumber !== src.phoneNumber) dest.phoneNumber = src.phoneNumber;
            if (dest.profilePicture !== src.profilePicture) dest.profilePicture = src.profilePicture;
        };

        // Stream Management
        let currentStream: ServerStreamingCall<StreamRequest, StreamResponse> | undefined = undefined;

        /**
         * Handle a single stream response.
         */
        const handleStreamResponse = async (
            resp: StreamResponse,
            foundMarkers: number[],
            activeChar: UserLike | null,
            livemap: LivemapSettings,
        ) => {
            error.value = undefined;
            if (!resp || !resp.data) return;
            if (resp.userOnDuty !== undefined) userOnDuty.value = resp.userOnDuty;
            logger.debug('Received change - oneofKind:', resp.data.oneofKind, resp.data);

            switch (resp.data.oneofKind) {
                case 'jobs':
                    handleJobsUpdate(resp.data.jobs);
                    break;
                case 'markers':
                    handleMarkersUpdate(resp.data.markers, foundMarkers);
                    break;
                case 'snapshot':
                    handleSnapshotUpdate(resp.data.snapshot);
                    break;
                case 'userDeletes':
                    handleUserDeletes(resp.data.userDeletes);
                    break;
                case 'userUpdates':
                    handleUserUpdates(resp.data.userUpdates, livemap, activeChar);
                    break;
                default:
                    logger.warn('Unknown data received - oneofKind:' + resp.data.oneofKind);
            }

            if (!userOnDuty.value) {
                if (markersUsers.value.size > 0) markersUsers.value.clear();
                logger.info('User is not on duty, clearing user markers');
            }
        };

        /**
         * Handle jobs update from stream response.
         */
        const handleJobsUpdate = (jobs: JobsList) => {
            logger.info('Jobs received. Users:', jobs.users.length, 'markers:', jobs.markers.length);
            jobsMarkers.value = jobs.markers;
            jobsUsers.value = jobs.users;
        };

        /**
         * Handle markers update from stream response.
         */
        const handleMarkersUpdate = (markers: MarkerMarkersUpdates, foundMarkers: number[]) => {
            markers.updated.forEach((v: MarkerMarker) => {
                if (!markers.partial) {
                    foundMarkers.push(v.id);
                }
                addOrUpdateMarkerMarker(v);
            });
            markers.deleted.forEach((id: number) => markersMarkers.value.delete(id));
            if (!markers.partial && markers.part <= 0) {
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
        };

        /**
         * Handle snapshot update from stream response.
         */
        const handleSnapshotUpdate = (snapshot: Snapshot) => {
            markersUsers.value.clear();
            snapshot.markers.forEach((marker: UserMarker) => addOrUpdateUserMarker(marker));
            initiated.value = true;
        };

        /**
         * Handle user deletions from stream response.
         */
        const handleUserDeletes = (deletes: UserDeletes) => {
            for (const userDelete of deletes.deletes) {
                const userId = userDelete.id;
                const um = markersUsers.value.get(userId);
                if (!um || um.job !== userDelete.job) continue;
                markersUsers.value.delete(userId);
                if (selectedMarker.value?.userId === userId) {
                    selectedMarker.value = undefined;
                    followMarker.value = false;
                }
                if (ownMarker.value?.userId === userId) {
                    ownMarker.value = undefined;
                }
                logger.debug('User marker deleted:', userId);
            }
        };

        /**
         * Handle user updates from stream response.
         */
        const handleUserUpdates = (updates: UserUpdates, livemap: LivemapSettings, activeChar: UserLike | null) => {
            for (const userUpdate of updates.updates) {
                addOrUpdateUserMarker(userUpdate);
                if (livemap.centerSelectedMarker && userUpdate.userId === selectedMarker.value?.userId) {
                    selectedMarker.value = userUpdate;
                }
                if (activeChar?.userId === userUpdate.userId) {
                    ownMarker.value = userUpdate;
                }
            }
        };

        /**
         * Start the gRPC stream and handle responses.
         */
        const startStream = async (): Promise<void> => {
            if (abort.value !== undefined) return;

            stopping.value = false;
            logger.debug('Starting Stream');

            const { activeChar } = useAuth();
            const settingsStore = useSettingsStore();
            const { livemap } = storeToRefs(settingsStore);

            abort.value = new AbortController();
            error.value = undefined;

            const foundMarkers: number[] = [];
            cleanupMarkerMarkers();

            const livemapLivemapClient = await getLivemapLivemapClient();

            try {
                currentStream = livemapLivemapClient.stream({}, { abort: abort.value.signal });
                for await (const respRaw of currentStream.responses) {
                    const resp = respRaw as StreamResponse;
                    await handleStreamResponse(resp, foundMarkers, activeChar.value, livemap.value);
                }
            } catch (e) {
                const err = e as RpcError;
                error.value = undefined;
                if (err.code === 'CANCELLED' || err.code === 'ABORTED') {
                    if (!abort.value?.signal.aborted) {
                        await restartStream();
                    }
                    return;
                }
                logger.error('Stream failed', err.code, err.message, err.cause);
                if (!abort.value?.signal.aborted) {
                    await restartStream();
                } else {
                    error.value = err;
                }
            }
            logger.debug('Stream ended');
        };

        /**
         * Stop the gRPC stream.
         */
        const stopStream = async (end?: boolean): Promise<void> => {
            if (end === true) stopping.value = true;
            if (abort.value) {
                abort.value.abort();
                logger.debug('Stopping Stream');
            }
            abort.value = undefined;
        };

        /**
         * Restart the gRPC stream with backoff logic.
         */
        const restartStream = async (): Promise<void> => {
            if (!abort.value || abort.value.signal.aborted) return;
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

        // Misc Actions
        /**
         * Move to the given coordinates and optionally set in-game waypoint.
         */
        const gotoCoords = async (loc: Coordinate, ingame = true): Promise<void> => {
            location.value = { x: loc.x, y: loc.y };
            if (ingame) return setWaypoint(loc.x, loc.y);
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
            followMarker,

            // Actions
            startStream,
            stopStream,
            restartStream,
            addOrUpdateMarkerMarker,
            addOrUpdateUserMarker,
            updateUserMarker,
            updateUserInfo,
            deleteMarkerMarker,
            gotoCoords,
        };
    },
    {
        persist: false,
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useLivemapStore, import.meta.hot));
}
