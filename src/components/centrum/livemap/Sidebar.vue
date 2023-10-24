<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { LControl } from '@vue-leaflet/vue-leaflet';
import { useDebounceFn, useIntervalFn, watchDebounced } from '@vueuse/core';
import { useSound } from '@vueuse/sound';
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
    MonitorIcon,
    RobotIcon,
    ToggleSwitchIcon,
    ToggleSwitchOffIcon,
} from 'mdi-vue3';
import { DefineComponent } from 'vue';
import { default as DispatchStatusUpdateModal } from '~/components/centrum/dispatches/StatusUpdateModal.vue';
import { default as DisponentsModal } from '~/components/centrum/disponents/Modal.vue';
import { dispatchStatusToBGColor, unitStatusToBGColor } from '~/components/centrum/helpers';
import { default as UnitDetails } from '~/components/centrum/units/Details.vue';
import { default as UnitStatusUpdateModal } from '~/components/centrum/units/StatusUpdateModal.vue';
import { useCentrumStore } from '~/store/centrum';
import { useNotificatorStore } from '~/store/notificator';
import { StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { CentrumMode } from '~~/gen/ts/resources/dispatch/settings';
import { StatusUnit } from '~~/gen/ts/resources/dispatch/units';
import DispatchEntry from './DispatchEntry.vue';
import DispatchesLayer from './DispatchesLayer.vue';
import JoinUnitModal from './JoinUnitModal.vue';
import TakeDispatchModal from './TakeDispatchModal.vue';

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const { $grpc } = useNuxtApp();

const centrumStore = useCentrumStore();
const { getCurrentMode, getOwnUnit, ownUnitId, dispatches, ownDispatches, pendingDispatches, disponents } =
    storeToRefs(centrumStore);
const { startStream, stopStream } = centrumStore;

const notifications = useNotificatorStore();

const actionsUnit: {
    icon: DefineComponent;
    name: string;
    action?: Function;
    class?: string;
    status?: StatusUnit;
}[] = [
    { icon: markRaw(CarBackIcon), name: 'Unavailable', status: StatusUnit.UNAVAILABLE },
    { icon: markRaw(CalendarCheckIcon), name: 'Available', status: StatusUnit.AVAILABLE },
    { icon: markRaw(CoffeeIcon), name: 'On Break', status: StatusUnit.ON_BREAK },
    { icon: markRaw(CalendarRemoveIcon), name: 'Busy', status: StatusUnit.BUSY },
    { icon: markRaw(ListStatusIcon), name: 'components.centrum.update_unit_status.title', class: 'bg-base-800' },
];

const actionsDispatch: {
    icon: DefineComponent;
    name: string;
    action?: Function;
    class?: string;
    status?: StatusDispatch;
}[] = [
    { icon: markRaw(CarBackIcon), name: 'En Route', status: StatusDispatch.EN_ROUTE },
    { icon: markRaw(MarkerCheckIcon), name: 'On Scene', status: StatusDispatch.ON_SCENE },
    { icon: markRaw(HelpCircleIcon), name: 'Need Assistance', status: StatusDispatch.NEED_ASSISTANCE },
    { icon: markRaw(CheckBoldIcon), name: 'Completed', status: StatusDispatch.COMPLETED },
    { icon: markRaw(ListStatusIcon), name: 'components.centrum.update_dispatch_status.title', class: 'bg-base-800' },
];

const canStream = can('CentrumService.Stream');

const joinUnitOpen = ref(false);

const selectedDispatch = ref<bigint | undefined>();
const openDispatchStatus = ref(false);
const openTakeDispatch = ref(false);

const openUnitDetails = ref(false);
const openUnitStatus = ref(false);

const openDisponents = ref(false);

async function updateDispatchStatus(dispatchId: bigint, status: StatusDispatch): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().updateDispatchStatus({
                dispatchId: dispatchId,
                status: status,
            });
            await call;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function updateDspStatus(dispatchId?: bigint, status?: StatusDispatch): Promise<void> {
    if (!dispatchId) {
        notifications.dispatchNotification({
            title: { key: 'notifications.centrum.sidebar.no_dispatch_selected.title', parameters: {} },
            content: { key: 'notifications.centrum.sidebar.no_dispatch_selected.content', parameters: {} },
            type: 'error',
        });
        return;
    }

    if (status === undefined) {
        openDispatchStatus.value = true;
        return;
    }

    updateDispatchStatus(dispatchId, status);
    notifications.dispatchNotification({
        title: { key: 'notifications.centrum.sidebar.dispatch_status_updated.title', parameters: {} },
        content: { key: 'notifications.centrum.sidebar.dispatch_status_updated.content', parameters: {} },
        type: 'success',
    });
}

async function updateUnitStatus(id: bigint, status: StatusUnit): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().updateUnitStatus({
                unitId: id,
                status: status,
            });
            await call;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function updateUtStatus(id: bigint, status?: StatusUnit): Promise<void> {
    if (status === undefined) {
        openUnitStatus.value = true;
        return;
    }

    updateUnitStatus(id, status);
    notifications.dispatchNotification({
        title: { key: 'notifications.centrum.sidebar.dispatch_status_updated.title', parameters: {} },
        content: { key: 'notifications.centrum.sidebar.dispatch_status_updated.content', parameters: {} },
        type: 'success',
    });
}

// Show unit sidebar when ownUnit is set/updated, otherwise it will be hidden (automagically)
watch(ownUnitId, async () => {
    if (ownUnitId.value !== undefined) {
        open.value = true;
    } else {
        open.value = false;
        joinUnitOpen.value = false;
    }
});

const ownUnitStatus = computed(() => unitStatusToBGColor(getOwnUnit.value?.status?.status));

async function ensureDispatchSelected(): Promise<void> {
    if (ownDispatches.value.length === 0) {
        selectedDispatch.value = undefined;
        return;
    }

    selectedDispatch.value = ownDispatches.value[ownDispatches.value.length - 1];
}

watchDebounced(
    selectedDispatch,
    () => {
        if (selectedDispatch.value !== undefined && isNUIAvailable()) {
            const dispatch = dispatches.value.get(selectedDispatch.value);
            if (dispatch !== undefined) {
                setWaypoint(dispatch.x, dispatch.y);
            }
        }
    },
    {
        debounce: 100,
        maxWait: 450,
    },
);

watchDebounced(ownDispatches.value, async () => ensureDispatchSelected(), {
    debounce: 150,
    maxWait: 600,
});

onBeforeMount(async () => {
    setTimeout(async () => canStream && startStream(), 250);
    ensureDispatchSelected();
});

onBeforeUnmount(async () => {
    stopStream();
    centrumStore.$reset();
});

const SEVENTEEN_MINUTES = 21 * 60 * 1000;

const attentionSound = useSound('/sounds/centrum/attention.mp3', {
    volume: 0.15,
    playbackRate: 1.25,
});

const debouncedPlay = useDebounceFn(() => attentionSound.play(), 950);

useIntervalFn(checkup, SEVENTEEN_MINUTES);

async function checkup(): Promise<void> {
    console.debug('Centrum: Sidebar - Running unit status checkup');
    const ownUnit = getOwnUnit.value;
    if (ownUnit === undefined || ownUnit.status === undefined) {
        return;
    }

    if (ownUnit.status.status === StatusUnit.AVAILABLE || ownUnit.status.status === StatusUnit.UNAVAILABLE) {
        return;
    }

    const notifications = useNotificatorStore();

    notifications.dispatchNotification({
        title: { key: 'notifications.centrum.unitUpdated.checkup.title', parameters: {} },
        content: { key: 'notifications.centrum.unitUpdated.checkup.content', parameters: {} },
        type: 'info',
        duration: 10000,
        callback: () => debouncedPlay(),
    });
}

const open = ref(false);
</script>

<template>
    <Livemap>
        <template v-slot:default v-if="canStream">
            <DispatchesLayer :show-all-dispatches="getCurrentMode === CentrumMode.SIMPLIFIED" @goto="$emit('goto', $event)" />

            <LControl position="bottomright">
                <button
                    type="button"
                    class="rounded-md bg-neutral text-black border-2 border-black/20 bg-clip-padding hover:bg-[#f4f4f4] focus:outline-none inset-0 inline-flex items-center justify-center"
                    @click="open = !open"
                >
                    <ToggleSwitchIcon v-if="open" class="h-6 w-6" aria-hidden="true" />
                    <span v-else class="inline-flex items-center justify-center">
                        <ToggleSwitchOffIcon
                            class="h-6 w-6"
                            :class="ownUnitId === undefined ? 'animate-pulse' : ''"
                            aria-hidden="true"
                        />
                        <span class="pr-0.5">
                            {{ $t('common.units') }}
                        </span>
                    </span>
                </button>
            </LControl>
        </template>
        <template v-slot:afterMap v-if="canStream">
            <div class="lg:inset-y-0 lg:flex lg:w-50 lg:flex-col">
                <!-- Dispatch -->
                <TakeDispatchModal
                    v-if="getOwnUnit !== undefined"
                    :open="openTakeDispatch"
                    @close="openTakeDispatch = false"
                    @goto="$emit('goto', $event)"
                />

                <DispatchStatusUpdateModal
                    v-if="selectedDispatch"
                    :open="openDispatchStatus"
                    @close="openDispatchStatus = false"
                    :dispatch-id="selectedDispatch"
                />

                <div class="h-full flex grow gap-y-5 overflow-y-auto overflow-x-hidden bg-base-600 px-4 py-0.5">
                    <nav v-if="open" class="flex flex-1 flex-col">
                        <ul role="list" class="flex flex-1 flex-col gap-y-2 divide-y divide-base-400">
                            <li class="-mx-2 -mb-1">
                                <DisponentsModal :open="openDisponents" @close="openDisponents = false" />

                                <button
                                    type="button"
                                    class="text-neutral hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-row items-center justify-center rounded-md p-1 text-xs mt-0.5"
                                    :class="
                                        getCurrentMode === CentrumMode.AUTO_ROUND_ROBIN
                                            ? 'bg-info-400/10 text-info-500 ring-info-400/20'
                                            : disponents.length === 0
                                            ? 'bg-warn-400/10 text-warn-500 ring-warn-400/20'
                                            : 'bg-success-500/10 text-success-400 ring-success-500/20'
                                    "
                                    @click="openDisponents = true"
                                >
                                    <template v-if="getCurrentMode !== CentrumMode.AUTO_ROUND_ROBIN">
                                        <MonitorIcon class="h-5 w-5 mr-1" aria-hidden="true" />
                                        <span class="truncate">
                                            {{ $t('common.disponent', disponents.length) }}
                                        </span>
                                    </template>
                                    <template v-else>
                                        <RobotIcon class="h-5 w-5 mr-1" aria-hidden="true" />
                                        <span class="truncate">
                                            {{ $t('enums.centrum.CentrumMode.AUTO_ROUND_ROBIN') }}
                                        </span>
                                    </template>
                                </button>
                            </li>
                            <li>
                                <ul role="list" class="-mx-2 mt-1 space-y-1">
                                    <li>
                                        <template v-if="getOwnUnit !== undefined">
                                            <button
                                                @click="openUnitDetails = true"
                                                type="button"
                                                class="text-neutral hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-1.5 text-xs my-0.5"
                                                :class="ownUnitStatus"
                                            >
                                                <InformationOutlineIcon class="h-5 w-5" aria-hidden="true" />
                                                <span class="mt-1 truncate">
                                                    <span class="font-semibold">{{ getOwnUnit.initials }}</span
                                                    >: {{ getOwnUnit.name }}</span
                                                >
                                                <span class="mt-1 truncate">
                                                    <span class="font-semibold">{{ $t('common.status') }}:</span>
                                                    {{
                                                        $t(
                                                            `enums.centrum.StatusUnit.${
                                                                StatusUnit[getOwnUnit.status?.status ?? 0]
                                                            }`,
                                                        )
                                                    }}
                                                </span>
                                            </button>

                                            <UnitDetails
                                                :unit="getOwnUnit"
                                                :open="openUnitDetails"
                                                @close="openUnitDetails = false"
                                            />
                                        </template>
                                        <button
                                            @click="joinUnitOpen = true"
                                            type="button"
                                            class="text-neutral bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-1.5 text-xs my-0.5"
                                        >
                                            <template v-if="getOwnUnit === undefined" class="flex w-full flex-col items-center">
                                                <InformationOutlineIcon class="h-5 w-5" aria-hidden="true" />
                                                <span class="mt-1 truncate">{{ $t('common.no_own_unit') }}</span>
                                            </template>
                                            <template v-else class="truncate">{{ $t('common.leave_unit') }}</template>
                                        </button>

                                        <JoinUnitModal :open="joinUnitOpen" @close="joinUnitOpen = false" />
                                    </li>
                                </ul>
                            </li>
                            <template v-if="getOwnUnit !== undefined">
                                <li>
                                    <ul role="list" class="-mx-2 space-y-1">
                                        <li>
                                            <Disclosure as="div" v-slot="{ open }">
                                                <DisclosureButton
                                                    class="flex w-full items-start justify-between text-left text-neutral"
                                                >
                                                    <span class="text-base-100 leading-7">
                                                        <div class="text-xs font-semibold leading-6 text-base-100">
                                                            {{ $t('common.unit') }}
                                                        </div>
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
                                                            <UnitStatusUpdateModal
                                                                :unit="getOwnUnit"
                                                                :open="openUnitStatus"
                                                                @close="openUnitStatus = false"
                                                            />

                                                            <button
                                                                v-for="(item, idx) in actionsUnit"
                                                                :key="item.name"
                                                                type="button"
                                                                class="text-neutral bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-1.5 text-xs my-0.5"
                                                                :class="[
                                                                    idx >= actionsUnit.length - 1 ? 'col-span-2' : '',
                                                                    item.status ? unitStatusToBGColor(item.status) : item.class,
                                                                    item.class,
                                                                ]"
                                                                @click="updateUtStatus(getOwnUnit.id, item.status)"
                                                            >
                                                                <component
                                                                    :is="item.icon ?? HoopHouseIcon"
                                                                    class="text-base-100 group-hover:text-neutral h-5 w-5 shrink-0"
                                                                    aria-hidden="true"
                                                                />
                                                                <span class="mt-1">
                                                                    {{
                                                                        item.status
                                                                            ? $t(
                                                                                  `enums.centrum.StatusUnit.${
                                                                                      StatusUnit[item.status ?? 0]
                                                                                  }`,
                                                                              )
                                                                            : $t(item.name)
                                                                    }}
                                                                </span>
                                                            </button>
                                                        </div>
                                                    </div>
                                                </DisclosurePanel>
                                            </Disclosure>
                                        </li>
                                    </ul>
                                </li>
                                <li>
                                    <ul role="list" class="-mx-2 space-y-1">
                                        <div class="text-xs font-semibold leading-6 text-base-100">
                                            {{ $t('common.dispatch') }} {{ $t('common.status') }}
                                        </div>
                                        <li>
                                            <div class="grid grid-cols-2 gap-0.5">
                                                <button
                                                    v-for="(item, idx) in actionsDispatch"
                                                    :key="item.name"
                                                    type="button"
                                                    class="text-neutral bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-1.5 text-xs my-0.5"
                                                    :class="[
                                                        idx >= actionsDispatch.length - 1 ? 'col-span-2' : '',
                                                        item.status ? dispatchStatusToBGColor(item.status) : item.class,
                                                        item.class,
                                                    ]"
                                                    @click="updateDspStatus(selectedDispatch, item.status)"
                                                >
                                                    <component
                                                        :is="item.icon ?? HoopHouseIcon"
                                                        class="text-base-100 group-hover:text-neutral h-5 w-5 shrink-0"
                                                        aria-hidden="true"
                                                    />
                                                    <span class="mt-1">
                                                        {{
                                                            item.status
                                                                ? $t(
                                                                      `enums.centrum.StatusDispatch.${
                                                                          StatusDispatch[item.status ?? 0]
                                                                      }`,
                                                                  )
                                                                : $t(item.name)
                                                        }}
                                                    </span>
                                                </button>
                                            </div>
                                        </li>
                                    </ul>
                                </li>
                                <li>
                                    <div class="text-xs font-semibold leading-6 text-base-100">
                                        {{ $t('common.your_dispatches') }}
                                    </div>
                                    <ul role="list" class="-mx-2 mt-1 space-y-1">
                                        <li v-if="ownDispatches.length === 0">
                                            <button
                                                type="button"
                                                class="text-neutral bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-1.5 text-xs my-0.5"
                                            >
                                                <CarEmergencyIcon class="h-5 w-5" aria-hidden="true" />
                                                <span class="mt-1 truncate">{{ $t('common.no_assigned_dispatches') }}</span>
                                            </button>
                                        </li>
                                        <template v-else>
                                            <DispatchEntry
                                                v-for="id in ownDispatches.slice().reverse()"
                                                :dispatch="dispatches.get(id)!"
                                                @goto="$emit('goto', $event)"
                                                v-model:selected-dispatch="selectedDispatch"
                                            />
                                        </template>
                                    </ul>
                                    <div
                                        class="mt-1 mb-0.5 leading-4 text-center text-xs text-neutral divide-y border-t border-base-400"
                                    >
                                        {{ $t('components.centrum.livemap.total_dispatches') }}: {{ dispatches.size }}
                                    </div>
                                </li>
                            </template>
                        </ul>
                    </nav>
                </div>

                <!-- "Take Dispatches" Button -->
                <span v-if="getOwnUnit !== undefined" class="fixed inline-flex z-30 bottom-2 right-1/2">
                    <span class="flex absolute h-3 w-3 top-0 right-0 -mt-1 -mr-1" v-if="pendingDispatches.length > 0">
                        <span
                            class="animate-ping absolute inline-flex h-full w-full rounded-full bg-error-400 opacity-75"
                        ></span>
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
            </div>
        </template>
    </Livemap>
</template>
