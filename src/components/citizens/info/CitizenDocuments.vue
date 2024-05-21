<script lang="ts" setup>
import Pagination from '~/components/partials/Pagination.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import { ListUserDocumentsResponse } from '~~/gen/ts/services/docstore/docstore';

const { t } = useI18n();

const props = defineProps<{
    userId: number;
}>();

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
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());

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
    <DataErrorBlock
        v-if="error"
        :title="$t('common.unable_to_load', [`${$t('common.citizen', 1)} ${$t('common.document', 2)}`])"
        :retry="refresh"
    />
    <template v-else>
        <UTable
            :loading="loading"
            :columns="columns"
            :rows="data?.relations"
            :empty-state="{
                icon: 'i-mdi-file-multiple',
                label: $t('common.not_found', [`${$t('common.citizen', 1)} ${$t('common.document', 2)}`]),
            }"
        >
            <template #document-data="{ row }">
                <UButton
                    variant="link"
                    :to="{
                        name: 'documents-id',
                        params: {
                            id: row.documentId,
                        },
                    }"
                >
                    <UBadge v-if="row.document?.category">
                        {{ row.document.category.name }}
                    </UBadge>

                    <div class="line-clamp-2 w-full whitespace-normal break-words hover:line-clamp-4">
                        {{ row.document?.title ?? $t('common.na') }}
                    </div>
                </UButton>
            </template>
            <template #closed-data="{ row }">
                <div class="inline-flex gap-1">
                    <template v-if="row.document?.closed">
                        <UIcon name="i-mdi-lock" class="size-5 text-error-400" />
                        <span class="text-sm font-medium text-error-400">
                            {{ $t('common.close', 2) }}
                        </span>
                    </template>
                    <template v-else>
                        <UIcon name="i-mdi-lock-open-variant" class="size-5 text-success-400" />
                        <span class="text-sm font-medium text-success-400">
                            {{ $t('common.open', 2) }}
                        </span>
                    </template>
                </div>
            </template>
            <template #relation-data="{ row }">
                <span class="font-medium">
                    {{ $t(`enums.docstore.DocRelation.${DocRelation[row.relation]}`) }}
                </span>
            </template>
            <template #date-data="{ row }">
                <GenericTime :value="row.createdAt" />
            </template>
            <template #creator-data="{ row }">
                <CitizenInfoPopover :user="row.sourceUser" />
            </template>
        </UTable>

        <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
    </template>
</template>
