<script lang="ts" setup>
import DocumentList from '~/components/documents/DocumentList.vue';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';

useHead({
    title: 'pages.documents.title',
});
definePageMeta({
    title: 'pages.documents.title',
    requiresAuth: true,
    permission: 'DocStoreService.ListDocuments',
});

const modal = useModal();
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.documents.title')">
                <template #right>
                    <UButtonGroup class="inline-flex">
                        <UButton
                            v-if="can('CompletorService.CompleteDocumentCategories')"
                            :to="{ name: 'documents-categories' }"
                            icon="i-mdi-shape"
                        >
                            {{ $t('common.category', 2) }}
                        </UButton>

                        <UButton
                            v-if="can('DocStoreService.ListTemplates')"
                            :to="{ name: 'documents-templates' }"
                            icon="i-mdi-file-code"
                        >
                            {{ $t('common.template', 2) }}
                        </UButton>

                        <UButton
                            :label="$t('common.create')"
                            trailing-icon="i-mdi-plus"
                            color="gray"
                            @click="modal.open(TemplatesModal, {})"
                        />
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>

            <DocumentList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
