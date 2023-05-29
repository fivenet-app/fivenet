<script lang="ts" setup>
import { RectangleGroupIcon } from '@heroicons/vue/24/outline';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import CitizenInfoActivityFeedEntry from '~/components/citizens/CitizenInfoActivityFeedEntry.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import { UserActivity } from '~~/gen/ts/resources/users/users';

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
        <button
            v-else-if="activities && activities.length === 0"
            type="button"
            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
            disabled
        >
            <RectangleGroupIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-gray-300">
                {{ $t('common.not_found', [`${$t('common.user', 1)} ${$t('common.activity')}`]) }}
            </span>
        </button>
        <ul v-else role="list" class="divide-y divide-gray-200">
            <li v-for="activity in activities" :key="activity.id" class="py-4">
                <CitizenInfoActivityFeedEntry :activity="activity" />
            </li>
        </ul>
    </div>
</template>
