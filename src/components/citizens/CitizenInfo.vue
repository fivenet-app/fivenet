<script lang="ts" setup>
import { RectangleGroupIcon, UserIcon, TruckIcon, DocumentTextIcon } from '@heroicons/vue/20/solid'
import { User } from '@arpanet/gen/resources/users/users_pb';
import CitizenInfoProfile from './CitizenInfoProfile.vue';
import CitizenInfoDocuments from './CitizenInfoDocuments.vue';
import CitizenActivityFeed from './CitizenActivityFeed.vue';
import VehiclesList from '../vehicles/VehiclesList.vue';
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from '@headlessui/vue';

const tabs = [
    { name: 'Profile', icon: UserIcon, permission: 'CitizenStoreService.FindUsers' },
    { name: 'Vehicles', icon: TruckIcon, permission: 'DMVService.FindVehicles' },
    { name: 'Documents', icon: DocumentTextIcon, permission: 'DocStoreService.GetUserDocuments' },
    { name: 'Activity', icon: RectangleGroupIcon, permission: 'CitizenStoreService.GetUserActivity' },
];

defineProps({
    user: {
        required: true,
        type: User,
    },
});
</script>

<template>
    <div class="sm:min-h-[52rem]">
        <div class="flex flex-row items-center gap-3">
            <h3 class="text-xl font-bold text-neutral sm:text-2xl inline-flex">
                {{ user?.getFirstname() }}, {{ user?.getLastname() }}
            </h3>
            <span
                class="inline-flex items-center rounded-full bg-base-100 px-2.5 py-0.5 text-sm font-medium text-base-800">{{
                    user.getJobLabel() }} (Rank: {{ user.getJobGradeLabel() }})
            </span>
            <span v-if="user.getProps()?.getWanted()"
                class="inline-flex items-center rounded-full bg-error-100 px-2.5 py-0.5 text-sm font-medium text-error-700">WANTED</span>
        </div>
        <TabGroup>
            <TabList class="border-b border-base-200">
                <Tab v-for="tab in tabs" :key="tab.name" v-slot="{ selected }" :v-can="tab.permission">
                    <button
                        :class="[selected ? 'border-primary-400 text-primary-500' : 'border-transparent text-base-500 hover:border-base-300 hover:text-base-700', 'group inline-flex items-center border-b-2 py-4 px-1 text-sm font-medium']"
                        :aria-current="selected ? 'page' : undefined">
                        <component :is="tab.icon"
                            :class="[selected ? 'text-primary-400' : 'text-base-400 group-hover:text-base-500', '-ml-0.5 mr-2 h-5 w-5']"
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
                    <VehiclesList :userId="user.getUserId()" :hide-owner="true" />
                </TabPanel>
                <TabPanel>
                    <CitizenInfoDocuments :userId="user.getUserId()" />
                </TabPanel>
                <TabPanel v-can="'CitizenStoreService.GetUserActivity'">
                    <CitizenActivityFeed :userId="user.getUserId()" />
                </TabPanel>
            </TabPanels>
        </TabGroup>
    </div>
</template>
