<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getDocumentsStampsClient } from '~~/gen/ts/clients';
import type { ListUsableStampsResponse } from '~~/gen/ts/services/documents/stamps';

useHead({
    title: 'pages.documents.stamps.title',
});

definePageMeta({
    title: 'pages.documents.stamps.title',
    requiresAuth: true,
    permission: 'documents.StampsService/ListUsableStamps',
});

defineProps<{
    itemsLeft: NavigationMenuItem[];
    itemsRight: NavigationMenuItem[];
}>();

const { can } = useAuth();

const stampsClient = await getDocumentsStampsClient();

const { data, status, error, refresh } = useLazyAsyncData(`documents-approvals-${JSON.stringify({})}`, () =>
    listApprovalTasks(),
);

async function listApprovalTasks(): Promise<ListUsableStampsResponse> {
    const call = stampsClient.listUsableStamps({
        pagination: {
            offset: 0,
        },
    });
    const { response } = await call;

    return response;
}

// TODO add logic for creating/updating stamps
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.stamp', 2)">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/documents" />

                    <UTooltip v-if="can('documents.StampsService/UpsertStamp').value" :text="$t('common.coming_soon')">
                        <UButton trailing-icon="i-mdi-plus" color="neutral" truncate disabled>
                            <span class="hidden truncate sm:block">
                                {{ $t('common.stamp', 1) }}
                            </span>
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('common.stamp', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataPendingBlock v-else-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.task', 2)])" />
            <DataNoDataBlock v-else-if="data?.stamps.length === 0" :type="$t('common.stamp', 2)" icon="i-mdi-stamper" />

            <div v-else class="flex justify-center">
                <UPageGrid>
                    <UPageCard v-for="stamp in data?.stamps" :key="stamp.id">
                        <template #title>
                            <span>{{ stamp.name }}</span>

                            <UTooltip
                                v-if="can('documents.StampsService/UpsertStamp').value"
                                class="flex-1"
                                :text="$t('common.edit')"
                            >
                                <UButton class="flex-1" block color="neutral" variant="ghost" icon="i-mdi-pencil" />
                            </UTooltip>
                        </template>

                        <template #description>
                            <!-- eslint-disable-next-line vue/no-v-html -->
                            <div v-html="stamp.svgTemplate" />
                        </template>
                    </UPageCard>
                </UPageGrid>
            </div>
        </template>

        <template v-if="itemsLeft.length > 1 || itemsRight.length > 1" #footer>
            <USeparator />

            <UDashboardToolbar>
                <div class="flex min-w-0 flex-1 flex-row flex-wrap justify-between gap-2">
                    <UNavigationMenu v-if="itemsLeft.length > 0" orientation="horizontal" :items="itemsLeft" class="-mx-1" />

                    <UNavigationMenu v-if="itemsRight.length > 0" orientation="horizontal" :items="itemsRight" class="-mx-1" />
                </div>
            </UDashboardToolbar>
        </template>
    </UDashboardPanel>
</template>
