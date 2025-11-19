<script lang="ts" setup>
import { fromDate, getLocalTimeZone } from '@internationalized/date';
import { watchDebounced } from '@vueuse/shared';
import { addDays } from 'date-fns';
import { z } from 'zod';
import ListEntry from '~/components/documents/ListEntry.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import InputDateRangePopover from '~/components/partials/InputDateRangePopover.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { useCompletorStore } from '~/stores/completor';
import { useSettingsStore } from '~/stores/settings';
import type { ToggleItem } from '~/utils/types';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { ListDocumentsRequest, ListDocumentsResponse } from '~~/gen/ts/services/documents/documents';
import CategoryBadge from '../partials/documents/CategoryBadge.vue';
import SelectMenu from '../partials/SelectMenu.vue';
import PinnedList from './PinnedList.vue';
import TemplateModal from './templates/TemplateModal.vue';
import { breakpointsTailwind } from '@vueuse/core';

const { t } = useI18n();

const { can } = useAuth();

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
    documentIds: z.coerce.string().max(16).optional(),
    title: z.coerce
        .string()
        .max(64)
        .optional()
        .default('')
        .transform((val) => val.slice(0, 64)),
    creators: z.coerce.number().array().max(5).default([]),
    date: z
        .object({
            start: z.coerce.date(),
            end: z.coerce.date().max(addDays(new Date(), 2)),
        })
        .optional(),
    closed: z.coerce.boolean().optional(),
    categories: z.coerce.number().array().max(3).default([]),
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

const { data, status, refresh, error } = useLazyAsyncData(`documents-${JSON.stringify(query.sorting)}-${query.page}`, () =>
    listDocuments(),
);

async function listDocuments(): Promise<ListDocumentsResponse> {
    const pagination = {
        offset: 0,
        pageSize: 16,
        end: 0,
        totalCount: query.page * 16,

        ...data.value?.pagination,
    };

    const req: ListDocumentsRequest = {
        pagination: {
            offset: calculateOffset(query.page, pagination),
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

    try {
        return await documentsDocuments.listDocuments(req);
    } catch (e) {
        handleGRPCError(e);
        throw e;
    }
}

const formRef = useTemplateRef('formRef');

watchDebounced(query, async () => (await formRef.value?.validate({})) && refresh(), { debounce: 200, maxWait: 1250 });

const isPinnedDocumentsVisible = ref(false);

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
                    <UButton
                        class="2xl:hidden"
                        trailing-icon="i-mdi-pin"
                        color="neutral"
                        truncate
                        @click="isPinnedDocumentsVisible = !isPinnedDocumentsVisible"
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.pinned') }}
                        </span>
                    </UButton>

                    <UFieldGroup class="inline-flex">
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
                    </UFieldGroup>

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
                <UForm
                    ref="formRef"
                    class="my-2 flex w-full flex-1 flex-col gap-2"
                    :schema="schema"
                    :state="query"
                    @submit="refresh()"
                >
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
                            class="flex min-w-40 shrink-0 flex-col"
                            name="onlyDrafts"
                            :label="$t('common.show')"
                            :ui="{ container: 'flex-1 flex' }"
                        >
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.onlyDrafts"
                                    :items="onlyDrafts"
                                    class="w-full"
                                    label-key="label"
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

                    <UCollapsible>
                        <UButton
                            class="group"
                            color="neutral"
                            variant="ghost"
                            trailing-icon="i-mdi-chevron-down"
                            :label="$t('common.advanced_search')"
                            :ui="{
                                trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                            }"
                            block
                        />

                        <template #content>
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
                                    <SelectMenu
                                        v-model="query.categories"
                                        multiple
                                        :filter-fields="['name']"
                                        class="w-full"
                                        :searchable="
                                            async (search: string) => {
                                                try {
                                                    return await completorStore.completeDocumentCategories(search);
                                                } catch (e) {
                                                    handleGRPCError(e as RpcError);
                                                    throw e;
                                                }
                                            }
                                        "
                                        searchable-key="completor-document-categories"
                                        :search-input="{ placeholder: $t('common.category', 1) }"
                                        value-key="id"
                                    >
                                        <template v-if="query.categories" #default="{ items }">
                                            <CategoryBadge
                                                v-for="category in query.categories"
                                                :key="category"
                                                :category="items?.find((c) => c.id === category)"
                                            />
                                        </template>

                                        <template #item-label="{ item }">
                                            <CategoryBadge :category="item" />
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.category', 2)]) }}
                                        </template>
                                    </SelectMenu>
                                </UFormField>

                                <UFormField class="flex-1" name="creator" :label="$t('common.creator')">
                                    <SelectMenu
                                        v-model="query.creators"
                                        multiple
                                        nullable
                                        class="w-full"
                                        :searchable="
                                            async (q: string): Promise<UserShort[]> =>
                                                await completorStore.completeCitizens({
                                                    search: q,
                                                    userIds: query.creators,
                                                })
                                        "
                                        searchable-key="completor-citizens"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                        :filter-fields="['firstname', 'lastname']"
                                        :placeholder="$t('common.creator')"
                                        trailing
                                        value-key="userId"
                                    >
                                        <template #item-label="{ item }">
                                            {{ userToLabel(item) }}
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                                        </template>
                                    </SelectMenu>
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
                                            <template #default>
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

                                            <template #item-label="{ item }">
                                                <div class="inline-flex items-center gap-1 truncate">
                                                    <template v-if="typeof item.value === 'boolean'">
                                                        <UIcon
                                                            v-if="!item.value"
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
                                    <InputDateRangePopover
                                        v-model="query.date"
                                        class="flex-1"
                                        :max-value="fromDate(addDays(new Date(), 1), getLocalTimeZone())"
                                        clearable
                                        time
                                        :range="false"
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
                    </UCollapsible>
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
                class="min-w-full divide-y divide-default"
                :class="[
                    design.documents.listStyle === 'double' ? '2xl:grid 2xl:grid-cols-2' : '',
                    isRequestPending(status) ? 'overflow-y-hidden' : '',
                ]"
                role="list"
            >
                <template v-if="isRequestPending(status)">
                    <li v-for="idx in 10" :key="idx" class="flex-initial p-1">
                        <div class="m-2 flex flex-col gap-1">
                            <USkeleton class="h-[129px] w-full" />
                        </div>
                    </li>
                </template>

                <template v-else>
                    <ListEntry v-for="doc in data?.documents" :key="doc.id" :document="doc" />
                </template>
            </ul>
        </template>

        <template #footer>
            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />

            <slot name="footer" />
        </template>
    </UDashboardPanel>

    <PinnedList v-model:open="isPinnedDocumentsVisible" />
</template>
