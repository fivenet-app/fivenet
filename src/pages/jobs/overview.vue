<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { default as TimeclockStatsBlock } from '~/components/jobs/timeclock/Stats.vue';
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

const { data: timeclockStats, pending, error } = useLazyAsyncData(`jobs-timeclock-stats`, () => getTimeclockStats());

async function getTimeclockStats(): Promise<TimeclockStats> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getJobsClient().timeclockStats({});
            const { response } = await call;

            if (!response.stats) return rej();
            return res(response.stats);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            throw e;
        }
    });
}
</script>

<template>
    <div>
        <div v-if="timeclockStats" class="grid grid-col-2 gap-2">
            <div class="sm:flex-auto">
                <Divider :label="$t('components.jobs.timeclock.Stats.title')" />
                <TimeclockStatsBlock :stats="timeclockStats ?? undefined" />
            </div>
        </div>
    </div>
</template>
