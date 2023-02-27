<script lang="ts">
import { defineComponent } from 'vue';

import Navbar from '../../components/Navbar.vue';
import Footer from '../../components/Footer.vue';

export default defineComponent({
  components: {
    Navbar,
    Footer,
  },
});
</script>

<route lang="json">
{
  "name": "documents",
  "meta": {
    "requiresAuth": true
  }
}
</route>

<template>
  <Navbar />
  <h1>DOCUMENTS</h1>
  <Footer />
</template>
