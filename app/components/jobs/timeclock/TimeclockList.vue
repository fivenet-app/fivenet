<script lang="ts" setup>
import { addDays, addWeeks, isBefore, isFuture, subDays, subWeeks } from 'date-fns';
import { z } from 'zod';
import DatePickerPopoverClient from '~/components/partials/DatePickerPopover.client.vue';
import DateRangePickerPopoverClient from '~/components/partials/DateRangePickerPopover.client.vue';
import Pagination from '~/components/partials/Pagination.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useCompletorStore } from '~/store/completor';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import { TimeclockMode, TimeclockUserMode } from '~~/gen/ts/resources/jobs/timeclock';
import type { ListTimeclockRequest, ListTimeclockResponse } from '~~/gen/ts/services/jobs/timeclock';
import ColleagueInfoPopover from '../colleagues/ColleagueInfoPopover.vue';
import TimeclockStatsBlock from './TimeclockStatsBlock.vue';

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

const { t } = useI18n();

const completorStore = useCompletorStore();

const canAccessAll = attr('JobsTimeclockService.ListTimeclock', 'Access', 'All');

const dateLowerLimit = new Date(2022, 1, 1);

const schema = z.object({
    userMode: z.nativeEnum(TimeclockUserMode),
    mode: z.nativeEnum(TimeclockMode),
    users: z.custom<Colleague>().array().max(5).optional(),
    date: z.object({
        start: z.date(),
        end: z.date(),
    }),
    perDay: z.boolean(),
});

type Schema = z.output<typeof schema>;

const route = useRoute();

const query = reactive<Schema>({
    userMode:
        TimeclockUserMode[(route.query?.mode as string | undefined)?.toUpperCase() as keyof typeof TimeclockUserMode] ??
        TimeclockUserMode.SELF,
    mode:
        TimeclockMode[(route.query?.view as string | undefined)?.toUpperCase() as keyof typeof TimeclockMode] ??
        TimeclockMode.WEEKLY,
    users: [],
    date: {
        start: subDays(new Date(), 7),
        end: new Date(),
    },
    perDay: true,
});

function setFromProps(): void {
    if (props.userId === undefined) {
        return;
    }

    query.userMode = TimeclockUserMode.ALL;
    query.users = [
        {
            userId: props.userId,
            firstname: '',
            lastname: '',
            dateofbirth: '',
            job: '',
            jobGrade: 0,
        },
    ];
}

watch(props, setFromProps);
setFromProps();

const usersLoading = ref(false);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const sort = ref<TableSortable>({
    column: 'date',
    direction: 'desc',
});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `jobs-timeclock-${sort.value.column}:${sort.value.direction}-${query.date.start.toDateString()}-${query.date.end.toDateString()}-${query.perDay}-${query.users?.map((u) => u.userId)}-${page.value}`,
    () => listTimeclockEntries(),
    {
        watch: [sort],
    },
);

