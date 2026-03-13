<script setup lang="ts">
import { CalendarDate, getLocalTimeZone, today } from '@internationalized/date';
import { subYears } from 'date-fns';
import type { Range } from './helpers';

const { t } = useI18n();

const { format: formatDate } = useDateFormatter();

const selected = defineModel<Range>({ required: true });

const ranges = [
    { label: t('common.date_range.last_7_days'), days: 7 },
    { label: t('common.date_range.last_14_days'), days: 14 },
    { label: t('common.date_range.last_30_days'), days: 30 },
    { label: t('common.date_range.last_90_days'), months: 3 },
    { label: t('common.date_range.last_180_days'), months: 6 },
    { label: t('common.date_range.last_365_days'), years: 1 },
];

const toCalendarDate = (date: Date) => {
    return new CalendarDate(date.getFullYear(), date.getMonth() + 1, date.getDate());
};

const calendarRange = computed({
    get: () => ({
        start: selected.value.start ? toCalendarDate(selected.value.start) : undefined,
        end: selected.value.end ? toCalendarDate(selected.value.end) : undefined,
    }),
    set: (newValue: { start: CalendarDate | null; end: CalendarDate | null }) => {
        selected.value = {
            start: newValue.start ? newValue.start.toDate(getLocalTimeZone()) : new Date(),
            end: newValue.end ? newValue.end.toDate(getLocalTimeZone()) : new Date(),
        };
    },
});

const isRangeSelected = (range: { days?: number; months?: number; years?: number }) => {
    if (!selected.value.start || !selected.value.end) return false;

    const currentDate = today(getLocalTimeZone());
    let startDate = currentDate.copy();

    if (range.days) {
        startDate = startDate.subtract({ days: range.days });
    } else if (range.months) {
        startDate = startDate.subtract({ months: range.months });
    } else if (range.years) {
        startDate = startDate.subtract({ years: range.years });
    }

    const selectedStart = toCalendarDate(selected.value.start);
    const selectedEnd = toCalendarDate(selected.value.end);

    return selectedStart.compare(startDate) === 0 && selectedEnd.compare(currentDate) === 0;
};

const selectRange = (range: { days?: number; months?: number; years?: number }) => {
    const endDate = today(getLocalTimeZone());
    let startDate = endDate.copy();

    if (range.days) {
        startDate = startDate.subtract({ days: range.days });
    } else if (range.months) {
        startDate = startDate.subtract({ months: range.months });
    } else if (range.years) {
        startDate = startDate.subtract({ years: range.years });
    }

    selected.value = {
        start: startDate.toDate(getLocalTimeZone()),
        end: endDate.toDate(getLocalTimeZone()),
    };
};

const now = new Date();
const lastYear = subYears(now, 1);
const minDate = new CalendarDate(lastYear.getFullYear(), lastYear.getMonth() + 1, lastYear.getDate());
const maxDate = new CalendarDate(now.getFullYear(), now.getMonth() + 1, now.getDate());
</script>

<template>
    <UPopover :content="{ align: 'start' }" :modal="true">
        <UButton color="neutral" variant="ghost" icon="i-lucide-calendar" class="group data-[state=open]:bg-elevated">
            <span class="truncate">
                <template v-if="selected.start">
                    <template v-if="selected.end"> {{ formatDate(selected.start) }} - {{ formatDate(selected.end) }} </template>
                    <template v-else>
                        {{ formatDate(selected.start) }}
                    </template>
                </template>
                <template v-else> {{ $t('common.pick_date') }} </template>
            </span>

            <template #trailing>
                <UIcon
                    name="i-lucide-chevron-down"
                    class="size-5 shrink-0 text-dimmed transition-transform duration-200 group-data-[state=open]:rotate-180"
                />
            </template>
        </UButton>

        <template #content>
            <div class="flex items-stretch divide-default sm:divide-x">
                <div class="hidden flex-col justify-center sm:flex">
                    <UButton
                        v-for="(range, index) in ranges"
                        :key="index"
                        :label="range.label"
                        color="neutral"
                        variant="ghost"
                        class="rounded-none px-4"
                        :class="[isRangeSelected(range) ? 'bg-elevated' : 'hover:bg-elevated/50']"
                        truncate
                        @click="selectRange(range)"
                    />
                </div>

                <UCalendar
                    v-model="calendarRange"
                    class="p-2"
                    :number-of-months="2"
                    range
                    :maximum-days="365"
                    :min-value="minDate"
                    :max-value="maxDate"
                />
            </div>
        </template>
    </UPopover>
</template>
