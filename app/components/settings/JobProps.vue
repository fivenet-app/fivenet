<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { vMaska } from 'maska/vue';
import { CodeDiff } from 'v-code-diff';
import { z } from 'zod';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import { useAuthStore } from '~/stores/auth';
import { useSettingsStore } from '~/stores/settings';
import { getSettingsSettingsClient } from '~~/gen/ts/clients';
import type { JobProps } from '~~/gen/ts/resources/jobs/job_props';
import { type DiscordSyncChange, UserInfoSyncUnemployedMode } from '~~/gen/ts/resources/jobs/job_settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import FileUpload from '../partials/elements/FileUpload.vue';
import FormatBuilder from '../partials/FormatBuilder.vue';
import NotSupportedTabletBlock from '../partials/NotSupportedTabletBlock.vue';

const { t } = useI18n();

const { can } = useAuth();

const settingsStore = useSettingsStore();
const { streamerMode, nuiEnabled } = storeToRefs(settingsStore);

const appConfig = useAppConfig();

const authStore = useAuthStore();

const notifications = useNotificationsStore();

const settingsSettingsClient = await getSettingsSettingsClient();

const schema = z.object({
    livemapMarkerColor: z.string().length(7),
    quickButtons: z.object({
        penaltyCalculator: z.coerce.boolean(),
        mathCalculator: z.coerce.boolean(),
    }),
    radioFrequency: z.string().max(24),
    discordGuildId: z.string().max(48),
    discordSyncSettings: z.object({
        dryRun: z.coerce.boolean(),
        userInfoSync: z.coerce.boolean(),
        userInfoSyncSettings: z.object({
            employeeRoleEnabled: z.coerce.boolean(),
            employeeRoleFormat: z.string().max(64),
            gradeRoleFormat: z.string().max(64),
            unemployedEnabled: z.coerce.boolean(),
            unemployedMode: z.nativeEnum(UserInfoSyncUnemployedMode),
            unemployedRoleName: z.string().max(64),
            syncNicknames: z.coerce.boolean(),
            groupMapping: z
                .object({
                    name: z.string().max(64),
                    fromGrade: z.coerce.number().min(0).max(99999),
                    toGrade: z.coerce.number().min(0).max(99999),
                })
                .array()
                .max(25)
                .default([]),
        }),
        statusLog: z.coerce.boolean(),
        statusLogSettings: z.object({
            channelId: z.string().max(64),
        }),
        jobsAbsence: z.coerce.boolean(),
        jobsAbsenceSettings: z.object({
            absenceRole: z.string().max(64),
        }),
        groupSyncSettings: z.object({
            ignoredRoleIds: z.string().max(64).array().max(20).default([]),
        }),
        qualificationsRoleFormat: z.string().max(64),
    }),
    settings: z.object({
        absencePastDays: z.coerce.number().int().nonnegative().min(0).max(31),
        absenceFutureDays: z.coerce.number().int().nonnegative().min(0).max(186),
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
        const call = settingsSettingsClient.getJobProps({});
        const { response } = await call;

        return response.jobProps!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { data: jobProps, status, refresh, error } = useLazyAsyncData(`settings-jobprops`, () => getJobProps());

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
        const { response } = await settingsSettingsClient.setJobProps({
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

            state.discordSyncSettings.userInfoSyncSettings.employeeRoleFormat =
                state.discordSyncSettings.userInfoSyncSettings.employeeRoleFormat.replaceAll('%s', '%job%');
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

const canEdit = can('settings.SettingsService/SetJobProps');

const dcConnectRequired = ref(false);
const { data: userGuilds } = useLazyAsyncData(`settings-userguilds`, () => listGuilds(), {
    immediate: appConfig.discord.botEnabled,
});

async function listGuilds() {
    if (!canEdit.value) return [];

    try {
        const call = settingsSettingsClient.listUserGuilds({});
        const { response } = await call;

        return response.guilds;
    } catch (e) {
        handleGRPCError(e as RpcError);
        if ((e as Error).message.includes('ErrDiscordConnectRequired')) {
            dcConnectRequired.value = true;
        } else {
            dcConnectRequired.value = false;
        }

        return [];
    }
}

async function searchChannels() {
    if (!canEdit.value) return [];

    try {
        const call = settingsSettingsClient.listDiscordChannels({});
        const { response } = await call;

        return response.channels.sort((a, b) => a.position - b.position);
    } catch (e) {
        handleGRPCError(e as RpcError);
        return [];
    }
}

const items = [
    {
        slot: 'jobprops' as const,
        label: t('components.settings.job_props.job_properties'),
        icon: 'i-mdi-settings',
        value: 'jobprops',
    },
    {
        slot: 'discord' as const,
        label: t('common.discord'),
        icon: 'i-simple-icons-discord',
        value: 'discord',
        disabled: !appConfig.discord.botEnabled,
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'jobprops';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
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
        class="flex min-h-dvh w-full max-w-full flex-1 flex-col overflow-y-auto"
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
                :retry="refresh"
            />

            <template v-else-if="isRequestPending(status) || jobProps">
                <UTabs v-model="selectedTab" class="w-full" :items="items">
                    <template #jobprops>
                        <div v-if="isRequestPending(status)" class="space-y-1 px-4">
                            <USkeleton v-for="idx in 5" :key="idx" class="h-20 w-full" />
                        </div>

                        <UDashboardPanelContent v-else>
                            <UPageCard
                                :title="$t('components.settings.job_props.job_properties')"
                                :description="$t('components.settings.job_props.your_job_properties')"
                            >
                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="logoFile"
                                    :label="$t('common.logo')"
                                    :ui="{ container: '' }"
                                >
                                    <FileUpload
                                        v-model="jobProps.logoFile"
                                        :disabled="!canSubmit || !canEdit"
                                        :upload-fn="(opts) => settingsSettingsClient.uploadJobLogo(opts)"
                                        :delete-fn="() => settingsSettingsClient.deleteJobLogo({})"
                                    />
                                </UFormField>

                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="livemapMarkerColor"
                                    :label="$t('components.settings.job_props.livemap_marker_color')"
                                    :ui="{ container: '' }"
                                >
                                    <ColorPickerClient v-model="state.livemapMarkerColor" :disabled="!canSubmit || !canEdit" />
                                </UFormField>

                                <UFormField
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
                                </UFormField>

                                <UFormField
                                    v-if="jobProps.quickButtons"
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="quickButtons"
                                    :label="$t('components.settings.job_props.quick_buttons')"
                                    :ui="{ container: '' }"
                                >
                                    <div class="flex flex-col gap-2">
                                        <div class="space-y-4">
                                            <div class="flex items-center gap-2">
                                                <USwitch
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
                                                <USwitch
                                                    v-model="state.quickButtons.mathCalculator"
                                                    :disabled="!canSubmit || !canEdit"
                                                />
                                                <span class="text-sm font-medium">{{
                                                    $t('components.mathcalculator.title')
                                                }}</span>
                                            </div>
                                        </div>
                                    </div>
                                </UFormField>
                            </UPageCard>

                            <USeparator class="mb-4" />

                            <UPageCard :title="$t('components.settings.job_props.settings.absence.title')">
                                <UFormField
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
                                </UFormField>

                                <UFormField
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
                                </UFormField>
                            </UPageCard>
                        </UDashboardPanelContent>
                    </template>

                    <template v-if="appConfig.discord.botEnabled" #discord>
                        <div v-if="isRequestPending(status)" class="space-y-1 px-4">
                            <USkeleton v-for="idx in 10" :key="idx" class="h-20 w-full" />
                        </div>

                        <UDashboardPanelContent v-else-if="jobProps.discordSyncSettings">
                            <UPageCard
                                :title="$t('components.settings.job_props.discord_sync_settings.title')"
                                :description="$t('components.settings.job_props.discord_sync_settings.subtitle')"
                            >
                                <template v-if="appConfig.discord.botEnabled" #links>
                                    <NotSupportedTabletBlock v-if="nuiEnabled" />
                                    <UButton
                                        v-else-if="canEdit"
                                        class="mt-1"
                                        block
                                        color="neutral"
                                        trailing-icon="i-mdi-robot"
                                        :disabled="!canSubmit || !canEdit"
                                        to="/api/discord/invite-bot"
                                        external
                                    >
                                        {{ $t('components.settings.job_props.discord_sync_settings.invite_bot') }}
                                    </UButton>
                                </template>

                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordGuildId"
                                    :label="$t('components.settings.job_props.discord_sync_settings.discord_guild_id')"
                                    :ui="{ container: '' }"
                                >
                                    <template v-if="nuiEnabled || dcConnectRequired">
                                        <NotSupportedTabletBlock v-if="nuiEnabled" />
                                        <UButton
                                            v-else-if="dcConnectRequired"
                                            icon="i-mdi-connection"
                                            :label="$t('common.connect')"
                                            @click="
                                                async () =>
                                                    await navigateTo(
                                                        generateDiscordConnectURL('discord', '/settings/props?tab=discord#'),
                                                        { external: true },
                                                    )
                                            "
                                        />

                                        <div v-if="state.discordGuildId" class="mt-2">{{ state.discordGuildId }}</div>
                                    </template>
                                    <USelectMenu
                                        v-else
                                        v-model="state.discordGuildId"
                                        :items="userGuilds"
                                        searchable
                                        :search-attributes="['name', 'id']"
                                        :placeholder="
                                            $t('components.settings.job_props.discord_sync_settings.discord_guild_id')
                                        "
                                        value-key="id"
                                        :disabled="!canSubmit || !canEdit || userGuilds?.length === 0"
                                        size="lg"
                                    >
                                        <template #item-label="{ item }">
                                            <div class="inline-flex items-center gap-2">
                                                <UAvatar :src="item?.icon" :alt="item?.name" />
                                                <span class="truncate">{{
                                                    item?.name ??
                                                    (state.discordGuildId !== '' ? state.discordGuildId : '&nbsp;')
                                                }}</span>
                                            </div>
                                        </template>

                                        <template #item="{ option }">
                                            <div class="inline-flex items-center gap-2">
                                                <UAvatar :src="option.icon" :alt="option.name" />
                                                <span class="truncate">{{ option.name }}</span>
                                            </div>
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.discord_guild', 2)]) }}
                                        </template>
                                    </USelectMenu>
                                    <p v-if="jobProps.discordLastSync" class="mt-2 text-xs">
                                        {{ $t('components.settings.job_props.discord_sync_settings.last_sync') }}:
                                        <GenericTime :value="jobProps.discordLastSync" />
                                    </p>
                                </UFormField>

                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordSyncSettings.dryRun"
                                    :label="$t('components.settings.job_props.discord_sync_settings.dry_run')"
                                    :ui="{ container: '' }"
                                >
                                    <USwitch v-model="state.discordSyncSettings.dryRun" :disabled="!canSubmit || !canEdit" />
                                </UFormField>

                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordSyncSettings.statusLog"
                                    :label="$t('components.settings.job_props.discord_sync_settings.status_log')"
                                    :ui="{ container: '' }"
                                >
                                    <USwitch v-model="state.discordSyncSettings.statusLog" :disabled="!canSubmit || !canEdit" />
                                </UFormField>

                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordSyncSettings.statusLogSettings.channelId"
                                    :label="
                                        $t('components.settings.job_props.discord_sync_settings.status_log_settings.channel_id')
                                    "
                                    :ui="{ container: '' }"
                                >
                                    <USelectMenu
                                        v-model="state.discordSyncSettings.statusLogSettings!.channelId"
                                        name="discordSyncSettings.statusLogSettings.channelId"
                                        :disabled="
                                            !state.discordGuildId ||
                                            !state.discordSyncSettings.statusLog ||
                                            !canSubmit ||
                                            !canEdit
                                        "
                                        :searchable="searchChannels"
                                        :search-attributes="['name']"
                                        searchable-lazy
                                        :searchable-placeholder="$t('common.search_field')"
                                        value-key="id"
                                        nullable
                                        :placeholder="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.status_log_settings.channel_id',
                                            )
                                        "
                                    >
                                        <template #item-label="{ item }">
                                            <span class="truncate">{{
                                                item
                                                    ? `${item.name} (${item.id})`
                                                    : state.discordSyncSettings.statusLogSettings!.channelId !== ''
                                                      ? state.discordSyncSettings.statusLogSettings!.channelId
                                                      : '&nbsp;'
                                            }}</span>
                                        </template>

                                        <template #item="{ option }">
                                            <span class="truncate">{{ option.name }} ({{ option.id }})</span>
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.channel', 1)]) }}
                                        </template>
                                    </USelectMenu>
                                </UFormField>

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
                            </UPageCard>

                            <USeparator class="mb-4" />

                            <UPageCard
                                :title="$t('components.settings.job_props.discord_sync_features.title')"
                                :description="$t('components.settings.job_props.discord_sync_features.subtitle')"
                            >
                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordSyncSettings.userInfoSync"
                                    :label="$t('components.settings.job_props.discord_sync_settings.user_info_sync')"
                                    :ui="{ container: '' }"
                                >
                                    <USwitch
                                        v-model="state.discordSyncSettings.userInfoSync"
                                        :disabled="!canSubmit || !canEdit"
                                    />
                                </UFormField>

                                <template v-if="jobProps.discordSyncSettings.userInfoSyncSettings">
                                    <UFormField
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.gradeRoleFormat"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.grade_role_format.title',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <FormatBuilder
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.gradeRoleFormat"
                                            :disabled="!canSubmit || !canEdit"
                                            :extensions="[
                                                {
                                                    label: $t(
                                                        'components.settings.job_props.discord_sync_settings.user_info_sync_settings.grade_role_format.grade_single',
                                                    ),
                                                    value: 'grade_single',
                                                },
                                                {
                                                    label: $t(
                                                        'components.settings.job_props.discord_sync_settings.user_info_sync_settings.grade_role_format.grade',
                                                    ),
                                                    value: 'grade',
                                                },
                                                {
                                                    label: $t(
                                                        'components.settings.job_props.discord_sync_settings.user_info_sync_settings.grade_role_format.grade_label',
                                                    ),
                                                    value: 'grade_label',
                                                },
                                            ]"
                                        />
                                    </UFormField>

                                    <UFormField
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.employeeRoleEnabled"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.employee_role_enabled',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <USwitch
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.employeeRoleEnabled"
                                            :disabled="!state.discordSyncSettings.userInfoSync || !canSubmit || !canEdit"
                                        />
                                    </UFormField>

                                    <UFormField
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.employeeRoleFormat"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.employee_role_format',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <FormatBuilder
                                            v-model="state.discordSyncSettings.userInfoSyncSettings!.employeeRoleFormat"
                                            :disabled="
                                                !state.discordSyncSettings.userInfoSync ||
                                                !state.discordSyncSettings.userInfoSyncSettings?.employeeRoleEnabled ||
                                                !canSubmit ||
                                                !canEdit
                                            "
                                            :extensions="[{ label: $t('common.job_name'), value: 'job' }]"
                                        />
                                    </UFormField>

                                    <UFormField
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.unemployedEnabled"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.unemployed_enabled',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <USwitch
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.unemployedEnabled"
                                            :disabled="!state.discordSyncSettings.userInfoSync"
                                        />
                                    </UFormField>

                                    <UFormField
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
                                                value-key="value"
                                                :items="[
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
                                                <template #item-label>
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

                                                <template #item="{ option }">
                                                    <span class="truncate">{{
                                                        $t(
                                                            `enums.settings.UserInfoSyncUnemployedMode.${UserInfoSyncUnemployedMode[option.value]}`,
                                                        )
                                                    }}</span>
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>
                                    </UFormField>

                                    <UFormField
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
                                    </UFormField>

                                    <UFormField
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.syncNicknames"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.sync_nicknames',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <USwitch
                                            v-model="state.discordSyncSettings.userInfoSyncSettings.syncNicknames"
                                            :disabled="!canSubmit || !canEdit"
                                        />
                                    </UFormField>

                                    <UFormField
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
                                                    <UFormField
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
                                                    </UFormField>

                                                    <div class="flex flex-row gap-1">
                                                        <UFormField
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
                                                        </UFormField>
                                                        <UFormField
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
                                                        </UFormField>
                                                    </div>
                                                </div>

                                                <UButton
                                                    v-if="canEdit"
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
                                    </UFormField>

                                    <UFormField
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.jobsAbsence"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.jobs_absence_settings.jobs_absence_role_enabled',
                                            )
                                        "
                                        :ui="{ container: '' }"
                                    >
                                        <USwitch
                                            v-model="state.discordSyncSettings.jobsAbsence"
                                            :disabled="!state.discordSyncSettings.userInfoSync || !canSubmit || !canEdit"
                                        />
                                    </UFormField>

                                    <template v-if="jobProps.discordSyncSettings.jobsAbsenceSettings">
                                        <UFormField
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
                                        </UFormField>
                                    </template>

                                    <UFormField
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
                                        <FormatBuilder
                                            v-model="state.discordSyncSettings.qualificationsRoleFormat"
                                            :disabled="!canSubmit || !canEdit"
                                            :extensions="[
                                                { label: $t('common.qualification_name'), value: 'name' },
                                                { label: $t('common.abbreviation'), value: 'abbr' },
                                            ]"
                                        />
                                    </UFormField>
                                </template>
                            </UPageCard>

                            <UPageCard v-if="jobProps.discordSyncChanges">
                                <UAccordion
                                    :items="[
                                        { label: $t('common.diff'), slot: 'diff' as const, icon: 'i-mdi-difference-left' },
                                    ]"
                                >
                                    <template #diff>
                                        <ClientOnly>
                                            <USelectMenu
                                                v-model="selectedChange"
                                                :items="jobProps.discordSyncChanges.changes"
                                                :searchable-placeholder="$t('common.search_field')"
                                            >
                                                <template #item-label>
                                                    <span class="truncate">{{
                                                        $d(toDate(selectedChange?.time), 'short')
                                                    }}</span>
                                                </template>

                                                <template #item="{ option }">
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
                            </UPageCard>

                            <UPageCard
                                v-if="jobProps.discordSyncSettings?.groupSyncSettings"
                                :title="$t('components.settings.job_props.discord_sync_settings.group_sync_settings.title')"
                                :description="
                                    $t('components.settings.job_props.discord_sync_settings.group_sync_settings.subtitle')
                                "
                            >
                                <UFormField
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
                                            <UFormField
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
                                            </UFormField>

                                            <UButton
                                                v-if="canEdit"
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
                                        :disabled="!canSubmit"
                                        icon="i-mdi-plus"
                                        @click="state.discordSyncSettings?.groupSyncSettings.ignoredRoleIds.push('')"
                                    />
                                </UFormField>
                            </UPageCard>
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
