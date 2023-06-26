<script lang="ts" setup>
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiChevronRight, mdiDotsVertical } from '@mdi/js';

const statuses: { [key: string]: string } = {
    offline: 'animate-pulse text-gray-500 bg-gray-100/10',
    online: 'text-green-400 bg-green-400/10',
    error: 'text-rose-400 bg-rose-400/10',
};

const projects = [
    {
        id: 1,
        href: '#',
        name: 'ios-app',
        job: 'LSPD',
        status: 'offline',
        statusText: 'Initiated 1m 32s ago',
        description: 'Deploys from GitHub',
        environment: 'Preview',
    },
    {
        id: 2,
        href: '#',
        name: 'mobile-api',
        job: 'LSPD',
        status: 'online',
        statusText: 'Deployed 3m ago',
        description: 'Deploys from GitHub',
        environment: 'Production',
    },
    {
        id: 3,
        href: '#',
        name: 'tailwindcss.com',
        job: 'LSPD',
        status: 'offline',
        statusText: 'Deployed 3h ago',
        description: 'Deploys from GitHub',
        environment: 'Preview',
    },
];
</script>

<template>
    <div class="mx-2">
        <ul role="list" class="divide-y divide-white/5">
            <li v-for="project in projects" :key="project.id" class="relative flex items-center space-x-4 py-2">
                <div class="min-w-0 flex-auto">
                    <div class="flex items-center gap-x-3">
                        <div :class="[statuses[project.status], 'flex-none rounded-full p-1']">
                            <div class="h-2 w-2 rounded-full bg-current" />
                        </div>
                        <h2 class="min-w-0 text-sm font-semibold leading-6 text-white">
                            <a :href="project.href" class="flex gap-x-2">
                                <span class="truncate">{{ project.job }}</span>
                                <span class="text-gray-400">/</span>
                                <span class="whitespace-nowrap">{{ project.name }}</span>
                                <span class="absolute inset-0" />
                            </a>
                        </h2>
                        <Menu as="div" class="relative flex-none">
                            <MenuButton class="-m-2.5 block p-2.5 text-gray-500 hover:text-gray-100">
                                <span class="sr-only">Open options</span>
                                <SvgIcon type="mdi" :path="mdiDotsVertical" class="h-5 w-5" aria-hidden="true" />
                            </MenuButton>
                            <transition
                                enter-active-class="transition ease-out duration-100"
                                enter-from-class="transform opacity-0 scale-95"
                                enter-to-class="transform opacity-100 scale-100"
                                leave-active-class="transition ease-in duration-75"
                                leave-from-class="transform opacity-100 scale-100"
                                leave-to-class="transform opacity-0 scale-95"
                            >
                                <MenuItems
                                    class="absolute right-0 z-10 mt-2 w-32 origin-top-right rounded-md bg-white py-2 shadow-lg ring-1 ring-gray-900/5 focus:outline-none"
                                >
                                    <MenuItem v-slot="{ active }">
                                        <a
                                            href="#"
                                            :class="[
                                                active ? 'bg-gray-50' : '',
                                                'block px-3 py-1 text-sm leading-6 text-gray-900',
                                            ]"
                                            >Edit<span class="sr-only">, {{ project.name }}</span></a
                                        >
                                    </MenuItem>
                                    <MenuItem v-slot="{ active }">
                                        <a
                                            href="#"
                                            :class="[
                                                active ? 'bg-gray-50' : '',
                                                'block px-3 py-1 text-sm leading-6 text-gray-900',
                                            ]"
                                            >Move<span class="sr-only">, {{ project.name }}</span></a
                                        >
                                    </MenuItem>
                                    <MenuItem v-slot="{ active }">
                                        <a
                                            href="#"
                                            :class="[
                                                active ? 'bg-gray-50' : '',
                                                'block px-3 py-1 text-sm leading-6 text-gray-900',
                                            ]"
                                            >Delete<span class="sr-only">, {{ project.name }}</span></a
                                        >
                                    </MenuItem>
                                </MenuItems>
                            </transition>
                        </Menu>
                    </div>
                    <div class="mt-2 flex items-center gap-x-2.5 text-xs leading-5 text-gray-400">
                        <p class="truncate">{{ project.description }}</p>
                        <svg viewBox="0 0 2 2" class="h-0.5 w-0.5 flex-none fill-gray-300">
                            <circle cx="1" cy="1" r="1" />
                        </svg>
                        <p class="whitespace-nowrap">{{ project.statusText }}</p>
                    </div>
                </div>
                <SvgIcon type="mdi" :path="mdiChevronRight" class="h-5 w-5 flex-none text-gray-400" aria-hidden="true" />
            </li>
        </ul>
    </div>
</template>
