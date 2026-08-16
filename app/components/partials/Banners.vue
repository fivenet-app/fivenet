<script lang="ts" setup>
import { useResizeObserver } from '@vueuse/core';
import ImpersonatingBanner from '~/components/auth/ImpersonatingBanner.vue';
import BannerMessage from '~/components/partials/BannerMessage.vue';

const appConfig = useAppConfig();

const authSessionStore = useAuthSessionStore();
const { userInfo } = storeToRefs(authSessionStore);

const notificationStore = useNotificationsStore();
const { dismissedBannerMessageID } = storeToRefs(notificationStore);

const sidebarWidth = ref<number>(0);
const sidebarEl = ref<HTMLElement | null>(null);

function updateSidebarWidth(entries: ResizeObserverEntry[]) {
    sidebarWidth.value = entries[0]?.contentRect.width ?? 0;
}

const bannerMessageClosed = ref<boolean>(false);

watch(
    () => appConfig.system.bannerMessage,
    () => (bannerMessageClosed.value = false),
);

const bannersPaddingRequired = computed<number>(
    () =>
        (appConfig.system.bannerMessageEnabled &&
        appConfig.system.bannerMessage &&
        dismissedBannerMessageID.value !== appConfig.system.bannerMessage.id &&
        !bannerMessageClosed.value
            ? 16
            : 0) + (userInfo.value?.originalJob ? 5 : 0),
);

watch(
    bannersPaddingRequired,
    () =>
        document.documentElement.style.setProperty(
            '--dashboard-panel-bottom-offset',
            `calc(${Math.max(Math.min(bannersPaddingRequired.value, 21), 2)} * var(--spacing))`,
        ),
    { immediate: true },
);

useResizeObserver(sidebarEl, updateSidebarWidth);

onMounted(() => (sidebarEl.value = document.getElementById('dashboard-sidebar-default')));
</script>

<template>
    <div
        v-if="bannersPaddingRequired > 0"
        class="pointer-events-none fixed inset-x-0 bottom-0 z-[49] flex max-h-21 flex-col gap-0 sm:left-[var(--sidebar-width)]"
        :style="{ '--sidebar-width': `${sidebarWidth}px` }"
    >
        <BannerMessage
            v-if="appConfig.system.bannerMessageEnabled && appConfig.system.bannerMessage"
            :message="appConfig.system.bannerMessage"
            @close="() => (bannerMessageClosed = true)"
        />

        <ImpersonatingBanner v-if="userInfo?.originalJob" :job="userInfo?.job" :job-grade="userInfo?.jobGrade" />
    </div>
</template>
