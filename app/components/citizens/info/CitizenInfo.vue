<script lang="ts" setup>
import type { TabItem } from '#ui/types';
import CitizenActivityFeed from '~/components/citizens/info/CitizenActivityFeed.vue';
import CitizenDocuments from '~/components/citizens/info/CitizenDocuments.vue';
import CitizenProfile from '~/components/citizens/info/CitizenProfile.vue';
import CitizenVehicles from '~/components/citizens/info/CitizenVehicles.vue';
import AddToButton from '~/components/clipboard/AddToButton.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { useClipboardStore } from '~/stores/clipboard';
import { useNotificatorStore } from '~/stores/notificator';
import type { Perms } from '~~/gen/ts/perms';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { User } from '~~/gen/ts/resources/users/users';
import CitizenActions from './CitizenActions.vue';
import CitizenSetLabels from './props/CitizenSetLabels.vue';

const props = defineProps<{
    userId: number;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { attr, can } = useAuth();

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const items: TabItem[] = [
    {
        slot: 'profile',
        label: t('common.profile'),
        icon: 'i-mdi-account',
        permission: 'CitizenStoreService.ListCitizens' as Perms,
    },
    {
        slot: 'vehicles',
        label: t('common.vehicle', 2),
        icon: 'i-mdi-car',
        permission: 'DMVService.ListVehicles' as Perms,
    },
    {
        slot: 'documents',
        label: t('common.document', 2),
        icon: 'i-mdi-file-document-multiple',
        permission: 'DocStoreService.ListUserDocuments' as Perms,
    },
    {
        slot: 'activity',
        label: t('common.activity'),
        icon: 'i-mdi-pulse',
        permission: 'CitizenStoreService.ListUserActivity' as Perms,
    },
].flatMap((item) => (can(item.permission).value ? [item] : []));

const {
    data: user,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`citizen-${props.userId}`, () => getUser(props.userId));

async function getUser(userId: number): Promise<User> {
    try {
        const call = $grpc.citizenstore.citizenStore.getUser({ userId });
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
        timeout: 3250,
        type: NotificationType.INFO,
    });
}

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value]?.slot }, hash: '#' });
    },
});

const { game } = useAppConfig();

const isOpen = ref(false);
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel
            class="shrink-0 border-b border-gray-200 lg:w-[--width] lg:border-b-0 lg:border-r dark:border-gray-800"
            grow
        >
            <UDashboardNavbar :title="$t('pages.citizens.id.title')">
                <template #right>
                    <PartialsBackButton fallback-to="/citizens" />

                    <UButtonGroup class="inline-flex lg:hidden">
                        <IDCopyBadge
                            :id="userId"
                            prefix="CIT"
                            :title="{ key: 'notifications.citizen_info.copy_citizen_id.title', parameters: {} }"
                            :content="{ key: 'notifications.citizen_info.copy_citizen_id.content', parameters: {} }"
                        />

                        <AddToButton :title="$t('components.clipboard.clipboard_button.add')" :callback="addToClipboard" />
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardPanelContent>
                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.citizen', 1)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.citizen', 1)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!user" />

                <div v-else>
                    <div class="mb-4 flex items-center gap-2 px-4">
                        <ProfilePictureImg
                            :src="user?.props?.mugShot?.url"
                            :name="`${user.firstname} ${user.lastname}`"
                            :alt="$t('common.mug_shot')"
                            :enable-popup="true"
                            size="3xl"
                        />

                        <div class="w-full flex-1">
                            <div class="flex snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                                <h1 class="flex-1 break-words px-0.5 py-1 text-4xl font-bold sm:pl-1">
                                    {{ user?.firstname }} {{ user?.lastname }}
                                </h1>
                            </div>

                            <div class="inline-flex gap-2">
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

                        <UButton class="lg:hidden" icon="i-mdi-menu" @click="isOpen = true">
                            {{ $t('common.action', 2) }}
                        </UButton>
                    </div>

                    <UTabs v-model="selectedTab" class="w-full" :items="items" :unmount="true">
                        <template #profile>
                            <UContainer>
                                <CitizenProfile :user="user" />
                            </UContainer>
                        </template>

                        <template #vehicles>
                            <UContainer>
                                <CitizenVehicles :user-id="user.userId" />
                            </UContainer>
                        </template>

                        <template #documents>
                            <UContainer>
                                <CitizenDocuments :user-id="user.userId" />
                            </UContainer>
                        </template>

                        <template #activity>
                            <UContainer>
                                <CitizenActivityFeed :user-id="user.userId" />
                            </UContainer>
                        </template>
                    </UTabs>
                </div>
            </UDashboardPanelContent>
        </UDashboardPanel>

        <UDashboardPanel v-if="user" v-model="isOpen" class="max-w-72 flex-1" collapsible side="right">
            <UDashboardNavbar>
                <template #right>
                    <UButtonGroup class="hidden lg:inline-flex">
                        <IDCopyBadge
                            :id="userId"
                            prefix="CIT"
                            :title="{ key: 'notifications.citizen_info.copy_citizen_id.title', parameters: {} }"
                            :content="{ key: 'notifications.citizen_info.copy_citizen_id.content', parameters: {} }"
                        />

                        <AddToButton :title="$t('components.clipboard.clipboard_button.add')" :callback="addToClipboard" />
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardPanelContent>
                <div class="flex flex-1 flex-col">
                    <template v-if="user">
                        <UDashboardSection
                            :ui="{
                                wrapper: 'divide-y !divide-transparent space-y-0 *:pt-2 first:*:pt-2 first:*:pt-0 mb-6',
                            }"
                            :title="$t('common.action', 2)"
                        >
                            <!-- Register shortcuts for the citizens actions here as it will always be available not like the profile tab content -->
                            <CitizenActions
                                :user="user"
                                register-shortcuts
                                @update:wanted-status="user.props!.wanted = $event"
                                @update:job="
                                    user.job = $event.job.name;
                                    user.jobLabel = $event.job.label;
                                    user.jobGrade = $event.grade.grade;
                                    user.jobGradeLabel = $event.grade.label;
                                "
                                @update:traffic-infraction-points="user.props!.trafficInfractionPoints = $event"
                                @update:mug-shot="user.props!.mugShot = $event"
                            />
                        </UDashboardSection>

                        <UDashboardSection
                            v-if="
                                can('CitizenStoreService.GetUser').value &&
                                attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.Labels').value
                            "
                            :ui="{
                                wrapper: 'divide-y !divide-transparent space-y-0 *:pt-2 first:*:pt-2 first:*:pt-0 mb-6',
                            }"
                            :title="$t('common.label', 2)"
                        >
                            <CitizenSetLabels v-model="user.props!.labels" :user-id="user.userId" />
                        </UDashboardSection>
                    </template>
                </div>
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardPage>
</template>
