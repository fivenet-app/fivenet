<script lang="ts">
import Sidebar from './components/partials/Sidebar.vue';

import { defineComponent } from 'vue';
import store from './store';
import { NotificationProvider } from './components/notification';

export default defineComponent({
    components: {
        Sidebar,
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
                    <Sidebar :child="Component" />
                </div>
            </transition>
        </RouterView>
    </NotificationProvider>
</template>
