<script setup lang="ts">
import type { DateRange } from 'reka-ui';
import { addDays, addMonths } from 'date-fns';
import { z } from 'zod';
import type { CalendarDate, DateValue } from '@internationalized/date';
import type { BadgeProps, DropdownMenuItem } from '@nuxt/ui';
import { getCalendarEntriesClient } from '~~/gen/ts/clients';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';
import type { ListCalendarsResponse } from '~~/gen/ts/services/calendar/calendar';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import CreateOrUpdateModal from '~/components/calendar/calendar/CreateOrUpdateModal.vue';
import ViewSlideover from '~/components/calendar/calendar/ViewSlideover.vue';
import FindCalendarDrawer from '~/components/calendar/calendar/FindCalendarDrawer.vue';
import EntryCreateOrUpdateModal from '~/components/calendar/entry/EntryCreateOrUpdateModal.vue';
import EntryViewModal from '~/components/calendar/entry/EntryViewModal.vue';
import CalendarMonthView from '~/components/calendar/calendar/view/CalendarMonthView.vue';
import CalendarWeekView from '~/components/calendar/calendar/view/CalendarWeekView.vue';
import CalendarDayView from '~/components/calendar/calendar/view/CalendarDayView.vue';
import CalendarAgendaView from '~/components/calendar/calendar/view/CalendarAgendaView.vue';
import { useCalendarStore } from '~/stores/calendar';
import { checkCalendarAccess } from '~/components/calendar/helpers';
import { dateToCalendarDate } from '~/utils/time';
import { getCalendarEntryDisplayStartDate } from '~/utils/calendar';
import {
    fetchMonthsForRange,
    isCalendarView,
    normalizeCalendarView,
    routeForCalendarView,
    rangeTitle,
    viewDateRange,
} from '~/utils/calendar-view';
import { unsafeRoute } from '~/utils/route';
import Pagination from '~/components/partials/Pagination.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import { copyToClipboardWrapper } from '~/utils/clipboard';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access/access';

useHead({
    title: 'common.calendar',
});

definePageMeta({
    requiresAuth: true,
    validate(route) {
        const params = route.params as { view?: string; date?: string };
        const view = params.view;
        const date = params.date;

        if (typeof view !== 'string' || !isCalendarView(view)) return false;
        if (typeof date !== 'string') return false;

        return !!parseRouteDate(date ?? '');
    },
});

const route = useRoute();
const router = useRouter();
const overlay = useOverlay();

const { can } = useAuth();
const { t } = useI18n();

const calendarStore = useCalendarStore();
const { activeCalendarIds, currentDate, calendars, entries } = storeToRefs(calendarStore);

const calendarViewSlideover = overlay.create(ViewSlideover);
const calendarCreateOrUpdateModal = overlay.create(CreateOrUpdateModal);
const entryViewModal = overlay.create(EntryViewModal);
const entryCreateOrUpdateModal = overlay.create(EntryCreateOrUpdateModal);
const findCalendarsDrawer = overlay.create(FindCalendarDrawer);
const confirmModal = overlay.create(ConfirmModal);
const notifications = useNotificationsStore();

const page = ref(1);

function parseRouteDate(value: string | undefined): Date | undefined {
    const parsed = new Date(`${value}T00:00:00`);
    return Number.isNaN(parsed.getTime()) ? undefined : parsed;
}

const routeParams = computed(() => route.params as { view?: string; date?: string });

const view = computed(() => normalizeCalendarView(routeParams.value.view));
const selectedDate = computed(() => parseRouteDate(routeParams.value.date) ?? new Date());

watch(
    [view, selectedDate],
    () => {
        currentDate.value = {
            year: selectedDate.value.getFullYear(),
            month: selectedDate.value.getMonth() + 1,
        };
    },
    { immediate: true },
);

const title = computed(() => rangeTitle(view.value, selectedDate.value));

const writableCalendar = computed(() => {
    return calendars.value.find(
        (calendar) =>
            activeCalendarIds.value.includes(calendar.id) &&
            checkCalendarAccess(calendar.access, calendar.creator, AccessLevel.EDIT, calendar.job, calendar.creatorJob),
    );
});

