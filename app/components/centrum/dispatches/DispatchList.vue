<script lang="ts" setup>
import { UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { h } from 'vue';
import DispatchAssignModal from '~/components/centrum/dispatches/DispatchAssignModal.vue';
import DispatchDetailsByIDSlideover from '~/components/centrum/dispatches/DispatchDetailsByIDSlideover.vue';
import DispatchStatusUpdateModal from '~/components/centrum/dispatches/DispatchStatusUpdateModal.vue';
import { checkDispatchAccess, dispatchTimeToBadge } from '~/components/centrum/helpers';
import DispatchAttributes from '~/components/centrum/partials/DispatchAttributes.vue';
import DispatchStatusBreakdown from '~/components/centrum/partials/DispatchStatusBreakdown.vue';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { useCentrumStore } from '~/stores/centrum';
import { useLivemapStore } from '~/stores/livemap';
import { CentrumAccessLevel } from '~~/gen/ts/resources/centrum/access';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';
import DispatchStatusBadge from '../partials/DispatchStatusBadge.vue';

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

const overlay = useOverlay();

const { gotoCoords } = useLivemapStore();

const centrumStore = useCentrumStore();
const { getSortedDispatches, settings, abort, stopping } = storeToRefs(centrumStore);

const settingsStore = useSettingsStore();
const { centrum } = storeToRefs(settingsStore);

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
        id: 'actions',
        meta: {
            class: {
                td: 'py-px',
            },
        },
        cell: ({ row }) =>
            h('div', { class: 'flex items-center' }, [
                h(
                    UTooltip,
                    {
                        text: t('common.assign'),
                        vIf: !props.hideActions,
                    },
                    [
                        h(UButton, {
                            variant: 'link',
                            icon: 'i-mdi-account-multiple-plus',
                            ui: { base: 'px-0.5 py-1' },
                            disabled: !checkDispatchAccess(row.original.jobs, CentrumAccessLevel.DISPATCH),
                            onClick: () => {
                                dispatchAssignModal.open({
                                    dispatchId: row.original.id,
                                });
                            },
                        }),
                    ],
                ),
                h(
                    UTooltip,
                    {
                        text: t('common.go_to_location'),
                    },
                    [
                        h(UButton, {
                            variant: 'link',
                            icon: 'i-mdi-map-marker',
                            onClick: () => gotoCoords({ x: row.original.x, y: row.original.y }),
                        }),
                    ],
                ),
                h(
                    UTooltip,
                    {
                        text: t('common.status'),
                        vIf: !props.hideActions,
                    },
                    [
                        h(UButton, {
                            variant: 'link',
                            icon: 'i-mdi-refresh',
                            disabled: !checkDispatchAccess(row.original.jobs, CentrumAccessLevel.DISPATCH),
                            onClick: () => {
                                dispatchStatusUpdateModal.open({
                                    dispatchId: row.original.id,
                                });
                            },
                        }),
                    ],
                ),
                h(
                    UTooltip,
                    {
                        text: t('common.detail', 2),
                    },
                    [
                        h(UButton, {
                            variant: 'link',
                            icon: 'i-mdi-dots-vertical',
                            onClick: () => {
                                dispatchDetailsSlideover.open({
                                    dispatchId: row.original.id,
                                });
                            },
                        }),
                    ],
                ),
            ]),
    },
    {
        accessorKey: 'id',
        header: t('common.id'),
        meta: {
            class: {
                td: 'h-full text-left',
            },
        },
        cell: ({ row }) => h(IDCopyBadge, { id: row.original.id, prefix: 'DSP', disableTooltip: true, variant: 'link' }),
    },
    {
        accessorKey: 'createdAt',
        header: t('common.created'),
        meta: {
            class: {
                td: 'h-full text-center',
            },
        },
        cell: ({ row }) =>
            h(GenericTime, {
                value: row.original.createdAt,
                type: 'compact',
                badge: true,
                size: 'sm',
                updateCallback: () =>
                    dispatchTimeToBadge(
                        row.original.createdAt,
                        row.original.status?.status,
                        settings.value?.timings?.dispatchMaxWait,
                    ),
            }),
    },
    {
        accessorKey: 'status',
        header: t('common.status'),
        meta: {
            class: {
                td: 'h-full text-center',
            },
        },
        cell: ({ row }) => h(DispatchStatusBadge, { size: 'sm', status: row.original.status?.status }),
    },
    {
        accessorKey: 'postal',
        header: t('common.postal'),
        meta: {
            class: {
                td: 'h-full text-center',
            },
        },
        cell: ({ row }) => h('span', {}, row.original.postal ?? t('common.na')),
    },
    {
        accessorKey: 'units',
        header: t('common.unit'),
        cell: ({ row }) =>
            row.original.units.length === 0
                ? h('span', { class: 'italic' }, t('enums.centrum.StatusDispatch.UNASSIGNED'))
                : h(
                      'span',
                      { class: 'grid grid-flow-row auto-rows-auto gap-1 sm:grid-flow-col' },
                      row.original.units.map((unit) =>
                          h(UnitInfoPopover, {
                              unit: unit.unit,
                              initialsOnly: true,
                              badge: true,
                              assignment: unit,
                          }),
                      ),
                  ),
    },
    {
        accessorKey: 'creator',
        header: t('common.creator'),
        cell: ({ row }) =>
            row.original.anon
                ? h('span', {}, t('common.anon'))
                : row.original.creator
                  ? h(CitizenInfoPopover, { user: row.original.creator, trailing: false })
                  : h('span', {}, t('common.unknown')),
    },
    {
        accessorKey: 'attributes',
        header: t('common.attributes', 2),
        cell: ({ row }) => h(DispatchAttributes, { attributes: row.original.attributes }),
    },
    {
        accessorKey: 'message',
        header: t('common.message'),
        cell: ({ row }) => h('p', { class: 'line-clamp-2 hover:line-clamp-6' }, row.original.message),
    },
] as TableColumn<Dispatch>[];

