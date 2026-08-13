<script lang="ts" setup>
import { listEnumValues } from '@protobuf-ts/runtime';
import { z } from 'zod';
import ActivityFeedEntry from '~/components/jobs/colleagues/info/ActivityFeedEntry.vue';
import UserGroupSelector from '~/components/jobs/UserGroupSelector.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import type { Form } from '@nuxt/ui';
import { getJobsColleaguesClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { ColleagueActivityType } from '~~/gen/ts/resources/jobs/colleagues/activity/activity';
import type { ListColleagueActivityResponse } from '~~/gen/ts/services/jobs/colleagues';
import { jobsUserActivityTypeBGColor, jobsUserActivityTypeIcon } from './helpers';
import { userSelectorSchema } from '~/utils/validation';

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

const jobsColleaguesClient = await getJobsColleaguesClient();

const typesAttrs = computed(() =>
    (isSuperuser.value
        ? listEnumValues(ColleagueActivityType)
              .filter((t) => t.number !== 0)
              .map((t) => t.name)
        : attrStringList('jobs.ColleaguesService/ListColleagueActivity', 'Types').value
    ).map((t) => t.toUpperCase()),
);
const activityTypes = computed(() =>
    Object.keys(ColleagueActivityType)
        .filter((at) => typesAttrs.value.includes(at))
        .map((at) => ColleagueActivityType[at as keyof typeof ColleagueActivityType]),
);

const schema = z.object({
    users: userSelectorSchema,
    types: z.enum(ColleagueActivityType).array().max(typesAttrs.value.length).default(activityTypes.value),
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

const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { validatedQuery, commitValidatedQuery } = useFormSearchValidation<typeof schema>(query, formRef);

const activityKey = computed(
    () =>
        `jobs-colleague-${JSON.stringify(validatedQuery.value.sorting)}-${validatedQuery.value.page}-${validatedQuery.value.types.join(':')}-${JSON.stringify(validatedQuery.value.users)}-${props.userId}`,
);

if (props.userId !== undefined) {
    query.users = { userIds: [props.userId] };
}

const { data, status, refresh, error } = useLazyAsyncData(activityKey, () => listColleagueActivity(validatedQuery.value));

async function listColleagueActivity(values: Schema): Promise<ListColleagueActivityResponse> {
    try {
        const call = jobsColleaguesClient.listColleagueActivity({
            pagination: {
                offset: calculateOffset(values.page, data.value?.pagination),
            },
            sort: values.sorting,
            users: values.users,
            activityTypes: values.types,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const accessAttrs = attrStringList('jobs.ColleaguesService/GetColleague', 'Access');
const colleagueSearchAttrs = ['Own', 'Lower_Rank', 'Same_Rank', 'Any'];

watch(
    () => props.userId,
    async (userId) => {
        query.users = userId !== undefined ? { userIds: [userId] } : { userIds: [] };

        await commitValidatedQuery();
    },
);
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardToolbar
                v-if="userId === undefined || accessAttrs.some((a) => colleagueSearchAttrs.includes(a)) || isSuperuser"
            >
                <UForm
                    ref="formRef"
                    class="my-2 flex w-full gap-2"
                    :schema="schema"
                    :state="query"
                    @submit="commitValidatedQuery"
                >
                    <UFormField v-if="userId === undefined" class="flex-1" name="users" :label="$t('common.search')">
                        <UserGroupSelector v-model="query.users" class="w-full" />
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
                                :items="
                                    activityTypes.map((aType) => ({
                                        aType: aType,
                                        icon: jobsUserActivityTypeIcon(aType),
                                        ui: {
                                            itemLeadingIcon: jobsUserActivityTypeBGColor(aType),
                                        },
                                    }))
                                "
                                value-key="aType"
                                :search-input="{ placeholder: $t('common.type', 2) }"
                            >
                                <template #default>
                                    {{ $t('common.selected', query.types.length) }}
                                </template>

                                <template #item-label="{ item }">
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
