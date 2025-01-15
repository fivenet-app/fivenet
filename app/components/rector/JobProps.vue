<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { vMaska } from 'maska/vue';
import { CodeDiff } from 'v-code-diff';
import { z } from 'zod';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { JobProps } from '~~/gen/ts/resources/users/job_props';
import { type DiscordSyncChange, UserInfoSyncUnemployedMode } from '~~/gen/ts/resources/users/job_settings';
import StreamerModeAlert from '../partials/StreamerModeAlert.vue';

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
        mathCalculator: z.boolean(),
    }),
    radioFrequency: z.string().max(24),
    discordGuildId: z.string().max(48),
    discordSyncSettings: z.object({
        dryRun: z.boolean(),
        userInfoSync: z.boolean(),
        userInfoSyncSettings: z.object({
            employeeRoleEnabled: z.boolean(),
            employeeRoleFormat: z.string().max(64),
            gradeRoleFormat: z.string().max(64),
            unemployedEnabled: z.boolean(),
            unemployedMode: z.nativeEnum(UserInfoSyncUnemployedMode),
            unemployedRoleName: z.string().max(64),
            syncNicknames: z.boolean(),
            groupMapping: z
                .object({
                    name: z.string().max(64),
                    fromGrade: z.number().min(0).max(99999),
                    toGrade: z.number().min(0).max(99999),
                })
                .array()
                .max(25),
        }),
        statusLog: z.boolean(),
        statusLogSettings: z.object({
            channelId: z.string().max(64),
        }),
        jobsAbsence: z.boolean(),
        jobsAbsenceSettings: z.object({
            absenceRole: z.string().max(64),
        }),
        groupSyncSettings: z.object({
            ignoredRoleIds: z.string().max(64).array().max(20),
        }),
        qualificationsRoleFormat: z.string().max(64),
    }),
    logoUrl: zodFileSingleSchema(appConfig.fileUpload.fileSizes.images, appConfig.fileUpload.types.images, true).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    livemapMarkerColor: '',
    quickButtons: {
        penaltyCalculator: false,
        bodyCheckup: false,
        mathCalculator: false,
    },
    radioFrequency: '',
    discordGuildId: '',
    discordSyncSettings: {
        dryRun: true,
        userInfoSync: false,
        userInfoSyncSettings: {
            employeeRoleEnabled: false,
            employeeRoleFormat: '',
            gradeRoleFormat: '',
            unemployedEnabled: false,
            unemployedMode: UserInfoSyncUnemployedMode.GIVE_ROLE,
            unemployedRoleName: '',
            syncNicknames: true,
            groupMapping: [],
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
        qualificationsRoleFormat: '',
    },
    logoUrl: undefined,
});

async function getJobProps(): Promise<JobProps> {
    try {
        const call = getGRPCRectorClient().getJobProps({});
        const { response } = await call;

        return response.jobProps!;
    } catch (e) {
        handleGRPCError(e as RpcError);
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
    jobProps.value.discordGuildId = values.discordGuildId.trim().length > 0 ? values.discordGuildId : undefined;
    jobProps.value.discordSyncSettings = values.discordSyncSettings;
    if (values.logoUrl && values.logoUrl[0]) {
        jobProps.value.logoUrl = { data: new Uint8Array(await values.logoUrl[0].arrayBuffer()) };
    }

    try {
        const { response } = await getGRPCRectorClient().setJobProps({
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
        handleGRPCError(e as RpcError);
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
        state.discordSyncSettings.dryRun = jobProps.value.discordSyncSettings.dryRun;
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

        selectedChange.value = jobProps.value.discordSyncChanges?.changes.at(
            jobProps.value.discordSyncChanges?.changes.length - 1,
        );

        state.discordSyncSettings.qualificationsRoleFormat = jobProps.value.discordSyncSettings.qualificationsRoleFormat;
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
        const index = items.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value]?.slot }, hash: '#' });
    },
});

