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
// eslint-disable-next-line camelcase
import { max, min, regex, required, url, min_value, max_value, numeric } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { CheckIcon, CloseIcon, LoadingIcon, OfficeBuildingCogIcon, PlusIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useSettingsStore } from '~/store/settings';
import GenericContainerPanel from '~/components/partials/elements/GenericContainerPanel.vue';
import GenericContainerPanelEntry from '~/components/partials/elements/GenericContainerPanelEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { type GetAppConfigResponse } from '~~/gen/ts/services/rector/config';
import { useNotificatorStore } from '~/store/notificator';
import { useCompletorStore } from '~/store/completor';
import type { Perm } from '~~/gen/ts/resources/rector/config';
import { toDuration } from '~/utils/duration';

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

const completorStore = useCompletorStore();
const { listJobs } = completorStore;

const { data: jobs } = useLazyAsyncData(`rector-appconfig-jobs`, () => listJobs());

interface FormData {
    permsDefault: Perm[];

    websiteLinksPrivacyPolicy?: string;
    websiteLinksImprint?: string;

    jobInfoUnemployedName: string;
    jobInfoUnemployedGrade: number;

    userTrackerRefreshTime: string;
    userTrackerDbRefreshTime: string;

    discordBotInviteUrl?: string;
    discordSyncInterval: string;
}

async function updateAppConfig(values: FormData): Promise<void> {
    if (!data.value?.config) {
        return;
    }

    // Perms
    if (data.value.config.perms === undefined) {
        data.value.config.perms = {
            default: [],
        };
    }
    data.value.config.perms.default = values.permsDefault;

    // Website
    if (data.value.config.website === undefined) {
        data.value.config.website = {
            links: {},
        };
    }
    if (data.value.config.website.links === undefined) {
        data.value.config.website.links = {};
    }
    data.value.config.website.links.imprint = values.websiteLinksImprint;
    data.value.config.website.links.privacyPolicy = values.websiteLinksPrivacyPolicy;

    // Job Info
    if (data.value.config.jobInfo === undefined) {
        data.value.config.jobInfo = {
            hiddenJobs: [],
            publicJobs: [],
        };
    }
    if (data.value.config.jobInfo.unemployedJob === undefined) {
        data.value.config.jobInfo.unemployedJob = {
            name: '',
            grade: 0,
        };
    }
    data.value.config.jobInfo.unemployedJob.name = values.jobInfoUnemployedName;
    data.value.config.jobInfo.unemployedJob.grade = values.jobInfoUnemployedGrade;

    // User Tracker
    if (data.value.config.userTracker === undefined) {
        data.value.config.userTracker = {
            livemapJobs: [],
        };
    }
    if (values.userTrackerRefreshTime) {
        data.value.config.userTracker.refreshTime = toDuration(values.userTrackerRefreshTime);
    }
    if (values.userTrackerDbRefreshTime) {
        data.value.config.userTracker.dbRefreshTime = toDuration(values.userTrackerDbRefreshTime);
    }

    // Discord
    if (data.value.config.discord === undefined) {
        data.value.config.discord = {
            enabled: false,
        };
    }
    data.value.config.discord.inviteUrl = values.discordBotInviteUrl;
    if (data.value.config.discord.syncInterval === undefined && values.discordSyncInterval) {
        data.value.config.discord.syncInterval = toDuration(values.discordSyncInterval);
    }
    console.log(data.value);

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

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);
defineRule('regex', regex);
defineRule('url', url);
defineRule('min_value', min_value);
defineRule('max_value', max_value);
defineRule('numeric', numeric);

const { handleSubmit, meta, setValues } = useForm<FormData>({
    validationSchema: {
        permsDefault: { max: 25 },

        websitePrivacyPolicy: { required: false, max: 255, url: 'https://.*' },
        websiteLinksImprint: { required: false, max: 255, url: 'https://.*' },

        jobInfoUnemployedName: { required: true, max: 20 },
        jobInfoUnemployedGrade: { required: true, numeric: true, min_value: 1, max_value: 99 },

        userTrackerRefreshTime: { required: true, max: 5, regex: /^\d+(\.\d+)?s$/ },
        userTrackerDbRefreshTime: { required: true, max: 5, regex: /^\d+(\.\d+)?s$/ },

        discordBotInviteUrl: { required: false, url: 'https://discord.com/.*' },
        discordSyncInterval: { required: true, max: 5, regex: /^\d+(\.\d+)?s$/ },
    },
    validateOnMount: true,
});

