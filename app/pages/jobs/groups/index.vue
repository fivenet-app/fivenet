<script lang="ts" setup>
import type { DropdownMenuItem, TableColumn } from '@nuxt/ui';
import { z } from 'zod';

useHead({
    title: 'pages.jobs.groups.title',
});

definePageMeta({
    title: 'pages.jobs.groups.title',
    requiresAuth: true,
    permission: 'TODOService/TODOMethod',
});

type GroupStatus = 'active' | 'archived';
type GroupKind = 'manual' | 'smart' | 'mixed';

type JobGroup = {
    id: number;
    name: string;
    description: string;
    status: GroupStatus;
    kind: GroupKind;
    members: number;
    manualMembers: number;
    derivedMembers: number;
    excludedMembers: number;
    leaders: { name: string; avatar?: string }[];
    rules: string[];
    updatedAt: string;
};

const breadcrumbs = [
    { label: 'Jobs', to: '/jobs' },
    { label: 'Groups', to: '/jobs/groups' },
];

const formRef = useTemplateRef('formRef');

const schema = z.object({
    search: z.string().optional(),
    status: z.enum(['active', 'archived', 'all']),
    kind: z.enum(['manual', 'smart', 'mixed', 'all']),
});

const query = reactive<z.infer<typeof schema>>({
    search: '',
    status: 'active',
    kind: 'all',
});

function commitValidatedQuery() {
    // Push to route query here if wanted.
}

const statusItems = [
    { label: 'Active', value: 'active' },
    { label: 'Archived', value: 'archived' },
    { label: 'All', value: 'all' },
];

const kindItems = [
    { label: 'All types', value: 'all' },
    { label: 'Manual', value: 'manual' },
    { label: 'Smart', value: 'smart' },
    { label: 'Mixed', value: 'mixed' },
];

const groups = ref<JobGroup[]>([
    {
        id: 1,
        name: 'K9 Unit',
        description: 'Certified handlers and assigned support staff.',
        status: 'active',
        kind: 'mixed',
        members: 12,
        manualMembers: 3,
        derivedMembers: 10,
        excludedMembers: 1,
        leaders: [{ name: 'M. Schneider' }, { name: 'L. Weber' }],
        rules: ['Qualification: K9 Handler', 'Rank: Sergeant+'],
        updatedAt: '2026-06-15',
    },
    {
        id: 2,
        name: 'Field Training Officers',
        description: 'Employees allowed to supervise probationary colleagues.',
        status: 'active',
        kind: 'smart',
        members: 8,
        manualMembers: 0,
        derivedMembers: 8,
        excludedMembers: 0,
        leaders: [{ name: 'A. Klein' }],
        rules: ['Qualification: FTO'],
        updatedAt: '2026-06-12',
    },
    {
        id: 3,
        name: 'Command Staff',
        description: 'Leadership group for internal coordination.',
        status: 'active',
        kind: 'manual',
        members: 5,
        manualMembers: 5,
        derivedMembers: 0,
        excludedMembers: 0,
        leaders: [{ name: 'Chief Bauer' }],
        rules: [],
        updatedAt: '2026-06-08',
    },
    {
        id: 4,
        name: 'Old Patrol Division A',
        description: 'Archived historic patrol structure.',
        status: 'archived',
        kind: 'manual',
        members: 0,
        manualMembers: 0,
        derivedMembers: 0,
        excludedMembers: 0,
        leaders: [],
        rules: [],
        updatedAt: '2026-04-21',
    },
]);

const filteredGroups = computed(() => {
    return groups.value.filter((group) => {
        const matchesSearch = [group.name, group.description, ...group.rules]
            .join(' ')
            .toLowerCase()
            .includes((query.search ?? '').toLowerCase());

        const matchesStatus = query.status === 'all' || group.status === query.status;
        const matchesKind = query.kind === 'all' || group.kind === query.kind;

        return matchesSearch && matchesStatus && matchesKind;
    });
});

const stats = computed(() => {
    const active = groups.value.filter((group) => group.status === 'active');

    return [
        {
            label: 'Active groups',
            value: active.length,
            icon: 'i-mdi-account-group',
        },
        {
            label: 'Total memberships',
            value: active.reduce((sum, group) => sum + group.members, 0),
            icon: 'i-mdi-account-check',
        },
        {
            label: 'Smart groups',
            value: active.filter((group) => group.kind !== 'manual').length,
            icon: 'i-mdi-auto-fix',
        },
        {
            label: 'Manual exclusions',
            value: active.reduce((sum, group) => sum + group.excludedMembers, 0),
            icon: 'i-mdi-account-cancel',
        },
    ];
});

const columns: TableColumn<JobGroup>[] = [
    {
        accessorKey: 'name',
        header: 'Group',
    },
    {
        accessorKey: 'kind',
        header: 'Type',
    },
    {
        accessorKey: 'members',
        header: 'Members',
        meta: {
            class: {
                th: 'text-right',
                td: 'text-right',
            },
        },
    },
    {
        accessorKey: 'leaders',
        header: 'Leaders',
    },
    {
        accessorKey: 'rules',
        header: 'Rules',
    },
    {
        accessorKey: 'updatedAt',
        header: 'Updated',
    },
    {
        id: 'action',
        header: '',
    },
];

