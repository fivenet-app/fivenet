<script lang="ts" setup>
import type { DropdownMenuItem } from '@nuxt/ui';
import type { RpcError } from '@protobuf-ts/runtime-rpc';
import { useAuthStore } from '~/stores/auth';

const props = defineProps<{
    collapsed?: boolean;
}>();

const version = APP_VERSION;

const { t } = useI18n();

const authStore = useAuthStore();
const { selectSessionJob } = authStore;
const { activeChar, isSuperuser } = storeToRefs(authStore);

const isSwitchingWorkContext = ref(false);

const workContextJobs = computed(() => activeChar.value?.jobs ?? []);
const hasWorkContextSwitcher = computed(() => !isSuperuser.value && workContextJobs.value.length > 1);
const activeWorkContextLabel = computed(() => activeChar.value?.jobLabel || activeChar.value?.job || t('common.job'));
const title = computed(() => (activeChar.value ? `FiveNet – ${activeWorkContextLabel.value}` : 'FiveNet'));
const tooltip = computed(() =>
    hasWorkContextSwitcher.value ? `${t('pages.jobs.work_context')}: ${activeWorkContextLabel.value}` : version,
);

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
    <UDropdownMenu
        v-if="hasWorkContextSwitcher"
        :items="workContextItems"
        :content="{ align: 'start', collisionPadding: 12 }"
        :ui="{ content: props.collapsed ? 'w-64' : 'w-(--reka-dropdown-menu-trigger-width)' }"
    >
        <UTooltip :text="tooltip">
            <UButton
                class="w-full"
                color="neutral"
                variant="ghost"
                :avatar="{
                    src: '/images/logo.webp',
                    alt: 'FiveNet',
                }"
                :loading="isSwitchingWorkContext"
                :disabled="isSwitchingWorkContext"
                :label="props.collapsed ? undefined : title"
                :trailing-icon="props.collapsed ? undefined : 'i-mdi-chevron-down'"
                :ui="{
                    base: props.collapsed ? 'w-40' : 'w-(--reka-dropdown-menu-trigger-width)',
                    leadingAvatar: 'rounded-none',
                    label: 'min-w-0 flex-1 text-left font-semibold text-highlighted',
                }"
            />
        </UTooltip>
    </UDropdownMenu>

    <UTooltip v-else :text="tooltip">
        <UButton
            class="w-full"
            color="neutral"
            variant="ghost"
            :avatar="{
                src: '/images/logo.webp',
                alt: 'FiveNet',
            }"
            :label="props.collapsed ? undefined : title"
            :ui="{
                base: props.collapsed ? 'w-40' : 'w-(--reka-dropdown-menu-trigger-width)',
                leadingAvatar: 'rounded-none',
                label: 'min-w-0 flex-1 text-left font-semibold text-highlighted',
            }"
        />
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
