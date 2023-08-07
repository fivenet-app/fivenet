<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import {
    mdiCalendarCheck,
    mdiCalendarRemove,
    mdiCarBack,
    mdiCarEmergency,
    mdiCheckBold,
    mdiChevronDown,
    mdiCoffee,
    mdiHelpCircle,
    mdiHoopHouse,
    mdiInformationOutline,
    mdiListStatus,
    mdiMarkerCheck,
} from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { default as UnitDetails } from '~/components/centrum/units/Details.vue';
import DispatchEntry from '~/components/livemap/centrum/DispatchEntry.vue';
import { DISPATCH_STATUS, Dispatch, DispatchStatus } from '~~/gen/ts/resources/dispatch/dispatches';
import { Settings } from '~~/gen/ts/resources/dispatch/settings';
import { UNIT_STATUS, Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { UserShort } from '~~/gen/ts/resources/users/users';

const { $grpc } = useNuxtApp();

const settings = ref<Settings>();
const unit = ref<Unit>();
const feed = ref<(DispatchStatus | UnitStatus)[]>([]);
const controllers = ref<UserShort[]>([]);
const dispatches = ref<Array<Dispatch>>([]);

type Action = { icon: string; name: string; action?: Function; class?: string; status?: DISPATCH_STATUS | UNIT_STATUS };

const actionsDispatch: Action[] = [
    { icon: mdiCarBack, name: 'En Route', class: 'bg-info-600', status: DISPATCH_STATUS.EN_ROUTE },
    { icon: mdiMarkerCheck, name: 'On Scene', class: 'bg-primary-600', status: DISPATCH_STATUS.ON_SCENE },
    { icon: mdiHelpCircle, name: 'Need Assistance', class: 'bg-warn-600', status: DISPATCH_STATUS.NEED_ASSISTANCE },
    { icon: mdiCheckBold, name: 'Completed', class: 'bg-success-600', status: DISPATCH_STATUS.COMPLETED },
    { icon: mdiListStatus, name: 'Update Status', class: 'bg-base-800' },
];

const actionsUnit: Action[] = [
    { icon: mdiCarBack, name: 'Unavailable', class: 'bg-error-600', status: UNIT_STATUS.UNAVAILABLE },
    { icon: mdiCalendarCheck, name: 'Available', class: 'bg-success-600', status: UNIT_STATUS.AVAILABLE },
    { icon: mdiCoffee, name: 'On Break', class: 'bg-warn-600', status: UNIT_STATUS.ON_BREAK },
    { icon: mdiCalendarRemove, name: 'Busy', class: 'bg-info-600', status: UNIT_STATUS.BUSY },
    { icon: mdiListStatus, name: 'Update Status', class: 'bg-base-800' },
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

            console.debug('Centrum: Received change - Kind:', resp.change.oneofKind, resp.change);

            if (!dispatches.value) {
                continue;
            }

            if (resp.change.oneofKind === 'initial') {
                settings.value = resp.change.initial.settings;
                unit.value = resp.change.initial.unit;
                dispatches.value = resp.change.initial.dispatches;
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

const unitOpen = ref(false);
</script>

<template>
    <!-- Sidebar component, swap this element with another sidebar if you like -->
    <div class="h-full flex grow gap-y-5 overflow-y-auto bg-base-600 px-4 py-0.5">
        <nav class="flex flex-1 flex-col">
            <ul role="list" class="flex flex-1 flex-col gap-y-2 divide-y divide-base-400">
                <li>
                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                        <li v-if="unit">
                            <button
                                @click="unitOpen = true"
                                type="button"
                                class="text-white bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                            >
                                <SvgIcon type="mdi" :path="mdiInformationOutline" class="h-5 w-5" aria-hidden="true" />
                                <span class="mt-2 truncate">{{ unit.initials }}: {{ unit.name }}</span>
                            </button>
                            <UnitDetails :unit="unit" :open="unitOpen" @close="unitOpen = false" />
                        </li>
                        <li v-else>
                            <button
                                type="button"
                                class="text-white bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                            >
                                <SvgIcon type="mdi" :path="mdiInformationOutline" class="h-5 w-5" aria-hidden="true" />
                                <span class="mt-2 truncate">Not in any Unit.</span>
                            </button>
                        </li>
                    </ul>
                </li>
                <li v-if="unit">
                    <ul role="list" class="-mx-2 space-y-1">
                        <li>
                            <Disclosure as="div" v-slot="{ open }">
                                <DisclosureButton class="flex w-full items-start justify-between text-left text-white">
                                    <span class="text-base-200 leading-7">
                                        <div class="text-xs font-semibold leading-6 text-base-200">Unit</div>
                                    </span>
                                    <span class="ml-6 flex h-7 items-center">
                                        <SvgIcon
                                            :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                            aria-hidden="true"
                                            type="mdi"
                                            :path="mdiChevronDown"
                                        />
                                    </span>
                                </DisclosureButton>
                                <DisclosurePanel>
                                    <div class="flex flex-row gap-2">
                                        <div class="w-full grid grid-cols-2 gap-0.5">
                                            <button
                                                v-for="(item, idx) in actionsUnit"
                                                :key="item.name"
                                                type="button"
                                                class="text-white bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                                :class="[idx >= actionsUnit.length - 1 ? 'col-span-2' : '', item.class]"
                                            >
                                                <SvgIcon
                                                    type="mdi"
                                                    :path="item.icon ?? mdiHoopHouse"
                                                    class="text-base-200 group-hover:text-white h-5 w-5 shrink-0"
                                                    aria-hidden="true"
                                                />
                                                <span class="mt-1">{{ item.name }}</span>
                                            </button>
                                        </div>
                                    </div>
                                </DisclosurePanel>
                            </Disclosure>
                        </li>
                    </ul>
                </li>
                <li v-if="unit">
                    <ul role="list" class="-mx-2 space-y-1">
                        <div class="text-xs font-semibold leading-6 text-base-200">Dispatches</div>
                        <li>
                            <div class="grid grid-cols-2 gap-0.5">
                                <button
                                    v-for="(item, idx) in actionsDispatch"
                                    :key="item.name"
                                    type="button"
                                    class="text-white bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                    :class="[idx >= actionsDispatch.length - 1 ? 'col-span-2' : '', item.class]"
                                >
                                    <SvgIcon
                                        type="mdi"
                                        :path="item.icon ?? mdiHoopHouse"
                                        class="text-base-200 group-hover:text-white h-5 w-5 shrink-0"
                                        aria-hidden="true"
                                    />
                                    <span class="mt-1">{{ item.name }}</span>
                                </button>
                            </div>
                        </li>
                    </ul>
                </li>
                <li v-if="unit">
                    <div class="text-xs font-semibold leading-6 text-base-200">Your Dispatches</div>
                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                        <li v-if="!dispatches || dispatches.length === 0">
                            <button
                                type="button"
                                class="text-white bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                            >
                                <SvgIcon type="mdi" :path="mdiCarEmergency" class="h-5 w-5" aria-hidden="true" />
                                <span class="mt-2 truncate">No assigned Dispatches.</span>
                            </button>
                        </li>
                        <DispatchEntry v-else v-for="dispatch in dispatches" :dispatch="dispatch" />
                    </ul>
                </li>
            </ul>
        </nav>
    </div>
</template>
