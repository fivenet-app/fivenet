<script lang="ts" setup>
import DispatchFeedItem from '~/components/centrum/dispatches/DispatchFeedItem.vue';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import type { ListDispatchActivityResponse } from '~~/gen/ts/services/centrum/centrum';

const props = defineProps<{
    dispatchId?: number;
}>();

const centrumCentrumClient = await getCentrumCentrumClient();

const offset = ref(0);

const { data, refresh } = useLazyAsyncData(`centrum-dispatch-${props.dispatchId ?? 0}-activity-${offset.value}`, () =>
    listDispatchActivity(),
);

async function listDispatchActivity(): Promise<ListDispatchActivityResponse> {
    try {
        const call = centrumCentrumClient.listDispatchActivity({
            pagination: {
                offset: offset.value,
            },
            id: props.dispatchId ?? 0,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
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
    <div class="flex h-full flex-1 grow flex-col px-1">
        <div class="flex justify-between">
            <h2 class="inline-flex flex-1 items-center text-base font-semibold leading-6">{{ $t('common.feed') }}</h2>
        </div>

        <div class="flex flex-1 flex-col overflow-x-auto overflow-y-auto">
            <ul class="space-y-2" role="list">
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
</template>
