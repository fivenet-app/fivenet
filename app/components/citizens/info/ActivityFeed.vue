<script lang="ts" setup>
import { z } from 'zod';
import ActivityFeedEntry from '~/components/citizens/info/ActivityFeedEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { UserActivityType } from '~~/gen/ts/resources/users/activity';
import type { ListUserActivityResponse } from '~~/gen/ts/services/citizens/citizens';

const props = defineProps<{
    userId: number;
}>();

const { t } = useI18n();

const { attr, activeChar } = useAuth();

const citizensCitizensClient = await getCitizensCitizensClient();

const activityTypes = Object.keys(UserActivityType)
    .map((t) => UserActivityType[t as keyof typeof UserActivityType])
    .filter((at) => {
        if (typeof at === 'string') {
            return false;
        } else if (typeof at === 'number' && at < 3) {
            return false;
        }
        return true;
    });
const options = activityTypes.map((at) => ({ label: t(`enums.users.UserActivityType.${UserActivityType[at]}`), value: at }));

const schema = z.object({
    types: z.enum(UserActivityType).array().max(activityTypes.length).default(activityTypes),
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

const query = useSearchForm('citizen_activity', schema);

const { data, status, refresh, error } = useLazyAsyncData(
    () => `citizeninfo-activity-${JSON.stringify(query.sorting)}-${props.userId}-${query.page}`,
    () => listUserActivity(),
);

async function listUserActivity(): Promise<ListUserActivityResponse> {
    try {
        const call = citizensCitizensClient.listUserActivity({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sorting,
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

const denyView = computed(() => {
    return (
        props.userId === activeChar.value?.userId && attr('citizens.CitizensService/ListUserActivity', 'Fields', 'Own').value
    );
});

watchDebounced(query, async () => refresh(), {
    debounce: 500,
    maxWait: 1250,
});
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template v-if="!denyView" #header>
            <UDashboardToolbar>
                <template #default>
                    <UForm class="my-2 flex w-full flex-row gap-2" :schema="schema" :state="query" @submit="refresh()">
                        <UFormField class="flex-1 grow" name="types" :label="$t('common.type', 2)">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.types"
                                    class="min-w-40 flex-1"
                                    multiple
                                    :items="options"
                                    value-key="value"
                                    :search-input="{ placeholder: $t('common.type', 2) }"
                                >
                                    <template #default>
                                        {{ $t('common.selected', query.types.length) }}
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
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <UContainer v-if="denyView" class="my-2">
                <UAlert
                    variant="subtle"
                    color="error"
                    icon="i-mdi-denied"
                    :title="$t('components.citizens.CitizenInfoActivityFeed.own.title')"
                    :description="$t('components.citizens.CitizenInfoActivityFeed.own.message')"
                />
            </UContainer>

            <DataPendingBlock
                v-else-if="isRequestPending(status)"
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

            <div v-else class="relative flex-1">
                <ul class="min-w-full divide-y divide-default overflow-clip" role="list">
                    <ActivityFeedEntry v-for="activity in data?.activity" :key="activity.id" :activity="activity" />
                </ul>
            </div>
        </template>

        <template #footer>
            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
        </template>
    </UDashboardPanel>
</template>
