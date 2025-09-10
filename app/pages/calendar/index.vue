<script lang="ts" setup>
import type { ButtonProps } from '@nuxt/ui';
import { isFuture, isPast, isSameDay, isToday } from 'date-fns';
import type { DateRangeSource } from 'v-calendar/dist/types/src/utils/date/range.js';
import CalendarCreateOrUpdateModal from '~/components/calendar/CalendarCreateOrUpdateModal.vue';
import CalendarViewSlideover from '~/components/calendar/CalendarViewSlideover.vue';
import FindCalendarDrawer from '~/components/calendar/FindCalendarDrawer.vue';
import EntryCreateOrUpdateModal from '~/components/calendar/entry/EntryCreateOrUpdateModal.vue';
import EntryViewSlideover from '~/components/calendar/entry/EntryViewSlideover.vue';
import MonthCalendarClient from '~/components/partials/MonthCalendar.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useCalendarStore } from '~/stores/calendar';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';
import type { ListCalendarsResponse } from '~~/gen/ts/services/calendar/calendar';

useHead({
    title: 'common.calendar',
});

definePageMeta({
    title: 'common.calendar',
    requiresAuth: true,
});

const { t, d } = useI18n();

const { can } = useAuth();

const overlay = useOverlay();

const calendarStore = useCalendarStore();
const { activeCalendarIds, currentDate, view, calendars, entries, hasEditAccessToCalendar } = storeToRefs(calendarStore);

const calRef = useTemplateRef('calRef');

const page = useRouteQuery('page', '1', { transform: Number });

const {
    data: calendarsData,
    status: calendarsStatus,
    error: calendarsError,
    refresh: calendarsRefresh,
} = useLazyAsyncData(`calendars:${page.value}`, () => listCalendars());

async function listCalendars(): Promise<ListCalendarsResponse> {
    const response = await calendarStore.listCalendars({
        pagination: {
            offset: calculateOffset(page.value, calendarsData.value?.pagination),
        },
        onlyPublic: false,
    });

    if (activeCalendarIds.value.length === 0) {
        activeCalendarIds.value = response.calendars.map((c) => c.id);
    }

    refresh();

    return response;
}

const { refresh, status, error } = useLazyAsyncData(
    `calendar-entries-${currentDate.value.year}-${currentDate.value.month}-${activeCalendarIds.value.join(':')}`,
    () =>
        calendarStore.listCalendarEntries({
            year: currentDate.value.year,
            month: currentDate.value.month,
            calendarIds: activeCalendarIds.value,
        }),
    { immediate: false },
);

watchDebounced(currentDate.value, async () => refresh(), { debounce: 100, maxWait: 1000 });
watchDebounced(activeCalendarIds, async () => refresh());

function formatStartEndTime(entry: CalendarEntry): string {
    const start = toDate(entry.startTime);
    const end = entry.endTime ? toDate(entry.endTime) : undefined;

    if (!end) {
        return d(start, 'time');
    }

    return (
        d(start, 'time') +
        ' - ' +
        d(
            end,
            isSameDay(start, end)
                ? 'time'
                : {
                      month: '2-digit',
                      day: '2-digit',
                      hour: 'numeric',
                      minute: 'numeric',
                  },
        )
    );
}

type CalEntry = {
    key: string;
    customData: CalendarEntry & {
        color: string;
        isPast: boolean;
        multiDay: boolean;
        ongoing: boolean;
        time: string;
        timeEnd?: string;
    };
    dates: DateRangeSource | DateRangeSource[];
};

