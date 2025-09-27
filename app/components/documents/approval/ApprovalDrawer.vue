<script lang="ts" setup>
import { getDocumentsApprovalClient } from '~~/gen/ts/clients';
import type { ApprovalTask } from '~~/gen/ts/resources/documents/approval';

const props = defineProps<{
    documentId: number;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const approvalClient = await getDocumentsApprovalClient();

const { data } = useLazyAsyncData(`approval-drawer-${props.documentId}`, () => listTasks());

async function listTasks(): Promise<ApprovalTask[]> {
    if (!props.documentId) return [];

    const call = approvalClient.listTasks({
        pagination: { offset: 0 },
        documentId: props.documentId,
        statuses: [],
    });
    const { response } = await call;

    return response.tasks;
}

// TODO
</script>

<template>
    <UDrawer
        :title="$t('common.approve')"
        :overlay="false"
        handle-only
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{ content: 'min-h-[50%]', title: 'flex flex-row gap-2' }"
    >
        <template #title>
            <span class="flex-1">{{ $t('common.approve') }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="mx-auto w-full max-w-[80%] min-w-3/4">
                <div class="flex flex-1 flex-col sm:flex-row sm:gap-4">
                    <div class="basis-1/4">
                        <div class="mb-2 inline-flex items-center justify-between gap-2">
                            <p class="shrink-0 text-lg font-medium">1/2 Approvals</p>

                            <UButton icon="i-mdi-refresh" size="sm" :label="$t('common.refresh')" />
                        </div>

                        <p class="text-md">Approvals required</p>
                    </div>

                    <div class="basis-3/4">
                        <h3 class="mb-2 text-sm font-semibold">Your pending approvals</h3>

                        <UAlert color="gray" variant="soft" title="No pending approvals" class="mt-2" />

                        {{ data }}
                    </div>
                </div>
            </div>
        </template>

        <template #footer>
            <div class="mx-auto flex w-full max-w-[80%] min-w-3/4 flex-1 flex-col">
                <UButtonGroup class="w-full flex-1">
                    <UButton color="success" icon="i-mdi-check-bold" block size="lg" :label="$t('common.approve')" />
                    <UButton color="red" icon="i-mdi-close-bold" block size="lg" :label="$t('common.decline')" />
                </UButtonGroup>
            </div>
        </template>
    </UDrawer>
</template>
