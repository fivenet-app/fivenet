import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { StoreDefinition, defineStore } from 'pinia';
import { DISPATCH_STATUS, Dispatch, DispatchStatus } from '~~/gen/ts/resources/dispatch/dispatches';
import { CENTRUM_MODE, Settings } from '~~/gen/ts/resources/dispatch/settings';
import { UNIT_STATUS, Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { useAuthStore } from './auth';
import { useNotificationsStore } from './notifications';

const SEVEN_MINUTES = 7 * 60 * 1000;
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
    units: Unit[];
    dispatches: Dispatch[];
    ownUnit: Unit | undefined;
    ownDispatches: Dispatch[];
    pendingDispatches: Dispatch[];
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
            units: [] as Unit[],
            dispatches: [] as Dispatch[],
            ownUnit: undefined,
            ownDispatches: [] as Dispatch[],
            pendingDispatches: [] as Dispatch[],
        }) as CentrumState,
    persist: false,
    actions: {
        addOrUpdateUnit(unit: Unit): void {
            const idx = this.units.findIndex((d) => d.id === unit.id) ?? -1;
            if (idx === -1) {
                if (unit.status === undefined) {
                    unit.status = {
                        unitId: unit.id,
                        id: 0n,
                        status: UNIT_STATUS.UNKNOWN,
                    };
                }
                this.units.push(unit);
            } else {
                this.units[idx].job = unit.job;
                this.units[idx].createdAt = unit.createdAt;
                this.units[idx].updatedAt = unit.updatedAt;
                this.units[idx].name = unit.name;
                this.units[idx].initials = unit.initials;
                this.units[idx].color = unit.color;
                this.units[idx].description = unit.description;
                this.units[idx].users = unit.users;

                if (unit.users.length == 0) {
                    this.units[idx].users.length = 0;
                } else {
                    this.units[idx].users = unit.users;
                }

                if (unit.status !== undefined) {
                    if (this.units[idx].status === undefined) {
                        this.units[idx].status = unit.status;
                    } else {
                        this.units[idx].status!.id = unit.status.id;
                        this.units[idx].status!.createdAt = unit.status.createdAt;
                        this.units[idx].status!.unitId = unit.status.unitId;
                        this.units[idx].status!.status = unit.status.status;
                        this.units[idx].status!.reason = unit.status.reason;
                        this.units[idx].status!.code = unit.status.code;
                        this.units[idx].status!.userId = unit.status.userId;
                        this.units[idx].status!.user = unit.status.user;
                        this.units[idx].status!.x = unit.status.x;
                        this.units[idx].status!.y = unit.status.y;
                        this.units[idx].status!.creator = unit.status.creator;
                        this.units[idx].status!.creatorId = unit.status.creatorId;
                    }
                }
            }

            this.updateOwnUnit(unit);
        },
        updateOwnUnit(unit: Unit): void {
            if (unit.id !== this.ownUnit?.id) return;

            this.ownUnit.id = unit.id;
            this.ownUnit.createdAt = unit.createdAt;
            this.ownUnit.updatedAt = unit.updatedAt;
            this.ownUnit.job = unit.job;
            this.ownUnit.name = unit.name;
            this.ownUnit.initials = unit.initials;
            this.ownUnit.color = unit.color;
            this.ownUnit.description = unit.description;
            this.ownUnit.status = unit.status;
            this.ownUnit.users = unit.users;

            if (unit.users.length == 0) {
                this.ownUnit.users.length = 0;
            } else {
                this.ownUnit.users = unit.users;
            }

            if (unit.status !== undefined) {
                if (this.ownUnit.status === undefined) {
                    this.ownUnit.status = unit.status;
                } else {
                    this.ownUnit.status!.id = unit.status.id;
                    this.ownUnit.status!.createdAt = unit.status.createdAt;
                    this.ownUnit.status!.unitId = unit.status.unitId;
                    this.ownUnit.status!.status = unit.status.status;
                    this.ownUnit.status!.reason = unit.status.reason;
                    this.ownUnit.status!.code = unit.status.code;
                    this.ownUnit.status!.userId = unit.status.userId;
                    this.ownUnit.status!.user = unit.status.user;
                    this.ownUnit.status!.x = unit.status.x;
                    this.ownUnit.status!.y = unit.status.y;
                    this.ownUnit.status!.creator = unit.status.creator;
                    this.ownUnit.status!.creatorId = unit.status.creatorId;
                }
            }
        },
        removeUnit(unit: Unit): void {
            const idx = this.units?.findIndex((d) => d.id === unit.id) ?? -1;
            if (idx > -1) {
                this.units?.splice(idx, 1);
            }

            // User's unit has been deleted, reset it
            if (this.ownUnit !== undefined && this.ownUnit.id === unit.id) {
                this.ownUnit = undefined;
            }
        },

        checkIfUnitAssignedToDispatch(dsp: Dispatch, unit?: Unit): boolean {
            if (unit === undefined) return false;

            return dsp.units.findIndex((d) => d.unitId === unit.id) > -1;
        },
        addOrUpdateDispatch(dispatch: Dispatch): void {
            const idx = this.dispatches?.findIndex((d) => d.id === dispatch.id) ?? -1;
            if (idx === -1) {
                if (dispatch.status === undefined) {
                    dispatch.status = {
                        dispatchId: dispatch.id,
                        id: 0n,
                        status: DISPATCH_STATUS.NEW,
                    };
                }
                this.dispatches.push(dispatch);
            } else {
                this.dispatches[idx].createdAt = dispatch.createdAt;
                this.dispatches[idx].updatedAt = dispatch.updatedAt;
                this.dispatches[idx].job = dispatch.job;
                this.dispatches[idx].message = dispatch.message;
                this.dispatches[idx].description = dispatch.description;
                this.dispatches[idx].attributes = dispatch.attributes;
                this.dispatches[idx].x = dispatch.x;
                this.dispatches[idx].y = dispatch.y;
                this.dispatches[idx].anon = dispatch.anon;
                this.dispatches[idx].userId = dispatch.userId;
                this.dispatches[idx].user = dispatch.user;

                if (dispatch.units.length == 0) {
                    this.dispatches[idx].units.length = 0;
                } else {
                    this.dispatches[idx].units = dispatch.units;
                }

                if (dispatch.status !== undefined) {
                    if (this.dispatches[idx].status === undefined) {
                        this.dispatches[idx].status = dispatch.status;
                    } else {
                        this.dispatches[idx].status!.id = dispatch.status.id;
                        this.dispatches[idx].status!.createdAt = dispatch.status.createdAt;
                        this.dispatches[idx].status!.dispatchId = dispatch.status.dispatchId;
                        this.dispatches[idx].status!.unitId = dispatch.status.unitId;
                        this.dispatches[idx].status!.unit = dispatch.status.unit;
                        this.dispatches[idx].status!.status = dispatch.status.status;
                        this.dispatches[idx].status!.reason = dispatch.status.reason;
                        this.dispatches[idx].status!.code = dispatch.status.code;
                        this.dispatches[idx].status!.userId = dispatch.status.userId;
                        this.dispatches[idx].status!.user = dispatch.status.user;
                        this.dispatches[idx].status!.x = dispatch.status.x;
                        this.dispatches[idx].status!.y = dispatch.status.y;
                    }
                }
            }

            this.handleDispatchAssignment(dispatch);
        },
        removeDispatch(id: bigint): void {
            const idx = this.dispatches?.findIndex((d) => d.id === id) ?? -1;
            if (idx > -1) {
                this.dispatches?.splice(idx, 1);
            }

            this.removePendingDispatch(id);
            this.removeOwnDispatch(id);
        },
        addOrUpdateOwnDispatch(dispatch: Dispatch): void {
            const idx = this.ownDispatches?.findIndex((d) => d.id === dispatch.id) ?? -1;
            if (idx === -1) {
                this.ownDispatches.push(dispatch);
            } else {
                this.ownDispatches[idx].createdAt = dispatch.createdAt;
                this.ownDispatches[idx].updatedAt = dispatch.updatedAt;
                this.ownDispatches[idx].job = dispatch.job;
                this.ownDispatches[idx].status = dispatch.status;
                this.ownDispatches[idx].message = dispatch.message;
                this.ownDispatches[idx].description = dispatch.description;
                this.ownDispatches[idx].attributes = dispatch.attributes;
                this.ownDispatches[idx].x = dispatch.x;
                this.ownDispatches[idx].y = dispatch.y;
                this.ownDispatches[idx].anon = dispatch.anon;
                this.ownDispatches[idx].userId = dispatch.userId;
                this.ownDispatches[idx].user = dispatch.user;
                if (dispatch.units.length == 0) {
                    this.ownDispatches[idx].units.length = 0;
                } else {
                    this.ownDispatches[idx].units = dispatch.units;
                }
            }
        },
        removeOwnDispatch(id: bigint): void {
            const idx = this.ownDispatches?.findIndex((d) => d.id === id) ?? -1;
            if (idx > -1) {
                this.ownDispatches?.splice(idx, 1);
            }
        },
        handleDispatchAssignment(dispatch: Dispatch): void {
            if (this.ownUnit === undefined) return;

            if (
                dispatch.status?.status === DISPATCH_STATUS.UNIT_UNASSIGNED ||
                dispatch.status?.status === DISPATCH_STATUS.UNASSIGNED
            ) {
                // Handle unassigment of dispatches
                this.removePendingDispatch(dispatch.id);
                this.removeOwnDispatch(dispatch.id);
            } else {
                const assignment = dispatch.units.find((u) => u.unitId === this.ownUnit?.id);
                if (assignment === undefined) {
                    this.removePendingDispatch(dispatch.id);
                    this.removeOwnDispatch(dispatch.id);
                    return;
                }

                // When dispatch has expiration, it is a "pending" dispatch
                if (assignment?.expiresAt !== undefined) {
                    this.addOrUpdatePendingDispatch(dispatch);
                } else {
                    this.removePendingDispatch(dispatch.id);
                    this.addOrUpdateOwnDispatch(dispatch);
                }
            }
        },

        addOrUpdatePendingDispatch(dispatch: Dispatch): void {
            const idx = this.pendingDispatches?.findIndex((d) => d.id === dispatch.id) ?? -1;
            if (idx === -1) {
                this.pendingDispatches.push(dispatch);

                useNotificationsStore().dispatchNotification({
                    title: { key: 'notifications.centrum.store.assigned_dispatch.title', parameters: [] },
                    content: { key: 'notifications.centrum.store.assigned_dispatch.content', parameters: [] },
                    type: 'info',
                });
            }
        },

        removePendingDispatch(id: bigint): void {
            const tDIdx = this.pendingDispatches.findIndex((d) => d.id === id);
            if (tDIdx > -1) {
                this.pendingDispatches.splice(tDIdx, 1);
            }
        },

        async startStream(): Promise<void> {
            if (this.abort !== undefined) return;
            this.restarting = false;
            if (this.cleanupIntervalId !== undefined) clearInterval(this.cleanupIntervalId);
            this.cleanupIntervalId = setInterval(() => this.cleanup(), TWO_MINUTES);

            console.debug('Centrum: Starting Data Stream');

            const authStore = useAuthStore();
            const notifications = useNotificationsStore();

            try {
                this.abort = new AbortController();

                const { $grpc } = useNuxtApp();
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
                        this.ownUnit = resp.change.latestState.ownUnit;

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

                            this.restartStream();
                            break;
                        }
                    } else if (resp.change.oneofKind === 'unitAssigned') {
                        if (this.ownUnit !== undefined && resp.change.unitAssigned.id !== this.ownUnit?.id) {
                            console.warn('Received unit user assigned event for other unit'), resp.change.unitAssigned;
                            continue;
                        }

                        const idx = resp.change.unitAssigned.users.findIndex((u) => u.userId === authStore.activeChar?.userId);
                        if (idx === -1) {
                            // User has been removed from the unit
                            this.ownUnit = undefined;
                            this.ownDispatches.length = 0;
                            this.pendingDispatches.length = 0;

                            notifications.dispatchNotification({
                                title: { key: 'notifications.centrum.unitAssigned.removed.title', parameters: [] },
                                content: { key: 'notifications.centrum.unitAssigned.removed.content', parameters: [] },
                                type: 'success',
                            });
                        } else if (this.ownUnit !== undefined) {
                            this.updateOwnUnit(resp.change.unitAssigned);

                            if (this.ownUnit.id === resp.change.unitAssigned.id) return;

                            // User has been newly added to unit
                            notifications.dispatchNotification({
                                title: { key: 'notifications.centrum.unitAssigned.joined.title', parameters: [] },
                                content: { key: 'notifications.centrum.unitAssigned.joined.content', parameters: [] },
                                type: 'success',
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
                    } else if (resp.change.oneofKind === 'dispatchDeleted') {
                        this.removeDispatch(resp.change.dispatchDeleted.id);
                    } else if (resp.change.oneofKind === 'dispatchUpdated') {
                        this.addOrUpdateDispatch(resp.change.dispatchUpdated);
                    } else if (resp.change.oneofKind === 'dispatchStatus') {
                        const id = resp.change.dispatchStatus.id;

                        if (this.isDisponent && resp.change.dispatchStatus.status) {
                            this.feed.unshift(resp.change.dispatchStatus.status);
                        }

                        if (resp.change.dispatchStatus.status?.status === DISPATCH_STATUS.ARCHIVED) {
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
                        this.restartStream();
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
                            this.restartStream();
                        }
                    }
                }
            }

            console.debug('Centrum: Data Stream Ended');
        },
        async stopStream(): Promise<void> {
            if (this.abort !== undefined) this.abort.abort();
            this.abort = undefined;
            if (this.cleanupIntervalId) clearInterval(this.cleanupIntervalId);
            this.cleanupIntervalId = undefined;
            console.debug('Centrum: Stopping Data Stream');
        },
        async restartStream(): Promise<void> {
            this.restarting = true;
            console.debug('Centrum: Restarting Data Stream');
            await this.stopStream();

            setTimeout(async () => this.startStream(), 1000);
        },
        // Central "can user do that" method as we will take the dispatch center mode into account further
        canDo(action: canDoAction, dispatch?: Dispatch): boolean {
            // TODO check perms and dispatch center mode

            switch (action) {
                case 'TakeControl':
                    return can('CentrumService.TakeControl');

                case 'TakeDispatch':
                    return can('CentrumService.TakeDispatch') && this.settings.mode !== CENTRUM_MODE.CENTRAL_COMMAND;

                case 'AssignDispatch':
                    return can('CentrumService.AssignDispatch');

                case 'UpdateDispatchStatus':
                    return (
                        can('CentrumService.TakeDispatch') &&
                        dispatch !== undefined &&
                        this.checkIfUnitAssignedToDispatch(dispatch, this.ownUnit)
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
                pd.units.forEach((ua) => {
                    const expiresAt = toDate(ua.expiresAt);
                    if (now - expiresAt.getTime() > SEVEN_MINUTES) this.removePendingDispatch(pd.id);
                });
            });

            // Remove completed, cancelled and archived dispatches after the status is 5 minutes or older
            this.dispatches.forEach((d) => {
                if (
                    d.status?.status !== DISPATCH_STATUS.COMPLETED &&
                    d.status?.status !== DISPATCH_STATUS.CANCELLED &&
                    d.status?.status !== DISPATCH_STATUS.ARCHIVED
                )
                    return;

                if (now - toDate(d.status?.createdAt).getTime() > SEVEN_MINUTES) this.removeDispatch(d.id);
            });
        },
    },
});

type canDoAction = 'TakeControl' | 'TakeDispatch' | 'AssignDispatch' | 'UpdateDispatchStatus' | 'UpdateUnitStatus';

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCentrumStore as unknown as StoreDefinition, import.meta.hot));
}
