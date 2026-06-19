<script lang="ts" setup>
import ColleagueCard from '~/components/jobs/colleagues/ColleagueCard.vue';
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCompletorStore } from '~/stores/completor';
import { getJobsGroupsClient } from '~~/gen/ts/clients';
import type { GroupLeader } from '~~/gen/ts/resources/jobs/groups/group';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import type { ListGroupLeadersResponse } from '~~/gen/ts/services/jobs/groups';

const props = defineProps<{
    groupId: number;
    canManage: boolean;
}>();

const emit = defineEmits<{
    changed: [];
}>();

const { t } = useI18n();
const overlay = useOverlay();
const completorStore = useCompletorStore();
const jobsGroupsClient = await getJobsGroupsClient();
const confirmModal = overlay.create(ConfirmModal);

const page = ref(1);
const selectedLeader = ref<UserShort>();
const pendingAction = ref<string>();

const leadersKey = computed(() => `jobs-group-leaders-${props.groupId}-${page.value}`);

const {
    data: leadersData,
    status: leadersStatus,
    error: leadersError,
    refresh: refreshLeaders,
} = useLazyAsyncData(leadersKey, () => listGroupLeaders(), {
    watch: [() => props.groupId, page],
});

const leaders = computed<GroupLeader[]>(() => leadersData.value?.leaders ?? []);
const leaderIds = computed(() => new Set(leaders.value.map((leader) => leader.userId)));
const isMutating = computed(() => pendingAction.value !== undefined);

async function listGroupLeaders(): Promise<ListGroupLeadersResponse> {
    const { response } = await jobsGroupsClient.listGroupLeaders({
        groupId: props.groupId,
        pagination: {
            offset: calculateOffset(page.value, leadersData.value?.pagination),
        },
    });

    return response;
}

function resetLeaderForm(): void {
    selectedLeader.value = undefined;
}

async function runMutation(action: string, mutate: () => Promise<void>): Promise<void> {
    pendingAction.value = action;
    try {
        await mutate();
        await refreshLeaders();
        emit('changed');
    } catch (e) {
        handleGRPCError(e as RpcError);
    } finally {
        pendingAction.value = undefined;
    }
}

async function addLeader(): Promise<void> {
    if (!selectedLeader.value?.userId) return;

    await runMutation('leader', async () => {
        await jobsGroupsClient.addGroupLeader({
            groupId: props.groupId,
            userId: selectedLeader.value!.userId,
        });
        resetLeaderForm();
    });
}

async function removeLeader(userId: number): Promise<void> {
    confirmModal.open({
        title: t('common.remove'),
        confirm: async () =>
            await runMutation(`leader-${userId}`, async () => {
                await jobsGroupsClient.removeGroupLeader({
                    groupId: props.groupId,
                    userId,
                });
                if (selectedLeader.value?.userId === userId) resetLeaderForm();
            }),
    });
}

watch(
    () => props.groupId,
    () => {
        page.value = 1;
        resetLeaderForm();
    },
);
</script>

<template>
    <div class="grid gap-4">
        <UCard v-if="canManage" variant="subtle">
            <div class="grid gap-3 lg:grid-cols-[minmax(0,1fr)_auto] lg:items-end">
                <UFormField :label="$t('common.colleague', 1)">
                    <SelectMenu
                        v-model="selectedLeader"
                        class="w-full"
                        :searchable="
                            async (q: string) =>
                                await completorStore.completeColleagues(
                                    q,
                                    selectedLeader?.userId ? [selectedLeader.userId] : [],
                                )
                        "
                        searchable-key="jobs-group-leaders"
                        :filter-fields="['firstname', 'lastname']"
                        :search-input="{ placeholder: $t('common.search_field') }"
                        :placeholder="$t('common.colleague', 1)"
                        :disabled="isMutating || !canManage"
                    >
                        <template v-if="selectedLeader" #default>
                            {{ userToLabel(selectedLeader) }}
                        </template>
                        <template #item-label="{ item }">
                            {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                        </template>
                        <template #empty>
                            {{ $t('common.not_found', [$t('common.colleague', 2)]) }}
                        </template>
                    </SelectMenu>
                </UFormField>

                <UFieldGroup class="inline-flex w-full sm:w-auto">
                    <UButton
                        icon="i-mdi-account-star"
                        :label="$t('common.add')"
                        :loading="pendingAction === 'leader'"
                        :disabled="isMutating || !canManage || !selectedLeader?.userId || leaderIds.has(selectedLeader.userId)"
                        @click="addLeader"
                    />
                    <UButton
                        v-if="selectedLeader"
                        color="neutral"
                        variant="outline"
                        icon="i-mdi-close"
                        :label="$t('common.cancel')"
                        :disabled="isMutating || !canManage"
                        @click="resetLeaderForm"
                    />
                </UFieldGroup>
            </div>
        </UCard>

        <DataPendingBlock
            v-if="isRequestPending(leadersStatus)"
            :message="$t('common.loading', [$t('components.jobs.groups.leaders')])"
        />
        <DataErrorBlock
            v-else-if="leadersError"
            :title="$t('common.unable_to_load', [$t('components.jobs.groups.leaders')])"
            :error="leadersError"
            :retry="refreshLeaders"
        />
        <DataNoDataBlock
            v-else-if="leadersData?.pagination?.totalCount === 0"
            :type="$t('components.jobs.groups.leaders')"
            icon="i-mdi-account-star"
            :padded="false"
        />
        <template v-else>
            <ColleagueCard
                v-for="leader in leaders"
                :key="leader.userId"
                compact
                :colleague="leader.colleague"
                :user-id="leader.userId"
            >
                <template #badges>
                    <UBadge color="warning" variant="soft">
                        {{ $t('components.jobs.groups.details.leader') }}
                    </UBadge>
                </template>

                <template #footer>
                    <div class="flex flex-wrap items-center justify-between gap-2">
                        <p class="text-sm text-muted">
                            {{ $t('common.created_by') }}
                            <ColleagueInfoPopover :user="leader.createdBy" :user-id="leader.createdByUserId" hide-props />
                        </p>

                        <GenericTime v-if="leader.createdAt" class="text-sm text-muted" :value="leader.createdAt" />

                        <UButton
                            v-if="canManage"
                            color="error"
                            variant="outline"
                            icon="i-mdi-account-star-outline"
                            :label="$t('common.remove')"
                            :loading="pendingAction === `leader-${leader.userId}`"
                            :disabled="isMutating || !canManage"
                            @click="removeLeader(leader.userId)"
                        />
                    </div>
                </template>
            </ColleagueCard>

            <Pagination
                v-model="page"
                :pagination="leadersData?.pagination"
                :status="leadersStatus"
                :refresh="refreshLeaders"
            />
        </template>
    </div>
</template>
