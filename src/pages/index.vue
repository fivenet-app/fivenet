<script lang="ts" setup>
import { computed } from 'vue';
import { useStore } from '../store/store';
import Footer from '../components/partials/Footer.vue';
import HeroFull from '../components/partials/HeroFull.vue';
import './herofull-pattern.css';

const store = useStore();

const accessToken = computed(() => store.state.auth?.accessToken);
</script>

<route lang="json">
{
    "name": "Home",
    "meta": {
        "requiresAuth": false
    }
}
</route>

<template>
    <HeroFull>
        <h1 class="text-4xl font-bold tracking-tight text-white sm:text-6xl">Welcome to aRPaNet!</h1>
        <p class="mt-6 text-lg leading-8 text-gray-300">
            From searching the state's citizen and vehicles database, filling documents for investigations, court and a
            livemap of your colleagues and dispatches.
            <br />All that and more is (mostly) ready in this net, the aRPaNet.
        </p>
        <div class="mt-10 flex items-center justify-center gap-x-6">
            <router-link v-if="accessToken" :to="{ name: 'Overview' }"
                class="rounded-md bg-indigo-500 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-400">Overview</router-link>
            <router-link v-else :to="{ name: 'Login' }"
                class="rounded-md bg-indigo-500 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-400">Login</router-link>
        </div>
    </HeroFull>
    <Footer />
</template>
