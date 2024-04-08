<script lang="ts" setup>
import { CalendarIcon, ChevronLeftIcon, ChevronRightIcon } from 'mdi-vue3';
import { format } from 'date-fns';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import { TimeclockEntry } from '~~/gen/ts/resources/jobs/timeclock';
import TimeclockListEntry from '~/components/jobs/timeclock/TimeclockListEntry.vue';
import TimeclockStatsBlock from '~/components/jobs/timeclock/TimeclockStatsBlock.vue';
import { getWeekNumber } from '~/utils/time';
import type { ListTimeclockRequest, ListTimeclockResponse } from '~~/gen/ts/services/jobs/timeclock';
import GenericTable from '~/components/partials/elements/GenericTable.vue';
import DatePicker from '~/components/partials/DatePicker.vue';
import { useCompletorStore } from '~/store/completor';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import Pagination from '~/components/partials/Pagination.vue';

const { $grpc } = useNuxtApp();

const completorStore = useCompletorStore();

const canAccessAll = attr('JobsTimeclockService.ListTimeclock', 'Access', 'All');

const now = new Date();
const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
const currentDay = ref(new Date(today.getFullYear(), today.getMonth(), today.getDate()));

const futureDay = ref(new Date(currentDay.value.getFullYear(), currentDay.value.getMonth(), currentDay.value.getDate() + 1));
const previousDay = ref(new Date(currentDay.value.getFullYear(), currentDay.value.getMonth(), currentDay.value.getDate() - 1));

const query = ref<{
    user_ids?: Colleague[];
    user?: Colleague;
    from?: Date;
    to?: Date;
    perDay: boolean;
}>({
    from: currentDay.value,
    to: canAccessAll ? previousDay.value : undefined,
    perDay: canAccessAll,
});

const usersLoading = ref(false);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending, refresh, error } = useLazyAsyncData(
    `jobs-timeclock-${query.value.from}-${query.value.to}-${query.value.perDay}-${query.value.user ?? query.value.user_ids?.map((u) => u.userId)}-${page.value}`,
    () => listTimeclockEntries(),
);

async function listTimeclockEntries(): Promise<ListTimeclockResponse> {
    try {
        const req: ListTimeclockRequest = {
            pagination: {
                offset: offset.value,
            },
            userIds: query.value.user ? [query.value.user.userId] : query.value.user_ids?.map((u) => u.userId) ?? [],
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

type GroupedTimeClockEntries = { date: Date; key: string; entries: TimeclockEntry[] }[];
const grouped = computed(() => {
    const groups: GroupedTimeClockEntries = [];
    data.value?.entries.forEach((e) => {
        const date = toDate(e.date);
        const idx = groups.findIndex((g) => g.key === date.toString());
        if (idx === -1) {
            groups.push({
                date,
                entries: [e],
                key: date.toString(),
            });
        } else {
            groups[idx].entries.push(e);
        }
    });

    return groups;
});

watch(offset, async () => refresh());
watchDebounced(
    query.value,
    async () => {
        if (canAccessAll) {
            if (query.value.user !== undefined || (query.value.user_ids !== undefined && query.value.user_ids.length > 0)) {
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

const input = ref<{ input: HTMLInputElement }>();

defineShortcuts({
    '/': () => {
        input.value?.input?.focus();
    },
});
</script>

<template>
    <div>
        <UForm :schema="{}" :state="{}" @submit="refresh()">
            <UDashboardToolbar>
                <template #default>
                    <div class="flex w-full flex-col">
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
                                <UFormGroup v-if="canAccessAll" class="flex-1" name="user" :label="$t('common.colleague', 2)">
                                    <UInputMenu
                                        v-model="query.user"
                                        ref="input"
                                        :search="
                                            async (query: string) => {
                                                usersLoading = true;
                                                const colleagues = await completorStore.listColleagues({
                                                    pagination: { offset: 0 },
                                                    searchName: query,
                                                });
                                                usersLoading = false;
                                                return colleagues;
                                            }
                                        "
                                        :search-attributes="['firstname', 'lastname']"
                                        block
                                        :placeholder="
                                            query.user
                                                ? `${query.user?.firstname} ${query.user?.lastname} (${query.user?.dateofbirth})`
                                                : $t('common.owner')
                                        "
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
                                    class="flex-1"
                                    name="from"
                                    :label="
                                        query.perDay ? $t('common.date') : `${$t('common.time_range')} ${$t('common.from')}`
                                    "
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
                                    class="flex-1"
                                    name="to"
                                    :label="`${$t('common.time_range')} ${$t('common.to')}`"
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

                        <div v-if="query.perDay" class="flex flex-row gap-4 pt-2">
                            <div class="flex-1">
                                <UButton block :disabled="futureDay > today" @click="dayForward()">
                                    <ChevronLeftIcon class="size-5" />
                                    {{ $t('common.forward') }} - {{ $d(futureDay, 'date') }}
                                </UButton>
                            </div>

                            <div class="flex-initial">
                                <UButton
                                    disabled
                                    class="disabled relative flex w-full cursor-pointer flex-col place-content-end items-center rounded-md bg-base-500 px-3 py-2 text-sm font-semibold hover:bg-base-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-500"
                                >
                                    <span class="inline-flex flex-row items-center gap-1">
                                        <CalendarIcon class="size-5" />
                                        {{ $d(currentDay, 'date') }}
                                    </span>
                                    <span>{{ $t('common.calendar_week') }}: {{ getWeekNumber(currentDay) }}</span>
                                </UButton>
                            </div>

                            <div class="flex-1">
                                <UButton block @click="dayBackwards()">
                                    {{ $d(previousDay, 'date') }} - {{ $t('common.previous') }}
                                    <ChevronRightIcon class="size-5" />
                                </UButton>
                            </div>
                        </div>
                    </div>
                </template>
            </UDashboardToolbar>

            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.timeclock', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.timeclock', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="data && data.entries && data.entries.length === 0"
                            :message="$t('components.citizens.CitizensList.no_citizens')"
                        />
                        <template v-else>
                            <GenericTable class="min-w-full divide-y divide-base-600">
                                <template #thead>
                                    <tr>
                                        <th
                                            v-if="!query.perDay"
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1"
                                        >
                                            {{ $t('common.date') }}
                                        </th>
                                        <th
                                            v-if="
                                                query.user_ids === undefined ||
                                                query.user_ids.length === 0 ||
                                                query.user_ids?.length >= 1
                                            "
                                            scope="col"
                                            class="px-2 py-3.5 text-left text-sm font-semibold"
                                        >
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold">
                                            {{ $t('common.time') }}
                                        </th>
                                    </tr>
                                </template>
                                <template #tbody>
                                    <template v-for="group in grouped" :key="group.key">
                                        <TimeclockListEntry
                                            v-for="(entry, idx) in group.entries"
                                            :key="entry.userId + toDate(entry.date).toString()"
                                            :entry="entry"
                                            :first="idx === 0 ? group.date : undefined"
                                            :show-date="!query.perDay"
                                        />
                                    </template>
                                </template>
                            </GenericTable>
                        </template>

                        <Pagination v-model="page" :pagination="data?.pagination" />
                    </div>
                </div>
            </div>

            <div v-if="data && data.stats" class="mb-4 flow-root">
                <div class="mt-2 sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <TimeclockStatsBlock
                            :stats="data.stats"
                            :weekly="data.weekly"
                            :hide-header="true"
                            :failed="error !== null"
                            @refresh="refresh()"
                        />
                    </div>
                </div>
            </div>
        </UForm>
    </div>
</template>
