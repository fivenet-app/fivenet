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

async function selected(t: TemplateShort): Promise<void> {
    await navigateTo({ name: 'documents-templates-id', params: { id: t.id } });
}
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.documents.templates.title')">
                <template #right>
                    <UButton
                        v-if="can('DocStoreService.CreateTemplate')"
                        :to="{ name: 'documents-templates-create' }"
                        color="gray"
                        trailing-icon="i-mdi-plus"
                    >
                        {{ $t('pages.documents.templates.create_template') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <div class="inline-block min-w-full px-1 py-2 align-middle">
                <TemplatesList @selected="selected($event)" />
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>
