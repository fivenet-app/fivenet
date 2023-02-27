<script lang="ts">
import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';
import Home from '../components/Home.vue';

import { mapActions } from 'vuex';
import { defineComponent } from 'vue';

export default defineComponent({
  components: {
    Navbar,
    Footer,
    Home,
  },
  methods: {
    ...mapActions([
      'fetchAccessToken'
    ]),
  },
  created() {
    this.fetchAccessToken();
  }
});
</script>

<template>
  <Navbar />
  <Home />
  <Footer />
</template>
