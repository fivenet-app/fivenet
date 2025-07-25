<script lang="ts" setup>
import EmailInfoPopover from '~/components/mailer/EmailInfoPopover.vue';
import CharSexBadge from '~/components/partials/citizens/CharSexBadge.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import type { User } from '~~/gen/ts/resources/users/users';

defineProps<{
    user: User;
}>();

const { attr } = useAuth();
</script>

<template>
    <div class="w-full grow lg:flex">
        <div class="flex-1 px-4 py-5 sm:p-0">
            <dl class="2xl:grid 2xl:grid-cols-2">
                <div
                    class="hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-b border-gray-100 py-1 sm:flex sm:px-5 sm:py-4 dark:border-gray-800"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.date_of_birth') }}
                    </dt>
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        {{ user.dateofbirth }}
                    </dd>
                </div>

                <div
                    class="hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-b border-gray-100 py-1 sm:flex sm:px-5 sm:py-4 dark:border-gray-800"
                >
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

                <div
                    class="hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-b border-gray-100 py-1 sm:flex sm:px-5 sm:py-4 dark:border-gray-800"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.height') }}
                    </dt>
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        {{ user?.height ? user.height + 'cm' : $t('common.unknown') }}
                    </dd>
                </div>

                <div
                    v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'PhoneNumber').value"
                    class="hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-b border-gray-100 py-1 sm:flex sm:px-5 sm:py-4 dark:border-gray-800"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.phone_number') }}
                    </dt>
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        <PhoneNumberBlock :number="user.phoneNumber" />
                    </dd>
                </div>

                <div
                    v-if="user.visum"
                    class="hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-b border-gray-100 py-1 sm:flex sm:px-5 sm:py-4 dark:border-gray-800"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.visum') }}
                    </dt>
                    <dd class="mt-1 text-sm text-blue-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-blue-300">
                        {{ user?.visum }}
                    </dd>
                </div>

                <div
                    v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.BloodType').value"
                    class="hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-b border-gray-100 py-1 sm:flex sm:px-5 sm:py-4 dark:border-gray-800"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.blood_type') }}
                    </dt>
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        {{ user?.props?.bloodType ?? $t('common.na') }}
                    </dd>
                </div>

                <div
                    v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.TrafficInfractionPoints').value"
                    class="hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-b border-gray-100 py-1 sm:flex sm:px-5 sm:py-4 dark:border-gray-800"
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
                    v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.OpenFines').value"
                    class="hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-b border-gray-100 py-1 sm:flex sm:px-5 sm:py-4 dark:border-gray-800"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.fine') }}
                    </dt>
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        <span v-if="(user.props?.openFines ?? 0) <= 0">
                            {{ $t('common.no_open_fine') }}
                        </span>
                        <span v-else class="text-error-500">
                            {{ $n(user?.props?.openFines ?? 0, 'currency') }}
                        </span>
                    </dd>
                </div>

                <div
                    v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.Labels').value"
                    class="hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-b border-gray-100 py-1 sm:flex sm:px-5 sm:py-4 dark:border-gray-800"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.label', 2) }}
                    </dt>
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        <p v-if="!user.props?.labels?.list.length" class="text-sm leading-6">
                            {{ $t('common.none', [$t('common.label', 2)]) }}
                        </p>
                        <template v-else>
                            <div class="flex max-w-80 flex-row flex-wrap gap-1">
                                <UBadge
                                    v-for="label in user.props?.labels?.list"
                                    :key="label.name"
                                    class="justify-between gap-2"
                                    :class="isColorBright(hexToRgb(label.color, RGBBlack)!) ? '!text-black' : '!text-white'"
                                    :style="{ backgroundColor: label.color }"
                                    size="md"
                                >
                                    {{ label.name }}
                                </UBadge>
                            </div>
                        </template>
                    </dd>
                </div>

                <div
                    v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.Email').value"
                    class="hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-b border-gray-100 py-1 sm:flex sm:px-5 sm:py-4 dark:border-gray-800"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.mail', 1) }}
                    </dt>
                    <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        <EmailInfoPopover :email="user?.props?.email" />
                    </dd>
                </div>

                <div
                    v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'Licenses').value"
                    class="hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-b border-gray-100 py-1 sm:flex sm:px-5 sm:py-4 dark:border-gray-800"
                >
                    <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                        {{ $t('common.license', 2) }}
                    </dt>
                    <dd class="mt-1 w-full text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                        <span v-if="user?.licenses.length === 0">
                            {{ $t('common.no_licenses') }}
                        </span>
                        <ul v-else class="w-full divide-y divide-base-200 rounded-md border border-base-200" role="list">
                            <li
                                v-for="license in user?.licenses"
                                :key="license.type"
                                class="flex items-center justify-between py-3 pl-3 pr-4 text-sm"
                            >
                                <div class="flex flex-1 items-center">
                                    <UIcon class="size-5 shrink-0" name="i-mdi-license" />
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
