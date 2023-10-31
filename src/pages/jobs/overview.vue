<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import TimeclockStatsBlock from '~/components/jobs/timeclock/TimeclockStatsBlock.vue';
import Divider from '~/components/partials/elements/Divider.vue';
import { TimeclockStats } from '~~/gen/ts/resources/jobs/timeclock';

useHead({
    title: 'pages.jobs.overview.title',
});
definePageMeta({
    title: 'common.overview',
    requiresAuth: true,
    permission: 'JobsService.ColleaguesList',
});

const { $grpc } = useNuxtApp();

const { data: timeclockStats } = useLazyAsyncData(`jobs-timeclock-stats`, () => getTimeclockStats());

async function getTimeclockStats(): Promise<TimeclockStats> {
    try {
        const call = $grpc.getJobsClient().timeclockStats({});
        const { response } = await call;

        return response.stats!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <div>
        <div v-if="timeclockStats" class="grid grid-col-2 gap-2">
            <div class="sm:flex-auto">
                <Divider :label="$t('components.jobs.timeclock.Stats.title')" />
                <TimeclockStatsBlock :stats="timeclockStats" />
            </div>
        </div>
    </div>
</template>
