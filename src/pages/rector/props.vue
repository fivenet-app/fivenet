<script lang="ts" setup>
import JobProps from '~/components/rector/JobProps.vue';

useHead({
    title: 'components.rector.job_props.job_properties',
});
definePageMeta({
    title: 'components.rector.job_props.job_properties',
    requiresAuth: true,
    permission: 'RectorService.GetJobProps',
    showQuickButtons: false,
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('components.rector.job_props.job_properties')"> </UDashboardNavbar>

            <JobProps />
        </UDashboardPanel>
    </UDashboardPage>
</template>
