<script lang="ts" setup>
import { CalendarDate } from '@internationalized/date';
import type { TabsItem } from '@nuxt/ui';
import { addDays, addWeeks, isBefore, isFuture, subDays, subMonths, subWeeks } from 'date-fns';
import { z } from 'zod';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import InputDatePicker from '~/components/partials/InputDatePicker.vue';
import InputDateRangePopover from '~/components/partials/InputDateRangePopover.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/stores/completor';
import { getJobsTimeclockClient } from '~~/gen/ts/clients';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { TimeclockMode, TimeclockViewMode } from '~~/gen/ts/resources/jobs/timeclock';
import type { ListTimeclockRequest, ListTimeclockResponse } from '~~/gen/ts/services/jobs/timeclock';
import ColleagueInfoPopover from '../colleagues/ColleagueInfoPopover.vue';
import ColleagueName from '../colleagues/ColleagueName.vue';
import TimeclockTimeline from './TimeclockTimeline.vue';

const props = withDefaults(
    defineProps<{
        userMode?: TimeclockViewMode;
        userId?: number;
        showStats?: boolean;
        historicSubDays?: number;
        forceHistoricView?: boolean;
        hideDaily?: boolean;
    }>(),
    {
        userMode: undefined,
        userId: undefined,
        showStats: true,
        historicSubDays: 7,
        forceHistoricView: undefined,
        hideDaily: false,
    },
);

const { t } = useI18n();

const { attr } = useAuth();

const completorStore = useCompletorStore();

const jobsTimeclockClient = await getJobsTimeclockClient();

const canAccessAll = attr('jobs.TimeclockService/ListTimeclock', 'Access', 'All');

const route = useRoute();

const dateLowerLimit = new Date(2022, 1, 1);
const today = new Date();
const tomorrow = addDays(today, 1);
const minMonth = subMonths(today, 6);

const schema = z.object({
    viewMode: z
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
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'plate',
                        desc: false,
                    },
                ]),
        })
        .default({ columns: [{ id: 'plate', desc: false }] }),
    page: pageNumberSchema,
});

const query = useSearchForm('jobs_timeclock', schema);

function setFromProps(): void {
    if (props.userId === undefined) {
        return;
    }

    query.viewMode = TimeclockViewMode.ALL;
    query.users = [props.userId];
}

setFromProps();
watch(props, setFromProps);

const colleaguesSearchTerm = ref('');
const colleaguesSearchTermDebounced = refDebounced(colleaguesSearchTerm, 200);

const { data: colleagues, status: colleaguesStatus } = useLazyAsyncData(
    () => `jobs-timeclock-colleagues-${colleaguesSearchTerm.value}-${JSON.stringify(query.users)}`,
    () => completorStore.completeColleagues(colleaguesSearchTermDebounced.value),
    {
        watch: [colleaguesSearchTermDebounced, () => query.users],
        immediate: true,
    },
);

