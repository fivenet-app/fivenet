<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import DocumentCategoryBadge from '../partials/documents/DocumentCategoryBadge.vue';
import { relationToBadge } from './helpers';

const props = withDefaults(
    defineProps<{
        documentId: string;
        showDocument?: boolean;
        showSource?: boolean;
    }>(),
    {
        showDocument: true,
        showSource: true,
    },
);

const { t } = useI18n();

const {
    data: relations,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-relations`, () => getDocumentRelations());

async function getDocumentRelations(): Promise<DocumentRelation[]> {
    try {
        const call = getGRPCDocStoreClient().getDocumentRelations({
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
                  key: 'document',
                  label: t('common.document'),
              }
            : undefined,
        {
            key: 'targetUser',
            label: t('common.target'),
        },
        {
            key: 'relation',
            label: t('common.relation', 1),
        },
        props.showSource
            ? {
                  key: 'sourceUser',
                  label: t('common.creator'),
              }
            : undefined,
        {
            key: 'date',
            label: t('common.date'),
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);
</script>

<template>
    <div>
        <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.relation', 2)])" />
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
                <ul role="list" class="divide-y divide-gray-600 overflow-hidden rounded-lg sm:hidden">
                    <li v-for="relation in relations" :key="relation.id" class="block p-4 hover:bg-base-900">
                        <span class="flex items-center space-x-4">
                            <span class="flex flex-1 space-x-2 truncate">
                                <UIcon name="i-mdi-arrow-collapse" class="size-5 shrink-0" />
                                <span class="flex flex-col truncate text-sm">
                                    <span v-if="showDocument">
                                        <ULink
                                            :to="{
                                                name: 'documents-id',
                                                params: {
                                                    id: relation.documentId,
                                                },
                                            }"
                                            class="inline-flex items-center gap-1 truncate"
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
                                        {{ $t(`enums.docstore.DocRelation.${DocRelation[relation.relation]}`) }}
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
                                :loading="loading"
                                :columns="columns"
                                :rows="relations"
                                :empty-state="{
                                    icon: 'i-mdi-account',
                                    label: $t('common.not_found', [$t('common.relation', 2)]),
                                }"
                                sort-mode="auto"
                            >
                                <template v-if="showDocument" #document-data="{ row: relation }">
                                    <ULink
                                        :to="{
                                            name: 'documents-id',
                                            params: {
                                                id: relation.documentId,
                                            },
                                        }"
                                        class="inline-flex items-center gap-1 truncate"
                                    >
                                        <DocumentCategoryBadge :category="relation.document?.category" />

                                        <span>
                                            {{ relation.document?.title }}
                                        </span>
                                    </ULink>
                                </template>
                                <template #targetUser-data="{ row: relation }">
                                    <span class="inline-flex items-center gap-1">
                                        <CitizenInfoPopover :user="relation.targetUser" />
                                        ({{ relation.targetUser?.dateofbirth }})
                                    </span>
                                </template>
                                <template #relation-data="{ row: relation }">
                                    <UBadge :color="relationToBadge(relation.relation)">
                                        {{ $t(`enums.docstore.DocRelation.${DocRelation[relation.relation]}`) }}
                                    </UBadge>
                                </template>
                                <template v-if="showSource" #sourceUser-data="{ row: relation }">
                                    <CitizenInfoPopover :user="relation.sourceUser" />
                                </template>
                                <template #date-data="{ row: relation }">
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
