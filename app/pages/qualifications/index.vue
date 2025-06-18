<script lang="ts" setup>
import type { TabItem } from '#ui/types';
import QualificationsList from '~/components/qualifications/QualificationsList.vue';
import QualificationsRequestsList from '~/components/qualifications/QualificationsRequestsList.vue';
import QualificationsResultsList from '~/components/qualifications/QualificationsResultsList.vue';

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

const items: TabItem[] = [
    {
        key: 'yours',
        slot: 'yours',
        label: t('components.qualifications.user_qualifications'),
        icon: 'i-mdi-account-circle',
    },
    {
        key: 'all',
        slot: 'all',
        label: t('components.qualifications.all_qualifications'),
        icon: 'i-mdi-view-list',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value]?.slot }, hash: '#' });
    },
});

const qualifications = useQualifications();
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.qualifications.title')">
                <template #right>
                    <UTooltip
                        v-if="can('qualifications.QualificationsService/UpdateQualification').value"
                        :text="$t('common.create')"
                    >
                        <UButton trailing-icon="i-mdi-plus" color="gray" @click="qualifications.createQualification()">
                            <span class="hidden truncate sm:block">
                                {{ $t('common.qualification', 1) }}
                            </span>
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardNavbar>

            <UDashboardPanelContent class="p-0 sm:pb-0">
                <UTabs v-model="selectedTab" :items="items" :unmount="true" :ui="{ list: { rounded: '' } }">
                    <template #yours>
                        <UContainer>
                            <div class="flex flex-col gap-2">
                                <QualificationsResultsList />

                                <QualificationsRequestsList />
                            </div>
                        </UContainer>
                    </template>
                    <template #all>
                        <UContainer>
                            <QualificationsList />
                        </UContainer>
                    </template>
                </UTabs>
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
