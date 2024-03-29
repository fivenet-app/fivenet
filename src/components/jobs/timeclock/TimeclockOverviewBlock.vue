<script lang="ts" setup>
import type { GetTimeclockStatsResponse } from '~~/gen/ts/services/jobs/timeclock';
import TimeclockStatsBlock from '~/components/jobs/timeclock/TimeclockStatsBlock.vue';

const props = defineProps<{
    userId?: number;
}>();

const { $grpc } = useNuxtApp();

const { data: timeclockStats, error, refresh } = useLazyAsyncData(`jobs-timeclock-stats`, () => getTimeclockStats());

async function getTimeclockStats(): Promise<GetTimeclockStatsResponse> {
    try {
        const call = $grpc.getJobsTimeclockClient().getTimeclockStats({
            userId: props.userId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <TimeclockStatsBlock
        class="mt-4"
        :stats="timeclockStats?.stats"
        :weekly="timeclockStats?.weekly"
        :failed="error !== null"
        @refresh="refresh()"
    />
</template>
