<script lang="ts" setup>
import { UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import FeedItem from './FeedItem.vue';

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
                        <FeedItem
                            v-for="(activityItem, activityItemIdx) in activity"
                            :key="activityItem.id.toString()"
                            :activityLength="activity?.length ?? 0"
                            :activityItem="activityItem"
                            :activityItemIdx="activityItemIdx"
                        />
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
