<script lang="ts" setup>
import { ref, onBeforeMount } from 'vue';
import { UserCircleIcon } from '@heroicons/vue/20/solid'
import { getCitizenStoreClient } from '../../grpc/grpc';
import { GetUserActivityRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { UserActivity } from '@arpanet/gen/resources/users/users_pb';
import { getDateRelativeString } from '../../utils/time';
import { USER_ACTIVITY_TYPE_Util } from '@arpanet/gen/resources/users/users.pb_enums';
import { RectangleGroupIcon } from '@heroicons/vue/24/outline';

const activities = ref<Array<UserActivity>>([]);
const defaultIcon = UserCircleIcon;

const props = defineProps({
    userId: {
        required: true,
        type: Number,
    },
});

function getUserActivity() {
    const req = new GetUserActivityRequest();
    req.setUserId(props.userId);

    getCitizenStoreClient().
        getUserActivity(req, null).then((resp) => {
            activities.value = resp.getActivityList();
        });
}

onBeforeMount(() => {
    getUserActivity();
});
</script>

<template>
    <div class="mt-3">
        <button v-if="activities.length == 0" type="button"
            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
            disabled>
            <RectangleGroupIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-gray-300">No User Activity found</span>
        </button>
        <ul role="list" class="divide-y divide-gray-200">
            <li v-for="activity in activities" :key="activity.getId()" class="py-4">
                <div class="flex space-x-3">
                    <div class="h-6 w-6 rounded-full flex items-center justify-center bg-white">
                        <component :is="defaultIcon" />
                    </div>
                    <div class="flex-1 space-y-1">
                        <div class="flex items-center justify-between">
                            <h3 class="text-sm font-medium text-neutral">{{ activity.getSourceUser()?.getFirstname() }} {{
                                activity.getSourceUser()?.getLastname() }}</h3>
                            <p class="text-sm text-gray-400">{{ getDateRelativeString(activity.getCreatedAt()) }}</p>
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
