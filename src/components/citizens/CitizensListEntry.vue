<script lang="ts" setup>
import { User } from '@arpanet/gen/resources/users/users_pb';

defineProps({
    'user': {
        required: true,
        type: User,
    },
});
</script>

<template>
    <tr :key="user.getUserId()">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-white sm:pl-0">
            {{ user.getFirstname() }}, {{ user.getLastname() }}
            <span v-if="user.getProps()?.getWanted()"
                class="inline-flex items-center rounded-md bg-red-100 px-2.5 py-0.5 text-sm font-medium text-red-800">WANTED</span>
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ user.getJobLabel() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ user.getSex().toUpperCase() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ user.getDateofbirth() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ user.getHeight() }}cm
        </td>
        <td class="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <div v-can="'CitizenStoreService.FindUsers'">
                <router-link :to="{ path: '/citizens/:id', params: { id: user.getUserId().toString() } }"
                    class="text-indigo-400 hover:text-indigo-300">VIEW</router-link>
            </div>
        </td>
    </tr>
</template>
