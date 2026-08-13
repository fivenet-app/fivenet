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
import { type GroupType, GroupExclusionReason, type GroupMemberExclusion } from '~~/gen/ts/resources/jobs/groups/group';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import type { ListGroupMemberExclusionsResponse } from '~~/gen/ts/services/jobs/groups';
import { checkGroupAccess } from '../helpers';
import { groupTypeAllowsExclusions } from '../policy';

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
const selectedExclusionMember = ref<UserShort>();
const reasonType = ref(GroupExclusionReason.MANUAL);
const reason = ref('');
const editingExclusionMemberId = ref<number>();
const pendingAction = ref<string>();

const exclusionsKey = computed(() => `jobs-group-exclusions-${props.groupId}-${page.value}`);

const {
    data: exclusionsData,
    status: exclusionsStatus,
    error: exclusionsError,
    refresh: refreshExclusions,
} = useLazyAsyncData(exclusionsKey, () => listGroupMemberExclusions(), {
    watch: [() => props.groupId, page],
});

const exclusions = computed<GroupMemberExclusion[]>(() => exclusionsData.value?.exclusions ?? []);
const exclusionMemberIds = computed(() => new Set(exclusions.value.map((exclusion) => exclusion.userId)));
const isMutating = computed(() => pendingAction.value !== undefined);
const supportsExclusions = computed(() => groupTypeAllowsExclusions(props.groupType));
const canManageExclusions = computed(
    () => props.canManage && supportsExclusions.value && checkGroupAccess(props.access, GroupAccessLevel.EDIT),
);
const policyNotice = computed(() =>
    supportsExclusions.value ? undefined : 'components.jobs.groups.policy.exclusions_disabled',
);

const exclusionReasonItems = computed(() => [
    { label: t('enums.jobs.groups.GroupExclusionReason.MANUAL'), value: GroupExclusionReason.MANUAL },
    { label: t('enums.jobs.groups.GroupExclusionReason.TEMPORARY'), value: GroupExclusionReason.TEMPORARY },
    { label: t('enums.jobs.groups.GroupExclusionReason.NOT_ELIGIBLE'), value: GroupExclusionReason.NOT_ELIGIBLE },
    { label: t('enums.jobs.groups.GroupExclusionReason.OTHER'), value: GroupExclusionReason.OTHER },
]);

async function listGroupMemberExclusions(): Promise<ListGroupMemberExclusionsResponse> {
    const { response } = await jobsGroupsClient.listGroupMemberExclusions({
        groupId: props.groupId,
        pagination: {
            offset: calculateOffset(page.value, exclusionsData.value?.pagination),
        },
    });

    return response;
}

function resetExclusionForm(): void {
    editingExclusionMemberId.value = undefined;
    selectedExclusionMember.value = undefined;
    reasonType.value = GroupExclusionReason.MANUAL;
    reason.value = '';
}

function editExclusion(exclusion: GroupMemberExclusion): void {
    editingExclusionMemberId.value = exclusion.userId;
    selectedExclusionMember.value = {
        userId: exclusion.userId,
        job: exclusion.colleague?.job ?? '',
        jobGrade: exclusion.colleague?.jobGrade ?? 0,
        firstname: exclusion.colleague?.firstname ?? `${t('common.id')}: ${exclusion.userId}`,
        lastname: exclusion.colleague?.lastname ?? '',
        dateofbirth: exclusion.colleague?.dateofbirth ?? '',
    };
    reasonType.value = exclusion.reasonType || GroupExclusionReason.MANUAL;
    reason.value = exclusion.reason ?? '';
}

async function runMutation(action: string, mutate: () => Promise<void>): Promise<void> {
    pendingAction.value = action;
    try {
        await mutate();
        await refreshExclusions();
        emit('changed');
    } catch (e) {
        handleGRPCError(e as RpcError);
    } finally {
        pendingAction.value = undefined;
    }
}

async function addExclusion(): Promise<void> {
    if (!selectedExclusionMember.value?.userId) return;

    await runMutation('exclusion', async () => {
        await jobsGroupsClient.excludeGroupMember({
            groupId: props.groupId,
            userId: selectedExclusionMember.value!.userId,
            reasonType: reasonType.value,
            reason: reason.value.trim() || undefined,
        });
        resetExclusionForm();
    });
}

async function removeExclusion(userId: number): Promise<void> {
    confirmModal.open({
        title: t('common.remove'),
        confirm: async () =>
            await runMutation(`exclusion-${userId}`, async () => {
                await jobsGroupsClient.removeGroupMemberExclusion({
                    groupId: props.groupId,
                    userId,
                });
                if (editingExclusionMemberId.value === userId) resetExclusionForm();
            }),
    });
}

