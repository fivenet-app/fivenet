<script lang="ts" setup>
import { UBadge, UButton, ULink } from '#components';
import type { TableColumn } from '@nuxt/ui';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import CategoryBadge from '~/components/partials/documents/CategoryBadge.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { type DocumentReference, DocReference } from '~~/gen/ts/resources/documents/documents';
import { docReferenceToBadge } from './helpers';

const appConfig = useAppConfig();

const props = withDefaults(
    defineProps<{
        documentId: number;
        showSource?: boolean;
    }>(),
    {
        showSource: true,
    },
);

const { t } = useI18n();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const {
    data: references,
    status,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-references`, () => getDocumentReferences());

async function getDocumentReferences(): Promise<DocumentReference[]> {
    try {
        const call = documentsDocumentsClient.getDocumentReferences({
            documentId: props.documentId,
        });
        const { response } = await call;

        return response.references;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const columns = computed(() =>
    (
        [
            {
                accessorKey: 'targetDocument',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.target'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                cell: ({ row }) =>
                    h(
                        ULink,
                        {
                            class: 'inline-flex items-center gap-1 truncate',
                            to: {
                                name: 'documents-id',
                                params: { id: row.original.targetDocumentId },
                            },
                        },
                        [
                            h(CategoryBadge, { category: row.original.targetDocument?.category }),
                            h('span', row.original.targetDocument?.title),
                        ],
                    ),
            },
            {
                accessorKey: 'reference',
                header: t('common.reference', 1),
                cell: ({ row }) =>
                    h(
                        UBadge,
                        { color: docReferenceToBadge(row.original.reference) },
                        t(`enums.documents.DocReference.${DocReference[row.original.reference]}`),
                    ),
            },
            props.showSource
                ? {
                      accessorKey: 'sourceDocument',
                      header: t('common.source'),
                      cell: ({ row }) =>
                          h(
                              ULink,
                              {
                                  class: 'inline-flex items-center gap-1 truncate',
                                  to: {
                                      name: 'documents-id',
                                      params: { id: row.original.sourceDocumentId },
                                  },
                              },
                              [
                                  h(CategoryBadge, { category: row.original.sourceDocument?.category }),
                                  h('span', row.original.sourceDocument?.title),
                              ],
                          ),
                  }
                : undefined,
            {
                accessorKey: 'creator',
                header: t('common.creator'),
                cell: ({ row }) => h(CitizenInfoPopover, { user: row.original.creator }),
            },
            {
                accessorKey: 'createdAt',
                header: t('common.date'),
                cell: ({ row }) => h(GenericTime, { value: row.original.createdAt }),
            },
        ] as TableColumn<DocumentReference>[]
    ).flatMap((item) => (item !== undefined ? [item] : [])),
);
</script>

<template>
    <div>
        <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.reference', 2)])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.reference', 2)])"
            :error="error"
            :retry="refresh"
        />
        <DataNoDataBlock
            v-else-if="!references || references.length === 0"
            :type="`${$t('common.document', 1)} ${$t('common.reference', 2)}`"
            icon="i-mdi-file-document-multiple"
        />

        <template v-else>
            <!-- Relations list (smallest breakpoint only) -->
            <div class="sm:hidden">
                <ul class="divide-y divide-default rounded-lg sm:hidden" role="list">
                    <li v-for="reference in references" :key="reference.id">
                        <ULink
                            class="block p-4 hover:bg-neutral-900"
                            :to="{
                                name: 'documents-id',
                                params: { id: reference.targetDocumentId },
                            }"
                        >
                            <span class="flex items-center space-x-4">
                                <span class="flex flex-1 space-x-2 truncate">
                                    <UIcon class="size-5 shrink-0" name="i-mdi-arrow-collapse" />
                                    <span class="flex flex-col truncate text-sm">
                                        <span>
                                            {{ reference.targetDocument?.title
                                            }}<span v-if="reference.targetDocument?.category"
                                                >&nbsp;({{ $t('common.category', 1) }}:
                                                {{ reference.targetDocument?.category?.name }})</span
                                            >
                                        </span>
                                        <span class="font-medium">
                                            {{ $t(`enums.documents.DocReference.${DocReference[reference.reference]}`) }}
                                        </span>
                                        <span v-if="showSource" class="truncate">
                                            {{ reference.sourceDocument?.title
                                            }}<span v-if="reference.sourceDocument?.category">
                                                ({{ $t('common.category', 1) }}:
                                                {{ reference.sourceDocument?.category?.name }})</span
                                            >
                                        </span>
                                        <span>
                                            <CitizenInfoPopover :user="reference.sourceDocument?.creator" />
                                        </span>
                                        <GenericTime :value="reference.createdAt" />
                                    </span>
                                </span>
                            </span>
                        </ULink>
                    </li>
                </ul>
            </div>

            <!-- Relations table (small breakpoint and up) -->
            <div class="hidden sm:block">
                <div>
                    <div class="flex flex-col">
                        <div class="min-w-full overflow-hidden overflow-x-auto align-middle sm:rounded-lg">
                            <UTable
                                :columns="columns"
                                :data="references"
                                :loading="isRequestPending(status)"
                                :empty="$t('common.not_found', [$t('common.reference', 2)])"
                                :pagination-options="{ manualPagination: true }"
                                :sorting-options="{ manualSorting: true }"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </div>
</template>
