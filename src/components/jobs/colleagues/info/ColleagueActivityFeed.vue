<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { BulletinBoardIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import ColleagueActivityFeedEntry from '~/components/jobs/colleagues/info/ColleagueActivityFeedEntry.vue';
import type { ListColleagueActivityResponse } from '~~/gen/ts/services/jobs/jobs';

const props = defineProps<{
    userId: number;
}>();

const { $grpc } = useNuxtApp();

const offset = ref(0n);
const { data, pending, refresh, error } = useLazyAsyncData(`jobs-colleague-${props.userId}-${offset.value}`, () =>
    listColleagueActivity(props.userId),
);

async function listColleagueActivity(userId: number): Promise<ListColleagueActivityResponse> {
    try {
        const call = $grpc.getJobsClient().listColleagueActivity({
            userId,
            pagination: { offset: offset.value },
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(props, async () => refresh());
</script>

<template>
    <div>
        <div class="py-2 pb-14">
            <div class="px-1 sm:px-2 lg:px-4">
                <div class="flow-root">
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
                                :icon="markRaw(BulletinBoardIcon)"
                                :type="`${$t('common.colleague', 1)} ${$t('common.activity')}`"
                            />
                            <div v-else>
                                <ul role="list" class="divide-y divide-gray-200">
                                    <li v-for="activity in data?.activity" :key="activity.id" class="py-4">
                                        <ColleagueActivityFeedEntry :activity="activity" />
                                    </li>
                                </ul>

                                <TablePagination
                                    :pagination="data?.pagination"
                                    :refresh="refresh"
                                    @offset-change="offset = $event"
                                />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