function setSettingsValues(): void {
    if (!data.value) {
        return;
    }

    setValues({
        permsDefault: data.value.config?.perms?.default,

        websiteLinksPrivacyPolicy: data.value.config?.website?.links?.privacyPolicy,
        websiteLinksImprint: data.value.config?.website?.links?.imprint,

        jobInfoUnemployedName: data.value.config?.jobInfo?.unemployedJob?.name,
        jobInfoUnemployedGrade: data.value.config?.jobInfo?.unemployedJob?.grade,

        userTrackerRefreshTime:
            parseFloat(
                data.value.config?.userTracker?.refreshTime?.seconds.toString() +
                    '.' +
                    (data.value.config?.userTracker?.refreshTime?.nanos ?? 0) / 1000000,
            ).toString() + 's',
        userTrackerDbRefreshTime:
            parseFloat(
                data.value.config?.userTracker?.dbRefreshTime?.seconds.toString() +
                    '.' +
                    (data.value.config?.userTracker?.dbRefreshTime?.nanos ?? 0) / 1000000,
            ).toString() + 's',

        discordBotInviteUrl: data.value.config?.discord?.inviteUrl,
    });
}

watch(data, () => setSettingsValues());

setSettingsValues();

const queryJobsRaw = ref('');
const queryJobs = computed(() => queryJobsRaw.value.trim());

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await updateAppConfig(values).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const { remove, push, fields } = useFieldArray<Perm>('permsDefault');
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
                    <template #title>{{ $t('components.rector.app_config.auth.title') }}</template>
                    <template #description>{{ $t('components.rector.app_config.auth.description') }}</template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>{{ $t('components.rector.app_config.auth.sign_up') }}</template>
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
                    <template #title>{{ $t('components.rector.app_config.perms.title') }}</template>
                    <template #description>{{ $t('components.rector.app_config.perms.description') }}</template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>{{ $t('components.rector.app_config.perms.default_perms') }}</template>
                            <template #default>
                                <div class="flex flex-col gap-1">
                                    <div v-for="(field, idx) in fields" :key="field.key" class="flex gap-1 items-center">
                                        <div class="flex-1">
                                            <VeeField
                                                :name="`permsDefault[${idx}].category`"
                                                type="text"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                :placeholder="$t('common.category')"
                                                :label="$t('common.category')"
                                                :rules="required"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                            <VeeErrorMessage
                                                :name="`permsDefault[${idx}].category`"
                                                as="p"
                                                class="mt-2 text-sm text-error-400"
                                            />
                                        </div>
                                        <div class="flex-1">
                                            <VeeField
                                                :name="`permsDefault[${idx}].name`"
                                                type="text"
                                                class="flex-1 block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                :placeholder="$t('common.name')"
                                                :label="$t('common.name')"
                                                :rules="required"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                            <VeeErrorMessage
                                                :name="`permsDefault[${idx}].name`"
                                                as="p"
                                                class="mt-2 text-sm text-error-400"
                                            />
                                        </div>

                                        <button
                                            type="button"
                                            class="rounded-full bg-primary-500 p-1.5 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                            @click="remove(idx)"
                                        >
                                            <CloseIcon class="h-5 w-5" aria-hidden="true" />
                                        </button>
                                    </div>
                                </div>
                                <button
                                    type="button"
                                    class="mt-2 rounded-full p-1.5 text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                    :disabled="!canSubmit"
                                    :class="
                                        !canSubmit
                                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500'
                                    "
                                    @click="push({ category: '', name: '' })"
                                >
                                    <PlusIcon class="h-5 w-5" aria-hidden="true" />
                                </button>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel>
                    <template #title>{{ $t('components.rector.app_config.website.title') }}</template>
                    <template #description>{{ $t('components.rector.app_config.website.description') }}</template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>{{ $t('components.rector.app_config.website.links.title') }}</template>
                            <template #default>
                                <div class="flex-1 form-control">
                                    <label for="websiteLinksPrivacyPolicy">
                                        {{ $t('common.privacy_policy') }}
                                    </label>
                                    <VeeField
                                        type="text"
                                        name="websiteLinksPrivacyPolicy"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :value="data.config!.website!.links!.privacyPolicy"
                                        :placeholder="$t('common.privacy_policy')"
                                        :label="$t('common.privacy_policy')"
                                        maxlength="128"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                    <VeeErrorMessage
                                        name="websiteLinksPrivacyPolicy"
                                        as="p"
                                        class="mt-2 text-sm text-error-400"
                                    />
                                </div>
                                <div class="flex-1 form-control">
                                    <label for="websiteLinksImprint">
                                        {{ $t('common.imprint') }}
                                    </label>
                                    <VeeField
                                        type="text"
                                        name="websiteLinksImprint"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.imprint')"
                                        :label="$t('common.imprint')"
                                        :value="data.config!.website!.links!.imprint"
                                        maxlength="128"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                    <VeeErrorMessage name="websiteLinksImprint" as="p" class="mt-2 text-sm text-error-400" />
                                </div>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel>
                    <template #title>{{ $t('components.rector.app_config.job_info.title') }}</template>
                    <template #description>{{ $t('components.rector.app_config.job_info.description') }}</template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>{{ $t('components.rector.app_config.job_info.unemployed_job') }}</template>
                            <template #default>
                                <div class="flex-1 form-control">
                                    <label for="jobInfoUnemployedName"> {{ $t('common.job') }} {{ $t('common.name') }} </label>
                                    <VeeField
                                        type="text"
                                        name="jobInfoUnemployedName"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :value="data.config!.jobInfo!.unemployedJob!.name"
                                        :placeholder="$t('common.job')"
                                        :label="$t('common.job')"
                                        maxlength="128"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                    <VeeErrorMessage name="jobInfoUnemployedName" as="p" class="mt-2 text-sm text-error-400" />
                                </div>
                                <div class="flex-1 form-control">
                                    <label for="jobInfoUnemployedGrade">
                                        {{ $t('common.rank') }}
                                    </label>
                                    <VeeField
                                        type="number"
                                        min="1"
                                        max="99"
                                        :value="data.config!.jobInfo!.unemployedJob!.grade"
                                        name="jobInfoUnemployedGrade"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.rank')"
                                        :label="$t('common.rank')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                    <VeeErrorMessage name="jobInfoUnemployedGrade" as="p" class="mt-2 text-sm text-error-400" />
                                </div>
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>{{ $t('components.rector.app_config.job_info.public_jobs') }}</template>
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
                                                        {{ job.name }}
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
                            <template #title>{{ $t('components.rector.app_config.job_info.hidden_jobs') }}</template>
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
                                                        {{ job.name }}
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
                    <template #title>{{ $t('components.rector.app_config.user_tracker.title') }}</template>
                    <template #description>{{ $t('components.rector.app_config.user_tracker.description') }}</template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>{{ $t('components.rector.app_config.user_tracker.refresh_time') }}</template>
                            <template #default>
                                <VeeField
                                    name="userTrackerRefreshTime"
                                    type="text"
                                    class="flex-1 block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    :value="
                                        parseFloat(
                                            data.config?.userTracker?.refreshTime?.seconds.toString() +
                                                '.' +
                                                (data.config?.userTracker?.refreshTime?.nanos ?? 0) / 1000000,
                                        ).toString() + 's'
                                    "
                                    :placeholder="$t('common.duration')"
                                    :label="$t('common.duration')"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                                <VeeErrorMessage name="userTrackerRefreshTime" as="p" class="mt-2 text-sm text-error-400" />
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>{{ $t('components.rector.app_config.user_tracker.db_refresh_time') }}</template>
                            <template #default>
                                <VeeField
                                    name="userTrackerDbRefreshTime"
                                    type="text"
                                    class="flex-1 block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    :placeholder="$t('common.duration')"
                                    :label="$t('common.duration')"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                                <VeeErrorMessage name="userTrackerRefreshTime" as="p" class="mt-2 text-sm text-error-400" />
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>{{ $t('components.rector.app_config.user_tracker.livemap_jobs') }}</template>
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
                                                        {{ job.name }}
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
                    <template #title>{{ $t('common.discord') }}</template>
                    <template #description>{{ $t('components.rector.app_config.discord.description') }}</template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>{{ $t('common.enabled') }}</template>
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
                        <GenericContainerPanelEntry>
                            <template #title>{{ $t('components.rector.app_config.discord.sync_interval') }}</template>
                            <template #default>
                                <VeeField
                                    name="discordSyncInterval"
                                    type="text"
                                    class="flex-1 block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    :value="
                                        parseFloat(
                                            data.config?.discord?.syncInterval?.seconds.toString() +
                                                '.' +
                                                (data.config?.discord?.syncInterval?.nanos ?? 0) / 1000000,
                                        ).toString() + 's'
                                    "
                                    :placeholder="$t('common.duration')"
                                    :label="$t('common.duration')"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                                <VeeErrorMessage name="discordSyncInterval" as="p" class="mt-2 text-sm text-error-400" />
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>{{ $t('components.rector.app_config.discord.bot_invite_url') }}</template>
                            <template #default>
                                <div class="flex-1 form-control">
                                    <VeeField
                                        type="url"
                                        name="discordBotInviteUrl"
                                        :value="data.config!.discord!.inviteUrl"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('components.rector.app_config.discord.bot_invite_url')"
                                        :label="$t('components.rector.app_config.discord.bot_invite_url')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                    <VeeErrorMessage name="discordBotInviteUrl" as="p" class="mt-2 text-sm text-error-400" />
                                </div>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>

                <GenericContainerPanel>
                    <template #title>{{ $t('common.save', 1) }}</template>
                    <template #description>{{ $t('components.rector.app_config.save.description') }}</template>
                    <template #default>
                        <!-- Save button -->
                        <GenericContainerPanelEntry v-if="can('RectorService.SetJobProps')">
                            <template #default>
                                <button
                                    type="button"
                                    class="flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-neutral transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                    :class="[
                                        !canSubmit || !meta.valid
                                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                    ]"
                                    :disabled="!canSubmit || !meta.valid"
                                    @click="onSubmitThrottle"
                                >
                                    <template v-if="!canSubmit">
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
