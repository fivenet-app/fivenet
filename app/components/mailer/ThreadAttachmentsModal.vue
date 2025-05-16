<script lang="ts" setup>
import type { MessageAttachment } from '~~/gen/ts/resources/mailer/message';
import ThreadAttachmentsForm from './ThreadAttachmentsForm.vue';

defineProps<{
    attachments: MessageAttachment[];
    canSubmit: boolean;
}>();

defineEmits<{
    (e: 'update:attachments', attachments: MessageAttachment[]): void;
}>();

const { isOpen } = useModal();
</script>

<template>
    <UModal>
        <UCard
            :ui="{
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                base: 'flex flex-1 flex-col',
                body: { base: 'flex flex-1 flex-col' },
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.attachment', 2) }}
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <ThreadAttachmentsForm
                    :model-value="attachments"
                    :can-submit="canSubmit"
                    @update:model-value="$emit('update:attachments', $event)"
                />
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton class="flex-1" block color="black" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
