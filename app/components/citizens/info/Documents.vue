<script lang="ts" setup>
import { UBadge } from '#components';
import type { Form, TableColumn } from '@nuxt/ui';
import { listEnumValues } from '@protobuf-ts/runtime';
import { computed, h } from 'vue';
import { z } from 'zod';
import { docRelationToBadge, docRelationToIcon } from '~/components/documents/helpers';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import TableSortButton from '~/components/partials/TableSortButton.vue';
import type { ToggleItem } from '~/utils/types';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { DocRelation, type DocumentRelation } from '~~/gen/ts/resources/documents/relations/relations';
import type { ListUserDocumentsResponse } from '~~/gen/ts/services/documents/documents';

const props = defineProps<{
    userId: number;
}>();

const { t } = useI18n();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const openclose: ToggleItem<boolean | undefined>[] = [
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
        .enum(DocRelation)
        .array()
        .max(docRelations.length)
        .default(docRelationsEnum.map((r) => DocRelation[r.name as keyof typeof DocRelation])),
    includeCreated: z.coerce.boolean().optional().default(true),

    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'createdAt',
                        desc: true,
                    },
                ]),
        })
        .default({ columns: [{ id: 'createdAt', desc: true }] }),
    page: pageNumberSchema,
});

type Schema = z.output<typeof schema>;

const query = useSearchForm('citizen_documents', schema);

const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { validatedQuery, commitValidatedQuery } = useFormSearchValidation<typeof schema>(query, formRef);

const documentsKey = computed(() => `citizeninfo-documents-${props.userId}-${JSON.stringify(validatedQuery.value)}`);

const { data, status, refresh, error } = useLazyAsyncData(documentsKey, () => listUserDocuments(validatedQuery.value));

async function listUserDocuments(values: Schema): Promise<ListUserDocumentsResponse> {
    try {
        const call = documentsDocumentsClient.listUserDocuments({
            pagination: {
                offset: calculateOffset(values.page, data.value?.pagination),
            },
            sort: values.sorting,
            userId: props.userId,
            relations: values.relations,
            closed: values.closed,
            includeCreated: values.includeCreated,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const columns = computed(
    () =>
        [
            {
                accessorKey: 'document',
                header: t('common.document', 1),
                cell: ({ row }) => h(DocumentInfoPopover, { document: row.original.document, loadOnOpen: true, showId: true }),
            },
            {
                accessorKey: 'closed',
                header: t('common.close', 2),
                cell: ({ row }) => h(OpenClosedBadge, { closed: row.original.document?.meta?.closed, variant: 'subtle' }),
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
                    return h(TableSortButton, {
                        column,
                        label: t('common.created_at'),
                    });
                },
                cell: ({ row }) => h(GenericTime, { value: row.original.createdAt }),
            },
            {
                accessorKey: 'creator',
                header: t('common.creator'),
                cell: ({ row }) =>
                    row.original.sourceUser ? h(CitizenInfoPopover, { user: row.original.sourceUser }) : undefined,
            },
        ] as TableColumn<DocumentRelation>[],
);
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm
                ref="formRef"
                class="my-2 flex w-full flex-row gap-2"
                :state="query"
                :schema="schema"
                @submit="commitValidatedQuery"
            >
                <UFormField class="flex-1" name="closed" :label="$t('common.close', 2)">
                    <ClientOnly>
                        <USelectMenu
                            v-model="query.closed"
                            class="w-full"
                            :items="openclose"
                            value-key="value"
                            :search-input="{ placeholder: $t('common.search_field') }"
                        >
                            <template #default>
                                <div class="inline-flex items-center gap-1 truncate">
                                    <template v-if="typeof query.closed === 'boolean'">
                                        <UIcon
                                            v-if="!query.closed"
                                            class="size-4"
                                            name="i-mdi-lock-open-variant"
                                            color="success"
                                        />
                                        <UIcon v-else class="size-4" name="i-mdi-lock" color="error" />
                                    </template>

                                    {{
                                        query.closed === undefined
                                            ? openclose[0]!.label
                                            : (openclose.findLast((o) => o.value === query.closed)?.label ?? $t('common.na'))
                                    }}
                                </div>
                            </template>

                            <template #item-label="{ item }">
                                <div class="inline-flex items-center gap-1 truncate">
                                    <template v-if="typeof item.value === 'boolean'">
                                        <UIcon
                                            v-if="!item.value"
                                            class="size-4"
                                            name="i-mdi-lock-open-variant"
                                            color="success"
                                        />
                                        <UIcon v-else class="size-4" name="i-mdi-lock" color="error" />
                                    </template>

                                    {{ item.label }}
                                </div>
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UFormField class="flex-1" name="relation" :label="$t('common.relation')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="query.relations"
                            class="w-full"
                            multiple
                            :items="docRelations"
                            value-key="value"
                            :search-input="{ placeholder: $t('common.relation', 2) }"
                        >
                            <template #default>
                                {{ $t('common.selected', query.relations.length) }}
                            </template>

                            <template #item-label="{ item }">
                                <UBadge
                                    class="truncate"
                                    :color="docRelationToBadge(item.value)"
                                    :icon="docRelationToIcon(item.value)"
                                    :label="$t(`enums.documents.DocRelation.${DocRelation[item.value]}`)"
                                />
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UFormField name="includeCreated" :label="$t('components.citizens.documents.include_created.title')">
                    <div class="inline-flex items-center gap-2">
                        <USwitch v-model="query.includeCreated" />

                        <UTooltip :text="$t('components.citizens.documents.include_created.description')">
                            <UIcon class="size-4" name="i-mdi-information" />
                        </UTooltip>
                    </div>
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
        v-model:sorting="query.sorting.columns"
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
