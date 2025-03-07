<script lang="ts" setup>
import JobProps from '~/components/rector/JobProps.vue';

useHead({
    title: 'components.rector.job_props.job_properties',
});

definePageMeta({
    title: 'components.rector.job_props.job_properties',
    requiresAuth: true,
    permission: 'RectorService.GetJobProps',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <JobProps />
        </UDashboardPanel>
    </UDashboardPage>
</template>
