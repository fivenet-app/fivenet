<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import {
    CalendarCheckIcon,
    CalendarRemoveIcon,
    CarBackIcon,
    CarEmergencyIcon,
    CheckBoldIcon,
    ChevronDownIcon,
    CoffeeIcon,
    HelpCircleIcon,
    HoopHouseIcon,
    InformationOutlineIcon,
    ListStatusIcon,
    MarkerCheckIcon,
} from 'mdi-vue3';
import { DefineComponent } from 'vue';
import { default as DispatchDetails } from '~/components/centrum/dispatches/Details.vue';
import { default as UpdateDispatchStatus } from '~/components/centrum/dispatches/StatusUpdateModal.vue';
import { default as UnitDetails } from '~/components/centrum/units/Details.vue';
import { default as UpdateUnitStatus } from '~/components/centrum/units/StatusUpdateModal.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';
import { DISPATCH_STATUS, Dispatch, DispatchStatus } from '~~/gen/ts/resources/dispatch/dispatches';
import { Settings } from '~~/gen/ts/resources/dispatch/settings';
import { UNIT_STATUS, Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { dispatchStatusToBGColor } from './helpers';
import DispatchEntry from './sidebar/DispatchEntry.vue';
import JoinUnit from './sidebar/JoinUnitModal.vue';
import TakeDispatch from './sidebar/TakeDispatchModal.vue';

defineEmits<{
    (e: 'goto', loc: { x: number; y: number }): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const settings = ref<Settings>();
const disponents = ref<UserShort[]>([]);
const ownUnit = ref<Unit>();
const units = ref<Unit[]>([]);
const dispatches = ref<Dispatch[]>([]);
const feed = ref<(DispatchStatus | UnitStatus)[]>([]);
const takeDispatches = ref<Dispatch[]>([]);

const actionsDispatch: {
    icon: DefineComponent;
    name: string;
    action?: Function;
    class?: string;
    status?: DISPATCH_STATUS;
}[] = [
    { icon: markRaw(CarBackIcon), name: 'En Route', class: 'bg-info-600', status: DISPATCH_STATUS.EN_ROUTE },
    { icon: markRaw(MarkerCheckIcon), name: 'On Scene', class: 'bg-primary-600', status: DISPATCH_STATUS.ON_SCENE },
    { icon: markRaw(HelpCircleIcon), name: 'Need Assistance', class: 'bg-warn-600', status: DISPATCH_STATUS.NEED_ASSISTANCE },
    { icon: markRaw(CheckBoldIcon), name: 'Completed', class: 'bg-success-600', status: DISPATCH_STATUS.COMPLETED },
    { icon: markRaw(ListStatusIcon), name: 'Update Status', class: 'bg-base-800' },
];

const actionsUnit: {
    icon: DefineComponent;
    name: string;
    action?: Function;
    class?: string;
    status?: UNIT_STATUS;
}[] = [
    { icon: markRaw(CarBackIcon), name: 'Unavailable', class: 'bg-error-600', status: UNIT_STATUS.UNAVAILABLE },
    { icon: markRaw(CalendarCheckIcon), name: 'Available', class: 'bg-success-600', status: UNIT_STATUS.AVAILABLE },
    { icon: markRaw(CoffeeIcon), name: 'On Break', class: 'bg-warn-600', status: UNIT_STATUS.ON_BREAK },
    { icon: markRaw(CalendarRemoveIcon), name: 'Busy', class: 'bg-info-600', status: UNIT_STATUS.BUSY },
    { icon: markRaw(ListStatusIcon), name: 'Update Status', class: 'bg-base-800' },
];

function checkIfUnitAssignedToDispatch(dsp: Dispatch, unit?: Unit): boolean {
    if (unit === undefined) return false;

    return dsp.units.findIndex((d) => d.unitId === unit.id) > -1;
}

function addOrUpdateUnit(unit: Unit): void {
    const idx = units.value?.findIndex((d) => d.id === unit.id) ?? -1;
    if (idx === -1) {
        units.value?.unshift(unit);
    } else {
        units.value[idx].job = unit.job;
        units.value[idx].createdAt = unit.createdAt;
        units.value[idx].updatedAt = unit.updatedAt;
        units.value[idx].name = unit.name;
        units.value[idx].initials = unit.initials;
        units.value[idx].color = unit.color;
        units.value[idx].description = unit.description;
        units.value[idx].status = unit.status;
        units.value[idx].users = unit.users;
    }
}

function addOrUpdateDispatch(dispatch: Dispatch): void {
    const idx = dispatches.value?.findIndex((d) => d.id === dispatch.id) ?? -1;
    if (idx === -1) {
        dispatches.value?.unshift(dispatch);
    } else {
        dispatches.value[idx].createdAt = dispatch.createdAt;
        dispatches.value[idx].updatedAt = dispatch.updatedAt;
        dispatches.value[idx].job = dispatch.job;
        dispatches.value[idx].status = dispatch.status;
        dispatches.value[idx].message = dispatch.message;
        dispatches.value[idx].description = dispatch.description;
        dispatches.value[idx].attributes = dispatch.attributes;
        dispatches.value[idx].x = dispatch.x;
        dispatches.value[idx].y = dispatch.y;
        dispatches.value[idx].anon = dispatch.anon;
        dispatches.value[idx].userId = dispatch.userId;
        dispatches.value[idx].user = dispatch.user;
        if (dispatch.units.length == 0) {
            dispatches.value[idx].units.length = 0;
        } else {
            dispatches.value[idx].units = dispatch.units;
        }
    }
}

function addOrUpdateTakenDispatch(dispatch: Dispatch): void {
    const idx = takeDispatches.value?.findIndex((d) => d.id === dispatch.id) ?? -1;
    if (idx === -1) {
        takeDispatches.value?.unshift(dispatch);
    }
}

function removeUnit(unit: Unit): void {
    const idx = units.value?.findIndex((d) => d.id === unit.id) ?? -1;
    if (idx > -1) {
        units.value?.splice(idx, 1);
    }

    // User's unit has been deleted, reset it
    if (ownUnit.value !== undefined && ownUnit.value.id === unit.id) {
        ownUnit.value = undefined;
    }
}

function removeDispatchFromList(id: bigint): void {
    const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
    if (idx > -1) {
        dispatches.value?.splice(idx, 1);
    }

    removeTakenDispatch(id);
}

function removeTakenDispatch(id: bigint): void {
    const tDIdx = takeDispatches.value.findIndex((d) => d.id === id);
    if (tDIdx > -1) {
        takeDispatches.value.splice(tDIdx, 1);
    }
}

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

            if (resp.change.oneofKind === 'latestState') {
                settings.value = resp.change.latestState.settings;
                if (resp.change.latestState.unit !== undefined) {
                    ownUnit.value = resp.change.latestState.unit;
                } else {
                    ownUnit.value = undefined;
                }
                units.value = resp.change.latestState.units;
                dispatches.value = resp.change.latestState.dispatches;
            } else if (resp.change.oneofKind === 'settings') {
                settings.value = resp.change.settings;
            } else if (resp.change.oneofKind === 'disponents') {
                disponents.value = resp.change.disponents.disponents;
            } else if (resp.change.oneofKind === 'unitAssigned') {
                if (ownUnit.value !== undefined && resp.change.unitAssigned.id !== ownUnit.value?.id) {
                    console.warn('Received unit user assigned event for other unit'), resp.change.unitAssigned;
                    continue;
                }

                const idx = resp.change.unitAssigned.users.findIndex((u) => u.userId === activeChar.value?.userId);
                if (idx === -1) {
                    // User has been removed from the unit
                    ownUnit.value = undefined;

                    notifications.dispatchNotification({
                        title: { key: 'notifications.centrum.unitAssigned.removed.title', parameters: [] },
                        content: { key: 'notifications.centrum.unitAssigned.removed.content', parameters: [] },
                        type: 'success',
                    });
                } else {
                    // User has been added to unit
                    ownUnit.value = resp.change.unitAssigned;

                    notifications.dispatchNotification({
                        title: { key: 'notifications.centrum.unitAssigned.joined.title', parameters: [] },
                        content: { key: 'notifications.centrum.unitAssigned.joined.content', parameters: [] },
                        type: 'success',
                    });
                }
            } else if (resp.change.oneofKind === 'unitDeleted') {
                removeUnit(resp.change.unitDeleted);
            } else if (resp.change.oneofKind === 'unitUpdated') {
                addOrUpdateUnit(resp.change.unitUpdated);
            } else if (resp.change.oneofKind === 'unitStatus') {
                const id = resp.change.unitStatus.id;
                const unit = units.value.find((u) => u.id === id);
                if (unit) {
                    unit.status = resp.change.unitStatus.status;
                } else {
                    units.value.push(resp.change.unitStatus);
                }

                if (resp.change.unitStatus.status) {
                    feed.value.unshift(resp.change.unitStatus.status);
                }
            } else if (resp.change.oneofKind === 'dispatchCreated') {
                if (!checkIfUnitAssignedToDispatch(resp.change.dispatchCreated, ownUnit.value)) continue;

                const id = resp.change.dispatchCreated.id;
                const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    dispatches.value?.unshift(resp.change.dispatchCreated);
                } else {
                    dispatches.value[idx].units = resp.change.dispatchCreated.units;
                }
            } else if (resp.change.oneofKind === 'dispatchDeleted') {
                removeDispatchFromList(resp.change.dispatchDeleted.id);
            } else if (resp.change.oneofKind === 'dispatchUpdated') {
                if (!checkIfUnitAssignedToDispatch(resp.change.dispatchUpdated, ownUnit.value)) continue;

                addOrUpdateDispatch(resp.change.dispatchUpdated);
            } else if (resp.change.oneofKind === 'dispatchStatus') {
                const id = resp.change.dispatchStatus.id;
                let idx = dispatches.value.findIndex((d) => d.id === id);
                if (idx === -1) {
                    dispatches.value?.unshift(resp.change.dispatchStatus);
                } else {
                    dispatches.value[idx] = resp.change.dispatchStatus;
                }

                if (resp.change.dispatchStatus.status) {
                    feed.value.unshift(resp.change.dispatchStatus.status);
                }

                if (resp.change.dispatchStatus.status?.status === DISPATCH_STATUS.UNIT_ASSIGNED) {
                    if (ownUnit.value && ownUnit.value.id === resp.change.dispatchStatus.status.unitId) {
                        const assignment = resp.change.dispatchStatus.units.find((u) => u.unitId === ownUnit.value?.id);
                        // When dispatch has expiration, it needs to be "taken"
                        if (assignment?.expiresAt) {
                            addOrUpdateTakenDispatch(resp.change.dispatchStatus);
                        } else {
                            removeTakenDispatch(resp.change.dispatchStatus.id);
                        }
                    }
                } else if (
                    resp.change.dispatchStatus.status?.status === DISPATCH_STATUS.UNIT_UNASSIGNED ||
                    resp.change.dispatchStatus.status?.status === DISPATCH_STATUS.ARCHIVED
                ) {
                    removeDispatchFromList(id);
                }
            } else {
                console.warn('Centrum: Unknown change received - Kind: ', resp.change.oneofKind, resp.change);
            }

            if (resp.restart !== undefined && resp.restart) {
                stopStream();
                setTimeout(() => {
                    startStream();
                }, 250);
            }
        }
    } catch (e) {
        const err = e as RpcError;
        error.value = err.message;
        notifications.dispatchNotification({
            content: { key: err.message, parameters: [] },
            title: { key: '', parameters: [] },
        });
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
const selectUnitOpen = ref(false);

const unitStatusOpen = ref(false);
const unitStatusSelected = ref<UNIT_STATUS | undefined>();

// TODO add function to set the active dispatch (last clicked and manually selected)
const activeDispatch = ref<Dispatch | undefined>();
const dispatchStatusOpen = ref(false);
const dispatchStatusSelected = ref<DISPATCH_STATUS | undefined>();

const openTakeDispatch = ref(false);

const selectedDispatch = ref<Dispatch | undefined>();
const openDispatchDetails = ref(false);
</script>

<template>
    <template v-if="ownUnit">
        <UpdateUnitStatus :open="unitStatusOpen" @close="unitStatusOpen = false" :unit="ownUnit" :status="unitStatusSelected" />
        <UpdateDispatchStatus
            v-if="activeDispatch"
            :open="dispatchStatusOpen"
            @close="dispatchStatusOpen = false"
            :dispatch="activeDispatch"
            :status="dispatchStatusSelected"
        />
        <TakeDispatch
            :open="openTakeDispatch"
            @close="openTakeDispatch = false"
            :own-unit="ownUnit"
            :dispatches="takeDispatches"
            @goto="$emit('goto', $event)"
        />
    </template>

    <div class="h-full flex grow gap-y-5 overflow-y-auto bg-base-600 px-4 py-0.5">
        <nav class="flex flex-1 flex-col">
            <ul role="list" class="flex flex-1 flex-col gap-y-2 divide-y divide-base-400">
                <li>
                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                        <li>
                            <div v-if="ownUnit">
                                <button
                                    @click="unitOpen = true"
                                    type="button"
                                    class="text-white bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                >
                                    <InformationOutlineIcon class="h-5 w-5" aria-hidden="true" />
                                    <span class="mt-2 truncate">{{ ownUnit.initials }}: {{ ownUnit.name }}</span>
                                    <span class="mt-2 truncate">
                                        {{
                                            $t(
                                                `enums.centrum.UNIT_STATUS.${
                                                    UNIT_STATUS[ownUnit.status?.status ?? (0 as number)]
                                                }`,
                                            )
                                        }}
                                    </span>
                                </button>
                                <UnitDetails :unit="ownUnit" :ownUnit="ownUnit" :open="unitOpen" @close="unitOpen = false" />
                            </div>
                            <button
                                @click="selectUnitOpen = true"
                                type="button"
                                class="text-white bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                            >
                                <span v-if="!ownUnit" class="flex w-full flex-col items-center">
                                    <InformationOutlineIcon class="h-5 w-5" aria-hidden="true" />
                                    <span class="mt-2 truncate">{{ $t('common.no_own_unit') }}</span>
                                </span>
                                <span v-else class="truncate">{{ $t('common.leave_unit') }}</span>
                            </button>

                            <JoinUnit
                                :open="selectUnitOpen"
                                @close="selectUnitOpen = false"
                                @joined="ownUnit = $event"
                                :own-unit="ownUnit"
                                :units="units"
                            />
                        </li>
                    </ul>
                </li>
                <li v-if="ownUnit">
                    <ul role="list" class="-mx-2 space-y-1">
                        <li>
                            <Disclosure as="div" v-slot="{ open }">
                                <DisclosureButton class="flex w-full items-start justify-between text-left text-white">
                                    <span class="text-base-200 leading-7">
                                        <div class="text-xs font-semibold leading-6 text-base-200">{{ $t('common.unit') }}</div>
                                    </span>
                                    <span class="ml-6 flex h-7 items-center">
                                        <ChevronDownIcon
                                            :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                            aria-hidden="true"
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
                                                @click="
                                                    unitStatusSelected = item.status;
                                                    unitStatusOpen = true;
                                                "
                                            >
                                                <component
                                                    :is="item.icon ?? HoopHouseIcon"
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
                <li v-if="ownUnit">
                    <ul role="list" class="-mx-2 space-y-1">
                        <div class="text-xs font-semibold leading-6 text-base-200">
                            {{ $t('common.dispatch', 2) }}
                        </div>
                        <li>
                            <div class="grid grid-cols-2 gap-0.5">
                                <button
                                    v-for="(item, idx) in actionsDispatch"
                                    :key="item.name"
                                    type="button"
                                    class="text-white bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                    :class="[
                                        idx >= actionsDispatch.length - 1 ? 'col-span-2' : '',
                                        item.class,
                                        dispatchStatusToBGColor(item.status),
                                    ]"
                                    @click="
                                        dispatchStatusSelected = item.status;
                                        dispatchStatusOpen = true;
                                    "
                                >
                                    <component
                                        :is="item.icon ?? HoopHouseIcon"
                                        class="text-base-200 group-hover:text-white h-5 w-5 shrink-0"
                                        aria-hidden="true"
                                    />
                                    <span class="mt-1">{{ item.name }}</span>
                                </button>
                            </div>
                        </li>
                    </ul>
                </li>
                <li v-if="ownUnit">
                    <div class="text-xs font-semibold leading-6 text-base-200">{{ $t('common.your_dispatches') }}</div>
                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                        <li v-if="dispatches.length === 0">
                            <button
                                type="button"
                                class="text-white bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                            >
                                <CarEmergencyIcon class="h-5 w-5" aria-hidden="true" />
                                <span class="mt-2 truncate">{{ $t('common.no_assigned_dispatches') }}</span>
                            </button>
                        </li>
                        <DispatchEntry
                            v-else
                            v-for="dispatch in dispatches"
                            :dispatch="dispatch"
                            :units="units"
                            @goto="$emit('goto', $event)"
                            @details=""
                        />
                    </ul>
                </li>
            </ul>
        </nav>
    </div>

    <template v-if="selectedDispatch">
        <DispatchDetails
            @close="openDispatchDetails = false"
            :dispatch="selectedDispatch"
            :open="openDispatchDetails"
            @goto="$emit('goto', $event)"
            :units="units"
        />
    </template>

    <!-- "Take Dispatches" Button -->
    <span v-if="ownUnit" class="fixed inline-flex z-90 bottom-2 right-1/2">
        <span class="flex absolute h-3 w-3 top-0 right-0 -mt-1 -mr-1" v-if="takeDispatches.length > 0">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-error-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-3 w-3 bg-error-500"></span>
        </span>
        <button
            type="button"
            @click="openTakeDispatch = true"
            class="flex items-center justify-center w-12 h-12 rounded-full bg-primary-500 shadow-float text-neutral hover:bg-primary-400"
        >
            <CarEmergencyIcon class="w-10 h-auto" />
        </button>
    </span>
</template>
