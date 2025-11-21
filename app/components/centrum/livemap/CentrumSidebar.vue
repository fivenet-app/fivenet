<script lang="ts" setup>
import DispatcherInfo from '~/components/centrum/dispatchers/DispatcherInfo.vue';
import DispatchStatusUpdateModal from '~/components/centrum/dispatches/DispatchStatusUpdateModal.vue';
import {
    dispatchStatusToBadgeColor,
    dispatchStatuses,
    isStatusDispatchCompleted,
    unitStatusToBadgeColor,
    unitStatuses,
} from '~/components/centrum/helpers';
import DispatchLayer from '~/components/centrum/livemap/DispatchLayer.vue';
import JoinUnitSlideover from '~/components/centrum/livemap/JoinUnitSlideover.vue';
import OwnDispatchEntry from '~/components/centrum/livemap/OwnDispatchEntry.vue';
import TakeDispatchSlideover from '~/components/centrum/livemap/TakeDispatchSlideover.vue';
import DispatchStatusBreakdown from '~/components/centrum/partials/DispatchStatusBreakdown.vue';
import UnitDetailsSlideover from '~/components/centrum/units/UnitDetailsSlideover.vue';
import UnitStatusUpdateModal from '~/components/centrum/units/UnitStatusUpdateModal.vue';
import FollowMarker from '~/components/livemap/controls/FollowMarker.vue';
import LivemapBase from '~/components/livemap/LivemapBase.vue';
import { setWaypointPLZ } from '~/composables/nui';
import { useCentrumStore } from '~/stores/centrum';
import { useLivemapStore } from '~/stores/livemap';
import { useSettingsStore } from '~/stores/settings';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings';
import { StatusUnit } from '~~/gen/ts/resources/centrum/units';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { can, jobProps } = useAuth();

const centrumStore = useCentrumStore();
const { startStream, stopStream } = centrumStore;
const { getCurrentMode, getOwnUnit, dispatches, getSortedOwnDispatches, pendingDispatches, timeCorrection, settings } =
    storeToRefs(centrumStore);

const livemapStore = useLivemapStore();
const { userOnDuty } = storeToRefs(livemapStore);

const notifications = useNotificationsStore();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const overlay = useOverlay();

const centrumCentrumClient = await getCentrumCentrumClient();

const logger = useLogger('⛑️ Centrum');

const canStream = can('centrum.CentrumService/Stream');

const selectedDispatch = ref<number | undefined>();

const dispatchStatusUpdateModal = overlay.create(DispatchStatusUpdateModal);
const unitStatusUpdateModal = overlay.create(UnitStatusUpdateModal);
const joinUnitSlideover = overlay.create(JoinUnitSlideover);
const takeDispatchSlideover = overlay.create(TakeDispatchSlideover);
const unitDetailsSlideover = overlay.create(UnitDetailsSlideover);

async function updateDispatchStatus(dispatchId: number, status: StatusDispatch): Promise<void> {
    try {
        const call = centrumCentrumClient.updateDispatchStatus({ dispatchId, status });
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
        dispatchStatusUpdateModal.open({
            dispatchId: dispatchId,
            status: status,
        });
        return;
    }

    await updateDispatchStatus(dispatchId, status);
}

