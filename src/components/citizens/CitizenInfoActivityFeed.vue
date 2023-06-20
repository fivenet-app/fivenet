<script lang="ts" setup>
import { mdiBulletinBoard } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import CitizenInfoActivityFeedEntry from '~/components/citizens/CitizenInfoActivityFeedEntry.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { UserActivity } from '~~/gen/ts/resources/users/users';
import DataNoDataBlock from '../partials/DataNoDataBlock.vue';
import TablePagination from '../partials/TablePagination.vue';

const { $grpc } = useNuxtApp();

const props = defineProps<{
    userId: number;
}>();

const pagination = ref<PaginationResponse>();
const offset = ref(0n);

const {
    data: activities,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`citizeninfo-activity-${props.userId}-${offset.value}`, () => listUserActivity());

async function listUserActivity(): Promise<Array<UserActivity>> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCitizenStoreClient().listUserActivity({
                pagination: {
                    offset: offset.value,
                },
                userId: props.userId,
            });
            const { response } = await call;

            pagination.value = response.pagination;
            return res(response.activity);
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
            v-else-if="activities && activities.length === 0"
            :icon="mdiBulletinBoard"
            :type="`${$t('common.citizen', 1)} ${$t('common.activity')}`"
        />
        <div v-else>
            <ul role="list" class="divide-y divide-gray-200">
                <li v-for="activity in activities" :key="activity.id?.toString()" class="py-4">
                    <CitizenInfoActivityFeedEntry :activity="activity" />
                </li>
            </ul>

            <TablePagination :pagination="pagination" @offset-change="offset = $event" />
        </div>
    </div>
</template>
