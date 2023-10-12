<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { ListDispatchActivityResponse } from '~~/gen/ts/services/centrum/centrum';
import FeedItem from './FeedItem.vue';

const props = defineProps<{
    dispatchId?: bigint;
}>();

const { $grpc } = useNuxtApp();

const offset = ref(0n);

const { data, refresh } = useLazyAsyncData(
    `centrum-dispatch-${(props.dispatchId ?? 0n).toString()}-activity-${offset.value}`,
    () => listDispatchActivity(),
);

const timer = setInterval(async () => refresh(), 3500);

async function listDispatchActivity(): Promise<ListDispatchActivityResponse> {
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

            return res(response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

onBeforeUnmount(() => {
    if (timer) clearInterval(timer);
});
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h2 class="text-base font-semibold leading-6 text-gray-100">Feed</h2>
            </div>
        </div>
        <div class="mt-2 flow-root">
            <div class="-mx-2 -my-2 overflow-y-auto sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <ul role="list" class="space-y-2">
                        <FeedItem
                            v-for="(activityItem, activityItemIdx) in data?.activity"
                            :key="activityItem.id.toString()"
                            :activityLength="data?.activity?.length ?? 0"
                            :item="activityItem"
                            :activityItemIdx="activityItemIdx"
                        />
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
