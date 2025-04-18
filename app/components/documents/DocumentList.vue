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
import type { OpenClose } from '~/typings';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { ListDocumentsRequest, ListDocumentsResponse } from '~~/gen/ts/services/docstore/docstore';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const completorStore = useCompletorStore();

const settingsStore = useSettingsStore();
const { design } = storeToRefs(settingsStore);

const openclose: OpenClose[] = [
    { id: 0, label: t('common.not_selected'), closed: undefined },
    { id: 1, label: t('common.open', 2), closed: false },
    { id: 2, label: t('common.close', 2), closed: true },
];

const schema = z.object({
    documentIds: z.string().max(16).optional(),
    title: z.string().max(64).optional(),
    creators: z.custom<UserShort>().array().max(5),
    date: z
        .object({
            start: z.date(),
            end: z.date(),
        })
        .optional(),
    closed: z.boolean().optional(),
    categories: z.custom<Category>().array().max(3),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    title: '',
    date: undefined,
    creators: [],
    categories: [],
});

const usersLoading = ref(false);

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'createdAt',
    direction: 'desc',
});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`documents-${sort.value.column}:${sort.value.direction}-${page.value}`, () => listDocuments(), {
    watch: [sort],
});

async function listDocuments(): Promise<ListDocumentsResponse> {
    const req: ListDocumentsRequest = {
        pagination: {
            offset: offset.value,
        },
        sort: sort.value,
        search: query.title ?? '',
        categoryIds: query.categories.map((c) => c.id),
        creatorIds: query.creators.map((c) => c.userId),
        documentIds: [],
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

    try {
        const call = $grpc.docstore.docStore.listDocuments(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

const categoriesLoading = ref(false);

const inputRef = useTemplateRef('inputRef');

defineShortcuts({
    '/': () => inputRef.value?.input?.focus(),
});
</script>

<template>
    <UDashboardToolbar>
        <UForm :schema="schema" :state="query" class="w-full" @submit="refresh()">
            <UFormGroup name="title" :label="$t('common.search')">
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
                                >
                                    <template #label>
                                        <div v-if="query.categories.length > 0" class="inline-flex gap-1">
                                            <template v-for="category in query.categories" :key="category.id">
                                                <span class="inline-flex gap-1" :class="`bg-${category.color}-500`">
                                                    <component
                                                        :is="
                                                            availableIcons.find((item) => item.name === category?.icon) ??
                                                            fallbackIcon
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
                                                :is="availableIcons.find((item) => item.name === option.icon) ?? fallbackIcon"
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
                                        async (query: string): Promise<UserShort[]> => {
                                            usersLoading = true;
                                            const users = await completorStore.completeCitizens({
                                                search: query,
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
                                    by="userId"
                                >
                                    <template #label>
                                        <template v-if="query.creators.length">
                                            {{ usersToLabel(query.creators) }}
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
                        <UFormGroup name="closed" :label="$t('common.close', 2)" class="flex-1">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.closed"
                                    :options="openclose"
                                    value-attribute="closed"
                                    :searchable-placeholder="$t('common.search_field')"
                                >
                                    <template #label>
                                        <div class="inline-flex items-center gap-1 truncate">
                                            <template v-if="typeof query.closed === 'boolean'">
                                                <UIcon
                                                    v-if="!query.closed"
                                                    name="i-mdi-lock-open-variant"
                                                    color="green"
                                                    class="size-4"
                                                />
                                                <UIcon v-else name="i-mdi-lock" color="red" class="size-4" />
                                            </template>

                                            {{
                                                query.closed === undefined
                                                    ? openclose[0]!.label
                                                    : (openclose.findLast((o) => o.closed === query.closed)?.label ??
                                                      $t('common.na'))
                                            }}
                                        </div>
                                    </template>

                                    <template #option="{ option }">
                                        <div class="inline-flex items-center gap-1 truncate">
                                            <template v-if="typeof option.closed === 'boolean'">
                                                <UIcon
                                                    v-if="!option.closed"
                                                    name="i-mdi-lock-open-variant"
                                                    color="green"
                                                    class="size-4"
                                                />
                                                <UIcon v-else name="i-mdi-lock" color="red" class="size-4" />
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

                        <UFormGroup :label="$t('common.sort_by')" class="flex-1 grow-0 basis-40">
                            <SortButton
                                v-model="sort"
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
                role="list"
                class="my-1 flex flex-initial flex-col divide-y divide-gray-100 dark:divide-gray-800"
                :class="design.documents.listStyle === 'double' ? '2xl:grid 2xl:grid-cols-2' : ''"
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
                    <DocumentListEntry v-for="doc in data?.documents" :key="doc.id" :document="doc" />
                </template>
            </ul>
        </div>
    </UDashboardPanelContent>

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
