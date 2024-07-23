<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { DocReference, DocumentReference } from '~~/gen/ts/resources/documents/documents';
import { refToBadge } from './helpers';

const props = withDefaults(
    defineProps<{
        documentId: string;
        showSource?: boolean;
    }>(),
    {
        showSource: true,
    },
);

const { t } = useI18n();

const {
    data: references,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-references`, () => getDocumentReferences());

async function getDocumentReferences(): Promise<DocumentReference[]> {
    try {
        const call = getGRPCDocStoreClient().getDocumentReferences({
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
    [
        {
            key: 'targetDocument',
            label: t('common.target'),
        },
        {
            key: 'reference',
            label: t('common.reference', 1),
        },
        props.showSource
            ? {
                  key: 'sourceDocument',
                  label: t('common.source'),
              }
            : undefined,
        {
            key: 'creator',
            label: t('common.creator'),
        },
        {
            key: 'date',
            label: t('common.date'),
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);
</script>

<template>
    <div>
        <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.reference', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.reference', 2)])" :retry="refresh" />
        <DataNoDataBlock
            v-else-if="!references || references.length === 0"
            :type="`${$t('common.document', 1)} ${$t('common.reference', 2)}`"
            icon="i-mdi-file-document-multiple"
        />

        <template v-else>
            <!-- Relations list (smallest breakpoint only) -->
            <div class="sm:hidden">
                <ul role="list" class="divide-y divide-gray-600 rounded-lg sm:hidden">
                    <li v-for="reference in references" :key="reference.id">
                        <ULink
                            :to="{
                                name: 'documents-id',
                                params: { id: reference.targetDocumentId },
                            }"
                            class="block p-4 hover:bg-base-900"
                        >
                            <span class="flex items-center space-x-4">
                                <span class="flex flex-1 space-x-2 truncate">
                                    <UIcon name="i-mdi-arrow-collapse" class="size-5 shrink-0" />
                                    <span class="flex flex-col truncate text-sm">
                                        <span>
                                            {{ reference.targetDocument?.title
                                            }}<span v-if="reference.targetDocument?.category"
                                                >&nbsp;({{ $t('common.category', 1) }}:
                                                {{ reference.targetDocument?.category?.name }})</span
                                            >
                                        </span>
                                        <span class="font-medium">
                                            {{ $t(`enums.docstore.DocReference.${DocReference[reference.reference]}`) }}
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
                                :loading="loading"
                                :columns="columns"
                                :rows="references"
                                :empty-state="{
                                    icon: 'i-mdi-account',
                                    label: $t('common.not_found', [$t('common.reference', 2)]),
                                }"
                            >
                                <template #targetDocument-data="{ row: reference }">
                                    <ULink
                                        :to="{
                                            name: 'documents-id',
                                            params: {
                                                id: reference.targetDocumentId,
                                            },
                                        }"
                                        class="inline-flex items-center gap-1 truncate"
                                    >
                                        <UBadge v-if="reference.targetDocument?.category" class="inline-flex gap-1" size="md">
                                            <UIcon name="i-mdi-shape" class="size-5" />
                                            <span :title="reference.targetDocument?.category.description ?? $t('common.na')">
                                                {{ reference.targetDocument?.category.name }}
                                            </span>
                                        </UBadge>

                                        <span>
                                            {{ reference.targetDocument?.title }}
                                        </span>
                                    </ULink>
                                </template>
                                <template #reference-data="{ row: reference }">
                                    <UBadge :color="refToBadge(reference.reference)">
                                        {{ $t(`enums.docstore.DocReference.${DocReference[reference.reference]}`) }}
                                    </UBadge>
                                </template>
                                <template v-if="showSource" #sourceDocument-data="{ row: reference }">
                                    <ULink
                                        :to="{
                                            name: 'documents-id',
                                            params: {
                                                id: reference.sourceDocumentId,
                                            },
                                        }"
                                        class="inline-flex items-center gap-1 truncate"
                                    >
                                        <UBadge v-if="reference.sourceDocument?.category" class="inline-flex gap-1" size="md">
                                            <UIcon name="i-mdi-shape" class="size-5" />
                                            <span :title="reference.sourceDocument?.category.description ?? $t('common.na')">
                                                {{ reference.sourceDocument?.category.name }}
                                            </span>
                                        </UBadge>

                                        <span>
                                            {{ reference.sourceDocument?.title }}
                                        </span>
                                    </ULink>
                                </template>
                                <template #creator-data="{ row: reference }">
                                    <CitizenInfoPopover :user="reference.creator" />
                                </template>
                                <template #date-data="{ row: reference }">
                                    <GenericTime :value="reference.createdAt" />
                                </template>
                            </UTable>
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </div>
</template>
