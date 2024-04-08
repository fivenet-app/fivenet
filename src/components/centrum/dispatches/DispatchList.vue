<script lang="ts" setup>
import { useCentrumStore } from '~/store/centrum';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import DispatchStatusBreakdown from '../partials/DispatchStatusBreakdown.vue';
import { dispatchStatusAnimate, dispatchStatusToBGColor } from '../helpers';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import UnitInfoPopover from '../units/UnitInfoPopover.vue';
import DispatchAttributes from '../partials/DispatchAttributes.vue';
import DispatchDetailsSlideover from './DispatchDetailsSlideover.vue';
import DispatchStatusUpdateModal from './DispatchStatusUpdateModal.vue';
import DispatchAssignModal from './DispatchAssignModal.vue';
import { useSettingsStore } from '~/store/settings';

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

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const { t } = useI18n();

const modal = useModal();

const slideover = useSlideover();

const centrumStore = useCentrumStore();
const { getSortedDispatches } = storeToRefs(centrumStore);

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
            groups[idx].dispatches.push(e);
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
    },
    {
        key: 'creator',
        label: t('common.creator'),
    },
    {
        key: 'attributes',
        label: t('common.attributes', 2),
    },
    {
        key: 'message',
        label: t('common.message'),
    },
];
</script>

<template>
    <div class="flex size-full grow flex-col overflow-y-auto px-1">
        <div class="flex justify-between">
            <h2 class="inline-flex flex-1 items-center text-base font-semibold leading-6">
                {{ $t('common.dispatches') }}

                <UButton
                    v-if="showButton"
                    :to="{ name: 'centrum-dispatches' }"
                    :title="$t('common.dispatches')"
                    icon="i-mdi-archive"
                    variant="link"
                />
            </h2>

            <DispatchStatusBreakdown v-if="dispatches === undefined" class="text-base font-semibold text-gray-100" />
        </div>

        <div class="flex-1">
            <template v-for="(group, idx) in grouped" :key="group.key">
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
                        <UButtonGroup class="inline-flex w-full">
                            <UButton
                                v-if="!hideActions"
                                variant="link"
                                icon="i-mdi-account-multiple-plus"
                                :padded="false"
                                :title="$t('common.assign')"
                                @click="
                                    modal.open(DispatchAssignModal, {
                                        dispatch: dispatch,
                                    })
                                "
                            />

                            <UButton
                                size="xs"
                                variant="link"
                                icon="i-mdi-map-marker"
                                :title="$t('common.go_to_location')"
                                @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
                            />

                            <UButton
                                v-if="!hideActions"
                                variant="link"
                                icon="i-mdi-close-octagon"
                                :padded="false"
                                :title="$t('common.status')"
                                @click="
                                    modal.open(DispatchStatusUpdateModal, {
                                        dispatchId: dispatch.id,
                                    })
                                "
                            />

                            <UButton
                                variant="link"
                                icon="i-mdi-dots-vertical"
                                :padded="false"
                                :title="$t('common.detail', 2)"
                                @click="
                                    slideover.open(DispatchDetailsSlideover, {
                                        dispatch: dispatch,
                                        onGoto: (loc) => $emit('goto', loc),
                                    })
                                "
                            />
                        </UButtonGroup>
                    </template>
                    <template #createdAt-data="{ row: dispatch }">
                        <GenericTime :value="dispatch.createdAt" type="compact" />
                    </template>
                    <template #status-data="{ row: dispatch }">
                        <span
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
