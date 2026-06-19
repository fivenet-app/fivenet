<script lang="ts" setup>
import { UButton, UTooltip } from '#components';
import type { DropdownMenuItem, Form, TableColumn } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { getJobsGroupsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { type Group, GroupState, GroupType } from '~~/gen/ts/resources/jobs/groups/group';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ListGroupsResponse } from '~~/gen/ts/services/jobs/groups';
import GroupEditorModal from './GroupEditorModal.vue';
import {
    groupMembershipModeColor,
    groupMembershipModeLabel,
    groupStateColor,
    groupStateFilterItems,
    groupStateLabel,
    groupTypeColor,
    groupTypeFilterItems,
    groupTypeLabel,
} from './helpers';

const { t } = useI18n();

const overlay = useOverlay();

const notifications = useNotificationsStore();

const { can } = useAuth();

const canCreateGroup = can('jobs.GroupsService/CreateGroup');
const canArchiveGroup = can('jobs.GroupsService/ArchiveGroup');

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
                        id: 'sort_order',
                        desc: false,
                    },
                ]),
        })
        .default({ columns: [{ id: 'sort_order', desc: false }] }),
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

const groupEditorModal = overlay.create(GroupEditorModal);
const confirmModal = overlay.create(ConfirmModal);

const breadcrumbs = [
    { label: 'Jobs', to: '/jobs' },
    { label: 'Groups', to: '/jobs/groups' },
];

const sortFields = [
    { label: 'Sort order', value: 'sort_order' },
    { label: 'Name', value: 'name' },
    { label: 'Status', value: 'state' },
    { label: 'Updated', value: 'updated_at' },
    { label: 'Created', value: 'created_at' },
    { label: 'Members', value: 'members_count' },
    { label: 'Leaders', value: 'leaders_count' },
    { label: 'Rules', value: 'rules_count' },
    { label: 'Exclusions', value: 'exclusions_count' },
];

const columns: TableColumn<Group>[] = [
    {
        accessorKey: 'name',
        header: t('common.name'),
    },
    {
        accessorKey: 'type',
        header: 'Type',
    },
    {
        accessorKey: 'membershipMode',
        header: 'Membership mode',
    },
    {
        accessorKey: 'counts',
        header: 'Counts',
    },
    {
        accessorKey: 'updatedAt',
        header: t('common.updated_at'),
    },
    {
        id: 'actions',
        header: '',
    },
];

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
            includeInactive: values.status === 'all',
            includeArchived: values.status === 'all',
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

const stats = computed(() => {
    const groups = visibleGroups.value;

    return [
        {
            label: 'Visible groups',
            value: groups.length,
            icon: 'i-mdi-account-group',
        },
        {
            label: 'Members',
            value: groups.reduce((total, group) => total + group.membersCount, 0),
            icon: 'i-mdi-account-multiple',
        },
        {
            label: 'Leaders',
            value: groups.reduce((total, group) => total + group.leadersCount, 0),
            icon: 'i-mdi-account-star',
        },
        {
            label: 'Exclusions',
            value: groups.reduce((total, group) => total + group.exclusionsCount, 0),
            icon: 'i-mdi-account-cancel',
        },
    ];
});

function openCreateGroup(): void {
    if (!canCreateGroup.value) return;

    groupEditorModal.open({
        onCreated: async () => refresh(),
    });
}

function openEditGroup(group: Group): void {
    if (!canCreateGroup.value || group.state === GroupState.ARCHIVED) return;

    groupEditorModal.open({
        group,
        onUpdated: async () => refresh(),
    });
}

