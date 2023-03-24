<script lang="ts" setup>
import { onBeforeMount } from 'vue';
import { useStore } from '../store/store';
import { useRouter } from 'vue-router/auto';
import Footer from '../components/partials/Footer.vue';
import HeroFull from '../components/partials/HeroFull.vue';
import './herofull-pattern.css';

const store = useStore();
const router = useRouter();

onBeforeMount(() => {
    store.dispatch('auth/doLogout').
        then(() => {
            setTimeout(() => {
                router.push('/');
            }, 1500);
        });
});
</script>

<route lang="json">
{
    "name": "Logout",
    "meta": {
        "requiresAuth": false
    }
}
</route>

<template>
    <HeroFull>
        <h1 class="text-4xl font-bold tracking-tight text-white sm:text-6xl">
            Signed out
        </h1>
        <p class="mt-6 text-lg leading-8 text-gray-300">
            We will redirect you to the home page in a few seconds.
        </p>
    </HeroFull>
    <Footer />
</template>
