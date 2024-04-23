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

const tabs = ref<{ key: string; slot: string; label: string; icon: string }[]>([
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
]);
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.qualifications.title')">
                <template #right>
                    <UButton
                        v-if="can('QualificationsService.CreateQualification')"
                        :to="{ name: 'qualifications-create' }"
                        trailing-icon="i-mdi-plus"
                        color="gray"
                    >
                        {{ $t('components.qualifications.create_new_qualification') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <UTabs :items="tabs" :unmount="true">
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
