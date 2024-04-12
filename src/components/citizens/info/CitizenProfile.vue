<script lang="ts" setup>
import CharSexBadge from '~/components/partials/citizens/CharSexBadge.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { attr } from '~/composables/can';
import { User } from '~~/gen/ts/resources/users/users';

defineProps<{
    user: User;
}>();
</script>

<template>
    <div class="w-full grow lg:flex">
        <div class="flex-1 px-4 py-5 sm:p-0">
            <dl class="space-y-4 divide-y divide-gray-100 sm:space-y-0 dark:divide-gray-800">
                <div class="sm:flex sm:px-5 sm:py-4">
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.date_of_birth') }}
                    </dt>
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
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
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        {{ user?.height }}cm
                    </dd>
                </div>

                <div v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber')" class="sm:flex sm:px-5 sm:py-4">
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.phone_number') }}
                    </dt>
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        <PhoneNumberBlock :number="user.phoneNumber" />
                    </dd>
                </div>

                <div v-if="user.visum" class="sm:flex sm:px-5 sm:py-4">
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.visum') }}
                    </dt>
                    <dd class="mt-1 text-sm text-blue-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-blue-300">
                        {{ user?.visum }}
                    </dd>
                </div>

                <div
                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.BloodType')"
                    class="sm:flex sm:px-5 sm:py-4"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.blood_type') }}
                    </dt>
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        {{ user?.props?.bloodType ?? $t('common.na') }}
                    </dd>
                </div>

                <div
                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.TrafficInfractionPoints')"
                    class="sm:flex sm:px-5 sm:py-4"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.traffic_infraction_points', 2) }}
                    </dt>
                    <dd
                        class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300"
                        :class="(user?.props?.trafficInfractionPoints ?? 0) >= 10 ? 'text-error-500' : ''"
                    >
                        {{ $t('common.point', user?.props?.trafficInfractionPoints ?? 0) }}
                    </dd>
                </div>

                <div
                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.OpenFines')"
                    class="sm:flex sm:px-5 sm:py-4"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.fine') }}
                    </dt>
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        <span v-if="(user.props?.openFines ?? 0) <= 0">
                            {{ $t('common.no_open_fine') }}
                        </span>
                        <span v-else class="text-error-500">
                            {{ $n(parseInt((user?.props?.openFines ?? 0).toString()), 'currency') }}
                        </span>
                    </dd>
                </div>

                <div v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'Licenses')" class="sm:flex sm:px-5 sm:py-4">
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.license', 2) }}
                    </dt>
                    <dd class="mt-1 w-full text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        <span v-if="user?.licenses.length === 0">
                            {{ $t('common.no_licenses') }}
                        </span>
                        <ul v-else role="list" class="w-full divide-y divide-base-200 rounded-md border border-base-200">
                            <li
                                v-for="license in user?.licenses"
                                :key="license.type"
                                class="flex items-center justify-between py-3 pl-3 pr-4 text-sm"
                            >
                                <div class="flex flex-1 items-center">
                                    <UIcon name="i-mdi-license" class="size-5 shrink-0" />
                                    <span class="ml-2 flex-1 truncate" :title="`${license.type.toUpperCase()}`"
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
</template>
