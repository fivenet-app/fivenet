<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        modelValue?: string;
    }>(),
    {
        modelValue: '',
    },
);

defineOptions({
    inheritAttrs: true,
});

const emit = defineEmits(['update:modelValue']);
const data = useVModel(props, 'modelValue', emit);
</script>

<template>
    <UInput v-model="data" @focusin="focusTablet(true)" @focusout="focusTablet(false)" @keydown.esc="$event.target.blur()" />
</template>
