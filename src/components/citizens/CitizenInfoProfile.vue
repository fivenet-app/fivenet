<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { User, UserProps } from '@arpanet/gen/resources/users/users_pb';
import { getCitizenStoreClient } from '../../grpc/grpc';
import { SetUserPropsRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { dispatchNotification } from '../notification';
import CharSexBadge from '../misc/CharSexBadge.vue';
import { KeyIcon } from '@heroicons/vue/20/solid';
import { useClipboard } from '@vueuse/core';

const w = window;
const clipboard = useClipboard();

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

function toggleWantedStatus(): void {
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
    <div class="mx-auto w-full max-w-7xl grow lg:flex xl:px-2">
        <div class="flex-1 xl:flex">
            <div class="py-3 px-2 xl:flex-1">
                <div class="divide-y divide-base-200">
                    <div class="px-4 py-5 sm:px-0 sm:py-0">
                        <dl class="space-y-8 sm:space-y-0 sm:divide-y sm:divide-base-200">
                            <div class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    Date of Birth</dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                    {{ user?.getDateofbirth() }}
                                </dd>
                            </div>
                            <div class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    Sex</dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                    {{ user?.getSex().toUpperCase() }}
                                    {{ ' ' }}
                                    <CharSexBadge :sex="user?.getSex() ? user?.getSex() : ''" />
                                </dd>
                            </div>
                            <div class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    Height</dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">{{
                                    user?.getHeight() }}cm</dd>
                            </div>
                            <div class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    Phone Number</dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">{{
                                    user?.getPhoneNumber() }}</dd>
                            </div>
                            <div class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    Visum</dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                    {{ user?.getVisum() }}</dd>
                            </div>
                            <div v-can="'CitizenStoreService.FindUsers.Licenses'" class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    Licenses</dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                    <span v-if="user?.getLicensesList().length == 0">No Licenses.</span>
                                    <ul v-else role="list"
                                        class="divide-y divide-base-200 rounded-md border border-base-200">
                                        <li v-for="license in user?.getLicensesList()"
                                            class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                                            <div class="flex flex-1 items-center">
                                                <KeyIcon class="h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                                                <span class="ml-2 flex-1 truncate">{{
                                                    license.getLabel() }} ({{ license.getType().toUpperCase() }})</span>
                                            </div>
                                        </li>
                                    </ul>
                                </dd>
                            </div>
                        </dl>
                    </div>
                </div>
            </div>
        </div>

        <div class="shrink-0 py-4 px-2 lg:w-96 pr-2 flex flex-col gap-2">
            <div class="flex-initial">
                <button v-can="'CitizenStoreService.SetUserProps.Wanted'" type="button"
                    class="inline-flex w-full flex-shrink-0 items-center justify-center rounded-md bg-error-500 px-3 py-2 text-sm font-semibold text-neutral shadow-sm hover:bg-error-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1 transition-colors"
                    @click="toggleWantedStatus()">{{ wantedState ?
                        'Revoke Wanted Status' : 'Set Person Wanted' }}
                </button>
            </div>
            <div class="flex-initial">
                <button v-can="'CitizenStoreService.SetUserProps.Wanted'" type="button"
                    class="inline-flex w-full flex-shrink-0 items-center justify-center rounded-md bg-base-700 px-3 py-2 text-sm font-semibold text-neutral shadow-sm hover:bg-base-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1 transition-colors"
                    @click="clipboard.copy(w.location.href)">Copy
                    profile link
                </button>
            </div>
        </div>
    </div>
</template>
