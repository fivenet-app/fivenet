<script lang="ts" setup>
import { DocumentTemplateShort } from '@fivenet/gen/resources/documents/templates_pb';
import { ListTemplatesRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { ref, onBeforeMount } from 'vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/20/solid';
import { ArrowUpRightIcon } from '@heroicons/vue/24/solid';
import { RpcError } from 'grpc-web';

const { $grpc } = useNuxtApp();

const templates = ref<Array<DocumentTemplateShort>>([]);

defineEmits<{
    (e: 'selected', t: DocumentTemplateShort): void,
}>();

async function findTemplates(): Promise<void> {
    const req = new ListTemplatesRequest();

    try {
        const resp = await $grpc.getDocStoreClient().
            listTemplates(req, null);

        templates.value = resp.getTemplatesList();
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
        return;
    }
}

onBeforeMount(async () => {
    findTemplates();
});
</script>

<template>
    <div>
        <button v-if="templates.length == 0" type="button"
            class="relative block w-full p-12 text-center rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400">
            <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-base-200">
                No templates for your job and rank found.
            </span>
        </button>
        <div v-else
            class="overflow-hidden rounded-lg divide-y-4 bg-base-900 sm:grid sm:grid-cols-2 sm:gap-1 sm:max-w-6xl sm:mx-auto divide-base-900 sm:divide-y-0">
            <div v-for="(template, templateIdx) in templates" :key="template.getId()" :class="[
                templateIdx === 0 ? 'rounded-tl-lg rounded-tr-lg sm:rounded-tr-none' : '',
                templateIdx === 1 ? 'sm:rounded-tr-lg' : '',
                templateIdx === templates.length - 2 && templates.length % 2 === 0 ? 'sm:rounded-bl-lg' : '',
                templateIdx === templates.length - 1 && templates.length % 2 === 0 ? 'rounded-br-lg' : '',
                templateIdx === templates.length - 1 ? 'rounded-bl-lg sm:rounded-bl-none' : '',
                'group relative bg-base-700 p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-neutral',
            ]">
                <div class="mt-8">
                    <h3 class="text-base font-semibold leading-6 text-neutral">
                        <button @click="$emit('selected', template)" class="focus:outline-none">
                            <!-- Extend touch target to entire panel -->
                            <span class="absolute inset-0" aria-hidden="true" />
                            {{ template.getTitle() }}
                        </button>
                    </h3>
                    <p class="mt-2 text-sm text-base-200">
                        {{ template.getDescription() }}
                    </p>
                </div>
                <span class="absolute pointer-events-none top-6 right-6 text-base-300 group-hover:text-base-200"
                    aria-hidden="true">
                    <ArrowUpRightIcon class="w-6 h-6" />
                </span>
            </div>
        </div>
    </div>
</template>
