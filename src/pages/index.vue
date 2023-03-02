<script lang="ts">
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';

import { defineComponent } from 'vue';
import { mapState } from 'vuex';

import './index.css';
import HeroFull from '../components/HeroFull.vue';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        HeroFull
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
        "requiresAuth": false
    }
}
</route>

<template>
    <Navbar />
    <HeroFull>
        <h1 class="text-5xl font-bold text-white">Welcome to aRPaNet!</h1>
        <p class="py-6 text-white">aRPaNet is your access to FiveM roleplay. From searching the state's user
            database, filling documents for court and a livemap of your colleagues and dispatches.<br />It's all
            in here.</p>
        <router-link v-if="accessToken" to="/overview" class="btn btn-primary">Overview</router-link>
        <router-link v-else to="/login" class="btn btn-primary">Login</router-link>
    </HeroFull>
    <Footer />
</template>
