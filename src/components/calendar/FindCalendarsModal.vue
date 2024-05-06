<script lang="ts" setup>
import { useCalendarStore } from '~/store/calendar';
import type { ListCalendarsResponse } from '~~/gen/ts/services/calendar/calendar';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';

const { isOpen } = useModal();

const calendarStore = useCalendarStore();

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending: loading, error, refresh } = useLazyAsyncData(`calendars-${page.value}`, () => listCalendars());

async function listCalendars(): Promise<ListCalendarsResponse> {
    try {
        const response = await calendarStore.listCalendars({
            pagination: {
                offset: offset.value,
            },
            onlySubscribed: true,
        });

        return response;
    } catch (e) {
        throw e;
    }
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
                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.calendar', 2)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :message="$t('common.loading', [$t('common.calendar', 2)])"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="!data?.calendars || data?.calendars.length === 0"
                    :type="`${$t('common.calendar')} ${$t('common.calendar', 2)}`"
                    icon="i-mdi-calendar"
                />

                <template v-else>
                    <!-- TODO -->

                    {{ data.calendars }}
                </template>
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
