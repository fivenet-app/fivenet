<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { BulletinBoardIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { ListUserActivityResponse } from '~~/gen/ts/services/citizenstore/citizenstore';
import ActivityFeedEntry from './ActivityFeedEntry.vue';

const { $grpc } = useNuxtApp();

const props = defineProps<{
    userId: number;
}>();

const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`citizeninfo-activity-${props.userId}-${offset.value}`, () =>
    listUserActivity(),
);

async function listUserActivity(): Promise<ListUserActivityResponse> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCitizenStoreClient().listUserActivity({
                pagination: {
                    offset: offset.value,
                },
                userId: props.userId,
            });
            const { response } = await call;

            return res(response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watch(offset, async () => refresh());
</script>

<template>
    <div class="mt-2">
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
                <li v-for="activity in data?.activity" :key="activity.id?.toString()" class="py-4">
                    <ActivityFeedEntry :activity="activity" />
                </li>
            </ul>

            <TablePagination :pagination="data?.pagination" @offset-change="offset = $event" />
        </div>
    </div>
</template>
