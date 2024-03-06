<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { BulletinBoardIcon, ChevronDownIcon } from 'mdi-vue3';
import ColleaguesList from '~/components/jobs/colleagues/ColleaguesList.vue';
import ColleagueActivityFeed from '~/components/jobs/colleagues/info/ColleagueActivityFeed.vue';

useHead({
    title: 'pages.jobs.colleagues.title',
});
definePageMeta({
    title: 'pages.jobs.colleagues.title',
    requiresAuth: true,
    permission: 'JobsService.ListColleagues',
});
</script>

<template>
    <div>
        <ColleaguesList />

        <div v-if="can('JobsService.GetColleague')" class="py-2 pb-14">
            <div class="px-1 sm:px-2 lg:px-4">
                <Disclosure v-slot="{ open }" as="div" class="border-neutral/20 text-neutral hover:border-neutral/70">
                    <DisclosureButton
                        :class="[
                            open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                            'flex w-full items-start justify-between border-2 border-inherit p-2 text-left transition-colors',
                        ]"
                    >
                        <span class="inline-flex items-center text-base font-semibold leading-7">
                            <BulletinBoardIcon class="mr-2 w-5 h-auto" aria-hidden="true" />
                            {{ $t('common.activity') }}
                        </span>
                        <span class="ml-6 flex h-7 items-center">
                            <ChevronDownIcon
                                :class="[open ? 'upsidedown' : '', 'h-5 w-5 transition-transform']"
                                aria-hidden="true"
                            />
                        </span>
                    </DisclosureButton>
                    <DisclosurePanel class="rounded-b-lg border-2 border-t-0 border-inherit transition-colors">
                        <ColleagueActivityFeed :show-target-user="true" />
                    </DisclosurePanel>
                </Disclosure>
            </div>
        </div>
    </div>
</template>
