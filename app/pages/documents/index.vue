<script lang="ts" setup>
import DocumentList from '~/components/documents/DocumentList.vue';
import PinnedDocumentsList from '~/components/documents/PinnedDocumentsList.vue';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';

useHead({
    title: 'pages.documents.title',
});

definePageMeta({
    title: 'pages.documents.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService.ListDocuments',
});

const modal = useModal();

const { can } = useAuth();

const isOpen = ref(false);
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel class="shrink-0 border-b border-gray-200 lg:border-b-0 lg:border-r dark:border-gray-800" grow>
            <UDashboardNavbar :title="$t('pages.documents.title')">
                <template #right>
                    <UButton class="2xl:hidden" trailing-icon="i-mdi-pin" color="white" truncate @click="isOpen = true">
                        <span class="hidden truncate sm:block">
                            {{ $t('common.pinned') }}
                        </span>
                    </UButton>

                    <UButtonGroup class="inline-flex">
                        <UButton
                            v-if="can('completor.CompletorService.CompleteDocumentCategories').value"
                            :to="{ name: 'documents-categories' }"
                            icon="i-mdi-shape"
                            truncate
                        >
                            <span class="hidden truncate sm:block">
                                {{ $t('common.category', 2) }}
                            </span>
                        </UButton>

                        <UButton
                            v-if="can('documents.DocumentsService.ListTemplates').value"
                            :to="{ name: 'documents-templates' }"
                            icon="i-mdi-file-code"
                            truncate
                        >
                            <span class="hidden truncate sm:block">
                                {{ $t('common.template', 2) }}
                            </span>
                        </UButton>
                    </UButtonGroup>

                    <UButton
                        v-if="can('documents.DocumentsService.CreateDocument').value"
                        trailing-icon="i-mdi-plus"
                        color="gray"
                        truncate
                        @click="modal.open(TemplatesModal, {})"
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.create') }}
                        </span>
                    </UButton>
                </template>
            </UDashboardNavbar>

            <DocumentList />
        </UDashboardPanel>

        <UDashboardPanel
            id="documentspinnedlist"
            v-model="isOpen"
            class="overflow-x-hidden"
            collapsible
            side="right"
            breakpoint="2xl"
            :width="350"
            :resizable="{ min: 275, max: 600 }"
            :ui="{ collapsible: 'lg:!hidden 2xl:!flex', slideover: 'lg:!flex 2xl:hidden' }"
        >
            <PinnedDocumentsList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
