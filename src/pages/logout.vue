<script lang="ts">
import { defineComponent } from 'vue';
import { mapActions } from 'vuex';
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';

export default defineComponent({
    components: {
        Navbar,
        Footer,
    },
    methods: {
        ...mapActions([
            'doLogout',
        ]),
    },
    mounted: function () {
        this.doLogout().then(() => {
            setTimeout(() => {
                this.$router.push('/');
            }, 2150);
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
    <div class="hero min-h-screen bg-base-200">
        <div class="hero-content text-center">
            <div class="max-w-md">
                <h1 class="text-5xl font-bold">Signed out</h1>
                <p class="py-6">We will redirect you to the home page in a few seconds.</p>
            </div>
        </div>
    </div>
    <Footer />
</template>
