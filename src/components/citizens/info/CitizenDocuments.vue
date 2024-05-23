<script lang="ts" setup>
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import { ListUserDocumentsResponse } from '~~/gen/ts/services/docstore/docstore';

const { t } = useI18n();

const props = defineProps<{
    userId: number;
}>();

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`citizeninfo-documents-${props.userId}-${page.value}`, () => listUserDocuments());

async function listUserDocuments(): Promise<ListUserDocumentsResponse> {
    try {
        const call = getGRPCDocStoreClient().listUserDocuments({
            pagination: {
                offset: offset.value,
            },
            userId: props.userId,
            relations: [],
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());

const columns = [
    {
        key: 'document',
        label: t('common.document', 1),
    },
    {
        key: 'closed',
        label: t('common.close', 2),
    },
    {
        key: 'relation',
        label: t('common.relation', 1),
    },
    {
        key: 'date',
        label: t('common.date'),
    },
    {
        key: 'creator',
        label: t('common.creator'),
    },
];
</script>

<template>
    <DataErrorBlock
        v-if="error"
        :title="$t('common.unable_to_load', [`${$t('common.citizen', 1)} ${$t('common.document', 2)}`])"
        :retry="refresh"
    />
    <UTable
        v-else
        :loading="loading"
        :columns="columns"
        :rows="data?.relations"
        :empty-state="{
            icon: 'i-mdi-file-multiple',
            label: $t('common.not_found', [`${$t('common.citizen', 1)} ${$t('common.document', 2)}`]),
        }"
        class="flex-1"
    >
        <template #document-data="{ row: relation }">
            <DocumentInfoPopover :document="relation.document" />
        </template>
        <template #closed-data="{ row: relation }">
            <OpenClosedBadge :closed="relation.document?.closed" variant="subtle" />
        </template>
        <template #relation-data="{ row: relation }">
            <span class="font-medium">
                {{ $t(`enums.docstore.DocRelation.${DocRelation[relation.relation]}`) }}
            </span>
        </template>
        <template #date-data="{ row: relation }">
            <GenericTime :value="relation.createdAt" />
        </template>
        <template #creator-data="{ row: relation }">
            <CitizenInfoPopover :user="relation.sourceUser" />
        </template>
    </UTable>

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
