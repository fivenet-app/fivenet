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
import type { ToggleItem } from '~/typings';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import type { ListUserDocumentsResponse } from '~~/gen/ts/services/documents/documents';

const props = defineProps<{
    userId: number;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

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
    sort: z.custom<TableSortable>().default({
        column: 'createdAt',
        direction: 'desc',
    }),
    page: z.coerce.number().min(1).default(1),
});

const query = useSearchForm('citizen_documents', schema);

const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (query.page - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`citizeninfo-documents-${props.userId}-${query.page}`, () => listUserDocuments());

async function listUserDocuments(): Promise<ListUserDocumentsResponse> {
    try {
        const call = $grpc.documents.documents.listUserDocuments({
            pagination: {
                offset: offset.value,
            },
            sort: query.sort,
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

watch(offset, async () => refresh());
watchDebounced(query, async () => refresh(), { debounce: 250, maxWait: 1250 });

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
        key: 'createdAt',
        label: t('common.created_at'),
        sortable: true,
    },
    {
        key: 'creator',
        label: t('common.creator'),
    },
];
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm class="flex w-full flex-row gap-2" :state="query" :schema="schema">
                <UFormGroup class="flex-1" name="closed" :label="$t('common.close', 2)">
                    <ClientOnly>
                        <USelectMenu
                            v-model="query.closed"
                            :options="openclose"
                            value-attribute="closed"
                            :searchable-placeholder="$t('common.search_field')"
                        >
                            <template #label>
                                {{
                                    query.closed === undefined
                                        ? openclose[0]!.label
                                        : (openclose.findLast((o) => o.value === query.closed)?.label ?? $t('common.na'))
                                }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormGroup>

                <UFormGroup class="flex-1" name="relation" :label="$t('common.relation')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="query.relations"
                            multiple
                            :options="docRelations"
                            value-attribute="value"
                            :searchable-placeholder="$t('common.relation', 2)"
                        >
                            <template #label>
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
                </UFormGroup>
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
        :loading="loading"
        :columns="columns"
        :rows="data?.relations"
        :empty-state="{
            icon: 'i-mdi-file-multiple',
            label: $t('common.not_found', [`${$t('common.citizen', 1)} ${$t('common.document', 2)}`]),
        }"
    >
        <template #document-data="{ row: relation }">
            <DocumentInfoPopover :document="relation.document" load-on-open />
        </template>

        <template #closed-data="{ row: relation }">
            <OpenClosedBadge :closed="relation.document?.closed" variant="subtle" />
        </template>

        <template #relation-data="{ row: relation }">
            <UBadge
                class="font-semibold"
                :color="docRelationToBadge(relation.relation)"
                :icon="docRelationToIcon(relation.relation)"
            >
                {{ $t(`enums.documents.DocRelation.${DocRelation[relation.relation]}`) }}
            </UBadge>
        </template>

        <template #createdAt-data="{ row: relation }">
            <GenericTime :value="relation.createdAt" />
        </template>

        <template #creator-data="{ row: relation }">
            <CitizenInfoPopover :user="relation.sourceUser" />
        </template>
    </UTable>

    <Pagination v-model="query.page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
