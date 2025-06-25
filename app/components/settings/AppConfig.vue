<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import type { LocaleObject } from '@nuxtjs/i18n';
import { subDays } from 'date-fns';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import { useCompletorStore } from '~/stores/completor';
import { useSettingsStore } from '~/stores/settings';
import { toDuration } from '~/utils/duration';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { DiscordBotPresenceType } from '~~/gen/ts/resources/settings/config';
import type { GetAppConfigResponse } from '~~/gen/ts/services/settings/config';
import { grpcMethods, grpcServices } from '~~/gen/ts/svcs';
import DatePickerPopoverClient from '../partials/DatePickerPopover.client.vue';
import TiptapEditor from '../partials/editor/TiptapEditor.vue';

const { $grpc } = useNuxtApp();

const { t, locales } = useI18n();

const { game } = useAppConfig();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const notifications = useNotificationsStore();

const { data: config, pending: loading, refresh, error } = useLazyAsyncData(`settings-appconfig`, () => getAppConfig());

async function getAppConfig(): Promise<GetAppConfigResponse> {
    try {
        const call = $grpc.settings.config.getAppConfig({});
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const completorStore = useCompletorStore();
const { listJobs } = completorStore;

const { data: jobs } = useLazyAsyncData(`settings-appconfig-jobs`, () => listJobs());

const schema = z.object({
    defaultLocale: z.custom<LocaleObject>(),

    auth: z.object({
        signupEnabled: z.coerce.boolean(),
        lastCharLock: z.coerce.boolean(),
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
        statsPage: z.coerce.boolean(),
    }),
    jobInfo: z.object({
        unemployedJob: z.object({
            name: z.string().min(1).max(20),
            grade: z.coerce.number().min(0).max(99),
        }),
        publicJobs: z.string().array().max(99).default([]),
        hiddenJobs: z.string().array().max(99).default([]),
    }),
    userTracker: z.object({
        refreshTime: zodDurationSchema,
        dbRefreshTime: zodDurationSchema,
    }),
    discord: z.object({
        enabled: z.coerce.boolean(),
        syncInterval: zodDurationSchema,
        botId: z.string().optional(),
        botPermissions: z.coerce.number(),
        inviteUrl: z.union([
            z.string().min(1).max(255).url().startsWith('https://discord.com/'),
            z.string().length(0).optional(),
        ]),
        ignoredJobs: z.string().array().max(99).default([]),
        botPresence: z
            .object({
                type: z.nativeEnum(DiscordBotPresenceType),
                status: z.string().max(255).optional(),
                url: z.union([z.string().max(255).url(), z.string().length(0).optional()]),
            })
            .optional(),
    }),
    system: z.union([
        z.object({
            bannerMessageEnabled: z.literal(false),
            bannerMessage: z.object({
                title: z.string().max(512),
                expiresAt: z.date().optional(),
            }),
        }),
        z.object({
            bannerMessageEnabled: z.literal(true),
            bannerMessage: z.object({
                title: z.string().min(3).max(512),
                expiresAt: z.date().optional(),
            }),
        }),
    ]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    defaultLocale: locales.value[0]!,

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
            grade: game.startJobGrade,
        },
    },
    userTracker: {
        dbRefreshTime: 1.0,
        refreshTime: 3.35,
    },
    discord: {
        enabled: false,
        syncInterval: 9.0,
        botId: '',
        botPermissions: 0,
        inviteUrl: '',
        ignoredJobs: [],
        botPresence: {
            type: DiscordBotPresenceType.UNSPECIFIED,
        },
    },
    system: {
        bannerMessageEnabled: false,
        bannerMessage: {
            title: '',
        },
    },
});

async function updateAppConfig(values: Schema): Promise<void> {
    if (!config.value || !config.value?.config) {
        return;
    }

    // Update local version of retrieved config
    config.value.config.defaultLocale = values.defaultLocale.code;
    config.value.config.auth = values.auth;
    config.value.config.perms = values.perms;
    config.value.config.website = values.website;
    config.value.config.jobInfo = values.jobInfo;
    config.value.config.userTracker = {
        dbRefreshTime: toDuration(values.userTracker.dbRefreshTime),
        refreshTime: toDuration(values.userTracker.refreshTime),
    };
    config.value.config.discord = {
        enabled: values.discord.enabled,
        inviteUrl: values.discord.inviteUrl,
        botId: values.discord.botId ? values.discord.botId : undefined,
        botPermissions: values.discord.botPermissions,
        syncInterval: toDuration(values.discord.syncInterval),
        ignoredJobs: values.discord.ignoredJobs,
        botPresence: values.discord.botPresence,
    };
    config.value.config.system = {
        bannerMessageEnabled: values.system.bannerMessageEnabled,
        bannerMessage: {
            id: '',
            title: values.system.bannerMessage.title,
            expiresAt: values.system.bannerMessage.expiresAt ? toTimestamp(values.system.bannerMessage.expiresAt) : undefined,
        },
    };

    try {
        const { response } = await $grpc.settings.config.updateAppConfig({
            config: config.value?.config,
        });

        notifications.add({
            title: { key: 'notifications.settings.app_config.title', parameters: {} },
            description: { key: 'notifications.settings.app_config.content', parameters: {} },
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

    const locale = locales.value.find((l) => l.code === config.value?.config?.defaultLocale);
    if (locale) {
        state.defaultLocale = locale;
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
    }
    if (config.value.config.discord) {
        state.discord.enabled = config.value.config.discord.enabled;
        if (config.value.config.discord.syncInterval) {
            state.discord.syncInterval = fromDuration(config.value.config.discord.syncInterval);
        }
        state.discord.botId = config.value.config.discord.botId;
        state.discord.botPermissions = config.value.config.discord.botPermissions;
        state.discord.inviteUrl = config.value.config.discord.inviteUrl;
        state.discord.ignoredJobs = config.value.config.discord.ignoredJobs;
        state.discord.botPresence = config.value.config.discord.botPresence;
    }
    if (config.value.config.system) {
        state.system.bannerMessageEnabled = config.value.config.system.bannerMessageEnabled;
        state.system.bannerMessage = {
            title: config.value.config.system.bannerMessage?.title ?? '',
            expiresAt: config.value.config.system.bannerMessage?.expiresAt
                ? toDate(config.value.config.system.bannerMessage?.expiresAt)
                : undefined,
        };
    }
}

watch(config, () => setSettingsValues());

const items = [
    { slot: 'auth', label: t('components.settings.app_config.auth.title'), icon: 'i-mdi-login' },
    { slot: 'perms', label: t('components.settings.app_config.perms.title'), icon: 'i-mdi-user-access-control' },
    { slot: 'website', label: t('components.settings.app_config.website.title'), icon: 'i-mdi-spider-web' },
    { slot: 'jobInfo', label: t('components.settings.app_config.job_info.title'), icon: 'i-mdi-briefcase' },
    { slot: 'userTracker', label: t('components.settings.app_config.user_tracker.title'), icon: 'i-mdi-track-changes' },
    { slot: 'discord', label: t('common.discord'), icon: 'i-simple-icons-discord' },
    { slot: 'system', label: t('common.system'), icon: 'i-mdi-settings' },
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
    if (event.submitter?.getAttribute('role') === 'tab') {
        return;
    }

    canSubmit.value = false;
    await updateAppConfig(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
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
        <UDashboardNavbar :title="$t('pages.settings.settings.title')">
            <template #right>
                <PartialsBackButton fallback-to="/settings" />

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

        <UDashboardPanelContent class="p-0 sm:pb-0">
            <div v-if="loading" class="space-y-1 px-4">
                <USkeleton class="mb-6 h-11 w-full" />
                <USkeleton v-for="idx in 5" :key="idx" class="h-20 w-full" />
            </div>
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.setting', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="!config"
                icon="i-mdi-office-building-cog"
                :type="$t('common.setting', 2)"
                :retry="refresh"
            />

            <UTabs
                v-else
                v-model="selectedTab"
                class="flex flex-1 flex-col"
                :items="items"
                :ui="{
                    wrapper: 'space-y-0 overflow-y-hidden',
                    container: 'flex flex-1 flex-col overflow-y-hidden',
                    base: 'flex flex-1 flex-col overflow-y-hidden',
                    list: { rounded: '' },
                }"
                :unmount="false"
            >
                <template #auth>
                    <UDashboardPanelContent>
                        <UDashboardSection
                            :title="$t('components.settings.app_config.auth.title')"
                            :description="$t('components.settings.app_config.auth.description')"
                        >
                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="auth.signupEnabled"
                                :label="$t('components.settings.app_config.auth.sign_up')"
                                :ui="{ container: '' }"
                            >
                                <UToggle v-model="state.auth.signupEnabled">
                                    <span class="sr-only">
                                        {{ $t('components.settings.app_config.auth.sign_up') }}
                                    </span>
                                </UToggle>
                            </UFormGroup>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="auth.lastCharLock"
                                :label="$t('components.settings.app_config.auth.last_char_lock')"
                                :ui="{ container: '' }"
                            >
                                <UToggle v-model="state.auth.lastCharLock">
                                    <span class="sr-only">
                                        {{ $t('components.settings.app_config.auth.last_char_lock') }}
                                    </span>
                                </UToggle>
                            </UFormGroup>
                        </UDashboardSection>
                    </UDashboardPanelContent>
                </template>

                <template #perms>
                    <UDashboardPanelContent>
                        <UDashboardSection
                            :title="$t('components.settings.app_config.perms.title')"
                            :description="$t('components.settings.app_config.perms.description')"
                        >
                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="perms.default"
                                :label="$t('components.settings.app_config.perms.default_perms')"
                                :ui="{ container: '' }"
                            >
                                <div class="flex flex-col gap-1">
                                    <div v-for="(_, idx) in state.perms.default" :key="idx" class="flex items-center gap-1">
                                        <UFormGroup class="flex-1" :name="`perms.default.${idx}.category`">
                                            <ClientOnly>
                                                <USelectMenu
                                                    v-model="state.perms.default[idx]!.category"
                                                    searchable
                                                    :placeholder="$t('common.service')"
                                                    :options="grpcServices"
                                                >
                                                    <template #option-empty="{ query: search }">
                                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                                    </template>

                                                    <template #empty>
                                                        {{ $t('common.not_found', [$t('common.service')]) }}
                                                    </template>
                                                </USelectMenu>
                                            </ClientOnly>
                                        </UFormGroup>

                                        <UFormGroup class="flex-1" :name="`perms.default.${idx}.name`">
                                            <USelectMenu
                                                v-model="state.perms.default[idx]!.name"
                                                searchable
                                                :placeholder="$t('common.method')"
                                                :options="
                                                    grpcMethods
                                                        .filter((m) => m.startsWith(state.perms.default[idx]!.category + '/'))
                                                        .map((m) => m.split('/').at(1) ?? m)
                                                "
                                            >
                                                <template #option-empty="{ query: search }">
                                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                                </template>

                                                <template #empty>
                                                    {{ $t('common.not_found', [$t('common.method')]) }}
                                                </template>
                                            </USelectMenu>
                                        </UFormGroup>

                                        <UButton
                                            :ui="{ rounded: 'rounded-full' }"
                                            icon="i-mdi-close"
                                            @click="state.perms.default.splice(idx, 1)"
                                        />
                                    </div>
                                </div>

                                <UButton
                                    :class="state.perms.default.length ? 'mt-2' : ''"
                                    :ui="{ rounded: 'rounded-full' }"
                                    :disabled="!canSubmit"
                                    icon="i-mdi-plus"
                                    @click="state.perms.default.push({ category: '', name: '' })"
                                />
                            </UFormGroup>
                        </UDashboardSection>
                    </UDashboardPanelContent>
                </template>

                <template #website>
                    <UDashboardPanelContent>
                        <UDashboardSection
                            :title="$t('components.settings.app_config.website.title')"
                            :description="$t('components.settings.app_config.website.description')"
                        >
                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="defaultLocale"
                                :label="$t('common.default_lang')"
                                :ui="{ container: '' }"
                            >
                                <USelectMenu
                                    v-model="state.defaultLocale"
                                    :placeholder="$t('common.language', 1)"
                                    :options="locales"
                                >
                                    <template #label>
                                        <template v-if="state.defaultLocale">
                                            <UIcon
                                                class="size-4"
                                                :name="
                                                    locales.find((l) => l.code === state.defaultLocale.code)?.icon ??
                                                    'i-mdi-question-mark'
                                                "
                                            />
                                            <span class="truncate">{{
                                                state.defaultLocale.name ?? state.defaultLocale.code
                                            }}</span>
                                        </template>
                                        <template v-else>
                                            <span class="truncate">{{
                                                $t('common.none_selected', [$t('common.language')])
                                            }}</span>
                                        </template>
                                    </template>

                                    <template #option="{ option: locale }">
                                        <UIcon class="size-4" :name="locale.icon" />
                                        <span class="truncate">{{ locale.name }}</span>
                                    </template>
                                </USelectMenu>
                            </UFormGroup>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="website.links.privacyPolicy"
                                :label="$t('common.privacy_policy')"
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
                                class="grid grid-cols-2 items-center gap-2"
                                name="website.links.imprint"
                                :label="$t('common.imprint')"
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
                                class="grid grid-cols-2 items-center gap-2"
                                name="website.statsPage"
                                :label="$t('common.stats')"
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
                    <UDashboardPanelContent>
                        <UDashboardSection
                            :title="$t('components.settings.app_config.job_info.title')"
                            :description="$t('components.settings.app_config.job_info.description')"
                        >
                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="jobInfo.unemployedJob.name"
                                :label="`${$t('common.job')} ${$t('common.name')}`"
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
                                class="grid grid-cols-2 items-center gap-2"
                                name="jobInfo.unemployedJob.grade"
                                :label="$t('common.rank')"
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
                                class="grid grid-cols-2 items-center gap-2"
                                name="jobInfo.publicJobs"
                                :label="$t('components.settings.app_config.job_info.public_jobs')"
                                :ui="{ container: '' }"
                            >
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.jobInfo.publicJobs"
                                        multiple
                                        :options="jobs ?? []"
                                        searchable
                                        value-attribute="name"
                                        :searchable-placeholder="$t('common.search_field')"
                                        :search-attributes="['label', 'name']"
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
                                class="grid grid-cols-2 items-center gap-2"
                                name="jobInfo.hiddenJobs"
                                :label="$t('components.settings.app_config.job_info.hidden_jobs')"
                                :ui="{ container: '' }"
                            >
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.jobInfo.hiddenJobs"
                                        multiple
                                        :options="jobs ?? []"
                                        searchable
                                        value-attribute="name"
                                        :searchable-placeholder="$t('common.search_field')"
                                        :search-attributes="['label', 'name']"
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
                    <UDashboardPanelContent>
                        <UDashboardSection
                            :title="$t('components.settings.app_config.user_tracker.title')"
                            :description="$t('components.settings.app_config.user_tracker.description')"
                        >
                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="userTracker.refreshTime"
                                :label="$t('components.settings.app_config.user_tracker.refresh_time')"
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
                                class="grid grid-cols-2 items-center gap-2"
                                name="userTracker.dbRefreshTime"
                                :label="$t('components.settings.app_config.user_tracker.db_refresh_time')"
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
                        </UDashboardSection>
                    </UDashboardPanelContent>
                </template>

                <template #discord>
                    <UDashboardPanelContent>
                        <UDashboardSection
                            :title="$t('common.discord')"
                            :description="$t('components.settings.app_config.discord.description')"
                        >
                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="discordEnabled"
                                :label="$t('common.enabled')"
                                :ui="{ container: '' }"
                            >
                                <UToggle v-model="state.discord.enabled">
                                    <span class="sr-only">
                                        {{ $t('common.enabled') }}
                                    </span>
                                </UToggle>
                            </UFormGroup>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.syncInterval"
                                :label="$t('components.settings.app_config.discord.sync_interval')"
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
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.botId"
                                :label="$t('components.settings.app_config.discord.bot_id')"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.discord.botId"
                                    type="text"
                                    :placeholder="$t('components.settings.app_config.discord.bot_id')"
                                />
                            </UFormGroup>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.botPermissions"
                                :label="$t('components.settings.app_config.discord.bot_permissions')"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.discord.botPermissions"
                                    type="text"
                                    :placeholder="$t('components.settings.app_config.discord.bot_permissions')"
                                />
                            </UFormGroup>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.ignoredJobs"
                                :label="$t('components.settings.app_config.discord.ignored_jobs')"
                                :ui="{ container: '' }"
                            >
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.discord.ignoredJobs"
                                        multiple
                                        :options="jobs ?? []"
                                        searchable
                                        value-attribute="name"
                                        :searchable-placeholder="$t('common.search_field')"
                                        :search-attributes="['label', 'name']"
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
                            :title="$t('components.settings.app_config.discord.bot_presence.title')"
                            :description="$t('components.settings.app_config.discord.bot_presence.description')"
                        >
                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.botPresence.type"
                                :label="$t('components.settings.app_config.discord.bot_presence.type')"
                                :ui="{ container: '' }"
                            >
                                <USelectMenu
                                    v-model="state.discord.botPresence.type"
                                    :options="botPresenceTypes"
                                    value-attribute="mode"
                                    :placeholder="$t('components.settings.app_config.discord.bot_presence.type')"
                                >
                                    <template #label>
                                        <span class="truncate text-gray-900 dark:text-white">{{
                                            $t(
                                                `enums.settings.AppConfig.DiscordBotPresenceType.${DiscordBotPresenceType[state.discord.botPresence.type ?? 0]}`,
                                            )
                                        }}</span>
                                    </template>

                                    <template #option="{ option }">
                                        <span class="truncate">{{
                                            $t(
                                                `enums.settings.AppConfig.DiscordBotPresenceType.${DiscordBotPresenceType[option.mode ?? 0]}`,
                                            )
                                        }}</span>
                                    </template>
                                </USelectMenu>
                            </UFormGroup>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.botPresence.status"
                                :label="$t('components.settings.app_config.discord.bot_presence.status')"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.discord.botPresence.status"
                                    type="text"
                                    :placeholder="$t('common.status')"
                                />
                            </UFormGroup>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.botPresence.url"
                                :label="$t('components.settings.app_config.discord.bot_presence.url')"
                                :ui="{ container: '' }"
                            >
                                <UInput
                                    v-model="state.discord.botPresence.url"
                                    type="text"
                                    :placeholder="$t('components.settings.app_config.discord.bot_presence.url')"
                                />
                            </UFormGroup>
                        </UDashboardSection>
                    </UDashboardPanelContent>
                </template>

                <template #system>
                    <UDashboardPanelContent>
                        <UDashboardSection
                            :title="$t('components.settings.app_config.system.banner_message.title')"
                            :description="$t('components.settings.app_config.system.banner_message.subtitle')"
                        >
                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="system.bannerMessageEnabled"
                                :label="$t('common.enabled')"
                                :ui="{ container: '' }"
                            >
                                <UToggle v-model="state.system.bannerMessageEnabled">
                                    <span class="sr-only">
                                        {{ $t('common.enabled') }}
                                    </span>
                                </UToggle>
                            </UFormGroup>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="system.bannerMessage.title"
                                :label="$t('common.message')"
                                :ui="{ container: '' }"
                            >
                                <TiptapEditor v-model="state.system.bannerMessage.title" />
                            </UFormGroup>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="system.bannerMessage.expiresAt"
                                :label="$t('common.expires_at')"
                                :ui="{ container: '' }"
                            >
                                <DatePickerPopoverClient
                                    v-model="state.system.bannerMessage.expiresAt"
                                    date-format="dd.MM.yyyy HH:mm"
                                    :date-picker="{
                                        mode: 'dateTime',
                                        is24hr: true,
                                        clearable: true,
                                        disabledDates: [{ start: null, end: subDays(new Date(), 1) }],
                                    }"
                                />
                            </UFormGroup>
                        </UDashboardSection>
                    </UDashboardPanelContent>
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
