<script lang="ts" setup>
import List from '~/components/qualifications/List.vue';
import ListPage from '~/components/qualifications/ListPage.vue';

useHead({
    title: 'components.qualifications.all_qualifications',
});

definePageMeta({
    title: 'components.qualifications.all_qualifications',
    requiresAuth: true,
    permission: 'qualifications.QualificationsService/ListQualifications',
});
</script>

<template>
    <ListPage>
        <List />
    </ListPage>
</template>
