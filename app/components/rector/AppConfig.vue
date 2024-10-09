<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { toDuration } from '~/utils/duration';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { DiscordBotPresenceType } from '~~/gen/ts/resources/rector/config';
import type { GetAppConfigResponse } from '~~/gen/ts/services/rector/config';

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const notifications = useNotificatorStore();

const { data: config, pending: loading, refresh, error } = useLazyAsyncData(`rector-appconfig`, () => getAppConfig());

async function getAppConfig(): Promise<GetAppConfigResponse> {
    try {
        const call = getGRPCRectorConfigClient().getAppConfig({});
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const completorStore = useCompletorStore();
const { listJobs } = completorStore;

const { data: jobs } = useLazyAsyncData(`rector-appconfig-jobs`, () => listJobs());

const schema = z.object({
    auth: z.object({
        signupEnabled: z.boolean(),
        lastCharLock: z.boolean(),
    }),
    perms: z.object({
        default: z
            .array(
                z.object({
                    category: z.string().min(1).max(48),
                    name: z.string().min(1).max(48),
                }),
            )
            .max(25),
    }),
    website: z.object({
        links: z.object({
            privacyPolicy: z.union([z.string().min(1).max(255).url().startsWith('https://'), z.string().length(0).optional()]),
            imprint: z.union([z.string().min(1).max(255).url().startsWith('https://'), z.string().length(0).optional()]),
        }),
        statsPage: z.boolean(),
    }),
    jobInfo: z.object({
        unemployedJob: z.object({
            name: z.string().min(1).max(20),
            grade: z.coerce.number().min(1).max(99),
        }),
        publicJobs: z.string().array().max(99),
        hiddenJobs: z.string().array().max(99),
    }),
    userTracker: z.object({
        refreshTime: zodDurationSchema,
        dbRefreshTime: zodDurationSchema,
        livemapJobs: z.string().array().max(99),
    }),
    discord: z.object({
        enabled: z.boolean(),
        syncInterval: zodDurationSchema,
        inviteUrl: z.union([
            z.string().min(1).max(255).url().startsWith('https://discord.com/'),
            z.string().length(0).optional(),
        ]),
        ignoredJobs: z.string().array().max(99),
        botPresence: z
            .object({
                type: z.nativeEnum(DiscordBotPresenceType),
                status: z.string().max(255).optional(),
                url: z.union([z.string().max(255).url(), z.string().length(0).optional()]),
            })
            .optional(),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    auth: {
        signupEnabled: true,
        lastCharLock: false,
    },
    perms: {
        default: [],
    },
    website: {
        links: {},
        statsPage: false,
    },
    jobInfo: {
        hiddenJobs: [],
        publicJobs: [],
        unemployedJob: {
            name: '',
            grade: 1,
        },
    },
    userTracker: {
        dbRefreshTime: 1.0,
        refreshTime: 3.35,
        livemapJobs: [],
    },
    // Discord
    discord: {
        enabled: false,
        syncInterval: 9.0,
        inviteUrl: '',
        ignoredJobs: [],
        botPresence: {
            type: DiscordBotPresenceType.UNSPECIFIED,
        },
    },
});

async function updateAppConfig(values: Schema): Promise<void> {
    if (!config.value || !config.value?.config) {
        return;
    }

    // Update local version of retrieved config
    config.value.config.auth = values.auth;
    config.value.config.perms = values.perms;
    config.value.config.website = values.website;
    config.value.config.jobInfo = values.jobInfo;
    config.value.config.userTracker = {
        livemapJobs: values.userTracker.livemapJobs,
        dbRefreshTime: toDuration(values.userTracker.dbRefreshTime),
        refreshTime: toDuration(values.userTracker.refreshTime),
    };
    config.value.config.discord = {
        enabled: values.discord.enabled,
        inviteUrl: values.discord.inviteUrl,
        syncInterval: toDuration(values.discord.syncInterval),
        ignoredJobs: values.discord.ignoredJobs,
        botPresence: values.discord.botPresence,
    };

    try {
        const { response } = await getGRPCRectorConfigClient().updateAppConfig({
            config: config.value?.config,
        });

        notifications.add({
            title: { key: 'notifications.rector.app_config.title', parameters: {} },
            description: { key: 'notifications.rector.app_config.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (response.config) {
            config.value = response;
        } else {
            refresh();
        }
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function setSettingsValues(): void {
    if (!config.value || !config.value.config) {
        return;
    }

    if (config.value.config.auth) {
        state.auth = config.value.config.auth;
    }
    if (config.value.config.perms) {
        state.perms = config.value.config.perms;
    }
    if (config.value.config.website) {
        state.website.statsPage = config.value.config.website.statsPage;
        if (config.value.config.website.links) {
            state.website.links = config.value.config.website.links;
        }
    }
    if (config.value.config.jobInfo) {
        if (config.value.config.jobInfo.unemployedJob) {
            state.jobInfo.unemployedJob = config.value.config.jobInfo.unemployedJob;
        }
        state.jobInfo.hiddenJobs = config.value.config.jobInfo.hiddenJobs;
        state.jobInfo.publicJobs = config.value.config.jobInfo.publicJobs;
    }
    if (config.value.config.userTracker) {
        if (config.value.config.userTracker.dbRefreshTime) {
            state.userTracker.dbRefreshTime = fromDuration(config.value.config.userTracker.dbRefreshTime);
        }
        if (config.value.config.userTracker.refreshTime) {
            state.userTracker.refreshTime = fromDuration(config.value.config.userTracker.refreshTime);
        }
        state.userTracker.livemapJobs = config.value.config.userTracker.livemapJobs;
    }
    if (config.value.config.discord) {
        state.discord.enabled = config.value.config.discord.enabled;
        if (config.value.config.discord.syncInterval) {
            state.discord.syncInterval = fromDuration(config.value.config.discord.syncInterval);
        }
        state.discord.inviteUrl = config.value.config.discord.inviteUrl;
        state.discord.ignoredJobs = config.value.config.discord.ignoredJobs;
        state.discord.botPresence = config.value.config.discord.botPresence;
    }
}

watch(config, () => setSettingsValues());

const items = [
    { slot: 'auth', label: t('components.rector.app_config.auth.title'), icon: 'i-mdi-login' },
    { slot: 'perms', label: t('components.rector.app_config.perms.title'), icon: 'i-mdi-user-access-control' },
    { slot: 'website', label: t('components.rector.app_config.website.title'), icon: 'i-mdi-spider-web' },
    { slot: 'jobInfo', label: t('components.rector.app_config.job_info.title'), icon: 'i-mdi-briefcase' },
    { slot: 'userTracker', label: t('components.rector.app_config.user_tracker.title'), icon: 'i-mdi-track-changes' },
    { slot: 'discord', label: t('common.discord'), icon: 'i-simple-icons-discord' },
];

const botPresenceTypes = ref<{ mode: DiscordBotPresenceType }[]>([
    { mode: DiscordBotPresenceType.UNSPECIFIED },
    { mode: DiscordBotPresenceType.GAME },
    { mode: DiscordBotPresenceType.LISTENING },
    { mode: DiscordBotPresenceType.STREAMING },
    { mode: DiscordBotPresenceType.WATCH },
]);

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

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await updateAppConfig(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
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
            <UDashboardNavbar :title="$t('pages.rector.settings.title')">
                <template #right>
                    <UButton color="black" icon="i-mdi-arrow-back" to="/rector">
                        {{ $t('common.back') }}
                    </UButton>

                    <UButton
                        v-if="config"
                        type="submit"
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                    >
                        {{ $t('common.save', 1) }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <div v-if="loading" class="space-y-1 px-4">
                <USkeleton class="mb-6 h-11 w-full" />
                <USkeleton v-for="idx in 5" :key="idx" class="h-20 w-full" />
            </div>
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.setting', 2)])"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!config" icon="i-mdi-office-building-cog" :type="$t('common.setting', 2)" />

            <template v-else>
                <UTabs v-model="selectedTab" :items="items" class="w-full" :ui="{ list: { rounded: '' } }" :unmount="false">
                    <template #auth>
                        <UDashboardPanelContent class="pb-24">
                            <UDashboardSection
                                :title="$t('components.rector.app_config.auth.title')"
                                :description="$t('components.rector.app_config.auth.description')"
                            >
                                <UFormGroup
                                    name="auth.signupEnabled"
                                    :label="$t('components.rector.app_config.auth.sign_up')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.auth.signupEnabled">
                                        <span class="sr-only">
                                            {{ $t('components.rector.app_config.auth.sign_up') }}
                                        </span>
                                    </UToggle>
                                </UFormGroup>

                                <UFormGroup
                                    name="auth.lastCharLock"
                                    :label="$t('components.rector.app_config.auth.last_char_lock')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.auth.lastCharLock">
                                        <span class="sr-only">
                                            {{ $t('components.rector.app_config.auth.last_char_lock') }}
                                        </span>
                                    </UToggle>
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>

                    <template #perms>
                        <UDashboardPanelContent class="pb-24">
                            <UDashboardSection
                                :title="$t('components.rector.app_config.perms.title')"
                                :description="$t('components.rector.app_config.perms.description')"
                            >
                                <UFormGroup
                                    name="perms.default"
                                    :label="$t('components.rector.app_config.perms.default_perms')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <div class="flex flex-col gap-1">
                                        <div v-for="(_, idx) in state.perms.default" :key="idx" class="flex items-center gap-1">
                                            <UFormGroup :name="`perms.default.${idx}.category`" class="flex-1">
                                                <UInput
                                                    v-model="state.perms.default[idx]!.category"
                                                    type="text"
                                                    :placeholder="$t('common.category')"
                                                />
                                            </UFormGroup>

                                            <UFormGroup :name="`perms.default.${idx}.name`" class="flex-1">
                                                <UInput
                                                    v-model="state.perms.default[idx]!.name"
                                                    type="text"
                                                    :placeholder="$t('common.name')"
                                                />
                                            </UFormGroup>

                                            <UButton
                                                :ui="{ rounded: 'rounded-full' }"
                                                icon="i-mdi-close"
                                                @click="state.perms.default.splice(idx, 1)"
                                            />
                                        </div>
                                    </div>

                                    <UButton
                                        :ui="{ rounded: 'rounded-full' }"
                                        :disabled="!canSubmit"
                                        icon="i-mdi-plus"
                                        :class="state.perms.default.length ? 'mt-2' : ''"
                                        @click="state.perms.default.push({ category: '', name: '' })"
                                    />
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>

                    <template #website>
                        <UDashboardPanelContent class="pb-24">
                            <UDashboardSection
                                :title="$t('components.rector.app_config.website.title')"
                                :description="$t('components.rector.app_config.website.description')"
                            >
                                <UFormGroup
                                    name="website.links.privacyPolicy"
                                    :label="$t('common.privacy_policy')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.website.links.privacyPolicy"
                                        type="text"
                                        :placeholder="$t('common.privacy_policy')"
                                        maxlength="255"
                                    />
                                </UFormGroup>

                                <UFormGroup
                                    name="website.links.imprint"
                                    :label="$t('common.imprint')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.website.links.imprint"
                                        type="text"
                                        :placeholder="$t('common.imprint')"
                                        maxlength="255"
                                    />
                                </UFormGroup>
                            </UDashboardSection>

                            <UDashboardSection :title="$t('pages.stats.title')" :description="$t('pages.stats.subtitle')">
                                <UFormGroup
                                    name="website.statsPage"
                                    :label="$t('common.stats')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.website.statsPage">
                                        <span class="sr-only">
                                            {{ $t('common.enabled') }}
                                        </span>
                                    </UToggle>
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>

                    <template #jobInfo>
                        <UDashboardPanelContent class="pb-24">
                            <UDashboardSection
                                :title="$t('components.rector.app_config.job_info.title')"
                                :description="$t('components.rector.app_config.job_info.description')"
                            >
                                <UFormGroup
                                    name="jobInfo.unemployedJob.name"
                                    :label="`${$t('common.job')} ${$t('common.name')}`"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.jobInfo.unemployedJob.name"
                                        type="text"
                                        :placeholder="$t('common.job')"
                                        maxlength="255"
                                    />
                                </UFormGroup>

                                <UFormGroup
                                    name="jobInfo.unemployedJob.grade"
                                    :label="$t('common.rank')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.jobInfo.unemployedJob.grade"
                                        type="number"
                                        :min="1"
                                        :max="99"
                                        name="jobInfoUnemployedGrade"
                                        :placeholder="$t('common.rank')"
                                        :label="$t('common.rank')"
                                    />
                                </UFormGroup>

                                <UFormGroup
                                    name="jobInfo.publicJobs"
                                    :label="$t('components.rector.app_config.job_info.public_jobs')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <ClientOnly>
                                        <USelectMenu
                                            v-model="state.jobInfo.publicJobs"
                                            multiple
                                            :options="jobs ?? []"
                                            value-attribute="name"
                                            :searchable-placeholder="$t('common.search_field')"
                                        >
                                            <template #label>
                                                <template v-if="state.jobInfo.publicJobs.length">
                                                    <span class="truncate">{{ state.jobInfo.publicJobs.join(', ') }}</span>
                                                </template>
                                                <template v-else>
                                                    <span class="truncate">{{
                                                        $t('common.none_selected', [$t('common.job')])
                                                    }}</span>
                                                </template>
                                            </template>
                                            <template #option="{ option: job }">
                                                <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                            </template>
                                        </USelectMenu>
                                    </ClientOnly>
                                </UFormGroup>

                                <UFormGroup
                                    name="jobInfo.hiddenJobs"
                                    :label="$t('components.rector.app_config.job_info.hidden_jobs')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <ClientOnly>
                                        <USelectMenu
                                            v-model="state.jobInfo.hiddenJobs"
                                            multiple
                                            :options="jobs ?? []"
                                            value-attribute="name"
                                            :searchable-placeholder="$t('common.search_field')"
                                        >
                                            <template #label>
                                                <template v-if="state.jobInfo.hiddenJobs.length">
                                                    <span class="truncate">{{ state.jobInfo.hiddenJobs.join(', ') }}</span>
                                                </template>
                                                <template v-else>
                                                    <span class="truncate">{{
                                                        $t('common.none_selected', [$t('common.job')])
                                                    }}</span>
                                                </template>
                                            </template>
                                            <template #option="{ option: job }">
                                                <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                            </template>
                                        </USelectMenu>
                                    </ClientOnly>
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>

                    <template #userTracker>
                        <UDashboardPanelContent class="pb-24">
                            <UDashboardSection
                                :title="$t('components.rector.app_config.user_tracker.title')"
                                :description="$t('components.rector.app_config.user_tracker.description')"
                            >
                                <UFormGroup
                                    name="userTracker.refreshTime"
                                    :label="$t('components.rector.app_config.user_tracker.refresh_time')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.userTracker.refreshTime"
                                        type="number"
                                        :min="1"
                                        :step="0.01"
                                        :placeholder="$t('common.duration')"
                                    >
                                        <template #trailing>
                                            <span class="text-xs text-gray-500 dark:text-gray-400">s</span>
                                        </template>
                                    </UInput>
                                </UFormGroup>

                                <UFormGroup
                                    name="userTracker.dbRefreshTime"
                                    :label="$t('components.rector.app_config.user_tracker.db_refresh_time')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.userTracker.dbRefreshTime"
                                        type="number"
                                        :min="1"
                                        :step="0.01"
                                        :placeholder="$t('common.duration')"
                                    >
                                        <template #trailing>
                                            <span class="text-xs text-gray-500 dark:text-gray-400">s</span>
                                        </template>
                                    </UInput>
                                </UFormGroup>

                                <UFormGroup
                                    name="userTracker.livemapJobs"
                                    :label="$t('components.rector.app_config.user_tracker.livemap_jobs')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <ClientOnly>
                                        <ClientOnly>
                                            <USelectMenu
                                                v-model="state.userTracker.livemapJobs"
                                                multiple
                                                :options="jobs ?? []"
                                                value-attribute="name"
                                                :searchable-placeholder="$t('common.search_field')"
                                            >
                                                <template #label>
                                                    <template v-if="state.userTracker.livemapJobs.length">
                                                        <span class="truncate">{{
                                                            state.userTracker.livemapJobs.join(', ')
                                                        }}</span>
                                                    </template>
                                                    <template v-else>
                                                        <span class="truncate">{{
                                                            $t('common.none_selected', [$t('common.job')])
                                                        }}</span>
                                                    </template>
                                                </template>
                                                <template #option="{ option: job }">
                                                    <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>
                                    </ClientOnly>
                                </UFormGroup>
                            </UDashboardSection>
                        </UDashboardPanelContent>
                    </template>

                    <template #discord>
                        <UDashboardPanelContent class="pb-24">
                            <UDashboardSection
                                :title="$t('common.discord')"
                                :description="$t('components.rector.app_config.discord.description')"
                            >
                                <UFormGroup
                                    name="discordEnabled"
                                    :label="$t('common.enabled')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UToggle v-model="state.discord.enabled">
                                        <span class="sr-only">
                                            {{ $t('common.enabled') }}
                                        </span>
                                    </UToggle>
                                </UFormGroup>

                                <UFormGroup
                                    name="discord.syncInterval"
                                    :label="$t('components.rector.app_config.discord.sync_interval')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.discord.syncInterval"
                                        type="number"
                                        :min="1"
                                        :step="0.01"
                                        name="discord.syncInterval"
                                        :placeholder="$t('common.duration')"
                                    >
                                        <template #trailing>
                                            <span class="text-xs text-gray-500 dark:text-gray-400">s</span>
                                        </template>
                                    </UInput>
                                </UFormGroup>

                                <UFormGroup
                                    name="discord.inviteUrl"
                                    :label="$t('components.rector.app_config.discord.bot_invite_url')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.discord.inviteUrl"
                                        type="text"
                                        :placeholder="$t('components.rector.app_config.discord.bot_invite_url')"
                                    />
                                </UFormGroup>

                                <UFormGroup
                                    name="discord.ignoredJobs"
                                    :label="$t('components.rector.app_config.discord.ignored_jobs')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <ClientOnly>
                                        <USelectMenu
                                            v-model="state.discord.ignoredJobs"
                                            multiple
                                            :options="jobs ?? []"
                                            value-attribute="name"
                                            :searchable-placeholder="$t('common.search_field')"
                                        >
                                            <template #label>
                                                <template v-if="state.discord.ignoredJobs.length > 0">
                                                    <span class="truncate">{{ state.discord.ignoredJobs.join(', ') }}</span>
                                                </template>
                                                <template v-else>
                                                    <span class="truncate">{{
                                                        $t('common.none_selected', [$t('common.job')])
                                                    }}</span>
                                                </template>
                                            </template>
                                            <template #option="{ option: job }">
                                                <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                            </template>
                                        </USelectMenu>
                                    </ClientOnly>
                                </UFormGroup>
                            </UDashboardSection>

                            <UDashboardSection
                                v-if="state.discord.botPresence"
                                :title="$t('components.rector.app_config.discord.bot_presence.title')"
                                :description="$t('components.rector.app_config.discord.bot_presence.description')"
                            >
                                <UFormGroup
                                    name="discord.botPresence.type"
                                    :label="$t('components.rector.app_config.discord.bot_presence.type')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <USelectMenu
                                        v-model="state.discord.botPresence.type"
                                        :options="botPresenceTypes"
                                        value-attribute="mode"
                                        :placeholder="$t('components.rector.app_config.discord.bot_presence.type')"
                                    >
                                        <template #label>
                                            <span class="truncate text-gray-900 dark:text-white">{{
                                                $t(
                                                    `enums.rector.AppConfig.DiscordBotPresenceType.${DiscordBotPresenceType[state.discord.botPresence.type ?? 0]}`,
                                                )
                                            }}</span>
                                        </template>
                                        <template #option="{ option }">
                                            <span class="truncate">{{
                                                $t(
                                                    `enums.rector.AppConfig.DiscordBotPresenceType.${DiscordBotPresenceType[option.mode ?? 0]}`,
                                                )
                                            }}</span>
                                        </template>
                                    </USelectMenu>
                                </UFormGroup>

                                <UFormGroup
                                    name="discord.botPresence.status"
                                    :label="$t('components.rector.app_config.discord.bot_presence.status')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.discord.botPresence.status"
                                        type="text"
                                        :placeholder="$t('common.status')"
                                    />
                                </UFormGroup>

                                <UFormGroup
                                    name="discord.botPresence.url"
                                    :label="$t('components.rector.app_config.discord.bot_presence.url')"
                                    class="grid grid-cols-2 items-center gap-2"
                                    :ui="{ container: '' }"
                                >
                                    <UInput
                                        v-model="state.discord.botPresence.url"
                                        type="text"
                                        :placeholder="$t('components.rector.app_config.discord.bot_presence.url')"
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
