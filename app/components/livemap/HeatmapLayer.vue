<script lang="ts" setup>
import type { HeatLayer } from 'leaflet';
import HeatmapLegend from '~/components/livemap/controls/HeatmapLegend.vue';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import type { GetDispatchHeatmapResponse } from '~~/gen/ts/services/centrum/centrum';

const props = withDefaults(
    defineProps<{
        show?: boolean;
    }>(),
    {
        show: false,
    },
);

const centrumCentrumClient = await getCentrumCentrumClient();

const {
    data: heatmap,
    status,
    refresh,
} = useAuthedLazyAsyncData('userState', `centrum-heatmap`, ({ signal }) => getDispatchHeatmap(signal), {
    immediate: props.show,
});

async function getDispatchHeatmap(signal: AbortSignal): Promise<GetDispatchHeatmapResponse> {
    try {
        const call = centrumCentrumClient.getDispatchHeatmap(
            {
                status: [],
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

const heat = inject<Ref<HeatLayer | undefined>>('heat');
let updateGeneration = 0;

async function handleProps(show: boolean): Promise<void> {
    const generation = ++updateGeneration;

    if (show) {
        // If we have no heatmap data, or it's still pending, refresh it
        if (!heatmap.value || isRequestPending(status.value)) await refresh();

        // The heat layer is created asynchronously by BaseMap. If this runs
        // before that happens, the watcher will retry when the layer becomes ready.
        if (generation !== updateGeneration || !heat?.value) return;

        heat.value.setLatLngs(heatmap.value?.entries.map((e) => [e.y, e.x, e.w]) ?? []);
        return;
    }

    heat?.value?.setLatLngs([]);
}

watch([() => props.show, () => heat?.value], ([show]) => handleProps(show), { immediate: true });
</script>

<template>
    <HeatmapLegend :show="show" :max="heatmap?.maxEntries" />
</template>
