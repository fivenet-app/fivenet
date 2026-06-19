<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import type { RoutesNamedLocations } from '@typed-router';
import type { Perms } from '~~/gen/ts/perms';

useHead({
    title: 'pages.jobs.title',
});

definePageMeta({
    title: 'pages.jobs.title',
    requiresAuth: true,
    redirect: '/jobs/overview',
});

const { t } = useI18n();

const { can } = useAuth();

const route = useRoute();

const items = computed<NavigationMenuItem[]>(
    () =>
        [
            {
                label: t('common.overview'),
                icon: 'i-mdi-briefcase-outline',
                to: '/jobs/overview',
            },
            {
                label: t('pages.jobs.colleagues.title'),
                icon: 'i-mdi-account-group',
                to: '/jobs/colleagues',
                permission: 'jobs.ColleaguesService/ListColleagues' as Perms,
                active: route.name?.startsWith('jobs-colleagues'),
                children: [
                    {
                        label: t('pages.jobs.colleagues.stats.title'),
                        icon: 'i-mdi-chart-timeline-variant-shimmer',
                        to: '/jobs/colleagues/stats',
                        permission: 'jobs.StatsService/GetStats' as Perms,
                    },
                    {
                        label: t('pages.jobs.colleagues.labels.title'),
                        icon: 'i-mdi-label',
                        to: '/jobs/colleagues/labels',
                        permission: 'jobs.ColleaguesService/ManageLabels' as Perms,
                    },
                ].flatMap((item) => (item.permission === undefined || can(item.permission).value ? [item] : [])),
            },
            {
                label: t('common.activity'),
                icon: 'i-mdi-pulse',
                to: '/jobs/activity',
                permission: 'jobs.ColleaguesService/ListColleagueActivity' as Perms,
            },
            {
                label: t('pages.jobs.timeclock.title'),
                icon: 'i-mdi-timeline-clock',
                to: '/jobs/timeclock',
                permission: 'jobs.TimeclockService/ListTimeclock' as Perms,
                active: route.name?.startsWith('jobs-timeclock'),
                children: [
                    {
                        label: t('common.inactive_colleagues'),
                        icon: 'i-mdi-account-remove',
                        to: '/jobs/timeclock/inactive',
                        permission: 'jobs.TimeclockService/ListInactiveEmployees' as Perms,
                    },
                ].flatMap((item) => (item.permission === undefined || can(item.permission).value ? [item] : [])),
            },
            {
                label: t('pages.jobs.conduct.title'),
                icon: 'i-mdi-list-status',
                to: '/jobs/conduct',
                permission: 'jobs.ConductService/ListConductEntries' as Perms,
            },
            {
                label: t('common.group', 2),
                icon: 'i-mdi-users-group-outline',
                to: '/jobs/groups',
                permission: 'TODOService/TODOMethod' as Perms,
            },
        ].filter((t) => t.permission === undefined || can(t.permission).value) as (NavigationMenuItem & {
            permission?: Perms;
            to: RoutesNamedLocations;
        })[],
);

inject('links', items);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0 overflow-y-hidden' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.jobs.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/jobs/overview" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar :ui="{ root: 'overflow-x-visible' }">
                <UNavigationMenu class="-mx-1 flex-1" orientation="horizontal" :items="items" />
            </UDashboardToolbar>
        </template>

        <template #body>
            <NuxtPage />
        </template>
    </UDashboardPanel>
</template>
