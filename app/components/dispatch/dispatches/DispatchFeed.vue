<script lang="ts" setup>
import DispatchFeedItem from '~/components/dispatch/dispatches/DispatchFeedItem.vue';
import { dispatchStatuses, dispatchStatusToBGColor } from '~/components/dispatch/helpers';
import { getCentrumDispatchesClient } from '~~/gen/ts/clients';
import type { ListDispatchActivityResponse } from '~~/gen/ts/services/centrum/dispatches';

const props = defineProps<{
    dispatchId?: number | undefined;
}>();

const page = defineModel<number>('page', { default: 1 });

const centrumDispatchesClientPromise = getCentrumDispatchesClient();

const lastSuccessfulData = shallowRef<ListDispatchActivityResponse>();
const activityData = computed(() => data.value ?? lastSuccessfulData.value);
const pageSize = computed(() => activityData.value?.pagination?.pageSize ?? 10);
const offset = computed(() => (page.value - 1) * pageSize.value);
const dispatchId = computed(() => props.dispatchId ?? 0);
const hasDispatchId = computed(() => dispatchId.value > 0);

const timelineItems = computed(() =>
    (activityData.value?.activity ?? []).map((item) => ({
        ...item,
        icon: dispatchStatuses.find((status) => status.status === item.status)?.icon ?? 'i-mdi-info-circle',
        ui: { indicator: 'text-highlighted ' + dispatchStatusToBGColor(item.status) },
    })),
);

const activityKey = computed(() => `centrum-dispatch-${dispatchId.value}-activity-page-${page.value}`);

const { data, status, refresh } = useAuthedLazyAsyncData(
    'userState',
    activityKey,
    ({ signal }) => listDispatchActivity(signal),
    {
        default: (): ListDispatchActivityResponse => ({ activity: [] }),
        immediate: false,
    },
);

watch(
    data,
    (value) => {
        if (value !== undefined) {
            lastSuccessfulData.value = value;
        }
    },
    { immediate: true },
);

async function listDispatchActivity(signal: AbortSignal): Promise<ListDispatchActivityResponse> {
    if (!hasDispatchId.value) {
        return { activity: [] };
    }

    try {
        const centrumDispatchesClient = await centrumDispatchesClientPromise;
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

defineExpose({
    pagination: computed(() => activityData.value?.pagination),
    refresh,
    status,
});
</script>

<template>
    <UTimeline :items="timelineItems" size="xs" :ui="{ wrapper: '!mt-0 !pb-2' }">
        <template #wrapper="{ item }">
            <DispatchFeedItem :item="item" />
        </template>
    </UTimeline>
</template>
