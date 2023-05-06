<script lang="ts" setup>
import { TemplateShort } from '@fivenet/gen/resources/documents/templates_pb';
import { ListTemplatesRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { MagnifyingGlassIcon } from '@heroicons/vue/20/solid';
import { RpcError } from 'grpc-web';
import Cards from '~/components/partials/Cards.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import { CardElements } from '~/utils/types';

const { $grpc } = useNuxtApp();

defineEmits<{
    (e: 'selected', t: TemplateShort): void,
}>();

const { data: templates, pending, refresh, error } = useLazyAsyncData(`documents-templates`, () => listTemplates());

async function listTemplates(): Promise<Array<TemplateShort>> {
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

const items = ref<CardElements>([]);
watch(templates, () => templates.value?.forEach((v) => {
    items.value.push({ title: v?.getTitle(), description: v?.getDescription(), });
}));

function selected(idx: number): TemplateShort {
    return templates.value![idx];
}
</script>

<template>
    <div>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.template', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.template', 2)])"
            :retry="refresh" />
        <button v-else-if="templates && templates.length == 0" type="button"
            class="relative block w-full p-12 text-center rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400">
            <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-base-200">
                {{ $t('common.not_found', [$t('common.template', 2)]) }}
            </span>
        </button>
        <div v-else>
            <Cards :items="items" :show-icon="false" @selected="$emit('selected', selected($event))" />
        </div>
    </div>
</template>