const canCreateEntry = computed(() => writableCalendar.value !== undefined);

async function openCreateEntry(): Promise<void> {
    if (!writableCalendar.value) return;

    const response = await entryCreateOrUpdateModal.open({
        calendarId: writableCalendar.value.id,
    });
    if (response) {
        await loadEntries();
    }
}

async function openCreateEntryAt(range: { startTime: Date; endTime: Date }): Promise<void> {
    if (!writableCalendar.value) return;

    const response = await entryCreateOrUpdateModal.open({
        calendarId: writableCalendar.value.id,
        startTime: range.startTime,
        endTime: range.endTime,
    });

    if (response) {
        await loadEntries();
    }
}

async function openEditEntry(entry: CalendarEntry): Promise<void> {
    const response = await entryCreateOrUpdateModal.open({ entryId: entry.id, calendarId: entry.calendarId });
    if (response) {
        await loadEntries();
    }
}

const createActionItems = computed<DropdownMenuItem[]>(() => {
    const items: DropdownMenuItem[] = [];

    if (can('calendar.CalendarService/CreateCalendar').value) {
        items.push({
            label: t('common.calendar'),
            icon: 'i-mdi-calendar-plus',
            onSelect: async () => {
                const response = await calendarCreateOrUpdateModal.open({});
                if (response) {
                    await loadEntries();
                }
            },
        });
    }

    if (canCreateEntry.value) {
        items.push({
            label: t('common.entry', 1),
            icon: 'i-mdi-calendar-plus',
            onSelect: () => void openCreateEntry(),
        });
    }

    return items;
});
const canCreate = computed(() => createActionItems.value.length > 0);

useHead({
    title: () => `${title.value} - ${t('common.calendar')}`,
});

const {
    data: calendarsData,
    status: calendarsStatus,
    error: calendarsError,
    refresh: calendarsRefresh,
} = useLazyAsyncData(
    () => `calendars:${page.value}`,
    () => listCalendars(page.value),
);

async function listCalendars(currentPage: number): Promise<ListCalendarsResponse> {
    const response = await calendarStore.listCalendars({
        pagination: {
            offset: calculateOffset(currentPage, calendarsData.value?.pagination),
        },
        onlyPublic: false,
        calendarIds: [],
    });

    if (activeCalendarIds.value.length === 0) {
        activeCalendarIds.value = response.calendars.map((calendar) => calendar.id);
    }

    return response;
}

const enablePagination = computed<boolean>(() =>
    calendarsData.value?.pagination
        ? calendarsData.value.pagination.totalCount > calendarsData.value.pagination.pageSize
        : false,
);

const visibleEntries = ref<CalendarEntry[]>([]);
const entriesStatus = ref<'idle' | 'pending' | 'success' | 'error'>('idle');
const entriesError = ref<Error | undefined>();

async function loadEntries(): Promise<void> {
    if (activeCalendarIds.value.length === 0) {
        visibleEntries.value = [];
        entries.value = [];
        return;
    }

    entriesStatus.value = 'pending';
    entriesError.value = undefined;

    try {
        const range = viewDateRange(view.value, selectedDate.value);
        const months = fetchMonthsForRange(range);
        const calendarEntriesClient = await getCalendarEntriesClient();

        const responses = await Promise.all(
            months.map((month) =>
                calendarEntriesClient.listCalendarEntries({
                    year: month.getFullYear(),
                    month: month.getMonth() + 1,
                    calendarIds: activeCalendarIds.value,
                }),
            ),
        );

        const merged = new Map<string, CalendarEntry>();
        for (const response of responses) {
            for (const entry of response.response.entries) {
                const key = entry.occurrence?.key ?? String(entry.id);
                merged.set(key, entry);
            }
        }

        visibleEntries.value = Array.from(merged.values()).sort((left, right) => {
            const leftStart = getCalendarEntryDisplayStartDate(left).getTime();
            const rightStart = getCalendarEntryDisplayStartDate(right).getTime();
            if (leftStart !== rightStart) return leftStart - rightStart;
            if (left.calendarId !== right.calendarId) return left.calendarId - right.calendarId;
            return left.id - right.id;
        });

        entries.value = visibleEntries.value;
        entriesStatus.value = 'success';
    } catch (error) {
        entriesError.value = error as Error;
        entriesStatus.value = 'error';
        handleGRPCError(error as RpcError);
    }
}

