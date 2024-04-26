<script lang="ts" setup>
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import ColleagueActivityFeedEntry from '~/components/jobs/colleagues/info/ColleagueActivityFeedEntry.vue';
import type { ListColleagueActivityResponse } from '~~/gen/ts/services/jobs/jobs';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import { useCompletorStore } from '~/store/completor';
import Pagination from '~/components/partials/Pagination.vue';

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

const completorStore = useCompletorStore();

const usersLoading = ref(false);

const schema = z.object({
    colleagues: z.custom<Colleague>().array().max(10),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    colleagues: [],
});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const selectedUsersIds = computed(() => (props.userId !== undefined ? [props.userId] : query.colleagues.map((u) => u.userId)));

const { data, pending, refresh, error } = useLazyAsyncData(
    `jobs-colleague-${selectedUsersIds.value.join(',')}-${page.value}`,
    () => listColleagueActivity(selectedUsersIds.value),
);

async function listColleagueActivity(userIds: number[]): Promise<ListColleagueActivityResponse> {
    try {
        const call = $grpc.getJobsClient().listColleagueActivity({
            userIds,
            pagination: { offset: offset.value },
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
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
    <UDashboardToolbar v-if="userId === undefined && accessAttrs.some((a) => colleagueSearchAttrs.includes(a))">
        <UForm :schema="schema" :state="query" class="w-full" @submit="refresh()">
            <UFormGroup name="colleagues" :label="$t('common.search')" class="flex-1">
                <USelectMenu
                    v-model="query.colleagues"
                    multiple
                    :searchable="
                        async (query: string) => {
                            usersLoading = true;
                            const colleagues = await completorStore.listColleagues({
                                search: query,
                            });
                            usersLoading = false;
                            return colleagues;
                        }
                    "
                    :search-attributes="['firstname', 'lastname']"
                    block
                    :placeholder="$t('common.colleague', 2)"
                    trailing
                    by="userId"
                    :searchable-placeholder="$t('common.search_field')"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                >
                    <template #label>
                        <template v-if="query.colleagues.length">
                            {{ usersToLabel(query.colleagues) }}
                        </template>
                    </template>
                    <template #option="{ option: user }">
                        {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                    </template>
                    <template #option-empty="{ query: search }">
                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                    </template>
                    <template #empty> {{ $t('common.not_found', [$t('common.creator', 2)]) }} </template>
                </USelectMenu>
            </UFormGroup>
        </UForm>
    </UDashboardToolbar>

    <DataPendingBlock
        v-if="pending"
        :message="$t('common.loading', [`${$t('common.colleague', 1)} ${$t('common.activity')}`])"
    />
    <DataErrorBlock
        v-else-if="error"
        :title="$t('common.not_found', [`${$t('common.colleague', 1)} ${$t('common.activity')}`])"
        :retry="refresh"
    />
    <DataNoDataBlock
        v-else-if="data?.activity.length === 0"
        icon="i-mdi-bulletin-board"
        :type="`${$t('common.colleague', 1)} ${$t('common.activity')}`"
    />
    <ul v-else role="list" class="divide-y divide-gray-100 dark:divide-gray-800">
        <li v-for="activity in data?.activity" :key="activity.id" class="px-2 py-4">
            <ColleagueActivityFeedEntry :activity="activity" :show-target-user="showTargetUser" />
        </li>
    </ul>

    <Pagination v-model="page" :pagination="data?.pagination" />
</template>
