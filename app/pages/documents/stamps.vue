<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import { getDocumentsSigningClient } from '~~/gen/ts/clients';
import type { ListUsableStampsResponse } from '~~/gen/ts/services/documents/signing';

useHead({
    title: 'pages.documents.stamps.title',
});

definePageMeta({
    title: 'pages.documents.stamps.title',
    requiresAuth: true,
    permission: 'documents.SigningService/ListUsableStamps',
});

defineProps<{
    itemsLeft: NavigationMenuItem[];
    itemsRight: NavigationMenuItem[];
}>();

const signingClient = await getDocumentsSigningClient();

const { data, status, error, refresh } = useLazyAsyncData(
    () => `documents-approvals`,
    () => listApprovalTasks(),
);

async function listApprovalTasks(): Promise<ListUsableStampsResponse> {
    const call = signingClient.listUsableStamps({
        pagination: {
            offset: 0,
        },
    });
    const { response } = await call;

    return response;
}

// TODO
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.stamp', 2)">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right> </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div>
                Stamps Page - TODO

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
