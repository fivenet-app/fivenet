<script lang="ts">
import { defineComponent } from 'vue';
import { mapActions } from 'vuex';

import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';
import Login from '../components/Login.vue';

export default defineComponent({
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
});
</script>

<template>
  <Navbar />
  <Login />
  <Footer />
</template>
