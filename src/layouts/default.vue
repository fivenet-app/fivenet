<script lang="ts" setup>
import { NotificationProvider } from '~/components/notification';
import Sidebar from '~/components/partials/Sidebar.vue';
</script>

<template>
    <NotificationProvider>
        <Sidebar>
            <div class="h-full">
                <slot />
            </div>
        </Sidebar>
    </NotificationProvider>
</template>