async function listTimeclockEntries(): Promise<ListTimeclockResponse> {
    try {
        const req: ListTimeclockRequest = {
            pagination: {
                offset: offset.value,
            },
            sort: sort.value,
            userMode: query.userMode,
            mode: query.userMode === TimeclockUserMode.SELF ? TimeclockMode.RANGE : query.mode,
            date: {
                start: {
                    timestamp: googleProtobufTimestamp.Timestamp.fromDate(
                        query.mode === TimeclockMode.DAILY ? query.date.end : query.date.start,
                    ),
                },
                end: {
                    timestamp: googleProtobufTimestamp.Timestamp.fromDate(query.date.end),
                },
            },
            userIds: query.users?.map((u) => u.userId) ?? [],
            perDay: query.perDay,
        };

        const call = getGRPCJobsTimeclockClient().listTimeclock(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

const columns = computed(() =>
    [
        (query.mode !== TimeclockMode.DAILY && query.perDay) || query.userMode === TimeclockUserMode.SELF
            ? {
                  key: 'date',
                  label: t('common.date'),
                  sortable: true,
              }
            : undefined,
        props.userId === undefined && query.userMode === TimeclockUserMode.ALL
            ? {
                  key: 'name',
                  label: t('common.name'),
                  sortable: canAccessAll.value && props.userId === undefined,
              }
            : undefined,
        canAccessAll.value &&
        query.userMode === TimeclockUserMode.ALL &&
        query.users !== undefined &&
        query.users?.length > 0 &&
        props.userId === undefined
            ? {
                  key: 'rank',
                  label: t('common.rank'),
                  sortable: true,
              }
            : undefined,
        {
            key: 'time',
            label: t('common.time'),
            sortable: true,
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const selfPerDaySum = ref<number>(0);

const router = useRouter();

watch(query, () =>
    router.replace({
        query: {
            mode: TimeclockUserMode[query.userMode].toLowerCase(),
            view: TimeclockMode[query.mode].toLowerCase(),
        },
        hash: '#',
    }),
);

watch(data, () => {
    if (data.value?.entries.oneofKind !== 'range') {
        selfPerDaySum.value = 0;
        return;
    }

    selfPerDaySum.value = parseFloat(
        ((Math.round((data.value?.entries.range.entries[0]?.spentTime ?? 0) * 100) / 100) * 60 * 60).toPrecision(2),
    );
});

const selectedTab = computed({
    get() {
        const index = items.value.findIndex((item) => item.userMode === query.userMode);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        query.userMode = items.value[value]?.userMode ?? TimeclockUserMode.SELF;
    },
});

const items = computed(() =>
    [
        { slot: 'self', label: t('common.own'), icon: 'i-mdi-selfie', userMode: TimeclockUserMode.SELF },
        canAccessAll.value
            ? {
                  slot: 'colleagues',
                  label: t('components.jobs.timeclock.colleagues'),
                  icon: 'i-mdi-account-group',
                  userMode: TimeclockUserMode.ALL,
              }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const input = ref<{ input: HTMLInputElement }>();
</script>

<template>
    <div class="flex flex-1">
        <UTabs
            v-model="selectedTab"
            :items="items"
            class="flex-1"
            :unmount="true"
            :ui="{
                list: { base: props.userId !== undefined || items.length === 1 ? 'hidden' : undefined, rounded: '' },
                wrapper: 'flex-1 flex flex-col',
                container: 'relative w-full flex-1 flex',
                base: 'flex-1 flex flex-col',
            }"
        >
            <template #self>
                <UDashboardToolbar>
                    <template #default>
                        <UForm :schema="schema" :state="query" class="flex w-full flex-col gap-2" @submit="refresh()">
                            <div class="flex flex-1 justify-between gap-2">
                                <UFormGroup name="end" :label="$t('common.time_range')" class="flex-1">
                                    <DateRangePickerPopoverClient
                                        v-model="query.date"
                                        mode="date"
                                        class="flex-1"
                                        :popover="{ class: 'flex-1' }"
                                        disable-future
                                    />
                                </UFormGroup>

                                <UFormGroup name="day" :label="$t('common.per_day')">
                                    <UToggle v-model="query.perDay" class="mt-2" />
                                </UFormGroup>
                            </div>
                        </UForm>
                    </template>
                </UDashboardToolbar>

                <UCard v-if="!query.perDay">
                    <p
                        class="mt-2 flex w-full items-center gap-x-2 text-2xl font-semibold tracking-tight text-gray-900 dark:text-white"
                    >
                        {{
                            selfPerDaySum === 0
                                ? $t('common.not_found', [$t('common.entry', 2)])
                                : fromSecondsToFormattedDuration(selfPerDaySum, {
                                      seconds: false,
                                  })
                        }}
                    </p>
                </UCard>

                <UTable
                    v-else
                    v-model:sort="sort"
                    :loading="loading"
                    :columns="columns"
                    :rows="data?.entries.oneofKind === 'range' ? data.entries.range.entries : []"
                    :empty-state="{
                        icon: 'i-mdi-timeline-clock',
                        label: $t('common.not_found', [$t('common.entry', 2)]),
                    }"
                    sort-mode="manual"
                    class="flex-1"
                    :ui="{ caption: '' }"
                >
                    <template #caption>
                        <caption>
                            <p class="text-right">
                                <span class="font-semibold">{{ $t('common.sum') }}:</span>

                                {{
                                    fromSecondsToFormattedDuration(
                                        parseFloat(
                                            (
                                                (Math.round(
                                                    (data?.entries.oneofKind === 'range' ? data?.entries.range.sum : 0) * 100,
                                                ) /
                                                    100) *
                                                60 *
                                                60
                                            ).toPrecision(2),
                                        ),
                                        { seconds: false },
                                    )
                                }}
                            </p>
                        </caption>
                    </template>

                    <template #date-data="{ row: item }">
                        <div class="text-gray-900 dark:text-white">
                            {{ $d(toDate(item.date), 'date') }}
                        </div>
                    </template>

                    <template #rank-data="{ row: entry }">
                        {{ entry.user.jobGradeLabel }}<span v-if="entry.user.jobGrade > 0"> ({{ entry.user.jobGrade }})</span>
                    </template>

                    <template #time-data="{ row: entry }">
                        <div class="text-right">
                            {{
                                entry.spentTime > 0
                                    ? fromSecondsToFormattedDuration(
                                          parseFloat(((Math.round(entry.spentTime * 100) / 100) * 60 * 60).toPrecision(2)),
                                          { seconds: false },
                                      )
                                    : ''
                            }}

                            <UBadge v-if="entry.startTime !== undefined" color="green">
                                {{ $t('common.active') }}
                            </UBadge>
                        </div>
                    </template>
                </UTable>

                <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />

                <TimeclockStatsBlock
                    hide-header
                    :weekly="data?.statsWeekly"
                    :stats="data?.stats"
                    :loading="loading"
                    :ui="{ rounded: '' }"
                />
            </template>

            <template #colleagues>
                <UDashboardToolbar>
                    <template #default>
                        <UForm :schema="schema" :state="query" class="flex w-full flex-col gap-2" @submit="refresh()">
                            <div class="flex flex-1 justify-between gap-2">
                                <UButtonGroup class="inline-flex w-full">
                                    <UButton
                                        color="gray"
                                        block
                                        class="flex-1"
                                        :disabled="query.mode === TimeclockMode.WEEKLY"
                                        trailing-icon="i-mdi-view-week"
                                        @click="query.mode = TimeclockMode.WEEKLY"
                                    >
                                        {{ $t('common.time_ago.week') }}
                                    </UButton>

                                    <UButton
                                        v-if="!hideDaily"
                                        color="gray"
                                        block
                                        class="flex-1"
                                        :disabled="query.mode === TimeclockMode.DAILY"
                                        icon="i-mdi-view-day"
                                        @click="query.mode = TimeclockMode.DAILY"
                                    >
                                        {{ $t('common.time_ago.day') }}
                                    </UButton>

                                    <UButton
                                        color="gray"
                                        block
                                        class="flex-1"
                                        :disabled="query.mode === TimeclockMode.RANGE"
                                        trailing-icon="i-mdi-calendar-range"
                                        @click="query.mode = TimeclockMode.RANGE"
                                    >
                                        {{ $t('common.time_range') }}
                                    </UButton>
                                </UButtonGroup>

                                <UButton
                                    v-if="can('JobsTimeclockService.ListInactiveEmployees').value && userId === undefined"
                                    :to="{ name: 'jobs-timeclock-inactive' }"
                                    color="black"
                                    trailing-icon="i-mdi-arrow-right"
                                >
                                    {{ $t('common.inactive_colleagues') }}
                                </UButton>
                            </div>

                            <div class="flex w-full flex-row">
                                <div
                                    class="grid flex-1 gap-2"
                                    :class="canAccessAll && userId === undefined ? 'grid-cols-2' : 'grid-cols-1'"
                                >
                                    <UFormGroup
                                        v-if="canAccessAll && userId === undefined"
                                        name="users"
                                        :label="$t('common.search')"
                                    >
                                        <ClientOnly>
                                            <USelectMenu
                                                ref="input"
                                                v-model="query.users"
                                                multiple
                                                :searchable="
                                                    async (query: string) => {
                                                        usersLoading = true;
                                                        const colleagues = await completorStore.listColleagues({
                                                            search: query,
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
                                                by="userId"
                                            >
                                                <template #label>
                                                    <span v-if="query.users?.length" class="truncate">
                                                        {{ usersToLabel(query.users) }}
                                                    </span>
                                                </template>
                                                <template #option="{ option: user }">
                                                    <span class="truncate">
                                                        {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                                    </span>
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
                                            name="end"
                                            :label="
                                                query.mode === TimeclockMode.WEEKLY
                                                    ? $t('common.time_ago.week')
                                                    : query.mode === TimeclockMode.DAILY
                                                      ? $t('common.time_ago.day')
                                                      : query.mode === TimeclockMode.RANGE
                                                        ? $t('common.time_range')
                                                        : ''
                                            "
                                            class="flex-1"
                                        >
                                            <div
                                                v-if="query.mode === TimeclockMode.DAILY"
                                                class="flex flex-1 flex-col gap-1 sm:flex-row"
                                            >
                                                <UButton
                                                    square
                                                    icon="i-mdi-chevron-left"
                                                    class="flex-initial"
                                                    :disabled="isBefore(query.date.end, dateLowerLimit)"
                                                    @click="query.date.end = subDays(query.date.end, 1)"
                                                />

                                                <DatePickerPopoverClient
                                                    v-model="query.date.end"
                                                    :popover="{ class: 'flex-1' }"
                                                    disable-future
                                                />

                                                <UButton
                                                    square
                                                    icon="i-mdi-chevron-right"
                                                    class="flex-initial"
                                                    :disabled="isFuture(addDays(query.date.end, 1))"
                                                    @click="query.date.end = addDays(query.date.end, 1)"
                                                />
                                            </div>
                                            <div
                                                v-else-if="query.mode === TimeclockMode.WEEKLY"
                                                class="flex flex-1 flex-col gap-1 sm:flex-row"
                                            >
                                                <UButton
                                                    square
                                                    icon="i-mdi-chevron-left"
                                                    class="flex-initial"
                                                    :disabled="isBefore(query.date.end, dateLowerLimit)"
                                                    @click="query.date.end = subWeeks(query.date.end, 1)"
                                                />

                                                <DatePickerPopoverClient
                                                    v-model="query.date.end"
                                                    :popover="{ class: 'flex-1' }"
                                                    :date-format="`yyyy '${$t('common.calendar_week')}' w`"
                                                    disable-future
                                                />

                                                <UButton
                                                    square
                                                    icon="i-mdi-chevron-right"
                                                    class="flex-initial"
                                                    :disabled="isFuture(addWeeks(query.date.end, 1))"
                                                    @click="query.date.end = addWeeks(query.date.end, 1)"
                                                />
                                            </div>
                                            <DateRangePickerPopoverClient
                                                v-else
                                                v-model="query.date"
                                                mode="date"
                                                class="flex-1"
                                                :popover="{ class: 'flex-1' }"
                                                disable-future
                                            />
                                        </UFormGroup>

                                        <UFormGroup
                                            v-if="query.mode !== TimeclockMode.DAILY"
                                            name="day"
                                            :label="$t('common.per_day')"
                                        >
                                            <UToggle v-model="query.perDay" class="mt-2" />
                                        </UFormGroup>
                                    </div>
                                </div>
                            </div>
                        </UForm>
                    </template>
                </UDashboardToolbar>

                <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.entry', 2)])" :retry="refresh" />

                <template v-else>
                    <UTable
                        v-if="query.mode === TimeclockMode.WEEKLY"
                        v-model:sort="sort"
                        :loading="loading"
                        :columns="columns"
                        :rows="data?.entries.oneofKind === 'weekly' ? data.entries.weekly.entries : []"
                        :empty-state="{
                            icon: 'i-mdi-timeline-clock',
                            label: $t('common.not_found', [$t('common.entry', 2)]),
                        }"
                        sort-mode="manual"
                        class="flex-1"
                    >
                        <template #caption>
                            <caption>
                                <p class="text-right">
                                    <span class="font-semibold">{{ $t('common.sum') }}:</span>

                                    {{
                                        fromSecondsToFormattedDuration(
                                            parseFloat(
                                                (
                                                    (Math.round(
                                                        (data?.entries.oneofKind === 'weekly' ? data?.entries.weekly.sum : 0) *
                                                            100,
                                                    ) /
                                                        100) *
                                                    60 *
                                                    60
                                                ).toPrecision(2),
                                            ),
                                            { seconds: false },
                                        )
                                    }}
                                </p>
                            </caption>
                        </template>

                        <template #date-data="{ row: item }">
                            <div class="text-gray-900 dark:text-white">
                                {{ $d(toDate(item.date), 'date') }}
                            </div>
                        </template>

                        <template #name-data="{ row: item }">
                            <div class="inline-flex items-center gap-1">
                                <ProfilePictureImg
                                    :src="item.user?.avatar?.url"
                                    :name="`${item.user?.firstname} ${item.user?.lastname}`"
                                    size="xs"
                                />

                                <ColleagueInfoPopover :user="item.user" />
                            </div>
                        </template>

                        <template #rank-data="{ row: entry }">
                            {{ entry.user.jobGradeLabel
                            }}<span v-if="entry.user.jobGrade > 0"> ({{ entry.user.jobGrade }})</span>
                        </template>

                        <template #time-data="{ row: entry }">
                            <div class="text-right">
                                {{
                                    entry.spentTime > 0
                                        ? fromSecondsToFormattedDuration(
                                              parseFloat(((Math.round(entry.spentTime * 100) / 100) * 60 * 60).toPrecision(2)),
                                              { seconds: false },
                                          )
                                        : ''
                                }}

                                <UBadge v-if="entry.startTime !== undefined" color="green">
                                    {{ $t('common.active') }}
                                </UBadge>
                            </div>
                        </template>
                    </UTable>

                    <UTable
                        v-else
                        v-model:sort="sort"
                        :loading="loading"
                        :columns="columns"
                        :rows="
                            data?.entries.oneofKind === 'daily'
                                ? data.entries.daily.entries
                                : data?.entries.oneofKind === 'range'
                                  ? data.entries.range.entries
                                  : []
                        "
                        :empty-state="{
                            icon: 'i-mdi-timeline-clock',
                            label: $t('common.not_found', [$t('common.entry', 2)]),
                        }"
                        sort-mode="manual"
                        class="flex-1"
                    >
                        <template #caption>
                            <caption>
                                <p class="text-right">
                                    <span class="font-semibold">{{ $t('common.sum') }}:</span>

                                    {{
                                        fromSecondsToFormattedDuration(
                                            parseFloat(
                                                (
                                                    (Math.round(
                                                        (data?.entries.oneofKind === 'daily'
                                                            ? data?.entries.daily.sum
                                                            : data?.entries.oneofKind === 'range'
                                                              ? data?.entries.range.sum
                                                              : 0) * 100,
                                                    ) /
                                                        100) *
                                                    60 *
                                                    60
                                                ).toPrecision(2),
                                            ),
                                            { seconds: false },
                                        )
                                    }}
                                </p>
                            </caption>
                        </template>

                        <template #date-data="{ row: item }">
                            <div class="text-gray-900 dark:text-white">
                                {{ $d(toDate(item.date), 'date') }}
                            </div>
                        </template>

                        <template #name-data="{ row: entry }">
                            <div class="inline-flex items-center gap-1">
                                <ProfilePictureImg
                                    :src="entry.user?.avatar?.url"
                                    :name="`${entry.user?.firstname} ${entry.user?.lastname}`"
                                    size="xs"
                                />

                                <ColleagueInfoPopover :user="entry.user" />
                            </div>
                        </template>

                        <template #rank-data="{ row: entry }">
                            {{ entry.user.jobGradeLabel
                            }}<span v-if="entry.user.jobGrade > 0"> ({{ entry.user.jobGrade }})</span>
                        </template>

                        <template #time-data="{ row: entry }">
                            <div class="text-right">
                                {{
                                    entry.spentTime > 0
                                        ? fromSecondsToFormattedDuration(
                                              parseFloat(((Math.round(entry.spentTime * 100) / 100) * 60 * 60).toPrecision(2)),
                                              { seconds: false },
                                          )
                                        : ''
                                }}

                                <UBadge v-if="entry.startTime !== undefined" color="green">
                                    {{ $t('common.active') }}
                                </UBadge>
                            </div>
                        </template>
                    </UTable>

                    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
                </template>
            </template>
        </UTabs>
    </div>
</template>
