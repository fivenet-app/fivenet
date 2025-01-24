<script lang="ts" setup>
import type { ButtonColor } from '#ui/types';
import { isFuture, isPast, isSameDay, isToday } from 'date-fns';
import type { DateRangeSource } from 'v-calendar/dist/types/src/utils/date/range.js';
import CalendarCreateOrUpdateModal from '~/components/calendar/CalendarCreateOrUpdateModal.vue';
import CalendarViewSlideover from '~/components/calendar/CalendarViewSlideover.vue';
import FindCalendarsModal from '~/components/calendar/FindCalendarsModal.vue';
import EntryCreateOrUpdateModal from '~/components/calendar/entry/EntryCreateOrUpdateModal.vue';
import EntryViewSlideover from '~/components/calendar/entry/EntryViewSlideover.vue';
import MonthCalendarClient from '~/components/partials/MonthCalendar.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useCalendarStore } from '~/store/calendar';
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

const modal = useModal();
const slideover = useSlideover();

const calendarStore = useCalendarStore();
const { activeCalendarIds, currentDate, view, calendars, entries } = storeToRefs(calendarStore);

const calRef = ref<InstanceType<typeof MonthCalendarClient> | null>(null);

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() =>
    calendarsData.value?.pagination?.pageSize ? calendarsData.value?.pagination?.pageSize * (page.value - 1) : 0,
);

const {
    data: calendarsData,
    pending: calendarsLoading,
    error: calendarsError,
    refresh: calendarsRefresh,
} = useLazyAsyncData(`calendars:${page.value}`, () => listCalendars());

async function listCalendars(): Promise<ListCalendarsResponse> {
    const response = await calendarStore.listCalendars({
        pagination: {
            offset: offset.value,
        },
        onlyPublic: false,
    });

    if (activeCalendarIds.value.length === 0) {
        activeCalendarIds.value = response.calendars.map((c) => c.id);
    }

    refresh();

    return response;
}

const {
    refresh,
    pending: loading,
    error,
} = useLazyAsyncData(
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
        if (!activeCalendarIds.value.includes(calendarId)) {
            activeCalendarIds.value.push(calendarId);
        }
    } else {
        activeCalendarIds.value = activeCalendarIds.value.filter((cId) => cId !== calendarId);
    }
}

const entryIdQuery = useRouteQuery('entry_id', undefined, { transform: Number });

watch(entryIdQuery, () => {
    if (!entryIdQuery.value) {
        return;
    }

    slideover.open(EntryViewSlideover, {
        entryId: entryIdQuery.value,
    });
});

if (entryIdQuery.value) {
    slideover.open(EntryViewSlideover, {
        entryId: entryIdQuery.value,
    });
}

async function resetToToday(): Promise<void> {
    calRef.value?.calRef?.focusDate(new Date());
}

const loadingState = ref(false);
watch(loading, () => {
    if (loading.value) {
        loadingState.value = true;
    }
});
watchDebounced(
    loading,
    () => {
        if (!loading.value) {
            loadingState.value = false;
        }
    },
    {
        debounce: 750,
        maxWait: 1250,
    },
);

const viewOptions = [
    { label: t('common.time_ago.week'), icon: 'i-mdi-view-week', value: 'week' },
    { label: t('common.time_ago.month'), icon: 'i-mdi-view-headline', value: 'month' },
    { label: t('common.summary'), icon: 'i-mdi-view-agenda-outline', value: 'summary' },
];

