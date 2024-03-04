<script lang="ts" setup>
import { Switch } from '@headlessui/vue';
import { watchDebounced } from '@vueuse/core';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import ColleaguesListEntry from '~/components/jobs/colleagues/ColleaguesListEntry.vue';
import { useJobsStore } from '~/store/jobs';
import type { Perms } from '~~/gen/ts/perms';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import GenericTable from '~/components/partials/elements/GenericTable.vue';

const query = ref<{
    name: string;
    absent?: boolean;
}>({
    name: '',
    absent: false,
});
const offset = ref(0n);

const jobsStore = useJobsStore();
const { data, pending, refresh, error } = useLazyAsyncData(`jobs-colleagues-${offset.value}-${query.value.name}`, () =>
    jobsStore.listColleagues({
        pagination: {
            offset: offset.value,
        },
        searchName: query.value.name,
        absent: query.value.absent,
    }),
);

const searchInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

watch(offset, async () => refresh());
watchDebounced(query.value, () => refresh(), { debounce: 600, maxWait: 1400 });

function updateAbsenceDates(value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void {
    const colleague = data.value?.colleagues.find((c) => c.userId === value.userId);
    if (colleague === undefined) {
        return;
    }

    if (colleague.props === undefined) {
        colleague.props = {
            userId: colleague.userId,
            absenceBegin: value.absenceBegin,
            absenceEnd: value.absenceEnd,
        };
    } else {
        colleague.props.absenceBegin = value.absenceBegin;
        colleague.props.absenceEnd = value.absenceEnd;
    }
}
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="mx-auto flex flex-row gap-4">
                            <div class="form-control flex-1">
                                <label for="searchName" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.colleague', 1) }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <input
                                        ref="searchInput"
                                        v-model="query.name"
                                        type="text"
                                        name="searchName"
                                        :placeholder="$t('common.name')"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div class="form-control flex-initial">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.absent') }}
                                </label>
                                <div class="relative mt-3 flex items-center">
                                    <Switch
                                        v-model="query.absent"
                                        :class="[
                                            query.absent ? 'bg-error-500' : 'bg-base-700',
                                            'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2',
                                        ]"
                                    >
                                        <span class="sr-only">
                                            {{ $t('common.absent') }}
                                        </span>
                                        <span
                                            aria-hidden="true"
                                            :class="[
                                                query.absent ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-neutral ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </Switch>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.colleague', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.colleague', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="data?.colleagues.length === 0"
                            :focus="focusSearch"
                            :message="$t('components.citizens.citizens_list.no_citizens')"
                        />
                        <div v-else>
                            <GenericTable>
                                <template #thead>
                                    <tr>
                                        <th
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-1"
                                        ></th>
                                        <th
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-1"
                                        >
                                            {{ $t('common.name') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-2 py-3.5 text-left text-sm font-semibold text-neutral hidden lg:table-cell"
                                        >
                                            {{ $t('common.rank', 1) }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.absence_date') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.phone_number') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.date_of_birth') }}
                                        </th>
                                        <th
                                            v-if="can(['JobsService.GetColleague', 'JobsService.SetJobsUserProps'] as Perms[])"
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 text-right text-sm font-semibold text-neutral sm:pr-0"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </template>
                                <template #tbody>
                                    <ColleaguesListEntry
                                        v-for="colleague in data?.colleagues"
                                        :key="colleague.userId"
                                        :colleague="colleague"
                                        @update:absence-dates="updateAbsenceDates($event)"
                                    />
                                </template>
                            </GenericTable>

                            <TablePagination
                                :pagination="data?.pagination"
                                :refresh="refresh"
                                @offset-change="offset = $event"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
