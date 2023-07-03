<script lang="ts" setup>
import AttrView from '~/components/rector/attrs/AttrView.vue';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';

useHead({
    title: 'pages.rector.limiter.title',
});
definePageMeta({
    title: 'pages.rector.limiter.title',
    requiresAuth: true,
    permission: 'SuperUser',
    showQuickButtons: false,
    validate: async (route) => {
        route = route as TypedRouteFromName<'rector-limiter-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('rector-limiter-id');
</script>

<template>
    <div class="w-full">
        <AttrView :role-id="BigInt(route.params.id)" />
    </div>
</template>
