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
    redirect: { name: 'jobs-overview' },
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
                to: { name: 'jobs-overview' },
            },
            {
                label: t('pages.jobs.colleagues.title'),
                icon: 'i-mdi-account-group',
                to: { name: 'jobs-colleagues' },
                permission: 'jobs.JobsService/ListColleagues' as Perms,
                active: route.name?.startsWith('jobs-colleagues'),
            },
            {
                label: t('common.activity'),
                icon: 'i-mdi-pulse',
                to: { name: 'jobs-activity' },
                permission: 'jobs.JobsService/ListColleagueActivity' as Perms,
            },
            {
                label: t('pages.jobs.timeclock.title'),
                icon: 'i-mdi-timeline-clock',
                to: { name: 'jobs-timeclock' },
                permission: 'jobs.TimeclockService/ListTimeclock' as Perms,
                active: route.name?.startsWith('jobs-timeclock'),
            },
            {
                label: t('pages.jobs.conduct.title'),
                icon: 'i-mdi-list-status',
                to: { name: 'jobs-conduct' },
                permission: 'jobs.ConductService/ListConductEntries' as Perms,
            },
        ].filter((t) => t.permission === undefined || can(t.permission).value) as {
            label: string;
            to: RoutesNamedLocations;
            active?: boolean;
            icon: string;
            permission?: Perms;
        }[],
);

inject('links', items);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0 overflow-y-hidden' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.jobs.title')">
                <template #right>
                    <PartialsBackButton fallback-to="/jobs/overview" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <UNavigationMenu orientation="horizontal" :items="items" class="-mx-1 flex-1" />
            </UDashboardToolbar>
        </template>

        <template #body>
            <NuxtPage />
        </template>
    </UDashboardPanel>
</template>
