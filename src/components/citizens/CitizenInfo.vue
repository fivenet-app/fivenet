<script lang="ts" setup>
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccount, mdiBulletinBoard, mdiCar, mdiFileDocumentMultiple } from '@mdi/js';
import VehiclesList from '~/components/vehicles/VehiclesList.vue';
import { can } from '~/plugins/1.vCan';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificationsStore } from '~/store/notifications';
import { User } from '~~/gen/ts/resources/users/users';
import AddToClipboardButton from '../clipboard/AddToClipboardButton.vue';
import CitizenInfoActivityFeed from './CitizenInfoActivityFeed.vue';
import CitizenInfoDocuments from './CitizenInfoDocuments.vue';
import CitizenInfoProfile from './CitizenInfoProfile.vue';

const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();

const { t } = useI18n();

const tabs = [
    {
        name: t('common.profile'),
        icon: mdiAccount,
        permission: 'CitizenStoreService.ListCitizens',
    },
    {
        name: t('common.vehicle', 2),
        icon: mdiCar,
        permission: 'DMVService.ListVehicles',
    },
    {
        name: t('common.document', 2),
        icon: mdiFileDocumentMultiple,
        permission: 'DocStoreService.ListUserDocuments',
    },
    {
        name: t('common.activity'),
        icon: mdiBulletinBoard,
        permission: 'CitizenStoreService.ListUserActivity',
    },
];

const props = defineProps<{
    user: User;
}>();

function addToClipboard(): void {
    clipboardStore.addUser(props.user);

    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: [] },
        content: { key: 'notifications.clipboard.citizen_add.content', parameters: [] },
        duration: 3500,
        type: 'info',
    });
}
</script>

<template>
    <div class="py-2">
        <div class="flex flex-row items-center gap-3">
            <p class="text-xl font-bold text-neutral sm:text-4xl inline-flex">{{ user?.firstname }}, {{ user?.lastname }}</p>
            <span class="inline-flex items-center rounded-full bg-base-100 px-2.5 py-0.5 text-sm font-medium text-base-800">
                {{ user.jobLabel }} ({{ $t('common.rank') }}: {{ user.jobGradeLabel }})
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
                        <SvgIcon
                            :class="[
                                selected ? 'text-primary-400' : 'text-base-500 group-hover:text-base-300',
                                '-ml-0.5 mr-2 h-5 w-5',
                            ]"
                            aria-hidden="true"
                            type="mdi"
                            :path="tab.icon"
                        />
                        <span>
                            {{ tab.name }}
                        </span>
                    </button>
                </Tab>
            </TabList>
            <TabPanels>
                <TabPanel>
                    <CitizenInfoProfile :user="user" />
                </TabPanel>
                <TabPanel v-can="'DMVService.ListVehicles'">
                    <VehiclesList :userId="user.userId" :hide-owner="true" :hide-citizen-link="true" />
                </TabPanel>
                <TabPanel v-can="'DocStoreService.ListUserDocuments'">
                    <CitizenInfoDocuments :userId="user.userId" />
                </TabPanel>
                <TabPanel v-can="'CitizenStoreService.ListUserActivity'">
                    <CitizenInfoActivityFeed :userId="user.userId" />
                </TabPanel>
            </TabPanels>
        </TabGroup>
    </div>
    <AddToClipboardButton :callback="addToClipboard" :title="$t('components.clipboard.clipboard_button.add')" />
</template>
