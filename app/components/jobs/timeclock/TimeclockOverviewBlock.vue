<script lang="ts" setup>
import TimeclockStatsBlock from '~/components/jobs/timeclock/TimeclockStatsBlock.vue';
import type { GetTimeclockStatsResponse } from '~~/gen/ts/services/jobs/timeclock';

const props = defineProps<{
    userId?: number;
}>();

const { $grpc } = useNuxtApp();

const { data, error, pending: loading, refresh } = useLazyAsyncData(`jobs-timeclock-stats`, () => getTimeclockStats());

async function getTimeclockStats(): Promise<GetTimeclockStatsResponse> {
    try {
        const call = $grpc.jobs.timeclock.getTimeclockStats({
            userId: props.userId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canRefresh = ref(true);
const refreshThrottle = useThrottleFn(async () => {
    canRefresh.value = false;
    await refresh().finally(() => useTimeoutFn(() => (canRefresh.value = true), 400));
}, 2500);
</script>

<template>
    <TimeclockStatsBlock
        :stats="data?.stats"
        :weekly="data?.weekly"
        :failed="!!error"
        :loading="loading"
        @refresh="refreshThrottle"
    />
</template>
