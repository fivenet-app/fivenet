<script setup lang="ts">
import LoadingSpinner from './components/LoadingSpinner.vue';
import { NotificationProvider } from './components/notification';
import Sidebar from './components/partials/Sidebar.vue';
</script>

<template>
    <NotificationProvider>
        <RouterView v-slot="{ Component, route }">
            <LoadingSpinner />
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
