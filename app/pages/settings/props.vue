<script lang="ts" setup>
import JobProps from '~/components/settings/JobProps.vue';

useHead({
    title: 'components.settings.job_props.job_properties',
});

definePageMeta({
    title: 'components.settings.job_props.job_properties',
    requiresAuth: true,
    permission: 'settings.SettingsService.GetJobProps',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <JobProps />
        </UDashboardPanel>
    </UDashboardPage>
</template>
