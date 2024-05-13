<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { vMaska } from 'maska';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { JobProps, UserInfoSyncUnemployedMode } from '~~/gen/ts/resources/users/jobs';
import SquareImg from '~/components/partials/elements/SquareImg.vue';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const appConfig = useAppConfig();

const authStore = useAuthStore();

const notifications = useNotificatorStore();

const schema = z.object({
    livemapMarkerColor: z.string().length(7),
    quickButtons: z.object({
        penaltyCalculator: z.boolean(),
        bodyCheckup: z.boolean(),
    }),
    radioFrequency: z.string(),
    discordGuildId: z.string(),
    discordSyncSettings: z.object({
        userInfoSync: z.boolean(),
        userInfoSyncSettings: z.object({
            employeeRoleEnabled: z.boolean(),
            employeeRoleFormat: z.string(),
            gradeRoleFormat: z.string(),
            unemployedEnabled: z.boolean(),
            unemployedMode: z.nativeEnum(UserInfoSyncUnemployedMode),
            unemployedRoleName: z.string(),
        }),
        statusLog: z.boolean(),
        statusLogSettings: z.object({
            channelId: z.string(),
        }),
        jobsAbsence: z.boolean(),
        jobsAbsenceSettings: z.object({
            absenceRole: z.string(),
        }),
        groupSyncSettings: z.object({
            ignoredRoleIds: z.string().array().max(20),
        }),
    }),
    logoUrl: zodFileSingleSchema(appConfig.filestore.fileSizes.images, appConfig.filestore.types.images, true).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    livemapMarkerColor: '',
    quickButtons: {
        penaltyCalculator: false,
        bodyCheckup: false,
    },
    radioFrequency: '',
    discordGuildId: '',
    discordSyncSettings: {
        userInfoSync: false,
        userInfoSyncSettings: {
            employeeRoleEnabled: false,
            employeeRoleFormat: '',
            gradeRoleFormat: '',
            unemployedEnabled: false,
            unemployedMode: UserInfoSyncUnemployedMode.GIVE_ROLE,
            unemployedRoleName: '',
        },
        statusLog: false,
        statusLogSettings: {
            channelId: '',
        },
        jobsAbsence: false,
        jobsAbsenceSettings: {
            absenceRole: '',
        },
        groupSyncSettings: {
            ignoredRoleIds: [],
        },
    },
    logoUrl: undefined,
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

const { data: jobProps, pending: loading, refresh, error } = useLazyAsyncData(`rector-jobprops`, () => getJobProps());

async function setJobProps(values: Schema): Promise<void> {
    if (!jobProps.value) {
        return;
    }

    jobProps.value.livemapMarkerColor = values.livemapMarkerColor;
    jobProps.value.quickButtons = values.quickButtons;
    jobProps.value.radioFrequency = values.radioFrequency;
    jobProps.value.discordGuildId = values.discordGuildId;
    jobProps.value.discordSyncSettings = values.discordSyncSettings;
    if (values.logoUrl) {
        jobProps.value.logoUrl = { data: new Uint8Array(await values.logoUrl[0].arrayBuffer()) };
    }

    try {
        const { response } = await $grpc.getRectorClient().setJobProps({
            jobProps: jobProps.value,
        });

        notifications.add({
            title: { key: 'notifications.rector.job_props.title', parameters: {} },
            description: { key: 'notifications.rector.job_props.content', parameters: {} },
            type: NotificationType.SUCCESS,
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

function setSettingsValues(): void {
    if (!jobProps.value) {
        return;
    }

    if (jobProps.value.livemapMarkerColor) {
        state.livemapMarkerColor = jobProps.value.livemapMarkerColor;
    }
    if (jobProps.value.quickButtons) {
        state.quickButtons = jobProps.value.quickButtons;
    }
    if (jobProps.value.radioFrequency) {
        state.radioFrequency = jobProps.value.radioFrequency;
    }
    if (jobProps.value.discordGuildId) {
        state.discordGuildId = jobProps.value.discordGuildId;
    }
    if (jobProps.value.discordSyncSettings) {
        state.discordSyncSettings.statusLog = jobProps.value.discordSyncSettings.statusLog;
        if (jobProps.value.discordSyncSettings.statusLogSettings) {
            state.discordSyncSettings.statusLogSettings = jobProps.value.discordSyncSettings.statusLogSettings;
        }
        state.discordSyncSettings.userInfoSync = jobProps.value.discordSyncSettings.userInfoSync;
        if (jobProps.value.discordSyncSettings.userInfoSyncSettings) {
            state.discordSyncSettings.userInfoSyncSettings = jobProps.value.discordSyncSettings.userInfoSyncSettings;
        }
        state.discordSyncSettings.jobsAbsence = jobProps.value.discordSyncSettings.jobsAbsence;
        if (jobProps.value.discordSyncSettings.jobsAbsenceSettings) {
            state.discordSyncSettings.jobsAbsenceSettings = jobProps.value.discordSyncSettings.jobsAbsenceSettings;
        }
        if (jobProps.value.discordSyncSettings.groupSyncSettings) {
            state.discordSyncSettings.groupSyncSettings = jobProps.value.discordSyncSettings.groupSyncSettings;
        }
    }
}

watch(jobProps, () => setSettingsValues());

const items = [
    {
        slot: 'jobprops',
        label: t('components.rector.job_props.job_properties'),
        icon: 'i-mdi-settings',
    },
    { slot: 'discord', label: t('common.discord'), icon: 'i-simple-icons-discord' },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.findIndex((item) => item.label === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value].slot }, hash: '#' });
    },
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setJobProps(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <template v-if="streamerMode">
        <UDashboardNavbar :title="$t('pages.rector.settings.title')">
            <template #right>
                <UButton color="black" icon="i-mdi-arrow-back" to="/rector">
                    {{ $t('common.back') }}
                </UButton>
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent class="pb-24">
            <UDashboardSection
                :title="$t('system.streamer_mode.title')"
                :description="$t('system.streamer_mode.description')"
            />
        </UDashboardPanelContent>
    </template>
    <template v-else>
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UDashboardNavbar :title="$t('components.rector.job_props.job_properties')">
                <template #right>
                    <UButton color="black" icon="i-mdi-arrow-back" to="/rector">
                        {{ $t('common.back') }}
                    </UButton>

                    <UButton
                        v-if="!!jobProps"
                        type="submit"
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                    >
                        {{ $t('common.save', 1) }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <DataPendingBlock
                v-if="loading"
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
                <UTabs v-model="selectedTab" :items="items" class="w-full" :ui="{ list: { rounded: '' } }">
                    <template #default="{ item, selected }">
                        <div class="relative flex items-center gap-2 truncate">
                            <UIcon :name="item.icon" class="size-4 shrink-0" />

                            <span class="truncate">{{ item.label }}</span>

                            <span
                                v-if="selected"
                                class="bg-primary-500 dark:bg-primary-400 absolute -right-4 size-2 rounded-full"
                            />
                        </div>
                    </template>

                    <template #jobprops>
                        <UDashboardPanelContent class="pb-24">
                            <UDashboardSection
                                :title="$t('components.rector.job_props.job_properties')"
                                :description="$t('components.rector.job_props.your_job_properties')"
                            >
                                <UFormGroup
                                    name="livemapMarkerColor"
                                    :label="$t('components.rector.job_props.livemap_marker_color')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <ColorPicker v-model="state.livemapMarkerColor" />
                                </UFormGroup>

                                <UFormGroup
                                    name="radioFrequency"
                                    :label="$t('common.radio_frequency')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.radioFrequency"
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
                                                <UToggle v-model="state.quickButtons.penaltyCalculator">
                                                    <span
                                                        :class="[
                                                            state.quickButtons.penaltyCalculator
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
                                                <UToggle v-model="state.quickButtons.bodyCheckup">
                                                    <span
                                                        :class="[
                                                            state.quickButtons.bodyCheckup ? 'translate-x-5' : 'translate-x-0',
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
                                                @change="state.logoUrl = $event"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </template>

                                        <div class="mt-2 flex w-full items-center justify-center">
                                            <SquareImg
                                                v-if="jobProps.logoUrl?.url"
                                                size="3xl"
                                                :src="jobProps.logoUrl.url"
                                                :no-blur="true"
                                            />
                                        </div>
                                    </div>
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>

                    <template #discord>
                        <UDashboardPanelContent v-if="jobProps.discordSyncSettings" class="pb-24">
                            <UDashboardSection
                                :title="$t('components.rector.job_props.discord_sync_settings.title')"
                                :description="$t('components.rector.job_props.discord_sync_settings.subtitle')"
                            >
                                <UFormGroup
                                    name="discordGuildId"
                                    :label="$t('components.rector.job_props.discord_sync_settings.discord_guild_id')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.discordGuildId"
                                        type="text"
                                        :disabled="appConfig.discord.botInviteURL === undefined"
                                        :placeholder="$t('components.rector.job_props.discord_sync_settings.discord_guild_id')"
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
                                        {{ $t('components.rector.job_props.discord_sync_settings.invite_bot') }}
                                    </UButton>
                                    <p v-if="jobProps.discordLastSync" class="mt-2 text-xs">
                                        {{ $t('components.rector.job_props.discord_sync_settings.last_sync') }}:
                                        <GenericTime :value="jobProps.discordLastSync" />
                                    </p>
                                </UFormGroup>

                                <UFormGroup
                                    name="statusLog"
                                    :label="$t('components.rector.job_props.discord_sync_settings.status_log')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.discordSyncSettings.statusLog">
                                        <span
                                            :class="[
                                                state.discordSyncSettings.statusLog ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block size-5 rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </UToggle>
                                </UFormGroup>

                                <UFormGroup
                                    name="statusLog"
                                    :label="`${$t('components.rector.job_props.discord_sync_settings.status_log')} ${$t('components.rector.job_props.discord_sync_settings.status_log_settings.channel_id')}`"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.discordSyncSettings.statusLogSettings!.channelId"
                                        type="text"
                                        name="statusLogSettingsChannelId"
                                        :disabled="!state.discordSyncSettings.statusLog"
                                        :placeholder="
                                            $t(
                                                'components.rector.job_props.discord_sync_settings.status_log_settings.channel_id',
                                            )
                                        "
                                        maxlength="48"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </UFormGroup>

                                <UFormGroup
                                    name="userInfoSync"
                                    :label="$t('components.rector.job_props.discord_sync_settings.user_info_sync')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.discordSyncSettings.userInfoSync">
                                        <span class="sr-only">{{
                                            $t('components.rector.job_props.discord_sync_settings.user_info_sync')
                                        }}</span>
                                    </UToggle>
                                </UFormGroup>

                                <template v-if="jobProps.discordSyncSettings.userInfoSyncSettings">
                                    <UFormGroup
                                        name="userInfoSync"
                                        :label="
                                            $t(
                                                'components.rector.job_props.discord_sync_settings.user_info_sync_settings.grade_role_format',
                                            )
                                        "
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UInput
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.gradeRoleFormat"
                                            type="text"
                                            name="gradeRoleFormat"
                                            :disabled="!state.discordSyncSettings.userInfoSync"
                                            :placeholder="
                                                $t(
                                                    'components.rector.job_props.discord_sync_settings.user_info_sync_settings.grade_role_format',
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
                                                'components.rector.job_props.discord_sync_settings.user_info_sync_settings.employee_role_enabled',
                                            )
                                        "
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UToggle
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.employeeRoleEnabled"
                                            :disabled="!state.discordSyncSettings.userInfoSync"
                                        >
                                            <span class="sr-only">{{
                                                $t(
                                                    'components.rector.job_props.discord_sync_settings.user_info_sync_settings.employee_role_enabled',
                                                )
                                            }}</span>
                                        </UToggle>
                                    </UFormGroup>

                                    <UFormGroup
                                        name="userInfoSync"
                                        :label="
                                            $t(
                                                'components.rector.job_props.discord_sync_settings.user_info_sync_settings.employee_role_format',
                                            )
                                        "
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UInput
                                            v-model="state.discordSyncSettings.userInfoSyncSettings!.employeeRoleFormat"
                                            type="text"
                                            name="employeeRoleFormat"
                                            :disabled="
                                                !state.discordSyncSettings.userInfoSync ||
                                                !state.discordSyncSettings.userInfoSyncSettings?.employeeRoleEnabled
                                            "
                                            :placeholder="
                                                $t(
                                                    'components.rector.job_props.discord_sync_settings.user_info_sync_settings.employee_role_format',
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
                                                'components.rector.job_props.discord_sync_settings.user_info_sync_settings.unemployed_enabled',
                                            )
                                        "
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UToggle
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.unemployedEnabled"
                                            :disabled="!state.discordSyncSettings.userInfoSync"
                                        >
                                            <span class="sr-only">{{
                                                $t(
                                                    'components.rector.job_props.discord_sync_settings.user_info_sync_settings.unemployed_enabled',
                                                )
                                            }}</span>
                                        </UToggle>
                                    </UFormGroup>

                                    <UFormGroup
                                        name="userInfoSync"
                                        :label="
                                            $t(
                                                'components.rector.job_props.discord_sync_settings.user_info_sync_settings.unemployed_mode',
                                            )
                                        "
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <USelectMenu
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.unemployedMode"
                                            :disabled="
                                                !state.discordSyncSettings.userInfoSync ||
                                                !state.discordSyncSettings.userInfoSyncSettings.unemployedEnabled
                                            "
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
                                            :searchable-placeholder="$t('common.search_field')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        >
                                            <template #label>
                                                {{
                                                    $t(
                                                        `enums.rector.UserInfoSyncUnemployedMode.${
                                                            UserInfoSyncUnemployedMode[
                                                                state.discordSyncSettings.userInfoSyncSettings.unemployedMode ??
                                                                    0
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
                                        name="userInfoSync"
                                        :label="
                                            $t(
                                                'components.rector.job_props.discord_sync_settings.user_info_sync_settings.unemployed_role_name',
                                            )
                                        "
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UInput
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.unemployedRoleName"
                                            type="text"
                                            name="unemployedRoleName"
                                            :disabled="
                                                !state.discordSyncSettings.userInfoSync ||
                                                !state.discordSyncSettings.userInfoSyncSettings.unemployedEnabled
                                            "
                                            :placeholder="
                                                $t(
                                                    'components.rector.job_props.discord_sync_settings.user_info_sync_settings.unemployed_role_name',
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
                                                'components.rector.job_props.discord_sync_settings.jobs_absence_settings.jobs_absence_role_enabled',
                                            )
                                        "
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UToggle
                                            v-model="state.discordSyncSettings.jobsAbsence"
                                            :disabled="!state.discordSyncSettings.userInfoSync"
                                        >
                                            <span class="sr-only">{{
                                                $t(
                                                    'components.rector.job_props.discord_sync_settings.jobs_absence_settings.jobs_absence_role_enabled',
                                                )
                                            }}</span>
                                        </UToggle>
                                    </UFormGroup>

                                    <template v-if="jobProps.discordSyncSettings.jobsAbsenceSettings">
                                        <UFormGroup
                                            name="userInfoSync"
                                            :label="
                                                $t(
                                                    'components.rector.job_props.discord_sync_settings.jobs_absence_settings.jobs_absence_role_name',
                                                )
                                            "
                                            class="grid grid-cols-2 items-center gap-2"
                                            :ui="{ container: '' }"
                                        >
                                            <UInput
                                                v-model="state.discordSyncSettings.jobsAbsenceSettings.absenceRole"
                                                type="text"
                                                name="jobsAbsenceRole"
                                                :disabled="
                                                    !state.discordSyncSettings.userInfoSync ||
                                                    !state.discordSyncSettings.jobsAbsence
                                                "
                                                :placeholder="
                                                    $t(
                                                        'components.rector.job_props.discord_sync_settings.jobs_absence_settings.jobs_absence_role_name',
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

                            <UDashboardSection
                                v-if="jobProps.discordSyncSettings?.groupSyncSettings"
                                :title="$t('components.rector.job_props.discord_sync_settings.group_sync_settings.title')"
                                :description="
                                    $t('components.rector.job_props.discord_sync_settings.group_sync_settings.subtitle')
                                "
                            >
                                <UFormGroup
                                    name="groupSyncSettingsIgnoredIds"
                                    :label="
                                        $t(
                                            'components.rector.job_props.discord_sync_settings.group_sync_settings.ignored_role_ids.title',
                                        )
                                    "
                                    :description="
                                        $t(
                                            'components.rector.job_props.discord_sync_settings.group_sync_settings.ignored_role_ids.description',
                                        )
                                    "
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <div class="flex flex-col gap-1">
                                        <div
                                            v-for="(_, idx) in state.discordSyncSettings.groupSyncSettings.ignoredRoleIds"
                                            :key="idx"
                                            class="flex items-center gap-1"
                                        >
                                            <UFormGroup :name="`citizenAttributes.list.${idx}.name`" class="flex-1">
                                                <UInput
                                                    v-model="state.discordSyncSettings.groupSyncSettings.ignoredRoleIds[idx]"
                                                    :name="`groupSyncSettingsIgnoredIds.${idx}`"
                                                    type="text"
                                                    class="w-full"
                                                    :disabled="!state.discordSyncSettings.userInfoSync"
                                                    :placeholder="
                                                        $t(
                                                            'components.rector.job_props.discord_sync_settings.group_sync_settings.ignored_role_ids.field',
                                                        )
                                                    "
                                                    maxlength="24"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </UFormGroup>

                                            <UButton
                                                :ui="{ rounded: 'rounded-full' }"
                                                :disabled="!canSubmit"
                                                icon="i-mdi-close"
                                                @click="
                                                    state.discordSyncSettings?.groupSyncSettings.ignoredRoleIds.splice(idx, 1)
                                                "
                                            />
                                        </div>
                                    </div>

                                    <UButton
                                        :ui="{ rounded: 'rounded-full' }"
                                        :disabled="!canSubmit"
                                        :class="
                                            state.discordSyncSettings?.groupSyncSettings.ignoredRoleIds.length ? 'mt-2' : ''
                                        "
                                        icon="i-mdi-plus"
                                        @click="state.discordSyncSettings?.groupSyncSettings.ignoredRoleIds.push('')"
                                    />
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>
                </UTabs>
            </template>
        </UForm>
    </template>
</template>
