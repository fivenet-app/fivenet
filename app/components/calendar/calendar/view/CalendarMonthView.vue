<script setup lang="ts">
import { addDays, addHours, isSameMonth, isToday, startOfDay, startOfMonth, startOfWeek } from 'date-fns';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';
import { dateToDateString, getWeekNumber } from '~/utils/time';
import { groupEntriesByDay } from '~/utils/calendar-view';
import { sortCalendarEntriesForDisplay } from '~/utils/calendar';
import CalendarEntryChip from './CalendarEntryChip.vue';

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

const weeks = computed(() => {
    const start = startOfWeek(startOfMonth(props.date), { weekStartsOn: 1 });
    return Array.from({ length: 6 }, (_, weekIndex) =>
        Array.from({ length: 7 }, (_, dayIndex) => addDays(start, weekIndex * 7 + dayIndex)),
    );
});

const weekNumbers = computed(() => weeks.value.map((week) => getWeekNumber(week[0] ?? props.date)));

const entriesByDay = computed(() => groupEntriesByDay(props.entries));
const isEntryPopoverOpen = ref(false);

function dayEntries(day: Date): CalendarEntry[] {
    return sortCalendarEntriesForDisplay(entriesByDay.value.get(dateToDateString(day)) ?? []);
}

function openCreateAt(day: Date): void {
    if (!props.canCreate) return;

    const startTime = addHours(startOfDay(day), 9);
    emit('create', {
        startTime,
        endTime: addHours(startTime, 1),
    });
}

function handleEntryPopoverOpen(open: boolean): void {
    isEntryPopoverOpen.value = open;
}
</script>

<template>
    <div class="flex h-full min-h-0 min-w-0 flex-col overflow-hidden border border-default">
        <div
            class="grid grid-cols-[2rem_repeat(7,minmax(0,1fr))] border-b border-default bg-muted/30 text-xs font-medium tracking-wide text-muted uppercase"
        >
            <div class="border-e border-default px-1 py-2 text-center text-xs">
                {{ $t('common.calendar_week') }}
            </div>
            <div
                v-for="weekday in [
                    $t('common.week_days.monday'),
                    $t('common.week_days.tuesday'),
                    $t('common.week_days.wednesday'),
                    $t('common.week_days.thursday'),
                    $t('common.week_days.friday'),
                    $t('common.week_days.saturday'),
                    $t('common.week_days.sunday'),
                ]"
                :key="weekday"
                class="px-3 py-2"
            >
                {{ weekday }}
            </div>
        </div>

        <div class="flex min-h-0 flex-1 overflow-auto">
            <div class="grid h-full min-h-full w-full grid-rows-6">
                <div
                    v-for="(week, index) in weeks"
                    :key="dateToDateString(week[0]!)"
                    class="grid min-h-0 grid-cols-[2rem_repeat(7,minmax(0,1fr))] border-b border-default last:border-b-0"
                >
                    <div
                        class="flex flex-col items-center justify-center border-e border-default px-1 py-1 text-center text-xs font-medium text-muted"
                    >
                        <span class="leading-tight">
                            {{ weekNumbers[index] }}
                        </span>
                    </div>

                    <div
                        v-for="day in week"
                        :key="dateToDateString(day)"
                        class="relative flex min-h-0 min-w-0 flex-col border-e border-default p-2 last:border-e-0"
                        :class="[
                            isSameMonth(day, props.date) ? 'bg-transparent' : 'bg-muted/20 text-muted',
                            isToday(day) ? 'bg-primary/5 ring-1 ring-primary/40 ring-inset' : '',
                        ]"
                    >
                        <div
                            v-if="props.canCreate"
                            class="absolute inset-0 z-0 cursor-pointer"
                            :class="isEntryPopoverOpen ? 'pointer-events-none' : ''"
                            aria-hidden="true"
                            @click="openCreateAt(day)"
                        />

                        <div class="relative z-10 mb-2 flex items-center justify-between gap-2">
                            <span class="text-sm font-semibold">
                                {{ day.getDate() }}
                            </span>

                            <UBadge v-if="isToday(day)" size="xs" color="warning" :label="$t('common.today')" />
                        </div>

                        <div class="relative z-10 flex min-h-0 flex-col gap-1 overflow-hidden">
                            <CalendarEntryChip
                                v-for="entry in dayEntries(day).slice(0, 4)"
                                :key="entry.occurrence?.key ?? entry.id"
                                :entry="entry"
                                compact
                                :show-time="false"
                                @update:popover-open="handleEntryPopoverOpen"
                                @share="emit('share', $event)"
                                @delete="emit('delete', $event)"
                                @select="emit('select', $event)"
                                @edit="emit('edit', $event)"
                            />

                            <UPopover v-if="dayEntries(day).length > 4" :ui="{ content: 'w-80 p-3' }">
                                <UButton
                                    color="neutral"
                                    variant="ghost"
                                    size="xs"
                                    class="justify-start px-1.5 py-0.5 font-normal text-muted"
                                >
                                    +{{ dayEntries(day).length - 4 }} {{ $t('common.more') }}
                                </UButton>

                                <template #content>
                                    <div class="flex flex-col gap-2">
                                        <p class="text-sm font-semibold text-highlighted">
                                            {{ $d(day, 'date') }}
                                        </p>

                                        <div class="grid gap-1">
                                            <CalendarEntryChip
                                                v-for="entry in dayEntries(day)"
                                                :key="entry.occurrence?.key ?? entry.id"
                                                :entry="entry"
                                                :show-time="true"
                                                @update:popover-open="handleEntryPopoverOpen"
                                                @share="emit('share', $event)"
                                                @delete="emit('delete', $event)"
                                                @select="emit('select', $event)"
                                                @edit="emit('edit', $event)"
                                            />
                                        </div>
                                    </div>
                                </template>
                            </UPopover>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
