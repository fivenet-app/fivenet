<script lang="ts" setup>
import { JobProps } from '@fivenet/gen/resources/jobs/jobs_pb';
import { RpcError } from 'grpc-web';
import { GetJobPropsRequest, SetJobPropsRequest } from '@fivenet/gen/services/rector/rector_pb';
import DataPendingBlock from '../partials/DataPendingBlock.vue';
import DataErrorBlock from '../partials/DataErrorBlock.vue';
import { AdjustmentsVerticalIcon } from '@heroicons/vue/24/outline';

const { $grpc } = useNuxtApp();

async function getJobProps(): Promise<JobProps> {
    return new Promise(async (res, rej) => {
        const req = new GetJobPropsRequest();

        try {
            const resp = await $grpc.getRectorClient().
                getJobProps(req, null);

            return res(resp.getJobProps()!);
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const { data: jobProps, pending, refresh, error } = await useLazyAsyncData(`rector-jobprops`, () => getJobProps());

const props = ref<{ theme: string; livemapMarkerColor: string; }>({ theme: 'default', livemapMarkerColor: '#5C7AFF' });

async function saveJobProps(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new SetJobPropsRequest();
        const jProps = new JobProps();
        jProps.setTheme(props.value.theme);
        // Remove '#' from color code
        jProps.setLivemapMarkerColor(props.value.livemapMarkerColor.substring(1));
        req.setJobProps(jProps);

        try {
            await $grpc.getRectorClient().
                setJobProps(req, null);

            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watch(jobProps, () => props.value.livemapMarkerColor = '#' + jobProps.value?.getLivemapMarkerColor());
</script>

<template>
    <div class="py-2 mt-5 max-w-5xl mx-auto">
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [`${$t('common.job', 1)} ${$t('common.prop')}`])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [`${$t('common.job', 1)} ${$t('common.prop')}`])" :retry="refresh" />
        <button v-else-if="!jobProps" type="button"
            class="relative block w-full p-12 text-center border-2 border-gray-300 border-dashed rounded-lg hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
            <AdjustmentsVerticalIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold">

            </span>
        </button>
        <div v-else>
            <div class="overflow-hidden bg-base-800 shadow sm:rounded-lg text-neutral">
                <div class="px-4 py-5 sm:px-6">
                    <h3 class="text-base font-semibold leading-6">
                        {{ $t('components.rector.job_props.job_properties') }}
                    </h3>
                    <p class="mt-1 max-w-2xl text-sm">
                        {{ $t('components.rector.job_props.your_job_properties') }}
                    </p>
                </div>
                <div class="border-t border-base-400 px-4 py-5 sm:p-0">
                    <dl class="sm:divide-y sm:divide-base-400">
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('common.theme') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                {{ jobProps.getTheme() }}
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('components.rector.job_props.livemap_marker_color') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <input type="color" v-model="props.livemapMarkerColor" />
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5"
                            v-can="'RectorService.SetJobProps'">
                            <dt class="text-sm font-medium"></dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <button type="button" @click="saveJobProps()"
                                    class="rounded-md bg-green-600 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-green-400">
                                    {{ $t('common.save', 1) }}
                                </button>
                            </dd>
                        </div>
                    </dl>
                </div>
            </div>
        </div>
    </div>
</template>
