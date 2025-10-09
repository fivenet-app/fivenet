<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
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

const signingClient = await getDocumentsSigningClient();

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
                Signatures Page - TODO

                {{ data }}
            </div>
        </template>

        <template v-if="itemsLeft.length > 1 || itemsRight.length > 1" #footer>
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
