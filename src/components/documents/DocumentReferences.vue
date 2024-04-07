<script lang="ts" setup>
import { ArrowCollapseIcon, FileDocumentMultipleIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { DocReference, DocumentReference } from '~~/gen/ts/resources/documents/documents';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';

const { $grpc } = useNuxtApp();

const props = withDefaults(
    defineProps<{
        documentId: string;
        showSource?: boolean;
    }>(),
    {
        showSource: true,
    },
);

const {
    data: references,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId}-references`, () => getDocumentReferences());

async function getDocumentReferences(): Promise<DocumentReference[]> {
    try {
        const call = $grpc.getDocStoreClient().getDocumentReferences({
            documentId: props.documentId,
        });
        const { response } = await call;

        return response.references;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <div>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.reference', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.reference', 2)])" :retry="refresh" />
        <DataNoDataBlock
            v-else-if="references === null || references.length === 0"
            :type="`${$t('common.document', 1)} ${$t('common.reference', 2)}`"
            :mdi="FileDocumentMultipleIcon"
        />

        <template v-else>
            <!-- Relations list (smallest breakpoint only) -->
            <div class="text-neutral sm:hidden">
                <ul role="list" class="divide-y divide-gray-600 rounded-lg sm:hidden">
                    <li v-for="reference in references" :key="reference.id">
                        <NuxtLink
                            :to="{
                                name: 'documents-id',
                                params: { id: reference.targetDocumentId },
                            }"
                            class="block bg-base-800 p-4 hover:bg-base-700"
                        >
                            <span class="flex items-center space-x-4">
                                <span class="flex flex-1 space-x-2 truncate">
                                    <ArrowCollapseIcon class="size-5 shrink-0" />
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
                        </NuxtLink>
                    </li>
                </ul>
            </div>

            <!-- Relations table (small breakpoint and up) -->
            <div class="hidden sm:block">
                <div>
                    <div class="flex flex-col">
                        <div class="min-w-full overflow-hidden overflow-x-auto align-middle sm:rounded-lg">
                            <table class="bg-background min-w-full divide-y divide-base-600 border border-gray-600">
                                <thead>
                                    <tr>
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
                                            {{ $t('common.source') }}
                                        </th>
                                        <th class="hidden px-6 py-3 text-left text-sm font-semibold md:block" scope="col">
                                            {{ $t('common.creator') }}
                                        </th>
                                        <th class="px-6 py-3 text-right text-sm font-semibold" scope="col">
                                            {{ $t('common.date') }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-600 bg-base-700">
                                    <tr v-for="reference in references" :key="reference.id">
                                        <td class="px-6 py-4 text-sm">
                                            <div class="flex">
                                                <NuxtLink
                                                    :to="{
                                                        name: 'documents-id',
                                                        params: {
                                                            id: reference.targetDocumentId,
                                                        },
                                                    }"
                                                    class="group inline-flex max-w-xl space-x-2 truncate text-sm"
                                                >
                                                    <span
                                                        v-if="reference.targetDocument?.category"
                                                        class="bg-primary-400/10 text-primary-400 ring-primary-400/30 mr-1 inline-flex items-center rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset"
                                                    >
                                                        {{ reference.targetDocument?.category?.name }}
                                                    </span>
                                                    {{ reference.targetDocument?.title }}
                                                </NuxtLink>
                                            </div>
                                        </td>
                                        <td class="whitespace-nowrap px-6 py-4 text-right text-sm">
                                            <span class="font-medium">
                                                {{ $t(`enums.docstore.DocReference.${DocReference[reference.reference]}`) }}
                                            </span>
                                        </td>
                                        <td v-if="showSource" class="hidden whitespace-nowrap px-6 py-4 text-sm md:block">
                                            <div class="flex">
                                                <NuxtLink
                                                    :to="{
                                                        name: 'documents-id',
                                                        params: {
                                                            id: reference.sourceDocumentId,
                                                        },
                                                    }"
                                                    class="group inline-flex max-w-xl space-x-1 truncate text-sm"
                                                >
                                                    <span
                                                        v-if="reference.sourceDocument?.category"
                                                        class="bg-primary-400/10 text-primary-400 ring-primary-400/30 mr-1 inline-flex items-center rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset"
                                                    >
                                                        {{ reference.sourceDocument?.category?.name }}
                                                    </span>
                                                    {{ reference.sourceDocument?.title }}
                                                </NuxtLink>
                                            </div>
                                        </td>
                                        <td class="whitespace-nowrap px-6 py-4 text-right text-sm">
                                            <div class="flex">
                                                <NuxtLink :to="{ name: 'citizens-id', params: { id: reference.creatorId! } }">
                                                    <CitizenInfoPopover :user="reference.creator" />
                                                </NuxtLink>
                                            </div>
                                        </td>
                                        <td class="whitespace-nowrap px-6 py-4 text-right text-sm">
                                            <GenericTime :value="reference.createdAt" />
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
