<script lang="ts" setup>
import { UBadge, UButton, UTooltip } from '#components';
import type { Form, TableColumn } from '@nuxt/ui';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { getJobsGroupsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { type Group, GroupState, GroupType } from '~~/gen/ts/resources/jobs/groups/group';
import type { ListGroupsResponse } from '~~/gen/ts/services/jobs/groups';
import DetailsSlideover from './DetailsSlideover.vue';
import EditorModal from './EditorModal.vue';
import {
    groupMembershipModeColor,
    groupMembershipModeIcon,
    groupMembershipModeLabelKey,
    groupStateColor,
    groupStateFilterItems as groupStateFilterItemKeys,
    groupStateLabelKey,
    groupTypeColor,
    groupTypeFilterItems as groupTypeFilterItemKeys,
    groupTypeIcon,
    groupTypeLabelKey,
} from './helpers';

const { t } = useI18n();

const overlay = useOverlay();

const { can } = useAuth();

const canCreateGroup = can('jobs.GroupsService/CreateGroup');

const jobsGroupsClient = await getJobsGroupsClient();

const schema = z.object({
    search: z.coerce.string().max(100).default(''),
    status: z.enum(['active', 'inactive', 'archived', 'all']).default('active'),
    kind: z.enum(['all', 'manual', 'smart', 'mixed']).default('all'),
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'sort_rank',
                        desc: false,
                    },
                ]),
        })
        .default({ columns: [{ id: 'sort_rank', desc: false }] }),
    page: pageNumberSchema,
});

type Schema = z.output<typeof schema>;

const query = useSearchForm('jobs_groups', schema);
const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { validatedQuery, commitValidatedQuery } = useFormSearchValidation<typeof schema>(query, formRef);

const listKey = computed(
    () =>
        `jobs-groups-${validatedQuery.value.page}-${validatedQuery.value.status}-${validatedQuery.value.search}-${JSON.stringify(validatedQuery.value.sorting)}`,
);

const { data, status: requestStatus, refresh, error } = useLazyAsyncData(listKey, () => listGroups(validatedQuery.value));

const editorModal = overlay.create(EditorModal);
const detailsSlideover = overlay.create(DetailsSlideover);

const groupStateFilterItems = computed(() => groupStateFilterItemKeys.map((item) => ({ ...item, label: t(item.labelKey) })));
const groupTypeFilterItems = computed(() => groupTypeFilterItemKeys.map((item) => ({ ...item, label: t(item.labelKey) })));

const sortFields = computed(() => [
    { label: t('components.jobs.groups.sort_rank'), value: 'sort_rank' },
    { label: t('common.name'), value: 'name' },
    { label: t('common.status'), value: 'state' },
    { label: t('common.updated'), value: 'updated_at' },
    { label: t('common.created'), value: 'created_at' },
    { label: t('common.members', 2), value: 'members_count' },
    { label: t('components.jobs.groups.leaders'), value: 'leaders_count' },
    { label: t('components.jobs.groups.rules'), value: 'rules_count' },
    { label: t('components.jobs.groups.exclusions'), value: 'exclusions_count' },
]);

const columns = computed<TableColumn<Group>[]>(() => [
    {
        accessorKey: 'name',
        header: t('common.name'),
    },
    {
        accessorKey: 'type',
        header: t('common.type'),
        cell: ({ row }) =>
            h(UBadge, {
                color: groupTypeColor(row.original.type),
                icon: groupTypeIcon(row.original.type),
                label: t(groupTypeLabelKey(row.original.type)),
                variant: 'subtle',
            }),
    },
    {
        accessorKey: 'membershipMode',
        header: t('components.jobs.groups.membership_mode'),
        cell: ({ row }) =>
            h(UBadge, {
                color: groupMembershipModeColor(row.original.membershipMode),
                icon: groupMembershipModeIcon(row.original.membershipMode),
                label: t(groupMembershipModeLabelKey(row.original.membershipMode)),
                variant: 'subtle',
            }),
    },
    {
        accessorKey: 'counts',
        header: t('common.count'),
    },
    {
        accessorKey: 'updatedAt',
        header: t('common.updated_at'),
        cell: ({ row }) => h(GenericTime, { value: row.original.updatedAt ?? row.original.createdAt }),
    },
    {
        id: 'actions',
        header: '',
        cell: ({ row }) =>
            h(
                UTooltip,
                {
                    text: $t('components.jobs.groups.actions.details'),
                },
                [h(UButton, { icon: 'i-mdi-eye', variant: 'link', onClick: () => openGroupDetails(row.original) })],
            ),
    },
]);