async function archiveGroup(group: Group): Promise<void> {
    try {
        await jobsGroupsClient.archiveGroup({ id: group.id });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function restoreGroup(group: Group): Promise<void> {
    try {
        await jobsGroupsClient.restoreGroup({ id: group.id });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function getDropdownActions(group: Group): DropdownMenuItem[][] {
    const items: DropdownMenuItem[][] = [];

    if (group.state !== GroupState.ARCHIVED) {
        if (canCreateGroup.value) {
            items.push([
                {
                    label: 'Edit group',
                    icon: 'i-mdi-pencil',
                    onSelect: () => openEditGroup(group),
                },
            ]);
        }

        if (canArchiveGroup.value) {
            items.push([
                {
                    label: 'Archive group',
                    icon: 'i-mdi-archive',
                    color: 'warning',
                    onSelect: () =>
                        confirmModal.open({
                            title: 'Archive group',
                            description: `Archive "${group.name}"?`,
                            confirm: async () => archiveGroup(group),
                        }),
                },
            ]);
        }
    } else if (canArchiveGroup.value) {
        items.push([
            {
                label: 'Restore group',
                icon: 'i-mdi-restore',
                color: 'primary',
                onSelect: () =>
                    confirmModal.open({
                        title: 'Restore group',
                        description: `Restore "${group.name}"?`,
                        confirm: async () => restoreGroup(group),
                    }),
            },
        ]);
    }

    return items;
}
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardToolbar>
                <template #left>
                    <UBreadcrumb :items="breadcrumbs" />
                </template>

                <template #right>
                    <UTooltip v-if="canCreateGroup" :text="$t('common.create')">
                        <UButton label="New group" icon="i-mdi-plus" @click="openCreateGroup" />
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
                            <UFormField class="w-full" name="search" label="Search">
                                <UInput
                                    v-model="query.search"
                                    class="w-full"
                                    icon="i-mdi-magnify"
                                    placeholder="Search groups..."
                                />
                            </UFormField>

                            <UFormField class="w-full" name="status" label="Status">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="query.status"
                                        class="w-full"
                                        :items="groupStateFilterItems"
                                        value-key="value"
                                    />
                                </ClientOnly>
                            </UFormField>

                            <UFormField class="w-full" name="kind" label="Type">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="query.kind"
                                        class="w-full"
                                        :items="groupTypeFilterItems"
                                        value-key="value"
                                    />
                                </ClientOnly>
                            </UFormField>

                            <UFormField class="w-full" name="sorting" label="Sort">
                                <SortButton v-model="query.sorting" class="w-full" :fields="sortFields" />
                            </UFormField>
                        </div>

                        <UFormField label="&nbsp;">
                            <UButton type="submit" label="Apply" icon="i-mdi-filter" variant="soft" />
                        </UFormField>
                    </UForm>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <div class="flex min-h-0 flex-1 flex-col gap-4 overflow-auto p-4 sm:p-6">
                <div>
                    <h1 class="text-2xl font-semibold tracking-tight">Groups</h1>
                    <p class="mt-1 text-sm text-muted">
                        Manage job-local groups, their lifecycle state and the metadata used by other job views.
                    </p>
                </div>

                <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
                    <UCard v-for="stat in stats" :key="stat.label" variant="subtle">
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
                            <h2 class="text-base font-semibold">Group index</h2>
                            <p class="text-sm text-muted">
                                Search groups, inspect their cached membership counts and update metadata.
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
                            <div class="flex items-start gap-3">
                                <UAvatar
                                    :alt="row.original.name"
                                    :text="(row.original.shortName ?? row.original.name).slice(0, 2)"
                                    size="md"
                                />

                                <div class="min-w-0">
                                    <div class="flex flex-wrap items-center gap-2">
                                        <span class="font-medium text-highlighted">
                                            {{ row.original.name }}
                                        </span>

                                        <UBadge :color="groupStateColor(row.original.state)" variant="subtle" size="sm">
                                            {{ groupStateLabel(row.original.state) }}
                                        </UBadge>

                                        <UBadge v-if="row.original.shortName" color="neutral" variant="soft" size="sm">
                                            {{ row.original.shortName }}
                                        </UBadge>
                                    </div>

                                    <p class="mt-1 line-clamp-2 text-sm text-muted">
                                        {{ row.original.description || 'No description' }}
                                    </p>
                                </div>
                            </div>
                        </template>

                        <template #type-cell="{ row }">
                            <UBadge :color="groupTypeColor(row.original.type)" variant="subtle">
                                {{ groupTypeLabel(row.original.type) }}
                            </UBadge>
                        </template>

                        <template #membershipMode-cell="{ row }">
                            <UBadge :color="groupMembershipModeColor(row.original.membershipMode)" variant="subtle">
                                {{ groupMembershipModeLabel(row.original.membershipMode) }}
                            </UBadge>
                        </template>

                        <template #counts-cell="{ row }">
                            <div class="grid grid-cols-2 gap-x-3 gap-y-1 text-sm">
                                <span class="text-muted">Members</span>
                                <span class="text-right font-medium tabular-nums">{{ row.original.membersCount }}</span>

                                <span class="text-muted">Leaders</span>
                                <span class="text-right font-medium tabular-nums">{{ row.original.leadersCount }}</span>

                                <span class="text-muted">Rules</span>
                                <span class="text-right font-medium tabular-nums">{{ row.original.rulesCount }}</span>

                                <span class="text-muted">Exclusions</span>
                                <span class="text-right font-medium tabular-nums">{{ row.original.exclusionsCount }}</span>
                            </div>
                        </template>

                        <template #updatedAt-cell="{ row }">
                            <GenericTime :value="row.original.updatedAt ?? row.original.createdAt" />
                        </template>

                        <template #actions-cell="{ row }">
                            <UDropdownMenu
                                v-if="getDropdownActions(row.original).length"
                                :items="getDropdownActions(row.original)"
                            >
                                <UButton
                                    icon="i-mdi-dots-vertical"
                                    color="neutral"
                                    variant="ghost"
                                    aria-label="Group actions"
                                />
                            </UDropdownMenu>
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
