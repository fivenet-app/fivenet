<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import {
    mdiCalendarCheck,
    mdiCalendarRemove,
    mdiCarBack,
    mdiCarEmergency,
    mdiCheckBold,
    mdiCoffee,
    mdiHelpCircle,
    mdiHoopHouse,
    mdiInformationOutline,
    mdiListStatus,
    mdiMarkerCheck,
} from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { Dispatch, DispatchStatus } from '~~/gen/ts/resources/dispatch/dispatches';
import { Settings } from '~~/gen/ts/resources/dispatch/settings';
import { Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { UserShort } from '~~/gen/ts/resources/users/users';

const { $grpc } = useNuxtApp();

const settings = ref<Settings>();
const unit = ref<Unit>();
const feed = ref<(DispatchStatus | UnitStatus)[]>([]);
const controllers = ref<UserShort[]>([]);
const dispatches = ref<Array<Dispatch>>([]);

type Action = { icon?: string; name: string; action?: Function };

const actionsDispatch: Action[] = [
    { icon: mdiCarBack, name: 'Dispatch: En Route' },
    { icon: mdiMarkerCheck, name: 'Dispatch: At Scene' },
    { icon: mdiHelpCircle, name: 'Dispatch: Need Assistance' },
    { icon: mdiListStatus, name: 'Dispatch: Update Status' },
    { icon: mdiCheckBold, name: 'Dispatch: Complete' },
];

const actionsUnit: Action[] = [
    { icon: mdiCarBack, name: 'Unit: Unavailable' },
    { icon: mdiCalendarCheck, name: 'Unit: Available' },
    { icon: mdiCoffee, name: 'Unit: On Break' },
    { icon: mdiCalendarRemove, name: 'Unit: Busy' },
    { icon: mdiListStatus, name: 'Unit: Update Status', action: () => {} },
];

const abort = ref<AbortController | undefined>();
const error = ref<string | null>(null);
async function startStream(): Promise<void> {
    if (abort.value !== undefined) return;

    console.debug('Centrum: Starting Data Stream');
    try {
        abort.value = new AbortController();

        const call = $grpc.getCentrumClient().stream(
            {},
            {
                abort: abort.value.signal,
            },
        );

        for await (let resp of call.responses) {
            error.value = null;

            if (resp === undefined || !resp.change) {
                continue;
            }
            if (!dispatches.value) {
                continue;
            }

            console.debug('Centrum: Received change - Kind:', resp.change.oneofKind, resp.change);

            if (resp.change.oneofKind === 'initial') {
                settings.value = resp.change.initial.settings;
                unit.value = resp.change.initial.unit;
            } else if (resp.change.oneofKind === 'dispatchUpdate') {
                const id = resp.change.dispatchUpdate.id;
                const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    dispatches.value?.unshift(resp.change.dispatchUpdate);
                } else {
                    dispatches.value![idx] = resp.change.dispatchUpdate;
                }
            } else if (resp.change.oneofKind === 'dispatchStatus') {
                feed.value.unshift(resp.change.dispatchStatus);
            } else if (resp.change.oneofKind === 'dispatchUnassigned') {
                const id = resp.change.dispatchUnassigned.id;
                const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    dispatches.value?.unshift(resp.change.dispatchUnassigned);
                } else {
                    dispatches.value![idx].units = resp.change.dispatchUnassigned.units;
                }
            } else if (resp.change.oneofKind === 'dispatchAssigned') {
                const id = resp.change.dispatchAssigned.id;
                const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    dispatches.value?.unshift(resp.change.dispatchAssigned);
                } else {
                    dispatches.value![idx] = resp.change.dispatchAssigned;
                }
            } else if (resp.change.oneofKind === 'unitUpdate') {
                const id = resp.change.unitUpdate.id;
                if (!unit.value) continue;

                if (unit.value.id === id) {
                    unit.value = resp.change.unitUpdate;
                }
            } else if (resp.change.oneofKind === 'unitStatus') {
                feed.value.unshift(resp.change.unitStatus);
            } else if (resp.change.oneofKind === 'unitAssigned') {
                // TODO show popup and notification
                if (resp.change.unitAssigned.id === 0n) {
                    // User has been removed from the unit
                } else {
                    // User has been added to unit
                }
            } else if (resp.change.oneofKind === 'unitDeleted') {
                if (!unit.value) continue;

                const id = resp.change.unitDeleted;
                if (unit.value.id === id) {
                    unit.value = undefined;
                }

                // TODO User not in a unit anymore
            } else if (resp.change.oneofKind === 'controllers') {
                controllers.value = resp.change.controllers.controllers;
                // If user is part of controllers list, we need to restart the stream
                if (!resp.change.controllers.active) {
                    stopStream();
                    setTimeout(() => {
                        startStream();
                    }, 250);
                }
            } else if (resp.change.oneofKind === 'settings') {
                settings.value = resp.change.settings;
            } else {
                console.log('Centrum: Unknown change received - Kind: ', resp.change.oneofKind, resp.change);
            }
        }
    } catch (e) {
        const err = e as RpcError;
        error.value = err.message;
        stopStream();
    }

    console.debug('Centrum: Data Stream Ended');
}

