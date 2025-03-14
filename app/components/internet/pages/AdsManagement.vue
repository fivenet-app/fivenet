<script lang="ts" setup>
import type { Tab } from '~/stores/internet';

const props = defineProps<{
    modelValue?: Tab;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: Tab): void;
}>();

const { t } = useI18n();

const tab = useVModel(props, 'modelValue', emit);

function updateTabInfo(): void {
    if (!tab.value) {
        return;
    }

    tab.value.label = t('components.internet.pages.ads_management.title');
    tab.value.icon = undefined;
}
updateTabInfo();

// TODO
</script>

<template>
    <div>Advertisement Central</div>
</template>
