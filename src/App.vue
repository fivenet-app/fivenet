<script lang="ts">
import Navbar from './components/Navbar.vue';
import Footer from './components/Footer.vue';
import Welcome from './components/Home.vue';

import { mapActions } from 'vuex';
import { defineComponent } from 'vue';
import { store } from './store';

export default defineComponent({
  components: {
    Navbar,
    Footer,
    Welcome,
  },
  beforeCreate() {
    store.commit('initialiseStore');
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
