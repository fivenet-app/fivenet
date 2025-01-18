<script lang="ts" setup>
import type { Tab } from '~/store/internet';

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

    tab.value.label = t('components.internet.pages.nic_registrar.title');
    tab.value.icon = undefined;
}
updateTabInfo();

// TODO
</script>

<template>
    <div>Domain Nic Registrar Network Information Center</div>
</template>
