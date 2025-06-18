<script lang="ts" setup>
import CitizensLabelsModal from '~/components/citizens/CitizensLabelsModal.vue';
import CitizensList from '~/components/citizens/CitizensList.vue';

useHead({
    title: 'pages.citizens.title',
});

definePageMeta({
    title: 'pages.citizens.title',
    requiresAuth: true,
    permission: 'citizens.CitizensService/ListCitizens',
});

const { can } = useAuth();

const modal = useModal();
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.citizens.title')">
                <template #right>
                    <UButton
                        v-if="can('citizens.CitizensService/ManageLabels').value"
                        :label="$t('common.label', 2)"
                        icon="i-mdi-tag"
                        @click="modal.open(CitizensLabelsModal, {})"
                    />
                </template>
            </UDashboardNavbar>

            <CitizensList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
