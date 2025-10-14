<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getDocumentsApprovalClient } from '~~/gen/ts/clients';
import { ApprovalTaskStatus } from '~~/gen/ts/resources/documents/approval';
import type { DeleteApprovalTasksResponse, ListApprovalTasksResponse } from '~~/gen/ts/services/documents/approval';
import { approvalTaskStatusToColor } from './helpers';

const props = defineProps<{
    documentId: number;
}>();

const { can } = useAuth();

const approvalClient = await getDocumentsApprovalClient();

const { data, status, error, refresh } = useLazyAsyncData(
    () => `documents-approval-tasks-${props.documentId}`,
    () => listApprovalTasks(),
);

async function listApprovalTasks(): Promise<ListApprovalTasksResponse> {
    const call = approvalClient.listApprovalTasks({
        documentId: props.documentId,
        statuses: [ApprovalTaskStatus.PENDING, ApprovalTaskStatus.EXPIRED, ApprovalTaskStatus.CANCELLED],
    });
    const { response } = await call;

    return response;
}

async function removeTask(id: number): Promise<DeleteApprovalTasksResponse> {
    try {
        const call = approvalClient.deleteApprovalTasks({
            documentId: props.documentId,
            taskIds: [id],
            deleteAllPending: false,
        });
        const { response } = await call;

        await refresh();

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <UCard :ui="{ body: 'p-0 sm:p-0', footer: 'p-0 sm:px-2' }">
        <template #header>
            <div class="flex items-center justify-between gap-1">
                <h3 class="flex-1 text-base leading-6 font-semibold">
                    {{ $t('components.documents.approval.tasks') }}
                </h3>

                <slot name="header" :refresh="refresh" />

                <UTooltip :text="$t('common.refresh')">
                    <UButton variant="link" icon="i-mdi-refresh" @click="() => refresh()" />
                </UTooltip>
            </div>
        </template>

        <div>
            <DataPendingBlock
                v-if="isRequestPending(status)"
                :message="$t('common.loading', [$t('components.documents.approval.tasks')])"
            />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('components.documents.approval.tasks')])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="data?.tasks.length === 0"
                icon="i-mdi-approval"
                :type="$t('components.documents.approval.tasks')"
            />

            <ul v-else class="divide-y divide-default" role="list">
                <li
                    v-for="task in data?.tasks"
                    :key="task.id"
                    class="relative flex justify-between border-default p-2 hover:border-primary-500/25 hover:bg-primary-100/50 sm:px-4 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
                >
                    <div class="flex min-w-0 flex-1 gap-x-2">
                        <div class="min-w-0 flex-1 flex-auto">
                            <p class="text-sm leading-6 font-semibold text-toned">
                                <span class="absolute inset-x-0 -top-px bottom-0" />
                                <span class="inline-flex items-center gap-2 text-highlighted">
                                    <UBadge
                                        :label="$t(`enums.documents.ApprovalTaskStatus.${ApprovalTaskStatus[task.status]}`)"
                                        :color="approvalTaskStatusToColor(task.status)"
                                    />

                                    <CitizenInfoPopover v-if="task.userId" :user="task.user" :user-id="task.userId" />
                                    <p v-else>
                                        {{ task.jobLabel }}
                                        -
                                        {{ task.jobGradeLabel ?? task.minimumGrade }}
                                        ({{ $t('common.number') }}: {{ task.slotNo }})
                                    </p>
                                </span>
                            </p>
                            <p class="inline-flex gap-1 text-xs leading-6 font-semibold text-toned">
                                <span class="font-semibold">{{ $t('common.comment') }}:</span>
                                <span>
                                    {{ task.comment || $t('common.no_comment') }}
                                </span>
                            </p>
                        </div>

                        <div
                            v-if="can('documents.ApprovalService/DeleteApprovalTasks').value"
                            class="inline-flex items-center gap-2"
                        >
                            <UTooltip :text="$t('common.remove')">
                                <UButton icon="i-mdi-trash-can" color="red" variant="link" @click="removeTask(task.id)" />
                            </UTooltip>
                        </div>
                    </div>
                </li>
            </ul>
        </div>
    </UCard>
</template>
