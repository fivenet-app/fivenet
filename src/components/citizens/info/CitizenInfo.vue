<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/vue';
import { AccountIcon, BulletinBoardIcon, CarIcon, FileDocumentMultipleIcon } from 'mdi-vue3';
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

const props = defineProps<{
    id: string;
}>();

const { $grpc } = useNuxtApp();

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const { t } = useI18n();

const tabs = [
    {
        name: t('common.profile'),
        icon: markRaw(AccountIcon),
        permission: 'CitizenStoreService.ListCitizens',
    },
    {
        name: t('common.vehicle', 2),
        icon: markRaw(CarIcon),
        permission: 'DMVService.ListVehicles',
    },
    {
        name: t('common.document', 2),
        icon: markRaw(FileDocumentMultipleIcon),
        permission: 'DocStoreService.ListUserDocuments',
    },
    {
        name: t('common.activity'),
        icon: markRaw(BulletinBoardIcon),
        permission: 'CitizenStoreService.ListUserActivity',
    },
];

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
                <div class="px-4">
                    <h1 class="flex text-4xl font-bold text-neutral">{{ user?.firstname }}, {{ user?.lastname }}</h1>
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

                <TabGroup>
                    <TabList class="flex flex-row border-b-2 border-neutral/20">
                        <Tab
                            v-for="tab in tabs.filter((tab) => can(tab.permission))"
                            :key="tab.name"
                            v-slot="{ selected }"
                            class="flex-1"
                        >
                            <button
                                :class="[
                                    selected
                                        ? 'border-primary-400 text-primary-500'
                                        : 'border-transparent text-base-200 hover:border-base-300 hover:text-base-300',
                                    'group inline-flex w-full items-center justify-center border-b-2 px-1 py-4 text-sm font-medium transition-colors',
                                ]"
                                :aria-current="selected ? 'page' : undefined"
                            >
                                <component
                                    :is="tab.icon"
                                    :class="[
                                        selected ? 'text-primary-400' : 'text-base-200 group-hover:text-base-300',
                                        '-ml-0.5 mr-2 h-5 w-5',
                                    ]"
                                    aria-hidden="true"
                                />
                                <span>
                                    {{ tab.name }}
                                </span>
                            </button>
                        </Tab>
                    </TabList>
                    <TabPanels class="bg-transparent">
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
