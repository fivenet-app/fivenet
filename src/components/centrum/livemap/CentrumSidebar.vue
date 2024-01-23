<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { LControl } from '@vue-leaflet/vue-leaflet';
import { useDebounceFn, useIntervalFn, useThrottleFn, useTimeoutFn, watchDebounced } from '@vueuse/core';
import { useSound } from '@vueuse/sound';
import {
    CarEmergencyIcon,
    ChevronDownIcon,
    HomeFloorBIcon,
    HoopHouseIcon,
    InformationOutlineIcon,
    LoadingIcon,
    MonitorIcon,
    RobotIcon,
    ToggleSwitchIcon,
    ToggleSwitchOffIcon,
} from 'mdi-vue3';
import DispatchStatusUpdateModal from '~/components/centrum/dispatches/DispatchStatusUpdateModal.vue';
import DisponentsModal from '~/components/centrum/disponents/DisponentsModal.vue';
import {
    dispatchStatusToBGColor,
    dispatchStatuses,
    unitStatusToBGColor,
    unitStatuses,
    isStatusDispatchCompleted,
} from '~/components/centrum/helpers';
import UnitDetails from '~/components/centrum/units/UnitDetails.vue';
import UnitStatusUpdateModal from '~/components/centrum/units/UnitStatusUpdateModal.vue';
import { useCentrumStore } from '~/store/centrum';
import { useNotificatorStore } from '~/store/notificator';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';
import { StatusUnit } from '~~/gen/ts/resources/centrum/units';
import DispatchEntry from '~/components/centrum/livemap/DispatchEntry.vue';
import DispatchesLayer from '~/components/centrum/livemap/DispatchesLayer.vue';
import JoinUnitModal from '~/components/centrum/livemap/JoinUnitModal.vue';
import TakeDispatchModal from '~/components/centrum/livemap/TakeDispatchModal.vue';
import { useAuthStore } from '~/store/auth';
import LivemapBase from '~/components/livemap/LivemapBase.vue';
import { setWaypointPLZ } from '~/composables/nui';

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { jobProps } = storeToRefs(authStore);

const centrumStore = useCentrumStore();
const { getCurrentMode, getOwnUnit, dispatches, getSortedOwnDispatches, pendingDispatches, disponents, timeCorrection } =
    storeToRefs(centrumStore);
const { startStream, stopStream } = centrumStore;

const notifications = useNotificatorStore();

const canStream = can('CentrumService.Stream');

const joinUnitOpen = ref(false);

const selectedDispatch = ref<string | undefined>();
const openDispatchStatus = ref(false);
const openTakeDispatch = ref(false);

const openUnitDetails = ref(false);
const openUnitStatus = ref(false);

const openDisponents = ref(false);

async function updateDispatchStatus(dispatchId: string, status: StatusDispatch): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().updateDispatchStatus({ dispatchId, status });
        await call;

        notifications.dispatchNotification({
            title: { key: 'notifications.centrum.sidebar.dispatch_status_updated.title', parameters: {} },
            content: { key: 'notifications.centrum.sidebar.dispatch_status_updated.content', parameters: {} },
            type: 'success',
        });
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function updateDspStatus(dispatchId?: string, status?: StatusDispatch): Promise<void> {
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

    await updateDispatchStatus(dispatchId, status);
}

async function updateUnitStatus(id: string, status: StatusUnit): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().updateUnitStatus({
            unitId: id,
            status,
        });
        await call;

        notifications.dispatchNotification({
            title: { key: 'notifications.centrum.sidebar.unit_status_updated.title', parameters: {} },
            content: { key: 'notifications.centrum.sidebar.unit_status_updated.content', parameters: {} },
            type: 'success',
        });
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function updateUtStatus(id: string, status?: StatusUnit): Promise<void> {
    if (status === undefined) {
        openUnitStatus.value = true;
        return;
    }

    await updateUnitStatus(id, status);
}

const open = ref(false);

