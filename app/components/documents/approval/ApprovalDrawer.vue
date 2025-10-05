<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getDocumentsApprovalClient } from '~~/gen/ts/clients';
import { ApprovalRuleKind, OnEditBehavior, type ApprovalPolicy } from '~~/gen/ts/resources/documents/approval';
import PolicyForm from './PolicyForm.vue';
import TaskDecideDrawer from './TaskDecideDrawer.vue';
import TaskForm from './TaskForm.vue';
import TaskList from './TaskList.vue';

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
        :ui="{ container: 'flex-1', content: 'min-h-[50%]', title: 'flex flex-row gap-2', body: 'h-full' }"
    >
        <template #title>
            <span class="flex-1">{{ $t('common.approve') }}</span>
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
                        <div class="basis-1/4">
                            <UCard :ui="{ body: 'p-2 sm:p-2', footer: 'p-2 sm:px-2' }">
                                <div v-if="data" class="flex flex-col gap-2">
                                    <div class="inline-flex items-center justify-between gap-2">
                                        <p class="shrink-0 text-lg font-medium">
                                            {{ data?.approvedCount }}/{{ data?.requiredCount }} {{ $t('common.approvals') }}
                                        </p>

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

                                    <UBadge
                                        :label="$t(`enums.documents.ApprovalRuleKind.${ApprovalRuleKind[data?.ruleKind ?? 0]}`)"
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
                                    <div class="flex flex-1">
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
                                    </div>
                                </template>
                            </UCard>
                        </div>

                        <div class="basis-3/4">
                            <TaskList :document-id="documentId" :policy-id="data?.id ?? 0">
                                <template #header>
                                    <UButton
                                        :disabled="!data"
                                        variant="link"
                                        :label="$t('common.create')"
                                        trailing-icon="i-mdi-task-add"
                                        @click="taskFormDrawer.open({ policyId: data?.id ?? 0 })"
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
                <UButtonGroup class="w-full flex-1">
                    <TaskDecideDrawer :document-id="documentId" :policy-id="data?.id ?? 0" :approve="true">
                        <UButton color="success" icon="i-mdi-check-bold" block size="lg" :label="$t('common.approve')" />
                    </TaskDecideDrawer>

                    <TaskDecideDrawer :document-id="documentId" :policy-id="data?.id ?? 0" :approve="false">
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
