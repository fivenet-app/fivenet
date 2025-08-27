<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import { isSameDay } from 'date-fns';
import EntryCreateOrUpdateModal from '~/components/calendar/entry/EntryCreateOrUpdateModal.vue';
import { checkCalendarAccess } from '~/components/calendar/helpers';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import HTMLContent from '~/components/partials/content/HTMLContent.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import { useCalendarStore } from '~/stores/calendar';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import EntryRSVPList from './EntryRSVPList.vue';

const props = defineProps<{
    entryId: number;
}>();

defineEmits<{
    close: [boolean];
}>();

const overlay = useOverlay();

const { can } = useAuth();

const calendarStore = useCalendarStore();
const { calendars } = storeToRefs(calendarStore);

const notifications = useNotificationsStore();

const w = window;

const { data, status, refresh, error } = useLazyAsyncData(`calendar-entry:${props.entryId}`, () =>
    calendarStore.getCalendarEntry({ entryId: props.entryId }),
);

const entry = computed(() => data.value?.entry);

const color = computed(() => (entry.value?.calendar?.color ?? 'primary') as BadgeProps['color']);

function copyLinkToClipboard(): void {
    copyToClipboardWrapper(`${w.location.href}?entry_id=${props.entryId}`);

    notifications.add({
        title: { key: 'notifications.clipboard.link_copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.link_copied.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

const canDo = computed(() => ({
    share: checkCalendarAccess(data.value?.entry?.calendar?.access, entry.value?.creator, AccessLevel.SHARE),
    edit: checkCalendarAccess(data.value?.entry?.calendar?.access, entry.value?.creator, AccessLevel.EDIT),
    manage: checkCalendarAccess(data.value?.entry?.calendar?.access, entry.value?.creator, AccessLevel.MANAGE),
}));

const confirmModal = overlay.create(ConfirmModal);
const entryCreateOrUpdateModal = overlay.create(EntryCreateOrUpdateModal);
</script>

<template>
    <USlideover :title="entry?.title ?? $t('common.appointment', 1)" :overlay="false">
        <template #actions>
            <div class="flex items-center justify-between gap-2">
                <UTooltip
                    v-if="entry && can('calendar.CalendarService/CreateCalendar').value && canDo.edit"
                    :text="$t('common.edit')"
                >
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

                <UTooltip v-if="entry && canDo.manage" :text="$t('common.delete')">
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

                <div class="inline-flex gap-2">
                    <UButton icon="i-mdi-share" @click="copyLinkToClipboard()" />
                </div>
            </div>
        </template>

        <template #body>
            <div class="flex h-full flex-1 flex-col">
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
                        <UBadge class="inline-flex items-center gap-1" color="neutral" size="lg">
                            <UIcon class="size-5" name="i-mdi-access-time" />
                            <span>
                                {{ $t('common.date') }}
                                <GenericTime :value="entry?.startTime" type="long" />
                                <template v-if="entry.endTime">
                                    -
                                    <GenericTime
                                        :value="entry?.endTime"
                                        :type="isSameDay(toDate(entry?.startTime), toDate(entry?.endTime)) ? 'time' : 'long'"
                                    />
                                </template>
                            </span>
                        </UBadge>

                        <UBadge class="inline-flex items-center gap-1" color="neutral" size="md">
                            <UIcon class="size-5" name="i-mdi-calendar" />
                            <span>
                                {{ $t('common.calendar') }}
                                <UBadge :color="color" size="lg" />

                                {{ entry.calendar?.name ?? $t('common.na') }}
                            </span>
                        </UBadge>
                    </div>

                    <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-2">
                        <OpenClosedBadge :closed="entry.closed" />

                        <UBadge class="inline-flex gap-1" color="neutral" size="md">
                            <UIcon class="size-5" name="i-mdi-account" />
                            <span class="inline-flex items-center gap-1">
                                <span class="text-sm font-medium">{{ $t('common.created_by') }}</span>
                                <CitizenInfoPopover :user="entry.creator" show-avatar-in-name />
                            </span>
                        </UBadge>

                        <UBadge class="inline-flex gap-1" color="neutral" size="md">
                            <UIcon class="size-5" name="i-mdi-calendar" />
                            <span>
                                {{ $t('common.created_at') }}
                                <GenericTime :value="entry.createdAt" type="long" />
                            </span>
                        </UBadge>

                        <UBadge v-if="entry.updatedAt" class="inline-flex gap-1" color="neutral" size="md">
                            <UIcon class="size-5" name="i-mdi-calendar-edit" />
                            <span>
                                {{ $t('common.updated_at') }}
                                <GenericTime :value="entry.updatedAt" type="long" />
                            </span>
                        </UBadge>
                    </div>

                    <USeparator />

                    <template v-if="entry.rsvpOpen">
                        <EntryRSVPList
                            v-model="entry.rsvp"
                            :entry-id="entry.id"
                            :rsvp-open="entry.rsvpOpen"
                            :disabled="entry.closed"
                            :show-remove="!calendars.find((c) => c.id === entry?.calendarId)"
                            :can-share="canDo.share"
                        />

                        <USeparator />
                    </template>

                    <div class="dark:bg-base-900 mx-auto w-full max-w-(--breakpoint-xl) rounded-lg bg-neutral-100 break-words">
                        <HTMLContent v-if="entry.content?.content" class="px-4 py-2" :value="entry.content.content" />
                    </div>
                </template>
            </div>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block @click="$emit('close', false)">
                    {{ $t('common.close', 1) }}
                </UButton>
            </UButtonGroup>
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
