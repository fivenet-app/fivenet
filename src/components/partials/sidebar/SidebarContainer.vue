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
import { type RoutesNamedLocations } from '@typed-router';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';
import LanguageSwitcherMenu from '~/components/partials/sidebar/LanguageSwitcherMenu.vue';
import { useAuthStore } from '~/store/auth';
import type { Perms } from '~~/gen/ts/perms';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';

const authStore = useAuthStore();
const { accessToken, activeChar } = storeToRefs(authStore);

const router = useRouter();

const sidebarNavigation = ref<
    {
        name: string;
        to: RoutesNamedLocations;
        activePath?: string;
        permission?: Perms;
        icon: DefineComponent;
        position: 'top' | 'bottom';
        current: boolean;
        loggedIn?: boolean;
        charSelected?: boolean;
    }[]
>([
    {
        name: 'common.overview',
        to: { name: 'overview' },
        icon: markRaw(HomeIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.citizen',
        to: { name: 'citizens' },
        permission: 'CitizenStoreService.ListCitizens',
        icon: markRaw(AccountMultipleIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.vehicle',
        to: { name: 'vehicles' },
        permission: 'DMVService.ListVehicles',
        icon: markRaw(CarIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.document',
        to: { name: 'documents' },
        permission: 'DocStoreService.ListDocuments',
        icon: markRaw(FileDocumentMultipleIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.job',
        to: { name: 'jobs-overview' },
        activePath: '/jobs',
        permission: 'JobsService.ListColleagues',
        icon: markRaw(BriefcaseIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.livemap',
        to: { name: 'livemap' },
        permission: 'LivemapperService.Stream',
        icon: markRaw(MapIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.dispatch_center',
        to: { name: 'centrum' },
        permission: 'CentrumService.TakeControl',
        icon: markRaw(CarEmergencyIcon),
        position: 'top',
        current: false,
    },
    {
        name: 'common.control_panel',
        to: { name: 'rector' },
        permission: 'RectorService.GetRoles',
        icon: markRaw(CogIcon),
        position: 'top',
        current: false,
    },

    {
        name: 'common.about',
        to: { name: 'about' },
        icon: markRaw(HelpCircleIcon),
        position: 'bottom',
        current: false,
    },
]);
const userNavigation = ref<{ name: string; to: RoutesNamedLocations; permission?: Perms }[]>([
    { name: 'components.auth.login.title', to: { name: 'auth-login' } },
    { name: 'components.auth.registration_form.title', to: { name: 'auth-registration' } },
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
            to: { name: 'auth-character-selector' },
        });
    }
    if (accessToken.value) {
        userNavigation.value.push(
            {
                name: 'components.auth.account_info.title',
                to: { name: 'auth-account-info' },
            },
            {
                name: 'components.auth.settings_panel.title',
                to: { name: 'settings' },
            },
            { name: 'common.sign_out', to: { name: 'auth-logout' } },
        );
    }
    if (userNavigation.value.length === 0) {
        userNavigation.value = [
            { name: 'components.auth.login.title', to: { name: 'auth-login' } },
            { name: 'components.auth.registration_form.title', to: { name: 'auth-registration' } },
        ];
    }
}

function updateActiveItem(): void {
    const route = router.currentRoute.value;
    if (route.path) {
        sidebarNavigation.value.forEach((e) => {
            const itemRoute = useRouter().resolve(e.to);
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
    <div class="h-dscreen flex">
        <!-- Sidebar -->
        <div class="hidden w-28 overflow-y-auto defaultTheme:bg-accent-600 md:block">
            <div class="flex size-full flex-col items-center py-6">
                <div class="flex shrink-0 items-center">
                    <NuxtLink :to="{ name: accessToken ? 'overview' : 'index' }" aria-current-value="page">
                        <FiveNetLogo class="h-12 w-auto" />
                    </NuxtLink>
                </div>
                <div class="mt-6 w-full grow space-y-1 px-2 text-center">
                    <template v-if="!accessToken && !activeChar">
                        <NuxtLink
                            :to="{ name: 'index' }"
                            active-class="bg-accent-100/20 font-bold"
                            class="group my-2 flex w-full flex-col items-center rounded-md p-3 text-xs font-medium text-accent-100 hover:bg-accent-100/10 hover:transition-all"
                            exact-active-class="text-accent-200"
                            aria-current-value="page"
                        >
                            <HomeIcon class="h-auto w-6" />
                            <span class="mt-2">{{ $t('common.home') }}</span>
                        </NuxtLink>
                        <NuxtLink
                            :to="{ name: 'auth-login' }"
                            active-class="bg-accent-100/20 font-bold"
                            class="group my-2 flex w-full flex-col items-center rounded-md p-3 text-xs font-medium text-accent-100 hover:bg-accent-100/10 hover:transition-all"
                            exact-active-class="text-accent-200"
                            aria-current-value="page"
                        >
                            <LoginIcon class="h-auto w-6" />
                            <span class="mt-2">{{ $t('components.auth.login.title') }}</span>
                        </NuxtLink>
                        <NuxtLink
                            :to="{ name: 'auth-registration' }"
                            active-class="bg-accent-100/20 font-bold"
                            class="group my-2 flex w-full flex-col items-center rounded-md p-3 text-xs font-medium text-accent-100 hover:bg-accent-100/10 hover:transition-all"
                            exact-active-class="text-accent-200"
                            aria-current-value="page"
                        >
                            <AccountPlusIcon class="h-auto w-6" />
                            <span class="mt-2">{{ $t('components.auth.registration_form.title') }}</span>
                        </NuxtLink>
                    </template>
                    <template v-if="accessToken && !activeChar">
                        <NuxtLink
                            :to="{ name: 'index' }"
                            active-class="bg-accent-100/20 font-bold"
                            class="group my-2 flex w-full flex-col items-center rounded-md p-3 text-xs font-medium text-accent-100 hover:bg-accent-100/10 hover:transition-all"
                            exact-active-class="text-accent-200"
                            aria-current-value="page"
                        >
                            <HomeIcon class="h-auto w-6" />
                            <span class="mt-2">{{ $t('common.home') }}</span>
                        </NuxtLink>
                        <NuxtLink
                            :to="{ name: 'auth-character-selector' }"
                            active-class="bg-accent-100/20 font-bold"
                            class="group my-2 flex w-full flex-col items-center rounded-md p-3 text-xs font-medium text-accent-100 hover:bg-accent-100/10 hover:transition-all"
                            exact-active-class="text-accent-200"
                            aria-current-value="page"
                        >
                            <UnfoldMoreHorizontalIcon class="h-auto w-6" />
                            <span class="mt-2">{{ $t('components.auth.character_selector.title') }}</span>
                        </NuxtLink>
                    </template>
                    <template v-else-if="accessToken && activeChar">
                        <NuxtLink
                            v-for="item in sidebarNavigation.filter(
                                (e) => e.position === 'top' && (e.permission === undefined || can(e.permission)),
                            )"
                            :key="item.name"
                            :to="item.to"
                            :class="[
                                item.current
                                    ? 'bg-accent-100/20 font-bold'
                                    : 'font-medium text-accent-100 hover:bg-accent-100/10 ',
                                'group my-2 flex w-full flex-col items-center rounded-md p-3 text-xs hover:transition-all',
                            ]"
                            :aria-current="item.current ? 'page' : undefined"
                        >
                            <component
                                :is="item.icon"
                                :class="[item.current ? 'text-neutral' : 'text-accent-100 group-', 'h-auto w-6']"
                            />
                            <span class="mt-2">{{ $t(item.name) }}</span>
                        </NuxtLink>
                    </template>
                </div>
                <div class="w-full flex-initial space-y-1 px-2 text-center">
                    <NuxtLink
                        v-for="item in sidebarNavigation.filter(
                            (e) => e.position === 'bottom' && (e.permission === undefined || can(e.permission)),
                        )"
                        :key="item.name"
                        :to="item.to"
                        :class="[
                            item.current ? 'bg-accent-100/20 font-bold' : 'font-medium text-accent-100 hover:bg-accent-100/10 ',
                            'group my-2 flex w-full flex-col items-center rounded-md p-3 text-xs hover:transition-all',
                        ]"
                        :aria-current="item.current ? 'page' : undefined"
                    >
                        <component
                            :is="item.icon"
                            :class="[item.current ? 'text-neutral' : 'text-accent-100 group-', 'h-auto w-6']"
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
                    <div class="fixed inset-0 bg-base-900/75" />
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
                        <DialogPanel class="relative flex w-full max-w-xs flex-1 flex-col pb-4 pt-5 defaultTheme:bg-accent-600">
                            <TransitionChild
                                as="template"
                                enter="ease-in-out duration-300"
                                enter-from="opacity-0"
                                enter-to="opacity-100"
                                leave="ease-in-out duration-300"
                                leave-from="opacity-100"
                                leave-to="opacity-0"
                            >
                                <div class="absolute -right-3 top-1 -mr-14 p-1">
                                    <UButton
                                        class="flex size-12 items-center justify-center rounded-full ring-2 ring-neutral"
                                        @click="mobileMenuOpen = false"
                                    >
                                        <CloseIcon class="h-auto w-6" />
                                        <span class="sr-only">{{ $t('components.partials.sidebar.close_sidebar') }}</span>
                                    </UButton>
                                </div>
                            </TransitionChild>
                            <div class="flex shrink-0 items-center px-4">
                                <FiveNetLogo class="mx-auto size-16" />
                            </div>
                            <div class="mt-5 h-0 grow overflow-y-auto px-2">
                                <nav class="flex h-full flex-col">
                                    <div class="space-y-1">
                                        <template v-if="!accessToken && !activeChar">
                                            <NuxtLink
                                                :to="{ name: 'index' }"
                                                active-class="bg-accent-100/20 font-bold"
                                                class="group flex items-center rounded-md px-3 py-2 text-sm font-medium text-accent-100 hover:bg-accent-100/10"
                                                exact-active-class="text-accent-200"
                                                aria-current-value="page"
                                                @click="mobileMenuOpen = false"
                                            >
                                                <HomeIcon class="mr-3 h-auto w-6 text-accent-100 group-" />
                                                <span>{{ $t('common.home') }}</span>
                                            </NuxtLink>
                                            <NuxtLink
                                                :to="{ name: 'auth-login' }"
                                                active-class="bg-accent-100/20 font-bold"
                                                class="group flex items-center rounded-md px-3 py-2 text-sm font-medium text-accent-100 hover:bg-accent-100/10"
                                                exact-active-class="text-accent-200"
                                                aria-current-value="page"
                                                @click="mobileMenuOpen = false"
                                            >
                                                <LoginIcon class="mr-3 h-auto w-6 text-accent-100 group-" />
                                                <span>{{ $t('components.auth.login.title') }}</span>
                                            </NuxtLink>
                                            <NuxtLink
                                                :to="{ name: 'auth-registration' }"
                                                active-class="bg-accent-100/20 font-bold"
                                                class="group flex items-center rounded-md px-3 py-2 text-sm font-medium text-accent-100 hover:bg-accent-100/10"
                                                exact-active-class="text-accent-200"
                                                aria-current-value="page"
                                                @click="mobileMenuOpen = false"
                                            >
                                                <AccountPlusIcon class="mr-3 h-auto w-6 text-accent-100 group-" />
                                                <span>{{ $t('components.auth.registration_form.title') }}</span>
                                            </NuxtLink>
                                        </template>
                                        <NuxtLink
                                            v-for="item in sidebarNavigation.filter(
                                                (e) =>
                                                    e.position === 'top' && (e.permission === undefined || can(e.permission)),
                                            )"
                                            v-else-if="accessToken && activeChar"
                                            :key="item.name"
                                            :to="item.to"
                                            :class="[
                                                item.current
                                                    ? 'bg-accent-100/20 font-bold'
                                                    : 'font-medium text-accent-100 hover:bg-accent-100/10 ',
                                                'group flex items-center rounded-md px-3 py-2 text-sm',
                                            ]"
                                            :aria-current="item.current ? 'page' : undefined"
                                            @click="mobileMenuOpen = false"
                                        >
                                            <component
                                                :is="item.icon"
                                                :class="[
                                                    item.current ? 'text-neutral' : 'text-accent-100 group-',
                                                    'mr-3 h-auto w-6',
                                                ]"
                                            />
                                            <span>{{ $t(item.name, 2) }}</span>
                                        </NuxtLink>
                                    </div>
                                </nav>
                            </div>
                            <div class="flex-initial overflow-y-auto px-2">
                                <nav class="flex h-full flex-col">
                                    <div class="space-y-1">
                                        <NuxtLink
                                            v-for="item in sidebarNavigation.filter(
                                                (e) =>
                                                    e.position === 'bottom' &&
                                                    (e.permission === undefined || can(e.permission)),
                                            )"
                                            :key="item.name"
                                            :to="item.to"
                                            :class="[
                                                item.current
                                                    ? 'bg-accent-100/20 font-bold'
                                                    : 'font-medium text-accent-100 hover:bg-accent-100/10 ',
                                                'group flex items-center rounded-md px-3 py-2 text-sm',
                                            ]"
                                            :aria-current="item.current ? 'page' : undefined"
                                            @click="mobileMenuOpen = false"
                                        >
                                            <component
                                                :is="item.icon"
                                                :class="[
                                                    item.current ? 'text-neutral' : 'text-accent-100 group-',
                                                    'mr-3 h-auto w-6',
                                                ]"
                                            />
                                            <span>{{ $t(item.name, 2) }}</span>
                                        </NuxtLink>
                                    </div>
                                </nav>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                    <div class="w-14 shrink-0"></div>
                </div>
            </Dialog>
        </TransitionRoot>

        <!-- Content area -->
        <div class="flex flex-1 flex-col overflow-hidden">
            <header class="w-full">
                <div class="relative z-30 flex h-16 shrink-0 bg-base-800">
                    <UButton
                        class="px-4 focus:ring-2 focus:ring-inset focus:ring-primary-500 md:hidden"
                        @click="mobileMenuOpen = true"
                    >
                        <span class="sr-only">{{ $t('components.partials.sidebar.open_sidebar') }}</span>
                        <MenuIcon class="h-auto w-6" />
                    </UButton>
                    <div class="flex flex-1 justify-between px-2 sm:px-6">
                        <div class="flex flex-1">
                            <nav class="flex" aria-label="Breadcrumb">
                                <ol role="list" class="flex items-center space-x-2 sm:space-x-4">
                                    <li>
                                        <div>
                                            <NuxtLink
                                                :to="{
                                                    name: accessToken ? 'overview' : 'index',
                                                }"
                                                class="text-base-400 hover:transition-colors"
                                            >
                                                <HomeIcon class="size-5 shrink-0" />
                                                <span class="sr-only">{{ $t('common.home') }}</span>
                                            </NuxtLink>
                                        </div>
                                    </li>
                                    <template v-for="(page, key) in breadcrumbs" :key="key">
                                        <li v-if="key !== 0 && page.to !== undefined">
                                            <div class="flex items-center">
                                                <ChevronRightIcon class="size-5 shrink-0 text-base-300" />
                                                <!-- @vue-ignore the route should be valid, as we construct it based on our pages -->
                                                <NuxtLink
                                                    :to="page.to"
                                                    :class="[
                                                        key === breadcrumbs.length - 1
                                                            ? 'font-bold text-accent-200'
                                                            : 'font-medium text-base-300',
                                                        'ml-2 max-w-20 truncate text-sm  hover:transition-colors sm:ml-4 lg:max-w-full',
                                                    ]"
                                                    :aria-current="key === breadcrumbs.length - 1 ? 'page' : undefined"
                                                >
                                                    {{ $t(page.title as string) }}
                                                </NuxtLink>
                                            </div>
                                        </li>
                                    </template>
                                </ol>
                            </nav>
                        </div>
                        <div class="ml-2 flex items-center space-x-3 sm:ml-2 sm:space-x-4">
                            <template v-if="activeChar">
                                <div v-if="activeChar" class="hidden text-center text-sm font-medium text-accent-200 sm:block">
                                    <span class="line-clamp-3">
                                        <span class="line-clamp-2">{{ activeChar.firstname }} {{ activeChar.lastname }}</span>
                                        ({{ activeChar.jobLabel }})
                                    </span>
                                </div>
                            </template>
                            <LanguageSwitcherMenu />

                            <!-- Account dropdown -->
                            <Menu as="div" class="relative shrink-0">
                                <MenuButton
                                    class="flex rounded-full bg-base-800 text-sm ring-2 ring-base-600 focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                >
                                    <span class="sr-only">
                                        {{ $t('components.partials.sidebar.open_usermenu') }}
                                    </span>
                                    <AccountIcon
                                        v-if="!activeChar?.avatar?.url"
                                        class="h-10 w-auto rounded-full bg-base-800 fill-base-300 text-base-300 hover:fill-base-100 hover:transition-colors"
                                    />
                                    <ProfilePictureImg
                                        v-else
                                        :url="activeChar.avatar.url"
                                        :name="`${activeChar.firstname} ${activeChar.lastname}`"
                                        size="md"
                                        :rounded="true"
                                        :no-blur="true"
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
                                        class="absolute right-0 z-30 mt-2 w-48 origin-top-right rounded-md bg-base-800 py-1 shadow-float ring-1 ring-base-100/5"
                                    >
                                        <MenuItem
                                            v-for="item in userNavigation.filter(
                                                (e) => e.permission === undefined || can(e.permission),
                                            )"
                                            :key="item.name"
                                            v-slot="{ close }"
                                        >
                                            <NuxtLink
                                                :to="item.to"
                                                class="block px-4 py-2 text-sm hover:transition-colors"
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
                    <slot />
                </section>
            </main>
        </div>
    </div>
</template>
