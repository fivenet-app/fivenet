<script lang="ts" setup>
import type * as L from 'leaflet';
import HeatmapLegend from '~/components/livemap/controls/HeatmapLegend.vue';
import type { GetDispatchHeatmapResponse } from '~~/gen/ts/services/centrum/centrum';

const props = withDefaults(
    defineProps<{
        show?: boolean;
    }>(),
    {
        show: false,
    },
);

const { $grpc } = useNuxtApp();

const {
    data: heatmap,
    status,
    refresh,
} = useLazyAsyncData(`centrum-heatmap`, () => getDispatchHeatmap(), { immediate: props.show });

async function getDispatchHeatmap(): Promise<GetDispatchHeatmapResponse> {
    try {
        const call = $grpc.centrum.centrum.getDispatchHeatmap({
            status: [],
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const heat = inject<Ref<L.HeatLayer | undefined>>('heat');

async function handleProps(show: boolean): Promise<void> {
    if (show) {
        // If we have no heatmap data, or it's pending, refresh it
        if (!heatmap.value || isRequestPending(status.value)) {
            await refresh();
        }

        heatmap.value?.entries.forEach((e) => heat?.value?.addLatLng([e.y, e.x, e.w]));
    } else {
        heat?.value?.setLatLngs([]);
    }
}

watch(() => props.show, handleProps);

onBeforeMount(() => handleProps(props.show));
</script>

<template>
    <HeatmapLegend :show="show" :max="heatmap?.maxEntries" />
</template>
