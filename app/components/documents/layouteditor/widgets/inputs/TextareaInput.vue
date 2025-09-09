<script setup lang="ts">
const props = defineProps<{ modelValue: string; binding?: any }>();
const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>();
function onInput(e: Event) {
    const v = (e.target as HTMLTextAreaElement).value ?? '';
    emit('update:modelValue', applyMax(v, props.binding?.maxChars));
}
function applyMax(v: string, n?: number) {
    return n ? v.slice(0, n) : v;
}
</script>

<template>
    <textarea
        class="h-full w-full resize-none bg-transparent outline-none"
        :rows="binding?.rows ?? 3"
        :value="modelValue"
        :placeholder="binding?.placeholder"
        @input="onInput"
    />
</template>
