<script lang="ts" setup>
import type { FormSubmitEvent, TabsItem } from '@nuxt/ui';
import type { LocaleObject } from '@nuxtjs/i18n';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import { useCompletorStore } from '~/stores/completor';
import { useSettingsStore } from '~/stores/settings';
import { zodProtoDurationSchema } from '~/utils/validation';
import { getSettingsConfigClient } from '~~/gen/ts/clients';
import { GRPCServiceMethods, GRPCServices } from '~~/gen/ts/perms';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { DiscordBotPresenceType } from '~~/gen/ts/resources/settings/config';
import type { GetAppConfigResponse } from '~~/gen/ts/services/settings/config';
import TiptapEditor from '../partials/editor/TiptapEditor.vue';
import InputDatePicker from '../partials/InputDatePicker.vue';
import InputDurationPicker from '../partials/InputDurationPicker.vue';
import { currencies, intlLocales } from './helpers';

const { t, locales } = useI18n();

const { auth, display, game } = useAppConfig();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const notifications = useNotificationsStore();

const settingsConfigClient = await getSettingsConfigClient();

const { data: config, status, refresh, error } = useLazyAsyncData(`settings-appconfig`, () => getAppConfig());

async function getAppConfig(): Promise<GetAppConfigResponse> {
    try {
        const call = settingsConfigClient.getAppConfig({});
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

const botPresenceTypes = ref<{ mode: DiscordBotPresenceType }[]>([
    { mode: DiscordBotPresenceType.UNSPECIFIED },
    { mode: DiscordBotPresenceType.GAME },
    { mode: DiscordBotPresenceType.LISTENING },
    { mode: DiscordBotPresenceType.STREAMING },
    { mode: DiscordBotPresenceType.WATCH },
]);

const minDiscordSyncInterval = secondsToDuration(60);
const maxDiscordSyncInterval = secondsToDuration(24 * 60 * 60);

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
                    category: z.coerce.string().min(1).max(48),
                    name: z.coerce.string().min(1).max(48),
                }),
            )
            .max(25),
    }),
    website: z.object({
        links: z.object({
            privacyPolicy: z.union([z.url().min(1).max(255).startsWith('https://'), z.string().length(0).optional()]),
            imprint: z.union([z.url().min(1).max(255).startsWith('https://'), z.string().length(0).optional()]),
        }),
        statsPage: z.coerce.boolean(),
    }),
    jobInfo: z.object({
        unemployedJob: z.object({
            name: z.coerce.string().min(1).max(20),
            grade: z.coerce.number().min(0).max(99),
        }),
        publicJobs: z.coerce.string().array().max(99).default([]),
        hiddenJobs: z.coerce.string().array().max(99).default([]),
    }),
    userTracker: z.object({
        refreshTime: zodProtoDurationSchema({
            required: true,
            min: secondsToDuration(1),
            max: secondsToDuration(15 * 60 * 60),
        }),
        dbRefreshTime: zodProtoDurationSchema({
            required: true,
            min: secondsToDuration(1),
            max: secondsToDuration(15 * 60 * 60),
        }),
    }),
    discord: z.object({
        enabled: z.coerce.boolean(),
        syncInterval: zodProtoDurationSchema({
            required: true,
            min: minDiscordSyncInterval,
            max: maxDiscordSyncInterval,
        }),
        botId: z.string().optional(),
        botPermissions: z.coerce.number(),
        inviteUrl: z.union([z.url().min(1).max(255).startsWith('https://discord.com/'), z.string().length(0).optional()]),
        ignoredJobs: z.coerce.string().array().max(99).default([]),
        botPresence: z
            .object({
                type: z.enum(DiscordBotPresenceType),
                status: z.string().max(255).optional(),
                url: z.union([z.url({ protocol: /^https$/ }).max(255), z.string().length(0).optional()]),
            })
            .optional(),
    }),
    system: z.union([
        z.object({
            bannerMessageEnabled: z.literal(false),
            bannerMessage: z.object({
                title: z.coerce.string().max(512),
                expiresAt: z.date().optional(),
            }),
        }),
        z.object({
            bannerMessageEnabled: z.literal(true),
            bannerMessage: z.object({
                title: z.coerce.string().min(3).max(512),
                expiresAt: z.date().optional(),
            }),
        }),
    ]),
    display: z.object({
        intlLocale: z.string().default('en-US'),
        currencyName: z.string().default('USD'),
    }),
    quickButtons: z.object({
        penaltyCalculator: z
            .object({
                detentionTimeUnit: z
                    .object({
                        singular: z.string().max(32).optional(),
                        plural: z.string().max(32).optional(),
                    })
                    .optional(),
                maxCount: z.coerce.number().int().min(1).max(100).nonnegative().optional(),
                maxLeeway: z.coerce.number().int().min(0).max(1).nonnegative().optional(),
                warnSettings: z
                    .object({
                        enabled: z.coerce.boolean(),
                        fine: z.coerce.number().int().min(0).max(999_999_999_999).optional(),
                        detentionTime: z.coerce.number().int().min(0).max(999_999_999_999).optional(),
                        stvoPoints: z.coerce.number().int().min(0).max(999_999_999_999).optional(),
                        warnMessage: z.string().max(512).optional(),
                    })
                    .optional(),
            })
            .default({
                detentionTimeUnit: {},
                maxCount: 10,
                maxLeeway: 0.25,
                warnSettings: { enabled: false },
            }),
    }),
    livemap: z.object({
        enableCayoPerico: z.coerce.boolean().default(true),
    }),
    game: z.object({
        maxWantedDurationUserEnabled: z.coerce.boolean().default(false),
        maxWantedDurationUser: zodProtoDurationSchema({
            min: secondsToDuration(24 * 60 * 60),
            max: secondsToDuration(3650 * 24 * 60 * 60),
        }),
        maxWantedDurationVehicleEnabled: z.coerce.boolean().default(false),
        maxWantedDurationVehicle: zodProtoDurationSchema({
            min: secondsToDuration(24 * 60 * 60),
            max: secondsToDuration(3650 * 24 * 60 * 60),
        }),
    }),
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
        dbRefreshTime: secondsToDuration(1),
        refreshTime: secondsToDuration(3.35),
    },
    discord: {
        enabled: false,
        syncInterval: secondsToDuration(900),
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
    display: {
        intlLocale: 'en-US',
        currencyName: 'USD',
    },
    quickButtons: {
        penaltyCalculator: {
            detentionTimeUnit: {
                singular: t('common.month', 1),
                plural: t('common.month', 2),
            },
            maxCount: 10,
            warnSettings: { enabled: false },
        },
    },
    livemap: {
        enableCayoPerico: true,
    },
    game: {
        maxWantedDurationUserEnabled: false,
        maxWantedDurationUser: undefined,
        maxWantedDurationVehicleEnabled: false,
        maxWantedDurationVehicle: undefined,
    },
});

