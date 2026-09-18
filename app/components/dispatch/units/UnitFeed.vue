<script lang="ts" setup>
import UnitFeedItem from '~/components/dispatch/units/UnitFeedItem.vue';
import { unitStatuses, unitStatusToBGColor } from '~/components/dispatch/helpers';
import { getCentrumUnitsClient } from '~~/gen/ts/clients';
import type { ListUnitActivityResponse } from '~~/gen/ts/services/centrum/units';

const props = defineProps<{
    unitId: number;
}>();

const page = defineModel<number>('page', { default: 1 });

const centrumUnitsClientPromise = getCentrumUnitsClient();

const lastSuccessfulData = shallowRef<ListUnitActivityResponse>();
const activityData = computed(() => data.value ?? lastSuccessfulData.value);
const pageSize = computed(() => activityData.value?.pagination?.pageSize ?? 10);
const offset = computed(() => (page.value - 1) * pageSize.value);

const timelineItems = computed(() =>
    (activityData.value?.activity ?? []).map((item) => ({
        ...item,
        icon: unitStatuses.find((status) => status.status === item.status)?.icon ?? 'i-mdi-info-circle',
        ui: { indicator: 'text-highlighted ' + unitStatusToBGColor(item.status) },
    })),
);

const activityKey = computed(() => `centrum-unit-${props.unitId}-activity-page-${page.value}`);

const { data, status, refresh } = useAuthedLazyAsyncData('userState', activityKey, ({ signal }) => listUnitActivity(signal));

watch(
    data,
    (value) => {
        if (value !== undefined) {
            lastSuccessfulData.value = value;
        }
    },
    { immediate: true },
);

async function listUnitActivity(signal: AbortSignal): Promise<ListUnitActivityResponse> {
    try {
        const centrumUnitsClient = await centrumUnitsClientPromise;
        const call = centrumUnitsClient.listUnitActivity(
            {
                pagination: {
                    offset: offset.value,
                },
                id: props.unitId,
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

const { pause, resume } = useIntervalFn(async () => {
    pause();
    try {
        await refresh();
    } finally {
        resume();
    }
}, 3500);

defineExpose({
    pagination: computed(() => activityData.value?.pagination),
    refresh,
    status,
});
</script>

<template>
    <UTimeline :items="timelineItems" size="xs" :ui="{ wrapper: '!mt-0 !pb-2' }">
        <template #wrapper="{ item }">
            <UnitFeedItem :item="item" />
        </template>
    </UTimeline>
</template>
