<script lang="ts">
import Navbar from './components/Navbar.vue';
import Footer from './components/Footer.vue';
import Login from './components/Login.vue';

import { mapActions } from 'vuex';

export default {
  name: 'app',
  components: {
    Navbar,
    Footer,
    Login,
  },
  methods: {
    ...mapActions([
      'fetchAccessToken'
    ]),
  },
  created() {
    this.fetchAccessToken();
  }
};
</script>

<template>
  <Navbar />
  <Login />
  <Footer />
</template>
