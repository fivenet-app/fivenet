<script lang="ts" setup>
import type { ContextMenuItem } from '@nuxt/ui';
import { watchDebounced } from '@vueuse/shared';
import { addDays } from 'date-fns';
import { z } from 'zod';
import DocumentListEntry from '~/components/documents/DocumentListEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DateRangePickerPopoverClient from '~/components/partials/DateRangePickerPopover.client.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { useCompletorStore } from '~/stores/completor';
import { useSettingsStore } from '~/stores/settings';
import type { ToggleItem } from '~/utils/types';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { ListDocumentsRequest, ListDocumentsResponse } from '~~/gen/ts/services/documents/documents';
import DocumentCategoryBadge from '../partials/documents/DocumentCategoryBadge.vue';
import PinnedDocumentList from './PinnedDocumentList.vue';
import TemplateModal from './templates/TemplateModal.vue';

const { t } = useI18n();

const { can, attr, isSuperuser } = useAuth();

const overlay = useOverlay();

const completorStore = useCompletorStore();

const settingsStore = useSettingsStore();
const { design } = storeToRefs(settingsStore);

const documentsDocuments = await useDocumentsDocuments();

const openclose: ToggleItem[] = [
    { id: 0, label: t('common.not_selected'), value: undefined },
    { id: 1, label: t('common.open', 2), value: false },
    { id: 2, label: t('common.close', 2), value: true },
];

const onlyDrafts: ToggleItem[] = [
    { id: 0, label: t('common.all_documents'), value: undefined },
    { id: 1, label: t('common.only_published'), value: false },
    { id: 2, label: t('common.only_drafts'), value: true },
];

const schema = z.object({
    documentIds: z.string().max(16).optional(),
    title: z.string().max(64).optional().default(''),
    creators: z.coerce.number().array().max(5).default([]),
    date: z
        .object({
            start: z.coerce.date(),
            end: z.coerce.date(),
        })
        .optional(),
    closed: z.coerce.boolean().optional(),
    categories: z.number().array().max(3).default([]),
    onlyDrafts: z.coerce.boolean().optional(),
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'createdAt',
                        desc: true,
                    },
                ]),
        })
        .default({ columns: [{ id: 'createdAt', desc: true }] }),
    page: pageNumberSchema,
});

const query = useSearchForm('documents', schema);

const usersLoading = ref(false);

const { data, status, refresh, error } = useLazyAsyncData(
    () => `documents-${JSON.stringify(query.sorting)}-${query.page}`,
    () => listDocuments(),
);

async function listDocuments(): Promise<ListDocumentsResponse> {
    const req: ListDocumentsRequest = {
        pagination: {
            offset: calculateOffset(query.page, data.value?.pagination),
        },
        sort: query.sorting,
        search: query.title ?? '',
        categoryIds: query.categories,
        creatorIds: query.creators,
        documentIds: [],
        onlyDrafts: query.onlyDrafts,
    };

    if (query.documentIds) {
        const id = parseInt(query.documentIds.trim().replaceAll('-', '').replace(/\D/g, ''));
        if (id > 0) {
            req.documentIds.push(id);
        }
    }
    if (query.date) {
        req.from = {
            timestamp: googleProtobufTimestamp.Timestamp.fromDate(query.date.start),
        };
        req.to = {
            timestamp: googleProtobufTimestamp.Timestamp.fromDate(query.date.end),
        };
    }
    if (query.closed !== undefined) {
        req.closed = query.closed;
    }

    return documentsDocuments.listDocuments(req);
}

watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

const categoriesLoading = ref(false);

