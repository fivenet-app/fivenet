<script setup lang="ts">
import { Calendar as VCalendar } from 'v-calendar';
// @ts-ignore
import type { DatePickerDate, DatePickerRangeObject } from 'v-calendar/dist/types/src/use/datePicker';
import 'v-calendar/dist/style.css';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        modelValue: DatePickerDate | DatePickerRangeObject | null;
    }>(),
    {
        modelValue: null,
    },
);

const emits = defineEmits<{
    (e: 'update:model-value', value: Date | null): void;
    (e: 'selected', entry: CalendarEntry): void;
}>();

const date = computed({
    get: () => props.modelValue,
    set: (value) => emits('update:model-value', value),
});

const attrs = {
    transparent: true,
    borderless: true,
    color: 'primary',
    'is-dark': { selector: 'html', darkClass: 'dark' },
    'first-day-of-week': 2,
    'show-weeknumbers': true,
    'trim-weeks': true,
};

const masks = {
    weekdays: 'WWW',
};
</script>

<template>
    <div class="custom-calendar" :class="$attrs.class">
        <VCalendar v-model="date" view="monthly" :columns="1" :rows="1" :masks="masks" v-bind="{ ...attrs, ...$attrs }">
            <template #day-content="{ day, attributes }">
                <div class="z-10 flex h-full flex-col overflow-hidden">
                    <span class="day-label text-sm text-gray-900 dark:text-white">{{ day.day }}</span>
                    <div class="flex-grow overflow-x-auto overflow-y-auto">
                        <button
                            v-for="attr in attributes"
                            :key="attr.key"
                            class="vc-day-entry mb-1 mt-0 rounded-sm p-1 text-left text-xs leading-tight"
                            :class="attr.customData.class"
                            @click="$emit('selected', attr.customData)"
                        >
                            {{ attr.customData.title }}
                            <template v-if="attr.customData.time">
                                <br />
                                {{ attr.customData.time }}
                            </template>
                        </button>
                    </div>
                </div>
            </template>
        </VCalendar>
    </div>
</template>

<style>
:root {
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

<style lang="postcss" scoped>
::-webkit-scrollbar {
    width: 0px;
}

::-webkit-scrollbar-track {
    display: none;
}

.custom-calendar:deep(.vc-container) {
    --day-border: 1px solid rgb(var(--color-gray-700));
    --day-border-highlight: 1px solid rgb(var(--color-gray-600));
    --day-width: 120px;
    --day-height: 120px;
    --weekday-border: 1px solid rgb(var(--color-primary-900));

    border-radius: 0;
    width: 100%;

    & .vc-weekday {
        background-color: rgb(var(--color-primary-400));
        border-bottom: var(--weekday-border);
        border-top: var(--weekday-border);
        padding: 5px 0;
    }
    & .vc-day {
        padding: 0 5px 3px 5px;
        text-align: left;
        height: var(--day-height);
        min-width: var(--day-width);
        background-color: rgb(var(--color-gray-900));
        & :hover {
            background-color: rgb(var(--color-gray-800));
        }

        &.weekday-1,
        &.weekday-7 {
            background-color: rgb(var(--color-gray-800));
            & :hover {
                background-color: rgb(var(--color-gray-700));
            }
        }

        &:not(.on-bottom) {
            border-bottom: var(--day-border);
            &.weekday-1 {
                border-bottom: var(--day-border-highlight);
            }
        }
    }
    & .vc-day-dots {
        margin-bottom: 5px;
    }
}
</style>
