<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import { ListUserDocumentsResponse } from '~~/gen/ts/services/docstore/docstore';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const props = defineProps<{
    userId: number;
}>();

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending, refresh, error } = useLazyAsyncData(`citizeninfo-documents-${props.userId}-${page.value}`, () =>
    listUserDocuments(),
);

async function listUserDocuments(): Promise<ListUserDocumentsResponse> {
    try {
        const call = $grpc.getDocStoreClient().listUserDocuments({
            pagination: {
                offset: offset.value,
            },
            userId: props.userId,
            relations: [],
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
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
    <div>
        <DataErrorBlock
            v-if="error"
            :title="$t('common.unable_to_load', [`${$t('common.citizen', 1)} ${$t('common.document', 2)}`])"
            :retry="refresh"
        />
        <template v-else>
            <UTable
                :loading="pending"
                :columns="columns"
                :rows="data?.relations"
                :empty-state="{
                    icon: 'i-mdi-file-multiple',
                    label: $t('common.not_found', [`${$t('common.citizen', 1)} ${$t('common.document', 2)}`]),
                }"
            >
                <template #document-data="{ row }">
                    <UBadge v-if="row.document?.category">
                        {{ row.document.category.name }}
                    </UBadge>
                    <UButton
                        variant="link"
                        truncate
                        class="w-full truncate sm:max-w-sm md:max-w-md lg:max-w-full"
                        :to="{
                            name: 'documents-id',
                            params: {
                                id: row.documentId,
                            },
                        }"
                    >
                        {{ row.document?.title ?? $t('common.na') }}
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

            <div class="flex justify-end border-t border-gray-200 px-3 py-3.5 dark:border-gray-700">
                <UPagination
                    v-model="page"
                    :page-count="data?.pagination?.pageSize ?? 0"
                    :total="data?.pagination?.totalCount ?? 0"
                />
            </div>
        </template>
    </div>
</template>
