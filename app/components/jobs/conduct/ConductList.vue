<script lang="ts" setup>
import { UBadge, UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { h } from 'vue';
import { z } from 'zod';
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import ConductCreateOrUpdateModal from '~/components/jobs/conduct/ConductCreateOrUpdateModal.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCompletorStore } from '~/stores/completor';
import { getJobsConductClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { type ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import type { ListConductEntriesResponse } from '~~/gen/ts/services/jobs/conduct';
import ColleagueName from '../colleagues/ColleagueName.vue';
import ConductViewSlideover from './ConductViewSlideover.vue';
import { conductTypesToBadgeColor } from './helpers';

const props = defineProps<{
    userId?: number;
    hideUserSearch?: boolean;
}>();

const { t } = useI18n();

const { can } = useAuth();

const overlay = useOverlay();

const completorStore = useCompletorStore();

const jobsConductClient = await getJobsConductClient();

const availableTypes = ref<{ status: ConductType }[]>([
    { status: ConductType.NOTE },
    { status: ConductType.NEUTRAL },
    { status: ConductType.POSITIVE },
    { status: ConductType.NEGATIVE },
    { status: ConductType.WARNING },
    { status: ConductType.SUSPENSION },
]);

const schema = z.object({
    id: z.union([z.string().optional(), z.coerce.number().min(1).optional()]),
    types: z.nativeEnum(ConductType).array().max(10).default([]),
    showExpired: z.coerce.boolean().default(false),
    user: z.coerce.number().min(1).optional(),
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'id',
                        desc: true,
                    },
                ]),
        })
        .default({ columns: [{ id: 'id', desc: true }] }),
    page: pageNumberSchema,
});

const query = useSearchForm('jobs_conduct', schema);

const { data, status, refresh, error } = useLazyAsyncData(
    () =>
        `jobs-conduct-${JSON.stringify(query.sorting)}-${query.page}-${query.types.join(',')}-${query.showExpired}-${query.id}`,
    () => listConductEntries(),
);

