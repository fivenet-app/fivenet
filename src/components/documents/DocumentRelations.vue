<script lang="ts" setup>
import { AccountMultipleIcon, ArrowCollapseIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { DocRelation, DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';

const { $grpc } = useNuxtApp();

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

const {
    data: relations,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-relations`, () => getDocumentRelations());

async function getDocumentRelations(): Promise<DocumentRelation[]> {
    try {
        const call = $grpc.getDocStoreClient().getDocumentRelations({
            documentId: props.documentId,
        });
        const { response } = await call;

        return response.relations;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <div>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.relation', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.relation', 2)])" :retry="refresh" />
        <DataNoDataBlock
            v-if="!relations || relations.length === 0"
            :type="`${$t('common.document', 1)} ${$t('common.relation', 2)}`"
            :icon="AccountMultipleIcon"
        />

        <template v-else>
            <!-- Relations list (smallest breakpoint only) -->
            <div class="text-neutral sm:hidden">
                <ul role="list" class="divide-y divide-gray-600 overflow-hidden rounded-lg sm:hidden">
                    <li v-for="relation in relations" :key="relation.id" class="block bg-base-800 p-4 hover:bg-base-700">
                        <span class="flex items-center space-x-4">
                            <span class="flex flex-1 space-x-2 truncate">
                                <ArrowCollapseIcon class="size-5 shrink-0 text-gray-400" aria-hidden="true" />
                                <span class="flex flex-col truncate text-sm">
                                    <span v-if="showDocument">
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
                                                class="mr-1 inline-flex items-center rounded-md bg-primary-400/10 px-2 py-1 text-xs font-medium text-primary-400 ring-1 ring-inset ring-primary-400/30"
                                            >
                                                {{ relation.document?.category?.name }}
                                            </span>
                                            {{ relation.document?.title }}
                                        </NuxtLink>
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
                            <table class="w-full divide-y divide-base-400 bg-background text-neutral">
                                <thead>
                                    <tr>
                                        <th v-if="showDocument" class="px-6 py-3 text-left text-sm font-semibold" scope="col">
                                            {{ $t('common.document', 1) }}
                                        </th>
                                        <th class="px-6 py-3 text-left text-sm font-semibold" scope="col">
                                            {{ $t('common.target') }}
                                        </th>
                                        <th class="px-6 py-3 text-right text-sm font-semibold" scope="col">
                                            {{ $t('common.relation', 1) }}
                                        </th>
                                        <th
                                            v-if="showSource"
                                            class="hidden px-6 py-3 text-left text-sm font-semibold md:block"
                                            scope="col"
                                        >
                                            {{ $t('common.creator') }}
                                        </th>
                                        <th class="px-6 py-3 text-right text-sm font-semibold" scope="col">
                                            {{ $t('common.date') }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-600 bg-base-700 text-neutral">
                                    <tr v-for="relation in relations" :key="relation.id">
                                        <td v-if="showDocument" class="px-6 py-4 text-sm">
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
                                                    class="mr-1 inline-flex items-center rounded-md bg-primary-400/10 px-2 py-1 text-xs font-medium text-primary-400 ring-1 ring-inset ring-primary-400/30"
                                                >
                                                    {{ relation.document?.category?.name }}
                                                </span>
                                                {{ relation.document?.title }}
                                            </NuxtLink>
                                        </td>
                                        <td class="px-6 py-4 text-sm">
                                            <div class="flex">
                                                <span class="inline-flex items-center gap-1">
                                                    <CitizenInfoPopover :user="relation.targetUser" />
                                                    ({{ relation.targetUser?.dateofbirth }})
                                                </span>
                                            </div>
                                        </td>
                                        <td class="whitespace-nowrap px-6 py-4 text-right text-sm">
                                            <span class="font-medium">
                                                {{ $t(`enums.docstore.DocRelation.${DocRelation[relation.relation]}`) }}
                                            </span>
                                        </td>
                                        <td v-if="showSource" class="hidden whitespace-nowrap px-6 py-4 text-sm md:block">
                                            <div class="flex">
                                                <CitizenInfoPopover :user="relation.sourceUser" />
                                            </div>
                                        </td>
                                        <td class="whitespace-nowrap px-6 py-4 text-right text-sm">
                                            <GenericTime :value="relation.createdAt" />
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </div>
</template>
