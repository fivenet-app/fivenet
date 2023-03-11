<script lang="ts">
import { defineComponent } from 'vue';
import { BoltIcon, ChatBubbleLeftEllipsisIcon, TagIcon, UserCircleIcon } from '@heroicons/vue/20/solid'
import { getUsersClient, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';
import { GetUserActivityRequest } from '@arpanet/gen/services/users/users_pb';
import { UserActivity } from '@arpanet/gen/resources/users/users_pb';

export default defineComponent({
    components: {
        ChatBubbleLeftEllipsisIcon,
        TagIcon,
        UserCircleIcon,
        BoltIcon
    },
    data() {
        return {
            activities: [] as Array<UserActivity>,
            defaultIcon: UserCircleIcon,
        };
    },
    methods: {
        getUserActivity() {
            const req = new GetUserActivityRequest();
            req.setUserid(this.userID);

            getUsersClient().
                getUserActivity(req, null).then((resp) => {
                    this.activities = resp.getActivityList();
                }).catch((err: RpcError) => {
                    handleGRPCError(err, this.$route);
                });
        },
    },
    props: {
        userID: {
            required: true,
            type: Number,
        },
    },
    mounted() {
        this.getUserActivity();
    }
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
                            <h3 class="text-sm font-medium text-white">{{ activity.getCauseuser()?.getFirstname() }} {{
                                activity.getCauseuser()?.getLastname() }}</h3>
                            <p class="text-sm text-gray-400">{{ activity.getCreatedat() }}</p>
                        </div>
                        <p class="text-sm text-gray-300">{{ activity.getType() }} {{ activity.getKey() }}: {{
                            activity.getOldvalue() }} ⇒ {{ activity.getNewvalue() }}</p>
                    </div>
                </div>
            </li>
        </ul>
    </div>
</template>
