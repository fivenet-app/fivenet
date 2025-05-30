<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { vMaska } from 'maska/vue';
import { CodeDiff } from 'v-code-diff';
import { z } from 'zod';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import { useAuthStore } from '~/stores/auth';
import { useNotificatorStore } from '~/stores/notificator';
import { useSettingsStore } from '~/stores/settings';
import type { JobProps } from '~~/gen/ts/resources/jobs/job_props';
import { type DiscordSyncChange, UserInfoSyncUnemployedMode } from '~~/gen/ts/resources/jobs/job_settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import FileUpload from '../partials/elements/FileUpload.vue';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { can } = useAuth();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const appConfig = useAppConfig();

const authStore = useAuthStore();

const notifications = useNotificatorStore();

const schema = z.object({
    livemapMarkerColor: z.string().length(7),
    quickButtons: z.object({
        penaltyCalculator: z.boolean(),
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
    settings: z.object({
        absencePastDays: z.number().int().nonnegative().min(0).max(31),
        absenceFutureDays: z.number().int().nonnegative().min(0).max(186),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    livemapMarkerColor: '',
    quickButtons: {
        penaltyCalculator: false,
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
    settings: {
        absencePastDays: 7,
        absenceFutureDays: 93,
    },
});

async function getJobProps(): Promise<JobProps> {
    try {
        const call = $grpc.settings.settings.getJobProps({});
        const { response } = await call;

        return response.jobProps!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { data: jobProps, pending: loading, refresh, error } = useLazyAsyncData(`settings-jobprops`, () => getJobProps());

async function setJobProps(values: Schema): Promise<void> {
    if (!jobProps.value) {
        return;
    }

    jobProps.value.livemapMarkerColor = values.livemapMarkerColor;
    jobProps.value.quickButtons = values.quickButtons;
    jobProps.value.radioFrequency = values.radioFrequency;
    jobProps.value.discordGuildId = values.discordGuildId.trim().length > 0 ? values.discordGuildId : undefined;
    jobProps.value.discordSyncSettings = values.discordSyncSettings;
    if (!jobProps.value.settings) {
        jobProps.value.settings = {
            absencePastDays: 7,
            absenceFutureDays: 93,
        };
    }
    jobProps.value.settings.absencePastDays = values.settings.absencePastDays;
    jobProps.value.settings.absenceFutureDays = values.settings.absenceFutureDays;

    try {
        const { response } = await $grpc.settings.settings.setJobProps({
            jobProps: jobProps.value,
        });

        notifications.add({
            title: { key: 'notifications.settings.job_props.title', parameters: {} },
            description: { key: 'notifications.settings.job_props.content', parameters: {} },
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

        if (jobProps.value.settings) {
            state.settings = jobProps.value.settings;
        }
    }
}

watch(jobProps, () => setSettingsValues());

const canEdit = can('settings.SettingsService.SetJobProps');

const items = [
    {
        slot: 'jobprops',
        label: t('components.settings.job_props.job_properties'),
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
        <UDashboardNavbar :title="$t('pages.settings.settings.title')">
            <template #right>
                <PartialsBackButton fallback-to="/settings" />
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent>
            <StreamerModeAlert />
        </UDashboardPanelContent>
    </template>
    <UForm
        v-else
        class="min-h-dscreen flex w-full max-w-full flex-1 flex-col overflow-y-auto"
        :schema="schema"
        :state="state"
        @submit="onSubmitThrottle"
    >
        <UDashboardNavbar :title="$t('components.settings.job_props.job_properties')">
            <template #right>
                <PartialsBackButton fallback-to="/settings" />

                <UButton
                    v-if="!!jobProps && canEdit"
                    type="submit"
                    trailing-icon="i-mdi-content-save"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                >
                    {{ $t('common.save', 1) }}
                </UButton>
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent class="p-0 sm:pb-0">
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('components.settings.job_props.job_properties')])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="!jobProps"
                icon="i-mdi-tune"
                :type="$t('components.settings.job_props.job_properties')"
            />

            <template v-else-if="loading || jobProps">
                <UTabs v-model="selectedTab" class="w-full" :items="items" :ui="{ list: { rounded: '' } }">
                    <template #jobprops>
                        <div v-if="loading" class="space-y-1 px-4">
                            <USkeleton v-for="idx in 5" :key="idx" class="h-20 w-full" />
                        </div>

                        <UDashboardPanelContent v-else>
                            <UDashboardSection
                                :title="$t('components.settings.job_props.job_properties')"
                                :description="$t('components.settings.job_props.your_job_properties')"
                            >
                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="logoFile"
                                    :label="$t('common.logo')"
                                    :ui="{ container: '' }"
                                >
                                    <FileUpload
                                        v-model="jobProps.logoFile"
                                        :disabled="!canSubmit || !canEdit"
                                        :upload-fn="(opts) => $grpc.settings.settings.uploadJobLogo(opts)"
                                        :delete-fn="() => $grpc.settings.settings.deleteJobLogo({})"
                                    />
                                </UFormGroup>

                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="livemapMarkerColor"
                                    :label="$t('components.settings.job_props.livemap_marker_color')"
                                    :ui="{ container: '' }"
                                >
                                    <ColorPickerClient v-model="state.livemapMarkerColor" :disabled="!canSubmit || !canEdit" />
                                </UFormGroup>

                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="radioFrequency"
                                    :label="$t('common.radio_frequency')"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.radioFrequency"
                                        v-maska
                                        type="text"
                                        :disabled="!canSubmit || !canEdit"
                                        :placeholder="$t('common.radio_frequency')"
                                        :label="$t('common.radio_frequency')"
                                        data-maska="0.9"
                                        data-maska-tokens="0:\d:multiple|9:\d:multiple"
                                    />
                                </UFormGroup>

                                <UFormGroup
                                    v-if="jobProps.quickButtons"
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="quickButtons"
                                    :label="$t('components.settings.job_props.quick_buttons')"
                                    :ui="{ container: '' }"
                                >
                                    <div class="flex flex-col gap-2">
                                        <div class="space-y-4">
                                            <div class="flex items-center gap-2">
                                                <UToggle
                                                    v-model="state.quickButtons.penaltyCalculator"
                                                    :disabled="!canSubmit || !canEdit"
                                                />
                                                <span class="text-sm font-medium">{{
                                                    $t('components.penaltycalculator.title')
                                                }}</span>
                                            </div>
                                        </div>

                                        <div class="space-y-4">
                                            <div class="flex items-center gap-2">
                                                <UToggle
                                                    v-model="state.quickButtons.mathCalculator"
                                                    :disabled="!canSubmit || !canEdit"
                                                />
                                                <span class="text-sm font-medium">{{
                                                    $t('components.mathcalculator.title')
                                                }}</span>
                                            </div>
                                        </div>
                                    </div>
                                </UFormGroup>
                            </UDashboardSection>

                            <UDivider class="mb-4" />

                            <UDashboardSection :title="$t('components.settings.job_props.settings.absence.title')">
                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="settings.absencePastDays"
                                    :label="$t('components.settings.job_props.settings.absence.past_days')"
                                    :ui="{ container: '' }"
                                >
                                    <div class="flex items-center gap-1">
                                        <UInput
                                            v-model="state.settings.absencePastDays"
                                            class="flex-1"
                                            type="number"
                                            :disabled="!canSubmit || !canEdit"
                                            :min="0"
                                            :placeholder="$t('common.day', 2)"
                                            :label="$t('common.day', 2)"
                                        />
                                        <span>{{ $t('common.day', 2) }}</span>
                                    </div>
                                </UFormGroup>

                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="settings.absenceFutureDays"
                                    :label="$t('components.settings.job_props.settings.absence.future_days')"
                                    :ui="{ container: '' }"
                                >
                                    <div class="flex items-center gap-1">
                                        <UInput
                                            v-model="state.settings.absenceFutureDays"
                                            class="flex-1"
                                            type="number"
                                            :disabled="!canSubmit || !canEdit"
                                            :min="7"
                                            :max="186"
                                            :placeholder="$t('common.day', 2)"
                                            :label="$t('common.day', 2)"
                                        />
                                        <span>{{ $t('common.day', 2) }}</span>
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
                                :title="$t('components.settings.job_props.discord_sync_settings.title')"
                                :description="$t('components.settings.job_props.discord_sync_settings.subtitle')"
                            >
                                <template #links>
                                    <UButton
                                        v-if="appConfig.discord.botInviteURL !== undefined"
                                        class="mt-1"
                                        block
                                        color="white"
                                        trailing-icon="i-mdi-robot"
                                        :disabled="!canSubmit || !canEdit"
                                        :to="appConfig.discord.botInviteURL"
                                        :external="true"
                                    >
                                        {{ $t('components.settings.job_props.discord_sync_settings.invite_bot') }}
                                    </UButton>
                                </template>

                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordGuildId"
                                    :label="$t('components.settings.job_props.discord_sync_settings.discord_guild_id')"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.discordGuildId"
                                        type="text"
                                        :disabled="appConfig.discord.botInviteURL === undefined || !canSubmit || !canEdit"
                                        :placeholder="
                                            $t('components.settings.job_props.discord_sync_settings.discord_guild_id')
                                        "
                                        maxlength="70"
                                    />
                                    <p v-if="jobProps.discordLastSync" class="mt-2 text-xs">
                                        {{ $t('components.settings.job_props.discord_sync_settings.last_sync') }}:
                                        <GenericTime :value="jobProps.discordLastSync" />
                                    </p>
                                </UFormGroup>

                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordSyncSettings.dryRun"
                                    :label="$t('components.settings.job_props.discord_sync_settings.dry_run')"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.discordSyncSettings.dryRun" :disabled="!canSubmit || !canEdit" />
                                </UFormGroup>

                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordSyncSettings.statusLog"
                                    :label="$t('components.settings.job_props.discord_sync_settings.status_log')"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.discordSyncSettings.statusLog" :disabled="!canSubmit || !canEdit" />
                                </UFormGroup>

                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordSyncSettings.statusLogSettings.channelId"
                                    :label="`${$t('components.settings.job_props.discord_sync_settings.status_log')} ${$t('components.settings.job_props.discord_sync_settings.status_log_settings.channel_id')}`"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.discordSyncSettings.statusLogSettings!.channelId"
                                        type="text"
                                        name="discordSyncSettings.statusLogSettings.channelId"
                                        :disabled="!state.discordSyncSettings.statusLog || !canSubmit || !canEdit"
                                        :placeholder="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.status_log_settings.channel_id',
                                            )
                                        "
                                    />
                                </UFormGroup>

                                <UAlert
                                    :ui="{
                                        icon: { base: 'size-6' },
                                    }"
                                    icon="i-mdi-information-outline"
                                    :description="
                                        $t(
                                            'components.settings.job_props.discord_sync_settings.discord_sync_command_info.description',
                                        )
                                    "
                                    variant="subtle"
                                />
                            </UDashboardSection>

                            <UDivider class="mb-4" />

                            <UDashboardSection
                                :title="$t('components.settings.job_props.discord_sync_features.title')"
                                :description="$t('components.settings.job_props.discord_sync_features.subtitle')"
                            >
                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordSyncSettings.userInfoSync"
                                    :label="$t('components.settings.job_props.discord_sync_settings.user_info_sync')"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle
                                        v-model="state.discordSyncSettings.userInfoSync"
                                        :disabled="!canSubmit || !canEdit"
                                    />
                                </UFormGroup>

                                <template v-if="jobProps.discordSyncSettings.userInfoSyncSettings">
                                    <UFormGroup
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.gradeRoleFormat"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.grade_role_format',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <UInput
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.gradeRoleFormat"
                                            type="text"
                                            name="gradeRoleFormat"
                                            :disabled="!state.discordSyncSettings.userInfoSync || !canSubmit || !canEdit"
                                            :placeholder="
                                                $t(
                                                    'components.settings.job_props.discord_sync_settings.user_info_sync_settings.grade_role_format',
                                                )
                                            "
                                        />
                                    </UFormGroup>

                                    <UFormGroup
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.employeeRoleEnabled"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.employee_role_enabled',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <UToggle
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.employeeRoleEnabled"
                                            :disabled="!state.discordSyncSettings.userInfoSync || !canSubmit || !canEdit"
                                        />
                                    </UFormGroup>

                                    <UFormGroup
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.employeeRoleFormat"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.employee_role_format',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <UInput
                                            v-model="state.discordSyncSettings.userInfoSyncSettings!.employeeRoleFormat"
                                            type="text"
                                            name="employeeRoleFormat"
                                            :disabled="
                                                !state.discordSyncSettings.userInfoSync ||
                                                !state.discordSyncSettings.userInfoSyncSettings?.employeeRoleEnabled ||
                                                !canSubmit ||
                                                !canEdit
                                            "
                                            :placeholder="
                                                $t(
                                                    'components.settings.job_props.discord_sync_settings.user_info_sync_settings.employee_role_format',
                                                )
                                            "
                                        />
                                    </UFormGroup>

                                    <UFormGroup
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.unemployedEnabled"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.unemployed_enabled',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <UToggle
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.unemployedEnabled"
                                            :disabled="!state.discordSyncSettings.userInfoSync"
                                        />
                                    </UFormGroup>

                                    <UFormGroup
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.unemployedMode"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.unemployed_mode',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <ClientOnly>
                                            <USelectMenu
                                                v-model="state.discordSyncSettings.userInfoSyncSettings.unemployedMode"
                                                :disabled="
                                                    !state.discordSyncSettings.userInfoSync ||
                                                    !state.discordSyncSettings.userInfoSyncSettings.unemployedEnabled ||
                                                    !canSubmit ||
                                                    !canEdit
                                                "
                                                value-attribute="value"
                                                :options="[
                                                    {
                                                        label: $t('enums.settings.UserInfoSyncUnemployedMode.GIVE_ROLE'),
                                                        value: UserInfoSyncUnemployedMode.GIVE_ROLE,
                                                    },
                                                    {
                                                        label: $t('enums.settings.UserInfoSyncUnemployedMode.GIVE_ROLE'),
                                                        value: UserInfoSyncUnemployedMode.KICK,
                                                    },
                                                ]"
                                                :searchable-placeholder="$t('common.search_field')"
                                            >
                                                <template #label>
                                                    {{
                                                        $t(
                                                            `enums.settings.UserInfoSyncUnemployedMode.${
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
                                                            `enums.settings.UserInfoSyncUnemployedMode.${UserInfoSyncUnemployedMode[option.value]}`,
                                                        )
                                                    }}</span>
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>
                                    </UFormGroup>

                                    <UFormGroup
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.unemployedRoleName"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.unemployed_role_name',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <UInput
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.unemployedRoleName"
                                            type="text"
                                            name="unemployedRoleName"
                                            :disabled="
                                                !state.discordSyncSettings.userInfoSync ||
                                                !state.discordSyncSettings.userInfoSyncSettings.unemployedEnabled ||
                                                !canSubmit ||
                                                !canEdit
                                            "
                                            :placeholder="
                                                $t(
                                                    'components.settings.job_props.discord_sync_settings.user_info_sync_settings.unemployed_role_name',
                                                )
                                            "
                                        />
                                    </UFormGroup>

                                    <UFormGroup
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.syncNicknames"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.sync_nicknames',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <UToggle
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.syncNicknames"
                                            :disabled="!canSubmit || !canEdit"
                                        />
                                    </UFormGroup>

                                    <UFormGroup
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.groupMapping"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.title',
                                            )
                                        "
                                        :description="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.description',
                                            )
                                        "
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
                                                        class="flex-1"
                                                        :name="`discordSyncSettings.userInfoSyncSettings.groupMapping.${idx}.name`"
                                                        :label="
                                                            $t(
                                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.name',
                                                            )
                                                        "
                                                    >
                                                        <UInput
                                                            v-model="
                                                                state.discordSyncSettings.userInfoSyncSettings.groupMapping[
                                                                    idx
                                                                ]!.name
                                                            "
                                                            class="w-full"
                                                            :name="`userInfoSyncSettings.${idx}.name`"
                                                            type="text"
                                                            :disabled="
                                                                !state.discordSyncSettings.userInfoSync ||
                                                                !canSubmit ||
                                                                !canEdit
                                                            "
                                                            :placeholder="
                                                                $t(
                                                                    'components.settings.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.name',
                                                                )
                                                            "
                                                        />
                                                    </UFormGroup>

                                                    <div class="flex flex-row gap-1">
                                                        <UFormGroup
                                                            class="flex-1"
                                                            :name="`discordSyncSettings.userInfoSyncSettings.groupMapping.${idx}.fromGrade`"
                                                            :label="
                                                                $t(
                                                                    'components.settings.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.from_grade',
                                                                )
                                                            "
                                                        >
                                                            <UInput
                                                                v-model="
                                                                    state.discordSyncSettings.userInfoSyncSettings.groupMapping[
                                                                        idx
                                                                    ]!.fromGrade
                                                                "
                                                                class="w-full"
                                                                :name="`discordSyncSettings.userInfoSyncSettings.${idx}.fromGrade`"
                                                                type="number"
                                                                :min="0"
                                                                :disabled="
                                                                    !state.discordSyncSettings.userInfoSync ||
                                                                    !canSubmit ||
                                                                    !canEdit
                                                                "
                                                                :placeholder="
                                                                    $t(
                                                                        'components.settings.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.from_grade',
                                                                    )
                                                                "
                                                            />
                                                        </UFormGroup>
                                                        <UFormGroup
                                                            class="flex-1"
                                                            :name="`discordSyncSettings.userInfoSyncSettings.groupMapping.${idx}.toGrade`"
                                                            :label="
                                                                $t(
                                                                    'components.settings.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.to_grade',
                                                                )
                                                            "
                                                        >
                                                            <UInput
                                                                v-model="
                                                                    state.discordSyncSettings.userInfoSyncSettings.groupMapping[
                                                                        idx
                                                                    ]!.toGrade
                                                                "
                                                                class="w-full"
                                                                :name="`userInfoSyncSettings.${idx}.toGrade`"
                                                                type="number"
                                                                :min="0"
                                                                :disabled="
                                                                    !state.discordSyncSettings.userInfoSync ||
                                                                    !canSubmit ||
                                                                    !canEdit
                                                                "
                                                                :placeholder="
                                                                    $t(
                                                                        'components.settings.job_props.discord_sync_settings.user_info_sync_settings.group_mapping.to_grade',
                                                                    )
                                                                "
                                                            />
                                                        </UFormGroup>
                                                    </div>
                                                </div>

                                                <UButton
                                                    v-if="canEdit"
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
                                            v-if="canEdit"
                                            :class="
                                                state.discordSyncSettings?.userInfoSyncSettings.groupMapping.length
                                                    ? 'mt-2'
                                                    : ''
                                            "
                                            :ui="{ rounded: 'rounded-full' }"
                                            :disabled="!canSubmit"
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
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.jobsAbsence"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.jobs_absence_settings.jobs_absence_role_enabled',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <UToggle
                                            v-model="state.discordSyncSettings.jobsAbsence"
                                            :disabled="!state.discordSyncSettings.userInfoSync || !canSubmit || !canEdit"
                                        />
                                    </UFormGroup>

                                    <template v-if="jobProps.discordSyncSettings.jobsAbsenceSettings">
                                        <UFormGroup
                                            class="grid grid-cols-2 items-center gap-2"
                                            name="discordSyncSettings.jobsAbsenceSettings.absenceRole"
                                            :label="
                                                $t(
                                                    'components.settings.job_props.discord_sync_settings.jobs_absence_settings.jobs_absence_role_name',
                                                )
                                            "
                                            :ui="{ container: '' }"
                                        >
                                            <UInput
                                                v-model="state.discordSyncSettings.jobsAbsenceSettings.absenceRole"
                                                type="text"
                                                name="discordSyncSettings.jobsAbsenceSettings.absenceRole"
                                                :disabled="
                                                    !state.discordSyncSettings.userInfoSync ||
                                                    !state.discordSyncSettings.jobsAbsence ||
                                                    !canSubmit ||
                                                    !canEdit
                                                "
                                                :placeholder="
                                                    $t(
                                                        'components.settings.job_props.discord_sync_settings.jobs_absence_settings.jobs_absence_role_name',
                                                    )
                                                "
                                            />
                                        </UFormGroup>
                                    </template>

                                    <UFormGroup
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.qualificationsRoleFormat"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.qualifications_role_format.title',
                                            )
                                        "
                                        :description="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.qualifications_role_format.description',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <UInput
                                            v-model="state.discordSyncSettings.qualificationsRoleFormat"
                                            type="text"
                                            name="discordSyncSettings.qualificationsRoleFormat"
                                            :disabled="!canSubmit || !canEdit"
                                            :placeholder="
                                                $t(
                                                    'components.settings.job_props.discord_sync_settings.qualifications_role_format.title',
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
                                :title="$t('components.settings.job_props.discord_sync_settings.group_sync_settings.title')"
                                :description="
                                    $t('components.settings.job_props.discord_sync_settings.group_sync_settings.subtitle')
                                "
                            >
                                <UFormGroup
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="groupSyncSettingsIgnoredIds"
                                    :label="
                                        $t(
                                            'components.settings.job_props.discord_sync_settings.group_sync_settings.ignored_role_ids.title',
                                        )
                                    "
                                    :description="
                                        $t(
                                            'components.settings.job_props.discord_sync_settings.group_sync_settings.ignored_role_ids.description',
                                        )
                                    "
                                    :ui="{ container: '' }"
                                >
                                    <div class="flex flex-col gap-1">
                                        <div
                                            v-for="(_, idx) in state.discordSyncSettings.groupSyncSettings.ignoredRoleIds"
                                            :key="idx"
                                            class="flex items-center gap-1"
                                        >
                                            <UFormGroup
                                                class="flex-1"
                                                :name="`discordSyncSettings.groupSyncSettings.ignoredRoleIds.${idx}.name`"
                                            >
                                                <UInput
                                                    v-model="state.discordSyncSettings.groupSyncSettings.ignoredRoleIds[idx]"
                                                    class="w-full"
                                                    :name="`groupSyncSettingsIgnoredIds.${idx}`"
                                                    type="text"
                                                    :disabled="!canSubmit || !canEdit"
                                                    :placeholder="
                                                        $t(
                                                            'components.settings.job_props.discord_sync_settings.group_sync_settings.ignored_role_ids.field',
                                                        )
                                                    "
                                                />
                                            </UFormGroup>

                                            <UButton
                                                v-if="canEdit"
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
                                        v-if="canEdit"
                                        :class="
                                            state.discordSyncSettings?.groupSyncSettings.ignoredRoleIds.length ? 'mt-2' : ''
                                        "
                                        :ui="{ rounded: 'rounded-full' }"
                                        :disabled="!canSubmit"
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
