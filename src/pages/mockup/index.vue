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
    UserIcon,
} from '@heroicons/vue/24/outline'
import { MagnifyingGlassIcon, ChevronRightIcon, HomeIcon as HomeIconSolid } from '@heroicons/vue/20/solid'

export default defineComponent({
    components: {
        Bars3BottomLeftIcon,
        XMarkIcon,
        MagnifyingGlassIcon,
        UserIcon,
        ChevronRightIcon,
        HomeIconSolid,
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
                    href: '/mockup',
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
            breadcrumbs: [
                { name: 'Example', href: '#', current: false },
                { name: 'Example 2', href: '#', current: false }
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
        <!-- Sidebar -->
        <div class="hidden overflow-y-auto bg-accent-600 w-28 md:block">
            <div class="flex flex-col items-center w-full py-6">
                <div class="flex items-center flex-shrink-0">
                    <img class="w-auto h-12" src="/images/logo.png" alt="aRPaNet" />
                </div>
                <div class="flex-1 w-full px-2 mt-6 space-y-1">
                    <router-link v-for="item in sidebarNavigation" v-can="item.permission" :key="item.name" :to="item.href"
                        :class="[item.current ? 'bg-accent-100/20 text-neutral font-bold' : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium', 'hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2']"
                        :aria-current="item.current ? 'page' : undefined">
                        <component :is="item.icon"
                            :class="[item.current ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral', 'h-6 w-6']"
                            aria-hidden="true" />
                        <span class="mt-2">{{ item.name }}</span>
                    </router-link>
                </div>
            </div>
        </div>

        <!-- Mobile Sidebar -->
        <TransitionRoot as="template" :show="mobileMenuOpen">
            <Dialog as="div" class="relative z-20 md:hidden" @close="mobileMenuOpen = false">
                <TransitionChild as="template" enter="transition-opacity ease-linear duration-300" enter-from="opacity-0"
                    enter-to="opacity-100" leave="transition-opacity ease-linear duration-300" leave-from="opacity-100"
                    leave-to="opacity-0">
                    <div class="fixed inset-0 bg-opacity-75 bg-base-900" />
                </TransitionChild>

                <div class="fixed inset-0 z-40 flex">
                    <TransitionChild as="template" enter="transition ease-in-out duration-300 transform"
                        enter-from="-translate-x-full" enter-to="translate-x-0"
                        leave="transition ease-in-out duration-300 transform" leave-from="translate-x-0"
                        leave-to="-translate-x-full">
                        <DialogPanel class="relative flex flex-col flex-1 w-full max-w-xs pt-5 pb-4 bg-accent-600">
                            <TransitionChild as="template" enter="ease-in-out duration-300" enter-from="opacity-0"
                                enter-to="opacity-100" leave="ease-in-out duration-300" leave-from="opacity-100"
                                leave-to="opacity-0">
                                <div class="absolute p-1 -right-3 top-1 -mr-14">
                                    <button type="button"
                                        class="flex items-center justify-center w-12 h-12 rounded-full focus:outline-none ring-2 ring-neutral"
                                        @click="mobileMenuOpen = false">
                                        <XMarkIcon class="w-6 h-6 text-neutral" aria-hidden="true" />
                                        <span class="sr-only">Close sidebar</span>
                                    </button>
                                </div>
                            </TransitionChild>
                            <div class="flex items-center flex-shrink-0 px-4">
                                <img class="w-auto h-12" src="/images/logo.png" alt="aRPaNet" />
                            </div>
                            <div class="flex-1 h-0 px-2 mt-5 overflow-y-auto">
                                <nav class="flex flex-col h-full">
                                    <div class="space-y-1">
                                        <router-link v-for="item in sidebarNavigation" :key="item.name" :to="item.href"
                                            :class="[item.current ? 'bg-accent-100/20 text-neutral font-bold' : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium', 'group flex items-center rounded-md py-2 px-3 text-sm']"
                                            :aria-current="item.current ? 'page' : undefined">
                                            <component :is="item.icon"
                                                :class="[item.current ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral', 'mr-3 h-6 w-6']"
                                                aria-hidden="true" />
                                            <span>{{ item.name }}</span>
                                        </router-link>
                                    </div>
                                </nav>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                    <div class="flex-shrink-0 w-14" aria-hidden="true"></div>
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
                            <nav class="flex" aria-label="Breadcrumb">
                                <ol role="list" class="flex items-center space-x-4">
                                    <li>
                                        <div>
                                            <router-link to="/mockup"
                                                class="text-base-400 hover:text-neutral hover:transition-colors">
                                                <HomeIconSolid class="h-5 w-5 flex-shrink-0" aria-hidden="true" />
                                                <span class="sr-only">Home</span>
                                            </router-link>
                                        </div>
                                    </li>
                                    <li v-for="page in breadcrumbs" :key="page.name">
                                        <div class="flex items-center">
                                            <ChevronRightIcon class="h-5 w-5 flex-shrink-0 text-base-400"
                                                aria-hidden="true" />
                                            <router-link :to="page.href"
                                                class="ml-4 text-sm font-medium text-base-400 hover:text-neutral hover:transition-colors"
                                                :aria-current="page.current ? 'page' : undefined">{{ page.name
                                                }}</router-link>
                                        </div>
                                    </li>
                                </ol>
                            </nav>
                        </div>
                        <div class="flex items-center ml-2 space-x-4 sm:ml-6 sm:space-x-6">
                            <!-- Profile dropdown -->
                            <Menu as="div" class="relative flex-shrink-0">
                                <div>
                                    <MenuButton
                                        class="flex text-sm rounded-full bg-base-850 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                                        <span class="sr-only">Open user menu</span>
                                        <UserIcon
                                            class="w-auto h-10 hover:transition-colors rounded-full text-base-300 bg-base-800 fill-base-300 hover:text-base-200 hover:fill-base-200" />
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
                                            :class="[active ? 'bg-base-800' : '', 'block px-4 py-2 text-sm text-neutral hover:transition-colors']">{{
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
