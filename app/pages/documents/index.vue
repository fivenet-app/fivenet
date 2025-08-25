<script lang="ts" setup>
import DocumentList from '~/components/documents/DocumentList.vue';
import PinnedDocumentList from '~/components/documents/PinnedDocumentList.vue';
import TemplateModal from '~/components/documents/templates/TemplateModal.vue';

useHead({
    title: 'pages.documents.title',
});

definePageMeta({
    title: 'pages.documents.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService/ListDocuments',
});

const overlay = useOverlay();

const { can } = useAuth();

const isOpen = ref(false);

const templateModal = overlay.create(TemplateModal);
</script>

<template>
    <UDashboardPanel class="shrink-0 border-b border-gray-200 lg:border-r lg:border-b-0 dark:border-gray-800">
        <template #header>
            <UDashboardNavbar :title="$t('pages.documents.title')">
                <template #right>
                    <UButton class="2xl:hidden" trailing-icon="i-mdi-pin" color="neutral" truncate @click="isOpen = true">
                        <span class="hidden truncate sm:block">
                            {{ $t('common.pinned') }}
                        </span>
                    </UButton>

                    <UButtonGroup class="inline-flex">
                        <UButton
                            v-if="can('completor.CompletorService/CompleteDocumentCategories').value"
                            :to="{ name: 'documents-categories' }"
                            icon="i-mdi-shape"
                            truncate
                        >
                            <span class="hidden truncate sm:block">
                                {{ $t('common.category', 2) }}
                            </span>
                        </UButton>

                        <UButton
                            v-if="can('documents.DocumentsService/ListTemplates').value"
                            :to="{ name: 'documents-templates' }"
                            icon="i-mdi-file-code"
                            truncate
                        >
                            <span class="hidden truncate sm:block">
                                {{ $t('common.template', 2) }}
                            </span>
                        </UButton>
                    </UButtonGroup>

                    <UTooltip v-if="can('documents.DocumentsService/UpdateDocument').value" :text="$t('common.create')">
                        <UButton trailing-icon="i-mdi-plus" color="neutral" truncate @click="templateModal.open({})">
                            <span class="hidden truncate sm:block">
                                {{ $t('common.document', 1) }}
                            </span>
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <DocumentList />
        </template>
    </UDashboardPanel>

    <UDashboardPanel
        id="documents-pinnedlist"
        v-model:open="isOpen"
        class="overflow-x-hidden"
        side="right"
        breakpoint="2xl"
        :width="350"
        :min-size="25"
        :max-size="50"
        :ui="{ collapsible: 'lg:hidden! 2xl:flex!', slideover: 'lg:flex! 2xl:hidden' }"
    >
        <PinnedDocumentList />
    </UDashboardPanel>
</template>
