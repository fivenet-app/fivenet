<script lang="ts" setup>
import type { DateRangeSource } from 'v-calendar/dist/types/src/utils/date/range.js';
import { z } from 'zod';
import MonthCalendarClient from '~/components/partials/MonthCalendar.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';
import type { ListCalendarsResponse } from '~~/gen/ts/services/calendar/calendar';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import EntryViewSlideover from '~/components/calendar/entry/EntryViewSlideover.vue';
import EntryCreateOrUpdateModal from '~/components/calendar/entry/EntryCreateOrUpdateModal.vue';
import { useCalendarStore } from '~/store/calendar';
import CalendarCreateOrUpdateModal from '~/components/calendar/CalendarCreateOrUpdateModal.vue';
import CalendarViewSlideover from '~/components/calendar/CalendarViewSlideover.vue';

useHead({
    title: 'common.calendar',
});
definePageMeta({
    title: 'common.calendar',
    requiresAuth: true,
});

const { d } = useI18n();

const modal = useModal();
const slideover = useSlideover();

const calendarStore = useCalendarStore();

const schema = z.object({
    year: z.number(),
    month: z.number(),
    calendarIds: z.string().array().max(25),
});

type Schema = z.output<typeof schema>;

const date = ref(new Date());
watch(date, () => {
    query.year = date.value.getFullYear();
    query.month = date.value.getMonth();
});

const query = reactive<Schema>({
    year: date.value.getFullYear(),
    month: date.value.getMonth(),
    calendarIds: [],
});

const page = ref(1);
const offset = computed(() =>
    calendars.value?.pagination?.pageSize ? calendars.value?.pagination?.pageSize * (page.value - 1) : 0,
);

const {
    data: calendars,
    pending: calendarsLoading,
    error: calendarsError,
} = useLazyAsyncData(`calendars-${query.year}-${query.month}`, () => listCalendars());

async function listCalendars(): Promise<ListCalendarsResponse> {
    try {
        const response = await calendarStore.listCalendars({
            pagination: {
                offset: offset.value,
            },
        });

        if (query.calendarIds.length === 0) {
            query.calendarIds = response.calendars.map((c) => c.id);
        }

        refresh();

        return response;
    } catch (e) {
        throw e;
    }
}

const {
    data: calendarEntries,
    refresh,
    error,
} = useLazyAsyncData(
    `calendar-entries-${query.year}-${query.month}-${query.calendarIds.join(':')}`,
    () =>
        calendarStore.listCalendarEntries({
            year: query.year,
            month: query.month,
            calendarIds: query.calendarIds,
        }),
    { immediate: false },
);

watchDebounced(query, () => refresh(), { debounce: 200, maxWait: 1250 });

function formatStartEndTime(entry: CalendarEntry): string {
    const start = toDate(entry.startTime);
    const end = entry.endTime ? toDate(entry.endTime) : undefined;

    if (!end) {
        return d(start, 'time');
    }

    return d(start, 'time') + ' - ' + d(end, 'time');
}

type CalEntry = { key: string; customData: CalendarEntry & { color: string; time: string }; dates: DateRangeSource[] };

const transformedCalendarEntries = computed(() =>
    calendarEntries.value?.entries.map((entry) => {
        const color = calendars.value?.calendars.find((c) => c.id === entry.calendarId)?.color ?? 'primary';
        return {
            key: entry.id,
            customData: {
                ...entry,
                color: color,
                time: formatStartEndTime(entry),
            },
            dates: [
                {
                    start: toDate(entry.startTime),
                    end: entry.endTime ? toDate(entry.endTime) : undefined,
                },
            ] as DateRangeSource[],
        };
    }),
);

type GroupedCalendarEntries = { key: string; date: Date; entries: CalEntry[] }[];

const groupedCalendarEntries = computed(() => {
    const groups: GroupedCalendarEntries = [];
    transformedCalendarEntries.value?.forEach((entry) => {
        const date = toDate(entry.customData.startTime);
        const idx = groups.findIndex((g) => g.key === toDate(entry.customData.startTime).toDateString());
        if (idx === -1) {
            groups.push({
                key: date.toDateString(),
                date: date,
                entries: [entry],
            });
        } else {
            groups[idx].entries.push(entry);
        }
    });

    return groups;
});

function calendarIdChange(calendarId: string, state: boolean): void {
    if (state) {
        if (!query.calendarIds.includes(calendarId)) {
            query.calendarIds.push(calendarId);
        }
    } else {
        query.calendarIds = query.calendarIds.filter((cId) => cId !== calendarId);
    }
}
</script>

