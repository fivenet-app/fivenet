<script lang="ts" setup>
import QualificationsList from '~/components/jobs/qualifications/QualificationsList.vue';
import QualificationsRequestsList from '~/components/jobs/qualifications/QualificationsRequestsList.vue';
import QualificationsResultsList from '~/components/jobs/qualifications/QualificationsResultsList.vue';

useHead({
    title: 'pages.qualifications.title',
});
definePageMeta({
    title: 'pages.qualifications.title',
    requiresAuth: true,
    permission: 'QualificationsService.ListQualifications',
});
</script>

<template>
    <div>
        <div class="py-2 pb-14">
            <div class="px-1 sm:px-2 lg:px-4">
                <div class="grid gap-2 md:grid-cols-2">
                    <QualificationsList />

                    <QualificationsResultsList />

                    <QualificationsRequestsList />
                </div>
            </div>
        </div>
    </div>
</template>
