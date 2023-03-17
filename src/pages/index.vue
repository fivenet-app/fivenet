<script lang="ts">
import { defineComponent } from 'vue';
import { mapState } from 'vuex';
import Footer from '../components/partials/Footer.vue';
import HeroFull from '../components/partials/HeroFull.vue';

import './herofull-pattern.css';

export default defineComponent({
    components: {
        Footer,
        HeroFull,
    },
    computed: {
        ...mapState({
            accessToken: 'accessToken',
        }),
    },
});
</script>

<route lang="json">
{
    "name": "index",
    "meta": {
        "requiresAuth": false,
        "breadCrumbs": null
    }
}
</route>

<template>
    <HeroFull>
        <h1 class="text-4xl font-bold tracking-tight text-white sm:text-6xl">Welcome to aRPaNet!</h1>
        <p class="mt-6 text-lg leading-8 text-gray-300">
            aRPaNet is your access to FiveM roleplay. From searching the state's user
            database, filling documents for court and a livemap of your colleagues and dispatches.<br />It's all
            in here.
        </p>
        <div class="mt-10 flex items-center justify-center gap-x-6">
            <router-link v-if="accessToken" to="/overview"
                class="rounded-md bg-indigo-500 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-400">Overview</router-link>
            <router-link v-else to="/login"
                class="rounded-md bg-indigo-500 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-400">Login</router-link>
        </div>
    </HeroFull>
    <Footer />
</template>
