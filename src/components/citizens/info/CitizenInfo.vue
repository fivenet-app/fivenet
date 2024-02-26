<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/vue';
import { AccountIcon, BulletinBoardIcon, CarIcon, CloseIcon, FileDocumentMultipleIcon, MenuIcon } from 'mdi-vue3';
import AddToButton from '~/components/clipboard/AddToButton.vue';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { User } from '~~/gen/ts/resources/users/users';
import CitizenActivityFeed from '~/components/citizens/info/CitizenActivityFeed.vue';
import CitizenDocuments from '~/components/citizens/info/CitizenDocuments.vue';
import CitizenProfile from '~/components/citizens/info/CitizenProfile.vue';
import CitizenVehicles from '~/components/citizens/info/CitizenVehicles.vue';
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import type { Perms } from '~~/gen/ts/perms';
import AvatarImg from '~/components/partials/citizens/AvatarImg.vue';

const props = defineProps<{
    id: string;
}>();

const { $grpc } = useNuxtApp();

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const { t } = useI18n();

const tabs = [
    {
        id: 'profile',
        name: t('common.profile'),
        icon: markRaw(AccountIcon),
        permission: 'CitizenStoreService.ListCitizens' as Perms,
    },
    {
        id: 'vehicles',
        name: t('common.vehicle', 2),
        icon: markRaw(CarIcon),
        permission: 'DMVService.ListVehicles' as Perms,
    },
    {
        id: 'documents',
        name: t('common.document', 2),
        icon: markRaw(FileDocumentMultipleIcon),
        permission: 'DocStoreService.ListUserDocuments' as Perms,
    },
    {
        id: 'activity',
        name: t('common.activity'),
        icon: markRaw(BulletinBoardIcon),
        permission: 'CitizenStoreService.ListUserActivity' as Perms,
    },
].filter((tab) => can(tab.permission));

const { data: user, pending, refresh, error } = useLazyAsyncData(`citizen-${props.id}`, () => getUser(parseInt(props.id, 10)));

async function getUser(userId: number): Promise<User> {
    try {
        const call = $grpc.getCitizenStoreClient().getUser({ userId });
        const { response } = await call;

        if (response.user?.props === undefined) {
            response.user!.props = {
                userId: response.user!.userId,
            };
        }

        return response.user!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function addToClipboard(): void {
    if (user.value === null) {
        return;
    }

    clipboardStore.addUser(user.value);

    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: {} },
        content: { key: 'notifications.clipboard.citizen_add.content', parameters: {} },
        duration: 3250,
        type: 'info',
    });
}

const selectedTab = ref(0);

function changeTab(index: number) {
    selectedTab.value = index;
}

const open = ref(false);
</script>

