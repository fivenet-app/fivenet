import { defineStore, type StoreDefinition } from 'pinia';
import { statusOrder } from '~/components/centrum/helpers';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { Dispatch, DispatchStatus, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { CentrumMode, Settings } from '~~/gen/ts/resources/centrum/settings';
import { StatusUnit, Unit, UnitStatus } from '~~/gen/ts/resources/centrum/units';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import { UserShort } from '~~/gen/ts/resources/users/users';

const cleanupInterval = 40 * 1000; // 40 seconds
const dispatchEndOfLifeTime = 2 * 60 * 60 * 1000; // 2 hours

// In seconds
const initialReconnectBackoffTime = 0.75;

export interface CentrumState {
    error: RpcError | undefined;
    abort: AbortController | undefined;
    cleanupIntervalId: NodeJS.Timeout | undefined;
    reconnecting: boolean;
    reconnectBackoffTime: number;

    timeCorrection: number;

    settings: Settings | undefined;
    isDisponent: boolean;
    disponents: UserShort[];
    feed: (DispatchStatus | UnitStatus)[];
    units: Map<string, Unit>;
    dispatches: Map<string, Dispatch>;
    ownUnitId: string | undefined;
    ownDispatches: string[];
    pendingDispatches: string[];
}

export type canDoAction = 'TakeControl' | 'TakeDispatch' | 'AssignDispatch' | 'UpdateDispatchStatus' | 'UpdateUnitStatus';

export const useCentrumStore = defineStore('centrum', {
    state: () =>
        ({
            error: undefined,
            abort: undefined,
            reconnecting: false,
            reconnectBackoffTime: initialReconnectBackoffTime,

            timeCorrection: 0,

            settings: undefined,
            isDisponent: false,
            disponents: [] as UserShort[],
            feed: [] as (DispatchStatus | UnitStatus)[],
            units: new Map<string, Unit>(),
            dispatches: new Map<string, Dispatch>(),
            ownUnitId: undefined,
            ownDispatches: [] as string[],
            pendingDispatches: [] as string[],
        }) as CentrumState,
    persist: false,
    getters: {
        getCurrentMode: (state: CentrumState) =>
            state.disponents.length > 0
                ? state.settings?.mode ?? CentrumMode.UNSPECIFIED
                : state.settings?.fallbackMode ?? CentrumMode.UNSPECIFIED,
        getOwnUnit: (state: CentrumState) => (state.ownUnitId !== undefined ? state.units.get(state.ownUnitId) : undefined),
        getSortedUnits: (state: CentrumState) => {
            const filtered: Unit[] = [];
            state.units.forEach((u) => filtered.push(u));
            return filtered.sort(
                (a, b) =>
                    a.name.localeCompare(b.name) -
                    statusOrder.indexOf(a.status?.status ?? 0) -
                    statusOrder.indexOf(b.status?.status ?? 0),
            );
        },
        getSortedDispatches: (state: CentrumState) => {
            return Array.from(state.dispatches, ([_, dsp]) => dsp).sort((a, b) => a.id.localeCompare(b.id));
        },
        getSortedOwnDispatches: (state: CentrumState) => {
            return state.ownDispatches.sort((a, b) => b.localeCompare(a));
        },
    },
    actions: {
        // General
        updateSettings(settings: Settings): void {
            if (this.settings !== undefined) {
                this.settings.enabled = settings.enabled;
                this.settings.job = settings.job;
                this.settings.mode = settings.mode;
                this.settings.fallbackMode = settings.fallbackMode;
            } else {
                this.settings = settings;
            }
        },

        // Units
        addOrUpdateUnit(unit: Unit): void {
            const u = this.units.get(unit.id);

            if (u === undefined) {
                if (unit.status === undefined) {
                    unit.status = {
                        unitId: unit.id,
                        id: '0',
                        status: StatusUnit.UNKNOWN,
                    };
                }
                this.units.set(unit.id, unit);
            } else {
                u.job = unit.job;
                u.createdAt = unit.createdAt;
                u.updatedAt = unit.updatedAt;
                u.name = unit.name;
                u.initials = unit.initials;
                u.color = unit.color;
                u.description = unit.description;

                if (unit.users.length === 0) {
                    u.users.length = 0;
                } else {
                    u.users.length = 0;
                    u.users.push(...unit.users);
                }

                this.updateUnitStatus(unit.status);
            }
        },
        updateUnitStatus(status: UnitStatus | undefined): void {
            if (status === undefined) {
                return;
            }

            const u = this.units.get(status.unitId);
            if (u === undefined) {
                console.warn('Centrum: Processed Unit Status for unknown Unit', status.unitId);
                return;
            }

            if (u.status === undefined) {
                u.status = status;
            } else {
                if (status.status === StatusUnit.USER_ADDED || status.status === StatusUnit.USER_REMOVED) {
                    return;
                }

                u.status!.id = status.id;
                u.status!.createdAt = status.createdAt;
                u.status!.unitId = status.unitId;
                u.status!.status = status.status;
                u.status!.reason = status.reason;
                u.status!.code = status.code;
                u.status!.userId = status.userId;
                u.status!.user = status.user;
                u.status!.x = status.x;
                u.status!.y = status.y;
                u.status!.postal = status.postal;
                u.status!.creator = status.creator;
                u.status!.creatorId = status.creatorId;
            }
        },
        setOwnUnit(id: string | undefined): void {
            if (id === undefined) {
                this.ownUnitId = undefined;
            } else {
                this.ownUnitId = id;
            }
        },
        removeUnit(id: string): void {
            // User's unit has been deleted, reset it
            if (this.ownUnitId === id) {
                this.setOwnUnit(undefined);
            }

            this.units.delete(id);
        },

        // Dispatches
        checkIfUnitAssignedToDispatch(dsp: Dispatch, unit: string | undefined): boolean {
            if (unit === undefined) return false;

            return dsp.units.findIndex((d) => d.unitId === unit) > -1;
        },
        addOrUpdateDispatch(dispatch: Dispatch): void {
            const d = this.dispatches.get(dispatch.id);

            if (d === undefined) {
                if (dispatch.status === undefined) {
                    dispatch.status = {
                        dispatchId: dispatch.id,
                        id: '0',
                        status: StatusDispatch.NEW,
                    };
                }
                this.dispatches.set(dispatch.id, dispatch);
            } else {
                d!.createdAt = dispatch.createdAt;
                d!.updatedAt = dispatch.updatedAt;
                d!.job = dispatch.job;
                d!.message = dispatch.message;
                d!.description = dispatch.description;
                d!.attributes = dispatch.attributes;
                d!.references = dispatch.references;
                d!.x = dispatch.x;
                d!.y = dispatch.y;
                d!.anon = dispatch.anon;
                d!.creatorId = dispatch.creatorId;
                d!.creator = dispatch.creator;

                if (dispatch.units.length === 0) {
                    d!.units.length = 0;
                } else {
                    d!.units.length = 0;
                    d!.units.push(...dispatch.units);
                }

                this.updateDispatchStatus(dispatch.status);
            }

            this.handleDispatchAssignment(dispatch);
        },
        updateDispatchStatus(status: DispatchStatus | undefined): void {
            if (status === undefined) {
                return;
            }

            const d = this.dispatches.get(status.dispatchId);
            if (d === undefined) {
                console.warn('Centrum: Processed Dispatch Status for unknown Dispatch', status.dispatchId, status);
                return;
            }

            if (d!.status === undefined) {
                d!.status = status;
            } else {
                d!.status!.id = status.id;
                d!.status!.createdAt = status.createdAt;
                d!.status!.dispatchId = status.dispatchId;
                d!.status!.unitId = status.unitId;
                d!.status!.unit = status.unit;
                d!.status!.status = status.status;
                d!.status!.reason = status.reason;
                d!.status!.code = status.code;
                d!.status!.userId = status.userId;
                d!.status!.user = status.user;
                d!.status!.x = status.x;
                d!.status!.y = status.y;
                d!.status!.postal = status.postal;

                // Make sure the unit is removed (quickly)
                if (d.status.status === StatusDispatch.UNIT_UNASSIGNED) {
                    const idx = d.units.findIndex((ua) => ua.unitId === status.unitId);
                    if (idx > -1) {
                        d.units.splice(idx, 1);
                    }
                }
            }
        },
        removeDispatch(id: string): void {
            this.removePendingDispatch(id);
            this.removeOwnDispatch(id);

            this.dispatches.delete(id);
        },
        addOrUpdateOwnDispatch(id: string): void {
            const idx = this.ownDispatches?.findIndex((d) => d === id) ?? -1;
            if (idx === -1) {
                this.ownDispatches.push(id);
            }
        },
        removeOwnDispatch(id: string): void {
            const idx = this.ownDispatches?.findIndex((d) => d === id) ?? -1;
            if (idx > -1) {
                this.ownDispatches?.splice(idx, 1);
            }
        },
        handleDispatchAssignment(dispatch: Dispatch): void {
            if (this.ownUnitId === undefined) {
                return;
            }

            const assignment = dispatch.units.find((ua) => ua.unitId === this.ownUnitId);
            if (assignment === undefined) {
                // If we don't have such a dispatch in our lists, probably not for us
                if (
                    this.pendingDispatches.find((d) => d === dispatch.id) === undefined &&
                    this.ownDispatches.find((d) => d === dispatch.id) === undefined
                ) {
                    return;
                }

                // Handle unassigment of dispatches
                this.removePendingDispatch(dispatch.id);
                this.removeOwnDispatch(dispatch.id);
            } else {
                // If dispatch is cancelled/completed, just remove it from our pending list
                if (
                    dispatch.status?.status === StatusDispatch.CANCELLED ||
                    dispatch.status?.status === StatusDispatch.COMPLETED
                ) {
                    this.removePendingDispatch(dispatch.id);
                    return;
                }

                // When dispatch has no expiration, it's an accepted/assigned dispatch
                if (assignment.expiresAt === undefined) {
                    this.removePendingDispatch(dispatch.id);
                    this.addOrUpdateOwnDispatch(dispatch.id);
                } else {
                    this.addOrUpdatePendingDispatch(dispatch.id);
                }
            }
        },

        addOrUpdatePendingDispatch(id: string): void {
            const idx = this.pendingDispatches?.findIndex((d) => d === id) ?? -1;
            if (idx === -1) {
                this.pendingDispatches.push(id);

                useNotificatorStore().add({
                    title: { key: 'notifications.centrum.store.assigned_dispatch.title', parameters: {} },
                    description: { key: 'notifications.centrum.store.assigned_dispatch.content', parameters: {} },
                    type: NotificationType.INFO,
                });

                useSound().play({ name: 'centrum/message-incoming' });
            }
        },
        removePendingDispatch(id: string): void {
            const idx = this.pendingDispatches.findIndex((d) => d === id);
            if (idx > -1) {
                this.pendingDispatches.splice(idx, 1);
            }
        },

        // disponents
        checkIfDisponent(userId?: number): boolean {
            return !!this.disponents.find((d) => d.userId === userId);
        },

        // Stream
        async startStream(isCenter?: boolean): Promise<void> {
            if (this.abort !== undefined) {
                return;
            }

            if (this.cleanupIntervalId === undefined) {
                this.cleanupIntervalId = setInterval(() => this.cleanup(), cleanupInterval);
            }

            console.debug('Centrum: Starting Data Stream');

            const authStore = useAuthStore();
            const notifications = useNotificatorStore();
            const { $grpc } = useNuxtApp();

            this.abort = new AbortController();
            this.error = undefined;
            this.reconnecting = false;

            try {
                const call = $grpc.getCentrumClient().stream(
                    {},
                    {
                        abort: this.abort.signal,
                    },
                );

                for await (const resp of call.responses) {
                    this.error = undefined;

                    if (resp === undefined || !resp.change) {
                        continue;
                    }

                    console.debug('Centrum: Received change - Kind:', resp.change.oneofKind, resp.change);

                    if (resp.change.oneofKind === 'latestState') {
                        if (resp.change.latestState.serverTime !== undefined) {
                            this.calculateTimeCorrection(resp.change.latestState.serverTime);
                        }

                        if (resp.change.latestState.settings !== undefined) {
                            this.updateSettings(resp.change.latestState.settings);
                        }
                        this.disponents.length = 0;
                        this.disponents.push(...resp.change.latestState.disponents);
                        this.isDisponent = this.checkIfDisponent(authStore.activeChar?.userId);

                        const foundUnits: string[] = [];
                        resp.change.latestState.units.forEach((u) => {
                            foundUnits.push(u.id);
                            this.addOrUpdateUnit(u);
                        });
                        // Remove units not found in latest state
                        let removedUnits = 0;
                        this.units.forEach((_, id) => {
                            if (!foundUnits.includes(id)) {
                                this.removeUnit(id);
                                removedUnits++;
                            }
                        });
                        console.debug(`Centrum: Removed ${removedUnits} old units`);
                        this.setOwnUnit(resp.change.latestState.ownUnitId);

                        const foundDispatches: string[] = [];
                        resp.change.latestState.dispatches.forEach((d) => {
                            foundDispatches.push(d.id);
                            this.addOrUpdateDispatch(d);
                        });
                        // Remove dispatches not found in latest state
                        let removedDispatches = 0;
                        this.dispatches.forEach((_, id) => {
                            if (!foundDispatches.includes(id)) {
                                this.removeDispatch(id);
                                removedDispatches++;
                            }
                        });
                        console.debug(`Centrum: Removed ${removedDispatches} old dispatches`);
                    } else if (resp.change.oneofKind === 'settings') {
                        this.updateSettings(resp.change.settings);
                    } else if (resp.change.oneofKind === 'disponents') {
                        this.disponents.length = 0;
                        this.disponents.push(...resp.change.disponents.disponents);

                        // If user is not part of disponents list anymore
                        this.isDisponent = this.checkIfDisponent(authStore.activeChar?.userId);
                        const idx = this.disponents.findIndex((d) => d.userId === authStore.activeChar?.userId);
                        if (idx > -1) {
                            this.isDisponent = true;
                        } else {
                            this.isDisponent = false;
                        }
                    } else if (resp.change.oneofKind === 'unitCreated') {
                        this.addOrUpdateUnit(resp.change.unitCreated);
                    } else if (resp.change.oneofKind === 'unitDeleted') {
                        this.removeUnit(resp.change.unitDeleted.id);
                    } else if (resp.change.oneofKind === 'unitUpdated') {
                        this.addOrUpdateUnit(resp.change.unitUpdated);

                        // User added/in this unit
                        const idx = resp.change.unitUpdated.users.findIndex((u) => u.userId === authStore.activeChar?.userId);
                        if (idx > -1) {
                            // User already in that unit
                            if (this.ownUnitId === resp.change.unitUpdated.id) {
                                continue;
                            }

                            this.setOwnUnit(resp.change.unitUpdated.id);

                            // User has been newly added to unit
                            notifications.add({
                                title: { key: 'notifications.centrum.unitUpdated.joined.title', parameters: {} },
                                description: { key: 'notifications.centrum.unitUpdated.joined.content', parameters: {} },
                                type: NotificationType.SUCCESS,
                            });

                            this.dispatches.forEach((d) => this.handleDispatchAssignment(d));
                        } else {
                            if (this.ownUnitId === undefined || this.ownUnitId !== resp.change.unitUpdated.id) {
                                continue;
                            }

                            notifications.add({
                                title: { key: 'notifications.centrum.unitUpdated.removed.title', parameters: {} },
                                description: { key: 'notifications.centrum.unitUpdated.removed.content', parameters: {} },
                                type: NotificationType.WARNING,
                            });

                            // User has been removed from the unit
                            this.setOwnUnit(undefined);
                            this.ownDispatches.length = 0;
                            this.pendingDispatches.length = 0;
                        }

                        if (isCenter && resp.change.unitUpdated.status !== undefined) {
                            this.addFeedItem(resp.change.unitUpdated.status);
                        }
                    } else if (resp.change.oneofKind === 'unitStatus') {
                        this.updateUnitStatus(resp.change.unitStatus);

                        if (isCenter) {
                            this.addFeedItem(resp.change.unitStatus);
                            continue;
                        }

                        // Check if the unit status is about us
                        if (resp.change.unitStatus.userId !== authStore.activeChar?.userId) {
                            continue;
                        }

                        if (resp.change.unitStatus.status === StatusUnit.USER_ADDED) {
                            // User already in unit
                            if (this.ownUnitId === resp.change.unitStatus.unitId) {
                                continue;
                            }

                            this.setOwnUnit(resp.change.unitStatus.unitId);

                            // User has been newly added to unit
                            notifications.add({
                                title: { key: 'notifications.centrum.unitUpdated.joined.title', parameters: {} },
                                description: { key: 'notifications.centrum.unitUpdated.joined.content', parameters: {} },
                                type: NotificationType.SUCCESS,
                            });

                            this.dispatches.forEach((d) => this.handleDispatchAssignment(d));
                        } else if (resp.change.unitStatus.status === StatusUnit.USER_REMOVED) {
                            if (this.ownUnitId === undefined || this.ownUnitId !== resp.change.unitStatus.unitId) {
                                continue;
                            }

                            notifications.add({
                                title: { key: 'notifications.centrum.unitUpdated.removed.title', parameters: {} },
                                description: { key: 'notifications.centrum.unitUpdated.removed.content', parameters: {} },
                                type: NotificationType.WARNING,
                            });

                            // User has been removed from the unit
                            this.setOwnUnit(undefined);
                            this.ownDispatches.length = 0;
                            this.pendingDispatches.length = 0;
                        }
                    } else if (resp.change.oneofKind === 'dispatchCreated') {
                        this.addOrUpdateDispatch(resp.change.dispatchCreated);

                        if (isCenter && resp.change.dispatchCreated.status !== undefined) {
                            this.addFeedItem(resp.change.dispatchCreated.status);
                        }
                    } else if (resp.change.oneofKind === 'dispatchDeleted') {
                        this.removeDispatch(resp.change.dispatchDeleted.id);
                    } else if (resp.change.oneofKind === 'dispatchUpdated') {
                        this.addOrUpdateDispatch(resp.change.dispatchUpdated);

                        if (isCenter && resp.change.dispatchUpdated.status !== undefined) {
                            this.addFeedItem(resp.change.dispatchUpdated.status);
                        }
                    } else if (resp.change.oneofKind === 'dispatchStatus') {
                        const status = resp.change.dispatchStatus;
                        this.updateDispatchStatus(status);

                        if (isCenter) {
                            this.addFeedItem(status);
                        }

                        if (status.status === StatusDispatch.ARCHIVED) {
                            // If dispatch has been archived, remove from the main list
                            this.removeDispatch(status.id);
                            continue;
                        } else if (status.status === StatusDispatch.NEED_ASSISTANCE) {
                            useSound().play({ name: 'centrum/morse-sos' });
                        }

                        // Only handle dispatch status affecting pending or own dispatches, when
                        // they are from our unit
                        if (this.ownUnitId === status.unitId) {
                            if (status.status === StatusDispatch.UNIT_ACCEPTED) {
                                this.removePendingDispatch(status.dispatchId);
                                this.addOrUpdateOwnDispatch(status.dispatchId);
                            } else if (
                                status.status === StatusDispatch.UNIT_DECLINED ||
                                status.status === StatusDispatch.UNIT_UNASSIGNED
                            ) {
                                const dispatch = this.dispatches.get(status.dispatchId);
                                if (dispatch === undefined) {
                                    continue;
                                }

                                this.removeDispatchAssignments(dispatch, status.unitId);
                            }
                        }
                    } else {
                        console.warn('Centrum: Unknown change received - Kind: ' + resp.change.oneofKind);
                    }
                }
            } catch (e) {
                const error = e as RpcError;
                if (error) {
                    // Only restart when not cancelled and abort is still valid
                    if (error.code !== 'CANCELLED' && error.code !== 'ABORTED') {
                        console.error('Centrum: Data Stream Failed', error.code, error.message, error.cause);

                        // Only set error if we don't need to restart
                        if (this.abort !== undefined && !this.abort?.signal.aborted) {
                            this.restartStream(isCenter);
                        } else {
                            this.error = error;
                        }
                    } else {
                        this.error = undefined;
                    }
                }
            }

            console.debug('Centrum: Data Stream Ended');
        },
        async stopStream(): Promise<void> {
            if (this.abort !== undefined) {
                this.abort.abort();
                this.abort = undefined;

                console.debug('Centrum: Stopping Data Stream');
            }

            if (!this.reconnecting) {
                if (this.cleanupIntervalId !== undefined) {
                    clearInterval(this.cleanupIntervalId);
                    this.cleanupIntervalId = undefined;
                }
            }
        },
        async restartStream(isCenter?: boolean): Promise<void> {
            this.reconnecting = true;

            // Reset back off time when over 10 seconds
            if (this.reconnectBackoffTime > 7) {
                this.reconnectBackoffTime = initialReconnectBackoffTime;
            } else {
                this.reconnectBackoffTime += initialReconnectBackoffTime;
            }

            console.debug('Centrum: Restart back off time in', this.reconnectBackoffTime, 'seconds');
            await this.stopStream();

            setTimeout(async () => {
                if (this.reconnecting) {
                    this.startStream(isCenter);
                }
            }, this.reconnectBackoffTime * 1000);
        },

        // Utilities
        addFeedItem(item: DispatchStatus | UnitStatus): void {
            const idx = this.feed.findIndex((fi) => fi.id === item.id);
            if (idx === -1) {
                this.feed.unshift(item);
            }
        },
        calculateTimeCorrection(serverTime: Timestamp): void {
            const now = new Date().getTime();
            const st = toDate(serverTime).getTime();

            const correction = now - st;
            this.timeCorrection = Math.floor(correction);

            console.debug(
                'Centrum: Calculated time correction - Now: ',
                now,
                'Server Time:',
                st,
                'Time Correction (seconds): ',
                this.timeCorrection / 1000,
            );
        },
        // Central "can user do that" method as we will take the dispatch center mode into account further
        canDo(action: canDoAction, dispatch?: Dispatch): boolean {
            switch (action) {
                case 'TakeControl':
                    return can('CentrumService.TakeControl');

                case 'TakeDispatch':
                    return can('CentrumService.TakeDispatch') && this.getCurrentMode !== CentrumMode.CENTRAL_COMMAND;

                case 'AssignDispatch':
                    return can('CentrumService.TakeControl');

                case 'UpdateDispatchStatus':
                    return (
                        can('CentrumService.TakeDispatch') &&
                        dispatch !== undefined &&
                        this.checkIfUnitAssignedToDispatch(dispatch, this.ownUnitId)
                    );

                case 'UpdateUnitStatus':
                    return can('CentrumService.TakeDispatch');

                default:
                    return false;
            }
        },
        async cleanup(): Promise<void> {
            console.debug('Centrum: Running cleanup tasks');
            const now = new Date().getTime() - this.timeCorrection;

            // Cleanup pending dispatches
            this.pendingDispatches.forEach((pd) => {
                if (!this.dispatches.has(pd)) {
                    this.removePendingDispatch(pd);
                } else {
                    this.dispatches.get(pd)?.units.forEach((ua) => {
                        if (now - toDate(ua.expiresAt).getTime() >= cleanupInterval) this.removePendingDispatch(pd);
                    });
                }
            });

            let count = 0;
            let skipped = 0;
            this.dispatches.forEach((dispatch) => {
                // Remove dispatches older than cleanup time
                const endTime = now - toDate(dispatch.status?.createdAt ?? dispatch.createdAt).getTime();

                if (endTime >= dispatchEndOfLifeTime) {
                    this.removeDispatch(dispatch.id);
                    count++;
                    return;
                }

                if (
                    dispatch.status?.status !== StatusDispatch.COMPLETED &&
                    dispatch.status?.status !== StatusDispatch.CANCELLED &&
                    dispatch.status?.status !== StatusDispatch.ARCHIVED
                ) {
                    skipped++;
                    return;
                }

                // Remove completed/cancelled/archived dispatches after their status is consider "old"
                if (endTime >= cleanupInterval) {
                    this.removeDispatch(dispatch.id);
                    count++;
                    return;
                }

                this.removeDispatchAssignments(dispatch);
                skipped++;
            });

            // Cut down feed if necessary to a maximum of 100 entries
            if (this.feed.length > 100) {
                this.feed.length = 100;
            }

            console.info('Centrum: Cleaned up dispatches, count:', count, 'skipped:', skipped);
        },

        removeDispatchAssignments(dispatch: Dispatch, unitId?: string): void {
            this.removeOwnDispatch(dispatch.id);
            this.removePendingDispatch(dispatch.id);

            const now = new Date().getTime() - this.timeCorrection;

            // Remove stale expired unit assignements
            dispatch.units.forEach((ua, idx) => {
                if (ua.unitId === unitId) {
                    dispatch.units.splice(idx, 1);
                    return;
                }

                if (ua.expiresAt === undefined) {
                    return;
                }

                if (now - toDate(ua.expiresAt).getTime() >= cleanupInterval) {
                    dispatch.units.splice(idx, 1);
                }
            });
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCentrumStore as unknown as StoreDefinition, import.meta.hot));
}