const { data, status, refresh, error } = useLazyAsyncData(
    () =>
        `jobs-timeclock-${JSON.stringify(query.sorting)}-${query.date.start.toDateString()}-${query.date.end.toDateString()}-${query.perDay}-${query.users.join(',')}-${query.page}`,
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
            sort: query.sorting,
            userMode: query.viewMode,
            mode: query.mode,
            date: {
                start: {
                    timestamp: googleProtobufTimestamp.Timestamp.fromDate(
                        query.viewMode === TimeclockViewMode.ALL && query.mode === TimeclockMode.DAILY
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

        const call = jobsTimeclockClient.listTimeclock(req);
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
        accessorKey: 'date',
        header: t('common.date'),
        sortable: true,
        meta: {
            class: {
                td:
                    query.viewMode === TimeclockViewMode.SELF || (query.mode !== TimeclockMode.DAILY && query.perDay)
                        ? ''
                        : 'hidden',
                th:
                    query.viewMode === TimeclockViewMode.SELF || (query.mode !== TimeclockMode.DAILY && query.perDay)
                        ? ''
                        : 'hidden',
            },
        },
    },
    {
        accessorKey: 'name',
        header: t('common.name'),
        sortable: canAccessAll.value && props.userId === undefined,
        meta: {
            class: {
                td: props.userId === undefined && query.viewMode === TimeclockViewMode.ALL ? '' : 'hidden',
                th: props.userId === undefined && query.viewMode === TimeclockViewMode.ALL ? '' : 'hidden',
            },
        },
    },
    {
        accessorKey: 'rank',
        header: t('common.rank'),
        sortable: true,
        meta: {
            class: {
                td:
                    canAccessAll.value &&
                    query.viewMode === TimeclockViewMode.ALL &&
                    (query.users === undefined || query.users?.length === 0) &&
                    props.userId === undefined
                        ? ''
                        : 'hidden',
                th:
                    canAccessAll.value &&
                    query.viewMode === TimeclockViewMode.ALL &&
                    (query.users === undefined || query.users?.length === 0) &&
                    props.userId === undefined
                        ? ''
                        : 'hidden',
            },
        },
    },
    {
        accessorKey: 'time',
        header: t('common.time'),
        sortable: true,
    },
]);

const items = computed<TabsItem[]>(() =>
    [
        {
            slot: 'self' as const,
            label: t('common.own'),
            icon: 'i-mdi-selfie',
            value: TimeclockViewMode.SELF,
            viewMode: TimeclockViewMode.SELF,
        },
        canAccessAll.value
            ? {
                  slot: 'colleagues' as const,
                  label: t('components.jobs.timeclock.colleagues'),
                  icon: 'i-mdi-account-group',
                  viewMode: TimeclockViewMode.ALL,
                  value: TimeclockViewMode.ALL,
              }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const timeRangeModes = computed<TabsItem[]>(() => [
    { label: t('common.day_view'), icon: 'i-mdi-view-day', value: TimeclockMode.DAILY, mode: TimeclockMode.DAILY },
    { label: t('common.week_view'), icon: 'i-mdi-view-week', value: TimeclockMode.WEEKLY, mode: TimeclockMode.WEEKLY },
    { label: t('common.time_range'), icon: 'i-mdi-calendar-range', value: TimeclockMode.RANGE, mode: TimeclockMode.RANGE },
    { label: t('common.timeline'), icon: 'i-mdi-chart-timeline', value: TimeclockMode.TIMELINE, mode: TimeclockMode.TIMELINE },
]);

const { game } = useAppConfig();
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardToolbar>
                <template #default>
                    <UForm :schema="schema" :state="query" class="flex flex-1 flex-col gap-1" @submit="refresh">
                        <div class="flex min-w-0 flex-col justify-between lg:flex-row">
                            <UTabs
                                v-if="props.userId === undefined && items.length > 1"
                                v-model="query.viewMode"
                                :items="items"
                                variant="link"
                            />

                            <UTabs
                                v-if="query.viewMode === TimeclockViewMode.SELF"
                                v-model="query.mode"
                                :items="timeRangeModes.filter((m) => m.mode >= TimeclockMode.RANGE)"
                                variant="link"
                            />

                            <UTabs
                                v-else-if="query.viewMode === TimeclockViewMode.ALL"
                                v-model="query.mode"
                                :items="timeRangeModes"
                                variant="link"
                            />
                        </div>

                        <div v-if="query.viewMode === TimeclockViewMode.SELF" class="mb-2 flex flex-1 justify-between gap-2">
                            <UFormField class="flex-1" name="date" :label="$t('common.time_range')">
                                <InputDateRangePopover
                                    v-model="query.date"
                                    class="flex-1"
                                    :min-value="
                                        new CalendarDate(minMonth.getFullYear(), minMonth.getMonth() + 1, minMonth.getDate())
                                    "
                                    :max-value="
                                        new CalendarDate(tomorrow.getFullYear(), tomorrow.getMonth() + 1, tomorrow.getDate())
                                    "
                                />
                            </UFormField>

                            <UFormField
                                v-if="query.mode !== TimeclockMode.TIMELINE"
                                class="flex flex-initial flex-col"
                                name="perDay"
                                :label="$t('common.per_day')"
                                :ui="{ container: 'flex-1 flex' }"
                            >
                                <div class="flex flex-1 items-center">
                                    <USwitch v-model="query.perDay" />
                                </div>
                            </UFormField>
                        </div>

                        <div v-else-if="query.viewMode === TimeclockViewMode.ALL" class="mb-2 flex flex-1 flex-row">
                            <div
                                class="grid flex-1 gap-2"
                                :class="canAccessAll && userId === undefined ? 'grid-cols-2' : 'grid-cols-1'"
                            >
                                <UFormField
                                    v-if="canAccessAll && userId === undefined"
                                    name="users"
                                    :label="$t('common.search')"
                                >
                                    <USelectMenu
                                        v-model="query.users"
                                        v-model:search-term="colleaguesSearchTerm"
                                        :items="colleagues"
                                        multiple
                                        :loading="isRequestPending(colleaguesStatus)"
                                        :search-input="{
                                            placeholder: $t('common.search_field'),
                                        }"
                                        :filter-fields="['firstname', 'lastname']"
                                        :placeholder="$t('common.colleague', 2)"
                                        ignore-filter
                                        leading-icon="i-mdi-search"
                                        value-key="userId"
                                        class="w-full"
                                    >
                                        <template #item-label="{ item }">
                                            <span v-if="item" class="truncate">
                                                {{ userToLabel(item) }}
                                            </span>
                                        </template>

                                        <template #item="{ item }">
                                            <ColleagueName class="truncate" :colleague="item" birthday />
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                                        </template>
                                    </USelectMenu>
                                </UFormField>

                                <div class="flex flex-1 flex-row gap-1">
                                    <UFormField
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
                                        <div v-if="query.mode === TimeclockMode.DAILY" class="flex flex-1 flex-row gap-1">
                                            <UButton
                                                class="flex-initial"
                                                square
                                                icon="i-mdi-chevron-left"
                                                :disabled="isBefore(query.date.end, dateLowerLimit)"
                                                @click="query.date.end = subDays(query.date.end, 1)"
                                            />

                                            <InputDatePicker
                                                v-model="query.date.end"
                                                :min-value="
                                                    new CalendarDate(
                                                        minMonth.getFullYear(),
                                                        minMonth.getMonth() + 1,
                                                        minMonth.getDate(),
                                                    )
                                                "
                                                :max-value="
                                                    new CalendarDate(
                                                        tomorrow.getFullYear(),
                                                        tomorrow.getMonth() + 1,
                                                        tomorrow.getDate(),
                                                    )
                                                "
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

                                            <InputDatePicker
                                                v-model="query.date.end"
                                                :date-format="`yyyy '${$t('common.calendar_week')}' w`"
                                                :min-value="
                                                    new CalendarDate(
                                                        minMonth.getFullYear(),
                                                        minMonth.getMonth() + 1,
                                                        minMonth.getDate(),
                                                    )
                                                "
                                                :max-value="
                                                    new CalendarDate(
                                                        tomorrow.getFullYear(),
                                                        tomorrow.getMonth() + 1,
                                                        tomorrow.getDate(),
                                                    )
                                                "
                                            />

                                            <UButton
                                                class="flex-initial"
                                                square
                                                icon="i-mdi-chevron-right"
                                                :disabled="isFuture(addWeeks(query.date.end, 1))"
                                                @click="query.date.end = addWeeks(query.date.end, 1)"
                                            />
                                        </div>

                                        <InputDateRangePopover
                                            v-else
                                            v-model="query.date"
                                            class="flex-1"
                                            :min-value="
                                                new CalendarDate(
                                                    minMonth.getFullYear(),
                                                    minMonth.getMonth() + 1,
                                                    minMonth.getDate(),
                                                )
                                            "
                                            :max-value="
                                                new CalendarDate(
                                                    tomorrow.getFullYear(),
                                                    tomorrow.getMonth() + 1,
                                                    tomorrow.getDate(),
                                                )
                                            "
                                        />
                                    </UFormField>

                                    <UFormField
                                        v-if="query.mode !== TimeclockMode.DAILY && query.mode !== TimeclockMode.TIMELINE"
                                        class="flex flex-initial flex-col"
                                        name="perDay"
                                        :label="$t('common.per_day')"
                                        :ui="{ container: 'flex-1 flex' }"
                                    >
                                        <div class="flex flex-1 items-center">
                                            <USwitch v-model="query.perDay" />
                                        </div>
                                    </UFormField>
                                </div>
                            </div>
                        </div>
                    </UForm>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <div v-if="error" class="flex-1">
                <DataErrorBlock :title="$t('common.unable_to_load', [$t('common.entry', 2)])" :error="error" :retry="refresh" />
            </div>

            <UCard v-else-if="query.viewMode === TimeclockViewMode.SELF && !query.perDay">
                <p class="mt-2 flex w-full items-center gap-x-2 text-2xl font-semibold tracking-tight text-highlighted">
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
                v-model:sorting="query.sorting.columns"
                :loading="isRequestPending(status)"
                :columns="columns"
                :data="entries"
                :empty="$t('common.not_found', [$t('common.entry', 2)])"
                :sorting-options="{ manualSorting: true }"
                sticky
            >
                <template #caption>
                    <caption>
                        <p class="px-1 text-right">
                            <span class="font-semibold">{{ $t('common.sum') }}:</span>

                            {{ fromSecondsToFormattedDuration(totalTimeSum, { seconds: false }) }}
                        </p>
                    </caption>
                </template>

                <template #date-cell="{ row }">
                    <div class="text-highlighted">
                        {{ $d(toDate(row.original.date), 'date') }}
                    </div>
                </template>

                <template #name-cell="{ row }">
                    <div class="inline-flex items-center gap-1">
                        <ProfilePictureImg
                            :src="row.original.user?.profilePicture"
                            :name="`${row.original.user?.firstname} ${row.original.user?.lastname}`"
                            size="xs"
                        />

                        <ColleagueInfoPopover :user="row.original.user" />
                    </div>
                </template>

                <template #rank-cell="{ row }">
                    {{ row.original.user?.jobGradeLabel }}
                    <template v-if="row.original.user?.job !== game.unemployedJobName">
                        ({{ row.original.user?.jobGrade }})</template
                    >
                </template>

                <template #time-cell="{ row }">
                    {{
                        row.original.spentTime > 0
                            ? fromSecondsToFormattedDuration(Math.round(row.original.spentTime * 60 * 60), {
                                  seconds: false,
                              })
                            : ''
                    }}

                    <UBadge v-if="row.original.startTime !== undefined && row.original.endTime === undefined" color="green">
                        {{ $t('common.active') }}
                    </UBadge>
                </template>
            </UTable>

            <template v-else>
                <div v-if="query.viewMode !== TimeclockViewMode.SELF && query.users.length === 0" class="flex-1">
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
        </template>

        <template #footer>
            <Pagination
                v-model="query.page"
                :pagination="data?.pagination"
                :status="status"
                :refresh="refresh"
                :hide-text="query.mode === TimeclockMode.TIMELINE"
                :hide-buttons="query.mode === TimeclockMode.TIMELINE"
            >
                <template #default>
                    <div>
                        <div v-if="showStats && data && data.stats" class="inline-flex min-w-0 items-center">
                            <UButtonGroup class="truncate">
                                <UButton
                                    color="neutral"
                                    variant="subtle"
                                    icon="i-mdi-account-remove"
                                    :label="$t('common.inactive_colleagues')"
                                    to="/jobs/timeclock/inactive"
                                    :ui="{ leadingIcon: 'hidden sm:block' }"
                                />

                                <UDrawer>
                                    <UButton
                                        :label="$t('common.stats')"
                                        color="neutral"
                                        variant="subtle"
                                        icon="i-mdi-graph-line"
                                        :ui="{ leadingIcon: 'hidden sm:block' }"
                                    />

                                    <template #content>
                                        <div class="p-4">
                                            <LazyJobsTimeclockStatsBlock
                                                :weekly="data?.statsWeekly"
                                                :stats="data?.stats"
                                                :hide-header="true"
                                                :failed="!!error"
                                                :loading="isRequestPending(status)"
                                            />
                                        </div>
                                    </template>
                                </UDrawer>
                            </UButtonGroup>
                        </div>

                        <UTooltip
                            v-if="query.mode === TimeclockMode.TIMELINE"
                            :text="$t('components.jobs.timeclock.timeline.tooltip')"
                            :kbds="['CTRL', '🖱']"
                        >
                            <UIcon class="size-4" name="i-mdi-information-outline" />
                        </UTooltip>
                    </div>
                </template>
            </Pagination>
        </template>
    </UDashboardPanel>
</template>
