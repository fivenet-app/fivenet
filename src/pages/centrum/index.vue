<script lang="ts" setup>
import Livemap from '~/components/livemap/Livemap.vue';
import NotificationsList from '~/components/partials/notification/NotificationsList.vue';

useHead({
    title: 'common.dispatch_center',
});
definePageMeta({
    title: 'common.dispatch_center',
    requiresAuth: true,
    permission: 'CentrumService.Stream',
    showQuickButtons: false,
});
</script>

<template>
    <div class="flex-col h-full">
        <div class="relative w-full h-full z-0 flex">
            <!-- Left column -->
            <div class="flex-1">
                <!-- Main area -->
                <Livemap />
            </div>

            <div class="flex-1">
                <!-- Right column area -->
                <NotificationsList />
            </div>
        </div>
    </div>
</template>
