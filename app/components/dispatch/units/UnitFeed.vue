<script lang="ts" setup>
import UnitFeedItem from '~/components/dispatch/units/UnitFeedItem.vue';
import { unitStatuses, unitStatusToBGColor } from '~/components/dispatch/helpers';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import { getCentrumUnitsClient } from '~~/gen/ts/clients';
import type { ListUnitActivityResponse } from '~~/gen/ts/services/centrum/units';

const props = defineProps<{
    unitId: number;
}>();

const centrumUnitsClient = await getCentrumUnitsClient();

const offset = ref(0);

const timelineItems = computed(() =>
    (data.value?.activity ?? []).map((item) => ({
        ...item,
        icon: unitStatuses.find((status) => status.status === item.status)?.icon ?? 'i-mdi-info-circle',
        ui: { indicator: 'text-highlighted ' + unitStatusToBGColor(item.status) },
    })),
);

const { data, status, refresh } = useAuthedLazyAsyncData(
    'userState',
    `centrum-unit-${props.unitId}-activity-${offset.value}`,
    ({ signal }) => listUnitActivity(signal),
);

async function listUnitActivity(signal: AbortSignal): Promise<ListUnitActivityResponse> {
    try {
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
</script>

<template>
    <div class="my-1 flex flex-col gap-2 px-1">
        <div class="flex justify-between">
            <h2 class="inline-flex flex-1 items-center text-base leading-6 font-semibold">{{ $t('common.feed') }}</h2>

            <RefreshButton icon-only :loading="isRequestPending(status)" @click="() => refresh()" />
        </div>

        <div class="flex flex-col">
            <UTimeline :items="timelineItems" size="xs" :ui="{ wrapper: '!mt-0 !pb-2' }">
                <template #wrapper="{ item }">
                    <UnitFeedItem :item="item" />
                </template>
            </UTimeline>
        </div>
    </div>
</template>
