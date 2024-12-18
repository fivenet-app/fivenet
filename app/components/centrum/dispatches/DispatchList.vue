<script lang="ts" setup>
import DispatchAssignModal from '~/components/centrum/dispatches/DispatchAssignModal.vue';
import DispatchDetailsByIDSlideover from '~/components/centrum/dispatches/DispatchDetailsByIDSlideover.vue';
import DispatchStatusUpdateModal from '~/components/centrum/dispatches/DispatchStatusUpdateModal.vue';
import { dispatchStatusAnimate, dispatchStatusToBGColor, dispatchTimeToTextColor } from '~/components/centrum/helpers';
import DispatchAttributes from '~/components/centrum/partials/DispatchAttributes.vue';
import DispatchStatusBreakdown from '~/components/centrum/partials/DispatchStatusBreakdown.vue';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

const props = withDefaults(
    defineProps<{
        dispatches?: Dispatch[];
        showButton?: boolean;
        hideActions?: boolean;
        alwaysShowDay?: boolean;
    }>(),
    {
        dispatches: undefined,
        showButton: false,
        hideActions: false,
        alwaysShowDay: false,
    },
);

const { t } = useI18n();

const modal = useModal();

const slideover = useSlideover();

const { goto } = useLivemapStore();

const centrumStore = useCentrumStore();
const { getSortedDispatches, settings, abort, reconnecting } = storeToRefs(centrumStore);

type GroupedDispatches = { date: Date; key: string; dispatches: Dispatch[] }[];

const grouped = computedAsync(async () => {
    const groups: GroupedDispatches = [];
    (props.dispatches ?? getSortedDispatches.value).forEach((e) => {
        const date = toDate(e.createdAt);
        const idx = groups.findIndex((g) => g.key === dateToDateString(date));
        if (idx === -1) {
            groups.push({
                date,
                dispatches: [e],
                key: dateToDateString(date),
            });
        } else {
            groups[idx]!.dispatches.push(e);
        }
    });

    return groups;
});

const columns = [
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
    {
        key: 'id',
        label: t('common.id'),
    },
    {
        key: 'createdAt',
        label: t('common.created'),
    },
    {
        key: 'status',
        label: t('common.status'),
    },
    {
        key: 'postal',
        label: t('common.postal'),
    },
    {
        key: 'units',
        label: t('common.unit'),
        sortable: false,
    },
    {
        key: 'creator',
        label: t('common.creator'),
        sortable: false,
    },
    {
        key: 'attributes',
        label: t('common.attributes', 2),
        sortable: false,
    },
    {
        key: 'message',
        label: t('common.message'),
        sortable: false,
    },
];
</script>

