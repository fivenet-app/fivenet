<script lang="ts" setup>
import Login from '~/components/auth/Login.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';

useHead({
    title: 'Login',
});
definePageMeta({
    title: 'Login',
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
