<script lang="ts" setup>
import { useClipboard } from '@vueuse/core';
import { LicenseIcon } from 'mdi-vue3';
import { ref } from 'vue';
import CharSexBadge from '~/components/citizens/CharSexBadge.vue';
import JobModal from '~/components/citizens/info/JobModal.vue';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { useClipboardStore } from '~/store/clipboard';
import { User } from '~~/gen/ts/resources/users/users';
import TrafficPointsModal from './TrafficPointsModal.vue';
import WantedModal from './WantedModal.vue';

const clipboardStore = useClipboardStore();

const w = window;
const clipboard = useClipboard();

const props = defineProps<{
    user: User;
}>();

const templatesOpen = ref(false);

function openTemplates(): void {
    clipboardStore.addUser(props.user);

    templatesOpen.value = true;
}

const setJobModal = ref(false);
const setWantedModal = ref(false);
const trafficPointsModal = ref(false);
</script>

<template>
    <TemplatesModal :open="templatesOpen" @close="templatesOpen = false" :auto-fill="true" />
    <WantedModal :open="setWantedModal" @close="setWantedModal = false" :user="user" />
    <JobModal :open="setJobModal" @close="setJobModal = false" :user="user" />
    <TrafficPointsModal :open="trafficPointsModal" @close="trafficPointsModal = false" :user="user" />
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
                            <div
                                v-if="can('CitizenStoreService.ListCitizens.Fields.PhoneNumber')"
                                class="sm:flex sm:px-6 sm:py-5"
                            >
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    {{ $t('common.phone') }}
                                    {{ $t('common.number') }}
                                </dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                    <span v-for="part in (user?.phoneNumber ?? '').match(/.{1,3}/g)" class="mr-1">{{
                                        part
                                    }}</span>
                                </dd>
                            </div>
                            <div class="sm:flex sm:px-6 sm:py-5">
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    {{ $t('common.visum') }}
                                </dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6 text-blue-400">
                                    {{ user?.visum }}
                                </dd>
                            </div>
                            <div
                                v-if="can('CitizenStoreService.ListCitizens.Fields.UserProps.TrafficInfractionPoints')"
                                class="sm:flex sm:px-6 sm:py-5"
                            >
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    {{ $t('common.traffic_infraction_points') }}
                                </dt>
                                <dd
                                    class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6"
                                    :class="(user?.props?.trafficInfractionPoints ?? 0n) >= 10 ? 'text-error-500' : ''"
                                >
                                    {{ $t('common.point', parseInt((user?.props?.trafficInfractionPoints ?? 0n).toString())) }}
                                </dd>
                            </div>
                            <div
                                v-if="can('CitizenStoreService.ListCitizens.Fields.UserProps.OpenFines')"
                                class="sm:flex sm:px-6 sm:py-5"
                            >
                                <dt class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48">
                                    {{ $t('common.fine') }}
                                </dt>
                                <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                    <span v-if="(user.props?.openFines ?? 0n) <= 0n">
                                        {{ $t('common.no_open_fine') }}
                                    </span>
                                    <span v-else class="text-error-500">
                                        {{ $n(parseInt((user?.props?.openFines ?? 0n).toString()), 'currency') }}
                                    </span>
                                </dd>
                            </div>
                            <div v-if="can('CitizenStoreService.ListCitizens.Fields.Licenses')" class="sm:flex sm:px-6 sm:py-5">
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
                                                <LicenseIcon class="flex-shrink-0 w-5 h-5 text-base-400" aria-hidden="true" />
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
            <div class="flex-initial" v-if="can('CitizenStoreService.SetUserProps.Fields.Wanted')">
                <button
                    type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-error-500 text-neutral hover:bg-error-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="setWantedModal = true"
                >
                    {{
                        user.props?.wanted
                            ? $t('components.citizens.citizen_info_profile.revoke_wanted')
                            : $t('components.citizens.citizen_info_profile.set_wanted')
                    }}
                </button>
            </div>
            <div class="flex-initial">
                <button
                    v-if="can('CitizenStoreService.SetUserProps.Fields.Job')"
                    type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="setJobModal = true"
                >
                    {{ $t('components.citizens.citizen_info_profile.set_job') }}
                </button>
            </div>
            <div class="flex-initial" v-if="can('CitizenStoreService.SetUserProps.Fields.TrafficInfractionPoints')">
                <button
                    type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-500 text-neutral hover:bg-secondary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="trafficPointsModal = true"
                >
                    {{ $t('components.citizens.citizen_info_profile.set_traffic_points') }}
                </button>
            </div>
            <div class="flex-initial" v-if="can('DocStoreService.CreateDocument')">
                <button
                    type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-base-700 text-neutral hover:bg-base-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="openTemplates()"
                >
                    {{ $t('components.citizens.citizen_info_profile.create_new_document') }}
                </button>
            </div>
            <div class="flex-initial">
                <button
                    type="button"
                    class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-base-700 text-neutral hover:bg-base-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:flex-1"
                    @click="clipboard.copy(w.location.href)"
                >
                    {{ $t('components.citizens.citizen_info_profile.copy_profile_link') }}
                </button>
            </div>
            <div class="flex-initial">
                <IDCopyBadge
                    :id="user.userId.toString()"
                    prefix="CIT"
                    :title="{ key: 'notifications.citizen_info.copy_citizen_id.title', parameters: [] }"
                    :content="{ key: 'notifications.citizen_info.copy_citizen_id.content', parameters: [] }"
                />
            </div>
        </div>
    </div>
</template>
