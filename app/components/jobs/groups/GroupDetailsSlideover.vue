<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getJobsGroupsClient } from '~~/gen/ts/clients';
import { type Group, GroupState } from '~~/gen/ts/resources/jobs/groups/group';
import type { GetGroupResponse } from '~~/gen/ts/services/jobs/groups';
import GroupDetailsActivityPanel from '~/components/jobs/groups/GroupDetailsActivityPanel.vue';
import GroupDetailsExclusionsPanel from '~/components/jobs/groups/GroupDetailsExclusionsPanel.vue';
import GroupDetailsLeadersPanel from '~/components/jobs/groups/GroupDetailsLeadersPanel.vue';
import GroupDetailsManualMembersPanel from '~/components/jobs/groups/GroupDetailsManualMembersPanel.vue';
import GroupDetailsMembersPanel from '~/components/jobs/groups/GroupDetailsMembersPanel.vue';
import GroupDetailsRulesPanel from '~/components/jobs/groups/GroupDetailsRulesPanel.vue';
import {
    groupMembershipModeColor,
    groupMembershipModeIcon,
    groupMembershipModeLabelKey,
    groupStateColor,
    groupStateLabelKey,
    groupTypeColor,
    groupTypeIcon,
    groupTypeLabelKey,
} from './helpers';

const props = defineProps<{
    group: Group;
}>();

const { t } = useI18n();
const { can } = useAuth();
const jobsGroupsClient = await getJobsGroupsClient();

const selectedTab = ref('members');

const tabs = computed(() => [
    {
        slot: 'members' as const,
        label: t('common.members', 2),
        icon: 'i-mdi-account-multiple',
        value: 'members',
    },
    {
        slot: 'rules' as const,
        label: t('components.jobs.groups.rules'),
        icon: 'i-mdi-filter-cog',
        value: 'rules',
    },
    {
        slot: 'manualMembers' as const,
        label: t('components.jobs.groups.manual_members'),
        icon: 'i-mdi-account-plus',
        value: 'manualMembers',
    },
    {
        slot: 'leaders' as const,
        label: t('components.jobs.groups.leaders'),
        icon: 'i-mdi-account-star',
        value: 'leaders',
    },
    {
        slot: 'exclusions' as const,
        label: t('components.jobs.groups.exclusions'),
        icon: 'i-mdi-account-cancel',
        value: 'exclusions',
    },
    {
        slot: 'activity' as const,
        label: t('common.activity'),
        icon: 'i-mdi-history',
        value: 'activity',
    },
]);

const detailKey = computed(() => `jobs-group-details-${props.group.id}`);

const {
    data: detail,
    status: detailStatus,
    error: detailError,
    refresh: refreshDetail,
} = useLazyAsyncData(detailKey, () => getGroupDetails(props.group.id), {
    watch: [() => props.group.id],
});

const currentGroup = computed(() => detail.value?.group ?? props.group);
const canManageGroups = can('jobs.GroupsService/CreateGroup');
const canManageLeaders = can('jobs.GroupsService/AddGroupLeader');
const groupIsArchived = computed(() => currentGroup.value.state === GroupState.ARCHIVED);
const canMutateGroup = computed(() => canManageGroups.value && !groupIsArchived.value);
const canMutateLeaders = computed(() => canManageLeaders.value && !groupIsArchived.value);

async function getGroupDetails(groupId: number): Promise<GetGroupResponse> {
    const { response } = await jobsGroupsClient.getGroup({
        id: groupId,
        includeRules: false,
        includeLeaders: false,
        includeManualMembers: false,
        includeExclusions: false,
        includeResolvedMembers: false,
        includeArchived: true,
    });

    return response;
}

async function handlePanelChanged(): Promise<void> {
    await refreshDetail();
}
</script>

