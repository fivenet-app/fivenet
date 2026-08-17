<script setup lang="ts">
import { nextTick } from 'vue';
import { addDays, isToday, startOfDay } from 'date-fns';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';
import { dateToDateString } from '~/utils/time';
import { isCalendarEntryPast, getCalendarEntryRangeEnd } from '~/utils/calendar-view';
import { getCalendarEntryDisplayStartDate, sortCalendarEntriesForDisplay } from '~/utils/calendar';
import CalendarEntryChip from './CalendarEntryChip.vue';

const props = defineProps<{
    date: Date;
    entries: CalendarEntry[];
}>();

const emit = defineEmits<{
    (e: 'select', entry: CalendarEntry): void;
    (e: 'edit', entry: CalendarEntry): void;
    (e: 'share', entry: CalendarEntry): void;
    (e: 'delete', entry: CalendarEntry): void;
}>();

const scrollContainer = useTemplateRef('scrollContainer');

const groupedEntries = computed(() => {
    const groups = new Map<string, { date: Date; past: CalendarEntry[]; upcoming: CalendarEntry[] }>();

    for (const entry of props.entries) {
        const start = startOfDay(getCalendarEntryDisplayStartDate(entry));
        const end = getCalendarEntryRangeEnd(entry);

        for (let day = start; day < end; day = addDays(day, 1)) {
            const key = dateToDateString(day);
            const group = groups.get(key) ?? { date: day, past: [], upcoming: [] };

            if (isCalendarEntryPast(entry)) {
                group.past.push(entry);
            } else {
                group.upcoming.push(entry);
            }

            groups.set(key, group);
        }
    }

    return Array.from(groups.values())
        .map((group) => ({
            ...group,
            past: sortCalendarEntriesForDisplay(group.past),
            upcoming: sortCalendarEntriesForDisplay(group.upcoming),
        }))
        .sort((left, right) => left.date.getTime() - right.date.getTime());
});

function scrollToSelectedDate(): void {
    const container = scrollContainer.value;
    if (!container) return;

    const key = dateToDateString(props.date);
    const target = container.querySelector<HTMLElement>(`[data-day-key="${key}"]`);
    if (!target) return;

    target.scrollIntoView({
        block: 'start',
        behavior: 'auto',
    });
}

watch(
    () => [props.date, groupedEntries.value.length],
    async () => {
        await nextTick();
        scrollToSelectedDate();
    },
    { immediate: true },
);
</script>

<template>
    <div ref="scrollContainer" class="overflow-auto border border-default">
        <div
            v-for="group in groupedEntries"
            :key="group.date.toISOString()"
            :data-day-key="dateToDateString(group.date)"
            class="border-b border-default last:border-b-0"
        >
            <div class="flex items-center justify-between gap-2 border-b border-default bg-muted/20 px-4 py-3">
                <div class="flex items-center gap-2">
                    <p class="font-semibold">
                        {{ $d(group.date, 'date') }}
                    </p>
                    <UBadge v-if="isToday(group.date)" size="xs" color="warning" :label="$t('common.today')" />
                </div>
            </div>

            <div class="grid gap-2 p-4">
                <CalendarEntryChip
                    v-for="entry in group.past"
                    :key="entry.occurrence?.key ?? entry.id"
                    :entry="entry"
                    :show-time="true"
                    @share="emit('share', $event)"
                    @delete="emit('delete', $event)"
                    @select="emit('select', $event)"
                    @edit="emit('edit', $event)"
                />

                <div v-if="isToday(group.date)" class="relative my-1 flex items-center gap-2">
                    <div class="h-px flex-1 bg-error/60" />
                    <UBadge size="xs" color="error" variant="soft" :label="$t('common.now')" />
                    <div class="h-px flex-1 bg-error/60" />
                </div>

                <USeparator v-else-if="group.past.length && group.upcoming.length" size="sm" />

                <CalendarEntryChip
                    v-for="entry in group.upcoming"
                    :key="entry.occurrence?.key ?? entry.id"
                    :entry="entry"
                    :show-time="true"
                    @share="emit('share', $event)"
                    @delete="emit('delete', $event)"
                    @select="emit('select', $event)"
                    @edit="emit('edit', $event)"
                />
            </div>
        </div>
    </div>
</template>
