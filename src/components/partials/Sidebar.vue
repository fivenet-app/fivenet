<script lang="ts" setup>
import { computed, ref, FunctionalComponent, onMounted } from 'vue';
import { useAuthStore } from '~/store/auth';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';
import { toTitleCase } from '~/utils/strings';
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
    UserIcon,
    TruckIcon,
    QuestionMarkCircleIcon,
    Cog8ToothIcon,
} from '@heroicons/vue/24/outline';
import { ChevronRightIcon, HomeIcon as HomeIconSolid } from '@heroicons/vue/20/solid';
import SidebarJobSwitcher from '~/components/partials/SidebarJobSwitcher.vue';
import SidebarLanguageSwitcher from '~/components/partials/SidebarLanguageSwitcher.vue';

const { t } = useI18n();
const store = useAuthStore();
const router = useRouter();

const accessToken = computed(() => store.$state.accessToken);
const activeChar = computed(() => store.$state.activeChar);

const sidebarNavigation: { name: string, href: RoutesNamedLocations, permission: string, icon: FunctionalComponent, position: 'top' | 'bottom', current: boolean }[] = [
    {
        name: t('common.overview'),
        href: { name: 'overview' },
        permission: '',
        icon: HomeIcon,
        position: 'top',
        current: false,
    },
    {
        name: t('common.citizen', 2),
        href: { name: 'citizens' },
        permission: 'CitizenStoreService.FindUsers',
        icon: UsersIcon,
        position: 'top',
        current: false,
    },
    {
        name: t('common.vehicle', 2),
        href: { name: 'vehicles' },
        permission: 'DMVService.FindVehicles',
        icon: TruckIcon,
        position: 'top',
        current: false,
    },
    {
        name: t('common.document', 2),
        href: { name: 'documents' },
        permission: 'DocStoreService.FindDocuments',
        icon: DocumentTextIcon,
        position: 'top',
        current: false,
    },
    {
        name: t('common.job', 2),
        href: { name: 'jobs' },
        permission: 'Jobs.View',
        icon: BriefcaseIcon,
        position: 'top',
        current: false,
    },
    {
        name: t('common.livemap'),
        href: { name: 'livemap' },
        permission: 'LivemapperService.Stream',
        icon: MapIcon,
        position: 'top',
        current: false,
    },
    {
        name: t('common.control_panel'),
        href: { name: 'rector' },
        permission: 'RectorService.GetRoles',
        icon: Cog8ToothIcon,
        position: 'top',
        current: false,
    },
    {
        name: t('common.about'),
        href: { name: 'about' },
        permission: '',
        icon: QuestionMarkCircleIcon,
        position: 'bottom',
        current: false,
    },
];
const userNavigation = ref<{ name: string, href: RoutesNamedLocations, permission?: string }[]>([
    { name: t('common.login'), href: { name: 'auth-login' } },
]);
const breadcrumbs = ref<{ name: string, href: string, current: boolean }[]>([]);
const mobileMenuOpen = ref(false);

onMounted(() => {
    updateUserNav();
    updateActiveItem();
    updateBread();
});

function updateUserNav(): void {
    userNavigation.value.length = 0;
    if (activeChar.value) {
        userNavigation.value.push(
            { name: t('components.partials.sidebar.change_character'), href: { name: 'auth-character-selector' } },
        );
    }
    if (accessToken.value) {
        userNavigation.value.push(
            { name: t('components.partials.sidebar.account_info'), href: { name: 'auth-account-info' } },
            { name: t('common.sign_out'), href: { name: 'auth-logout' } },
        );
    }
    if (userNavigation.value.length === 0) {
        userNavigation.value = [
            { name: t('common.login'), href: { name: 'auth-login' } },
        ];
    }
}

function updateActiveItem(): void {
    const route = router.currentRoute.value;
    if (route.name) {
        sidebarNavigation.forEach(e => {
            if (route.name.toLowerCase().includes(e.href.name.toLowerCase())) {
                e.current = true;
            } else {
                e.current = false;
            }
        });
    } else {
        sidebarNavigation.forEach(e => e.current = false);
    }
}

function updateBread(): void {
    // Clear current breadcrumbs
    breadcrumbs.value.length = 0;
    const currentRoute = router.currentRoute.value;

    const pathSplit = currentRoute.path.split('/').filter(e => e !== '');
    pathSplit.forEach((breadcrumb, idx) => {
        breadcrumb = "/" + breadcrumb;
        if (idx > 0) {
            breadcrumb = "/" + pathSplit.slice(0, idx).join("/") + breadcrumb;
        }
        const route = router.getRoutes().find(r => r.path.toLowerCase() === breadcrumb.toLowerCase());
        if (route === undefined) {
            return;
        }

        if (route.path.toLowerCase() === currentRoute.path.toLowerCase()) {
            return;
        }

        breadcrumbs.value.push({
            name: route.meta.title ?? toTitleCase(breadcrumb),
            href: route.path,
            current: false,
        });
    });

    const breadcrumbIdx = breadcrumbs.value.findIndex((b) => {
        return b.href === currentRoute.path;
    })
    if (breadcrumbIdx === -1) {
        if (currentRoute.name !== 'index' && currentRoute.name !== 'overview') {
            breadcrumbs.value.push({
                name: toTitleCase(currentRoute.meta.title ?? currentRoute.name),
                href: '#',
                current: true,
            });
        }
    }
}

