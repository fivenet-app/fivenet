<script lang="ts" setup>
import DispatchesList from '~/components/centrum/DispatchesList.vue';
import SquadList from '~/components/centrum/SquadList.vue';
import Livemap from '~/components/livemap/Livemap.vue';

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
            <div class="flex-1 flex flex-col">
                <!-- Main area -->
                <div class="h-full">
                    <Livemap />
                </div>
            </div>

            <div class="flex-1 flex flex-col">
                <!-- Right column area -->
                <div class="basis-2/3">
                    <DispatchesList />
                </div>
                <div class="basis-1/3">
                    <SquadList />
                </div>
            </div>
        </div>
    </div>
</template>
