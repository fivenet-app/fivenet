<script lang="ts" setup>
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';
import Sidebar from '~/components/partials/Sidebar.vue';
import LoadingBar from '~/components/partials/LoadingBar.vue';
import Notificator from '~/components/partials/Notificator.vue';
</script>

<template>
    <NotificationProvider>
        <Sidebar>
            <div class="h-full">
                <LoadingBar />
                <slot />
            </div>
        </Sidebar>
        <Notificator />
    </NotificationProvider>
</template>