async function updateUnitStatus(id: number, status: StatusUnit): Promise<void> {
    try {
        const call = centrumCentrumClient.updateUnitStatus({
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
        if (!getOwnUnit.value) return;

        unitStatusUpdateModal.open({
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
        overlay.closeAll();

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
        overlay.closeAll();
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
        const result = joinUnitSlideover.open({});
        result.then(() => {
            if (getOwnUnit.value === undefined) {
                // User closed the slideover without joining an unit, so close the sidebar again
                open.value = false;
            }
        });
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

const ownUnitStatus = computed(() => unitStatusToBadgeColor(getOwnUnit.value?.status?.status));

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
        if (!isStatusDispatchCompleted(dispatch?.status?.status ?? StatusDispatch.UNSPECIFIED)) return;
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
    if (!settings.value?.enabled) return;

    useIntervalFn(() => checkup(), 1 * 60 * 1000);
    toggleSidebarBasedOnUnit();
    toggleRequireUnitNotification();
});

onBeforeMount(async () => {
    if (!canStream.value) return;

    useTimeoutFn(async () => {
        try {
            startStream();
        } catch (e) {
            logger.error('exception during start of centrum stream', e);
        }
    }, 500);
});

onBeforeRouteLeave(async (to) => {
    // Don't end centrum stream if user is switching to center or livemap page
    if (to.path.startsWith('/livemap') || to.path === '/centrum') return;

    await stopStream(true);
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
    if (ownUnit === undefined || ownUnit.status === undefined) return;

    if (ownUnit.status.status === StatusUnit.AVAILABLE || ownUnit.status.status === StatusUnit.UNAVAILABLE) return;

    const now = new Date();
    // If unit status is younger than time X, ignore and continue
    if (now.getTime() - toDate(ownUnit.status.createdAt, timeCorrection.value).getTime() <= unitCheckupStatusAge) return;

    if (
        lastCheckupNotification.value !== undefined &&
        now.getTime() - lastCheckupNotification.value.getTime() <= unitCheckupStatusReping
    )
        return;

    notifications.add({
        title: { key: 'notifications.centrum.unitUpdated.checkup.title', parameters: {} },
        description: { key: 'notifications.centrum.unitUpdated.checkup.content', parameters: {} },
        type: NotificationType.INFO,
        duration: 15000,
        callback: () => attentionDebouncedPlay(),
    });

    lastCheckupNotification.value = now;
}

function sendRequireUnitNotification(): void {
    if (!userOnDuty.value) return;

    notifications.add({
        title: { key: 'notifications.centrum.unitUpdated.require_unit.title', parameters: {} },
        description: { key: 'notifications.centrum.unitUpdated.require_unit.content', parameters: {} },
        type: NotificationType.WARNING,
        duration: 12500,
    });

    attentionSound.play();
}

function openTakeDispatches(): void {
    takeDispatchSlideover.open({});
}

defineShortcuts({
    'm-d': () => getOwnUnit.value && openTakeDispatches(),
    'm-h': () => getOwnUnit.value?.homePostal && setWaypointPLZ(getOwnUnit.value.homePostal),
    'c-u': () => getOwnUnit.value && onSubmitUnitStatusThrottle(getOwnUnit.value.id),
    'c-d': () => getOwnUnit.value && onSubmitDispatchStatusThrottle(),
});
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.livemap')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <DispatcherInfo v-if="canStream && settings?.enabled" hide-join />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div class="relative z-0 size-full">
                <LivemapBase>
                    <template #default>
                        <template v-if="canStream">
                            <DispatchLayer
                                :show-all-dispatches="livemap.showAllDispatches || getCurrentMode === CentrumMode.SIMPLIFIED"
                            />

                            <LControl position="bottomright">
                                <UChip
                                    v-if="settings?.enabled"
                                    :show="getSortedOwnDispatches.length > 0"
                                    :text="getSortedOwnDispatches.length"
                                    color="error"
                                    size="lg"
                                    position="top-left"
                                >
                                    <UButton
                                        class="inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-clip-padding text-black"
                                        size="xs"
                                        :icon="open ? 'i-mdi-chevron-double-right' : 'i-mdi-chevron-double-left'"
                                        :color="!getOwnUnit ? 'primary' : 'neutral'"
                                        @click="open = !open"
                                    >
                                        <span v-if="!open" class="inline-flex items-center justify-center">
                                            {{ !getOwnUnit ? $t('common.unit', 2) : $t('common.your_dispatches') }}
                                        </span>
                                    </UButton>
                                </UChip>
                            </LControl>
                        </template>

                        <FollowMarker />
                    </template>

                    <template v-if="canStream && settings?.enabled" #afterMap>
                        <!-- "Take Dispatches" Button -->
                        <span v-if="getOwnUnit !== undefined" class="absolute right-1/2 bottom-2 z-30 inline-flex">
                            <UChip
                                :ui="{
                                    base: 'absolute rounded-full ring-0 ring-white dark:ring-gray-900 flex items-center justify-center text-white dark:text-gray-900 font-medium whitespace-nowrap animate-ping duration-750',
                                }"
                                position="top-left"
                                size="xl"
                                color="error"
                                :show="pendingDispatches.length > 0"
                            >
                                <UTooltip :text="$t('components.centrum.take_dispatch.title')" :kbds="['M', 'D']">
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
                                :kbds="['M', 'H']"
                            >
                                <UButton
                                    class="flex size-12 items-center justify-center rounded-r-full"
                                    size="xl"
                                    icon="i-mdi-home-floor-b"
                                    @click="setWaypointPLZ(getOwnUnit.homePostal)"
                                />
                            </UTooltip>
                        </span>
                    </template>
                </LivemapBase>
            </div>
        </template>
    </UDashboardPanel>

    <UDashboardPanel
        v-if="canStream && open"
        resizable
        :min-size="14.0"
        :max-size="25"
        :default-size="15.25"
        class="max-w-[25rem]"
        :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0 border-b border-default' }"
    >
        <template #header>
            <UDashboardNavbar :ui="{ root: 'px-1 sm:px-1', center: 'flex flex-1', toggle: 'hidden' }">
                <template #default>
                    <div class="flex flex-1 flex-col items-center">
                        <UButton
                            v-if="getOwnUnit !== undefined"
                            class="inline-flex flex-col rounded-b-none"
                            :color="ownUnitStatus"
                            icon="i-mdi-information-outline"
                            block
                            @click="
                                unitDetailsSlideover.open({
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
                                {{ $t(`enums.centrum.StatusUnit.${StatusUnit[getOwnUnit.status?.status ?? 0]}`) }}
                            </span>
                        </UButton>

                        <UFieldGroup class="w-full" orientation="vertical">
                            <UButton
                                :class="getOwnUnit !== undefined ? 'rounded-t-none' : ''"
                                variant="soft"
                                color="primary"
                                size="xs"
                                block
                                :icon="getOwnUnit === undefined ? 'i-mdi-information-outline' : undefined"
                                @click="joinUnitSlideover.open({})"
                            >
                                <span v-if="getOwnUnit === undefined" class="truncate">{{ $t('common.no_own_unit') }}</span>
                                <span v-else class="truncate">{{ $t('common.leave_unit') }}</span>
                            </UButton>

                            <UButton
                                v-if="getOwnUnit === undefined"
                                variant="solid"
                                color="success"
                                size="xs"
                                block
                                icon="i-mdi-account-plus"
                                :label="$t('common.join_unit')"
                                @click="joinUnitSlideover.open({})"
                            />
                        </UFieldGroup>
                    </div>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div class="overflow-x-hidden overflow-y-auto p-0 sm:pb-0">
                <div class="flex flex-1 flex-col gap-y-2" :class="open || getOwnUnit !== undefined ? 'px-1' : ''">
                    <template v-if="getOwnUnit !== undefined">
                        <ul role="list">
                            <li class="inline-flex items-center gap-1 text-xs leading-6 font-semibold">
                                <span>{{ $t('common.units') }}</span>
                                <UIcon v-if="!canSubmitUnitStatus" class="size-4 animate-spin" name="i-mdi-loading" />
                            </li>

                            <li>
                                <div class="grid grid-cols-2 gap-0.5">
                                    <UButton
                                        v-for="item in unitStatuses"
                                        :key="item.name"
                                        :color="unitStatusToBadgeColor(item.status)"
                                        size="xs"
                                        :disabled="!canSubmitUnitStatus"
                                        :icon="item.icon"
                                        truncate
                                        :label="
                                            item.status
                                                ? $t(`enums.centrum.StatusUnit.${StatusUnit[item.status ?? 0]}`)
                                                : $t(item.name)
                                        "
                                        :ui="{ label: 'line-clamp-2' }"
                                        @click="onSubmitUnitStatusThrottle(getOwnUnit.id!, item.status)"
                                    />

                                    <UTooltip
                                        class="col-span-2"
                                        :text="$t('components.centrum.update_unit_status.title')"
                                        :kbds="['S', 'U']"
                                    >
                                        <UButton
                                            variant="soft"
                                            color="primary"
                                            size="xs"
                                            block
                                            :label="$t('components.centrum.update_unit_status.title')"
                                            @click="onSubmitUnitStatusThrottle(getOwnUnit.id)"
                                        />
                                    </UTooltip>
                                </div>
                            </li>
                        </ul>

                        <USeparator class="my-0.25" />

                        <ul class="" role="list">
                            <li class="inline-flex items-center gap-1 text-xs leading-6 font-semibold">
                                <span>{{ $t('common.dispatch') }} {{ $t('common.status') }}</span>
                                <UIcon v-if="!canSubmitDispatchStatus" class="size-4 animate-spin" name="i-mdi-loading" />
                            </li>

                            <li>
                                <div class="grid grid-cols-2 gap-0.5">
                                    <UButton
                                        v-for="item in dispatchStatuses.filter((s) => s.status !== StatusDispatch.CANCELLED)"
                                        :key="item.name"
                                        :color="dispatchStatusToBadgeColor(item.status)"
                                        size="xs"
                                        :disabled="!canSubmitDispatchStatus"
                                        :icon="item.icon"
                                        :ui="{ label: 'line-clamp-2' }"
                                        @click="onSubmitDispatchStatusThrottle(selectedDispatch, item.status)"
                                    >
                                        <span class="line-clamp-2">
                                            {{
                                                item.status
                                                    ? $t(`enums.centrum.StatusDispatch.${StatusDispatch[item.status ?? 0]}`)
                                                    : $t(item.name)
                                            }}
                                        </span>
                                    </UButton>

                                    <UTooltip
                                        class="col-span-2"
                                        :text="$t('components.centrum.update_dispatch_status.title')"
                                        :kbds="['S', 'D']"
                                    >
                                        <UButton
                                            variant="soft"
                                            color="primary"
                                            size="xs"
                                            block
                                            @click="updateDspStatus(selectedDispatch)"
                                        >
                                            {{ $t('components.centrum.update_dispatch_status.title') }}
                                        </UButton>
                                    </UTooltip>
                                </div>
                            </li>
                        </ul>

                        <ul role="list">
                            <li class="inline-flex items-center text-xs leading-6 font-semibold">
                                {{ $t('common.your_dispatches') }}
                            </li>

                            <li v-if="getSortedOwnDispatches.length === 0">
                                <UButton variant="soft" color="neutral" icon="i-mdi-car-emergency" size="xs" block>
                                    {{ $t('common.no_assigned_dispatches') }}
                                </UButton>
                            </li>

                            <template v-else>
                                <template v-for="dispatch in getSortedOwnDispatches.toReversed()" :key="dispatch">
                                    <OwnDispatchEntry
                                        v-if="dispatches.get(dispatch) !== undefined"
                                        v-model="selectedDispatch"
                                        :dispatch="dispatches.get(dispatch)!"
                                    />
                                </template>
                            </template>
                        </ul>
                    </template>
                </div>
            </div>
        </template>

        <template #footer>
            <div class="mx-2 my-1">
                <DispatchStatusBreakdown block popover-class="w-full" size="xs" />
            </div>
        </template>
    </UDashboardPanel>
</template>
