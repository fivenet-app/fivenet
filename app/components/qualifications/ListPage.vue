<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';

const { t } = useI18n();

const { can } = useAuth();

const items = computed<NavigationMenuItem[]>(() => [
    {
        label: t('components.qualifications.user_qualifications'),
        icon: 'i-mdi-account-circle',
        to: '/qualifications',
    },
    {
        label: t('components.qualifications.all_qualifications'),
        icon: 'i-mdi-view-list',
        to: '/qualifications/all',
    },
]);

const qualifications = await useQualifications();
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.qualifications.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <UTooltip
                        v-if="can('qualifications.QualificationsService/UpdateQualification').value"
                        :text="$t('common.create')"
                    >
                        <UButton
                            trailing-icon="i-mdi-plus"
                            color="neutral"
                            variant="outline"
                            @click="qualifications.createQualification()"
                        >
                            <span class="hidden truncate sm:block">
                                {{ $t('common.qualification', 1) }}
                            </span>
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <UNavigationMenu orientation="horizontal" :items="items" class="-mx-1 flex-1" />
            </UDashboardToolbar>
        </template>

        <template #body>
            <UContainer class="p-4 sm:p-4">
                <slot />
            </UContainer>
        </template>
    </UDashboardPanel>
</template>
