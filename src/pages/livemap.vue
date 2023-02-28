<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';
import Livemap from '../components/Livemap.vue';

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
  "name": "livemap",
  "meta": {
    "requiresAuth": true
  }
}
</route>

<template>
  <Navbar />
  <div class="container w-screen">
    <Livemap />
  </div>
  <Footer />
</template>
