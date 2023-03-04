<script lang="ts">
import { Character } from '@arpanet/gen/common/character_pb';
import { defineComponent } from 'vue';
import { KeyIcon } from '@heroicons/vue/20/solid';
import CitizenInfoFeed from './CitizenActivityFeed.vue';

export default defineComponent({
    components: {
        CitizenInfoFeed,
        KeyIcon,
    },
    props: {
        'char': {
            required: true,
            type: Character,
        },
    },
});
</script>

<template>
    <h1 class="text-4xl">{{ char.getFirstname() }}, {{ char.getLastname() }}</h1>
    <div class="overflow-hidden bg-white shadow sm:rounded-lg">
        <div class="border-t border-gray-200">
            <dl>
                <div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                    <dt class="text-sm font-medium text-gray-500">Full name</dt>
                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">{{ char.getFirstname() }}, {{
                        char.getLastname() }}</dd>
                </div>
                <div class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                    <dt class="text-sm font-medium text-gray-500">Date of Birth</dt>
                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">{{ char.getDateofbirth() }}</dd>
                </div>
                <div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                    <dt class="text-sm font-medium text-gray-500">Email address</dt>
                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">{{ char.getHeight() }}</dd>
                </div>
                <div class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                    <dt class="text-sm font-medium text-gray-500">Sex</dt>
                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">{{ char.getSex() }}</dd>
                </div>
                <div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                    <dt class="text-sm font-medium text-gray-500">Occupation</dt>
                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">{{ char.getJob() }}</dd>
                </div>
                <div class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                    <dt class="text-sm font-medium text-gray-500">Visum</dt>
                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">{{ char.getVisum() }}</dd>
                </div>
                <div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                    <dt class="text-sm font-medium text-gray-500">Licenses</dt>
                    <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                        <span v-if="!char.getLicensesList()">No Licenses</span>
                        <ul v-else role="list" class="divide-y divide-gray-200 rounded-md border border-gray-200">
                            <li v-for="license in char.getLicensesList()" class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                                <div class="flex w-0 flex-1 items-center">
                                    <KeyIcon class="h-5 w-5 flex-shrink-0 text-gray-400" aria-hidden="true" />
                                    <span class="ml-2 w-0 flex-1 truncate">{{ license.getName().toUpperCase() }}</span>
                                </div>
                            </li>
                        </ul>
                    </dd>
                </div>
                <div class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                    <CitizenInfoFeed />
                </div>
            </dl>
        </div>
    </div>
</template>
