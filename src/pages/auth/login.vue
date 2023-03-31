<script lang="ts" setup>
import Login from '../../components/auth/Login.vue';
import HeroFull from '../../components/partials/HeroFull.vue';

useHead({
    title: 'Login',
});
definePageMeta({
    requiresAuth: false,
});
</script>

<template>
    <HeroFull>
        <Login />
    </HeroFull>
</template>
