<script lang="ts" setup>
import CardsList from '~/components/partials/CardsList.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import type { CardElements } from '~/utils/types';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import type { TemplateShort } from '~~/gen/ts/resources/documents/templates';

const props = defineProps<{
    link?: boolean;
}>();

defineEmits<{
    (e: 'selected', t: TemplateShort | undefined): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const { data: templates, status, refresh, error } = useLazyAsyncData(`documents-templates`, () => listTemplates());

defineExpose({
    status,
    refresh,
});

const documentsDocumentsClient = await getDocumentsDocumentsClient();

async function listTemplates(): Promise<TemplateShort[]> {
    try {
        const call = documentsDocumentsClient.listTemplates({});
        const { response } = await call;

        return response.templates;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const items = computed<CardElements>(
    () =>
        templates.value?.map((v) => ({
            title: v?.title,
            description: v?.description,
            icon: v.icon ?? 'i-mdi-file-outline',
            color: v.color ?? 'primary',
            to: props.link ? { name: 'documents-templates-id', params: { id: v?.id } } : undefined,
        })) ?? [],
);

function selected(idx: number): TemplateShort | undefined {
    return templates.value?.at(idx);
}
</script>

<template>
    <div v-if="isRequestPending(status)" class="flex justify-center">
        <UPageGrid>
            <UPageCard v-for="idx in 6" :key="idx">
                <template #title>
                    <USkeleton class="h-6 w-[275px]" />
                </template>
                <template #description>
                    <USkeleton class="h-11 w-[350px]" />
                </template>
            </UPageCard>
        </UPageGrid>
    </div>

    <DataErrorBlock
        v-else-if="error"
        :title="$t('common.unable_to_load', [$t('common.template', 2)])"
        :error="error"
        :retry="refresh"
    />
    <DataNoDataBlock v-else-if="!templates || templates.length === 0" :type="$t('common.template', 2)" />

    <div v-else class="flex justify-center">
        <CardsList :class="$attrs.class" :items="items" @selected="$emit('selected', selected($event))" />
    </div>
</template>
