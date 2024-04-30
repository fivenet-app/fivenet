<script lang="ts" setup>
import type { DateRangeSource } from 'v-calendar/dist/types/src/utils/date/range.js';
import { z } from 'zod';
import MonthCalendarClient from '~/components/partials/MonthCalendar.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import CalendarEntryModal from '~/components/calendar/CalendarEntryModal.vue';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';
import type { ListCalendarEntriesResponse, ListCalendarsResponse } from '~~/gen/ts/services/calendar/calendar';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';

useHead({
    title: 'common.calendar',
});
definePageMeta({
    title: 'common.calendar',
    requiresAuth: true,
});

const { $grpc } = useNuxtApp();

const { d } = useI18n();

const modal = useModal();

const schema = z.object({
    year: z.number(),
    month: z.number(),
    calendarIds: z.string().array().max(25),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    year: 2024,
    month: 4,
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
        const call = $grpc.getCalendarClient().listCalendars({
            pagination: {
                offset: offset.value,
            },
        });
        const { response } = await call;

        if (query.calendarIds.length === 0) {
            query.calendarIds = response.calendars.map((c) => c.id);
        }

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const {
    data: calendarEntries,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`calendar-entries-${query.year}-${query.month}`, () => listCalendarEntries(), { immediate: false });

async function listCalendarEntries(): Promise<ListCalendarEntriesResponse> {
    try {
        const call = $grpc.getCalendarClient().listCalendarEntries({
            year: 2024,
            month: 4,
            calendarIds: query.calendarIds,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const date = ref(new Date());

watchDebounced(query, () => refresh(), { debounce: 200, maxWait: 1250 });

function formatStartEndTime(entry: CalendarEntry): string {
    const start = toDate(entry.startTime);
    const end = entry.endTime ? toDate(entry.endTime) : undefined;

    if (!end) {
        return '';
    }

    return d(start, 'time') + ' - ' + d(end, 'time');
}

type CalEntry = { key: string; customData: CalendarEntry & { class: string; time: string }; dates: DateRangeSource[] };

const transformedCalendarEntries = computed(() =>
    calendarEntries.value?.entries.map((entry) => {
        return {
            key: entry.id,
            customData: { ...entry, class: 'bg-blue-500 hover:!bg-blue-400 text-white', time: formatStartEndTime(entry) },
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
</script>

<template>
    <PagesJobsLayout>
        <template #default>
            <UContainer :ui="{ constrained: 'max-w-5xl' }" class="w-full p-2">
                <UAccordion class="" :items="[{ slot: 'calendar', label: $t('common.calendar', 2), icon: 'i-mdi-calendar' }]">
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
                                    <UTooltip
                                        v-for="calendar in calendars?.calendars"
                                        :key="calendar.id"
                                        :text="calendar.description"
                                    >
                                        <UCheckbox :label="calendar.name" disabled :model-value="true" class="truncate" />
                                    </UTooltip>
                                </div>
                            </template>
                        </div>
                    </template>
                </UAccordion>
            </UContainer>

            <DataPendingBlock v-if="loading || calendarsLoading" :message="$t('common.loading', [$t('common.calendar', 2)])" />
            <DataErrorBlock v-else-if="error" :retry="refresh" />

            <div v-else class="overflow-x-auto">
                <MonthCalendarClient
                    v-model="date"
                    class="hidden md:flex md:flex-1"
                    :attributes="transformedCalendarEntries"
                    @selected="
                        modal.open(CalendarEntryModal, {
                            entry: $event,
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

                            <ul role="list" class="list-disc">
                                <li v-for="attr in entries.entries" :key="attr.key">
                                    <ULink
                                        class="inline-flex w-full items-center justify-between gap-1"
                                        @click="
                                            modal.open(CalendarEntryModal, {
                                                entry: attr.customData,
                                            })
                                        "
                                    >
                                        <span class="inline-flex items-center gap-1">
                                            <UBadge :class="attr.customData.class" label="&nbsp;" />

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
