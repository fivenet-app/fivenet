<script lang="ts" setup>
import { createReusableTemplate, useMediaQuery } from '@vueuse/core';
import type { TabsItem } from '@nuxt/ui';
import ClipboardCitizens from '~/components/clipboard/modal/ClipboardCitizens.vue';
import ClipboardDocuments from '~/components/clipboard/modal/ClipboardDocuments.vue';
import ClipboardVehicles from '~/components/clipboard/modal/ClipboardVehicles.vue';
import { CLIPBOARD_MAX_ITEMS } from '~/stores/clipboard';

const props = withDefaults(
    defineProps<{
        guideInteractive?: boolean;
        open?: boolean;
    }>(),
    {
        guideInteractive: false,
        open: false,
    },
);

const emits = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { t } = useI18n();

const isOpen = computed({
    get: () => props.open,
    set: (v) => {
        if (!v) {
            emits('close', false);
        }
    },
});

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

const isDesktop = useMediaQuery('(min-width: 768px)');

const [DefineBodyTemplate, ReuseBodyTemplate] = createReusableTemplate();
const [DefineFooterTemplate, ReuseFooterTemplate] = createReusableTemplate();
</script>

<template>
    <DefineBodyTemplate>
        <div data-tour="clipboard-modal">
            <UTabs v-model="selectedTab" data-tour="clipboard-tabs" :items="items" variant="pill">
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
        </div>
    </DefineBodyTemplate>

    <DefineFooterTemplate>
        <UFieldGroup class="inline-flex w-full">
            <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

            <UButton
                class="flex-1"
                block
                color="error"
                data-tour="clipboard-clear"
                :label="$t('components.clipboard.clipboard_modal.clear')"
                @click="clipboardStore.clear()"
            />
        </UFieldGroup>
    </DefineFooterTemplate>

    <UModal
        v-if="isDesktop"
        v-model:open="isOpen"
        :title="$t('components.clipboard.clipboard_modal.title')"
        :modal="!props.guideInteractive"
        :overlay="!props.guideInteractive"
        :ui="{ body: 'min-h-90' }"
    >
        <template #body>
            <ReuseBodyTemplate />
        </template>

        <template #footer>
            <ReuseFooterTemplate />
        </template>
    </UModal>

    <UDrawer
        v-else
        v-model:open="isOpen"
        :title="$t('components.clipboard.clipboard_modal.title')"
        :modal="!props.guideInteractive"
        :overlay="!props.guideInteractive"
        :ui="{ body: 'min-h-90' }"
    >
        <template #body>
            <ReuseBodyTemplate />
        </template>

        <template #footer>
            <ReuseFooterTemplate />
        </template>
    </UDrawer>
</template>