// Show unit sidebar when ownUnit is set/updated, otherwise it will be hidden (automagically)
watch(getOwnUnit, async () => {
    if (getOwnUnit.value !== undefined) {
        // User has joined an unit
        open.value = true;
        joinUnitOpen.value = false;

        if (
            jobProps.value !== undefined &&
            jobProps.value?.radioFrequency !== undefined &&
            jobProps.value.radioFrequency.length > 0
        ) {
            setRadioFrequency(parseInt(jobProps.value.radioFrequency));
        }
    } else {
        // User not in an unit anymore
        open.value = false;
        joinUnitOpen.value = false;
    }
});

watch(open, async () => {
    if (open.value === true && getOwnUnit.value === undefined) {
        joinUnitOpen.value = true;
    }
});

const canSubmitUnitStatus = ref(true);
const onSubmitUnitStatusThrottle = useThrottleFn(async (unitId: string, status?: StatusUnit) => {
    canSubmitUnitStatus.value = false;
    await updateUtStatus(unitId, status).finally(() => setTimeout(() => (canSubmitUnitStatus.value = true), 300));
}, 1000);

const canSubmitDispatchStatus = ref(true);
const onSubmitDispatchStatusThrottle = useThrottleFn(async (dispatchId?: string, status?: StatusDispatch) => {
    canSubmitDispatchStatus.value = false;
    await updateDspStatus(dispatchId, status).finally(() => setTimeout(() => (canSubmitDispatchStatus.value = true), 300));
}, 1000);

const ownUnitStatus = computed(() => unitStatusToBGColor(getOwnUnit.value?.status?.status));

function ensureOwnDispatchSelected(): void {
    if (getSortedOwnDispatches.value.length === 0) {
        selectedDispatch.value = undefined;
        return;
    }

    // If the selected dispatch is still our own dispatch, don't do anything
    if (
        selectedDispatch.value !== undefined &&
        getSortedOwnDispatches.value.find((dispatchId) => dispatchId === selectedDispatch.value) !== undefined
    ) {
        const dispatch = dispatches.value.get(selectedDispatch.value);
        if (!isStatusDispatchCompleted(dispatch?.status?.status ?? StatusDispatch.UNSPECIFIED)) {
            return;
        }
    }

    // otherwise select that current first one
    if (getSortedOwnDispatches.value.length > 1) {
        for (let index = 0; index < getSortedOwnDispatches.value.length; ++index) {
            const od = getSortedOwnDispatches.value[index];
            if (od === selectedDispatch.value) {
                continue;
            }

            const dispatch = dispatches.value.get(od);
            if (isStatusDispatchCompleted(dispatch?.status?.status ?? StatusDispatch.UNSPECIFIED)) {
                continue;
            }

            selectedDispatch.value = od;
            break;
        }
    } else {
        selectedDispatch.value = getSortedOwnDispatches.value[0];
    }
}

watchDebounced(
    selectedDispatch,
    () => {
        if (selectedDispatch.value !== undefined && getOwnUnit.value !== undefined) {
            const dispatch = dispatches.value.get(selectedDispatch.value);
            if (dispatch !== undefined) {
                setWaypoint(dispatch.x, dispatch.y);
                console.debug('Centrum: Sidebar - Set Dispatch waypoint, id:', dispatch.id);
            }
        }
    },
    {
        debounce: 75,
        maxWait: 400,
    },
);

watchDebounced(getSortedOwnDispatches.value, () => ensureOwnDispatchSelected(), {
    debounce: 75,
    maxWait: 200,
});

const { resume, pause } = useIntervalFn(() => checkup(), 1 * 60 * 1000);

const { start, stop } = useTimeoutFn(async () => startStream(), 650);

onBeforeMount(async () => {
    if (canStream) {
        start();
        resume();
    }
});

onBeforeUnmount(async () => {
    pause();
    stop();
    stopStream();
    centrumStore.$reset();
});

const unitCheckupStatusAge = 20 * 60 * 1000;
const unitCheckupStatusReping = 15 * 60 * 1000;

