import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { StoreDefinition, defineStore } from 'pinia';
import { DISPATCH_STATUS, Dispatch, DispatchStatus } from '~~/gen/ts/resources/dispatch/dispatches';
import { Settings } from '~~/gen/ts/resources/dispatch/settings';
import { Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { useAuthStore } from './auth';
import { useNotificationsStore } from './notifications';

export interface CentrumState {
    error: RpcError | undefined;
    abort: AbortController | undefined;
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
                this.units.unshift(unit);
            } else {
                this.units[idx].job = unit.job;
                this.units[idx].createdAt = unit.createdAt;
                this.units[idx].updatedAt = unit.updatedAt;
                this.units[idx].name = unit.name;
                this.units[idx].initials = unit.initials;
                this.units[idx].color = unit.color;
                this.units[idx].description = unit.description;
                this.units[idx].status = unit.status;
                this.units[idx].users = unit.users;
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
                this.dispatches?.unshift(dispatch);
            } else {
                this.dispatches[idx].createdAt = dispatch.createdAt;
                this.dispatches[idx].updatedAt = dispatch.updatedAt;
                this.dispatches[idx].job = dispatch.job;
                this.dispatches[idx].status = dispatch.status;
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
                this.ownDispatches?.unshift(dispatch);
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
            if (
                dispatch.status?.status === DISPATCH_STATUS.UNIT_UNASSIGNED ||
                dispatch.status?.status === DISPATCH_STATUS.UNASSIGNED
            ) {
                // Handle unassigment of dispatches
                this.removePendingDispatch(dispatch.id);
                this.removeOwnDispatch(dispatch.id);
            } else {
                const assignment = dispatch.units.find((u) => u.unitId === this.ownUnit?.id);
                if (assignment === undefined) return;
                // When dispatch has expiration, it is a "pending" dispatch
                if (assignment?.expiresAt) {
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
                this.pendingDispatches?.unshift(dispatch);
            }
        },

        removePendingDispatch(id: bigint): void {
            const tDIdx = this.pendingDispatches.findIndex((d) => d.id === id);
            if (tDIdx > -1) {
                this.pendingDispatches.splice(tDIdx, 1);
            }
        },

        async startStream(isController?: boolean): Promise<void> {
            if (this.abort !== undefined) return;

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
                        this.units = resp.change.latestState.units;
                        resp.change.latestState.dispatches.forEach((d) => this.addOrUpdateDispatch(d));
                    } else if (resp.change.oneofKind === 'settings') {
                        this.settings = resp.change.settings;
                    } else if (resp.change.oneofKind === 'disponents') {
                        this.disponents = resp.change.disponents.disponents;
                        // If user is not part of disponents list anymore
                        const idx = this.disponents.findIndex((d) => d.userId === authStore.activeChar?.userId);
                        if (idx === -1) {
                            this.stopStream();
                            setTimeout(() => {
                                this.startStream();
                            }, 250);
                        }
                    } else if (resp.change.oneofKind === 'unitAssigned') {
                        // Ignore, doesn't matter for controllers
                        if (!isController) {
                            if (this.ownUnit !== undefined && resp.change.unitAssigned.id !== this.ownUnit?.id) {
                                console.warn('Received unit user assigned event for other unit'), resp.change.unitAssigned;
                                continue;
                            }

                            const idx = resp.change.unitAssigned.users.findIndex(
                                (u) => u.userId === authStore.activeChar?.userId,
                            );
                            if (idx === -1) {
                                // User has been removed from the unit
                                this.ownUnit = undefined;

                                notifications.dispatchNotification({
                                    title: { key: 'notifications.centrum.unitAssigned.removed.title', parameters: [] },
                                    content: { key: 'notifications.centrum.unitAssigned.removed.content', parameters: [] },
                                    type: 'success',
                                });
                            } else {
                                // User has been added to unit
                                this.ownUnit = resp.change.unitAssigned;

                                notifications.dispatchNotification({
                                    title: { key: 'notifications.centrum.unitAssigned.joined.title', parameters: [] },
                                    content: { key: 'notifications.centrum.unitAssigned.joined.content', parameters: [] },
                                    type: 'success',
                                });
                            }
                        }
                    } else if (resp.change.oneofKind === 'unitDeleted') {
                        this.removeUnit(resp.change.unitDeleted);
                    } else if (resp.change.oneofKind === 'unitUpdated') {
                        this.addOrUpdateUnit(resp.change.unitUpdated);
                    } else if (resp.change.oneofKind === 'unitStatus') {
                        const id = resp.change.unitStatus.id;
                        let idx = this.units.findIndex((d) => d.id === id);
                        if (idx === -1) {
                            this.units?.unshift(resp.change.unitStatus);
                        } else {
                            this.units[idx] = resp.change.unitStatus;
                        }

                        if (resp.change.unitStatus.status) {
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

                        if (resp.change.dispatchStatus.status) {
                            this.feed.unshift(resp.change.dispatchStatus.status);
                        }

                        if (resp.change.dispatchStatus.status?.status === DISPATCH_STATUS.ARCHIVED) {
                            // If dispatch has been archived, remove from the main list
                            this.removeDispatch(id);
                        } else {
                            this.addOrUpdateDispatch(resp.change.dispatchStatus);
                        }
                    } else {
                        console.warn('Centrum: Unknown change received - Kind: ', resp.change.oneofKind, resp.change);
                    }

                    if (resp.restart !== undefined && resp.restart) {
                        this.stopStream();
                        setTimeout(() => {
                            this.startStream();
                        }, 250);
                    }
                }
            } catch (e) {
                this.error = e as RpcError;

                this.stopStream();
                notifications.dispatchNotification({
                    content: { key: this.error.message, parameters: [] },
                    title: { key: '', parameters: [] },
                });
            }

            console.debug('Centrum: Data Stream Ended');
        },
        async stopStream(): Promise<void> {
            console.debug('Centrum: Stopping Data Stream');
            if (this.abort) this.abort.abort();
            this.abort = undefined;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCentrumStore as unknown as StoreDefinition, import.meta.hot));
}
