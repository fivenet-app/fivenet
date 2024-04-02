<script lang="ts" setup>
import { type RoutesNamedLocations } from '@typed-router';
import type { Perms } from '~~/gen/ts/perms';

useHead({
    title: 'pages.jobs.title',
});
definePageMeta({
    title: 'pages.jobs.title',
    requiresAuth: true,
    permission: 'JobsService.ListColleagues',
    redirect: { name: 'jobs-overview' },
});

const { t } = useI18n();

const tabs: { label: string; to: RoutesNamedLocations; permission?: Perms; icon: string }[] = [
    {
        label: t('common.overview'),
        to: { name: 'jobs-overview' },
        permission: 'JobsService.ListColleagues' as Perms,
        icon: 'i-mdi-briefcase',
    },
    {
        label: t('pages.jobs.colleagues.title'),
        to: { name: 'jobs-colleagues' },
        permission: 'JobsService.ListColleagues' as Perms,
        icon: 'i-mdi-account-group',
    },
    {
        label: t('common.activity'),
        to: { name: 'jobs-activity' },
        permission: 'JobsService.ListColleagueActivity' as Perms,
        icon: 'i-mdi-bulletin-board',
    },
    {
        label: t('pages.jobs.timeclock.title'),
        to: { name: 'jobs-timeclock' },
        permission: 'JobsTimeclockService.ListTimeclock' as Perms,
        icon: 'i-mdi-timeline-clock',
    },
    {
        label: t('pages.qualifications.title'),
        to: { name: 'jobs-qualifications' },
        permission: 'QualificationsService.ListQualifications' as Perms,
        icon: 'i-mdi-school',
    },
    {
        label: t('pages.jobs.conduct.title'),
        to: { name: 'jobs-conduct' },
        permission: 'JobsConductService.ListConductEntries' as Perms,
        icon: 'i-mdi-list-status',
    },
].filter((t) => t.permission === undefined || can(t.permission));

function onChange(index: number) {
    const item = tabs[index];

    navigateTo(item.to);
}
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.jobs.title')"> </UDashboardNavbar>

            <UTabs :items="tabs" class="w-full" @change="onChange">
                <template #default="{ item, selected }">
                    <div class="flex items-center gap-2 relative truncate">
                        <UIcon :name="item.icon" class="w-4 h-4 flex-shrink-0" />

                        <span class="truncate">{{ item.label }}</span>

                        <span
                            v-if="selected"
                            class="absolute -right-4 w-2 h-2 rounded-full bg-primary-500 dark:bg-primary-400"
                        />
                    </div>
                </template>
            </UTabs>

            <main>
                <div class="mx-auto max-w-7xl py-4 sm:px-6 lg:px-8">
                    <NuxtLayout name="blank">
                        <NuxtPage
                            :transition="{
                                name: 'page',
                                mode: 'out-in',
                            }"
                        />
                    </NuxtLayout>
                </div>
            </main>
        </UDashboardPanel>
    </UDashboardPage>
</template>