watch(accessToken, () => updateUserNav());
watch(activeChar, () => updateUserNav());

watch(router.currentRoute, () => {
    updateActiveItem();
    updateBread();
});

const appVersion = activeChar ? (' v' + __APP_VERSION__ + (import.meta.env.DEV ? '-dev' : '-prod')) : '';
</script>

<template>
    <div class="flex h-screen">
        <!-- Sidebar -->
        <div class="hidden overflow-y-auto bg-accent-600 w-28 md:block">
            <div class="flex flex-col items-center w-full py-6 h-full">
                <div class="flex items-center flex-shrink-0">
                    <NuxtLink :to="{ name: accessToken ? 'overview' : 'index' }" aria-current-value="page">
                        <nuxt-img class="w-auto h-12" src="/images/logo.png" alt="FiveNet Logo"
                            :title="'FiveNet' + appVersion" />
                    </NuxtLink>
                </div>
                <div class="flex-grow w-full px-2 mt-6 space-y-1 text-center">
                    <NuxtLink :to="{ name: 'index' }" v-if="!accessToken || !activeChar"
                        active-class="bg-accent-100/20 text-neutral font-bold"
                        class="text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2"
                        exact-active-class="text-neutral" aria-current-value="page">
                        <HomeIcon class="h-6 w-6" aria-hidden="true" />
                        <span class="mt-2">{{ $t('common.home') }}</span>
                    </NuxtLink>
                    <NuxtLink v-for="item in sidebarNavigation.filter(e => e.position === 'top')"
                        v-else-if="accessToken && activeChar" :key="item.name" :to="item.href" v-can="item.permission"
                        :class="[item.current ? 'bg-accent-100/20 text-neutral font-bold' : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium', 'hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2']"
                        :aria-current="item.current ? 'page' : undefined">
                        <component :is="item.icon"
                            :class="[item.current ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral', 'h-6 w-6']"
                            aria-hidden="true" />
                        <span class="mt-2">{{ item.name }}</span>
                    </NuxtLink>
                </div>
                <div class="flex-initial w-full px-2 space-y-1 text-center">
                    <NuxtLink v-for="item in sidebarNavigation.filter(e => e.position === 'bottom')" :key="item.name"
                        :to="item.href" v-can="item.permission"
                        :class="[item.current ? 'bg-accent-100/20 text-neutral font-bold' : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium', 'hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2']"
                        :aria-current="item.current ? 'page' : undefined">
                        <component :is="item.icon"
                            :class="[item.current ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral', 'h-6 w-6']"
                            aria-hidden="true" />
                        <span class="mt-2">{{ item.name }}</span>
                    </NuxtLink>
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
                                        <span class="sr-only">{{ $t('components.partials.sidebar.close_sidebar') }}</span>
                                    </button>
                                </div>
                            </TransitionChild>
                            <div class="flex items-center flex-shrink-0 px-4">
                                <nuxt-img class="w-16 h-16 mx-auto" src="/images/logo.png" alt="FiveNet Logo" />
                            </div>
                            <div class="flex-grow h-0 px-2 mt-5 overflow-y-auto">
                                <nav class="flex flex-col h-full">
                                    <div class="space-y-1">
                                        <NuxtLink v-if="!accessToken || !activeChar" :to="{ name: 'index' }"
                                            class="text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium group flex items-center rounded-md py-2 px-3 text-sm">
                                            <HomeIcon class="text-accent-100 group-hover:text-neutral mr-3 h-6 w-6"
                                                aria-hidden="true" />
                                            <span>{{ $t('common.home') }}</span>
                                        </NuxtLink>
                                        <NuxtLink v-for="item in sidebarNavigation.filter(e => e.position === 'top')"
                                            v-else-if="accessToken && activeChar" :key="item.name" :to="item.href"
                                            v-can="item.permission"
                                            :class="[item.current ? 'bg-accent-100/20 text-neutral font-bold' : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium', 'group flex items-center rounded-md py-2 px-3 text-sm']"
                                            :aria-current="item.current ? 'page' : undefined">
                                            <component :is="item.icon"
                                                :class="[item.current ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral', 'mr-3 h-6 w-6']"
                                                aria-hidden="true" />
                                            <span>{{ item.name }}</span>
                                        </NuxtLink>
                                    </div>
                                </nav>
                            </div>
                            <div class="flex-initial px-2 overflow-y-auto">
                                <nav class="flex flex-col h-full">
                                    <div class="space-y-1">
                                        <NuxtLink v-for="item in sidebarNavigation.filter(e => e.position === 'bottom')"
                                            :key="item.name" :to="item.href" v-can="item.permission"
                                            :class="[item.current ? 'bg-accent-100/20 text-neutral font-bold' : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium', 'group flex items-center rounded-md py-2 px-3 text-sm']"
                                            :aria-current="item.current ? 'page' : undefined">
                                            <component :is="item.icon"
                                                :class="[item.current ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral', 'mr-3 h-6 w-6']"
                                                aria-hidden="true" />
                                            <span>{{ item.name }}</span>
                                        </NuxtLink>
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
                <div class="relative z-10 flex flex-shrink-0 h-16 bg-base-850">
                    <button type="button"
                        class="px-4 text-neutral focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500 md:hidden"
                        @click="mobileMenuOpen = true">
                        <span class="sr-only">{{ $t('components.partials.sidebar.open_sidebar') }}</span>
                        <Bars3BottomLeftIcon class="w-6 h-6" aria-hidden="true" />
                    </button>
                    <div class="flex justify-between flex-1 px-4 sm:px-6">
                        <div class="flex flex-1">
                            <nav class="flex" aria-label="Breadcrumb">
                                <ol role="list" class="flex items-center space-x-4">
                                    <li>
                                        <div>
                                            <NuxtLink :to="{ name: accessToken ? 'overview' : 'index' }"
                                                class="text-base-400 hover:text-neutral hover:transition-colors">
                                                <HomeIconSolid class="flex-shrink-0 w-5 h-5" aria-hidden="true" />
                                                <span class="sr-only">Home</span>
                                            </NuxtLink>
                                        </div>
                                    </li>
                                    <li v-for="page in breadcrumbs" :key="page.name">
                                        <div class="flex items-center">
                                            <ChevronRightIcon class="flex-shrink-0 w-5 h-5 text-base-400"
                                                aria-hidden="true" />
                                            <NuxtLink :to="{ path: page.href }"
                                                :class="[page.current ? 'font-bold text-base-200' : 'font-medium text-base-400', 'ml-4 text-sm hover:text-neutral hover:transition-colors']"
                                                :aria-current="page.current ? 'page' : undefined">{{ page.name
                                                }}</NuxtLink>
                                        </div>
                                    </li>
                                </ol>
                            </nav>
                        </div>
                        <div class="flex items-center ml-2 space-x-4 sm:ml-6 sm:space-x-6">
                            <div v-can="'AuthService.SetJob'">
                                <SidebarJobSwitcher v-if="activeChar" />
                            </div>
                            <div v-if="activeChar" class="text-sm font-medium text-base-400">
                                {{ activeChar.getFirstname() }}, {{ activeChar.getLastname() }}
                                ({{ activeChar.getJobLabel() }})
                            </div>
                            <SidebarLanguageSwitcher />
                            <!-- Profile dropdown -->
                            <Menu as="div" class="relative flex-shrink-0">
                                <div>
                                    <MenuButton
                                        class="flex text-sm rounded-full bg-base-850 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                                        <span class="sr-only">{{ $t('components.partials.sidebar.open_usermenu') }}</span>
                                        <UserIcon
                                            class="w-auto h-10 rounded-full hover:transition-colors text-base-300 bg-base-800 fill-base-300 hover:text-base-100 hover:fill-base-100" />
                                    </MenuButton>
                                </div>
                                <transition enter-active-class="transition duration-100 ease-out"
                                    enter-from-class="transform scale-95 opacity-0"
                                    enter-to-class="transform scale-100 opacity-100"
                                    leave-active-class="transition duration-75 ease-in"
                                    leave-from-class="transform scale-100 opacity-100"
                                    leave-to-class="transform scale-95 opacity-0">
                                    <MenuItems
                                        class="absolute right-0 w-48 py-1 mt-2 origin-top-right rounded-md shadow-float bg-base-850 ring-1 ring-base-100 ring-opacity-5 focus:outline-none z-40">
                                        <MenuItem v-for="item in userNavigation" :key="item.name" v-slot="{ active }"
                                            v-can="item.permission">
                                        <NuxtLink :to="item.href"
                                            :class="[active ? 'bg-base-800' : '', 'block px-4 py-2 text-sm text-neutral hover:transition-colors']">
                                            {{ item.name }}
                                        </NuxtLink>
                                        </MenuItem>
                                    </MenuItems>
                                </transition>
                            </Menu>
                        </div>
                    </div>
                </div>
            </header>

            <main class="h-full overflow-y-auto">
                <section aria-labelledby="primary-heading" class="h-full min-w-0 lg:order-last">
                    <slot></slot>
                </section>
            </main>
        </div>
    </div>
</template>
