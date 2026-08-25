<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { getJobsGroupsClient } from '~~/gen/ts/clients';
import {
    GroupMemberSource,
    type GroupMembershipReason,
    GroupMembershipReasonType,
    type GroupResolvedMember,
} from '~~/gen/ts/resources/jobs/groups/group';
import type { ListGroupMembersResponse } from '~~/gen/ts/services/jobs/groups';
import ColleagueCard from '~/components/jobs/colleagues/ColleagueCard.vue';

const props = defineProps<{
    groupId: number;
    canView: boolean;
}>();

const { t } = useI18n();
const jobsGroupsClient = await getJobsGroupsClient();

const page = ref(1);
const search = ref('');
const selectedSources = ref<GroupMemberSource[]>([]);

const membersKey = computed(
    () => `jobs-group-members-${props.groupId}-${page.value}-${search.value}-${selectedSources.value.join(',')}`,
);

const {
    data: membersData,
    status: membersStatus,
    error: membersError,
    refresh: refreshMembers,
} = useLazyAsyncData(membersKey, () => listGroupMembers(), {
    watch: [() => props.groupId, page],
});

const members = computed<GroupResolvedMember[]>(() => membersData.value?.members ?? []);

const sourceItems = computed(() =>
    Object.values(GroupMemberSource)
        .filter((value): value is GroupMemberSource => typeof value === 'number' && value !== GroupMemberSource.UNSPECIFIED)
        .map((value) => ({
            label: t(sourceLabelKey(value)),
            value,
        })),
);

async function listGroupMembers(): Promise<ListGroupMembersResponse> {
    const { response } = await jobsGroupsClient.listGroupMembers({
        groupId: props.groupId,
        pagination: {
            offset: calculateOffset(page.value, membersData.value?.pagination),
        },
        sort: { columns: [{ id: 'user_id', desc: false }] },
        search: search.value.trim() || undefined,
        includeExcluded: true,
        includeLeaders: true,
        includeReasons: true,
        sources: selectedSources.value,
    });

    return response;
}

async function applyFilters(): Promise<void> {
    if (page.value === 1) {
        await refreshMembers();
        return;
    }

    page.value = 1;
}

async function clearFilters(): Promise<void> {
    search.value = '';
    selectedSources.value = [];
    if (page.value === 1) {
        await refreshMembers();
        return;
    }

    page.value = 1;
}

function sourceLabelKey(source: GroupMemberSource): string {
    return `enums.jobs.groups.GroupMemberSource.${GroupMemberSource[source] ?? 'UNSPECIFIED'}`;
}

function reasonLabelKey(reason: GroupMembershipReason): string {
    return `enums.jobs.groups.GroupMembershipReasonType.${GroupMembershipReasonType[reason.type] ?? 'UNSPECIFIED'}`;
}

function memberTone(member: GroupResolvedMember): 'success' | 'warning' | 'neutral' | 'error' {
    if (member.isExcluded) return 'error';
    if (member.isLeader) return 'warning';
    if (member.isMember) return 'success';
    return 'neutral';
}

function hasManualReason(member: GroupResolvedMember): boolean {
    return member.reasons.some((reason) => reason.type === GroupMembershipReasonType.MANUAL);
}

watch(
    () => props.groupId,
    () => (page.value = 1),
);
</script>

<template>
    <div v-if="canView" class="grid gap-4">
        <UCard variant="subtle">
            <div class="grid gap-3 lg:grid-cols-[minmax(0,1fr)_minmax(240px,320px)_auto] lg:items-end">
                <UFormField :label="$t('common.search')">
                    <UInput
                        v-model="search"
                        class="w-full"
                        icon="i-mdi-magnify"
                        :placeholder="$t('common.search_field')"
                        @keyup.enter="applyFilters"
                    />
                </UFormField>

                <UFormField :label="$t('common.type')">
                    <USelectMenu
                        v-model="selectedSources"
                        class="w-full"
                        multiple
                        :items="sourceItems"
                        value-key="value"
                        :search-input="{ placeholder: $t('common.search_field') }"
                    />
                </UFormField>

                <UFieldGroup class="inline-flex w-full sm:w-auto">
                    <UButton
                        color="neutral"
                        variant="outline"
                        icon="i-mdi-filter-remove"
                        :label="$t('common.clear')"
                        @click="clearFilters"
                    />
                    <UButton icon="i-mdi-filter" :label="$t('common.apply')" @click="applyFilters" />
                </UFieldGroup>
            </div>
        </UCard>

        <DataPendingBlock v-if="isRequestPending(membersStatus)" :message="$t('common.loading', [$t('common.members', 2)])" />
        <DataErrorBlock
            v-else-if="membersError"
            :title="$t('common.unable_to_load', [$t('common.members', 2)])"
            :error="membersError"
            :retry="refreshMembers"
        />
        <DataNoDataBlock
            v-else-if="members.length === 0"
            :type="$t('common.members', 2)"
            icon="i-mdi-account-multiple"
            :padded="false"
        />
        <div v-else class="grid gap-2">
            <ColleagueCard
                v-for="member in members"
                :key="member.userId"
                compact
                :colleague="member.colleague"
                :user-id="member.userId"
            >
                <template #badges>
                    <UBadge :color="memberTone(member)" variant="subtle">
                        {{
                            member.isExcluded
                                ? $t('components.jobs.groups.details.excluded')
                                : member.isMember
                                  ? $t('components.jobs.groups.details.member')
                                  : member.isLeader
                                    ? $t('components.jobs.groups.details.leader')
                                    : $t('common.info')
                        }}
                    </UBadge>

                    <UBadge v-if="hasManualReason(member)" color="neutral" variant="soft">
                        {{ $t('components.jobs.groups.details.manual') }}
                    </UBadge>
                </template>

                <div v-if="member.reasons.length" class="line-clamp-2 max-w-xl text-sm text-muted">
                    <span class="font-semibold">{{ $t('common.reason') }}:</span>

                    <p v-for="(reason, idx) in member.reasons" :key="`${member.userId}-${idx}`">
                        {{ $t(reasonLabelKey(reason)) }}
                        <template v-if="reason.ruleId"> #{{ reason.ruleId }} </template>
                        <template v-if="reason.detail"> - {{ reason.detail }} </template>
                    </p>
                </div>
            </ColleagueCard>
        </div>

        <Pagination v-model="page" :pagination="membersData?.pagination" :status="membersStatus" :refresh="refreshMembers" />
    </div>
    <DataNoDataBlock v-else :message="$t('common.no_access')" icon="i-mdi-lock" :padded="false" />
</template>