async function listGroups(values: Schema): Promise<ListGroupsResponse> {
    try {
        const statusStates =
            values.status === 'active'
                ? [GroupState.ACTIVE]
                : values.status === 'inactive'
                  ? [GroupState.INACTIVE]
                  : values.status === 'archived'
                    ? [GroupState.ARCHIVED]
                    : [];

        const { response } = await jobsGroupsClient.listGroups({
            pagination: {
                offset: calculateOffset(values.page, data.value?.pagination),
            },
            sort: values.sorting,
            search: values.search.trim() ? values.search.trim() : undefined,
            states: statusStates,
            includeCounts: true,
            includeInactive: values.status === 'all' || values.status === 'inactive',
            includeArchived: values.status === 'all' || values.status === 'archived',
            groupIds: [],
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function groupMatchesKind(group: Group): boolean {
    if (query.kind === 'all') return true;

    if (query.kind === 'manual') return group.type === GroupType.MANUAL;
    if (query.kind === 'smart') return group.type === GroupType.SMART;
    if (query.kind === 'mixed') return group.type === GroupType.MIXED;

    return true;
}

const visibleGroups = computed(() => (data.value?.groups ?? []).filter((group) => groupMatchesKind(group)));

const stats = computed(() => [
    {
        label: t('components.jobs.groups.visible_groups'),
        value: visibleGroups.value.length,
        icon: 'i-mdi-account-group',
    },
    {
        label: t('common.members', 2),
        value: visibleGroups.value.reduce((total, group) => total + group.membersCount, 0),
        icon: 'i-mdi-account-multiple',
    },
    {
        label: t('components.jobs.groups.leaders'),
        value: visibleGroups.value.reduce((total, group) => total + group.leadersCount, 0),
        icon: 'i-mdi-account-star',
    },
    {
        label: t('components.jobs.groups.exclusions'),
        value: visibleGroups.value.reduce((total, group) => total + group.exclusionsCount, 0),
        icon: 'i-mdi-account-cancel',
    },
]);

function openCreateGroup(): void {
    if (!canCreateGroup.value) return;

    editorModal.open({ onCreated: async () => refresh() });
}

function openGroupDetails(group: Group): void {
    detailsSlideover.open({
        group,
        onChanged: async () => refresh(),
    });
}
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardToolbar>
                <template #right>
                    <UTooltip v-if="canCreateGroup" :text="$t('common.create')">
                        <UButton :label="$t('components.jobs.groups.actions.new')" icon="i-mdi-plus" @click="openCreateGroup" />
                    </UTooltip>
                </template>
            </UDashboardToolbar>

            <UDashboardToolbar>
                <template #default>
                    <UForm
                        ref="formRef"
                        class="my-2 flex w-full flex-col gap-2 lg:flex-row lg:items-end lg:justify-between"
                        :schema="schema"
                        :state="query"
                        @submit="commitValidatedQuery"
                    >
                        <div class="grid flex-1 gap-2 sm:grid-cols-2 xl:grid-cols-4">
                            <UFormField class="w-full" name="search" :label="$t('common.search')">
                                <UInput
                                    v-model="query.search"
                                    class="w-full"
                                    icon="i-mdi-magnify"
                                    :placeholder="$t('components.jobs.groups.search_placeholder')"
                                />
                            </UFormField>

                            <UFormField class="w-full" name="status" :label="$t('common.status')">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="query.status"
                                        class="w-full"
                                        :items="groupStateFilterItems"
                                        value-key="value"
                                    />
                                </ClientOnly>
                            </UFormField>

                            <UFormField class="w-full" name="kind" :label="$t('common.type')">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="query.kind"
                                        class="w-full"
                                        :items="groupTypeFilterItems"
                                        value-key="value"
                                    />
                                </ClientOnly>
                            </UFormField>

                            <UFormField class="w-full" name="sorting" :label="$t('components.jobs.groups.sort')">
                                <SortButton v-model="query.sorting" class="w-full" :fields="sortFields" />
                            </UFormField>
                        </div>

                        <UFormField label="&nbsp;">
                            <UButton type="submit" :label="$t('common.apply')" icon="i-mdi-filter" variant="soft" />
                        </UFormField>
                    </UForm>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <div class="flex min-h-0 flex-1 flex-col gap-4 overflow-auto p-4 sm:p-6">
                <div>
                    <h1 class="text-2xl font-semibold tracking-tight">{{ $t('pages.jobs.groups.title') }}</h1>
                    <p class="mt-1 text-sm text-muted">
                        {{ $t('components.jobs.groups.subtitle') }}
                    </p>
                </div>

                <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
                    <UCard
                        v-for="stat in stats"
                        :key="stat.label"
                        orientation="horizontal"
                        variant="subtle"
                        :ui="{ body: 'p-2 sm:p-4' }"
                    >
                        <div class="flex items-center gap-3">
                            <div class="flex size-10 items-center justify-center rounded-lg bg-elevated">
                                <UIcon class="size-5 text-muted" :name="stat.icon" />
                            </div>

                            <div>
                                <p class="text-sm text-muted">
                                    {{ stat.label }}
                                </p>
                                <p class="text-2xl font-semibold tabular-nums">
                                    {{ stat.value }}
                                </p>
                            </div>
                        </div>
                    </UCard>
                </div>

                <UCard :ui="{ body: 'p-0 sm:p-0' }">
                    <template #header>
                        <div>
                            <h2 class="text-base font-semibold">{{ $t('components.jobs.groups.index_title') }}</h2>
                            <p class="text-sm text-muted">
                                {{ $t('components.jobs.groups.index_description') }}
                            </p>
                        </div>
                    </template>

                    <DataPendingBlock
                        v-if="isRequestPending(requestStatus)"
                        :message="$t('common.loading', [$t('common.group', 2)])"
                    />
                    <DataErrorBlock
                        v-else-if="error"
                        :title="$t('common.unable_to_load', [$t('common.group', 2)])"
                        :error="error"
                        :retry="refresh"
                    />
                    <DataNoDataBlock v-else-if="visibleGroups.length === 0" :type="$t('common.group', 2)" :retry="refresh" />
                    <UTable v-else class="min-h-96" :data="visibleGroups" :columns="columns">
                        <template #name-cell="{ row }">
                            <div class="flex items-center gap-3">
                                <UAvatar
                                    :alt="row.original.name"
                                    :text="(row.original.shortName ?? row.original.name).slice(0, 2)"
                                    size="md"
                                />

                                <div class="w-full min-w-0">
                                    <div class="flex w-full flex-wrap items-center gap-2">
                                        <span class="font-medium text-highlighted">
                                            {{ row.original.name }}
                                        </span>

                                        <UBadge
                                            v-if="row.original.shortName && row.original.shortName !== row.original.name"
                                            color="neutral"
                                            :label="row.original.shortName"
                                            variant="soft"
                                            size="sm"
                                        />

                                        <div class="flex-1" />

                                        <UBadge :color="groupStateColor(row.original.state)" variant="subtle" size="sm">
                                            {{ $t(groupStateLabelKey(row.original.state)) }}
                                        </UBadge>
                                    </div>

                                    <p class="mt-1 line-clamp-2 text-sm text-muted">
                                        {{ row.original.description || $t('components.jobs.groups.no_description') }}
                                    </p>
                                </div>
                            </div>
                        </template>

                        <template #counts-cell="{ row }">
                            <span class="text-right text-sm font-medium tabular-nums">
                                {{ $t('common.member', row.original.membersCount) }}
                            </span>

                            <div class="flex flex-row flex-wrap gap-x-1 gap-y-1 text-xs">
                                <span class="text-muted">{{ $t('components.jobs.groups.leaders') }}</span>
                                <span class="text-right font-medium tabular-nums">{{ row.original.leadersCount }}</span>
                                <span>·</span>

                                <span class="text-muted">{{ $t('components.jobs.groups.rules') }}</span>
                                <span class="text-right font-medium tabular-nums">{{ row.original.rulesCount }}</span>
                                <span>·</span>

                                <span class="text-muted">{{ $t('components.jobs.groups.exclusions') }}</span>
                                <span class="text-right font-medium tabular-nums">{{ row.original.exclusionsCount }}</span>
                            </div>
                        </template>
                    </UTable>
                </UCard>
            </div>
        </template>

        <template #footer>
            <Pagination v-model="query.page" :pagination="data?.pagination" :status="requestStatus" :refresh="refresh" />
        </template>
    </UDashboardPanel>
</template>
