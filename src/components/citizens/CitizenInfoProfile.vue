<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { User, UserProps } from '@arpanet/gen/resources/users/users_pb';
import { getCitizenStoreClient } from '../../grpc/grpc';
import { SetUserPropsRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { getSecondsFormattedAsDuration } from '../../utils/time';
import { dispatchNotification } from '../notification';
import CharSexBadge from '../misc/CharSexBadge.vue';
import { EllipsisVerticalIcon, KeyIcon } from '@heroicons/vue/20/solid';
import {
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
} from '@headlessui/vue';

const wantedState = ref(false);

const props = defineProps({
    user: {
        required: true,
        type: User,
    },
});

onMounted(() => {
    const userProps = props.user.getProps();
    if (!userProps) return;

    wantedState.value = userProps.getWanted();
});

function toggleWantedStatus() {
    if (!props.user) return;

    wantedState.value = !props.user.getProps()?.getWanted();

    const req = new SetUserPropsRequest();
    let userProps = props.user?.getProps();
    if (!userProps) {
        userProps = new UserProps();
        userProps.setUserId(props.user.getUserId());

        props.user.setProps(userProps);
    }

    userProps?.setWanted(wantedState.value);
    req.setProps(userProps);

    getCitizenStoreClient().
        setUserProps(req, null)
        .then((resp) => {
            dispatchNotification({ title: 'Success!', content: 'Your action was successfully submitted', type: 'success' });
        });
}
</script>

<template>
    <div class="divide-y divide-gray-200">
        <div class="pb-6">
            <div class="lg:-mt-15 -mt-12 flow-root px-4 sm:-mt-8 sm:flex sm:items-end sm:px-6">
                <div class="mt-6 sm:ml-6 sm:flex-1">
                    <div class="mt-5 flex flex-wrap space-y-3 sm:space-y-0 sm:space-x-3">
                        <button v-can="'CitizenStoreService.SetUserProps.Wanted'" type="button"
                            class="inline-flex w-full flex-shrink-0 items-center justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:flex-1"
                            @click="toggleWantedStatus()">{{ wantedState ?
                                'Revoke Wanted Status' : 'Set Person Wanted' }}
                        </button>
                        <div class="ml-3 inline-flex sm:ml-0">
                            <Menu as="div" class="relative inline-block text-left">
                                <MenuButton
                                    class="inline-flex items-center rounded-md bg-white p-2 text-gray-400 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
                                    <span class="sr-only">Open options menu</span>
                                    <EllipsisVerticalIcon class="h-5 w-5" aria-hidden="true" />
                                </MenuButton>
                                <transition enter-active-class="transition ease-out duration-100"
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
                                                :class="[active ? 'bg-gray-100 text-gray-300' : 'text-gray-700', 'block px-4 py-2 text-sm']">View
                                                profile</a>
                                            </MenuItem>
                                            <MenuItem v-slot="{ active }">
                                            <a href="#"
                                                :class="[active ? 'bg-gray-100 text-gray-300' : 'text-gray-700', 'block px-4 py-2 text-sm']">Copy
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
                    <dt class="text-sm font-medium text-white sm:w-40 sm:flex-shrink-0 lg:w-48">
                        Date of Birth</dt>
                    <dd class="mt-1 text-sm text-gray-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                        {{ user?.getDateofbirth() }}
                    </dd>
                </div>
                <div class="sm:flex sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium text-white sm:w-40 sm:flex-shrink-0 lg:w-48">
                        Sex</dt>
                    <dd class="mt-1 text-sm text-gray-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                        {{ user?.getSex().toUpperCase() }}
                        {{ ' ' }}
                        <CharSexBadge :sex="user?.getSex() ? user?.getSex() : ''" />
                    </dd>
                </div>
                <div class="sm:flex sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium text-white sm:w-40 sm:flex-shrink-0 lg:w-48">
                        Height</dt>
                    <dd class="mt-1 text-sm text-gray-300 sm:col-span-2 sm:mt-0 sm:ml-6">{{
                        user?.getHeight() }}cm</dd>
                </div>
                <div class="sm:flex sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium text-white sm:w-40 sm:flex-shrink-0 lg:w-48">
                        Phone Number</dt>
                    <dd class="mt-1 text-sm text-gray-300 sm:col-span-2 sm:mt-0 sm:ml-6">{{
                        user?.getPhoneNumber() }}</dd>
                </div>
                <div class="sm:flex sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium text-white sm:w-40 sm:flex-shrink-0 lg:w-48">
                        Visum</dt>
                    <dd class="mt-1 text-sm text-gray-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                        {{ user?.getVisum() }}</dd>
                </div>
                <div class="sm:flex sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium text-white sm:w-40 sm:flex-shrink-0 lg:w-48">
                        Playtime</dt>
                    <dd class="mt-1 text-sm text-gray-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                        {{ getSecondsFormattedAsDuration(user?.getPlaytime()) }}
                    </dd>
                </div>
                <div v-can="'CitizenStoreService.FindUsers.Licenses'" class="sm:flex sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium text-white sm:w-40 sm:flex-shrink-0 lg:w-48">
                        Licenses</dt>
                    <dd class="mt-1 text-sm text-gray-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                        <span v-if="user?.getLicensesList().length == 0">No Licenses.</span>
                        <ul v-else role="list" class="divide-y divide-gray-200 rounded-md border border-gray-200">
                            <li v-for="license in user?.getLicensesList()"
                                class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                                <div class="flex flex-1 items-center">
                                    <KeyIcon class="h-5 w-5 flex-shrink-0 text-gray-400" aria-hidden="true" />
                                    <span class="ml-2 flex-1 truncate">{{
                                        license.getType().toUpperCase() }}</span>
                                </div>
                            </li>
                        </ul>
                    </dd>
                </div>
            </dl>
        </div>
    </div>
</template>
