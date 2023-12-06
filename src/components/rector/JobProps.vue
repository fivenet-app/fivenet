<script lang="ts" setup>
import { Switch, SwitchGroup, SwitchLabel } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useThrottleFn } from '@vueuse/core';
import { LoadingIcon, TuneIcon } from 'mdi-vue3';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Time from '~/components/partials/elements/Time.vue';
import { useNotificatorStore } from '~/store/notificator';
import { JobProps } from '~~/gen/ts/resources/users/jobs';

const { $grpc } = useNuxtApp();

const appConfig = useAppConfig();

const notifications = useNotificatorStore();

async function getJobProps(): Promise<JobProps> {
    try {
        const call = $grpc.getRectorClient().getJobProps({});
        const { response } = await call;

        if (response.jobProps) {
            if (response.jobProps.quickButtons === undefined) {
                response.jobProps.quickButtons = {
                    penaltyCalculator: false,
                    bodyCheckup: false,
                };
            }
            if (response.jobProps.discordSyncSettings === undefined) {
                response.jobProps.discordSyncSettings = {
                    userInfoSync: false,
                };
            }
        }

        return response.jobProps!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { data: jobProps, pending, refresh, error } = useLazyAsyncData(`rector-jobprops`, () => getJobProps());

async function saveJobProps(): Promise<void> {
    if (!jobProps.value) {
        return;
    }

    try {
        await $grpc.getRectorClient().setJobProps({
            jobProps: jobProps.value,
        });

        notifications.dispatchNotification({
            title: { key: 'notifications.rector.job_props.title', parameters: {} },
            content: { key: 'notifications.rector.job_props.content', parameters: {} },
            type: 'success',
        });
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (_) => {
    canSubmit.value = false;
    await saveJobProps().finally(() => setTimeout(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div class="p-2">
        <div class="max-w-5xl mx-auto">
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
                                    <ColorInput
                                        v-model="jobProps.livemapMarkerColor"
                                        disable-alpha
                                        format="hex"
                                        position="top"
                                    />
                                </dd>
                            </div>
                            <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium">
                                    {{ $t('common.radio_frequency') }}
                                </dt>
                                <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                    <input
                                        v-model="jobProps.radioFrequency"
                                        type="text"
                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.radio_frequency')"
                                        :label="$t('common.radio_frequency')"
                                        maxlength="6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </dd>
                            </div>
                            <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium">
                                    {{ $t('components.rector.job_props.quick_buttons') }}
                                </dt>
                                <dd v-if="jobProps.quickButtons" class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                    <fieldset class="flex flex-col gap-4">
                                        <div class="space-y-5">
                                            <SwitchGroup as="div" class="flex items-center">
                                                <Switch
                                                    v-model="jobProps.quickButtons.penaltyCalculator"
                                                    :class="[
                                                        jobProps.quickButtons.penaltyCalculator
                                                            ? 'bg-indigo-600'
                                                            : 'bg-gray-200',
                                                        'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2',
                                                    ]"
                                                >
                                                    <span
                                                        aria-hidden="true"
                                                        :class="[
                                                            jobProps.quickButtons.penaltyCalculator
                                                                ? 'translate-x-5'
                                                                : 'translate-x-0',
                                                            'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                        ]"
                                                    />
                                                </Switch>
                                                <SwitchLabel as="span" class="ml-3 text-sm">
                                                    <span class="font-medium text-gray-300">{{
                                                        $t('components.penaltycalculator.title')
                                                    }}</span>
                                                </SwitchLabel>
                                            </SwitchGroup>
                                        </div>
                                        <div class="space-y-5">
                                            <SwitchGroup as="div" class="flex items-center">
                                                <Switch
                                                    v-model="jobProps.quickButtons.bodyCheckup"
                                                    :class="[
                                                        jobProps.quickButtons.bodyCheckup ? 'bg-indigo-600' : 'bg-gray-200',
                                                        'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2',
                                                    ]"
                                                >
                                                    <span
                                                        aria-hidden="true"
                                                        :class="[
                                                            jobProps.quickButtons.bodyCheckup
                                                                ? 'translate-x-5'
                                                                : 'translate-x-0',
                                                            'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                        ]"
                                                    />
                                                </Switch>
                                                <SwitchLabel as="span" class="ml-3 text-sm">
                                                    <span class="font-medium text-gray-300">{{
                                                        $t('components.bodycheckup.title')
                                                    }}</span>
                                                </SwitchLabel>
                                            </SwitchGroup>
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
                                        v-model="jobProps.discordGuildId"
                                        type="text"
                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :class="appConfig.discord.botInviteURL === undefined ? 'disabled' : ''"
                                        :disabled="appConfig.discord.botInviteURL === undefined"
                                        :placeholder="$t('components.rector.job_props.discord_guild_id')"
                                        :label="$t('components.rector.job_props.discord_guild_id')"
                                        maxlength="70"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                    <NuxtLink
                                        v-if="appConfig.discord.botInviteURL !== undefined"
                                        :to="appConfig.discord.botInviteURL"
                                        :external="true"
                                        class="mt-2 flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300 bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500"
                                    >
                                        {{ $t('components.rector.job_props.invite_bot') }}
                                    </NuxtLink>
                                    <p v-if="jobProps.discordLastSync" class="mt-2 text-base text-xs">
                                        {{ $t('components.rector.job_props.last_sync') }}:
                                        <Time :value="jobProps.discordLastSync" />
                                    </p>
                                </dd>
                            </div>
                            <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium">
                                    {{ $t('components.rector.job_props.discord_sync_settings') }}
                                </dt>
                                <dd v-if="jobProps.discordSyncSettings" class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                    <SwitchGroup as="div" class="flex items-center">
                                        <Switch
                                            v-model="jobProps.discordSyncSettings.userInfoSync"
                                            :class="[
                                                jobProps.discordSyncSettings.userInfoSync ? 'bg-indigo-600' : 'bg-gray-200',
                                                'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2',
                                            ]"
                                        >
                                            <span
                                                aria-hidden="true"
                                                :class="[
                                                    jobProps.discordSyncSettings.userInfoSync
                                                        ? 'translate-x-5'
                                                        : 'translate-x-0',
                                                    'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                ]"
                                            />
                                        </Switch>
                                        <SwitchLabel as="span" class="ml-3 text-sm">
                                            <span class="font-medium text-gray-300">{{
                                                $t('components.rector.job_props.user_info_sync')
                                            }}</span>
                                        </SwitchLabel>
                                    </SwitchGroup>
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
                                        class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                        :class="[
                                            !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-success-600 hover:bg-success-400 focus-visible:outline-success-500',
                                        ]"
                                        :disabled="!canSubmit"
                                        @click="onSubmitThrottle"
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="animate-spin h-5 w-5 mr-2" />
                                        </template>
                                        {{ $t('common.save', 1) }}
                                    </button>
                                </dd>
                            </div>
                        </dl>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
