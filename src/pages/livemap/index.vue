<script lang="ts" setup>
import DisponentsInfo from '~/components/centrum/disponents/DisponentsInfo.vue';
import LivemapHolder from '~/components/livemap/LivemapHolder.vue';

useHead({
    title: 'common.livemap',
});
definePageMeta({
    title: 'common.livemap',
    requiresAuth: true,
    permission: 'LivemapperService.Stream',
});
</script>

<template>
    <UDashboardPage class="size-full">
        <UDashboardPanel grow class="size-full">
            <UDashboardNavbar :title="$t('common.livemap')">
                <template #right>
                    <DisponentsInfo :hide-join="true" />
                </template>
            </UDashboardNavbar>

            <LivemapHolder class="size-full" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
