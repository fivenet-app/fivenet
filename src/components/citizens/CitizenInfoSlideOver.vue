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
import { EllipsisVerticalIcon, KeyIcon } from '@heroicons/vue/20/solid';
import { Character } from '@arpanet/gen/common/character_pb';
import CitizenActivityFeed from './CitizenActivityFeed.vue';
import { UsersServiceClient } from '@arpanet/gen/users/UsersServiceClientPb';
import { RpcError } from 'grpc-web';
import config from '../../config';
import { clientAuthOptions, handleGRPCError } from '../../grpc';
import { SetUserPropsRequest } from '@arpanet/gen/users/users_pb';

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
        KeyIcon,
        CitizenActivityFeed,
    },
    data() {
        return {
            client: new UsersServiceClient(config.apiProtoURL, null, clientAuthOptions),
        };
    },
    props: {
        'user': {
            required: true,
            type: Character,
        },
        'open': {
            required: true,
            type: Boolean,
        },
    },
    methods: {
        handleClose() {
            this.$emit('close');
        },
        getTimeInHoursAndMins(timeInsSeconds: number): string {
            if (timeInsSeconds && timeInsSeconds > 0) {
                const minsTemp = timeInsSeconds / 60;
                let hours = Math.floor(minsTemp / 60);
                const mins = minsTemp % 60;
                const hoursText = 'hrs';
                const minsText = 'mins';

                if (hours !== 0 && mins !== 0) {
                    if (mins >= 59) {
                        hours += 1;
                        return `${hours} ${hoursText} `;
                    } else {
                        return `${hours} ${hoursText} ${mins.toFixed(0)} ${minsText}`;
                    }
                } else if (hours === 0 && mins !== 0) {
                    return `${mins.toFixed(0)} ${minsText}`;
                } else if (hours !== 0 && mins === 0) {
                    return `${hours} ${hoursText}`;
                }
            }
            return '-';
        },
        toggleWantedStatus(event: any) {
            const wantedState = !this.user.getProps()?.getWanted();
            event.target.message = wantedState ? 'Revoke Wanted Status' : 'Set Person as Wanted'

            const req = new SetUserPropsRequest();
            req.setUserid(this.user.getUserid());
            req.setWanted(wantedState);
            this.client.setUserProps(req, null)
                .then((resp) => {
                    this.user.getProps()?.setWanted(wantedState);
                }).catch((err: RpcError) => {
                    handleGRPCError(err, this.$route);
                });
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
                            <DialogPanel class="pointer-events-auto w-screen max-w-5xl">
                                <div class="flex h-full flex-col overflow-y-scroll bg-white shadow-xl">
                                    <div class="px-4 py-6 sm:px-6">
                                        <div class="flex items-start justify-between">
                                            <DialogTitle class="text-base font-semibold leading-6 text-gray-900">
                                                Citizen Info
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
                                            <div
                                                class="lg:-mt-15 -mt-12 flow-root px-4 sm:-mt-8 sm:flex sm:items-end sm:px-6">
                                                <div class="mt-6 sm:ml-6 sm:flex-1">
                                                    <div>
                                                        <div class="flex items-center">
                                                            <h3 class="text-xl font-bold text-gray-900 sm:text-2xl">
                                                                {{ user.getFirstname() }}, {{ user.getLastname() }}
                                                                <span v-if="user.getProps()?.getWanted()"
                                                                    class="inline-flex items-center rounded-md bg-red-100 px-2.5 py-0.5 text-sm font-medium text-red-800">WANTED</span>
                                                            </h3>
                                                        </div>
                                                        <p class="text-sm text-gray-500">
                                                            <span
                                                                class="inline-flex items-center rounded-md bg-gray-100 px-2.5 py-0.5 text-sm font-medium text-gray-800">{{
                                                                    user.getJob() }} (Rank: {{ user.getJobgrade() }})
                                                            </span>
                                                        </p>
                                                    </div>
                                                    <div class="mt-5 flex flex-wrap space-y-3 sm:space-y-0 sm:space-x-3">
                                                        <button v-can="'users-setuserprops-wanted'" type="button"
                                                            class="inline-flex w-full flex-shrink-0 items-center justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:flex-1"
                                                            @click="toggleWantedStatus($event)">Set
                                                            Wanted Status</button>
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
                                                        Date of Birth</dt>
                                                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        {{ user.getDateofbirth() }}
                                                    </dd>
                                                </div>
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-gray-500 sm:w-40 sm:flex-shrink-0 lg:w-48">
                                                        Sex</dt>
                                                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0 sm:ml-6">{{
                                                        user.getSex() }}cm</dd>
                                                </div>
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-gray-500 sm:w-40 sm:flex-shrink-0 lg:w-48">
                                                        Height</dt>
                                                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0 sm:ml-6">{{
                                                        user.getHeight() }}cm</dd>
                                                </div>
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-gray-500 sm:w-40 sm:flex-shrink-0 lg:w-48">
                                                        Visum</dt>
                                                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        {{ user.getVisum() }}</dd>
                                                </div>
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-gray-500 sm:w-40 sm:flex-shrink-0 lg:w-48">
                                                        Playtime</dt>
                                                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        {{ getTimeInHoursAndMins(user.getPlaytime()) }}
                                                    </dd>
                                                </div>
                                                <div v-can="'users-findusers-licenses'" class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-gray-500 sm:w-40 sm:flex-shrink-0 lg:w-48">
                                                        Licenses</dt>
                                                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        <span v-if="user.getLicensesList().length == 0">No Licenses.</span>
                                                        <ul v-else role="list"
                                                            class="divide-y divide-gray-200 rounded-md border border-gray-200">
                                                            <li v-for="license in user.getLicensesList()"
                                                                class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                                                                <div class="flex flex-1 items-center">
                                                                    <KeyIcon class="h-5 w-5 flex-shrink-0 text-gray-400"
                                                                        aria-hidden="true" />
                                                                    <span class="ml-2 flex-1 truncate">{{
                                                                        license.getName().toUpperCase() }}</span>
                                                                </div>
                                                            </li>
                                                        </ul>
                                                    </dd>
                                                </div>
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <CitizenActivityFeed :userID="user.getUserid()" />
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