watch(
    () => props.groupId,
    () => {
        page.value = 1;
        resetExclusionForm();
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

        <UCard v-if="canManageExclusions" variant="subtle">
            <div class="grid gap-3">
                <div class="grid gap-3 lg:grid-cols-[minmax(0,1fr)_minmax(220px,320px)_auto] lg:items-end">
                    <UFormField :label="$t('common.colleague', 1)">
                        <SelectMenu
                            v-model="selectedExclusionMember"
                            class="w-full"
                            :searchable="
                                async (q: string) =>
                                    await completorStore.completeColleagues(
                                        q,
                                        selectedExclusionMember?.userId ? [selectedExclusionMember.userId] : [],
                                    )
                            "
                            searchable-key="jobs-group-exclusion-members"
                            :filter-fields="['firstname', 'lastname']"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :placeholder="$t('common.colleague', 1)"
                            :disabled="isMutating || !canManageExclusions"
                        >
                            <template v-if="selectedExclusionMember" #default>
                                {{ userToLabel(selectedExclusionMember) }}
                            </template>
                            <template #item-label="{ item }">
                                {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                            </template>
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.colleague', 2)]) }}
                            </template>
                        </SelectMenu>
                    </UFormField>

                    <UFormField :label="$t('common.reason', 1)">
                        <USelectMenu
                            v-model="reasonType"
                            class="w-full"
                            :items="exclusionReasonItems"
                            value-key="value"
                            :disabled="isMutating || !canManageExclusions"
                        />
                    </UFormField>

                    <UFieldGroup class="inline-flex w-full sm:w-auto">
                        <UButton
                            color="error"
                            variant="outline"
                            :icon="editingExclusionMemberId ? 'i-mdi-content-save' : 'i-mdi-account-cancel'"
                            :label="
                                editingExclusionMemberId
                                    ? $t('common.save', 1)
                                    : $t('components.jobs.groups.details.add_exclusion')
                            "
                            :loading="pendingAction === 'exclusion'"
                            :disabled="
                                isMutating ||
                                !canManageExclusions ||
                                !selectedExclusionMember?.userId ||
                                (exclusionMemberIds.has(selectedExclusionMember.userId) &&
                                    editingExclusionMemberId !== selectedExclusionMember.userId)
                            "
                            @click="addExclusion"
                        />
                        <UButton
                            v-if="editingExclusionMemberId"
                            color="neutral"
                            variant="outline"
                            icon="i-mdi-close"
                            :label="$t('common.cancel')"
                            :disabled="isMutating || !canManageExclusions"
                            @click="resetExclusionForm"
                        />
                    </UFieldGroup>
                </div>

                <UFormField :label="$t('common.description')">
                    <UTextarea
                        v-model="reason"
                        class="w-full"
                        :rows="2"
                        :placeholder="$t('common.reason', 1)"
                        :disabled="isMutating || !canManage"
                    />
                </UFormField>
            </div>
        </UCard>

        <DataPendingBlock
            v-if="isRequestPending(exclusionsStatus)"
            :message="$t('common.loading', [$t('components.jobs.groups.exclusions')])"
        />
        <DataErrorBlock
            v-else-if="exclusionsError"
            :title="$t('common.unable_to_load', [$t('components.jobs.groups.exclusions')])"
            :error="exclusionsError"
            :retry="refreshExclusions"
        />
        <DataNoDataBlock
            v-else-if="exclusionsData?.pagination?.totalCount === 0"
            :type="$t('components.jobs.groups.exclusions')"
            icon="i-mdi-account-cancel"
            :padded="false"
        />
        <template v-else>
            <ColleagueCard
                v-for="exclusion in exclusions"
                :key="exclusion.userId"
                compact
                :colleague="exclusion.colleague"
                :user-id="exclusion.userId"
            >
                <template #badges>
                    <UBadge color="error" variant="subtle">
                        {{
                            $t(
                                `enums.jobs.groups.GroupExclusionReason.${
                                    GroupExclusionReason[exclusion.reasonType] ?? 'UNSPECIFIED'
                                }`,
                            )
                        }}
                    </UBadge>
                </template>

                <p v-if="exclusion.reason">{{ exclusion.reason }}</p>
                <p>
                    {{ $t('common.created_by') }}
                    <ColleagueInfoPopover :user="exclusion.createdBy" :user-id="exclusion.createdByUserId" hide-props />
                </p>

                <template #footer>
                    <div class="flex flex-wrap items-center justify-end gap-2">
                        <div>
                            <GenericTime v-if="exclusion.createdAt" class="text-sm text-muted" :value="exclusion.createdAt" />
                        </div>
                        <UButton
                            v-if="canManageExclusions"
                            color="neutral"
                            variant="outline"
                            icon="i-mdi-pencil"
                            :label="$t('common.edit')"
                            :disabled="isMutating || !canManageExclusions"
                            @click="editExclusion(exclusion)"
                        />
                        <UButton
                            v-if="canManageExclusions"
                            color="error"
                            variant="outline"
                            icon="i-mdi-account-check"
                            :label="$t('common.remove')"
                            :loading="pendingAction === `exclusion-${exclusion.userId}`"
                            :disabled="isMutating || !canManageExclusions"
                            @click="removeExclusion(exclusion.userId)"
                        />
                    </div>
                </template>
            </ColleagueCard>

            <Pagination
                v-model="page"
                :pagination="exclusionsData?.pagination"
                :status="exclusionsStatus"
                :refresh="refreshExclusions"
            />
        </template>
    </div>
    <DataNoDataBlock v-else :message="$t('common.no_access')" icon="i-mdi-lock" :padded="false" />
</template>