const transformedCalendarEntries = computedAsync(async () =>
    entries.value
        .filter((e) => activeCalendarIds.value.includes(e.calendarId))
        .map((entry) => {
            const startTime = toDate(entry.startTime);
            const endTime = entry.endTime ? toDate(entry.endTime) : undefined;
            const past = endTime ? isPast(endTime) : isPast(startTime);

            return {
                key: `${startTime.toISOString()}-${entry.id}-${entry.calendarId}`,
                customData: {
                    ...entry,
                    color: entry.calendar?.color ?? 'primary',
                    isPast: past,
                    multiDay: !!endTime && !isSameDay(startTime, endTime),
                    ongoing: !!endTime && isPast(startTime) && isFuture(endTime),
                    time: formatStartEndTime(entry),
                    timeEnd:
                        endTime && !isSameDay(startTime, endTime)
                            ? d(startTime, {
                                  month: '2-digit',
                                  day: '2-digit',
                                  hour: 'numeric',
                                  minute: 'numeric',
                              }) +
                              ' - ' +
                              d(endTime, 'time')
                            : undefined,
                },
                dates: {
                    start: startTime,
                    end: endTime,
                    repeat: entry.recurring
                        ? {
                              every: [entry.recurring.count, entry.recurring.every],
                              until: entry.recurring?.until ? toDate(entry.recurring?.until) : undefined,
                          }
                        : undefined,
                } as DateRangeSource,
            };
        })
        .sort((a, b) => a.key.localeCompare(b.key) + (b.customData.id - a.customData.id)),
);

type GroupedCalendarEntries = {
    key: string;
    date: Date;
    isToday: boolean;
    entries: { past: CalEntry[]; upcoming: CalEntry[] };
}[];

const groupedCalendarEntries = computedAsync(async () => {
    const groups: GroupedCalendarEntries = [];
    transformedCalendarEntries.value?.forEach((entry) => {
        const date = toDate(entry.customData.startTime);
        let idx = groups.findIndex((g) => g.key === toDate(entry.customData.startTime).toDateString());
        if (idx === -1) {
            idx = groups.push({
                key: date.toDateString(),
                date: date,
                isToday: isToday(date),
                entries: {
                    past: [],
                    upcoming: [],
                },
            });
            idx = idx - 1;
        }

        if (entry.customData.isPast) {
            groups[idx]!.entries.past.push(entry);
        } else {
            groups[idx]!.entries.upcoming.push(entry);
        }
    });

    if (!groups.find((g) => g.isToday)) {
        const now = new Date();
        groups.push({
            key: now.toDateString(),
            date: now,
            isToday: true,
            entries: {
                past: [],
                upcoming: [],
            },
        });
    }

    return groups.sort((a, b) => b.date.getTime() - a.date.getTime());
});

function calendarIdChange(calendarId: number, state: boolean): void {
    if (state) {
        if (!activeCalendarIds.value.includes(calendarId)) activeCalendarIds.value.push(calendarId);
    } else {
        activeCalendarIds.value = activeCalendarIds.value.filter((cId) => cId !== calendarId);
    }
}

const entryIdQuery = useRouteQuery('entry_id', undefined, { transform: Number });

const calendarViewSlideover = overlay.create(CalendarViewSlideover);
const calendarCreateOrUpdateModal = overlay.create(CalendarCreateOrUpdateModal);
const entryViewSlideover = overlay.create(EntryViewSlideover);
const entryCreateOrUpdateModal = overlay.create(EntryCreateOrUpdateModal);
const findCalendarsDrawer = overlay.create(FindCalendarDrawer);

watch(entryIdQuery, () => {
    if (!entryIdQuery.value) {
        return;
    }

    entryViewSlideover.open({
        entryId: entryIdQuery.value,
    });
});

if (entryIdQuery.value) {
    entryViewSlideover.open({
        entryId: entryIdQuery.value,
    });
}

async function resetToToday(): Promise<void> {
    calRef.value?.calRef?.focusDate(new Date());
}

const loadingState = ref(false);
watch(status, () => {
    if (isRequestPending(status.value)) {
        loadingState.value = true;
    }
});
watchDebounced(
    status,
    () => {
        if (status.value === 'success' || status.value === 'error') {
            loadingState.value = false;
        }
    },
    {
        debounce: 750,
        maxWait: 1250,
    },
);

