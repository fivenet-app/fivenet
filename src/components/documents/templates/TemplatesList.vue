<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import Cards from '~/components/partials/Cards.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import { CardElements } from '~/utils/types';
import { TemplateShort } from '~~/gen/ts/resources/documents/templates';

const { $grpc } = useNuxtApp();

defineEmits<{
    (e: 'selected', t: TemplateShort): void;
}>();

const { data: templates, pending, refresh, error } = useLazyAsyncData(`documents-templates`, () => listTemplates());

async function listTemplates(): Promise<Array<TemplateShort>> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getDocStoreClient().listTemplates({});
            const { response } = await call;

            return res(response.templates);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const items = ref<CardElements>([]);
watch(templates, () =>
    templates.value?.forEach((v) => {
        items.value.push({ title: v?.title, description: v?.description });
    })
);

function selected(idx: number): TemplateShort {
    return templates.value![idx];
}
</script>

<template>
    <div>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.template', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.template', 2)])" :retry="refresh" />
        <DataNoDataBlock v-else-if="templates && templates.length === 0" :type="$t('common.template', 2)" />
        <div v-else>
            <Cards :items="items" :show-icon="false" @selected="$emit('selected', selected($event))" />
        </div>
    </div>
</template>
