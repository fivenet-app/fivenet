<script lang="ts" setup>
import { mdiCarEmergency, mdiHoopHouse } from '@mdi/js';
import SvgIcon from '@jamescoyle/vue-icon';

const navigation = [
    { name: 'Update Unit Status' },
    { name: 'Update Dispatch Status' },
    { name: 'TODO' },
    { name: 'TODO' },
    { name: 'TODO' },
    { name: 'TODO' },
    { name: 'Dispatch: Completed' },
];
const dispatches = [
    { id: 1, name: 'DSP-2019', color: 'H' },
    { id: 2, name: 'DSP-12345', color: 'T' },
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
                            <button
                                type="button"
                                class="text-accent-100 bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2"
                            >
                                <SvgIcon
                                    type="mdi"
                                    :path="mdiHoopHouse"
                                    class="text-base-200 group-hover:text-white h-6 w-6 shrink-0"
                                    aria-hidden="true"
                                />
                                <span class="mt-2">{{ item.name }}</span>
                            </button>
                        </li>
                    </ul>
                </li>
                <li>
                    <div class="text-xs font-semibold leading-6 text-base-200">Your Dispatches</div>
                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                        <li v-for="dispatch in dispatches" :key="dispatch.name">
                            <button
                                type="button"
                                class="text-accent-100 bg-error-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-3 text-xs my-2"
                            >
                                <SvgIcon type="mdi" :path="mdiCarEmergency" class="h-6 w-6" aria-hidden="true" />
                                <span class="mt-2 truncate">{{ dispatch.name }}</span>
                            </button>
                        </li>
                    </ul>
                </li>
            </ul>
        </nav>
    </div>
</template>