const dispatchAssignModal = overlay.create(DispatchAssignModal);
const dispatchStatusUpdateModal = overlay.create(DispatchStatusUpdateModal);
const dispatchDetailsSlideover = overlay.create(DispatchDetailsByIDSlideover);
</script>

<template>
    <div class="flex h-full flex-1 grow flex-col">
        <div class="flex justify-between px-1">
            <h2 class="inline-flex flex-1 items-center text-base leading-6 font-semibold">
                {{ $t('common.dispatches') }}

                <UTooltip v-if="showButton" :text="$t('common.archive')">
                    <UButton :to="{ name: 'centrum-dispatches' }" icon="i-mdi-archive" variant="link" />
                </UTooltip>
            </h2>

            <UFormField class="grid grid-cols-2 items-center gap-2" name="cards" :label="$t('common.card_view')">
                <div class="flex flex-1 items-center">
                    <USwitch v-model="centrum.dispatchListCardStyle" />
                </div>
            </UFormField>

            <DispatchStatusBreakdown v-if="dispatches === undefined" class="justify-end font-semibold text-toned" />
        </div>

        <div class="flex flex-1 flex-col overflow-x-auto overflow-y-auto">
            <div v-if="!dispatches && abort === undefined && stopping" class="space-y-1">
                <USkeleton v-for="idx in 7" :key="idx" class="h-9 w-full" />
            </div>

            <template v-for="(group, idx) in grouped" v-else :key="group.key">
                <h3 v-if="alwaysShowDay || idx !== 0"><GenericTime :value="group.date" type="date" /></h3>

                <UTable
                    v-if="!centrum.dispatchListCardStyle"
                    class="overflow-x-visible"
                    :columns="columns"
                    :data="group.dispatches"
                    :empty="$t('common.not_found', [$t('common.dispatch', 2)])"
                    :sorting-options="{ manualSorting: true }"
                    :pagination-options="{ manualPagination: true }"
                />

                <div v-else class="grid grid-cols-1 gap-1 p-1">
                    <UCard
                        v-for="dispatch in group.dispatches"
                        :key="dispatch.id"
                        :title="dispatch.message"
                        class="px-px"
                        :ui="{ header: 'p-1 sm:px-1', body: 'p-1 sm:p-1', footer: 'p-1 sm:px-1' }"
                    >
                        <template #header>
                            <div class="flex items-center justify-between">
                                <div class="flex flex-1 items-center gap-2">
                                    <div>
                                        <UTooltip v-if="!hideActions" :text="$t('common.assign')">
                                            <UButton
                                                variant="link"
                                                icon="i-mdi-account-multiple-plus"
                                                :ui="{ base: 'px-0.5 py-1' }"
                                                @click="
                                                    () => {
                                                        dispatchAssignModal.open({
                                                            dispatchId: dispatch.id,
                                                        });
                                                    }
                                                "
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('common.go_to_location')">
                                            <UButton
                                                variant="link"
                                                icon="i-mdi-map-marker"
                                                @click="() => gotoCoords({ x: dispatch.x, y: dispatch.y })"
                                            />
                                        </UTooltip>

                                        <UTooltip v-if="!hideActions" :text="$t('common.status')">
                                            <UButton
                                                variant="link"
                                                icon="i-mdi-refresh"
                                                :ui="{ base: 'px-0.5 py-1' }"
                                                @click="
                                                    () => {
                                                        dispatchStatusUpdateModal.open({
                                                            dispatchId: dispatch.id,
                                                        });
                                                    }
                                                "
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('common.detail', 2)">
                                            <UButton
                                                variant="link"
                                                icon="i-mdi-dots-vertical"
                                                :ui="{ base: 'px-0.5 py-1' }"
                                                @click="
                                                    () => {
                                                        dispatchDetailsSlideover.open({
                                                            dispatchId: dispatch.id,
                                                        });
                                                    }
                                                "
                                            />
                                        </UTooltip>
                                    </div>

                                    <div class="flex flex-1 items-center justify-center gap-2">
                                        <GenericTime
                                            :value="dispatch.createdAt"
                                            type="compact"
                                            badge
                                            size="sm"
                                            class="text-highlighted"
                                            :update-callback="
                                                () =>
                                                    dispatchTimeToBadge(
                                                        dispatch.createdAt,
                                                        dispatch.status?.status,
                                                        settings?.timings?.dispatchMaxWait,
                                                    )
                                            "
                                        />

                                        <DispatchStatusBadge :status="dispatch.status?.status" size="sm" />
                                    </div>

                                    <span class="text-sm">
                                        {{ $t('common.postal') }}: {{ dispatch.postal ?? $t('common.na') }}
                                    </span>

                                    <IDCopyBadge :id="dispatch.id" prefix="DSP" disable-tooltip variant="link" />
                                </div>
                            </div>
                        </template>

                        <div class="flex flex-col gap-1">
                            <div class="flex flex-row">
                                <p class="line-clamp-2 flex-1 hover:line-clamp-6">
                                    <span class="mr-1 font-semibold">{{ $t('common.message') }}:</span>
                                    <span>{{ dispatch.message }}</span>
                                </p>

                                <p class="inline-flex items-center gap-1">
                                    <span class="mr-1 font-semibold">{{ $t('common.creator') }}:</span>
                                    <span v-if="dispatch.anon">
                                        {{ $t('common.anon') }}
                                    </span>
                                    <span v-else-if="dispatch.creator">
                                        <CitizenInfoPopover :user="dispatch.creator" />
                                    </span>
                                    <span v-else>
                                        {{ $t('common.unknown') }}
                                    </span>
                                </p>
                            </div>
                        </div>

                        <template #footer>
                            <div class="flex items-center justify-between gap-2">
                                <span v-if="dispatch.units.length === 0" class="italic">{{
                                    $t('enums.centrum.StatusDispatch.UNASSIGNED')
                                }}</span>
                                <span v-else class="grid grid-flow-row auto-rows-auto gap-1 sm:grid-flow-col">
                                    <UnitInfoPopover
                                        v-for="unit in dispatch.units"
                                        :key="unit.unitId"
                                        :unit="unit.unit"
                                        :assignment="unit"
                                        initials-only
                                        badge
                                    />
                                </span>

                                <div class="truncate">
                                    {{ dispatch.jobs?.jobs.map((j) => j.label ?? j.name).join(',') }}
                                </div>

                                <DispatchAttributes :attributes="dispatch.attributes" />
                            </div>
                        </template>
                    </UCard>
                </div>
            </template>

            <div class="flex-1" />
        </div>
    </div>
</template>
