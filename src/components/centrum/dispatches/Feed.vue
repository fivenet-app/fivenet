<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import Time from '~/components/partials/elements/Time.vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { DispatchStatus, DISPATCH_STATUS } from '~~/gen/ts/resources/dispatch/dispatch';
import {
    mdiAccountCancel,
    mdiAccountPlus,
    mdiAccountRemove,
    mdiCar,
    mdiCheck,
    mdiCloseCircle,
    mdiHelp,
    mdiMapMarker,
    mdiNewBox,
    mdiUpdate,
} from '@mdi/js';

const props = defineProps<{
    dispatchId?: bigint;
}>();

const { $grpc } = useNuxtApp();

const pagination = ref<PaginationResponse>();
const offset = ref(0n);

const {
    data: activity,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`centrum-dispatch-${(props.dispatchId ?? 0n).toString()}-activity-${offset.value}`, () =>
    listDispatchActivity(),
);

async function listDispatchActivity(): Promise<Array<DispatchStatus>> {
    return new Promise(async (res, rej) => {
        try {
            const req = {
                pagination: {
                    offset: offset.value,
                },
                id: props.dispatchId ?? 0n,
            };

            const call = $grpc.getCentrumClient().listDispatchActivity(req);
            const { response } = await call;

            pagination.value = response.pagination;
            return res(response.activity);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h1 class="text-base font-semibold leading-6 text-gray-100">Feed</h1>
            </div>
        </div>
        <div class="mt-2 flow-root">
            <div class="-mx-2 -my-2 overflow-y-scroll sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <ul role="list" class="space-y-2">
                        <li
                            v-for="(activityItem, activityItemIdx) in activity"
                            :key="activityItem.id.toString()"
                            class="relative flex gap-x-2"
                        >
                            <div
                                :class="[
                                    activity !== null && activityItemIdx === activity.length - 1 ? 'h-6' : '-bottom-6',
                                    'absolute left-0 top-0 flex w-6 justify-center',
                                ]"
                            >
                                <div class="w-px bg-gray-200" />
                            </div>
                            <template v-if="activityItem.status === DISPATCH_STATUS.NEW">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon type="mdi" :path="mdiNewBox" class="h-6 w-6 text-primary-600" aria-hidden="true" />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">Dispatch created</p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === DISPATCH_STATUS.DECLINED">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon
                                        type="mdi"
                                        :path="mdiAccountCancel"
                                        class="h-6 w-6 text-primary-600"
                                        aria-hidden="true"
                                    />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Dispatch declined by
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === DISPATCH_STATUS.UNASSIGNED">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon
                                        type="mdi"
                                        :path="mdiAccountRemove"
                                        class="h-6 w-6 text-primary-600"
                                        aria-hidden="true"
                                    />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Dispatch unassigned by
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === DISPATCH_STATUS.UNIT_ASSIGNED">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon
                                        type="mdi"
                                        :path="mdiAccountPlus"
                                        class="h-6 w-6 text-primary-600"
                                        aria-hidden="true"
                                    />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Dispatch accepted
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === DISPATCH_STATUS.EN_ROUTE">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon type="mdi" :path="mdiCar" class="h-6 w-6 text-primary-600" aria-hidden="true" />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    En Route to Dispatch
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === DISPATCH_STATUS.AT_SCENE">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon
                                        type="mdi"
                                        :path="mdiMapMarker"
                                        class="h-6 w-6 text-primary-600"
                                        aria-hidden="true"
                                    />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Arrived on scene at Dispatch
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === DISPATCH_STATUS.NEED_ASSISTANCE">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon type="mdi" :path="mdiHelp" class="h-6 w-6 text-primary-600" aria-hidden="true" />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Need Assistance
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === DISPATCH_STATUS.COMPLETED">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon type="mdi" :path="mdiCheck" class="h-6 w-6 text-primary-600" aria-hidden="true" />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Dispatch completed
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