async function stopStream(): Promise<void> {
    console.debug('Centrum: Stopping Data Stream');
    abort.value?.abort();
    abort.value = undefined;
}

onMounted(() => {
    startStream();
});

onBeforeUnmount(() => {
    stopStream();
});
</script>

<template>
    <!-- Sidebar component, swap this element with another sidebar if you like -->
    <div class="h-full flex grow gap-y-5 overflow-y-auto bg-base-600 px-6 py-2">
        <nav class="flex flex-1 flex-col">
            <ul role="list" class="flex flex-1 flex-col gap-y-2 divide-y divide-base-400">
                <li>
                    <!-- <div class="text-xs font-semibold leading-6 text-base-200">Your Unit</div> -->
                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                        <li v-if="unit">
                            <button
                                type="button"
                                class="text-accent-100 bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-2"
                            >
                                <SvgIcon type="mdi" :path="mdiInformationOutline" class="h-6 w-6" aria-hidden="true" />
                                <span class="mt-2 truncate">{{ unit.initials }}: {{ unit.name }}</span>
                            </button>
                        </li>
                        <li v-else>
                            <button
                                type="button"
                                class="text-accent-100 bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-2"
                            >
                                <SvgIcon type="mdi" :path="mdiInformationOutline" class="h-6 w-6" aria-hidden="true" />
                                <span class="mt-2 truncate">You are not in any Unit.</span>
                            </button>
                        </li>
                    </ul>
                </li>
                <li>
                    <ul role="list" class="-mx-2 space-y-1">
                        <li v-for="item in actionsDispatch" :key="item.name">
                            <button
                                type="button"
                                class="text-accent-100 bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-2"
                            >
                                <SvgIcon
                                    type="mdi"
                                    :path="item.icon ?? mdiHoopHouse"
                                    class="text-base-200 group-hover:text-white h-6 w-6 shrink-0"
                                    aria-hidden="true"
                                />
                                <span class="mt-2">{{ item.name }}</span>
                            </button>
                        </li>
                    </ul>
                </li>
                <li>
                    <ul role="list" class="-mx-2 space-y-1">
                        <li v-for="item in actionsUnit" :key="item.name">
                            <button
                                type="button"
                                class="text-accent-100 bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-2"
                            >
                                <SvgIcon
                                    type="mdi"
                                    :path="item.icon ?? mdiHoopHouse"
                                    class="text-base-200 group-hover:text-white h-6 w-6 shrink-0"
                                    aria-hidden="true"
                                />
                                <span class="mt-2">{{ item.name }}</span>
                            </button>
                        </li>
                    </ul>
                </li>
                <li>
                    <!-- <div class="text-xs font-semibold leading-6 text-base-200">Your Dispatches</div> -->
                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                        <li v-if="!dispatches || dispatches.length === 0">
                            <button
                                type="button"
                                class="text-accent-100 bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-2"
                            >
                                <SvgIcon type="mdi" :path="mdiCarEmergency" class="h-6 w-6" aria-hidden="true" />
                                <span class="mt-2 truncate">No assigned Dispatches.</span>
                            </button>
                        </li>
                        <li v-else v-for="dispatch in dispatches" :key="dispatch.id.toString()">
                            <button
                                type="button"
                                class="text-accent-100 bg-error-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-2"
                            >
                                <SvgIcon type="mdi" :path="mdiCarEmergency" class="h-6 w-6" aria-hidden="true" />
                                <span class="mt-2 truncate">DSP-{{ dispatch.id }}</span>
                            </button>
                        </li>
                    </ul>
                </li>
            </ul>
        </nav>
    </div>
</template>
