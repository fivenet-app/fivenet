<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import '~/assets/css/herofull-pattern.css';

useHead({
    title: 'Error occured - FiveNet',
});

const props = defineProps<{
    error: Error | Object;
}>();

const route = useRoute();

const buttonDisabled = ref(true);

const { start } = useTimeoutFn(() => (buttonDisabled.value = false), 2000, { immediate: true });

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

    copyToClipboardWrapper(jsonStringify(props.error));
}
</script>

<template>
    <div class="h-dscreen">
        <NuxtLoadingIndicator />

        <div class="hero h-full bg-base-900">
            <div class="hero-overlay h-full">
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
                                        // @ts-ignore
                                        error.statusMessage
                                    "
                                />
                                <!-- @vue-ignore -->
                                <pre
                                    v-else-if="
                                        // @ts-ignore
                                        error.message
                                    "
                                    v-text="
                                        // @ts-ignore
                                        error.message
                                    "
                                />
                                <pre v-else>Unable to get error message</pre>
                            </span>
                            <span v-else>
                                <pre>Unknown error</pre>
                            </span>
                        </div>

                        <div class="flex justify-center gap-4">
                            <UButton size="xl" :disabled="buttonDisabled" @click="handleError()">
                                {{ $t('common.home') }}
                            </UButton>

                            <UButton size="xl" color="green" :disabled="buttonDisabled" @click="handleError(route.fullPath)">
                                {{ $t('common.retry') }}
                            </UButton>

                            <!-- @vue-ignore -->
                            <UButton
                                v-if="error && (error.statusMessage || error.message)"
                                size="xl"
                                color="amber"
                                @click="copyError"
                            >
                                {{ $t ? $t('pages.error.copy_error') : 'Copy Error message' }}
                            </UButton>
                        </div>
                    </div>
                </main>
            </div>
        </div>
    </div>
</template>
