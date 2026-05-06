<script lang="ts" setup>
useHead({
    title: 'common.privacy_policy',
});

definePageMeta({
    title: 'common.privacy_policy',
    layout: 'landing',
    requiresAuth: false,
    redirectIfAuthed: false,
    showCookieOptions: true,
});

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const { website } = useAppConfig();

if (website.links?.privacyPolicy === undefined) {
    navigateTo('/');
} else {
    useTimeoutFn(() => navigateTo(website.links!.privacyPolicy!, { external: true }), 1750);
}
</script>

<template>
    <div class="flex h-dvh flex-col">
        <div class="hero absolute inset-0 z-[-1]" />

        <div class="flex w-full flex-1 items-center justify-center">
            <UCard class="w-full max-w-xl bg-white/75 backdrop-blur-sm dark:bg-white/5">
                <template #header>
                    <h3 class="text-2xl leading-6 font-semibold">
                        {{ $t('common.redirecting_to', [$t('common.privacy_policy')]) }}
                    </h3>
                </template>

                <div class="space-y-4">
                    <UButton
                        class="mb-2"
                        trailing-icon="i-mdi-link-variant"
                        size="lg"
                        :to="nuiEnabled ? undefined : website.links!.privacyPolicy"
                        external
                        rel="noreferrer"
                        block
                        :label="$t('common.privacy_policy')"
                        @click="nuiEnabled ? openURLInWindow(website.links!.privacyPolicy!) : undefined"
                    />
                </div>
            </UCard>
        </div>
    </div>
</template>
