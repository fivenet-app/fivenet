<script lang="ts" setup>
import type { TabsItem } from '@nuxt/ui';
import ClipboardCitizens from '~/components/clipboard/modal/ClipboardCitizens.vue';
import ClipboardDocuments from '~/components/clipboard/modal/ClipboardDocuments.vue';
import ClipboardVehicles from '~/components/clipboard/modal/ClipboardVehicles.vue';
import { CLIPBOARD_MAX_ITEMS } from '~/stores/clipboard';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { t } = useI18n();

const clipboardStore = useClipboardStore();
const { users, vehicles, documents } = storeToRefs(clipboardStore);

const items = computed<TabsItem[]>(() => [
    {
        slot: 'citizens' as const,
        label: t('common.citizen', 2),
        icon: 'i-mdi-account-multiple',
        value: 'citizens',
        badge: { color: 'neutral', variant: 'soft', size: 'sm', label: `${users.value.length}/${CLIPBOARD_MAX_ITEMS}` },
    },
    {
        slot: 'vehicles' as const,
        label: t('common.vehicle', 2),
        icon: 'i-mdi-car',
        value: 'vehicles',
        badge: { color: 'neutral', variant: 'soft', size: 'sm', label: `${vehicles.value.length}/${CLIPBOARD_MAX_ITEMS}` },
    },
    {
        slot: 'documents' as const,
        label: t('common.document', 2),
        icon: 'i-mdi-file-document-multiple',
        value: 'documents',
        badge: { color: 'neutral', variant: 'soft', size: 'sm', label: `${documents.value.length}/${CLIPBOARD_MAX_ITEMS}` },
    },
]);

const selectedTab = ref('citizens');
</script>

<template>
    <UModal :title="$t('components.clipboard.clipboard_modal.title')" :ui="{ body: 'min-h-90' }">
        <template #body>
            <UTabs v-model="selectedTab" :items="items" variant="pill">
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
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    color="error"
                    :label="$t('components.clipboard.clipboard_modal.clear')"
                    @click="clipboardStore.clear()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