const viewOptions = [
    { label: t('common.week_view'), icon: 'i-mdi-view-week', value: 'week' },
    { label: t('common.monthly_view'), icon: 'i-mdi-view-headline', value: 'month' },
    { label: t('common.summary'), icon: 'i-mdi-view-agenda-outline', value: 'summary' },
];
</script>

<!-- eslint-disable vue/no-multiple-template-root -->
<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.calendar')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <UButtonGroup
                        v-if="can('calendar.CalendarService/CreateCalendar').value || hasEditAccessToCalendar"
                        class="inline-flex w-full xl:hidden"
                    >
                        <UButton
                            v-if="can('calendar.CalendarService/CreateCalendar').value"
                            class="flex-1"
                            block
                            color="neutral"
                            trailing-icon="i-mdi-plus"
                            :label="$t('common.calendar')"
                            @click="calendarCreateOrUpdateModal.open({})"
                        />

                        <UButton
                            v-if="hasEditAccessToCalendar"
                            class="flex-1"
                            block
                            color="neutral"
                            trailing-icon="i-mdi-plus"
                            :label="$t('common.entry', 1)"
                            @click="entryCreateOrUpdateModal.open({})"
                        />
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <template #default>
                    <div class="flex flex-1 items-center justify-between">
                        <UPopover :content="{ side: 'bottom', align: 'start' }">
                            <UButton
                                color="neutral"
                                icon="i-mdi-calendar"
                                trailing-icon="i-mdi-chevron-down"
                                :loading="isRequestPending(calendarsStatus)"
                                :label="$t('common.calendar')"
                            />

                            <template #content>
                                <div class="p-4">
                                    <DataPendingBlock
                                        v-if="isRequestPending(calendarsStatus)"
                                        :message="$t('common.loading', [$t('common.calendar')])"
                                        class="max-w-60"
                                    />
                                    <DataErrorBlock
                                        v-else-if="calendarsError"
                                        :title="$t('common.unable_to_load', [$t('common.calendar')])"
                                        :error="calendarsError"
                                        :retry="calendarsRefresh"
                                    />

                                    <div v-else class="grid grid-cols-1 gap-2 sm:grid-cols-2">
                                        <div
                                            v-for="calendar in calendars"
                                            :key="calendar.id"
                                            class="inline-flex items-center gap-2"
                                        >
                                            <USwitch
                                                class="truncate"
                                                :model-value="activeCalendarIds.includes(calendar.id)"
                                                @update:model-value="
                                                    ($event) => calendarIdChange(calendar.id, $event as boolean)
                                                "
                                            />

                                            <UButton
                                                :color="calendar.color as ButtonProps['color']"
                                                size="sm"
                                                truncate
                                                :label="calendar.name"
                                                @click="calendarViewSlideover.open({ calendarId: calendar.id })"
                                            />
                                        </div>
                                    </div>
                                </div>
                            </template>
                        </UPopover>

                        <UButton
                            icon="i-mdi-calendar-today"
                            :disabled="isRequestPending(calendarsStatus)"
                            @click="resetToToday"
                        >
                            {{ $t('common.today') }}
                        </UButton>
                    </div>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataErrorBlock v-if="error" :error="error" :retry="refresh" />

            <div v-else class="relative flex flex-1 overflow-x-auto">
                <MonthCalendarClient
                    v-if="view !== 'summary'"
                    ref="calRef"
                    class="flex flex-1"
                    :view="view === 'week' ? 'weekly' : 'monthly'"
                    :attributes="transformedCalendarEntries"
                    @selected="
                        entryViewSlideover.open({
                            entryId: $event.id,
                        })
                    "
                    @did-move="
                        currentDate.year = $event[0].year;
                        currentDate.month = $event[0].month;
                    "
                />

                <UContainer v-else class="flex flex-1 flex-col py-2">
                    <DataErrorBlock
                        v-if="error"
                        :title="$t('common.unable_to_load', [$t('common.entry', 2)])"
                        :error="error"
                        :retry="refresh"
                    />

                    <template v-else>
                        <template v-for="calendarEntries in groupedCalendarEntries" :key="calendarEntries.key">
                            <USeparator>
                                <div class="inline-flex items-center gap-1">
                                    <span class="text-lg font-semibold">
                                        {{ $d(calendarEntries.date, 'date') }}
                                    </span>
                                    <UBadge
                                        v-if="calendarEntries.isToday"
                                        id="today"
                                        size="xs"
                                        color="warning"
                                        :label="$t('common.today')"
                                    />
                                </div>
                            </USeparator>

                            <ul role="list">
                                <li v-for="attr in calendarEntries.entries.past" :key="attr.key">
                                    <ULink
                                        class="inline-flex w-full items-center justify-between gap-1"
                                        @click="
                                            entryViewSlideover.open({
                                                entryId: attr.customData.id,
                                            })
                                        "
                                    >
                                        <span class="inline-flex items-center gap-1">
                                            <UBadge :color="attr.customData.color as ButtonProps['color']" size="lg" />

                                            <template v-if="attr.customData.time">
                                                {{ attr.customData.time }}
                                            </template>
                                            <span>-</span>

                                            {{ attr.customData.title }}
                                        </span>

                                        <UButton variant="link" icon="i-mdi-eye" />
                                    </ULink>
                                </li>

                                <li>
                                    <USeparator
                                        v-if="
                                            calendarEntries.isToday &&
                                            (calendarEntries.entries.past.length > 0 ||
                                                calendarEntries.entries.upcoming.length > 0)
                                        "
                                        class="my-1"
                                        size="sm"
                                        :ui="{ border: 'border-red-300 dark:border-red-600' }"
                                    />
                                </li>

                                <li v-for="attr in calendarEntries.entries.upcoming" :key="attr.key">
                                    <ULink
                                        class="inline-flex w-full items-center justify-between gap-1"
                                        @click="
                                            entryViewSlideover.open({
                                                entryId: attr.customData.id,
                                            })
                                        "
                                    >
                                        <span class="inline-flex items-center gap-1">
                                            <UBadge :color="attr.customData.color as ButtonProps['color']" size="lg" />

                                            <template v-if="attr.customData.time">
                                                {{ attr.customData.time }}
                                            </template>
                                            <span>-</span>

                                            <UIcon
                                                v-if="attr.customData.ongoing"
                                                class="size-3 text-amber-800"
                                                name="i-mdi-timer-sand"
                                            />

                                            {{ attr.customData.title }}
                                        </span>

                                        <UButton variant="link" icon="i-mdi-eye" />
                                    </ULink>
                                </li>
                            </ul>
                        </template>
                    </template>
                </UContainer>
            </div>
        </template>

        <template #footer>
            <div
                class="flex justify-between border-t border-b-0 border-neutral-200 px-3 py-3.5 xl:hidden dark:border-neutral-700"
            >
                <UFormField
                    class="flex flex-row items-center gap-2"
                    :label="$t('common.view')"
                    :ui="{ container: '', label: 'hidden md:inline-flex' }"
                >
                    <ClientOnly>
                        <USelectMenu v-model="view" :items="viewOptions" value-key="value">
                            <template #default>
                                <UIcon
                                    class="size-5"
                                    :name="viewOptions.find((o) => o.value === view)?.icon ?? 'i-mdi-view-'"
                                />

                                {{ viewOptions.find((o) => o.value === view)?.label ?? $t('common.na') }}
                            </template>

                            <template #item="{ item }">
                                <UIcon class="size-5" :name="item.icon" />

                                <span class="truncate">{{ item.label }}</span>
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <div>
                    <UTooltip :text="$t('common.refresh')">
                        <UButton
                            icon="i-mdi-refresh"
                            variant="outline"
                            :disabled="isRequestPending(status) || loadingState"
                            :loading="isRequestPending(status) || loadingState"
                            :label="$t('common.refresh')"
                            @click="refresh()"
                        />
                    </UTooltip>
                </div>

                <div>
                    <UButton
                        class="font-semibold"
                        icon="i-mdi-calendar-search"
                        :label="$t('components.calendar.FindCalendarDrawer.title')"
                        @click="findCalendarsDrawer.open({})"
                    />
                </div>
            </div>
        </template>
    </UDashboardPanel>

    <UDashboardPanel :ui="{ root: 'hidden xl:flex max-w-90', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar>
                <template #right>
                    <UButtonGroup
                        v-if="can('calendar.CalendarService/CreateCalendar').value || hasEditAccessToCalendar"
                        class="inline-flex w-full"
                    >
                        <UButton
                            v-if="can('calendar.CalendarService/CreateCalendar').value"
                            class="flex-1"
                            block
                            color="neutral"
                            trailing-icon="i-mdi-plus"
                            :label="$t('common.calendar')"
                            @click="calendarCreateOrUpdateModal.open({})"
                        />

                        <UButton
                            v-if="hasEditAccessToCalendar"
                            class="flex-1"
                            block
                            color="neutral"
                            trailing-icon="i-mdi-plus"
                            :label="$t('common.entry', 1)"
                            @click="entryCreateOrUpdateModal.open({})"
                        />
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div class="mx-2 mb-2 flex h-full flex-col gap-2">
                <div>
                    <p class="font-semibold">{{ $t('common.calendar') }}</p>

                    <DataPendingBlock
                        v-if="isRequestPending(calendarsStatus)"
                        :message="$t('common.loading', [$t('common.calendar')])"
                    />
                    <DataErrorBlock
                        v-else-if="calendarsError"
                        :title="$t('common.unable_to_load', [$t('common.calendar', 1)])"
                        :error="calendarsError"
                        :retry="calendarsRefresh"
                    />

                    <div v-else class="grid grid-cols-1 gap-2">
                        <div v-for="calendar in calendars" :key="calendar.id" class="inline-flex items-center gap-2 truncate">
                            <USwitch
                                :model-value="activeCalendarIds.includes(calendar.id)"
                                @update:model-value="($event) => calendarIdChange(calendar.id, $event as boolean)"
                            />

                            <UButton
                                :color="calendar.color as ButtonProps['color']"
                                variant="solid"
                                size="sm"
                                truncate
                                :label="calendar.name"
                                @click="calendarViewSlideover.open({ calendarId: calendar.id })"
                            />
                        </div>
                    </div>
                </div>

                <div class="flex-1" />
            </div>
        </template>

        <template #footer>
            <USeparator class="sticky bottom-0" />

            <div class="flex flex-col gap-2 p-2">
                <UFormField
                    class="flex flex-1 flex-row items-center gap-2"
                    :label="$t('common.view')"
                    :ui="{ container: 'flex-1' }"
                >
                    <ClientOnly>
                        <USelectMenu v-model="view" class="w-full min-w-44" :items="viewOptions" value-key="value">
                            <template #default>
                                <UIcon
                                    class="size-5"
                                    :name="viewOptions.find((o) => o.value === view)?.icon ?? 'i-mdi-view-'"
                                />

                                {{ viewOptions.find((o) => o.value === view)?.label ?? $t('common.na') }}
                            </template>

                            <template #item="{ item }">
                                <UIcon class="size-5" :name="item.icon" />
                                <span class="truncate">{{ item.label }}</span>
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UTooltip class="inline-flex w-full" :text="$t('common.refresh')">
                    <UButton
                        class="w-full"
                        icon="i-mdi-refresh"
                        variant="outline"
                        :disabled="isRequestPending(status) || loadingState"
                        :loading="isRequestPending(status) || loadingState"
                        @click="refresh()"
                    >
                        {{ $t('common.refresh') }}
                    </UButton>
                </UTooltip>

                <UButton
                    class="font-semibold"
                    icon="i-mdi-calendar-search"
                    :label="$t('components.calendar.FindCalendarDrawer.title')"
                    @click="findCalendarsDrawer.open({})"
                />
            </div>
        </template>
    </UDashboardPanel>
</template>
