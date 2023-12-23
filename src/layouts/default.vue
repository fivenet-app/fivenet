<script lang="ts" setup>
import CommandPalette from '~/components/partials/CommandPalette.vue';
import LoadingBar from '~/components/partials/LoadingBar.vue';
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';
import NotificatorProvider from '~/components/partials/notification/NotificatorProvider.vue';
import SidebarContainer from '~/components/partials/sidebar/SidebarContainer.vue';
import { useAuthStore } from '~/store/auth';

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);
</script>

<template>
    <NotificationProvider>
        <SidebarContainer>
            <div class="h-full">
                <LoadingBar />
                <slot />
            </div>
        </SidebarContainer>
        <CommandPalette v-if="activeChar" />
        <NotificatorProvider />
    </NotificationProvider>
</template>
