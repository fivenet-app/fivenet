<script lang="ts" setup>
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

const modal = useModal();

const slideover = useSlideover();

const { goto } = useLivemapStore();

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
    <div class="flex h-full flex-1 grow flex-col px-1">
        <div class="flex justify-between">
            <h2 class="inline-flex flex-1 items-center text-base font-semibold leading-6">
                {{ $t('common.dispatches') }}

                <UTooltip v-if="showButton" :text="$t('common.archive')">
                    <UButton :to="{ name: 'centrum-dispatches' }" icon="i-mdi-archive" variant="link" />
                </UTooltip>
            </h2>

            <UFormGroup class="grid grid-cols-2 items-center gap-2" name="cards" :label="$t('common.card_view')">
                <div class="flex flex-1 items-center">
                    <UToggle v-model="centrum.dispatchListCardStyle" />
                </div>
            </UFormGroup>

            <DispatchStatusBreakdown v-if="dispatches === undefined" class="justify-end font-semibold text-gray-100" />
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
                                    :disabled="!checkDispatchAccess(dispatch.jobs, CentrumAccessLevel.DISPATCH)"
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
                                    icon="i-mdi-refresh"
                                    :disabled="!checkDispatchAccess(dispatch.jobs, CentrumAccessLevel.DISPATCH)"
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
                        <span class="text-gray-900 dark:text-white">
                            <GenericTime
                                :value="dispatch.createdAt"
                                type="compact"
                                badge
                                size="xs"
                                :update-callback="
                                    () =>
                                        dispatchTimeToBadge(
                                            dispatch.createdAt,
                                            dispatch.status?.status,
                                            settings?.timings?.dispatchMaxWait,
                                        )
                                "
                            />
                        </span>
                    </template>

                    <template #status-data="{ row: dispatch }">
                        <DispatchStatusBadge :status="dispatch.status?.status" />
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

                <div v-else class="grid grid-cols-1 gap-2">
                    <UCard
                        v-for="dispatch in group.dispatches"
                        :key="dispatch.id"
                        :title="dispatch.message"
                        class="px-px"
                        :ui="{
                            header: {
                                padding: 'px-2 py-1 sm:px-2',
                            },
                            body: {
                                padding: 'px-2 py-1 sm:px-2 sm:p-1',
                            },
                            footer: {
                                padding: 'px-2 py-1 sm:px-2',
                            },
                        }"
                    >
                        <template #header>
                            <div class="flex items-center justify-between">
                                <div class="flex flex-1 items-center gap-2">
                                    <div>
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
                                                icon="i-mdi-refresh"
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

                                    <div class="flex flex-1 items-center justify-center gap-2">
                                        <span class="text-gray-900 dark:text-white">
                                            <GenericTime
                                                :value="dispatch.createdAt"
                                                type="compact"
                                                badge
                                                size="sm"
                                                :update-callback="
                                                    () =>
                                                        dispatchTimeToBadge(
                                                            dispatch.createdAt,
                                                            dispatch.status?.status,
                                                            settings?.timings?.dispatchMaxWait,
                                                        )
                                                "
                                            />
                                        </span>

                                        <DispatchStatusBadge :status="dispatch.status?.status" />
                                    </div>

                                    <span> {{ $t('common.postal') }}: {{ dispatch.postal ?? $t('common.na') }} </span>

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
                                        :initials-only="true"
                                        :badge="true"
                                        :assignment="unit"
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
