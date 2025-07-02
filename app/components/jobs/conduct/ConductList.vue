<script lang="ts" setup>
import { z } from 'zod';
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import ConductCreateOrUpdateModal from '~/components/jobs/conduct/ConductCreateOrUpdateModal.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/stores/completor';
import { type ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import type { ListConductEntriesResponse } from '~~/gen/ts/services/jobs/conduct';
import ColleagueName from '../colleagues/ColleagueName.vue';
import ConductViewSlideover from './ConductViewSlideover.vue';
import { conductTypesToBadgeColor, conductTypesToBGColor } from './helpers';

const props = defineProps<{
    userId?: number;
    hideUserSearch?: boolean;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { can } = useAuth();

const modal = useModal();

const slideover = useSlideover();

const completorStore = useCompletorStore();

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
    sort: z.custom<TableSortable>().default({
        column: 'id',
        direction: 'desc',
    }),
    page: pageNumberSchema,
});

const query = useSearchForm('jobs_conduct', schema);

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `jobs-conduct-${query.sort.column}:${query.sort.direction}-${query.page}-${query.types.join(',')}-${query.showExpired}-${query.id}`,
    () => listConductEntries(),
    {
        transform: (input) => ({ ...input, entries: wrapRows(input?.entries, columns) }),
    },
);

async function listConductEntries(): Promise<ListConductEntriesResponse> {
    const entryIds: number[] = [];
    if (query.id) {
        entryIds.push(typeof query.id === 'string' ? parseInt(query.id, 10) : query.id);
    }

    const userIds = props.userId ? [props.userId] : query.user ? [query.user] : [];
    try {
        const call = $grpc.jobs.conduct.listConductEntries({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sort,
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
        const call = $grpc.jobs.conduct.deleteConductEntry({ id });
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
        key: 'id',
        label: t('common.id'),
        sortable: true,
    },
    {
        key: 'createdAt',
        label: t('common.created_at'),
    },
    {
        key: 'expiresAt',
        label: t('common.expires_at'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    {
        key: 'type',
        label: t('common.type'),
        sortable: true,
    },
    {
        key: 'message',
        label: t('common.message'),
    },
    {
        key: 'target',
        label: t('common.target'),
    },
    {
        key: 'creator',
        label: t('common.creator'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    {
        key: 'actions',
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
                    <UFormGroup v-if="hideUserSearch !== true" class="flex-1" name="user" :label="$t('common.search')">
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
                                value-attribute="userId"
                                @keydown.esc="$event.target.blur()"
                            >
                                <template #label="{ selected }">
                                    <span v-if="selected" class="truncate">
                                        {{ userToLabel(selected) }}
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
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="types" :label="$t('common.type')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="query.types"
                                multiple
                                nullable
                                :options="availableTypes"
                                value-attribute="status"
                                :placeholder="$t('common.na')"
                                :searchable-placeholder="$t('common.search_field')"
                                @keydown.esc="$event.target.blur()"
                            >
                                <template #label>
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
                    </UFormGroup>

                    <UFormGroup class="flex-initial" name="id" :label="$t('common.id')">
                        <UInput v-model="query.id" type="text" name="id" :placeholder="$t('common.id')" />
                    </UFormGroup>

                    <UFormGroup
                        class="flex flex-initial flex-col"
                        name="showExpired"
                        :label="$t('components.jobs.conduct.List.show_expired')"
                        :ui="{ container: 'flex-1 flex' }"
                    >
                        <div class="flex flex-1 items-center">
                            <UToggle v-model="query.showExpired">
                                <span class="sr-only">
                                    {{ $t('components.jobs.conduct.List.show_expired') }}
                                </span>
                            </UToggle>
                        </div>
                    </UFormGroup>

                    <UFormGroup
                        v-if="can('jobs.ConductService/CreateConductEntry').value"
                        class="flex-initial"
                        :label="$t('common.create')"
                    >
                        <UButton
                            trailing-icon="i-mdi-plus"
                            color="gray"
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
                    </UFormGroup>
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
        v-model:sort="query.sort"
        class="flex-1"
        :loading="loading"
        :columns="columns"
        :rows="data?.entries"
        :empty-state="{ icon: 'i-mdi-list-status', label: $t('common.not_found', [$t('common.entry', 2)]) }"
        sort-mode="manual"
    >
        <template #createdAt-data="{ row: conduct }">
            <GenericTime :value="conduct.createdAt" />
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.expires_at') }}</dt>
                <dd class="mt-1 truncate">
                    <GenericTime v-if="conduct.expiresAt?.value" class="font-semibold" :value="conduct.expiresAt.value" />
                    <span v-else>
                        {{ $t('components.jobs.conduct.List.no_expiration') }}
                    </span>
                </dd>
            </dl>
        </template>

        <template #expiresAt-data="{ row: conduct }">
            <GenericTime v-if="conduct.expiresAt?.value" class="font-semibold" type="date" :value="conduct.expiresAt.value" />
            <span v-else>
                {{ $t('components.jobs.conduct.List.no_expiration') }}
            </span>
        </template>

        <template #type-data="{ row: conduct }">
            <UBadge :color="conductTypesToBadgeColor(conduct.type)">
                {{ $t(`enums.jobs.ConductType.${ConductType[conduct.type ?? 0]}`) }}
            </UBadge>
        </template>

        <template #message-data="{ row: conduct }">
            <p class="line-clamp-2 w-full max-w-sm whitespace-normal break-all hover:line-clamp-6">
                {{ conduct.message }}
            </p>
        </template>

        <template #target-data="{ row: conduct }">
            <ColleagueInfoPopover :user="conduct.targetUser" />
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.creator') }}</dt>
                <dd class="mt-1 truncate">
                    <ColleagueInfoPopover :user="conduct.creator.value" />
                </dd>
            </dl>
        </template>

        <template #creator-data="{ row: conduct }">
            <ColleagueInfoPopover :user="conduct.creator.value" :hide-props="true" />
        </template>

        <template #actions-data="{ row: conduct }">
            <div :key="conduct.id">
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
                                confirm: async () => deleteConductEntry(conduct.id),
                            })
                        "
                    />
                </UTooltip>
            </div>
        </template>
    </UTable>

    <Pagination v-model="query.page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