<template>
    <div>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.citizen', 1)])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.citizen', 1)])"
            :message="$t(error.message)"
            :retry="refresh"
        />
        <DataNoDataBlock v-else-if="user === null" />

        <template v-else>
            <ClipboardButton />
            <div class="mb-14">
                <div class="flex gap-4 my-4 px-4">
                    <div class="">
                        <AvatarImg
                            :url="user.props?.mugShot?.url"
                            :name="`${user.firstname} ${user.lastname}`"
                            size="xl"
                            :rounded="false"
                        />
                    </div>
                    <div>
                        <div class="flex snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto">
                            <h1 class="flex-1 break-words py-1 pl-0.5 pr-0.5 text-4xl font-bold text-neutral sm:pl-1">
                                {{ user?.firstname }} {{ user?.lastname }}
                            </h1>
                            <IDCopyBadge
                                :id="user.userId"
                                prefix="CIT"
                                :title="{ key: 'notifications.citizen_info.copy_citizen_id.title', parameters: {} }"
                                :content="{ key: 'notifications.citizen_info.copy_citizen_id.content', parameters: {} }"
                                class="min-h-9 self-end"
                            />
                        </div>
                        <div class="my-2 flex flex-row items-center gap-2">
                            <span
                                class="inline-flex items-center rounded-full bg-base-100 px-2.5 py-0.5 text-sm font-medium text-base-800"
                            >
                                {{ user.jobLabel }}
                                <span v-if="user.jobGrade > 0">&nbsp;({{ $t('common.rank') }}: {{ user.jobGradeLabel }})</span>
                            </span>
                            <span
                                v-if="user.props?.wanted"
                                class="inline-flex items-center rounded-full bg-error-100 px-2.5 py-0.5 text-sm font-medium text-error-700"
                            >
                                {{ $t('common.wanted').toUpperCase() }}
                            </span>
                        </div>
                    </div>
                </div>

                <nav class="-ml-2 bg-base-700 lg:rounded-lg">
                    <div class="mx-auto ml-2 max-w-7xl px-4 sm:px-6 lg:px-8">
                        <div class="flex h-16 items-center justify-between">
                            <div class="flex items-center md:overflow-x-scroll">
                                <div class="-ml-2 flex md:hidden">
                                    <!-- Mobile menu button -->
                                    <button
                                        type="button"
                                        class="relative inline-flex items-center justify-center rounded-md bg-base-500 p-2 text-accent-200 hover:bg-base-400 hover:bg-opacity-75 hover:text-neutral focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2 focus:ring-offset-base-600"
                                        @click="open = !open"
                                    >
                                        <span class="absolute -inset-0.5" />
                                        <span class="sr-only">{{ $t('components.partials.sidebar.open_navigation') }}</span>
                                        <MenuIcon v-if="!open" class="block h-5 w-5" aria-hidden="true" />
                                        <CloseIcon v-else class="block h-5 w-5" aria-hidden="true" />
                                    </button>
                                </div>
                                <div class="md:block hidden">
                                    <div class="flex items-baseline space-x-2">
                                        <template v-for="(tab, index) in tabs" :key="tab.id">
                                            <span class="flex-1">
                                                <button
                                                    type="button"
                                                    class="group flex shrink-0 items-center gap-2 rounded-md p-3 text-sm font-medium text-accent-100 hover:bg-accent-100/10 hover:text-neutral hover:transition-all"
                                                    :class="
                                                        selectedTab === index
                                                            ? 'bg-accent-100/20 font-bold text-primary-300'
                                                            : ''
                                                    "
                                                    @click="selectedTab = index"
                                                >
                                                    <component
                                                        :is="tab.icon"
                                                        :class="[
                                                            selectedTab === index ? '' : 'group-hover:text-base-300',
                                                            'h-5 w-5',
                                                        ]"
                                                        aria-hidden="true"
                                                    />
                                                    <span>
                                                        {{ tab.name }}
                                                    </span>
                                                </button>
                                            </span>
                                        </template>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="-ml-3 -mr-3 md:hidden" :class="open ? 'block' : 'hidden'">
                            <div class="space-y-1 px-2 pb-3 pt-2 sm:px-3">
                                <template v-for="(tab, index) in tabs" :key="tab.id">
                                    <button
                                        type="button"
                                        class="group flex w-full shrink-0 items-center items-center gap-2 rounded-md p-2 text-sm font-medium text-accent-100 hover:bg-accent-100/10 hover:text-neutral hover:transition-all"
                                        :class="selectedTab === index ? 'bg-accent-100/20 font-bold text-primary-300' : ''"
                                        @click="selectedTab = index"
                                    >
                                        <component
                                            :is="tab.icon"
                                            :class="[selectedTab === index ? '' : 'group-hover:text-base-300', 'h-5 w-5']"
                                            aria-hidden="true"
                                        />
                                        <span>
                                            {{ tab.name }}
                                        </span>
                                    </button>
                                </template>
                            </div>
                        </div>
                    </div>
                </nav>

                <TabGroup :selected-index="selectedTab" @change="changeTab">
                    <TabList class="hidden">
                        <Tab v-for="tab in tabs" :key="tab.id"></Tab>
                    </TabList>
                    <TabPanels class="mt-2 bg-transparent">
                        <TabPanel>
                            <CitizenProfile
                                :user="user"
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
                        </TabPanel>
                        <TabPanel v-if="can('DMVService.ListVehicles')">
                            <CitizenVehicles :user-id="user.userId" />
                        </TabPanel>
                        <TabPanel v-if="can('DocStoreService.ListUserDocuments')">
                            <CitizenDocuments :user-id="user.userId" />
                        </TabPanel>
                        <TabPanel v-if="can('CitizenStoreService.ListUserActivity')">
                            <CitizenActivityFeed :user-id="user.userId" />
                        </TabPanel>
                    </TabPanels>
                </TabGroup>
            </div>
        </template>
    </div>

    <AddToButton :callback="addToClipboard" :title="$t('components.clipboard.clipboard_button.add')" />
</template>
