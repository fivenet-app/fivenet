<script lang="ts">
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';
import Home from '../components/Home.vue';

import { defineComponent } from 'vue';

export default defineComponent({
  components: {
    Navbar,
    Footer,
    Home,
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
  <Home />
  <Footer />
</template>
