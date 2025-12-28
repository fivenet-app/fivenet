<script lang="ts" setup>
import { LazyPartialsCodeDiff } from '#components';
import type { FormSubmitEvent } from '@nuxt/ui';
import { vMaska } from 'maska/vue';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import { useAuthStore } from '~/stores/auth';
import { useSettingsStore } from '~/stores/settings';
import { getSettingsSettingsClient } from '~~/gen/ts/clients';
import type { Guild } from '~~/gen/ts/resources/discord/discord';
import type { JobProps } from '~~/gen/ts/resources/jobs/job_props';
import { type DiscordSyncChange, UserInfoSyncUnemployedMode } from '~~/gen/ts/resources/jobs/job_settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import ConfirmModal from '../partials/ConfirmModal.vue';
import GenericImg from '../partials/elements/GenericImg.vue';
import FormatBuilder from '../partials/FormatBuilder.vue';
import NotSupportedTabletBlock from '../partials/NotSupportedTabletBlock.vue';
import SelectMenu from '../partials/SelectMenu.vue';

const { t } = useI18n();

const overlay = useOverlay();

const { can } = useAuth();

const settingsStore = useSettingsStore();
const { streamerMode, nuiEnabled } = storeToRefs(settingsStore);

const appConfig = useAppConfig();

const authStore = useAuthStore();

const notifications = useNotificationsStore();

const settingsSettingsClient = await getSettingsSettingsClient();

