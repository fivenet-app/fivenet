<script lang="ts" setup>
import { mdiHoopHouse } from '@mdi/js';
import SvgIcon from '@jamescoyle/vue-icon';

const navigation = [
    { name: 'Dashboard', href: '#', current: true },
    { name: 'Team', href: '#', current: false },
    { name: 'Projects', href: '#', current: false },
    { name: 'Calendar', href: '#', current: false },
    { name: 'Documents', href: '#', current: false },
    { name: 'Reports', href: '#', current: false },
];
const dispatches = [
    { id: 1, name: 'Heroicons', href: '#', initial: 'H', current: false },
    { id: 2, name: 'Tailwind Labs', href: '#', initial: 'T', current: false },
    { id: 3, name: 'Workcation', href: '#', initial: 'W', current: false },
];
</script>

<template>
    <!-- Sidebar component, swap this element with another sidebar if you like -->
    <div class="h-full flex grow gap-y-5 overflow-y-auto bg-base-600 px-6 py-4">
        <nav class="flex flex-1 flex-col">
            <ul role="list" class="flex flex-1 flex-col gap-y-7">
                <li>
                    <ul role="list" class="-mx-2 space-y-1">
                        <li v-for="item in navigation" :key="item.name">
                            <a
                                :href="item.href"
                                :class="[
                                    item.current
                                        ? 'bg-base-700 text-white'
                                        : 'text-base-200 hover:text-white hover:bg-base-700',
                                    'group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold',
                                ]"
                            >
                                <SvgIcon
                                    type="mdi"
                                    :path="mdiHoopHouse"
                                    :class="[
                                        item.current ? 'text-white' : 'text-base-200 group-hover:text-white',
                                        'h-6 w-6 shrink-0',
                                    ]"
                                    aria-hidden="true"
                                />
                                {{ item.name }}
                            </a>
                        </li>
                    </ul>
                </li>
                <li>
                    <div class="text-xs font-semibold leading-6 text-base-200">Your Dispatches</div>
                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                        <li v-for="dispatch in dispatches" :key="dispatch.name">
                            <a
                                :href="dispatch.href"
                                :class="[
                                    dispatch.current
                                        ? 'bg-base-700 text-white'
                                        : 'text-base-200 hover:text-white hover:bg-base-700',
                                    'group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold',
                                ]"
                            >
                                <span
                                    class="flex h-6 w-6 shrink-0 items-center justify-center rounded-lg border border-base-400 bg-base-500 text-[0.625rem] font-medium text-white"
                                    >{{ dispatch.initial }}</span
                                >
                                <span class="truncate">{{ dispatch.name }}</span>
                            </a>
                        </li>
                    </ul>
                </li>
            </ul>
        </nav>
    </div>
</template>
