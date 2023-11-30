<script lang="ts" setup>
import { Dialog, DialogPanel, Menu, MenuButton, MenuItem, MenuItems, TransitionChild, TransitionRoot } from '@headlessui/vue';
import {
    AccountPlusIcon,
    AccountIcon,
    AccountMultipleIcon,
    BriefcaseIcon,
    CarEmergencyIcon,
    CarIcon,
    ChevronRightIcon,
    CloseIcon,
    CogIcon,
    FileDocumentMultipleIcon,
    HelpCircleIcon,
    HomeIcon,
    LoginIcon,
    MapIcon,
    MenuIcon,
    UnfoldMoreHorizontalIcon,
} from 'mdi-vue3';
import { type DefineComponent } from 'vue';
import QuickButtons from '~/components/partials/QuickButtons.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import JobSwitcher from '~/components/partials/sidebar/JobSwitcher.vue';
import LanguageSwitcher from '~/components/partials/sidebar/LanguageSwitcher.vue';
import Notifications from '~/components/partials/sidebar/Notifications.vue';
import { useAuthStore } from '~/store/auth';
import { type RoutesNamedLocations } from '@typed-router';

const authStore = useAuthStore();
const { accessToken, activeChar } = storeToRefs(authStore);
const router = useRouter();
const route = useRoute();

const sidebarNavigation = ref<
    {
        name: string;
        href: RoutesNamedLocations;
        activePath?: string;
        permission: string;
        icon: DefineComponent;
        position: 'top' | 'bottom';
        current: boolean;
        loggedIn?: boolean;
        charSelected?: boolean;
    }[]
