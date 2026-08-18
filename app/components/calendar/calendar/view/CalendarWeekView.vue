<script setup lang="ts">
import { addDays, format, isToday } from 'date-fns';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';
import {
    HOUR_HEIGHT,
    getCalendarCreateRangeFromClick,
    layoutAllDayEntries,
    layoutDayEntries,
    viewDateRange,
} from '~/utils/calendar-view';
import { dateToDateString, getWeekNumber } from '~/utils/time';
import { isCalendarEntryAllDay } from '~/utils/calendar';
import CalendarEntryChip from './CalendarEntryChip.vue';
import CalendarNowIndicator from '~/components/calendar/calendar/view/CalendarNowIndicator.vue';

const ALL_DAY_LANE_HEIGHT = 28;
const ALL_DAY_ROW_BORDER = 1;

const props = defineProps<{
    date: Date;
    entries: CalendarEntry[];
    canCreate?: boolean;
}>();

const emit = defineEmits<{
    (e: 'create', value: { startTime: Date; endTime: Date }): void;
    (e: 'select', entry: CalendarEntry): void;
    (e: 'edit', entry: CalendarEntry): void;
    (e: 'share', entry: CalendarEntry): void;
    (e: 'delete', entry: CalendarEntry): void;
}>();

const days = computed(() => {
    const { start } = viewDateRange('week', props.date);
    return Array.from({ length: 7 }, (_, index) => addDays(start, index));
});

const allDayBaseEntries = computed(() => props.entries.filter((entry) => isCalendarEntryAllDay(entry)));
const timedBaseEntries = computed(() => props.entries.filter((entry) => !isCalendarEntryAllDay(entry)));

const allDayEntries = computed(() => layoutAllDayEntries(allDayBaseEntries.value, days.value));
const timedEntries = computed(() =>
    days.value.map((day) => layoutDayEntries(timedBaseEntries.value, day, { preferWideSolitaryEvents: true })),
);

const hours = Array.from({ length: 24 }, (_, index) => index);
const weekNumber = computed(() => getWeekNumber(days.value[0] ?? props.date));
const allDayHeight = computed(() =>
    allDayEntries.value.length
        ? (Math.max(...allDayEntries.value.map((entry) => entry.lane)) + 1) * ALL_DAY_LANE_HEIGHT + ALL_DAY_ROW_BORDER
        : 0,
);
const isEntryPopoverOpen = ref(false);

function openCreateAt(day: Date, event: MouseEvent): void {
    if (!props.canCreate) return;
    const range = getCalendarCreateRangeFromClick(day, event);
    if (!range) return;

    emit('create', range);
}

function handleEntryPopoverOpen(open: boolean): void {
    isEntryPopoverOpen.value = open;
}
</script>

<template>
    <div class="flex h-full min-h-0 min-w-0 flex-col overflow-hidden border border-default">
        <div class="shrink-0 border-b border-default bg-muted/30">
            <div class="grid grid-cols-[3.5rem_repeat(7,minmax(0,1fr))]">
                <div
                    class="flex items-center justify-center border-e border-default px-2 py-2 text-xs font-semibold text-muted"
                >
                    {{ $t('common.calendar_week') }} {{ weekNumber }}
                </div>
                <div
                    v-for="day in days"
                    :key="dateToDateString(day)"
                    class="border-s border-default px-3 py-2 text-center"
                    :class="isToday(day) ? 'bg-primary/5' : ''"
                >
                    <p class="text-xs tracking-wide text-muted uppercase">
                        {{ format(day, 'EEE') }}
                    </p>
                    <div class="mt-0.5 flex items-center justify-center gap-1">
                        <p class="font-semibold">
                            {{ day.getDate() }}
                        </p>
                        <UBadge v-if="isToday(day)" size="xs" color="warning" :label="$t('common.today')" />
                    </div>
                </div>
            </div>

            <div
                v-if="allDayEntries.length"
                class="grid grid-cols-[3.5rem_repeat(7,minmax(0,1fr))] border-t border-default bg-elevated/60"
            >
                <div class="truncate border-e border-default px-2 py-2 text-[10px] tracking-wide text-muted uppercase">
                    <UTooltip :text="$t('common.all_day')">
                        {{ $t('common.all_day') }}
                    </UTooltip>
                </div>
                <div class="col-span-7 grid grid-cols-7 gap-y-1 px-1 py-2" :style="{ minHeight: `${allDayHeight}px` }">
                    <CalendarEntryChip
                        v-for="{ entry, colStart, colSpan, lane } in allDayEntries"
                        :key="entry.occurrence?.key ?? entry.id"
                        :entry="entry"
                        compact
                        :show-time="false"
                        class="mx-0.5"
                        :style="{ gridColumn: `${colStart + 1} / span ${colSpan}`, gridRow: lane + 1 }"
                        @update:popover-open="handleEntryPopoverOpen"
                        @share="emit('share', $event)"
                        @delete="emit('delete', $event)"
                        @select="emit('select', $event)"
                        @edit="emit('edit', $event)"
                    />
                </div>
            </div>
        </div>

        <div class="flex min-h-0 flex-1 overflow-auto">
            <div class="grid h-[calc(24*64px)] min-h-0 w-full grid-cols-[3.5rem_repeat(7,minmax(0,1fr))]">
                <div class="relative border-e border-default">
                    <div
                        v-for="hour in hours"
                        :key="hour"
                        class="absolute end-2 -translate-y-1/2 text-[11px] text-muted tabular-nums"
                        :style="{ top: `${hour * HOUR_HEIGHT}px` }"
                    >
                        {{ format(new Date(2000, 0, 1, hour, 0), 'p') }}
                    </div>
                </div>

                <div
                    v-for="(day, index) in days"
                    :key="dateToDateString(day)"
                    class="relative border-s border-default"
                    :style="{ height: `${24 * HOUR_HEIGHT}px` }"
                >
                    <div
                        v-if="props.canCreate"
                        class="absolute inset-0 z-0 cursor-pointer"
                        :class="isEntryPopoverOpen ? 'pointer-events-none' : ''"
                        aria-hidden="true"
                        @click="openCreateAt(day, $event)"
                    />

                    <div
                        v-for="hour in 23"
                        :key="hour"
                        class="absolute inset-x-0 z-0 border-t border-default"
                        :style="{ top: `${hour * HOUR_HEIGHT}px` }"
                    />

                    <div
                        v-for="positioned in timedEntries[index] ?? []"
                        :key="positioned.entry.occurrence?.key ?? positioned.entry.id"
                        class="absolute z-10"
                        :style="{
                            top: `${positioned.top + 2}px`,
                            height: `${Math.max(positioned.height - 3, 24)}px`,
                            insetInlineStart: `calc(${positioned.left}% + 1px)`,
                            width: `calc(${positioned.width}% - 2px)`,
                        }"
                    >
                        <CalendarEntryChip
                            :entry="positioned.entry"
                            class="h-full"
                            :show-time="true"
                            :stacked="positioned.height >= 52"
                            @update:popover-open="handleEntryPopoverOpen"
                            @share="emit('share', $event)"
                            @delete="emit('delete', $event)"
                            @select="emit('select', $event)"
                            @edit="emit('edit', $event)"
                        />
                    </div>

                    <ClientOnly>
                        <CalendarNowIndicator v-if="isToday(day)" />
                    </ClientOnly>
                </div>
            </div>
        </div>
    </div>
</template>
