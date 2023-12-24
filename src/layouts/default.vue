<script lang="ts" setup>
import CommandPalette from '~/components/partials/CommandPalette.vue';
import LoadingBar from '~/components/partials/LoadingBar.vue';
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';
import NotificatorProvider from '~/components/partials/notification/NotificatorProvider.vue';
import SidebarContainer from '~/components/partials/sidebar/SidebarContainer.vue';
import { useAuthStore } from '~/store/auth';

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

// Use client date to show any event overlays
const now = new Date();
const showSnowflakes = now.getMonth() + 1 === 12 && (now.getDate() > 20 || now.getDate() < 26);
</script>

<template>
    <NotificationProvider>
        <SidebarContainer>
            <!-- Events-->
            <LazyPartialsEventsSnowflakesContainer v-if="showSnowflakes" />

            <div class="h-full">
                <LoadingBar />
                <slot />
            </div>
        </SidebarContainer>
        <CommandPalette v-if="activeChar" />
        <NotificatorProvider />
    </NotificationProvider>
</template>
