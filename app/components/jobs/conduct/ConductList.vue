<script lang="ts" setup>
import { z } from 'zod';
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import ConductCreateOrUpdateModal from '~/components/jobs/conduct/ConductCreateOrUpdateModal.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/stores/completor';
import { getJobsConductClient } from '~~/gen/ts/clients';
import { type ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import type { ListConductEntriesResponse } from '~~/gen/ts/services/jobs/conduct';
import ColleagueName from '../colleagues/ColleagueName.vue';
import ConductViewSlideover from './ConductViewSlideover.vue';
import { conductTypesToBadgeColor, conductTypesToBGColor } from './helpers';

const props = defineProps<{
    userId?: number;
    hideUserSearch?: boolean;
}>();

const { t } = useI18n();

const { can } = useAuth();

const modal = useOverlay();

const slideover = useOverlay();

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
        .custom<SortByColumn>()
        .array()
        .max(3)
        .default([
            {
                id: 'id',
                desc: true,
            },
        ]),
    page: pageNumberSchema,
});

const query = useSearchForm('jobs_conduct', schema);

const { data, status, refresh, error } = useLazyAsyncData(
    `jobs-conduct-${query.sorting.column}:${query.sorting.direction}-${query.page}-${query.types.join(',')}-${query.showExpired}-${query.id}`,
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
            sort: { columns: query.sorting },
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

const usersLoading = ref(false);

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

const columns = [
    {
        accessorKey: 'id',
        label: t('common.id'),
        sortable: true,
    },
    {
        accessorKey: 'createdAt',
        label: t('common.created_at'),
    },
    {
        accessorKey: 'expiresAt',
        label: t('common.expires_at'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    {
        accessorKey: 'type',
        label: t('common.type'),
        sortable: true,
    },
    {
        accessorKey: 'message',
        label: t('common.message'),
    },
    {
        accessorKey: 'target',
        label: t('common.target'),
    },
    {
        accessorKey: 'creator',
        label: t('common.creator'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    {
        accessorKey: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
];
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm class="w-full" :schema="schema" :state="query" @submit="refresh()">
                <div class="flex flex-row gap-2">
                    <UFormField v-if="hideUserSearch !== true" class="flex-1" name="user" :label="$t('common.search')">
                        <ClientOnly>
                            <USelectMenu
                                ref="input"
                                v-model="query.user"
                                :searchable="
                                    async (q: string) => {
                                        usersLoading = true;
                                        const colleagues = await completorStore.listColleagues({
                                            search: q,
                                            labelIds: [],
                                            userIds: query.user ? [query.user] : [],
                                        });
                                        usersLoading = false;
                                        return colleagues;
                                    }
                                "
                                searchable-lazy
                                :searchable-placeholder="$t('common.search_field')"
                                :search-attributes="['firstname', 'lastname']"
                                block
                                :placeholder="$t('common.colleague')"
                                trailing
                                leading-icon="i-mdi-search"
                                value-key="userId"
                                @keydown.esc="$event.target.blur()"
                            >
                                <template #item-label="{ item }">
                                    <span v-if="item" class="truncate">
                                        {{ userToLabel(item) }}
                                    </span>
                                </template>

                                <template #option="{ option: colleague }">
                                    <ColleagueName class="truncate" :colleague="colleague" birthday />
                                </template>

                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
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
                                :searchable-placeholder="$t('common.search_field')"
                                @keydown.esc="$event.target.blur()"
                            >
                                <template #item-label>
                                    {{ $t('common.selected', query.types.length) }}
                                </template>

                                <template #option="{ option }">
                                    <span class="truncate" :class="conductTypesToBGColor(option.status)">
                                        {{ $t(`enums.jobs.ConductType.${ConductType[option.status]}`) }}
                                    </span>
                                </template>

                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
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
                            <USwitch v-model="query.showExpired">
                                <span class="sr-only">
                                    {{ $t('components.jobs.conduct.List.show_expired') }}
                                </span>
                            </USwitch>
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
                                modal.open(ConductCreateOrUpdateModal, {
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

    <DataErrorBlock
        v-if="error"
        :title="$t('common.unable_to_load', [$t('common.conduct_register')])"
        :error="error"
        :retry="refresh"
    />
    <UTable
        v-else
        v-model:sorting="query.sorting"
        class="flex-1"
        :loading="isRequestPending(status)"
        :columns="columns"
        :data="data?.entries"
        :empty-state="{ icon: 'i-mdi-list-status', label: $t('common.not_found', [$t('common.entry', 2)]) }"
        sort-mode="manual"
    >
        <template #createdAt-cell="{ row: conduct }">
            <GenericTime :value="conduct.createdAt" />
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.expires_at') }}</dt>
                <dd class="mt-1 truncate">
                    <GenericTime
                        v-if="conduct.original.expiresAt?.value"
                        class="font-semibold"
                        :value="conduct.original.expiresAt.value"
                    />
                    <span v-else>
                        {{ $t('components.jobs.conduct.List.no_expiration') }}
                    </span>
                </dd>
            </dl>
        </template>

        <template #expiresAt-cell="{ row: conduct }">
            <GenericTime
                v-if="conduct.original.expiresAt?.value"
                class="font-semibold"
                type="date"
                :value="conduct.original.expiresAt.value"
            />
            <span v-else>
                {{ $t('components.jobs.conduct.List.no_expiration') }}
            </span>
        </template>

        <template #type-cell="{ row: conduct }">
            <UBadge :color="conductTypesToBadgeColor(conduct.original.type)">
                {{ $t(`enums.jobs.ConductType.${ConductType[conduct.original.type ?? 0]}`) }}
            </UBadge>
        </template>

        <template #message-cell="{ row: conduct }">
            <p class="line-clamp-2 w-full max-w-sm break-all whitespace-normal hover:line-clamp-6">
                {{ conduct.original.message }}
            </p>
        </template>

        <template #target-cell="{ row: conduct }">
            <ColleagueInfoPopover :user="conduct.original.targetUser" />
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.creator') }}</dt>
                <dd class="mt-1 truncate">
                    <ColleagueInfoPopover :user="conduct.original.creator.value" />
                </dd>
            </dl>
        </template>

        <template #creator-cell="{ row: conduct }">
            <ColleagueInfoPopover :user="conduct.original.creator.value" :hide-props="true" />
        </template>

        <template #actions-cell="{ row: conduct }">
            <div :key="conduct.original.id">
                <UTooltip :text="$t('common.show')">
                    <UButton
                        variant="link"
                        icon="i-mdi-eye"
                        @click="
                            slideover.open(ConductViewSlideover, {
                                entry: conduct,
                                onRefresh: async () => refresh(),
                            })
                        "
                    />
                </UTooltip>

                <UTooltip v-if="can('jobs.ConductService/UpdateConductEntry').value" :text="$t('common.update')">
                    <UButton
                        variant="link"
                        icon="i-mdi-pencil"
                        @click="
                            modal.open(ConductCreateOrUpdateModal, {
                                entry: conduct,
                                userId: userId,
                                onCreated: ($event) => data?.entries.unshift($event),
                                onUpdated: ($event) => updateEntryInPlace($event),
                            })
                        "
                    />
                </UTooltip>

                <UTooltip v-if="can('jobs.ConductService/DeleteConductEntry').value" :text="$t('common.delete')">
                    <UButton
                        variant="link"
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteConductEntry(conduct.original.id),
                            })
                        "
                    />
                </UTooltip>
            </div>
        </template>
    </UTable>

    <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
</template>
