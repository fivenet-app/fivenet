<script lang="ts" setup>
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue';
import { mimes, size } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { vMaska } from 'maska';
import { CheckIcon, ChevronDownIcon, LoadingIcon, TuneIcon } from 'mdi-vue3';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import { defineRule } from 'vee-validate';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { availableThemes, useSettingsStore } from '~/store/settings';
import { JobProps, UserInfoSyncUnemployedMode } from '~~/gen/ts/resources/users/jobs';
import GenericContainerPanel from '~/components/partials/elements/GenericContainerPanel.vue';
import GenericContainerPanelEntry from '~/components/partials/elements/GenericContainerPanelEntry.vue';
import SquareImg from '~/components/partials/elements/SquareImg.vue';

const { $grpc } = useNuxtApp();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

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

interface FormData {
    jobLogo?: Blob;
}

async function setJobProps(values: FormData): Promise<void> {
    if (!jobProps.value) {
        return;
    }

    if (values.jobLogo) {
        jobProps.value.logoUrl = { data: new Uint8Array(await values.jobLogo.arrayBuffer()) };
    }

    try {
        const { response } = await $grpc.getRectorClient().setJobProps({
            jobProps: jobProps.value,
        });

        notifications.add({
            title: { key: 'notifications.rector.job_props.title', parameters: {} },
            description: { key: 'notifications.rector.job_props.content', parameters: {} },
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

defineRule('mimes', mimes);
defineRule('size', size);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        jobLogo: { required: false, mimes: ['image/jpeg', 'image/jpg', 'image/png'], size: 2000 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> => await setJobProps(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
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
            <DataPendingBlock
                v-if="pending"
                :message="$t('common.loading', [$t('components.rector.job_props.job_properties')])"
            />
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
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        >
                                            <span class="block truncate">
                                                {{ availableThemes.find((t) => t.key === jobProps?.theme)?.name }}
                                            </span>
                                            <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                                                <ChevronDownIcon class="size-5 text-gray-400" />
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
                                                            'relative cursor-default select-none py-2 pl-8 pr-4',
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
                                                            <CheckIcon class="size-5" />
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
                                <UInput
                                    v-model="jobProps.radioFrequency"
                                    v-maska
                                    data-maska="0.9"
                                    data-maska-tokens="0:\d:multiple|9:\d:multiple"
                                    type="text"
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                        <div class="flex items-center">
                                            <UToggle v-model="jobProps.quickButtons.penaltyCalculator">
                                                <span
                                                    :class="[
                                                        jobProps.quickButtons.penaltyCalculator
                                                            ? 'translate-x-5'
                                                            : 'translate-x-0',
                                                        'pointer-events-none inline-block size-5 rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                    ]"
                                                />
                                            </UToggle>
                                            <span class="ml-3 text-sm font-medium text-gray-300">{{
                                                $t('components.penaltycalculator.title')
                                            }}</span>
                                        </div>
                                    </div>
                                    <div class="space-y-5">
                                        <div class="flex items-center">
                                            <UToggle v-model="jobProps.quickButtons.bodyCheckup">
                                                <span
                                                    :class="[
                                                        jobProps.quickButtons.bodyCheckup ? 'translate-x-5' : 'translate-x-0',
                                                        'pointer-events-none inline-block size-5 rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                    ]"
                                                />
                                            </UToggle>
                                            <span class="ml-3 text-sm font-medium text-gray-300">{{
                                                $t('components.bodycheckup.title')
                                            }}</span>
                                        </div>
                                    </div>
                                </fieldset>
                            </template>
                        </GenericContainerPanelEntry>
                        <GenericContainerPanelEntry>
                            <template #title>
                                {{ $t('common.logo') }}
                            </template>
                            <template #default>
                                <div class="flex flex-col">
                                    <template v-if="isNUIAvailable()">
                                        <p class="text-sm">
                                            {{ $t('system.not_supported_on_tablet.title') }}
                                        </p>
                                    </template>
                                    <template v-else>
                                        <VeeField
                                            v-slot="{ handleChange, handleBlur }"
                                            name="jobLogo"
                                            :placeholder="$t('common.image')"
                                            :label="$t('common.image')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        >
                                            <UInput
                                                type="file"
                                                accept="image/jpeg,image/jpg,image/png"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                @change="handleChange"
                                                @blur="handleBlur"
                                            />
                                        </VeeField>
                                        <VeeErrorMessage name="jobLogo" as="p" class="text-sm text-error-400" />
                                    </template>

                                    <SquareImg
                                        v-if="jobProps.logoUrl?.url"
                                        size="xl"
                                        :url="jobProps.logoUrl.url"
                                        :no-blur="true"
                                        class="mt-2"
                                    />
                                </div>
                            </template>
                        </GenericContainerPanelEntry>

                        <!-- Save button -->
                        <GenericContainerPanelEntry v-if="can('RectorService.SetJobProps')">
                            <template #default>
                                <UButton
                                    class="flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                    :class="[
                                        !canSubmit || !meta.valid
                                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                            : 'bg-primary-500 hover:bg-primary-400',
                                    ]"
                                    :disabled="!canSubmit || !meta.valid"
                                    @click="onSubmitThrottle"
                                >
                                    <template v-if="!canSubmit">
                                        <LoadingIcon class="mr-2 size-5 animate-spin" />
                                    </template>
                                    {{ $t('common.save', 1) }}
                                </UButton>
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
                                <UInput
                                    v-model="jobProps.discordGuildId"
                                    type="text"
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                    class="mt-2 flex w-full justify-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold transition-colors hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                >
                                    {{ $t('components.rector.job_props.invite_bot') }}
                                </NuxtLink>
                                <p v-if="jobProps.discordLastSync" class="mt-2 text-xs">
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
                                <div class="mb-1 flex items-center">
                                    <UToggle v-model="jobProps.discordSyncSettings.statusLog">
                                        <span
                                            :class="[
                                                jobProps.discordSyncSettings.statusLog ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block size-5 rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </UToggle>
                                    <span class="ml-3 text-sm font-medium text-gray-300">{{ $t('common.enabled') }}</span>
                                </div>

                                <template v-if="jobProps.discordSyncSettings.statusLog">
                                    <label for="statusLogSettingsChannelId">
                                        {{ $t('components.rector.job_props.status_log_settings.channel_id') }}:
                                    </label>
                                    <UInput
                                        v-model="jobProps.discordSyncSettings.statusLogSettings!.channelId"
                                        type="text"
                                        name="statusLogSettingsChannelId"
                                        :disabled="!jobProps.discordSyncSettings.statusLog"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                <div class="mb-1 flex items-center">
                                    <UToggle v-model="jobProps.discordSyncSettings.userInfoSync">
                                        <span class="sr-only">{{ $t('components.rector.job_props.user_info_sync') }}</span>
                                    </UToggle>
                                    <span class="ml-3 text-sm font-medium text-gray-300">{{ $t('common.enabled') }}</span>
                                </div>

                                <template v-if="jobProps.discordSyncSettings.userInfoSync">
                                    <label for="gradeRoleFormat">
                                        {{ $t('components.rector.job_props.user_info_sync_settings.grade_role_format') }}:
                                    </label>
                                    <UInput
                                        v-model="jobProps.discordSyncSettings.userInfoSyncSettings!.gradeRoleFormat"
                                        type="text"
                                        name="gradeRoleFormat"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="
                                            $t('components.rector.job_props.user_info_sync_settings.grade_role_format')
                                        "
                                        :label="$t('components.rector.job_props.user_info_sync_settings.grade_role_format')"
                                        maxlength="48"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />

                                    <!-- UserInfo Sync Settings -->
                                    <div
                                        v-if="jobProps.discordSyncSettings.userInfoSyncSettings !== undefined"
                                        class="mb-1 mt-2 flex items-center"
                                    >
                                        <UToggle
                                            v-model="jobProps.discordSyncSettings.userInfoSyncSettings.employeeRoleEnabled"
                                        >
                                            <span class="sr-only">{{
                                                $t('components.rector.job_props.user_info_sync_settings.employee_role_enabled')
                                            }}</span>
                                        </UToggle>
                                        <span class="ml-3 text-sm font-medium text-gray-300">{{
                                            $t('components.rector.job_props.user_info_sync_settings.employee_role_enabled')
                                        }}</span>
                                    </div>

                                    <div v-if="jobProps.discordSyncSettings.userInfoSyncSettings?.employeeRoleEnabled">
                                        <label for="employeeRoleFormat">
                                            {{
                                                $t('components.rector.job_props.user_info_sync_settings.employee_role_format')
                                            }}:
                                        </label>
                                        <UInput
                                            v-model="jobProps.discordSyncSettings.userInfoSyncSettings!.employeeRoleFormat"
                                            type="text"
                                            name="employeeRoleFormat"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                        <div class="mb-1 mt-2 flex items-center">
                                            <UToggle
                                                v-model="jobProps.discordSyncSettings.userInfoSyncSettings.unemployedEnabled"
                                            >
                                                <span class="sr-only">{{
                                                    $t('components.rector.job_props.user_info_sync_settings.unemployed_enabled')
                                                }}</span>
                                            </UToggle>
                                            <span class="ml-3 text-sm font-medium text-gray-300">{{
                                                $t('components.rector.job_props.user_info_sync_settings.unemployed_enabled')
                                            }}</span>
                                        </div>
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
                                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                                                <ChevronDownIcon class="size-5 text-gray-400" />
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
                                                                            'relative cursor-default select-none py-2 pl-8 pr-4',
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
                                                                            <CheckIcon class="size-5" />
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
                                                <UInput
                                                    v-model="
                                                        jobProps.discordSyncSettings.userInfoSyncSettings!.unemployedRoleName
                                                    "
                                                    type="text"
                                                    name="unemployedRoleName"
                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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

                                        <div class="mb-1 mt-2 flex items-center">
                                            <UToggle v-model="jobProps.discordSyncSettings.jobsAbsence">
                                                <span class="sr-only">{{
                                                    $t(
                                                        'components.rector.job_props.jobs_absence_settings.jobs_absence_role_enabled',
                                                    )
                                                }}</span>
                                            </UToggle>
                                            <span class="ml-3 text-sm font-medium text-gray-300">{{
                                                $t(
                                                    'components.rector.job_props.jobs_absence_settings.jobs_absence_role_enabled',
                                                )
                                            }}</span>
                                        </div>

                                        <template v-if="jobProps.discordSyncSettings.jobsAbsence">
                                            <div>
                                                <label for="jobsAbsenceRole">
                                                    {{
                                                        $t(
                                                            'components.rector.job_props.jobs_absence_settings.jobs_absence_role_name',
                                                        )
                                                    }}:
                                                </label>
                                                <UInput
                                                    v-model="jobProps.discordSyncSettings.jobsAbsenceSettings!.absenceRole"
                                                    type="text"
                                                    name="jobsAbsenceRole"
                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    :placeholder="
                                                        $t(
                                                            'components.rector.job_props.jobs_absence_settings.jobs_absence_role_name',
                                                        )
                                                    "
                                                    :label="
                                                        $t(
                                                            'components.rector.job_props.jobs_absence_settings.jobs_absence_role_name',
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
                                <UButton
                                    class="flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                    :class="[
                                        !canSubmit || !meta.valid
                                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                            : 'bg-primary-500 hover:bg-primary-400',
                                    ]"
                                    :disabled="!canSubmit || !meta.valid"
                                    @click="onSubmitThrottle"
                                >
                                    <template v-if="!canSubmit">
                                        <LoadingIcon class="mr-2 size-5 animate-spin" />
                                    </template>
                                    {{ $t('common.save', 1) }}
                                </UButton>
                            </template>
                        </GenericContainerPanelEntry>
                    </template>
                </GenericContainerPanel>
            </template>
        </template>
    </div>
</template>
