<script lang="ts" setup>
import type { DateRangeSource } from 'v-calendar/dist/types/src/utils/date/range.js';
import { z } from 'zod';
import MonthCalendarClient from '~/components/partials/MonthCalendar.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import CalendarEntryModal from '~/components/calendar/CalendarEntryModal.vue';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';
import type { ListCalendarEntriesResponse } from '~~/gen/ts/services/calendar/calendar';

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

    return '(' + d(start, 'time') + ' - ' + d(end, 'time') + ')';
}

const calendarEntries = computed(() => {
    return data.value?.entries.map((entry) => {
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
    });
});
</script>

<template>
    <PagesJobsLayout>
        <template #default>
            <DataErrorBlock v-if="error" :retry="async () => $emit('refresh')" />

            <div v-else class="overflow-x-auto">
                <MonthCalendarClient
                    v-model="date"
                    :attributes="calendarEntries"
                    @selected="
                        modal.open(CalendarEntryModal, {
                            entry: $event,
                        })
                    "
                />
            </div>
        </template>
    </PagesJobsLayout>
</template>
