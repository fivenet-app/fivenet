<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import { getDocumentsSigningClient } from '~~/gen/ts/clients';
import { SignatureTaskStatus } from '~~/gen/ts/resources/documents/signing';
import type { ListSignatureTasksResponse } from '~~/gen/ts/services/documents/signing';

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

const { data, status, error, refresh } = useLazyAsyncData(
    () => `documents-signatures-`,
    () => listSignatureTasks(),
);

async function listSignatureTasks(): Promise<ListSignatureTasksResponse> {
    const call = signingClient.listSignatureTasks({
        documentId: 0,
        statuses: [SignatureTaskStatus.PENDING, SignatureTaskStatus.EXPIRED],
    });
    const { response } = await call;

    return response;
}

// TODO
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.documents.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right> </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div>Signatures Page - TODO</div>
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