const links = computed(() =>
    (
        [
            [
                {
                    label: t('common.open'),
                    icon: 'i-mdi-eye',
                    onSelect: () => navigateTo(`/documents/0`),
                },
                isSuperuser.value // && selectedDocument.value?.deletedAt
                    ? {
                          label: t('common.restore'),
                          icon: 'i-mdi-restore',
                          onSelect: () => navigateTo(`/documents/0`),
                      }
                    : undefined,
            ].filter((l) => l != undefined),
            [
                ...(can('documents.DocumentsService/ToggleDocumentPin').value
                    ? [
                          {
                              label: `${t('common.pin')}: ${t('common.personal')}`,
                              icon: 'i-mdi-playlist-plus',
                              to: '/components/vertical-navigation',
                          },
                          attr('documents.DocumentsService/ToggleDocumentPin', 'Types', 'JobWide').value
                              ? {
                                    label: `${t('common.pin')}: ${t('common.job')}`,
                                    icon: 'i-mdi-pin',
                                    to: '/components/vertical-navigation',
                                }
                              : undefined,
                      ].filter((l) => l != undefined)
                    : []),
            ],
        ] as ContextMenuItem[][]
    ).filter((l) => l.length > 0),
);

const inputRef = useTemplateRef('inputRef');

const templateModal = overlay.create(TemplateModal);

