<script setup lang="ts">
import { format, isSameDay } from 'date-fns';
import type { ReferenceElement } from '@floating-ui/vue';
import type { ButtonProps } from '@nuxt/ui';
import { getCalendarEntryColor, getCalendarEntryIcon, getCalendarEntryTimeLabel } from '~/utils/calendar-view';
import { getCalendarEntryDisplayEndDate, getCalendarEntryDisplayStartDate, isCalendarEntryAllDay } from '~/utils/calendar';
import { convertComponentIconNameToDynamic } from '~/utils/icons';
import {
    isSystemManagedCalendarEntry,
    checkCalendarAccess,
    isValidCalendarEntryRecurring,
} from '~/components/calendar/helpers';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access/access';
import { CalendarEntryRecurringEvery, type CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';
import { toDate } from '~/utils/time';
import CustomContentRenderer from '~/components/partials/content/CustomContentRenderer.vue';
import EntryRSVPList from '~/components/calendar/entry/EntryRSVPList.vue';
import EntryActionButtons from '~/components/calendar/entry/EntryActionButtons.vue';
import { useCalendarEntryShortcutState } from '~/composables/useCalendarEntryShortcutState';
import { useCalendarStore } from '~/stores/calendar';

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        entry: CalendarEntry;
        showTime?: boolean;
        compact?: boolean;
        stacked?: boolean;
    }>(),
    {
        showTime: true,
        compact: false,
        stacked: undefined,
    },
);

const emit = defineEmits<{
    (e: 'select', entry: CalendarEntry): void;
    (e: 'edit', entry: CalendarEntry): void;
    (e: 'share', entry: CalendarEntry): void;
    (e: 'delete', entry: CalendarEntry): void;
    (e: 'update:popover-open', open: boolean): void;
}>();

const calendarStore = useCalendarStore();
const { calendars } = storeToRefs(calendarStore);

const chipRef = useTemplateRef('chipRef');
const popoverPoint = ref<{ x: number; y: number } | null>(null);
const popoverOpen = ref(false);

const shortcutState = useCalendarEntryShortcutState();
const { width: chipWidth, height: chipHeight } = useElementSize(chipRef);
const rsvpEntry = ref(props.entry.rsvp);

watch(
    () => props.entry.rsvp,
    (next) => (rsvpEntry.value = next),
);

const color = computed(() => getCalendarEntryColor(props.entry));
const icon = computed(() => getCalendarEntryIcon(props.entry));
const time = computed(() => getCalendarEntryTimeLabel(props.entry));
const isAllDay = computed(() => isCalendarEntryAllDay(props.entry));
const startDate = computed(() => getCalendarEntryDisplayStartDate(props.entry));
const endDate = computed(() => getCalendarEntryDisplayEndDate(props.entry));
const durationMinutes = computed(() => {
    if (!endDate.value) return undefined;

    return Math.max(0, Math.round((endDate.value.getTime() - startDate.value.getTime()) / 60000));
});
const spansMultipleDays = computed(() => !!endDate.value && !isSameDay(startDate.value, endDate.value));
const timeRange = computed(() => {
    if (isAllDay.value) {
        if (!endDate.value) {
            return format(startDate.value, 'P');
        }

        if (isSameDay(startDate.value, endDate.value)) {
            return format(startDate.value, 'P');
        }

        return `${format(startDate.value, 'P')} - ${format(endDate.value, 'P')}`;
    }

    if (!endDate.value) {
        return format(startDate.value, 'p');
    }

    if (spansMultipleDays.value) {
        return `${format(startDate.value, 'P p')}\n${format(endDate.value, 'P p')}`;
    }

    return `${format(startDate.value, 'p')} - ${format(endDate.value, 'p')}`;
});
const popoverTime = computed(() => (spansMultipleDays.value ? timeRange.value : time.value));
const calendar = computed(
    () => calendars.value.find((candidate) => candidate.id === props.entry.calendarId) ?? props.entry.calendar,
);
const isSystemManaged = computed(() => isSystemManagedCalendarEntry(calendar.value, props.entry));
const showStackedTime = computed(
    () =>
        props.stacked ??
        (props.showTime &&
            !props.compact &&
            !isAllDay.value &&
            (spansMultipleDays.value ||
                chipHeight.value >= 44 ||
                (durationMinutes.value !== undefined && durationMinutes.value <= 120))),
);
const showTrailingTime = computed(
    () => props.showTime && !props.compact && !isAllDay.value && !showStackedTime.value && chipWidth.value >= 150,
);
const isNarrowChip = computed(() => chipWidth.value > 0 && chipWidth.value < 160);
const recurringLabel = computed(() => {
    const recurring = props.entry.recurring;
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
        ? ` · ${t('components.calendar.EntryCreateOrUpdateModal.recurring.until')} ${format(toDate(recurring.until), 'P')}`
        : '';

    return `${t('components.calendar.EntryCreateOrUpdateModal.recurring.every')} ${recurring.count} ${everyUnit}${until}`;
});
const popoverContent = {
    side: 'bottom' as const,
    align: 'start' as const,
    sideOffset: 8,
    collisionPadding: 12,
    updatePositionStrategy: 'always' as const,
};
const popoverReference = computed<ReferenceElement | undefined>(() => {
    if (popoverPoint.value) {
        return {
            getBoundingClientRect: () => new DOMRect(popoverPoint.value?.x ?? 0, popoverPoint.value?.y ?? 0, 0, 0),
            contextElement: chipRef.value ?? undefined,
        };
    }

    return chipRef.value ?? undefined;
});

