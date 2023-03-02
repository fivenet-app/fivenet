<script lang="ts">
import Navbar from '../../components/Navbar.vue';
import Footer from '../../components/Footer.vue';

import { defineComponent } from 'vue';

import './index.css';

export default defineComponent({
    components: {
        Navbar,
        Footer,
    },
});
</script>

<route lang="json">
{
    "name": "dispatches",
    "meta": {
        "requiresAuth": true
    }
}
</route>

<template>
    <Navbar />
    <h1>DISPATCHES SYSTEM</h1>
    <Footer />
</template>
