<script lang="ts" setup>
import { AdjustmentsVerticalIcon } from '@heroicons/vue/24/outline';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { useNotificationsStore } from '~/store/notifications';
import { JobProps } from '~~/gen/ts/resources/jobs/jobs';
import DataErrorBlock from '../partials/DataErrorBlock.vue';
import DataPendingBlock from '../partials/DataPendingBlock.vue';

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const { t } = useI18n();

const properties = ref<{
    theme: string;
    livemapMarkerColor: string;
    quickButtons: {
        PenaltyCalculator: boolean;
    };
}>({
    theme: 'default',
    livemapMarkerColor: '#5C7AFF',
    quickButtons: {
        PenaltyCalculator: false,
    },
});

async function getJobProps(): Promise<JobProps> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().getJobProps({});
            const { response } = await call;

            if (response.jobProps) {
                properties.value.livemapMarkerColor = '#' + response.jobProps?.livemapMarkerColor;

                const components = response.jobProps!.quickButtons.split(';').filter((v) => v !== '');
                components.forEach((v) => {
                    switch (v) {
                        case 'PenaltyCalculator':
                            properties.value.quickButtons.PenaltyCalculator = true;
                            break;
                        default:
                            break;
                    }
                });
            }

            return res(response.jobProps!);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const { data: jobProps, pending, refresh, error } = useLazyAsyncData(`rector-jobprops`, () => getJobProps());

async function saveJobProps(): Promise<void> {
    return new Promise(async (res, rej) => {
        // How scuffed do you want this code to be: "Yes"
        let quickButtons = '';
        if (properties.value.quickButtons.PenaltyCalculator) {
            quickButtons += 'PenaltyCalculator;';
        }

        const jProps: JobProps = {
            job: '',
            theme: properties.value.theme,
            // Remove '#' from color code
            livemapMarkerColor: properties.value.livemapMarkerColor.substring(1),
            quickButtons: '',
        };

        jProps.quickButtons = quickButtons.replace(/;$/, '');

        try {
            await $grpc.getRectorClient().setJobProps({
                jobProps: jProps,
            });

            notifications.dispatchNotification({
                title: t('notifications.rector.job_props.title'),
                content: t('notifications.rector.job_props.content'),
                type: 'success',
            });

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div class="py-2 mt-5 max-w-5xl mx-auto">
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [`${$t('common.job', 1)} ${$t('common.prop')}`])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [`${$t('common.job', 1)} ${$t('common.prop')}`])"
            :retry="refresh"
        />
        <button
            v-else-if="!jobProps"
            type="button"
            class="relative block w-full p-12 text-center border-2 border-gray-300 border-dashed rounded-lg hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
        >
            <AdjustmentsVerticalIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold"> </span>
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
                                {{ jobProps.theme }}
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('components.rector.job_props.livemap_marker_color') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <input type="color" v-model="properties.livemapMarkerColor" />
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('components.rector.job_props.quick_buttons') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <fieldset>
                                    <div class="space-y-5">
                                        <div class="relative flex items-start">
                                            <div class="flex h-6 items-center">
                                                <input
                                                    aria-describedby="comments-description"
                                                    name="comments"
                                                    type="checkbox"
                                                    v-model="properties.quickButtons.PenaltyCalculator"
                                                    class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-600"
                                                />
                                            </div>
                                            <div class="ml-3 text-sm leading-6">
                                                <label for="comments" class="font-medium text-white">
                                                    {{ $t('components.penaltycalculator.title') }}
                                                </label>
                                            </div>
                                        </div>
                                    </div>
                                </fieldset>
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5" v-can="'RectorService.SetJobProps'">
                            <dt class="text-sm font-medium"></dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <button
                                    type="button"
                                    @click="saveJobProps()"
                                    class="rounded-md bg-green-600 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-green-400"
                                >
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
