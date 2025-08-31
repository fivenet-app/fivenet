<script lang="ts" setup>
import { getJobsTimeclockClient } from '~~/gen/ts/clients';
import type { GetTimeclockStatsResponse } from '~~/gen/ts/services/jobs/timeclock';

const props = withDefaults(
    defineProps<{
        userId?: number;
        hideHeader?: boolean;
        loading?: boolean;
    }>(),
    {
        userId: undefined,
        hideHeader: false,
        loading: false,
    },
);

defineEmits<{
    (e: 'refresh'): void;
}>();

const jobsTimeclockClient = await getJobsTimeclockClient();

const { data, error, status, refresh } = useLazyAsyncData(`jobs-timeclock-stats`, () => getTimeclockStats());

async function getTimeclockStats(): Promise<GetTimeclockStatsResponse> {
    try {
        const call = jobsTimeclockClient.getTimeclockStats({
            userId: props.userId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canRefresh = ref(true);
const refreshThrottle = useThrottleFn(async () => {
    canRefresh.value = false;
    await refresh().finally(() => useTimeoutFn(() => (canRefresh.value = true), 400));
}, 2500);

const loadingState = ref(false);
watch(
    () => props.loading,
    () => {
        if (props.loading) {
            loadingState.value = true;
        }
    },
);
watchDebounced(
    () => props.loading,
    () => {
        if (!props.loading) {
            loadingState.value = false;
        }
    },
    {
        debounce: 750,
        maxWait: 1250,
    },
);
</script>

<template>
    <UCard>
        <template v-if="!hideHeader" #header>
            <h2 class="inline-flex w-full items-center justify-between text-lg font-semibold">
                {{ $t('common.timeclock') }}

                <UTooltip :text="$t('common.refresh')">
                    <UButton
                        variant="link"
                        icon="i-mdi-refresh"
                        :disabled="loading || loadingState"
                        :loading="loading || loadingState"
                        :label="$t('common.refresh')"
                        @click="$emit('refresh')"
                    />
                </UTooltip>
            </h2>
        </template>

        <LazyJobsTimeclockStatsBlock
            :stats="data?.stats"
            :weekly="data?.weekly"
            :failed="!!error"
            :loading="isRequestPending(status)"
            @refresh="refreshThrottle"
        />
    </UCard>
</template>
