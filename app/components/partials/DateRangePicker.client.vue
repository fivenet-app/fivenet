<script setup lang="ts">
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';
import { DatePicker as VCalendarDatePicker } from 'v-calendar';
import 'v-calendar/dist/style.css';

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        modelValue: { start: Date; end: Date } | undefined;
        clearable?: boolean;
    }>(),
    {
        modelValue: undefined,
        clearable: false,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', value: typeof props.modelValue): void;
    (e: 'close'): void;
}>();

const date = computed({
    get: () => (props.modelValue === undefined ? null : props.modelValue),
    set: (value) => emit('update:modelValue', value === null ? undefined : value),
});

const breakpoints = useBreakpoints(breakpointsTailwind);

const smallerThanSm = breakpoints.smaller('sm');

const attrs = {
    transparent: true,
    borderless: true,
    color: 'primary',
    'is-dark': { selector: 'html', darkClass: 'dark' },
    'first-day-of-week': 2,
    'show-weeknumbers': true,
};
</script>

<template>
    <VCalendarDatePicker
        v-model.range="date"
        :columns="smallerThanSm ? 1 : 2"
        :rows="smallerThanSm ? 2 : 1"
        v-bind="{ ...attrs, ...$attrs }"
        @update:model-value="date = $event"
        @close="$emit('close')"
    >
        <template v-if="clearable" #footer>
            <div class="w-full px-4 pb-3">
                <UButton block @click="date = null">
                    {{ $t('common.clear') }}
                </UButton>
            </div>
        </template>
    </VCalendarDatePicker>
</template>

<style>
:root {
    /* Font Family override */
    --vc-font-family: DM Sans, ui-sans-serif, system-ui, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
        'Noto Color Emoji', DM Sans, ui-sans-serif, system-ui, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji',
        'Segoe UI Symbol', 'Noto Color Emoji';
    --vc-gray-50: rgb(var(--color-gray-50));
    --vc-gray-100: rgb(var(--color-gray-100));
    --vc-gray-200: rgb(var(--color-gray-200));
    --vc-gray-300: rgb(var(--color-gray-300));
    --vc-gray-400: rgb(var(--color-gray-400));
    --vc-gray-500: rgb(var(--color-gray-500));
    --vc-gray-600: rgb(var(--color-gray-600));
    --vc-gray-700: rgb(var(--color-gray-700));
    --vc-gray-800: rgb(var(--color-gray-800));
    --vc-gray-900: rgb(var(--color-gray-900));
}

.vc-primary {
    --vc-accent-50: rgb(var(--color-primary-50));
    --vc-accent-100: rgb(var(--color-primary-100));
    --vc-accent-200: rgb(var(--color-primary-200));
    --vc-accent-300: rgb(var(--color-primary-300));
    --vc-accent-400: rgb(var(--color-primary-400));
    --vc-accent-500: rgb(var(--color-primary-500));
    --vc-accent-600: rgb(var(--color-primary-600));
    --vc-accent-700: rgb(var(--color-primary-700));
    --vc-accent-800: rgb(var(--color-primary-800));
    --vc-accent-900: rgb(var(--color-primary-900));
}
</style>
