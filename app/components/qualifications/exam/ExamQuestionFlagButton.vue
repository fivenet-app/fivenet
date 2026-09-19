<script lang="ts" setup>
withDefaults(
    defineProps<{
        flagged?: boolean;
        disabled?: boolean;
    }>(),
    {
        flagged: false,
        disabled: false,
    },
);

const emit = defineEmits<{
    (e: 'toggle'): void;
}>();
</script>

<template>
    <UButton
        type="button"
        size="sm"
        variant="ghost"
        :color="flagged ? 'warning' : 'neutral'"
        :icon="flagged ? 'i-mdi-flag' : 'i-mdi-flag-outline'"
        :aria-label="$t(flagged ? 'components.qualifications.exam_view.unflag' : 'components.qualifications.exam_view.flag')"
        :disabled="disabled"
        @click="emit('toggle')"
    />
</template>
