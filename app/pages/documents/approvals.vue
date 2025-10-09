<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { getDocumentsApprovalClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { ApprovalTaskStatus } from '~~/gen/ts/resources/documents/approval';
import type { ListApprovalTasksInboxResponse } from '~~/gen/ts/services/documents/approval';

useHead({
    title: 'pages.documents.approvals.title',
});

definePageMeta({
    title: 'pages.documents.approvals.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService/ListDocuments',
});

defineProps<{
    itemsLeft: NavigationMenuItem[];
    itemsRight: NavigationMenuItem[];
}>();

const approvalClient = await getDocumentsApprovalClient();

const statuses = computed(() => [
    { label: 'enums.documents.ApprovalTaskStatus.PENDING', value: ApprovalTaskStatus.PENDING },
    { label: 'enums.documents.ApprovalTaskStatus.APPROVED', value: ApprovalTaskStatus.APPROVED },
    { label: 'enums.documents.ApprovalTaskStatus.DECLINED', value: ApprovalTaskStatus.DECLINED },
    { label: 'enums.documents.ApprovalTaskStatus.EXPIRED', value: ApprovalTaskStatus.EXPIRED },
    { label: 'enums.documents.ApprovalTaskStatus.CANCELLED', value: ApprovalTaskStatus.CANCELLED },
]);

const schema = z.object({
    statuses: z.enum(ApprovalTaskStatus).array().optional().default([ApprovalTaskStatus.PENDING, ApprovalTaskStatus.CANCELLED]),
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
    () => `documents-approvals-${JSON.stringify(query)}`,
    () => listApprovalTasksInbox(),
);

async function listApprovalTasksInbox(): Promise<ListApprovalTasksInboxResponse> {
    const call = approvalClient.listApprovalTasksInbox({
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
            <UDashboardNavbar :title="$t('pages.documents.approvals.title')">
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
                        <UFormField
                            class="flex min-w-40 shrink-0 flex-col"
                            name="onlyDrafts"
                            :label="$t('common.status')"
                            :ui="{ container: 'flex-1 flex' }"
                        >
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

            <div v-else>
                Approvals Page - TODO

                {{ data }}
            </div>
        </template>

        <template v-if="itemsLeft.length > 1 || itemsRight.length > 1" #footer>
            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />

            <USeparator />

            <UDashboardToolbar>
                <div class="flex min-w-0 flex-1 flex-col justify-between gap-2 lg:flex-row">
                    <UNavigationMenu v-if="itemsLeft.length > 0" orientation="horizontal" :items="itemsLeft" class="-mx-1" />

                    <UNavigationMenu v-if="itemsRight.length > 0" orientation="horizontal" :items="itemsRight" class="-mx-1" />
                </div>
            </UDashboardToolbar>
        </template>
    </UDashboardPanel>
</template>
