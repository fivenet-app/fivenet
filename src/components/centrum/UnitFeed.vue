<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import Time from '~/components/partials/elements/Time.vue';
import { UNIT_STATUS, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { mdiAccountPlus, mdiAccountRemove, mdiBriefcase, mdiCoffee, mdiHelp, mdiPlay, mdiStop } from '@mdi/js';

const props = defineProps<{
    unitId: bigint;
}>();

const { $grpc } = useNuxtApp();

const pagination = ref<PaginationResponse>();
const offset = ref(0n);

const {
    data: activity,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`centrum-unit-${props.unitId.toString()}-activity-${offset.value}`, () => listUnitActivity());

async function listUnitActivity(): Promise<Array<UnitStatus>> {
    return new Promise(async (res, rej) => {
        try {
            const req = {
                pagination: {
                    offset: offset.value,
                },
                id: props.unitId,
            };

            const call = $grpc.getCentrumClient().listUnitActivity(req);
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
    <div class="px-4 sm:px-6 lg:px-8 h-full overflow-y-scroll">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h1 class="text-base font-semibold leading-6 text-gray-100">Feed</h1>
            </div>
        </div>
        <div class="mt-2 flow-root">
            <div class="-mx-2 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
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
                            <template v-if="activityItem.status === UNIT_STATUS.UNKNOWN">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon type="mdi" :path="mdiHelp" class="h-6 w-6 text-primary-600" aria-hidden="true" />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">Unit created</p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === UNIT_STATUS.USER_ADDED">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon
                                        type="mdi"
                                        :path="mdiAccountPlus"
                                        class="h-6 w-6 text-primary-600"
                                        aria-hidden="true"
                                    />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Colleague added to Unit
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === UNIT_STATUS.USER_REMOVED">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon
                                        type="mdi"
                                        :path="mdiAccountRemove"
                                        class="h-6 w-6 text-primary-600"
                                        aria-hidden="true"
                                    />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Colleague removed from Unit
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === UNIT_STATUS.UNAVAILABLE">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon type="mdi" :path="mdiStop" class="h-6 w-6 text-primary-600" aria-hidden="true" />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Unit unavailable
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === UNIT_STATUS.AVAILABLE">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon type="mdi" :path="mdiPlay" class="h-6 w-6 text-primary-600" aria-hidden="true" />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Unit available
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === UNIT_STATUS.ON_BREAK">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon type="mdi" :path="mdiCoffee" class="h-6 w-6 text-primary-600" aria-hidden="true" />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Unit on break
                                    <span class="font-medium text-gray-400">
                                        {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                                    </span>
                                </p>
                                <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                                    <Time :value="activityItem.createdAt" />
                                </span>
                            </template>
                            <template v-else-if="activityItem.status === UNIT_STATUS.BUSY">
                                <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                                    <SvgIcon
                                        type="mdi"
                                        :path="mdiBriefcase"
                                        class="h-6 w-6 text-primary-600"
                                        aria-hidden="true"
                                    />
                                </div>
                                <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                                    Unit busy
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
