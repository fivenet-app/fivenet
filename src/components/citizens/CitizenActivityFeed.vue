<script lang="ts">
import { defineComponent } from 'vue';
import { BoltIcon, ChatBubbleLeftEllipsisIcon, TagIcon, UserCircleIcon } from '@heroicons/vue/20/solid'
import { UsersServiceClient } from '@arpanet/gen/users/UsersServiceClientPb';
import config from '../../config';
import { clientAuthOptions, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';
import { GetUserActivityRequest, UserActivity } from '@arpanet/gen/users/users_pb';

export default defineComponent({
    components: {
        ChatBubbleLeftEllipsisIcon,
        TagIcon,
        UserCircleIcon,
        BoltIcon
    },
    data() {
        return {
            'activities': [] as Array<UserActivity>,
            'client': new UsersServiceClient(config.apiProtoURL, null, clientAuthOptions),
        };
    },
    methods: {
        getUserActivity() {
            const req = new GetUserActivityRequest();
            req.setUserid(this.userID);
            this.client.getUserActivity(req, null).then((resp) => {
                console.log(resp.getActivityList())
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
    <div class="flow-root">
        <ul role="list" class="-mb-8">
            <li v-for="(activity, activityIdx) in activities" :key="activity.getId()">
                <div class="relative pb-8">
                    <span v-if="activityIdx !== activities.length - 1"
                        class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" aria-hidden="true" />
                    <div class="relative flex space-x-3">
                        <div>
                            <span class="h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white">
                                <BoltIcon class="h-5 w-5 text-white" aria-hidden="true" />
                            </span>
                        </div>
                        <div class="flex min-w-0 flex-1 justify-between space-x-4 pt-1.5">
                            <div>
                                <p class="text-sm text-gray-500">
                                    {{ activity.getCauseuser() }} <a href="#" class="font-medium text-gray-900">{{
                                        activity.getTargetuser() }}</a>
                                </p>
                            </div>
                            <div class="whitespace-nowrap text-right text-sm text-gray-500">
                                <time :datetime="activity.getCreatedat()">{{ activity.getCreatedat() }}</time>
                            </div>
                        </div>
                    </div>
                </div>
            </li>
        </ul>
    </div>
</template>
