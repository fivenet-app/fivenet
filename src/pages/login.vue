<script lang="ts">
import { defineComponent } from 'vue';

import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';
import Login from '../components/Login.vue';

export default defineComponent({
  components: {
    Navbar,
    Footer,
    Login,
  },
});
</script>

<route lang="json">
{
  "name": "login",
  "meta": {
    "requiresAuth": true
  }
}
</route>

<template>
  <Navbar />
  <Login />
  <Footer />
</template>