defineShortcuts({
    '/': () => inputRef.value?.inputRef?.focus(),
});
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.documents.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <UButton class="2xl:hidden" trailing-icon="i-mdi-pin" color="neutral" truncate @click="isOpen = true">
                        <span class="hidden truncate sm:block">
                            {{ $t('common.pinned') }}
                        </span>
                    </UButton>

                    <UButtonGroup class="inline-flex">
                        <UButton
                            v-if="can('completor.CompletorService/CompleteDocumentCategories').value"
                            :to="{ name: 'documents-categories' }"
                            icon="i-mdi-shape"
                            truncate
                        >
                            <span class="hidden truncate sm:block">
                                {{ $t('common.category', 2) }}
                            </span>
                        </UButton>

                        <UButton
                            v-if="can('documents.DocumentsService/ListTemplates').value"
                            :to="{ name: 'documents-templates' }"
                            icon="i-mdi-file-code"
                            truncate
                        >
                            <span class="hidden truncate sm:block">
                                {{ $t('common.template', 2) }}
                            </span>
                        </UButton>
                    </UButtonGroup>

                    <UTooltip v-if="can('documents.DocumentsService/UpdateDocument').value" :text="$t('common.create')">
                        <UButton trailing-icon="i-mdi-plus" color="neutral" truncate @click="templateModal.open({})">
                            <span class="hidden truncate sm:block">
                                {{ $t('common.document', 1) }}
                            </span>
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <UForm class="mt-2 w-full" :schema="schema" :state="query" @submit="refresh()">
                    <div class="flex flex-1 flex-row gap-2">
                        <UFormField class="flex-1" name="title" :label="$t('common.search')">
                            <UInput
                                ref="inputRef"
                                v-model="query.title"
                                type="text"
                                name="title"
                                :placeholder="$t('common.title')"
                                class="w-full"
                                leading-icon="i-mdi-search"
                            >
                                <template #trailing>
                                    <UKbd value="/" />
                                </template>
                            </UInput>
                        </UFormField>

                        <UFormField
                            class="flex min-w-32 shrink-0 flex-col"
                            name="onlyDrafts"
                            :label="$t('common.show')"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.onlyDrafts"
                                    :items="onlyDrafts"
                                    class="w-full"
                                    value-key="value"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                >
                                    <template #default="{ modelValue }">
                                        {{
                                            modelValue === undefined
                                                ? $t('common.all_documents')
                                                : onlyDrafts.find((item) => item.value === modelValue)?.label
                                        }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>
                    </div>

                    <UAccordion
                        class="mt-2"
                        color="neutral"
                        variant="soft"
                        size="sm"
                        :items="[{ label: $t('common.advanced_search'), slot: 'search' as const }]"
                    >
                        <template #search>
                            <div class="flex flex-row flex-wrap gap-1">
                                <UFormField
                                    class="flex-1"
                                    name="documentIds"
                                    :label="`${$t('common.document')} ${$t('common.id')}`"
                                >
                                    <UInput
                                        v-model="query.documentIds"
                                        type="text"
                                        name="documentIds"
                                        placeholder="DOC-..."
                                        class="w-full"
                                    />
                                </UFormField>

                                <UFormField class="flex-1" name="category" :label="$t('common.category', 1)">
                                    <ClientOnly>
                                        <USelectMenu
                                            v-model="query.categories"
                                            multiple
                                            :filter-fields="['name']"
                                            class="w-full"
                                            :searchable="
                                                async (search: string) => {
                                                    try {
                                                        categoriesLoading = true;
                                                        const categories =
                                                            await completorStore.completeDocumentCategories(search);
                                                        categoriesLoading = false;
                                                        return categories;
                                                    } catch (e) {
                                                        handleGRPCError(e as RpcError);
                                                        throw e;
                                                    } finally {
                                                        categoriesLoading = false;
                                                    }
                                                }
                                            "
                                            :searchable-placeholder="$t('common.category', 1)"
                                            value-key="id"
                                        >
                                            <template #item-label="{ item }">
                                                <div v-if="item.length > 0" class="inline-flex gap-1">
                                                    <template v-for="category in item" :key="category.id">
                                                        <DocumentCategoryBadge :category="category" />
                                                    </template>
                                                </div>
                                                <span v-else> &nbsp; </span>
                                            </template>

                                            <template #item="{ item }">
                                                <DocumentCategoryBadge :category="item" />
                                            </template>

                                            <template #empty>
                                                {{ $t('common.not_found', [$t('common.category', 2)]) }}
                                            </template>
                                        </USelectMenu>
                                    </ClientOnly>
                                </UFormField>

                                <UFormField class="flex-1" name="creator" :label="$t('common.creator')">
                                    <ClientOnly>
                                        <USelectMenu
                                            v-model="query.creators"
                                            multiple
                                            nullable
                                            class="w-full"
                                            :searchable="
                                                async (q: string): Promise<UserShort[]> => {
                                                    usersLoading = true;
                                                    const users = await completorStore.completeCitizens({
                                                        search: q,
                                                        userIds: query.creators,
                                                    });
                                                    usersLoading = false;
                                                    return users;
                                                }
                                            "
                                            :search-input="{ placeholder: $t('common.search_field') }"
                                            :filter-fields="['firstname', 'lastname']"
                                            :placeholder="$t('common.creator')"
                                            trailing
                                            value-key="userId"
                                        >
                                            <template #item-label="{ item }">
                                                <template v-if="item.length">
                                                    {{ usersToLabel(item) }}
                                                </template>
                                            </template>

                                            <template #item="{ item }">
                                                {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                                            </template>

                                            <template #empty>
                                                {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                                            </template>
                                        </USelectMenu>
                                    </ClientOnly>
                                </UFormField>
                            </div>

                            <div class="flex flex-row flex-wrap gap-2">
                                <UFormField class="flex-1" name="closed" :label="$t('common.close', 2)">
                                    <ClientOnly>
                                        <USelectMenu
                                            v-model="query.closed"
                                            :items="openclose"
                                            value-key="value"
                                            class="w-full"
                                            :search-input="{ placeholder: $t('common.search_field') }"
                                        >
                                            <template #item-label>
                                                <div class="inline-flex items-center gap-1 truncate">
                                                    <template v-if="typeof query.closed === 'boolean'">
                                                        <UIcon
                                                            v-if="!query.closed"
                                                            class="size-4"
                                                            name="i-mdi-lock-open-variant"
                                                            color="green"
                                                        />
                                                        <UIcon v-else class="size-4" name="i-mdi-lock" color="error" />
                                                    </template>

                                                    {{
                                                        query.closed === undefined
                                                            ? openclose[0]!.label
                                                            : (openclose.findLast((o) => o.value === query.closed)?.label ??
                                                              $t('common.na'))
                                                    }}
                                                </div>
                                            </template>

                                            <template #item="{ item }">
                                                <div class="inline-flex items-center gap-1 truncate">
                                                    <template v-if="typeof item.closed === 'boolean'">
                                                        <UIcon
                                                            v-if="!item.closed"
                                                            class="size-4"
                                                            name="i-mdi-lock-open-variant"
                                                            color="green"
                                                        />
                                                        <UIcon v-else class="size-4" name="i-mdi-lock" color="error" />
                                                    </template>

                                                    {{ item.label }}
                                                </div>
                                            </template>
                                        </USelectMenu>
                                    </ClientOnly>
                                </UFormField>

                                <UFormField class="flex-1" name="date" :label="$t('common.time_range')">
                                    <DateRangePickerPopoverClient
                                        v-model="query.date"
                                        class="flex-1"
                                        date-format="dd.MM.yyyy HH:mm"
                                        :popover="{ class: 'flex-1' }"
                                        :date-picker="{
                                            mode: 'dateTime',
                                            disabledDates: [{ start: addDays(new Date(), 1), end: null }],
                                            is24Hr: true,
                                            clearable: true,
                                        }"
                                    />
                                </UFormField>

                                <UFormField class="flex-1 grow-0 basis-40" :label="$t('common.sort_by')">
                                    <SortButton
                                        v-model="query.sorting"
                                        :fields="[
                                            { label: $t('common.created_at'), value: 'createdAt' },
                                            { label: $t('common.title'), value: 'title' },
                                        ]"
                                    />
                                </UFormField>
                            </div>
                        </template>
                    </UAccordion>
                </UForm>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('common.document', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="data?.documents.length === 0" :type="$t('common.document', 2)" />

            <ul
                v-else-if="data?.documents || isRequestPending(status)"
                class="min-w-full divide-y divide-default overflow-clip"
                :class="design.documents.listStyle === 'double' ? '2xl:grid 2xl:grid-cols-2' : ''"
                role="list"
            >
                <template v-if="isRequestPending(status)">
                    <li v-for="idx in 8" :key="idx" class="flex-initial">
                        <div class="m-2">
                            <div class="flex flex-row gap-2 truncate">
                                <div class="flex flex-1 flex-row items-center justify-start">
                                    <USkeleton class="h-7 w-full max-w-[125px]" />
                                </div>

                                <USkeleton class="h-7 w-full max-w-[125px]" />

                                <div class="flex flex-1 flex-row items-center justify-end gap-1">
                                    <USkeleton class="h-7 w-full max-w-[125px]" />
                                </div>
                            </div>

                            <div class="flex flex-row gap-2 truncate">
                                <div class="inline-flex items-center gap-1 truncate">
                                    <h2 class="truncate py-2 pr-3 text-xl font-medium">
                                        <USkeleton class="h-7 w-full max-w-[650px]" />
                                    </h2>
                                </div>

                                <div class="flex flex-1 flex-row items-center justify-end">
                                    <USkeleton class="h-6 w-full max-w-[250px]" />
                                </div>
                            </div>

                            <div class="flex flex-row gap-2">
                                <div class="flex flex-1 flex-row items-center justify-start">
                                    <USkeleton class="h-6 w-full max-w-[150px]" />
                                </div>

                                <div class="flex flex-1 flex-row items-center justify-center">
                                    <USkeleton class="h-6 w-full max-w-[150px]" />
                                </div>

                                <div class="flex flex-1 flex-row items-center justify-end">
                                    <USkeleton class="h-6 w-full max-w-[250px]" />
                                </div>
                            </div>
                        </div>
                    </li>
                </template>

                <template v-else>
                    <DocumentListEntry v-for="doc in data?.documents" :key="doc.id" :document="doc" />
                </template>
            </ul>
            <UContextMenu :items="links"> </UContextMenu>
        </template>

        <template #footer>
            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
        </template>
    </UDashboardPanel>

    <UDashboardPanel id="documents-pinnedlist" class="overflow-x-hidden" resizable :width="15" :min-size="15" :max-size="40">
        <!-- TODO the sidebar should be closeable -->
        <PinnedDocumentList />
    </UDashboardPanel>
</template>
