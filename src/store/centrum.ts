import { RpcError } from '@protobuf-ts/runtime-rpc';
import { defineStore, type StoreDefinition } from 'pinia';
import { statusOrder } from '~/components/centrum/helpers';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { Dispatch, DispatchStatus, StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { CentrumMode, Settings } from '~~/gen/ts/resources/dispatch/settings';
import { StatusUnit, Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { UserShort } from '~~/gen/ts/resources/users/users';

const cleanupInterval = 1 * 40 * 1000;
const dispatchEndOfLifeTime = 2 * 60 * 60 * 1000;

// In seconds
const initialBackoffTime = 0.75;

export interface CentrumState {
    error: RpcError | undefined;
    abort: AbortController | undefined;
    cleanupIntervalId: NodeJS.Timeout | undefined;
    restarting: boolean;
    restartBackoffTime: number;

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

type canDoAction = 'TakeControl' | 'TakeDispatch' | 'AssignDispatch' | 'UpdateDispatchStatus' | 'UpdateUnitStatus';

export const useCentrumStore = defineStore('centrum', {
    state: () =>
        ({
            error: undefined,
            abort: undefined,
            restarting: false,
            restartBackoffTime: initialBackoffTime,

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
            return filtered
                .sort(
                    (a, b) =>
                        statusOrder.indexOf(b.status?.status ?? 0) -
                        statusOrder.indexOf(a.status?.status ?? 0) +
                        b.name.localeCompare(a.name),
                )
                .reverse();
        },
        getSortedDispatches: (state: CentrumState) => {
            return Array.from(state.dispatches, ([_, dsp]) => dsp).sort((a, b) => (a.id < b.id ? -1 : a.id > b.id ? 1 : 0));
        },
    },
    actions: {
        addFeedItem(item: DispatchStatus | UnitStatus): void {
            const idx = this.feed.findIndex((fi) => fi.id === item.id);
            if (idx === -1) {
                this.feed.unshift(item);
            }
        },
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
                console.error('Centrum: Processed Unit Status for unknown Unit', status.unitId);
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
        removeUnit(unit: Unit): void {
            // User's unit has been deleted, reset it
            if (this.ownUnitId !== undefined && this.ownUnitId === unit.id) {
                this.setOwnUnit(undefined);
            }

            this.units.delete(unit.id);
        },

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
                console.error('Centrum: Processed Dispatch Status for unknown Dispatch', status.dispatchId, status);
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

            const assignment = dispatch.units.find((u) => u.unitId === this.ownUnitId);
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
                // If dispatch is cancelled/completed, we can go ahead and remove it
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

                useNotificatorStore().dispatchNotification({
                    title: { key: 'notifications.centrum.store.assigned_dispatch.title', parameters: {} },
                    content: { key: 'notifications.centrum.store.assigned_dispatch.content', parameters: {} },
                    type: 'info',
                });
            }
        },

        removePendingDispatch(id: string): void {
            const idx = this.pendingDispatches.findIndex((d) => d === id);
            if (idx > -1) {
                this.pendingDispatches.splice(idx, 1);
            }
        },

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
            this.restarting = false;

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
                        if (resp.change.latestState.settings !== undefined) {
                            this.updateSettings(resp.change.latestState.settings);
                        }
                        this.disponents.length = 0;
                        this.disponents.push(...resp.change.latestState.disponents);
                        this.isDisponent = resp.change.latestState.isDisponent;
                        this.setOwnUnit(resp.change.latestState.ownUnit?.id);

                        // TODO handle removed units and dispatches
                        resp.change.latestState.units.forEach((u) => this.addOrUpdateUnit(u));
                        resp.change.latestState.dispatches.forEach((d) => this.addOrUpdateDispatch(d));
                    } else if (resp.change.oneofKind === 'settings') {
                        this.updateSettings(resp.change.settings);
                    } else if (resp.change.oneofKind === 'disponents') {
                        this.disponents.length = 0;
                        this.disponents.push(...resp.change.disponents.disponents);

                        // If user is not part of disponents list anymore
                        const idx = this.disponents.findIndex((d) => d.userId === authStore.activeChar?.userId);
                        if (idx === -1) {
                            this.isDisponent = false;
                            if (resp.restart !== undefined && !resp.restart) {
                                resp.restart = true;
                            }
                        }
                    } else if (resp.change.oneofKind === 'unitCreated') {
                        this.addOrUpdateUnit(resp.change.unitCreated);
                    } else if (resp.change.oneofKind === 'unitDeleted') {
                        this.removeUnit(resp.change.unitDeleted);
                    } else if (resp.change.oneofKind === 'unitUpdated') {
                        this.addOrUpdateUnit(resp.change.unitUpdated);
                        if (isCenter) {
                            continue;
                        }

                        // User added/in this unit
                        const idx = resp.change.unitUpdated.users.findIndex((u) => u.userId === authStore.activeChar?.userId);
                        if (idx > -1) {
                            // User already in unit
                            if (this.ownUnitId === resp.change.unitUpdated.id) {
                                continue;
                            }

                            this.setOwnUnit(resp.change.unitUpdated.id);

                            // User has been newly added to unit
                            notifications.dispatchNotification({
                                title: { key: 'notifications.centrum.unitUpdated.joined.title', parameters: {} },
                                content: { key: 'notifications.centrum.unitUpdated.joined.content', parameters: {} },
                                type: 'success',
                            });

                            this.dispatches.forEach((d) => this.handleDispatchAssignment(d));
                        } else {
                            if (this.ownUnitId === undefined) {
                                continue;
                            }

                            notifications.dispatchNotification({
                                title: { key: 'notifications.centrum.unitUpdated.removed.title', parameters: {} },
                                content: { key: 'notifications.centrum.unitUpdated.removed.content', parameters: {} },
                                type: 'warning',
                            });

                            // User has been removed from the unit
                            this.setOwnUnit(undefined);
                            this.ownDispatches.length = 0;
                            this.pendingDispatches.length = 0;
                        }
                    } else if (resp.change.oneofKind === 'unitStatus') {
                        this.addOrUpdateUnit(resp.change.unitStatus);

                        if (isCenter) {
                            continue;
                        }

                        if (this.isDisponent && resp.change.unitStatus.status) {
                            this.addFeedItem(resp.change.unitStatus.status);
                        }

                        // User added/in this unit
                        const idx = resp.change.unitStatus.users.findIndex((u) => u.userId === authStore.activeChar?.userId);
                        if (idx > -1) {
                            // User already in unit
                            if (this.ownUnitId === resp.change.unitStatus.id) {
                                continue;
                            }

                            this.setOwnUnit(resp.change.unitStatus.id);

                            // User has been newly added to unit
                            notifications.dispatchNotification({
                                title: { key: 'notifications.centrum.unitUpdated.joined.title', parameters: {} },
                                content: { key: 'notifications.centrum.unitUpdated.joined.content', parameters: {} },
                                type: 'success',
                            });

                            this.dispatches.forEach((d) => this.handleDispatchAssignment(d));
                        } else {
                            if (this.ownUnitId === undefined || this.ownUnitId !== resp.change.unitStatus.id) {
                                continue;
                            }

                            notifications.dispatchNotification({
                                title: { key: 'notifications.centrum.unitUpdated.removed.title', parameters: {} },
                                content: { key: 'notifications.centrum.unitUpdated.removed.content', parameters: {} },
                                type: 'warning',
                            });

                            // User has been removed from the unit
                            this.setOwnUnit(undefined);
                            this.ownDispatches.length = 0;
                            this.pendingDispatches.length = 0;
                        }
                    } else if (resp.change.oneofKind === 'dispatchCreated') {
                        this.addOrUpdateDispatch(resp.change.dispatchCreated);

                        if (resp.change.dispatchCreated.status !== undefined) {
                            this.addFeedItem(resp.change.dispatchCreated.status);
                        }
                    } else if (resp.change.oneofKind === 'dispatchDeleted') {
                        this.removeDispatch(resp.change.dispatchDeleted);
                    } else if (resp.change.oneofKind === 'dispatchUpdated') {
                        this.addOrUpdateDispatch(resp.change.dispatchUpdated);
                    } else if (resp.change.oneofKind === 'dispatchStatus') {
                        if (this.isDisponent && resp.change.dispatchStatus.status) {
                            this.addFeedItem(resp.change.dispatchStatus.status);
                        }

                        if (resp.change.dispatchStatus.status?.status === StatusDispatch.ARCHIVED) {
                            // If dispatch has been archived, remove from the main list
                            this.removeDispatch(resp.change.dispatchStatus.id);
                        } else {
                            this.addOrUpdateDispatch(resp.change.dispatchStatus);
                        }
                    } else {
                        console.warn('Centrum: Unknown change received - Kind: ', resp.change.oneofKind, resp.change);
                    }

                    if (resp.restart !== undefined && resp.restart) {
                        this.restartBackoffTime = initialBackoffTime;
                        this.restartStream(isCenter);
                        break;
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
            }

            if (!this.restarting) {
                if (this.cleanupIntervalId !== undefined) {
                    clearInterval(this.cleanupIntervalId);
                    this.cleanupIntervalId = undefined;
                }
            }

            console.debug('Centrum: Stopping Data Stream');
        },
        async restartStream(isCenter?: boolean): Promise<void> {
            this.restarting = true;

            // Reset back off time when over 10 seconds
            if (this.restartBackoffTime > 10) {
                this.restartBackoffTime = initialBackoffTime;
            } else {
                this.restartBackoffTime += initialBackoffTime;
            }

            console.debug('Centrum: Restart back off time in', this.restartBackoffTime, 'seconds');
            await this.stopStream();

            setTimeout(async () => {
                if (this.restarting) {
                    this.startStream(isCenter);
                }
            }, this.restartBackoffTime * 1000);
        },
        // Central "can user do that" method as we will take the dispatch center mode into account further
        canDo(action: canDoAction, dispatch?: Dispatch): boolean {
            // TODO check perms and dispatch center mode

            switch (action) {
                case 'TakeControl':
                    return can('CentrumService.TakeControl');

                case 'TakeDispatch':
                    return can('CentrumService.TakeDispatch') && this.getCurrentMode !== CentrumMode.CENTRAL_COMMAND;

                case 'AssignDispatch':
                    return can('CentrumService.AssignDispatch');

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
            const now = new Date().getTime();

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

            this.dispatches.forEach((d) => {
                // Remove dispatches older than 2 hours
                const endTime = now - toDate(d.status?.createdAt ?? d.createdAt).getTime();

                if (endTime >= dispatchEndOfLifeTime) {
                    this.removeDispatch(d.id);
                    return;
                }

                if (
                    d.status?.status !== StatusDispatch.COMPLETED &&
                    d.status?.status !== StatusDispatch.CANCELLED &&
                    d.status?.status !== StatusDispatch.ARCHIVED
                ) {
                    return;
                }

                // Remove completed/cancelled/archived dispatches after their status is consider "old"
                if (endTime >= cleanupInterval) {
                    this.removeDispatch(d.id);
                    return;
                }

                // Remove stale expired unit assignements
                d.units.forEach((ua, idx) => {
                    if (ua.expiresAt === undefined) {
                        return;
                    }

                    if (now - toDate(ua.expiresAt).getTime() >= cleanupInterval) {
                        d.units.splice(idx, 1);
                    }
                });
            });
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCentrumStore as unknown as StoreDefinition, import.meta.hot));
}
