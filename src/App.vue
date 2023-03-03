<script lang="ts">
import Navbar from './components/Navbar.vue';
import Footer from './components/Footer.vue';

import { defineComponent } from 'vue';
import store from './store';
import { NotificationProvider } from './components/Notification';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        NotificationProvider,
    },
    beforeCreate() {
        store.commit('initialiseStore');
    },
});
</script>

<template>
    <NotificationProvider>
        <RouterView v-slot="{ Component, route }">
            <transition name="fade" mode="out-in">
                <div :key="route.path">
                    <component :is="Component" />
                </div>
            </transition>
        </RouterView>
    </NotificationProvider>
</template>
