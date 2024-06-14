<script lang="ts" setup>
import QualificationsList from '~/components/qualifications/QualificationsList.vue';
import QualificationsRequestsList from '~/components/qualifications/QualificationsRequestsList.vue';
import QualificationsResultsList from '~/components/qualifications/QualificationsResultsList.vue';

useHead({
    title: 'pages.qualifications.title',
});
definePageMeta({
    title: 'pages.qualifications.title',
    requiresAuth: true,
    permission: 'QualificationsService.ListQualifications',
});

const { t } = useI18n();

const items: { key: string; slot: string; label: string; icon: string }[] = [
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
        router.replace({ query: { tab: items[value].slot }, hash: '#' });
    },
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.qualifications.title')">
                <template #right>
                    <UButton
                        v-if="can('QualificationsService.CreateQualification').value"
                        :to="{ name: 'qualifications-create' }"
                        trailing-icon="i-mdi-plus"
                        color="gray"
                    >
                        {{ $t('components.qualifications.create_new_qualification') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <UTabs v-model="selectedTab" :items="items" :unmount="true" :ui="{ list: { rounded: '' } }">
                <template #default="{ item, selected }">
                    <div class="relative flex items-center gap-2 truncate">
                        <UIcon :name="item.icon" class="size-4 shrink-0" />

                        <span class="truncate">{{ item.label }}</span>

                        <span
                            v-if="selected"
                            class="bg-primary-500 dark:bg-primary-400 absolute -right-4 size-2 rounded-full"
                        />
                    </div>
                </template>

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
        </UDashboardPanel>
    </UDashboardPage>
</template>
