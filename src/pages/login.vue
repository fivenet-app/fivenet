<script lang="ts">
import { defineComponent } from 'vue';

import Navbar from '../components/Navbar.vue';
import Footer from '../components/Footer.vue';
import Login from '../components/Login.vue';
import CharacterSelector from '../components/CharacterSelector.vue';
import { mapState } from 'vuex';

export default defineComponent({
  components: {
    Navbar,
    Footer,
    Login,
    CharacterSelector
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
  "name": "login",
  "meta": {
    "requiresAuth": false
  }
}
</route>

<template>
  <Navbar />
  <Login v-if="!accessToken" />
  <div v-if="accessToken" class="container mx-auto py-8">
    <CharacterSelector />
  </div>
  <Footer />
</template>
