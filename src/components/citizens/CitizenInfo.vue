<script lang="ts">
import { defineComponent } from 'vue';
import { DocumentIcon, RectangleGroupIcon, UserIcon, } from '@heroicons/vue/20/solid'
import { User } from '@arpanet/gen/resources/users/users_pb';
import CitizenInfoProfile from './CitizenInfoProfile.vue';
import CitizenInfoDocuments from './CitizenInfoDocuments.vue';
import CitizenActivityFeed from './CitizenActivityFeed.vue';

export default defineComponent({
    components: {
        DocumentIcon,
        RectangleGroupIcon,
        UserIcon,
        CitizenInfoProfile,
        CitizenInfoDocuments,
        CitizenActivityFeed,
    },
    data() {
        return {
            currentTab: 'Profile' as string,
            tabs: [
                { name: 'Profile', icon: UserIcon, permission: 'CitizenStoreService.FindUsers' },
                { name: 'Activity', icon: RectangleGroupIcon, permission: 'CitizenStoreService.GetUserActivity' },
                { name: 'Documents', icon: DocumentIcon, permission: 'CitizenStoreService.GetUserDocuments' },
            ],
        };
    },
    props: {
        user: {
            required: true,
            type: User,
        },
    },
    methods: {
        setTab: function (tabName: string) {
            this.currentTab = tabName;
        },
    },
});
</script>

<template>
    <div>
        <div class="flex items-center">
            <h3 class="text-xl font-bold text-gray-300 sm:text-2xl">
                {{ user?.getFirstname() }}, {{ user?.getLastname() }}
                <span v-if="user.getProps()?.getWanted()"
                    class="inline-flex items-center rounded-md bg-red-100 px-2.5 py-0.5 text-sm font-medium text-red-800">WANTED</span>
            </h3>
        </div>
        <p class="text-sm text-white">
            <span class="inline-flex items-center rounded-md bg-gray-100 px-2.5 py-0.5 text-sm font-medium text-gray-800">{{
                user?.getJob() }} (Rank: {{ user?.getJobGrade() }})
            </span>
        </p>
    </div>
    <div>
        <div class="sm:hidden">
            <label for="tabs" class="sr-only">Select a tab</label>
            <select id="tabs" name="tabs"
                class="block w-full rounded-md border-gray-300 focus:border-indigo-500 focus:ring-indigo-500"
                v-model="currentTab">
                <option v-for="tab in tabs" :key="tab.name" :selected="tab.name === currentTab" v-can="tab.permission">{{ tab.name }}</option>
            </select>
        </div>
        <div class="hidden sm:block">
            <div class="border-b border-gray-200">
                <nav class="-mb-px flex space-x-8" aria-label="Tabs">
                    <button v-for="tab in tabs" :key="tab.name" href="#" @click="setTab(tab.name)"
                        :class="[tab.name === currentTab ? 'border-indigo-500 text-indigo-600' : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700', 'group inline-flex items-center border-b-2 py-4 px-1 text-sm font-medium']"
                        :aria-current="tab.name === currentTab ? 'page' : undefined">
                        <component :is="tab.icon"
                            :class="[tab.name === currentTab ? 'text-indigo-500' : 'text-gray-400 group-hover:text-gray-500', '-ml-0.5 mr-2 h-5 w-5']"
                            aria-hidden="true" />
                        <span>{{ tab.name }}</span>
                    </button>
                </nav>
            </div>
        </div>
    </div>
    <div class="p-3 mt-6">
        <div v-if="currentTab === 'Profile'" v-can="'CitizenStoreService.FindUsers'">
            <CitizenInfoProfile :user="user" />
        </div>
        <div v-if="currentTab === 'Activity'" v-can="'CitizenStoreService.GetUserActivity'">
            <CitizenActivityFeed :userId="user.getUserId()" />
        </div>
        <div v-if="currentTab === 'Documents'" v-can="'CitizenStoreService.GetUserDocuments'">
            <CitizenInfoDocuments :userId="user.getUserId()" />
        </div>
    </div>
</template>
