<script lang="ts" setup>
import { useAuthStore } from '~/stores/auth';
import { getStatsStatsClient } from '~~/gen/ts/clients';
import type { Perms } from '~~/gen/ts/perms';
import type { Stat } from '~~/gen/ts/resources/stats/stats';

useHead({
    title: 'pages.stats.title',
});

definePageMeta({
    title: 'pages.stats.title',
    requiresAuth: true,
    authTokenOnly: true,
    redirectIfAuthed: false,
});

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const { can } = useAuth();

const statsStatsClient = await getStatsStatsClient();

const CountUp = defineAsyncComponent(async () => {
    const m = await import('vue-countup-v3');
    return m.default;
});

type Stats = { [key: string]: Stat & { unit?: string; icon?: string; permission?: Perms; to?: string } };

const defaultStats: Stats = {
    users_registered: {
        icon: 'i-mdi-user',
    },
    documents_created: {
        icon: 'i-mdi-file-document-box-multiple',
        permission: 'documents.DocumentsService/ListDocuments' as Perms,
        to: '/documents',
    },
    dispatches_created: {
        icon: 'i-mdi-car-emergency',
        permission: 'livemap.LivemapService/Stream' as Perms,
        to: '/livemap',
    },
    citizen_activity: {
        icon: 'i-mdi-pulse',
        permission: 'citizens.CitizensService/ListCitizens' as Perms,
        to: '/citizens',
    },
    timeclock_tracked: {
        unit: 'common.time_ago.year',
        icon: 'i-mdi-timeline-clock',
        permission: 'jobs.TimeclockService/ListTimeclock' as Perms,
        to: '/jobs/timeclock',
    },
    citizens_total: {
        icon: 'i-mdi-user-group',
        permission: 'citizens.CitizensService/ListCitizens' as Perms,
        to: '/citizens',
    },
};

type StatsState = { stats: Stats; fetchedAt?: number };

const state = useState<StatsState>('stats', () => ({ stats: defaultStats, fetchedAt: undefined }));

const { data: stats, status } = useLazyAsyncData<StatsState>(
    'stats',
    async () => {
        const nextState: StatsState = {
            stats: await getPublicStats(),
            fetchedAt: Date.now(),
        };
        state.value = nextState;
        return nextState;
    },
    {
        default: () => state.value,
        getCachedData(_, ctx) {
            if (!state.value.fetchedAt || ctx.cause === 'refresh:manual') return undefined;

            const expireDate = new Date(state.value.fetchedAt);
            expireDate.setTime(expireDate.getTime() + 60 * 1000);
            if (expireDate.getTime() < Date.now()) return undefined;

            return state.value;
        },
    },
);

const displayStats = computed<Stats>(() => stats.value?.stats ?? defaultStats);

async function getPublicStats(): Promise<Stats> {
    try {
        const call = statsStatsClient.getPublicStats({});
        const { response } = await call;

        const stats = { ...defaultStats };
        for (const key in response.stats) {
            if (defaultStats[key]) {
                stats[key] = { ...response.stats[key], ...defaultStats[key] };
            }
        }
        return stats;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { website } = useAppConfig();

onBeforeMount(async () => {
    if (website.statsPage) return;

    if (activeChar.value === null) {
        await navigateTo('/');
    } else {
        await navigateTo('/overview');
    }
});
</script>

<template>
    <UDashboardPanel :ui="{ root: 'pb-(--page-content-bottom-offset)' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.stats.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div class="grid grid-cols-2 gap-4">
                <UPageCard
                    v-for="(stat, key) in displayStats"
                    :key="key"
                    :title="$t(`pages.stats.stats.${key}`)"
                    :icon="stat.icon"
                    :to="stat.permission && can(stat.permission).value ? stat.to : undefined"
                    :ui="{ leadingIcon: 'size-12' }"
                >
                    <template #description>
                        <p class="flex w-full items-center gap-x-1 text-2xl font-semibold tracking-tight text-highlighted">
                            <USkeleton v-if="isRequestPending(status) || stat.value === undefined" class="h-8 w-[175px]" />
                            <ClientOnly v-else>
                                <CountUp
                                    :start-val="0"
                                    :end-val="stat.value"
                                    :options="{ enableScrollSpy: true, scrollSpyOnce: true }"
                                />

                                <span v-if="stat.unit !== undefined">
                                    {{ $t(stat.unit ?? 'common.time_ago.week', 2) }}
                                </span>
                            </ClientOnly>
                        </p>
                    </template>
                </UPageCard>
            </div>
        </template>
    </UDashboardPanel>
</template>