const canEdit = computed(
    () =>
        !isSystemManaged.value &&
        checkCalendarAccess(
            calendar.value?.access,
            props.entry.creator,
            AccessLevel.EDIT,
            calendar.value?.job,
            calendar.value?.creatorJob,
        ),
);
const canShare = computed(
    () =>
        !isSystemManaged.value &&
        checkCalendarAccess(
            calendar.value?.access,
            props.entry.creator,
            AccessLevel.SHARE,
            calendar.value?.job,
            calendar.value?.creatorJob,
        ),
);
const canDelete = computed(
    () =>
        !isSystemManaged.value &&
        checkCalendarAccess(
            calendar.value?.access,
            props.entry.creator,
            AccessLevel.MANAGE,
            calendar.value?.job,
            calendar.value?.creatorJob,
        ),
);
const attrs = useAttrs();
const { t } = useI18n();
const popoverAttrs = computed(() => {
    const { class: _class, style: _style, ...rest } = attrs as Record<string, unknown>;
    return rest;
});

function capturePopoverPoint(event: PointerEvent): void {
    popoverPoint.value = {
        x: event.clientX,
        y: event.clientY,
    };
}

watch(popoverOpen, (open) => {
    emit('update:popover-open', open);
    shortcutState.isPopoverOpen.value = open;
});

onUnmounted(() => {
    shortcutState.isPopoverOpen.value = false;
});

const attrsClass = computed(
    () => attrs.class as string | Record<string, boolean> | (string | Record<string, boolean>)[] | undefined,
);
const attrsStyle = computed(() => attrs.style as string | Record<string, string | number> | undefined);
</script>

