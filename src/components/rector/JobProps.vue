<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { TuneIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useConfigStore } from '~/store/config';
import { useNotificatorStore } from '~/store/notificator';
import { JobProps } from '~~/gen/ts/resources/users/jobs';

const { $grpc } = useNuxtApp();

const configStore = useConfigStore();
const { appConfig } = storeToRefs(configStore);

const notifications = useNotificatorStore();

const properties = ref<{
    theme: string;
    livemapMarkerColor: string;
    quickButtons: {
        PenaltyCalculator: boolean;
    };
    discordGuildId?: string;
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

                if (response.jobProps.discordGuildId && response.jobProps.discordGuildId > 0) {
                    properties.value.discordGuildId = response.jobProps.discordGuildId?.toString();
                }
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
            discordGuildId:
                properties.value.discordGuildId && properties.value.discordGuildId !== ''
                    ? BigInt(properties.value.discordGuildId)
                    : undefined,
        };

        jProps.quickButtons = quickButtons.replace(/;$/, '');

        try {
            await $grpc.getRectorClient().setJobProps({
                jobProps: jProps,
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.rector.job_props.title', parameters: [] },
                content: { key: 'notifications.rector.job_props.content', parameters: [] },
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
        <DataNoDataBlock v-else-if="!jobProps" :icon="TuneIcon" :type="`${$t('common.job', 1)} ${$t('common.prop')}`" />
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
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('components.rector.job_props.discord_guild_id') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <input
                                    type="text"
                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    :class="appConfig.discord.botInviteURL ? 'disabled' : ''"
                                    :disabled="appConfig.discord.botInviteURL === undefined"
                                    :placeholder="$t('components.rector.job_props.discord_guild_id')"
                                    :label="$t('components.rector.job_props.discord_guild_id')"
                                    maxlength="70"
                                    v-model="properties.discordGuildId"
                                />
                                <NuxtLink
                                    v-if="appConfig.discord.botInviteURL !== undefined"
                                    :to="appConfig.discord.botInviteURL"
                                    :external="true"
                                    class="mt-2 flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300 bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500"
                                >
                                    {{ $t('components.rector.job_props.invite_bot') }}
                                </NuxtLink>
                            </dd>
                        </div>
                        <div
                            v-if="can('RectorService.SetJobProps')"
                            class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5"
                        >
                            <dt class="text-sm font-medium"></dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <button
                                    type="button"
                                    @click="saveJobProps()"
                                    class="rounded-md bg-success-600 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-success-400"
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
