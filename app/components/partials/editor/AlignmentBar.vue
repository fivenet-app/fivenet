<script setup lang="ts">
defineProps<{ modelValue: 'left' | 'center' | 'right' }>();

defineEmits(['update:modelValue']);
</script>

<template>
    <UButtonGroup>
        <UButton
            :disabled="modelValue === 'left'"
            :title="$t('components.partials.TiptapEditor.align_left')"
            color="white"
            @click="$emit('update:modelValue', 'left')"
        >
            <UIcon name="i-mdi-align-horizontal-left" />
        </UButton>
        <UButton
            :disabled="modelValue === 'center'"
            :title="$t('components.partials.TiptapEditor.align_center')"
            color="white"
            @click="$emit('update:modelValue', 'center')"
        >
            <UIcon name="i-mdi-align-horizontal-center" />
        </UButton>
        <UButton
            :disabled="modelValue === 'right'"
            :title="$t('components.partials.TiptapEditor.align_right')"
            color="white"
            @click="$emit('update:modelValue', 'right')"
        >
            <UIcon name="i-mdi-align-horizontal-right" />
        </UButton>
    </UButtonGroup>
</template>
