<script lang="ts" setup>
import { ColorPicker } from 'vue3-colorpicker';
import 'vue3-colorpicker/style.css';

const props = defineProps<{
    modelValue: string;
    disabled?: boolean;
    hideIcon?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:model-value', value: string): void;
    (e: 'close'): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const color = computed({
    get: () => props.modelValue,
    set: (value) => {
        emit('update:model-value', value);
        emit('close');
    },
});
</script>

<template>
    <UPopover :popper="{ placement: 'bottom-start' }">
        <UButton
            variant="outline"
            color="white"
            :disabled="disabled"
            block
            :icon="!hideIcon ? 'i-mdi-palette' : ''"
            :label="!hideIcon ? '' : '&nbsp;'"
            :style="{ backgroundColor: color }"
            :class="$attrs.class"
        />

        <template #panel>
            <ColorPicker v-model:pureColor="color" is-widget format="hex" picker-type="chrome" disable-alpha disable-history />
        </template>
    </UPopover>
</template>

<style>
.vc-input-toggle {
    display: none !important;
}
</style>
