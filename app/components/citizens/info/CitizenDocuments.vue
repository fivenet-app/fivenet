<script lang="ts" setup>
import { z } from 'zod';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { OpenClose } from '~/typings';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import { ListUserDocumentsResponse } from '~~/gen/ts/services/docstore/docstore';

const { t } = useI18n();

const props = defineProps<{
    userId: number;
}>();

const openclose: OpenClose[] = [
    { id: 0, label: t('common.not_selected') },
    { id: 1, label: t('common.open', 2), closed: false },
    { id: 2, label: t('common.close', 2), closed: true },
];

const schema = z.object({
    closed: z.boolean().optional(),
});

type Schema = z.output<typeof schema>;

const query = ref<Schema>({});

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
            closed: query.value.closed,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), { debounce: 200, maxWait: 1250 });

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
    <div class="flex flex-row flex-wrap gap-2">
        <UFormGroup name="closed" :label="$t('common.close', 2)" class="flex-1">
            <USelectMenu
                v-model="query.closed"
                :options="openclose"
                value-attribute="closed"
                :searchable-placeholder="$t('common.search_field')"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
        </UFormGroup>
    </div>

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
