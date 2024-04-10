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
        <UContainer>
            <div class="px-1 sm:px-2">
                <div class="flex flex-col gap-2">
                    <QualificationsResultsList />

                    <QualificationsRequestsList />

                    <QualificationsList />
                </div>
            </div>
        </UContainer>
    </div>
</template>
