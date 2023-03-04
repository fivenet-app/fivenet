<script lang="ts">
import { defineComponent } from 'vue';
import { mapActions } from 'vuex';
import Navbar from '../components/partials/Navbar.vue';
import Footer from '../components/partials/Footer.vue';
import HeroFull from '../components/partials/HeroFull.vue';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        HeroFull
    },
    methods: {
        ...mapActions([
            'doLogout',
        ]),
    },
    beforeMount: function () {
        this.doLogout().then(() => {
            setTimeout(() => {
                this.$router.push('/');
            }, 2000);
        });
    },
});
</script>

<route lang="json">
{
    "name": "logout",
    "meta": {
        "requiresAuth": false
    }
}
</route>

<template>
    <Navbar />
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
