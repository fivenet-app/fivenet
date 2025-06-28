<script lang="ts" setup>
import UnitList from '~/components/centrum/settings/UnitList.vue';

const { t } = useI18n();

useHead({
    title: () => t('common.unit', 2),
});

definePageMeta({
    title: 'common.unit',
    requiresAuth: true,
    permission: 'centrum.CentrumService/CreateOrUpdateUnit',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UnitList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
