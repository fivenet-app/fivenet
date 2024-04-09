<script lang="ts" setup>
import CharSexBadge from '~/components/partials/citizens/CharSexBadge.vue';
import CitizenSetJobModal from '~/components/citizens/info/props/CitizenSetJobModal.vue';
import TemplatesModal from '~/components/documents/templates/TemplatesModal.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { attr } from '~/composables/can';
import { useClipboardStore } from '~/store/clipboard';
import { User } from '~~/gen/ts/resources/users/users';
import CitizenSetTrafficPointsModal from '~/components/citizens/info/props/CitizenSetTrafficPointsModal.vue';
import CitizenSetWantedModal from '~/components/citizens/info/props/CitizenSetWantedModal.vue';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import { useNotificatorStore } from '~/store/notificator';
import CitizenSetMugShotModal from '~/components/citizens/info/props/CitizenSetMugShotModal.vue';
import type { File } from '~~/gen/ts/resources/filestore/file';

const props = defineProps<{
    user: User;
}>();

const emits = defineEmits<{
    (e: 'update:wantedStatus', value: boolean): void;
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
    (e: 'update:trafficInfractionPoints', value: number): void;
    (e: 'update:mugShot', value?: File): void;
}>();

const w = window;

const clipboardStore = useClipboardStore();

const notifications = useNotificatorStore();

const modal = useModal();

function openTemplates(): void {
    clipboardStore.addUser(props.user);

    modal.open(TemplatesModal, {});
}

function copyLinkToClipboard(): void {
    copyToClipboardWrapper(w.location.href);

    notifications.add({
        title: { key: 'notifications.clipboard.link_copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.link_copied.content', parameters: {} },
        timeout: 3250,
        type: 'info',
    });
}

