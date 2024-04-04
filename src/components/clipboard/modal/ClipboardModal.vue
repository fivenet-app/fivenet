<script lang="ts" setup>
import { useClipboardStore } from '~/store/clipboard';
import ClipboardCitizens from '~/components/clipboard/modal/ClipboardCitizens.vue';
import ClipboardDocuments from '~/components/clipboard/modal/ClipboardDocuments.vue';
import ClipboardVehicles from '~/components/clipboard/modal/ClipboardVehicles.vue';

const { isOpen } = useModal();

const clipboardStore = useClipboardStore();
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.clipboard.clipboard_modal.title') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <ClipboardCitizens />
                <ClipboardVehicles />
                <ClipboardDocuments />
            </div>

            <template #footer>
                <div class="isolate inline-flex w-full rounded-md pr-4 shadow-sm">
                    <UButton @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton color="red" @click="clipboardStore.clear()">
                        {{ $t('components.clipboard.clipboard_modal.clear') }}
                    </UButton>
                </div>
            </template>
        </UCard>
    </UModal>
</template>
