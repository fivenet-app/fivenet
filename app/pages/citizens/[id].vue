<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import type { TypedRouteFromName } from '@typed-router';
import Header from '~/components/citizens/info/Header.vue';
import AddToButton from '~/components/clipboard/AddToButton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { Perms } from '~~/gen/ts/perms';
import { ObjectType } from '~~/gen/ts/resources/notifications/client_view';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { User } from '~~/gen/ts/resources/users/users';

useHead({
    title: 'pages.citizens.id.title',
});

definePageMeta({
    title: 'pages.citizens.id.title',
    requiresAuth: true,
    permission: 'citizens.CitizensService/GetUser',
    validate: async (route) => {
        route = route as TypedRouteFromName<'citizens-id'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const { t } = useI18n();

const { can } = useAuth();

const clipboardStore = useClipboardStore();

const notifications = useNotificationsStore();

const route = useRoute('citizens-id');

const citizensCitizensClient = await getCitizensCitizensClient();

const {
    data: user,
    status,
    refresh,
    error,
} = useLazyAsyncData(`citizen-${route.params.id}`, () => getUser(parseInt(route.params.id)), {
    watch: [() => route.params.id],
});

async function getUser(userId: number): Promise<User> {
    try {
        const call = citizensCitizensClient.getUser({ userId });
        const { response } = await call;

        if (response.user?.props === undefined) {
            response.user!.props = {
                userId: response.user!.userId,
            };
        }

        return response.user!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function addToClipboard(): void {
    if (!user.value) {
        return;
    }

    clipboardStore.addUser(user.value);

    notifications.add({
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: {} },
        description: { key: 'notifications.clipboard.citizen_add.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

// Handle the client update event
const { sendClientView } = useClientUpdate(ObjectType.CITIZEN, () =>
    notifications.add({
        title: { key: 'notifications.citizens.client_view_update.title', parameters: {} },
        description: { key: 'notifications.citizens.client_view_update.content', parameters: {} },
        duration: 7500,
        type: NotificationType.INFO,
        actions: [
            {
                label: { key: 'common.refresh', parameters: {} },
                icon: 'i-mdi-refresh',
                onClick: () => refresh(),
            },
        ],
    }),
);
watch(user, () => user.value && sendClientView(user.value.userId));

const items = computed<NavigationMenuItem[]>(() =>
    [
        {
            label: t('common.profile'),
            icon: 'i-mdi-account',
            permission: 'citizens.CitizensService/ListCitizens' as Perms,
            to: '/citizens/' + route.params.id,
            exact: true,
        },
        {
            label: t('common.vehicle', 2),
            icon: 'i-mdi-car',
            permission: 'vehicles.VehiclesService/ListVehicles' as Perms,
            to: '/citizens/' + route.params.id + '/vehicles',
        },
        {
            label: t('common.document', 2),
            icon: 'i-mdi-file-document-multiple',
            permission: 'documents.DocumentsService/ListUserDocuments' as Perms,
            to: '/citizens/' + route.params.id + '/documents',
        },
        {
            label: t('common.activity'),
            icon: 'i-mdi-pulse',
            to: '/citizens/' + route.params.id + '/activity',
            permission: 'citizens.CitizensService/ListUserActivity' as Perms,
        },
    ].flatMap((item) => (item.permission === undefined || can(item.permission).value ? [item] : [])),
);
</script>

<template>
    <UDashboardPanel class="shrink-0 border-b border-gray-200 lg:w-(--width) lg:border-r lg:border-b-0 dark:border-gray-800">
        <template #header>
            <UDashboardNavbar :title="$t('pages.citizens.id.title')">
                <template #right>
                    <PartialsBackButton fallback-to="/citizens" />

                    <UButton
                        icon="i-mdi-refresh"
                        :label="$t('common.refresh')"
                        :loading="isRequestPending(status)"
                        @click="() => refresh()"
                    />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar v-if="user">
                <div class="flex flex-1 flex-row items-center gap-1">
                    <Header :user="user" class="flex-1" />

                    <div>
                        <UButtonGroup v-if="user">
                            <IDCopyBadge
                                :id="user.userId"
                                prefix="CIT"
                                :title="{ key: 'notifications.citizens.copy_citizen_id.title', parameters: {} }"
                                :content="{ key: 'notifications.citizens.copy_citizen_id.content', parameters: {} }"
                            />

                            <AddToButton :title="$t('components.clipboard.clipboard_button.add')" :callback="addToClipboard" />
                        </UButtonGroup>
                    </div>
                </div>
            </UDashboardToolbar>

            <UDashboardToolbar>
                <UNavigationMenu orientation="horizontal" :items="items" class="-mx-1 flex-1" />
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.citizen', 1)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.citizen', 1)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!user" />

            <NuxtPage v-else v-model:user="user" />
        </template>
    </UDashboardPanel>
</template>
