<script lang="ts" setup>
import { ref } from 'vue'
import { User, UserProps } from '@arpanet/gen/resources/users/users_pb';
import { getCitizenStoreClient } from '../../grpc/grpc';
import { SetUserPropsRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { dispatchNotification } from '../notification';
import CharSexBadge from '../misc/CharSexBadge.vue';
import { KeyIcon } from '@heroicons/vue/20/solid';
import { useClipboard } from '@vueuse/core';
import TemplatesModal from '../documents/TemplatesModal.vue';

const w = window;
const clipboard = useClipboard();

const props = defineProps({
    user: {
        required: true,
        type: User,
    },
});

const wantedState = ref(props.user.getProps() ? props.user.getProps()?.getWanted() : false);

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

const templatesOpen = ref(false);
</script>

<template>
    <TemplatesModal :open="templatesOpen" @close="templatesOpen = false" />
    <div class="w-full mx-auto max-w-7xl grow lg:flex xl:px-2">
        <div class="flex-1 xl:flex">
            <div class="px-2 py-3 xl:flex-1">
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
                                        class="border divide-y rounded-md divide-base-200 border-base-200">
                                        <li v-for="license in user?.getLicensesList()"
                                            class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                                            <div class="flex items-center flex-1">
                                                <KeyIcon class="flex-shrink-0 w-5 h-5 text-base-400" aria-hidden="true" />
                                                <span class="flex-1 ml-2 truncate">{{
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

        <div class="flex flex-col gap-2 px-2 py-4 pr-2 shrink-0 lg:w-96">
            <div class="flex-initial">
                <button v-can="'DocStoreService.CreateDocument'" type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-base-700 text-neutral hover:bg-base-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="templatesOpen = true">
                    Create New Document
                </button>
            </div>
            <div class="flex-initial">
                <button v-can="'CitizenStoreService.SetUserProps.Wanted'" type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-base-700 text-neutral hover:bg-base-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="clipboard.copy(w.location.href)">
                    Copy profile link
                </button>
            </div>
            <div class="flex-initial">
                <button v-can="'CitizenStoreService.SetUserProps.Wanted'" type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-error-500 text-neutral hover:bg-error-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="toggleWantedStatus()">{{ wantedState ?
                        'Revoke Wanted Status' : 'Set Person Wanted' }}
                </button>
            </div>
        </div>
    </div>
</template>
