<script lang="ts" setup>
import { format } from 'date-fns';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import { TimeclockEntry } from '~~/gen/ts/resources/jobs/timeclock';
import TimeclockStatsBlock from '~/components/jobs/timeclock/TimeclockStatsBlock.vue';
import { getWeekNumber } from '~/utils/time';
import type { ListTimeclockRequest, ListTimeclockResponse } from '~~/gen/ts/services/jobs/timeclock';
import DatePicker from '~/components/partials/DatePicker.vue';
import { useCompletorStore } from '~/store/completor';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import Pagination from '~/components/partials/Pagination.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import ColleagueInfoPopover from '../colleagues/ColleagueInfoPopover.vue';
import type { UserShort } from '~~/gen/ts/resources/users/users';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const completorStore = useCompletorStore();

const canAccessAll = attr('JobsTimeclockService.ListTimeclock', 'Access', 'All');

const now = new Date();
const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
const currentDay = ref(new Date(today.getFullYear(), today.getMonth(), today.getDate()));

const futureDay = ref(new Date(currentDay.value.getFullYear(), currentDay.value.getMonth(), currentDay.value.getDate() + 1));
const previousDay = ref(new Date(currentDay.value.getFullYear(), currentDay.value.getMonth(), currentDay.value.getDate() - 1));

const schema = z.object({
    userIds: z.custom<Colleague>().array().max(5).optional(),
    user: z.custom<Colleague>().optional(),
    from: z.date().optional(),
    to: z.date().optional(),
    perDay: z.boolean(),
});

type Schema = z.output<typeof schema>;

const query = ref<Schema>({
    from: currentDay.value,
    to: canAccessAll ? previousDay.value : undefined,
    perDay: canAccessAll,
});

const usersLoading = ref(false);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `jobs-timeclock-${query.value.from}-${query.value.to}-${query.value.perDay}-${query.value.user ?? query.value.userIds?.map((u) => u.userId)}-${page.value}`,
    () => listTimeclockEntries(),
);

