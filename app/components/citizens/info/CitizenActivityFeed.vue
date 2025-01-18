<script lang="ts" setup>
import { z } from 'zod';
import CitizenActivityFeedEntry from '~/components/citizens/info/CitizenActivityFeedEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { UserActivityType } from '~~/gen/ts/resources/users/activity';
import type { ListUserActivityResponse } from '~~/gen/ts/services/citizenstore/citizenstore';

const props = defineProps<{
    userId: number;
}>();

const { attr, activeChar } = useAuth();

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
    types: z.nativeEnum(UserActivityType).array().max(activityTypes.length),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    types: activityTypes,
});

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'createdAt',
    direction: 'desc',
});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `citizeninfo-activity-${sort.value.column}:${sort.value.direction}-${props.userId}-${page.value}`,
    () => listUserActivity(),
    {
        watch: [sort],
    },
);

async function listUserActivity(): Promise<ListUserActivityResponse> {
    try {
        const call = getGRPCCitizenStoreClient().listUserActivity({
            pagination: {
                offset: offset.value,
            },
            sort: sort.value,
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

watch(offset, async () => refresh());
watchDebounced(query, async () => refresh(), {
    debounce: 500,
    maxWait: 1250,
});
</script>

<template>
    <UAlert
        v-if="userId === activeChar?.userId && !attr('CitizenStoreService.ListUserActivity', 'Fields', 'Own').value"
        variant="subtle"
        color="red"
        icon="i-mdi-denied"
        :title="$t('components.citizens.CitizenInfoActivityFeed.own.title')"
        :description="$t('components.citizens.CitizenInfoActivityFeed.own.message')"
    />

    <div v-else>
        <UForm :schema="schema" :state="query" class="flex w-full gap-2" @submit="refresh()">
            <UFormGroup name="types" :label="$t('common.type', 2)" class="flex-1 grow">
                <ClientOnly>
                    <USelectMenu
                        v-model="query.types"
                        class="min-w-40 flex-1"
                        multiple
                        block
                        trailing
                        option-attribute="aType"
                        :options="options"
                        value-attribute="aType"
                        :searchable-placeholder="$t('common.type', 2)"
                    >
                        <template #label>
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
            </UFormGroup>

            <UFormGroup label="&nbsp;">
                <SortButton v-model="sort" :fields="[{ label: $t('common.created_at'), value: 'createdAt' }]" />
            </UFormGroup>
        </UForm>

        <div class="relative mt-2 flex-1">
            <DataPendingBlock
                v-if="loading"
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
                <ul role="list" class="divide-y divide-gray-100 dark:divide-gray-800">
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

        <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
    </div>
</template>
