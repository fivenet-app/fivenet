<script lang="ts" setup>
import { LicenseIcon } from 'mdi-vue3';
import { ref } from 'vue';
import CharSexBadge from '~/components/citizens/CharSexBadge.vue';
import JobModal from '~/components/citizens/info/JobModal.vue';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import PhoneNumber from '~/components/partials/citizens/PhoneNumber.vue';
import { attr } from '~/composables/can';
import { useClipboardStore } from '~/store/clipboard';
import { User } from '~~/gen/ts/resources/users/users';
import TrafficPointsModal from '~/components/citizens/info/TrafficPointsModal.vue';
import WantedModal from '~/components/citizens/info/WantedModal.vue';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';

const clipboardStore = useClipboardStore();

const w = window;

const props = defineProps<{
    user: User;
}>();

defineEmits<{
    (e: 'update:wantedStatus', value: boolean): void;
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
    (e: 'update:trafficInfractionPoints', value: number): void;
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
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="flow-root">
                <div class="mx-0 -my-2 overflow-x-auto">
                    <div class="inline-block min-w-full align-middle px-1">
                        <TemplatesModal :open="templatesOpen" :auto-fill="true" @close="templatesOpen = false" />
                        <WantedModal
                            :open="setWantedModal"
                            :user="user"
                            @close="setWantedModal = false"
                            @update:wanted-status="$emit('update:wantedStatus', $event)"
                        />
                        <JobModal
                            :user="user"
                            :open="setJobModal"
                            @close="setJobModal = false"
                            @update:job="$emit('update:job', $event)"
                        />
                        <TrafficPointsModal
                            :open="trafficPointsModal"
                            :user="user"
                            @close="trafficPointsModal = false"
                            @update:traffic-infraction-points="$emit('update:trafficInfractionPoints', $event)"
                        />

                        <div class="w-full grow lg:flex xl:px-2">
                            <div class="flex-1 xl:flex">
                                <div class="px-2 py-3 xl:flex-1">
                                    <div class="divide-y divide-base-200">
                                        <div class="px-4 py-5 sm:px-0 sm:py-0">
                                            <dl class="space-y-8 sm:space-y-0 sm:divide-y sm:divide-base-200">
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.date_of_birth') }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        {{ user?.dateofbirth }}
                                                    </dd>
                                                </div>
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.sex') }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        {{ user?.sex!.toUpperCase() }}
                                                        {{ ' ' }}
                                                        <CharSexBadge :sex="user?.sex ? user?.sex : ''" />
                                                    </dd>
                                                </div>
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.height') }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        {{ user?.height }}cm
                                                    </dd>
                                                </div>
                                                <div
                                                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber')"
                                                    class="sm:flex sm:px-6 sm:py-5"
                                                >
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.phone_number') }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        <PhoneNumber :number="user.phoneNumber" />
                                                    </dd>
                                                </div>
                                                <div class="sm:flex sm:px-6 sm:py-5">
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.visum') }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6 text-blue-400"
                                                    >
                                                        {{ user?.visum }}
                                                    </dd>
                                                </div>
                                                <div
                                                    v-if="
                                                        attr(
                                                            'CitizenStoreService.ListCitizens',
                                                            'Fields',
                                                            'UserProps.TrafficInfractionPoints',
                                                        )
                                                    "
                                                    class="sm:flex sm:px-6 sm:py-5"
                                                >
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.traffic_infraction_points') }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6"
                                                        :class="
                                                            (user?.props?.trafficInfractionPoints ?? 0) >= 10
                                                                ? 'text-error-500'
                                                                : ''
                                                        "
                                                    >
                                                        {{ $t('common.point', user?.props?.trafficInfractionPoints ?? 0) }}
                                                    </dd>
                                                </div>
                                                <div
                                                    v-if="
                                                        attr(
                                                            'CitizenStoreService.ListCitizens',
                                                            'Fields',
                                                            'UserProps.OpenFines',
                                                        )
                                                    "
                                                    class="sm:flex sm:px-6 sm:py-5"
                                                >
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.fine') }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        <span v-if="(user.props?.openFines ?? 0n) <= 0n">
                                                            {{ $t('common.no_open_fine') }}
                                                        </span>
                                                        <span v-else class="text-error-500">
                                                            {{
                                                                $n(
                                                                    parseInt((user?.props?.openFines ?? 0n).toString(), 10),
                                                                    'currency',
                                                                )
                                                            }}
                                                        </span>
                                                    </dd>
                                                </div>
                                                <div
                                                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'Licenses')"
                                                    class="sm:flex sm:px-6 sm:py-5"
                                                >
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.license', 2) }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:mt-0 sm:ml-6">
                                                        <span v-if="user?.licenses.length === 0">
                                                            {{ $t('common.no_licenses') }}
                                                        </span>
                                                        <ul
                                                            v-else
                                                            role="list"
                                                            class="border divide-y rounded-md divide-base-200 border-base-200"
                                                        >
                                                            <li
                                                                v-for="license in user?.licenses"
                                                                :key="license.type"
                                                                class="flex items-center justify-between py-3 pl-3 pr-4 text-sm"
                                                            >
                                                                <div class="flex items-center flex-1">
                                                                    <LicenseIcon
                                                                        class="flex-shrink-0 w-5 h-5 text-base-400"
                                                                        aria-hidden="true"
                                                                    />
                                                                    <span class="flex-1 ml-2 truncate"
                                                                        >{{ license.label }} ({{
                                                                            license.type.toUpperCase()
                                                                        }})</span
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
                                <div v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'Wanted')" class="flex-initial">
                                    <button
                                        type="button"
                                        class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-error-500 text-neutral hover:bg-error-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500 sm:flex-1"
                                        @click="setWantedModal = true"
                                    >
                                        {{
                                            user.props?.wanted
                                                ? $t('components.citizens.citizen_info_profile.revoke_wanted')
                                                : $t('components.citizens.citizen_info_profile.set_wanted')
                                        }}
                                    </button>
                                </div>
                                <div v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'Job')" class="flex-initial">
                                    <button
                                        type="button"
                                        class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500 sm:flex-1"
                                        @click="setJobModal = true"
                                    >
                                        {{ $t('components.citizens.citizen_info_profile.set_job') }}
                                    </button>
                                </div>
                                <div
                                    v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'TrafficInfractionPoints')"
                                    class="flex-initial"
                                >
                                    <button
                                        type="button"
                                        class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-500 text-neutral hover:bg-secondary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500 sm:flex-1"
                                        @click="trafficPointsModal = true"
                                    >
                                        {{ $t('components.citizens.citizen_info_profile.set_traffic_points') }}
                                    </button>
                                </div>
                                <div v-if="can('DocStoreService.CreateDocument')" class="flex-initial">
                                    <button
                                        type="button"
                                        class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-base-700 text-neutral hover:bg-base-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500 sm:flex-1"
                                        @click="openTemplates()"
                                    >
                                        {{ $t('components.citizens.citizen_info_profile.create_new_document') }}
                                    </button>
                                </div>
                                <div class="flex-initial">
                                    <button
                                        type="button"
                                        class="inline-flex items-center justify-center flex-shrink-0 w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-base-700 text-neutral hover:bg-base-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500 sm:flex-1"
                                        @click="copyToClipboardWrapper(w.location.href)"
                                    >
                                        {{ $t('components.citizens.citizen_info_profile.copy_profile_link') }}
                                    </button>
                                </div>
                                <div class="flex-initial">
                                    <IDCopyBadge
                                        :id="user.userId"
                                        prefix="CIT"
                                        :title="{ key: 'notifications.citizen_info.copy_citizen_id.title', parameters: {} }"
                                        :content="{ key: 'notifications.citizen_info.copy_citizen_id.content', parameters: {} }"
                                    />
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