watch(
    [view, selectedDate, activeCalendarIds],
    () => {
        void loadEntries();
    },
    { immediate: true, deep: true },
);

const selectedDay = computed(() => selectedDate.value);

function navigateToDate(date: Date, nextView = view.value): void {
    void router.push(unsafeRoute(routeForCalendarView(nextView, date)));
}

function shiftRange(direction: -1 | 1): void {
    if (view.value === 'day') {
        navigateToDate(addDays(selectedDate.value, direction));
        return;
    }

    if (view.value === 'week') {
        navigateToDate(addDays(selectedDate.value, 7 * direction));
        return;
    }

    navigateToDate(addMonths(selectedDate.value, direction));
}

function goToToday(): void {
    navigateToDate(new Date());
}

function changeView(nextView: string): void {
    if (!isCalendarView(nextView)) return;
    navigateToDate(selectedDate.value, nextView);
}

function calendarIdChange(calendarId: number, state: boolean): void {
    if (state) {
        if (!activeCalendarIds.value.includes(calendarId)) activeCalendarIds.value.push(calendarId);
    } else {
        activeCalendarIds.value = activeCalendarIds.value.filter((candidate) => candidate !== calendarId);
    }
}

function openSelectedEntry(entry: CalendarEntry): void {
    entryViewModal.open({ entry });
}

function openSelectedCalendar(calendarId: number): void {
    calendarViewSlideover.open({ calendarId });
}

