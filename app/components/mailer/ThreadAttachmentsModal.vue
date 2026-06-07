<script lang="ts" setup>
import type { MessageAttachment } from '~~/gen/ts/resources/mailer/messages/message';
import ThreadAttachmentsForm from './ThreadAttachmentsForm.vue';

const props = defineProps<{
    attachments: MessageAttachment[];
    canSubmit: boolean;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:attachments', attachments: MessageAttachment[]): void;
}>();

const attachmentsSnapshot = computed(() =>
    props.attachments.map((attachment) => ({
        type: attachment.data.oneofKind,
        documentId: attachment.data.oneofKind === 'document' ? attachment.data.document.id : null,
        documentTitle: attachment.data.oneofKind === 'document' ? attachment.data.document.title : null,
    })),
);

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(attachmentsSnapshot);

syncSnapshot();

async function closeModal(): Promise<void> {
    if (!props.canSubmit) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UModal :title="$t('common.attachment', 2)" :close="false" :dismissible="!hasUnsavedChanges && props.canSubmit">
        <template #header>
            <div class="flex w-full items-center justify-between gap-1.5">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('common.attachment', 2) }}
                </h3>

                <UButton color="neutral" variant="ghost" icon="i-mdi-close" :disabled="!props.canSubmit" @click="closeModal" />
            </div>
        </template>

        <template #body>
            <ThreadAttachmentsForm
                :model-value="attachments"
                :can-submit="canSubmit"
                @update:model-value="emit('update:attachments', $event)"
            />
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="closeModal" />
            </UFieldGroup>
        </template>
    </UModal>
</template>
