<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import EntryRSVPList from './EntryRSVPList.vue';
import { useCalendarStore } from '~/store/calendar';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import EntryCreateOrUpdateModal from './EntryCreateOrUpdateModal.vue';
import { checkCalendarAccess } from '../helpers';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';

const props = defineProps<{
    calendarId: string;
    entryId: string;
}>();

const modal = useModal();
const { isOpen } = useSlideover();

const calendarStore = useCalendarStore();

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`calendar-${props.calendarId}-entry:${props.entryId}`, () =>
    calendarStore.getCalendarEntry({ calendarId: props.calendarId, entryId: props.entryId }),
);

const entry = computed(() => data.value?.entry);
</script>

<template>
    <USlideover :ui="{ width: 'w-full sm:max-w-2xl' }" :overlay="false">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex flex-col gap-1">
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex gap-2 text-2xl font-semibold leading-6">
                            <span>{{ entry?.title ?? $t('common.appointment', 1) }}</span>

                            <UButton
                                v-if="entry && checkCalendarAccess(entry?.access, entry?.creator, AccessLevel.EDIT)"
                                variant="link"
                                :padded="false"
                                icon="i-mdi-pencil"
                                @click="
                                    modal.open(EntryCreateOrUpdateModal, {
                                        calendar: entry?.calendar,
                                        calendarId: entry?.calendarId,
                                        entryId: entry?.id,
                                    })
                                "
                            />

                            <UButton
                                v-if="entry && checkCalendarAccess(entry?.access, entry?.creator, AccessLevel.MANAGE)"
                                variant="link"
                                :padded="false"
                                icon="i-mdi-trash-can"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () => calendarStore.deleteCalendarEntry(entry?.calendarId!, entry?.id!),
                                    })
                                "
                            />
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>

                    <p class="flex-1">
                        {{ $d(toDate(entry?.startTime), 'long') }} -
                        {{ $d(toDate(entry?.endTime), 'long') }}
                    </p>

                    <div class="flex flex-row items-center gap-2">
                        <span>{{ $t('common.creator') }}:</span>
                        <CitizenInfoPopover :user="entry?.creator" show-avatar-in-name />
                    </div>
                </div>
            </template>

            <div>
                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.entry', 1)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.entry', 1)])"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!entry" :type="$t('common.entry', 1)" icon="i-mdi-comment-text-multiple" />

                <template v-else>
                    <p v-html="entry.content"></p>

                    <template v-if="entry.rsvpOpen">
                        <UDivider class="mb-2 mt-2" />

                        <EntryRSVPList :entry-id="entry.id" :rsvp-open="entry.rsvpOpen" />
                    </template>
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
