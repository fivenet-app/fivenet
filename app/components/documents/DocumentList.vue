<script lang="ts" setup>
import { watchDebounced } from '@vueuse/shared';
import { z } from 'zod';
import DocumentListEntry from '~/components/documents/DocumentListEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import type { OpenClose } from '~/typings';
import { useCompletorStore } from '~/store/completor';
import { useSettingsStore } from '~/store/settings';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { ListDocumentsRequest, ListDocumentsResponse } from '~~/gen/ts/services/docstore/docstore';
import DatePickerPopoverClient from '../partials/DatePickerPopover.client.vue';

const { t } = useI18n();

const completorStore = useCompletorStore();

const settingsStore = useSettingsStore();
const { design } = storeToRefs(settingsStore);

const openclose: OpenClose[] = [
    { id: 0, label: t('common.not_selected') },
    { id: 1, label: t('common.open', 2), closed: false },
    { id: 2, label: t('common.close', 2), closed: true },
];

const schema = z.object({
    documentIds: z.string().max(16).optional(),
    title: z.string().max(64),
    creators: z.custom<UserShort>().array().max(5),
    from: z.date().optional(),
    to: z.date().optional(),
    closed: z.boolean().optional(),
    category: z.custom<Category>().optional(),
});

type Schema = z.output<typeof schema>;

const query = ref<Schema>({
    title: '',
    creators: [],
});

const usersLoading = ref(false);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending: loading, refresh, error } = useLazyAsyncData(`documents-${page.value}`, () => listDocuments());

async function listDocuments(): Promise<ListDocumentsResponse> {
    const req: ListDocumentsRequest = {
        pagination: {
            offset: offset.value,
        },
        orderBy: [],
        search: query.value.title ?? '',
        categoryIds: [],
        creatorIds: [],
        documentIds: [],
    };
    if (query.value.category) {
        req.categoryIds.push(query.value.category.id);
    }
    if (query.value.creators) {
        query.value.creators.forEach((c) => req.creatorIds.push(c.userId));
    }
    if (query.value.documentIds) {
        const id = query.value.documentIds.trim().replaceAll('-', '').replace(/\D/g, '');
        if (id.length > 0) {
            req.documentIds.push(id);
        }
    }
    if (query.value.from) {
        req.from = {
            timestamp: googleProtobufTimestamp.Timestamp.fromDate(query.value.from!),
        };
    }
    if (query.value.to) {
        req.to = {
            timestamp: googleProtobufTimestamp.Timestamp.fromDate(query.value.to!),
        };
    }
    if (query.value.closed !== undefined) {
        req.closed = query.value.closed;
    }

    try {
        const call = getGRPCDocStoreClient().listDocuments(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), { debounce: 200, maxWait: 1250 });

const categoriesLoading = ref(false);

const input = ref<{ input: HTMLInputElement }>();

defineShortcuts({
    '/': () => {
        input.value?.input?.focus();
    },
});
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm :schema="schema" :state="query" class="w-full" @submit="refresh()">
                <UFormGroup name="title" :label="$t('common.search')">
                    <UInput
                        ref="input"
                        v-model="query.title"
                        type="text"
                        name="title"
                        :placeholder="$t('common.title')"
                        block
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
                            <UFormGroup
                                class="flex-1"
                                name="documentIds"
                                :label="`${$t('common.document')} ${$t('common.id')}`"
                            >
                                <UInput
                                    v-model="query.documentIds"
                                    type="text"
                                    name="documentIds"
                                    placeholder="DOC-..."
                                    block
                                />
                            </UFormGroup>

                            <UFormGroup class="flex-1" name="category" :label="$t('common.category', 1)">
                                <UInputMenu
                                    v-model="query.category"
                                    nullable
                                    option-attribute="name"
                                    :search-attributes="['name']"
                                    block
                                    by="name"
                                    :search="
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
                                    search-lazy
                                    :search-placeholder="$t('common.category', 1)"
                                >
                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>
                                    <template #empty> {{ $t('common.not_found', [$t('common.category', 2)]) }} </template>
                                </UInputMenu>
                            </UFormGroup>

                            <UFormGroup class="flex-1" name="creator" :label="$t('common.creator')">
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
                            </UFormGroup>
                        </div>

                        <div class="flex flex-row flex-wrap gap-2">
                            <UFormGroup name="closed" :label="$t('common.close', 2)" class="flex-1">
                                <USelectMenu
                                    v-model="query.closed"
                                    :options="openclose"
                                    value-attribute="closed"
                                    :searchable-placeholder="$t('common.search_field')"
                                />
                            </UFormGroup>

                            <UFormGroup class="flex-1" name="from" :label="`${$t('common.time_range')} ${$t('common.from')}`">
                                <DatePickerPopoverClient
                                    v-model="query.from"
                                    :popover="{ popper: { placement: 'bottom-start' } }"
                                    :date-picker="{ mode: 'dateTime', is24hr: true, clearable: true }"
                                />
                            </UFormGroup>

                            <UFormGroup class="flex-1" name="to" :label="`${$t('common.time_range')} ${$t('common.to')}`">
                                <DatePickerPopoverClient
                                    v-model="query.to"
                                    :popover="{ popper: { placement: 'bottom-start' } }"
                                    :date-picker="{ mode: 'dateTime', is24hr: true, clearable: true }"
                                />
                            </UFormGroup>
                        </div>
                    </template>
                </UAccordion>
            </UForm>
        </template>
    </UDashboardToolbar>

    <UDashboardPanelContent class="p-0">
        <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.document', 2)])" :retry="refresh" />
        <DataNoDataBlock v-else-if="data?.documents.length === 0" :type="$t('common.document', 2)" />

        <div v-else-if="data?.documents || loading" class="relative overflow-x-auto">
            <ul
                role="list"
                class="my-1 flex flex-initial flex-col divide-y divide-gray-100 dark:divide-gray-800"
                :class="design.documents.listStyle === 'double' ? '2xl:grid 2xl:grid-cols-2' : ''"
            >
                <li v-for="_ in 8" v-if="loading" class="flex-initial">
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

                <DocumentListEntry v-for="doc in data?.documents" v-else :key="doc.id" :document="doc" />
            </ul>
        </div>
    </UDashboardPanelContent>

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
