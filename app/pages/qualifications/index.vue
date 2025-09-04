<script lang="ts" setup>
import type { TabsItem } from '@nuxt/ui';
import QualificationList from '~/components/qualifications/QualificationList.vue';
import QualificationRequestList from '~/components/qualifications/QualificationRequestList.vue';
import QualificationResultList from '~/components/qualifications/QualificationResultList.vue';

useHead({
    title: 'pages.qualifications.title',
});

definePageMeta({
    title: 'pages.qualifications.title',
    requiresAuth: true,
    permission: 'qualifications.QualificationsService/ListQualifications',
});

const { t } = useI18n();

const { can } = useAuth();

const items: TabsItem[] = [
    {
        key: 'yours',
        slot: 'yours' as const,
        label: t('components.qualifications.user_qualifications'),
        icon: 'i-mdi-account-circle',
        value: 'yours',
    },
    {
        key: 'all',
        slot: 'all' as const,
        label: t('components.qualifications.all_qualifications'),
        icon: 'i-mdi-view-list',
        value: 'all',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'yours';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
    },
});

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
                        <UButton trailing-icon="i-mdi-plus" color="neutral" @click="qualifications.createQualification()">
                            <span class="hidden truncate sm:block">
                                {{ $t('common.qualification', 1) }}
                            </span>
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <UTabs v-model="selectedTab" :items="items" variant="link">
                <template #yours>
                    <UContainer>
                        <div class="flex flex-col gap-2">
                            <QualificationResultList />

                            <QualificationRequestList />
                        </div>
                    </UContainer>
                </template>
                <template #all>
                    <UContainer>
                        <QualificationList />
                    </UContainer>
                </template>
            </UTabs>
        </template>
    </UDashboardPanel>
</template>
