<script lang="ts" setup>
import CitizensAttributesModal from '~/components/citizens/CitizensAttributesModal.vue';
import CitizensList from '~/components/citizens/CitizensList.vue';

useHead({
    title: 'pages.citizens.title',
});
definePageMeta({
    title: 'pages.citizens.title',
    requiresAuth: true,
    permission: 'CitizenStoreService.ListCitizens',
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
                        v-if="can('CitizenStoreService.ManageCitizenAttributes').value"
                        :label="$t('common.attributes', 2)"
                        icon="i-mdi-tag"
                        @click="modal.open(CitizensAttributesModal, {})"
                    />
                </template>
            </UDashboardNavbar>

            <CitizensList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
