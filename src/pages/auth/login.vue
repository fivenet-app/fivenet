<script lang="ts" setup>
import Login from '~/components/auth/Login.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';

useHead({
    title: 'pages.auth.login.title',
});
definePageMeta({
    title: 'pages.auth.login.title',
    requiresAuth: false,
});
</script>

<template>
    <HeroFull>
        <ContentCenterWrapper>
            <Login />
        </ContentCenterWrapper>
    </HeroFull>
</template>
