<script lang="ts" setup>
import { RectangleGroupIcon, UserIcon, TruckIcon, DocumentTextIcon } from '@heroicons/vue/20/solid'
import { User } from '@fivenet/gen/resources/users/users_pb';
import CitizenInfoProfile from './CitizenInfoProfile.vue';
import CitizenInfoDocuments from './CitizenInfoDocuments.vue';
import CitizenInfoActivityFeed from './CitizenInfoActivityFeed.vue';
import VehiclesList from '~/components/vehicles/VehiclesList.vue';
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/vue';
import { useClipboardStore } from '~/store/clipboard';
import { PlusIcon } from '@heroicons/vue/24/solid';
import { dispatchNotification } from '~/components/partials/notification';

const tabs = [
    { name: 'Profile', icon: UserIcon, permission: 'CitizenStoreService.FindUsers' },
    { name: 'Vehicles', icon: TruckIcon, permission: 'DMVService.FindVehicles' },
    { name: 'Documents', icon: DocumentTextIcon, permission: 'DocStoreService.GetUserDocuments' },
    { name: 'Activity', icon: RectangleGroupIcon, permission: 'CitizenStoreService.GetUserActivity' },
];

const store = useClipboardStore();

const props = defineProps({
    user: {
        required: true,
        type: User,
    },
});

function addToClipboard(): void {
    store.addUser(props.user);
    dispatchNotification({ title: 'Clipboard: Citizen added', content: 'Citizen has been added to clipboard', duration: 3500 });
}
</script>

<template>
    <div class="py-2">
        <div class="flex flex-row items-center gap-3">
            <p class="text-xl font-bold text-neutral sm:text-4xl inline-flex">
                {{ user?.getFirstname() }}, {{ user?.getLastname() }}
            </p>
            <span
                class="inline-flex items-center rounded-full bg-base-100 px-2.5 py-0.5 text-sm font-medium text-base-800">{{
                    user.getJobLabel() }} (Rank: {{ user.getJobGradeLabel() }})
            </span>
            <span v-if="user.getProps()?.getWanted()"
                class="inline-flex items-center rounded-full bg-error-100 px-2.5 py-0.5 text-sm font-medium text-error-700">WANTED</span>
        </div>
        <TabGroup>
            <TabList class="border-b border-base-200 flex flex-row">
                <Tab v-for="tab in tabs" :key="tab.name" v-slot="{ selected }" v-can="tab.permission" class="flex-1">
                    <button
                        :class="[selected ? 'border-primary-400 text-primary-500' : 'border-transparent text-base-500 hover:border-base-300 hover:text-base-300', 'w-full justify-center group inline-flex items-center border-b-2 py-4 px-1 text-sm font-medium transition-colors']"
                        :aria-current="selected ? 'page' : undefined">
                        <component :is="tab.icon"
                            :class="[selected ? 'text-primary-400' : 'text-base-500 group-hover:text-base-300', '-ml-0.5 mr-2 h-5 w-5']"
                            aria-hidden="true" />
                        <span>{{ tab.name }}</span>
                    </button>
                </Tab>
            </TabList>
            <TabPanels>
                <TabPanel>
                    <CitizenInfoProfile :user="user" />
                </TabPanel>
                <TabPanel>
                    <VehiclesList :userId="user.getUserId()" :hide-owner="true" :hide-citizen-link="true" />
                </TabPanel>
                <TabPanel>
                    <CitizenInfoDocuments :userId="user.getUserId()" />
                </TabPanel>
                <TabPanel v-can="'CitizenStoreService.GetUserActivity'">
                    <CitizenInfoActivityFeed :userId="user.getUserId()" />
                </TabPanel>
            </TabPanels>
        </TabGroup>
    </div>
    <button title="Add to Clipboard" @click="addToClipboard()"
        class="fixed flex items-center justify-center w-12 h-12 rounded-full z-90 bottom-24 right-8 bg-primary-500 shadow-float text-neutral hover:bg-primary-400">
        <PlusIcon class="w-10 h-auto" />
    </button>
</template>
