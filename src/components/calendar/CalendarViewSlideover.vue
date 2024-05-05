<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import CalendarCreateOrUpdateModal from './CalendarCreateOrUpdateModal.vue';
import { checkCalendarAccess } from './helpers';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access';
import ConfirmModal from '../partials/ConfirmModal.vue';
import { useCalendarStore } from '~/store/calendar';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';

const props = defineProps<{
    calendarId: string;
}>();

const modal = useModal();
const { isOpen } = useSlideover();

const calendarStore = useCalendarStore();

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`calendar-${props.calendarId}`, () => calendarStore.getCalendar({ calendarId: props.calendarId }));

const calendar = computed(() => data.value?.calendar);
</script>

<template>
    <USlideover :ui="{ width: 'w-full sm:max-w-2xl' }" :overlay="false">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex flex-col gap-1">
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex gap-2 text-2xl font-semibold leading-6">
                            <span>{{ $t('common.calendar') }}: {{ calendar?.name ?? $t('common.calendar', 1) }}</span>

                            <UButton
                                v-if="calendar && checkCalendarAccess(calendar?.access, calendar?.creator, AccessLevel.EDIT)"
                                variant="link"
                                :padded="false"
                                icon="i-mdi-pencil"
                                @click="
                                    modal.open(CalendarCreateOrUpdateModal, {
                                        calendarId: calendar?.id,
                                    })
                                "
                            />

                            <UButton
                                v-if="calendar && checkCalendarAccess(calendar?.access, calendar?.creator, AccessLevel.MANAGE)"
                                variant="link"
                                :padded="false"
                                icon="i-mdi-trash-can"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () => calendarStore.deleteCalendar(calendar?.id!),
                                    })
                                "
                            />
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>

                    <div class="flex flex-row items-center gap-2">
                        <span>{{ $t('common.creator') }}:</span>
                        <CitizenInfoPopover :user="calendar?.creator" show-avatar-in-name />
                    </div>
                </div>
            </template>

            <div>
                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.calendar', 1)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.calendar', 1)])"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!calendar" :type="$t('common.calendar', 1)" icon="i-mdi-comment-text-multiple" />

                <template v-else>
                    <p class="flex-1">
                        {{ $t('common.description') }}:
                        {{ calendar.description ?? $t('common.na') }}
                    </p>
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
    </USlideover>
</template>
