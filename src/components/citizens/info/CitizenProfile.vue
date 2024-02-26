<script lang="ts" setup>
import {
    AccountAlertIcon,
    AccountCancelIcon,
    BriefcaseIcon,
    CameraIcon,
    CounterIcon,
    FileDocumentPlusIcon,
    LicenseIcon,
    LinkIcon,
} from 'mdi-vue3';
import { ref } from 'vue';
import CharSexBadge from '~/components/partials/citizens/CharSexBadge.vue';
import CitizenSetJobModal from '~/components/citizens/info/CitizenSetJobModal.vue';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { attr } from '~/composables/can';
import { useClipboardStore } from '~/store/clipboard';
import { User } from '~~/gen/ts/resources/users/users';
import CitizenSetTrafficPointsModal from '~/components/citizens/info/CitizenSetTrafficPointsModal.vue';
import CitizenSetWantedModal from '~/components/citizens/info/CitizenSetWantedModal.vue';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import { useNotificatorStore } from '~/store/notificator';
import CitizenSetMugShotModal from '~/components/citizens/info/CitizenSetMugShotModal.vue';
import type { File } from '~~/gen/ts/resources/filestore/file';

const props = defineProps<{
    user: User;
}>();

defineEmits<{
    (e: 'update:wantedStatus', value: boolean): void;
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
    (e: 'update:trafficInfractionPoints', value: number): void;
    (e: 'update:mugShot', value?: File): void;
}>();

const w = window;

const clipboardStore = useClipboardStore();

const notifications = useNotificatorStore();

const templatesOpen = ref(false);

function openTemplates(): void {
    clipboardStore.addUser(props.user);

    templatesOpen.value = true;
}

function copyLinkToClipboard(): void {
    copyToClipboardWrapper(w.location.href);

    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.link_copied.title', parameters: {} },
        content: { key: 'notifications.clipboard.link_copied.content', parameters: {} },
        duration: 3250,
        type: 'info',
    });
}

