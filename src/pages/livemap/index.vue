<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../../components/partials/Navbar.vue';
import Footer from '../../components/partials/Footer.vue';
import Livemap from '../../components/Livemap.vue';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        Livemap,
    },
});
</script>

<route lang="json">
{
    "name": "Livemap",
    "meta": {
        "requiresAuth": true,
        "permission": "livemap-view",
        "breadCrumbs": null
    }
}
</route>

<template>
    <Navbar />
    <div class="w-full">
        <Livemap />
    </div>
    <Footer />
</template>
