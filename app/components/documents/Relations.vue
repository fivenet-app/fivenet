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
import { type DocumentRelation, DocRelation } from '~~/gen/ts/resources/documents/documents';
import { docRelationToBadge } from './helpers';

const props = withDefaults(
    defineProps<{
        documentId: number;
        showDocument?: boolean;
        showSource?: boolean;
    }>(),
    {
        showDocument: true,
        showSource: true,
    },
);

const appConfig = useAppConfig();

const { t } = useI18n();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const {
    data: relations,
    status,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-relations`, () => getDocumentRelations());

async function getDocumentRelations(): Promise<DocumentRelation[]> {
    try {
        const call = documentsDocumentsClient.getDocumentRelations({
            documentId: props.documentId,
        });
        const { response } = await call;

        return response.relations;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const columns = computed(() =>
    (
        [
            props.showDocument
                ? {
                      accessorKey: 'document',
                      header: ({ column }) => {
                          const isSorted = column.getIsSorted();

                          return h(UButton, {
                              color: 'neutral',
                              variant: 'ghost',
                              label: t('common.document'),
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
                                      params: { id: row.original.documentId },
                                  },
                              },
                              [
                                  h(CategoryBadge, { category: row.original.document?.category }),
                                  h('span', row.original.document?.title),
                              ],
                          ),
                  }
                : undefined,
            {
                accessorKey: 'targetUser',
                header: t('common.target'),
                cell: ({ row }) =>
                    h('span', { class: 'inline-flex items-center gap-1' }, [
                        h(CitizenInfoPopover, { user: row.original.targetUser }),
                        `(${row.original.targetUser?.dateofbirth})`,
                    ]),
            },
            {
                accessorKey: 'relation',
                header: t('common.relation', 1),
                cell: ({ row }) =>
                    h(
                        UBadge,
                        { color: docRelationToBadge(row.original.relation) },
                        t(`enums.documents.DocRelation.${DocRelation[row.original.relation]}`),
                    ),
            },
            props.showSource
                ? {
                      accessorKey: 'sourceUser',
                      header: t('common.creator'),
                      cell: ({ row }) => h(CitizenInfoPopover, { user: row.original.sourceUser }),
                  }
                : undefined,
            {
                accessorKey: 'date',
                header: t('common.date'),
                cell: ({ row }) => h(GenericTime, { value: row.original.createdAt }),
            },
        ] as TableColumn<DocumentRelation>[]
    ).flatMap((item) => (item !== undefined ? [item] : [])),
);
</script>

<template>
    <div>
        <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.relation', 2)])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.relation', 2)])"
            :error="error"
            :retry="refresh"
        />
        <DataNoDataBlock
            v-if="!relations || relations.length === 0"
            :type="`${$t('common.document', 1)} ${$t('common.relation', 2)}`"
            icon="i-mdi-account-multiple"
        />

        <template v-else>
            <!-- Relations list (smallest breakpoint only) -->
            <div class="sm:hidden">
                <ul class="divide-y divide-default overflow-hidden rounded-lg sm:hidden" role="list">
                    <li v-for="relation in relations" :key="relation.id" class="block p-4 hover:bg-neutral-900">
                        <span class="flex items-center space-x-4">
                            <span class="flex flex-1 space-x-2 truncate">
                                <UIcon class="size-5 shrink-0" name="i-mdi-arrow-collapse" />
                                <span class="flex flex-col truncate text-sm">
                                    <span v-if="showDocument">
                                        <ULink
                                            class="inline-flex items-center gap-1 truncate"
                                            :to="{
                                                name: 'documents-id',
                                                params: {
                                                    id: relation.documentId,
                                                },
                                            }"
                                        >
                                            <CategoryBadge :category="relation.document?.category" />

                                            <span>
                                                {{ relation.document?.title }}
                                            </span>
                                        </ULink>
                                    </span>
                                    <span>
                                        <span class="inline-flex items-center gap-1">
                                            <CitizenInfoPopover :user="relation.targetUser" />
                                            ({{ relation.targetUser?.dateofbirth }})
                                        </span>
                                    </span>
                                    <span class="font-medium">
                                        {{ $t(`enums.documents.DocRelation.${DocRelation[relation.relation]}`) }}
                                    </span>
                                    <span v-if="showSource" class="truncate">
                                        <CitizenInfoPopover :user="relation.sourceUser" />
                                    </span>
                                    <GenericTime :value="relation.createdAt" />
                                </span>
                            </span>
                        </span>
                    </li>
                </ul>
            </div>

            <!-- Relations table (small breakpoint and up) -->
            <div class="hidden sm:block">
                <div>
                    <div class="flex flex-col">
                        <div class="w-full overflow-hidden overflow-x-auto align-middle">
                            <UTable
                                :columns="columns"
                                :data="relations"
                                :loading="isRequestPending(status)"
                                :empty="$t('common.not_found', [$t('common.relation', 2)])"
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
