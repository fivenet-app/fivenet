<script lang="ts" setup>
import DispatchFeed from '~/components/centrum/DispatchFeed.vue';
import DispatchesList from '~/components/centrum/DispatchesList.vue';
import Livemap from '~/components/centrum/Livemap.vue';
import UnitsList from '~/components/centrum/UnitsList.vue';

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
            <div class="flex flex-col basis-1/3 divide-x">
                <div class="h-full">
                    <Livemap />
                </div>
            </div>

            <!-- Right column -->
            <div class="flex flex-col basis-2/3 divide-y">
                <div class="basis-3/5">
                    <DispatchesList />
                </div>
                <div class="basis-1/5">
                    <DispatchFeed />
                </div>
                <div class="basis-1/5">
                    <UnitsList />
                </div>
            </div>
        </div>
    </div>
</template>
