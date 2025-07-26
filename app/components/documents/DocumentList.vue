<script lang="ts" setup>
import { watchDebounced } from '@vueuse/shared';
import { addDays } from 'date-fns';
import { z } from 'zod';
import DocumentListEntry from '~/components/documents/DocumentListEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DateRangePickerPopoverClient from '~/components/partials/DateRangePickerPopover.client.vue';
import { availableIcons, fallbackIcon } from '~/components/partials/icons';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { useCompletorStore } from '~/stores/completor';
import { useSettingsStore } from '~/stores/settings';
import type { ToggleItem } from '~/utils/types';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import type { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { ListDocumentsRequest, ListDocumentsResponse } from '~~/gen/ts/services/documents/documents';

const { t } = useI18n();

const { can, attr, isSuperuser } = useAuth();

const completorStore = useCompletorStore();

const settingsStore = useSettingsStore();
const { design } = storeToRefs(settingsStore);

const documentsDocuments = useDocumentsDocuments();

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
    sort: z.custom<TableSortable>().default({
        column: 'createdAt',
        direction: 'desc',
    }),
    page: pageNumberSchema,
});

const query = useSearchForm('documents', schema);

const usersLoading = ref(false);

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`documents-${query.sort.column}:${query.sort.direction}-${query.page}`, () => listDocuments());

