<script lang="ts" setup>
import AccessBadges from '~/components/partials/access/AccessBadges.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import { useCalendarStore } from '~/stores/calendar';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access';
import CalendarCreateOrUpdateModal from './CalendarCreateOrUpdateModal.vue';
import { checkCalendarAccess } from './helpers';

const props = defineProps<{
    calendarId: number;
}>();

defineEmits<{
    close: [boolean];
}>();

const { can } = useAuth();

const calendarStore = useCalendarStore();

const overlay = useOverlay();

const confirmModal = overlay.create(ConfirmModal);
const calendarCreateOrUpdateModal = overlay.create(CalendarCreateOrUpdateModal);

const { data, status, refresh, error } = useLazyAsyncData(`calendar-${props.calendarId}`, () =>
    calendarStore.getCalendar({ calendarId: props.calendarId }),
);

const calendar = computed(() => data.value?.calendar);
</script>

<template>
    <USlideover :title="`${$t('common.calendar')}: ${calendar?.name ?? $t('common.calendar')}`" :overlay="false">
        <template #actions>
            <div class="flex items-center justify-between gap-2">
                <UTooltip
                    v-if="
                        calendar &&
                        can('calendar.CalendarService/CreateCalendar').value &&
                        checkCalendarAccess(calendar?.access, calendar?.creator, AccessLevel.EDIT)
                    "
                    :text="$t('common.edit')"
                >
                    <UButton
                        variant="link"
                        icon="i-mdi-pencil"
                        @click="
                            calendarCreateOrUpdateModal.open({
                                calendarId: calendar?.id,
                            })
                        "
                    />
                </UTooltip>

                <UTooltip
                    v-if="calendar && checkCalendarAccess(calendar?.access, calendar?.creator, AccessLevel.MANAGE)"
                    :text="$t('common.delete')"
                >
                    <UButton
                        variant="link"
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            confirmModal.open({
                                confirm: async () => calendarStore.deleteCalendar(calendar?.id!),
                            })
                        "
                    />
                </UTooltip>
            </div>
        </template>

        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.calendar')])" />
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

                    <UBadge class="inline-flex gap-1" color="neutral" icon="i-mdi-account">
                        <span class="text-sm font-medium">{{ $t('common.created_by') }}</span>
                        <CitizenInfoPopover :user="calendar.creator" show-avatar-in-name />
                    </UBadge>

                    <UBadge
                        class="inline-flex gap-1"
                        color="neutral"
                        :icon="calendar.public ? 'i-mdi-public' : 'i-mdi-calendar-lock'"
                        :label="
                            calendar.public
                                ? $t('components.calendar.CalendarCreateOrUpdateModal.public')
                                : $t('components.calendar.CalendarCreateOrUpdateModal.private')
                        "
                    />
                </div>

                <p>
                    <span class="font-semibold">{{ $t('common.description') }}:</span>
                    {{
                        calendar.description === undefined || calendar.description === ''
                            ? $t('common.na')
                            : calendar.description
                    }}
                </p>
            </template>

            <UCollapsible
                v-if="calendar?.access && (calendar?.access?.jobs.length > 0 || calendar?.access?.users.length > 0)"
                class="group flex flex-col gap-2"
            >
                <UButton
                    color="neutral"
                    variant="subtle"
                    :label="$t('common.access')"
                    icon="i-mdi-lock"
                    class="w-full"
                    trailing-icon="i-mdi-chevron-down"
                    block
                    :ui="{
                        trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                    }"
                />

                <template #content>
                    <AccessBadges
                        :access-level="AccessLevel"
                        :jobs="calendar?.access.jobs"
                        :users="calendar?.access.users"
                        i18n-key="enums.calendar"
                    />
                </template>
            </UCollapsible>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UButtonGroup>
        </template>
    </USlideover>
</template>
