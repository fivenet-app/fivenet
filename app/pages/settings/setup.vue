<script lang="ts" setup>
import type { FormSubmitEvent, StepperItem } from '@nuxt/ui';
import type { LocaleObject } from '@nuxtjs/i18n';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';
import { useCompletorStore } from '~/stores/completor';
import { useSettingsStore } from '~/stores/settings';
import { deepToRaw } from '~/utils/deepToRaw';
import { currencies, intlLocales } from '~/components/settings/helpers';
import { getSettingsConfigClient } from '~~/gen/ts/clients';
import { GRPCServiceMethods, GRPCServices } from '~~/gen/ts/perms';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { AppConfig as SettingsAppConfig } from '~~/gen/ts/resources/settings/config';
import type { GetAppConfigResponse } from '~~/gen/ts/services/settings/config';
import UserMenu from '~/components/UserMenu.vue';
import type { RoutePathSchema } from '@typed-router';

useHead({
    title: 'pages.settings.setup.title',
});

definePageMeta({
    title: 'pages.settings.setup.title',
    requiresAuth: true,
    authTokenOnly: true,
    layout: 'empty',
});

const { t, locales } = useI18n();
const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const { $reloadConfigFromServer, $applyClientConfig } = useNuxtApp();

type SetupStep = 'access' | 'locale' | 'site' | 'review';

const setupStepValues: SetupStep[] = ['locale', 'access', 'site', 'review'];
const route = useRoute();
const router = useRouter();

const selectedStep = computed<SetupStep>({
    get() {
        const candidate = route.query.step;

        if (typeof candidate === 'string' && setupStepValues.includes(candidate as SetupStep)) {
            return candidate as SetupStep;
        }

        return 'locale';
    },
    set(step) {
        if (route.query.step === step) return;

        void router.replace({
            query: {
                ...route.query,
                step,
            },
        });
    },
});

watchEffect(() => {
    const candidate = route.query.step;

    if (candidate !== undefined && (typeof candidate !== 'string' || !setupStepValues.includes(candidate as SetupStep))) {
        void router.replace({
            query: {
                ...route.query,
                step: 'locale',
            },
        });
    }
});

const setupSteps = computed<StepperItem[]>(() => [
    {
        value: 'locale',
        slot: 'locale',
        icon: 'i-mdi-translate',
        title: t('pages.settings.setup.steps.locale.title'),
        description: t('pages.settings.setup.steps.locale.description'),
    },
    {
        value: 'access',
        slot: 'access',
        icon: 'i-mdi-shield-account',
        title: t('pages.settings.setup.steps.access.title'),
        description: t('pages.settings.setup.steps.access.description'),
    },
    {
        value: 'site',
        slot: 'site',
        icon: 'i-mdi-web',
        title: t('pages.settings.setup.steps.site.title'),
        description: t('pages.settings.setup.steps.site.description'),
    },
    {
        value: 'review',
        slot: 'review',
        icon: 'i-mdi-clipboard-check-outline',
        title: t('pages.settings.setup.steps.review.title'),
        description: t('pages.settings.setup.steps.review.description'),
    },
]);

const currentStepIndex = computed(() => setupStepValues.indexOf(selectedStep.value));
const currentStep = computed(() => setupSteps.value[currentStepIndex.value] ?? setupSteps.value[0]);
const isFirstStep = computed(() => currentStepIndex.value <= 0);
const isLastStep = computed(() => currentStepIndex.value === setupSteps.value.length - 1);

const stepIndicator = computed(() =>
    t('pages.settings.setup.step_indicator', [currentStepIndex.value + 1, setupSteps.value.length]),
);

