<script lang="ts" setup>
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import ConfirmModalWithReason from '~/components/partials/ConfirmModalWithReason.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getJobsGroupsClient } from '~~/gen/ts/clients';
import { AccessLevel as GroupAccessLevel } from '~~/gen/ts/resources/jobs/groups/access/access';
import { type Group, GroupState } from '~~/gen/ts/resources/jobs/groups/group';
import type { GetGroupResponse } from '~~/gen/ts/services/jobs/groups';
import EditorModal from '~/components/jobs/groups/EditorModal.vue';
import ActivityPanel from '~/components/jobs/groups/details/ActivityPanel.vue';
import ExclusionsPanel from '~/components/jobs/groups/details/ExclusionsPanel.vue';
import LeadersPanel from '~/components/jobs/groups/details/LeadersPanel.vue';
import ManualMembersPanel from '~/components/jobs/groups/details/ManualMembersPanel.vue';
import MembersPanel from '~/components/jobs/groups/details/MembersPanel.vue';
import RulesPanel from '~/components/jobs/groups/details/RulesPanel.vue';
import { isLegacyGroupPolicyState } from './policy';
import {
    checkGroupAccess,
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

const emit = defineEmits<{
    changed: [];
}>();

const { t } = useI18n();
const { can } = useAuth();
const overlay = useOverlay();
const jobsGroupsClient = await getJobsGroupsClient();
const editorModal = overlay.create(EditorModal);
const confirmModal = overlay.create(ConfirmModal);
const confirmModalWithReason = overlay.create(ConfirmModalWithReason);

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
const groupIsArchived = computed(() => currentGroup.value.state === GroupState.ARCHIVED);
const legacyPolicyState = computed(() => isLegacyGroupPolicyState(currentGroup.value));
const currentGroupAccess = computed(() => detail.value?.access);
const canViewGroup = computed(
    () => can('jobs.GroupsService/ListGroups').value && checkGroupAccess(currentGroupAccess.value, GroupAccessLevel.VIEW),
);
const canMutateGroup = computed(
    () =>
        !groupIsArchived.value &&
        can('jobs.GroupsService/ListGroups').value &&
        checkGroupAccess(currentGroupAccess.value, GroupAccessLevel.EDIT),
);
const canMutateLeaders = computed(
    () =>
        !groupIsArchived.value &&
        can('jobs.GroupsService/ListGroups').value &&
        checkGroupAccess(currentGroupAccess.value, GroupAccessLevel.MANAGE),
);
const canManageArchiveState = computed(
    () => can('jobs.GroupsService/ArchiveGroup').value && checkGroupAccess(currentGroupAccess.value, GroupAccessLevel.MANAGE),
);
const slideoverTitle = computed(() => (canViewGroup.value ? currentGroup.value.name : t('common.no_access')));

function emitChanged(): void {
    emit('changed');
}

async function openEditGroup(): Promise<void> {
    if (!canMutateGroup.value) return;

    editorModal.open({
        group: currentGroup.value,
        access: currentGroupAccess.value,
        onUpdated: async () => {
            await refreshDetail();
            emitChanged();
        },
    });
}

async function archiveGroup(): Promise<void> {
    if (!canManageArchiveState.value || groupIsArchived.value) return;

    confirmModalWithReason.open({
        title: t('components.jobs.groups.actions.archive'),
        description: t('components.jobs.groups.confirm.archive', { name: currentGroup.value.name }),
        confirm: async (reason: string) => {
            await jobsGroupsClient.archiveGroup({
                id: currentGroup.value.id,
                reason,
            });

            await refreshDetail();
            emitChanged();
        },
    });
}

async function restoreGroup(): Promise<void> {
    if (!canManageArchiveState.value || !groupIsArchived.value) return;

    confirmModal.open({
        title: t('components.jobs.groups.actions.restore'),
        description: t('components.jobs.groups.confirm.restore', { name: currentGroup.value.name }),
        confirm: async () => {
            await jobsGroupsClient.restoreGroup({ id: currentGroup.value.id });

            await refreshDetail();
            emitChanged();
        },
    });
}

async function getGroupDetails(groupId: number): Promise<GetGroupResponse> {
    const { response } = await jobsGroupsClient.getGroup({
        id: groupId,
        includeArchived: true,
    });

    return response;
}

async function handlePanelChanged(): Promise<void> {
    await refreshDetail();
}
</script>

<template>
    <USlideover :title="slideoverTitle" :overlay="false" :ui="{ content: 'max-w-5xl' }">
        <template #body>
            <div class="flex flex-col gap-4">
                <DataPendingBlock v-if="isRequestPending(detailStatus)" :message="$t('common.loading', [$t('common.group')])" />
                <DataErrorBlock
                    v-else-if="detailError"
                    :title="$t('common.unable_to_load', [$t('common.group')])"
                    :error="detailError"
                    :retry="refreshDetail"
                />
                <template v-else-if="canViewGroup">
                    <UCard>
                        <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
                            <div class="min-w-0">
                                <div class="flex flex-wrap items-center gap-2">
                                    <h2 class="text-xl font-semibold text-highlighted">
                                        {{ currentGroup.name }}
                                    </h2>

                                    <UBadge
                                        v-if="currentGroup.shortName && currentGroup.shortName !== currentGroup.name"
                                        color="neutral"
                                        :label="currentGroup.shortName"
                                        variant="soft"
                                        size="sm"
                                    />

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

                            <div class="flex flex-col gap-3">
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
                                        <p class="text-lg font-semibold tabular-nums">
                                            {{ $t('components.jobs.groups.approx_member_count', currentGroup.membersCount) }}
                                        </p>
                                    </div>
                                    <div class="rounded-lg bg-elevated px-3 py-2">
                                        <p class="text-muted">{{ $t('components.jobs.groups.rules') }}</p>
                                        <p class="text-lg font-semibold tabular-nums">{{ currentGroup.rulesCount }}</p>
                                    </div>
                                </div>

                                <UFieldGroup
                                    v-if="canMutateGroup || canManageArchiveState"
                                    class="inline-flex justify-end gap-2"
                                >
                                    <UButton
                                        v-if="canMutateGroup"
                                        color="neutral"
                                        variant="outline"
                                        icon="i-mdi-pencil"
                                        :label="$t('common.edit')"
                                        @click="openEditGroup"
                                    />

                                    <UButton
                                        v-if="!groupIsArchived && canManageArchiveState"
                                        color="warning"
                                        variant="outline"
                                        icon="i-mdi-archive"
                                        :label="$t('components.jobs.groups.actions.archive')"
                                        @click="archiveGroup"
                                    />

                                    <UButton
                                        v-else-if="groupIsArchived && canManageArchiveState"
                                        color="primary"
                                        variant="outline"
                                        icon="i-mdi-restore"
                                        :label="$t('components.jobs.groups.actions.restore')"
                                        @click="restoreGroup"
                                    />
                                </UFieldGroup>
                            </div>
                        </div>
                    </UCard>

                    <UAlert
                        v-if="legacyPolicyState"
                        color="warning"
                        icon="i-mdi-alert-circle"
                        :title="$t('components.jobs.groups.policy.legacy_state_title')"
                        :description="$t('components.jobs.groups.policy.legacy_state_content')"
                    />

                    <UTabs v-model="selectedTab" :items="tabs" :unmount-on-hide="true" variant="link" :ui="{ content: 'pt-2' }">
                        <template #members>
                            <MembersPanel :group-id="currentGroup.id" :can-view="canViewGroup" />
                        </template>

                        <template #rules>
                            <RulesPanel
                                :group-id="currentGroup.id"
                                :group-type="currentGroup.type"
                                :access="currentGroupAccess"
                                :can-view="canViewGroup"
                                :can-manage="canMutateGroup"
                                @changed="handlePanelChanged"
                            />
                        </template>

                        <template #manualMembers>
                            <ManualMembersPanel
                                :group-id="currentGroup.id"
                                :group-type="currentGroup.type"
                                :access="currentGroupAccess"
                                :can-view="canViewGroup"
                                :can-manage="canMutateGroup"
                                @changed="handlePanelChanged"
                            />
                        </template>

                        <template #leaders>
                            <LeadersPanel
                                :group-id="currentGroup.id"
                                :access="currentGroupAccess"
                                :can-view="canViewGroup"
                                :can-manage="canMutateLeaders"
                                @changed="handlePanelChanged"
                            />
                        </template>

                        <template #exclusions>
                            <ExclusionsPanel
                                :group-id="currentGroup.id"
                                :group-type="currentGroup.type"
                                :access="currentGroupAccess"
                                :can-view="canViewGroup"
                                :can-manage="canMutateGroup"
                                @changed="handlePanelChanged"
                            />
                        </template>

                        <template #activity>
                            <ActivityPanel :group-id="currentGroup.id" :can-view="canViewGroup" />
                        </template>
                    </UTabs>
                </template>
                <DataNoDataBlock v-else :message="$t('common.no_access')" icon="i-mdi-lock" :padded="false" />
            </div>
        </template>
    </USlideover>
</template>
