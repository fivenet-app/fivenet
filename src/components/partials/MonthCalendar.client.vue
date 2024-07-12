<script setup lang="ts">
import { Calendar as VCalendar } from 'v-calendar';
// @ts-ignore
import 'v-calendar/dist/style.css';
import type { CalendarView } from 'v-calendar/dist/types/src/use/calendar.js';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';
import MonthCalendarDay from './MonthCalendarDay.vue';

defineOptions({
    inheritAttrs: false,
});

withDefaults(
    defineProps<{
        view?: CalendarView;
    }>(),
    {
        view: 'monthly',
    },
);

const emits = defineEmits<{
    (e: 'selected', entry: CalendarEntry): void;
}>();

const calRef = ref<InstanceType<typeof VCalendar> | null>(null);

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

defineExpose({
    calRef,
});
</script>

<template>
    <div class="custom-calendar" :class="$attrs.class">
        <VCalendar ref="calRef" :view="view" :columns="1" :rows="1" :masks="masks" v-bind="{ ...attrs, ...$attrs }">
            <template #day-content="{ day, attributes }">
                <MonthCalendarDay :day="day" :attributes="attributes" @selected="$emit('selected', $event)" />
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
    --day-width: 110px;
    --day-height: 115px;
    --weekday-border: 1px solid rgb(var(--color-primary-900));

    border-radius: 0;
    width: 100%;
    height: 100%;

    & .vc-pane-container,
    & .vc-pane-layout,
    & .vc-pane {
        display: flex;
        flex-direction: column;
        flex: 1;
    }

    & .vc-header {
        margin-top: 2px;
        padding-left: 0;
        padding-right: 0;
    }
    & .vc-weeks {
        padding-left: 0;
        padding-right: 0;
        padding: 0;

        display: flex;
        flex-direction: column;
        flex: 1;
    }

    & .vc-week {
        flex: 1;
    }

    & .vc-weekday {
        background-color: rgb(var(--color-primary-400));
        border-bottom: var(--weekday-border);
        border-top: var(--weekday-border);
        padding: 5px 0;
    }
    & .vc-weekday.vc-weekday-1,
    & .vc-weekday.vc-weekday-7 {
        background-color: rgb(var(--color-primary-600));
    }

    .vc-day:hover {
        background-color: rgb(var(--color-gray-800));

        &.weekday-1:hover,
        &.weekday-7:hover {
            background-color: rgb(var(--color-gray-700));
        }
    }

    & .vc-day {
        padding: 0 5px 3px 5px;
        text-align: left;
        min-height: var(--day-height);
        min-width: var(--day-width);
        background-color: rgb(var(--color-gray-900));

        height: 100%;

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
.custom-calendar:deep(.vc-container.vc-weekly) {
    & .vc-day {
        &:not(.on-bottom) {
            border-bottom: none;
            &.weekday-1 {
                border-bottom: none;
            }
        }
    }
}
</style>
