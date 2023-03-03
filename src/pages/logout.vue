<script lang="ts">
import { defineComponent } from 'vue';
import { mapActions } from 'vuex';
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';
import HeroFull from '../components/HeroFull.vue';

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
        <h1 class="text-5xl font-bold">Signed out</h1>
        <p class="py-6">We will redirect you to the home page in a few seconds.</p>
    </HeroFull>
    <Footer />
</template>
