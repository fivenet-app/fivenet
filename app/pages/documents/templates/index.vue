<script lang="ts" setup>
import TemplatesList from '~/components/documents/templates/TemplatesList.vue';
import Pagination from '~/components/partials/Pagination.vue';
import type { TemplateShort } from '~~/gen/ts/resources/documents/templates';

useHead({
    title: 'pages.documents.templates.title',
});

definePageMeta({
    title: 'pages.documents.templates.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService.ListTemplates',
});

const { can } = useAuth();

async function selected(t: TemplateShort | undefined): Promise<void> {
    if (!t) {
        return;
    }

    await navigateTo({ name: 'documents-templates-id', params: { id: t.id } });
}

const templatesListRef = useTemplateRef('templatesListRef');
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.documents.templates.title')">
                <template #right>
                    <PartialsBackButton to="/documents" />

                    <UButtonGroup class="inline-flex">
                        <UButton
                            v-if="can('documents.DocumentsService.CreateTemplate').value"
                            :to="{ name: 'documents-templates-create' }"
                            color="gray"
                            trailing-icon="i-mdi-plus"
                        >
                            {{ $t('pages.documents.templates.create_template') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardPanelContent>
                <TemplatesList ref="templatesListRef" @selected="selected($event)" />
            </UDashboardPanelContent>

            <Pagination :loading="templatesListRef?.loading" :refresh="templatesListRef?.refresh" hide-buttons hide-text />
        </UDashboardPanel>
    </UDashboardPage>
</template>