<template>
    <div ref="chipRef" class="block h-full w-full min-w-0" :class="attrsClass" :style="attrsStyle">
        <UPopover
            v-bind="popoverAttrs"
            v-model:open="popoverOpen"
            :reference="popoverReference"
            :content="popoverContent"
            :ui="{ content: 'w-[30rem] max-w-[calc(100vw-1rem)] p-5 sm:w-[36rem] sm:p-6' }"
        >
            <UButton
                :color="color as ButtonProps['color']"
                :variant="entry.deletedAt ? 'subtle' : 'soft'"
                :size="compact ? 'xs' : 'sm'"
                :icon="icon"
                :aria-label="
                    showStackedTime
                        ? `${entry.title}, ${popoverTime}`
                        : showTrailingTime
                          ? `${entry.title}, ${time}`
                          : entry.title
                "
                class="relative flex h-full w-full max-w-full min-w-0 overflow-hidden rounded-xs pr-2 pl-3 text-start shadow-none backdrop-blur-sm transition-colors"
                :class="[
                    showStackedTime
                        ? 'min-h-9 flex-col items-start justify-start gap-0.5 py-1'
                        : 'flex-row items-center gap-1.5 py-0.5',
                ]"
                @pointerdown="capturePopoverPoint"
            >
                <span class="absolute inset-y-0.5 start-0.5 w-1 rounded-full bg-current/35" />

                <span
                    class="min-w-0"
                    :class="
                        showStackedTime
                            ? 'w-full text-[11px] leading-tight font-medium'
                            : isNarrowChip
                              ? 'line-clamp-2 flex-1 basis-0 text-[10px] leading-tight font-medium break-words whitespace-normal'
                              : 'flex-1 basis-0 truncate text-[11px] leading-tight font-medium'
                    "
                >
                    {{ entry.title }}
                </span>

                <span v-if="showStackedTime" class="w-full min-w-0 text-[11px] whitespace-pre-line text-muted tabular-nums">
                    {{ timeRange }}
                </span>
                <span v-else-if="showTrailingTime" class="shrink-0 text-[10px] text-muted tabular-nums">
                    {{ time }}
                </span>
            </UButton>

            <template #content>
                <div class="flex flex-col gap-4">
                    <div class="flex items-start justify-between gap-3">
                        <div class="min-w-0">
                            <div class="flex min-w-0 items-center gap-1.5">
                                <UIcon
                                    v-if="entry.icon"
                                    class="size-4 shrink-0 text-muted"
                                    :name="convertComponentIconNameToDynamic(entry.icon)"
                                />
                                <p class="truncate text-sm font-semibold text-highlighted">
                                    {{ entry.title }}
                                </p>
                            </div>
                            <div class="flex flex-wrap items-center gap-1 text-muted">
                                <UIcon name="i-mdi-access-time" class="size-4" />

                                <p class="text-xs leading-5 whitespace-pre-line text-muted">
                                    {{ popoverTime }}
                                </p>
                            </div>
                        </div>

                        <UBadge
                            :color="stringToButtonColor(color ?? 'primary')"
                            :label="calendar?.name ?? $t('common.calendar')"
                            :icon="calendar?.icon ? convertComponentIconNameToDynamic(calendar.icon) : undefined"
                            variant="subtle"
                            size="md"
                        />
                    </div>

                    <div class="flex flex-wrap gap-2 text-xs">
                        <UBadge v-if="recurringLabel" class="inline-flex gap-1" color="neutral" icon="i-mdi-repeat" size="sm">
                            {{ recurringLabel }}
                        </UBadge>

                        <UBadge v-if="entry.creator" class="inline-flex gap-1" color="neutral" icon="i-mdi-account" size="sm">
                            <span>{{ $t('common.created_by') }}</span>
                            <CitizenInfoPopover
                                :user="entry.creator"
                                :show-avatar-in-name="true"
                                text-class="text-xs"
                                size="xs"
                            />
                        </UBadge>

                        <UBadge
                            v-if="entry.updatedAt"
                            class="inline-flex gap-1"
                            color="neutral"
                            icon="i-mdi-calendar-edit"
                            size="sm"
                        >
                            {{ $t('common.updated_at') }}
                            <GenericTime :value="entry.updatedAt" type="short" />
                        </UBadge>
                    </div>

                    <EntryRSVPList
                        v-if="entry.rsvpOpen"
                        v-model="rsvpEntry"
                        :entry-id="entry.id"
                        :occurrence-key="entry.occurrence?.key"
                        :rsvp-open="entry.rsvpOpen"
                        :disabled="entry.closed"
                        :can-share="canShare"
                    />

                    <div v-if="entry.content" class="line-clamp-6 overflow-hidden">
                        <CustomContentRenderer :value="entry.content" :placeholder="$t('common.na')" />
                    </div>

                    <EntryActionButtons
                        :entry="entry"
                        :can-edit="canEdit"
                        :can-share="canShare"
                        :can-delete="canDelete"
                        :show-open="true"
                        @open="emit('select', $event)"
                        @edit="emit('edit', $event)"
                        @share="emit('share', $event)"
                        @delete="emit('delete', $event)"
                    />
                </div>
            </template>
        </UPopover>
    </div>
</template>
