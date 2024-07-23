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

const { links } = useAppConfig();

if (links.imprint === undefined) {
    navigateTo({ name: 'index' });
} else {
    useTimeoutFn(() => navigateTo(links.imprint!, { external: true }), 1750);
}
</script>

<template>
    <div class="flex flex-col">
        <div class="hero absolute inset-0 z-[-1] [mask-image:radial-gradient(100%_100%_at_top,white,transparent)]" />

        <div class="w-full flex-1">
            <ULandingHero
                :title="$t('common.redirecting_to', [$t('common.imprint')])"
                :links="[{ label: $t('common.imprint'), icon: 'i-mdi-link-variant', size: 'lg', to: links.imprint }]"
            />
        </div>
    </div>
</template>
