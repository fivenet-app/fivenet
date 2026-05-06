<script lang="ts" setup>
useHead({
    title: 'common.imprint',
});

definePageMeta({
    title: 'common.imprint',
    layout: 'landing',
    requiresAuth: false,
    redirectIfAuthed: false,
    showCookieOptions: true,
});

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const { website } = useAppConfig();

if (website.links?.imprint === undefined) {
    navigateTo('/');
} else {
    useTimeoutFn(() => navigateTo(website.links!.imprint!, { external: true }), 1750);
}
</script>

<template>
    <div class="flex h-dvh flex-col">
        <div class="hero absolute inset-0 z-[-1]" />

        <div class="flex w-full flex-1 items-center justify-center">
            <UCard class="w-full max-w-xl bg-white/75 backdrop-blur-sm dark:bg-white/5">
                <template #header>
                    <h3 class="text-2xl leading-6 font-semibold">
                        {{ $t('common.redirecting_to', [$t('common.imprint')]) }}
                    </h3>
                </template>

                <div class="space-y-4">
                    <UButton
                        class="mb-2"
                        trailing-icon="i-mdi-link-variant"
                        size="lg"
                        :to="nuiEnabled ? undefined : website.links!.imprint"
                        external
                        rel="noreferrer"
                        block
                        :label="$t('common.imprint')"
                        @click="nuiEnabled ? openURLInWindow(website.links!.imprint!) : undefined"
                    />
                </div>
            </UCard>
        </div>
    </div>
</template>
