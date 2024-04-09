<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { vMaska } from 'maska';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { JobProps, UserInfoSyncUnemployedMode } from '~~/gen/ts/resources/users/jobs';
import SquareImg from '~/components/partials/elements/SquareImg.vue';

const { $grpc } = useNuxtApp();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const appConfig = useAppConfig();

const authStore = useAuthStore();

const notifications = useNotificatorStore();

const schema = z.object({
    jobLogo: zodFileSingleSchema(appConfig.filestore.fileSizes.images, appConfig.filestore.types.images),
});

type Schema = z.output<typeof schema>;

const state = reactive({
    jobLogo: undefined,
});

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

async function setJobProps(values: Schema): Promise<void> {
    if (!jobProps.value) {
        return;
    }

    if (values.jobLogo) {
        jobProps.value.logoUrl = { data: new Uint8Array(await values.jobLogo[0].arrayBuffer()) };
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

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setJobProps(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div>
        <template v-if="streamerMode">
            <UDashboardPanelContent class="pb-2">
                <UDashboardSection
                    :title="$t('system.streamer_mode.title')"
                    :description="$t('system.streamer_mode.description')"
                />
            </UDashboardPanelContent>
        </template>
        <template v-else>
            <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
                <DataPendingBlock
                    v-if="pending"
                    :message="$t('common.loading', [$t('components.rector.job_props.job_properties')])"
                />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('components.rector.job_props.job_properties')])"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="jobProps === null"
                    icon="i-mdi-tune"
                    :type="$t('components.rector.job_props.job_properties')"
                />

                <template v-else>
                    <UDashboardNavbar :title="$t('components.rector.job_props.job_properties')">
                        <template #right>
                            <UButton
                                trailing-icon="i-mdi-content-save"
                                :disabled="!canSubmit"
                                :loading="!canSubmit"
                                @click="onSubmitThrottle"
                            >
                                {{ $t('common.save', 1) }}
                            </UButton>
                        </template>
                    </UDashboardNavbar>

                    <UTabs
                        :items="[
                            { slot: 'jobprops', label: $t('components.rector.job_props.job_properties') },
                            { slot: 'discord', label: $t('common.discord') },
                        ]"
                        class="w-full"
                    >
                        <template #jobprops>
                            <UDashboardPanelContent class="pb-2">
                                <UDashboardSection
                                    :title="$t('components.rector.job_props.job_properties')"
                                    :description="$t('components.rector.job_props.your_job_properties')"
                                >
                                    <UFormGroup
                                        name="theme"
                                        :label="$t('common.theme')"
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        {{ jobProps.theme }}
                                    </UFormGroup>

                                    <UFormGroup
                                        name="livemapMarkerColor"
                                        :label="$t('components.rector.job_props.livemap_marker_color')"
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <ColorInput
                                            v-model="jobProps.livemapMarkerColor"
                                            disable-alpha
                                            format="hex"
                                            position="top"
                                        />
                                    </UFormGroup>

                                    <UFormGroup
                                        name="radioFrequency"
                                        :label="$t('common.radio_frequency')"
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UInput
                                            v-model="jobProps.radioFrequency"
                                            v-maska
                                            data-maska="0.9"
                                            data-maska-tokens="0:\d:multiple|9:\d:multiple"
                                            type="text"
                                            :placeholder="$t('common.radio_frequency')"
                                            :label="$t('common.radio_frequency')"
                                            maxlength="24"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </UFormGroup>

                                    <UFormGroup
                                        v-if="jobProps.quickButtons"
                                        name="quickButtons"
                                        :label="$t('components.rector.job_props.quick_buttons')"
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <div class="flex flex-col gap-2">
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
                                                    <span class="ml-3 text-sm font-medium">{{
                                                        $t('components.penaltycalculator.title')
                                                    }}</span>
                                                </div>
                                            </div>
                                            <div class="space-y-4">
                                                <div class="flex items-center">
                                                    <UToggle v-model="jobProps.quickButtons.bodyCheckup">
                                                        <span
                                                            :class="[
                                                                jobProps.quickButtons.bodyCheckup
                                                                    ? 'translate-x-5'
                                                                    : 'translate-x-0',
                                                                'pointer-events-none inline-block size-5 rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                            ]"
                                                        />
                                                    </UToggle>
                                                    <span class="ml-3 text-sm font-medium">{{
                                                        $t('components.bodycheckup.title')
                                                    }}</span>
                                                </div>
                                            </div>
                                        </div>
                                    </UFormGroup>

                                    <UFormGroup
                                        name="jobLogo"
                                        :label="$t('common.logo')"
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <div class="flex flex-col">
                                            <template v-if="isNUIAvailable()">
                                                <p class="text-sm">
                                                    {{ $t('system.not_supported_on_tablet.title') }}
                                                </p>
                                            </template>
                                            <template v-else>
                                                <UInput
                                                    name="jobLogo"
                                                    type="file"
                                                    accept="image/jpeg,image/jpg,image/png"
                                                    block
                                                    :placeholder="$t('common.image')"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </template>

                                            <div class="mt-2 flex w-full items-center justify-center">
                                                <SquareImg
                                                    v-if="jobProps.logoUrl?.url"
                                                    size="3xl"
                                                    :url="jobProps.logoUrl.url"
                                                    :no-blur="true"
                                                />
                                            </div>
                                        </div>
                                    </UFormGroup>
                                </UDashboardSection>
                            </UDashboardPanelContent>
                        </template>

                        <template #discord>
                            <UDashboardPanelContent v-if="jobProps.discordSyncSettings" class="pb-2">
                                <UDashboardSection
                                    :title="$t('components.rector.job_props.discord_sync_settings.title')"
                                    :description="$t('components.rector.job_props.discord_sync_settings.subtitle')"
                                >
                                    <UFormGroup
                                        name="discordGuildId"
                                        :label="$t('components.rector.job_props.discord_guild_id')"
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UInput
                                            v-model="jobProps.discordGuildId"
                                            type="text"
                                            :disabled="appConfig.discord.botInviteURL === undefined"
                                            :placeholder="$t('components.rector.job_props.discord_guild_id')"
                                            :label="$t('components.rector.job_props.discord_guild_id')"
                                            maxlength="70"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <UButton
                                            v-if="appConfig.discord.botInviteURL !== undefined"
                                            block
                                            class="mt-1"
                                            :to="appConfig.discord.botInviteURL"
                                            :external="true"
                                        >
                                            {{ $t('components.rector.job_props.invite_bot') }}
                                        </UButton>
                                        <p v-if="jobProps.discordLastSync" class="mt-2 text-xs">
                                            {{ $t('components.rector.job_props.last_sync') }}:
                                            <GenericTime :value="jobProps.discordLastSync" />
                                        </p>
                                    </UFormGroup>

                                    <UFormGroup
                                        name="statusLog"
                                        :label="$t('components.rector.job_props.status_log')"
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UToggle v-model="jobProps.discordSyncSettings.statusLog">
                                            <span
                                                :class="[
                                                    jobProps.discordSyncSettings.statusLog ? 'translate-x-5' : 'translate-x-0',
                                                    'pointer-events-none inline-block size-5 rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                                ]"
                                            />
                                        </UToggle>
                                    </UFormGroup>

                                    <UFormGroup
                                        v-if="jobProps.discordSyncSettings.statusLog"
                                        name="statusLog"
                                        :label="`${$t('components.rector.job_props.status_log')} ${$t('components.rector.job_props.status_log_settings.channel_id')}`"
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UInput
                                            v-model="jobProps.discordSyncSettings.statusLogSettings!.channelId"
                                            type="text"
                                            name="statusLogSettingsChannelId"
                                            :disabled="!jobProps.discordSyncSettings.statusLog"
                                            :placeholder="$t('components.rector.job_props.status_log_settings.channel_id')"
                                            :label="$t('components.rector.job_props.status_log_settings.channel_id')"
                                            maxlength="48"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </UFormGroup>

                                    <UFormGroup
                                        name="userInfoSync"
                                        :label="$t('components.rector.job_props.user_info_sync')"
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UToggle v-model="jobProps.discordSyncSettings.userInfoSync">
                                            <span class="sr-only">{{ $t('components.rector.job_props.user_info_sync') }}</span>
                                        </UToggle>
                                    </UFormGroup>

                                    <template
                                        v-if="
                                            jobProps.discordSyncSettings.userInfoSync &&
                                            jobProps.discordSyncSettings.userInfoSyncSettings
                                        "
                                    >
                                        <UFormGroup
                                            name="userInfoSync"
                                            :label="$t('components.rector.job_props.user_info_sync_settings.grade_role_format')"
                                            class="grid grid-cols-2 items-center gap-2"
                                            :ui="{ container: '' }"
                                        >
                                            <UInput
                                                v-model="jobProps.discordSyncSettings.userInfoSyncSettings.gradeRoleFormat"
                                                type="text"
                                                name="gradeRoleFormat"
                                                :placeholder="
                                                    $t('components.rector.job_props.user_info_sync_settings.grade_role_format')
                                                "
                                                :label="
                                                    $t('components.rector.job_props.user_info_sync_settings.grade_role_format')
                                                "
                                                maxlength="48"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </UFormGroup>

                                        <UFormGroup
                                            name="userInfoSync"
                                            :label="
                                                $t('components.rector.job_props.user_info_sync_settings.employee_role_enabled')
                                            "
                                            class="grid grid-cols-2 items-center gap-2"
                                            :ui="{ container: '' }"
                                        >
                                            <UToggle
                                                v-model="jobProps.discordSyncSettings.userInfoSyncSettings.employeeRoleEnabled"
                                            >
                                                <span class="sr-only">{{
                                                    $t(
                                                        'components.rector.job_props.user_info_sync_settings.employee_role_enabled',
                                                    )
                                                }}</span>
                                            </UToggle>
                                        </UFormGroup>

                                        <UFormGroup
                                            v-if="jobProps.discordSyncSettings.userInfoSyncSettings?.employeeRoleEnabled"
                                            name="userInfoSync"
                                            :label="
                                                $t('components.rector.job_props.user_info_sync_settings.employee_role_format')
                                            "
                                            class="grid grid-cols-2 items-center gap-2"
                                            :ui="{ container: '' }"
                                        >
                                            <UInput
                                                v-model="jobProps.discordSyncSettings.userInfoSyncSettings!.employeeRoleFormat"
                                                type="text"
                                                name="employeeRoleFormat"
                                                :placeholder="
                                                    $t(
                                                        'components.rector.job_props.user_info_sync_settings.employee_role_format',
                                                    )
                                                "
                                                :label="
                                                    $t(
                                                        'components.rector.job_props.user_info_sync_settings.employee_role_format',
                                                    )
                                                "
                                                maxlength="48"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </UFormGroup>

                                        <UFormGroup
                                            name="userInfoSync"
                                            :label="
                                                $t('components.rector.job_props.user_info_sync_settings.unemployed_enabled')
                                            "
                                            class="grid grid-cols-2 items-center gap-2"
                                            :ui="{ container: '' }"
                                        >
                                            <UToggle
                                                v-model="jobProps.discordSyncSettings.userInfoSyncSettings.unemployedEnabled"
                                            >
                                                <span class="sr-only">{{
                                                    $t('components.rector.job_props.user_info_sync_settings.unemployed_enabled')
                                                }}</span>
                                            </UToggle>
                                        </UFormGroup>

                                        <UFormGroup
                                            name="userInfoSync"
                                            :label="$t('components.rector.job_props.user_info_sync_settings.unemployed_mode')"
                                            class="grid grid-cols-2 items-center gap-2"
                                            :ui="{ container: '' }"
                                        >
                                            <USelectMenu
                                                v-model="jobProps.discordSyncSettings.userInfoSyncSettings.unemployedMode"
                                                value-attribute="value"
                                                :options="[
                                                    {
                                                        label: $t('enums.rector.UserInfoSyncUnemployedMode.GIVE_ROLE'),
                                                        value: UserInfoSyncUnemployedMode.GIVE_ROLE,
                                                    },
                                                    {
                                                        label: $t('enums.rector.UserInfoSyncUnemployedMode.GIVE_ROLE'),
                                                        value: UserInfoSyncUnemployedMode.KICK,
                                                    },
                                                ]"
                                            >
                                                <template #label>
                                                    {{
                                                        $t(
                                                            `enums.rector.UserInfoSyncUnemployedMode.${
                                                                UserInfoSyncUnemployedMode[
                                                                    jobProps.discordSyncSettings.userInfoSyncSettings
                                                                        .unemployedMode ?? 0
                                                                ]
                                                            }`,
                                                        )
                                                    }}
                                                </template>

                                                <template #option="{ option }">
                                                    <span class="truncate">{{
                                                        $t(
                                                            `enums.rector.UserInfoSyncUnemployedMode.${UserInfoSyncUnemployedMode[option.value]}`,
                                                        )
                                                    }}</span>
                                                </template>
                                            </USelectMenu>
                                        </UFormGroup>

                                        <UFormGroup
                                            v-if="jobProps.discordSyncSettings.userInfoSyncSettings.unemployedEnabled"
                                            name="userInfoSync"
                                            :label="
                                                $t('components.rector.job_props.user_info_sync_settings.unemployed_role_name')
                                            "
                                            class="grid grid-cols-2 items-center gap-2"
                                            :ui="{ container: '' }"
                                        >
                                            <UInput
                                                v-model="jobProps.discordSyncSettings.userInfoSyncSettings.unemployedRoleName"
                                                type="text"
                                                name="unemployedRoleName"
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
                                        </UFormGroup>

                                        <UFormGroup
                                            name="userInfoSync"
                                            :label="
                                                $t(
                                                    'components.rector.job_props.jobs_absence_settings.jobs_absence_role_enabled',
                                                )
                                            "
                                            class="grid grid-cols-2 items-center gap-2"
                                            :ui="{ container: '' }"
                                        >
                                            <UToggle v-model="jobProps.discordSyncSettings.jobsAbsence">
                                                <span class="sr-only">{{
                                                    $t(
                                                        'components.rector.job_props.jobs_absence_settings.jobs_absence_role_enabled',
                                                    )
                                                }}</span>
                                            </UToggle>
                                        </UFormGroup>

                                        <template v-if="jobProps.discordSyncSettings.jobsAbsenceSettings">
                                            <UFormGroup
                                                v-if="jobProps.discordSyncSettings.jobsAbsence"
                                                name="userInfoSync"
                                                :label="
                                                    $t(
                                                        'components.rector.job_props.jobs_absence_settings.jobs_absence_role_name',
                                                    )
                                                "
                                                class="grid grid-cols-2 items-center gap-2"
                                                :ui="{ container: '' }"
                                            >
                                                <UInput
                                                    v-model="jobProps.discordSyncSettings.jobsAbsenceSettings.absenceRole"
                                                    type="text"
                                                    name="jobsAbsenceRole"
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
                                            </UFormGroup>
                                        </template>
                                    </template>
                                </UDashboardSection>
                            </UDashboardPanelContent>
                        </template>
                    </UTabs>
                </template>
            </UForm>
        </template>
    </div>
</template>
