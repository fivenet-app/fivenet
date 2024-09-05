<script lang="ts" setup>
import { type RoutesNamedLocations, type TypedRouteFromName } from '@typed-router';
import ColleagueInfo from '~/components/jobs/colleagues/info/ColleagueInfo.vue';
import PagesJobsLayout from '~/components/jobs/PagesJobsLayout.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import type { Perms } from '~~/gen/ts/perms';
import type { GetColleagueResponse } from '~~/gen/ts/services/jobs/jobs';

useHead({
    title: 'pages.jobs.colleagues.single.title',
});
definePageMeta({
    title: 'pages.jobs.colleagues.single.title',
    requiresAuth: true,
    permission: 'JobsService.GetColleague',
    redirect: { name: 'jobs-colleagues-id-info' },
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-info'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return idParamRegex.test(route.params.id as string);
    },
});

const { t } = useI18n();

const route = useRoute('jobs-colleagues-id-info');

const {
    data: colleague,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`jobs-colleague-${route.params.id as string}`, () => getColleague(parseInt(route.params.id as string)));

async function getColleague(userId: number): Promise<GetColleagueResponse> {
    try {
        const call = getGRPCJobsClient().getColleague({
            userId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const links = [
    {
        label: t('common.info'),
        to: { name: 'jobs-colleagues-id-info' },
        icon: 'i-mdi-information-slab-circle',
        permission: 'JobsService.GetColleague' as Perms,
    },
    {
        label: t('common.activity'),
        to: { name: 'jobs-colleagues-id-activity' },
        icon: 'i-mdi-pulse',
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
].filter((tab) => can(tab.permission).value) as { label: string; to: RoutesNamedLocations; icon: string; permission: Perms }[];
</script>

<template>
    <PagesJobsLayout>
        <template #default>
            <UDashboardPanelContent>
                <DataPendingBlock v-if="!colleague && loading" :message="$t('common.loading', [$t('common.colleague', 1)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.colleague', 1)])"
                    :message="$t(error.message)"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!colleague || !colleague.colleague" />

                <template v-else>
                    <ColleagueInfo :colleague="colleague.colleague" />

                    <UDashboardToolbar class="overflow-x-auto px-1.5 py-0">
                        <UHorizontalNavigation :links="links" />
                    </UDashboardToolbar>

                    <NuxtPage :colleague="colleague.colleague" @refresh="refresh()" />
                </template>
            </UDashboardPanelContent>
        </template>
    </PagesJobsLayout>
</template>
