<script lang="ts" setup>
import { useStore } from '../store/store';
import HeroFull from '../components/partials/HeroFull.vue';
import { computed } from 'vue';

const store = useStore();
const accessToken = computed(() => store.state.auth?.accessToken);
</script>

<route lang="json">
{
    "name": "404 Not Found",
    "meta": {
        "requiresAuth": false
    }
}
</route>

<template>
    <HeroFull>
        <h1 class="text-5xl font-bold text-neutral">Page Not Found</h1>
        <h2 class="text-4xl text-neutral">4😵‍💫4 Error</h2>
        <p class="py-6 text-neutral">Seems like you got lost. Let's get you back to the start!</p>
        <router-link v-if="accessToken" :to="{ name: 'Overview' }" class="btn btn-primary">Overview</router-link>
        <router-link v-else :to="{ name: 'Login' }" class="btn btn-primary">Login</router-link>
    </HeroFull>
</template>
