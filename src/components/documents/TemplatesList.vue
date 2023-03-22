<script lang="ts" setup>
import { DocumentTemplate, DocumentTemplateShort } from '@arpanet/gen/resources/documents/documents_pb';
import { GetTemplateRequest, ListTemplatesRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { ref, onBeforeMount } from 'vue';
import { getDocStoreClient } from '../../grpc/grpc';

const templates = ref<Array<DocumentTemplateShort>>([]);
const templateObj = ref<undefined | DocumentTemplate>(undefined);

function findTemplates(): void {
    const req = new ListTemplatesRequest();

    getDocStoreClient().
        listTemplates(req, null).then((resp) => {
            templates.value = resp.getTemplatesList();
        });
}
function getTemplate(id: number): void {
    const req = new GetTemplateRequest();
    req.setTemplateId(id);

    getDocStoreClient().
        getTemplate(req, null).then((resp) => {
            templateObj.value = resp.getTemplate();
        });
}

onBeforeMount(() => {
    findTemplates();
});
</script>

<template>
    <div
        class="divide-y divide-gray-200 overflow-hidden rounded-lg bg-gray-200 shadow sm:grid sm:grid-cols-2 sm:gap-px sm:divide-y-0">
        <div v-for="(template, templateIdx) in templates" :key="template.getId()"
            :class="[templateIdx === 0 ? 'rounded-tl-lg rounded-tr-lg sm:rounded-tr-none' : '', templateIdx === 1 ? 'sm:rounded-tr-lg' : '', templateIdx === templates.length - 2 ? 'sm:rounded-bl-lg' : '', templateIdx === templates.length - 1 ? 'rounded-bl-lg rounded-br-lg sm:rounded-bl-none' : '', 'group relative bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-500']">
            <div class="mt-8">
                <h3 class="text-base font-semibold leading-6 text-gray-900">
                    <button @click="getTemplate(template.getId())" class="focus:outline-none">
                        <!-- Extend touch target to entire panel -->
                        <span class="absolute inset-0" aria-hidden="true" />
                        {{ template.getTitle() }}
                    </button>
                </h3>
                <p class="mt-2 text-sm text-gray-500">{{ template.getDescription() }}</p>
            </div>
            <span class="pointer-events-none absolute top-6 right-6 text-gray-300 group-hover:text-gray-400"
                aria-hidden="true">
                <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
                    <path
                        d="M20 4h1a1 1 0 00-1-1v1zm-1 12a1 1 0 102 0h-2zM8 3a1 1 0 000 2V3zM3.293 19.293a1 1 0 101.414 1.414l-1.414-1.414zM19 4v12h2V4h-2zm1-1H8v2h12V3zm-.707.293l-16 16 1.414 1.414 16-16-1.414-1.414z" />
                </svg>
            </span>
        </div>
    </div>
</template>
