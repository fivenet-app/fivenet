<script lang="ts" setup>
import DocumentList from '~/components/documents/DocumentList.vue';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';
import Pagination from '~/components/partials/Pagination.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import type { ListDocumentPinsResponse } from '~~/gen/ts/services/docstore/docstore';

useHead({
    title: 'pages.documents.title',
});

definePageMeta({
    title: 'pages.documents.title',
    requiresAuth: true,
    permission: 'DocStoreService.ListDocuments',
});

const { $grpc } = useNuxtApp();

const modal = useModal();

const { can } = useAuth();

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending: loading, error, refresh } = useLazyAsyncData(`calendars-${page.value}`, () => listCalendars());

async function listCalendars(): Promise<ListDocumentPinsResponse> {
    const call = $grpc.docstore.docStore.listDocumentPins({
        pagination: {
            offset: offset.value,
        },
    });
    const { response } = await call;

    return response;
}

const isOpen = ref(false);
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel class="shrink-0 border-b border-gray-200 lg:border-b-0 lg:border-r dark:border-gray-800" grow>
            <UDashboardNavbar :title="$t('pages.documents.title')">
                <template #right>
                    <UButtonGroup class="inline-flex 2xl:hidden">
                        <UButton
                            v-if="can('CompletorService.CompleteDocumentCategories').value"
                            :to="{ name: 'documents-categories' }"
                            icon="i-mdi-shape"
                            truncate
                        >
                            <span class="hidden truncate sm:block">
                                {{ $t('common.category', 2) }}
                            </span>
                        </UButton>

                        <UButton
                            v-if="can('DocStoreService.ListTemplates').value"
                            :to="{ name: 'documents-templates' }"
                            icon="i-mdi-file-code"
                            truncate
                        >
                            <span class="hidden truncate sm:block">
                                {{ $t('common.template', 2) }}
                            </span>
                        </UButton>
                    </UButtonGroup>

                    <UButton class="2xl:hidden" trailing-icon="i-mdi-pin" color="gray" truncate @click="isOpen = true">
                        <span class="hidden truncate sm:block">
                            {{ $t('common.pinned') }}
                        </span>
                    </UButton>

                    <UButton
                        v-if="can('DocStoreService.CreateDocument').value"
                        trailing-icon="i-mdi-plus"
                        color="gray"
                        truncate
                        @click="modal.open(TemplatesModal, {})"
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.create') }}
                        </span>
                    </UButton>
                </template>
            </UDashboardNavbar>

            <DocumentList />
        </UDashboardPanel>

        <UDashboardPanel
            v-model="isOpen"
            class="max-w-72"
            collapsible
            side="right"
            breakpoint="2xl"
            :ui="{ collapsible: 'lg:!hidden 2xl:!flex', slideover: 'lg:!flex 2xl:hidden' }"
        >
            <UDashboardNavbar>
                <template #toggle>
                    <UDashboardNavbarToggle class="lg:block 2xl:hidden" />
                </template>

                <template #right>
                    <UButtonGroup class="hidden truncate 2xl:inline-flex">
                        <UButton
                            v-if="can('CompletorService.CompleteDocumentCategories').value"
                            :to="{ name: 'documents-categories' }"
                            icon="i-mdi-shape"
                            truncate
                        >
                            <span class="truncate">
                                {{ $t('common.category', 2) }}
                            </span>
                        </UButton>

                        <UButton
                            v-if="can('DocStoreService.ListTemplates').value"
                            :to="{ name: 'documents-templates' }"
                            icon="i-mdi-file-code"
                            truncate
                        >
                            <span class="truncate">
                                {{ $t('common.template', 2) }}
                            </span>
                        </UButton>
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardPanelContent class="p-2">
                <UDashboardSection
                    :ui="{
                        wrapper: 'divide-y space-y-0 *:pt-2 first:*:pt-2 first:*:pt-2 mb-6',
                    }"
                    :title="$t('common.pinned_document', 2)"
                >
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

                        <template v-else-if="loading">
                            <USkeleton v-for="idx in 10" :key="idx" class="h-16 w-full" />
                        </template>

                        <template v-else>
                            <DocumentInfoPopover
                                v-for="doc in data?.documents"
                                :key="doc.id"
                                :document="doc"
                                button-class="line-clamp-3 hyphens-auto break-words flex flex-col items-start"
                                hide-trailing
                            />
                        </template>
                    </div>
                </UDashboardSection>
            </UDashboardPanelContent>

            <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