const schema = z.object({
    defaultLocale: z.custom<LocaleObject>(),

    auth: z.object({
        signupEnabled: z.coerce.boolean(),
        lastCharLock: z.coerce.boolean(),
        configAdminGroups: z.coerce.string().array().max(10).default([]),
        configAdminUsers: z.coerce.string().array().max(10).default([]),
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
    display: z.object({
        intlLocale: z.string().default('en-US'),
        currencyName: z.string().default('USD'),
    }),
    jobInfo: z.object({
        unemployedJob: z.object({
            name: z.coerce.string().min(1).max(20),
            grade: z.coerce.number().min(0).max(99),
        }),
        publicJobs: z.coerce.string().array().max(99).default([]),
        hiddenJobs: z.coerce.string().array().max(99).default([]),
    }),
    website: z.object({
        links: z.object({
            privacyPolicy: z.union([z.url().min(1).max(255).startsWith('https://'), z.string().length(0).optional()]),
            imprint: z.union([z.url().min(1).max(255).startsWith('https://'), z.string().length(0).optional()]),
        }),
        statsPage: z.coerce.boolean(),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    defaultLocale: locales.value[0]!,
    auth: {
        signupEnabled: true,
        lastCharLock: false,
        configAdminGroups: [],
        configAdminUsers: [],
    },
    perms: {
        default: [],
    },
    display: {
        intlLocale: 'en-US',
        currencyName: 'USD',
    },
    jobInfo: {
        unemployedJob: {
            name: 'unemployed',
            grade: 1,
        },
        publicJobs: [],
        hiddenJobs: [],
    },
    website: {
        links: {
            privacyPolicy: '',
            imprint: '',
        },
        statsPage: false,
    },
});

const settingsConfigClient = await getSettingsConfigClient();
const { data: config, status, refresh, error } = useLazyAsyncData(`settings-setup-appconfig`, () => getAppConfig());

const completorStore = useCompletorStore();
const { listJobs } = completorStore;
const { data: jobs } = useLazyAsyncData(`settings-setup-jobs`, () => listJobs());

const notifications = useNotificationsStore();
const isSubmitting = ref(false);

const redirectTarget = computed<RoutePathSchema>(() => {
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/overview';
    return getRedirectPath(redirect) as RoutePathSchema;
});

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

function setSettingsValues(): void {
    if (!config.value?.config) return;

    const current = config.value.config;
    const locale = locales.value.find((l) => l.code === current.defaultLocale);

    if (locale) {
        state.defaultLocale = locale;
    }

    state.auth = {
        signupEnabled: current.auth?.signupEnabled ?? true,
        lastCharLock: current.auth?.lastCharLock ?? false,
        configAdminGroups: current.auth?.configAdminGroups ?? [],
        configAdminUsers: current.auth?.configAdminUsers ?? [],
    };

    state.perms = {
        default: current.perms?.default?.map((perm) => ({ ...perm })) ?? [],
    };

    state.display = {
        intlLocale: current.display?.intlLocale ?? 'en-US',
        currencyName: current.display?.currencyName ?? 'USD',
    };

    state.jobInfo = {
        unemployedJob: {
            name: current.jobInfo?.unemployedJob?.name ?? 'unemployed',
            grade: current.jobInfo?.unemployedJob?.grade ?? 1,
        },
        publicJobs: current.jobInfo?.publicJobs ?? [],
        hiddenJobs: current.jobInfo?.hiddenJobs ?? [],
    };

    state.website = {
        links: {
            privacyPolicy: current.website?.links?.privacyPolicy ?? '',
            imprint: current.website?.links?.imprint ?? '',
        },
        statsPage: current.website?.statsPage ?? false,
    };
}

watch(config, () => setSettingsValues(), { immediate: true });

async function saveAppConfig(values: Schema): Promise<void> {
    if (!config.value?.config) return;

    const payload = structuredClone(deepToRaw(config.value.config)) as SettingsAppConfig;
    payload.setupComplete = true;
    payload.defaultLocale = values.defaultLocale.code;
    payload.auth = {
        ...(payload.auth ?? {
            signupEnabled: true,
            lastCharLock: false,
            jobAdminGroups: [],
            jobAdminUsers: [],
            configAdminGroups: [],
            configAdminUsers: [],
        }),
        signupEnabled: values.auth.signupEnabled,
        lastCharLock: values.auth.lastCharLock,
        configAdminGroups: values.auth.configAdminGroups,
        configAdminUsers: values.auth.configAdminUsers,
    };
    payload.perms = values.perms;
    payload.display = values.display;
    payload.jobInfo = {
        ...(payload.jobInfo ?? {
            publicJobs: [],
            hiddenJobs: [],
            unemployedJob: {
                name: 'unemployed',
                grade: 1,
            },
        }),
        unemployedJob: values.jobInfo.unemployedJob,
        publicJobs: values.jobInfo.publicJobs,
        hiddenJobs: values.jobInfo.hiddenJobs,
    };
    payload.website = {
        ...(payload.website ?? {
            links: {},
            statsPage: false,
        }),
        links: {
            ...(payload.website?.links ?? {}),
            privacyPolicy: values.website.links.privacyPolicy?.length ? values.website.links.privacyPolicy : undefined,
            imprint: values.website.links.imprint?.length ? values.website.links.imprint : undefined,
        },
        statsPage: values.website.statsPage,
    };

    try {
        const { response } = await settingsConfigClient.updateAppConfig({
            config: payload,
        });

        notifications.add({
            title: { key: 'notifications.settings.app_config.title', parameters: {} },
            description: { key: 'notifications.settings.app_config.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (response.config) {
            config.value = response;
        } else {
            await refresh();
        }

        const clientConfig = await $reloadConfigFromServer();
        $applyClientConfig(clientConfig);

        await navigateTo(redirectTarget.value);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    isSubmitting.value = true;

    await saveAppConfig(event.data).finally(() =>
        useTimeoutFn(() => {
            canSubmit.value = true;
            isSubmitting.value = false;
        }, 400),
    );
}, 1000);

function goPrev(): void {
    if (isFirstStep.value) return;

    selectedStep.value = setupStepValues[currentStepIndex.value - 1]!;
}

function goNext(): void {
    if (isLastStep.value) return;

    selectedStep.value = setupStepValues[currentStepIndex.value + 1]!;
}

const selectedLocaleName = computed(() => state.defaultLocale.name);
const selectedCurrencyName = computed(
    () => currencies.find((currency) => currency.code === state.display.currencyName)?.name ?? state.display.currencyName,
);
const selectedIntlLocaleName = computed(
    () => intlLocales.find((locale) => locale.code === state.display.intlLocale)?.name ?? state.display.intlLocale,
);
const selectedJobSummary = computed(
    () =>
        `${state.jobInfo.unemployedJob.name} (${state.jobInfo.unemployedJob.grade}) · ${state.jobInfo.publicJobs.length} ${t('common.job', 2)} ${t('common.public')} · ${state.jobInfo.hiddenJobs.length} ${t('common.hidden')}`,
);
const permissionServiceItems = computed(() =>
    GRPCServices.map((service) => ({ label: t(`perms.${service}.service`), value: service })),
);

function permissionMethodItems(category: string): Array<{ label: string; value: string }> {
    if (!category) return [];

    return GRPCServiceMethods.filter((method) => method.startsWith(`${category}/`))
        .map((method) => method.split('/').at(1) ?? method)
        .map((name) => ({ label: t(`perms.${category}.${name}.key`), value: name }));
}

function addDefaultPermission(): void {
    state.perms.default.push({ category: '', name: '' });
}

function removeDefaultPermission(idx: number): void {
    state.perms.default.splice(idx, 1);
}

const open = ref<boolean>(false);
</script>

<template>
    <UDashboardGroup unit="rem" storage="local">
        <UDashboardSidebar
            id="default"
            v-model:open="open"
            class="bg-elevated/25"
            :default-size="16.5"
            :min-size="13.5"
            :max-size="23"
            collapsible
            resizable
            :ui="{ footer: 'lg:border-t lg:border-default' }"
        >
            <template #header="{ collapsed }">
                <TopLogoDropdown :collapsed="collapsed" hide-notifications />
            </template>

            <template #footer="{ collapsed }">
                <UserMenu :collapsed="collapsed" hide-user-menus />
            </template>
        </UDashboardSidebar>

        <UDashboardPanel :ui="{ body: 'p-0 sm:p-0' }">
            <template #body>
                <div class="relative isolate min-h-dvh">
                    <div
                        class="pointer-events-none fixed inset-0 -z-10 bg-[radial-gradient(circle_at_top_left,rgba(14,165,233,0.16),transparent_28%),radial-gradient(circle_at_bottom_right,rgba(99,102,241,0.12),transparent_32%)]"
                    />

                    <div
                        class="pointer-events-none fixed inset-0 -z-10 bg-[linear-gradient(rgba(15,23,42,0.04)_1px,transparent_1px),linear-gradient(90deg,rgba(15,23,42,0.04)_1px,transparent_1px)] bg-[size:32px_32px] opacity-40 dark:opacity-20"
                    />

                    <UContainer class="relative py-8 sm:py-10 lg:py-14">
                        <div class="mx-auto flex w-full max-w-6xl flex-col gap-6 lg:gap-8">
                            <div class="max-w-3xl space-y-3">
                                <UBadge color="primary" variant="soft" icon="i-mdi-creation">
                                    {{ $t('pages.settings.setup.eyebrow') }}
                                </UBadge>

                                <div class="space-y-2">
                                    <h1 class="text-3xl font-semibold tracking-tight text-gray-950 sm:text-4xl dark:text-white">
                                        {{ $t('pages.settings.setup.title') }}
                                    </h1>

                                    <p class="max-w-2xl text-sm leading-6 text-gray-600 dark:text-gray-300">
                                        {{ $t('pages.settings.setup.subtitle') }}
                                    </p>
                                </div>
                            </div>

                            <UCard
                                class="border border-gray-200/80 bg-white/90 shadow-2xl shadow-slate-950/10 backdrop-blur dark:border-gray-800 dark:bg-slate-950/75"
                            >
                                <template #default>
                                    <UForm
                                        v-if="!streamerMode"
                                        ref="formRef"
                                        class="flex flex-col"
                                        :schema="schema"
                                        :state="state"
                                        @submit="onSubmitThrottle"
                                    >
                                        <div
                                            v-if="isRequestPending(status)"
                                            class="space-y-1 px-4 py-4 sm:px-6 sm:py-6 lg:px-8"
                                        >
                                            <USkeleton class="mb-6 h-11 w-full" />
                                            <USkeleton v-for="idx in 4" :key="idx" class="h-36 w-full" />
                                        </div>

                                        <DataErrorBlock
                                            v-else-if="error"
                                            class="px-4 py-4 sm:px-6 sm:py-6 lg:px-8"
                                            :title="$t('common.unable_to_load', [$t('common.setting', 2)])"
                                            :error="error"
                                            :retry="refresh"
                                        />

                                        <DataNoDataBlock
                                            v-else-if="!config"
                                            class="px-4 py-4 sm:px-6 sm:py-6 lg:px-8"
                                            icon="i-mdi-office-building-cog"
                                            :type="$t('common.setting', 2)"
                                            :retry="refresh"
                                        />

                                        <div v-else class="space-y-6 p-4 sm:p-6 lg:p-8">
                                            <div class="space-y-3">
                                                <div class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
                                                    <div class="space-y-1">
                                                        <p class="text-xs font-medium tracking-[0.28em] text-primary uppercase">
                                                            {{ stepIndicator }}
                                                        </p>

                                                        <h2 class="text-xl font-semibold text-gray-950 dark:text-white">
                                                            {{ currentStep?.title }}
                                                        </h2>

                                                        <p class="max-w-3xl text-sm leading-6 text-gray-600 dark:text-gray-300">
                                                            {{ currentStep?.description }}
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>

                                            <UStepper
                                                v-model="selectedStep"
                                                class="w-full"
                                                :items="setupSteps"
                                                color="primary"
                                                size="md"
                                            >
                                                <template #locale>
                                                    <div class="flex flex-col gap-4">
                                                        <UPageCard
                                                            :title="$t('pages.settings.setup.cards.locale.language.title')"
                                                            :description="
                                                                $t('pages.settings.setup.cards.locale.language.description')
                                                            "
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
                                                        </UPageCard>

                                                        <UPageCard
                                                            :title="$t('pages.settings.setup.cards.locale.formatting.title')"
                                                            :description="
                                                                $t('pages.settings.setup.cards.locale.formatting.description')
                                                            "
                                                        >
                                                            <UFormField
                                                                class="grid grid-cols-2 items-center gap-2"
                                                                name="display.intlLocale"
                                                                :label="
                                                                    $t('components.settings.app_config.display.intl_locale')
                                                                "
                                                            >
                                                                <USelectMenu
                                                                    v-model="state.display.intlLocale"
                                                                    class="w-full"
                                                                    :placeholder="
                                                                        $t('components.settings.app_config.display.intl_locale')
                                                                    "
                                                                    :items="intlLocales"
                                                                    label-key="name"
                                                                    value-key="code"
                                                                    :icon="
                                                                        state.display.intlLocale
                                                                            ? intlLocales?.find(
                                                                                  (l) => l.code === state.display.intlLocale,
                                                                              )?.icon
                                                                            : undefined
                                                                    "
                                                                />
                                                            </UFormField>
                                                        </UPageCard>

                                                        <UPageCard
                                                            :title="$t('pages.settings.setup.cards.locale.currency.title')"
                                                            :description="
                                                                $t('pages.settings.setup.cards.locale.currency.description')
                                                            "
                                                        >
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
                                                                            ? currencies?.find(
                                                                                  (c) => c.code === state.display.currencyName,
                                                                              )?.flag
                                                                            : undefined
                                                                    "
                                                                />
                                                            </UFormField>
                                                        </UPageCard>
                                                    </div>
                                                </template>

                                                <template #access>
                                                    <div class="flex flex-col gap-4">
                                                        <UPageCard
                                                            :title="$t('pages.settings.setup.cards.access.admins.title')"
                                                            :description="
                                                                $t('pages.settings.setup.cards.access.admins.description')
                                                            "
                                                        >
                                                            <UFormField
                                                                class="grid grid-cols-2 items-center gap-2"
                                                                name="auth.configAdminGroups"
                                                                :label="
                                                                    $t(
                                                                        'components.settings.app_config.auth.admin_groups.config_admin_groups',
                                                                    )
                                                                "
                                                                :description="
                                                                    $t(
                                                                        'components.settings.app_config.auth.admin_groups.config_admin_groups_description',
                                                                    )
                                                                "
                                                            >
                                                                <UInputTags
                                                                    v-model="state.auth.configAdminGroups"
                                                                    class="w-full"
                                                                />
                                                            </UFormField>

                                                            <UFormField
                                                                class="grid grid-cols-2 items-center gap-2"
                                                                name="auth.configAdminUsers"
                                                                :label="
                                                                    $t(
                                                                        'components.settings.app_config.auth.admin_groups.config_admin_users',
                                                                    )
                                                                "
                                                                :description="
                                                                    $t(
                                                                        'components.settings.app_config.auth.admin_groups.config_admin_users_description',
                                                                    )
                                                                "
                                                            >
                                                                <UInputTags
                                                                    v-model="state.auth.configAdminUsers"
                                                                    class="w-full"
                                                                />
                                                            </UFormField>
                                                        </UPageCard>

                                                        <UPageCard
                                                            :title="$t('pages.settings.setup.cards.access.policy.title')"
                                                            :description="
                                                                $t('pages.settings.setup.cards.access.policy.description')
                                                            "
                                                        >
                                                            <UFormField
                                                                class="grid grid-cols-2 items-center gap-2"
                                                                name="auth.signupEnabled"
                                                                :label="$t('components.settings.app_config.auth.sign_up.title')"
                                                                :description="
                                                                    $t(
                                                                        'components.settings.app_config.auth.sign_up.description',
                                                                    )
                                                                "
                                                            >
                                                                <USwitch v-model="state.auth.signupEnabled" />
                                                            </UFormField>

                                                            <UFormField
                                                                class="grid grid-cols-2 items-center gap-2"
                                                                name="auth.lastCharLock"
                                                                :label="
                                                                    $t('components.settings.app_config.auth.last_char_lock')
                                                                "
                                                            >
                                                                <USwitch v-model="state.auth.lastCharLock" />
                                                            </UFormField>
                                                        </UPageCard>

                                                        <UPageCard
                                                            :title="$t('pages.settings.setup.cards.access.perms.title')"
                                                            :description="
                                                                $t('pages.settings.setup.cards.access.perms.description')
                                                            "
                                                        >
                                                            <UFormField name="perms.default">
                                                                <div class="flex flex-col gap-2">
                                                                    <p
                                                                        v-if="!state.perms.default.length"
                                                                        class="text-sm leading-6 text-gray-500 dark:text-gray-400"
                                                                    >
                                                                        {{
                                                                            $t('common.none_selected', [
                                                                                $t('common.permission', 2),
                                                                            ])
                                                                        }}
                                                                    </p>

                                                                    <div
                                                                        v-for="(_, idx) in state.perms.default"
                                                                        :key="idx"
                                                                        class="flex flex-col gap-2 sm:flex-row"
                                                                    >
                                                                        <UFormField
                                                                            class="flex-1"
                                                                            :name="`perms.default.${idx}.category`"
                                                                        >
                                                                            <ClientOnly>
                                                                                <USelectMenu
                                                                                    v-model="state.perms.default[idx]!.category"
                                                                                    class="w-full"
                                                                                    :filter-fields="['label', 'value']"
                                                                                    :items="permissionServiceItems"
                                                                                    :placeholder="$t('common.service')"
                                                                                    value-key="value"
                                                                                >
                                                                                    <template
                                                                                        v-if="
                                                                                            state.perms.default[idx]!.category
                                                                                        "
                                                                                        #default
                                                                                    >
                                                                                        {{
                                                                                            $t(
                                                                                                `perms.${state.perms.default[idx]!.category}.service`,
                                                                                            )
                                                                                        }}
                                                                                    </template>

                                                                                    <template #empty>
                                                                                        {{
                                                                                            $t('common.not_found', [
                                                                                                $t('common.service'),
                                                                                            ])
                                                                                        }}
                                                                                    </template>
                                                                                </USelectMenu>
                                                                            </ClientOnly>
                                                                        </UFormField>

                                                                        <UFormField
                                                                            class="flex-1"
                                                                            :name="`perms.default.${idx}.name`"
                                                                        >
                                                                            <USelectMenu
                                                                                v-model="state.perms.default[idx]!.name"
                                                                                class="w-full"
                                                                                :disabled="!state.perms.default[idx]?.category"
                                                                                :filter-fields="['label', 'value']"
                                                                                :items="
                                                                                    permissionMethodItems(
                                                                                        state.perms.default[idx]!.category,
                                                                                    )
                                                                                "
                                                                                :placeholder="$t('common.method')"
                                                                                value-key="value"
                                                                            >
                                                                                <template
                                                                                    v-if="state.perms.default[idx]!.name"
                                                                                    #default
                                                                                >
                                                                                    {{
                                                                                        $t(
                                                                                            `perms.${state.perms.default[idx]!.category}.${state.perms.default[idx]!.name}.key`,
                                                                                        )
                                                                                    }}
                                                                                </template>

                                                                                <template #item-label="{ item }">
                                                                                    {{
                                                                                        $t(
                                                                                            `perms.${state.perms.default[idx]!.category}.${item}.key`,
                                                                                        )
                                                                                    }}
                                                                                </template>

                                                                                <template #empty>
                                                                                    {{
                                                                                        $t('common.not_found', [
                                                                                            $t('common.method'),
                                                                                        ])
                                                                                    }}
                                                                                </template>
                                                                            </USelectMenu>
                                                                        </UFormField>

                                                                        <div>
                                                                            <UTooltip :text="$t('common.remove')">
                                                                                <UButton
                                                                                    type="button"
                                                                                    color="red"
                                                                                    icon="i-mdi-close"
                                                                                    :disabled="!canSubmit || isSubmitting"
                                                                                    @click="removeDefaultPermission(idx)"
                                                                                />
                                                                            </UTooltip>
                                                                        </div>
                                                                    </div>
                                                                </div>

                                                                <UButton
                                                                    :class="state.perms.default.length ? 'mt-2' : ''"
                                                                    type="button"
                                                                    icon="i-mdi-plus"
                                                                    :disabled="!canSubmit || isSubmitting"
                                                                    @click="addDefaultPermission"
                                                                />
                                                            </UFormField>
                                                        </UPageCard>
                                                    </div>
                                                </template>

                                                <template #site>
                                                    <div class="flex flex-col gap-4">
                                                        <UPageCard
                                                            :title="$t('pages.settings.setup.cards.site.jobs.title')"
                                                            :description="
                                                                $t('pages.settings.setup.cards.site.jobs.description')
                                                            "
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
                                                                        jobs?.find(
                                                                            (j) => j.name === state.jobInfo.unemployedJob.name,
                                                                        )?.grades?.[0]?.grade ?? 0
                                                                    "
                                                                    :max="
                                                                        jobs
                                                                            ?.find(
                                                                                (j) =>
                                                                                    j.name === state.jobInfo.unemployedJob.name,
                                                                            )
                                                                            ?.grades?.at(-1)?.grade ?? 99
                                                                    "
                                                                    name="jobInfoUnemployedGrade"
                                                                    :placeholder="$t('common.rank')"
                                                                    :label="$t('common.rank')"
                                                                />
                                                            </UFormField>
                                                        </UPageCard>

                                                        <UPageCard
                                                            :title="$t('components.settings.app_config.jobs_users.title')"
                                                            :description="
                                                                $t('components.settings.app_config.jobs_users.description')
                                                            "
                                                        >
                                                            <UFormField
                                                                class="grid grid-cols-2 items-center gap-2"
                                                                name="jobInfo.publicJobs"
                                                                :label="
                                                                    $t('components.settings.app_config.jobs_users.public_jobs')
                                                                "
                                                            >
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
                                                                        <span
                                                                            v-if="state.jobInfo.publicJobs.length"
                                                                            class="truncate"
                                                                        >
                                                                            {{ state.jobInfo.publicJobs.join(', ') }}
                                                                        </span>
                                                                        <span v-else class="truncate">
                                                                            {{ $t('common.none_selected', [$t('common.job')]) }}
                                                                        </span>
                                                                    </template>

                                                                    <template #item-label="{ item }">
                                                                        {{ item.label }} ({{ item.name }})
                                                                    </template>
                                                                </USelectMenu>
                                                            </UFormField>

                                                            <UFormField
                                                                class="grid grid-cols-2 items-center gap-2"
                                                                name="jobInfo.hiddenJobs"
                                                                :label="
                                                                    $t('components.settings.app_config.jobs_users.hidden_jobs')
                                                                "
                                                            >
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
                                                                        <span
                                                                            v-if="state.jobInfo.hiddenJobs.length"
                                                                            class="truncate"
                                                                        >
                                                                            {{ state.jobInfo.hiddenJobs.join(', ') }}
                                                                        </span>
                                                                        <span v-else class="truncate">
                                                                            {{ $t('common.none_selected', [$t('common.job')]) }}
                                                                        </span>
                                                                    </template>

                                                                    <template #item-label="{ item }">
                                                                        {{ item.label }} ({{ item.name }})
                                                                    </template>
                                                                </USelectMenu>
                                                            </UFormField>
                                                        </UPageCard>

                                                        <UPageCard
                                                            class="lg:col-span-2"
                                                            :title="$t('pages.settings.setup.cards.site.website.title')"
                                                            :description="
                                                                $t('pages.settings.setup.cards.site.website.description')
                                                            "
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

                                                            <UFormField
                                                                class="grid grid-cols-2 items-center gap-2"
                                                                name="website.statsPage"
                                                                :label="$t('pages.stats.title')"
                                                            >
                                                                <USwitch v-model="state.website.statsPage" />
                                                            </UFormField>
                                                        </UPageCard>
                                                    </div>
                                                </template>

                                                <template #review>
                                                    <div class="space-y-4">
                                                        <UPageCard
                                                            :title="$t('pages.settings.setup.cards.review.summary.title')"
                                                            :description="
                                                                $t('pages.settings.setup.cards.review.summary.description')
                                                            "
                                                        >
                                                            <dl
                                                                class="space-y-3 text-sm leading-6 text-gray-600 dark:text-gray-300"
                                                            >
                                                                <div class="flex items-start justify-between gap-4">
                                                                    <dt class="text-gray-500 dark:text-gray-400">
                                                                        {{ $t('common.default_lang') }}
                                                                    </dt>
                                                                    <dd class="font-medium text-gray-950 dark:text-white">
                                                                        {{ selectedLocaleName }}
                                                                    </dd>
                                                                </div>

                                                                <div class="flex items-start justify-between gap-4">
                                                                    <dt class="text-gray-500 dark:text-gray-400">
                                                                        {{
                                                                            $t(
                                                                                'components.settings.app_config.display.intl_locale',
                                                                            )
                                                                        }}
                                                                    </dt>
                                                                    <dd class="font-medium text-gray-950 dark:text-white">
                                                                        {{ selectedIntlLocaleName }}
                                                                    </dd>
                                                                </div>

                                                                <div class="flex items-start justify-between gap-4">
                                                                    <dt class="text-gray-500 dark:text-gray-400">
                                                                        {{ $t('common.currency') }}
                                                                    </dt>
                                                                    <dd class="font-medium text-gray-950 dark:text-white">
                                                                        {{ selectedCurrencyName }}
                                                                    </dd>
                                                                </div>

                                                                <div class="flex items-start justify-between gap-4">
                                                                    <dt class="text-gray-500 dark:text-gray-400">
                                                                        {{ $t('pages.settings.setup.cards.site.jobs.title') }}
                                                                    </dt>
                                                                    <dd class="font-medium text-gray-950 dark:text-white">
                                                                        {{ selectedJobSummary }}
                                                                    </dd>
                                                                </div>
                                                            </dl>
                                                        </UPageCard>

                                                        <UAlert
                                                            :title="$t('pages.settings.setup.cards.review.discord.title')"
                                                            :description="
                                                                $t('pages.settings.setup.cards.review.discord.description')
                                                            "
                                                            color="info"
                                                            icon="i-simple-icons-discord"
                                                            variant="subtle"
                                                        />
                                                    </div>
                                                </template>
                                            </UStepper>

                                            <div
                                                class="flex flex-col gap-3 border-t border-gray-200/70 pt-6 sm:flex-row sm:items-center sm:justify-between dark:border-gray-800"
                                            >
                                                <p class="max-w-2xl text-sm leading-6 text-gray-500 dark:text-gray-400">
                                                    {{ $t('pages.settings.setup.footer') }}
                                                </p>

                                                <div class="flex flex-col gap-3 sm:flex-row">
                                                    <UButton
                                                        color="neutral"
                                                        variant="soft"
                                                        leading-icon="i-mdi-arrow-left"
                                                        :label="$t('common.previous')"
                                                        :disabled="isFirstStep || isSubmitting"
                                                        @click="goPrev"
                                                    />

                                                    <UButton
                                                        :color="isLastStep ? 'success' : 'primary'"
                                                        trailing-icon="i-mdi-arrow-right"
                                                        :label="
                                                            isLastStep
                                                                ? $t('pages.settings.setup.review.finish_action')
                                                                : $t('common.next')
                                                        "
                                                        :type="isLastStep ? 'submit' : 'button'"
                                                        :loading="isSubmitting"
                                                        :disabled="!canSubmit || isSubmitting"
                                                        @click="!isLastStep ? goNext() : undefined"
                                                    />
                                                </div>
                                            </div>
                                        </div>
                                    </UForm>

                                    <StreamerModeAlert v-else />
                                </template>
                            </UCard>
                        </div>
                    </UContainer>
                </div>
            </template>
        </UDashboardPanel>
    </UDashboardGroup>
</template>
