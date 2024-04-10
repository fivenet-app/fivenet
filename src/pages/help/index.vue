<script setup lang="ts">
useHead({
    title: 'common.help',
});
definePageMeta({
    title: 'common.help',
    layout: 'help',
    requiresAuth: false,
    showCookieOptions: true,
});
</script>

<template>
    <div>
        <ULandingHero :title="$t('help.hero.title')" :description="$t('help.hero.description')">
            <template #title>
                <MDC :value="$t('help.hero.title')" />
            </template>
        </ULandingHero>

        <NuxtLink to="/help/citizens"> Citizens </NuxtLink>
        <NuxtLink to="/help/vehicles"> Vehicles </NuxtLink>
    </div>
</template>
