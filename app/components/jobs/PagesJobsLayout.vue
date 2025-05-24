<script lang="ts" setup>
import type { RoutesNamedLocations } from '@typed-router';
import type { Perms } from '~~/gen/ts/perms';

const { t } = useI18n();

const { can } = useAuth();

const links = [
    {
        label: t('common.overview'),
        to: { name: 'jobs-overview' },
        icon: 'i-mdi-briefcase',
    },
    {
        label: t('pages.jobs.colleagues.title'),
        to: { name: 'jobs-colleagues' },
        icon: 'i-mdi-account-group',
        permission: 'jobs.JobsService.ListColleagues' as Perms,
    },
    {
        label: t('common.activity'),
        to: { name: 'jobs-activity' },
        icon: 'i-mdi-pulse',
        permission: 'jobs.JobsService.ListColleagueActivity' as Perms,
    },
    {
        label: t('pages.jobs.timeclock.title'),
        to: { name: 'jobs-timeclock' },
        icon: 'i-mdi-timeline-clock',
        permission: 'jobs.TimeclockService.ListTimeclock' as Perms,
    },
    {
        label: t('pages.jobs.conduct.title'),
        to: { name: 'jobs-conduct' },
        icon: 'i-mdi-list-status',
        permission: 'jobs.ConductService.ListConductEntries' as Perms,
    },
].filter((t) => t.permission === undefined || can(t.permission).value) as {
    label: string;
    to: RoutesNamedLocations;
    icon: string;
    permission?: Perms;
}[];
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.jobs.title')" />

            <UDashboardToolbar class="overflow-x-auto px-1.5 py-0">
                <UHorizontalNavigation :links="links" />
            </UDashboardToolbar>

            <UDashboardPanelContent class="p-0 sm:pb-0">
                <slot />
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
