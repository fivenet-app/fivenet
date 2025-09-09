<script setup lang="ts">
const props = defineProps<{ modelValue: string; binding?: any }>();
const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>();
function onInput(e: Event) {
    const v = (e.target as HTMLInputElement).value ?? '';
    emit('update:modelValue', applyMax(v, props.binding?.maxChars));
}
function applyMax(v: string, n?: number) {
    return n ? v.slice(0, n) : v;
}
</script>

<template>
    <input
        class="h-full w-full bg-transparent outline-none"
        :value="modelValue"
        :placeholder="binding?.placeholder"
        @input="onInput"
    />
</template>
