<script lang="ts">
import { defineComponent } from 'vue';
import { ref } from 'vue'
import {
    Dialog,
    DialogPanel,
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue'
import {
    Bars3BottomLeftIcon,
    XMarkIcon,
    UsersIcon,
    DocumentTextIcon,
    BellAlertIcon,
    BriefcaseIcon,
    MapIcon,
    HomeIcon,
    Square2StackIcon,
    UserIcon
} from '@heroicons/vue/24/outline'
import { MagnifyingGlassIcon } from '@heroicons/vue/20/solid'

export default defineComponent({
    components: {
        Bars3BottomLeftIcon,
        XMarkIcon,
        MagnifyingGlassIcon,
        UserIcon,
        Dialog,
        DialogPanel,
        Menu,
        MenuButton,
        MenuItem,
        MenuItems,
        TransitionChild,
        TransitionRoot,
    },
    data() {
        return {
            sidebarNavigation: [
                {
                    name: 'Home',
                    href: '/mockup/',
                    permission: '',
                    icon: HomeIcon,
                    current: true,
                },
                {
                    name: 'Overview',
                    href: '/mockup/overview',
                    permission: 'overview-view',
                    icon: Square2StackIcon,
                    current: false,
                },
                {
                    name: 'Citizens',
                    href: '/mockup/citizens',
                    permission: 'users-findusers',
                    icon: UsersIcon,
                    current: false,
                },
                {
                    name: 'Documents',
                    href: '/mockup/documents',
                    permission: 'documents-view',
                    icon: DocumentTextIcon,
                    current: false,
                },
                {
                    name: 'Dispatches',
                    href: '/mockup/dispatches',
                    permission: 'dispatches-view',
                    icon: BellAlertIcon,
                    current: false,
                },
                {
                    name: 'Job',
                    href: '/mockup/job',
                    permission: 'job-view',
                    icon: BriefcaseIcon,
                    current: false,
                },
                {
                    name: 'Livemap',
                    href: '/mockup/livemap',
                    permission: 'livemap-stream',
                    icon: MapIcon,
                    current: false,
                },
            ],
            userNavigation: [
                { name: 'Your Profile', href: '#' },
                { name: 'Sign out', href: '#' },
            ],
            mobileMenuOpen: ref(false)
        }
    }
});
</script>

<route lang="json">
{
    "name": "mockup",
    "meta": {
        "requiresAuth": false,
        "breadCrumbs": []
    }
}
</route>

<template>
    <div class="flex h-screen">
        <!-- Narrow sidebar -->
        <div class="hidden overflow-y-auto bg-accent-600 w-28 md:block">
            <div class="flex flex-col items-center w-full py-6">
                <div class="flex items-center flex-shrink-0">
                    <img class="w-auto h-12" src="/images/logo.png" alt="aRPaNet" />
                </div>
                <div class="flex-1 w-full px-2 mt-6 space-y-1">
                    <router-link v-for="item in sidebarNavigation" v-can="item.permission" :key="item.name" :to="item.href"
                        :class="[item.current ? 'bg-accent-100/20 text-neutral font-bold' : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium', 'ransition-colors group flex w-full flex-col items-center rounded-md p-3 text-xs my-2']"
                        :aria-current="item.current ? 'page' : undefined">
                        <component :is="item.icon"
                            :class="[item.current ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral', 'h-6 w-6']"
                            aria-hidden="true" />
                        <span class="mt-2">{{ item.name }}</span>
                    </router-link>
                </div>
            </div>
        </div>

        <!-- Mobile menu -->
        <TransitionRoot as="template" :show="mobileMenuOpen">
            <Dialog as="div" class="relative z-20 md:hidden" @close="mobileMenuOpen = false">
                <TransitionChild as="template" enter="transition-opacity ease-linear duration-300" enter-from="opacity-0"
                    enter-to="opacity-100" leave="transition-opacity ease-linear duration-300" leave-from="opacity-100"
                    leave-to="opacity-0">
                    <div class="fixed inset-0 bg-gray-600 bg-opacity-75" />
                </TransitionChild>

                <div class="fixed inset-0 z-40 flex">
                    <TransitionChild as="template" enter="transition ease-in-out duration-300 transform"
                        enter-from="-translate-x-full" enter-to="translate-x-0"
                        leave="transition ease-in-out duration-300 transform" leave-from="translate-x-0"
                        leave-to="-translate-x-full">
                        <DialogPanel class="relative flex flex-col flex-1 w-full max-w-xs pt-5 pb-4 bg-indigo-700">
                            <TransitionChild as="template" enter="ease-in-out duration-300" enter-from="opacity-0"
                                enter-to="opacity-100" leave="ease-in-out duration-300" leave-from="opacity-100"
                                leave-to="opacity-0">
                                <div class="absolute right-0 p-1 top-1 -mr-14">
                                    <button type="button"
                                        class="flex items-center justify-center w-12 h-12 rounded-full focus:outline-none focus:ring-2 focus:ring-white"
                                        @click="mobileMenuOpen = false">
                                        <XMarkIcon class="w-6 h-6 text-white" aria-hidden="true" />
                                        <span class="sr-only">Close sidebar</span>
                                    </button>
                                </div>
                            </TransitionChild>
                            <div class="flex items-center flex-shrink-0 px-4">
                                <img class="w-auto h-12" src="/images/logo.png" alt="Your Company" />
                            </div>
                            <div class="flex-1 h-0 px-2 mt-5 overflow-y-auto">
                                <nav class="flex flex-col h-full">
                                    <div class="space-y-1">
                                        <router-link v-for="item in sidebarNavigation" :key="item.name" :to="item.href"
                                            :class="[item.current ? 'bg-indigo-800 text-white' : 'text-indigo-100 hover:bg-indigo-800 hover:text-white', 'group flex items-center rounded-md py-2 px-3 text-sm font-medium']"
                                            :aria-current="item.current ? 'page' : undefined">
                                            <component :is="item.icon"
                                                :class="[item.current ? 'text-white' : 'text-indigo-300 group-hover:text-white', 'mr-3 h-6 w-6']"
                                                aria-hidden="true" />
                                            <span>{{ item.name }}</span>
                                        </router-link>
                                    </div>
                                </nav>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                    <div class="flex-shrink-0 w-14" aria-hidden="true">
                        <!-- Dummy element to force sidebar to shrink to fit close icon -->
                    </div>
                </div>
            </Dialog>
        </TransitionRoot>

        <!-- Content area -->
        <div class="flex flex-col flex-1 overflow-hidden">
            <header class="w-full">
                <div class="relative z-10 flex flex-shrink-0 h-16 shadow-sm bg-base-850">
                    <button type="button"
                        class="px-4 text-neutral focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500 md:hidden"
                        @click="mobileMenuOpen = true">
                        <span class="sr-only">Open sidebar</span>
                        <Bars3BottomLeftIcon class="w-6 h-6" aria-hidden="true" />
                    </button>
                    <div class="flex justify-between flex-1 px-4 sm:px-6">
                        <div class="flex flex-1">
                            <form class="flex w-full max-w-2xl px-5 my-3 rounded-full md:ml-0 bg-base-800" action="#"
                                method="GET">
                                <label for="search-field" class="sr-only">Search all files</label>
                                <div class="relative w-full transition-colors text-base-300 focus-within:text-neutral">
                                    <div class="absolute inset-y-0 left-0 flex items-center pointer-events-none">
                                        <MagnifyingGlassIcon class="flex-shrink-0 w-5 h-5" aria-hidden="true" />
                                    </div>
                                    <input name="search-field" id="search-field"
                                        class="w-full h-full py-2 pl-8 pr-3 transition-colors border-0 bg-inherit text-base-300 focus:outline-none focus:ring-0 focus:placeholder:text-neutral focus:text-neutral sm:text-sm"
                                        placeholder="Search" type="search" />
                                </div>
                            </form>
                        </div>
                        <div class="flex items-center ml-2 space-x-4 sm:ml-6 sm:space-x-6">
                            <!-- Profile dropdown -->
                            <Menu as="div" class="relative flex-shrink-0">
                                <div>
                                    <MenuButton
                                        class="flex text-sm rounded-full bg-base-850 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                                        <span class="sr-only">Open user menu</span>
                                        <UserIcon class="w-auto h-10 transition-colors rounded-full text-base-300 bg-base-800 fill-base-300 hover:text-base-200 hover:fill-base-200"  />
                                    </MenuButton>
                                </div>
                                <transition enter-active-class="transition duration-100 ease-out"
                                    enter-from-class="transform scale-95 opacity-0"
                                    enter-to-class="transform scale-100 opacity-100"
                                    leave-active-class="transition duration-75 ease-in"
                                    leave-from-class="transform scale-100 opacity-100"
                                    leave-to-class="transform scale-95 opacity-0">
                                    <MenuItems
                                        class="absolute right-0 z-10 w-48 py-1 mt-2 origin-top-right rounded-md shadow-float bg-base-850 ring-1 ring-base-100 ring-opacity-5 focus:outline-none">
                                        <MenuItem v-for="item in userNavigation" :key="item.name" v-slot="{ active }">
                                        <a :href="item.href"
                                            :class="[active ? 'bg-base-800' : '', 'block px-4 py-2 text-sm text-neutral transition-colors']">{{
                                                item.name }}</a>
                                        </MenuItem>
                                    </MenuItems>
                                </transition>
                            </Menu>
                        </div>
                    </div>
                </div>
            </header>

            <!-- Main content -->
            <div class="flex items-stretch flex-1 overflow-hidden">
                <main class="flex-1 overflow-y-auto">
                    <!-- Primary column -->
                    <section aria-labelledby="primary-heading" class="flex flex-col flex-1 h-full min-w-0 lg:order-last">
                        <h1 id="primary-heading" class="sr-only">Photos</h1>
                        <!-- Your content -->
                    </section>
                </main>

                <!-- Secondary column (hidden on smaller screens) -->
                <aside class="hidden overflow-y-auto border-l-8 w-96 border-base-850 bg-base-900 lg:block">
                    <!-- Your content -->
                </aside>
            </div>
        </div>
    </div>
</template>
