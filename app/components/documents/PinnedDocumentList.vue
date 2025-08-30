<script lang="ts" setup>
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import type { ListDocumentPinsResponse, ToggleDocumentPinResponse } from '~~/gen/ts/services/documents/documents';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DocumentCategoryBadge from '../partials/documents/DocumentCategoryBadge.vue';
import DocumentInfoPopover from '../partials/documents/DocumentInfoPopover.vue';
import IDCopyBadge from '../partials/IDCopyBadge.vue';
import Pagination from '../partials/Pagination.vue';

defineEmits<{
    (e: 'close'): void;
}>();

const { attr, can } = useAuth();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const page = useRouteQuery('page', '1', { transform: Number });

const { data, status, error, refresh } = useLazyAsyncData(`calendars-${page.value}`, () => listDocumentPins(), {
    immediate: can('documents.DocumentsService/ToggleDocumentPin').value,
});

async function listDocumentPins(): Promise<ListDocumentPinsResponse> {
    const call = documentsDocumentsClient.listDocumentPins({
        pagination: {
            offset: calculateOffset(page.value, data.value?.pagination),
        },
    });
    const { response } = await call;

    return response;
}

async function togglePin(documentId: number, state: boolean, personal: boolean): Promise<ToggleDocumentPinResponse> {
    try {
        const call = documentsDocumentsClient.toggleDocumentPin({
            documentId: documentId,
            state: state,
            personal: personal,
        });
        const { response } = await call;

        const idx = data.value?.documents.findIndex((d) => d.id === documentId);
        if (idx && idx > -1 && data.value?.documents[idx]) {
            if (!response.pin?.job && !response.pin?.userId) {
                data.value.documents.splice(idx, 1);
            } else {
                data.value.documents[idx].pin = response.pin;
            }
        }

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const editing = ref(false);
</script>

<template>
    <UDashboardPanel
        id="documents-pinnedlist"
        class="overflow-x-hidden"
        side="right"
        breakpoint="2xl"
        :width="350"
        :min-size="25"
        :max-size="50"
    >
        <template #header>
            <UDashboardNavbar :title="$t('common.pinned_document', 2)">
                <template #toggle>
                    <UButton
                        class="lg:block 2xl:hidden"
                        icon="i-mdi-menu"
                        variant="ghost"
                        color="neutral"
                        @click="$emit('close')"
                    />
                </template>

                <template #right>
                    <UTooltip
                        v-if="can('documents.DocumentsService/ToggleDocumentPin').value"
                        :text="editing ? $t('common.save') : $t('common.edit')"
                    >
                        <UButton
                            variant="link"
                            :icon="editing ? 'i-mdi-content-save' : 'i-mdi-pencil'"
                            @click="editing = !editing"
                        />
                    </UTooltip>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div class="flex flex-col gap-2">
                <DataErrorBlock
                    v-if="error"
                    :title="$t('common.unable_to_load', [$t('common.pinned_document', 2)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="!data || data.documents.length === 0"
                    icon="i-mdi-pin-outline"
                    :type="$t('common.pinned_document', 0)"
                />

                <template v-else-if="isRequestPending(status)">
                    <USkeleton v-for="idx in 10" :key="idx" class="h-16 w-full" />
                </template>
                <div v-else class="flex flex-col gap-2">
                    <div v-for="doc in data?.documents" :key="doc.id" class="flex flex-row gap-1 divide-x divide-default">
                        <UButtonGroup
                            v-if="editing && can('documents.DocumentsService/ToggleDocumentPin').value"
                            class="inline-flex items-center gap-1"
                            orientation="vertical"
                        >
                            <UTooltip :text="doc.pin?.state && doc.pin?.userId ? $t('common.pin', 1) : $t('common.unpin')">
                                <UButton
                                    class="shrink-0 flex-col text-center"
                                    variant="link"
                                    size="xs"
                                    :color="doc.pin?.state && doc.pin?.userId ? 'error' : 'primary'"
                                    @click="togglePin(doc.id, !doc.pin?.userId, true)"
                                >
                                    <UIcon
                                        class="size-5"
                                        :name="
                                            doc.pin?.state && doc.pin?.userId ? 'i-mdi-playlist-remove' : 'i-mdi-playlist-plus'
                                        "
                                    />
                                    {{ $t('common.personal') }}
                                </UButton>
                                <UTooltip
                                    v-if="attr('documents.DocumentsService/ToggleDocumentPin', 'Types', 'JobWide').value"
                                    :text="doc.pin?.state && doc.pin?.job ? $t('common.pin', 1) : $t('common.unpin')"
                                >
                                    <UButton
                                        class="shrink-0 flex-col text-center"
                                        variant="link"
                                        size="xs"
                                        :color="doc.pin?.state && doc.pin?.job ? 'error' : 'primary'"
                                        @click="togglePin(doc.id, !doc.pin?.job, false)"
                                    >
                                        <UIcon
                                            class="size-5"
                                            :name="doc.pin?.state && doc.pin?.job ? 'i-mdi-pin-off' : 'i-mdi-pin'"
                                        />
                                        {{ $t('common.job') }}
                                    </UButton>
                                </UTooltip>
                            </UTooltip>
                        </UButtonGroup>

                        <div class="flex-1 pr-1">
                            <DocumentInfoPopover
                                v-if="doc.createdAt !== undefined"
                                class="flex-1"
                                :document="doc"
                                button-class="hyphens-auto  flex-col items-start"
                                load-on-open
                                hide-trailing
                            >
                                <template #title="{ document }">
                                    <div class="inline-flex items-center gap-1 overflow-hidden">
                                        <IDCopyBadge :id="document?.id" prefix="DOC" size="xs" disable-tooltip />
                                        <DocumentCategoryBadge v-if="document?.category" :category="document?.category" />
                                    </div>

                                    <span class="line-clamp-2 text-left break-words hover:line-clamp-4">{{
                                        document?.title
                                    }}</span>
                                </template>
                            </DocumentInfoPopover>

                            <DocumentInfoPopover
                                v-else
                                class="flex-1"
                                :document-id="doc.id"
                                disable-tooltip
                                load-on-open
                                hide-trailing
                                button-class="flex-col items-start"
                            >
                                <template #title>
                                    <IDCopyBadge :id="doc?.id" prefix="DOC" size="xs" disable-tooltip />

                                    <UBadge :label="$t('common.no_access_to_document')" color="error" size="md" />
                                </template>
                            </DocumentInfoPopover>
                        </div>
                    </div>
                </div>
            </div>

            <Pagination v-model="page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
        </template>
    </UDashboardPanel>
</template>
