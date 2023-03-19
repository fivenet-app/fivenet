<script setup lang="ts">
import { useStore } from 'vuex';
import { NotificationProvider } from './components/notification';
import Sidebar from './components/partials/Sidebar.vue';

const store = useStore();
store.commit('initialiseStore');
</script>

<template>
    <NotificationProvider>
        <RouterView v-slot="{ Component, route }">
            <transition name="fade" mode="out-in">
                <div :key="route.path">
                    <Sidebar>
                        <component :is="Component"></component>
                    </Sidebar>
                </div>
            </transition>
        </RouterView>
    </NotificationProvider>
</template>
