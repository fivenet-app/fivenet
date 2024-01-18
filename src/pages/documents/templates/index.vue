<script lang="ts" setup>
import List from '~/components/documents/templates/List.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
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
    <ContentWrapper>
        <div class="py-2">
            <div class="px-1 sm:px-2 lg:px-4">
                <div v-if="'DocStoreService.CreateTemplate'" class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <div class="flex flex-row items-center gap-2 sm:mx-auto">
                            <div v-if="can('DocStoreService.CreateTemplate')" class="form-control flex-1">
                                <NuxtLink
                                    :to="{ name: 'documents-templates-create' }"
                                    class="inline-flex w-full justify-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
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
                            <List @selected="selected($event)" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </ContentWrapper>
</template>