async function updateAppConfig(values: Schema): Promise<void> {
    if (!config.value || !config.value?.config) return;

    // Update local version of retrieved config
    config.value.config.defaultLocale = values.defaultLocale.code;
    config.value.config.auth = values.auth;
    config.value.config.perms = values.perms;
    config.value.config.website = values.website;
    config.value.config.jobInfo = values.jobInfo;
    config.value.config.userTracker = {
        dbRefreshTime: values.userTracker.dbRefreshTime,
        refreshTime: values.userTracker.refreshTime,
    };
    config.value.config.discord = {
        enabled: values.discord.enabled,
        inviteUrl: values.discord.inviteUrl,
        botId: values.discord.botId ? values.discord.botId : undefined,
        botPermissions: values.discord.botPermissions,
        syncInterval: values.discord.syncInterval,
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
    config.value.config.display = values.display;
    config.value.config.quickButtons = values.quickButtons;
    config.value.config.livemap = values.livemap;
    config.value.config.game = {
        maxWantedDurationUserEnabled: values.game.maxWantedDurationUserEnabled,
        maxWantedDurationUser: values.game.maxWantedDurationUser,
        maxWantedDurationVehicleEnabled: values.game.maxWantedDurationVehicleEnabled,
        maxWantedDurationVehicle: values.game.maxWantedDurationVehicle,
    };

    try {
        const { response } = await settingsConfigClient.updateAppConfig({
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
    if (!config.value || !config.value.config) return;

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
            state.userTracker.dbRefreshTime = config.value.config.userTracker.dbRefreshTime;
        }
        if (config.value.config.userTracker.refreshTime) {
            state.userTracker.refreshTime = config.value.config.userTracker.refreshTime;
        }
    }
    if (config.value.config.discord) {
        state.discord.enabled = config.value.config.discord.enabled;
        if (config.value.config.discord.syncInterval) {
            state.discord.syncInterval = config.value.config.discord.syncInterval;
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
    if (config.value.config.display) {
        state.display = {
            intlLocale: config.value.config.display.intlLocale ?? 'en-US',
            currencyName: config.value.config.display.currencyName ?? 'USD',
        };
    }
    if (config.value.config.quickButtons) {
        if (config.value.config.quickButtons.penaltyCalculator) {
            if (!config.value.config.quickButtons.penaltyCalculator.detentionTimeUnit)
                config.value.config.quickButtons.penaltyCalculator.detentionTimeUnit = {
                    singular: t('common.month', 1),
                    plural: t('common.month', 2),
                };

            if (!config.value.config.quickButtons.penaltyCalculator.warnSettings)
                config.value.config.quickButtons.penaltyCalculator.warnSettings = { enabled: false };

            state.quickButtons.penaltyCalculator = {
                detentionTimeUnit: config.value.config.quickButtons.penaltyCalculator.detentionTimeUnit,
                maxCount: config.value.config.quickButtons.penaltyCalculator.maxCount ?? 10,
                maxLeeway: (config.value.config.quickButtons.penaltyCalculator.maxLeeway ?? 25) / 100,
                warnSettings: config.value.config.quickButtons.penaltyCalculator.warnSettings,
            };
        }
    }

    if (config.value.config.livemap) {
        state.livemap = config.value.config.livemap;
    }

    if (config.value.config.game) {
        state.game.maxWantedDurationUserEnabled = config.value.config.game.maxWantedDurationUserEnabled;
        state.game.maxWantedDurationUser = config.value.config.game.maxWantedDurationUser;

        state.game.maxWantedDurationVehicleEnabled = config.value.config.game.maxWantedDurationVehicleEnabled;
        state.game.maxWantedDurationVehicle = config.value.config.game.maxWantedDurationVehicle;
    }
}

watch(config, () => setSettingsValues());

const items = computed<TabsItem[]>(() => [
    { slot: 'auth' as const, label: t('components.settings.app_config.auth.title'), icon: 'i-mdi-account-key', value: 'auth' },
    {
        slot: 'jobInfo' as const,
        label: t('components.settings.app_config.jobs_users.tab'),
        icon: 'i-mdi-briefcase-search',
        value: 'jobInfo',
    },
    { slot: 'discord' as const, label: t('common.discord'), icon: 'i-simple-icons-discord', value: 'discord' },
    {
        slot: 'game' as const,
        label: t('components.settings.app_config.game.tab'),
        icon: 'i-mdi-details',
        value: 'game',
    },
    { slot: 'website' as const, label: t('components.settings.app_config.website.title'), icon: 'i-mdi-web', value: 'website' },
    { slot: 'system' as const, label: t('common.system'), icon: 'i-mdi-cog', value: 'system' },
]);

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'auth';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
    },
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (event.submitter?.getAttribute('role') === 'tab') return;

    canSubmit.value = false;
    await updateAppConfig(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.settings.settings.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/settings" />

                    <UButton
                        v-if="!streamerMode && config"
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
                <div v-if="isRequestPending(status)" class="space-y-1 px-4">
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
                    variant="link"
                    :unmount-on-hide="false"
                    :ui="{ content: 'p-4 flex flex-col gap-4 max-w-(--ui-container) mx-auto' }"
                >
                    <template #auth>
                        <UPageCard
                            :title="$t('components.settings.app_config.auth.title')"
                            :description="$t('components.settings.app_config.auth.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="auth.signupEnabled"
                                :label="$t('components.settings.app_config.auth.sign_up')"
                            >
                                <USwitch v-model="state.auth.signupEnabled" />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="auth.lastCharLock"
                                :label="$t('components.settings.app_config.auth.last_char_lock')"
                            >
                                <USwitch v-model="state.auth.lastCharLock" />
                            </UFormField>
                        </UPageCard>

                        <UPageCard
                            :title="$t('components.settings.app_config.perms.default_perms.title')"
                            :description="$t('components.settings.app_config.perms.default_perms.description')"
                        >
                            <UFormField name="perms.default">
                                <div class="flex flex-col gap-1">
                                    <div
                                        v-for="(_, idx) in state.perms.default"
                                        :key="idx"
                                        class="flex flex-col gap-1 sm:flex-row"
                                    >
                                        <UFormField class="flex-1" :name="`perms.default.${idx}.category`">
                                            <ClientOnly>
                                                <USelectMenu
                                                    v-model="state.perms.default[idx]!.category"
                                                    class="w-full"
                                                    :placeholder="$t('common.service')"
                                                    :items="GRPCServices"
                                                >
                                                    <template v-if="state.perms.default[idx]!.category" #default>
                                                        {{ $t(`perms.${state.perms.default[idx]!.category}.category`) }}
                                                    </template>

                                                    <template #item-label="{ item }">
                                                        {{ $t(`perms.${item}.category`) }}
                                                    </template>

                                                    <template #empty>
                                                        {{ $t('common.not_found', [$t('common.service')]) }}
                                                    </template>
                                                </USelectMenu>
                                            </ClientOnly>
                                        </UFormField>

                                        <UFormField class="flex-1" :name="`perms.default.${idx}.name`">
                                            <USelectMenu
                                                v-model="state.perms.default[idx]!.name"
                                                class="w-full"
                                                :placeholder="$t('common.method')"
                                                :items="
                                                    GRPCServiceMethods.filter((m) =>
                                                        m.startsWith(state.perms.default[idx]!.category + '/'),
                                                    ).map((m) => m.split('/').at(1) ?? m)
                                                "
                                                :disabled="!state.perms.default[idx]?.category"
                                            >
                                                <template v-if="state.perms.default[idx]!.name" #default>
                                                    {{
                                                        $t(
                                                            `perms.${state.perms.default[idx]!.category}.${state.perms.default[idx]!.name}.key`,
                                                        )
                                                    }}
                                                </template>

                                                <template #item-label="{ item }">
                                                    {{ $t(`perms.${state.perms.default[idx]!.category}.${item}.key`) }}
                                                </template>

                                                <template #empty>
                                                    {{ $t('common.not_found', [$t('common.method')]) }}
                                                </template>
                                            </USelectMenu>
                                        </UFormField>

                                        <div>
                                            <UTooltip :text="$t('common.remove')">
                                                <UButton
                                                    color="red"
                                                    icon="i-mdi-close"
                                                    @click="state.perms.default.splice(idx, 1)"
                                                />
                                            </UTooltip>
                                        </div>
                                    </div>
                                </div>

                                <UButton
                                    :class="state.perms.default.length ? 'mt-2' : ''"
                                    :disabled="!canSubmit"
                                    icon="i-mdi-plus"
                                    @click="state.perms.default.push({ category: '', name: '' })"
                                />
                            </UFormField>
                        </UPageCard>

                        <UPageCard
                            :title="$t('components.settings.app_config.auth.social_login_providers.title')"
                            :description="$t('components.settings.app_config.auth.social_login_providers.description')"
                        >
                            <UPageGrid class="lg:grid-cols-2">
                                <UCard
                                    v-for="provider in auth.providers"
                                    :key="provider.name"
                                    :ui="{
                                        header: 'flex flex-col',
                                        body: 'flex-1 flex flex-col',
                                    }"
                                >
                                    <template #header>
                                        <div class="flex flex-1 gap-2">
                                            <div class="inline-flex flex-1 gap-2">
                                                <NuxtImg
                                                    v-if="!provider.icon?.startsWith('i-')"
                                                    class="size-10"
                                                    :src="provider.icon"
                                                    :alt="provider.name"
                                                    placeholder-class="size-10"
                                                    loading="lazy"
                                                />
                                                <UIcon
                                                    v-else
                                                    class="size-10"
                                                    :name="provider.icon"
                                                    :style="provider.name === 'discord' && { color: '#7289da' }"
                                                />

                                                <div
                                                    class="flex items-center gap-1.5 truncate text-base font-semibold text-highlighted"
                                                >
                                                    {{ provider.label }}
                                                </div>
                                            </div>
                                        </div>
                                    </template>
                                    <template #footer>
                                        <UButton
                                            size="xs"
                                            variant="link"
                                            color="neutral"
                                            :label="$t('components.auth.SocialLogins.connection_website')"
                                            external
                                            :to="provider.homepage"
                                            target="_blank"
                                            trailing-icon="i-mdi-external-link"
                                        />
                                    </template>
                                </UCard>
                            </UPageGrid>
                        </UPageCard>
                    </template>

                    <template #jobInfo>
                        <UPageCard
                            :title="$t('components.settings.app_config.jobs_users.unemployed_job.title')"
                            :description="$t('components.settings.app_config.jobs_users.unemployed_job.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="jobInfo.unemployedJob.name"
                                :label="`${$t('common.job')} ${$t('common.name')}`"
                            >
                                <USelectMenu
                                    v-model="state.jobInfo.unemployedJob.name"
                                    class="w-full"
                                    :placeholder="$t('common.job')"
                                    :items="jobs ?? []"
                                    value-key="name"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :filter-fields="['label', 'name']"
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="jobInfo.unemployedJob.grade"
                                :label="$t('common.rank')"
                            >
                                <UInputNumber
                                    v-model="state.jobInfo.unemployedJob.grade"
                                    class="w-full"
                                    :min="
                                        jobs?.find((j) => j.name === state.jobInfo.unemployedJob.name)?.grades?.[0]?.grade ?? 0
                                    "
                                    :max="
                                        jobs?.find((j) => j.name === state.jobInfo.unemployedJob.name)?.grades?.at(-1)?.grade ??
                                        99
                                    "
                                    name="jobInfoUnemployedGrade"
                                    :placeholder="$t('common.rank')"
                                    :label="$t('common.rank')"
                                />
                            </UFormField>
                        </UPageCard>

                        <UPageCard
                            :title="$t('components.settings.app_config.jobs_users.title')"
                            :description="$t('components.settings.app_config.jobs_users.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="jobInfo.publicJobs"
                                :label="$t('components.settings.app_config.jobs_users.public_jobs')"
                            >
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.jobInfo.publicJobs"
                                        class="w-full"
                                        multiple
                                        :items="jobs ?? []"
                                        value-key="name"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                        :filter-fields="['label', 'name']"
                                    >
                                        <template #default>
                                            <span v-if="state.jobInfo.publicJobs.length" class="truncate">{{
                                                state.jobInfo.publicJobs.join(', ')
                                            }}</span>
                                            <span v-else class="truncate">{{
                                                $t('common.none_selected', [$t('common.job')])
                                            }}</span>
                                        </template>

                                        <template #item-label="{ item }"> {{ item.label }} ({{ item.name }}) </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="jobInfo.hiddenJobs"
                                :label="$t('components.settings.app_config.jobs_users.hidden_jobs')"
                            >
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.jobInfo.hiddenJobs"
                                        class="w-full"
                                        multiple
                                        :items="jobs ?? []"
                                        value-key="name"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                        :filter-fields="['label', 'name']"
                                    >
                                        <template #default>
                                            <span v-if="state.jobInfo.hiddenJobs.length" class="truncate">{{
                                                state.jobInfo.hiddenJobs.join(', ')
                                            }}</span>

                                            <span v-else class="truncate">{{
                                                $t('common.none_selected', [$t('common.job')])
                                            }}</span>
                                        </template>

                                        <template #item-label="{ item }"> {{ item.label }} ({{ item.name }}) </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>
                        </UPageCard>

                        <UPageCard
                            :title="$t('components.settings.app_config.jobs_users.user_tracker.title')"
                            :description="$t('components.settings.app_config.jobs_users.user_tracker.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="userTracker.refreshTime"
                                :label="$t('components.settings.app_config.jobs_users.user_tracker.refresh_time')"
                            >
                                <InputDurationPicker
                                    v-model="state.userTracker.refreshTime"
                                    class="w-full"
                                    :units="['second']"
                                    :min="secondsToDuration(1)"
                                    :step="0.01"
                                    :max="secondsToDuration(15 * 60 * 60)"
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="userTracker.dbRefreshTime"
                                :label="$t('components.settings.app_config.jobs_users.user_tracker.db_refresh_time')"
                            >
                                <InputDurationPicker
                                    v-model="state.userTracker.dbRefreshTime"
                                    class="w-full"
                                    :units="['second']"
                                    :min="secondsToDuration(1)"
                                    :step="0.01"
                                    :max="secondsToDuration(15 * 60 * 60)"
                                />
                            </UFormField>
                        </UPageCard>
                    </template>

                    <template #discord>
                        <UPageCard
                            :title="$t('common.discord')"
                            :description="$t('components.settings.app_config.discord.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="discordEnabled"
                                :label="$t('common.enabled')"
                            >
                                <USwitch v-model="state.discord.enabled" />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.syncInterval"
                                :label="$t('components.settings.app_config.discord.sync_interval')"
                            >
                                <InputDurationPicker
                                    v-model="state.discord.syncInterval"
                                    class="w-full"
                                    :units="['minute', 'second']"
                                    :step="1"
                                    :min="minDiscordSyncInterval"
                                    :max="maxDiscordSyncInterval"
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.botId"
                                :label="$t('components.settings.app_config.discord.bot_id')"
                            >
                                <UInput
                                    v-model="state.discord.botId"
                                    class="w-full"
                                    type="text"
                                    :placeholder="$t('components.settings.app_config.discord.bot_id')"
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.botPermissions"
                                :label="$t('components.settings.app_config.discord.bot_permissions')"
                            >
                                <UInput
                                    v-model="state.discord.botPermissions"
                                    class="w-full"
                                    type="text"
                                    :placeholder="$t('components.settings.app_config.discord.bot_permissions')"
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.ignoredJobs"
                                :label="$t('components.settings.app_config.discord.ignored_jobs')"
                            >
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.discord.ignoredJobs"
                                        class="w-full"
                                        multiple
                                        :items="jobs ?? []"
                                        value-key="name"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                        :filter-fields="['label', 'name']"
                                    >
                                        <template #default>
                                            <span v-if="state.discord.ignoredJobs.length > 0" class="truncate">{{
                                                state.discord.ignoredJobs.join(', ')
                                            }}</span>
                                            <span v-else class="truncate">{{
                                                $t('common.none_selected', [$t('common.job')])
                                            }}</span>
                                        </template>

                                        <template #item-label="{ item }"> {{ item.label }} ({{ item.name }}) </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>
                        </UPageCard>

                        <UPageCard
                            v-if="state.discord.botPresence"
                            :title="$t('components.settings.app_config.discord.bot_presence.title')"
                            :description="$t('components.settings.app_config.discord.bot_presence.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.botPresence.type"
                                :label="$t('components.settings.app_config.discord.bot_presence.type')"
                            >
                                <USelectMenu
                                    v-model="state.discord.botPresence.type"
                                    class="w-full"
                                    :items="botPresenceTypes"
                                    value-key="mode"
                                    :placeholder="$t('components.settings.app_config.discord.bot_presence.type')"
                                >
                                    <template #default>
                                        {{
                                            $t(
                                                `enums.settings.AppConfig.DiscordBotPresenceType.${DiscordBotPresenceType[state.discord.botPresence.type ?? 0]}`,
                                            )
                                        }}
                                    </template>

                                    <template #item-label="{ item }">
                                        {{
                                            $t(
                                                `enums.settings.AppConfig.DiscordBotPresenceType.${DiscordBotPresenceType[item.mode ?? 0]}`,
                                            )
                                        }}
                                    </template>
                                </USelectMenu>
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.botPresence.status"
                                :label="$t('components.settings.app_config.discord.bot_presence.status')"
                            >
                                <UInput
                                    v-model="state.discord.botPresence.status"
                                    class="w-full"
                                    type="text"
                                    :placeholder="$t('common.status')"
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="discord.botPresence.url"
                                :label="$t('components.settings.app_config.discord.bot_presence.url')"
                            >
                                <UInput
                                    v-model="state.discord.botPresence.url"
                                    class="w-full"
                                    type="text"
                                    :placeholder="$t('components.settings.app_config.discord.bot_presence.url')"
                                />
                            </UFormField>
                        </UPageCard>
                    </template>

                    <template #game>
                        <UPageCard
                            :title="$t('components.settings.app_config.game.max_wanted_duration.title')"
                            :description="$t('components.settings.app_config.game.max_wanted_duration.description')"
                        >
                            <UFormField
                                :label="$t('components.settings.app_config.game.max_wanted_duration.enabled_user')"
                                name="game.maxWantedDurationUserEnabled"
                            >
                                <USwitch v-model="state.game.maxWantedDurationUserEnabled" />
                            </UFormField>

                            <UFormField
                                :label="$t('components.settings.app_config.game.max_wanted_duration.max')"
                                name="game.maxWantedDurationUser"
                            >
                                <InputDurationPicker
                                    v-model="state.game.maxWantedDurationUser"
                                    class="w-full"
                                    :units="['day']"
                                    :min="secondsToDuration(24 * 60 * 60)"
                                    :step="1"
                                    :max="secondsToDuration(3650 * 24 * 60 * 60)"
                                />
                            </UFormField>

                            <UFormField
                                :label="$t('components.settings.app_config.game.max_wanted_duration.enabled_vehicle')"
                                name="game.maxWantedDurationVehicleEnabled"
                            >
                                <USwitch v-model="state.game.maxWantedDurationVehicleEnabled" />
                            </UFormField>

                            <UFormField
                                :label="$t('components.settings.app_config.game.max_wanted_duration.max')"
                                name="game.maxWantedDurationVehicle"
                            >
                                <InputDurationPicker
                                    v-model="state.game.maxWantedDurationVehicle"
                                    class="w-full"
                                    :units="['day']"
                                    :min="secondsToDuration(24 * 60 * 60)"
                                    :step="1"
                                    :max="secondsToDuration(3650 * 24 * 60 * 60)"
                                />
                            </UFormField>
                        </UPageCard>

                        <UPageCard
                            :title="$t('components.settings.app_config.quick_buttons.penalty_calculator.title')"
                            :description="$t('components.settings.app_config.quick_buttons.penalty_calculator.description')"
                        >
                            <UFormField
                                :label="$t('components.settings.app_config.quick_buttons.penalty_calculator.max_count')"
                                name="quickButtons.penaltyCalculator.maxCount"
                            >
                                <UInputNumber
                                    v-model="state.quickButtons.penaltyCalculator.maxCount"
                                    class="w-full"
                                    :min="1"
                                    :max="100"
                                    :placeholder="
                                        $t('components.settings.app_config.quick_buttons.penalty_calculator.max_count')
                                    "
                                />
                            </UFormField>

                            <UFormField
                                :label="$t('components.settings.app_config.quick_buttons.penalty_calculator.max_leeway')"
                                name="quickButtons.penaltyCalculator.maxLeeway"
                            >
                                <UInputNumber
                                    v-model="state.quickButtons.penaltyCalculator.maxLeeway"
                                    class="w-full"
                                    :min="1"
                                    :step="1"
                                    :max="99"
                                    :format-options="{ style: 'percent', minimumFractionDigits: 0, maximumFractionDigits: 0 }"
                                />
                            </UFormField>

                            <UFormField
                                :label="
                                    $t(
                                        'components.settings.app_config.quick_buttons.penalty_calculator.detention_time_unit.title',
                                    )
                                "
                                :description="
                                    $t(
                                        'components.settings.app_config.quick_buttons.penalty_calculator.detention_time_unit.description',
                                    )
                                "
                                :ui="{ container: 'flex w-full flex-row gap-2' }"
                            >
                                <UFormField
                                    class="flex-1"
                                    :label="$t('common.singular')"
                                    name="quickButtons.penaltyCalculator.detentionTimeUnit.singular"
                                >
                                    <UInput
                                        v-model="state.quickButtons.penaltyCalculator.detentionTimeUnit!.singular"
                                        class="w-full"
                                        type="text"
                                    />
                                </UFormField>

                                <UFormField
                                    class="flex-1"
                                    :label="$t('common.plural')"
                                    name="quickButtons.penaltyCalculator.detentionTimeUnit.plural"
                                >
                                    <UInput
                                        v-model="state.quickButtons.penaltyCalculator.detentionTimeUnit!.plural"
                                        class="w-full"
                                        type="text"
                                    />
                                </UFormField>
                            </UFormField>

                            <UFormField
                                class="flex-1"
                                :label="
                                    $t('components.settings.app_config.quick_buttons.penalty_calculator.warn_settings.title')
                                "
                                :description="
                                    $t(
                                        'components.settings.app_config.quick_buttons.penalty_calculator.warn_settings.description',
                                    )
                                "
                            >
                                <UFormField name="quickButtons.penaltyCalculator.warnSettings.enabled">
                                    <USwitch v-model="state.quickButtons.penaltyCalculator.warnSettings!.enabled" />
                                </UFormField>

                                <UFormField
                                    :label="
                                        $t(
                                            'components.settings.app_config.quick_buttons.penalty_calculator.warn_settings.thresholds.title',
                                        )
                                    "
                                    :description="
                                        $t(
                                            'components.settings.app_config.quick_buttons.penalty_calculator.warn_settings.thresholds.description',
                                        )
                                    "
                                >
                                    <div class="flex gap-2 sm:flex-row">
                                        <UFormField
                                            class="flex-1"
                                            :label="$t('common.fine', 2)"
                                            name="quickButtons.penaltyCalculator.warnSettings.fine"
                                        >
                                            <UInputNumber
                                                v-model="state.quickButtons.penaltyCalculator.warnSettings!.fine"
                                                class="w-full"
                                                :min="0"
                                                :step="1000"
                                                :format-options="{
                                                    style: 'currency',
                                                    currency: display.currencyName,
                                                    currencyDisplay: 'code',
                                                    currencySign: 'accounting',
                                                    maximumFractionDigits: 0,
                                                }"
                                            />
                                        </UFormField>

                                        <UFormField
                                            class="flex-1"
                                            :label="$t('common.detention_time', 2)"
                                            name="quickButtons.penaltyCalculator.warnSettings.detentionTime"
                                        >
                                            <UInputNumber
                                                v-model="state.quickButtons.penaltyCalculator.warnSettings!.detentionTime"
                                                class="w-full"
                                                :min="0"
                                                :step="1"
                                            />
                                        </UFormField>

                                        <UFormField
                                            class="flex-1"
                                            :label="$t('common.traffic_infraction_points', 2)"
                                            name="quickButtons.penaltyCalculator.warnSettings.stvoPoints"
                                        >
                                            <UInputNumber
                                                v-model="state.quickButtons.penaltyCalculator.warnSettings!.stvoPoints"
                                                class="w-full"
                                                :min="0"
                                                :step="1"
                                            />
                                        </UFormField>
                                    </div>
                                </UFormField>

                                <UFormField
                                    :label="
                                        $t(
                                            'components.settings.app_config.quick_buttons.penalty_calculator.warn_settings.warn_message.title',
                                        )
                                    "
                                    :description="
                                        $t(
                                            'components.settings.app_config.quick_buttons.penalty_calculator.warn_settings.warn_message.description',
                                        )
                                    "
                                    name="quickButtons.penaltyCalculator.warnSettings.warnMessage"
                                >
                                    <UTextarea
                                        v-model="state.quickButtons.penaltyCalculator.warnSettings!.warnMessage"
                                        class="w-full"
                                        autoresize
                                        :rows="5"
                                    />
                                </UFormField>
                            </UFormField>
                        </UPageCard>

                        <UPageCard
                            :title="$t('components.settings.app_config.livemap.title')"
                            :description="$t('components.settings.app_config.livemap.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="livemap.enableCayoPerico"
                                :label="$t('components.settings.app_config.livemap.enable_cayo_perico')"
                            >
                                <USwitch v-model="state.livemap.enableCayoPerico" />
                            </UFormField>
                        </UPageCard>
                    </template>

                    <template #website>
                        <UPageCard
                            :title="$t('components.settings.app_config.display.title')"
                            :description="$t('components.settings.app_config.display.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="defaultLocale"
                                :label="$t('common.default_lang')"
                            >
                                <USelectMenu
                                    v-model="state.defaultLocale"
                                    class="w-full"
                                    :placeholder="$t('common.language', 1)"
                                    :items="locales"
                                    label-key="name"
                                    :icon="state.defaultLocale?.icon"
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="display.intlLocale"
                                :label="$t('components.settings.app_config.display.intl_locale')"
                            >
                                <USelectMenu
                                    v-model="state.display.intlLocale"
                                    class="w-full"
                                    :placeholder="$t('components.settings.app_config.display.intl_locale')"
                                    :items="intlLocales"
                                    label-key="name"
                                    value-key="code"
                                    :icon="
                                        state.display.intlLocale
                                            ? intlLocales?.find((l) => l.code === state.display.intlLocale)?.icon
                                            : undefined
                                    "
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="display.currencyName"
                                :label="$t('common.currency')"
                            >
                                <USelectMenu
                                    v-model="state.display.currencyName"
                                    class="w-full"
                                    :placeholder="$t('common.currency')"
                                    :items="currencies"
                                    label-key="name"
                                    value-key="code"
                                    :icon="
                                        state.display.currencyName
                                            ? currencies?.find((c) => c.code === state.display.currencyName)?.flag
                                            : undefined
                                    "
                                />
                            </UFormField>
                        </UPageCard>

                        <UPageCard
                            :title="$t('components.settings.app_config.website.title')"
                            :description="$t('components.settings.app_config.website.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="website.links.privacyPolicy"
                                :label="$t('common.privacy_policy')"
                            >
                                <UInput
                                    v-model="state.website.links.privacyPolicy"
                                    class="w-full"
                                    type="text"
                                    :placeholder="$t('common.privacy_policy')"
                                    maxlength="255"
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="website.links.imprint"
                                :label="$t('common.imprint')"
                            >
                                <UInput
                                    v-model="state.website.links.imprint"
                                    class="w-full"
                                    type="text"
                                    :placeholder="$t('common.imprint')"
                                    maxlength="255"
                                />
                            </UFormField>
                        </UPageCard>

                        <UPageCard :title="$t('pages.stats.title')" :description="$t('pages.stats.subtitle')">
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="website.statsPage"
                                :label="$t('common.enabled')"
                            >
                                <USwitch v-model="state.website.statsPage" />
                            </UFormField>
                        </UPageCard>
                    </template>

                    <template #system>
                        <UPageCard
                            :title="$t('components.settings.app_config.system.banner_message.title')"
                            :description="$t('components.settings.app_config.system.banner_message.subtitle')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="system.bannerMessageEnabled"
                                :label="$t('common.enabled')"
                            >
                                <USwitch v-model="state.system.bannerMessageEnabled" />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="system.bannerMessage.title"
                                :label="$t('common.message')"
                                :ui="{ container: '', error: 'hidden' }"
                            >
                                <TiptapEditor
                                    v-model="state.system.bannerMessage.title"
                                    name="system.bannerMessage.title"
                                    :limit="1024"
                                />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="system.bannerMessage.expiresAt"
                                :label="$t('common.expires_at')"
                            >
                                <InputDatePicker v-model="state.system.bannerMessage.expiresAt" clearable time />
                            </UFormField>
                        </UPageCard>
                    </template>
                </UTabs>
            </UForm>
        </template>
    </UDashboardPanel>
</template>
