<script lang="ts" setup>
import List from '~/components/documents/templates/List.vue';
import Pagination from '~/components/partials/Pagination.vue';
import type { TemplateShort } from '~~/gen/ts/resources/documents/templates';

useHead({
    title: 'pages.documents.templates.title',
});

definePageMeta({
    title: 'pages.documents.templates.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService/ListTemplates',
});

const { can } = useAuth();

async function selected(t: TemplateShort | undefined): Promise<void> {
    if (!t) return;

    await navigateTo({ name: 'documents-templates-id', params: { id: t.id } });
}

const templatesListRef = useTemplateRef('templatesListRef');
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('pages.documents.templates.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton to="/documents" />

                    <UButton
                        v-if="can('TODOService/TODOMethod').value"
                        :to="{ name: 'documents-templates-forms' }"
                        icon="i-mdi-form"
                        truncate
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.form', 2) }}
                        </span>
                    </UButton>

                    <UTooltip v-if="can('documents.DocumentsService/CreateTemplate').value" :text="$t('common.create')">
                        <UButton :to="{ name: 'documents-templates-create' }" color="neutral" trailing-icon="i-mdi-plus">
                            <span class="hidden truncate sm:block">
                                {{ $t('common.template') }}
                            </span>
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <List ref="templatesListRef" @selected="selected($event)" />
        </template>

        <template #footer>
            <Pagination
                :status="templatesListRef?.status ?? 'pending'"
                :refresh="templatesListRef?.refresh"
                hide-buttons
                hide-text
            />
        </template>
    </UDashboardPanel>
</template>
