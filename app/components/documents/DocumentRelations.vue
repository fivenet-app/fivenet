<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DocumentCategoryBadge from '~/components/partials/documents/DocumentCategoryBadge.vue';
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
    [
        props.showDocument
            ? {
                  accessorKey: 'document',
                  label: t('common.document'),
              }
            : undefined,
        {
            accessorKey: 'targetUser',
            label: t('common.target'),
        },
        {
            accessorKey: 'relation',
            label: t('common.relation', 1),
        },
        props.showSource
            ? {
                  accessorKey: 'sourceUser',
                  label: t('common.creator'),
              }
            : undefined,
        {
            accessorKey: 'date',
            label: t('common.date'),
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
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
                <ul class="divide-y divide-gray-600 overflow-hidden rounded-lg sm:hidden" role="list">
                    <li v-for="relation in relations" :key="relation.id" class="hover:bg-base-900 block p-4">
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
                                            <DocumentCategoryBadge :category="relation.document?.category" />

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
                                :loading="isRequestPending(status)"
                                :columns="columns"
                                :data="relations"
                                :empty-state="{
                                    icon: 'i-mdi-account',
                                    label: $t('common.not_found', [$t('common.relation', 2)]),
                                }"
                                sort-mode="auto"
                            >
                                <template v-if="showDocument" #document-cell="{ row: relation }">
                                    <ULink
                                        class="inline-flex items-center gap-1 truncate"
                                        :to="{
                                            name: 'documents-id',
                                            params: {
                                                id: relation.documentId,
                                            },
                                        }"
                                    >
                                        <DocumentCategoryBadge :category="relation.document?.category" />

                                        <span>
                                            {{ relation.document?.title }}
                                        </span>
                                    </ULink>
                                </template>
                                <template #targetUser-cell="{ row: relation }">
                                    <span class="inline-flex items-center gap-1">
                                        <CitizenInfoPopover :user="relation.targetUser" />
                                        ({{ relation.targetUser?.dateofbirth }})
                                    </span>
                                </template>
                                <template #relation-cell="{ row: relation }">
                                    <UBadge :color="docRelationToBadge(relation.relation)">
                                        {{ $t(`enums.documents.DocRelation.${DocRelation[relation.relation]}`) }}
                                    </UBadge>
                                </template>
                                <template v-if="showSource" #sourceUser-cell="{ row: relation }">
                                    <CitizenInfoPopover :user="relation.sourceUser" />
                                </template>
                                <template #date-cell="{ row: relation }">
                                    <GenericTime :value="relation.createdAt" />
                                </template>
                            </UTable>
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </div>
</template>
