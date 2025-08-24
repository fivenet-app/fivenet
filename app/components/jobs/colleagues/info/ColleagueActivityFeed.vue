<script lang="ts" setup>
import { listEnumValues } from '@protobuf-ts/runtime';
import { z } from 'zod';
import ColleagueActivityFeedEntry from '~/components/jobs/colleagues/info/ColleagueActivityFeedEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { useCompletorStore } from '~/stores/completor';
import { getJobsJobsClient } from '~~/gen/ts/clients';
import { ColleagueActivityType } from '~~/gen/ts/resources/jobs/activity';
import type { ListColleagueActivityResponse } from '~~/gen/ts/services/jobs/jobs';
import ColleagueName from '../ColleagueName.vue';

const props = withDefaults(
    defineProps<{
        userId?: number;
        showTargetUser?: boolean;
    }>(),
    {
        userId: undefined,
        showTargetUser: false,
    },
);

const { attrStringList, isSuperuser } = useAuth();

const completorStore = useCompletorStore();

const jobsJobsClient = await getJobsJobsClient();

const usersLoading = ref(false);

const typesAttrs = (
    isSuperuser.value
        ? listEnumValues(ColleagueActivityType)
              .filter((t) => t.number !== 0)
              .map((t) => t.name)
        : attrStringList('jobs.JobsService/ListColleagueActivity', 'Types').value
).map((t) => t.toUpperCase());
const activityTypes = Object.keys(ColleagueActivityType)
    .filter((aType) => typesAttrs.includes(aType))
    .map((aType) => ColleagueActivityType[aType as keyof typeof ColleagueActivityType]);

const schema = z.object({
    colleagues: z.coerce
        .number()
        .array()
        .max(10)
        .default(props.userId ? [props.userId] : []),
    types: z.nativeEnum(ColleagueActivityType).array().max(typesAttrs.length).default(activityTypes),
    sorting: z
        .custom<SortByColumn>()
        .array()
        .max(3)
        .default([
            {
                id: 'createdAt',
                desc: true,
            },
        ]),
    page: pageNumberSchema,
});

type Schema = z.output<typeof schema>;

const query = useSearchForm('jobs_colleagues_activity' + (props.userId !== undefined ? '' : '_individual'), schema);

if (props.userId !== undefined) {
    query.colleagues = [props.userId];
}

const { data, status, refresh, error } = useLazyAsyncData(
    `jobs-colleague-${query.sorting.column}:${query.sorting.direction}-${query.page}-${query.colleagues.join(',')}-${query.types.join(':')}-${props.userId}`,
    () => listColleagueActivity(query),
);

async function listColleagueActivity(values: Schema): Promise<ListColleagueActivityResponse> {
    try {
        const call = jobsJobsClient.listColleagueActivity({
            pagination: {
                offset: calculateOffset(values.page, data.value?.pagination),
            },
            sort: values.sort,
            userIds: values.colleagues,
            activityTypes: values.types,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

const accessAttrs = attrStringList('jobs.JobsService/GetColleague', 'Access');
const colleagueSearchAttrs = ['Own', 'Lower_Rank', 'Same_Rank', 'Any'];

watch(props, async () => refresh());
</script>

<template>
    <UDashboardToolbar v-if="userId === undefined || accessAttrs.some((a) => colleagueSearchAttrs.includes(a)) || isSuperuser">
        <UForm class="flex w-full gap-2" :schema="schema" :state="query" @submit="refresh()">
            <UFormField v-if="userId === undefined" class="flex-1" name="colleagues" :label="$t('common.search')">
                <ClientOnly>
                    <USelectMenu
                        v-model="query.colleagues"
                        multiple
                        :searchable="
                            async (q: string) => {
                                usersLoading = true;
                                const colleagues = await completorStore.listColleagues({
                                    search: q,
                                    labelIds: [],
                                    userIds: query.colleagues,
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
                        leading-icon="i-mdi-search"
                        value-key="userId"
                    >
                        <template #item-label="{ selected }">
                            <template v-if="selected.length">
                                {{ usersToLabel(selected) }}
                            </template>
                        </template>

                        <template #option="{ option: colleague }">
                            <ColleagueName :colleague="colleague" birthday />
                        </template>

                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.colleague', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormField>
            <div v-else class="flex-1" />

            <UFormField
                v-if="isSuperuser || accessAttrs.some((a) => colleagueSearchAttrs.includes(a))"
                name="types"
                :label="$t('common.type', 2)"
            >
                <ClientOnly>
                    <USelectMenu
                        v-model="query.types"
                        class="w-48 min-w-40 flex-initial"
                        multiple
                        block
                        trailing
                        option-attribute="aType"
                        :items="activityTypes.map((aType) => ({ aType: aType }))"
                        value-key="aType"
                        :searchable-placeholder="$t('common.type', 2)"
                    >
                        <template #item-label>
                            {{ $t('common.selected', query.types.length) }}
                        </template>

                        <template #option="{ option }">
                            {{ $t(`enums.jobs.ColleagueActivityType.${ColleagueActivityType[option.aType]}`) }}
                        </template>

                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormField>

            <UFormField label="&nbsp;">
                <SortButton v-model="query.sorting" :fields="[{ label: $t('common.created_at'), value: 'createdAt' }]" />
            </UFormField>
        </UForm>
    </UDashboardToolbar>

    <div class="relative flex-1 overflow-x-auto">
        <DataErrorBlock
            v-if="error"
            class="w-full"
            :title="$t('common.not_found', [`${$t('common.colleague', 1)} ${$t('common.activity')}`])"
            :error="error"
            :retry="refresh"
        />
        <DataNoDataBlock
            v-else-if="data?.activity.length === 0"
            class="w-full"
            icon="i-mdi-pulse"
            :type="`${$t('common.colleague', 1)} ${$t('common.activity')}`"
        />

        <div v-else-if="isRequestPending(status) || data?.activity">
            <ul class="divide-y divide-gray-100 dark:divide-gray-800" role="list">
                <template v-if="isRequestPending(status)">
                    <li v-for="idx in 10" :key="idx" class="px-2 py-4">
                        <div class="flex space-x-3">
                            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                                <USkeleton class="size-full" />
                            </div>

                            <div class="flex-1 space-y-1">
                                <div class="flex items-center justify-between">
                                    <h3 class="text-sm font-medium">
                                        <USkeleton class="h-5 w-[350px]" />
                                    </h3>

                                    <p>
                                        <USkeleton class="h-5 w-[175px]" />
                                    </p>
                                </div>

                                <div class="flex items-center justify-between">
                                    <p class="flex flex-col gap-1 text-sm">
                                        <USkeleton class="h-8 w-[200px]" />
                                    </p>
                                    <p class="inline-flex items-center gap-1 text-sm">
                                        <USkeleton class="h-5 w-[175px]" />
                                    </p>
                                </div>
                            </div>
                        </div>
                    </li>
                </template>

                <template v-else>
                    <li
                        v-for="activity in data?.activity"
                        :key="activity.id"
                        class="border-white px-2 py-4 hover:border-primary-500/25 hover:bg-primary-100/50 dark:border-gray-900 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
                    >
                        <ColleagueActivityFeedEntry :activity="activity" :show-target-user="showTargetUser" />
                    </li>
                </template>
            </ul>
        </div>
    </div>

    <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
</template>
