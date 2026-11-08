<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import View from '~/components/qualifications/View.vue';

useHead({
    title: 'pages.qualifications.id.title',
});

definePageMeta({
    title: 'pages.qualifications.id.title',
    requiresAuth: true,
    permission: 'qualifications.QualificationsService/ListQualifications',
    validate: async (route) => {
        route = route as TypedRouteFromName<'qualifications-id'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('qualifications-id');
const qualificationId = computed(() => (typeof route.params.id === 'string' ? parseInt(route.params.id) : 0));
</script>

<template>
    <View :qualification-id="qualificationId" />
</template>
