<script lang="ts" setup>
import { addDays, isFuture, subDays } from 'date-fns';
import { z } from 'zod';
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import TimeclockStatsBlock from '~/components/jobs/timeclock/TimeclockStatsBlock.vue';
import DatePickerPopoverClient from '~/components/partials/DatePickerPopover.client.vue';
import Pagination from '~/components/partials/Pagination.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useCompletorStore } from '~/store/completor';
import { getWeekNumber } from '~/utils/time';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { TimeclockEntry } from '~~/gen/ts/resources/jobs/timeclock';
import type { ListTimeclockRequest, ListTimeclockResponse } from '~~/gen/ts/services/jobs/timeclock';

const props = withDefaults(
    defineProps<{
        userId?: number;
        showStats?: boolean;
    }>(),
    {
        showStats: true,
    },
);

const { t } = useI18n();

const completorStore = useCompletorStore();

const canAccessAll = attr('JobsTimeclockService.ListTimeclock', 'Access', 'All');

const schema = z.object({
    users: z.custom<Colleague>().array().max(5).optional(),
    from: z.date(),
    to: z.date().optional(),
    perDay: z.boolean(),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    users: [],
    from: new Date(),
    to: subDays(new Date(), 1),
    perDay: true,
});

function setFromProps(): void {
    if (!props.userId) {
        return;
    }

    query.users?.push({
        userId: props.userId,
        firstname: '',
        lastname: '',
        dateofbirth: '',
        job: '',
        jobGrade: 0,
    });
}

watch(props, setFromProps);

const futureDay = computed(() => addDays(query.from, 1));
const previousDay = computed(() => subDays(query.from, 1));

const usersLoading = ref(false);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `jobs-timeclock-${query.from}-${query.to}-${query.perDay}-${query.users?.map((u) => u.userId)}-${page.value}`,
    () => listTimeclockEntries(),
);

const perDayView = computed(() => !canAccessAll.value || !(query.users !== undefined && query.users.length > 0));

