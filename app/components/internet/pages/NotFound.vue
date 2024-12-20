<script lang="ts" setup>
import { useInternetStore, type Tab } from '~/store/internet';

const props = defineProps<{
    modelValue?: Tab;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: Tab): void;
}>();

const { t } = useI18n();

const tab = useVModel(props, 'modelValue', emit);

const internetStore = useInternetStore();

function updateTabInfo(): void {
    if (!tab.value) {
        return;
    }

    tab.value.label = t('components.internet.not_found.title');
    tab.value.icon = 'i-mdi-error-outline';
}

updateTabInfo();
watch(tab, () => updateTabInfo());
</script>

<template>
    <UContainer class="mt-4">
        <ULandingCard
            v-if="modelValue"
            :title="$t('components.internet.not_found.title')"
            icon="i-mdi-information-slab-circle-outline"
            class="w-screen max-w-md"
            orientation="vertical"
        >
            <template #description>
                {{ $t('components.internet.not_found.description', { domain: modelValue.domain }) }}
            </template>

            <div>
                <UButton :label="$t('common.back')" color="black" icon="i-mdi-arrow-back" @click="internetStore.back()" />
            </div>
        </ULandingCard>
    </UContainer>
</template>