>([
    {
        name: 'common.overview',
        href: { name: 'overview' },
        permission: '',
        icon: markRaw(HomeIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.citizen',
        href: { name: 'citizens' },
        permission: 'CitizenStoreService.ListCitizens',
        icon: markRaw(AccountMultipleIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.vehicle',
        href: { name: 'vehicles' },
        permission: 'DMVService.ListVehicles',
        icon: markRaw(CarIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.document',
        href: { name: 'documents' },
        permission: 'DocStoreService.ListDocuments',
        icon: markRaw(FileDocumentMultipleIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.job',
        href: { name: 'jobs-overview' },
        activePath: '/jobs',
        permission: 'JobsService.ColleaguesList',
        icon: markRaw(BriefcaseIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.livemap',
        href: { name: 'livemap' },
        permission: 'LivemapperService.Stream',
        icon: markRaw(MapIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.dispatch_center',
        href: { name: 'centrum' },
        permission: 'CentrumService.TakeControl',
        icon: markRaw(CarEmergencyIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.control_panel',
        href: { name: 'rector' },
        permission: 'RectorService.GetRoles',
        icon: markRaw(CogIcon),
        position: 'top',
        current: false,
    },

    {
        name: 'common.about',
        href: { name: 'about' },
        permission: '',
        icon: markRaw(HelpCircleIcon),
        position: 'bottom',
        current: false,
    },
]);
const userNavigation = ref<{ name: string; href: RoutesNamedLocations; permission?: string }[]>([
    { name: 'components.auth.login.title', href: { name: 'auth-login' } },
    { name: 'components.auth.registration_form.title', href: { name: 'auth-registration' } },
]);
const breadcrumbs = useBreadcrumbs();
const mobileMenuOpen = ref(false);

onMounted(() => {
    updateUserNav();
    updateActiveItem();
});

function updateUserNav(): void {
    userNavigation.value.length = 0;
    if (activeChar.value) {
        userNavigation.value.push({
            name: 'components.partials.sidebar.change_character',
            href: { name: 'auth-character-selector' },
        });
    }
    if (accessToken.value) {
        userNavigation.value.push(
            {
                name: 'components.partials.sidebar.account_info',
                href: { name: 'auth-account-info' },
            },
            { name: 'common.sign_out', href: { name: 'auth-logout' } },
        );
    }
    if (userNavigation.value.length === 0) {
        userNavigation.value = [
            { name: 'components.auth.login.title', href: { name: 'auth-login' } },
            { name: 'components.auth.registration_form.title', href: { name: 'auth-registration' } },
        ];
    }
}

function updateActiveItem(): void {
    const route = router.currentRoute.value;
    if (route.path) {
        sidebarNavigation.value.forEach((e) => {
            const itemRoute = useRouter().resolve(e.href);
            if (
                route.path.toLowerCase().startsWith(itemRoute.path.toLowerCase()) ||
                (e.activePath && route.path.toLowerCase().startsWith(e.activePath))
            ) {
                e.current = true;
            } else {
                e.current = false;
            }
        });
    } else {
        sidebarNavigation.value.forEach((e) => (e.current = false));
    }
}

watch(accessToken, () => updateUserNav());
watch(activeChar, () => updateUserNav());

watch(router.currentRoute, () => updateActiveItem());
</script>

<template>
    <div class="flex h-dscreen">
        <!-- Sidebar -->
        <div class="hidden overflow-y-auto bg-accent-600 w-28 md:block">
            <div class="flex flex-col items-center w-full py-6 h-full">
                <div class="flex items-center flex-shrink-0">
                    <NuxtLink :to="{ name: accessToken ? 'overview' : 'index' }" aria-current-value="page">
                        <FiveNetLogo class="w-auto h-12" />
                    </NuxtLink>
                </div>
                <div class="flex-grow w-full px-2 mt-6 space-y-1 text-center">
                    <template v-if="!accessToken && !activeChar">
                        <NuxtLink
                            :to="{ name: 'index' }"
                            active-class="bg-accent-100/20 text-neutral font-bold"
                            class="text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2"
                            exact-active-class="text-neutral"
                            aria-current-value="page"
                        >
                            <HomeIcon class="h-6 w-6" aria-hidden="true" />
                            <span class="mt-2">{{ $t('common.home') }}</span>
                        </NuxtLink>
                        <NuxtLink
                            :to="{ name: 'auth-login' }"
                            active-class="bg-accent-100/20 text-neutral font-bold"
                            class="text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2"
                            exact-active-class="text-neutral"
                            aria-current-value="page"
                        >
                            <LoginIcon class="h-6 w-6" aria-hidden="true" />
                            <span class="mt-2">{{ $t('components.auth.login.title') }}</span>
                        </NuxtLink>
                        <NuxtLink
                            :to="{ name: 'auth-registration' }"
                            active-class="bg-accent-100/20 text-neutral font-bold"
                            class="text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2"
                            exact-active-class="text-neutral"
                            aria-current-value="page"
                        >
                            <AccountPlusIcon class="h-6 w-6" aria-hidden="true" />
                            <span class="mt-2">{{ $t('components.auth.registration_form.title') }}</span>
                        </NuxtLink>
                    </template>
                    <template v-if="accessToken && !activeChar">
                        <NuxtLink
                            :to="{ name: 'index' }"
                            active-class="bg-accent-100/20 text-neutral font-bold"
                            class="text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2"
                            exact-active-class="text-neutral"
                            aria-current-value="page"
                        >
                            <HomeIcon class="h-6 w-6" aria-hidden="true" />
                            <span class="mt-2">{{ $t('common.home') }}</span>
                        </NuxtLink>
                        <NuxtLink
                            :to="{ name: 'auth-character-selector' }"
                            active-class="bg-accent-100/20 text-neutral font-bold"
                            class="text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2"
                            exact-active-class="text-neutral"
                            aria-current-value="page"
                        >
                            <UnfoldMoreHorizontalIcon class="h-6 w-6" aria-hidden="true" />
                            <span class="mt-2">{{ $t('components.auth.character_selector.title') }}</span>
                        </NuxtLink>
                    </template>
                    <template v-else-if="accessToken && activeChar">
                        <NuxtLink
                            v-for="item in sidebarNavigation.filter((e) => e.position === 'top' && can(e.permission))"
                            :key="item.name"
                            :to="item.href"
                            :class="[
                                item.current
                                    ? 'bg-accent-100/20 text-neutral font-bold'
                                    : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium',
                                'hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2',
                            ]"
                            :aria-current="item.current ? 'page' : undefined"
                        >
                            <component
                                :is="item.icon"
                                :class="[item.current ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral', 'h-6 w-6']"
                                aria-hidden="true"
                            />
                            <span class="mt-2">{{ $t(item.name) }}</span>
                        </NuxtLink>
                    </template>
                </div>
                <div class="flex-initial w-full px-2 space-y-1 text-center">
                    <NuxtLink
                        v-for="item in sidebarNavigation.filter((e) => e.position === 'bottom' && can(e.permission))"
                        :key="item.name"
                        :to="item.href"
                        :class="[
                            item.current
                                ? 'bg-accent-100/20 text-neutral font-bold'
                                : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium',
                            'hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2',
                        ]"
                        :aria-current="item.current ? 'page' : undefined"
                    >
                        <component
                            :is="item.icon"
                            :class="[item.current ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral', 'h-6 w-6']"
                            aria-hidden="true"
                        />
                        <span class="mt-2">{{ $t(item.name) }}</span>
                    </NuxtLink>
                </div>
            </div>
        </div>

        <!-- Mobile Sidebar -->
        <TransitionRoot as="template" :show="mobileMenuOpen">
            <Dialog as="div" class="relative z-30 md:hidden" @close="mobileMenuOpen = false">
                <TransitionChild
                    as="template"
                    enter="transition-opacity ease-linear duration-300"
                    enter-from="opacity-0"
                    enter-to="opacity-100"
                    leave="transition-opacity ease-linear duration-300"
                    leave-from="opacity-100"
                    leave-to="opacity-0"
                >
                    <div class="fixed inset-0 bg-opacity-75 bg-base-900/10" />
                </TransitionChild>

                <div class="fixed inset-0 z-30 flex">
                    <TransitionChild
                        as="template"
                        enter="transition ease-in-out duration-300 transform"
                        enter-from="-translate-x-full"
                        enter-to="translate-x-0"
                        leave="transition ease-in-out duration-300 transform"
                        leave-from="translate-x-0"
                        leave-to="-translate-x-full"
                    >
                        <DialogPanel class="relative flex flex-col flex-1 w-full max-w-xs pt-5 pb-4 bg-accent-600">
                            <TransitionChild
                                as="template"
                                enter="ease-in-out duration-300"
                                enter-from="opacity-0"
                                enter-to="opacity-100"
                                leave="ease-in-out duration-300"
                                leave-from="opacity-100"
                                leave-to="opacity-0"
                            >
                                <div class="absolute p-1 -right-3 top-1 -mr-14">
                                    <button
                                        type="button"
                                        class="flex items-center justify-center w-12 h-12 rounded-full focus:outline-none ring-2 ring-neutral"
                                        @click="mobileMenuOpen = false"
                                    >
                                        <CloseIcon class="w-6 h-6 text-neutral" aria-hidden="true" />
                                        <span class="sr-only">{{ $t('components.partials.sidebar.close_sidebar') }}</span>
                                    </button>
                                </div>
                            </TransitionChild>
                            <div class="flex items-center flex-shrink-0 px-4">
                                <FiveNetLogo class="w-16 h-16 mx-auto" />
                            </div>
                            <div class="flex-grow h-0 px-2 mt-5 overflow-y-auto">
                                <nav class="flex flex-col h-full">
                                    <div class="space-y-1">
                                        <template v-if="!accessToken && !activeChar">
                                            <NuxtLink
                                                :to="{ name: 'index' }"
                                                active-class="bg-accent-100/20 text-neutral font-bold"
                                                class="text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium group flex items-center rounded-md py-2 px-3 text-sm"
                                                exact-active-class="text-neutral"
                                                aria-current-value="page"
                                                @click="mobileMenuOpen = false"
                                            >
                                                <HomeIcon
                                                    class="text-accent-100 group-hover:text-neutral mr-3 h-6 w-6"
                                                    aria-hidden="true"
                                                />
                                                <span>{{ $t('common.home') }}</span>
                                            </NuxtLink>
                                            <NuxtLink
                                                :to="{ name: 'auth-login' }"
                                                active-class="bg-accent-100/20 text-neutral font-bold"
                                                class="text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium group flex items-center rounded-md py-2 px-3 text-sm"
                                                exact-active-class="text-neutral"
                                                aria-current-value="page"
                                                @click="mobileMenuOpen = false"
                                            >
                                                <LoginIcon
                                                    class="text-accent-100 group-hover:text-neutral mr-3 h-6 w-6"
                                                    aria-hidden="true"
                                                />
                                                <span>{{ $t('components.auth.login.title') }}</span>
                                            </NuxtLink>
                                            <NuxtLink
                                                :to="{ name: 'auth-registration' }"
                                                active-class="bg-accent-100/20 text-neutral font-bold"
                                                class="text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium group flex items-center rounded-md py-2 px-3 text-sm"
                                                exact-active-class="text-neutral"
                                                aria-current-value="page"
                                                @click="mobileMenuOpen = false"
                                            >
                                                <AccountPlusIcon
                                                    class="text-accent-100 group-hover:text-neutral mr-3 h-6 w-6"
                                                    aria-hidden="true"
                                                />
                                                <span>{{ $t('components.auth.registration_form.title') }}</span>
                                            </NuxtLink>
                                        </template>
                                        <NuxtLink
                                            v-for="item in sidebarNavigation.filter(
                                                (e) => e.position === 'top' && can(e.permission),
                                            )"
                                            v-else-if="accessToken && activeChar"
                                            :key="item.name"
                                            :to="item.href"
                                            :class="[
                                                item.current
                                                    ? 'bg-accent-100/20 text-neutral font-bold'
                                                    : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium',
                                                'group flex items-center rounded-md py-2 px-3 text-sm',
                                            ]"
                                            :aria-current="item.current ? 'page' : undefined"
                                            @click="mobileMenuOpen = false"
                                        >
                                            <component
                                                :is="item.icon"
                                                :class="[
                                                    item.current ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral',
                                                    'mr-3 h-6 w-6',
                                                ]"
                                                aria-hidden="true"
                                            />
                                            <span>{{ $t(item.name, 2) }}</span>
                                        </NuxtLink>
                                    </div>
                                </nav>
                            </div>
                            <div class="flex-initial px-2 overflow-y-auto">
                                <nav class="flex flex-col h-full">
                                    <div class="space-y-1">
                                        <NuxtLink
                                            v-for="item in sidebarNavigation.filter(
                                                (e) => e.position === 'bottom' && can(e.permission),
                                            )"
                                            :key="item.name"
                                            :to="item.href"
                                            :class="[
                                                item.current
                                                    ? 'bg-accent-100/20 text-neutral font-bold'
                                                    : 'text-accent-100 hover:bg-accent-100/10 hover:text-neutral font-medium',
                                                'group flex items-center rounded-md py-2 px-3 text-sm',
                                            ]"
                                            :aria-current="item.current ? 'page' : undefined"
                                            @click="mobileMenuOpen = false"
                                        >
                                            <component
                                                :is="item.icon"
                                                :class="[
                                                    item.current ? 'text-neutral' : 'text-accent-100 group-hover:text-neutral',
                                                    'mr-3 h-6 w-6',
                                                ]"
                                                aria-hidden="true"
                                            />
                                            <span>{{ $t(item.name, 2) }}</span>
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
                <div class="relative z-30 flex flex-shrink-0 h-16 bg-base-800">
                    <button
                        type="button"
                        class="px-4 text-neutral focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500 md:hidden"
                        @click="mobileMenuOpen = true"
                    >
                        <span class="sr-only">{{ $t('components.partials.sidebar.open_sidebar') }}</span>
                        <MenuIcon class="w-6 h-6" aria-hidden="true" />
                    </button>
                    <div class="flex justify-between flex-1 px-2 sm:px-6">
                        <div class="flex flex-1">
                            <nav class="flex" aria-label="Breadcrumb">
                                <ol role="list" class="flex items-center space-x-2 sm:space-x-4">
                                    <li>
                                        <div>
                                            <NuxtLink
                                                :to="{
                                                    name: accessToken ? 'overview' : 'index',
                                                }"
                                                class="text-base-400 hover:text-neutral hover:transition-colors"
                                            >
                                                <HomeIcon class="flex-shrink-0 w-5 h-5" aria-hidden="true" />
                                                <span class="sr-only">Home</span>
                                            </NuxtLink>
                                        </div>
                                    </li>
                                    <template v-for="(item, key) in breadcrumbs" :key="key">
                                        <li v-if="key !== 0">
                                            <div class="flex items-center">
                                                <ChevronRightIcon
                                                    class="flex-shrink-0 w-5 h-5 text-base-400"
                                                    aria-hidden="true"
                                                />
                                                <NuxtLink
                                                    :to="item.to"
                                                    :class="[
                                                        key === breadcrumbs.length - 1
                                                            ? 'font-bold text-base-200'
                                                            : 'font-medium text-base-400',
                                                        'ml-2 sm:ml-4 text-sm hover:text-neutral hover:transition-colors truncate max-w-[5rem] lg:max-w-full',
                                                    ]"
                                                    :aria-current="key === breadcrumbs.length - 1 ? 'page' : undefined"
                                                >
                                                    {{ $t(item.title as string) }}
                                                </NuxtLink>
                                            </div>
                                        </li>
                                    </template>
                                </ol>
                            </nav>
                        </div>
                        <div class="flex items-center ml-2 space-x-3 sm:ml-2 sm:space-x-4">
                            <JobSwitcher v-if="can('SuperUser') && activeChar" />
                            <div v-if="activeChar" class="hidden sm:block text-sm font-medium text-base-200">
                                {{ activeChar.firstname }}, {{ activeChar.lastname }} ({{ activeChar.jobLabel }})
                            </div>
                            <Notifications v-if="activeChar" />
                            <LanguageSwitcher />

                            <!-- Account dropdown -->
                            <Menu as="div" class="relative flex-shrink-0">
                                <MenuButton
                                    class="flex text-sm rounded-full bg-base-800 ring-2 ring-base-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                >
                                    <span class="sr-only">
                                        {{ $t('components.partials.sidebar.open_usermenu') }}
                                    </span>
                                    <AccountIcon
                                        class="w-auto h-10 rounded-full hover:transition-colors text-base-300 bg-base-800 fill-base-300 hover:text-base-100 hover:fill-base-100"
                                    />
                                </MenuButton>
                                <transition
                                    enter-active-class="transition duration-100 ease-out"
                                    enter-from-class="transform scale-95 opacity-0"
                                    enter-to-class="transform scale-100 opacity-100"
                                    leave-active-class="transition duration-75 ease-in"
                                    leave-from-class="transform scale-100 opacity-100"
                                    leave-to-class="transform scale-95 opacity-0"
                                >
                                    <MenuItems
                                        class="absolute z-30 right-0 w-48 py-1 mt-2 origin-top-right rounded-md shadow-float bg-base-800 ring-1 ring-base-100 ring-opacity-5 focus:outline-none"
                                    >
                                        <MenuItem
                                            v-for="item in userNavigation.filter(
                                                (e) => e.permission === undefined || can(e.permission),
                                            )"
                                            :key="item.name"
                                            v-slot="{ close }"
                                        >
                                            <NuxtLink
                                                :to="item.href"
                                                class="block px-4 py-2 text-sm text-neutral hover:transition-colors"
                                                active-class="bg-primary-500"
                                                @mouseup="close"
                                            >
                                                {{ $t(item.name) }}
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

                    <QuickButtons
                        v-if="activeChar && (route.meta.showQuickButtons === undefined || route.meta.showQuickButtons)"
                    />
                </section>
            </main>
        </div>
    </div>
</template>
