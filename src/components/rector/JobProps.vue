<script lang="ts" setup>
import { Switch, SwitchGroup, SwitchLabel, Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useThrottleFn } from '@vueuse/core';
import { vMaska } from 'maska';
import { CheckIcon, ChevronDownIcon, LoadingIcon, TuneIcon } from 'mdi-vue3';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { availableThemes } from '~/store/settings';
import { JobProps, UserInfoSyncUnemployedMode } from '~~/gen/ts/resources/users/jobs';
import GenericContainerPanel from '~/components/partials/GenericContainerPanel.vue';
import GenericContainerPanelEntry from '~/components/partials/GenericContainerPanelEntry.vue';
import SquareImg from '~/components/partials/SquareImg.vue';
import { File } from '~~/gen/ts/resources/filestore/file';

const { $grpc } = useNuxtApp();

const appConfig = useAppConfig();

const authStore = useAuthStore();

const notifications = useNotificatorStore();

async function getJobProps(): Promise<JobProps> {
    try {
        const call = $grpc.getRectorClient().getJobProps({});
        const { response } = await call;

        return response.jobProps!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { data: jobProps, pending, refresh, error } = useLazyAsyncData(`rector-jobprops`, () => getJobProps());

async function setJobProps(): Promise<void> {
    if (!jobProps.value) {
        return;
    }

    try {
        const { response } = await $grpc.getRectorClient().setJobProps({
            jobProps: jobProps.value,
        });

        notifications.dispatchNotification({
            title: { key: 'notifications.rector.job_props.title', parameters: {} },
            content: { key: 'notifications.rector.job_props.content', parameters: {} },
            type: 'success',
        });

        if (response.jobProps) {
            jobProps.value = response.jobProps;
            authStore.setJobProps(jobProps.value);
        }
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (_) => {
    canSubmit.value = false;
    await setJobProps().finally(() => setTimeout(() => (canSubmit.value = true), 400));
}, 1000);

const fileUploadRef = ref<HTMLInputElement | null>(null);
async function loadImage(): Promise<void> {
    if (!jobProps.value || !fileUploadRef.value || !fileUploadRef.value.files || fileUploadRef.value.files.length <= 0) {
        return;
    }

    if (fileUploadRef.value.files[0].size > 2097152) {
        return;
    }

    if (jobProps.value.logoUrl === undefined) {
        jobProps.value.logoUrl = {} as File;
    }

    jobProps.value.logoUrl.data = new Uint8Array(await fileUploadRef.value.files[0].arrayBuffer());
}
</script>

<template>
    <div class="p-2">
        <div class="mx-auto max-w-5xl">
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [`${$t('common.job', 1)} ${$t('common.prop')}`])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [`${$t('common.job', 1)} ${$t('common.prop')}`])"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="jobProps === null"
                :icon="TuneIcon"
                :type="`${$t('common.job', 1)} ${$t('common.prop')}`"
            />

            <template v-else>
                <GenericContainerPanel>
                    <template #title>
                        {{ $t('components.rector.job_props.job_properties') }}
                    </template>
                    <template #description>
                        {{ $t('components.rector.job_props.your_job_properties') }}
                    </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('common.theme') }}
                            </template>
                            <template #default>
                                <Listbox v-model="jobProps.theme" as="div">
                                    <div class="relative">
                                        <ListboxButton
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        >
                                            <span class="block truncate">
                                                {{ availableThemes.find((t) => t.key === jobProps?.theme)?.name }}
                                            </span>
                                            <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                                                <ChevronDownIcon class="h-5 w-5 text-gray-400" aria-hidden="true" />
                                            </span>
                                        </ListboxButton>

                                        <transition
                                            leave-active-class="transition duration-100 ease-in"
                                            leave-from-class="opacity-100"
                                            leave-to-class="opacity-0"
                                        >
                                            <ListboxOptions
                                                class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                            >
                                                <ListboxOption
                                                    v-for="theme in availableThemes"
                                                    :key="theme.key"
                                                    v-slot="{ active, selected }"
                                                    as="template"
                                                    :value="theme.key"
                                                >
                                                    <li
                                                        :class="[
                                                            active ? 'bg-primary-500' : '',
                                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                        ]"
                                                    >
                                                        <span
                                                            :class="[
                                                                selected ? 'font-semibold' : 'font-normal',
                                                                'block truncate',
                                                            ]"
                                                        >
                                                            {{ theme.name }}
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
                                                </ListboxOption>
                                            </ListboxOptions>
                                        </transition>
                                    </div>
                                </Listbox>
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('components.rector.job_props.livemap_marker_color') }}
                            </template>
                            <template #default>
                                <ColorInput v-model="jobProps.livemapMarkerColor" disable-alpha format="hex" position="top" />
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('common.radio_frequency') }}
                            </template>
                            <template #default>
                                <input
                                    v-model="jobProps.radioFrequency"
                                    v-maska
                                    data-maska="0.9"
                                    data-maska-tokens="0:\d:multiple|9:\d:multiple"
                                    type="text"
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    :placeholder="$t('common.radio_frequency')"
                                    :label="$t('common.radio_frequency')"
                                    maxlength="24"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry v-if="jobProps.quickButtons">
                            <template #title>
                                {{ $t('components.rector.job_props.quick_buttons') }}
                            </template>
                            <template #default>
                                <fieldset class="flex flex-col gap-2">
                                    <div class="space-y-4">
                                        <SwitchGroup as="div" class="flex items-center">
                                            <Switch
                                                v-model="jobProps.quickButtons.penaltyCalculator"
                                                :class="[
                                                    jobProps.quickButtons.penaltyCalculator ? 'bg-indigo-600' : 'bg-gray-200',
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
                                                        jobProps.quickButtons.bodyCheckup ? 'translate-x-5' : 'translate-x-0',
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
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('common.logo') }}
                            </template>
                            <template #default>
                                <div class="flex flex-col gap-3">
                                    <input
                                        ref="fileUploadRef"
                                        type="file"
                                        accept="image/x-png,image/jpeg"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @change="loadImage()"
                                    />

                                    <SquareImg v-if="jobProps.logoUrl?.url" size="xl" :url="jobProps.logoUrl.url" />
                                </div>
                            </template>
                        </GenericContainerPanelEntry>

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
                                    <template v-if="!canSubmit">
                                        <LoadingIcon class="mr-2 h-5 w-5 animate-spin" />
                                    </template>
                                    {{ $t('common.save', 1) }}
                                </button>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
                <GenericContainerPanel v-if="jobProps.discordSyncSettings">
                    <template #title>
                        {{ $t('components.rector.job_props.discord_sync_settings.title') }}
                    </template>
                    <template #description>
                        {{ $t('components.rector.job_props.discord_sync_settings.subtitle') }}
                    </template>
                    <template #default>
                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('components.rector.job_props.discord_guild_id') }}
                            </template>
                            <template #default>
                                <input
                                    v-model="jobProps.discordGuildId"
                                    type="text"
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                    class="mt-2 flex w-full justify-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300 focus-visible:outline-primary-500"
                                >
                                    {{ $t('components.rector.job_props.invite_bot') }}
                                </NuxtLink>
                                <p v-if="jobProps.discordLastSync" class="mt-2 text-base text-xs">
                                    {{ $t('components.rector.job_props.last_sync') }}:
                                    <GenericTime :value="jobProps.discordLastSync" />
                                </p>
                            </template>
                        </GenericContainerPanelEntry>

                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('components.rector.job_props.status_log') }}
                            </template>
                            <template #default>
                                <SwitchGroup as="div" class="flex items-center mb-1">
                                    <Switch
                                        v-model="jobProps.discordSyncSettings.statusLog"
                                        :class="[
                                            jobProps.discordSyncSettings.statusLog ? 'bg-indigo-600' : 'bg-gray-200',
                                            'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2',
                                        ]"
                                    >
                                        <span
                                            aria-hidden="true"
                                            :class="[
                                                jobProps.discordSyncSettings.statusLog ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </Switch>
                                    <SwitchLabel as="span" class="ml-3 text-sm">
                                        <span class="font-medium text-gray-300">{{ $t('common.enabled') }}</span>
                                    </SwitchLabel>
                                </SwitchGroup>

                                <template v-if="jobProps.discordSyncSettings.statusLog">
                                    <label for="statusLogSettingsChannelId">
                                        {{ $t('components.rector.job_props.status_log_settings.channel_id') }}:
                                    </label>
                                    <input
                                        v-model="jobProps.discordSyncSettings.statusLogSettings!.channelId"
                                        type="text"
                                        name="statusLogSettingsChannelId"
                                        :disabled="!jobProps.discordSyncSettings.statusLog"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :class="!jobProps.discordSyncSettings.statusLog ? 'disabled' : ''"
                                        :placeholder="$t('components.rector.job_props.status_log_settings.channel_id')"
                                        :label="$t('components.rector.job_props.status_log_settings.channel_id')"
                                        maxlength="48"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </template>
                            </template>
                        </GenericContainerPanelEntry>

                        <!-- User Info Sync Settings -->
                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('components.rector.job_props.user_info_sync') }}
                            </template>
                            <template #default>
                                <SwitchGroup as="div" class="flex items-center mb-1">
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
                                                jobProps.discordSyncSettings.userInfoSync ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </Switch>
                                    <SwitchLabel as="span" class="ml-3 text-sm">
                                        <span class="font-medium text-gray-300">{{ $t('common.enabled') }}</span>
                                    </SwitchLabel>
                                </SwitchGroup>

                                <template v-if="jobProps.discordSyncSettings.userInfoSync">
                                    <label for="gradeRoleFormat">
                                        {{ $t('components.rector.job_props.user_info_sync_settings.grade_role_format') }}:
                                    </label>
                                    <input
                                        v-model="jobProps.discordSyncSettings.userInfoSyncSettings!.gradeRoleFormat"
                                        type="text"
                                        name="gradeRoleFormat"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="
                                            $t('components.rector.job_props.user_info_sync_settings.grade_role_format')
                                        "
                                        :label="$t('components.rector.job_props.user_info_sync_settings.grade_role_format')"
                                        maxlength="48"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />

                                    <!-- UserInfo Sync Settings -->
                                    <SwitchGroup
                                        v-if="jobProps.discordSyncSettings.userInfoSyncSettings !== undefined"
                                        as="div"
                                        class="flex items-center mt-2 mb-1"
                                    >
                                        <Switch
                                            v-model="jobProps.discordSyncSettings.userInfoSyncSettings.employeeRoleEnabled"
                                            :class="[
                                                jobProps.discordSyncSettings.userInfoSyncSettings.employeeRoleEnabled
                                                    ? 'bg-indigo-600'
                                                    : 'bg-gray-200',
                                                'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2',
                                            ]"
                                        >
                                            <span
                                                aria-hidden="true"
                                                :class="[
                                                    jobProps.discordSyncSettings.userInfoSyncSettings.employeeRoleEnabled
                                                        ? 'translate-x-5'
                                                        : 'translate-x-0',
                                                    'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                ]"
                                            />
                                        </Switch>
                                        <SwitchLabel as="span" class="ml-3 text-sm">
                                            <span class="font-medium text-gray-300">{{
                                                $t('components.rector.job_props.user_info_sync_settings.employee_role_enabled')
                                            }}</span>
                                        </SwitchLabel>
                                    </SwitchGroup>

                                    <div v-if="jobProps.discordSyncSettings.userInfoSyncSettings?.employeeRoleEnabled">
                                        <label for="employeeRoleFormat">
                                            {{
                                                $t('components.rector.job_props.user_info_sync_settings.employee_role_format')
                                            }}:
                                        </label>
                                        <input
                                            v-model="jobProps.discordSyncSettings.userInfoSyncSettings!.employeeRoleFormat"
                                            type="text"
                                            name="employeeRoleFormat"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="
                                                $t('components.rector.job_props.user_info_sync_settings.employee_role_format')
                                            "
                                            :label="
                                                $t('components.rector.job_props.user_info_sync_settings.employee_role_format')
                                            "
                                            maxlength="48"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </div>

                                    <template v-if="jobProps.discordSyncSettings.userInfoSyncSettings !== undefined">
                                        <SwitchGroup as="div" class="flex items-center mt-2 mb-1">
                                            <Switch
                                                v-model="jobProps.discordSyncSettings.userInfoSyncSettings.unemployedEnabled"
                                                :class="[
                                                    jobProps.discordSyncSettings.userInfoSyncSettings.unemployedEnabled
                                                        ? 'bg-indigo-600'
                                                        : 'bg-gray-200',
                                                    'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2',
                                                ]"
                                            >
                                                <span
                                                    aria-hidden="true"
                                                    :class="[
                                                        jobProps.discordSyncSettings.userInfoSyncSettings.unemployedEnabled
                                                            ? 'translate-x-5'
                                                            : 'translate-x-0',
                                                        'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                    ]"
                                                />
                                            </Switch>
                                            <SwitchLabel as="span" class="ml-3 text-sm">
                                                <span class="font-medium text-gray-300">{{
                                                    $t('components.rector.job_props.user_info_sync_settings.unemployed_enabled')
                                                }}</span>
                                            </SwitchLabel>
                                        </SwitchGroup>
                                        <template v-if="jobProps.discordSyncSettings.userInfoSyncSettings.unemployedEnabled">
                                            <div>
                                                <label for="unemployedMode">
                                                    {{
                                                        $t(
                                                            'components.rector.job_props.user_info_sync_settings.unemployed_mode',
                                                        )
                                                    }}:
                                                </label>

                                                <Listbox
                                                    v-model="jobProps.discordSyncSettings.userInfoSyncSettings!.unemployedMode"
                                                    as="div"
                                                >
                                                    <div class="relative">
                                                        <ListboxButton
                                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        >
                                                            <span class="block truncate">
                                                                {{
                                                                    UserInfoSyncUnemployedMode[
                                                                        jobProps.discordSyncSettings.userInfoSyncSettings!
                                                                            .unemployedMode ?? 0
                                                                    ]
                                                                }}
                                                            </span>
                                                            <span
                                                                class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2"
                                                            >
                                                                <ChevronDownIcon
                                                                    class="h-5 w-5 text-gray-400"
                                                                    aria-hidden="true"
                                                                />
                                                            </span>
                                                        </ListboxButton>

                                                        <transition
                                                            leave-active-class="transition duration-100 ease-in"
                                                            leave-from-class="opacity-100"
                                                            leave-to-class="opacity-0"
                                                        >
                                                            <ListboxOptions
                                                                class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                                            >
                                                                <ListboxOption
                                                                    v-for="mode in [
                                                                        UserInfoSyncUnemployedMode.GIVE_ROLE,
                                                                        UserInfoSyncUnemployedMode.KICK,
                                                                    ]"
                                                                    :key="mode"
                                                                    v-slot="{ active, selected }"
                                                                    as="template"
                                                                    :value="mode"
                                                                >
                                                                    <li
                                                                        :class="[
                                                                            active ? 'bg-primary-500' : '',
                                                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                                        ]"
                                                                    >
                                                                        <span
                                                                            :class="[
                                                                                selected ? 'font-semibold' : 'font-normal',
                                                                                'block truncate',
                                                                            ]"
                                                                        >
                                                                            {{ UserInfoSyncUnemployedMode[mode] }}
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
                                                                </ListboxOption>
                                                            </ListboxOptions>
                                                        </transition>
                                                    </div>
                                                </Listbox>
                                            </div>

                                            <div>
                                                <label for="unemployedRoleName">
                                                    {{
                                                        $t(
                                                            'components.rector.job_props.user_info_sync_settings.unemployed_role_name',
                                                        )
                                                    }}:
                                                </label>
                                                <input
                                                    v-model="
                                                        jobProps.discordSyncSettings.userInfoSyncSettings!.unemployedRoleName
                                                    "
                                                    type="text"
                                                    name="unemployedRoleName"
                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    :placeholder="
                                                        $t(
                                                            'components.rector.job_props.user_info_sync_settings.unemployed_role_name',
                                                        )
                                                    "
                                                    :label="
                                                        $t(
                                                            'components.rector.job_props.user_info_sync_settings.unemployed_role_name',
                                                        )
                                                    "
                                                    maxlength="48"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </div>
                                        </template>
                                    </template>
                                </template>
                            </template>
                        </GenericContainerPanelEntry>
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
                                    <template v-if="!canSubmit">
                                        <LoadingIcon class="mr-2 h-5 w-5 animate-spin" />
                                    </template>
                                    {{ $t('common.save', 1) }}
                                </button>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
            </template>
        </div>
    </div>
</template>
