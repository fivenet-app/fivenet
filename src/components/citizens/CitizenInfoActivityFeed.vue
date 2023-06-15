<script lang="ts" setup>
import { mdiBulletinBoard } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import CitizenInfoActivityFeedEntry from '~/components/citizens/CitizenInfoActivityFeedEntry.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import { UserActivity } from '~~/gen/ts/resources/users/users';
import DataNoDataBlock from '../partials/DataNoDataBlock.vue';

const { $grpc } = useNuxtApp();

const props = defineProps<{
    userId: number;
}>();

const {
    data: activities,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`citizeninfo-activity-${props.userId}`, () => listUserActivity());

async function listUserActivity(): Promise<Array<UserActivity>> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCitizenStoreClient().listUserActivity({
                userId: props.userId,
            });
            const { response } = await call;

            return res(response.activity);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
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
        <ul v-else role="list" class="divide-y divide-gray-200">
            <li v-for="activity in activities" :key="activity.id?.toString()" class="py-4">
                <CitizenInfoActivityFeedEntry :activity="activity" />
            </li>
        </ul>
    </div>
</template>
