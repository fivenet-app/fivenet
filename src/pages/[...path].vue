<script setup lang="ts">
import { useStore } from '../store/store';
import HeroFull from '../components/partials/HeroFull.vue';
import Footer from '../components/partials/Footer.vue';
import './herofull-pattern.css';
import { computed } from 'vue';

const store = useStore();
const accessToken = computed(() => store.state.auth?.accessToken);
</script>

<route lang="json">
{
    "meta": {
        "requiresAuth": false
    }
}
</route>

<template>
    <HeroFull>
        <h1 class="text-5xl font-bold text-white">Page Not Found</h1>
        <h2 class="text-4xl text-white">4😵‍💫4 Error</h2>
        <p class="py-6 text-white">Seems like you got lost. Let's get you back to the start!</p>
        <router-link v-if="accessToken" to="/overview" class="btn btn-primary">Overview</router-link>
        <router-link v-else to="/login" class="btn btn-primary">Login</router-link>
    </HeroFull>
    <Footer />
</template>
