<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getDocumentsApprovalClient } from '~~/gen/ts/clients';
import {
    ApprovalRuleKind,
    ApprovalTaskStatus,
    OnEditBehavior,
    type ApprovalPolicy,
} from '~~/gen/ts/resources/documents/approval';
import type { DocumentMeta, DocumentShort } from '~~/gen/ts/resources/documents/documents';
import ApprovalList from './ApprovalList.vue';
import PolicyForm from './PolicyForm.vue';
import TaskDecideDrawer from './TaskDecideDrawer.vue';
import TaskForm from './TaskForm.vue';
import TaskList from './TaskList.vue';
import TaskStatusBadge from './TaskStatusBadge.vue';

const props = defineProps<{
    documentId: number;
    doc: DocumentShort;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const docMeta = defineModel<DocumentMeta | undefined>('docMeta');

const overlay = useOverlay();

const { can, activeChar } = useAuth();

const approvalClient = await getDocumentsApprovalClient();

const {
    data: policy,
    status,
    error,
    refresh,
} = useLazyAsyncData(`documents-approval-policy-${props.documentId}`, () => getPolicy());

async function getPolicy(): Promise<ApprovalPolicy | undefined> {
    const call = approvalClient.listApprovalPolicies({
        documentId: props.documentId,
    });
    const { response } = await call;

    if (!response.policy) return undefined;

    if (docMeta.value === undefined) {
        docMeta.value = {
            documentId: props.documentId,
            approved: false,
            closed: false,
            draft: false,
            public: false,
            state: '',
        };
    }

    const policy = response.policy;

    if (
        policy.approvedCount >= (response.policy.requiredCount ?? 1)
            ? ApprovalTaskStatus.APPROVED
            : policy.anyDeclined
              ? ApprovalTaskStatus.PENDING
              : ApprovalTaskStatus.DECLINED
    ) {
        docMeta.value.approved = true;
    } else if (
        !policy.anyDeclined && policy.approvedCount > policy.assignedCount
            ? ApprovalTaskStatus.APPROVED
            : ApprovalTaskStatus.DECLINED
    ) {
        docMeta.value.approved = true;
    }

    return policy;
}

async function recomputeApprovalPolicyCounters() {
    if (!policy.value) return;

    try {
        const call = approvalClient.recomputeApprovalPolicyCounters({
            documentId: props.documentId,
        });
        await call;

        refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
    }
}

const policyForm = overlay.create(PolicyForm, {
    props: {
        documentId: props.documentId,
    },
});

const taskFormDrawer = overlay.create(TaskForm);
</script>

<template>
    <UDrawer
        :title="$t('common.approve')"
        :overlay="false"
        handle-only
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{ container: 'flex-1', content: 'min-h-[50%]', title: 'flex flex-row gap-2 justify-between', body: 'h-full' }"
    >
        <template #title>
            <div class="inline-flex items-center gap-2">
                <span>{{ $t('common.approve') }}</span>

                <div v-if="policy?.ruleKind !== undefined && policy?.approvedCount !== undefined">
                    <TaskStatusBadge
                        v-if="policy.ruleKind === ApprovalRuleKind.REQUIRE_ALL"
                        :status="
                            !policy.anyDeclined && policy.approvedCount > policy.assignedCount
                                ? ApprovalTaskStatus.APPROVED
                                : ApprovalTaskStatus.DECLINED
                        "
                    />
                    <TaskStatusBadge
                        v-else
                        :status="
                            policy.approvedCount >= (policy.requiredCount ?? 1)
                                ? ApprovalTaskStatus.APPROVED
                                : policy.anyDeclined
                                  ? ApprovalTaskStatus.PENDING
                                  : ApprovalTaskStatus.DECLINED
                        "
                    />
                </div>
            </div>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="mx-auto w-full max-w-[80%] min-w-3/4">
                <div class="flex flex-1 flex-col sm:flex-row sm:gap-4">
                    <DataPendingBlock
                        v-if="isRequestPending(status)"
                        :message="$t('common.loading', [$t('common.approvals', 2)])"
                    />
                    <DataErrorBlock
                        v-else-if="error"
                        :title="$t('common.unable_to_load', [$t('common.approvals', 2)])"
                        :error="error"
                        :retry="refresh"
                    />

                    <template v-else>
                        <div class="basis-1/4 gap-4 sm:gap-0">
                            <UCard :ui="{ header: 'p-2 sm:px-2', body: 'p-2 sm:p-2', footer: 'p-2 sm:px-2' }">
                                <template v-if="policy" #header>
                                    <div class="flex flex-row flex-wrap gap-1">
                                        <UCollapsible class="flex flex-1 flex-col gap-1">
                                            <template #default>
                                                <div class="group flex flex-1 flex-row flex-wrap items-center gap-1">
                                                    <UTooltip :text="$t('common.expand_collapse')">
                                                        <UButton
                                                            class="order-last"
                                                            size="sm"
                                                            icon="i-mdi-chevron-double-down"
                                                            variant="link"
                                                            :ui="{
                                                                leadingIcon:
                                                                    'group-data-[state=open]:rotate-180 transition-transform duration-200',
                                                            }"
                                                        />
                                                    </UTooltip>

                                                    <p class="flex-1 shrink-0 text-lg font-medium">
                                                        <UTooltip :text="$t('common.approved')"
                                                            ><span class="text-success">{{
                                                                policy?.approvedCount
                                                            }}</span></UTooltip
                                                        >/<UTooltip :text="$t('common.declined')"
                                                            ><span class="text-error">{{
                                                                policy?.declinedCount
                                                            }}</span></UTooltip
                                                        >/<UTooltip :text="$t('enums.documents.ApprovalTaskStatus.PENDING')"
                                                            ><span class="text-info">{{ policy?.pendingCount }}</span></UTooltip
                                                        ><template v-if="policy?.ruleKind === ApprovalRuleKind.QUORUM_ANY"
                                                            >/<UTooltip :text="$t('common.required')"
                                                                ><span>{{ policy?.requiredCount }}</span></UTooltip
                                                            ></template
                                                        >
                                                        {{ $t('common.approvals', 2) }}
                                                    </p>
                                                </div>
                                            </template>

                                            <template #content>
                                                <div class="flex flex-col gap-1">
                                                    <p class="text-muted-foreground text-sm">
                                                        {{ $t('common.approved') }}:
                                                        <span class="text-success">{{ policy?.approvedCount }}</span>
                                                    </p>
                                                    <p class="text-muted-foreground text-sm">
                                                        {{ $t('common.declined') }}:
                                                        <span class="text-error">{{ policy?.declinedCount }}</span>
                                                    </p>
                                                    <p class="text-muted-foreground text-sm">
                                                        {{ $t('enums.documents.ApprovalTaskStatus.PENDING') }}:
                                                        <span class="text-info">{{ policy?.pendingCount }}</span>
                                                    </p>
                                                </div>
                                            </template>
                                        </UCollapsible>

                                        <div>
                                            <UTooltip :text="$t('common.refresh')">
                                                <UButton
                                                    icon="i-mdi-refresh"
                                                    size="sm"
                                                    variant="link"
                                                    :loading="!policy && isRequestPending(status)"
                                                    @click="() => refresh()"
                                                />
                                            </UTooltip>
                                        </div>
                                    </div>
                                </template>

                                <div v-if="policy" class="flex flex-col gap-2">
                                    <UBadge
                                        :label="
                                            $t(`enums.documents.ApprovalRuleKind.${ApprovalRuleKind[policy?.ruleKind ?? 0]}`)
                                        "
                                        color="neutral"
                                        variant="outline"
                                    />

                                    <UBadge
                                        v-if="policy.ruleKind !== ApprovalRuleKind.REQUIRE_ALL"
                                        :label="`${$t('common.required')}: ${(policy?.requiredCount ?? 0) > 0 ? policy?.requiredCount : $t('common.all')} ${$t('common.approvals', (policy?.requiredCount ?? 0) > 0 ? (policy?.requiredCount ?? 0) : 2)}`"
                                        color="neutral"
                                        variant="outline"
                                    />

                                    <UBadge
                                        :label="
                                            $t(`enums.documents.OnEditBehavior.${OnEditBehavior[policy?.onEditBehavior ?? 0]}`)
                                        "
                                        color="info"
                                        variant="outline"
                                    />
                                </div>

                                <div v-else>
                                    <DataNoDataBlock icon="i-mdi-approval" :type="$t('common.policy')" :padded="false" />
                                </div>

                                <template
                                    v-if="
                                        can('documents.ApprovalService/UpsertApprovalPolicy').value ||
                                        can('documents.ApprovalService/RevokeApproval').value
                                    "
                                    #footer
                                >
                                    <UButtonGroup class="flex w-full flex-1">
                                        <UButton
                                            v-if="can('documents.ApprovalService/UpsertApprovalPolicy').value"
                                            block
                                            :label="$t('common.policy')"
                                            :trailing-icon="policy ? 'i-mdi-pencil' : 'i-mdi-plus'"
                                            @click="
                                                policyForm.open({
                                                    documentId: props.documentId,
                                                    modelValue: policy,
                                                    'onUpdate:modelValue': () => refresh(),
                                                })
                                            "
                                        />

                                        <UButton
                                            v-if="can('documents.ApprovalService/RevokeApproval').value"
                                            icon="i-mdi-calculator-variant"
                                            variant="outline"
                                            @click="() => recomputeApprovalPolicyCounters()"
                                        />
                                    </UButtonGroup>
                                </template>
                            </UCard>
                        </div>

                        <div class="flex flex-1 basis-3/4 flex-col gap-4 overflow-x-hidden p-0.5">
                            <ApprovalList :document-id="documentId" />

                            <TaskList :document-id="documentId">
                                <template
                                    v-if="can('documents.ApprovalService/UpsertApprovalTasks').value"
                                    #header="{ refresh: tasksRefresh }"
                                >
                                    <UButton
                                        :disabled="!policy"
                                        variant="link"
                                        :label="$t('common.create')"
                                        trailing-icon="i-mdi-task-add"
                                        @click="
                                            taskFormDrawer.open({
                                                documentId: documentId,
                                                policy: policy,
                                                onClose: (val) => val && tasksRefresh(),
                                            })
                                        "
                                    />
                                </template>
                            </TaskList>
                        </div>
                    </template>
                </div>
            </div>
        </template>

        <template #footer>
            <div class="mx-auto flex w-full max-w-[80%] min-w-3/4 flex-1 flex-col gap-4">
                <!-- RevokeApproval / ReopenApprovalTask perms are indicators for being able to do ad-hoc approval, otherwise a policy and a matching task is required -->
                <UButtonGroup v-if="!doc || doc.creatorId !== activeChar?.userId" class="w-full flex-1">
                    <TaskDecideDrawer
                        v-model:policy="policy"
                        :document-id="documentId"
                        :approve="true"
                        @close="(val) => val && refresh()"
                    >
                        <UButton color="success" icon="i-mdi-check-bold" block size="lg" :label="$t('common.approve')" />
                    </TaskDecideDrawer>

                    <TaskDecideDrawer
                        v-model:policy="policy"
                        :document-id="documentId"
                        :approve="false"
                        @close="(val) => val && refresh()"
                    >
                        <UButton color="red" icon="i-mdi-close-bold" block size="lg" :label="$t('common.decline')" />
                    </TaskDecideDrawer>
                </UButtonGroup>

                <UButtonGroup class="w-full flex-1">
                    <UButton
                        class="flex-1"
                        color="neutral"
                        block
                        :label="$t('common.close', 1)"
                        @click="$emit('close', false)"
                    />
                </UButtonGroup>
            </div>
        </template>
    </UDrawer>
</template>
