<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import type { TypedRouteFromName } from '@typed-router';
import CitizenActions from '~/components/citizens/info/CitizenActions.vue';
import SetLabels from '~/components/citizens/info/props/SetLabels.vue';
import AddToButton from '~/components/clipboard/AddToButton.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { Perms } from '~~/gen/ts/perms';
import { ObjectType } from '~~/gen/ts/resources/notifications/client_view';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { User } from '~~/gen/ts/resources/users/users';

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

const { game } = useAppConfig();

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

useHead({
    title: () =>
        user.value
            ? `${user.value.firstname} ${user.value.lastname} (${user.value.dateofbirth}) - ${t('pages.citizens.id.title')}`
            : t('pages.citizens.id.title'),
});

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

const isOpen = ref(false);
</script>

<!-- eslint-disable vue/no-multiple-template-root -->
<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.citizens.id.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/citizens" />

                    <UButton
                        icon="i-mdi-refresh"
                        :label="$t('common.refresh')"
                        :loading="isRequestPending(status)"
                        :ui="{ label: 'hidden sm:inline-flex' }"
                        @click="() => refresh()"
                    />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar v-if="user">
                <div class="my-2 flex flex-1 flex-row items-center gap-1">
                    <div class="flex flex-1 items-center gap-2">
                        <ProfilePictureImg
                            :src="user?.props?.mugshot?.filePath"
                            :name="`${user.firstname} ${user.lastname}`"
                            :alt="$t('common.mugshot')"
                            :enable-popup="true"
                            size="3xl"
                            class="shrink-0"
                        />

                        <div class="flex-1">
                            <div class="flex snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                                <h2 class="flex-1 px-0.5 py-1 text-4xl font-bold break-words sm:pl-1">
                                    {{ user?.firstname }} {{ user?.lastname }}
                                </h2>
                            </div>

                            <div class="inline-flex flex-col gap-2 lg:flex-row">
                                <UBadge>
                                    {{ user.jobLabel }}
                                    <template v-if="user.job !== game.unemployedJobName">
                                        ({{ $t('common.rank') }}: {{ user.jobGradeLabel }})
                                    </template>
                                    {{ user.props?.jobName || user.props?.jobGradeNumber ? '*' : '' }}
                                </UBadge>

                                <UBadge v-if="user?.props?.wanted" color="error">
                                    {{ $t('common.wanted').toUpperCase() }}
                                </UBadge>
                            </div>
                        </div>

                        <div class="flex flex-col gap-1 sm:flex-row">
                            <UButton
                                :label="$t('common.action', 2)"
                                class="lg:hidden"
                                icon="i-mdi-menu"
                                @click="isOpen = true"
                            />
                        </div>
                    </div>

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

            <UDashboardToolbar v-if="user">
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

            <NuxtPage v-else v-model:user="user" @refresh="() => refresh()" />
        </template>
    </UDashboardPanel>

    <UDashboardSidebar id="citizen-id-actions" v-model:open="isOpen" side="right" class="bg-elevated/25">
        <template #header>
            <div class="flex min-w-0 items-center gap-1.5">
                <h1 class="flex items-center gap-1.5 truncate font-semibold text-highlighted">
                    {{ $t('common.action', 2) }}
                </h1>
            </div>
        </template>

        <!-- Register kbds for the citizens actions here as it will always be available not like the profile tab content -->
        <CitizenActions
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

        <USeparator />

        <template v-if="user">
            <div class="flex shrink-0 items-center gap-1.5">
                <div class="flex min-w-0 items-center gap-1.5">
                    <h1 class="flex items-center gap-1.5 truncate font-semibold text-highlighted">
                        {{ $t('common.label', 2) }}
                    </h1>
                </div>
            </div>

            <SetLabels v-model="user.props!.labels" :user-id="user.userId" class="flex-1" />
        </template>
    </UDashboardSidebar>
</template>
