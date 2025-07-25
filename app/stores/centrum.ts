import type { RpcError, ServerStreamingCall } from '@protobuf-ts/runtime-rpc';
import { defineStore } from 'pinia';
import { statusOrder } from '~/components/centrum/helpers';
import type { NotificationActionI18n } from '~/utils/notifications';
import type { Dispatchers } from '~~/gen/ts/resources/centrum/dispatchers';
import { type Dispatch, type DispatchStatus, StatusDispatch, TakeDispatchResp } from '~~/gen/ts/resources/centrum/dispatches';
import { type Settings, CentrumMode, CentrumType } from '~~/gen/ts/resources/centrum/settings';
import { type Unit, type UnitStatus, StatusUnit } from '~~/gen/ts/resources/centrum/units';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { JobAccess, StreamRequest, StreamResponse } from '~~/gen/ts/services/centrum/centrum';

export const logger = useLogger('⛑️ Centrum');

const cleanupInterval = 40 * 1000; // 40 seconds
const dispatchEndOfLifeTime = 2 * 60 * 60 * 1000; // 2 hours

// In seconds
const maxBackOffTime = 7;
const initialReconnectBackoffTime = 0.75;

export type canDoAction = 'TakeControl' | 'TakeDispatch' | 'AssignDispatch' | 'UpdateDispatchStatus' | 'UpdateUnitStatus';

