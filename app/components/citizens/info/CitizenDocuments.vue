<script lang="ts" setup>
import { listEnumValues } from '@protobuf-ts/runtime';
import { z } from 'zod';
import { docRelationToBadge, docRelationToColor, docRelationToIcon } from '~/components/documents/helpers';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { ToggleItem } from '~/utils/types';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import type { ListUserDocumentsResponse } from '~~/gen/ts/services/documents/documents';

const props = defineProps<{
    userId: number;
}>();

const { t } = useI18n();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const openclose: ToggleItem[] = [
    { id: 0, label: t('common.not_selected'), value: undefined },
    { id: 1, label: t('common.open', 2), value: false },
    { id: 2, label: t('common.close', 2), value: true },
];

const docRelationsEnum = listEnumValues(DocRelation).filter((r) => r.number !== 0);
const docRelations = docRelationsEnum.map((r) => ({
    label: t(`enums.documents.DocRelation.${r.name}`),
    value: DocRelation[r.name as keyof typeof DocRelation],
}));

const schema = z.object({
    closed: z.coerce.boolean().optional(),
    relations: z
        .nativeEnum(DocRelation)
        .array()
        .max(docRelations.length)
        .default(docRelationsEnum.map((r) => DocRelation[r.name as keyof typeof DocRelation])),
    sorting: z
        .custom<SortByColumn>()
        .array()
        .max(3)
        .default([
            {
                id: 'createdAt',
                desc: true,
            },
        ]),
    page: pageNumberSchema,
});

const query = useSearchForm('citizen_documents', schema);

const { data, status, refresh, error } = useLazyAsyncData(`citizeninfo-documents-${props.userId}-${query.page}`, () =>
    listUserDocuments(),
);

async function listUserDocuments(): Promise<ListUserDocumentsResponse> {
    try {
        const call = documentsDocumentsClient.listUserDocuments({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: { columns: query.sorting },
            userId: props.userId,
            relations: query.relations,
            closed: query.closed,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watchDebounced(query, async () => refresh(), { debounce: 250, maxWait: 1250 });

const columns = [
    {
        accessorKey: 'document',
        label: t('common.document', 1),
    },
    {
        accessorKey: 'closed',
        label: t('common.close', 2),
    },
    {
        accessorKey: 'relation',
        label: t('common.relation', 1),
    },
    {
        accessorKey: 'createdAt',
        label: t('common.created_at'),
        sortable: true,
    },
    {
        accessorKey: 'creator',
        label: t('common.creator'),
    },
];
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm class="flex w-full flex-row gap-2" :state="query" :schema="schema">
                <UFormField class="flex-1" name="closed" :label="$t('common.close', 2)">
                    <ClientOnly>
                        <USelectMenu
                            v-model="query.closed"
                            :items="openclose"
                            value-key="closed"
                            :searchable-placeholder="$t('common.search_field')"
                        >
                            <template #item-label>
                                {{
                                    query.closed === undefined
                                        ? openclose[0]!.label
                                        : (openclose.findLast((o) => o.value === query.closed)?.label ?? $t('common.na'))
                                }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UFormField class="flex-1" name="relation" :label="$t('common.relation')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="query.relations"
                            multiple
                            :items="docRelations"
                            value-key="value"
                            :searchable-placeholder="$t('common.relation', 2)"
                        >
                            <template #item-label>
                                {{ $t('common.selected', query.relations.length) }}
                            </template>

                            <template #option="{ option }">
                                <span class="inline-flex gap-1" :class="`bg-${docRelationToColor(option.value)}-500`">
                                    <UIcon class="size-4" :name="docRelationToIcon(option.value)" />
                                    <span class="truncate">
                                        {{ $t(`enums.documents.DocRelation.${DocRelation[option.value]}`) }}
                                    </span>
                                </span>
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>
            </UForm>
        </template>
    </UDashboardToolbar>

    <DataErrorBlock
        v-if="error"
        :title="$t('common.unable_to_load', [`${$t('common.citizen', 1)} ${$t('common.document', 2)}`])"
        :error="error"
        :retry="refresh"
    />
    <UTable
        v-else
        class="flex-1"
        :loading="isRequestPending(status)"
        :columns="columns"
        :data="data?.relations"
        :empty-state="{
            icon: 'i-mdi-file-multiple',
            label: $t('common.not_found', [`${$t('common.citizen', 1)} ${$t('common.document', 2)}`]),
        }"
    >
        <template #document-cell="{ row: relation }">
            <DocumentInfoPopover :document="relation.document" load-on-open />
        </template>

        <template #closed-cell="{ row: relation }">
            <OpenClosedBadge :closed="relation.document?.closed" variant="subtle" />
        </template>

        <template #relation-cell="{ row: relation }">
            <UBadge
                class="font-semibold"
                :color="docRelationToBadge(relation.relation)"
                :icon="docRelationToIcon(relation.relation)"
            >
                {{ $t(`enums.documents.DocRelation.${DocRelation[relation.relation]}`) }}
            </UBadge>
        </template>

        <template #createdAt-cell="{ row: relation }">
            <GenericTime :value="relation.createdAt" />
        </template>

        <template #creator-cell="{ row: relation }">
            <CitizenInfoPopover :user="relation.sourceUser" />
        </template>
    </UTable>

    <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
</template>
