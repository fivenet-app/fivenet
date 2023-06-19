<script lang="ts" setup>
import CommandPalette from '~/components/partials/CommandPalette.vue';
import LoadingBar from '~/components/partials/LoadingBar.vue';
import Notificator from '~/components/partials/Notificator.vue';
import Sidebar from '~/components/partials/Sidebar.vue';
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';
</script>

<template>
    <NotificationProvider>
        <Sidebar>
            <div class="h-full">
                <LoadingBar />
                <slot />
            </div>
        </Sidebar>
        <CommandPalette />
        <Notificator />
    </NotificationProvider>
</template>
