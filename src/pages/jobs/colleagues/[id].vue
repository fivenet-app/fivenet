<script lang="ts" setup>
import { type RoutesNamedLocations, type TypedRouteFromName } from '@typed-router';
import ColleagueInfo from '~/components/jobs/colleagues/info/ColleagueInfo.vue';
import type { Perms } from '~~/gen/ts/perms';

useHead({
    title: 'pages.jobs.colleagues.single.title',
});
definePageMeta({
    title: 'pages.jobs.colleagues.single.title',
    requiresAuth: true,
    permission: 'JobsService.GetColleague',
    redirect: { name: 'jobs-colleagues-id-actvitiy' },
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-actvitiy'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const { t } = useI18n();

const route = useRoute('jobs-colleagues-id-actvitiy');

const links = [
    {
        label: t('common.activity'),
        to: { name: 'jobs-colleagues-id-actvitiy' },
        icon: 'i-mdi-bulletin-board',
        permission: 'JobsService.ListColleagueActivity' as Perms,
    },
    {
        label: t('common.timeclock'),
        to: { name: 'jobs-colleagues-id-timeclock' },
        icon: 'i-mdi-timeline-clock',
        permission: 'JobsTimeclockService.ListTimeclock' as Perms,
    },
    {
        label: t('pages.qualifications.title'),
        to: { name: 'jobs-colleagues-id-qualifications' },
        icon: 'i-mdi-school',
        permission: 'QualificationsService.ListQualifications' as Perms,
    },
    {
        label: t('pages.jobs.conduct.title'),
        to: { name: 'jobs-colleagues-id-conduct' },
        icon: 'i-mdi-list-status',
        permission: 'JobsConductService.ListConductEntries' as Perms,
    },
].filter((tab) => can(tab.permission)) as { label: string; to: RoutesNamedLocations; icon: string; permission: Perms }[];
</script>

<template>
    <PagesJobsLayout>
        <template #default>
            <UDashboardPanelContent>
                <ColleagueInfo :user-id="parseInt(route.params.id)" />

                <UDashboardToolbar class="overflow-x-auto px-1.5 py-0">
                    <UHorizontalNavigation :links="links" />
                </UDashboardToolbar>

                <NuxtPage />
            </UDashboardPanelContent>
        </template>
    </PagesJobsLayout>
</template>
