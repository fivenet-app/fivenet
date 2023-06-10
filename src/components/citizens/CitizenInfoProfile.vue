<script lang="ts" setup>
import { KeyIcon } from '@heroicons/vue/20/solid';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { useClipboard } from '@vueuse/core';
import { ref } from 'vue';
import CharSexBadge from '~/components/citizens/CharSexBadge.vue';
import CitizenInfoJobModal from '~/components/citizens/CitizenInfoJobModal.vue';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificationsStore } from '~/store/notifications';
import { User } from '~~/gen/ts/resources/users/users';
import CitizenInfoReasonModal from './CitizenInfoReasonModal.vue';

const { $grpc } = useNuxtApp();
const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();

const w = window;
const clipboard = useClipboard();

const props = defineProps<{
    user: User;
}>();

const wantedState = ref(props.user.props ? props.user.props?.wanted : false);
const reason = ref<string>('');
const jobModal = ref<boolean>(false);

async function toggleWantedStatus(): Promise<void> {
    return new Promise(async (res, rej) => {
        if (!props.user) {
            return res();
        }

        wantedState.value = !props.user.props?.wanted;

        let userProps = props.user?.props;
        if (!userProps) {
            userProps = {
                userId: props.user.userId,
            };

            props.user.props = userProps;
        }

        userProps.wanted = wantedState.value;

        try {
            await $grpc.getCitizenStoreClient().setUserProps({
                props: userProps,
                reason: reason.value,
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.action_successfull.title', parameters: [] },
                content: { key: 'notifications.action_successfull.content', parameters: [] },
                type: 'success',
            });

            reason.value = '';
            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const templatesOpen = ref(false);

const reasonOpen = ref(false);

function openTemplates(): void {
    clipboardStore.addUser(props.user);

    templatesOpen.value = true;
}
</script>

<template>
    <TemplatesModal :open="templatesOpen" @close="templatesOpen = false" :auto-fill="true" />
    <CitizenInfoReasonModal
        :open="reasonOpen"
        @close="reasonOpen = false"
        @submit="
            toggleWantedStatus();
            jobModal = false;
        "
        v-model:reason="reason"
    />
    <CitizenInfoJobModal :open="jobModal" @close="jobModal = false" :user="user" @submit="jobModal = false" />
    <div class="w-full mx-auto max-w-7xl grow lg:flex xl:px-2">
        <div class="flex-1 xl:flex">
            <div class="px-2 py-3 xl:flex-1">
                <div class="divide-y divide-base-200">
                    <div class="px-4 py-5 sm:px-0 sm:py-0">
                        <dl class="space-y-8 sm:space-y-0 sm:divide-y sm:divide-base-200">
                            <div class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    {{ $t('common.date_of_birth') }}
                                </dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                    {{ user?.dateofbirth }}
                                </dd>
                            </div>
                            <div class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    {{ $t('common.sex') }}
                                </dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                    {{ user?.sex!.toUpperCase() }}
                                    {{ ' ' }}
                                    <CharSexBadge :sex="user?.sex ? user?.sex : ''" />
                                </dd>
                            </div>
                            <div class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    {{ $t('common.height') }}
                                </dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">{{ user?.height }}cm</dd>
                            </div>
                            <div v-can="'CitizenStoreService.ListCitizens.Fields.PhoneNumber'" class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    {{ $t('common.phone') }}
                                    {{ $t('common.number') }}
                                </dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                    {{ user?.phoneNumber }}
                                </dd>
                            </div>
                            <div class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    {{ $t('common.visum') }}
                                </dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                    {{ user?.visum }}
                                </dd>
                            </div>
                            <div v-can="'CitizenStoreService.ListCitizens.Fields.Licenses'" class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    {{ $t('common.license', 2) }}
                                </dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                    <span v-if="user?.licenses.length === 0">
                                        {{ $t('common.no_licenses') }}
                                    </span>
                                    <ul v-else role="list" class="border divide-y rounded-md divide-base-200 border-base-200">
                                        <li
                                            v-for="license in user?.licenses"
                                            class="flex items-center justify-between py-3 pl-3 pr-4 text-sm"
                                        >
                                            <div class="flex items-center flex-1">
                                                <KeyIcon class="flex-shrink-0 w-5 h-5 text-base-400" aria-hidden="true" />
                                                <span class="flex-1 ml-2 truncate"
                                                    >{{ license.label }} ({{ license.type.toUpperCase() }})</span
                                                >
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
                <button
                    v-can="'CitizenStoreService.SetUserProps.Fields.Job'"
                    type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="jobModal = true"
                >
                    {{ $t('components.citizens.citizen_info_profile.set_job') }}
                </button>
            </div>
            <div class="flex-initial" v-can="'CitizenStoreService.SetUserProps.Fields.Wanted'">
                <button
                    type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-error-500 text-neutral hover:bg-error-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="reasonOpen = true"
                >
                    {{
                        wantedState
                            ? $t('components.citizens.citizen_info_profile.revoke_wanted')
                            : $t('components.citizens.citizen_info_profile.set_wanted')
                    }}
                </button>
            </div>
            <div class="flex-initial">
                <button
                    v-can="'DocStoreService.CreateDocument'"
                    type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-base-700 text-neutral hover:bg-base-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="openTemplates()"
                >
                    {{ $t('components.citizens.citizen_info_profile.create_new_document') }}
                </button>
            </div>
            <div class="flex-initial">
                <button
                    v-can="'CitizenStoreService.SetUserProps.Fields.Wanted'"
                    type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-base-700 text-neutral hover:bg-base-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="clipboard.copy(w.location.href)"
                >
                    {{ $t('components.citizens.citizen_info_profile.copy_profile_link') }}
                </button>
            </div>
        </div>
    </div>
</template>
