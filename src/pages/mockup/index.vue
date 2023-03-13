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
        <div class="hidden w-28 overflow-y-auto bg-indigo-700 md:block">
            <div class="flex w-full flex-col items-center py-6">
                <div class="flex flex-shrink-0 items-center">
                    <img class="h-12 w-auto" src="/images/logo.png" alt="aRPaNet" />
                </div>
                <div class="mt-6 w-full flex-1 space-y-1 px-2">
                    <router-link v-for="item in sidebarNavigation" v-can="item.permission" :key="item.name" :to="item.href"
                        :class="[item.current ? 'bg-indigo-800 text-neutral' : 'text-indigo-100 hover:bg-indigo-800 hover:text-neutral', 'transition-colors group flex w-full flex-col items-center rounded-md p-3 text-xs font-medium']"
                        :aria-current="item.current ? 'page' : undefined">
                        <component :is="item.icon"
                            :class="[item.current ? 'text-neutral' : 'text-indigo-300 group-hover:text-neutral', 'h-6 w-6']"
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
                        <DialogPanel class="relative flex w-full max-w-xs flex-1 flex-col bg-indigo-700 pt-5 pb-4">
                            <TransitionChild as="template" enter="ease-in-out duration-300" enter-from="opacity-0"
                                enter-to="opacity-100" leave="ease-in-out duration-300" leave-from="opacity-100"
                                leave-to="opacity-0">
                                <div class="absolute top-1 right-0 -mr-14 p-1">
                                    <button type="button"
                                        class="flex h-12 w-12 items-center justify-center rounded-full focus:outline-none focus:ring-2 focus:ring-white"
                                        @click="mobileMenuOpen = false">
                                        <XMarkIcon class="h-6 w-6 text-white" aria-hidden="true" />
                                        <span class="sr-only">Close sidebar</span>
                                    </button>
                                </div>
                            </TransitionChild>
                            <div class="flex flex-shrink-0 items-center px-4">
                                <img class="h-12 w-auto" src="/images/logo.png" alt="Your Company" />
                            </div>
                            <div class="mt-5 h-0 flex-1 overflow-y-auto px-2">
                                <nav class="flex h-full flex-col">
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
                    <div class="w-14 flex-shrink-0" aria-hidden="true">
                        <!-- Dummy element to force sidebar to shrink to fit close icon -->
                    </div>
                </div>
            </Dialog>
        </TransitionRoot>

        <!-- Content area -->
        <div class="flex flex-1 flex-col overflow-hidden">
            <header class="w-full">
                <div class="relative z-10 flex h-16 flex-shrink-0 bg-base-850 shadow-sm">
                    <button type="button"
                        class="px-4 text-neutral focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500 md:hidden"
                        @click="mobileMenuOpen = true">
                        <span class="sr-only">Open sidebar</span>
                        <Bars3BottomLeftIcon class="h-6 w-6" aria-hidden="true" />
                    </button>
                    <div class="flex flex-1 justify-between px-4 sm:px-6">
                        <div class="flex flex-1">
                            <form class="flex w-full md:ml-0 bg-base-800 rounded-full px-5 my-3 max-w-2xl" action="#"
                                method="GET">
                                <label for="search-field" class="sr-only">Search all files</label>
                                <div class="relative w-full text-base-300 focus-within:text-neutral transition-colors">
                                    <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center">
                                        <MagnifyingGlassIcon class="h-5 w-5 flex-shrink-0" aria-hidden="true" />
                                    </div>
                                    <input name="search-field" id="search-field"
                                        class="h-full w-full border-0 py-2 pl-8 pr-3 bg-inherit transition-colors text-base-300 focus:outline-none focus:ring-0 focus:placeholder:text-neutral focus:text-neutral sm:text-sm"
                                        placeholder="Search" type="search" />
                                </div>
                            </form>
                        </div>
                        <div class="ml-2 flex items-center space-x-4 sm:ml-6 sm:space-x-6">
                            <!-- Profile dropdown -->
                            <Menu as="div" class="relative flex-shrink-0">
                                <div>
                                    <MenuButton
                                        class="flex rounded-full bg-base-850 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                                        <span class="sr-only">Open user menu</span>
                                        <UserIcon class="h-10 w-auto text-base-300 bg-base-800 fill-base-300 hover:text-base-200 hover:fill-base-200 transition-colors rounded-full"  />
                                    </MenuButton>
                                </div>
                                <transition enter-active-class="transition ease-out duration-100"
                                    enter-from-class="transform opacity-0 scale-95"
                                    enter-to-class="transform opacity-100 scale-100"
                                    leave-active-class="transition ease-in duration-75"
                                    leave-from-class="transform opacity-100 scale-100"
                                    leave-to-class="transform opacity-0 scale-95">
                                    <MenuItems
                                        class="absolute shadow-float right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-base-850 py-1 ring-1 ring-base-100 ring-opacity-5 focus:outline-none">
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
            <div class="flex flex-1 items-stretch overflow-hidden">
                <main class="flex-1 overflow-y-auto">
                    <!-- Primary column -->
                    <section aria-labelledby="primary-heading" class="flex h-full min-w-0 flex-1 flex-col lg:order-last">
                        <h1 id="primary-heading" class="sr-only">Photos</h1>
                        <!-- Your content -->
                    </section>
                </main>

                <!-- Secondary column (hidden on smaller screens) -->
                <aside class="hidden w-96 overflow-y-auto border-l-8 border-base-850 bg-base-900 lg:block">
                    <!-- Your content -->
                </aside>
            </div>
        </div>
    </div>
</template>
