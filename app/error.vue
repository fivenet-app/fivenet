<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">
import '~/assets/css/herofull-pattern.css';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

useHead({
    title: 'Error occured - FiveNet',
});

const props = defineProps<{
    error: Error | object;
}>();

const router = useRouter();
const route = router.currentRoute;

const buttonDisabled = ref(true);

const { start } = useTimeoutFn(() => (buttonDisabled.value = false), 2000);

function handleError(url?: string): void {
    start();

    if (url === undefined) {
        url = '/';
    }

    clearError();
    reloadNuxtApp({
        path: url,
        persistState: false,
        ttl: 2000,
    });
}

function copyError(): void {
    if (!props.error) {
        return;
    }

    copyToClipboardWrapper(JSON.stringify(props.error));
}

const isDev = import.meta.dev;
</script>

<template>
    <div class="h-dscreen">
        <NuxtLoadingIndicator color="repeating-linear-gradient(to right, #d72638 0%, #ac1e2d 50%, #d72638 100%)" />

        <div class="h-full">
            <div class="hero absolute inset-0 z-[-1] [mask-image:radial-gradient(100%_100%_at_top,white,transparent)]" />

            <div class="h-full">
                <main class="mx-auto flex size-full max-w-3xl text-center">
                    <div class="my-auto py-2 max-sm:w-full sm:mx-auto lg:mx-auto">
                        <FiveNetLogo class="mx-auto mb-2 h-auto w-36" />

                        <h1 class="text-5xl font-bold">
                            {{ $t ? $t('pages.error.title') : 'Error occured' }}
                        </h1>
                        <h2 class="text-xl">
                            {{ $t ? $t('pages.error.subtitle') : 'A fatal error occured, please try again in a few seconds.' }}
                        </h2>

                        <div class="mb-4 py-2">
                            <p class="py-2 font-semibold">
                                {{ $t ? $t('pages.error.error_message') : 'Error message:' }}
                            </p>
                            <span v-if="error">
                                <!-- @vue-ignore -->
                                <pre
                                    v-if="error.statusMessage"
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

                        <div class="inline-flex w-full gap-2">
                            <UButton
                                color="primary"
                                block
                                class="flex-1"
                                size="lg"
                                :disabled="buttonDisabled"
                                @click="handleError()"
                            >
                                {{ $t('common.home') }}
                            </UButton>

                            <UButton
                                block
                                class="flex-1"
                                size="lg"
                                color="green"
                                :disabled="buttonDisabled"
                                @click="handleError(route.fullPath)"
                            >
                                {{ $t('common.retry') }}
                            </UButton>

                            <!-- @vue-ignore -->
                            <UButton
                                v-if="error && (error.statusMessage || error.message)"
                                block
                                class="flex-1"
                                size="lg"
                                color="amber"
                                @click="copyError"
                            >
                                {{ $t ? $t('pages.error.copy_error') : 'Copy Error message' }}
                            </UButton>
                        </div>

                        <UButton
                            v-if="isDev"
                            class="mt-4"
                            @click="
                                updateAppConfig({ version: 'UNKNOWN' });
                                clearError();
                            "
                        >
                            Set Dev App Config
                        </UButton>
                    </div>
                </main>
            </div>
        </div>
    </div>
</template>
