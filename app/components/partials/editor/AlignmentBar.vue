<script setup lang="ts">
import type { ImageAlign } from '~/composables/tiptap/extensions/EnhancedImage';

const align = defineModel<ImageAlign>({ required: true });
</script>

<template>
    <UFieldGroup>
        <UTooltip :text="$t('components.partials.tiptap_editor.align_left')">
            <UButton
                :disabled="modelValue === 'left'"
                color="neutral"
                icon="i-mdi-align-horizontal-left"
                square
                variant="subtle"
                @click="align = 'left'"
            />
        </UTooltip>

        <UTooltip :text="$t('components.partials.tiptap_editor.align_center')">
            <UButton
                :disabled="modelValue === 'center'"
                color="neutral"
                icon="i-mdi-align-horizontal-center"
                square
                variant="subtle"
                @click="align = 'center'"
            >
            </UButton>
        </UTooltip>

        <UTooltip :text="$t('components.partials.tiptap_editor.align_right')">
            <UButton
                :disabled="modelValue === 'right'"
                color="neutral"
                icon="i-mdi-align-horizontal-right"
                square
                variant="subtle"
                @click="align = 'right'"
            />
        </UTooltip>
    </UFieldGroup>
</template>
