<script lang="ts" setup>
import { useResizeObserver } from '@vueuse/core';

const props = defineProps<{
    job?: string;
    jobGrade?: number;
}>();

const completorStore = useCompletorStore();
const { listJobs } = completorStore;
const { jobs } = storeToRefs(completorStore);

const { isSuperuser } = useAuth();

const authStore = useAuthStore();
const { impersonateJob, setSuperuserMode } = authStore;

const foundJob = computed(() => jobs.value.find((j) => j.name === props.job));
const foundJobGrade = computed(() => foundJob.value?.grades.find((g) => g.grade === props.jobGrade));

const bannerRef = useTemplateRef<{ $el: HTMLElement }>('bannerRef');
const bannerEl = computed(() => bannerRef.value?.$el ?? null);

const bannerImpersonatingBottomOffsetVar = '--banner-impersonating-bottom-offset';

function setBannerImpersonatingBottomOffset(height: number): void {
    if (typeof document === 'undefined') return;

    document.documentElement.style.setProperty(bannerImpersonatingBottomOffsetVar, `${height}px`);
}

function syncBannerImpersonatingBottomOffset(): void {
    setBannerImpersonatingBottomOffset(bannerEl.value?.clientHeight ?? 0);
}

onBeforeMount(async () => listJobs());

useResizeObserver(bannerEl, syncBannerImpersonatingBottomOffset);

watch(bannerEl, syncBannerImpersonatingBottomOffset, { immediate: true });

onBeforeUnmount(() => setBannerImpersonatingBottomOffset(0));
</script>

<template>
    <UBanner
        ref="bannerRef"
        :title="
            $t('common.impersonation_active', {
                job: `${foundJobGrade?.label ?? jobGrade}${isSuperuser ? `&nbsp;-&nbsp;${foundJob?.label ?? job}` : ''} (${jobGrade})`,
            })
        "
        icon="i-mdi-drama-masks"
        :ui="{ root: 'w-full pointer-events-auto', container: 'h-5', title: 'text-xs', icon: 'size-4', close: 'text-xs' }"
        :close="{
            icon: undefined,
            trailingIcon: 'i-mdi-exit-run',
            label: $t('common.stop_impersonation'),
            ui: {
                trailingIcon: 'size-4',
            },
        }"
        @close="() => (isSuperuser ? setSuperuserMode(false) : impersonateJob(-1))"
    />
</template>
