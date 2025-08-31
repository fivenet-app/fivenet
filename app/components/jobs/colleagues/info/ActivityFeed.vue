<script lang="ts" setup>
import { listEnumValues } from '@protobuf-ts/runtime';
import { z } from 'zod';
import ActivityFeedEntry from '~/components/jobs/colleagues/info/ActivityFeedEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { useCompletorStore } from '~/stores/completor';
import { getJobsJobsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
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

const typesAttrs = computed(() =>
    (isSuperuser.value
        ? listEnumValues(ColleagueActivityType)
              .filter((t) => t.number !== 0)
              .map((t) => t.name)
        : attrStringList('jobs.JobsService/ListColleagueActivity', 'Types').value
    ).map((t) => t.toUpperCase()),
);
const activityTypes = computed(() =>
    Object.keys(ColleagueActivityType)
        .filter((at) => typesAttrs.value.includes(at))
        .map((at) => ColleagueActivityType[at as keyof typeof ColleagueActivityType]),
);

const schema = z.object({
    colleagues: z.coerce
        .number()
        .array()
        .max(10)
        .default(props.userId ? [props.userId] : []),
    types: z.nativeEnum(ColleagueActivityType).array().max(typesAttrs.value.length).default(activityTypes.value),
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'createdAt',
                        desc: true,
                    },
                ]),
        })
        .default({ columns: [{ id: 'createdAt', desc: true }] }),
    page: pageNumberSchema,
});

type Schema = z.output<typeof schema>;

const query = useSearchForm('jobs_colleagues_activity' + (props.userId !== undefined ? '' : '_individual'), schema);

if (props.userId !== undefined) {
    query.colleagues = [props.userId];
}

const { data, status, refresh, error } = useLazyAsyncData(
    () =>
        `jobs-colleague-${JSON.stringify(query.sorting)}-${query.page}-${query.colleagues.join(',')}-${query.types.join(':')}-${props.userId}`,
    () => listColleagueActivity(query),
);

async function listColleagueActivity(values: Schema): Promise<ListColleagueActivityResponse> {
    try {
        const call = jobsJobsClient.listColleagueActivity({
            pagination: {
                offset: calculateOffset(values.page, data.value?.pagination),
            },
            sort: values.sorting,
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
    <UDashboardPanel :ui="{ root: 'min-h-0', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardToolbar
                v-if="userId === undefined || accessAttrs.some((a) => colleagueSearchAttrs.includes(a)) || isSuperuser"
            >
                <UForm class="my-2 flex w-full gap-2" :schema="schema" :state="query" @submit="refresh()">
                    <UFormField v-if="userId === undefined" class="flex-1" name="colleagues" :label="$t('common.search')">
                        <SelectMenu
                            v-model="query.colleagues"
                            multiple
                            class="w-full"
                            :searchable="
                                async (q: string) =>
                                    await completorStore.listColleagues({
                                        search: q,
                                        labelIds: [],
                                        userIds: query.colleagues,
                                    })
                            "
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :filter-fields="['firstname', 'lastname']"
                            :placeholder="$t('common.colleague', 2)"
                            leading-icon="i-mdi-search"
                            value-key="userId"
                        >
                            <template #item-label="{ item }">
                                <template v-if="item">
                                    {{ userToLabel(item) }}
                                </template>
                            </template>

                            <template #item="{ item }">
                                <ColleagueName v-if="item" :colleague="item" birthday />
                            </template>

                            <template #empty> {{ $t('common.not_found', [$t('common.colleague', 2)]) }} </template>
                        </SelectMenu>
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
                                :items="activityTypes.map((aType) => ({ aType: aType }))"
                                value-key="aType"
                                :search-input="{ placeholder: $t('common.type', 2) }"
                            >
                                <template #item-label>
                                    {{ $t('common.selected', query.types.length) }}
                                </template>

                                <template #item="{ item }">
                                    {{ $t(`enums.jobs.ColleagueActivityType.${ColleagueActivityType[item.aType]}`) }}
                                </template>

                                <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormField>

                    <UFormField label="&nbsp;">
                        <SortButton
                            v-model="query.sorting"
                            :fields="[{ label: $t('common.created_at'), value: 'createdAt' }]"
                        />
                    </UFormField>
                </UForm>
            </UDashboardToolbar>
        </template>

        <template #body>
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
                    <ul class="divide-y divide-default" role="list">
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
                            <ActivityFeedEntry
                                v-for="activity in data?.activity"
                                :key="activity.id"
                                :activity="activity"
                                :show-target-user="showTargetUser"
                            />
                        </template>
                    </ul>
                </div>
            </div>
        </template>

        <template #footer>
            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
        </template>
    </UDashboardPanel>
</template>
