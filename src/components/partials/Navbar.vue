<script lang="ts">
import { themeChange } from 'theme-change';
import { defineComponent } from 'vue';
import { mapState } from 'vuex';

import { Disclosure, DisclosureButton, DisclosurePanel, Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue';
import { Bars3Icon, BellIcon, XMarkIcon, UserCircleIcon } from '@heroicons/vue/24/outline';
import { HomeIcon } from '@heroicons/vue/20/solid';

export default defineComponent({
    components: {
        Disclosure,
        DisclosureButton,
        DisclosurePanel,
        Menu,
        MenuButton,
        MenuItem,
        MenuItems,
        Bars3Icon,
        BellIcon,
        XMarkIcon,
        UserCircleIcon,
        HomeIcon,
    },
    computed: {
        ...mapState({
            accessToken: 'accessToken',
            activeChar: 'activeChar',
            activeCharID: 'activeCharID',
        }),
    },
    data() {
        return {
            navigation: [
                { name: 'Overview', href: '/overview', permission: 'overview-view', },
                { name: 'Citizens', href: '/citizens', permission: 'users-view' },
                { name: 'Documents', href: '/documents', permission: 'documents-view' },
                { name: 'Dispatches', href: '/dispatches', permission: 'dispatches-view' },
                { name: 'Job', href: '/job', permission: 'job-view' },
                { name: 'Livemap', href: '/livemap', permission: 'livemap-stream' },
            ],
            userNavigation: [
                { name: 'Change Characters', href: '/login', },
                { name: 'Sign out', href: '/logout', },
            ],
        };
    },
    mounted() {
        themeChange(false);
    },
});
</script>

<template>
    <Disclosure as="nav" class="bg-gray-800" v-slot="{ open }">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <div class="flex h-16 items-center justify-between">
                <div class="flex items-center">
                    <div class="flex-shrink-0">
                        <img class="h-8 w-8" src="/images/logo.png" alt="aRPaNet Logo" />
                    </div>
                    <div v-if="accessToken" class="hidden md:block">
                        <div class="ml-10 flex items-baseline space-x-4">
                            <router-link v-for="item in navigation" :key="item.name" :to="item.href" v-can="item.permission"
                                class="text-gray-300 hover:bg-gray-700 hover:text-white rounded-md px-3 py-2 text-sm font-medium"
                                active-class="bg-gray-900 text-white"
                                :aria-current="$route.name == item.href ? 'page' : undefined">{{ item.name
                                }}</router-link>
                        </div>
                    </div>
                </div>
                <div v-if="accessToken" class="hidden md:block">
                    <div class="ml-4 flex items-center md:ml-6">
                        <button type="button"
                            class="rounded-full bg-gray-800 p-1 text-gray-400 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                            <span class="sr-only">View notifications</span>
                            <BellIcon class="h-6 w-6" aria-hidden="true" />
                        </button>

                        <!-- Profile dropdown -->
                        <Menu as="div" class="relative ml-3">
                            <div class="text-gray-400">
                                <MenuButton
                                    class="flex max-w-xs items-center rounded-full bg-gray-800 text-sm focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                                    <span class="sr-only">Open user menu</span>
                                    <UserCircleIcon class="h-6 w-6" aria-hidden="true" />
                                </MenuButton>
                            </div>
                            <transition enter-active-class="transition ease-out duration-100"
                                enter-from-class="transform opacity-0 scale-95"
                                enter-to-class="transform opacity-100 scale-100"
                                leave-active-class="transition ease-in duration-75"
                                leave-from-class="transform opacity-100 scale-100"
                                leave-to-class="transform opacity-0 scale-95">
                                <MenuItems
                                    class="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                                    <MenuItem v-for="item in userNavigation" :key="item.name" v-slot="{ active }">
                                    <router-link :to="item.href"
                                        :class="[active ? 'bg-gray-100' : '', 'block px-4 py-2 text-sm text-gray-700']">{{
                                            item.name }}</router-link>
                                    </MenuItem>
                                </MenuItems>
                            </transition>
                        </Menu>
                    </div>
                </div>
                <div class="-mr-2 flex md:hidden">
                    <!-- Mobile menu button -->
                    <button
                        class="inline-flex items-center justify-center rounded-md bg-gray-800 p-2 text-gray-400 hover:bg-gray-700 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                        <span class="sr-only">Open main menu</span>
                        <Bars3Icon v-if="!open" class="block h-6 w-6" aria-hidden="true" />
                        <XMarkIcon v-else class="block h-6 w-6" aria-hidden="true" />
                    </button>
                </div>
            </div>
        </div>

        <DisclosurePanel class="md:hidden" v-if="activeChar">
            <div class="space-y-1 px-2 pt-2 pb-3 sm:px-3">
                <router-link v-for="item in navigation" :key="item.name" as="a" :to="item.href"
                    :class="[$route.name == item.href ? 'bg-gray-900 text-white' : 'text-gray-300 hover:bg-gray-700 hover:text-white', 'block rounded-md px-3 py-2 text-base font-medium']"
                    :aria-current="$route.name == item.href ? 'page' : undefined">{{ item.name }}</router-link>
            </div>
            <div class="border-t border-gray-700 pt-4 pb-3">
                <div class="flex items-center px-5">
                    <div class="flex-shrink-0">
                        <UserCircleIcon class="h-6 w-6" aria-hidden="true" />
                    </div>
                    <div class="ml-3">
                        <div class="text-base font-medium leading-none text-white">{{ activeChar.getFirstname() }}</div>
                        <div class="text-sm font-medium leading-none text-gray-400">{{ activeChar.getLastname() }}</div>
                    </div>
                    <button type="button"
                        class="ml-auto flex-shrink-0 rounded-full bg-gray-800 p-1 text-gray-400 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                        <span class="sr-only">View notifications</span>
                        <BellIcon class="h-6 w-6" aria-hidden="true" />
                    </button>
                </div>
                <div class="mt-3 space-y-1 px-2">
                    <router-link v-for="item in userNavigation" :key="item.name" :to="item.href"
                        class="block rounded-md px-3 py-2 text-base font-medium text-gray-400 hover:bg-gray-700 hover:text-white">
                        {{ item.name }}</router-link>
                </div>
            </div>
        </DisclosurePanel>
    </Disclosure>
</template>
