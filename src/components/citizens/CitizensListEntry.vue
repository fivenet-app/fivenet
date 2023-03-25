<script lang="ts" setup>
import { User } from '@arpanet/gen/resources/users/users_pb';
import { EyeIcon } from '@heroicons/vue/24/solid';

defineProps({
    user: {
        required: true,
        type: User,
    },
});
</script>

<template>
    <tr :key="user.getUserId()">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ user.getFirstname() }}, {{ user.getLastname() }}
            <span v-if="user.getProps()?.getWanted()"
                class="inline-flex items-center rounded-full bg-error-100 px-2.5 py-0.5 text-sm font-medium text-error-700 ml-1">WANTED</span>
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ user.getJobLabel() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ user.getSex().toUpperCase() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ user.getDateofbirth() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ user.getHeight() }}cm
        </td>
        <td class="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <div v-can="'CitizenStoreService.FindUsers'">
                <router-link :to="{ name: 'Citizens: Info', params: { id: user.getUserId() ?? 0 } }"
                    class="text-primary-500 hover:text-primary-400"><EyeIcon class="w-6 h-auto ml-auto mr-2.5" /></router-link>
            </div>
        </td>
    </tr>
</template>
