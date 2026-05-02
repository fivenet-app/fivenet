<script lang="ts" setup>
import type { MessageAttachment } from '~~/gen/ts/resources/mailer/messages/message';
import ThreadAttachmentsForm from './ThreadAttachmentsForm.vue';

defineProps<{
    attachments: MessageAttachment[];
    canSubmit: boolean;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:attachments', attachments: MessageAttachment[]): void;
}>();
</script>

<template>
    <UModal :title="$t('common.attachment', 2)">
        <template #body>
            <ThreadAttachmentsForm
                :model-value="attachments"
                :can-submit="canSubmit"
                @update:model-value="$emit('update:attachments', $event)"
            />
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UFieldGroup>
        </template>
    </UModal>
</template>
