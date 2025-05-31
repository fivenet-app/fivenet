<script lang="ts" setup>
import type { TabItem } from '#ui/types';
import ClipboardCitizens from '~/components/clipboard/modal/ClipboardCitizens.vue';
import ClipboardDocuments from '~/components/clipboard/modal/ClipboardDocuments.vue';
import ClipboardVehicles from '~/components/clipboard/modal/ClipboardVehicles.vue';

const { t } = useI18n();

const { isOpen } = useModal();

const clipboardStore = useClipboardStore();

const items: TabItem[] = [
    {
        key: 'citizens',
        slot: 'citizens',
        label: t('common.citizen', 2),
        icon: 'i-mdi-account-multiple',
    },
    {
        key: 'vehicles',
        slot: 'vehicles',
        label: t('common.vehicle', 2),
        icon: 'i-mdi-car',
    },
    {
        key: 'documents',
        slot: 'documents',
        label: t('common.document', 2),
        icon: 'i-mdi-file-document-multiple',
    },
];

const selectedTab = ref(0);
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

            <UTabs v-model="selectedTab" :items="items" :unmount="true">
                <template #citizens>
                    <ClipboardCitizens hide-header @close="isOpen = false" />
                </template>
                <template #vehicles>
                    <ClipboardVehicles hide-header @close="isOpen = false" />
                </template>
                <template #documents>
                    <ClipboardDocuments hide-header @close="isOpen = false" />
                </template>
            </UTabs>

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
