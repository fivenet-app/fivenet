<script lang="ts" setup>
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

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const selectedUsers = ref<Colleague[]>([]);
const selectedUsersIds = computed(() =>
    props.userId !== undefined ? [props.userId] : selectedUsers.value.map((u) => u.userId),
);

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

const queryColleagueNameRaw = ref<string>('');
const queryColleagueName = computed(() => queryColleagueNameRaw.value.trim());

const { data: colleagues, refresh: refreshColleagues } = useLazyAsyncData(
    `jobs-colleagues-${offset.value}-${queryColleagueName.value}`,
    async () => {
        try {
            const call = $grpc.getJobsClient().listColleagues({
                pagination: {
                    offset: offset.value,
                },
                searchName: queryColleagueName.value,
            });
            const { response } = await call;

            return response;
        } catch (e) {
            $grpc.handleError(e as RpcError);
            throw e;
        }
    },
);

watchDebounced(
    queryColleagueName,
    async () => {
        await refreshColleagues();
        if (props.userId === undefined && selectedUsers.value) {
            colleagues.value?.colleagues.unshift(...selectedUsers.value);
        }
    },
    {
        debounce: 500,
        maxWait: 1250,
    },
);

const accessAttrs = attrList('JobsService.GetColleague', 'Access');
const colleagueSearchAttrs = ['own', 'lower_rank', 'same_rank', 'any'];

watch(props, async () => refresh());
watchDebounced(selectedUsers, async () => refresh(), {
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
        <div class="px-1 sm:px-2 lg:px-4">
            <div
                v-if="userId === undefined && accessAttrs.some((a) => colleagueSearchAttrs.includes(a))"
                class="mb-4 sm:flex sm:items-center"
            >
                <div class="sm:flex-auto">
                    <UForm :schema="{}" :state="{}" @submit="refresh()">
                        <div class="flex flex-row gap-2">
                            <div class="flex-1">
                                <UFormGroup name="selectedUsers" :label="$t('common.colleague', 2)" class="flex-1">
                                    <USelectMenu
                                        v-model="selectedUsers"
                                        multiple
                                        :searchable="
                                            async (query: string) => {
                                                usersLoading = true;
                                                const colleagues = await completorStore.completeCitizens({
                                                    search: query,
                                                });
                                                usersLoading = false;
                                                return colleagues;
                                            }
                                        "
                                        :search-attributes="['firstname', 'lastname']"
                                        block
                                        :placeholder="selectedUsers ? charsGetDisplayValue(selectedUsers) : $t('common.owner')"
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
