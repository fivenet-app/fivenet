<script lang="ts" setup>
import AccessBadges from '~/components/partials/access/AccessBadges.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import { useCalendarStore } from '~/stores/calendar';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access/access';
import CreateOrUpdateModal from './CreateOrUpdateModal.vue';
import { checkCalendarAccess, isSystemManagedCalendar } from '../helpers';
import CustomContentRenderer from '~/components/partials/content/CustomContentRenderer.vue';

const props = defineProps<{
    calendarId: number;
}>();

const emits = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { can } = useAuth();

const calendarStore = useCalendarStore();

const overlay = useOverlay();

const confirmModal = overlay.create(ConfirmModal);
const calendarCreateOrUpdateModal = overlay.create(CreateOrUpdateModal);

const { data, status, refresh, error } = useLazyAsyncData(`calendar-${props.calendarId}`, () =>
    calendarStore.getCalendar({ calendarId: props.calendarId }),
);

const calendar = computed(() => data.value?.calendar);
const isSystemManaged = computed(() => isSystemManagedCalendar(calendar.value));

const canDo = computed(() => ({
    edit:
        can('calendar.CalendarService/CreateCalendar').value &&
        checkCalendarAccess(
            calendar.value?.access,
            calendar.value?.creator,
            AccessLevel.EDIT,
            calendar.value?.job,
            calendar.value?.creatorJob,
        ),
    manage:
        !isSystemManaged.value &&
        checkCalendarAccess(
            calendar.value?.access,
            calendar.value?.creator,
            AccessLevel.MANAGE,
            calendar.value?.job,
            calendar.value?.creatorJob,
        ),
}));

async function openUpdateModal(): Promise<void> {
    if (!calendar.value) return;

    const response = await calendarCreateOrUpdateModal.open({
        calendarId: calendar.value.id,
        systemManaged: isSystemManaged.value,
    });
    if (response) {
        refresh();
    }
}

async function openDeleteConfirmModal(): Promise<void> {
    if (!calendar.value) return;

    const response = await confirmModal.open({
        confirm: async () => calendar.value && calendarStore.deleteCalendar(calendar.value.id),
    });
    if (response) {
        emits('close', false);
    }
}

defineShortcuts({
    e: () => openUpdateModal(),
    d: () => openDeleteConfirmModal(),
});
</script>

<template>
    <USlideover :title="`${$t('common.calendar')}: ${calendar?.name ?? $t('common.calendar')}`" :overlay="false">
        <template #actions>
            <div v-if="calendar" class="flex items-center justify-between gap-2">
                <UTooltip v-if="canDo.edit" :text="$t('common.edit')" :kbds="['E']">
                    <UButton variant="link" icon="i-mdi-pencil" @click="openUpdateModal" />
                </UTooltip>

                <UTooltip v-if="canDo.manage" :text="$t('common.delete')" :kbds="['D']">
                    <UButton variant="link" icon="i-mdi-delete" color="error" @click="openDeleteConfirmModal" />
                </UTooltip>
            </div>
        </template>

        <template #body>
            <div class="flex h-full w-full flex-1 flex-col gap-2">
                <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.calendar')])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.calendar')])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!calendar" :type="$t('common.calendar')" icon="i-mdi-calendar" />

                <template v-else>
                    <template v-if="!isSystemManaged">
                        <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-2">
                            <OpenClosedBadge :closed="calendar.closed" />

                            <UBadge class="inline-flex gap-1 text-xs" color="neutral" icon="i-mdi-account">
                                <span class="font-medium">{{ $t('common.created_by') }}</span>
                                <CitizenInfoPopover :user="calendar.creator" show-avatar-in-name text-class="text-xs" />
                            </UBadge>

                            <UBadge
                                class="inline-flex gap-1"
                                color="neutral"
                                :icon="calendar.public ? 'i-mdi-public' : 'i-mdi-calendar-lock'"
                                :label="
                                    calendar.public
                                        ? $t('components.calendar.calendar.CreateOrUpdateModal.public')
                                        : $t('components.calendar.calendar.CreateOrUpdateModal.private')
                                "
                            />
                        </div>

                        <div class="mx-auto w-full max-w-(--breakpoint-xl) break-words">
                            <div class="rounded-lg bg-neutral-100 p-4 dark:bg-neutral-800">
                                <CustomContentRenderer :value="calendar.description" :placeholder="$t('common.na')" />
                            </div>
                        </div>

                        <UCollapsible
                            v-if="calendar.access && (calendar.access?.jobs.length > 0 || calendar.access?.users.length > 0)"
                            class="group flex flex-col gap-2"
                        >
                            <UButton
                                class="w-full"
                                color="neutral"
                                variant="subtle"
                                :label="$t('common.access')"
                                icon="i-mdi-lock"
                                trailing-icon="i-mdi-chevron-down"
                                block
                                :ui="{
                                    trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                                }"
                            />

                            <template #content>
                                <AccessBadges
                                    :access-level="AccessLevel"
                                    :jobs="calendar.access.jobs"
                                    :users="calendar.access.users"
                                    i18n-key="enums.calendar"
                                />
                            </template>
                        </UCollapsible>
                    </template>

                    <p v-else class="text-sm text-neutral-500 dark:text-neutral-400">
                        {{ $t('common.read_only') }}
                    </p>
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
