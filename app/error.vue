<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">
import '~/assets/css/herofull-pattern.css';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

useHead({
    title: 'Error occured - FiveNet',
});

const props = defineProps<{
    error: Error | object | undefined;
}>();

const router = useRouter();
const route = router.currentRoute;

// This page must still render when app plugins (including i18n) failed to initialize.
const nuxtApp = useNuxtApp();
function translate(key: string, fallback: string): string {
    try {
        const translator = (nuxtApp as { $t?: unknown }).$t;
        return typeof translator === 'function' ? String(translator(key)) : fallback;
    } catch {
        return fallback;
    }
}

const buttonsDisabled = ref(true);
const handlingError = ref(false);

async function handleError(url = '/'): Promise<void> {
    if (handlingError.value) return;

    handlingError.value = true;
    try {
        await clearError();
        reloadNuxtApp({
            path: url,
            persistState: false,
            ttl: 2000,
        });
    } catch {
        handlingError.value = false;
    }
}

const version = APP_VERSION;

function copyError(): void {
    if (!props.error) return;

    void copyToClipboardWrapper(`**App Error occured - ${new Date().toLocaleString()}**
\`\`\`
${props.error ? JSON.stringify(props.error) : 'Unknown error'}
\`\`\`
**Version:** ${version}
`).catch(() => undefined);
}

function setDevConfig(): void {
    updateAppConfig({ version: 'UNKNOWN' });
    clearError();
}

const kbdBlockClasses =
    'inline-flex items-center rounded-sm bg-neutral-100 px-1 text-gray-900 ring-1 ring-inset ring-gray-300 dark:bg-neutral-800 dark:text-white dark:ring-gray-700';

const showClearSiteData = ref<boolean>(false);

onMounted(() => {
    useTimeoutFn(() => (buttonsDisabled.value = false), 2000);
    useTimeoutFn(() => (showClearSiteData.value = true), 6500);
});

const isDev = import.meta.dev;
</script>

<!-- eslint-disable tailwindcss/no-custom-classname -->
<template>
    <div class="relative isolate min-h-dvh overflow-hidden">
        <div class="hero pointer-events-none absolute inset-0 z-[-1]" />

        <NuxtLoadingIndicator color="repeating-linear-gradient(to right, #d72638 0%, #ac1e2d 50%, #d72638 100%)" />

        <div class="flex min-h-dvh flex-col items-center justify-center">
            <UButton
                class="absolute top-4 z-10"
                icon="i-mdi-home"
                :label="translate('common.home', 'Home')"
                to="/"
                color="neutral"
            />

            <UCard class="w-full max-w-md bg-white/75 pt-10 pb-4 backdrop-blur-sm dark:bg-white/5">
                <template #header>
                    <FiveNetLogo class="mx-auto mb-2 h-auto w-20" />

                    <h1 class="text-center text-4xl font-bold">
                        {{ translate('pages.error.title', 'Error occured') }}
                    </h1>

                    <p class="text-center text-lg">
                        {{ translate('pages.error.subtitle', 'A fatal error occured, please try again in a few seconds.') }}
                    </p>
                </template>

                <div class="flex flex-col items-center gap-1">
                    <div class="inline-flex flex-col gap-1">
                        <p class="text-center font-semibold">{{ translate('components.debug_info.version', 'Version') }}:</p>

                        <pre class="text-wrap" :class="kbdBlockClasses">{{ version }}</pre>
                    </div>

                    <div class="inline-flex flex-col gap-1">
                        <p class="text-center font-semibold">{{ translate('pages.error.error_message', 'Error message') }}:</p>

                        <span v-if="error">
                            <!-- @vue-ignore -->
                            <pre
                                v-if="error.statusMessage"
                                class="text-wrap"
                                :class="kbdBlockClasses"
                                v-text="
                                    // @ts-expect-error
                                    error.statusMessage
                                "
                            />
                            <!-- @vue-ignore -->
                            <pre
                                v-else-if="
                                    // @ts-expect-error
                                    error.message
                                "
                                class="text-wrap"
                                :class="kbdBlockClasses"
                                v-text="
                                    // @ts-expect-error
                                    error.message
                                "
                            />
                            <pre v-else>Unable to get error message</pre>
                        </span>
                        <span v-else>
                            <pre>Unknown error</pre>
                        </span>
                    </div>
                </div>

                <template #footer>
                    <div class="flex flex-col gap-2">
                        <div class="grid w-full grid-cols-3 gap-2">
                            <UButton
                                class=""
                                color="primary"
                                icon="i-mdi-home"
                                size="lg"
                                :disabled="buttonsDisabled || handlingError"
                                :loading="handlingError"
                                :label="translate('common.home', 'Home')"
                                @click="() => handleError()"
                            />

                            <UButton
                                class="col-span-2 truncate"
                                color="success"
                                icon="i-mdi-refresh"
                                size="lg"
                                :disabled="buttonsDisabled || handlingError"
                                :loading="handlingError"
                                :label="translate('common.retry', 'Retry')"
                                @click="() => handleError(route.fullPath)"
                            />
                        </div>

                        <!-- @vue-ignore -->
                        <UButton
                            v-if="error && (error.statusMessage || error.message)"
                            class="col-span-1 truncate"
                            color="warning"
                            icon="i-mdi-content-copy"
                            size="lg"
                            :label="translate('pages.error.copy_error', 'Copy Error message')"
                            @click="() => copyError()"
                        />

                        <USeparator v-if="showClearSiteData || isDev" class="my-1" />

                        <UButton
                            v-if="showClearSiteData"
                            class="col-span-1 truncate"
                            color="error"
                            icon="i-mdi-restart-alert"
                            size="lg"
                            :label="translate('components.debug_info.factory_reset', 'Factory Reset FiveNet App')"
                            variant="soft"
                            external
                            to="/api/clear-site-data"
                        />

                        <UButton v-if="isDev" label="Set Dev App Config" @click="() => setDevConfig()" />
                    </div>
                </template>
            </UCard>
        </div>
    </div>
</template>
