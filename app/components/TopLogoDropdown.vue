<script lang="ts" setup>
import type { DropdownMenuItem } from '@nuxt/ui';
import type { RpcError } from '@protobuf-ts/runtime-rpc';
import { breakpointsTailwind, createReusableTemplate, useBreakpoints } from '@vueuse/core';
import { useAuthStore } from '~/stores/auth';

const props = defineProps<{
    collapsed?: boolean;
}>();

const { t } = useI18n();

const [DefineLogoButton, ReuseLogoButton] = createReusableTemplate();
const breakpoints = useBreakpoints(breakpointsTailwind);
const isMobile = breakpoints.smaller('lg');
const { isDashboardSidebarSlideoverOpen } = useDashboard();

const authStore = useAuthStore();
const { selectSessionJob } = authStore;
const { activeChar, isSuperuser, jobProps } = storeToRefs(authStore);

const isSwitchingWorkContext = ref(false);

const fivenetLogo = '/images/logo.webp';
const workContextJobs = computed(() => activeChar.value?.jobs ?? []);
const hasWorkContextSwitcher = computed(() => !isSuperuser.value && workContextJobs.value.length > 1);
const activeWorkContextJob = computed(() => activeChar.value?.job);
const activeWorkContextLabel = computed(() => activeChar.value?.jobLabel || activeChar.value?.job || t('common.job'));
const jobLogo = computed(() =>
    jobProps.value?.logoFile?.filePath ? `/api/filestore/${jobProps.value.logoFile.filePath}` : undefined,
);
const jobIdentityLogo = computed(() => jobLogo.value || fivenetLogo);
const tooltip = computed(() =>
    hasWorkContextSwitcher.value ? `${t('pages.jobs.work_context')}: ${activeWorkContextLabel.value}` : undefined,
);

const isShowingJobIdentity = ref(false);
const isShowingJobLabel = computed(() => isShowingJobIdentity.value);
const identityTransition = computed(() => ({
    appear: true,
    mode: 'out-in' as const,
    enterActiveClass: `transition-all duration-[200ms] ease-out`,
    leaveActiveClass: `transition-all duration-[200ms] ease-in`,
    enterFromClass: 'translate-x-1 opacity-0',
    leaveToClass: '-translate-x-1 opacity-0',
}));
let identityTimer: ReturnType<typeof setTimeout> | undefined;

function resetJobIdentity(): void {
    if (identityTimer) clearTimeout(identityTimer);
    identityTimer = undefined;
    isShowingJobIdentity.value = false;
}

function scheduleJobIdentity(): void {
    resetJobIdentity();

    // On mobile, the sidebar is a drawer. Do not run the handoff until it is visible.
    if (!activeWorkContextJob.value || (isMobile.value && !isDashboardSidebarSlideoverOpen.value)) return;

    identityTimer = setTimeout(
        () => {
            isShowingJobIdentity.value = true;
            identityTimer = undefined;
        },
        isMobile.value ? 650 : 800,
    );
}

watch(
    [isMobile, isDashboardSidebarSlideoverOpen, isSuperuser, activeWorkContextJob, jobLogo, activeWorkContextLabel],
    scheduleJobIdentity,
);
onMounted(scheduleJobIdentity);
onBeforeUnmount(resetJobIdentity);

async function switchWorkContext(job: string): Promise<void> {
    if (isSwitchingWorkContext.value || job === activeChar.value?.job) return;

    isSwitchingWorkContext.value = true;
    try {
        await selectSessionJob(job);
    } catch (e) {
        handleGRPCError(e as RpcError);
    } finally {
        isSwitchingWorkContext.value = false;
    }
}

const workContextItems = computed<DropdownMenuItem[][]>(() => [
    workContextJobs.value.map((job) => ({
        label: job.jobLabel || job.job,
        description: job.isPrimary ? t('pages.jobs.work_context_primary') : t('pages.jobs.work_context_session'),
        icon: job.job === activeChar.value?.job ? 'i-mdi-check' : 'i-mdi-briefcase-outline',
        disabled: isSwitchingWorkContext.value,
        onSelect: () => switchWorkContext(job.job),
    })),
]);
</script>

<template>
    <DefineLogoButton>
        <UButton
            class="w-full font-bold"
            color="neutral"
            variant="ghost"
            :square="props.collapsed"
            :loading="isSwitchingWorkContext"
            :disabled="isSwitchingWorkContext"
            :trailing-icon="props.collapsed || !hasWorkContextSwitcher ? undefined : 'i-mdi-chevron-down'"
            :ui="{
                base: props.collapsed ? 'w-40' : 'w-(--reka-dropdown-menu-trigger-width)',
                label: 'min-w-0 flex-1 text-left text-highlighted',
                trailingIcon: 'ml-auto',
            }"
        >
            <template #leading>
                <Transition v-bind="identityTransition">
                    <NuxtImg
                        :key="isShowingJobIdentity ? `job-${activeWorkContextJob}` : 'fivenet'"
                        class="size-5 shrink-0 rounded-none object-contain"
                        :src="isShowingJobIdentity ? jobIdentityLogo : fivenetLogo"
                        :alt="isShowingJobIdentity ? activeWorkContextLabel : 'FiveNet'"
                    />
                </Transition>
            </template>

            <template #default="{ ui }">
                <Transition v-if="!props.collapsed" v-bind="identityTransition">
                    <span :key="isShowingJobLabel ? activeWorkContextLabel : 'FiveNet'" :class="ui.label()">
                        {{ isShowingJobLabel ? activeWorkContextLabel : 'FiveNet' }}
                    </span>
                </Transition>
            </template>
        </UButton>
    </DefineLogoButton>

    <UDropdownMenu
        v-if="hasWorkContextSwitcher"
        :items="workContextItems"
        :content="{ align: 'start', collisionPadding: 12 }"
        :ui="{ content: props.collapsed ? 'w-64' : 'w-(--reka-dropdown-menu-trigger-width)' }"
    >
        <UTooltip :text="tooltip">
            <ReuseLogoButton />
        </UTooltip>
    </UDropdownMenu>

    <UTooltip v-else :text="tooltip">
        <ReuseLogoButton />
    </UTooltip>

    <Teleport to="body">
        <div
            v-if="isSwitchingWorkContext"
            class="fixed inset-0 z-100 flex items-center justify-center bg-default/60 backdrop-blur-[1px]"
            role="status"
            aria-live="polite"
            :aria-label="$t('common.loading', [$t('pages.jobs.work_context')])"
        >
            <UIcon name="i-mdi-loading" class="size-26 animate-spin text-white" />
        </div>
    </Teleport>
</template>
