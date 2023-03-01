<script lang="ts">
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';

import { defineComponent } from 'vue';
import { mapState } from 'vuex';

export default defineComponent({
  components: {
    Navbar,
    Footer,
  },
  computed: {
    ...mapState({
      accessToken: 'accessToken',
    }),
  },
});
</script>

<route lang="json">
{
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
        <h1 class="text-5xl font-bold">Page Not Found</h1>
        <h2 class="text-4xl">4😵‍💫4 Error</h2>
        <p class="py-6">Seems like you got lost. Let's get you back to the start!</p>
        <router-link v-if="!accessToken" to="/login" class="btn btn-primary">Login</router-link>
        <router-link v-if="accessToken" to="/overview" class="btn btn-primary">Overview</router-link>
      </div>
    </div>
  </div>
  <Footer />
</template>