async function listDocuments(): Promise<ListDocumentsResponse> {
    const req: ListDocumentsRequest = {
        pagination: {
            offset: calculateOffset(query.page, data.value?.pagination),
        },
        sort: query.sort,
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

const { x, y } = useMouse();
const { y: windowY } = useWindowScroll();

const isOpen = ref(false);
const virtualElement = ref({ getBoundingClientRect: () => ({}) });

const selectedDocument = ref<DocumentShort | undefined>(undefined);

function onContextMenu(doc: DocumentShort) {
    selectedDocument.value = doc;
    const top = unref(y) - unref(windowY);
    const left = unref(x);

    virtualElement.value.getBoundingClientRect = () => ({
        width: 0,
        height: 0,
        top,
        left,
    });

    isOpen.value = true;
}

const links = computed(() =>
    [
        [
            {
                label: t('common.open'),
                icon: 'i-mdi-eye',
                to: {
                    name: 'documents-id',
                    params: { id: selectedDocument.value?.id ?? 0 },
                },
            },
            isSuperuser.value && selectedDocument.value?.deletedAt
                ? {
                      label: t('common.restore'),
                      icon: 'i-mdi-restore',
                      to: {
                          name: 'documents-id',
                          params: { id: selectedDocument.value?.id ?? 0 },
                      },
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
    ].filter((L) => L.length > 0),
);

const inputRef = useTemplateRef('inputRef');

defineShortcuts({
    '/': () => inputRef.value?.input?.focus(),
});
</script>

<template>
    <UDashboardToolbar>
        <UForm class="w-full" :schema="schema" :state="query" @submit="refresh()">
            <div class="flex flex-1 flex-row gap-2">
                <UFormGroup class="flex-1" name="title" :label="$t('common.search')">
                    <UInput
                        ref="inputRef"
                        v-model="query.title"
                        type="text"
                        name="title"
                        :placeholder="$t('common.title')"
                        block
                        leading-icon="i-mdi-search"
                        @keydown.esc="$event.target.blur()"
                    >
                        <template #trailing>
                            <UKbd value="/" />
                        </template>
                    </UInput>
                </UFormGroup>

                <UFormGroup
                    class="flex shrink-0 grow-0 flex-col"
                    name="onlyDrafts"
                    :label="$t('common.show')"
                    :ui="{ container: 'flex-1 flex' }"
                >
                    <ClientOnly>
                        <USelectMenu
                            v-model="query.onlyDrafts"
                            :options="onlyDrafts"
                            value-attribute="value"
                            option-attribute="label"
                            :searchable-placeholder="$t('common.search_field')"
                        >
                            <template #label>
                                {{
                                    query.onlyDrafts === undefined
                                        ? onlyDrafts[0]!.label
                                        : (onlyDrafts.findLast((o) => o.value === query.onlyDrafts)?.label ?? $t('common.na'))
                                }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormGroup>
            </div>

            <UAccordion
                class="mt-2"
                color="white"
                variant="soft"
                size="sm"
                :items="[{ label: $t('common.advanced_search'), slot: 'search' }]"
            >
                <template #search>
                    <div class="flex flex-row flex-wrap gap-1">
                        <UFormGroup class="flex-1" name="documentIds" :label="`${$t('common.document')} ${$t('common.id')}`">
                            <UInput v-model="query.documentIds" type="text" name="documentIds" placeholder="DOC-..." block />
                        </UFormGroup>

                        <UFormGroup class="flex-1" name="category" :label="$t('common.category', 1)">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.categories"
                                    multiple
                                    option-attribute="name"
                                    :search-attributes="['name']"
                                    block
                                    by="name"
                                    :searchable="
                                        async (search: string) => {
                                            try {
                                                categoriesLoading = true;
                                                const categories = await completorStore.completeDocumentCategories(search);
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
                                    searchable-lazy
                                    :searchable-placeholder="$t('common.category', 1)"
                                    value-attribute="id"
                                >
                                    <template #label="{ selected }">
                                        <div v-if="selected.length > 0" class="inline-flex gap-1">
                                            <template v-for="category in selected" :key="category.id">
                                                <span class="inline-flex gap-1" :class="`bg-${category.color}-500`">
                                                    <component
                                                        :is="
                                                            availableIcons.find((item) => item.name === category?.icon)
                                                                ?.component ?? fallbackIcon.component
                                                        "
                                                        v-if="category.icon"
                                                        class="size-5"
                                                    />
                                                    <span class="truncate">{{ category.name }}</span>
                                                </span>
                                            </template>
                                        </div>
                                        <span v-else> &nbsp; </span>
                                    </template>

                                    <template #option="{ option }">
                                        <span class="inline-flex gap-1" :class="`bg-${option.color}-500`">
                                            <component
                                                :is="
                                                    availableIcons.find((item) => item.name === option.icon)?.component ??
                                                    fallbackIcon.component
                                                "
                                                v-if="option.icon"
                                                class="size-5"
                                            />
                                            <span class="truncate">{{ option.name }}</span>
                                        </span>
                                    </template>

                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>

                                    <template #empty> {{ $t('common.not_found', [$t('common.category', 2)]) }} </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormGroup>

                        <UFormGroup class="flex-1" name="creator" :label="$t('common.creator')">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.creators"
                                    multiple
                                    nullable
                                    block
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
                                    searchable-lazy
                                    :searchable-placeholder="$t('common.search_field')"
                                    :search-attributes="['firstname', 'lastname']"
                                    :placeholder="$t('common.creator')"
                                    trailing
                                    value-attribute="userId"
                                >
                                    <template #label="{ selected }">
                                        <template v-if="selected.length">
                                            {{ usersToLabel(selected) }}
                                        </template>
                                    </template>

                                    <template #option="{ option: user }">
                                        {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                    </template>

                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>

                                    <template #empty> {{ $t('common.not_found', [$t('common.creator', 2)]) }} </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormGroup>
                    </div>

                    <div class="flex flex-row flex-wrap gap-2">
                        <UFormGroup class="flex-1" name="closed" :label="$t('common.close', 2)">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.closed"
                                    :options="openclose"
                                    value-attribute="value"
                                    option-attribute="label"
                                    :searchable-placeholder="$t('common.search_field')"
                                >
                                    <template #label>
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

                                    <template #option="{ option }">
                                        <div class="inline-flex items-center gap-1 truncate">
                                            <template v-if="typeof option.closed === 'boolean'">
                                                <UIcon
                                                    v-if="!option.closed"
                                                    class="size-4"
                                                    name="i-mdi-lock-open-variant"
                                                    color="green"
                                                />
                                                <UIcon v-else class="size-4" name="i-mdi-lock" color="error" />
                                            </template>

                                            {{ option.label }}
                                        </div>
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormGroup>

                        <UFormGroup class="flex-1" name="date" :label="$t('common.time_range')">
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
                        </UFormGroup>

                        <UFormGroup class="flex-1 grow-0 basis-40" :label="$t('common.sort_by')">
                            <SortButton
                                v-model="query.sort"
                                :fields="[
                                    { label: $t('common.created_at'), value: 'createdAt' },
                                    { label: $t('common.title'), value: 'title' },
                                ]"
                            />
                        </UFormGroup>
                    </div>
                </template>
            </UAccordion>
        </UForm>
    </UDashboardToolbar>

    <UDashboardPanelContent class="p-0 sm:pb-0">
        <DataErrorBlock
            v-if="error"
            :title="$t('common.unable_to_load', [$t('common.document', 2)])"
            :error="error"
            :retry="refresh"
        />
        <DataNoDataBlock v-else-if="data?.documents.length === 0" :type="$t('common.document', 2)" />

        <div v-else-if="data?.documents || loading" class="relative overflow-x-auto">
            <ul
                class="my-1 flex flex-initial flex-col divide-y divide-gray-100 dark:divide-gray-800"
                :class="design.documents.listStyle === 'double' ? '2xl:grid 2xl:grid-cols-2' : ''"
                role="list"
            >
                <template v-if="loading">
                    <li v-for="idx in 8" :key="idx" class="flex-initial">
                        <div class="m-2">
                            <div class="flex flex-row gap-2 truncate">
                                <div class="flex flex-1 flex-row items-center justify-start">
                                    <USkeleton class="h-7 w-[125px]" />
                                </div>

                                <USkeleton class="h-7 w-[125px]" />

                                <div class="flex flex-1 flex-row items-center justify-end gap-1">
                                    <USkeleton class="h-7 w-[125px]" />
                                </div>
                            </div>

                            <div class="flex flex-row gap-2 truncate">
                                <div class="inline-flex items-center gap-1 truncate">
                                    <h2 class="truncate py-2 pr-3 text-xl font-medium">
                                        <USkeleton class="h-7 w-[650px]" />
                                    </h2>
                                </div>

                                <div class="flex flex-1 flex-row items-center justify-end">
                                    <USkeleton class="h-6 w-[250px]" />
                                </div>
                            </div>

                            <div class="flex flex-row gap-2">
                                <div class="flex flex-1 flex-row items-center justify-start">
                                    <USkeleton class="h-6 w-[150px]" />
                                </div>

                                <div class="flex flex-1 flex-row items-center justify-center">
                                    <USkeleton class="h-6 w-[150px]" />
                                </div>

                                <div class="flex flex-1 flex-row items-center justify-end">
                                    <USkeleton class="h-6 w-[250px]" />
                                </div>
                            </div>
                        </div>
                    </li>
                </template>

                <template v-else>
                    <DocumentListEntry
                        v-for="doc in data?.documents"
                        :key="doc.id"
                        :document="doc"
                        @contextmenu.prevent="onContextMenu(doc)"
                    />
                </template>
            </ul>

            <UContextMenu v-model="isOpen" :virtual-element="virtualElement">
                <UVerticalNavigation :links="links" />
            </UContextMenu>
        </div>
    </UDashboardPanelContent>

    <Pagination v-model="query.page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
