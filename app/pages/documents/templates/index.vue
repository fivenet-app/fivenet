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
    permission: 'documents.DocumentsService/ListTemplates',
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

                    <UTooltip v-if="can('documents.DocumentsService/CreateTemplate').value" :text="$t('common.create')">
                        <UButton :to="{ name: 'documents-templates-create' }" color="gray" trailing-icon="i-mdi-plus">
                            <span class="hidden truncate sm:block">
                                {{ $t('common.template') }}
                            </span>
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardNavbar>

            <UDashboardPanelContent>
                <TemplatesList ref="templatesListRef" @selected="selected($event)" />
            </UDashboardPanelContent>

            <Pagination
                :loading="isRequestPending(templatesListRef?.status ?? 'pending')"
                :refresh="templatesListRef?.refresh"
                hide-buttons
                hide-text
            />
        </UDashboardPanel>
    </UDashboardPage>
</template>
