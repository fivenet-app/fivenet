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

const { links } = useAppConfig();

if (links.privacyPolicy === undefined) {
    navigateTo({ name: 'index' });
} else {
    useTimeoutFn(() => navigateTo(links.privacyPolicy!, { external: true }), 1750);
}
</script>

<template>
    <div class="hero flex flex-col">
        <div class="w-full flex-1 bg-black/50">
            <ULandingHero
                :title="$t('common.redirecting_to', [$t('common.privacy_policy')])"
                :links="[{ label: $t('common.privacy_policy'), icon: 'i-mdi-link-variant', size: 'lg', to: links.imprint }]"
            />
        </div>
    </div>
</template>
