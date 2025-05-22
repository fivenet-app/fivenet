<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import DispatchStatusUpdateModal from '~/components/centrum/dispatches/DispatchStatusUpdateModal.vue';
import DisponentsInfo from '~/components/centrum/disponents/DisponentsInfo.vue';
import {
    dispatchStatusToBGColor,
    dispatchStatuses,
    isStatusDispatchCompleted,
    unitStatusToBGColor,
    unitStatuses,
} from '~/components/centrum/helpers';
import DispatchesLayer from '~/components/centrum/livemap/DispatchesLayer.vue';
import JoinUnitSlideover from '~/components/centrum/livemap/JoinUnitSlideover.vue';
import OwnDispatchEntry from '~/components/centrum/livemap/OwnDispatchEntry.vue';
import TakeDispatchSlideover from '~/components/centrum/livemap/TakeDispatchSlideover.vue';
import DispatchStatusBreakdown from '~/components/centrum/partials/DispatchStatusBreakdown.vue';
import UnitDetailsSlideover from '~/components/centrum/units/UnitDetailsSlideover.vue';
import UnitStatusUpdateModal from '~/components/centrum/units/UnitStatusUpdateModal.vue';
import LivemapBase from '~/components/livemap/LivemapBase.vue';
import { setWaypointPLZ } from '~/composables/nui';
import { useCentrumStore } from '~/stores/centrum';
import { useLivemapStore } from '~/stores/livemap';
import { useNotificatorStore } from '~/stores/notificator';
import { useSettingsStore } from '~/stores/settings';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';
import { StatusUnit } from '~~/gen/ts/resources/centrum/units';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { $grpc } = useNuxtApp();

const modal = useModal();

const slideover = useSlideover();

const { can, jobProps } = useAuth();

const centrumStore = useCentrumStore();
const { startStream, stopStream } = centrumStore;
const { getCurrentMode, getOwnUnit, dispatches, getSortedOwnDispatches, pendingDispatches, timeCorrection, settings } =
    storeToRefs(centrumStore);

const livemapStore = useLivemapStore();
const { userOnDuty } = storeToRefs(livemapStore);

const notifications = useNotificatorStore();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const logger = useLogger('⛑️ Centrum');

const canStream = can('CentrumService.Stream');

const selectedDispatch = ref<number | undefined>();

