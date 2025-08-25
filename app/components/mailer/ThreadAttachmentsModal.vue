<script lang="ts" setup>
import type { MessageAttachment } from '~~/gen/ts/resources/mailer/message';
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
    <UModal>
        <template #title>
            <h3 class="text-2xl leading-6 font-semibold">
                {{ $t('common.attachment', 2) }}
            </h3>
        </template>

        <template #body>
            <ThreadAttachmentsForm
                :model-value="attachments"
                :can-submit="canSubmit"
                @update:model-value="$emit('update:attachments', $event)"
            />
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" @click="$emit('close', false)">
                    {{ $t('common.close', 1) }}
                </UButton>
            </UButtonGroup>
        </template>
    </UModal>
</template>
