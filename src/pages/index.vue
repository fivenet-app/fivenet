<script lang="ts">
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';

import { defineComponent } from 'vue';
import { mapState } from 'vuex';

import './index.css';

export default defineComponent({
    components: {
        Navbar,
        Footer,
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
    <div class="hero min-h-screen bg-base-200 index">
        <div class="hero-overlay bg-opacity-60"></div>
        <div class="hero-content text-center w-">
            <div class="max-w-md">
                <h1 class="mb-5 text-5xl font-bold text-white">Welcome to aRPaNet!</h1>
                <p class="py-6 text-white">aRPaNet is your access to FiveM roleplay. From searching the state's user database, filling documents for court and a livemap of your colleagues and dispatches.<br />It's all in here.</p>
                <router-link v-if="!accessToken" to="/login" class="btn btn-primary">Login</router-link>
                <router-link v-if="accessToken" to="/overview" class="btn btn-primary">Overview</router-link>
            </div>
        </div>
    </div>
    <Footer />
</template>
