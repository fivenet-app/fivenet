<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">
import '~/assets/css/herofull-pattern.css';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import PageFooter from './components/partials/PageFooter.vue';

useHead({
    title: 'Error occured - FiveNet',
});

const props = defineProps<{
    error: Error | object | undefined;
}>();

const router = useRouter();
const route = router.currentRoute;

const buttonDisabled = ref(true);

onMounted(() => useTimeoutFn(() => (buttonDisabled.value = false), 2000));

async function handleError(url?: string): Promise<void> {
    if (url === undefined) {
        url = '/';
    }

    await clearError();
    reloadNuxtApp({
        path: url,
        persistState: false,
        ttl: 2000,
    });
}

const version = APP_VERSION;

function copyError(): void {
    if (!props.error) {
        return;
    }

    copyToClipboardWrapper(`**App Error occured - ${new Date().toLocaleString()}**
\`\`\`
${props.error ? JSON.stringify(props.error) : 'Unknown error'}
\`\`\`
**Version:** ${version}
`);
}

const kbdBlockClasses =
    'inline-flex items-center rounded-sm bg-neutral-100 px-1 text-gray-900 ring-1 ring-inset ring-gray-300 dark:bg-neutral-800 dark:text-white dark:ring-gray-700';

const isDev = import.meta.dev;
</script>

<!-- eslint-disable tailwindcss/no-custom-classname -->
<template>
    <div class="h-dvh">
        <div class="hero absolute inset-0 z-[-1] mask-[radial-gradient(100%_100%_at_top,white,transparent)]" />
        <NuxtLoadingIndicator color="repeating-linear-gradient(to right, #d72638 0%, #ac1e2d 50%, #d72638 100%)" />

        <div class="flex h-full flex-col items-center justify-center">
            <UButton class="absolute top-4 z-10" icon="i-mdi-home" :label="$t('common.home')" to="/" color="neutral" />

            <UCard class="w-full max-w-md bg-white/75 backdrop-blur-sm dark:bg-white/5">
                <template #header>
                    <FiveNetLogo class="mx-auto mb-2 h-auto w-20" />

                    <h1 class="text-center text-4xl font-bold">
                        {{ $t !== undefined ? $t('pages.error.title') : 'Error occured' }}
                    </h1>
                </template>

                <div class="flex flex-col gap-1">
                    <p class="text-base">
                        {{
                            $t !== undefined
                                ? $t('pages.error.subtitle')
                                : 'A fatal error occured, please try again in a few seconds.'
                        }}
                    </p>
                </div>

                <div class="flex flex-col gap-1">
                    <div class="flex flex-row gap-1">
                        <p>
                            <span class="font-semibold"
                                >{{ $t !== undefined ? $t('components.debug_info.version') : 'Version' }}:</span
                            >
                        </p>

                        <pre class="text-wrap" :class="kbdBlockClasses">{{ version }}</pre>
                    </div>

                    <p class="font-semibold">{{ $t !== undefined ? $t('pages.error.error_message') : 'Error message' }}:</p>
                    <span v-if="error">
                        <!-- @vue-ignore -->
                        <pre
                            v-if="error.statusMessage"
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

                <template #footer>
                    <div class="flex flex-col gap-2">
                        <div class="flex w-full gap-2">
                            <UButton
                                class="flex-1"
                                color="primary"
                                block
                                size="lg"
                                :disabled="buttonDisabled"
                                @click="() => handleError()"
                            >
                                {{ $t('common.home') }}
                            </UButton>

                            <UButton
                                class="flex-1"
                                block
                                size="lg"
                                color="green"
                                :disabled="buttonDisabled"
                                @click="() => handleError(route.fullPath)"
                            >
                                {{ $t('common.retry') }}
                            </UButton>

                            <!-- @vue-ignore -->
                            <UButton
                                v-if="error && (error.statusMessage || error.message)"
                                class="flex-1"
                                block
                                size="lg"
                                color="warning"
                                @click="() => copyError()"
                            >
                                {{ $t !== undefined ? $t('pages.error.copy_error') : 'Copy Error message' }}
                            </UButton>
                        </div>

                        <UButton
                            v-if="isDev"
                            @click="
                                updateAppConfig({ version: 'UNKNOWN' });
                                clearError();
                            "
                        >
                            Set Dev App Config
                        </UButton>
                    </div>
                </template>
            </UCard>
        </div>

        <PageFooter />
    </div>
</template>
