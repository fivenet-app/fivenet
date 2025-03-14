<script lang="ts" setup>
import { z } from 'zod';
import ColleagueActivityFeedEntry from '~/components/jobs/colleagues/info/ColleagueActivityFeedEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { useCompletorStore } from '~/stores/completor';
import { JobsUserActivityType } from '~~/gen/ts/resources/jobs/activity';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
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

const { $grpc } = useNuxtApp();

const { attrList } = useAuth();

const completorStore = useCompletorStore();

const usersLoading = ref(false);

const typesAttrs = attrList('JobsService.ListColleagueActivity', 'Types').value.map((t) => t.toUpperCase());
const activityTypes = Object.keys(JobsUserActivityType)
    .filter((aType) => typesAttrs.includes(aType))
    .map((aType) => JobsUserActivityType[aType as keyof typeof JobsUserActivityType]);

const schema = z.object({
    colleagues: z.custom<Colleague>().array().max(10),
    types: z.nativeEnum(JobsUserActivityType).array().max(typesAttrs.length),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    colleagues: [],
    types: Object.keys(JobsUserActivityType)
        .filter((aType) => typesAttrs.includes(aType))
        .map((aType) => JobsUserActivityType[aType as keyof typeof JobsUserActivityType]),
});

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'createdAt',
    direction: 'desc',
});

const selectedUsersIds = computed(() => (props.userId !== undefined ? [props.userId] : query.colleagues.map((u) => u.userId)));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `jobs-colleague-${sort.value.column}:${sort.value.direction}-${page.value}-${selectedUsersIds.value.join(',')}-${query.types.join(':')}`,
    () => listColleagueActivity(selectedUsersIds.value, query.types),
    {
        watch: [sort],
    },
);

async function listColleagueActivity(
    userIds: number[],
    activityTypes: JobsUserActivityType[],
): Promise<ListColleagueActivityResponse> {
    try {
        const call = $grpc.jobs.jobs.listColleagueActivity({
            pagination: {
                offset: offset.value,
            },
            sort: sort.value,
            userIds: userIds,
            activityTypes: activityTypes,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());

const accessAttrs = attrList('JobsService.GetColleague', 'Access');
const colleagueSearchAttrs = ['own', 'lower_rank', 'same_rank', 'any'];

watch(props, async () => refresh());
watchDebounced(query, async () => refresh(), {
    debounce: 500,
    maxWait: 1250,
});
</script>

<template>
    <UDashboardToolbar v-if="userId === undefined || accessAttrs.some((a) => colleagueSearchAttrs.includes(a))">
        <UForm :schema="schema" :state="query" class="flex w-full gap-2" @submit="refresh()">
            <UFormGroup v-if="userId === undefined" name="colleagues" :label="$t('common.search')" class="flex-1">
                <ClientOnly>
                    <USelectMenu
                        v-model="query.colleagues"
                        multiple
                        :searchable="
                            async (query: string) => {
                                usersLoading = true;
                                const colleagues = await completorStore.listColleagues({
                                    search: query,
                                    labelIds: [],
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
                    >
                        <template #label>
                            <template v-if="query.colleagues.length">
                                {{ usersToLabel(query.colleagues) }}
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
            </UFormGroup>
            <div v-else class="flex-1" />

            <UFormGroup
                v-if="accessAttrs.some((a) => colleagueSearchAttrs.includes(a))"
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
                        :options="activityTypes.map((aType) => ({ aType: aType }))"
                        value-attribute="aType"
                        :searchable-placeholder="$t('common.type', 2)"
                    >
                        <template #label>
                            {{ $t('common.selected', query.types.length) }}
                        </template>

                        <template #option="{ option }">
                            {{ $t(`enums.jobs.JobsUserActivityType.${JobsUserActivityType[option.aType]}`) }}
                        </template>

                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormGroup>

            <UFormGroup label="&nbsp;">
                <SortButton v-model="sort" :fields="[{ label: $t('common.created_at'), value: 'createdAt' }]" />
            </UFormGroup>
        </UForm>
    </UDashboardToolbar>

    <div class="relative flex-1 overflow-x-auto">
        <DataErrorBlock
            v-if="error"
            :title="$t('common.not_found', [`${$t('common.colleague', 1)} ${$t('common.activity')}`])"
            :error="error"
            :retry="refresh"
            class="w-full"
        />
        <DataNoDataBlock
            v-else-if="data?.activity.length === 0"
            icon="i-mdi-pulse"
            :type="`${$t('common.colleague', 1)} ${$t('common.activity')}`"
            class="w-full"
        />

        <div v-else-if="loading || data?.activity">
            <ul role="list" class="divide-y divide-gray-100 dark:divide-gray-800">
                <template v-if="loading">
                    <li v-for="idx in 10" :key="idx" class="px-2 py-4">
                        <div class="flex space-x-3">
                            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                                <USkeleton class="size-full" :ui="{ rounded: 'rounded-full' }" />
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
                        class="hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-white px-2 py-4 dark:border-gray-900"
                    >
                        <ColleagueActivityFeedEntry :activity="activity" :show-target-user="showTargetUser" />
                    </li>
                </template>
            </ul>
        </div>
    </div>

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
