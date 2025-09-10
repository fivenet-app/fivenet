<script setup lang="ts">
import type { DisplayBinding, InputWidget } from '../types';
import BoxedTextInput from './inputs/BoxedTextInput.vue';
import TextInput from './inputs/TextInput.vue';
import TextareaInput from './inputs/TextareaInput.vue';

const props = defineProps<{
    path: string;
    value: string | number | null | undefined;
    binding: DisplayBinding;
    widget?: InputWidget; // optional override
}>();
const emit = defineEmits<{ (e: 'update', value: string): void }>();

const chosen = computed<InputWidget>(() => props.widget ?? 'text');
</script>

<template>
    <component
        :is="mapWidget(chosen)"
        :model-value="value ?? ''"
        :binding="binding"
        @update:model-value="(val) => emit('update', val)"
    />
</template>

<script lang="ts">
function mapWidget(w: string) {
    switch (w) {
        case 'textarea':
            return TextareaInput;
        case 'boxed-text':
            return BoxedTextInput;
        default:
            return TextInput;
    }
}
</script>
