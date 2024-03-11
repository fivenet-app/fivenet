<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

useHead({
    title: 'pages.documents.edit.title',
});
definePageMeta({
    title: 'pages.documents.edit.title',
    requiresAuth: true,
    permission: 'JobsQualificationsService.UpdateQualification',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-qualifications-id-edit'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('jobs-qualifications-id-edit');
</script>

<template>
    <ContentWrapper>
        <!-- TODO -->
        EDIT QUALIFICATION
    </ContentWrapper>
</template>
