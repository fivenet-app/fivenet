<script lang="ts" setup>
import DocumentList from '~/components/documents/DocumentList.vue';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';
import Pagination from '~/components/partials/Pagination.vue';
import DashboardPanel from '~/components/partials/dashboard/DashboardPanel.vue';
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

const modal = useModal();

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending: loading, error, refresh } = useLazyAsyncData(`calendars-${page.value}`, () => listCalendars());

async function listCalendars(): Promise<ListDocumentPinsResponse> {
    try {
        const call = getGRPCDocStoreClient().listDocumentPins({
            pagination: {
                offset: offset.value,
            },
        });
        const { response } = await call;

        return response;
    } catch (e) {
        throw e;
    }
}

const isOpen = ref(false);
</script>

<template>
    <UDashboardPage class="h-full">
        <UDashboardPanel
            class="h-full flex-shrink-0 border-b border-gray-200 lg:border-b-0 lg:border-r dark:border-gray-800"
            grow
        >
            <UDashboardNavbar :title="$t('pages.documents.title')">
                <template #right>
                    <UButtonGroup class="inline-flex 2xl:hidden">
                        <UButton
                            v-if="can('CompletorService.CompleteDocumentCategories')"
                            :to="{ name: 'documents-categories' }"
                            icon="i-mdi-shape"
                        >
                            {{ $t('common.category', 2) }}
                        </UButton>

                        <UButton
                            v-if="can('DocStoreService.ListTemplates')"
                            :to="{ name: 'documents-templates' }"
                            icon="i-mdi-file-code"
                        >
                            {{ $t('common.template', 2) }}
                        </UButton>
                    </UButtonGroup>

                    <UButton
                        :label="$t('common.pinned')"
                        trailing-icon="i-mdi-pin"
                        color="gray"
                        class="2xl:hidden"
                        @click="isOpen = true"
                    />

                    <UButton
                        :label="$t('common.create')"
                        trailing-icon="i-mdi-plus"
                        color="gray"
                        @click="modal.open(TemplatesModal, {})"
                    />
                </template>
            </UDashboardNavbar>

            <DocumentList />
        </UDashboardPanel>

        <DashboardPanel
            v-model="isOpen"
            collapsible
            side="right"
            class="max-w-72 flex-1"
            :ui="{ collapsible: 'lg:!hidden 2xl:!flex', slideover: 'lg:!flex 2xl:hidden' }"
        >
            <UDashboardNavbar>
                <template #toggle>
                    <UDashboardNavbarToggle class="lg:block 2xl:hidden" />
                </template>

                <template #right>
                    <UButtonGroup class="hidden 2xl:inline-flex">
                        <UButton
                            v-if="can('CompletorService.CompleteDocumentCategories')"
                            :to="{ name: 'documents-categories' }"
                            icon="i-mdi-shape"
                        >
                            {{ $t('common.category', 2) }}
                        </UButton>

                        <UButton
                            v-if="can('DocStoreService.ListTemplates')"
                            :to="{ name: 'documents-templates' }"
                            icon="i-mdi-file-code"
                        >
                            {{ $t('common.template', 2) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardPanelContent>
                <div class="flex flex-1 flex-col">
                    <UDashboardSection
                        :ui="{
                            wrapper: 'divide-y space-y-0 *:pt-0 first:*:pt-0 first:*:pt-0 mb-6',
                        }"
                        :title="$t('common.pinned_document', 2)"
                    >
                        <div>
                            <DataErrorBlock
                                v-if="error"
                                :title="$t('common.unable_to_load', [$t('common.pinned_document', 2)])"
                                :retry="refresh"
                            />
                            <DataNoDataBlock
                                v-else-if="!data || data.documents.length === 0"
                                icon="i-mdi-pin-outline"
                                :type="$t('common.pinned_document', 2)"
                            />

                            <ul
                                v-else-if="loading"
                                role="list"
                                class="my-1 flex flex-col gap-1 divide-y divide-gray-100 dark:divide-gray-800"
                            >
                                <li v-if="loading" v-for="_ in 10">
                                    <USkeleton class="h-16 w-full" />
                                </li>

                                <li v-else v-for="document in data?.documents" class="flex flex-col">
                                    <DocumentInfoPopover :document="document" class="line-clamp-3 text-left" />
                                </li>
                            </ul>
                        </div>
                    </UDashboardSection>
                </div>
            </UDashboardPanelContent>

            <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
        </DashboardPanel>
    </UDashboardPage>
</template>
