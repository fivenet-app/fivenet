<script lang="ts" setup>
import { DocumentTemplateShort } from '@fivenet/gen/resources/documents/templates_pb';
import { ListTemplatesRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { ref, onBeforeMount, FunctionalComponent } from 'vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/20/solid';
import { ArrowUpRightIcon } from '@heroicons/vue/24/solid';
import { RpcError } from 'grpc-web';
import Cards from '~/components/partials/Cards.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';
import { CardElements } from '~~/src/utils/types';

const { $grpc } = useNuxtApp();

defineEmits<{
    (e: 'selected', t: DocumentTemplateShort): void,
}>();

const { data: templates, pending, refresh, error } = await useLazyAsyncData(`documents-templates`, () => findTemplates());
const items = ref<CardElements>([]);

async function findTemplates(): Promise<Array<DocumentTemplateShort>> {
    return new Promise(async (res, rej) => {
        const req = new ListTemplatesRequest();

        try {
            const resp = await $grpc.getDocStoreClient().
                listTemplates(req, null);

            return res(resp.getTemplatesList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watch(templates, () => templates.value?.forEach((v) => {
    items.value.push({ title: v?.getTitle(), description: v?.getDescription(), });
}));

function selected(idx: number): DocumentTemplateShort {
    return templates.value![idx];
}
</script>

<template>
    <div>
        <DataPendingBlock v-if="pending" message="Loading templates..." />
        <DataErrorBlock v-else-if="error" title="Unable to load templates!" :retry="refresh" />
        <button v-else-if="templates && templates.length == 0" type="button"
            class="relative block w-full p-12 text-center rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400">
            <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-base-200">
                No templates for your job and rank found.
            </span>
        </button>
        <div v-else>
            <Cards :items="items" :show-icon="false" @selected="$emit('selected', selected($event))" />
        </div>
    </div>
</template>
