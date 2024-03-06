<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import {
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions,
    Switch,
    SwitchGroup,
    SwitchLabel,
} from '@headlessui/vue';
import { useThrottleFn } from '@vueuse/core';
import { CheckIcon, LoadingIcon, OfficeBuildingCogIcon } from 'mdi-vue3';
import { useSettingsStore } from '~/store/settings';
import GenericContainerPanel from '~/components/partials/elements/GenericContainerPanel.vue';
import GenericContainerPanelEntry from '~/components/partials/elements/GenericContainerPanelEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { type GetAppConfigResponse } from '~~/gen/ts/services/rector/config';
import { useNotificatorStore } from '~/store/notificator';
import { useCompletorStore } from '~/store/completor';

const { $grpc } = useNuxtApp();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const notifications = useNotificatorStore();

const { data, pending, refresh, error } = useLazyAsyncData(`rector-jobprops`, () => getAppConfig());

async function getAppConfig(): Promise<GetAppConfigResponse> {
    try {
        const call = $grpc.getRectorConfigClient().getAppConfig({});
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function updateAppConfig(): Promise<void> {
    if (!data.value?.config) {
        return;
    }

    try {
        const { response } = await $grpc.getRectorConfigClient().updateAppConfig({
            config: data.value.config,
        });

        notifications.dispatchNotification({
            title: { key: 'notifications.rector.app_config.title', parameters: {} },
            content: { key: 'notifications.rector.app_config.content', parameters: {} },
            type: 'success',
        });

        if (response.config) {
            data.value.config = response.config;
        } else {
            refresh();
        }
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const completorStore = useCompletorStore();
const { listJobs } = completorStore;

const { data: jobs } = useLazyAsyncData(`rector-appconfig-jobs`, () => listJobs());

const queryJobsRaw = ref('');
const queryJobs = computed(() => queryJobsRaw.value.trim());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (_) => {
    canSubmit.value = false;
    await updateAppConfig().finally(() => setTimeout(() => (canSubmit.value = true), 400));
}, 1000);

watch(data, () => console.log('data', data.value?.config));
</script>

<template>
    <div class="mx-auto max-w-5xl py-2">
        <template v-if="streamerMode">
            <GenericContainerPanel>
                <template #title>
                    {{ $t('system.streamer_mode.title') }}
                </template>
                <template #description>
                    {{ $t('system.streamer_mode.description') }}
                </template>
            </GenericContainerPanel>
        </template>
        <template v-else>
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.setting', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.setting', 2)])"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="data === null" :icon="OfficeBuildingCogIcon" :type="$t('common.setting', 2)" />

            <template v-else>
                <GenericContainerPanel>
                    <template #title> Auth </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>Sign-up Enabled</template>
                            <template #default>
                                <SwitchGroup as="div" class="flex items-center">
                                    <Switch
                                        v-model="data.config!.auth!.signupEnabled"
                                        :class="[
                                            data.config!.auth!.signupEnabled ? 'bg-primary-600' : 'bg-gray-200',
                                            'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-600 focus:ring-offset-2',
                                        ]"
                                    >
                                        <span
                                            aria-hidden="true"
                                            :class="[
                                                data.config!.auth!.signupEnabled ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </Switch>
                                    <SwitchLabel as="span" class="ml-3 text-sm">
                                        <span class="font-medium text-gray-300">
                                            <template v-if="data.config!.auth!.signupEnabled">
                                                {{ $t('common.enabled') }}
                                            </template>
                                            <template v-else>
                                                {{ $t('common.disabled') }}
                                            </template>
                                        </span>
                                    </SwitchLabel>
                                </SwitchGroup>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel>
                    <template #title> Permissions </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>Default Permissions</template>
                            <template #default> TODO </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel>
                    <template #title> Website </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>Links</template>
                            <template #default>
                                <div class="flex-1 form-control">
                                    <label for="privacyPolicy">
                                        {{ $t('common.privacy_policy') }}
                                    </label>
                                    <input
                                        v-model="data.config!.website!.links!.privacyPolicy"
                                        type="text"
                                        name="privacyPolicy"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.privacy_policy')"
                                        :label="$t('common.privacy_policy')"
                                        maxlength="128"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                                <div class="flex-1 form-control">
                                    <label for="imprint">
                                        {{ $t('common.imprint') }}
                                    </label>
                                    <input
                                        v-model="data.config!.website!.links!.imprint"
                                        type="text"
                                        name="imprint"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.imprint')"
                                        :label="$t('common.imprint')"
                                        maxlength="128"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel>
                    <template #title> Job Info </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>Unemployed Job</template>
                            <template #default>
                                <div class="flex-1 form-control">
                                    <label for="unemployedJobName"> {{ $t('common.job') }} {{ $t('common.name') }} </label>
                                    <input
                                        v-model="data.config!.jobInfo!.unemployedJob!.name"
                                        type="text"
                                        name="unemployedJobName"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.privacy_policy')"
                                        :label="$t('common.privacy_policy')"
                                        maxlength="128"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                                <div class="flex-1 form-control">
                                    <label for="unemployedJobGrade">
                                        {{ $t('common.rank') }}
                                    </label>
                                    <input
                                        v-model="data.config!.jobInfo!.unemployedJob!.grade"
                                        type="number"
                                        min="1"
                                        max="99999"
                                        name="unemployedJobGrade"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.rank')"
                                        :label="$t('common.rank')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>Public Jobs</template>
                            <template #default>
                                <Combobox v-model="data.config!.jobInfo!.publicJobs" as="div" multiple nullable>
                                    <div class="relative">
                                        <ComboboxButton as="div">
                                            <ComboboxInput
                                                autocomplete="off"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                :display-value="(js: any) => (js ? js.join(', ') : $t('common.na'))"
                                                :placeholder="$t('common.job', 2)"
                                                @change="queryJobsRaw = $event.target.value"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </ComboboxButton>

                                        <ComboboxOptions
                                            v-if="jobs !== null && jobs.length > 0"
                                            class="absolute z-30 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                        >
                                            <ComboboxOption
                                                v-for="job in jobs.filter(
                                                    (j) => j.label.includes(queryJobs) || j.name.includes(queryJobs),
                                                )"
                                                v-slot="{ active, selected }"
                                                :key="job.name"
                                                :value="job.name"
                                                as="template"
                                            >
                                                <li
                                                    :class="[
                                                        'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                        active ? 'bg-primary-500' : '',
                                                    ]"
                                                >
                                                    <span :class="['block truncate', selected && 'font-semibold']">
                                                        {{ job.label }}
                                                    </span>

                                                    <span
                                                        v-if="selected"
                                                        :class="[
                                                            active ? 'text-neutral' : 'text-primary-500',
                                                            'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                        ]"
                                                    >
                                                        <CheckIcon class="h-5 w-5" aria-hidden="true" />
                                                    </span>
                                                </li>
                                            </ComboboxOption>
                                        </ComboboxOptions>
                                    </div>
                                </Combobox>
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>Hidden Jobs</template>
                            <template #default>
                                <Combobox v-model="data.config!.jobInfo!.hiddenJobs" as="div" multiple nullable>
                                    <div class="relative">
                                        <ComboboxButton as="div">
                                            <ComboboxInput
                                                autocomplete="off"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                :display-value="(js: any) => (js ? js.join(', ') : $t('common.na'))"
                                                :placeholder="$t('common.job', 2)"
                                                @change="queryJobsRaw = $event.target.value"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </ComboboxButton>

                                        <ComboboxOptions
                                            v-if="jobs !== null && jobs.length > 0"
                                            class="absolute z-30 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                        >
                                            <ComboboxOption
                                                v-for="job in jobs.filter(
                                                    (j) => j.label.includes(queryJobs) || j.name.includes(queryJobs),
                                                )"
                                                v-slot="{ active, selected }"
                                                :key="job.name"
                                                :value="job.name"
                                                as="template"
                                            >
                                                <li
                                                    :class="[
                                                        'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                        active ? 'bg-primary-500' : '',
                                                    ]"
                                                >
                                                    <span :class="['block truncate', selected && 'font-semibold']">
                                                        {{ job.label }}
                                                    </span>

                                                    <span
                                                        v-if="selected"
                                                        :class="[
                                                            active ? 'text-neutral' : 'text-primary-500',
                                                            'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                        ]"
                                                    >
                                                        <CheckIcon class="h-5 w-5" aria-hidden="true" />
                                                    </span>
                                                </li>
                                            </ComboboxOption>
                                        </ComboboxOptions>
                                    </div>
                                </Combobox>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel>
                    <template #title> User Tracker / Livemap </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>Refresh Times</template>
                            <template #default>
                                {{
                                    parseFloat(
                                        data.config?.userTracker?.refreshTime?.seconds.toString() +
                                            '.' +
                                            (data.config?.userTracker?.refreshTime?.nanos ?? 0) / 1000000,
                                    ).toString()
                                }}s
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>DB Refresh Times</template>
                            <template #default>
                                {{
                                    parseFloat(
                                        data.config?.userTracker?.dbRefreshTime?.seconds.toString() +
                                            '.' +
                                            (data.config?.userTracker?.dbRefreshTime?.nanos ?? 0) / 1000000,
                                    ).toString()
                                }}s
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>Livemap Jobs</template>
                            <template #default>
                                <Combobox v-model="data.config!.userTracker!.livemapJobs" as="div" multiple nullable>
                                    <div class="relative">
                                        <ComboboxButton as="div">
                                            <ComboboxInput
                                                autocomplete="off"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                :display-value="(js: any) => (js ? js.join(', ') : $t('common.na'))"
                                                :placeholder="$t('common.job', 2)"
                                                @change="queryJobsRaw = $event.target.value"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </ComboboxButton>

                                        <ComboboxOptions
                                            v-if="jobs !== null && jobs.length > 0"
                                            class="absolute z-30 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                        >
                                            <ComboboxOption
                                                v-for="job in jobs.filter(
                                                    (j) => j.label.includes(queryJobs) || j.name.includes(queryJobs),
                                                )"
                                                v-slot="{ active, selected }"
                                                :key="job.name"
                                                :value="job.name"
                                                as="template"
                                            >
                                                <li
                                                    :class="[
                                                        'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                        active ? 'bg-primary-500' : '',
                                                    ]"
                                                >
                                                    <span :class="['block truncate', selected && 'font-semibold']">
                                                        {{ job.label }}
                                                    </span>

                                                    <span
                                                        v-if="selected"
                                                        :class="[
                                                            active ? 'text-neutral' : 'text-primary-500',
                                                            'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                        ]"
                                                    >
                                                        <CheckIcon class="h-5 w-5" aria-hidden="true" />
                                                    </span>
                                                </li>
                                            </ComboboxOption>
                                        </ComboboxOptions>
                                    </div>
                                </Combobox>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel>
                    <template #title> Discord </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>Enabled</template>
                            <template #default>
                                <SwitchGroup as="div" class="flex items-center">
                                    <Switch
                                        v-model="data.config!.discord!.enabled"
                                        :class="[
                                            data.config!.discord!.enabled ? 'bg-primary-600' : 'bg-gray-200',
                                            'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-600 focus:ring-offset-2',
                                        ]"
                                    >
                                        <span
                                            aria-hidden="true"
                                            :class="[
                                                data.config!.discord!.enabled ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </Switch>
                                    <SwitchLabel as="span" class="ml-3 text-sm">
                                        <span class="font-medium text-gray-300">
                                            <template v-if="data.config!.discord!.enabled">
                                                {{ $t('common.enabled') }}
                                            </template>
                                            <template v-else>
                                                {{ $t('common.disabled') }}
                                            </template>
                                        </span>
                                    </SwitchLabel>
                                </SwitchGroup>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>

                <GenericContainerPanel>
                    <template #title>{{ $t('common.save', 1) }}</template>
                    <template #description>Make sure to double check any config options before saving the config.</template>
                    <template #default>
                        <!-- Save button -->
                        <GenericContainerPanelEntry v-if="can('RectorService.SetJobProps')">
                            <template #default>
                                <button
                                    type="button"
                                    class="flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-neutral transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                    :class="[
                                        !canSubmit
                                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                    ]"
                                    :disabled="!canSubmit"
                                    @click="onSubmitThrottle"
                                >
                                    <template v-if="!canSubmit && false">
                                        <LoadingIcon class="mr-2 h-5 w-5 animate-spin" aria-hidden="true" />
                                    </template>
                                    {{ $t('common.save', 1) }}
                                </button>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
            </template>
        </template>
    </div>
</template>