const attentionSound = useSound('/sounds/centrum/attention.mp3', {
    volume: 0.15,
    playbackRate: 1.25,
});

const debouncedPlay = useDebounceFn(() => attentionSound.play(), 950);

const lastCheckupNotification = ref<Date | undefined>();

async function checkup(): Promise<void> {
    console.debug('Centrum: Sidebar - Running checkup');
    const ownUnit = getOwnUnit.value;
    if (ownUnit === undefined || ownUnit.status === undefined) {
        return;
    }

    if (ownUnit.status.status === StatusUnit.AVAILABLE || ownUnit.status.status === StatusUnit.UNAVAILABLE) {
        return;
    }

    const now = new Date();
    // If unit status is younger than time X, ignore and continue
    if (now.getTime() - toDate(ownUnit.status.createdAt, timeCorrection.value).getTime() <= unitCheckupStatusAge) {
        return;
    }

    if (
        lastCheckupNotification.value !== undefined &&
        now.getTime() - lastCheckupNotification.value.getTime() <= unitCheckupStatusReping
    ) {
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

    lastCheckupNotification.value = now;
}
</script>

<template>
    <LivemapBase @goto="$emit('goto', $event)">
        <template v-if="canStream" #default>
            <DispatchesLayer :show-all-dispatches="getCurrentMode === CentrumMode.SIMPLIFIED" @goto="$emit('goto', $event)" />

            <LControl position="bottomright">
                <button
                    type="button"
                    class="inset-0 inline-flex items-center justify-center rounded-md border-2 border-black/20 bg-neutral bg-clip-padding text-black hover:bg-[#f4f4f4] focus:outline-none"
                    @click="open = !open"
                >
                    <ToggleSwitchIcon v-if="open" class="h-5 w-5" aria-hidden="true" />
                    <span v-else class="inline-flex items-center justify-center">
                        <ToggleSwitchOffIcon
                            class="h-5 w-5"
                            :class="getOwnUnit === undefined ? 'animate-pulse' : ''"
                            aria-hidden="true"
                        />
                        <span class="pr-0.5">
                            {{ $t('common.units') }}
                        </span>
                    </span>
                </button>
            </LControl>
        </template>

        <template v-if="canStream" #afterMap>
            <div class="lg:w-50 lg:inset-y-0 lg:flex lg:flex-col">
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
                    :dispatch-id="selectedDispatch"
                    @close="openDispatchStatus = false"
                />

                <div class="flex h-full grow gap-y-5 overflow-y-auto overflow-x-hidden bg-base-600 px-4 py-0.5">
                    <nav v-if="open" class="flex flex-1 flex-col">
                        <ul role="list" class="flex flex-1 flex-col gap-y-2 divide-y divide-base-400">
                            <li class="-mx-2 -mb-1">
                                <DisponentsModal :open="openDisponents" @close="openDisponents = false" />

                                <button
                                    type="button"
                                    class="group mt-0.5 flex w-full flex-row items-center justify-center rounded-md p-1 text-xs font-medium text-neutral hover:bg-primary-100/10 hover:text-neutral hover:transition-all"
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
                                        <MonitorIcon class="mr-1 h-5 w-5" aria-hidden="true" />
                                        <span class="truncate">
                                            {{ $t('common.disponent', disponents.length) }}
                                        </span>
                                    </template>
                                    <template v-else>
                                        <RobotIcon class="mr-1 h-5 w-5" aria-hidden="true" />
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
                                                type="button"
                                                class="group flex w-full flex-col items-center rounded-md p-1.5 text-xs font-medium text-neutral hover:bg-primary-100/10 hover:text-neutral hover:transition-all"
                                                :class="ownUnitStatus"
                                                @click="openUnitDetails = true"
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
                                                @goto="$emit('goto', $event)"
                                            />
                                        </template>
                                        <button
                                            type="button"
                                            class="group my-0.5 flex flex w-full w-full flex-col flex-col items-center items-center rounded-md bg-info-700 p-1.5 text-xs font-medium text-neutral hover:bg-primary-100/10 hover:text-neutral hover:transition-all"
                                            @click="joinUnitOpen = true"
                                        >
                                            <template v-if="getOwnUnit === undefined">
                                                <InformationOutlineIcon class="h-5 w-5" aria-hidden="true" />
                                                <span class="mt-1 truncate">{{ $t('common.no_own_unit') }}</span>
                                            </template>
                                            <template v-else>
                                                <span class="truncate">{{ $t('common.leave_unit') }}</span>
                                            </template>
                                        </button>

                                        <JoinUnitModal :open="joinUnitOpen" @close="joinUnitOpen = false" />
                                    </li>
                                </ul>
                            </li>
                            <template v-if="getOwnUnit !== undefined">
                                <li>
                                    <ul role="list" class="-mx-2 space-y-1">
                                        <li>
                                            <Disclosure as="div">
                                                <DisclosureButton
                                                    class="flex w-full items-start justify-between text-left text-neutral"
                                                >
                                                    <span class="leading-7 text-base-100">
                                                        <div
                                                            class="inline-flex items-center text-xs font-semibold leading-6 text-base-100"
                                                        >
                                                            {{ $t('common.unit') }}
                                                            <LoadingIcon
                                                                v-if="!canSubmitUnitStatus"
                                                                class="ml-1 h-4 w-4 animate-spin"
                                                            />
                                                        </div>
                                                    </span>
                                                    <span class="ml-6 flex h-7 items-center">
                                                        <ChevronDownIcon
                                                            :class="[open ? 'upsidedown' : '', 'h-5 w-5 transition-transform']"
                                                            aria-hidden="true"
                                                        />
                                                    </span>
                                                </DisclosureButton>
                                                <DisclosurePanel>
                                                    <div class="flex flex-row gap-2">
                                                        <div class="grid w-full grid-cols-2 gap-0.5">
                                                            <UnitStatusUpdateModal
                                                                :unit="getOwnUnit"
                                                                :open="openUnitStatus"
                                                                @close="openUnitStatus = false"
                                                            />

                                                            <button
                                                                v-for="item in unitStatuses"
                                                                :key="item.name"
                                                                type="button"
                                                                class="bg-primary group my-0.5 flex w-full flex-col items-center rounded-md p-1.5 text-xs font-medium text-neutral hover:bg-primary-100/10 hover:text-neutral hover:transition-all"
                                                                :disabled="!canSubmitUnitStatus"
                                                                :class="[
                                                                    !canSubmitUnitStatus ? 'disabled' : '',
                                                                    item.status ? unitStatusToBGColor(item.status) : item.class,
                                                                    item.class,
                                                                ]"
                                                                @click="onSubmitUnitStatusThrottle(getOwnUnit.id!, item.status)"
                                                            >
                                                                <component
                                                                    :is="item.icon ?? HoopHouseIcon"
                                                                    class="h-5 w-5 shrink-0 text-base-100 group-hover:text-neutral"
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
                                                            <button
                                                                type="button"
                                                                class="group col-span-2 my-0.5 flex w-full flex-col items-center rounded-md bg-base-800 p-1.5 text-xs font-medium text-neutral hover:bg-primary-100/10 hover:text-neutral hover:transition-all"
                                                                @click="updateUtStatus(getOwnUnit.id)"
                                                            >
                                                                {{ $t('components.centrum.update_unit_status.title') }}
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
                                        <div class="inline-flex items-center text-xs font-semibold leading-6 text-base-100">
                                            {{ $t('common.dispatch') }} {{ $t('common.status') }}
                                            <LoadingIcon v-if="!canSubmitDispatchStatus" class="ml-1 h-4 w-4 animate-spin" />
                                        </div>
                                        <li>
                                            <div class="grid grid-cols-2 gap-0.5">
                                                <button
                                                    v-for="item in dispatchStatuses.filter(
                                                        (s) => s.status !== StatusDispatch.CANCELLED,
                                                    )"
                                                    :key="item.name"
                                                    type="button"
                                                    class="bg-primary group my-0.5 flex w-full flex-col items-center rounded-md p-1.5 text-xs font-medium text-neutral hover:bg-primary-100/10 hover:text-neutral hover:transition-all"
                                                    :disabled="!canSubmitDispatchStatus"
                                                    :class="[
                                                        !canSubmitDispatchStatus ? 'disabled' : '',
                                                        item.status ? dispatchStatusToBGColor(item.status) : item.class,
                                                        item.class,
                                                    ]"
                                                    @click="onSubmitDispatchStatusThrottle(selectedDispatch, item.status)"
                                                >
                                                    <component
                                                        :is="item.icon ?? HoopHouseIcon"
                                                        class="h-5 w-5 shrink-0 text-base-100 group-hover:text-neutral"
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
                                                <button
                                                    type="button"
                                                    class="group col-span-2 my-0.5 flex w-full flex-col items-center rounded-md bg-base-800 p-1.5 text-xs font-medium text-neutral hover:bg-primary-100/10 hover:text-neutral hover:transition-all"
                                                    @click="updateDspStatus(selectedDispatch)"
                                                >
                                                    {{ $t('components.centrum.update_dispatch_status.title') }}
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
                                        <li v-if="getSortedOwnDispatches.length === 0">
                                            <button
                                                type="button"
                                                class="group my-0.5 flex w-full flex-col items-center rounded-md bg-primary-100/10 p-1.5 text-xs font-medium text-neutral hover:text-neutral hover:transition-all"
                                            >
                                                <CarEmergencyIcon class="h-5 w-5" aria-hidden="true" />
                                                <span class="mt-1 truncate">{{ $t('common.no_assigned_dispatches') }}</span>
                                            </button>
                                        </li>
                                        <template v-else>
                                            <DispatchEntry
                                                v-for="id in getSortedOwnDispatches.slice().reverse()"
                                                :key="id"
                                                v-model:selected-dispatch="selectedDispatch"
                                                :dispatch="dispatches.get(id)!"
                                                @goto="$emit('goto', $event)"
                                            />
                                        </template>
                                    </ul>
                                    <div
                                        class="mb-0.5 mt-1 divide-y border-t border-base-400 text-center text-xs leading-4 text-neutral"
                                    >
                                        {{ $t('components.centrum.livemap.total_dispatches') }}: {{ dispatches.size }}
                                    </div>
                                </li>
                            </template>
                        </ul>
                    </nav>
                </div>

                <!-- "Take Dispatches" Button -->
                <template v-if="getOwnUnit !== undefined">
                    <span class="fixed bottom-2 right-1/2 z-30 inline-flex">
                        <span>
                            <span v-if="pendingDispatches.length > 0" class="absolute left-0 top-0 -mr-1 -mt-1 flex h-3 w-3">
                                <span
                                    class="absolute inline-flex h-full w-full animate-ping rounded-full bg-error-400 opacity-75"
                                ></span>
                                <span class="relative inline-flex h-3 w-3 rounded-full bg-error-500"></span>
                            </span>
                            <button
                                type="button"
                                class="flex h-12 w-12 items-center justify-center bg-primary-500 text-neutral hover:bg-primary-400"
                                :class="getOwnUnit.homePostal !== undefined ? 'rounded-l-full' : 'rounded-full'"
                                @click="openTakeDispatch = true"
                            >
                                <CarEmergencyIcon class="h-auto w-10" />
                            </button>
                        </span>
                        <button
                            v-if="getOwnUnit.homePostal !== undefined"
                            type="button"
                            class="flex h-12 w-12 items-center justify-center rounded-r-full bg-primary-500 text-neutral hover:bg-primary-400"
                            @click="setWaypointPLZ(getOwnUnit.homePostal)"
                        >
                            <HomeFloorBIcon class="h-auto w-10" />
                        </button>
                    </span>
                </template>
            </div>
        </template>
    </LivemapBase>
</template>
