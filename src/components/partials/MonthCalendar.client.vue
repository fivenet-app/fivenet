<script setup lang="ts">
import { useBreakpoints, breakpointsTailwind } from '@vueuse/core';
import { DatePicker as VCalendarDatePicker } from 'v-calendar';
// @ts-ignore
import type { DatePickerDate, DatePickerRangeObject } from 'v-calendar/dist/types/src/use/datePicker';
import 'v-calendar/dist/style.css';

defineOptions({
    inheritAttrs: false,
});

const props = defineProps({
    modelValue: {
        type: [Date, Object] as PropType<DatePickerDate | DatePickerRangeObject | null>,
        default: null,
    },
    clearable: {
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits(['update:model-value', 'close']);

const date = computed({
    get: () => props.modelValue,
    set: (value) => {
        emit('update:model-value', value);
        emit('close');
    },
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
    'trim-weeks': true,
};

const masks = {
    weekdays: 'WWW',
};

const month = new Date().getMonth();
const year = new Date().getFullYear();

const attributes = [
    {
        key: 1,
        customData: {
            title: 'Lunch with mom.',
            time: '17:30',
            class: 'bg-red-600 text-white',
        },
        dates: new Date(year, month, 1),
    },
    {
        key: 2,
        customData: {
            title: 'Take Noah to basketball practice',
            time: '17:30',
            class: 'bg-blue-500 text-white',
        },
        dates: new Date(year, month, 2),
    },
    {
        key: 3,
        customData: {
            title: "Noah's basketball game.",
            time: '17:30',
            class: 'bg-blue-500 text-white',
        },
        dates: new Date(year, month, 5),
    },
    {
        key: 4,
        customData: {
            title: 'Take car to the shop',
            time: '17:30',
            class: 'bg-indigo-500 text-white',
        },
        dates: new Date(year, month, 5),
    },
    {
        key: 4,
        customData: {
            title: 'Meeting with new client.',
            time: '17:30',
            class: 'bg-teal-500 text-white',
        },
        dates: new Date(year, month, 7),
    },
    {
        key: 5,
        customData: {
            title: "Mia's gymnastics practice.",
            time: '17:30',
            class: 'bg-pink-500 text-white',
        },
        dates: new Date(year, month, 11),
    },
];
</script>

<template>
    <div class="custom-calendar">
        <VCalendarDatePicker
            v-if="date && (date as DatePickerRangeObject)?.start && (date as DatePickerRangeObject)?.end"
            v-model.range="date"
            class="custom-calendar"
            :columns="smallerThanSm ? 1 : 2"
            :rows="smallerThanSm ? 2 : 1"
            :attributes="attributes"
            :masks="masks"
            v-bind="{ ...attrs, ...$attrs }"
        >
            <template #day-content="{ day, attributes }">
                <div class="z-10 flex h-full flex-col overflow-hidden">
                    <span class="day-label text-sm text-gray-900 dark:text-white">{{ day.day }}</span>
                    <div class="flex-grow overflow-x-auto overflow-y-auto">
                        <p
                            v-for="attr in attributes"
                            :key="attr.key"
                            class="mb-1 mt-0 rounded-sm p-1 text-xs leading-tight"
                            :class="attr.customData.class"
                        >
                            {{ attr.customData.title }}
                        </p>
                    </div>
                </div>
            </template>
        </VCalendarDatePicker>
        <VCalendarDatePicker
            v-else
            v-model="date"
            v-bind="{ ...attrs, ...$attrs }"
            class="custom-calendar"
            :attributes="attributes"
            :masks="masks"
        >
            <template #day-content="{ day, attributes }">
                <div class="z-10 flex h-full flex-col overflow-hidden">
                    <span class="day-label text-sm text-gray-900 dark:text-white">{{ day.day }}</span>
                    <div class="flex-grow overflow-x-auto overflow-y-auto">
                        <p
                            v-for="attr in attributes"
                            :key="attr.key"
                            class="mb-1 mt-0 rounded-sm p-1 text-xs leading-tight"
                            :class="attr.customData.class"
                        >
                            {{ attr.customData.title }} -
                            {{ attr.customData.time }}
                        </p>
                    </div>
                </div>
            </template>
        </VCalendarDatePicker>
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
    --day-width: 90px;
    --day-height: 90px;
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
        &.weekday-1,
        &.weekday-7 {
            background-color: rgb(var(--color-gray-800));
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
