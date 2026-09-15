<script lang="ts" setup>
import DispatchFeedItem from '~/components/dispatch/dispatches/DispatchFeedItem.vue';
import { dispatchStatuses, dispatchStatusToBGColor } from '~/components/dispatch/helpers';
import { getCentrumDispatchesClient } from '~~/gen/ts/clients';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import type { ListDispatchActivityResponse } from '~~/gen/ts/services/centrum/dispatches';

const props = defineProps<{
    dispatchId?: number | undefined;
}>();

const centrumDispatchesClient = await getCentrumDispatchesClient();

const offset = ref(0);
const dispatchId = computed(() => props.dispatchId ?? 0);
const hasDispatchId = computed(() => dispatchId.value > 0);

const timelineItems = computed(() =>
    (data.value?.activity ?? []).map((item) => ({
        ...item,
        icon: dispatchStatuses.find((status) => status.status === item.status)?.icon ?? 'i-mdi-info-circle',
        ui: { indicator: 'text-highlighted ' + dispatchStatusToBGColor(item.status) },
    })),
);

const activityKey = computed(() => `centrum-dispatch-${dispatchId.value}-activity-${offset.value}`);

const { data, status, refresh } = useAuthedLazyAsyncData(
    'userState',
    activityKey,
    ({ signal }) => listDispatchActivity(signal),
    {
        default: () => ({ activity: [] }),
        immediate: false,
    },
);

async function listDispatchActivity(signal: AbortSignal): Promise<ListDispatchActivityResponse> {
    if (!hasDispatchId.value) {
        return { activity: [] };
    }

    try {
        const call = centrumDispatchesClient.listDispatchActivity(
            {
                pagination: {
                    offset: offset.value,
                },
                id: dispatchId.value,
            },
            { abort: signal },
        );
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
    <div class="my-1 flex flex-col gap-2 px-1">
        <div class="flex justify-between">
            <h2 class="inline-flex flex-1 items-center text-base leading-6 font-semibold">{{ $t('common.feed') }}</h2>

            <RefreshButton icon-only :disabled="!hasDispatchId" :loading="isRequestPending(status)" @click="() => refresh()" />
        </div>

        <div class="flex flex-col">
            <UTimeline :items="timelineItems" size="xs" :ui="{ wrapper: '!mt-0 !pb-2' }">
                <template #wrapper="{ item }">
                    <DispatchFeedItem :item="item" />
                </template>
            </UTimeline>
        </div>
    </div>
</template>