function copyEntryLink(entry: CalendarEntry): void {
    const url = new URL(window.location.href);

    if (entry.occurrence?.key) {
        url.searchParams.set('entryKey', entry.occurrence.key);
        url.searchParams.delete('entryId');
    } else {
        url.searchParams.set('entryId', String(entry.id));
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

async function openDeleteEntry(entry: CalendarEntry): Promise<void> {
    const response = await confirmModal.open({
        confirm: async () => calendarStore.deleteCalendarEntry(entry.id),
    });

    if (response) {
        await loadEntries();
    }
}

const schema = z.object({
    entryId: z.coerce.number().nonnegative().optional(),
    entryKey: z.string().optional(),
});

const query = useSearchForm('calendar', schema);

function openEntryFromQuery(): void {
    if (query.entryKey) {
        const entry = visibleEntries.value.find((candidate) => candidate.occurrence?.key === query.entryKey);
        if (entry) {
            openSelectedEntry(entry);
            query.entryKey = undefined;
            query.entryId = undefined;
        }
        return;
    }

    if (!query.entryId) return;

    const entry = visibleEntries.value.find((candidate) => candidate.id === query.entryId);
    if (entry) {
        openSelectedEntry(entry);
    } else {
        entryViewModal.open({
            entryId: query.entryId,
        });
    }

    query.entryId = undefined;
    query.entryKey = undefined;
}

watch([toRef(query, 'entryId'), toRef(query, 'entryKey'), visibleEntries], () => openEntryFromQuery(), { deep: true });
onMounted(() => openEntryFromQuery());

const viewOptions = [
    { label: t('common.day'), value: 'day', icon: 'i-mdi-view-day' },
    { label: t('common.week_view'), value: 'week', icon: 'i-mdi-view-week' },
    { label: t('common.monthly_view'), value: 'month', icon: 'i-mdi-view-module' },
    { label: t('common.summary'), value: 'summary', icon: 'i-mdi-view-agenda-outline' },
];

const miniDate = computed<CalendarDate>(() => dateToCalendarDate(selectedDay.value)!);

function handleMiniCalendarSelect(value: DateValue | DateRange | DateValue[] | null | undefined): void {
    if (!value || Array.isArray(value) || 'start' in value || !('year' in value)) return;
    navigateToDate(new Date(value.year, value.month - 1, value.day));
}

const showLoading = computed(() => entriesStatus.value === 'pending' && visibleEntries.value.length === 0);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="title">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #title>
                    <div class="flex items-center gap-2">
                        <div class="flex flex-col gap-1">
                            <span class="truncate text-sm font-semibold">{{ title }}</span>
                            <span class="hidden text-xs text-muted sm:block">
                                {{ $t('common.calendar') }} - {{ viewOptions.find((o) => o.value === view)?.label }}
                            </span>
                        </div>

                        <div v-if="canCreate" class="flex justify-end">
                            <UDropdownMenu :items="createActionItems" arrow :content="{ align: 'end' }">
                                <UButton
                                    color="neutral"
                                    variant="outline"
                                    icon="i-mdi-plus"
                                    square
                                    :aria-label="$t('common.create')"
                                />
                            </UDropdownMenu>
                        </div>
                    </div>
                </template>

                <template #right>
                    <div class="inline-flex items-center gap-2">
                        <UButton icon="i-mdi-chevron-left" color="neutral" variant="ghost" @click="shiftRange(-1)" />
                        <UButton
                            icon="i-mdi-calendar-today"
                            color="neutral"
                            variant="outline"
                            :label="$t('common.today')"
                            @click="goToToday"
                        />
                        <UButton icon="i-mdi-chevron-right" color="neutral" variant="ghost" @click="shiftRange(1)" />
                    </div>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div class="grid h-full min-h-0 min-w-0 xl:grid-cols-[18rem_minmax(0,1fr)]">
                <aside
                    class="flex flex-col gap-4 border-b border-default px-4 py-4 pb-(--page-content-bottom-offset) xl:sticky xl:top-0 xl:h-[calc(100svh-var(--ui-header-height))] xl:border-e xl:border-b-0 xl:px-4"
                >
                    <div class="flex flex-col space-y-3 lg:flex-1">
                        <div class="hidden items-center justify-between lg:flex">
                            <h2 class="text-sm font-semibold">{{ $t('common.calendar', 2) }}</h2>

                            <div class="inline-flex items-center">
                                <UButton
                                    color="neutral"
                                    :label="$t('components.calendar.calendar.FindCalendarDrawer.title')"
                                    size="xs"
                                    trailing-icon="i-mdi-calendar-search"
                                    variant="ghost"
                                    @click="findCalendarsDrawer.open({})"
                                />

                                <RefreshButton
                                    :disabled="isRequestPending(calendarsStatus) || showLoading"
                                    :loading="isRequestPending(calendarsStatus) || showLoading"
                                    icon-only
                                    size="xs"
                                    @click="loadEntries"
                                />
                            </div>
                        </div>

                        <div class="relative hidden lg:block">
                            <div class="space-y-1.5">
                                <div v-for="calendar in calendars" :key="calendar.id" class="flex items-center gap-1">
                                    <UCheckbox
                                        :model-value="activeCalendarIds.includes(calendar.id)"
                                        :color="stringToButtonColor(calendar.color)"
                                        @update:model-value="($event) => calendarIdChange(calendar.id, $event as boolean)"
                                    />

                                    <UButton
                                        class="min-w-0 flex-1 justify-start truncate"
                                        :color="stringToButtonColor(calendar.color)"
                                        :icon="calendar.deletedAt ? 'i-mdi-delete' : undefined"
                                        variant="ghost"
                                        size="xs"
                                        @click="openSelectedCalendar(calendar.id)"
                                    >
                                        <UBadge
                                            :color="calendar.color as BadgeProps['color']"
                                            size="xs"
                                            :icon="calendar.icon ? convertComponentIconNameToDynamic(calendar.icon) : undefined"
                                            :class="calendar.icon ? '' : 'size-[14px]'"
                                        />
                                        <span class="truncate">{{ calendar.name }}</span>
                                    </UButton>
                                </div>
                            </div>

                            <div
                                v-if="isRequestPending(calendarsStatus) || calendarsError"
                                class="absolute inset-0 z-10 flex items-center justify-center rounded-lg bg-elevated/70 px-2 py-4 backdrop-blur-[1px]"
                            >
                                <DataPendingBlock
                                    v-if="isRequestPending(calendarsStatus)"
                                    :message="$t('common.loading', [$t('common.calendar')])"
                                />
                                <DataErrorBlock
                                    v-else
                                    :title="$t('common.unable_to_load', [$t('common.calendar')])"
                                    :error="calendarsError"
                                    :retry="calendarsRefresh"
                                />
                            </div>
                        </div>

                        <UFieldGroup class="flex lg:hidden">
                            <UPopover :content="{ side: 'bottom', align: 'start' }" arrow class="w-full">
                                <UButton
                                    block
                                    class="w-full"
                                    color="neutral"
                                    icon="i-mdi-calendar"
                                    trailing-icon="i-mdi-chevron-down"
                                    :loading="isRequestPending(calendarsStatus)"
                                    :label="$t('common.calendar')"
                                />

                                <template #content>
                                    <div class="p-2">
                                        <DataPendingBlock
                                            v-if="isRequestPending(calendarsStatus)"
                                            class="max-w-60"
                                            :message="$t('common.loading', [$t('common.calendar')])"
                                        />
                                        <DataErrorBlock
                                            v-else-if="calendarsError"
                                            :title="$t('common.unable_to_load', [$t('common.calendar')])"
                                            :error="calendarsError"
                                            :retry="calendarsRefresh"
                                        />
                                        <DataNoDataBlock
                                            v-else-if="!calendars || calendars.length === 0"
                                            :type="$t('common.calendar')"
                                            icon="i-mdi-calendar"
                                        />

                                        <div v-else class="grid grid-cols-1 gap-2 sm:grid-cols-2">
                                            <div
                                                v-for="calendar in calendars"
                                                :key="calendar.id"
                                                class="inline-flex items-center gap-2"
                                            >
                                                <UCheckbox
                                                    class="truncate"
                                                    :model-value="activeCalendarIds.includes(calendar.id)"
                                                    :color="stringToButtonColor(calendar.color)"
                                                    @update:model-value="
                                                        ($event) => calendarIdChange(calendar.id, $event as boolean)
                                                    "
                                                />

                                                <UButton
                                                    :color="stringToButtonColor(calendar.color)"
                                                    :icon="calendar.deletedAt ? 'i-mdi-delete' : undefined"
                                                    size="xs"
                                                    truncate
                                                    :variant="calendar.deletedAt ? 'subtle' : 'solid'"
                                                    @click="calendarViewSlideover.open({ calendarId: calendar.id })"
                                                >
                                                    <UBadge
                                                        :class="calendar.icon ? '' : 'size-[14px]'"
                                                        :color="calendar.color as BadgeProps['color']"
                                                        size="lg"
                                                        :label="calendar.icon ? undefined : ''"
                                                        :icon="
                                                            calendar.icon
                                                                ? convertComponentIconNameToDynamic(calendar.icon)
                                                                : undefined
                                                        "
                                                    />
                                                    <span class="truncate">{{ calendar.name }}</span>
                                                </UButton>
                                            </div>
                                        </div>
                                    </div>
                                </template>
                            </UPopover>

                            <UButton
                                color="neutral"
                                :label="$t('components.calendar.calendar.FindCalendarDrawer.title')"
                                size="xs"
                                trailing-icon="i-mdi-calendar-search"
                                variant="ghost"
                                @click="findCalendarsDrawer.open({})"
                            />

                            <RefreshButton
                                :disabled="isRequestPending(calendarsStatus) || showLoading"
                                :loading="isRequestPending(calendarsStatus) || showLoading"
                                size="xs"
                                @click="loadEntries"
                            />
                        </UFieldGroup>

                        <Pagination
                            v-if="enablePagination"
                            v-model="page"
                            :compact="true"
                            :status="calendarsStatus"
                            :pagination="calendarsData?.pagination"
                            :refresh="calendarsRefresh"
                            hide-text
                            hide-refresh
                        />
                    </div>

                    <UFormField name="view" :ui="{ container: '' }">
                        <ClientOnly>
                            <USelectMenu
                                class="w-full"
                                :model-value="view"
                                :items="viewOptions"
                                :icon="viewOptions.find((o) => o.value === view)?.icon"
                                value-key="value"
                                @update:model-value="($event) => changeView($event)"
                            />
                        </ClientOnly>
                    </UFormField>

                    <USeparator />

                    <ClientOnly>
                        <div class="hidden lg:block">
                            <UCalendar
                                class="w-full"
                                :model-value="miniDate"
                                :week-starts-on="1"
                                :year-controls="false"
                                size="xs"
                                fixed-weeks
                                @update:model-value="handleMiniCalendarSelect"
                            />
                        </div>

                        <UCollapsible class="group flex flex-col gap-2 lg:hidden">
                            <UButton
                                block
                                color="neutral"
                                variant="subtle"
                                icon="i-mdi-calendar"
                                trailing-icon="i-mdi-chevron-down"
                                :label="$t('common.calendar')"
                                :ui="{
                                    trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                                }"
                            />

                            <template #content>
                                <div class="pt-2">
                                    <UCalendar
                                        class="w-full"
                                        :model-value="miniDate"
                                        :week-starts-on="1"
                                        :year-controls="false"
                                        size="xs"
                                        fixed-weeks
                                        @update:model-value="handleMiniCalendarSelect"
                                    />
                                </div>
                            </template>
                        </UCollapsible>
                    </ClientOnly>
                </aside>

                <section class="flex min-h-0 min-w-0 flex-col overflow-hidden pb-(--page-content-bottom-offset) xl:px-0">
                    <div class="relative flex min-h-0 flex-1 flex-col overflow-hidden">
                        <div class="flex min-h-0 flex-1 flex-col overflow-hidden">
                            <DataErrorBlock
                                v-if="entriesStatus === 'error'"
                                :title="$t('common.unable_to_load', [$t('common.entry', 2)])"
                                :error="entriesError"
                                :retry="loadEntries"
                            />

                            <CalendarMonthView
                                v-else-if="view === 'month'"
                                :date="selectedDate"
                                :entries="visibleEntries"
                                :can-create="canCreateEntry"
                                @create="openCreateEntryAt"
                                @select="openSelectedEntry"
                                @edit="openEditEntry"
                                @share="copyEntryLink"
                                @delete="openDeleteEntry"
                            />

                            <CalendarWeekView
                                v-else-if="view === 'week'"
                                :date="selectedDate"
                                :entries="visibleEntries"
                                :can-create="canCreateEntry"
                                @create="openCreateEntryAt"
                                @select="openSelectedEntry"
                                @edit="openEditEntry"
                                @share="copyEntryLink"
                                @delete="openDeleteEntry"
                            />

                            <CalendarDayView
                                v-else-if="view === 'day'"
                                :date="selectedDate"
                                :entries="visibleEntries"
                                :can-create="canCreateEntry"
                                @create="openCreateEntryAt"
                                @select="openSelectedEntry"
                                @edit="openEditEntry"
                                @share="copyEntryLink"
                                @delete="openDeleteEntry"
                            />

                            <CalendarAgendaView
                                v-else
                                :date="selectedDate"
                                :entries="visibleEntries"
                                @select="openSelectedEntry"
                                @edit="openEditEntry"
                                @share="copyEntryLink"
                                @delete="openDeleteEntry"
                            />
                        </div>

                        <div
                            v-if="showLoading"
                            class="absolute inset-0 z-20 flex items-center justify-center bg-default/30 backdrop-blur-[1px]"
                        >
                            <div
                                class="flex items-center gap-3 rounded-full border border-default bg-elevated/80 px-4 py-2 shadow-lg"
                            >
                                <UIcon name="i-mdi-loader-circle" class="size-5 animate-spin text-primary" />
                                <span class="text-sm font-medium text-highlighted">
                                    {{ $t('common.loading', [$t('common.entry', 2)]) }}
                                </span>
                            </div>
                        </div>
                    </div>
                </section>
            </div>
        </template>
    </UDashboardPanel>
</template>