async function listConductEntries(): Promise<ListConductEntriesResponse> {
    const entryIds: number[] = [];
    if (query.id) {
        entryIds.push(typeof query.id === 'string' ? parseInt(query.id, 10) : query.id);
    }

    const userIds = props.userId ? [props.userId] : query.user ? [query.user] : [];
    try {
        const call = jobsConductClient.listConductEntries({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sorting,
            types: query.types,
            userIds: userIds,
            showExpired: query.showExpired,
            ids: entryIds,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteConductEntry(id: number): Promise<void> {
    try {
        const call = jobsConductClient.deleteConductEntry({ id });
        await call;

        refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

async function updateEntryInPlace(entry: ConductEntry): Promise<void> {
    if (data.value === null) {
        return refresh();
    }

    const idx = data.value?.entries.findIndex((e) => e.id === entry.id);
    if (idx !== undefined && idx > -1) {
        data.value?.entries.splice(idx, 1, entry);
    }

    refresh();
}

const conductViewSlideover = overlay.create(ConductViewSlideover);
const conductCreateOrUpdateModal = overlay.create(ConductCreateOrUpdateModal);
const confirmModal = overlay.create(ConfirmModal);

const appConfig = useAppConfig();

const columns = computed(
    () =>
        [
            {
                accessorKey: 'id',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.id'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                sortable: true,
            },
            {
                accessorKey: 'createdAt',
                header: t('common.created_at'),
                cell: ({ row }) => h(GenericTime, { value: row.original.createdAt }),
            },
            {
                accessorKey: 'expiresAt',
                header: t('common.expires_at'),
                meta: {
                    class: {
                        td: 'hidden lg:table-cell',
                        th: 'hidden lg:table-cell',
                    },
                },
                cell: ({ row }) =>
                    row.original.expiresAt
                        ? h(GenericTime, { value: row.original.expiresAt, type: 'date', class: 'font-semibold' })
                        : t('components.jobs.conduct.List.no_expiration'),
            },
            {
                accessorKey: 'type',
                header: t('common.type'),
                cell: ({ row }) =>
                    h(UBadge, { color: conductTypesToBadgeColor(row.original.type) }, () =>
                        t(`enums.jobs.ConductType.${ConductType[row.original.type ?? 0]}`),
                    ),
                sortable: true,
            },
            {
                accessorKey: 'message',
                header: t('common.message'),
                cell: ({ row }) =>
                    h(
                        'p',
                        { class: 'line-clamp-2 w-full max-w-sm break-all whitespace-normal hover:line-clamp-6' },
                        row.original.message,
                    ),
            },
            {
                accessorKey: 'target',
                header: t('common.target'),
                cell: ({ row }) => h(ColleagueInfoPopover, { user: row.original.targetUser }),
            },
            {
                accessorKey: 'creator',
                header: t('common.creator'),
                meta: {
                    class: {
                        td: 'hidden lg:table-cell',
                        th: 'hidden lg:table-cell',
                    },
                },
                cell: ({ row }) => h(ColleagueInfoPopover, { user: row.original.creator, hideProps: true }),
            },
            {
                id: 'actions',
                cell: ({ row }) =>
                    h('div', [
                        h(UTooltip, { text: t('common.show') }, () =>
                            h(UButton, {
                                variant: 'link',
                                icon: 'i-mdi-eye',
                                onClick: () => {
                                    conductViewSlideover.open({
                                        entry: row.original,
                                        onRefresh: async () => refresh(),
                                    });
                                },
                            }),
                        ),
                        can('jobs.ConductService/UpdateConductEntry').value &&
                            h(UTooltip, { text: t('common.update') }, () =>
                                h(UButton, {
                                    variant: 'link',
                                    icon: 'i-mdi-pencil',
                                    onClick: () => {
                                        conductCreateOrUpdateModal.open({
                                            entry: row.original,
                                            userId: props.userId,
                                            onCreated: ($event) => data.value?.entries.unshift($event),
                                            onUpdated: ($event) => updateEntryInPlace($event),
                                        });
                                    },
                                }),
                            ),
                        can('jobs.ConductService/DeleteConductEntry').value &&
                            h(UTooltip, { text: t('common.delete') }, () =>
                                h(UButton, {
                                    variant: 'link',
                                    icon: 'i-mdi-delete',
                                    color: 'error',
                                    onClick: () => {
                                        confirmModal.open({
                                            confirm: async () => deleteConductEntry(row.original.id),
                                        });
                                    },
                                }),
                            ),
                    ]),
            },
        ] as TableColumn<ConductEntry>[],
);
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardToolbar>
                <template #default>
                    <UForm class="my-2 w-full" :schema="schema" :state="query" @submit="refresh()">
                        <div class="flex flex-row gap-2">
                            <UFormField v-if="hideUserSearch !== true" class="flex-1" name="user" :label="$t('common.search')">
                                <SelectMenu
                                    ref="input"
                                    v-model="query.user"
                                    :searchable="
                                        async (q: string) => {
                                            const colleagues = await completorStore.listColleagues({
                                                search: q,
                                                labelIds: [],
                                                userIds: query.user ? [query.user] : [],
                                            });
                                            return colleagues;
                                        }
                                    "
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :filter-fields="['firstname', 'lastname']"
                                    block
                                    :placeholder="$t('common.colleague')"
                                    trailing
                                    leading-icon="i-mdi-search"
                                    value-key="userId"
                                >
                                    <template #item-label="{ item }">
                                        <span v-if="item" class="truncate">
                                            {{ userToLabel(item) }}
                                        </span>
                                    </template>

                                    <template #item="{ item }">
                                        <ColleagueName class="truncate" :colleague="item" birthday />
                                    </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                                    </template>
                                </SelectMenu>
                            </UFormField>

                            <UFormField class="flex-1" name="types" :label="$t('common.type')">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="query.types"
                                        multiple
                                        nullable
                                        :items="availableTypes"
                                        value-key="status"
                                        :placeholder="$t('common.na')"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                    >
                                        <template #item-label>
                                            {{ $t('common.selected', query.types.length) }}
                                        </template>

                                        <template #item="{ item }">
                                            <UBadge class="truncate" :color="conductTypesToBadgeColor(item.status)">
                                                {{ $t(`enums.jobs.ConductType.${ConductType[item.status]}`) }}
                                            </UBadge>
                                        </template>

                                        <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>

                            <UFormField class="flex-initial" name="id" :label="$t('common.id')">
                                <UInput v-model="query.id" type="text" name="id" :placeholder="$t('common.id')" />
                            </UFormField>

                            <UFormField
                                class="flex flex-initial flex-col"
                                name="showExpired"
                                :label="$t('components.jobs.conduct.List.show_expired')"
                                :ui="{ container: 'flex-1 flex' }"
                            >
                                <div class="flex flex-1 items-center">
                                    <USwitch v-model="query.showExpired" />
                                </div>
                            </UFormField>

                            <UFormField
                                v-if="can('jobs.ConductService/CreateConductEntry').value"
                                class="flex-initial"
                                :label="$t('common.create')"
                            >
                                <UButton
                                    trailing-icon="i-mdi-plus"
                                    color="neutral"
                                    truncate
                                    @click="
                                        conductCreateOrUpdateModal.open({
                                            onCreated: ($event) => data?.entries.unshift($event),
                                            onUpdated: ($event) => updateEntryInPlace($event),
                                        })
                                    "
                                >
                                    {{ $t('common.create') }}
                                </UButton>
                            </UFormField>
                        </div>
                    </UForm>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('common.conduct_register')])"
                :error="error"
                :retry="refresh"
            />

            <UTable
                v-else
                v-model:sorting="query.sorting.columns"
                class="flex-1"
                :loading="isRequestPending(status)"
                :columns="columns"
                :data="data?.entries"
                :empty="$t('common.not_found', [$t('common.entry', 2)])"
                :pagination-options="{ manualPagination: true }"
                :sorting-options="{ manualSorting: true }"
                sticky
            />
        </template>

        <template #footer>
            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
        </template>
    </UDashboardPanel>
</template>