async function listTimeclockEntries(): Promise<ListTimeclockResponse> {
    try {
        const req: ListTimeclockRequest = {
            pagination: {
                offset: offset.value,
            },
            userIds: query.users?.map((u) => u.userId) ?? [],
        };

        req.from = {
            timestamp: googleProtobufTimestamp.Timestamp.fromDate(query.from),
        };

        req.perDay = perDayView.value;
        if (req.perDay) {
            req.to = {
                timestamp: googleProtobufTimestamp.Timestamp.fromDate(subDays(query.from, 1)),
            };
        } else if (query.to) {
            req.to = {
                timestamp: googleProtobufTimestamp.Timestamp.fromDate(query.to),
            };
        }

        const call = getGRPCJobsTimeclockClient().listTimeclock(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

type GroupTimeClockEntry = { date?: Date; entry: TimeclockEntry }[];

const grouped = computed(() => {
    const groups: GroupTimeClockEntry = [];

    data.value?.entries.forEach((e) => {
        const date = toDate(e.date);
        const idx = groups.findIndex((g) => g.date === date);
        if (idx === -1) {
            groups.push({
                date: date,
                entry: e,
            });
        } else {
            groups.push({
                date: date,
                entry: e,
            });
        }
    });

    return groups;
});

watch(offset, async () => refresh());
watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

function dayForward(): void {
    query.from = addDays(query.from, 1);
    query.to = addDays(query.to ?? new Date(), 1);
}

function dayBackwards(): void {
    query.from = subDays(query.from, 1);
    query.to = subDays(query.to ?? new Date(), 1);
}

// Update date to something reasonable when per day view is actived
watch(perDayView, () => {
    if (canAccessAll && !perDayView.value) {
        const from = query.from;
        const to = subDays(query.to ?? new Date(), 7);
        query.from = to;
        query.to = from;
    }
});

const columns = computed(() =>
    [
        !perDayView.value
            ? {
                  key: 'date',
                  label: t('common.date'),
              }
            : undefined,
        {
            key: 'name',
            label: t('common.name'),
        },
        {
            key: 'time',
            label: t('common.time'),
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const input = ref<{ input: HTMLInputElement }>();
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm :schema="schema" :state="query" class="flex-1" @submit="refresh()">
                <div class="flex w-full flex-col gap-2">
                    <div class="flex w-full flex-col">
                        <UButton
                            v-if="can('JobsTimeclockService.ListInactiveEmployees').value && userId === undefined"
                            :to="{ name: 'jobs-timeclock-inactive' }"
                            class="mb-2 place-self-end"
                            trailing-icon="i-mdi-arrow-right"
                        >
                            {{ $t('common.inactive_colleagues') }}
                        </UButton>

                        <div class="flex flex-row gap-2">
                            <UFormGroup
                                v-if="canAccessAll && userId === undefined"
                                name="users"
                                :label="$t('common.search')"
                                class="flex-1"
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
                                        block
                                        :placeholder="$t('common.colleague', 2)"
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
                                        <template #empty> {{ $t('common.not_found', [$t('common.creator', 2)]) }} </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormGroup>

                            <UFormGroup
                                name="from"
                                :label="perDayView ? $t('common.date') : `${$t('common.time_range')} ${$t('common.from')}`"
                                class="flex-1"
                            >
                                <DatePickerPopoverClient
                                    v-model="query.from"
                                    :popover="{ popper: { placement: 'bottom-start' } }"
                                />
                            </UFormGroup>

                            <UFormGroup
                                v-if="!perDayView"
                                name="to"
                                :label="`${$t('common.time_range')} ${$t('common.to')}`"
                                class="flex-1"
                            >
                                <DatePickerPopoverClient
                                    v-model="query.to"
                                    :popover="{ popper: { placement: 'bottom-start' } }"
                                />
                            </UFormGroup>
                        </div>
                    </div>

                    <div v-if="perDayView" class="flex flex-row gap-2">
                        <UButton
                            block
                            class="flex-1"
                            :disabled="isFuture(futureDay)"
                            icon="i-mdi-chevron-left"
                            @click="dayForward()"
                        >
                            {{ $t('common.forward') }} - {{ $d(futureDay, 'date') }}
                        </UButton>

                        <UButton
                            disabled
                            icon="i-mdi-calendar"
                            class="flex flex-initial cursor-pointer flex-col place-content-end items-center"
                        >
                            <span>
                                {{ $d(query.from, 'date') }}
                            </span>
                            <span>{{ $t('common.calendar_week') }}: {{ getWeekNumber(query.from) }}</span>
                        </UButton>

                        <UButton class="flex-1" block trailing-icon="i-mdi-chevron-right" @click="dayBackwards()">
                            {{ $d(previousDay, 'date') }} - {{ $t('common.previous') }}
                        </UButton>
                    </div>
                </div>
            </UForm>
        </template>
    </UDashboardToolbar>

    <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.entry', 2)])" :retry="refresh" />

    <UTable
        v-else
        :loading="loading"
        :columns="columns"
        :rows="grouped"
        :empty-state="{
            icon: 'i-mdi-timeline-clock',
            label: $t('common.not_found', [$t('common.entry', 2)]),
        }"
        class="flex-1"
    >
        <template #date-data="{ row: entry }">
            <div class="inline-flex items-center text-gray-900 dark:text-white">
                {{ $d(entry.date, 'date') }}
            </div>
        </template>

        <template #name-data="{ row: entry }">
            <div class="inline-flex items-center gap-1">
                <ProfilePictureImg
                    :src="entry.entry.user?.avatar?.url"
                    :name="`${entry.entry.user?.firstname} ${entry.entry.user?.lastname}`"
                    size="sm"
                />

                <ColleagueInfoPopover :user="entry.entry.user" />
            </div>
        </template>

        <template #time-data="{ row: entry }">
            <div class="text-right">
                {{
                    entry.entry.spentTime > 0
                        ? fromSecondsToFormattedDuration(
                              parseFloat(((Math.round(entry.entry.spentTime * 100) / 100) * 60 * 60).toPrecision(2)),
                              { seconds: false },
                          )
                        : ''
                }}

                <UBadge v-if="entry.entry.startTime !== undefined" color="green">
                    {{ $t('common.active') }}
                </UBadge>
            </div>
        </template>
    </UTable>

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />

    <UAccordion
        v-if="showStats && data && data.stats"
        :items="[{ slot: 'stats', label: $t('common.stats') }]"
        class="px-3 py-0.5"
    >
        <template #stats>
            <TimeclockStatsBlock
                :stats="data.stats"
                :weekly="data.weekly"
                :hide-header="true"
                :failed="!!error"
                :loading="loading"
            />
        </template>
    </UAccordion>
</template>
