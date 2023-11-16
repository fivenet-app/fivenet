<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { BulletinBoardIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { ListUserActivityResponse } from '~~/gen/ts/services/citizenstore/citizenstore';
import ActivityFeedEntry from '~/components/citizens/info/ActivityFeedEntry.vue';

const { $grpc } = useNuxtApp();

const props = defineProps<{
    userId: number;
}>();

const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`citizeninfo-activity-${props.userId}-${offset.value}`, () =>
    listUserActivity(),
);

async function listUserActivity(): Promise<ListUserActivityResponse> {
    try {
        const call = $grpc.getCitizenStoreClient().listUserActivity({
            pagination: {
                offset: offset.value,
            },
            userId: props.userId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="flow-root">
                <div class="mx-0 -my-2 overflow-x-auto">
                    <div class="inline-block min-w-full align-middle px-1">
                        <DataPendingBlock
                            v-if="pending"
                            :message="$t('common.loading', [`${$t('common.user', 1)} ${$t('common.activity')}`])"
                        />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.not_found', [`${$t('common.user', 1)} ${$t('common.activity')}`])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="data?.activity.length === 0"
                            :icon="BulletinBoardIcon"
                            :type="`${$t('common.citizen', 1)} ${$t('common.activity')}`"
                        />
                        <div v-else>
                            <ul role="list" class="divide-y divide-gray-200">
                                <li v-for="activity in data?.activity" :key="activity.id" class="py-4">
                                    <ActivityFeedEntry :activity="activity" />
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
</template>