const isOpen = ref(false);
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel class="shrink-0 border-b border-gray-200 lg:border-b-0 lg:border-r dark:border-gray-800" grow>
            <UDashboardNavbar :title="$t('common.calendar')">
                <template #right>
                    <UButtonGroup
                        v-if="
                            can('CalendarService.CreateOrUpdateCalendarEntry').value ||
                            can('CalendarService.CreateOrUpdateCalendar').value
                        "
                        class="inline-flex w-full xl:hidden"
                    >
                        <UButton
                            v-if="can('CalendarService.CreateOrUpdateCalendar').value"
                            block
                            color="gray"
                            trailing-icon="i-mdi-plus"
                            class="flex-1"
                            @click="modal.open(CalendarCreateOrUpdateModal, {})"
                        >
                            {{ $t('common.calendar') }}
                        </UButton>

                        <UButton
                            v-if="can('CalendarService.CreateOrUpdateCalendarEntry').value && calendars.length > 0"
                            block
                            color="gray"
                            trailing-icon="i-mdi-plus"
                            class="flex-1"
                            @click="modal.open(EntryCreateOrUpdateModal, {})"
                        >
                            {{ $t('common.entry', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <template #default>
                    <div class="flex flex-1 items-center justify-between">
                        <UPopover :popper="{ placement: 'bottom-end', offsetDistance: 10 }">
                            <UButton
                                color="white"
                                icon="i-mdi-calendar"
                                trailing-icon="i-mdi-chevron-down"
                                :loading="calendarsLoading"
                            >
                                {{ $t('common.calendar') }}
                            </UButton>

                            <template #panel>
                                <div class="p-4">
                                    <DataPendingBlock
                                        v-if="calendarsLoading"
                                        :message="$t('common.loading', [$t('common.calendar')])"
                                    />
                                    <DataErrorBlock
                                        v-else-if="calendarsError"
                                        :title="$t('common.unable_to_load', [$t('common.calendar')])"
                                        :error="calendarsError"
                                        :retry="calendarsRefresh"
                                    />

                                    <div v-else class="flex flex-col gap-4">
                                        <div class="grid grid-cols-1 gap-2 sm:grid-cols-2 lg:grid-cols-3">
                                            <div
                                                v-for="calendar in calendars"
                                                :key="calendar.id"
                                                class="inline-flex items-center gap-2"
                                            >
                                                <UCheckbox
                                                    :model-value="activeCalendarIds.includes(calendar.id)"
                                                    class="truncate"
                                                    @change="calendarIdChange(calendar.id, $event)"
                                                />

                                                <UBadge
                                                    :color="calendar.color as ButtonColor"
                                                    :ui="{ rounded: 'rounded-full' }"
                                                    size="lg"
                                                />

                                                <UButton
                                                    :color="calendar.color as ButtonColor"
                                                    size="sm"
                                                    variant="link"
                                                    :padded="false"
                                                    truncate
                                                    :label="calendar.name"
                                                    @click="slideover.open(CalendarViewSlideover, { calendarId: calendar.id })"
                                                />
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </template>
                        </UPopover>

                        <UButton icon="i-mdi-calendar-today" :disabled="calendarsLoading" @click="resetToToday">
                            {{ $t('common.today') }}
                        </UButton>
                    </div>
                </template>
            </UDashboardToolbar>

            <DataErrorBlock v-if="error" :error="error" :retry="refresh" />

            <div v-else class="relative flex flex-1 overflow-x-auto">
                <MonthCalendarClient
                    v-if="view !== 'summary'"
                    ref="calRef"
                    class="flex flex-1"
                    :view="view === 'week' ? 'weekly' : 'monthly'"
                    :attributes="transformedCalendarEntries"
                    @selected="
                        slideover.open(EntryViewSlideover, {
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
                            <UDivider>
                                <div class="inline-flex items-center gap-1">
                                    <span class="text-lg font-semibold">
                                        {{ $d(calendarEntries.date, 'date') }}
                                    </span>
                                    <UBadge
                                        v-if="calendarEntries.isToday"
                                        id="today"
                                        size="xs"
                                        color="amber"
                                        :label="$t('common.today')"
                                    />
                                </div>
                            </UDivider>

                            <ul role="list">
                                <li v-for="attr in calendarEntries.entries.past" :key="attr.key">
                                    <ULink
                                        class="inline-flex w-full items-center justify-between gap-1"
                                        @click="
                                            slideover.open(EntryViewSlideover, {
                                                entryId: attr.customData.id,
                                            })
                                        "
                                    >
                                        <span class="inline-flex items-center gap-1">
                                            <UBadge
                                                :color="attr.customData.color as ButtonColor"
                                                :ui="{ rounded: 'rounded-full' }"
                                                size="lg"
                                            />

                                            <template v-if="attr.customData.time">
                                                {{ attr.customData.time }}
                                            </template>
                                            <span>-</span>

                                            {{ attr.customData.title }}
                                        </span>

                                        <UButton :padded="false" variant="link" icon="i-mdi-eye" />
                                    </ULink>
                                </li>

                                <li>
                                    <UDivider
                                        v-if="
                                            calendarEntries.isToday &&
                                            (calendarEntries.entries.past.length > 0 ||
                                                calendarEntries.entries.upcoming.length > 0)
                                        "
                                        size="sm"
                                        :ui="{ border: { base: 'border-red-300 dark:border-red-600' } }"
                                        class="my-1"
                                    />
                                </li>

                                <li v-for="attr in calendarEntries.entries.upcoming" :key="attr.key">
                                    <ULink
                                        class="inline-flex w-full items-center justify-between gap-1"
                                        @click="
                                            slideover.open(EntryViewSlideover, {
                                                entryId: attr.customData.id,
                                            })
                                        "
                                    >
                                        <span class="inline-flex items-center gap-1">
                                            <UBadge
                                                :color="attr.customData.color as ButtonColor"
                                                :ui="{ rounded: 'rounded-full' }"
                                                size="lg"
                                            />

                                            <template v-if="attr.customData.time">
                                                {{ attr.customData.time }}
                                            </template>
                                            <span>-</span>

                                            <UIcon
                                                v-if="attr.customData.ongoing"
                                                name="i-mdi-timer-sand"
                                                class="size-3 text-amber-800"
                                            />

                                            {{ attr.customData.title }}
                                        </span>

                                        <UButton :padded="false" variant="link" icon="i-mdi-eye" />
                                    </ULink>
                                </li>
                            </ul>
                        </template>
                    </template>
                </UContainer>
            </div>

            <div class="flex justify-between border-b-0 border-t border-gray-200 px-3 py-3.5 xl:hidden dark:border-gray-700">
                <UFormGroup
                    :label="$t('common.view')"
                    :ui="{ container: '', label: { base: 'hidden md:inline-flex' } }"
                    class="flex flex-row items-center gap-2"
                >
                    <ClientOnly>
                        <USelectMenu v-model="view" :options="viewOptions" value-attribute="value">
                            <template #label>
                                <UIcon
                                    :name="viewOptions.find((o) => o.value === view)?.icon ?? 'i-mdi-view-'"
                                    class="size-5"
                                />

                                {{ viewOptions.find((o) => o.value === view)?.label ?? $t('common.na') }}
                            </template>

                            <template #option="{ option }">
                                <UIcon :name="option.icon" class="size-5" />
                                <span class="truncate">{{ option.label }}</span>
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormGroup>

                <UButton
                    icon="i-mdi-refresh"
                    variant="outline"
                    :title="$t('common.refresh')"
                    :disabled="loading || loadingState"
                    :loading="loading || loadingState"
                    @click="refresh()"
                >
                    {{ $t('common.refresh') }}
                </UButton>

                <UButton icon="i-mdi-search" class="font-semibold" @click="modal.open(FindCalendarsModal, {})">
                    {{ $t('components.calendar.FindCalendarsModal.title') }}
                </UButton>
            </div>
        </UDashboardPanel>

        <UDashboardPanel v-model="isOpen" collapsible side="right" class="!hidden max-w-64 flex-1 xl:!flex">
            <UDashboardNavbar>
                <template #right>
                    <UButtonGroup
                        v-if="
                            can('CalendarService.CreateOrUpdateCalendarEntry').value ||
                            can('CalendarService.CreateOrUpdateCalendar').value
                        "
                        class="inline-flex w-full"
                    >
                        <UButton
                            v-if="can('CalendarService.CreateOrUpdateCalendar').value"
                            block
                            color="gray"
                            trailing-icon="i-mdi-plus"
                            class="flex-1"
                            @click="modal.open(CalendarCreateOrUpdateModal, {})"
                        >
                            {{ $t('common.calendar') }}
                        </UButton>

                        <UButton
                            v-if="can('CalendarService.CreateOrUpdateCalendarEntry').value"
                            block
                            color="gray"
                            trailing-icon="i-mdi-plus"
                            class="flex-1"
                            @click="modal.open(EntryCreateOrUpdateModal, {})"
                        >
                            {{ $t('common.entry', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>

            <div class="mx-2 mb-2 flex h-full flex-col gap-2">
                <div>
                    <p class="font-semibold">{{ $t('common.calendar') }}</p>

                    <DataPendingBlock v-if="calendarsLoading" :message="$t('common.loading', [$t('common.calendar')])" />
                    <DataErrorBlock
                        v-else-if="calendarsError"
                        :title="$t('common.unable_to_load', [$t('common.calendar', 1)])"
                        :error="calendarsError"
                        :retry="calendarsRefresh"
                    />

                    <div v-else class="grid grid-cols-1 gap-2">
                        <div v-for="calendar in calendars" :key="calendar.id" class="inline-flex items-center gap-2 truncate">
                            <UCheckbox
                                :model-value="activeCalendarIds.includes(calendar.id)"
                                @change="calendarIdChange(calendar.id, $event)"
                            />

                            <UBadge :color="calendar.color as ButtonColor" :ui="{ rounded: 'rounded-full' }" />

                            <UButton
                                :color="calendar.color as ButtonColor"
                                :padded="false"
                                variant="link"
                                size="sm"
                                truncate
                                :label="calendar.name"
                                @click="slideover.open(CalendarViewSlideover, { calendarId: calendar.id })"
                            />
                        </div>
                    </div>
                </div>

                <div class="flex-1" />

                <UDivider class="sticky bottom-0" />

                <UFormGroup :label="$t('common.view')" class="flex flex-row items-center gap-2">
                    <ClientOnly>
                        <USelectMenu v-model="view" :options="viewOptions" value-attribute="value" class="min-w-44">
                            <template #label>
                                <UIcon
                                    :name="viewOptions.find((o) => o.value === view)?.icon ?? 'i-mdi-view-'"
                                    class="size-5"
                                />

                                {{ viewOptions.find((o) => o.value === view)?.label ?? $t('common.na') }}
                            </template>

                            <template #option="{ option }">
                                <UIcon :name="option.icon" class="size-5" />
                                <span class="truncate">{{ option.label }}</span>
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormGroup>

                <UTooltip :text="$t('common.refresh')" class="inline-flex w-full">
                    <UButton
                        icon="i-mdi-refresh"
                        variant="outline"
                        class="w-full"
                        :disabled="loading || loadingState"
                        :loading="loading || loadingState"
                        @click="refresh()"
                    >
                        {{ $t('common.refresh') }}
                    </UButton>
                </UTooltip>

                <UButton icon="i-mdi-search" class="font-semibold" @click="modal.open(FindCalendarsModal, {})">
                    {{ $t('components.calendar.FindCalendarsModal.title') }}
                </UButton>
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>
