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
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.qualifications.title')">
                <template #right>
                    <UButton
                        v-if="can('QualificationsService.CreateQualification')"
                        :to="{ name: 'jobs-qualifications-create' }"
                        trailing-icon="i-mdi-plus"
                        color="gray"
                    >
                        {{ $t('components.qualifications.create_new_qualification') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <div class="flex flex-col gap-2">
                <QualificationsResultsList />

                <QualificationsRequestsList />

                <QualificationsList />
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>