<template>
    <div class="flex size-full grow flex-col overflow-y-auto px-1">
        <div class="flex justify-between">
            <h2 class="inline-flex flex-1 items-center text-base font-semibold leading-6">
                {{ $t('common.dispatches') }}

                <UTooltip v-if="showButton" :text="$t('common.dispatches')">
                    <UButton
                        :to="{ name: 'centrum-dispatches' }"
                        icon="i-mdi-archive"
                        variant="link"
                    />
                </UTooltip>
            </h2>

            <DispatchStatusBreakdown v-if="dispatches === undefined" class="font-semibold text-gray-100" />
        </div>

        <div class="flex-1">
            <div v-if="!dispatches && abort === undefined && !reconnecting" class="space-y-1">
                <USkeleton v-for="idx in 7" :key="idx" class="h-9 w-full" />
            </div>

            <template v-for="(group, idx) in grouped" v-else :key="group.key">
                <h3 v-if="alwaysShowDay || idx !== 0"><GenericTime :value="group.date" type="date" /></h3>
                <UTable
                    :columns="columns"
                    :rows="group.dispatches"
                    :empty-state="{
                        icon: 'i-mdi-car-emergency',
                        label: $t('common.not_found', [$t('common.dispatch', 2)]),
                    }"
                    :ui="{ th: { padding: 'px-0.5 py-0.5' }, td: { padding: 'px-1 py-0.5' } }"
                >
                    <template #actions-data="{ row: dispatch }">
                        <div :key="dispatch.id">
                            <UTooltip v-if="!hideActions" :text="$t('common.assign')">
                                <UButton
                                    variant="link"
                                    icon="i-mdi-account-multiple-plus"
                                    @click="
                                        () =>
                                            modal.open(DispatchAssignModal, {
                                                dispatchId: dispatch.id,
                                            })
                                    "
                                />
                            </UTooltip>

                            <UTooltip :text="$t('common.go_to_location')">
                                <UButton
                                    variant="link"
                                    icon="i-mdi-map-marker"
                                    @click="() => goto({ x: dispatch.x, y: dispatch.y })"
                                />
                            </UTooltip>

                            <UTooltip v-if="!hideActions" :text="$t('common.status')">
                                <UButton
                                    variant="link"
                                    icon="i-mdi-close-octagon"
                                    @click="
                                        () =>
                                            modal.open(DispatchStatusUpdateModal, {
                                                dispatchId: dispatch.id,
                                            })
                                    "
                                />
                            </UTooltip>

                            <UTooltip :text="$t('common.detail', 2)">
                                <UButton
                                    variant="link"
                                    icon="i-mdi-dots-vertical"
                                    @click="
                                        () =>
                                            slideover.open(DispatchDetailsByIDSlideover, {
                                                dispatchId: dispatch.id,
                                            })
                                    "
                                />
                            </UTooltip>
                        </div>
                    </template>
                    <template #createdAt-data="{ row: dispatch }">
                        <GenericTime
                            :value="dispatch.createdAt"
                            type="compact"
                            :update-callback="
                                () =>
                                    dispatchTimeToTextColor(
                                        dispatch.createdAt,
                                        dispatch.status.status,
                                        settings?.timings?.dispatchMaxWait,
                                    )
                            "
                        />
                    </template>
                    <template #status-data="{ row: dispatch }">
                        <span
                            class="text-gray-900 dark:text-white"
                            :class="[
                                dispatchStatusToBGColor(dispatch.status?.status),
                                dispatchStatusAnimate(dispatch.status?.status) ? 'animate-pulse' : '',
                            ]"
                        >
                            {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`) }}
                        </span>
                    </template>
                    <template #postal-data="{ row: dispatch }">
                        {{ dispatch.postal ?? $t('common.na') }}
                    </template>
                    <template #units-data="{ row: dispatch }">
                        <span v-if="dispatch.units.length === 0" class="italic">{{
                            $t('enums.centrum.StatusDispatch.UNASSIGNED')
                        }}</span>
                        <span v-else class="grid grid-flow-row auto-rows-auto gap-1 sm:grid-flow-col">
                            <UnitInfoPopover
                                v-for="unit in dispatch.units"
                                :key="unit.unitId"
                                :unit="unit.unit"
                                :initials-only="true"
                                :badge="true"
                                :assignment="unit"
                            />
                        </span>
                    </template>
                    <template #creator-data="{ row: dispatch }">
                        <span v-if="dispatch.anon">
                            {{ $t('common.anon') }}
                        </span>
                        <span v-else-if="dispatch.creator">
                            <CitizenInfoPopover :user="dispatch.creator" :trailing="false" />
                        </span>
                        <span v-else>
                            {{ $t('common.unknown') }}
                        </span>
                    </template>
                    <template #attributes-data="{ row: dispatch }">
                        <DispatchAttributes :attributes="dispatch.attributes" />
                    </template>
                    <template #message-data="{ row: dispatch }">
                        <p class="line-clamp-2 hover:line-clamp-6">
                            {{ dispatch.message }}
                        </p>
                    </template>
                </UTable>
            </template>
        </div>
    </div>
</template>
