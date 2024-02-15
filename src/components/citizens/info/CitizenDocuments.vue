<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { ArrowExpandIcon, ChevronRightIcon, FileDocumentMultipleIcon, LockIcon, LockOpenVariantIcon } from 'mdi-vue3';
import { ref } from 'vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
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
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 align-middle">
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
                            <div v-if="data?.relations.length > 0" class="text-neutral sm:hidden">
                                <ul role="list" class="mt-2 divide-y divide-gray-600 overflow-hidden rounded-lg sm:hidden">
                                    <li
                                        v-for="relation in data?.relations"
                                        :key="relation.id"
                                        class="block bg-base-800 px-4 py-4 hover:bg-base-700"
                                    >
                                        <span class="flex items-center space-x-4">
                                            <span class="flex flex-1 space-x-2 truncate">
                                                <ArrowExpandIcon
                                                    class="h-5 w-5 flex-shrink-0 text-gray-400"
                                                    aria-hidden="true"
                                                />
                                                <span class="flex flex-col truncate text-sm">
                                                    <span>
                                                        <NuxtLink
                                                            :to="{
                                                                name: 'documents-id',
                                                                params: {
                                                                    id: relation.documentId,
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
                                                            class="flex flex-initial flex-row gap-1 rounded-full bg-error-100 px-2 py-1"
                                                        >
                                                            <LockIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
                                                            <span class="text-sm font-medium text-error-700">
                                                                {{ $t('common.close', 2) }}
                                                            </span>
                                                        </div>
                                                        <div
                                                            v-else
                                                            class="flex flex-initial flex-row gap-1 rounded-full bg-success-100 px-2 py-1"
                                                        >
                                                            <LockOpenVariantIcon
                                                                class="h-5 w-5 text-success-500"
                                                                aria-hidden="true"
                                                            />
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
                                                            class="group inline-flex space-x-2 truncate text-sm"
                                                        >
                                                            {{
                                                                relation.targetUser?.firstname +
                                                                ', ' +
                                                                relation.targetUser?.lastname
                                                            }}
                                                        </NuxtLink>
                                                    </span>
                                                    <span class="font-medium">
                                                        {{ $t(`enums.docstore.DocRelation.${DocRelation[relation.relation]}`) }}
                                                    </span>
                                                    <span class="truncate">
                                                        {{ relation.sourceUser?.firstname }}
                                                        {{ relation.sourceUser?.lastname }}
                                                    </span>
                                                    <GenericTime :value="relation.createdAt" :ago="true" />
                                                </span>
                                            </span>
                                            <ChevronRightIcon
                                                class="h-5 w-5 flex-shrink-0 text-accent-200"
                                                aria-hidden="true"
                                            />
                                        </span>
                                    </li>
                                </ul>

                                <TablePagination
                                    :pagination="data?.pagination"
                                    :refresh="refresh"
                                    @offset-change="offset = $event"
                                />
                            </div>

                            <!-- Relations table (small breakpoint and up) -->
                            <div v-if="data?.relations.length > 0" class="hidden sm:block">
                                <div>
                                    <div class="mt-2 flex flex-col">
                                        <div class="min-w-full overflow-hidden overflow-x-auto align-middle sm:rounded-lg">
                                            <table class="mb-2 min-w-full bg-base-700 text-neutral">
                                                <thead>
                                                    <tr>
                                                        <th class="px-6 py-3 text-left text-sm font-semibold" scope="col">
                                                            {{ $t('common.document', 1) }}
                                                        </th>
                                                        <th class="px-6 py-3 text-left text-sm font-semibold" scope="col">
                                                            {{ $t('common.close', 2) }}
                                                        </th>
                                                        <th class="px-6 py-3 text-right text-sm font-semibold" scope="col">
                                                            {{ $t('common.relation', 1) }}
                                                        </th>
                                                        <th class="px-6 py-3 text-right text-sm font-semibold" scope="col">
                                                            {{ $t('common.date') }}
                                                        </th>
                                                        <th
                                                            class="hidden px-6 py-3 text-left text-sm font-semibold md:block"
                                                            scope="col"
                                                        >
                                                            {{ $t('common.creator') }}
                                                        </th>
                                                    </tr>
                                                </thead>
                                                <tbody class="divide-y divide-gray-600 bg-base-800 text-neutral">
                                                    <tr v-for="relation in data?.relations" :key="relation.id">
                                                        <td class="px-6 py-4 text-sm">
                                                            <NuxtLink
                                                                :to="{
                                                                    name: 'documents-id',
                                                                    params: {
                                                                        id: relation.documentId,
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
                                                                class="flex flex-initial flex-row gap-1 rounded-full bg-error-100 px-2 py-1"
                                                            >
                                                                <LockIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
                                                                <span class="text-sm font-medium text-error-700">
                                                                    {{ $t('common.close', 2) }}
                                                                </span>
                                                            </div>
                                                            <div
                                                                v-else
                                                                class="flex flex-initial flex-row gap-1 rounded-full bg-success-100 px-2 py-1"
                                                            >
                                                                <LockOpenVariantIcon
                                                                    class="h-5 w-5 text-success-500"
                                                                    aria-hidden="true"
                                                                />
                                                                <span class="text-sm font-medium text-success-700">
                                                                    {{ $t('common.open', 2) }}
                                                                </span>
                                                            </div>
                                                        </td>
                                                        <td class="whitespace-nowrap px-6 py-4 text-right text-sm">
                                                            <span class="font-medium">
                                                                {{
                                                                    $t(
                                                                        `enums.docstore.DocRelation.${
                                                                            DocRelation[relation.relation]
                                                                        }`,
                                                                    )
                                                                }}
                                                            </span>
                                                        </td>
                                                        <td class="whitespace-nowrap px-6 py-4 text-right text-sm">
                                                            <GenericTime :value="relation.createdAt" />
                                                        </td>
                                                        <td class="hidden whitespace-nowrap px-6 py-4 text-sm md:block">
                                                            <div class="flex">
                                                                <NuxtLink
                                                                    :to="{
                                                                        name: 'citizens-id',
                                                                        params: {
                                                                            id: relation.sourceUserId,
                                                                        },
                                                                    }"
                                                                    class="group inline-flex space-x-2 truncate text-sm"
                                                                >
                                                                    {{
                                                                        relation.sourceUser?.firstname +
                                                                        ', ' +
                                                                        relation.sourceUser?.lastname
                                                                    }}
                                                                </NuxtLink>
                                                            </div>
                                                        </td>
                                                    </tr>
                                                </tbody>
                                            </table>

                                            <TablePagination
                                                :pagination="data?.pagination"
                                                :refresh="refresh"
                                                @offset-change="offset = $event"
                                            />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
