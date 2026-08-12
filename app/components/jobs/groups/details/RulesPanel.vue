<script lang="ts" setup>
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import ConfirmModalWithReason from '~/components/partials/ConfirmModalWithReason.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCompletorStore } from '~/stores/completor';
import { getJobsGroupsClient, getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { Access } from '~~/gen/ts/resources/access/access';
import { AccessLevel as GroupAccessLevel } from '~~/gen/ts/resources/jobs/groups/access/access';
import {
    GroupGradeRuleType,
    type GroupQualificationRule,
    GroupQualificationRuleType,
    type GroupRule,
    GroupRuleType,
} from '~~/gen/ts/resources/jobs/groups/group';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import { QualificationExamMode } from '~~/gen/ts/resources/qualifications/exam/exam';
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { GroupRuleInput, ListGroupRulesResponse } from '~~/gen/ts/services/jobs/groups';
import { checkGroupAccess, groupRuleLabel } from '../helpers';

const props = defineProps<{
    groupId: number;
    canView: boolean;
    canManage: boolean;
    access?: Access;
}>();

const emit = defineEmits<{
    changed: [];
}>();

const { t } = useI18n();
const { activeChar } = useAuth();
const overlay = useOverlay();
const completorStore = useCompletorStore();

const jobsGroupsClient = await getJobsGroupsClient();
const qualificationsQualificationsClient = await getQualificationsQualificationsClient();
const confirmModalWithReason = overlay.create(ConfirmModalWithReason);

const page = ref(1);
const selectedQualifications = ref<QualificationShort[]>([]);
const jobs = ref<Job[]>([]);
const editingRuleId = ref<number>();
const pendingAction = ref<string>();

const ruleForm = reactive({
    enabled: true,
    type: GroupRuleType.GRADE,
    gradeType: GroupGradeRuleType.MINIMUM,
    grade: 0,
    minGrade: 0,
    maxGrade: 0,
    qualificationType: GroupQualificationRuleType.ALL,
    requireCompleted: true,
    reason: '',
});

const rulesKey = computed(() => `jobs-group-rules-${props.groupId}-${page.value}`);

const {
    data: rulesData,
    status: rulesStatus,
    error: rulesError,
    refresh: refreshRules,
} = useLazyAsyncData(rulesKey, () => listGroupRules(), {
    watch: [() => props.groupId, page],
});

const rules = computed(() => rulesData.value?.rules ?? []);
const isMutating = computed(() => pendingAction.value !== undefined);
const canManageRules = computed(() => props.canManage && checkGroupAccess(props.access, GroupAccessLevel.EDIT));
const activeJob = computed(() => jobs.value.find((job) => job.name === activeChar.value?.job));
const activeJobGrades = computed(() => activeJob.value?.grades ?? []);
const rangeMinGradeItems = computed(() =>
    activeJobGrades.value.filter((grade) => ruleForm.maxGrade === 0 || grade.grade <= ruleForm.maxGrade),
);
const rangeMaxGradeItems = computed(() => activeJobGrades.value.filter((grade) => grade.grade >= ruleForm.minGrade));
const ruleGradeValid = computed(() => {
    if (ruleForm.type !== GroupRuleType.GRADE) return true;
    if (activeJobGrades.value.length === 0) return false;

    if (ruleForm.gradeType === GroupGradeRuleType.RANGE) {
        return (
            activeJobGrades.value.some((grade) => grade.grade === ruleForm.minGrade) &&
            activeJobGrades.value.some((grade) => grade.grade === ruleForm.maxGrade) &&
            ruleForm.minGrade <= ruleForm.maxGrade
        );
    }

    return activeJobGrades.value.some((grade) => grade.grade === ruleForm.grade);
});

const ruleTypeItems = computed(() => [
    { label: t('enums.jobs.groups.GroupRuleType.GRADE'), value: GroupRuleType.GRADE },
    { label: t('enums.jobs.groups.GroupRuleType.QUALIFICATION'), value: GroupRuleType.QUALIFICATION },
]);

const gradeRuleTypeItems = computed(() => [
    { label: t('enums.jobs.groups.GroupGradeRuleType.MINIMUM'), value: GroupGradeRuleType.MINIMUM },
    { label: t('enums.jobs.groups.GroupGradeRuleType.EXACT'), value: GroupGradeRuleType.EXACT },
    { label: t('enums.jobs.groups.GroupGradeRuleType.RANGE'), value: GroupGradeRuleType.RANGE },
]);

const qualificationRuleTypeItems = computed(() => [
    { label: t('enums.jobs.groups.GroupQualificationRuleType.ALL'), value: GroupQualificationRuleType.ALL },
    { label: t('enums.jobs.groups.GroupQualificationRuleType.ANY'), value: GroupQualificationRuleType.ANY },
]);

async function listGroupRules(): Promise<ListGroupRulesResponse> {
    const { response } = await jobsGroupsClient.listGroupRules({
        groupId: props.groupId,
        pagination: {
            offset: calculateOffset(page.value, rulesData.value?.pagination),
        },
    });

    return response;
}

function gradeLabel(grade: number): string {
    const jobGrade = activeJobGrades.value.find((jobGrade) => jobGrade.grade === grade);
    if (!jobGrade) return `${grade}`;

    return `${jobGrade.label} (${jobGrade.grade})`;
}

function firstJobGrade(): number {
    return activeJobGrades.value[0]?.grade ?? 0;
}

function lastJobGrade(): number {
    return activeJobGrades.value.at(-1)?.grade ?? firstJobGrade();
}

function normalizeJobGrade(value: number, fallback: number): number {
    if (activeJobGrades.value.length === 0) return value;
    if (activeJobGrades.value.some((grade) => grade.grade === value)) return value;

    return fallback;
}

function normalizeRuleGradeForm(): void {
    if (ruleForm.type !== GroupRuleType.GRADE || activeJobGrades.value.length === 0) return;

    if (ruleForm.gradeType === GroupGradeRuleType.RANGE) {
        ruleForm.minGrade = normalizeJobGrade(ruleForm.minGrade, firstJobGrade());
        ruleForm.maxGrade = normalizeJobGrade(ruleForm.maxGrade, lastJobGrade());
        if (ruleForm.minGrade > ruleForm.maxGrade) {
            ruleForm.maxGrade = ruleForm.minGrade;
        }
        return;
    }

    ruleForm.grade = normalizeJobGrade(ruleForm.grade, firstJobGrade());
}

async function ensureRuleJobs(): Promise<void> {
    if (jobs.value.length > 0) return;

    jobs.value = await completorStore.listJobs();
}

function resetRuleForm(): void {
    editingRuleId.value = undefined;
    ruleForm.enabled = true;
    ruleForm.type = GroupRuleType.GRADE;
    ruleForm.gradeType = GroupGradeRuleType.MINIMUM;
    ruleForm.grade = 0;
    ruleForm.minGrade = 0;
    ruleForm.maxGrade = 0;
    ruleForm.qualificationType = GroupQualificationRuleType.ALL;
    ruleForm.requireCompleted = true;
    ruleForm.reason = '';
    selectedQualifications.value = [];
    normalizeRuleGradeForm();
}

function qualificationPlaceholder(id: number): QualificationShort {
    return {
        id,
        job: '',
        weight: 0,
        closed: false,
        draft: false,
        public: false,
        abbreviation: `#${id}`,
        title: t('components.jobs.groups.details.qualification_id', { id }),
        description: '',
        creatorJob: '',
        requirements: [],
        examMode: QualificationExamMode.DISABLED,
    };
}

async function loadQualificationShorts(ids: number[]): Promise<QualificationShort[]> {
    return await Promise.all(
        ids.map(async (id) => {
            try {
                const { response } = await qualificationsQualificationsClient.getQualification({
                    qualificationId: id,
                    withExam: false,
                });
                return (response.qualification ?? qualificationPlaceholder(id)) as QualificationShort;
            } catch {
                return qualificationPlaceholder(id);
            }
        }),
    );
}

async function editRule(rule: GroupRule): Promise<void> {
    editingRuleId.value = rule.id;
    ruleForm.enabled = rule.enabled;
    ruleForm.reason = '';

    if (rule.rule.oneofKind === 'grade') {
        await ensureRuleJobs();
        ruleForm.type = GroupRuleType.GRADE;
        ruleForm.gradeType = rule.rule.grade.type || GroupGradeRuleType.MINIMUM;
        ruleForm.grade = rule.rule.grade.grade ?? 0;
        ruleForm.minGrade = rule.rule.grade.minGrade ?? 0;
        ruleForm.maxGrade = rule.rule.grade.maxGrade ?? 0;
        selectedQualifications.value = [];
        normalizeRuleGradeForm();
        return;
    }

    if (rule.rule.oneofKind === 'qualification') {
        const qualification = rule.rule.qualification;
        ruleForm.type = GroupRuleType.QUALIFICATION;
        ruleForm.qualificationType = qualification.type || GroupQualificationRuleType.ALL;
        ruleForm.requireCompleted = qualification.requireCompleted;
        selectedQualifications.value = qualification.qualificationIds.map(qualificationPlaceholder);
        selectedQualifications.value = await loadQualificationShorts(qualification.qualificationIds);
    }
}

function buildRuleInput(): GroupRuleInput | undefined {
    if (ruleForm.type === GroupRuleType.GRADE) {
        if (!ruleGradeValid.value) return undefined;

        return {
            enabled: ruleForm.enabled,
            rule: {
                oneofKind: 'grade',
                grade: {
                    type: ruleForm.gradeType,
                    grade: ruleForm.gradeType === GroupGradeRuleType.RANGE ? undefined : ruleForm.grade,
                    minGrade: ruleForm.gradeType === GroupGradeRuleType.RANGE ? ruleForm.minGrade : undefined,
                    maxGrade: ruleForm.gradeType === GroupGradeRuleType.RANGE ? ruleForm.maxGrade : undefined,
                },
            },
        };
    }

    const qualificationIds = selectedQualifications.value.map((qualification) => qualification.id).filter((id) => id > 0);
    if (qualificationIds.length === 0) return undefined;

    const qualification: GroupQualificationRule = {
        type: ruleForm.qualificationType,
        qualificationIds,
        requireCompleted: ruleForm.requireCompleted,
    };

    return {
        enabled: ruleForm.enabled,
        rule: {
            oneofKind: 'qualification',
            qualification,
        },
    };
}

async function runMutation(action: string, mutate: () => Promise<void>): Promise<void> {
    pendingAction.value = action;
    try {
        await mutate();
        await refreshRules();
        emit('changed');
    } catch (e) {
        handleGRPCError(e as RpcError);
    } finally {
        pendingAction.value = undefined;
    }
}

async function saveRule(): Promise<void> {
    const input = buildRuleInput();
    if (!input) return;

    await runMutation('rule', async () => {
        const request = {
            groupId: props.groupId,
            rule: input,
            reason: ruleForm.reason.trim() || undefined,
        };

        if (editingRuleId.value) {
            await jobsGroupsClient.updateGroupRule({ ...request, ruleId: editingRuleId.value });
        } else {
            await jobsGroupsClient.createGroupRule(request);
        }

        resetRuleForm();
    });
}

async function deleteRule(rule: GroupRule): Promise<void> {
    confirmModalWithReason.open({
        title: t('common.delete'),
        confirm: async (reason: string) =>
            await runMutation(`rule-${rule.id}`, async () => {
                await jobsGroupsClient.deleteGroupRule({
                    groupId: props.groupId,
                    ruleId: rule.id,
                    reason,
                });
                if (editingRuleId.value === rule.id) resetRuleForm();
            }),
    });
}

watch(
    () => [ruleForm.type, ruleForm.gradeType, activeChar.value?.job] as const,
    async () => {
        if (ruleForm.type !== GroupRuleType.GRADE) return;

        await ensureRuleJobs();
        normalizeRuleGradeForm();
    },
    { immediate: true },
);

watch(
    () => ruleForm.minGrade,
    (minGrade) => {
        if (ruleForm.type !== GroupRuleType.GRADE || ruleForm.gradeType !== GroupGradeRuleType.RANGE) return;
        if (minGrade > ruleForm.maxGrade) {
            ruleForm.maxGrade = minGrade;
        }
    },
);

watch(
    () => ruleForm.maxGrade,
    (maxGrade) => {
        if (ruleForm.type !== GroupRuleType.GRADE || ruleForm.gradeType !== GroupGradeRuleType.RANGE) return;
        if (maxGrade < ruleForm.minGrade) {
            ruleForm.minGrade = maxGrade;
        }
    },
);

watch(
    () => props.groupId,
    () => {
        page.value = 1;
        resetRuleForm();
    },
);
</script>

<template>
    <div v-if="canView" class="grid gap-4">
        <UCard v-if="canManageRules" variant="subtle">
            <div class="grid gap-3">
                <div class="flex flex-1 gap-3 lg:flex-row lg:items-end">
                    <UFormField class="w-full" :label="$t('common.type')">
                        <USelectMenu
                            v-model="ruleForm.type"
                            class="w-full"
                            :items="ruleTypeItems"
                            value-key="value"
                            :disabled="isMutating || !canManageRules"
                        />
                    </UFormField>

                    <UFormField :label="$t('common.status')">
                        <USwitch
                            v-model="ruleForm.enabled"
                            :label="$t('common.enabled')"
                            :disabled="isMutating || !canManageRules"
                        />
                    </UFormField>

                    <UFieldGroup class="inline-flex w-full sm:w-auto">
                        <UButton
                            :icon="editingRuleId ? 'i-mdi-content-save' : 'i-mdi-plus'"
                            :label="editingRuleId ? $t('common.save', 1) : $t('common.add')"
                            :loading="pendingAction === 'rule'"
                            :disabled="
                                isMutating ||
                                !canManageRules ||
                                (ruleForm.type === GroupRuleType.GRADE && !ruleGradeValid) ||
                                (ruleForm.type === GroupRuleType.QUALIFICATION && selectedQualifications.length === 0)
                            "
                            @click="saveRule"
                        />
                        <UButton
                            color="neutral"
                            variant="outline"
                            icon="i-mdi-close"
                            :label="$t('common.cancel')"
                            :disabled="isMutating || !canManageRules"
                            @click="resetRuleForm"
                        />
                    </UFieldGroup>
                </div>

                <div v-if="ruleForm.type === GroupRuleType.GRADE" class="grid gap-3 sm:grid-cols-3">
                    <UFormField :label="$t('components.jobs.groups.details.rule_type')">
                        <USelectMenu
                            v-model="ruleForm.gradeType"
                            class="w-full"
                            :items="gradeRuleTypeItems"
                            value-key="value"
                            :disabled="isMutating || !canManage"
                        />
                    </UFormField>

                    <UFormField v-if="ruleForm.gradeType !== GroupGradeRuleType.RANGE" :label="$t('common.rank')">
                        <USelectMenu
                            v-model="ruleForm.grade"
                            class="w-full"
                            :items="activeJobGrades"
                            value-key="grade"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :disabled="isMutating || !canManageRules || activeJobGrades.length === 0"
                        >
                            <template v-if="activeJobGrades.length > 0" #default>
                                {{ gradeLabel(ruleForm.grade) }}
                            </template>

                            <template #item-label="{ item }"> {{ item.label }} ({{ item.grade }}) </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.job_grade')]) }}
                            </template>
                        </USelectMenu>
                    </UFormField>

                    <template v-else>
                        <UFormField :label="$t('common.min')">
                            <USelectMenu
                                v-model="ruleForm.minGrade"
                                class="w-full"
                                :items="rangeMinGradeItems"
                                value-key="grade"
                                :search-input="{ placeholder: $t('common.search_field') }"
                                :disabled="isMutating || !canManageRules || activeJobGrades.length === 0"
                            >
                                <template v-if="activeJobGrades.length > 0" #default>
                                    {{ gradeLabel(ruleForm.minGrade) }}
                                </template>

                                <template #item-label="{ item }"> {{ item.label }} ({{ item.grade }}) </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.job_grade')]) }}
                                </template>
                            </USelectMenu>
                        </UFormField>
                        <UFormField :label="$t('common.max')">
                            <USelectMenu
                                v-model="ruleForm.maxGrade"
                                class="w-full"
                                :items="rangeMaxGradeItems"
                                value-key="grade"
                                :search-input="{ placeholder: $t('common.search_field') }"
                                :disabled="isMutating || !canManageRules || activeJobGrades.length === 0"
                            >
                                <template v-if="activeJobGrades.length > 0" #default>
                                    {{ gradeLabel(ruleForm.maxGrade) }}
                                </template>

                                <template #item-label="{ item }"> {{ item.label }} ({{ item.grade }}) </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.job_grade')]) }}
                                </template>
                            </USelectMenu>
                        </UFormField>
                    </template>
                </div>

                <div v-else class="grid gap-3">
                    <div class="grid gap-3 sm:grid-cols-2">
                        <UFormField :label="$t('components.jobs.groups.details.rule_type')">
                            <USelectMenu
                                v-model="ruleForm.qualificationType"
                                class="w-full"
                                :items="qualificationRuleTypeItems"
                                value-key="value"
                                :disabled="isMutating || !canManageRules"
                            />
                        </UFormField>

                        <UFormField :label="$t('common.status')">
                            <UCheckbox
                                v-model="ruleForm.requireCompleted"
                                :label="$t('components.jobs.groups.details.require_completed')"
                                :disabled="isMutating || !canManageRules"
                            />
                        </UFormField>
                    </div>

                    <UFormField :label="$t('common.qualification', 2)">
                        <SelectMenu
                            v-model="selectedQualifications"
                            class="w-full"
                            multiple
                            :searchable="
                                async (q: string) => {
                                    const { response } = await qualificationsQualificationsClient.listQualifications({
                                        pagination: { offset: 0 },
                                        search: q,
                                    });
                                    return (response?.qualifications ?? []) as QualificationShort[];
                                }
                            "
                            searchable-key="jobs-group-rule-qualifications"
                            :filter-fields="['abbreviation', 'title']"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :placeholder="$t('common.qualification', 2)"
                            :disabled="isMutating || !canManageRules"
                        >
                            <template #item-label="{ item }">
                                {{ `${item?.abbreviation}: ${item?.title}` }}
                            </template>
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.qualification', 2)]) }}
                            </template>
                        </SelectMenu>
                    </UFormField>
                </div>

                <UFormField :label="$t('common.reason', 1)">
                    <UTextarea
                        v-model="ruleForm.reason"
                        class="w-full"
                        :rows="2"
                        :placeholder="$t('common.reason', 1)"
                        :disabled="isMutating || !canManageRules"
                    />
                </UFormField>
            </div>
        </UCard>

        <DataPendingBlock
            v-if="isRequestPending(rulesStatus)"
            :message="$t('common.loading', [$t('components.jobs.groups.rules')])"
        />
        <DataErrorBlock
            v-else-if="rulesError"
            :title="$t('common.unable_to_load', [$t('components.jobs.groups.rules')])"
            :error="rulesError"
            :retry="refreshRules"
        />
        <DataNoDataBlock
            v-else-if="rulesData?.pagination?.totalCount === 0"
            :type="$t('components.jobs.groups.rules')"
            icon="i-mdi-filter-cog"
            :padded="false"
        />
        <template v-else>
            <UCard v-for="rule in rules" :key="rule.id" variant="subtle">
                <div class="flex flex-row items-start justify-between gap-3">
                    <div class="flex flex-1 flex-col gap-1">
                        <p class="font-medium">
                            <span class="font-bold">#{{ rule.id }}</span> - {{ groupRuleLabel(rule, t) }}
                        </p>

                        <p class="mt-1 text-sm text-muted">
                            {{ $t('common.created_by') }}
                            <ColleagueInfoPopover :user="rule.createdBy" :user-id="rule.createdByUserId" hide-props />
                        </p>
                    </div>

                    <div class="flex shrink-0 flex-col gap-1">
                        <div class="place-self-end">
                            <UBadge :color="rule.enabled ? 'success' : 'neutral'" variant="subtle">
                                {{ rule.enabled ? $t('common.enabled') : $t('common.disabled') }}
                            </UBadge>
                        </div>

                        <GenericTime v-if="rule.createdAt" class="text-sm text-muted" :value="rule.createdAt" />
                    </div>
                </div>

                <div v-if="canManageRules" class="mt-3 flex justify-end gap-2">
                    <UButton
                        color="neutral"
                        variant="outline"
                        icon="i-mdi-pencil"
                        :label="$t('common.edit')"
                        :disabled="isMutating"
                        @click="editRule(rule)"
                    />
                    <UButton
                        color="error"
                        variant="outline"
                        icon="i-mdi-delete"
                        :label="$t('common.delete')"
                        :loading="pendingAction === `rule-${rule.id}`"
                        :disabled="isMutating"
                        @click="deleteRule(rule)"
                    />
                </div>
            </UCard>

            <Pagination v-model="page" :pagination="rulesData?.pagination" :status="rulesStatus" :refresh="refreshRules" />
        </template>
    </div>
    <DataNoDataBlock v-else :message="$t('common.no_access')" icon="i-mdi-lock" :padded="false" />
</template>
