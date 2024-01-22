<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import TimeclockStatsBlock from '~/components/jobs/timeclock/TimeclockStatsBlock.vue';
import GenericContainer from '~/components/partials/elements/GenericContainer.vue';
import GenericDivider from '~/components/partials/elements/GenericDivider.vue';
import { useAuthStore } from '~/store/auth';
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

const authStore = useAuthStore();

const { activeChar, jobProps } = storeToRefs(authStore);

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
        <div class="py-2 pb-14">
            <div class="px-1 sm:px-2 lg:px-4">
                <div class="grid-col-2 grid gap-2">
                    <div class="sm:flex-auto">
                        <div class="flex flex-row gap-2">
                            <GenericContainer class="flex-1">
                                <h1 class="text-3xl font-semibold leading-6 text-neutral">
                                    {{ activeChar?.jobLabel }}
                                </h1>
                                <h2 class="mt-2 text-xl font-semibold leading-6 text-neutral">
                                    {{ $t('common.rank') }}: {{ activeChar?.jobGradeLabel }}
                                </h2>
                            </GenericContainer>

                            <GenericContainer v-if="jobProps?.radioFrequency" class="text-neutral">
                                <h3 class="text-lg font-semibold">{{ $t('common.radio_frequency') }}</h3>
                                <p class="text-center text-lg font-bold">{{ jobProps?.radioFrequency }}.00</p>
                            </GenericContainer>
                        </div>

                        <GenericDivider :label="$t('components.jobs.timeclock.Stats.title')" />
                        <TimeclockStatsBlock :stats="timeclockStats" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
