<script setup lang="ts">
import { computed, ref, watch } from 'vue';

const props = defineProps<{
    modelValue: string;
    binding?: {
        boxCharCount?: number; // explicit number of boxes
        widthMm?: number; // optional: derive box count by width/font (future)
        letterSpacing?: number; // px or mm gap between boxes
        maxChars?: number;
    };
}>();
const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>();

const raw = ref(props.modelValue ?? '');
watch(
    () => props.modelValue,
    (v) => (raw.value = v ?? ''),
);

const count = computed(() => props.binding?.boxCharCount ?? props.binding?.maxChars ?? 10);

function setCharAt(s: string, i: number, c: string) {
    const arr = (s ?? '').split('');
    arr[i] = c;
    return arr.join('');
}

function onKey(i: number, e: KeyboardEvent) {
    const t = e.target as HTMLInputElement;
    // Allow navigation
    if (e.key === 'ArrowLeft' || e.key === 'ArrowRight' || e.key === 'Tab' || e.key === 'Backspace' || e.key.length > 1) {
        if (e.key === 'Backspace' && !t.value) {
            focusBox(i - 1);
        }
        return;
    }
    // Accept only a single printable char; enforce regex if needed (A-Z0-9.)
    e.preventDefault();
    const ch = e.key.slice(0, 1);
    const next = setCharAt(padRight(raw.value, count.value), i, ch);
    update(next);
    focusBox(i + 1);
}

function update(v: string) {
    const trimmed = v.slice(0, count.value).replace(/\n/g, '');
    raw.value = trimmed;
    emit('update:modelValue', trimmed);
}

function padRight(s: string, n: number) {
    return (s ?? '').padEnd(n, ' ');
}

function focusBox(i: number) {
    const el = document.querySelector<HTMLInputElement>(`[data-box-idx="${i}"]`);
    if (el) el.focus();
}
</script>

<template>
    <div class="flex h-full w-full items-center gap-1" :style="{ letterSpacing: (binding?.letterSpacing ?? 0) + 'px' }">
        <template v-for="i in count" :key="i">
            <input
                :data-box-idx="i - 1"
                class="h-9 w-7 rounded-sm border border-gray-400 bg-transparent text-center outline-none"
                :value="
                    (raw || '')
                        .padEnd(count, ' ')
                        .charAt(i - 1)
                        .trim()
                "
                @keydown="onKey(i - 1, $event as KeyboardEvent)"
            />
        </template>
    </div>
</template>
