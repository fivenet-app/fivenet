<script lang="ts" setup>
import { z } from 'zod';
import CitizenActivityFeedEntry from '~/components/citizens/info/CitizenActivityFeedEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import { UserActivityType } from '~~/gen/ts/resources/users/activity';
import type { ListUserActivityResponse } from '~~/gen/ts/services/citizens/citizens';

const props = defineProps<{
    userId: number;
}>();

const { attr, activeChar } = useAuth();

const citizensCitizensClient = await getCitizensCitizensClient();

const activityTypes = Object.keys(UserActivityType)
    .map((aType) => UserActivityType[aType as keyof typeof UserActivityType])
    .filter((aType) => {
        if (typeof aType === 'string') {
            return false;
        } else if (typeof aType === 'number' && aType < 3) {
            return false;
        }
        return true;
    });
const options = activityTypes.map((aType) => ({ aType: aType }));

const schema = z.object({
    types: z.nativeEnum(UserActivityType).array().max(activityTypes.length).default(activityTypes),
    sort: z.custom<TableSortable>().default({
        column: 'createdAt',
        direction: 'desc',
    }),
    page: pageNumberSchema,
});

const query = useSearchForm('citizen_activity', schema);

const { data, status, refresh, error } = useLazyAsyncData(
    `citizeninfo-activity-${query.sort.column}:${query.sort.direction}-${props.userId}-${query.page}`,
    () => listUserActivity(),
);

async function listUserActivity(): Promise<ListUserActivityResponse> {
    try {
        const call = citizensCitizensClient.listUserActivity({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sort,
            userId: props.userId,
            types: query.types,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watchDebounced(query, async () => refresh(), {
    debounce: 500,
    maxWait: 1250,
});
</script>

<template>
    <UAlert
        v-if="userId === activeChar?.userId && !attr('citizens.CitizensService/ListUserActivity', 'Fields', 'Own').value"
        variant="subtle"
        color="error"
        icon="i-mdi-denied"
        :title="$t('components.citizens.CitizenInfoActivityFeed.own.title')"
        :description="$t('components.citizens.CitizenInfoActivityFeed.own.message')"
    />

    <div v-else>
        <UDashboardToolbar>
            <template #default>
                <UForm class="flex w-full flex-row gap-2" :schema="schema" :state="query" @submit="refresh()">
                    <UFormField class="flex-1 grow" name="types" :label="$t('common.type', 2)">
                        <ClientOnly>
                            <USelectMenu
                                v-model="query.types"
                                class="min-w-40 flex-1"
                                multiple
                                block
                                trailing
                                option-attribute="aType"
                                :items="options"
                                value-key="aType"
                                :searchable-placeholder="$t('common.type', 2)"
                            >
                                <template #item-label>
                                    {{ $t('common.selected', query.types.length) }}
                                </template>

                                <template #option="{ option }">
                                    {{ $t(`enums.users.UserActivityType.${UserActivityType[option.aType]}`) }}
                                </template>

                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>

                                <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormField>

                    <UFormField label="&nbsp;">
                        <SortButton v-model="query.sort" :fields="[{ label: $t('common.created_at'), value: 'createdAt' }]" />
                    </UFormField>
                </UForm>
            </template>
        </UDashboardToolbar>

        <div class="relative mt-2 flex-1">
            <DataPendingBlock
                v-if="isRequestPending(status)"
                :message="$t('common.loading', [`${$t('common.citizen', 1)} ${$t('common.activity')}`])"
            />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.not_found', [`${$t('common.citizen', 1)} ${$t('common.activity')}`])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="!data || data?.activity.length === 0"
                :type="`${$t('common.citizen', 1)} ${$t('common.activity')}`"
                icon="i-mdi-pulse"
            />

            <div v-else>
                <ul class="divide-y divide-gray-100 dark:divide-gray-800" role="list">
                    <li
                        v-for="activity in data?.activity"
                        :key="activity.id"
                        class="hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-white py-2 dark:border-gray-900"
                    >
                        <CitizenActivityFeedEntry :activity="activity" />
                    </li>
                </ul>
            </div>
        </div>

        <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
    </div>
</template>
