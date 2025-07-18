<script lang="ts" setup>
import { addDays, addWeeks, isBefore, isFuture, subDays, subMonths, subWeeks } from 'date-fns';
import { z } from 'zod';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DatePickerPopoverClient from '~/components/partials/DatePickerPopover.client.vue';
import DateRangePickerPopoverClient from '~/components/partials/DateRangePickerPopover.client.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/stores/completor';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import { TimeclockMode, TimeclockViewMode } from '~~/gen/ts/resources/jobs/timeclock';
import type { ListTimeclockRequest, ListTimeclockResponse } from '~~/gen/ts/services/jobs/timeclock';
import ColleagueInfoPopover from '../colleagues/ColleagueInfoPopover.vue';
import ColleagueName from '../colleagues/ColleagueName.vue';
import TimeclockStatsBlock from './TimeclockStatsBlock.vue';
import TimeclockTimeline from './TimeclockTimeline.vue';

const props = withDefaults(
    defineProps<{
        userId?: number;
        showStats?: boolean;
        historicSubDays?: number;
        forceHistoricView?: boolean;
        hideDaily?: boolean;
    }>(),
    {
        userId: undefined,
        showStats: true,
        historicSubDays: 7,
        forceHistoricView: undefined,
        hideDaily: false,
    },
);

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { attr, can } = useAuth();

const completorStore = useCompletorStore();

const canAccessAll = attr('jobs.TimeclockService/ListTimeclock', 'Access', 'All');

const route = useRoute();

const dateLowerLimit = new Date(2022, 1, 1);

const schema = z.object({
    userMode: z
        .nativeEnum(TimeclockViewMode)
        .default(
            TimeclockViewMode[(route.query?.mode as string | undefined)?.toUpperCase() as keyof typeof TimeclockViewMode] ??
                TimeclockViewMode.SELF,
        ),
    mode: z
        .nativeEnum(TimeclockMode)
        .default(
            TimeclockMode[(route.query?.view as string | undefined)?.toUpperCase() as keyof typeof TimeclockMode] ??
                (props.hideDaily ? TimeclockMode.WEEKLY : TimeclockMode.RANGE),
        ),
    users: z.coerce.number().array().max(10).default([]),
    date: z
        .object({
            start: z.coerce.date(),
            end: z.coerce.date(),
        })
        .default({
            start: subDays(new Date(), 7),
            end: new Date(),
        }),
    perDay: z.coerce.boolean().default(true),
    sort: z.custom<TableSortable>().default({
        column: 'id',
        direction: 'desc',
    }),
    page: pageNumberSchema,
});

const query = useSearchForm('jobs_timeclock', schema);

function setFromProps(): void {
    if (props.userId === undefined) {
        return;
    }

    query.userMode = TimeclockViewMode.ALL;
    query.users = [props.userId];
}

setFromProps();
watch(props, setFromProps);

const usersLoading = ref(false);

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `jobs-timeclock-${query.sort.column}:${query.sort.direction}-${query.date.start.toDateString()}-${query.date.end.toDateString()}-${query.perDay}-${query.users.join(',')}-${query.page}`,
    () => listTimeclockEntries(),
);

