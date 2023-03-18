<script lang="ts">
import { defineComponent } from 'vue';
import Footer from '../../components/partials/Footer.vue';
import Livemap from '../../components/Livemap.vue';

export default defineComponent({
    components: {
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
        "permission": "LivemapperService.Stream",
        "breadCrumbs": null
    }
}
</route>

<template>
    <div class="w-full">
        <Livemap />
    </div>
    <Footer />
</template>
