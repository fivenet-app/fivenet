<script lang="ts" setup>
import QualificationsList from '~/components/jobs/qualifications/QualificationsList.vue';
import QualificationsResultsList from '~/components/jobs/qualifications/QualificationsResultsList.vue';

useHead({
    title: 'pages.jobs.qualifications.title',
});
definePageMeta({
    title: 'pages.jobs.qualifications.title',
    requiresAuth: true,
    permission: 'QualificationsService.ListQualifications',
});
</script>

<template>
    <div>
        <div class="py-2 pb-14">
            <div class="px-1 sm:px-2 lg:px-4">
                <div class="grid-col-2 grid gap-2">
                    <div class="sm:flex-auto">
                        <QualificationsList />

                        <QualificationsResultsList class="mt-4" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
