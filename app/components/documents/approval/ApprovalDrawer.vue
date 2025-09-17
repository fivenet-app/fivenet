<script lang="ts" setup>
import type { ApprovalPanelSnapshot } from '~~/gen/ts/services/documents/approval';
import ApprovalDrawerContents from './ApprovalDrawerContents.vue';

const props = defineProps<{
    documentId: number;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { activeChar } = useAuth();

const panel = ref<ApprovalPanelSnapshot>({
    allRequiredStagesSatisfied: false,
    anyDeclined: false,
    currentOrder: 0,
    documentId: props.documentId,
    pendingTasks: [],
    stages: [],
});

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
                        <p class="mb-2 text-lg font-medium">Approvals required</p>

                        <p class="text-md">1/2 Approvals</p>
                    </div>

                    <div class="basis-3/4">
                        <ApprovalDrawerContents :document-id="documentId" :panel="panel" :me="activeChar!.userId" />
                    </div>
                </div>
            </div>
        </template>

        <template #footer>
            <div class="mx-auto flex w-full max-w-[80%] min-w-3/4 flex-1 flex-col">
                <UButtonGroup class="w-full flex-1">
                    <UButton color="success" icon="i-mdi-check-bold" block size="xl" :label="$t('common.approve')" />
                    <UButton color="red" icon="i-mdi-close-bold" block size="xl" :label="$t('common.decline')" />
                </UButtonGroup>
            </div>
        </template>
    </UDrawer>
</template>