const setJobModal = ref(false);
const setWantedModal = ref(false);
const trafficPointsModal = ref(false);
const mugShotModal = ref(false);
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 align-middle">
                        <TemplatesModal :open="templatesOpen" :auto-fill="true" @close="templatesOpen = false" />
                        <CitizenSetWantedModal
                            :open="setWantedModal"
                            :user="user"
                            @close="setWantedModal = false"
                            @update:wanted-status="$emit('update:wantedStatus', $event)"
                        />
                        <CitizenSetJobModal
                            :user="user"
                            :open="setJobModal"
                            @close="setJobModal = false"
                            @update:job="$emit('update:job', $event)"
                        />
                        <CitizenSetTrafficPointsModal
                            :open="trafficPointsModal"
                            :user="user"
                            @close="trafficPointsModal = false"
                            @update:traffic-infraction-points="$emit('update:trafficInfractionPoints', $event)"
                        />
                        <CitizenSetMugShotModal
                            :open="mugShotModal"
                            :user="user"
                            @close="mugShotModal = false"
                            @update:mug-shot="$emit('update:mugShot', $event)"
                        />

                        <div class="w-full grow lg:flex xl:px-2">
                            <div class="flex-1 xl:flex">
                                <div class="px-2 py-3 xl:flex-1">
                                    <div class="divide-y divide-base-200">
                                        <div class="px-4 py-5 sm:px-0 sm:py-0">
                                            <dl class="space-y-4 sm:space-y-0 sm:divide-y sm:divide-base-200">
                                                <div class="sm:flex sm:px-5 sm:py-4">
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.date_of_birth') }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:ml-6 sm:mt-0">
                                                        {{ user?.dateofbirth }}
                                                    </dd>
                                                </div>
                                                <div class="sm:flex sm:px-5 sm:py-4">
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.sex') }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 inline-flex items-center gap-2 text-sm text-base-300 sm:col-span-2 sm:ml-6 sm:mt-0"
                                                    >
                                                        <span>{{ user?.sex!.toUpperCase() }} </span>
                                                        <CharSexBadge :sex="user?.sex ? user?.sex : ''" />
                                                    </dd>
                                                </div>
                                                <div class="sm:flex sm:px-5 sm:py-4">
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.height') }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:ml-6 sm:mt-0">
                                                        {{ user?.height }}cm
                                                    </dd>
                                                </div>
                                                <div
                                                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber')"
                                                    class="sm:flex sm:px-5 sm:py-4"
                                                >
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.phone_number') }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:ml-6 sm:mt-0">
                                                        <PhoneNumberBlock :number="user.phoneNumber" />
                                                    </dd>
                                                </div>
                                                <div class="sm:flex sm:px-5 sm:py-4">
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.visum') }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-blue-400 sm:col-span-2 sm:ml-6 sm:mt-0">
                                                        {{ user?.visum }}
                                                    </dd>
                                                </div>
                                                <div
                                                    v-if="
                                                        attr(
                                                            'CitizenStoreService.ListCitizens',
                                                            'Fields',
                                                            'UserProps.BloodType',
                                                        )
                                                    "
                                                    class="sm:flex sm:px-5 sm:py-4"
                                                >
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.blood_type') }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:ml-6 sm:mt-0">
                                                        {{ user?.props?.bloodType ?? $t('common.na') }}
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
                                                    class="sm:flex sm:px-5 sm:py-4"
                                                >
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.traffic_infraction_points', 2) }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 text-sm text-base-300 sm:col-span-2 sm:ml-6 sm:mt-0"
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
                                                    class="sm:flex sm:px-5 sm:py-4"
                                                >
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.fine') }}
                                                    </dt>
                                                    <dd class="mt-1 text-sm text-base-300 sm:col-span-2 sm:ml-6 sm:mt-0">
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
                                                    class="sm:flex sm:px-5 sm:py-4"
                                                >
                                                    <dt
                                                        class="text-sm font-medium text-neutral sm:w-40 sm:flex-shrink-0 lg:w-48"
                                                    >
                                                        {{ $t('common.license', 2) }}
                                                    </dt>
                                                    <dd class="w-full mt-1 text-sm text-base-300 sm:col-span-2 sm:ml-6 sm:mt-0">
                                                        <span v-if="user?.licenses.length === 0">
                                                            {{ $t('common.no_licenses') }}
                                                        </span>
                                                        <ul
                                                            v-else
                                                            role="list"
                                                            class="w-full divide-y divide-base-200 rounded-md border border-base-200"
                                                        >
                                                            <li
                                                                v-for="license in user?.licenses"
                                                                :key="license.type"
                                                                class="flex items-center justify-between py-3 pl-3 pr-4 text-sm"
                                                            >
                                                                <div class="flex flex-1 items-center">
                                                                    <LicenseIcon
                                                                        class="h-5 w-5 flex-shrink-0"
                                                                        aria-hidden="true"
                                                                    />
                                                                    <span
                                                                        class="ml-2 flex-1 truncate"
                                                                        :title="`${license.type.toUpperCase()}`"
                                                                        >{{ license.label }}
                                                                    </span>
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

                            <div class="flex shrink-0 flex-col gap-2 px-2 py-4 pr-2 lg:w-96">
                                <div v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'Wanted')" class="flex-initial">
                                    <button
                                        type="button"
                                        class="inline-flex w-full flex-shrink-0 items-center justify-center rounded-md bg-error-500 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-error-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500 sm:flex-1"
                                        @click="setWantedModal = true"
                                    >
                                        <AccountAlertIcon v-if="user.props?.wanted" class="w-5 h-auto mr-1.5" />
                                        <AccountCancelIcon v-else class="w-5 h-auto mr-1.5" />
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
                                        class="inline-flex w-full flex-shrink-0 items-center justify-center rounded-md bg-secondary-500 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-secondary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-secondary-500 sm:flex-1"
                                        @click="setJobModal = true"
                                    >
                                        <BriefcaseIcon class="w-5 h-auto mr-1.5" />
                                        {{ $t('components.citizens.citizen_info_profile.set_job') }}
                                    </button>
                                </div>
                                <div
                                    v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'TrafficInfractionPoints')"
                                    class="flex-initial"
                                >
                                    <button
                                        type="button"
                                        class="inline-flex w-full flex-shrink-0 items-center justify-center rounded-md bg-secondary-500 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-secondary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-secondary-500 sm:flex-1"
                                        @click="trafficPointsModal = true"
                                    >
                                        <CounterIcon class="w-5 h-auto mr-1.5" />
                                        {{ $t('components.citizens.citizen_info_profile.set_traffic_points') }}
                                    </button>
                                </div>
                                <div v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'MugShot')" class="flex-initial">
                                    <button
                                        type="button"
                                        class="inline-flex w-full flex-shrink-0 items-center justify-center rounded-md bg-secondary-500 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-secondary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-secondary-500 sm:flex-1"
                                        @click="mugShotModal = true"
                                    >
                                        <CameraIcon class="w-5 h-auto mr-1.5" />
                                        {{ $t('components.citizens.citizen_info_profile.set_mug_shot') }}
                                    </button>
                                </div>
                                <div v-if="can('DocStoreService.CreateDocument')" class="flex-initial">
                                    <button
                                        type="button"
                                        class="inline-flex w-full flex-shrink-0 items-center justify-center rounded-md bg-base-700 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-base-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500 sm:flex-1"
                                        @click="openTemplates()"
                                    >
                                        <FileDocumentPlusIcon class="w-5 h-auto mr-1.5" />
                                        {{ $t('components.citizens.citizen_info_profile.create_new_document') }}
                                    </button>
                                </div>
                                <div class="flex-initial">
                                    <button
                                        type="button"
                                        class="inline-flex w-full flex-shrink-0 items-center justify-center rounded-md bg-base-700 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-base-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500 sm:flex-1"
                                        @click="copyLinkToClipboard()"
                                    >
                                        <LinkIcon class="w-5 h-auto mr-1.5" />
                                        {{ $t('components.citizens.citizen_info_profile.copy_profile_link') }}
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
