<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getDocumentsApprovalClient } from '~~/gen/ts/clients';
import { ApprovalTaskStatus } from '~~/gen/ts/resources/documents/approval';
import type { ListApprovalTasksResponse } from '~~/gen/ts/services/documents/approval';
import { approvalTaskStatusToColor } from './helpers';

const props = defineProps<{
    documentId: number;
}>();

const approvalClient = await getDocumentsApprovalClient();

const { data, status, error, refresh } = useLazyAsyncData(
    () => `approval-drawer-${props.documentId}-tasklist`,
    () => listApprovalTasks(),
);

async function listApprovalTasks(): Promise<ListApprovalTasksResponse> {
    const call = approvalClient.listApprovalTasks({
        documentId: props.documentId,
        statuses: [],
    });
    const { response } = await call;

    return response;
}
</script>

<template>
    <UCard :ui="{ body: 'p-0 sm:p-0', footer: 'p-0 sm:px-2' }">
        <template #header>
            <div class="flex items-center justify-between gap-1">
                <h3 class="flex-1 text-base leading-6 font-semibold">
                    {{ $t('components.documents.approval.tasks') }}
                </h3>

                <slot name="header" />

                <UButton variant="link" icon="i-mdi-refresh" @click="() => refresh()" />
            </div>
        </template>

        <div>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.approvals')])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.approvals')])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="data?.tasks.length === 0" icon="i-mdi-approval" :type="$t('common.task', 2)" />

            <ul v-else class="divide-y divide-default" role="list">
                <li
                    v-for="task in data?.tasks"
                    :key="task.id"
                    class="relative flex justify-between border-default p-2 hover:border-primary-500/25 hover:bg-primary-100/50 sm:px-4 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
                >
                    <div class="flex min-w-0 gap-x-2">
                        <div class="min-w-0 flex-auto">
                            <p class="text-sm leading-6 font-semibold text-toned">
                                <span class="absolute inset-x-0 -top-px bottom-0" />
                                <span class="inline-flex items-center gap-2 text-highlighted">
                                    <UBadge
                                        :label="$t(`enums.documents.ApprovalTaskStatus.${ApprovalTaskStatus[task.status]}`)"
                                        :color="approvalTaskStatusToColor(task.status)"
                                    />

                                    <CitizenInfoPopover v-if="task.userId" :user="task.user" :user-id="task.userId" />
                                    <span v-else> {{ task.jobLabel }} ({{ task.jobGradeLabel ?? task.minimumGrade }}) </span>
                                </span>
                            </p>
                        </div>
                    </div>
                </li>
            </ul>
        </div>
    </UCard>
</template>
