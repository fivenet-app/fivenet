import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { StoreDefinition, defineStore } from 'pinia';
import { Dispatch, DispatchStatus, StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { CentrumMode, Settings } from '~~/gen/ts/resources/dispatch/settings';
import { StatusUnit, Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { useAuthStore } from './auth';
import { useNotificatorStore } from './notificator';

const FIVE_MINUTES = 7 * 60 * 1000;
const TWO_MINUTES = 2 * 60 * 1000;

export interface CentrumState {
    error: RpcError | undefined;
    abort: AbortController | undefined;
    cleanupIntervalId: NodeJS.Timeout | undefined;
    restarting: boolean;

    settings: Settings;
    isDisponent: boolean;
    disponents: UserShort[];
    feed: (DispatchStatus | UnitStatus)[];
    units: Map<bigint, Unit>;
    dispatches: Map<bigint, Dispatch>;
    ownUnitId: bigint | undefined;
    ownDispatches: bigint[];
    pendingDispatches: bigint[];
}

export const useCentrumStore = defineStore('centrum', {
    state: () =>
        ({
            error: undefined,
            abort: undefined,
            restarting: false,
            settings: {},
            isDisponent: false,
            disponents: [] as UserShort[],
            feed: [] as (DispatchStatus | UnitStatus)[],
            units: new Map<bigint, Unit>(),
            dispatches: new Map<bigint, Dispatch>(),
            ownUnitId: undefined,
            ownDispatches: [] as bigint[],
            pendingDispatches: [] as bigint[],
        }) as CentrumState,
    persist: false,
    actions: {
        addOrUpdateUnit(unit: Unit): void {
            const u = this.units.get(unit.id);
            if (u === undefined) {
                if (unit.status === undefined) {
                    unit.status = {
                        unitId: unit.id,
                        id: 0n,
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

                if (unit.status !== undefined) {
                    if (u.status === undefined) {
                        u.status = unit.status;
                    } else {
                        u.status!.id = unit.status.id;
                        u.status!.createdAt = unit.status.createdAt;
                        u.status!.unitId = unit.status.unitId;
                        u.status!.status = unit.status.status;
                        u.status!.reason = unit.status.reason;
                        u.status!.code = unit.status.code;
                        u.status!.userId = unit.status.userId;
                        u.status!.user = unit.status.user;
                        u.status!.x = unit.status.x;
                        u.status!.y = unit.status.y;
                        u.status!.postal = unit.status.postal;
                        u.status!.creator = unit.status.creator;
                        u.status!.creatorId = unit.status.creatorId;
                    }
                }
            }
        },
        setOwnUnit(id: bigint | undefined): void {
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

        checkIfUnitAssignedToDispatch(dsp: Dispatch, unit: bigint | undefined): boolean {
            if (unit === undefined) return false;

            return dsp.units.findIndex((d) => d.unitId === unit) > -1;
        },
        addOrUpdateDispatch(dispatch: Dispatch): void {
            if (!this.dispatches.has(dispatch.id)) {
                if (dispatch.status === undefined) {
                    dispatch.status = {
                        dispatchId: dispatch.id,
                        id: 0n,
                        status: StatusDispatch.NEW,
                    };
                }
                this.dispatches.set(dispatch.id, dispatch);
            } else {
                const d = this.dispatches.get(dispatch.id);
                d!.createdAt = dispatch.createdAt;
                d!.updatedAt = dispatch.updatedAt;
                d!.job = dispatch.job;
                d!.message = dispatch.message;
                d!.description = dispatch.description;
                d!.attributes = dispatch.attributes;
                d!.x = dispatch.x;
                d!.y = dispatch.y;
                d!.anon = dispatch.anon;
                d!.userId = dispatch.userId;
                d!.user = dispatch.user;

                if (dispatch.units.length === 0) {
                    d!.units.length = 0;
                } else {
                    d!.units.length = 0;
                    d!.units.concat(dispatch.units);
                }

                if (dispatch.status !== undefined) {
                    if (d!.status === undefined) {
                        d!.status = dispatch.status;
                    } else {
                        d!.status!.id = dispatch.status.id;
                        d!.status!.createdAt = dispatch.status.createdAt;
                        d!.status!.dispatchId = dispatch.status.dispatchId;
                        d!.status!.unitId = dispatch.status.unitId;
                        d!.status!.unit = dispatch.status.unit;
                        d!.status!.status = dispatch.status.status;
                        d!.status!.reason = dispatch.status.reason;
                        d!.status!.code = dispatch.status.code;
                        d!.status!.userId = dispatch.status.userId;
                        d!.status!.user = dispatch.status.user;
                        d!.status!.x = dispatch.status.x;
                        d!.status!.y = dispatch.status.y;
                        d!.status!.postal = dispatch.status.postal;
                    }
                }
            }

            this.handleDispatchAssignment(dispatch);
        },
        removeDispatch(id: bigint): void {
            this.removePendingDispatch(id);
            this.removeOwnDispatch(id);

            this.dispatches.delete(id);
        },
        addOrUpdateOwnDispatch(id: bigint): void {
            const idx = this.ownDispatches?.findIndex((d) => d === id) ?? -1;
            if (idx === -1) {
                this.ownDispatches.push(id);
            }
        },
        removeOwnDispatch(id: bigint): void {
            const idx = this.ownDispatches?.findIndex((d) => d === id) ?? -1;
            if (idx > -1) {
                this.ownDispatches?.splice(idx, 1);
            }
        },
        handleDispatchAssignment(dispatch: Dispatch): void {
            if (this.ownUnitId === undefined) return;

            if (
                dispatch.status?.status === StatusDispatch.UNIT_UNASSIGNED ||
                dispatch.status?.status === StatusDispatch.UNASSIGNED
            ) {
                // Handle unassigment of dispatches
                this.removePendingDispatch(dispatch.id);
                this.removeOwnDispatch(dispatch.id);
            } else {
                const assignment = dispatch.units.find((u) => u.unitId === this.ownUnitId);
                if (assignment === undefined) {
                    this.removePendingDispatch(dispatch.id);
                    this.removeOwnDispatch(dispatch.id);
                    return;
                }

                // When dispatch has expiration, it is a "pending" dispatch
                if (assignment?.expiresAt !== undefined) {
                    this.addOrUpdatePendingDispatch(dispatch.id);
                } else {
                    this.removePendingDispatch(dispatch.id);
                    this.addOrUpdateOwnDispatch(dispatch.id);
                }
            }
        },

        addOrUpdatePendingDispatch(id: bigint): void {
            const idx = this.pendingDispatches?.findIndex((d) => d === id) ?? -1;
            if (idx === -1) {
                this.pendingDispatches.push(id);

                useNotificatorStore().dispatchNotification({
                    title: { key: 'notifications.centrum.store.assigned_dispatch.title', parameters: [] },
                    content: { key: 'notifications.centrum.store.assigned_dispatch.content', parameters: [] },
                    type: 'info',
                });
            }
        },

        removePendingDispatch(id: bigint): void {
            const tDIdx = this.pendingDispatches.findIndex((d) => d === id);
            if (tDIdx > -1) {
                this.pendingDispatches.splice(tDIdx, 1);
            }
        },

        async startStream(isCenter?: boolean): Promise<void> {
            if (this.abort !== undefined) return;
            if (this.cleanupIntervalId !== undefined) clearInterval(this.cleanupIntervalId);
            this.cleanupIntervalId = setInterval(() => this.cleanup(), TWO_MINUTES);

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

                for await (let resp of call.responses) {
                    this.error = undefined;

                    if (resp === undefined || !resp.change) {
                        continue;
                    }

                    console.debug('Centrum: Received change - Kind:', resp.change.oneofKind, resp.change);

                    if (resp.change.oneofKind === 'latestState') {
                        if (resp.change.latestState.settings !== undefined) {
                            this.settings = resp.change.latestState.settings;
                        }
                        this.disponents = resp.change.latestState.disponents;
                        this.isDisponent = resp.change.latestState.isDisponent;
                        this.setOwnUnit(resp.change.latestState.ownUnit?.id);

                        resp.change.latestState.units.forEach((u) => this.addOrUpdateUnit(u));
                        resp.change.latestState.dispatches.forEach((d) => this.addOrUpdateDispatch(d));
                    } else if (resp.change.oneofKind === 'settings') {
                        this.settings = resp.change.settings;
                    } else if (resp.change.oneofKind === 'disponents') {
                        this.disponents = resp.change.disponents.disponents;
                        // If user is not part of disponents list anymore
                        const idx = this.disponents.findIndex((d) => d.userId === authStore.activeChar?.userId);
                        if (idx === -1) {
                            this.isDisponent = false;
                            if (resp.restart !== undefined && !resp.restart) resp.restart = true;
                        }
                    } else if (resp.change.oneofKind === 'unitAssigned') {
                        this.addOrUpdateUnit(resp.change.unitAssigned);
                        if (isCenter) continue;

                        // Ignore unit assignments for other units
                        if (this.ownUnitId !== undefined && this.ownUnitId !== resp.change.unitAssigned.id) {
                            continue;
                        }

                        // User added/in this unit
                        const idx = resp.change.unitAssigned.users.findIndex((u) => u.userId === authStore.activeChar?.userId);
                        if (idx > -1) {
                            // User already in unit
                            if (this.ownUnitId === resp.change.unitAssigned.id) continue;

                            this.setOwnUnit(resp.change.unitAssigned.id);

                            // User has been newly added to unit
                            notifications.dispatchNotification({
                                title: { key: 'notifications.centrum.unitAssigned.joined.title', parameters: [] },
                                content: { key: 'notifications.centrum.unitAssigned.joined.content', parameters: [] },
                                type: 'success',
                            });
                        } else {
                            if (this.ownUnitId === undefined) return;

                            // User has been removed from the unit
                            this.setOwnUnit(undefined);
                            this.ownDispatches.length = 0;
                            this.pendingDispatches.length = 0;

                            notifications.dispatchNotification({
                                title: { key: 'notifications.centrum.unitAssigned.removed.title', parameters: [] },
                                content: { key: 'notifications.centrum.unitAssigned.removed.content', parameters: [] },
                                type: 'warning',
                            });
                        }
                    } else if (resp.change.oneofKind === 'unitDeleted') {
                        this.removeUnit(resp.change.unitDeleted);
                    } else if (resp.change.oneofKind === 'unitUpdated') {
                        this.addOrUpdateUnit(resp.change.unitUpdated);
                    } else if (resp.change.oneofKind === 'unitStatus') {
                        this.addOrUpdateUnit(resp.change.unitStatus);

                        if (this.isDisponent && resp.change.unitStatus.status) {
                            this.feed.unshift(resp.change.unitStatus.status);
                        }
                    } else if (resp.change.oneofKind === 'dispatchCreated') {
                        this.addOrUpdateDispatch(resp.change.dispatchCreated);

                        if (resp.change.dispatchCreated.status !== undefined) {
                            this.feed.unshift(resp.change.dispatchCreated.status);
                        }
                    } else if (resp.change.oneofKind === 'dispatchDeleted') {
                        this.removeDispatch(resp.change.dispatchDeleted);
                    } else if (resp.change.oneofKind === 'dispatchUpdated') {
                        this.addOrUpdateDispatch(resp.change.dispatchUpdated);
                    } else if (resp.change.oneofKind === 'dispatchStatus') {
                        const id = resp.change.dispatchStatus.id;

                        if (this.isDisponent && resp.change.dispatchStatus.status) {
                            this.feed.unshift(resp.change.dispatchStatus.status);
                        }

                        if (resp.change.dispatchStatus.status?.status === StatusDispatch.ARCHIVED) {
                            // If dispatch has been archived, remove from the main list
                            this.removeDispatch(id);
                        } else {
                            this.addOrUpdateDispatch(resp.change.dispatchStatus);
                        }
                    } else if (resp.change.oneofKind === 'ping') {
                        console.debug('Centrum: Ping received');
                    } else {
                        console.warn('Centrum: Unknown change received - Kind: ', resp.change.oneofKind, resp.change);
                    }

                    if (resp.restart !== undefined && resp.restart) {
                        this.restartStream(isCenter);
                        break;
                    }
                }
            } catch (e) {
                this.error = e as RpcError;
                if (this.error) {
                    // Only restart when not cancelled and abort is still valid
                    if (this.error.code != 'CANCELLED') {
                        console.error('Centrum: Data Stream Failed', this.error.code, this.error.message, this.error.cause);

                        if (this.abort !== undefined && !this.abort?.signal.aborted) {
                            this.restartStream(isCenter);
                        }
                    }
                }
            }

            console.debug('Centrum: Data Stream Ended');
        },
        async stopStream(): Promise<void> {
            if (this.abort !== undefined) this.abort.abort();
            this.abort = undefined;
            if (this.cleanupIntervalId !== undefined) clearInterval(this.cleanupIntervalId);
            this.cleanupIntervalId = undefined;

            console.debug('Centrum: Stopping Data Stream');
        },
        async restartStream(isCenter?: boolean): Promise<void> {
            this.restarting = true;
            console.debug('Centrum: Restarting Data Stream');
            await this.stopStream();

            setTimeout(async () => this.startStream(isCenter), 1000);
        },
        // Central "can user do that" method as we will take the dispatch center mode into account further
        canDo(action: canDoAction, dispatch?: Dispatch): boolean {
            // TODO check perms and dispatch center mode

            switch (action) {
                case 'TakeControl':
                    return can('CentrumService.TakeControl');

                case 'TakeDispatch':
                    return can('CentrumService.TakeDispatch') && this.settings.mode !== CentrumMode.CENTRAL_COMMAND;

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
            this.pendingDispatches.forEach((pd, index) => {
                if (!this.dispatches.has(pd)) {
                    this.pendingDispatches.splice(index, 1);
                } else {
                    this.dispatches.get(pd)?.units.forEach((ua) => {
                        const expiresAt = toDate(ua.expiresAt);
                        if (now - expiresAt.getTime() > FIVE_MINUTES) this.removePendingDispatch(pd);
                    });
                }
            });

            // Remove completed, cancelled and archived dispatches after the status is 5 minutes or older
            this.dispatches.forEach((d) => {
                if (
                    d.status?.status !== StatusDispatch.COMPLETED &&
                    d.status?.status !== StatusDispatch.CANCELLED &&
                    d.status?.status !== StatusDispatch.ARCHIVED
                )
                    return;

                if (now - toDate(d.status?.createdAt).getTime() > FIVE_MINUTES) this.removeDispatch(d.id);
            });
        },
    },
});

type canDoAction = 'TakeControl' | 'TakeDispatch' | 'AssignDispatch' | 'UpdateDispatchStatus' | 'UpdateUnitStatus';

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCentrumStore as unknown as StoreDefinition, import.meta.hot));
}
