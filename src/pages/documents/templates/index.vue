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
            <UDashboardNavbar :title="$t('pages.documents.templates.title')"> </UDashboardNavbar>

            <div v-if="'DocStoreService.CreateTemplate'" class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <div class="flex flex-row items-center gap-2 sm:mx-auto">
                        <div v-if="can('DocStoreService.CreateTemplate')" class="flex-1">
                            <NuxtLink
                                :to="{ name: 'documents-templates-create' }"
                                class="bg-primary-500 hover:bg-primary-400 inline-flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                            >
                                {{ $t('pages.documents.templates.create_template') }}
                            </NuxtLink>
                        </div>
                    </div>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
                        <TemplatesList @selected="selected($event)" />
                    </div>
                </div>
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>