async function updateDispatchStatus(dispatchId: number, status: StatusDispatch): Promise<void> {
    try {
        const call = $grpc.centrum.centrum.updateDispatchStatus({ dispatchId, status });
        await call;

        notifications.add({
            title: { key: 'notifications.centrum.sidebar.dispatch_status_updated.title', parameters: {} },
            description: { key: 'notifications.centrum.sidebar.dispatch_status_updated.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function updateDspStatus(dispatchId?: number, status?: StatusDispatch): Promise<void> {
    if (!dispatchId) {
        notifications.add({
            title: { key: 'notifications.centrum.sidebar.no_dispatch_selected.title', parameters: {} },
            description: { key: 'notifications.centrum.sidebar.no_dispatch_selected.content', parameters: {} },
            type: NotificationType.ERROR,
        });
        return;
    }

    if (status === undefined) {
        modal.open(DispatchStatusUpdateModal, {
            dispatchId: dispatchId,
            status: status,
        });
        return;
    }

    await updateDispatchStatus(dispatchId, status);
}

async function updateUnitStatus(id: number, status: StatusUnit): Promise<void> {
    try {
        const call = $grpc.centrum.centrum.updateUnitStatus({
            unitId: id,
            status,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.centrum.sidebar.unit_status_updated.title', parameters: {} },
            description: { key: 'notifications.centrum.sidebar.unit_status_updated.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function updateUtStatus(id: number, status?: StatusUnit): Promise<void> {
    if (status === undefined) {
        if (!getOwnUnit.value) {
            return;
        }

        modal.open(UnitStatusUpdateModal, {
            unit: getOwnUnit.value,
        });
        return;
    }

    await updateUnitStatus(id, status);
}

const open = ref(false);

async function toggleSidebarBasedOnUnit(): Promise<void> {
    if (getOwnUnit.value !== undefined) {
        // User has joined an unit
        open.value = true;
        slideover.close();

        if (
            jobProps.value !== undefined &&
            jobProps.value?.radioFrequency !== undefined &&
            jobProps.value.radioFrequency.length > 0
        ) {
            setRadioFrequency(jobProps.value.radioFrequency);
        }
    } else {
        // User not in an unit anymore
        open.value = false;
        slideover.close();
    }
}

const requireUnitInterval = computed(() => settings.value?.timings?.requireUnitReminderSeconds ?? 900 * 1000);
const { pause, resume } = useIntervalFn(
    () => sendRequireUnitNotification(),
    () => requireUnitInterval.value * 1000,
    {
        immediate: false,
    },
);

function toggleRequireUnitNotification(): void {
    if (canStream.value && settings.value?.enabled) {
        if (settings.value?.timings?.requireUnit === true && getOwnUnit.value === undefined) {
            resume();
        } else {
            pause();
        }
    }
}

// Show unit sidebar when ownUnit is set/updated, otherwise it will be hidden (automagically)
watch(getOwnUnit, () => {
    toggleSidebarBasedOnUnit();
    toggleRequireUnitNotification();
});

watch(open, async () => {
    if (open.value === true && getOwnUnit.value === undefined) {
        slideover.open(JoinUnitSlideover, {});
    }
});

const canSubmitUnitStatus = ref(true);
const onSubmitUnitStatusThrottle = useThrottleFn(async (unitId: number, status?: StatusUnit) => {
    canSubmitUnitStatus.value = false;
    await updateUtStatus(unitId, status).finally(() => useTimeoutFn(() => (canSubmitUnitStatus.value = true), 300));
}, 1000);

const canSubmitDispatchStatus = ref(true);
const onSubmitDispatchStatusThrottle = useThrottleFn(async (dispatchId?: number, status?: StatusDispatch) => {
    canSubmitDispatchStatus.value = false;
    await updateDspStatus(dispatchId, status).finally(() => useTimeoutFn(() => (canSubmitDispatchStatus.value = true), 300));
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
            const ownedDsp = getSortedOwnDispatches.value[index];
            if (!ownedDsp || ownedDsp === selectedDispatch.value) {
                continue;
            }

            const dispatch = dispatches.value.get(ownedDsp);
            if (isStatusDispatchCompleted(dispatch?.status?.status ?? StatusDispatch.UNSPECIFIED)) {
                continue;
            }

            selectedDispatch.value = ownedDsp;
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
                logger.debug('Centrum: Sidebar - Set Dispatch waypoint, id:', dispatch.id);
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

watch(settings, () => {
    if (!settings.value?.enabled) {
        return;
    }

    useIntervalFn(() => checkup(), 1 * 60 * 1000);
    toggleSidebarBasedOnUnit();
    toggleRequireUnitNotification();
});

onBeforeMount(async () => {
    if (!canStream.value) {
        return;
    }

    useTimeoutFn(async () => {
        try {
            startStream();
        } catch (e) {
            logger.error('exception during centrum stream', e);
        }
    }, 400);
});

onBeforeRouteLeave(async () => {
    await stopStream();
});

const attentionSound = useSounds('/sounds/centrum/attention.mp3', { playbackRate: 1.85 });

const unitCheckupStatusAge = 12.5 * 60 * 1000;
const unitCheckupStatusReping = 15 * 60 * 1000;

const debouncedPlay = useDebounceFn(async () => {
    attentionSound.play();
}, 950);

const attentionDebouncedPlay = useDebounceFn(async () => debouncedPlay(), 950);

const lastCheckupNotification = ref<Date | undefined>();

async function checkup(): Promise<void> {
    logger.debug('Centrum: Sidebar - Running checkup');
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

    notifications.add({
        title: { key: 'notifications.centrum.unitUpdated.checkup.title', parameters: {} },
        description: { key: 'notifications.centrum.unitUpdated.checkup.content', parameters: {} },
        type: NotificationType.INFO,
        timeout: 15000,
        callback: () => attentionDebouncedPlay(),
    });

    lastCheckupNotification.value = now;
}

function sendRequireUnitNotification(): void {
    if (!userOnDuty.value) {
        return;
    }

    useNotificatorStore().add({
        title: { key: 'notifications.centrum.unitUpdated.require_unit.title', parameters: {} },
        description: { key: 'notifications.centrum.unitUpdated.require_unit.content', parameters: {} },
        type: NotificationType.WARNING,
        timeout: 12500,
    });

    attentionSound.play();
}

function openTakeDispatches(): void {
    slideover.open(TakeDispatchSlideover, {});
}

defineShortcuts({
    'm-d': () => getOwnUnit.value && openTakeDispatches(),
    'm-h': () => getOwnUnit.value?.homePostal && setWaypointPLZ(getOwnUnit.value.homePostal),
    'c-u': () => getOwnUnit.value && onSubmitUnitStatusThrottle(getOwnUnit.value.id),
    'c-d': () => getOwnUnit.value && onSubmitDispatchStatusThrottle(),
});
</script>

<template>
    <UDashboardPanel grow>
        <UDashboardNavbar :title="$t('common.livemap')">
            <template #right>
                <DisponentsInfo v-if="canStream && settings?.enabled" :hide-join="true" />
            </template>
        </UDashboardNavbar>

        <UMain>
            <div class="relative z-0 size-full">
                <LivemapBase>
                    <template v-if="canStream" #default>
                        <DispatchesLayer
                            :show-all-dispatches="livemap.showAllDispatches || getCurrentMode === CentrumMode.SIMPLIFIED"
                        />

                        <LControl position="bottomright">
                            <UButton
                                v-if="settings?.enabled"
                                class="inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-clip-padding text-black hover:bg-[#f4f4f4]"
                                size="2xs"
                                :icon="open ? 'i-mdi-chevron-double-right' : 'i-mdi-chevron-double-left'"
                                @click="open = !open"
                            >
                                <span v-if="!open" class="inline-flex items-center justify-center">
                                    {{ $t('common.unit', 2) }}
                                </span>
                            </UButton>
                        </LControl>
                    </template>

                    <template v-if="canStream && settings?.enabled" #afterMap>
                        <div>
                            <Transition
                                enter-active-class="transform transition ease-in-out duration-100 sm:duration-200"
                                enter-from-class="translate-x-full"
                                enter-to-class="translate-x-0"
                                leave-active-class="transform transition ease-in-out duration-100 sm:duration-200"
                                leave-from-class="translate-x-0"
                                leave-to-class="translate-x-full"
                            >
                                <div
                                    v-if="open"
                                    class="bg-background flex h-full grow gap-y-5 overflow-y-auto overflow-x-hidden py-1"
                                    :class="open || getOwnUnit !== undefined ? 'px-2' : ''"
                                >
                                    <nav class="flex min-w-48 max-w-48 flex-1 flex-col md:min-w-64 md:max-w-64">
                                        <ul class="flex flex-1 flex-col gap-y-2 divide-y divide-base-400" role="list">
                                            <li>
                                                <ul class="-mx-1 space-y-0.5" role="list">
                                                    <li>
                                                        <UButton
                                                            v-if="getOwnUnit !== undefined"
                                                            class="flex flex-col"
                                                            :class="ownUnitStatus"
                                                            icon="i-mdi-information-outline"
                                                            block
                                                            @click="
                                                                slideover.open(UnitDetailsSlideover, {
                                                                    unit: getOwnUnit,
                                                                })
                                                            "
                                                        >
                                                            <span class="truncate">
                                                                <span class="font-semibold">{{ getOwnUnit.initials }}:</span>
                                                                {{ getOwnUnit.name }}</span
                                                            >
                                                            <span class="truncate text-xs">
                                                                <span class="font-semibold">{{ $t('common.status') }}:</span>
                                                                {{
                                                                    $t(
                                                                        `enums.centrum.StatusUnit.${
                                                                            StatusUnit[getOwnUnit.status?.status ?? 0]
                                                                        }`,
                                                                    )
                                                                }}
                                                            </span>
                                                        </UButton>

                                                        <UButtonGroup class="w-full" orientation="vertical">
                                                            <UButton
                                                                variant="soft"
                                                                color="primary"
                                                                size="xs"
                                                                block
                                                                :icon="
                                                                    getOwnUnit === undefined
                                                                        ? 'i-mdi-information-outline'
                                                                        : undefined
                                                                "
                                                                @click="slideover.open(JoinUnitSlideover, {})"
                                                            >
                                                                <template v-if="getOwnUnit === undefined">
                                                                    <span class="truncate">{{ $t('common.no_own_unit') }}</span>
                                                                </template>
                                                                <template v-else>
                                                                    <span class="truncate">{{ $t('common.leave_unit') }}</span>
                                                                </template>
                                                            </UButton>

                                                            <UButton
                                                                v-if="getOwnUnit === undefined"
                                                                variant="solid"
                                                                color="green"
                                                                size="xs"
                                                                block
                                                                icon="i-mdi-account-plus"
                                                                @click="slideover.open(JoinUnitSlideover, {})"
                                                            >
                                                                <span class="truncate">{{ $t('common.join_unit') }}</span>
                                                            </UButton>
                                                        </UButtonGroup>
                                                    </li>
                                                </ul>
                                            </li>

                                            <template v-if="getOwnUnit !== undefined">
                                                <li>
                                                    <ul class="-mx-1 space-y-0.5" role="list">
                                                        <li class="inline-flex items-center text-xs font-semibold leading-6">
                                                            {{ $t('common.units') }}
                                                            <UIcon
                                                                v-if="!canSubmitUnitStatus"
                                                                class="ml-1 size-4 animate-spin"
                                                                name="i-mdi-loading"
                                                            />
                                                        </li>

                                                        <li>
                                                            <div class="grid grid-cols-2 gap-0.5">
                                                                <UButton
                                                                    v-for="item in unitStatuses"
                                                                    :key="item.name"
                                                                    :class="[item.status && unitStatusToBGColor(item.status)]"
                                                                    :ui="{
                                                                        gap: { xs: 'gap-x-0.5' },
                                                                        padding: { xs: 'px-1.5 py-1.5' },
                                                                    }"
                                                                    size="xs"
                                                                    :disabled="!canSubmitUnitStatus"
                                                                    :icon="item.icon"
                                                                    truncate
                                                                    @click="
                                                                        onSubmitUnitStatusThrottle(getOwnUnit.id!, item.status)
                                                                    "
                                                                >
                                                                    <span class="line-clamp-2">
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
                                                                </UButton>

                                                                <UTooltip
                                                                    class="col-span-2"
                                                                    :text="$t('components.centrum.update_unit_status.title')"
                                                                    :shortcuts="['S', 'U']"
                                                                >
                                                                    <UButton
                                                                        variant="soft"
                                                                        color="primary"
                                                                        size="xs"
                                                                        block
                                                                        @click="onSubmitUnitStatusThrottle(getOwnUnit.id)"
                                                                    >
                                                                        {{ $t('components.centrum.update_unit_status.title') }}
                                                                    </UButton>
                                                                </UTooltip>
                                                            </div>
                                                        </li>
                                                    </ul>
                                                </li>

                                                <li>
                                                    <ul class="-mx-1 space-y-0.5" role="list">
                                                        <li class="inline-flex items-center text-xs font-semibold leading-6">
                                                            {{ $t('common.dispatch') }} {{ $t('common.status') }}
                                                            <UIcon
                                                                v-if="!canSubmitDispatchStatus"
                                                                class="ml-1 size-4 animate-spin"
                                                                name="i-mdi-loading"
                                                            />
                                                        </li>

                                                        <li>
                                                            <div class="grid grid-cols-2 gap-0.5">
                                                                <UButton
                                                                    v-for="item in dispatchStatuses.filter(
                                                                        (s) => s.status !== StatusDispatch.CANCELLED,
                                                                    )"
                                                                    :key="item.name"
                                                                    :class="[
                                                                        item.status && dispatchStatusToBGColor(item.status),
                                                                    ]"
                                                                    :ui="{
                                                                        gap: { xs: 'gap-x-0.5' },
                                                                        padding: { xs: 'px-1.5 py-1.5' },
                                                                    }"
                                                                    size="xs"
                                                                    :disabled="!canSubmitDispatchStatus"
                                                                    :icon="item.icon"
                                                                    @click="
                                                                        onSubmitDispatchStatusThrottle(
                                                                            selectedDispatch,
                                                                            item.status,
                                                                        )
                                                                    "
                                                                >
                                                                    <span class="mt-0.5 line-clamp-2">
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
                                                                </UButton>

                                                                <UTooltip
                                                                    class="col-span-2"
                                                                    :text="
                                                                        $t('components.centrum.update_dispatch_status.title')
                                                                    "
                                                                    :shortcuts="['S', 'D']"
                                                                >
                                                                    <UButton
                                                                        variant="soft"
                                                                        color="primary"
                                                                        size="xs"
                                                                        block
                                                                        @click="updateDspStatus(selectedDispatch)"
                                                                    >
                                                                        {{
                                                                            $t(
                                                                                'components.centrum.update_dispatch_status.title',
                                                                            )
                                                                        }}
                                                                    </UButton>
                                                                </UTooltip>
                                                            </div>
                                                        </li>
                                                    </ul>
                                                </li>

                                                <li>
                                                    <ul class="-mx-1 space-y-0.5" role="list">
                                                        <li class="inline-flex items-center text-xs font-semibold leading-6">
                                                            {{ $t('common.your_dispatches') }}
                                                        </li>

                                                        <li v-if="getSortedOwnDispatches.length === 0">
                                                            <UButton
                                                                variant="soft"
                                                                color="white"
                                                                icon="i-mdi-car-emergency"
                                                                block
                                                            >
                                                                {{ $t('common.no_assigned_dispatches') }}
                                                            </UButton>
                                                        </li>

                                                        <template v-else>
                                                            <template
                                                                v-for="id in getSortedOwnDispatches.slice().reverse()"
                                                                :key="id"
                                                            >
                                                                <OwnDispatchEntry
                                                                    v-if="dispatches.get(id) !== undefined"
                                                                    v-model:selected-dispatch="selectedDispatch"
                                                                    :dispatch="dispatches.get(id)!"
                                                                />
                                                            </template>
                                                        </template>
                                                    </ul>
                                                </li>

                                                <li>
                                                    <div class="mb-0.5 mt-1 flex w-full">
                                                        <DispatchStatusBreakdown size="xs" />
                                                    </div>
                                                </li>
                                            </template>
                                        </ul>
                                    </nav>
                                </div>
                            </Transition>

                            <!-- "Take Dispatches" Button -->
                            <span v-if="open && getOwnUnit !== undefined" class="fixed bottom-2 right-1/2 z-30 inline-flex">
                                <UChip
                                    :ui="{
                                        base: 'absolute rounded-full ring-0 ring-white dark:ring-gray-900 flex items-center justify-center text-white dark:text-gray-900 font-medium whitespace-nowrap animate-ping duration-750',
                                    }"
                                    position="top-left"
                                    size="xl"
                                    color="error"
                                    :show="pendingDispatches.length > 0"
                                >
                                    <UTooltip :text="$t('components.centrum.take_dispatch.title')" :shortcuts="['M', 'D']">
                                        <UButton
                                            class="flex size-12 items-center justify-center"
                                            :class="[getOwnUnit.homePostal !== undefined ? 'rounded-l-full' : 'rounded-full']"
                                            :color="pendingDispatches.length > 0 ? 'error' : 'primary'"
                                            size="xl"
                                            icon="i-mdi-car-emergency"
                                            @click="openTakeDispatches"
                                        />
                                    </UTooltip>
                                </UChip>

                                <UTooltip
                                    v-if="getOwnUnit.homePostal !== undefined"
                                    :text="`${$t('common.mark')}: ${$t('common.department_postal')}`"
                                    :shortcuts="['M', 'H']"
                                >
                                    <UButton
                                        class="flex size-12 items-center justify-center rounded-r-full"
                                        size="xl"
                                        icon="i-mdi-home-floor-b"
                                        @click="setWaypointPLZ(getOwnUnit.homePostal)"
                                    />
                                </UTooltip>
                            </span>
                        </div>
                    </template>
                </LivemapBase>
            </div>
        </UMain>
    </UDashboardPanel>
</template>
