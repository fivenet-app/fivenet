<script lang="ts" setup>
defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        label?: string;
        icon?: string;
    }>(),
    {
        label: '',
        icon: 'i-mdi-question-mark-circle-outline',
    },
);

const { isHelpSlideoverOpen } = useDashboard();

function openHelp(): void {
    isHelpSlideoverOpen.value = true;
}
</script>

<template>
    <UButton v-bind="$attrs" :trailing-icon="props.icon" color="neutral" variant="ghost" @click="openHelp">
        <span class="hidden truncate sm:block">
            {{ props.label || $t('common.help') }}
        </span>
    </UButton>
</template>
