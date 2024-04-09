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

watch(offset, () => refresh());

const accessAttrs = attrList('JobsService.GetColleague', 'Access');
const colleagueSearchAttrs = ['own', 'lower_rank', 'same_rank', 'any'];

watch(props, async () => refresh());
watchDebounced(query, async () => refresh(), {
    debounce: 500,
    maxWait: 1250,
});

function charsGetDisplayValue(chars: Colleague[]): string {
    const cs: string[] = [];
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname}`));

    return cs.join(', ');
}
</script>

<template>
    <div class="py-2 pb-4">
        <div class="px-1 sm:px-2">
            <div
                v-if="userId === undefined && accessAttrs.some((a) => colleagueSearchAttrs.includes(a))"
                class="mb-4 sm:flex sm:items-center"
            >
                <div class="sm:flex-auto">
                    <UForm :schema="schema" :state="query" @submit="refresh()">
                        <div class="flex flex-row gap-2">
                            <div class="flex-1">
                                <UFormGroup name="colleagues" :label="$t('common.colleague', 2)" class="flex-1">
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
                                        :placeholder="
                                            query.colleagues ? charsGetDisplayValue(query.colleagues) : $t('common.owner')
                                        "
                                        trailing
                                        by="userId"
                                    >
                                        <template #option="{ option: user }">
                                            {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                        </template>
                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>
                                        <template #empty> {{ $t('common.not_found', [$t('common.creator', 2)]) }} </template>
                                    </USelectMenu>
                                </UFormGroup>
                            </div>
                        </div>
                    </UForm>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 align-middle">
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
                        <div v-else>
                            <ul role="list" class="divide-y divide-gray-200">
                                <li v-for="activity in data?.activity" :key="activity.id" class="py-4">
                                    <ColleagueActivityFeedEntry :activity="activity" :show-target-user="showTargetUser" />
                                </li>
                            </ul>

                            <Pagination v-model="page" :pagination="data?.pagination" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
