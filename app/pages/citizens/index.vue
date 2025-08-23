<script lang="ts" setup>
import CitizenLabelModal from '~/components/citizens/CitizenLabelModal.vue';
import CitizenList from '~/components/citizens/CitizenList.vue';

useHead({
    title: 'pages.citizens.title',
});

definePageMeta({
    title: 'pages.citizens.title',
    requiresAuth: true,
    permission: 'citizens.CitizensService/ListCitizens',
});

const { can } = useAuth();

const modal = useOverlay();
</script>

<template>
    <UDashboardPanel>
        <UDashboardNavbar :title="$t('pages.citizens.title')">
            <template #right>
                <UButton
                    v-if="can('citizens.CitizensService/ManageLabels').value"
                    :label="$t('common.label', 2)"
                    icon="i-mdi-tag"
                    @click="modal.open(CitizenLabelModal, {})"
                />
            </template>
        </UDashboardNavbar>

        <CitizenList />
    </UDashboardPanel>
</template>
