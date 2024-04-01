<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel, Switch } from '@headlessui/vue';
import { vMaska } from 'maska';
import { ChevronDownIcon } from 'mdi-vue3';
import CitizensList from '~/components/citizens/CitizensList.vue';
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import GenericInput from '~/composables/partials/forms/GenericInput.vue';

useHead({
    title: 'pages.citizens.title',
});
definePageMeta({
    title: 'pages.citizens.title',
    requiresAuth: true,
    permission: 'CitizenStoreService.ListCitizens',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.citizens.title')"> </UDashboardNavbar>

            <UDashboardToolbar>
                <template #default>
                    <div class="flex flex-col">
                        <div class="mx-auto flex flex-row gap-4">
                            <div class="flex-1">
                                <label for="searchName" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.citizen', 1) }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <GenericInput
                                        ref="searchInput"
                                        type="text"
                                        name="searchName"
                                        :placeholder="`${$t('common.citizen', 1)} ${$t('common.name')}`"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    />
                                </div>
                            </div>
                            <div class="flex-1">
                                <label for="dateofbirth" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.date_of_birth') }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <input
                                        v-maska
                                        type="text"
                                        name="dateofbirth"
                                        data-maska="##.##.####"
                                        :placeholder="`${$t('common.date_of_birth')} (DD.MM.YYYY)`"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div
                                v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.Wanted')"
                                class="flex-initial"
                            >
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('components.citizens.citizens_list.only_wanted') }}
                                </label>
                                <div class="relative mt-3 flex items-center">
                                    <Switch
                                        :class="[
                                            false ? 'bg-primary-600' : 'bg-gray-200',
                                            'relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2',
                                        ]"
                                    >
                                        <span class="sr-only">
                                            {{ $t('components.citizens.citizens_list.only_wanted') }}
                                        </span>
                                        <span
                                            aria-hidden="true"
                                            :class="[
                                                false ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block size-5 rounded-full bg-neutral ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </Switch>
                                </div>
                            </div>
                        </div>
                        <Disclosure v-slot="{ open }" as="div" class="pt-2">
                            <DisclosureButton class="flex w-full items-start justify-between text-left text-sm text-neutral">
                                <span class="leading-7 text-accent-200">{{ $t('common.advanced_search') }}</span>
                                <span class="ml-6 flex h-7 items-center">
                                    <ChevronDownIcon
                                        :class="[open ? 'upsidedown' : '', 'size-5 transition-transform']"
                                        aria-hidden="true"
                                    />
                                </span>
                            </DisclosureButton>
                            <DisclosurePanel class="mt-2 pr-4">
                                <div class="flex flex-row gap-2">
                                    <div
                                        v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber')"
                                        class="flex-1"
                                    >
                                        <label for="searchPhone" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.search') }}
                                            {{ $t('common.phone_number') }}
                                        </label>
                                        <div class="relative mt-2 flex items-center">
                                            <input
                                                type="tel"
                                                name="searchPhone"
                                                :placeholder="$t('common.phone_number')"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </div>
                                    </div>
                                    <div class="flex-1">
                                        <label
                                            for="trafficInfractionPoints"
                                            class="block text-sm font-medium leading-6 text-neutral"
                                        >
                                            {{ $t('common.search') }}
                                            {{ $t('common.traffic_infraction_points', 2) }}
                                        </label>
                                        <div class="relative mt-2 flex items-center">
                                            <input
                                                type="number"
                                                name="trafficInfractionPoints"
                                                :placeholder="`${$t('common.traffic_infraction_points')}`"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </div>
                                    </div>
                                    <div
                                        v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.OpenFines')"
                                        class="flex-initial"
                                    >
                                        <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('components.citizens.citizens_list.open_fine') }}
                                        </label>
                                        <div class="relative mt-2 flex items-center">
                                            <input
                                                type="number"
                                                name="fine"
                                                :placeholder="`${$t('common.fine')}`"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </div>
                                    </div>
                                </div>
                            </DisclosurePanel>
                        </Disclosure>
                    </div>
                </template>

                <template #right> </template>
            </UDashboardToolbar>

            <CitizensList />
            <ClipboardButton />
        </UDashboardPanel>
    </UDashboardPage>
</template>