export const useCentrumStore = defineStore(
    'centrum',
    () => {
        const { $grpc } = useNuxtApp();
        const notifications = useNotificationsStore();

        // State
        const error = ref<RpcError | undefined>(undefined);
        const abort = ref<AbortController | undefined>(undefined);
        const cleanupIntervalId = ref<ReturnType<typeof setInterval> | undefined>(undefined);
        const stopping = ref<boolean>(false);
        const reconnectBackoffTime = ref<number>(initialReconnectBackoffTime);

        const timeCorrection = ref<number>(0);

        const dispatchesJobs = ref<JobAccess | undefined>(undefined);
        const settings = ref<Settings | undefined>(undefined);
        const isDispatcher = ref<boolean>(false);
        const dispatchers = ref<Dispatchers[]>([]);
        const feed = ref<(DispatchStatus | UnitStatus)[]>([]);
        const isCenter = ref<boolean>(false);

        const units = ref<Map<number, Unit>>(new Map());
        const dispatches = ref<Map<number, Dispatch>>(new Map());

        const ownUnitId = ref<number | undefined>(undefined);
        const ownDispatches = ref<number[]>([]);
        const pendingDispatches = ref<number[]>([]);

        const messageIncomdingSound = useSounds('/sounds/centrum/message-incoming.mp3');
        const sosSound = useSounds('/sounds/centrum/morse-sos.mp3');

        // Helpers
        const removeOwnDispatch = (id: number): void => {
            const idx = ownDispatches.value.findIndex((d) => d === id);
            if (idx > -1) {
                ownDispatches.value.splice(idx, 1);
            }
        };

        const removePendingDispatch = (id: number): void => {
            const idx = pendingDispatches.value.findIndex((d) => d === id);
            if (idx > -1) {
                pendingDispatches.value.splice(idx, 1);
            }
        };

        // Getters
        const getCurrentMode = computed<CentrumMode>(() => {
            return dispatchers.value.length > 0
                ? (settings.value?.mode ?? CentrumMode.UNSPECIFIED)
                : (settings.value?.fallbackMode ?? CentrumMode.UNSPECIFIED);
        });

        const getJobDispatchers = computed((): Dispatchers | undefined => {
            const { activeChar } = useAuth();
            return dispatchers.value.find((d) => d.job === activeChar.value?.job);
        });

        const anyDispatchersActive = computed(() => {
            return (
                dispatchers.value !== undefined &&
                dispatchers.value.length > 0 &&
                dispatchers.value.map((d) => d.dispatchers.length).reduce((sum, d) => d + sum) > 0
            );
        });

        const getOwnUnit = computed<Unit | undefined>(() => {
            return ownUnitId.value !== undefined ? units.value.get(ownUnitId.value) : undefined;
        });

        const getSortedUnits = computed<Unit[]>(() => {
            const { activeChar } = useAuth();

            const array: Unit[] = [];
            units.value.forEach((u) => u.job === activeChar.value?.job && array.push(u));
            return array.sort(
                (a, b) =>
                    a.name.localeCompare(b.name) -
                    statusOrder.indexOf(a.status?.status ?? 0) -
                    statusOrder.indexOf(b.status?.status ?? 0),
            );
        });

        const getSortedDispatches = computed<Dispatch[]>(() => {
            return Array.from(dispatches.value, ([_, dsp]) => dsp).sort((a, b) => a.id - b.id);
        });

        const getSortedOwnDispatches = computed<number[]>(() => {
            // Sort descending
            return ownDispatches.value.sort((a, b) => b - a);
        });

        // Actions

        const setOrUpdateSettings = (newSettings: Settings): void => {
            if (settings.value !== undefined) {
                settings.value.enabled = newSettings.enabled;
                settings.value.job = newSettings.job;
                settings.value.mode = newSettings.mode;
                settings.value.fallbackMode = newSettings.fallbackMode;
                settings.value.timings = newSettings.timings;
            } else {
                settings.value = newSettings;
            }
        };

        // Units
        const addOrUpdateUnit = (unit: Unit): void => {
            const existing = units.value.get(unit.id);
            if (!existing) {
                if (!unit.access) {
                    unit.access = {
                        jobs: [],
                        qualifications: [],
                    };
                }
                if (!unit.status) {
                    unit.status = {
                        unitId: unit.id,
                        id: 0,
                        status: StatusUnit.UNKNOWN,
                    };
                }
                units.value.set(unit.id, unit);
            } else {
                existing.job = unit.job;
                existing.createdAt = unit.createdAt;
                existing.updatedAt = unit.updatedAt;
                existing.name = unit.name;
                existing.initials = unit.initials;
                existing.color = unit.color;
                existing.icon = unit.icon;
                existing.description = unit.description;

                existing.users.length = 0;
                if (unit.users.length > 0) {
                    existing.users.push(...unit.users);
                }

                if (unit.access) {
                    existing.access = unit.access;
                }

                updateUnitStatus(unit.status);
            }
        };

        const updateUnitStatus = (status: UnitStatus | undefined): void => {
            if (!status) {
                return;
            }
            const u = units.value.get(status.unitId);
            if (!u) {
                logger.warn('Processed Unit Status for unknown unit:', status.unitId);
                return;
            }
            if (!status.unit) {
                status.unit = u;
            }
            if (!u.status) {
                u.status = status;
            } else {
                // User added / removed
                if (status.status === StatusUnit.USER_ADDED || status.status === StatusUnit.USER_REMOVED) {
                    return;
                }

                // Normal status update
                u.status.id = status.id;
                u.status.createdAt = status.createdAt;
                u.status.unitId = status.unitId;
                u.status.status = status.status;
                u.status.reason = status.reason;
                u.status.code = status.code;
                u.status.userId = status.userId;
                u.status.user = status.user;
                u.status.x = status.x;
                u.status.y = status.y;
                u.status.postal = status.postal;
                u.status.creator = status.creator;
                u.status.creatorId = status.creatorId;
            }
        };

        const setOwnUnit = (id: number | undefined): void => {
            ownUnitId.value = id;
        };

        const removeUnit = (id: number): void => {
            if (ownUnitId.value === id) {
                setOwnUnit(undefined);
            }
            units.value.delete(id);
        };

        // Dispatches
        const checkIfUnitAssignedToDispatch = (dsp: Dispatch, unit: number | undefined): boolean => {
            if (!unit) return false;
            return dsp.units.findIndex((u) => u.unitId === unit) > -1;
        };

        const addOrUpdateDispatch = (dispatchObj: Dispatch): void => {
            const existing = dispatches.value.get(dispatchObj.id);
            if (!existing) {
                // Ensure the dispatch has a status
                if (!dispatchObj.status) {
                    dispatchObj.status = {
                        dispatchId: dispatchObj.id,
                        id: 0,
                        status: StatusDispatch.NEW,
                    };
                }
                dispatches.value.set(dispatchObj.id, dispatchObj);
            } else {
                existing.createdAt = dispatchObj.createdAt;
                existing.updatedAt = dispatchObj.updatedAt;
                existing.job = dispatchObj.job;
                existing.jobs = dispatchObj.jobs;
                existing.message = dispatchObj.message;
                existing.description = dispatchObj.description;
                existing.attributes = dispatchObj.attributes;
                existing.references = dispatchObj.references;
                existing.x = dispatchObj.x;
                existing.y = dispatchObj.y;
                existing.anon = dispatchObj.anon;
                existing.creatorId = dispatchObj.creatorId;
                existing.creator = dispatchObj.creator;

                existing.units.length = 0;
                if (dispatchObj.units.length > 0) {
                    existing.units.push(...dispatchObj.units);
                }

                updateDispatchStatus(dispatchObj.status);
            }
            handleDispatchAssignment(dispatchObj);
        };

        const updateDispatchStatus = (status: DispatchStatus | undefined): void => {
            if (!status) {
                return;
            }
            const disp = dispatches.value.get(status.dispatchId);
            if (!disp) {
                logger.warn('Processed Dispatch Status for unknown dispatch:', status.dispatchId, status);
                return;
            }
            if (!disp.status) {
                disp.status = status;
            } else {
                disp.status.id = status.id;
                disp.status.createdAt = status.createdAt;
                disp.status.dispatchId = status.dispatchId;
                disp.status.unitId = status.unitId;
                disp.status.unit = status.unit;
                disp.status.status = status.status;
                disp.status.reason = status.reason;
                disp.status.code = status.code;
                disp.status.userId = status.userId;
                disp.status.user = status.user;
                disp.status.x = status.x;
                disp.status.y = status.y;
                disp.status.postal = status.postal;

                // If unit got unassigned, remove it from the dispatch's units
                if (disp.status.status === StatusDispatch.UNIT_UNASSIGNED) {
                    const idx = disp.units.findIndex((ua) => ua.unitId === status.unitId);
                    if (idx > -1) {
                        disp.units.splice(idx, 1);
                    }
                }
            }
        };

        const removeDispatch = (id: number): void => {
            removePendingDispatch(id);
            removeOwnDispatch(id);
            dispatches.value.delete(id);
        };

        const addOrUpdateOwnDispatch = (id: number): void => {
            if (!ownDispatches.value.includes(id)) {
                ownDispatches.value.push(id);
            }
        };

        const handleDispatchAssignment = (dsp: Dispatch): void => {
            if (!ownUnitId.value) {
                return;
            }
            const assignment = dsp.units.find((ua) => ua.unitId === ownUnitId.value);
            if (!assignment) {
                // If not assigned, remove from pending/own if present
                if (!pendingDispatches.value.includes(dsp.id) && !ownDispatches.value.includes(dsp.id)) {
                    return;
                }
                removePendingDispatch(dsp.id);
                removeOwnDispatch(dsp.id);
            } else {
                // If dispatch is ended, remove from pending
                if (dsp.status?.status === StatusDispatch.CANCELLED || dsp.status?.status === StatusDispatch.COMPLETED) {
                    removePendingDispatch(dsp.id);
                    return;
                }
                // If no expiration => accepted
                if (!assignment.expiresAt) {
                    removePendingDispatch(dsp.id);
                    addOrUpdateOwnDispatch(dsp.id);
                } else {
                    // else => it's pending
                    addOrUpdatePendingDispatch(dsp.id);
                }
            }
        };

        const addOrUpdatePendingDispatch = (id: number): void => {
            if (!pendingDispatches.value.includes(id)) {
                pendingDispatches.value.push(id);
                notifications.add({
                    title: { key: 'notifications.centrum.store.assigned_dispatch.title', parameters: {} },
                    description: { key: 'notifications.centrum.store.assigned_dispatch.content', parameters: {} },
                    type: NotificationType.INFO,
                    actions: getNotificationActions(),
                });

                messageIncomdingSound.play();
            }
        };

        // Dispatchers
        const checkIfDispatcher = (userId?: number): boolean => {
            return !!dispatchers.value.find((d) => d.dispatchers.find((c) => c.userId === userId));
        };

        // Stream
        let currentStream: ServerStreamingCall<StreamRequest, StreamResponse> | undefined = undefined;

        const startStream = async (): Promise<void> => {
            if (abort.value !== undefined) return;
            stopping.value = false;

            logger.debug('Starting Stream');

            const { activeChar } = useAuth();

            abort.value = new AbortController();
            error.value = undefined;

            if (!cleanupIntervalId.value) {
                cleanupIntervalId.value = setInterval(() => cleanup(), cleanupInterval);
            }

            try {
                currentStream = $grpc.centrum.centrum.stream({}, { abort: abort.value.signal });

                for await (const respRaw of currentStream.responses) {
                    // The gRPC stream may yield unknown, so cast to the expected type
                    const resp = respRaw as StreamResponse;
                    error.value = undefined;

                    if (!resp || !resp.change) {
                        continue;
                    }

                    logger.debug('Received change - oneofKind:', resp.change.oneofKind, resp.change);

                    if (resp.change.oneofKind === 'handshake') {
                        if (resp.change.handshake.serverTime) {
                            calculateTimeCorrection(resp.change.handshake.serverTime);
                        }

                        dispatchesJobs.value = resp.change.handshake.jobAccess;

                        if (resp.change.handshake.settings) {
                            setOrUpdateSettings(resp.change.handshake.settings);
                        }
                    } else if (resp.change.oneofKind === 'latestState') {
                        logger.info(
                            'Latest state received. Dispatches:',
                            resp.change.latestState.dispatches.length,
                            'units:',
                            resp.change.latestState.units.length,
                        );

                        dispatchers.value.length = 0;
                        dispatchers.value.push(...(resp.change.latestState.dispatchers?.dispatchers ?? []));
                        isDispatcher.value = checkIfDispatcher(activeChar.value?.userId);

                        const foundUnits: number[] = [];
                        resp.change.latestState.units.forEach((u) => {
                            foundUnits.push(u.id);
                            addOrUpdateUnit(u);
                        });
                        // Remove missing units
                        let removedUnits = 0;
                        units.value.forEach((_, id) => {
                            if (!foundUnits.includes(id)) {
                                removeUnit(id);
                                removedUnits++;
                            }
                        });
                        logger.debug(`Removed ${removedUnits} old units`);
                        setOwnUnit(resp.change.latestState.ownUnitId);

                        const foundDispatches: number[] = [];
                        resp.change.latestState.dispatches.forEach((d) => {
                            foundDispatches.push(d.id);
                            addOrUpdateDispatch(d);
                        });
                        // Remove missing dispatches
                        let removedDispatches = 0;
                        dispatches.value.forEach((_, id) => {
                            if (!foundDispatches.includes(id)) {
                                removeDispatch(id);
                                removedDispatches++;
                            }
                        });
                        logger.debug(`Removed ${removedDispatches} old dispatches`);
                    } else if (resp.change.oneofKind === 'settings') {
                        setOrUpdateSettings(resp.change.settings);
                    } else if (resp.change.oneofKind === 'dispatchers') {
                        const idx = dispatchers.value.findIndex(
                            (d) => resp.change.oneofKind === 'dispatchers' && d.job === resp.change.dispatchers.job,
                        );
                        if (idx > -1) {
                            dispatchers.value[idx] = resp.change.dispatchers;
                        } else {
                            dispatchers.value.push(resp.change.dispatchers);
                        }

                        isDispatcher.value = checkIfDispatcher(activeChar.value?.userId);
                    } else if (resp.change.oneofKind === 'unitDeleted') {
                        removeUnit(resp.change.unitDeleted);
                    } else if (resp.change.oneofKind === 'unitUpdated') {
                        addOrUpdateUnit(resp.change.unitUpdated);

                        // Check if user is in that unit
                        const idx = resp.change.unitUpdated.users.findIndex((u) => u.userId === activeChar.value?.userId);
                        if (idx > -1) {
                            // Already in unit
                            if (ownUnitId.value === resp.change.unitUpdated.id) {
                                continue;
                            }

                            setOwnUnit(resp.change.unitUpdated.id);

                            notifications.add({
                                title: { key: 'notifications.centrum.unitUpdated.joined.title', parameters: {} },
                                description: { key: 'notifications.centrum.unitUpdated.joined.content', parameters: {} },
                                type: NotificationType.SUCCESS,
                                actions: getNotificationActions(),
                            });

                            dispatches.value.forEach((d) => handleDispatchAssignment(d));
                        } else {
                            // User was removed from that unit
                            if (ownUnitId.value === resp.change.unitUpdated.id) {
                                notifications.add({
                                    title: { key: 'notifications.centrum.unitUpdated.removed.title', parameters: {} },
                                    description: { key: 'notifications.centrum.unitUpdated.removed.content', parameters: {} },
                                    type: NotificationType.WARNING,
                                    actions: getNotificationActions(),
                                });

                                setOwnUnit(undefined);
                                ownDispatches.value.length = 0;
                                pendingDispatches.value.length = 0;
                            }
                        }

                        if (isCenter.value && resp.change.unitUpdated.status) {
                            addFeedItem(resp.change.unitUpdated.status);
                        }
                    } else if (resp.change.oneofKind === 'unitStatus') {
                        updateUnitStatus(resp.change.unitStatus);

                        if (isCenter.value) {
                            addFeedItem(resp.change.unitStatus);
                            continue;
                        }

                        // Check if the status is relevant to the current user
                        if (resp.change.unitStatus.userId !== activeChar.value?.userId) {
                            continue;
                        }

                        if (resp.change.unitStatus.status === StatusUnit.USER_ADDED) {
                            // User is in unit
                            if (ownUnitId.value === resp.change.unitStatus.unitId) {
                                continue;
                            }

                            setOwnUnit(resp.change.unitStatus.unitId);

                            notifications.add({
                                title: { key: 'notifications.centrum.unitUpdated.joined.title', parameters: {} },
                                description: { key: 'notifications.centrum.unitUpdated.joined.content', parameters: {} },
                                type: NotificationType.SUCCESS,
                                actions: getNotificationActions(),
                            });

                            dispatches.value.forEach((d) => handleDispatchAssignment(d));
                        } else if (resp.change.unitStatus.status === StatusUnit.USER_REMOVED) {
                            if (!ownUnitId.value || ownUnitId.value !== resp.change.unitStatus.unitId) {
                                continue;
                            }

                            notifications.add({
                                title: { key: 'notifications.centrum.unitUpdated.removed.title', parameters: {} },
                                description: { key: 'notifications.centrum.unitUpdated.removed.content', parameters: {} },
                                type: NotificationType.WARNING,
                                actions: getNotificationActions(),
                            });

                            setOwnUnit(undefined);
                            ownDispatches.value.length = 0;
                            pendingDispatches.value.length = 0;
                        }
                    } else if (resp.change.oneofKind === 'dispatchDeleted') {
                        removeDispatch(resp.change.dispatchDeleted);
                    } else if (resp.change.oneofKind === 'dispatchUpdated') {
                        addOrUpdateDispatch(resp.change.dispatchUpdated);

                        if (isCenter.value && resp.change.dispatchUpdated.status) {
                            addFeedItem(resp.change.dispatchUpdated.status);
                        }
                    } else if (resp.change.oneofKind === 'dispatchStatus') {
                        const ds = resp.change.dispatchStatus;
                        updateDispatchStatus(ds);

                        if (isCenter.value) {
                            addFeedItem(ds);
                        }

                        if (ds.status === StatusDispatch.ARCHIVED) {
                            removeDispatch(ds.id);
                            continue;
                        } else if (ds.status === StatusDispatch.NEED_ASSISTANCE) {
                            sosSound.play();
                        }

                        // If update is from the user
                        if (ownUnitId.value === ds.unitId) {
                            if (ds.status === StatusDispatch.UNIT_ACCEPTED) {
                                removePendingDispatch(ds.dispatchId);
                                addOrUpdateOwnDispatch(ds.dispatchId);
                            } else if (
                                ds.status === StatusDispatch.UNIT_DECLINED ||
                                ds.status === StatusDispatch.UNIT_UNASSIGNED
                            ) {
                                const d = dispatches.value.get(ds.dispatchId);
                                if (!d) {
                                    continue;
                                }
                                removeDispatchAssignments(d, ds.unitId);
                            }
                        }
                    } else {
                        logger.warn('Unknown change received - oneofKind:' + resp.change.oneofKind);
                    }
                }
            } catch (e) {
                const rpcError = e as RpcError;

                // Handle stream cancellation or abortion
                if (rpcError.code === 'CANCELLED' || rpcError.code === 'ABORTED') {
                    error.value = undefined;
                    if (!abort.value?.signal.aborted) {
                        await restartStream();
                    }
                    return;
                }

                // Handle specific disabled error
                if (rpcError.code === 'INVALID_ARGUMENT' && rpcError.message.includes('CentrumService.ErrDisabled')) {
                    settings.value = {
                        enabled: false,
                        type: CentrumType.DISPATCH,
                        mode: CentrumMode.UNSPECIFIED,
                        fallbackMode: CentrumMode.UNSPECIFIED,
                        job: '',
                        timings: undefined,
                        public: false,
                    };

                    notifications.add({
                        title: { key: 'notifications.centrum.disabled.title', parameters: {} },
                        description: { key: 'notifications.centrum.disabled.content', parameters: {} },
                        type: NotificationType.INFO,
                        actions: getNotificationActions(),
                    });
                    logger.info('Centrum is disabled for job, stopping stream.');
                    return;
                }

                // Log and handle other errors
                logger.error('Stream failed', rpcError.code, rpcError.message, rpcError.cause);

                if (abort.value && !abort.value.signal.aborted) {
                    await restartStream();
                } else {
                    error.value = rpcError;
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
            if (stopping.value && cleanupIntervalId.value) {
                clearInterval(cleanupIntervalId.value);
                cleanupIntervalId.value = undefined;
            }

            // Remove all dispatch layers from the settings
            const settingsStore = useSettingsStore();
            dispatchesJobs.value?.dispatches.forEach((job) => settingsStore.removeLivemapLayer(`dispatches_job_${job.job}`));

            abort.value = undefined;
        };

        const restartStream = async (): Promise<void> => {
            if (!abort.value || abort.value.signal.aborted) {
                return;
            }

            if (reconnectBackoffTime.value > maxBackOffTime) {
                reconnectBackoffTime.value = initialReconnectBackoffTime;
            } else {
                reconnectBackoffTime.value += initialReconnectBackoffTime;
            }

            logger.debug('Restart back off time in', reconnectBackoffTime.value, 'seconds');
            await stopStream();
            useTimeoutFn(() => {
                if (!stopping.value) {
                    startStream();
                }
            }, reconnectBackoffTime.value * 1000);
        };

        // Helpers
        const calculateTimeCorrection = (serverTime: Timestamp): void => {
            const now = new Date().getTime();
            const st = toDate(serverTime).getTime();

            const correction = now - st;
            timeCorrection.value = Math.floor(correction);

            logger.debug(
                'Calculated time correction - Now:',
                now,
                'Server Time:',
                st,
                'Time Correction (seconds):',
                timeCorrection.value / 1000,
            );
        };

        const addFeedItem = (item: DispatchStatus | UnitStatus): void => {
            const idx = feed.value.findIndex((fi) => fi.id === item.id);
            if (idx === -1) {
                feed.value.unshift(item);
            }
        };

        const canDo = (action: canDoAction, dispatchParam?: Dispatch): boolean => {
            const { can } = useAuth();

            switch (action) {
                case 'TakeControl':
                    return can('centrum.CentrumService/TakeControl').value;
                case 'TakeDispatch':
                    return (
                        can('centrum.CentrumService/TakeDispatch').value && getCurrentMode.value !== CentrumMode.CENTRAL_COMMAND
                    );
                case 'AssignDispatch':
                    return can('centrum.CentrumService/TakeControl').value;
                case 'UpdateDispatchStatus':
                    return (
                        dispatchParam !== undefined &&
                        can('centrum.CentrumService/TakeDispatch').value &&
                        checkIfUnitAssignedToDispatch(dispatchParam, ownUnitId.value)
                    );
                case 'UpdateUnitStatus':
                    return can('centrum.CentrumService/TakeDispatch').value;
                default:
                    return false;
            }
        };

        const cleanup = async (): Promise<void> => {
            logger.debug('Running cleanup tasks');
            const now = new Date().getTime() - timeCorrection.value;

            // Cleanup pending dispatches
            pendingDispatches.value.forEach((pd) => {
                if (!dispatches.value.has(pd)) {
                    removePendingDispatch(pd);
                } else {
                    const dsp = dispatches.value.get(pd);
                    dsp?.units.forEach((ua) => {
                        if (ua.expiresAt && now - toDate(ua.expiresAt).getTime() >= cleanupInterval) {
                            removePendingDispatch(pd);
                        }
                    });
                }
            });

            let count = 0;
            let skipped = 0;
            dispatches.value.forEach((d) => {
                const endTime = now - toDate(d.status?.createdAt ?? d.createdAt).getTime();
                if (endTime >= dispatchEndOfLifeTime) {
                    removeDispatch(d.id);
                    count++;
                    return;
                }

                if (
                    d.status?.status !== StatusDispatch.COMPLETED &&
                    d.status?.status !== StatusDispatch.CANCELLED &&
                    d.status?.status !== StatusDispatch.ARCHIVED
                ) {
                    skipped++;
                    return;
                }

                if (endTime >= cleanupInterval) {
                    removeDispatch(d.id);
                    count++;
                    return;
                }

                removeDispatchAssignments(d);
                skipped++;
            });

            if (feed.value.length > 100) {
                feed.value.length = 100;
            }

            logger.info('Cleaned up dispatches, count:', count, 'skipped:', skipped);
        };

        const removeDispatchAssignments = (dispatchObj: Dispatch, unitId?: number): void => {
            removeOwnDispatch(dispatchObj.id);
            removePendingDispatch(dispatchObj.id);

            const now = new Date().getTime() - timeCorrection.value;
            // If a particular unit was unassigned
            if (unitId) {
                dispatchObj.units = dispatchObj.units.filter((u) => u.unitId !== unitId);
                return;
            }

            // Remove stale pending assignments
            dispatchObj.units = dispatchObj.units.filter((ua) => {
                if (!ua.expiresAt) {
                    return true;
                }
                return now - toDate(ua.expiresAt).getTime() < cleanupInterval;
            });
        };

        const getNotificationActions = (): NotificationActionI18n[] => {
            const route = useRoute();
            if (route.name !== 'centrum' && route.name !== 'livemap') {
                return [
                    {
                        label: { key: 'common.click_here' },
                        to: '/livemap',
                    },
                ];
            }
            return [];
        };

        const selfAssign = async (id: number): Promise<void> => {
            if (ownUnitId.value === undefined) {
                notifications.add({
                    title: { key: 'notifications.centrum.unitUpdated.not_in_unit.title' },
                    description: { key: 'notifications.centrum.unitUpdated.not_in_unit.content' },
                    type: NotificationType.ERROR,
                });

                return;
            }

            try {
                const call = $grpc.centrum.centrum.takeDispatch({
                    dispatchIds: [id],
                    resp: TakeDispatchResp.ACCEPTED,
                });
                await call;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        const updateDispatchers = async (toRemove: number[]): Promise<void> => {
            try {
                const call = $grpc.centrum.centrum.updateDispatchers({
                    toRemove: toRemove,
                });
                await call;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        return {
            // State
            error,
            abort,
            cleanupIntervalId,
            stopping,
            reconnectBackoffTime,
            timeCorrection,
            settings,
            dispatchesJobs,
            isDispatcher,
            dispatchers,
            feed,
            units,
            dispatches,
            ownUnitId,
            ownDispatches,
            pendingDispatches,
            isCenter,

            // Getters
            getCurrentMode,
            getJobDispatchers,
            anyDispatchersActive,
            getOwnUnit,
            getSortedUnits,
            getSortedDispatches,
            getSortedOwnDispatches,

            // Actions
            setOrUpdateSettings,
            addOrUpdateUnit,
            updateUnitStatus,
            setOwnUnit,
            removeUnit,
            checkIfUnitAssignedToDispatch,
            addOrUpdateDispatch,
            updateDispatchStatus,
            removeDispatch,
            addOrUpdateOwnDispatch,
            handleDispatchAssignment,
            addOrUpdatePendingDispatch,
            removePendingDispatch,
            checkIfDispatcher,
            startStream,
            stopStream,
            restartStream,
            calculateTimeCorrection,
            addFeedItem,
            canDo,
            cleanup,
            removeDispatchAssignments,
            getNotificationActions,
            selfAssign,
            updateDispatchers,
        };
    },
    {
        persist: false,
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCentrumStore, import.meta.hot));
}
