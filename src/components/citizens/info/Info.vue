<script lang="ts" setup>
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/vue';
import { AccountIcon, BulletinBoardIcon, CarIcon, FileDocumentMultipleIcon } from 'mdi-vue3';
import AddToButton from '~/components/clipboard/AddToButton.vue';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { User } from '~~/gen/ts/resources/users/users';
import ActivityFeed from './ActivityFeed.vue';
import Documents from './Documents.vue';
import Profile from './Profile.vue';
import Vehicles from './Vehicles.vue';

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

const props = defineProps<{
    user: User;
}>();

function addToClipboard(): void {
    clipboardStore.addUser(props.user);

    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: {} },
        content: { key: 'notifications.clipboard.citizen_add.content', parameters: {} },
        duration: 3500,
        type: 'info',
    });
}
</script>

<template>
    <div class="py-2 pb-14">
        <div class="flex flex-row items-center gap-3">
            <h3 class="text-xl font-bold text-neutral sm:text-4xl inline-flex lg:px-4">
                {{ user?.firstname }}, {{ user?.lastname }}
            </h3>
            <span class="inline-flex items-center rounded-full bg-base-100 px-2.5 py-0.5 text-sm font-medium text-base-800">
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
        <TabGroup>
            <TabList class="border-b border-base-200 flex flex-row">
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
                                : 'border-transparent text-base-500 hover:border-base-300 hover:text-base-300',
                            'w-full justify-center group inline-flex items-center border-b-2 py-4 px-1 text-sm font-medium transition-colors',
                        ]"
                        :aria-current="selected ? 'page' : undefined"
                    >
                        <component
                            :is="tab.icon"
                            :class="[
                                selected ? 'text-primary-400' : 'text-base-500 group-hover:text-base-300',
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
            <TabPanels>
                <TabPanel>
                    <Profile :user="user" />
                </TabPanel>
                <TabPanel v-if="can('DMVService.ListVehicles')">
                    <Vehicles :userId="user.userId" />
                </TabPanel>
                <TabPanel v-if="can('DocStoreService.ListUserDocuments')">
                    <Documents :userId="user.userId" />
                </TabPanel>
                <TabPanel v-if="can('CitizenStoreService.ListUserActivity')">
                    <ActivityFeed :userId="user.userId" />
                </TabPanel>
            </TabPanels>
        </TabGroup>
    </div>
    <AddToButton :callback="addToClipboard" :title="$t('components.clipboard.clipboard_button.add')" />
</template>
