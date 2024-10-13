<script lang="ts" setup>
import type { BadgeColor } from '#ui/types';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCalendarStore } from '~/store/calendar';
import type { ListCalendarsResponse, SubscribeToCalendarResponse } from '~~/gen/ts/services/calendar/calendar';

const { isOpen } = useModal();

const calendarStore = useCalendarStore();
const { currentDate } = storeToRefs(calendarStore);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending: loading, error, refresh } = useLazyAsyncData(`calendars-${page.value}`, () => listCalendars());

async function listCalendars(): Promise<ListCalendarsResponse> {
    const response = await calendarStore.listCalendars({
        pagination: {
            offset: offset.value,
        },
        onlyPublic: true,
    });

    return response;
}

async function subscribeToCalendar(calendarId: string, subscribe: boolean): Promise<SubscribeToCalendarResponse> {
    const call = getGRPCCalendarClient().subscribeToCalendar({
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
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.calendar.FindCalendarsModal.title') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.calendar')])" />
                <DataErrorBlock v-else-if="error" :title="$t('common.not_found', [$t('common.calendar')])" :retry="refresh" />
                <DataNoDataBlock
                    v-else-if="!data?.calendars || data?.calendars.length === 0"
                    :type="`${$t('common.calendar')} ${$t('common.calendar')}`"
                    icon="i-mdi-calendar"
                />

                <ul v-else role="list" class="my-1 flex flex-col divide-y divide-gray-100 dark:divide-gray-800">
                    <li
                        v-for="calendar in data?.calendars"
                        :key="calendar.id"
                        class="hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 flex flex-initial items-center justify-between gap-1 border-white py-1 dark:border-gray-900"
                    >
                        <div class="inline-flex gap-1">
                            <UBadge :color="calendar.color as BadgeColor" :ui="{ rounded: 'rounded-full' }" size="lg" />

                            <span>{{ calendar.name }}</span>
                            <span v-if="calendar.description" class="hidden sm:block"
                                >({{ $t('common.description') }}: {{ calendar.description }})</span
                            >

                            <CitizenInfoPopover v-if="calendar.creator" :user="calendar.creator" />
                        </div>

                        <div>
                            <UButton v-if="calendar.subscription" color="red" @click="subscribeToCalendar(calendar.id, false)">
                                {{ $t('common.unsubscribe') }}
                            </UButton>
                            <UButton v-else color="amber" @click="subscribeToCalendar(calendar.id, true)">
                                {{ $t('common.subscribe') }}
                            </UButton>
                        </div>
                    </li>
                </ul>

                <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