const schema = z.object({
    livemapMarkerColor: z.coerce.string().length(7),
    quickButtons: z.object({
        penaltyCalculator: z.coerce.boolean(),
    }),
    radioFrequency: z.coerce.string().max(24),
    discordGuildId: z.coerce.string().max(48),
    discordSyncSettings: z.object({
        dryRun: z.coerce.boolean(),
        userInfoSync: z.coerce.boolean(),
        userInfoSyncSettings: z.object({
            employeeRoleEnabled: z.coerce.boolean(),
            employeeRoleFormat: z.coerce.string().max(64),
            gradeRoleFormat: z.coerce.string().max(64),
            unemployedEnabled: z.coerce.boolean(),
            unemployedMode: z.enum(UserInfoSyncUnemployedMode),
            unemployedRoleName: z.coerce.string().max(64),
            syncNicknames: z.coerce.boolean(),
            groupMapping: z
                .object({
                    name: z.coerce.string().max(64),
                    fromGrade: z.coerce.number().min(0).max(99999),
                    toGrade: z.coerce.number().min(0).max(99999),
                })
                .array()
                .max(25)
                .default([]),
        }),
        statusLog: z.coerce.boolean(),
        statusLogSettings: z.object({
            channelId: z.coerce.string().max(64),
        }),
        jobsAbsence: z.coerce.boolean(),
        jobsAbsenceSettings: z.object({
            absenceRole: z.coerce.string().max(64),
        }),
        groupSyncSettings: z.object({
            ignoredRoleIds: z.coerce.string().max(64).array().max(20).default([]),
        }),
        qualificationsRoleFormat: z.coerce.string().max(64),
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
    if (!jobProps.value) return;

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
    if (!jobProps.value) return;

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
    transform: (guilds) =>
        guilds.map((guild) => ({
            ...guild,
            avatar: {
                src: guild.icon ? `https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png` : undefined,
                alt: guild.name,
            },
        })) ?? [],
    default: () => [] as Guild[],
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
    if (!userGuilds.value || dcConnectRequired.value) return [];

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

const { resizeAndUpload } = useFileUploader((opts) => settingsSettingsClient.uploadJobLogo(opts), 'jobprops', 0);

const selectedChange = ref<DiscordSyncChange | undefined>();

const formRef = useTemplateRef('formRef');

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (event.submitter?.getAttribute('role') === 'tab') return;

    canSubmit.value = false;
    await setJobProps(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('components.settings.job_props.job_properties')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/settings" />

                    <UButton
                        v-if="!!jobProps && canEdit"
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                        :label="$t('common.save', 1)"
                        @click="() => formRef?.submit()"
                    />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <StreamerModeAlert v-if="streamerMode" />
            <UForm v-else ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
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
                    <UTabs
                        v-model="selectedTab"
                        class="w-full"
                        :items="items"
                        variant="link"
                        :ui="{ content: 'p-4 flex flex-col gap-4 max-w-(--ui-container) mx-auto' }"
                        :unmount-on-hide="false"
                    >
                        <template #jobprops>
                            <div v-if="isRequestPending(status)" class="space-y-1 px-4">
                                <USkeleton v-for="idx in 5" :key="idx" class="h-20 w-full" />
                            </div>

                            <UPageCard
                                :title="$t('components.settings.job_props.job_properties')"
                                :description="$t('components.settings.job_props.your_job_properties')"
                            >
                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="logoFile"
                                    :label="$t('common.logo')"
                                >
                                    <div v-if="jobProps.logoFileId" class="flex w-full flex-1 items-center justify-center">
                                        <GenericImg
                                            :src="`/api/filestore/jobprops/${jobProps.job}`"
                                            :alt="`${jobProps.job} ${$t('common.logo')}`"
                                            class="mb-2 size-full max-h-40 min-h-40 max-w-40"
                                        />
                                    </div>

                                    <NotSupportedTabletBlock v-if="nuiEnabled" />
                                    <div v-else class="flex flex-col gap-2 md:flex-row">
                                        <UFileUpload
                                            class="w-full flex-1 grow"
                                            :disabled="!canSubmit || !canEdit"
                                            :accept="appConfig.fileUpload.types.images.join(',')"
                                            :placeholder="$t('common.image')"
                                            :label="$t('common.file_upload_label')"
                                            :description="$t('common.allowed_file_types')"
                                            @update:model-value="($event) => $event && resizeAndUpload($event)"
                                        />

                                        <UButton
                                            v-if="jobProps.logoFileId"
                                            variant="outline"
                                            color="red"
                                            trailing-icon="i-mdi-clear"
                                            :label="$t('common.clear')"
                                            class="grow-0"
                                            @click="
                                                () => {
                                                    confirmModal.open({
                                                        confirm: () => settingsSettingsClient.deleteJobLogo({}),
                                                    });
                                                }
                                            "
                                        />
                                    </div>
                                </UFormField>

                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="livemapMarkerColor"
                                    :label="$t('components.settings.job_props.livemap_marker_color')"
                                >
                                    <ColorPicker
                                        v-model="state.livemapMarkerColor"
                                        :disabled="!canSubmit || !canEdit"
                                        class="w-full"
                                    />
                                </UFormField>

                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="radioFrequency"
                                    :label="$t('common.radio_frequency')"
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
                                        class="w-full"
                                    />
                                </UFormField>

                                <UFormField
                                    v-if="jobProps.quickButtons"
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="quickButtons"
                                    :label="$t('components.settings.job_props.quick_buttons')"
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
                                    </div>
                                </UFormField>
                            </UPageCard>

                            <UPageCard :title="$t('components.settings.job_props.settings.absence.title')">
                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="settings.absencePastDays"
                                    :label="$t('components.settings.job_props.settings.absence.past_days')"
                                    :description="$t('common.day', 2)"
                                >
                                    <UInputNumber
                                        v-model="state.settings.absencePastDays"
                                        :disabled="!canSubmit || !canEdit"
                                        :min="0"
                                        :placeholder="$t('common.day', 2)"
                                        :label="$t('common.day', 2)"
                                    />
                                </UFormField>

                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="settings.absenceFutureDays"
                                    :label="$t('components.settings.job_props.settings.absence.future_days')"
                                    :description="$t('common.day', 2)"
                                >
                                    <UInputNumber
                                        v-model="state.settings.absenceFutureDays"
                                        :disabled="!canSubmit || !canEdit"
                                        :min="7"
                                        :max="186"
                                        :placeholder="$t('common.day', 2)"
                                        :label="$t('common.day', 2)"
                                    />
                                </UFormField>
                            </UPageCard>
                        </template>

                        <template v-if="appConfig.discord.botEnabled" #discord>
                            <div v-if="isRequestPending(status)" class="space-y-1 px-4">
                                <USkeleton v-for="idx in 10" :key="idx" class="h-20 w-full" />
                            </div>

                            <UPageCard
                                v-else-if="jobProps.discordSyncSettings"
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
                                >
                                    <template v-if="nuiEnabled || dcConnectRequired">
                                        <NotSupportedTabletBlock v-if="nuiEnabled" />
                                        <UButton
                                            v-else-if="dcConnectRequired"
                                            icon="i-mdi-connection"
                                            :label="$t('common.connect')"
                                            @click="
                                                async () => {
                                                    await navigateTo(
                                                        generateDiscordConnectURL('discord', '/settings/props?tab=discord#'),
                                                        { external: true },
                                                    );
                                                }
                                            "
                                        />

                                        <div v-if="state.discordGuildId" class="mt-2">{{ state.discordGuildId }}</div>
                                    </template>
                                    <USelectMenu
                                        v-else
                                        v-model="state.discordGuildId"
                                        :items="userGuilds"
                                        :filter-fields="['name', 'id']"
                                        :placeholder="
                                            $t('components.settings.job_props.discord_sync_settings.discord_guild_id')
                                        "
                                        value-key="id"
                                        :disabled="!canSubmit || !canEdit || userGuilds?.length === 0"
                                        size="lg"
                                    >
                                        <template #default>
                                            <div class="inline-flex items-center gap-2">
                                                <UAvatar
                                                    :src="userGuilds?.find((g) => g.id === state.discordGuildId)?.icon"
                                                    :alt="userGuilds?.find((g) => g.id === state.discordGuildId)?.name"
                                                />

                                                <span class="truncate">{{
                                                    userGuilds?.find((g) => g.id === state.discordGuildId)?.name ??
                                                    (state.discordGuildId !== '' ? state.discordGuildId : '&nbsp;')
                                                }}</span>
                                            </div>
                                        </template>

                                        <template #item-label="{ item }">
                                            <div class="inline-flex items-center gap-2">
                                                <UAvatar :src="item.icon" :alt="item.name" />

                                                <span class="truncate">{{ item.name }}</span>
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
                                >
                                    <USwitch v-model="state.discordSyncSettings.dryRun" :disabled="!canSubmit || !canEdit" />
                                </UFormField>

                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordSyncSettings.statusLog"
                                    :label="$t('components.settings.job_props.discord_sync_settings.status_log')"
                                >
                                    <USwitch v-model="state.discordSyncSettings.statusLog" :disabled="!canSubmit || !canEdit" />
                                </UFormField>

                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordSyncSettings.statusLogSettings.channelId"
                                    :label="
                                        $t('components.settings.job_props.discord_sync_settings.status_log_settings.channel_id')
                                    "
                                >
                                    <SelectMenu
                                        v-model="state.discordSyncSettings.statusLogSettings!.channelId"
                                        name="discordSyncSettings.statusLogSettings.channelId"
                                        :disabled="
                                            !state.discordGuildId ||
                                            !state.discordSyncSettings.statusLog ||
                                            !canSubmit ||
                                            !canEdit
                                        "
                                        :searchable="
                                            () =>
                                                searchChannels().then((channels) =>
                                                    channels.map((c) => ({
                                                        id: c.id,
                                                        type: 'item',
                                                        label: `${c.name} (${c.id})`,
                                                        item: c,
                                                    })),
                                                )
                                        "
                                        searchable-key="settings-jobprops-discord-channels"
                                        :filter-fields="['name']"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                        value-key="id"
                                        class="w-full"
                                        nullable
                                        :placeholder="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.status_log_settings.channel_id',
                                            )
                                        "
                                    >
                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.channel', 1)]) }}
                                        </template>
                                    </SelectMenu>
                                </UFormField>

                                <UAlert
                                    :ui="{
                                        icon: 'size-6',
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

                            <UPageCard
                                :title="$t('components.settings.job_props.discord_sync_features.title')"
                                :description="$t('components.settings.job_props.discord_sync_features.subtitle')"
                            >
                                <UFormField
                                    class="grid grid-cols-2 items-center gap-2"
                                    name="discordSyncSettings.userInfoSync"
                                    :label="$t('components.settings.job_props.discord_sync_settings.user_info_sync')"
                                >
                                    <USwitch
                                        v-model="state.discordSyncSettings.userInfoSync"
                                        :disabled="!canSubmit || !canEdit"
                                    />
                                </UFormField>

                                <template v-if="jobProps.discordSyncSettings?.userInfoSyncSettings">
                                    <UFormField
                                        class="grid grid-cols-2 items-center gap-2"
                                        name="discordSyncSettings.userInfoSyncSettings.gradeRoleFormat"
                                        :label="
                                            $t(
                                                'components.settings.job_props.discord_sync_settings.user_info_sync_settings.grade_role_format.title',
                                            )
                                        "
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
                                                :search-input="{ placeholder: $t('common.search_field') }"
                                                class="w-full"
                                            >
                                                <template #default>
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

                                                <template #item-label="{ item }">
                                                    {{
                                                        $t(
                                                            `enums.settings.UserInfoSyncUnemployedMode.${UserInfoSyncUnemployedMode[item.value]}`,
                                                        )
                                                    }}
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
                                            class="w-full"
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
                                                class="w-full"
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
                                                :search-input="{ placeholder: $t('common.search_field') }"
                                            >
                                                <template #default>
                                                    {{ $d(toDate(selectedChange?.time), 'short') }}
                                                </template>

                                                <template #item-label="{ item }">
                                                    {{ $d(toDate(item.time), 'short') }}
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>

                                        <LazyPartialsCodeDiff
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
                        </template>
                    </UTabs>
                </template>
            </UForm>
        </template>
    </UDashboardPanel>
</template>

<style scoped>
.codediff:deep(.diff-commandbar) {
    display: none;
}
</style>