<template>
    <USlideover :title="currentGroup.name" :overlay="false" :ui="{ content: 'max-w-5xl' }">
        <template #body>
            <div class="flex flex-col gap-4">
                <UCard>
                    <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
                        <div class="min-w-0">
                            <div class="flex flex-wrap items-center gap-2">
                                <h2 class="text-xl font-semibold text-highlighted">
                                    {{ currentGroup.name }}
                                </h2>

                                <UBadge
                                    :color="groupTypeColor(currentGroup.type)"
                                    :icon="groupTypeIcon(currentGroup.type)"
                                    :label="$t(groupTypeLabelKey(currentGroup.type))"
                                    variant="subtle"
                                />

                                <UBadge
                                    :color="groupMembershipModeColor(currentGroup.membershipMode)"
                                    :icon="groupMembershipModeIcon(currentGroup.membershipMode)"
                                    :label="$t(groupMembershipModeLabelKey(currentGroup.membershipMode))"
                                    variant="subtle"
                                />

                                <UBadge
                                    :color="groupStateColor(currentGroup.state)"
                                    :label="$t(groupStateLabelKey(currentGroup.state))"
                                    variant="subtle"
                                />
                            </div>

                            <p class="mt-2 text-sm text-muted">
                                {{ currentGroup.description || $t('components.jobs.groups.no_description') }}
                            </p>
                        </div>

                        <div class="grid min-w-48 grid-cols-2 gap-2 text-sm">
                            <div class="rounded-lg bg-elevated px-3 py-2">
                                <p class="text-muted">{{ $t('components.jobs.groups.leaders') }}</p>
                                <p class="text-lg font-semibold tabular-nums">{{ currentGroup.leadersCount }}</p>
                            </div>
                            <div class="rounded-lg bg-elevated px-3 py-2">
                                <p class="text-muted">{{ $t('components.jobs.groups.exclusions') }}</p>
                                <p class="text-lg font-semibold tabular-nums">{{ currentGroup.exclusionsCount }}</p>
                            </div>

                            <div class="rounded-lg bg-elevated px-3 py-2">
                                <p class="text-muted">{{ $t('common.members', 2) }}</p>
                                <p class="text-lg font-semibold tabular-nums">{{ currentGroup.membersCount }}</p>
                            </div>
                            <div class="rounded-lg bg-elevated px-3 py-2">
                                <p class="text-muted">{{ $t('components.jobs.groups.rules') }}</p>
                                <p class="text-lg font-semibold tabular-nums">{{ currentGroup.rulesCount }}</p>
                            </div>
                        </div>
                    </div>
                </UCard>

                <DataPendingBlock v-if="isRequestPending(detailStatus)" :message="$t('common.loading', [$t('common.group')])" />
                <DataErrorBlock
                    v-else-if="detailError"
                    :title="$t('common.unable_to_load', [$t('common.group')])"
                    :error="detailError"
                    :retry="refreshDetail"
                />
                <UTabs
                    v-else
                    v-model="selectedTab"
                    :items="tabs"
                    :unmount-on-hide="true"
                    variant="link"
                    :ui="{ content: 'pt-2' }"
                >
                    <template #members>
                        <GroupDetailsMembersPanel :group-id="currentGroup.id" />
                    </template>

                    <template #rules>
                        <GroupDetailsRulesPanel
                            :group-id="currentGroup.id"
                            :can-manage="canMutateGroup"
                            @changed="handlePanelChanged"
                        />
                    </template>

                    <template #manualMembers>
                        <GroupDetailsManualMembersPanel
                            :group-id="currentGroup.id"
                            :can-manage="canMutateGroup"
                            @changed="handlePanelChanged"
                        />
                    </template>

                    <template #leaders>
                        <GroupDetailsLeadersPanel
                            :group-id="currentGroup.id"
                            :can-manage="canMutateLeaders"
                            @changed="handlePanelChanged"
                        />
                    </template>

                    <template #exclusions>
                        <GroupDetailsExclusionsPanel
                            :group-id="currentGroup.id"
                            :can-manage="canMutateGroup"
                            @changed="handlePanelChanged"
                        />
                    </template>

                    <template #activity>
                        <GroupDetailsActivityPanel :group-id="currentGroup.id" />
                    </template>
                </UTabs>
            </div>
        </template>
    </USlideover>
</template>
