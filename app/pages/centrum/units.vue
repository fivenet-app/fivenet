<script lang="ts" setup>
import UnitsList from '~/components/centrum/settings/UnitsList.vue';

const { t } = useI18n();

useHead({
    title: () => t('common.unit', 2),
});

definePageMeta({
    title: 'common.unit',
    requiresAuth: true,
    permission: 'CentrumService.CreateOrUpdateUnit',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UnitsList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
