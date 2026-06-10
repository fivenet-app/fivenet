<script setup lang="ts">
import type { ImageAlign } from '~/composables/tiptap/extensions/EnhancedImage';

const align = defineModel<ImageAlign>({ required: true });
</script>

<template>
    <UFieldGroup>
        <UButton
            :disabled="modelValue === 'left'"
            :title="$t('components.partials.tiptap_editor.align_left')"
            color="neutral"
            variant="subtle"
            icon="i-mdi-align-horizontal-left"
            square
            @click="align = 'left'"
        />

        <UButton
            :disabled="modelValue === 'center'"
            :title="$t('components.partials.tiptap_editor.align_center')"
            color="neutral"
            variant="subtle"
            icon="i-mdi-align-horizontal-center"
            square
            @click="align = 'center'"
        >
        </UButton>

        <UButton
            :disabled="modelValue === 'right'"
            :title="$t('components.partials.tiptap_editor.align_right')"
            color="neutral"
            variant="subtle"
            icon="i-mdi-align-horizontal-right"
            square
            @click="align = 'right'"
        />
    </UFieldGroup>
</template>
