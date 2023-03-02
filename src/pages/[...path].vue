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
    "meta": {
        "requiresAuth": false
    }
}
</route>

<template>
    <Navbar />
    <HeroFull>
        <h1 class="text-5xl font-bold text-white">Page Not Found</h1>
        <h2 class="text-4xl text-white">4😵‍💫4 Error</h2>
        <p class="py-6 text-white">Seems like you got lost. Let's get you back to the start!</p>
        <router-link v-if="accessToken" to="/overview" class="btn btn-primary">Overview</router-link>
        <router-link v-else to="/login" class="btn btn-primary">Login</router-link>
    </HeroFull>
    <Footer />
</template>
