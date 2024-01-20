<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import TimeclockStatsBlock from '~/components/jobs/timeclock/TimeclockStatsBlock.vue';
import GenericDivider from '~/components/partials/elements/GenericDivider.vue';
import { TimeclockStats } from '~~/gen/ts/resources/jobs/timeclock';

useHead({
    title: 'pages.jobs.overview.title',
});
definePageMeta({
    title: 'common.overview',
    requiresAuth: true,
    permission: 'JobsService.ListColleagues',
});

const { $grpc } = useNuxtApp();

const { data: timeclockStats } = useLazyAsyncData(`jobs-timeclock-stats`, () => getTimeclockStats());

async function getTimeclockStats(): Promise<TimeclockStats> {
    try {
        const call = $grpc.getJobsTimeclockClient().getTimeclockStats({});
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
        <div class="grid-col-2 grid gap-2">
            <div class="sm:flex-auto">
                <GenericDivider :label="$t('components.jobs.timeclock.Stats.title')" />
                <TimeclockStatsBlock :stats="timeclockStats" />
            </div>
        </div>
    </div>
</template>
