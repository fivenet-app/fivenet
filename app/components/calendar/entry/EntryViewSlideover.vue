<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import { isSameDay } from 'date-fns';
import EntryCreateOrUpdateModal from '~/components/calendar/entry/EntryCreateOrUpdateModal.vue';
import {
    checkCalendarAccess,
    isBirthdayEntry as isBirthdayCalendarEntry,
    isSystemManagedCalendarEntry,
} from '~/components/calendar/helpers';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import CustomContentRenderer from '~/components/partials/content/CustomContentRenderer.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import { useCalendarStore } from '~/stores/calendar';
import { getCalendarEntryDisplayEndDate, getCalendarEntryDisplayStartDate } from '~/utils/calendar';
import { toDate } from '~/utils/time';
import { CalendarEntryRecurringEvery, type CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access/access';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import EntryRSVPList from './EntryRSVPList.vue';
import { emojiBlast } from 'emoji-blast';

const props = defineProps<{
    entryId?: number;
    entry?: CalendarEntry;
}>();

defineEmits<{
    close: [boolean];
}>();

const overlay = useOverlay();

const { can } = useAuth();
const { t, d } = useI18n();

const calendarStore = useCalendarStore();
const { calendars } = storeToRefs(calendarStore);

const notifications = useNotificationsStore();

const w = window;

const fetched = props.entry
    ? undefined
    : useLazyAsyncData(`calendar-entry:${props.entryId}`, () => calendarStore.getCalendarEntry({ entryId: props.entryId! }));

const calendarDetails = props.entry
    ? useLazyAsyncData(`calendar-entry-calendar:${props.entry.calendarId}`, () =>
          calendarStore.getCalendar({ calendarId: props.entry!.calendarId }),
      )
    : undefined;

const entry = computed(() => props.entry ?? fetched?.data.value?.entry);
const status = computed(() => (props.entry ? 'success' : (fetched?.status.value ?? 'idle')));
const error = computed(() => (props.entry ? undefined : fetched?.error.value));
const refresh = async (): Promise<void> => {
    if (fetched) await fetched.refresh();
};

const calendar = computed(
    () =>
        calendarDetails?.data.value?.calendar ??
        entry.value?.calendar ??
        calendars.value.find((candidate) => candidate.id === entry.value?.calendarId),
);
const isSystemManaged = computed(() => isSystemManagedCalendarEntry(calendar.value, entry.value));
const isBirthdayEntry = computed(() => isBirthdayCalendarEntry(entry.value));

const displayStartTime = computed(() => (entry.value ? getCalendarEntryDisplayStartDate(entry.value) : new Date()));
const displayEndTime = computed(() => (entry.value ? getCalendarEntryDisplayEndDate(entry.value) : undefined));

const color = computed(() => (calendar.value?.color ?? 'primary') as BadgeProps['color']);

const recurringLabel = computed(() => {
    const recurring = entry.value?.recurring;
    if (!recurring || recurring.every === CalendarEntryRecurringEvery.UNSPECIFIED) return '';

    const everyUnit = (() => {
        switch (recurring.every) {
            case CalendarEntryRecurringEvery.DAY:
                return t('common.time_ago.day', recurring.count);
            case CalendarEntryRecurringEvery.WEEK:
                return t('common.time_ago.week', recurring.count);
            case CalendarEntryRecurringEvery.MONTH:
                return t('common.time_ago.month', recurring.count);
            case CalendarEntryRecurringEvery.YEAR:
                return t('common.time_ago.year', recurring.count);
            default:
                return '';
        }
    })();

    const until = recurring.until
        ? ` · ${t('components.calendar.EntryCreateOrUpdateModal.recurring.until')} ${d(toDate(recurring.until), 'date')}`
        : '';

    return `${t('components.calendar.EntryCreateOrUpdateModal.recurring.every')} ${recurring.count} ${everyUnit}${until}`;
});

function copyLinkToClipboard(): void {
    const url = new URL(w.location.href);

    if (entry.value?.occurrence?.key) {
        url.searchParams.set('entryKey', entry.value.occurrence.key);
        url.searchParams.delete('entryId');
    } else {
        url.searchParams.set('entryId', String(props.entryId ?? entry.value?.id ?? 0));
        url.searchParams.delete('entryKey');
    }

    copyToClipboardWrapper(url.toString());

    notifications.add({
        title: { key: 'notifications.clipboard.link_copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.link_copied.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

const canDo = computed(() => ({
    share:
        !isSystemManaged.value &&
        checkCalendarAccess(
            calendar.value?.access,
            entry.value?.creator,
            AccessLevel.SHARE,
            calendar.value?.job,
            calendar.value?.creatorJob,
        ),
    edit:
        !isSystemManaged.value &&
        checkCalendarAccess(
            calendar.value?.access,
            entry.value?.creator,
            AccessLevel.EDIT,
            calendar.value?.job,
            calendar.value?.creatorJob,
        ),
    manage:
        !isSystemManaged.value &&
        checkCalendarAccess(
            calendar.value?.access,
            entry.value?.creator,
            AccessLevel.MANAGE,
            calendar.value?.job,
            calendar.value?.creatorJob,
        ),
}));

const confirmModal = overlay.create(ConfirmModal);
const entryCreateOrUpdateModal = overlay.create(EntryCreateOrUpdateModal);
</script>

<template>
    <USlideover :title="entry?.title ?? $t('common.appointment', 1)" :overlay="false">
        <template #actions>
            <div v-if="entry" class="flex items-center justify-between gap-2">
                <UTooltip v-if="can('calendar.CalendarService/CreateCalendar').value && canDo.edit" :text="$t('common.edit')">
                    <UButton
                        variant="link"
                        icon="i-mdi-pencil"
                        @click="
                            entryCreateOrUpdateModal.open({
                                calendarId: entry?.calendarId,
                                entryId: entry?.id,
                            })
                        "
                    />
                </UTooltip>

                <UTooltip v-if="canDo.manage" :text="$t('common.delete')">
                    <UButton
                        variant="link"
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            confirmModal.open({
                                confirm: async () => calendarStore.deleteCalendarEntry(entry?.id!),
                            })
                        "
                    />
                </UTooltip>

                <UTooltip :text="$t('common.share')">
                    <UButton variant="link" icon="i-mdi-share" @click="copyLinkToClipboard()" />
                </UTooltip>
            </div>
        </template>

        <template #body>
            <div class="flex h-full w-full flex-1 flex-col gap-2">
                <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.entry', 1)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.entry', 1)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!entry" :type="$t('common.entry', 1)" icon="i-mdi-calendar" />

                <template v-else>
                    <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-2">
                        <UBadge class="inline-flex items-center gap-1" color="neutral" size="lg" icon="i-mdi-access-time">
                            {{ $t('common.date') }}
                            <GenericTime :value="displayStartTime" :type="entry?.occurrence?.allDay ? 'date' : 'long'" />
                            <template v-if="displayEndTime && !entry?.occurrence?.allDay">
                                -
                                <GenericTime
                                    :value="displayEndTime"
                                    :type="isSameDay(displayStartTime, displayEndTime) ? 'time' : 'long'"
                                />
                            </template>
                        </UBadge>

                        <UBadge class="inline-flex items-center gap-1" color="neutral" icon="i-mdi-calendar">
                            {{ $t('common.calendar') }}
                            <UBadge :color="color" size="lg" />

                            {{ entry.calendar?.name ?? $t('common.na') }}
                        </UBadge>
                    </div>

                    <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-2">
                        <OpenClosedBadge v-if="!isSystemManaged" :closed="entry.closed" />

                        <UBadge v-if="recurringLabel" class="inline-flex gap-1" color="neutral" icon="i-mdi-repeat">
                            {{ recurringLabel }}
                        </UBadge>

                        <UBadge
                            v-if="!isSystemManaged && entry.creator"
                            class="inline-flex gap-1"
                            color="neutral"
                            icon="i-mdi-account"
                        >
                            <span>{{ $t('common.created_by') }}</span>
                            <CitizenInfoPopover :user="entry.creator" :show-avatar-in-name="false" text-class="text-xs" />
                        </UBadge>

                        <UBadge
                            v-else-if="!isSystemManaged"
                            class="inline-flex gap-1"
                            color="neutral"
                            icon="i-mdi-cog"
                            :label="$t('components.calendar.system_generated_entry')"
                        />

                        <template v-if="!isSystemManaged">
                            <UBadge class="inline-flex gap-1" color="neutral" icon="i-mdi-calendar">
                                {{ $t('common.created_at') }}
                                <GenericTime :value="entry.createdAt" type="long" />
                            </UBadge>

                            <UBadge v-if="entry.updatedAt" class="inline-flex gap-1" color="neutral" icon="i-mdi-calendar-edit">
                                {{ $t('common.updated_at') }}
                                <GenericTime :value="entry.updatedAt" type="long" />
                            </UBadge>
                        </template>
                    </div>

                    <USeparator />

                    <template v-if="entry.rsvpOpen">
                        <EntryRSVPList
                            v-model="entry.rsvp"
                            :entry-id="entry.id"
                            :occurrence-key="entry.occurrence?.key"
                            :rsvp-open="entry.rsvpOpen"
                            :disabled="entry.closed"
                            :show-remove="!calendars.find((c) => c.id === entry?.calendarId)"
                            :can-share="canDo.share"
                        />

                        <USeparator />
                    </template>

                    <div class="mx-auto w-full max-w-(--breakpoint-xl) break-words">
                        <UAlert
                            v-if="isBirthdayEntry"
                            class="rounded-lg"
                            icon="i-mdi-birthday-cake"
                            variant="subtle"
                            :title="$t('components.calendar.birthday_entry_block.title')"
                            :description="$t('components.calendar.birthday_entry_block.description')"
                            :actions="[
                                {
                                    label: $t('components.calendar.birthday_entry_action'),
                                    icon: 'i-mdi-party-popper',
                                    variant: 'subtle',
                                    onClick: () => {
                                        emojiBlast({
                                            emojis: ['🎂', '🎁', '🍰', '🎈', '🎉', '🥳', '🎊', '✨'],
                                        });
                                    },
                                },
                            ]"
                            :ui="{ icon: 'size-6' }"
                        />
                        <div v-else class="rounded-lg bg-neutral-100 p-4 dark:bg-neutral-800">
                            <CustomContentRenderer v-if="entry.content" :value="entry.content" />
                            <p v-else>{{ $t('common.na') }}</p>
                        </div>
                    </div>
                </template>
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UFieldGroup>
        </template>
    </USlideover>
</template>

<style scoped>
.contentView:deep(.prose) {
    * {
        margin-top: 4px;
        margin-bottom: 4px;
    }

    input[type='checkbox']:checked {
        opacity: 1;
    }
}
</style>