async function listTimeclockEntries(): Promise<ListTimeclockResponse> {
    try {
        if (!isBefore(query.date.start, query.date.end)) {
            query.date.start = query.mode > TimeclockMode.DAILY ? subWeeks(query.date.end, 1) : query.date.start;
        }

        const req: ListTimeclockRequest = {
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sort,
            userMode: query.userMode,
            mode: query.mode,
            date: {
                start: {
                    timestamp: googleProtobufTimestamp.Timestamp.fromDate(
                        query.userMode === TimeclockViewMode.ALL && query.mode === TimeclockMode.DAILY
                            ? query.date.end
                            : query.date.start,
                    ),
                },
                end: {
                    timestamp: googleProtobufTimestamp.Timestamp.fromDate(query.date.end),
                },
            },
            userIds: query.users,
            perDay: query.perDay,
        };

        const call = $grpc.jobs.timeclock.listTimeclock(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

const entries = computed(() => {
    if (data.value?.entries.oneofKind === 'daily') {
        return data.value.entries.daily.entries;
    } else if (data.value?.entries.oneofKind === 'weekly') {
        return data.value.entries.weekly.entries;
    } else if (data.value?.entries.oneofKind === 'range') {
        return data.value.entries.range.entries;
    }

    return [];
});

const totalTimeSum = computed(() => {
    let sum = 0;
    if (data.value?.entries.oneofKind === 'daily') {
        sum = data.value.entries.daily.sum;
    } else if (data.value?.entries.oneofKind === 'weekly') {
        sum = data.value.entries.weekly.sum;
    } else if (data.value?.entries.oneofKind === 'range') {
        sum = data.value.entries.range.sum;
    }

    return sum;
});

const columns = computed(() => [
    {
        key: 'date',
        label: t('common.date'),
        sortable: true,
        class:
            query.userMode === TimeclockViewMode.SELF || (query.mode !== TimeclockMode.DAILY && query.perDay) ? '' : 'hidden',
        rowClass:
            query.userMode === TimeclockViewMode.SELF || (query.mode !== TimeclockMode.DAILY && query.perDay) ? '' : 'hidden',
    },
    {
        key: 'name',
        label: t('common.name'),
        sortable: canAccessAll.value && props.userId === undefined,
        class: props.userId === undefined && query.userMode === TimeclockViewMode.ALL ? '' : 'hidden',
        rowClass: props.userId === undefined && query.userMode === TimeclockViewMode.ALL ? '' : 'hidden',
    },
    {
        key: 'rank',
        label: t('common.rank'),
        sortable: true,
        class:
            canAccessAll.value &&
            query.userMode === TimeclockViewMode.ALL &&
            (query.users === undefined || query.users?.length === 0) &&
            props.userId === undefined
                ? ''
                : 'hidden',
        rowClass:
            canAccessAll.value &&
            query.userMode === TimeclockViewMode.ALL &&
            (query.users === undefined || query.users?.length === 0) &&
            props.userId === undefined
                ? ''
                : 'hidden',
    },
    {
        key: 'time',
        label: t('common.time'),
        sortable: true,
    },
]);

const items = computed(() =>
    [
        { slot: 'self', label: t('common.own'), icon: 'i-mdi-selfie', userMode: TimeclockViewMode.SELF },
        canAccessAll.value
            ? {
                  slot: 'colleagues',
                  label: t('components.jobs.timeclock.colleagues'),
                  icon: 'i-mdi-account-group',
                  userMode: TimeclockViewMode.ALL,
              }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const selectedUserMode = computed({
    get() {
        const index = items.value.findIndex((item) => item.userMode === query.userMode);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        query.userMode = items.value[value]?.userMode ?? TimeclockViewMode.SELF;

        // Select range mode when user self mode is selected
        if (query.userMode === TimeclockViewMode.SELF && query.mode < TimeclockMode.RANGE) {
            query.mode = TimeclockMode.RANGE;
        }
    },
});

const selfTimeRangeModes = computed(() => [
    { label: t('common.time_range'), icon: 'i-mdi-calendar-range', mode: TimeclockMode.RANGE },
    { label: t('common.timeline'), icon: 'i-mdi-chart-timeline', mode: TimeclockMode.TIMELINE },
]);

const selectedSelfMode = computed({
    get() {
        const index = selfTimeRangeModes.value.findIndex((item) => item.mode === query.mode);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        query.mode = selfTimeRangeModes.value[value]?.mode ?? TimeclockMode.RANGE;
    },
});

const timeRangeModes = computed(() => [
    { label: t('common.day_view'), icon: 'i-mdi-view-day', mode: TimeclockMode.DAILY },
    { label: t('common.week_view'), icon: 'i-mdi-view-week', mode: TimeclockMode.WEEKLY },
    { label: t('common.time_range'), icon: 'i-mdi-calendar-range', mode: TimeclockMode.RANGE },
    { label: t('common.timeline'), icon: 'i-mdi-chart-timeline', mode: TimeclockMode.TIMELINE },
]);

const selectedMode = computed({
    get() {
        const index = timeRangeModes.value.findIndex((item) => item.mode === query.mode);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        query.mode = timeRangeModes.value[value]?.mode ?? TimeclockMode.DAILY;
    },
});

const { game } = useAppConfig();
</script>

<template>
    <UTabs
        v-model="selectedUserMode"
        :items="items"
        :ui="{
            list: { base: props.userId !== undefined || items.length === 1 ? 'hidden' : undefined, rounded: '' },
        }"
    />

    <UDashboardToolbar>
        <UForm class="flex w-full flex-col gap-2" :schema="schema" :state="query" @submit="refresh()">
            <template v-if="query.userMode === TimeclockViewMode.SELF">
                <div class="flex flex-1 flex-col justify-between gap-2 sm:flex-row">
                    <UTabs
                        v-model="selectedSelfMode"
                        :items="selfTimeRangeModes.filter((m) => m.mode >= TimeclockMode.RANGE)"
                        :ui="{ wrapper: 'relative space-y-0 flex-1', container: '' }"
                    />
                </div>

                <div class="flex flex-1 justify-between gap-2">
                    <UFormGroup class="flex-1" name="date" :label="$t('common.time_range')">
                        <DateRangePickerPopoverClient
                            v-model="query.date"
                            class="flex-1"
                            mode="date"
                            :popover="{ class: 'flex-1' }"
                            :date-picker="{
                                disabledDates: [
                                    { start: addDays(new Date(), 1), end: null },
                                    { end: subMonths(new Date(), 6) },
                                ],
                            }"
                        />
                    </UFormGroup>

                    <UFormGroup
                        v-if="query.mode !== TimeclockMode.TIMELINE"
                        class="flex flex-initial flex-col"
                        name="perDay"
                        :label="$t('common.per_day')"
                        :ui="{ container: 'flex-1 flex' }"
                    >
                        <div class="flex flex-1 items-center">
                            <UToggle v-model="query.perDay" />
                        </div>
                    </UFormGroup>
                </div>
            </template>

            <template v-if="query.userMode === TimeclockViewMode.ALL">
                <div class="flex flex-1 flex-col justify-between gap-2 sm:flex-row">
                    <UTabs
                        v-model="selectedMode"
                        :items="timeRangeModes"
                        :ui="{ wrapper: 'relative space-y-0 flex-1', container: '' }"
                    />

                    <div class="flex items-center">
                        <UButton
                            v-if="can('jobs.TimeclockService/ListInactiveEmployees').value && userId === undefined"
                            :to="{ name: 'jobs-timeclock-inactive' }"
                            color="black"
                            trailing-icon="i-mdi-arrow-right"
                        >
                            {{ $t('common.inactive_colleagues') }}
                        </UButton>
                    </div>
                </div>

                <div class="flex w-full flex-row">
                    <div
                        class="grid flex-1 gap-2"
                        :class="canAccessAll && userId === undefined ? 'grid-cols-2' : 'grid-cols-1'"
                    >
                        <UFormGroup v-if="canAccessAll && userId === undefined" name="users" :label="$t('common.search')">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.users"
                                    multiple
                                    :searchable="
                                        async (q: string) => {
                                            usersLoading = true;
                                            const colleagues = await completorStore.listColleagues({
                                                search: q,
                                                labelIds: [],
                                                userIds: query.users,
                                            });
                                            usersLoading = false;
                                            return colleagues;
                                        }
                                    "
                                    searchable-lazy
                                    :searchable-placeholder="$t('common.search_field')"
                                    :search-attributes="['firstname', 'lastname']"
                                    :placeholder="$t('common.colleague', 2)"
                                    block
                                    trailing
                                    leading-icon="i-mdi-search"
                                    value-attribute="userId"
                                >
                                    <template #label="{ selected }">
                                        <span v-if="selected.length > 0" class="truncate">
                                            {{ usersToLabel(selected) }}
                                        </span>
                                    </template>

                                    <template #option="{ option: colleague }">
                                        <ColleagueName class="truncate" :colleague="colleague" birthday />
                                    </template>

                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormGroup>

                        <div class="flex flex-1 flex-row gap-1">
                            <UFormGroup
                                class="flex-1"
                                name="end"
                                :label="
                                    query.mode === TimeclockMode.WEEKLY
                                        ? $t('common.week_view')
                                        : query.mode === TimeclockMode.DAILY
                                          ? $t('common.day_view')
                                          : $t('common.time_range')
                                "
                            >
                                <div v-if="query.mode === TimeclockMode.DAILY" class="flex flex-1 flex-col gap-1 sm:flex-row">
                                    <UButton
                                        class="flex-initial"
                                        square
                                        icon="i-mdi-chevron-left"
                                        :disabled="isBefore(query.date.end, dateLowerLimit)"
                                        @click="query.date.end = subDays(query.date.end, 1)"
                                    />

                                    <DatePickerPopoverClient
                                        v-model="query.date.end"
                                        :popover="{ class: 'flex-1' }"
                                        :date-picker="{
                                            disabledDates: [
                                                { start: addDays(new Date(), 1), end: null },
                                                { end: subMonths(new Date(), 6) },
                                            ],
                                        }"
                                    />

                                    <UButton
                                        class="flex-initial"
                                        square
                                        icon="i-mdi-chevron-right"
                                        :disabled="isFuture(addDays(query.date.end, 1))"
                                        @click="query.date.end = addDays(query.date.end, 1)"
                                    />
                                </div>
                                <div
                                    v-else-if="query.mode === TimeclockMode.WEEKLY"
                                    class="flex flex-1 flex-col gap-1 sm:flex-row"
                                >
                                    <UButton
                                        class="flex-initial"
                                        square
                                        icon="i-mdi-chevron-left"
                                        :disabled="isBefore(query.date.end, dateLowerLimit)"
                                        @click="query.date.end = subWeeks(query.date.end, 1)"
                                    />

                                    <DatePickerPopoverClient
                                        v-model="query.date.end"
                                        :popover="{ class: 'flex-1' }"
                                        :date-format="`yyyy '${$t('common.calendar_week')}' w`"
                                        :date-picker="{
                                            disabledDates: [
                                                { start: addDays(new Date(), 1), end: null },
                                                { end: subMonths(new Date(), 6) },
                                            ],
                                        }"
                                    />

                                    <UButton
                                        class="flex-initial"
                                        square
                                        icon="i-mdi-chevron-right"
                                        :disabled="isFuture(addWeeks(query.date.end, 1))"
                                        @click="query.date.end = addWeeks(query.date.end, 1)"
                                    />
                                </div>
                                <DateRangePickerPopoverClient
                                    v-else
                                    v-model="query.date"
                                    class="flex-1"
                                    mode="date"
                                    :popover="{ class: 'flex-1' }"
                                    :date-picker="{
                                        disabledDates: [
                                            { start: addDays(new Date(), 1), end: null },
                                            { end: subMonths(new Date(), 6) },
                                        ],
                                    }"
                                />
                            </UFormGroup>

                            <UFormGroup
                                v-if="query.mode !== TimeclockMode.DAILY && query.mode !== TimeclockMode.TIMELINE"
                                class="flex flex-initial flex-col"
                                name="perDay"
                                :label="$t('common.per_day')"
                                :ui="{ container: 'flex-1 flex' }"
                            >
                                <div class="flex flex-1 items-center">
                                    <UToggle v-model="query.perDay" />
                                </div>
                            </UFormGroup>
                        </div>
                    </div>
                </div>
            </template>
        </UForm>
    </UDashboardToolbar>

    <div v-if="error" class="flex-1">
        <DataErrorBlock :title="$t('common.unable_to_load', [$t('common.entry', 2)])" :error="error" :retry="refresh" />
    </div>

    <UCard v-else-if="query.userMode === TimeclockViewMode.SELF && !query.perDay">
        <p class="mt-2 flex w-full items-center gap-x-2 text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">
            {{
                totalTimeSum === 0
                    ? $t('common.not_found', [$t('common.entry', 2)])
                    : fromSecondsToFormattedDuration(totalTimeSum, {
                          seconds: false,
                      })
            }}
        </p>
    </UCard>

    <UTable
        v-else-if="query.mode !== TimeclockMode.TIMELINE"
        v-model:sort="query.sort"
        class="flex-1"
        :loading="loading"
        :columns="columns"
        :rows="entries"
        :empty-state="{
            icon: 'i-mdi-timeline-clock',
            label: $t('common.not_found', [$t('common.entry', 2)]),
        }"
        sort-mode="manual"
    >
        <template #caption>
            <caption>
                <p class="px-1 text-right">
                    <span class="font-semibold">{{ $t('common.sum') }}:</span>

                    {{ fromSecondsToFormattedDuration(totalTimeSum, { seconds: false }) }}
                </p>
            </caption>
        </template>

        <template #date-data="{ row: entry }">
            <div class="text-gray-900 dark:text-white">
                {{ $d(toDate(entry.date), 'date') }}
            </div>
        </template>

        <template #name-data="{ row: entry }">
            <div class="inline-flex items-center gap-1">
                <ProfilePictureImg
                    :src="entry.user?.avatar"
                    :name="`${entry.user?.firstname} ${entry.user?.lastname}`"
                    size="xs"
                />

                <ColleagueInfoPopover :user="entry.user" />
            </div>
        </template>

        <template #rank-data="{ row: entry }">
            {{ entry.user.jobGradeLabel }}
            <template v-if="entry.user.job !== game.unemployedJobName"> ({{ entry.user.jobGrade }})</template>
        </template>

        <template #time-data="{ row: entry }">
            {{
                entry.spentTime > 0
                    ? fromSecondsToFormattedDuration(Math.round(entry.spentTime * 60 * 60), {
                          seconds: false,
                      })
                    : ''
            }}

            <UBadge v-if="entry.startTime !== undefined && entry.endTime === undefined" color="green">
                {{ $t('common.active') }}
            </UBadge>
        </template>
    </UTable>

    <template v-else>
        <div v-if="query.userMode !== TimeclockViewMode.SELF && query.users.length === 0" class="flex-1">
            <DataNoDataBlock :description="$t('components.jobs.timeclock.timeline.select_users')" />
        </div>

        <TimeclockTimeline v-else :data="entries" :from="query.date.start" :to="query.date.end">
            <template #caption>
                <p class="shrink-0 text-right">
                    <span class="font-semibold">{{ $t('common.sum') }}:</span>

                    {{
                        fromSecondsToFormattedDuration(totalTimeSum, {
                            seconds: false,
                        })
                    }}
                </p>
            </template>
        </TimeclockTimeline>
    </template>

    <div class="flex flex-row items-center">
        <Pagination
            v-model="query.page"
            class="flex-1"
            :pagination="data?.pagination"
            :loading="loading"
            :refresh="refresh"
            :hide-text="query.mode === TimeclockMode.TIMELINE"
            :hide-buttons="query.mode === TimeclockMode.TIMELINE"
        >
            <template #default>
                <div>
                    <UTooltip
                        v-if="query.mode === TimeclockMode.TIMELINE"
                        :text="$t('components.jobs.timeclock.timeline.tooltip')"
                        :shortcuts="['CTRL', '🖱']"
                        :popper="{ placement: 'left' }"
                        :ui="{ shortcuts: 'inline-flex' }"
                    >
                        <UIcon class="size-4" name="i-mdi-information-outline" />
                    </UTooltip>
                </div>
            </template>
        </Pagination>
    </div>

    <UAccordion
        v-if="showStats && data && data.stats"
        class="px-3 py-0.5"
        :items="[{ slot: 'stats', label: $t('common.stats') }]"
    >
        <template #stats>
            <TimeclockStatsBlock
                :weekly="data?.statsWeekly"
                :stats="data?.stats"
                :hide-header="true"
                :failed="!!error"
                :loading="loading"
                :ui="{ rounded: '' }"
            />
        </template>
    </UAccordion>
</template>
