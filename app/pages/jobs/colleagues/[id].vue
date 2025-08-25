<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import ColleagueInfo from '~/components/jobs/colleagues/info/ColleagueInfo.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getJobsJobsClient } from '~~/gen/ts/clients';
import type { Perms } from '~~/gen/ts/perms';
import { ObjectType } from '~~/gen/ts/resources/notifications/client_view';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { GetColleagueResponse } from '~~/gen/ts/services/jobs/jobs';

useHead({
    title: 'pages.jobs.colleagues.single.title',
});

definePageMeta({
    title: 'pages.jobs.colleagues.single.title',
    requiresAuth: true,
    permission: 'jobs.JobsService/GetColleague',
    redirect: { name: 'jobs-colleagues-id-info' },
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-info'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const { t } = useI18n();

const { attr, can } = useAuth();

const notifications = useNotificationsStore();

const route = useRoute('jobs-colleagues-id-info');

const jobsJobsClient = await getJobsJobsClient();

const {
    data: colleague,
    status,
    refresh,
    error,
} = useLazyAsyncData(`jobs-colleague-${route.params.id as string}`, () => getColleague(parseInt(route.params.id as string)));

async function getColleague(userId: number): Promise<GetColleagueResponse> {
    try {
        const call = jobsJobsClient.getColleague({
            userId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function updateColleageAbsence(value: { userId: number; absenceBegin?: Timestamp; absenceEnd?: Timestamp }): void {
    if (colleague.value?.colleague === undefined) {
        return;
    }

    if (colleague.value.colleague?.props === undefined) {
        colleague.value.colleague.props = {
            userId: colleague.value.colleague.userId,
            job: colleague.value.colleague.job,
        };
    }

    colleague.value.colleague.props.absenceBegin = value.absenceBegin;
    colleague.value.colleague.props.absenceEnd = value.absenceEnd;
}

// Handle the client update event
const { sendClientView } = useClientUpdate(ObjectType.JOBS_COLLEAGUE, () =>
    notifications.add({
        title: { key: 'notifications.jobs.colleague.client_view_update.title', parameters: {} },
        description: { key: 'notifications.jobs.colleague.client_view_update.content', parameters: {} },
        duration: 7500,
        type: NotificationType.INFO,
        actions: [
            {
                label: { key: 'common.refresh', parameters: {} },
                icon: 'i-mdi-refresh',
                onClick: () => refresh(),
            },
        ],
    }),
);
sendClientView(parseInt(route.params.id as string));

const links = computed(() =>
    [
        {
            label: t('common.info'),
            to: { name: 'jobs-colleagues-id-info', params: { id: route.params?.id ?? 0 } },
            icon: 'i-mdi-information-outline',
            permission: 'jobs.JobsService/GetColleague' as Perms,
        },
        {
            label: t('common.activity'),
            to: { name: 'jobs-colleagues-id-activity', params: { id: route.params?.id ?? 0 } },
            icon: 'i-mdi-pulse',
            permission: 'jobs.JobsService/ListColleagueActivity' as Perms,
        },
        {
            label: t('common.timeclock'),
            to: { name: 'jobs-colleagues-id-timeclock', params: { id: route.params?.id ?? 0 } },
            icon: 'i-mdi-timeline-clock',
            permission: 'jobs.TimeclockService/ListTimeclock' as Perms,
            check: attr('jobs.TimeclockService/ListTimeclock', 'Access', 'All').value,
        },
        {
            label: t('pages.qualifications.title'),
            to: { name: 'jobs-colleagues-id-qualifications', params: { id: route.params?.id ?? 0 } },
            icon: 'i-mdi-school',
            permission: 'qualifications.QualificationsService/ListQualifications' as Perms,
        },
        {
            label: t('pages.jobs.conduct.title'),
            to: { name: 'jobs-colleagues-id-conduct', params: { id: route.params?.id ?? 0 } },
            icon: 'i-mdi-list-status',
            permission: 'jobs.ConductService/ListConductEntries' as Perms,
        },
    ].filter((tab) => can(tab.permission).value && (tab.check === undefined || tab.check === true)),
);
</script>

<template>
    <UDashboardPanelContent>
        <DataPendingBlock
            v-if="!colleague && isRequestPending(status)"
            :message="$t('common.loading', [$t('common.colleague', 1)])"
        />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.colleague', 1)])"
            :error="error"
            :retry="refresh"
        />
        <DataNoDataBlock v-else-if="!colleague || !colleague.colleague" />

        <template v-else>
            <ColleagueInfo :colleague="colleague.colleague" @update:absence-dates="updateColleageAbsence($event)" />

            <UDashboardToolbar>
                <UNavigationMenu orientation="horizontal" :items="links" />
            </UDashboardToolbar>

            <NuxtPage :colleague="colleague.colleague" @refresh="refresh()" />
        </template>
    </UDashboardPanelContent>
</template>