defineShortcuts({
    'c-w': () => {
        if (!attr('CitizenStoreService.SetUserProps', 'Fields', 'Wanted')) {
            return;
        }

        modal.open(CitizenSetWantedModal, {
            user: props.user,
            'onUpdate:wantedStatus': ($event) => emits('update:wantedStatus', $event),
        });
    },
    'c-j': () => {
        if (!attr('CitizenStoreService.SetUserProps', 'Fields', 'Job')) {
            return;
        }

        modal.open(CitizenSetJobModal, {
            user: props.user,
            'onUpdate:job': ($event) => emits('update:job', $event),
        });
    },
    'c-p': () => {
        if (!attr('CitizenStoreService.SetUserProps', 'Fields', 'TrafficInfractionPoints')) {
            return;
        }

        modal.open(CitizenSetTrafficPointsModal, {
            user: props.user,
            'onUpdate:trafficInfractionPoints': ($event) => emits('update:trafficInfractionPoints', $event),
        });
    },
    'c-m': () => {
        if (!attr('CitizenStoreService.SetUserProps', 'Fields', 'MugShot')) {
            return;
        }

        modal.open(CitizenSetMugShotModal, {
            user: props.user,
            'onUpdate:mugShot': ($event) => emits('update:mugShot', $event),
        });
    },
    'c-d': () => {
        if (!can('DocStoreService.CreateDocument')) {
            return;
        }

        openTemplates();
    },
});
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2">
            <div class="flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 align-middle">
                        <div class="w-full grow lg:flex xl:px-2">
                            <div class="flex-1 xl:flex">
                                <div class="xl:flex-1">
                                    <div class="divide-y divide-base-200">
                                        <div class="px-4 py-5 sm:p-0">
                                            <dl class="space-y-4 sm:space-y-0 sm:divide-y sm:divide-base-200">
                                                <div class="sm:flex sm:px-5 sm:py-4">
                                                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                                                        {{ $t('common.date_of_birth') }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300"
                                                    >
                                                        {{ user.dateofbirth }}
                                                    </dd>
                                                </div>
                                                <div class="sm:flex sm:px-5 sm:py-4">
                                                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                                                        {{ $t('common.sex') }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 inline-flex items-center gap-2 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300"
                                                    >
                                                        <span>{{ user?.sex!.toUpperCase() }} </span>
                                                        <CharSexBadge :sex="user?.sex ? user?.sex : ''" />
                                                    </dd>
                                                </div>
                                                <div class="sm:flex sm:px-5 sm:py-4">
                                                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                                                        {{ $t('common.height') }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300"
                                                    >
                                                        {{ user?.height }}cm
                                                    </dd>
                                                </div>
                                                <div
                                                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber')"
                                                    class="sm:flex sm:px-5 sm:py-4"
                                                >
                                                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                                                        {{ $t('common.phone_number') }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300"
                                                    >
                                                        <PhoneNumberBlock :number="user.phoneNumber" />
                                                    </dd>
                                                </div>
                                                <div v-if="user.visum" class="sm:flex sm:px-5 sm:py-4">
                                                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                                                        {{ $t('common.visum') }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 text-sm text-blue-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-blue-300"
                                                    >
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
                                                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                                                        {{ $t('common.blood_type') }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300"
                                                    >
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
                                                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                                                        {{ $t('common.traffic_infraction_points', 2) }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300"
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
                                                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                                                        {{ $t('common.fine') }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300"
                                                    >
                                                        <span v-if="(user.props?.openFines ?? 0n) <= 0n">
                                                            {{ $t('common.no_open_fine') }}
                                                        </span>
                                                        <span v-else class="text-error-500">
                                                            {{
                                                                $n(
                                                                    parseInt((user?.props?.openFines ?? 0n).toString()),
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
                                                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                                                        {{ $t('common.license', 2) }}
                                                    </dt>
                                                    <dd
                                                        class="mt-1 w-full text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300"
                                                    >
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
                                                                    <UIcon name="i-mdi-license" class="size-5 shrink-0" />
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

                            <div class="flex shrink-0 flex-col gap-2 px-2 py-4 pr-0 lg:w-96">
                                <UTooltip
                                    v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'Wanted')"
                                    :text="
                                        user.props?.wanted
                                            ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                                            : $t('components.citizens.CitizenInfoProfile.set_wanted')
                                    "
                                    :shortcuts="['C', 'W']"
                                >
                                    <UButton
                                        color="red"
                                        block
                                        :icon="user.props?.wanted ? 'i-mdi-account-alert' : 'i-mdi-account-cancel'"
                                        @click="
                                            modal.open(CitizenSetWantedModal, {
                                                user: user,
                                                'onUpdate:wantedStatus': ($event) => $emit('update:wantedStatus', $event),
                                            })
                                        "
                                    >
                                        {{
                                            user.props?.wanted
                                                ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                                                : $t('components.citizens.CitizenInfoProfile.set_wanted')
                                        }}
                                    </UButton>
                                </UTooltip>

                                <UTooltip
                                    v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'Job')"
                                    :text="$t('components.citizens.CitizenInfoProfile.set_job')"
                                    :shortcuts="['C', 'J']"
                                >
                                    <UButton
                                        block
                                        icon="i-mdi-briefcase"
                                        @click="
                                            modal.open(CitizenSetJobModal, {
                                                user: user,
                                                'onUpdate:job': ($event) => $emit('update:job', $event),
                                            })
                                        "
                                    >
                                        {{ $t('components.citizens.CitizenInfoProfile.set_job') }}
                                    </UButton>
                                </UTooltip>

                                <UTooltip
                                    v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'TrafficInfractionPoints')"
                                    :text="$t('components.citizens.CitizenInfoProfile.set_traffic_points')"
                                    :shortcuts="['C', 'P']"
                                >
                                    <UButton
                                        block
                                        icon="i-mdi-counter"
                                        @click="
                                            modal.open(CitizenSetTrafficPointsModal, {
                                                user: user,
                                                'onUpdate:trafficInfractionPoints': ($event) =>
                                                    $emit('update:trafficInfractionPoints', $event),
                                            })
                                        "
                                    >
                                        {{ $t('components.citizens.CitizenInfoProfile.set_traffic_points') }}
                                    </UButton>
                                </UTooltip>

                                <UTooltip
                                    v-if="attr('CitizenStoreService.SetUserProps', 'Fields', 'MugShot')"
                                    :text="
                                        user.props?.wanted
                                            ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                                            : $t('components.citizens.CitizenInfoProfile.set_wanted')
                                    "
                                    :shortcuts="['C', 'M']"
                                >
                                    <UButton
                                        block
                                        icon="i-mdi-camera"
                                        @click="
                                            modal.open(CitizenSetMugShotModal, {
                                                user: user,
                                                'onUpdate:mugShot': ($event) => $emit('update:mugShot', $event),
                                            })
                                        "
                                    >
                                        {{ $t('components.citizens.CitizenInfoProfile.set_mug_shot') }}
                                    </UButton>
                                </UTooltip>

                                <UTooltip
                                    v-if="can('DocStoreService.CreateDocument')"
                                    :text="
                                        user.props?.wanted
                                            ? $t('components.citizens.CitizenInfoProfile.revoke_wanted')
                                            : $t('components.citizens.CitizenInfoProfile.set_wanted')
                                    "
                                    :shortcuts="['C', 'D']"
                                >
                                    <UButton block icon="i-mdi-file-document-plus" @click="openTemplates()">
                                        {{ $t('components.citizens.CitizenInfoProfile.create_new_document') }}
                                    </UButton>
                                </UTooltip>

                                <UButton block icon="i-mdi-link" @click="copyLinkToClipboard()">
                                    {{ $t('components.citizens.CitizenInfoProfile.copy_profile_link') }}
                                </UButton>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
