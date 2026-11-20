<script setup lang="ts">
import type { ImageAlign } from '~/composables/tiptap/extensions/EnhancedImage';

const align = defineModel<ImageAlign>({ required: true });

const { t } = useI18n();

const alignmentItems = computed(() => [
    {
        value: 'left' as const,
        icon: 'i-mdi-align-horizontal-left',
        'aria-label': t('components.partials.tiptap_editor.align_left'),
        activeClass: 'bg-accented',
    },
    {
        value: 'center' as const,
        icon: 'i-mdi-align-horizontal-center',
        'aria-label': t('components.partials.tiptap_editor.align_center'),
        activeClass: 'bg-accented',
    },
    {
        value: 'right' as const,
        icon: 'i-mdi-align-horizontal-right',
        'aria-label': t('components.partials.tiptap_editor.align_right'),
        activeClass: 'bg-accented',
    },
]);

function setAlignment(value: unknown): void {
    if (value !== 'left' && value !== 'center' && value !== 'right') return;

    align.value = value;
}
</script>

<template>
    <UFormField>
        <UTooltip :text="$t('components.partials.tiptap_editor.image_alignment')">
            <UToggleGroup
                :model-value="align"
                :items="alignmentItems"
                :aria-label="$t('components.partials.tiptap_editor.image_alignment')"
                type="single"
                attached
                color="neutral"
                variant="subtle"
                active-color="neutral"
                active-variant="subtle"
                @update:model-value="setAlignment"
            />
        </UTooltip>
    </UFormField>
</template>
