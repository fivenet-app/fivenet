<script setup lang="ts">
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import ContentHeroFull from '~/components/partials/ContentHeroFull.vue';
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
        <ContentHeroFull>
            <ContentCenterWrapper class="mx-auto max-w-3xl text-center">
                <FiveNetLogo class="mx-auto mb-2 h-auto w-36" />

                <h1 class="text-5xl font-bold text-neutral">
                    {{ $t ? $t('pages.error.title') : 'Error occured' }}
                </h1>
                <h2 class="text-xl text-neutral">
                    {{ $t ? $t('pages.error.subtitle') : 'A fatal error occured, please try again in a few seconds.' }}
                </h2>
                <div class="mb-4 py-2 text-neutral">
                    <p class="py-2 font-semibold">
                        {{ $t ? $t('pages.error.error_message') : 'Error message:' }}
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
                        class="w-60 rounded-md px-3.5 py-2.5 text-sm font-semibold text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
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
                        class="w-60 rounded-md px-3.5 py-2.5 text-sm font-semibold text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 sm:ml-4"
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
                        class="w-60 rounded-md bg-base-600 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-base-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-500 sm:ml-4"
                        @click="copyError"
                    >
                        {{ $t ? $t('pages.error.copy_error') : 'Copy Error message' }}
                    </button>
                </div>
            </ContentCenterWrapper>
        </ContentHeroFull>
    </div>
</template>
