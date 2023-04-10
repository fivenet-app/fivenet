<script lang="ts" setup>
import { UserCircleIcon } from '@heroicons/vue/20/solid'
import { GetUserActivityRequest } from '@fivenet/gen/services/citizenstore/citizenstore_pb';
import { UserActivity } from '@fivenet/gen/resources/users/users_pb';
import { toDateRelativeString } from '~/utils/time';
import { USER_ACTIVITY_TYPE_Util } from '@fivenet/gen/resources/users/users.pb_enums';
import { RectangleGroupIcon } from '@heroicons/vue/24/outline';
import { RpcError } from 'grpc-web';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';

const { $grpc } = useNuxtApp();

const defaultIcon = UserCircleIcon;

const props = defineProps({
    userId: {
        required: true,
        type: Number,
    },
});

const { data: activities, pending, refresh, error } = await useLazyAsyncData(`citizeninfo-activity-${props.userId}`, () => getUserActivity());

async function getUserActivity(): Promise<Array<UserActivity>> {
    return new Promise(async (res, rej) => {
        const req = new GetUserActivityRequest();
        req.setUserId(props.userId);

        try {
            const resp = await $grpc.getCitizenStoreClient().
                getUserActivity(req, null);

            return res(resp.getActivityList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div class="mt-3">
        <DataPendingBlock v-if="pending" message="Loading user activity..." />
        <DataErrorBlock v-else-if="error" title="Unable to load user activity!" :retry="refresh" />
        <button v-else-if="activities && activities.length == 0" type="button"
            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
            disabled>
            <RectangleGroupIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-gray-300">
                No User Activity found.
            </span>
        </button>
        <ul v-else role="list" class="divide-y divide-gray-200">
            <li v-for="activity in activities" :key="activity.getId()" class="py-4">
                <div class="flex space-x-3">
                    <div class="h-6 w-6 rounded-full flex items-center justify-center bg-white">
                        <component :is="defaultIcon" />
                    </div>
                    <div class="flex-1 space-y-1">
                        <div class="flex items-center justify-between">
                            <h3 class="text-sm font-medium text-neutral">{{ activity.getSourceUser()?.getFirstname() }} {{
                                activity.getSourceUser()?.getLastname() }}</h3>
                            <p class="text-sm text-gray-400">{{ toDateRelativeString(activity.getCreatedAt()) }}</p>
                        </div>
                        <p class="text-sm text-gray-300">{{ activity.getKey() }} <span class="font-bold">{{
                            USER_ACTIVITY_TYPE_Util.toEnumKey(activity.getType()) }}</span>: {{
        activity.getOldvalue() }} ⇒ {{ activity.getNewvalue() }}</p>
                    </div>
                </div>
            </li>
        </ul>
    </div>
</template>