const selectedChange = ref<DiscordSyncChange | undefined>();

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (event.submitter?.getAttribute('role') === 'tab') {
        return;
    }

    canSubmit.value = false;
    await setJobProps(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <template v-if="streamerMode">
        <UDashboardNavbar :title="$t('pages.rector.settings.title')">
            <template #right>
                <PartialsBackButton fallback-to="/rector" />
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent>
            <StreamerModeAlert />
        </UDashboardPanelContent>
    </template>
    <UForm
        v-else
        :schema="schema"
        :state="state"
        class="flex min-h-screen w-full max-w-full flex-1 flex-col overflow-y-auto"
        @submit="onSubmitThrottle"
    >
        <UDashboardNavbar :title="$t('components.rector.job_props.job_properties')">
            <template #right>
                <PartialsBackButton fallback-to="/rector" />

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

        <UDashboardPanelContent class="p-0">
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('components.rector.job_props.job_properties')])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!jobProps" icon="i-mdi-tune" :type="$t('components.rector.job_props.job_properties')" />

            <template v-else-if="loading || jobProps">
                <UTabs v-model="selectedTab" :items="items" class="w-full" :ui="{ list: { rounded: '' } }">
                    <template #jobprops>
                        <div v-if="loading" class="space-y-1 px-4">
                            <USkeleton v-for="idx in 5" :key="idx" class="h-20 w-full" />
                        </div>

                        <UDashboardPanelContent v-else>
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
                                    <ColorPickerClient v-model="state.livemapMarkerColor" />
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
                                            <div class="flex items-center gap-2">
                                                <UToggle v-model="state.quickButtons.penaltyCalculator" />
                                                <span class="text-sm font-medium">{{
                                                    $t('components.penaltycalculator.title')
                                                }}</span>
                                            </div>
                                        </div>

                                        <div class="space-y-4">
                                            <div class="flex items-center gap-2">
                                                <UToggle v-model="state.quickButtons.bodyCheckup" />
                                                <span class="text-sm font-medium">{{
                                                    $t('components.bodycheckup.title')
                                                }}</span>
                                            </div>
                                        </div>

                                        <div class="space-y-4">
                                            <div class="flex items-center gap-2">
                                                <UToggle v-model="state.quickButtons.mathCalculator" />
                                                <span class="text-sm font-medium">{{
                                                    $t('components.mathcalculator.title')
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
                                            />
                                        </template>

                                        <div class="mt-2 flex w-full items-center justify-center">
                                            <GenericImg
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
                        <div v-if="loading" class="space-y-1 px-4">
                            <USkeleton v-for="idx in 10" :key="idx" class="h-20 w-full" />
                        </div>

                        <UDashboardPanelContent v-else-if="jobProps.discordSyncSettings">
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
                                    name="dryRun"
                                    :label="$t('components.rector.job_props.discord_sync_settings.dry_run')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.discordSyncSettings.dryRun" />
                                </UFormGroup>

                                <UFormGroup
                                    name="statusLog"
                                    :label="$t('components.rector.job_props.discord_sync_settings.status_log')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.discordSyncSettings.statusLog" />
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
                                        name="discordSyncSettings.statusLogSettings.channelId"
                                        :disabled="!state.discordSyncSettings.statusLog"
                                        :placeholder="
                                            $t(
                                                'components.rector.job_props.discord_sync_settings.status_log_settings.channel_id',
                                            )
                                        "
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
                                        <ClientOnly>
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
                                            >
                                                <template #label>
                                                    {{
                                                        $t(
                                                            `enums.rector.UserInfoSyncUnemployedMode.${
                                                                UserInfoSyncUnemployedMode[
                                                                    state.discordSyncSettings.userInfoSyncSettings
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
                                        </ClientOnly>
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
                                        />
                                    </UFormGroup>

                                    <UFormGroup
                                        name="syncNicknames"
                                        :label="
                                            $t(
                                                'components.rector.job_props.discord_sync_settings.user_info_sync_settings.sync_nicknames',
                                            )
                                        "
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UToggle v-model="state.discordSyncSettings.userInfoSyncSettings.syncNicknames" />
                                    </UFormGroup>

                                    <UFormGroup
                                        name="discordSyncSettings.userInfoSyncSettings.groupMapping"
                                        :label="
                                            $t(
                                                'components.rector.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.title',
                                            )
                                        "
                                        :description="
                                            $t(
                                                'components.rector.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.description',
                                            )
                                        "
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <div class="flex flex-col gap-1">
                                            <div
                                                v-for="(_, idx) in state.discordSyncSettings.userInfoSyncSettings.groupMapping"
                                                :key="idx"
                                                class="flex items-center gap-1"
                                            >
                                                <div class="flex flex-col gap-1">
                                                    <UFormGroup
                                                        :name="`discordSyncSettings.userInfoSyncSettings.groupMapping.${idx}.name`"
                                                        :label="
                                                            $t(
                                                                'components.rector.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.name',
                                                            )
                                                        "
                                                        class="flex-1"
                                                    >
                                                        <UInput
                                                            v-model="
                                                                state.discordSyncSettings.userInfoSyncSettings.groupMapping[
                                                                    idx
                                                                ]!.name
                                                            "
                                                            :name="`userInfoSyncSettings.${idx}.name`"
                                                            type="text"
                                                            class="w-full"
                                                            :disabled="!state.discordSyncSettings.userInfoSync"
                                                            :placeholder="
                                                                $t(
                                                                    'components.rector.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.name',
                                                                )
                                                            "
                                                        />
                                                    </UFormGroup>

                                                    <div class="flex flex-row gap-1">
                                                        <UFormGroup
                                                            :name="`discordSyncSettings.userInfoSyncSettings.groupMapping.${idx}.fromGrade`"
                                                            :label="
                                                                $t(
                                                                    'components.rector.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.from_grade',
                                                                )
                                                            "
                                                            class="flex-1"
                                                        >
                                                            <UInput
                                                                v-model="
                                                                    state.discordSyncSettings.userInfoSyncSettings.groupMapping[
                                                                        idx
                                                                    ]!.fromGrade
                                                                "
                                                                :name="`discordSyncSettings.userInfoSyncSettings.${idx}.fromGrade`"
                                                                type="number"
                                                                class="w-full"
                                                                :disabled="!state.discordSyncSettings.userInfoSync"
                                                                :placeholder="
                                                                    $t(
                                                                        'components.rector.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.from_grade',
                                                                    )
                                                                "
                                                            />
                                                        </UFormGroup>
                                                        <UFormGroup
                                                            :name="`discordSyncSettings.userInfoSyncSettings.groupMapping.${idx}.toGrade`"
                                                            :label="
                                                                $t(
                                                                    'components.rector.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.to_grade',
                                                                )
                                                            "
                                                            class="flex-1"
                                                        >
                                                            <UInput
                                                                v-model="
                                                                    state.discordSyncSettings.userInfoSyncSettings.groupMapping[
                                                                        idx
                                                                    ]!.toGrade
                                                                "
                                                                :name="`userInfoSyncSettings.${idx}.toGrade`"
                                                                type="number"
                                                                class="w-full"
                                                                :disabled="!state.discordSyncSettings.userInfoSync"
                                                                :placeholder="
                                                                    $t(
                                                                        'components.rector.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.to_grade',
                                                                    )
                                                                "
                                                            />
                                                        </UFormGroup>
                                                    </div>
                                                </div>

                                                <UButton
                                                    :ui="{ rounded: 'rounded-full' }"
                                                    :disabled="!canSubmit"
                                                    icon="i-mdi-close"
                                                    @click="
                                                        state.discordSyncSettings?.userInfoSyncSettings.groupMapping.splice(
                                                            idx,
                                                            1,
                                                        )
                                                    "
                                                />
                                            </div>
                                        </div>

                                        <UButton
                                            :ui="{ rounded: 'rounded-full' }"
                                            :disabled="!canSubmit"
                                            :class="
                                                state.discordSyncSettings?.userInfoSyncSettings.groupMapping.length
                                                    ? 'mt-2'
                                                    : ''
                                            "
                                            icon="i-mdi-plus"
                                            @click="
                                                state.discordSyncSettings?.userInfoSyncSettings.groupMapping.push({
                                                    name: '',
                                                    fromGrade: appConfig.game.startJobGrade,
                                                    toGrade: appConfig.game.startJobGrade,
                                                })
                                            "
                                        />
                                    </UFormGroup>

                                    <UFormGroup
                                        name="discordSyncSettings.jobsAbsence"
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
                                                name="discordSyncSettings.jobsAbsenceSettings.absenceRole"
                                                :disabled="
                                                    !state.discordSyncSettings.userInfoSync ||
                                                    !state.discordSyncSettings.jobsAbsence
                                                "
                                                :placeholder="
                                                    $t(
                                                        'components.rector.job_props.discord_sync_settings.jobs_absence_settings.jobs_absence_role_name',
                                                    )
                                                "
                                            />
                                        </UFormGroup>
                                    </template>

                                    <UFormGroup
                                        name="discordSyncSettings.qualificationsRoleFormat"
                                        :label="
                                            $t(
                                                'components.rector.job_props.discord_sync_settings.qualifications_role_format.title',
                                            )
                                        "
                                        :description="
                                            $t(
                                                'components.rector.job_props.discord_sync_settings.qualifications_role_format.description',
                                            )
                                        "
                                        class="grid grid-cols-2 items-center gap-2"
                                        :ui="{ container: '' }"
                                    >
                                        <UInput
                                            v-model="state.discordSyncSettings.qualificationsRoleFormat"
                                            type="text"
                                            name="discordSyncSettings.qualificationsRoleFormat"
                                            :placeholder="
                                                $t(
                                                    'components.rector.job_props.discord_sync_settings.qualifications_role_format.title',
                                                )
                                            "
                                        />
                                    </UFormGroup>
                                </template>
                            </UDashboardSection>

                            <UDashboardSection v-if="jobProps.discordSyncChanges">
                                <UAccordion
                                    :items="[{ label: $t('common.diff'), slot: 'diff', icon: 'i-mdi-difference-left' }]"
                                >
                                    <template #diff>
                                        <ClientOnly>
                                            <USelectMenu
                                                v-model="selectedChange"
                                                :options="jobProps.discordSyncChanges.changes"
                                                :searchable-placeholder="$t('common.search_field')"
                                            >
                                                <template #label>
                                                    <span class="truncate">{{
                                                        $d(toDate(selectedChange?.time), 'short')
                                                    }}</span>
                                                </template>
                                                <template #option="{ option }">
                                                    <span class="truncate">{{ $d(toDate(option.time), 'short') }}</span>
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>

                                        <CodeDiff
                                            class="codediff"
                                            theme="dark"
                                            :old-string="jobProps.discordSyncChanges.changes[0]?.plan ?? ''"
                                            :new-string="selectedChange?.plan ?? ''"
                                            :filename="$d(toDate(jobProps.discordSyncChanges.changes[0]?.time), 'short')"
                                            :new-filename="$d(toDate(selectedChange?.time), 'short')"
                                            output-format="side-by-side"
                                            hide-stat
                                            trim
                                        />
                                    </template>
                                </UAccordion>
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
                                            <UFormGroup
                                                :name="`discordSyncSettings.groupSyncSettings.ignoredRoleIds.${idx}.name`"
                                                class="flex-1"
                                            >
                                                <UInput
                                                    v-model="state.discordSyncSettings.groupSyncSettings.ignoredRoleIds[idx]"
                                                    :name="`groupSyncSettingsIgnoredIds.${idx}`"
                                                    type="text"
                                                    class="w-full"
                                                    :placeholder="
                                                        $t(
                                                            'components.rector.job_props.discord_sync_settings.group_sync_settings.ignored_role_ids.field',
                                                        )
                                                    "
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
        </UDashboardPanelContent>
    </UForm>
</template>

<style scoped>
.codediff:deep(.diff-commandbar) {
    display: none;
}
</style>
