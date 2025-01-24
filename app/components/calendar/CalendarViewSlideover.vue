<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import { useCalendarStore } from '~/store/calendar';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access';
import AccessBadges from '../partials/access/AccessBadges.vue';
import CalendarCreateOrUpdateModal from './CalendarCreateOrUpdateModal.vue';
import { checkCalendarAccess } from './helpers';

const props = defineProps<{
    calendarId: number;
}>();

const { can } = useAuth();

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
    <USlideover :ui="{ width: 'w-screen sm:max-w-2xl' }" :overlay="false">
        <UCard
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex-1 min-h-[calc(100dvh-(2*var(--header-height)))] max-h-[calc(100dvh-(2*var(--header-height)))] overflow-y-auto',
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex flex-col gap-1">
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex gap-2 text-2xl font-semibold leading-6">
                            <span>{{ $t('common.calendar') }}: {{ calendar?.name ?? $t('common.calendar') }}</span>

                            <UButton
                                v-if="
                                    calendar &&
                                    can('CalendarService.CreateOrUpdateCalendar').value &&
                                    checkCalendarAccess(calendar?.access, calendar?.creator, AccessLevel.EDIT)
                                "
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
                                v-if="
                                    calendar &&
                                    can('CalendarService.DeleteCalendar').value &&
                                    checkCalendarAccess(calendar?.access, calendar?.creator, AccessLevel.MANAGE)
                                "
                                variant="link"
                                :padded="false"
                                icon="i-mdi-trash-can"
                                color="red"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () => calendarStore.deleteCalendar(calendar?.id!),
                                    })
                                "
                            />
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </div>
            </template>

            <div>
                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.calendar')])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.calendar')])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!calendar" :type="$t('common.calendar')" icon="i-mdi-comment-text-multiple" />

                <template v-else>
                    <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-2">
                        <OpenClosedBadge :closed="calendar.closed" />

                        <UBadge color="black" class="inline-flex gap-1" size="md">
                            <UIcon name="i-mdi-account" class="size-5" />
                            <span class="inline-flex items-center gap-1">
                                <span class="text-sm font-medium">{{ $t('common.created_by') }}</span>
                                <CitizenInfoPopover :user="calendar.creator" show-avatar-in-name />
                            </span>
                        </UBadge>

                        <UBadge color="black" class="inline-flex gap-1" size="md">
                            <UIcon :name="calendar.public ? 'i-mdi-public' : 'i-mdi-calendar-lock'" class="size-5" />
                            <span>
                                {{
                                    calendar.public
                                        ? $t('components.calendar.CalendarCreateOrUpdateModal.public')
                                        : $t('components.calendar.CalendarCreateOrUpdateModal.private')
                                }}
                            </span>
                        </UBadge>
                    </div>

                    <p>
                        {{ $t('common.description') }}:
                        {{ calendar.description ?? $t('common.na') }}
                    </p>
                </template>

                <UAccordion
                    v-if="calendar?.access && (calendar?.access?.jobs.length > 0 || calendar?.access?.users.length > 0)"
                    multiple
                    :items="[{ slot: 'access', label: $t('common.access'), icon: 'i-mdi-lock' }]"
                    :unmount="true"
                >
                    <template #access>
                        <UContainer>
                            <AccessBadges
                                :access-level="AccessLevel"
                                :jobs="calendar?.access.jobs"
                                :users="calendar?.access.users"
                                i18n-key="enums.calendar"
                            />
                        </UContainer>
                    </template>
                </UAccordion>
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
