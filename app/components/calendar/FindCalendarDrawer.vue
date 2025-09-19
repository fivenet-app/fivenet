<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCalendarStore } from '~/stores/calendar';
import { getCalendarCalendarClient } from '~~/gen/ts/clients';
import type { ListCalendarsResponse, SubscribeToCalendarResponse } from '~~/gen/ts/services/calendar/calendar';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const calendarStore = useCalendarStore();
const { currentDate } = storeToRefs(calendarStore);

const calendarCalendarClient = await getCalendarCalendarClient();

const page = useRouteQuery('page', '1', { transform: Number });

const { data, status, error, refresh } = useLazyAsyncData(`calendars-${page.value}`, () => listCalendars());

async function listCalendars(): Promise<ListCalendarsResponse> {
    const response = await calendarStore.listCalendars({
        pagination: {
            offset: calculateOffset(page.value, data.value?.pagination),
        },
        onlyPublic: true,
    });

    return response;
}

async function subscribeToCalendar(calendarId: number, subscribe: boolean): Promise<SubscribeToCalendarResponse> {
    const call = calendarCalendarClient.subscribeToCalendar({
        delete: !subscribe,
        sub: {
            calendarId: calendarId,
            confirmed: true,
            muted: false,
            userId: 0,
        },
    });
    const { response } = await call;

    const calendar = data.value?.calendars.find((c) => c.id === calendarId);
    if (calendar) {
        calendar.subscription = response.sub;

        // Update calendar list and entries if necessary after (un-)subscribing
        await calendarStore.listCalendars({
            pagination: {
                offset: 0,
            },
            onlyPublic: false,
        });

        calendarStore.listCalendarEntries({
            calendarIds: [calendarId],
            year: currentDate.value.year,
            month: currentDate.value.month,
        });
    }

    return response;
}
</script>

<template>
    <UDrawer
        :title="$t('components.calendar.FindCalendarDrawer.title')"
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{ title: 'flex' }"
    >
        <template #title>
            <span>{{ $t('components.calendar.FindCalendarDrawer.title') }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="mx-auto w-full max-w-(--breakpoint-xl)">
                <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.calendar')])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.not_found', [$t('common.calendar')])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="!data?.calendars || data?.calendars.length === 0"
                    :type="`${$t('common.calendar')} ${$t('common.calendar')}`"
                    icon="i-mdi-calendar"
                />

                <ul v-else class="my-1 flex min-h-40 flex-col divide-y divide-default" role="list">
                    <li
                        v-for="calendar in data?.calendars"
                        :key="calendar.id"
                        class="flex flex-initial items-center justify-between gap-1 border-default py-1 hover:border-primary-500/25 hover:bg-primary-100/50 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
                    >
                        <div class="inline-flex gap-1">
                            <UBadge :color="calendar.color as BadgeProps['color']" size="lg" />

                            <span>{{ calendar.name }}</span>
                            <span v-if="calendar.description" class="hidden sm:block"
                                >({{ $t('common.description') }}: {{ calendar.description }})</span
                            >

                            <CitizenInfoPopover v-if="calendar.creator" :user="calendar.creator" />
                        </div>

                        <div>
                            <UButton
                                v-if="calendar.subscription"
                                color="error"
                                @click="subscribeToCalendar(calendar.id, false)"
                            >
                                {{ $t('common.unsubscribe') }}
                            </UButton>
                            <UButton
                                v-else
                                color="warning"
                                @click="
                                    () => {
                                        subscribeToCalendar(calendar.id, true);
                                    }
                                "
                            >
                                {{ $t('common.subscribe') }}
                            </UButton>
                        </div>
                    </li>
                </ul>

                <Pagination v-model="page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
            </div>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UButtonGroup>
        </template>
    </UDrawer>
</template>
