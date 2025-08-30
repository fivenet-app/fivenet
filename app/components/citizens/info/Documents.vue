<script lang="ts" setup>
import { UBadge, UButton } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { listEnumValues } from '@protobuf-ts/runtime';
import { computed, h } from 'vue';
import { z } from 'zod';
import { docRelationToBadge, docRelationToIcon } from '~/components/documents/helpers';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import type { ToggleItem } from '~/utils/types';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { DocRelation, type DocumentRelation } from '~~/gen/ts/resources/documents/documents';
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
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'plate',
                        desc: false,
                    },
                ]),
        })
        .default({ columns: [{ id: 'plate', desc: false }] }),
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
            sort: query.sorting,
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

const appConfig = useAppConfig();

const columns = computed(
    () =>
        [
            {
                accessorKey: 'document',
                header: t('common.document', 1),
                cell: ({ row }) => h(DocumentInfoPopover, { document: row.original.document, loadOnOpen: true }),
            },
            {
                accessorKey: 'closed',
                header: t('common.close', 2),
                cell: ({ row }) => h(OpenClosedBadge, { closed: row.original.document?.closed, variant: 'subtle' }),
            },
            {
                accessorKey: 'relation',
                header: t('common.relation', 1),
                cell: ({ row }) =>
                    h(
                        UBadge,
                        {
                            class: 'font-semibold',
                            color: docRelationToBadge(row.original.relation),
                            icon: docRelationToIcon(row.original.relation),
                        },
                        () => t(`enums.documents.DocRelation.${DocRelation[row.original.relation]}`),
                    ),
            },
            {
                accessorKey: 'createdAt',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.created_at'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(isSorted === 'asc'),
                    });
                },
            },
            {
                accessorKey: 'creator',
                header: t('common.creator'),
                cell: ({ row }) => h(CitizenInfoPopover, { user: row.original.sourceUser }),
            },
        ] as TableColumn<DocumentRelation>[],
);
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
                            value-key="value"
                            :search-input="{ placeholder: $t('common.search_field') }"
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
                            :search-input="{ placeholder: $t('common.relation', 2) }"
                        >
                            <template #item-label>
                                {{ $t('common.selected', query.relations.length) }}
                            </template>

                            <template #item="{ item }">
                                <UBadge
                                    class="truncate"
                                    :color="docRelationToBadge(item.value)"
                                    :icon="docRelationToIcon(item.value)"
                                >
                                    {{ $t(`enums.documents.DocRelation.${DocRelation[item.value]}`) }}
                                </UBadge>
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
        :empty="$t('common.not_found', [`${$t('common.citizen', 1)} ${$t('common.document', 2)}`])"
        :pagination-options="{ manualPagination: true }"
        :sorting-options="{ manualSorting: true }"
        sticky
    />

    <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
</template>
