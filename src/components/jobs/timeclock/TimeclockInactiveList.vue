<script lang="ts" setup>
import { min, numeric, required } from '@vee-validate/rules';
import { watchDebounced } from '@vueuse/core';
import { ArrowLeftIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TimeclockInactiveListEntry from '~/components/jobs/timeclock/TimeclockInactiveListEntry.vue';
import type { Perms } from '~~/gen/ts/perms';
import GenericTable from '~/components/partials/elements/GenericTable.vue';
import type { ListInactiveEmployeesResponse } from '~~/gen/ts/services/jobs/timeclock';

const { $grpc } = useNuxtApp();

const query = ref<{
    days: number;
}>({
    days: 14,
});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * page.value : 0));

const { data, pending, refresh, error } = useLazyAsyncData(`jobs-timeclock-inactive-${page.value}-${query.value.days}`, () =>
    listInactiveEmployees(),
);

async function listInactiveEmployees(): Promise<ListInactiveEmployeesResponse> {
    try {
        const call = $grpc.getJobsTimeclockClient().listInactiveEmployees({
            pagination: {
                offset: offset.value,
            },
            days: query.value.days,
        });

        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

interface FormData {
    days: number;
}

defineRule('required', required);
defineRule('min', min);
defineRule('numeric', numeric);

const { meta } = useForm<FormData>({
    validationSchema: {
        days: { required: true, min: 1, numeric: true },
    },
});

const searchInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

watchDebounced(
    query.value,
    async () => {
        if (meta.value.valid) {
            refresh();
        }
    },
    { debounce: 600, maxWait: 1400 },
);
watch(offset, async () => refresh());
</script>

<template>
    <div class="py-2 pb-4">
        <div v-if="can('JobsTimeclockService.ListTimeclock')" class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <NuxtLink
                    :to="{ name: 'jobs-timeclock' }"
                    class="inline-flex rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                >
                    <ArrowLeftIcon class="mr-1 size-5" />
                    {{ $t('common.timeclock') }}
                </NuxtLink>
            </div>
        </div>
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="mx-auto flex flex-row gap-4">
                            <div class="flex-1">
                                <label for="days" class="block text-sm font-medium leading-6">
                                    {{ $t('common.time_ago.day', 2) }}
                                </label>
                                <div class="relative mt-2">
                                    <VeeField
                                        ref="searchInput"
                                        v-model="query.days"
                                        name="days"
                                        type="number"
                                        min="3"
                                        max="31"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :label="$t('common.time_ago.day', 2)"
                                        :placeholder="$t('common.time_ago.day', 2)"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                    <VeeErrorMessage name="days" as="p" class="mt-2 text-sm text-error-400" />
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
                        <template v-else>
                            <GenericTable>
                                <template #thead>
                                    <tr>
                                        <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1"></th>
                                        <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="hidden px-2 py-3.5 text-left text-sm font-semibold lg:table-cell"
                                        >
                                            {{ $t('common.rank', 1) }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold">
                                            {{ $t('common.phone_number') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold">
                                            {{ $t('common.date_of_birth') }}
                                        </th>
                                        <th
                                            v-if="can(['JobsService.GetColleague'] as Perms[])"
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 text-right text-sm font-semibold sm:pr-0"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </template>
                                <template #tbody>
                                    <TimeclockInactiveListEntry
                                        v-for="colleague in data?.colleagues"
                                        :key="colleague.userId"
                                        :colleague="colleague"
                                    />
                                </template>
                            </GenericTable>
                        </template>

                        <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
                            <UPagination
                                v-model="page"
                                :page-count="data?.pagination?.pageSize ?? 0"
                                :total="data?.pagination?.totalCount ?? 0"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
