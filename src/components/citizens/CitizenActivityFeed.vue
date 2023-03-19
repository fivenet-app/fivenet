<script setup lang="ts">
import { ref, defineProps, onBeforeMount } from 'vue';
import { BoltIcon, ChatBubbleLeftEllipsisIcon, TagIcon, UserCircleIcon } from '@heroicons/vue/20/solid'
import { getCitizenStoreClient, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';
import { GetUserActivityRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { UserActivity } from '@arpanet/gen/resources/users/users_pb';
import { useRoute } from 'vue-router/auto';

const route = useRoute();

const activities = ref<Array<UserActivity>>([]);
const defaultIcon = UserCircleIcon;

const $props = defineProps({
    userId: {
        required: true,
        type: Number,
    },
});

function getUserActivity() {
    const req = new GetUserActivityRequest();
    req.setUserId($props.userId);

    getCitizenStoreClient().
        getUserActivity(req, null).then((resp) => {
            activities.value = resp.getActivityList();
        }).catch((err: RpcError) => {
            handleGRPCError(err, route);
        });
}

onBeforeMount(() => {
    getUserActivity();
});
</script>

<template>
    <div>
        <span v-if="activities.length === 0">
            <p class="text-sm font-medium text-white">
                No Citizen Activities found.
            </p>
        </span>
        <ul role="list" class="divide-y divide-gray-200">
            <li v-for="activity in activities" :key="activity.getId()" class="py-4">
                <div class="flex space-x-3">
                    <div class="h-6 w-6 rounded-full flex items-center justify-center bg-white">
                        <component :is="defaultIcon" />
                    </div>
                    <div class="flex-1 space-y-1">
                        <div class="flex items-center justify-between">
                            <h3 class="text-sm font-medium text-white">{{ activity.getSourceUser()?.getFirstname() }} {{
                                activity.getSourceUser()?.getLastname() }}</h3>
                            <p class="text-sm text-gray-400">{{ activity.getCreatedAt() }}</p>
                        </div>
                        <p class="text-sm text-gray-300">{{ activity.getType() }} {{ activity.getKey() }}: {{
                            activity.getOldvalue() }} ⇒ {{ activity.getNewvalue() }}</p>
                    </div>
                </div>
            </li>
        </ul>
    </div>
</template>
