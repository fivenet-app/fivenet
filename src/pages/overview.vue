<script lang="ts">
import { defineComponent } from 'vue';
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';
import Overview from '../components/Overview.vue';

export default defineComponent({
  components: {
    Navbar,
    Footer,
    Overview,
  },
});
</script>

<route lang="json">
  {
    "name": "overview",
    "meta": {
      "requiresAuth": true
    }
  }
  </route>

<template>
  <Navbar />
  <Overview />
  <Footer />
</template>
