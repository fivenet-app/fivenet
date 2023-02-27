<script lang="ts">
import Navbar from './components/Navbar.vue';
import Footer from './components/Footer.vue';
import Welcome from './components/Home.vue';

import { mapActions } from 'vuex';
import { defineComponent } from 'vue';

export default defineComponent({
  components: {
    Navbar,
    Footer,
    Welcome,
  },
  methods: {
    ...mapActions([
      'fetchAccessToken'
    ]),
  },
  created() {
    this.fetchAccessToken();
  },
});
</script>

<template>
  <RouterView v-slot="{ Component, route }">
    <transition name="fade" mode="out-in">
      <div :key="route.path">
        <component :is="Component" />
      </div>
    </transition>
  </RouterView>
</template>
