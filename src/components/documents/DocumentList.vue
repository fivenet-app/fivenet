<script lang="ts" setup>
import { watchDebounced } from '@vueuse/shared';
import { format } from 'date-fns';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useCompletorStore } from '~/store/completor';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import { Category } from '~~/gen/ts/resources/documents/category';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { ListDocumentsRequest, ListDocumentsResponse } from '~~/gen/ts/services/docstore/docstore';
import DocumentListEntry from '~/components/documents/DocumentListEntry.vue';
import DatePicker from '~/components/partials/DatePicker.vue';

const { $grpc } = useNuxtApp();

const completorStore = useCompletorStore();

const { t } = useI18n();

type OpenClose = { id: number; label: string; closed?: boolean };
const openclose: OpenClose[] = [
    { id: 0, label: t('common.not_selected') },
    { id: 1, label: t('common.open', 2), closed: false },
    { id: 2, label: t('common.close', 2), closed: true },
];

const query = ref<{
    documentIds?: string;
    title?: string;
    creator?: UserShort;
    from?: Date;
    to?: Date;
    closed?: boolean;
    category?: Category;
}>({});

const queryClosed = ref<OpenClose>(openclose[0]);

const usersLoading = ref(false);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending, refresh, error } = useLazyAsyncData(`documents-${page.value}`, () => listDocuments());

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
    if (query.value.creator) {
        req.creatorIds.push(query.value.creator.userId);
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
        const call = $grpc.getDocStoreClient().listDocuments(req);
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), { debounce: 600, maxWait: 1400 });

watch(queryClosed, () => (query.value.closed = queryClosed.value.closed));
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm class="w-full" :state="{}" @submit="refresh()">
                <UFormGroup name="search" :label="$t('common.search')">
                    <UInput
                        v-model="query.title"
                        type="text"
                        name="search"
                        :placeholder="$t('common.title')"
                        block
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
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
                                    name="search"
                                    placeholder="DOC-..."
                                    block
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>

                            <UFormGroup class="flex-1" name="category" :label="$t('common.category', 1)">
                                <UInputMenu
                                    v-model="query.category"
                                    option-attribute="name"
                                    :search-attributes="['name']"
                                    block
                                    nullable
                                    :search="completorStore.completeDocumentCategories"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                >
                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>
                                    <template #empty> {{ $t('common.not_found', [$t('common.category', 2)]) }} </template>
                                </UInputMenu>
                            </UFormGroup>

                            <UFormGroup class="flex-1" name="creator" :label="$t('common.creator')">
                                <UInputMenu
                                    v-model="query.creator"
                                    :search="
                                        async (query: string) => {
                                            usersLoading = true;
                                            const users = await completorStore.completeCitizens({
                                                search: query,
                                            });
                                            usersLoading = false;
                                            return users;
                                        }
                                    "
                                    :search-attributes="['firstname', 'lastname']"
                                    block
                                    :placeholder="
                                        query.creator
                                            ? `${query.creator?.firstname} ${query.creator?.lastname} (${query.creator?.dateofbirth})`
                                            : $t('common.owner')
                                    "
                                    trailing
                                    by="userId"
                                >
                                    <template #option="{ option: user }">
                                        {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                    </template>
                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>
                                    <template #empty> {{ $t('common.not_found', [$t('common.creator', 2)]) }} </template>
                                </UInputMenu>
                            </UFormGroup>
                        </div>
                        <div class="flex flex-row flex-wrap gap-2">
                            <UFormGroup class="flex-1" :label="$t('common.close', 2)">
                                <USelectMenu v-model="queryClosed" :options="openclose" />
                            </UFormGroup>

                            <UFormGroup class="flex-1" name="from" :label="`${$t('common.time_range')} ${$t('common.from')}`">
                                <UPopover :popper="{ placement: 'bottom-start' }">
                                    <UButton
                                        variant="outline"
                                        color="gray"
                                        block
                                        icon="i-heroicons-calendar-days-20-solid"
                                        :label="query.from ? format(query.from, 'dd.MM.yyyy') : 'dd.mm.yyyy'"
                                    />

                                    <template #panel="{ close }">
                                        <DatePicker v-model="query.from" @close="close" />
                                    </template>
                                </UPopover>
                            </UFormGroup>
                            <UFormGroup class="flex-1" name="to" :label="`${$t('common.time_range')} ${$t('common.to')}`">
                                <UPopover :popper="{ placement: 'bottom-start' }">
                                    <UButton
                                        variant="outline"
                                        color="gray"
                                        block
                                        icon="i-heroicons-calendar-days-20-solid"
                                        :label="query.to ? format(query.to, 'dd.MM.yyyy') : 'dd.mm.yyyy'"
                                    />

                                    <template #panel="{ close }">
                                        <DatePicker v-model="query.to" @close="close" />
                                    </template>
                                </UPopover>
                            </UFormGroup>
                        </div>
                    </template>
                </UAccordion>
            </UForm>
        </template>
    </UDashboardToolbar>

    <div class="inline-block w-full max-w-full px-1 py-2 align-middle">
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.document', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.document', 2)])" :retry="refresh" />
        <DataNoDataBlock v-else-if="data?.documents.length === 0" :type="$t('common.document', 2)" />
        <template v-else>
            <ul role="list" class="flex flex-col">
                <DocumentListEntry v-for="doc in data?.documents" :key="doc.id" :doc="doc" />
            </ul>
        </template>

        <div class="flex justify-end border-t border-gray-200 px-3 py-3.5 dark:border-gray-700">
            <UPagination
                v-model="page"
                :page-count="data?.pagination?.pageSize ?? 0"
                :total="data?.pagination?.totalCount ?? 0"
            />
        </div>
    </div>
</template>
