<script lang="ts" setup>
import TemplatesList from '~/components/documents/templates/TemplatesList.vue';
import { TemplateShort } from '~~/gen/ts/resources/documents/templates';

useHead({
    title: 'pages.documents.templates.title',
});
definePageMeta({
    title: 'pages.documents.templates.title',
    requiresAuth: true,
    permission: 'DocStoreService.ListTemplates',
});

async function selected(t: TemplateShort | undefined): Promise<void> {
    if (!t) {
        return;
    }

    await navigateTo({ name: 'documents-templates-id', params: { id: t.id } });
}
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.documents.templates.title')">
                <template #right>
                    <UButtonGroup class="inline-flex">
                        <UButton color="black" icon="i-mdi-arrow-back" to="/documents">
                            {{ $t('common.back') }}
                        </UButton>

                        <UButton
                            v-if="can('DocStoreService.CreateTemplate').value"
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
                <TemplatesList @selected="selected($event)" />
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
