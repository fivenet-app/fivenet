<script lang="ts" setup>
import DispatchFeedItem from '~/components/dispatch/dispatches/DispatchFeedItem.vue';
import { getCentrumDispatchesClient } from '~~/gen/ts/clients';
import type { ListDispatchActivityResponse } from '~~/gen/ts/services/centrum/dispatches';

const props = defineProps<{
    dispatchId?: number | undefined;
}>();

const centrumDispatchesClient = await getCentrumDispatchesClient();

const offset = ref(0);
const dispatchId = computed(() => props.dispatchId ?? 0);
const hasDispatchId = computed(() => dispatchId.value > 0);

const activityKey = computed(() => `centrum-dispatch-${dispatchId.value}-activity-${offset.value}`);

const { data, refresh } = useLazyAsyncData(activityKey, () => listDispatchActivity(), {
    default: () => ({ activity: [] }),
    immediate: false,
});

async function listDispatchActivity(): Promise<ListDispatchActivityResponse> {
    if (!hasDispatchId.value) {
        return { activity: [] };
    }

    try {
        const call = centrumDispatchesClient.listDispatchActivity({
            pagination: {
                offset: offset.value,
            },
            id: dispatchId.value,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { pause, resume } = useIntervalFn(
    async () => {
        if (!hasDispatchId.value) {
            pause();
            return;
        }

        pause();
        try {
            await refresh();
        } finally {
            if (hasDispatchId.value) {
                resume();
            }
        }
    },
    3500,
    { immediate: false },
);

watch(
    dispatchId,
    async (id) => {
        pause();

        if (id <= 0) {
            return;
        }

        await refresh();
        resume();
    },
    { immediate: true },
);
</script>

<template>
    <div class="my-1 flex h-full flex-1 grow flex-col gap-2 px-1">
        <div class="flex justify-between">
            <h2 class="inline-flex flex-1 items-center text-base leading-6 font-semibold">{{ $t('common.feed') }}</h2>
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
