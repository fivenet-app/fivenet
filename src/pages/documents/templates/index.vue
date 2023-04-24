<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import TemplatesList from '~/components/documents/templates/TemplatesList.vue';
import { DocumentTemplateShort } from '@fivenet/gen/resources/documents/templates_pb';

const router = useRouter();

useHead({
    title: 'Templates',
});
definePageMeta({
    title: 'Templates',
    requiresAuth: true,
    permission: 'DocStoreService.FindDocuments',
});

async function selected(t: DocumentTemplateShort): Promise<void> {
    await router.push({ name: 'documents-templates-id', params: { id: t.getId() } });
}
</script>

<template>
    <ContentWrapper>
        <div class="py-2">
            <div class="px-2 sm:px-6 lg:px-8">
                <div v-if="'DocStoreService.CreateTemplate'" class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <div class="flex flex-row items-center gap-2 sm:mx-auto">
                            <div class="flex-initial form-control" v-can="'DocStoreService.CreateDocument'">
                                <NuxtLink :to="{ name: 'documents-templates-create' }"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">
                                    {{ $t('pages.documents.templates.create_template') }}
                                </NuxtLink>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="flow-root mt-2">
                    <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                        <TemplatesList @selected="selected($event)" />
                    </div>
                </div>
            </div>
        </div>
    </ContentWrapper>
</template>
