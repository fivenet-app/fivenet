<script setup lang="ts">
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
import LoadingBar from '~/components/partials/LoadingBar.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

useHead({
    title: 'Error occured - FiveNet',
});

const { $loading } = useNuxtApp();

const props = defineProps<{
    error: Error | Object;
}>();

const buttonDisabled = ref(true);

function handleError(url?: string): void {
    if ($loading !== undefined) {
        $loading.start();
    }
    startButtonTimer();

    if (url === undefined) {
        url = '/';
    }

    reloadNuxtApp({
        path: url,
        persistState: false,
        ttl: 2000,
    });
    clearError();
}

function copyError(): void {
    if (!props.error) {
        return;
    }

    copyToClipboardWrapper(jsonStringify(props.error));
}

function startButtonTimer(): void {
    buttonDisabled.value = true;

    setTimeout(() => (buttonDisabled.value = false), 2000);
    setTimeout(() => {
        if ($loading !== undefined) {
            $loading.errored();
        }
    }, 400);
}

onBeforeMount(() => {
    if ($loading !== undefined) {
        $loading.start();
    }
    startButtonTimer();
});
</script>

<template>
    <div class="h-dscreen">
        <LoadingBar />
        <HeroFull>
            <ContentCenterWrapper class="max-w-3xl mx-auto text-center">
                <FiveNetLogo class="h-auto mx-auto mb-2 w-36" />

                <h1 class="text-5xl font-bold text-neutral">
                    {{ $t('pages.error.title') }}
                </h1>
                <h2 class="text-xl text-neutral">
                    {{ $t('pages.error.subtitle') }}
                </h2>
                <div class="py-2 text-neutral mb-4">
                    <p class="py-2 font-semibold">
                        {{ $t('pages.error.error_message') }}
                    </p>
                    <span v-if="error">
                        <!-- @vue-expect-error -->
                        <pre
                            v-if="error.statusMessage"
                            v-text="
                                // @ts-ignore
                                error.statusMessage
                            "
                        />
                        <!-- @vue-expect-error -->
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

                <div class="flex justify-center">
                    <button
                        :disabled="buttonDisabled"
                        class="rounded-md w-60 px-3.5 py-2.5 text-sm font-semibold text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                        :class="[
                            buttonDisabled
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                        ]"
                        @click="handleError()"
                    >
                        {{ $t('common.home') }}
                    </button>

                    <button
                        :disabled="buttonDisabled"
                        class="rounded-md w-60 px-3.5 py-2.5 sm:ml-4 text-sm font-semibold text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                        :class="[
                            buttonDisabled
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-secondary-500 hover:bg-secondary-400 focus-visible:outline-secondary-500',
                        ]"
                        @click="handleError(useRoute().fullPath)"
                    >
                        {{ $t('common.retry') }}
                    </button>

                    <!-- @vue-expect-error -->
                    <button
                        v-if="error && (error.statusMessage || error.message)"
                        class="rounded-md w-60 bg-base-600 sm:ml-4 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-base-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-500"
                        @click="copyError"
                    >
                        {{ $t('pages.error.copy_error') }}
                    </button>
                </div>
            </ContentCenterWrapper>
        </HeroFull>
    </div>
</template>