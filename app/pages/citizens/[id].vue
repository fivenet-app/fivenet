<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import type { TypedRouteFromName } from '@typed-router';
import { breakpointsTailwind } from '@vueuse/core';
import CitizenActionsPanel from '~/components/citizens/info/CitizenActionsPanel.vue';
import Header from '~/components/citizens/info/Header.vue';
import AddToButton from '~/components/clipboard/AddToButton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { Perms } from '~~/gen/ts/perms';
import { ObjectType } from '~~/gen/ts/resources/notifications/clientview/clientview';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { User } from '~~/gen/ts/resources/users/user';

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
const { open: openClipboardModal } = useClipboardModal();

const notifications = useNotificationsStore();

const route = useRoute('citizens-id');

const citizensCitizensClient = await getCitizensCitizensClient();

const {
    data: user,
    status,
    refresh,
    error,
} = useAuthedLazyAsyncData(
    'userState',
    `citizen-${route.params.id}`,
    ({ signal }) => getUser(parseInt(route.params.id), signal),
    {
        watch: [() => route.params.id],
    },
);

async function getUser(userId: number, signal: AbortSignal): Promise<User> {
    try {
        const call = citizensCitizensClient.getUser({ userId }, { abort: signal });
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

useHead({
    title: () => (user.value ? `${userToLabel(user.value)} - ${t('pages.citizens.id.title')}` : t('pages.citizens.id.title')),
});

function addToClipboard(): void {
    if (!user.value) return;

    const added = clipboardStore.addUser(user.value);

    notifications.add({
        title: {
            key: added ? 'notifications.clipboard.citizen_add.title' : 'notifications.clipboard.limit_reached.title',
            parameters: {},
        },
        description: {
            key: added ? 'notifications.clipboard.citizen_add.content' : 'notifications.clipboard.limit_reached.content',
            parameters: {},
        },
        duration: 3250,
        type: added ? NotificationType.INFO : NotificationType.WARNING,
        actions: added
            ? []
            : [
                  {
                      label: { key: 'common.open', parameters: {} },
                      icon: 'i-mdi-clipboard-list-outline',
                      onClick: () => void openClipboardModal(),
                  },
              ],
    });
}

// Handle the client update event
const { sendClientView } = useClientUpdate(ObjectType.CITIZEN, () =>
    notifications.add({
        title: { key: 'notifications.citizens.clientview_update.title', parameters: {} },
        description: { key: 'notifications.citizens.clientview_update.content', parameters: {} },
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

const breakpoints = useBreakpoints(breakpointsTailwind);
const isMobile = breakpoints.smaller('lg');

const isOpen = ref<boolean>(false);
</script>

<template>
    <UDashboardPanel :ui="{ root: 'pb-(--page-content-bottom-offset)', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.citizens.id.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/citizens" />

                    <RefreshButton :loading="isRequestPending(status)" @click="() => refresh()" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar v-if="user" class="min-w-0">
                <div class="my-2 flex min-w-0 flex-1 flex-col gap-2 lg:flex-row lg:items-center lg:gap-3">
                    <Header :user="user" @toggle-actions="isOpen = true" />

                    <div class="shrink-0 self-end lg:self-auto">
                        <UFieldGroup v-if="user">
                            <IDCopyBadge
                                :id="user.userId"
                                prefix="CIT"
                                :title="{ key: 'notifications.citizens.copy_citizen_id.title', parameters: {} }"
                                :content="{ key: 'notifications.citizens.copy_citizen_id.content', parameters: {} }"
                            />

                            <AddToButton :title="$t('components.clipboard.clipboard_button.add')" :callback="addToClipboard" />
                        </UFieldGroup>
                    </div>
                </div>
            </UDashboardToolbar>

            <UDashboardToolbar v-if="user" class="overflow-x-auto">
                <UNavigationMenu class="-mx-1 min-w-max flex-1" orientation="horizontal" :items="items" />
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
            <DataNoDataBlock
                v-else-if="!user"
                icon="i-mdi-account"
                :message="$t('common.not_found', [$t('common.citizen', 1)])"
            />

            <NuxtPage v-else v-model:user="user" @refresh="() => refresh()" />
        </template>
    </UDashboardPanel>

    <UDashboardPanel
        v-if="!isMobile"
        id="citizen-id-actions"
        class="bg-elevated/25"
        resizable
        :default-size="23"
        :min-size="18"
        :max-size="34"
        :ui="{ root: 'pb-(--page-content-bottom-offset)', body: 'gap-2 sm:gap-2' }"
    >
        <template #header>
            <UDashboardNavbar :title="$t('common.action', 2)" :toggle="false" />
        </template>

        <template #body>
            <!-- Register kbds for the citizens actions here as it will always be available not like the profile tab content -->
            <CitizenActionsPanel
                v-if="user"
                :user="user"
                register-kbds
                @update:wanted-status="user.props!.wanted = $event"
                @update:job="
                    user.job = $event.job.name;
                    user.jobLabel = $event.job.label;
                    user.jobGrade = $event.grade.grade;
                    user.jobGradeLabel = $event.grade.label;
                "
                @update:traffic-infraction-points="user.props!.trafficInfractionPoints = $event"
                @update:mug-shot="user.props!.mugshot = $event"
            />
        </template>
    </UDashboardPanel>

    <ClientOnly>
        <USlideover v-if="isMobile" v-model:open="isOpen">
            <template #content>
                <UDashboardPanel
                    id="citizens-id-actions-mobile"
                    :ui="{ root: 'min-h-full pb-(--page-content-bottom-offset)', body: 'gap-2 sm:gap-2' }"
                >
                    <template #header>
                        <UDashboardNavbar :title="$t('common.action', 2)" :ui="{ toggle: 'hidden' }">
                            <template #leading>
                                <UButton
                                    class="-ms-1.5"
                                    icon="i-mdi-close"
                                    color="neutral"
                                    variant="ghost"
                                    @click="isOpen = false"
                                />
                            </template>
                        </UDashboardNavbar>
                    </template>

                    <template #body>
                        <!-- Register kbds for the citizens actions here as it will always be available not like the profile tab content -->
                        <CitizenActionsPanel
                            v-if="user"
                            :user="user"
                            register-kbds
                            @update:wanted-status="user.props!.wanted = $event"
                            @update:job="
                                user.job = $event.job.name;
                                user.jobLabel = $event.job.label;
                                user.jobGrade = $event.grade.grade;
                                user.jobGradeLabel = $event.grade.label;
                            "
                            @update:traffic-infraction-points="user.props!.trafficInfractionPoints = $event"
                            @update:mug-shot="user.props!.mugshot = $event"
                        />
                    </template>
                </UDashboardPanel>
            </template>
        </USlideover>
    </ClientOnly>
</template>
