<script lang="ts" setup>
import ClipboardCitizens from '~/components/clipboard/modal/ClipboardCitizens.vue';
import ClipboardDocuments from '~/components/clipboard/modal/ClipboardDocuments.vue';
import ClipboardVehicles from '~/components/clipboard/modal/ClipboardVehicles.vue';
import { useClipboardStore } from '~/stores/clipboard';

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

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <div class="mx-auto">
                <ClipboardCitizens @close="isOpen = false" />

                <ClipboardVehicles @close="isOpen = false" />

                <ClipboardDocuments @close="isOpen = false" />
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton class="flex-1" color="black" block @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>

                    <UButton class="flex-1" block color="error" @click="clipboardStore.clear()">
                        {{ $t('components.clipboard.clipboard_modal.clear') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
