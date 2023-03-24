<script lang="ts" setup>
import { computed, ref, FunctionalComponent, onMounted } from 'vue';
import { RouteNamedMap } from 'vue-router/auto/routes';
import { toTitleCase } from '../../utils/strings';
import {
    Dialog,
    DialogPanel,
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue';
import {
    Bars3BottomLeftIcon,
    XMarkIcon,
    UsersIcon,
    DocumentTextIcon,
    BriefcaseIcon,
    MapIcon,
    HomeIcon,
    Square2StackIcon,
    UserIcon,
    TruckIcon,
} from '@heroicons/vue/24/outline';
import { ChevronRightIcon, HomeIcon as HomeIconSolid } from '@heroicons/vue/20/solid';
import { useStore } from '../../store/store';
import { useRoute, useRouter } from 'vue-router/auto';

const store = useStore();
const route = useRoute();
const router = useRouter();

const accessToken = computed(() => store.state.auth?.accessToken);
const activeChar = computed(() => store.state.auth?.activeChar);

const sidebarNavigation = [
    {
        name: 'Home',
        href: 'Home',
        permission: '',
        icon: HomeIcon,
    },
    {
        name: 'Overview',
        href: 'Overview',
        permission: 'Overview.View',
        icon: Square2StackIcon,
    },
    {
        name: 'Citizens',
        href: 'Citizens',
        permission: 'CitizenStoreService.FindUsers',
        icon: UsersIcon,
    },
    {
        name: 'Vehicles',
        href: 'Vehicles',
        permission: 'DMVService.FindVehicles',
        icon: TruckIcon,
    },
    {
        name: 'Documents',
        href: 'Documents',
        permission: 'DocStoreService.FindDocuments',
        icon: DocumentTextIcon,
    },
    {
        name: 'Job',
        href: 'Jobs',
        permission: 'Jobs.View',
        icon: BriefcaseIcon,
    },
    {
        name: 'Livemap',
        href: 'Livemap',
        permission: 'LivemapperService.Stream',
        icon: MapIcon,
    },
] as { name: string, href: keyof RouteNamedMap, permission: string, icon: FunctionalComponent }[];
const currSidebar = ref('')
let userNavigation = [
    { name: 'Login', href: 'Login' }
] as { name: string, href: string }[];
const breadcrumbs = [] as { name: string, href: string, current: boolean }[];
const mobileMenuOpen = ref(false);

onMounted(() => {
    if (accessToken.value && activeChar.value) {
        sidebarNavigation.shift();
        userNavigation = [
            { name: 'Change Character', href: 'Character Selector' },
            { name: 'Sign out', href: 'Logout' }
        ];
    }

    if (route.name) {
        const sidebarIndex = sidebarNavigation.findIndex(e => e.href.toLowerCase() === route.name.toLowerCase());
        if (sidebarIndex !== -1) {
            currSidebar.value = sidebarNavigation[sidebarIndex].name;
        } else {
            currSidebar.value = sidebarNavigation[0].name;
        }
    } else {
        currSidebar.value = sidebarNavigation[0].name;
    }

    const pathSplit = route.path.split('/').filter(e => e !== '');
    pathSplit.forEach(breadcrumb => {
        const route = router.getRoutes().find(r => r.name?.toString().toLowerCase() === breadcrumb.toLowerCase());

        if (route === undefined) {
            return;
        }

        breadcrumbs.push({
            name: toTitleCase(breadcrumb),
            href: route ? route.path : '/',
            current: route.name?.toString().toLowerCase() === breadcrumb.toLowerCase()
        })
    })
});
</script>

<template>
    <div class="flex h-screen">
        <!-- Sidebar -->
        <div class="hidden overflow-y-auto bg-accent-600 w-28 md:block">
            <div class="flex flex-col items-center w-full py-6">
                <div class="flex items-center flex-shrink-0">
                    <img class="w-auto h-12" src="/images/logo.png" alt="aRPaNet" />
                </div>
                <div class="flex-1 w-full px-2 mt-6 space-y-1">
                    <router-link v-for="item in sidebarNavigation" :key="item.name" :to="{ name: item.href }"
                        v-can="item.permission"
                        :class="[currSidebar === item.name ? 'bg-accent-100/20 text-neutral font-bold' : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium', 'hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2']"
                        :aria-current="currSidebar === item.name ? 'page' : undefined">
                        <component :is="item.icon"
                            :class="[currSidebar === item.name ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral', 'h-6 w-6']"
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
                                            :class="[currSidebar === item.name ? 'bg-accent-100/20 text-neutral font-bold' : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium', 'group flex items-center rounded-md py-2 px-3 text-sm']"
                                            :aria-current="currSidebar === item.name ? 'page' : undefined">
                                            <component :is="item.icon"
                                                :class="[currSidebar === item.name ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral', 'mr-3 h-6 w-6']"
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
                                            <router-link :to="{ path: accessToken ? '/overview' : '/' }"
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
                                            <router-link :to="{ path: page.href }"
                                                class="ml-4 text-sm font-medium text-base-400 hover:text-neutral hover:transition-colors"
                                                :aria-current="page.current ? 'page' : undefined">{{ page.name
                                                }}</router-link>
                                        </div>
                                    </li>
                                </ol>
                            </nav>
                        </div>
                        <div class="flex items-center ml-2 space-x-4 sm:ml-6 sm:space-x-6">
                            <span v-if="activeChar" class="text-sm font-medium text-base-400">{{ activeChar.getFirstname() }}, {{
                                activeChar.getLastname() }}</span>
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
            <div class="flex flex-1 overflow-hidden items-center justify-center">
                <main class="flex-1 overflow-y-auto">
                    <!-- Primary column -->
                    <section aria-labelledby="primary-heading" class="flex flex-col flex-1 h-full min-w-0 lg:order-last">
                        <slot></slot>
                    </section>
                </main>

                <!-- <aside class="hidden overflow-y-auto border-l-8 w-96 border-base-850 bg-base-900 lg:block"></aside> -->
            </div>
        </div>
    </div>
</template>