<template>
    <PagesJobsLayout>
        <template #default>
            <UDashboardToolbar
                v-if="can('CalendarService.CreateOrUpdateCalendarEntry') || can('CalendarService.CreateOrUpdateCalendar')"
                class="flex w-full flex-row gap-2 p-2"
            >
                <template #default>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton
                            v-if="can('CalendarService.CreateOrUpdateCalendar')"
                            block
                            color="gray"
                            trailing-icon="i-mdi-plus"
                            class="flex-1"
                            @click="modal.open(CalendarCreateOrUpdateModal, {})"
                        >
                            {{ $t('common.calendar', 1) }}
                            {{ $t('common.create') }}
                        </UButton>

                        <UButton
                            v-if="can('CalendarService.CreateOrUpdateCalendarEntry')"
                            block
                            color="gray"
                            trailing-icon="i-mdi-plus"
                            class="flex-1"
                            @click="modal.open(EntryCreateOrUpdateModal, {})"
                        >
                            {{ $t('common.entry', 1) }}
                            {{ $t('common.create') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UDashboardToolbar>

            <UContainer :ui="{ constrained: 'max-w-5xl' }" class="mt-2 w-full">
                <UAccordion :items="[{ slot: 'calendar', label: $t('common.calendar', 2), icon: 'i-mdi-calendar' }]">
                    <template #calendar>
                        <div>
                            <DataPendingBlock
                                v-if="calendarsLoading"
                                :message="$t('common.loading', [$t('common.calendar', 2)])"
                            />
                            <DataErrorBlock
                                v-else-if="calendarsError"
                                :message="$t('common.loading', [$t('common.calendar', 2)])"
                                :retry="refresh"
                            />
                            <template v-else>
                                <div class="grid grid-cols-1 gap-2 sm:grid-cols-2 lg:grid-cols-3">
                                    <div
                                        v-for="calendar in calendars?.calendars"
                                        :key="calendar.id"
                                        class="inline-flex items-center gap-2"
                                    >
                                        <UCheckbox
                                            :model-value="query.calendarIds.includes(calendar.id)"
                                            class="truncate"
                                            @change="calendarIdChange(calendar.id, $event)"
                                        />
                                        <UButton
                                            :color="calendar.color"
                                            size="sm"
                                            truncate
                                            :label="calendar.name"
                                            @click="slideover.open(CalendarViewSlideover, { calendarId: calendar.id })"
                                        />
                                    </div>
                                </div>
                            </template>
                        </div>
                    </template>
                </UAccordion>
            </UContainer>

            <DataErrorBlock v-if="error" :retry="refresh" />

            <div v-else class="overflow-x-auto">
                <MonthCalendarClient
                    v-model="date"
                    class="hidden md:flex md:flex-1"
                    :attributes="transformedCalendarEntries"
                    @selected="
                        slideover.open(EntryViewSlideover, {
                            entryId: $event.id,
                            calendarId: $event.calendarId,
                        })
                    "
                />

                <UContainer class="flex flex-1 flex-col md:hidden">
                    <DataErrorBlock v-if="error" :message="$t('common.loading', [$t('common.entry', 2)])" :retry="refresh" />
                    <DataNoDataBlock
                        v-else-if="!groupedCalendarEntries || groupedCalendarEntries.length === 0"
                        :type="`${$t('common.calendar')} ${$t('common.entry', 2)}`"
                        icon="i-mdi-calendar"
                    />

                    <template v-else>
                        <template v-for="entries in groupedCalendarEntries" :key="entries.key">
                            <UDivider class="text-lg font-semibold">{{ $d(entries.date, 'date') }}</UDivider>

                            <ul role="list">
                                <li v-for="attr in entries.entries" :key="attr.key">
                                    <ULink
                                        class="inline-flex w-full items-center justify-between gap-1"
                                        @click="
                                            slideover.open(EntryViewSlideover, {
                                                entryId: attr.customData.id,
                                                calendarId: attr.customData.calendarId,
                                            })
                                        "
                                    >
                                        <span class="inline-flex items-center gap-1">
                                            <UBadge
                                                :color="attr.customData.color"
                                                :ui="{ rounded: 'rounded-full' }"
                                                label="&nbsp;"
                                            />

                                            <template v-if="attr.customData.time">
                                                {{ attr.customData.time }}
                                            </template>
                                            <span>-</span>

                                            {{ attr.customData.title }}
                                        </span>

                                        <UButton :padded="false" variant="link" icon="i-mdi-eye" />
                                    </ULink>
                                </li>
                            </ul>
                        </template>
                    </template>
                </UContainer>
            </div>
        </template>
    </PagesJobsLayout>
</template>
