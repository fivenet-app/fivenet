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
import ApprovalList from './ApprovalList.vue';
import PolicyForm from './PolicyForm.vue';
import TaskDecideDrawer from './TaskDecideDrawer.vue';
import TaskForm from './TaskForm.vue';
import TaskList from './TaskList.vue';
import TaskStatusBadge from './TaskStatusBadge.vue';

const props = defineProps<{
    documentId: number;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const overlay = useOverlay();

const approvalClient = await getDocumentsApprovalClient();

const { data, status, error, refresh } = useLazyAsyncData(
    () => `approval-drawer-${props.documentId}`,
    () => getPolicy(),
);

async function getPolicy(): Promise<ApprovalPolicy | undefined> {
    const call = approvalClient.listApprovalPolicies({
        documentId: props.documentId,
    });
    const { response } = await call;

    return response.policy;
}

async function recomputeApprovalPolicyCounters() {
    if (!data.value) return;

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

                <div v-if="data?.ruleKind !== undefined && data?.approvedCount !== undefined">
                    <TaskStatusBadge
                        v-if="data.ruleKind === ApprovalRuleKind.REQUIRE_ALL"
                        :status="data.anyDeclined ? ApprovalTaskStatus.DECLINED : ApprovalTaskStatus.APPROVED"
                    />
                    <TaskStatusBadge
                        v-else
                        :status="
                            data.approvedCount >= (data.requiredCount ?? 1)
                                ? ApprovalTaskStatus.APPROVED
                                : data.anyDeclined
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
                        :message="$t('common.loading', [$t('common.approvals')])"
                    />
                    <DataErrorBlock
                        v-else-if="error"
                        :title="$t('common.unable_to_load', [$t('common.approvals')])"
                        :error="error"
                        :retry="refresh"
                    />

                    <template v-else>
                        <div class="basis-1/4 gap-4 sm:gap-0">
                            <UCard :ui="{ header: 'p-2 sm:px-2', body: 'p-2 sm:p-2', footer: 'p-2 sm:px-2' }">
                                <template v-if="data" #header>
                                    <div class="flex flex-row flex-wrap gap-1">
                                        <UCollapsible class="flex flex-1 flex-col gap-1">
                                            <template #default>
                                                <div class="group flex flex-1 flex-row flex-wrap items-center gap-1">
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

                                                    <p class="flex-1 shrink-0 text-lg font-medium">
                                                        <UTooltip :text="$t('common.approved')"
                                                            ><span class="text-success">{{
                                                                data?.approvedCount
                                                            }}</span></UTooltip
                                                        >/<UTooltip :text="$t('common.declined')"
                                                            ><span class="text-error">{{ data?.declinedCount }}</span></UTooltip
                                                        >/<UTooltip :text="$t('enums.documents.ApprovalTaskStatus.PENDING')"
                                                            ><span class="text-info">{{ data?.pendingCount }}</span></UTooltip
                                                        ><template v-if="data?.ruleKind === ApprovalRuleKind.QUORUM_ANY"
                                                            >/<UTooltip :text="$t('common.required')"
                                                                ><span>{{ data?.requiredCount }}</span></UTooltip
                                                            ></template
                                                        >
                                                        {{ $t('common.approvals') }}
                                                    </p>
                                                </div>
                                            </template>

                                            <template #content>
                                                <div class="flex flex-col gap-1">
                                                    <p class="text-muted-foreground text-sm">
                                                        {{ $t('common.approved') }}:
                                                        <span class="text-success">{{ data?.approvedCount }}</span>
                                                    </p>
                                                    <p class="text-muted-foreground text-sm">
                                                        {{ $t('common.declined') }}:
                                                        <span class="text-error">{{ data?.declinedCount }}</span>
                                                    </p>
                                                    <p class="text-muted-foreground text-sm">
                                                        {{ $t('enums.documents.ApprovalTaskStatus.PENDING') }}:
                                                        <span class="text-info">{{ data?.pendingCount }}</span>
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
                                                    :loading="!data && isRequestPending(status)"
                                                    @click="() => refresh()"
                                                />
                                            </UTooltip>
                                        </div>
                                    </div>
                                </template>

                                <div v-if="data" class="flex flex-col gap-2">
                                    <UBadge
                                        :label="$t(`enums.documents.ApprovalRuleKind.${ApprovalRuleKind[data?.ruleKind ?? 0]}`)"
                                        color="neutral"
                                        variant="outline"
                                    />

                                    <UBadge
                                        v-if="data.ruleKind !== ApprovalRuleKind.REQUIRE_ALL"
                                        :label="`${$t('common.required')}: ${(data?.requiredCount ?? 0) > 0 ? data?.requiredCount : $t('common.all')} ${$t('common.approvals', (data?.requiredCount ?? 0) > 0 ? (data?.requiredCount ?? 0) : 2)}`"
                                        color="neutral"
                                        variant="outline"
                                    />

                                    <UBadge
                                        :label="
                                            $t(`enums.documents.OnEditBehavior.${OnEditBehavior[data?.onEditBehavior ?? 0]}`)
                                        "
                                        color="info"
                                        variant="outline"
                                    />
                                </div>

                                <div v-else>
                                    <DataNoDataBlock icon="i-mdi-approval" :type="$t('common.policy')" />
                                </div>

                                <template #footer>
                                    <UButtonGroup class="flex w-full flex-1">
                                        <UButton
                                            block
                                            :label="$t('common.policy')"
                                            :trailing-icon="data ? 'i-mdi-pencil' : 'i-mdi-plus'"
                                            @click="
                                                policyForm.open({
                                                    documentId: props.documentId,
                                                    modelValue: data,
                                                    'onUpdate:modelValue': () => refresh(),
                                                })
                                            "
                                        />

                                        <UButton
                                            icon="i-mdi-calculator-variant"
                                            variant="outline"
                                            @click="() => recomputeApprovalPolicyCounters()"
                                        />
                                    </UButtonGroup>
                                </template>
                            </UCard>
                        </div>

                        <div class="flex flex-1 basis-3/4 flex-col gap-4">
                            <ApprovalList :document-id="documentId" :policy-id="data?.id ?? 0" />

                            <TaskList :document-id="documentId" :policy-id="data?.id ?? 0">
                                <template #header="{ refresh: tasksRefresh }">
                                    <UButton
                                        :disabled="!data"
                                        variant="link"
                                        :label="$t('common.create')"
                                        trailing-icon="i-mdi-task-add"
                                        @click="
                                            taskFormDrawer.open({
                                                policyId: data?.id ?? 0,
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
                <UButtonGroup class="w-full flex-1">
                    <TaskDecideDrawer
                        :document-id="documentId"
                        :policy-id="data?.id ?? 0"
                        :approve="true"
                        @close="(val) => val && refresh()"
                    >
                        <UButton color="success" icon="i-mdi-check-bold" block size="lg" :label="$t('common.approve')" />
                    </TaskDecideDrawer>

                    <TaskDecideDrawer
                        :document-id="documentId"
                        :policy-id="data?.id ?? 0"
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
