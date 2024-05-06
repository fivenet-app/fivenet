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

const color = computed(() => entry.value?.calendar?.color ?? 'primary');
</script>

<template>
    <USlideover :ui="{ width: 'w-full sm:max-w-2xl' }" :overlay="false">
        <UCard
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex-1 min-h-[calc(100vh-(2*var(--header-height)))] max-h-[calc(100vh-(2*var(--header-height)))] overflow-y-auto',
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
                            <span>{{ entry?.title ?? $t('common.appointment', 1) }}</span>

                            <UButton
                                v-if="entry && checkCalendarAccess(entry?.access, entry?.creator, AccessLevel.EDIT)"
                                variant="link"
                                :padded="false"
                                icon="i-mdi-pencil"
                                @click="
                                    modal.open(EntryCreateOrUpdateModal, {
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
                </div>
            </template>

            <div>
                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.entry', 1)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.entry', 1)])"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!entry" :type="$t('common.entry', 1)" icon="i-mdi-calendar" />

                <template v-else>
                    <p>
                        <span class="font-semibold">{{ $t('common.calendar', 1) }}</span
                        >:
                        <span>
                            <UBadge :color="color" :ui="{ rounded: 'rounded-full' }" label="&nbsp;" />

                            {{ entry.calendar?.name }}
                        </span>
                    </p>

                    <p>
                        <span class="font-semibold">{{ $t('common.date') }}</span
                        >: {{ $d(toDate(entry?.startTime), 'long') }} -
                        {{ $d(toDate(entry?.endTime), 'long') }}
                    </p>

                    <div class="flex flex-row items-center gap-2">
                        <span class="font-semibold">{{ $t('common.creator') }}:</span>
                        <CitizenInfoPopover :user="entry?.creator" show-avatar-in-name />
                    </div>

                    <template v-if="entry.rsvpOpen">
                        <EntryRSVPList :entry-id="entry.id" :rsvp-open="entry.rsvpOpen" />
                    </template>

                    <div class="contentView mx-auto max-w-screen-xl break-words rounded-lg bg-base-900">
                        <!-- eslint-disable vue/no-v-html -->
                        <div class="prose prose-invert min-w-full px-4 py-2" v-html="entry.content"></div>
                    </div>
                </template>

                <UAccordion
                    v-if="!entry?.access && (entry?.access?.jobs.length > 0 || entry?.access?.users.length > 0)"
                    multiple
                    :items="[{ slot: 'access', label: $t('common.access'), icon: 'i-mdi-lock' }]"
                    :unmount="true"
                >
                    <template #access>
                        <UContainer>
                            <div class="flex flex-col gap-2">
                                <div class="flex flex-row flex-wrap gap-1">
                                    <UBadge
                                        v-for="item in entry?.access?.jobs"
                                        :key="item.id"
                                        color="black"
                                        class="inline-flex gap-1"
                                        size="md"
                                    >
                                        <span class="size-2 rounded-full bg-info-500" />
                                        <span>
                                            {{ item.jobLabel
                                            }}<span
                                                v-if="item.minimumGrade > 0"
                                                :title="`${item.jobLabel} - ${$t('common.rank')} ${item.minimumGrade}`"
                                            >
                                                ({{ item.jobGradeLabel }})</span
                                            >
                                            -
                                            {{ $t(`enums.calendar.AccessLevel.${AccessLevel[item.access]}`) }}
                                        </span>
                                    </UBadge>
                                </div>

                                <div class="flex flex-row flex-wrap gap-1">
                                    <UBadge
                                        v-for="item in entry?.access?.users"
                                        :key="item.id"
                                        color="black"
                                        class="inline-flex gap-1"
                                        size="md"
                                    >
                                        <span class="size-2 rounded-full bg-amber-500" />
                                        <span :title="`${$t('common.id')} ${item.userId}`">
                                            {{ item.user?.firstname }}
                                            {{ item.user?.lastname }} -
                                            {{ $t(`enums.calendar.AccessLevel.${AccessLevel[item.access]}`) }}
                                        </span>
                                    </UBadge>
                                </div>
                            </div>
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
