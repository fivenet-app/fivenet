<script lang="ts" setup>
import CountUp from 'vue-countup-v3';
import { useAuthStore } from '~/store/auth';
import type { Stat } from '~~/gen/ts/resources/stats/stats';

useHead({
    title: 'pages.stats.title',
});
definePageMeta({
    title: 'pages.stats.title',
    layout: 'landing',
    requiresAuth: true,
    authTokenOnly: true,
    showCookieOptions: true,
    redirectIfAuthed: false,
});

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

type Stats = { [key: string]: Stat & { unit?: string; icon?: string } };

const defaultStats: Stats = {
    users_registered: {
        icon: 'i-mdi-user',
    },
    documents_created: {
        icon: 'i-mdi-file-document-box-multiple',
    },
    dispatches_created: {
        icon: 'i-mdi-car-emergency',
    },
    citizen_activity: {
        icon: 'i-mdi-pulse',
    },
    timeclock_tracked: {
        unit: 'common.time_ago.year',
        icon: 'i-mdi-timeline-clock',
    },
    citizens_total: {
        icon: 'i-mdi-user-group',
    },
};

type StatsState = { stats: Stats; fetchedAt?: number };
const state = useState<StatsState>('stats', () => ({ stats: defaultStats, fetchedAt: undefined }));

const { data: stats, pending: loading } = useLazyAsyncData('stats', () => getStats(), {
    transform: (input): StatsState => ({
        stats: input,
        fetchedAt: new Date().getTime(),
    }),
    getCachedData() {
        if (!state.value.fetchedAt) {
            return undefined;
        }

        const expireDate = new Date(state.value.fetchedAt);
        expireDate.setTime(expireDate.getTime() + 60 * 1000);
        if (expireDate.getTime() < Date.now()) {
            return undefined;
        }

        return state.value;
    },
});

async function getStats(): Promise<Stats> {
    try {
        const call = getGRPCStatsClient().getStats({});
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
    if (website.statsPage) {
        return;
    }

    if (activeChar.value === null) {
        await navigateTo('/');
    } else {
        await navigateTo('/overview');
    }
});
</script>

<template>
    <UPage>
        <UDashboardPanel>
            <div class="flex flex-col justify-between">
                <div>
                    <div class="relative isolate px-6 py-20 lg:px-8">
                        <div
                            class="hero absolute inset-0 z-[-1] [mask-image:radial-gradient(100%_100%_at_top,white,transparent)]"
                        />

                        <div class="mx-auto max-w-2xl text-center">
                            <h2 class="text-4xl font-bold tracking-tight sm:text-6xl">
                                {{ $t('pages.stats.title') }}
                            </h2>
                            <p class="mt-6 text-lg leading-8">
                                {{ $t('pages.stats.subtitle') }}
                            </p>
                        </div>
                    </div>

                    <ULandingSection>
                        <UPageGrid>
                            <ULandingCard
                                v-for="(stat, key) in stats?.stats"
                                :key="key"
                                :title="$t(`pages.stats.stats.${key}`)"
                                :icon="stat?.icon"
                            >
                                <template #description>
                                    <p
                                        class="mt-2 flex w-full items-center gap-x-2 text-2xl font-semibold tracking-tight text-gray-900 dark:text-white"
                                    >
                                        <USkeleton v-if="loading || stat?.value === undefined" class="h-8 w-[175px]" />
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
                            </ULandingCard>
                        </UPageGrid>
                    </ULandingSection>
                </div>
            </div>
        </UDashboardPanel>
    </UPage>
</template>