async function listTimeclockEntries(): Promise<ListTimeclockResponse> {
    try {
        const req: ListTimeclockRequest = {
            pagination: {
                offset: offset.value,
            },
            userIds: query.value.user ? [query.value.user.userId] : query.value.userIds?.map((u) => u.userId) ?? [],
        };
        if (query.value.perDay !== undefined) {
            req.perDay = query.value.perDay;
        }
        if (query.value.from) {
            req.from = {
                timestamp: googleProtobufTimestamp.Timestamp.fromDate(query.value.from),
            };
        }
        if (query.value.to) {
            req.to = {
                timestamp: googleProtobufTimestamp.Timestamp.fromDate(query.value.to),
            };
        }

        const call = $grpc.getJobsTimeclockClient().listTimeclock(req);
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
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
watchDebounced(
    query.value,
    async () => {
        if (canAccessAll) {
            if (query.value.user !== undefined || (query.value.userIds !== undefined && query.value.userIds.length > 0)) {
                if (query.value.perDay) {
                    query.value.perDay = false;
                    query.value.to = undefined;
                }
            } else {
                query.value.perDay = true;
            }
        } else {
            query.value.perDay = false;
        }

        return refresh();
    },
    { debounce: 200, maxWait: 1250 },
);

function dayForward(): void {
    currentDay.value.setDate(currentDay.value.getDate() + 1);
    currentDay.value = new Date(currentDay.value);

    updateDates();
}

function dayBackwards(): void {
    currentDay.value.setDate(currentDay.value.getDate() - 1);
    currentDay.value = new Date(currentDay.value);

    updateDates();
}

function updateDates(): void {
    futureDay.value.setTime(
        new Date(currentDay.value.getFullYear(), currentDay.value.getMonth(), currentDay.value.getDate() + 1).getTime(),
    );
    futureDay.value = new Date(futureDay.value);
    previousDay.value.setTime(
        new Date(currentDay.value.getFullYear(), currentDay.value.getMonth(), currentDay.value.getDate() - 1).getTime(),
    );
    previousDay.value = new Date(previousDay.value);

    query.value.from = currentDay.value;
    query.value.to = previousDay.value;
}

function charsGetDisplayValue(chars: UserShort[]): string {
    const cs: string[] = [];
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname}`));

    return cs.join(', ');
}

const columns = computed(() =>
    [
        !query.value.perDay
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

defineShortcuts({
    '/': () => {
        input.value?.input?.focus();
    },
});
</script>

<template>
    <div>
        <UForm :schema="schema" :state="query" @submit="refresh()">
            <UDashboardToolbar>
                <template #default>
                    <div class="flex w-full flex-col gap-2">
                        <div class="flex w-full flex-col">
                            <UButton
                                v-if="can('JobsTimeclockService.ListInactiveEmployees')"
                                :to="{ name: 'jobs-timeclock-inactive' }"
                                class="place-self-end"
                                trailing-icon="i-mdi-arrow-right"
                            >
                                {{ $t('common.inactive_colleagues') }}
                            </UButton>

                            <div class="flex flex-row gap-2">
                                <UFormGroup v-if="canAccessAll" name="user" :label="$t('common.colleague', 2)" class="flex-1">
                                    <UInputMenu
                                        ref="input"
                                        v-model="query.userIds"
                                        nullable
                                        :search="
                                            async (query: string) => {
                                                usersLoading = true;
                                                const colleagues = await completorStore.listColleagues({
                                                    pagination: { offset: 0 },
                                                    search: query,
                                                });
                                                usersLoading = false;
                                                return colleagues;
                                            }
                                        "
                                        :search-attributes="['firstname', 'lastname']"
                                        block
                                        :placeholder="query.userIds ? charsGetDisplayValue(query.userIds) : $t('common.owner')"
                                        trailing
                                        by="userId"
                                    >
                                        <template #option="{ option: user }">
                                            {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                        </template>
                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>
                                        <template #empty> {{ $t('common.not_found', [$t('common.creator', 2)]) }} </template>

                                        <template #trailing>
                                            <UKbd value="/" />
                                        </template>
                                    </UInputMenu>
                                </UFormGroup>

                                <UFormGroup
                                    name="from"
                                    :label="
                                        query.perDay ? $t('common.date') : `${$t('common.time_range')} ${$t('common.from')}`
                                    "
                                    class="flex-1"
                                >
                                    <UPopover :popper="{ placement: 'bottom-start' }">
                                        <UButton
                                            variant="outline"
                                            color="gray"
                                            block
                                            icon="i-mdi-calendar-month"
                                            :label="query.from ? format(query.from, 'dd.MM.yyyy') : 'dd.mm.yyyy'"
                                        />

                                        <template #panel="{ close }">
                                            <DatePicker v-model="query.from" @close="close" />
                                        </template>
                                    </UPopover>
                                </UFormGroup>

                                <UFormGroup
                                    v-if="!query.perDay"
                                    name="to"
                                    :label="`${$t('common.time_range')} ${$t('common.to')}`"
                                    class="flex-1"
                                >
                                    <UPopover :popper="{ placement: 'bottom-start' }">
                                        <UButton
                                            variant="outline"
                                            color="gray"
                                            block
                                            icon="i-mdi-calendar-month"
                                            :label="query.to ? format(query.to, 'dd.MM.yyyy') : 'dd.mm.yyyy'"
                                        />

                                        <template #panel="{ close }">
                                            <DatePicker v-model="query.to" @close="close" />
                                        </template>
                                    </UPopover>
                                </UFormGroup>
                            </div>
                        </div>

                        <div v-if="query.perDay" class="flex flex-row gap-2">
                            <UButton
                                block
                                class="flex-1"
                                :disabled="futureDay > today"
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
                                    {{ $d(currentDay, 'date') }}
                                </span>
                                <span>{{ $t('common.calendar_week') }}: {{ getWeekNumber(currentDay) }}</span>
                            </UButton>

                            <UButton class="flex-1" block trailing-icon="i-mdi-chevron-right" @click="dayBackwards()">
                                {{ $d(previousDay, 'date') }} - {{ $t('common.previous') }}
                            </UButton>
                        </div>
                    </div>
                </template>
            </UDashboardToolbar>

            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
                        <DataErrorBlock
                            v-if="error"
                            :title="$t('common.unable_to_load', [$t('common.citizen', 2)])"
                            :retry="refresh"
                        />
                        <template v-else>
                            <UTable
                                :loading="loading"
                                :columns="columns"
                                :rows="grouped"
                                :empty-state="{
                                    icon: 'i-mdi-timeline-clock',
                                    label: $t('common.not_found', [$t('common.entry', 2)]),
                                }"
                            >
                                <template #date-data="{ row: entry }">
                                    <div class="inline-flex items-center">
                                        {{ $d(entry.date, 'date') }}
                                    </div>
                                </template>

                                <template #name-data="{ row: entry }">
                                    <div class="inline-flex items-center gap-1">
                                        <ProfilePictureImg
                                            :url="entry.entry.user?.avatar?.url"
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
                                                      parseFloat(
                                                          (
                                                              (Math.round(entry.entry.spentTime * 100) / 100) *
                                                              60 *
                                                              60
                                                          ).toPrecision(2),
                                                      ),
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

                            <Pagination v-model="page" :pagination="data?.pagination" />
                        </template>
                    </div>
                </div>
            </div>

            <TimeclockStatsBlock
                v-if="data && data.stats"
                :stats="data.stats"
                :weekly="data.weekly"
                :hide-header="true"
                :failed="error !== null"
                :loading="loading"
            />
        </UForm>
    </div>
</template>
