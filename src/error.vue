<script setup lang="ts">
import { useClipboard } from '@vueuse/core';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
import LoadingBar from '~/components/partials/LoadingBar.vue';

useHead({
    title: 'Error occured - FiveNet',
});

const { $loading } = useNuxtApp();
const clipboard = useClipboard();

const props = defineProps<{
    error: Error | Object;
}>();

const buttonDisabled = ref(true);

function handleError(): void {
    $loading.start();

    startButtonTimer();
    clearError({ redirect: '/' });
}

function copyError(): void {
    if (!props.error) {
        return;
    }

    clipboard.copy(JSON.stringify(props.error));
}

function startButtonTimer(): void {
    buttonDisabled.value = true;

    setTimeout(() => (buttonDisabled.value = false), 2000);
    setTimeout(() => $loading.errored(), 350);
}

$loading.start();
startButtonTimer();
</script>

<template>
    <div class="h-dscreen">
        <LoadingBar />
        <HeroFull>
            <ContentCenterWrapper class="max-w-3xl mx-auto text-center">
                <img class="h-auto mx-auto mb-2 w-36" src="/images/logo.png" alt="FiveNet Logo" />

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
                        <pre v-if="error.statusMessage">{{ error.statusMessage }}</pre>
                        <pre v-else-if="error.message">{{ error.message }}</pre>
                        <pre v-else>Unable to get error message</pre>
                    </span>
                    <span v-else>
                        <pre>Unknown error</pre>
                    </span>
                </div>

                <div class="flex justify-center">
                    <button
                        @click="handleError"
                        :disabled="buttonDisabled"
                        class="rounded-md w-60 px-3.5 py-2.5 text-sm font-semibold text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                        :class="[
                            buttonDisabled
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                        ]"
                    >
                        {{ $t('common.home') }}
                    </button>

                    <button
                        @click="copyError"
                        v-if="error && (error.statusMessage || error.message)"
                        class="rounded-md w-60 bg-base-600 sm:ml-4 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-base-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-500"
                    >
                        {{ $t('pages.error.copy_error') }}
                    </button>
                </div>
            </ContentCenterWrapper>
        </HeroFull>
    </div>
</template>
