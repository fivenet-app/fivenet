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
    redirect: { name: 'jobs-colleagues-id-activity' },
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const { t } = useI18n();

const route = useRoute('jobs-colleagues-id');

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
    <div>
        <ColleagueInfo :user-id="parseInt(route.params.id)" />

        <UHorizontalNavigation :links="links" class="border-b border-gray-200 dark:border-gray-800" />

        <NuxtLayout name="blank">
            <NuxtPage
                :transition="{
                    name: 'page',
                    mode: 'out-in',
                }"
            />
        </NuxtLayout>
    </div>
</template>
