<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import CategoryBadge from '~/components/partials/documents/CategoryBadge.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { getDocumentsSigningClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { SignatureTaskStatus } from '~~/gen/ts/resources/documents/signing';
import type { ListSignatureTasksInboxResponse } from '~~/gen/ts/services/documents/signing';

useHead({
    title: 'pages.documents.signatures.title',
});

definePageMeta({
    title: 'pages.documents.signatures.title',
    requiresAuth: true,
    permission: 'TODOService/TODOMethod',
});

defineProps<{
    itemsLeft: NavigationMenuItem[];
    itemsRight: NavigationMenuItem[];
}>();

const { t } = useI18n();

const signingClient = await getDocumentsSigningClient();

const statuses = computed(() => [
    { label: t('enums.documents.SignatureTaskStatus.PENDING'), value: SignatureTaskStatus.PENDING },
    { label: t('enums.documents.SignatureTaskStatus.SIGNED'), value: SignatureTaskStatus.SIGNED },
    { label: t('enums.documents.SignatureTaskStatus.EXPIRED'), value: SignatureTaskStatus.EXPIRED },
]);

const schema = z.object({
    statuses: z.enum(SignatureTaskStatus).array().optional().default([SignatureTaskStatus.PENDING]),
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

const query = useSearchForm('documents-approvals', schema);

const { data, status, error, refresh } = useLazyAsyncData(
    () => `documents-signatures-${JSON.stringify(query)}`,
    () => listSignatureTasksInbox(),
);

async function listSignatureTasksInbox(): Promise<ListSignatureTasksInboxResponse> {
    const call = signingClient.listSignatureTasksInbox({
        pagination: {
            offset: calculateOffset(query.page, data.value?.pagination),
        },
        statuses: query.statuses,
    });
    const { response } = await call;

    return response;
}

// TODO
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.documents.signatures.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/documents" />
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
                        <UFormField class="flex flex-1 shrink-0 flex-col" name="onlyDrafts" :label="$t('common.status')">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.statuses"
                                    :items="statuses"
                                    multiple
                                    class="w-full"
                                    label-key="label"
                                    value-key="value"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                >
                                    <template #default>
                                        {{ $t('common.selected', query.statuses.length) }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>
                    </div>
                </UForm>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('common.task', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataPendingBlock v-else-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.task', 2)])" />
            <DataNoDataBlock v-else-if="data?.tasks.length === 0" :type="$t('common.task', 2)" />
            <ul
                v-else
                class="min-w-full divide-y divide-default"
                :class="isRequestPending(status) ? 'overflow-y-hidden' : ''"
                role="list"
            >
                <li
                    v-for="task in data?.tasks"
                    :key="task.id"
                    class="flex-initial p-1 hover:border-primary-500/25 hover:bg-primary-100/50 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
                >
                    <ULink
                        :to="{
                            name: 'documents-id',
                            params: { id: task.documentId },
                            hash: '#signatures',
                        }"
                    >
                        <div class="m-2 flex flex-col gap-1">
                            <div class="flex flex-row justify-between gap-2">
                                <div class="flex items-center">
                                    <IDCopyBadge
                                        :id="task.documentId"
                                        prefix="DOC"
                                        :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                                        :content="{
                                            key: 'notifications.document_view.copy_document_id.content',
                                            parameters: {},
                                        }"
                                        size="xs"
                                    />
                                </div>

                                <UBadge
                                    v-if="task.document?.meta?.state"
                                    class="inline-flex gap-1"
                                    size="md"
                                    icon="i-mdi-note-check"
                                    :label="task.document.meta.state"
                                />

                                <div
                                    v-if="task.document?.deletedAt"
                                    class="flex flex-1 flex-row items-center justify-center gap-1.5 font-bold"
                                >
                                    <UIcon class="size-4 shrink-0" name="i-mdi-delete" />
                                    {{ $t('common.deleted') }}
                                </div>

                                <UBadge
                                    v-if="task.document?.meta?.apPoliciesActive"
                                    class="inline-flex gap-1"
                                    size="md"
                                    :color="task.document?.meta?.approved ? 'info' : 'warning'"
                                    icon="i-mdi-approval"
                                    :label="task.document?.meta?.approved ? $t('common.approved') : $t('common.unapproved')"
                                />

                                <UBadge
                                    v-if="task.document?.meta?.sigPoliciesActive"
                                    class="inline-flex gap-1"
                                    size="md"
                                    :color="task.document?.meta?.signed ? 'info' : 'warning'"
                                    icon="i-mdi-approval"
                                    :label="task.document?.meta?.signed ? $t('common.signed') : $t('common.unsigned')"
                                />

                                <div class="flex flex-row items-center gap-1">
                                    <OpenClosedBadge :closed="task.document?.meta?.closed" />
                                </div>
                            </div>

                            <div class="flex items-center gap-2 text-highlighted">
                                <TaskStatusBadge :status="task.status" size="lg" />

                                <p class="flex-1 text-lg leading-6 font-semibold text-toned">
                                    {{ task.comment || $t('common.no_comment') }}
                                </p>

                                <CitizenInfoPopover :user="task.creator" :user-id="task.creatorId" size="lg" />
                            </div>

                            <div class="flex max-w-full shrink flex-col gap-2">
                                <div class="flex flex-col gap-1 md:flex-row">
                                    <div>
                                        <CategoryBadge :category="task.document?.category" />
                                    </div>

                                    <h2
                                        class="line-clamp-2 flex-1 text-base font-medium break-words break-all text-highlighted hover:line-clamp-3 sm:text-xl md:line-clamp-1"
                                    >
                                        <span v-if="!task.document?.title" class="italic">
                                            {{ $t('common.untitled') }}
                                        </span>
                                        <span v-else>
                                            {{ task.document?.title }}
                                        </span>
                                    </h2>

                                    <UBadge
                                        v-if="task.document?.meta?.draft"
                                        class="inline-flex grow-0 gap-1 self-start"
                                        color="info"
                                        size="md"
                                        icon="i-mdi-pencil"
                                        :label="$t('common.draft')"
                                    />
                                </div>
                            </div>

                            <div class="flex justify-between gap-2">
                                <div class="flex-1" />

                                <div class="flex flex-1 flex-row items-center justify-end gap-1.5">
                                    <span>{{ task.document?.creatorJobLabel }}</span>
                                    <UIcon class="size-4 shrink-0" name="i-mdi-briefcase" />
                                </div>
                            </div>
                        </div>
                    </ULink>
                </li>
            </ul>
        </template>

        <template #footer>
            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />

            <template v-if="itemsLeft.length > 1 || itemsRight.length > 1">
                <USeparator />

                <UDashboardToolbar>
                    <div class="flex min-w-0 flex-1 flex-col justify-between gap-2 lg:flex-row">
                        <UNavigationMenu
                            v-if="itemsLeft.length > 0"
                            orientation="horizontal"
                            :items="itemsLeft"
                            class="-mx-1"
                        />

                        <UNavigationMenu
                            v-if="itemsRight.length > 0"
                            orientation="horizontal"
                            :items="itemsRight"
                            class="-mx-1"
                        />
                    </div>
                </UDashboardToolbar>
            </template>
        </template>
    </UDashboardPanel>
</template>
