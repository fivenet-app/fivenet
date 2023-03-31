<script lang="ts" setup>
import { useAuthStore } from '../store/auth';
import HeroFull from '../components/partials/HeroFull.vue';
import { computed } from 'vue';

useHead({
    title: '404 Not Found',
});
definePageMeta({
    requiresAuth: false,
});

const store = useAuthStore();
const accessToken = computed(() => store.$state.accessToken);
</script>

<template>
    <HeroFull>
        <h1 class="text-5xl font-bold text-neutral">Page Not Found</h1>
        <h2 class="text-4xl text-neutral">4😵‍💫4 Error</h2>
        <p class="py-6 text-neutral">Seems like you got lost. Let's get you back to the start!</p>
        <NuxtLink v-if="accessToken" :to="{ name: 'overview' }" class="btn btn-primary">Overview</NuxtLink>
        <NuxtLink v-else :to="{ name: 'auth-login' }" class="btn btn-primary">Login</NuxtLink>
    </HeroFull>
</template>
