<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import { isSameDay } from 'date-fns';
import EntryCreateOrUpdateModal from '~/components/calendar/entry/EntryCreateOrUpdateModal.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import {
    checkCalendarAccess,
    isBirthdayEntry as isBirthdayCalendarEntry,
    isValidCalendarEntryRecurring,
    isSystemManagedCalendarEntry,
} from '~/components/calendar/helpers';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import CustomContentRenderer from '~/components/partials/content/CustomContentRenderer.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import { copyToClipboardWrapper } from '~/utils/clipboard';
import { useCalendarStore } from '~/stores/calendar';
import { getCalendarEntryDisplayEndDate, getCalendarEntryDisplayStartDate, isCalendarEntryAllDay } from '~/utils/calendar';
import { toDate } from '~/utils/time';
import { useCalendarEntryShortcutState } from '~/composables/useCalendarEntryShortcutState';
import { CalendarEntryRecurringEvery, type CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access/access';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import EntryRSVPList from './EntryRSVPList.vue';
import EntryActionButtons from './EntryActionButtons.vue';
import { emojiBlast } from 'emoji-blast';

const props = defineProps<{
    entryId?: number;
    entry?: CalendarEntry;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const overlay = useOverlay();
const confirmModal = overlay.create(ConfirmModal);
const shortcutState = useCalendarEntryShortcutState();

const { t, d } = useI18n();
const notifications = useNotificationsStore();

const calendarStore = useCalendarStore();
const { calendars } = storeToRefs(calendarStore);

const entryId = props.entry?.id ?? props.entryId;

const {
    data: entry,
    refresh,
    error,
    status,
} = useLazyAsyncData(
    `calendar-entry:${entryId}`,
    async () => {
        if (!entryId) return props.entry;

        return await calendarStore.getCalendarEntry({ entryId });
    },
    {
        default: () => props.entry,
        immediate: !props.entry && !!entryId,
    },
);

const calendarDetails = props.entry
    ? useLazyAsyncData(`calendar-entry-calendar:${props.entry.calendarId}`, () =>
          calendarStore.getCalendar({ calendarId: props.entry!.calendarId }),
      )
    : undefined;

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
const isEntryLoading = computed(() => isRequestPending(status.value) && !entry.value);
const hasEntryError = computed(() => !!error && !entry.value);

const color = computed(() => (calendar.value?.color ?? 'primary') as BadgeProps['color']);

const recurringLabel = computed(() => {
    const recurring = entry.value?.recurring;
    if (!isValidCalendarEntryRecurring(recurring)) return '';

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
const entryCreateOrUpdateModal = overlay.create(EntryCreateOrUpdateModal);

async function openUpdateModal(): Promise<void> {
    if (!entry.value || !canDo.value.edit) return;

    const response = await entryCreateOrUpdateModal.open({
        calendarId: entry.value?.calendarId,
        entryId: entry.value?.id,
    });
    if (response) {
        await refresh();
    }
}

async function copyEntryLink(): Promise<void> {
    if (!entry.value) return;

    const url = new URL(window.location.href);

    if (entry.value.occurrence?.key) {
        url.searchParams.set('entryKey', entry.value.occurrence.key);
        url.searchParams.delete('entryId');
    } else {
        url.searchParams.set('entryId', String(entry.value.id));
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

async function openDeleteEntry(): Promise<void> {
    if (!entry.value || !canDo.value.manage) return;

    const response = await confirmModal.open({
        confirm: async () => calendarStore.deleteCalendarEntry(entry.value!.id),
    });

    if (response) {
        emit('close', false);
    }
}

onMounted(() => {
    shortcutState.isModalOpen.value = true;
});

onUnmounted(() => {
    shortcutState.isModalOpen.value = false;
});
</script>

<template>
    <UModal
        :title="entry?.title ?? $t('common.appointment', 1)"
        :close="false"
        :dismissible="!isEntryLoading"
        :ui="{ content: 'max-w-4xl', body: 'max-h-[80svh] overflow-y-auto' }"
    >
        <template #header>
            <div class="flex w-full items-start justify-between gap-3">
                <div class="min-w-0">
                    <h3 class="truncate text-lg font-semibold text-highlighted">
                        {{ entry?.title ?? $t('common.appointment', 1) }}
                    </h3>

                    <div class="mt-1 flex flex-wrap items-center gap-1 text-sm text-muted">
                        <UIcon name="i-mdi-access-time" class="size-4" />
                        <GenericTime :value="displayStartTime" :type="isCalendarEntryAllDay(entry) ? 'date' : 'long'" />
                        <template v-if="displayEndTime && !isCalendarEntryAllDay(entry)">
                            <span>-</span>
                            <GenericTime
                                :value="displayEndTime"
                                :type="isSameDay(displayStartTime, displayEndTime) ? 'time' : 'long'"
                            />
                        </template>
                        <template v-else-if="displayEndTime && !isSameDay(displayStartTime, displayEndTime)">
                            <span>-</span>
                            <GenericTime :value="displayEndTime" type="date" />
                        </template>
                    </div>
                </div>

                <div class="flex items-center gap-1">
                    <UBadge v-if="calendar?.name" :color="color" variant="subtle" size="sm">
                        {{ calendar.name }}
                    </UBadge>

                    <EntryActionButtons
                        v-if="entry"
                        :entry="entry"
                        mode="header"
                        :can-edit="canDo.edit"
                        :can-share="canDo.share"
                        :can-delete="canDo.manage"
                        @edit="openUpdateModal"
                        @share="copyEntryLink"
                        @delete="openDeleteEntry"
                    />

                    <UButton
                        color="neutral"
                        variant="ghost"
                        icon="i-mdi-close"
                        :aria-label="$t('common.close', 1)"
                        @click="emit('close', false)"
                    />
                </div>
            </div>
        </template>

        <template #body>
            <div class="flex w-full flex-1 flex-col gap-2">
                <DataPendingBlock v-if="isEntryLoading" :message="$t('common.loading', [$t('common.entry', 1)])" />
                <DataErrorBlock
                    v-else-if="hasEntryError"
                    :title="$t('common.unable_to_load', [$t('common.entry', 1)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!entry" :type="$t('common.entry', 1)" icon="i-mdi-calendar" />

                <template v-else>
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
                            <CustomContentRenderer :value="entry.content" :placeholder="$t('common.na')" />
                        </div>
                    </div>
                </template>
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="emit('close', false)" />
            </UFieldGroup>
        </template>
    </UModal>
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
