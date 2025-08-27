<script lang="ts" setup>
import type { TabsItem } from '@nuxt/ui';
import ClipboardCitizens from '~/components/clipboard/modal/ClipboardCitizens.vue';
import ClipboardDocuments from '~/components/clipboard/modal/ClipboardDocuments.vue';
import ClipboardVehicles from '~/components/clipboard/modal/ClipboardVehicles.vue';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { t } = useI18n();

const clipboardStore = useClipboardStore();

const items: TabsItem[] = [
    {
        slot: 'citizens' as const,
        label: t('common.citizen', 2),
        icon: 'i-mdi-account-multiple',
        value: 'citizens',
    },
    {
        slot: 'vehicles' as const,
        label: t('common.vehicle', 2),
        icon: 'i-mdi-car',
        value: 'vehicles',
    },
    {
        slot: 'documents' as const,
        label: t('common.document', 2),
        icon: 'i-mdi-file-document-multiple',
        value: 'documents',
    },
];

const selectedTab = ref('citizens');
</script>

<template>
    <UModal :title="$t('components.clipboard.clipboard_modal.title')">
        <template #body>
            <UTabs v-model="selectedTab" :items="items" :unmount="true">
                <template #citizens>
                    <ClipboardCitizens hide-header @close="$emit('close', false)" />
                </template>
                <template #vehicles>
                    <ClipboardVehicles hide-header @close="$emit('close', false)" />
                </template>
                <template #documents>
                    <ClipboardDocuments hide-header @close="$emit('close', false)" />
                </template>
            </UTabs>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block @click="$emit('close', false)">
                    {{ $t('common.close', 1) }}
                </UButton>

                <UButton class="flex-1" block color="error" @click="clipboardStore.clear()">
                    {{ $t('components.clipboard.clipboard_modal.clear') }}
                </UButton>
            </UButtonGroup>
        </template>
    </UModal>
</template>
