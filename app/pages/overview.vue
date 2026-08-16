<script lang="ts" setup>
import type { ContextMenuItem } from '@nuxt/ui';
import HintsBox from '~/components/HintsBox.vue';
import CardsList from '~/components/partials/CardsList.vue';
import QuickAccessList from '~/components/overview/QuickAccessList.vue';
import type { CardElement } from '~/utils/types';

useHead({
    title: 'common.overview',
});

definePageMeta({
    title: 'common.overview',
    requiresAuth: true,
});

const { t } = useI18n();

const items = useOverviewFeatures();

const settingsStore = useSettingsStore();
const { isOverviewQuickAccess, toggleOverviewQuickAccess } = settingsStore;

const getContextMenuItems = (item: CardElement): ContextMenuItem[][] => {
    if (!item.to) return [];

    const pinned = isOverviewQuickAccess(item.to);
    return [
        [
            {
                label: pinned ? t('common.unpin') : t('common.pin'),
                icon: pinned ? 'i-mdi-pin-off' : 'i-mdi-pin',
                onSelect: () => item.to && toggleOverviewQuickAccess(item.to),
            },
        ],
    ];
};
</script>

<template>
    <UDashboardPanel id="overview" :ui="{ root: 'pb-(--dashboard-panel-bottom-offset)' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.overview')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <QuickAccessList />

            <CardsList :items="items" :get-context-menu-items="getContextMenuItems" />

            <div class="max-w-(--breakpoint-lg) sm:mx-auto">
                <HintsBox />
            </div>
        </template>
    </UDashboardPanel>
</template>
