<script lang="ts">
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';

import { defineComponent } from 'vue';
import { mapState } from 'vuex';

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
    <div class="hero min-h-screen bg-base-200">
        <div class="hero-content text-center">
            <div class="max-w-md">
                <h1 class="text-5xl font-bold">Hello there</h1>
                <p class="py-6">Provident cupiditate voluptatem et in. Quaerat fugiat ut assumenda excepturi exercitationem
                    quasi. In deleniti eaque aut repudiandae et a id nisi.</p>
                <router-link v-if="!accessToken" to="/login" class="btn btn-primary">Login</router-link>
                <router-link v-if="accessToken" to="/overview" class="btn btn-primary">Overview</router-link>
            </div>
        </div>
    </div>
    <Footer />
</template>
