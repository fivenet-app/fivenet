<script lang="ts">
import { defineComponent } from 'vue'
import {
    Dialog,
    DialogPanel,
    DialogTitle,
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue';
import { XMarkIcon } from '@heroicons/vue/24/outline';
import { EllipsisVerticalIcon } from '@heroicons/vue/20/solid';

export default defineComponent({
    components: {
        Dialog,
        DialogPanel,
        DialogTitle,
        Menu,
        MenuButton,
        MenuItem,
        MenuItems,
        TransitionChild,
        TransitionRoot,
        XMarkIcon,
        EllipsisVerticalIcon,
    },
    props: {
        'open': {
            required: true,
            type: Boolean,
        },
    },
    methods: {
        handleClose() {
            this.$emit('close');
        },
    },
});
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="handleClose()">
            <div class="fixed inset-0" />

            <div class="fixed inset-0 overflow-hidden">
                <div class="absolute inset-0 overflow-hidden">
                    <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-full pl-10 sm:pl-16">
                        <TransitionChild as="template" enter="transform transition ease-in-out duration-500 sm:duration-700"
                            enter-from="translate-x-full" enter-to="translate-x-0"
                            leave="transform transition ease-in-out duration-500 sm:duration-700" leave-from="translate-x-0"
                            leave-to="translate-x-full">
                            <DialogPanel class="pointer-events-auto w-screen max-w-2xl">
                                <div class="flex h-full flex-col overflow-y-scroll bg-white shadow-xl">
                                    <div class="px-4 py-6 sm:px-6">
                                        <div class="flex items-start justify-between">
                                            <DialogTitle class="text-base font-semibold leading-6 text-gray-900">Profile
                                            </DialogTitle>
                                            <div class="ml-3 flex h-7 items-center">
                                                <button type="button"
                                                    class="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:ring-2 focus:ring-indigo-500"
                                                    @click="handleClose()">
                                                    <span class="sr-only">Close panel</span>
                                                    <XMarkIcon class="h-6 w-6" aria-hidden="true" />
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                    <!-- Main -->
                                    <div class="divide-y divide-gray-200">
                                        <div class="pb-6">
                                            <div class="h-24 bg-indigo-700 sm:h-20 lg:h-28" />
                                            <div
                                                class="lg:-mt-15 -mt-12 flow-root px-4 sm:-mt-8 sm:flex sm:items-end sm:px-6">
                                                <div>
                                                    <div class="-m-1 flex">
                                                        <div
                                                            class="inline-flex overflow-hidden rounded-lg border-4 border-white">
                                                            <img class="h-24 w-24 flex-shrink-0 sm:h-40 sm:w-40 lg:h-48 lg:w-48"
                                                                src="https://images.unsplash.com/photo-1501031170107-cfd33f0cbdcc?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=256&h=256&q=80"
                                                                alt="" />
                                                        </div>
                                                    </div>
                                                </div>
                                                <div class="mt-6 sm:ml-6 sm:flex-1">
                                                    <div>
                                                        <div class="flex items-center">
                                                            <h3 class="text-xl font-bold text-gray-900 sm:text-2xl">Ashley
                                                                Porter</h3>
                                                            <span
                                                                class="ml-2.5 inline-block h-2 w-2 flex-shrink-0 rounded-full bg-green-400">
                                                                <span class="sr-only">Online</span>
                                                            </span>
                                                        </div>
                                                        <p class="text-sm text-gray-500">@ashleyporter</p>
                                                    </div>
                                                    <div class="mt-5 flex flex-wrap space-y-3 sm:space-y-0 sm:space-x-3">
                                                        <button type="button"
                                                            class="inline-flex w-full flex-shrink-0 items-center justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:flex-1">Message</button>
                                                        <button type="button"
                                                            class="inline-flex w-full flex-1 items-center justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50">Call</button>
                                                        <div class="ml-3 inline-flex sm:ml-0">
                                                            <Menu as="div" class="relative inline-block text-left">
                                                                <MenuButton
                                                                    class="inline-flex items-center rounded-md bg-white p-2 text-gray-400 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
                                                                    <span class="sr-only">Open options menu</span>
                                                                    <EllipsisVerticalIcon class="h-5 w-5"
                                                                        aria-hidden="true" />
                                                                </MenuButton>
                                                                <transition
                                                                    enter-active-class="transition ease-out duration-100"
                                                                    enter-from-class="transform opacity-0 scale-95"
                                                                    enter-to-class="transform opacity-100 scale-100"
                                                                    leave-active-class="transition ease-in duration-75"
                                                                    leave-from-class="transform opacity-100 scale-100"
                                                                    leave-to-class="transform opacity-0 scale-95">
                                                                    <MenuItems
                                                                        class="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                                                                        <div class="py-1">
                                                                            <MenuItem v-slot="{ active }">
                                                                            <a href="#"
                                                                                :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'block px-4 py-2 text-sm']">View
                                                                                profile</a>
                                                                            </MenuItem>
                                                                            <MenuItem v-slot="{ active }">
                                                                            <a href="#"
                                                                                :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'block px-4 py-2 text-sm']">Copy
                                                                                profile link</a>
                                                                            </MenuItem>
                                                                        </div>
                                                                    </MenuItems>
                                                                </transition>
                                                            </Menu>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="px-4 py-5 sm:px-0 sm:py-0">
                                            <dl class="space-y-8 sm:space-y-0 sm:divide-y sm:divide-gray-200">
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-gray-500 sm:w-40 sm:flex-shrink-0 lg:w-48">
                                                        Bio</dt>
                                                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        <p>Enim feugiat ut ipsum, neque ut. Tristique mi id elementum
                                                            praesent. Gravida in tempus feugiat netus enim aliquet a, quam
                                                            scelerisque. Dictumst in convallis nec in bibendum aenean arcu.
                                                        </p>
                                                    </dd>
                                                </div>
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-gray-500 sm:w-40 sm:flex-shrink-0 lg:w-48">
                                                        Location</dt>
                                                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0 sm:ml-6">New
                                                        York, NY, USA</dd>
                                                </div>
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-gray-500 sm:w-40 sm:flex-shrink-0 lg:w-48">
                                                        Website</dt>
                                                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        ashleyporter.com</dd>
                                                </div>
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-gray-500 sm:w-40 sm:flex-shrink-0 lg:w-48">
                                                        Birthday</dt>
                                                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        <time datetime="1982-06-23">June 23, 1982</time>
                                                    </dd>
                                                </div>
                                            </dl>
                                        </div>
                                    </div>
                                </div>
                            </DialogPanel>
                        </TransitionChild>
                    </div>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