function kindColor(kind: GroupKind) {
    return {
        manual: 'neutral',
        smart: 'primary',
        mixed: 'success',
    }[kind] as 'neutral' | 'primary' | 'success';
}

function statusColor(status: GroupStatus) {
    return status === 'active' ? 'success' : 'neutral';
}

function getDropdownActions(group: JobGroup): DropdownMenuItem[][] {
    return [
        [
            {
                label: 'Open group',
                icon: 'i-mdi-arrow-right',
                to: `/jobs/groups/${group.id}`,
            },
            {
                label: 'Edit settings',
                icon: 'i-mdi-cog',
                to: `/jobs/groups/${group.id}/settings`,
            },
        ],
        [
            {
                label: 'View members',
                icon: 'i-mdi-account-group',
                to: `/jobs/groups/${group.id}/members`,
            },
            {
                label: 'View activity',
                icon: 'i-mdi-history',
                to: `/jobs/groups/${group.id}/activity`,
            },
        ],
        [
            {
                label: group.status === 'active' ? 'Archive group' : 'Restore group',
                icon: group.status === 'active' ? 'i-mdi-archive' : 'i-mdi-restore',
                color: group.status === 'active' ? 'warning' : 'primary',
            },
        ],
    ];
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
                    <UButton label="New group" icon="i-mdi-plus" to="/jobs/groups/new" />
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
                        <div class="grid flex-1 gap-2 sm:grid-cols-3">
                            <UFormField name="search" label="Search">
                                <UInput v-model="query.search" icon="i-mdi-magnify" placeholder="Search groups..." />
                            </UFormField>

                            <UFormField name="status" label="Status">
                                <USelect v-model="query.status" :items="statusItems" />
                            </UFormField>

                            <UFormField name="kind" label="Type">
                                <USelect v-model="query.kind" :items="kindItems" />
                            </UFormField>
                        </div>

                        <UButton type="submit" label="Apply" icon="i-mdi-filter" variant="soft" />
                    </UForm>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <div class="flex min-h-0 flex-1 flex-col gap-4 overflow-auto p-4 sm:p-6">
                <div>
                    <h1 class="text-2xl font-semibold tracking-tight">Groups</h1>
                    <p class="mt-1 text-sm text-muted">
                        Manage job-local groups for colleague lists, timeclock filters and delegated group leadership.
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
                                Search groups, inspect member source counts and jump into members or rules.
                            </p>
                        </div>
                    </template>

                    <UTable class="min-h-96" :data="filteredGroups" :columns="columns">
                        <template #name-cell="{ row }">
                            <div class="flex items-start gap-3">
                                <UAvatar :alt="row.original.name" :text="row.original.name.slice(0, 2)" size="md" />

                                <div class="min-w-0">
                                    <div class="flex flex-wrap items-center gap-2">
                                        <ULink
                                            class="font-medium text-highlighted hover:underline"
                                            :to="`/jobs/groups/${row.original.id}`"
                                        >
                                            {{ row.original.name }}
                                        </ULink>

                                        <UBadge :color="statusColor(row.original.status)" variant="subtle" size="sm">
                                            {{ row.original.status }}
                                        </UBadge>
                                    </div>

                                    <p class="mt-1 line-clamp-2 text-sm text-muted">
                                        {{ row.original.description }}
                                    </p>
                                </div>
                            </div>
                        </template>

                        <template #kind-cell="{ row }">
                            <UBadge class="capitalize" :color="kindColor(row.original.kind)" variant="subtle">
                                {{ row.original.kind }}
                            </UBadge>
                        </template>

                        <template #members-cell="{ row }">
                            <div class="space-y-1">
                                <p class="font-medium tabular-nums">
                                    {{ row.original.members }}
                                </p>

                                <p class="text-xs text-muted">
                                    {{ row.original.manualMembers }} manual · {{ row.original.derivedMembers }} derived ·
                                    {{ row.original.excludedMembers }} excluded
                                </p>
                            </div>
                        </template>

                        <template #leaders-cell="{ row }">
                            <div v-if="row.original.leaders.length" class="flex flex-wrap gap-1">
                                <UBadge
                                    v-for="leader in row.original.leaders"
                                    :key="leader.name"
                                    color="neutral"
                                    variant="soft"
                                >
                                    {{ leader.name }}
                                </UBadge>
                            </div>

                            <span v-else class="text-sm text-muted"> No leader </span>
                        </template>

                        <template #rules-cell="{ row }">
                            <div v-if="row.original.rules.length" class="flex flex-wrap gap-1">
                                <UBadge v-for="rule in row.original.rules" :key="rule" color="primary" variant="soft">
                                    {{ rule }}
                                </UBadge>
                            </div>

                            <span v-else class="text-sm text-muted"> Manual only </span>
                        </template>

                        <template #updatedAt-cell="{ row }">
                            <span class="text-sm text-muted">
                                {{ row.original.updatedAt }}
                            </span>
                        </template>

                        <template #action-cell="{ row }">
                            <UDropdownMenu :items="getDropdownActions(row.original)">
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
    </UDashboardPanel>
</template>
