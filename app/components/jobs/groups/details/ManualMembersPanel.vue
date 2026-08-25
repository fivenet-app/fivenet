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
import type { Access } from '~~/gen/ts/resources/access/access';
import { AccessLevel as GroupAccessLevel } from '~~/gen/ts/resources/jobs/groups/access/access';
import type { GroupType, GroupManualMember } from '~~/gen/ts/resources/jobs/groups/group';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import type { ListGroupManualMembersResponse } from '~~/gen/ts/services/jobs/groups';
import { checkGroupAccess } from '../helpers';
import { groupTypeAllowsManualMembers } from '../policy';

const props = defineProps<{
    groupId: number;
    groupType: GroupType;
    canView: boolean;
    canManage: boolean;
    access?: Access;
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
const selectedManualMember = ref<UserShort>();
const reason = ref('');
const editingManualMemberId = ref<number>();
const pendingAction = ref<string>();

const manualMembersKey = computed(() => `jobs-group-manual-members-${props.groupId}-${page.value}`);

const {
    data: manualMembersData,
    status: manualMembersStatus,
    error: manualMembersError,
    refresh: refreshManualMembers,
} = useLazyAsyncData(manualMembersKey, () => listGroupManualMembers(), {
    watch: [() => props.groupId, page],
});

const manualMembers = computed<GroupManualMember[]>(() => manualMembersData.value?.manualMembers ?? []);
const manualMemberIds = computed(() => new Set(manualMembers.value.map((member) => member.userId)));
const isMutating = computed(() => pendingAction.value !== undefined);
const supportsManualMembers = computed(() => groupTypeAllowsManualMembers(props.groupType));
const canManageMembers = computed(
    () => props.canManage && supportsManualMembers.value && checkGroupAccess(props.access, GroupAccessLevel.EDIT),
);
const policyNotice = computed(() =>
    supportsManualMembers.value ? undefined : 'components.jobs.groups.policy.manual_members_disabled',
);

async function listGroupManualMembers(): Promise<ListGroupManualMembersResponse> {
    const { response } = await jobsGroupsClient.listGroupManualMembers({
        groupId: props.groupId,
        pagination: {
            offset: calculateOffset(page.value, manualMembersData.value?.pagination),
        },
    });

    return response;
}

function resetManualMemberForm(): void {
    editingManualMemberId.value = undefined;
    selectedManualMember.value = undefined;
    reason.value = '';
}

function editManualMember(member: GroupManualMember): void {
    if (!canManageMembers.value) return;

    editingManualMemberId.value = member.userId;
    selectedManualMember.value = {
        userId: member.userId,
        job: member.colleague?.job ?? '',
        jobGrade: member.colleague?.jobGrade ?? 0,
        firstname: member.colleague?.firstname ?? `${t('common.id')}: ${member.userId}`,
        lastname: member.colleague?.lastname ?? '',
        dateofbirth: member.colleague?.dateofbirth ?? '',
    };
    reason.value = member.reason ?? '';
}

async function runMutation(action: string, mutate: () => Promise<void>): Promise<void> {
    pendingAction.value = action;
    try {
        await mutate();
        await refreshManualMembers();
        emit('changed');
    } catch (e) {
        handleGRPCError(e as RpcError);
    } finally {
        pendingAction.value = undefined;
    }
}

async function addManualMember(): Promise<void> {
    if (!canManageMembers.value) return;
    if (!selectedManualMember.value?.userId) return;

    await runMutation('manual-member', async () => {
        await jobsGroupsClient.addGroupMember({
            groupId: props.groupId,
            userId: selectedManualMember.value!.userId,
            reason: reason.value.trim() || undefined,
        });
        resetManualMemberForm();
    });
}

async function removeManualMember(userId: number): Promise<void> {
    if (!canManageMembers.value) return;

    confirmModal.open({
        title: t('common.remove'),
        confirm: async () =>
            await runMutation(`manual-member-${userId}`, async () => {
                await jobsGroupsClient.removeGroupMember({
                    groupId: props.groupId,
                    userId,
                });
                if (editingManualMemberId.value === userId) resetManualMemberForm();
            }),
    });
}

watch(
    () => props.groupId,
    () => {
        page.value = 1;
        resetManualMemberForm();
    },
);
</script>

<template>
    <div v-if="canView" class="grid gap-4">
        <UAlert
            v-if="policyNotice"
            color="warning"
            icon="i-mdi-alert-circle"
            :title="$t('components.jobs.groups.policy.unavailable_title')"
            :description="$t(policyNotice)"
        />

        <UCard v-if="canManageMembers" variant="subtle">
            <div class="grid gap-3">
                <div class="grid gap-3 lg:flex lg:flex-row lg:items-end">
                    <UFormField class="flex-1" :label="$t('common.colleague', 1)">
                        <SelectMenu
                            v-model="selectedManualMember"
                            class="w-full"
                            :searchable="
                                async (q: string) =>
                                    await completorStore.completeColleagues(
                                        q,
                                        selectedManualMember?.userId ? [selectedManualMember.userId] : [],
                                    )
                            "
                            searchable-key="jobs-group-manual-members"
                            :filter-fields="['firstname', 'lastname']"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :placeholder="$t('common.colleague', 1)"
                            :disabled="isMutating || !canManageMembers"
                        >
                            <template v-if="selectedManualMember" #default>
                                {{ userToLabel(selectedManualMember) }}
                            </template>
                            <template #item-label="{ item }">
                                {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                            </template>
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.colleague', 2)]) }}
                            </template>
                        </SelectMenu>
                    </UFormField>

                    <UFieldGroup class="inline-flex">
                        <UButton
                            :icon="editingManualMemberId ? 'i-mdi-content-save' : 'i-mdi-account-plus'"
                            :label="
                                editingManualMemberId
                                    ? $t('common.save', 1)
                                    : $t('components.jobs.groups.details.add_manual_member')
                            "
                            :loading="pendingAction === 'manual-member'"
                            :disabled="
                                isMutating ||
                                !canManageMembers ||
                                !selectedManualMember?.userId ||
                                (manualMemberIds.has(selectedManualMember.userId) &&
                                    editingManualMemberId !== selectedManualMember.userId)
                            "
                            @click="addManualMember"
                        />
                        <UButton
                            v-if="editingManualMemberId"
                            color="neutral"
                            variant="outline"
                            icon="i-mdi-close"
                            :label="$t('common.cancel')"
                            :disabled="isMutating || !canManageMembers"
                            @click="resetManualMemberForm"
                        />
                    </UFieldGroup>
                </div>

                <UFormField :label="$t('common.reason', 1)">
                    <UTextarea
                        v-model="reason"
                        class="w-full"
                        :rows="2"
                        :placeholder="$t('common.reason', 1)"
                        :disabled="isMutating || !canManageMembers"
                    />
                </UFormField>
            </div>
        </UCard>

        <DataPendingBlock
            v-if="isRequestPending(manualMembersStatus)"
            :message="$t('common.loading', [$t('components.jobs.groups.manual_members')])"
        />
        <DataErrorBlock
            v-else-if="manualMembersError"
            :title="$t('common.unable_to_load', [$t('components.jobs.groups.manual_members')])"
            :error="manualMembersError"
            :retry="refreshManualMembers"
        />
        <DataNoDataBlock
            v-else-if="manualMembersData?.pagination?.totalCount === 0"
            :type="$t('components.jobs.groups.manual_members')"
            icon="i-mdi-account-plus"
            :padded="false"
        />
        <template v-else>
            <ColleagueCard
                v-for="member in manualMembers"
                :key="member.userId"
                compact
                :colleague="member.colleague"
                :user-id="member.userId"
            >
                <p v-if="member.reason">{{ member.reason }}</p>

                <template #footer>
                    <div class="flex flex-wrap items-center justify-between gap-2">
                        <p class="text-sm text-muted">
                            {{ $t('common.created_by') }}
                            <ColleagueInfoPopover :user="member.createdBy" :user-id="member.createdByUserId" hide-props />
                        </p>

                        <GenericTime v-if="member.createdAt" class="text-sm text-muted" :value="member.createdAt" />

                        <UFieldGroup v-if="canManageMembers">
                            <UButton
                                color="neutral"
                                variant="outline"
                                icon="i-mdi-pencil"
                                :label="$t('common.edit')"
                                :disabled="isMutating || !canManageMembers"
                                @click="editManualMember(member)"
                            />
                            <UButton
                                color="error"
                                variant="outline"
                                icon="i-mdi-account-minus"
                                :label="$t('common.remove')"
                                :loading="pendingAction === `manual-member-${member.userId}`"
                                :disabled="isMutating || !canManageMembers"
                                @click="removeManualMember(member.userId)"
                            />
                        </UFieldGroup>
                    </div>
                </template>
            </ColleagueCard>

            <Pagination
                v-model="page"
                :pagination="manualMembersData?.pagination"
                :status="manualMembersStatus"
                :refresh="refreshManualMembers"
            />
        </template>
    </div>
    <DataNoDataBlock v-else :message="$t('common.no_access')" icon="i-mdi-lock" :padded="false" />
</template>
