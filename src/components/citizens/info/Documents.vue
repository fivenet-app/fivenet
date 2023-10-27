<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { ArrowExpandIcon, ChevronRightIcon, FileDocumentMultipleIcon, LockIcon, LockOpenVariantIcon } from 'mdi-vue3';
import { ref } from 'vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import Time from '~/components/partials/elements/Time.vue';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import { ListUserDocumentsResponse } from '~~/gen/ts/services/docstore/docstore';

const { $grpc } = useNuxtApp();

const props = defineProps<{
    userId: number;
}>();

const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`citizeninfo-documents-${props.userId}-${offset.value}`, () =>
    listUserDocuments(),
);

async function listUserDocuments(): Promise<ListUserDocumentsResponse> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getDocStoreClient().listUserDocuments({
                pagination: {
                    offset: offset.value,
                },
                userId: props.userId,
                relations: [],
            });
            const { response } = await call;

            return res(response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watch(offset, async () => refresh());
</script>

<template>
    <div class="mt-2">
        <DataPendingBlock
            v-if="pending"
            :message="$t('common.loading', [`${$t('common.user', 1)} ${$t('common.document', 2)}`])"
        />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [`${$t('common.user', 1)} ${$t('common.document', 2)}`])"
            :retry="refresh"
        />
        <DataNoDataBlock
            v-else-if="data?.relations.length === 0"
            :type="`${$t('common.document', 1)} ${$t('common.relation', 2)}`"
            :icon="FileDocumentMultipleIcon"
        />
        <div v-else-if="data?.relations">
            <!-- Relations list (smallest breakpoint only) -->
            <div v-if="data?.relations.length > 0" class="sm:hidden text-neutral">
                <ul role="list" class="mt-2 overflow-hidden divide-y divide-gray-600 rounded-lg sm:hidden">
                    <li v-for="relation in data?.relations" :key="relation.id?.toString()">
                        <a href="#" class="block px-4 py-4 bg-base-800 hover:bg-base-700">
                            <span class="flex items-center space-x-4">
                                <span class="flex flex-1 space-x-2 truncate">
                                    <ArrowExpandIcon class="flex-shrink-0 w-5 h-5 text-gray-400" aria-hidden="true" />
                                    <span class="flex flex-col text-sm truncate">
                                        <span>
                                            <NuxtLink
                                                :to="{
                                                    name: 'documents-id',
                                                    params: {
                                                        id: relation.documentId.toString(),
                                                    },
                                                }"
                                            >
                                                {{ relation.document?.title
                                                }}<span v-if="relation.document?.category">
                                                    (Category:
                                                    {{ relation.document?.category?.name }})</span
                                                >
                                            </NuxtLink>
                                        </span>
                                        <span>
                                            <div
                                                v-if="relation.document?.closed"
                                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100"
                                            >
                                                <LockIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                                                <span class="text-sm font-medium text-error-700">
                                                    {{ $t('common.close', 2) }}
                                                </span>
                                            </div>
                                            <div
                                                v-else
                                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100"
                                            >
                                                <LockOpenVariantIcon class="w-5 h-5 text-success-500" aria-hidden="true" />
                                                <span class="text-sm font-medium text-success-700">
                                                    {{ $t('common.open', 2) }}
                                                </span>
                                            </div>
                                        </span>
                                        <span>
                                            <NuxtLink
                                                :to="{
                                                    name: 'citizens-id',
                                                    params: {
                                                        id: relation.targetUserId,
                                                    },
                                                }"
                                                class="inline-flex space-x-2 text-sm truncate group"
                                            >
                                                {{ relation.targetUser?.firstname + ', ' + relation.targetUser?.lastname }}
                                            </NuxtLink>
                                        </span>
                                        <span class="font-medium">
                                            {{ $t(`enums.docstore.DocRelation.${DocRelation[relation.relation]}`) }}
                                        </span>
                                        <span class="truncate">
                                            {{ relation.sourceUser?.firstname }},
                                            {{ relation.sourceUser?.lastname }}
                                        </span>
                                        <Time :value="relation.createdAt" :ago="true" />
                                    </span>
                                </span>
                                <ChevronRightIcon class="flex-shrink-0 w-5 h-5 text-base-200" aria-hidden="true" />
                            </span>
                        </a>
                    </li>
                </ul>

                <TablePagination :pagination="data?.pagination" @offset-change="offset = $event" />
            </div>

            <!-- Relations table (small breakpoint and up) -->
            <div v-if="data?.relations.length > 0" class="hidden sm:block">
                <div>
                    <div class="flex flex-col mt-2">
                        <div class="min-w-full overflow-hidden overflow-x-auto align-middle sm:rounded-lg">
                            <table class="min-w-full bg-base-700 text-neutral mb-2">
                                <thead>
                                    <tr>
                                        <th class="px-6 py-3 text-sm font-semibold text-left" scope="col">
                                            {{ $t('common.document', 1) }}
                                        </th>
                                        <th class="px-6 py-3 text-sm font-semibold text-left" scope="col">
                                            {{ $t('common.close', 2) }}
                                        </th>
                                        <th class="px-6 py-3 text-sm font-semibold text-right" scope="col">
                                            {{ $t('common.relation', 1) }}
                                        </th>
                                        <th class="px-6 py-3 text-sm font-semibold text-right" scope="col">
                                            {{ $t('common.date') }}
                                        </th>
                                        <th class="hidden px-6 py-3 text-sm font-semibold text-left md:block" scope="col">
                                            {{ $t('common.creator') }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-600 bg-base-800 text-neutral">
                                    <tr v-for="relation in data?.relations" :key="relation.id?.toString()">
                                        <td class="px-6 py-4 text-sm">
                                            <NuxtLink
                                                :to="{
                                                    name: 'documents-id',
                                                    params: {
                                                        id: relation.documentId.toString(),
                                                    },
                                                }"
                                            >
                                                <span
                                                    v-if="relation.document?.category"
                                                    class="inline-flex items-center rounded-md bg-primary-400/10 px-2 py-1 text-xs font-medium text-primary-400 ring-1 ring-inset ring-primary-400/30"
                                                >
                                                    {{ relation.document.category.name }}
                                                </span>
                                                {{ relation.document?.title }}
                                            </NuxtLink>
                                        </td>
                                        <td class="px-6 py-4 text-sm">
                                            <div
                                                v-if="relation.document?.closed"
                                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100"
                                            >
                                                <LockIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                                                <span class="text-sm font-medium text-error-700">
                                                    {{ $t('common.close', 2) }}
                                                </span>
                                            </div>
                                            <div
                                                v-else
                                                class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100"
                                            >
                                                <LockOpenVariantIcon class="w-5 h-5 text-success-500" aria-hidden="true" />
                                                <span class="text-sm font-medium text-success-700">
                                                    {{ $t('common.open', 2) }}
                                                </span>
                                            </div>
                                        </td>
                                        <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                            <span class="font-medium">
                                                {{ $t(`enums.docstore.DocRelation.${DocRelation[relation.relation]}`) }}
                                            </span>
                                        </td>
                                        <td class="px-6 py-4 text-sm text-right whitespace-nowrap">
                                            <Time :value="relation.createdAt" />
                                        </td>
                                        <td class="hidden px-6 py-4 text-sm whitespace-nowrap md:block">
                                            <div class="flex">
                                                <NuxtLink
                                                    :to="{
                                                        name: 'citizens-id',
                                                        params: {
                                                            id: relation.sourceUserId,
                                                        },
                                                    }"
                                                    class="inline-flex space-x-2 text-sm truncate group"
                                                >
                                                    {{ relation.sourceUser?.firstname + ', ' + relation.sourceUser?.lastname }}
                                                </NuxtLink>
                                            </div>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>

                            <TablePagination :pagination="data?.pagination" @offset-change="offset = $event" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
