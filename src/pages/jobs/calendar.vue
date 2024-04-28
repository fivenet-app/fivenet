<script lang="ts" setup>
import type { DateRangeSource } from 'v-calendar/dist/types/src/utils/date/range.js';
import { z } from 'zod';
import MonthCalendarClient from '~/components/partials/MonthCalendar.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import CalendarEntryModal from '~/components/calendar/CalendarEntryModal.vue';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';
import type { ListCalendarEntriesResponse } from '~~/gen/ts/services/calendar/calendar';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';

const { $grpc } = useNuxtApp();

const { d } = useI18n();

const modal = useModal();

const schema = z.object({
    year: z.number(),
    month: z.number(),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    year: 2024,
    month: 4,
});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`calendar-${query.year}-${query.month}`, () => listCalendarEntries());

async function listCalendarEntries(): Promise<ListCalendarEntriesResponse> {
    try {
        const call = $grpc.getCalendarClient().listCalendarEntries({
            year: 2024,
            month: 4,
            calendarIds: [],
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

const calendarEntries = computed(() =>
    data.value?.entries.map((entry) => {
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
    calendarEntries.value?.forEach((entry) => {
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
            <DataErrorBlock v-if="error" :retry="refresh" />

            <div v-else class="overflow-x-auto">
                <MonthCalendarClient
                    v-model="date"
                    class="hidden md:flex md:flex-1"
                    :attributes="calendarEntries"
                    @selected="
                        modal.open(CalendarEntryModal, {
                            entry: $event,
                        })
                    "
                />

                <UContainer class="flex flex-1 flex-col md:hidden">
                    <DataNoDataBlock
                        v-if="!groupedCalendarEntries || groupedCalendarEntries.length === 0"
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
                                        <span>
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
