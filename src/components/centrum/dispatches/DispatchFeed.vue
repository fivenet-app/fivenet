<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useIntervalFn } from '@vueuse/core';
import { ListDispatchActivityResponse } from '~~/gen/ts/services/centrum/centrum';
import DispatchFeedItem from '~/components/centrum/dispatches/DispatchFeedItem.vue';

const props = defineProps<{
    dispatchId?: string;
}>();

const { $grpc } = useNuxtApp();

const offset = ref(0n);

const { data, refresh } = useLazyAsyncData(`centrum-dispatch-${props.dispatchId ?? '0'}-activity-${offset.value}`, () =>
    listDispatchActivity(),
);

async function listDispatchActivity(): Promise<ListDispatchActivityResponse> {
    try {
        const call = $grpc.getCentrumClient().listDispatchActivity({
            pagination: {
                offset: offset.value,
            },
            id: props.dispatchId ?? '0',
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { pause, resume } = useIntervalFn(async () => {
    pause();
    await refresh();
    resume();
}, 3500);
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
                        <DispatchFeedItem
                            v-for="(activityItem, activityItemIdx) in data?.activity"
                            :key="activityItem.id"
                            :activity-length="data?.activity?.length ?? 0"
                            :item="activityItem"
                            :activity-item-idx="activityItemIdx"
                        />
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
