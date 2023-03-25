<script lang="ts" setup>
import { DocumentTemplate, DocumentTemplateShort } from '@arpanet/gen/resources/documents/documents_pb';
import { GetTemplateRequest, ListTemplatesRequest } from '@arpanet/gen/services/docstore/docstore_pb';
import { ref, onBeforeMount, computed, toRaw } from 'vue';
import { getDocStoreClient } from '../../grpc/grpc';
import { MagnifyingGlassIcon } from '@heroicons/vue/20/solid';
import { useStore } from '../../store/store';
import { TemplateData } from '@arpanet/gen/resources/documents/templates/templates_pb';
import { User } from '@arpanet/gen/resources/users/users_pb';

const store = useStore();

const templates = ref<Array<DocumentTemplateShort>>([]);
const templateObj = ref<undefined | DocumentTemplate>(undefined);

const activeChar = computed(() => store.state.auth?.activeChar);

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
    req.setRender(true);

    const data = store.getters['clipboard/getTemplateData'] as TemplateData;
    data.setActivechar(activeChar.value!);
    if (data.getUsersList().length == 0) {
        data.setUsersList([activeChar.value!]);
    }
    req.setData(JSON.stringify(data.toObject()));

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
    <div>
        <button v-if="templates.length == 0" type="button"
            class="relative block w-full rounded-lg border-2 border-dashed border-base-300 p-12 text-center hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2">
            <MagnifyingGlassIcon class="text-neutral mx-auto h-12 w-12" />
            <span class="mt-2 block text-sm font-semibold text-gray-300">No templates for your job and rank
                found.</span>
        </button>
        <div v-else
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
        <hr>
        <div class="text-2xl text-white">
            <h2 class="text-3xl">{{ templateObj?.getContentTitle() }}</h2>
            <hr>
            <p>
                {{ templateObj?.getContent() }}
            </p>
        </div>
    </div>
</template>
