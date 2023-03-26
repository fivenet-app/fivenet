<script lang="ts" setup>
import { User } from '@arpanet/gen/resources/users/users_pb';
import { ClipboardDocumentIcon, EyeIcon } from '@heroicons/vue/24/solid';
import { useStore } from '../../store/store';

const store = useStore();

const props = defineProps({
    user: {
        required: true,
        type: User,
    },
});

function addToClipboard() {
    store.commit('clipboard/addUser', props.user);
}
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
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-sm font-medium sm:pr-0">
            <div class="flex flex-row justify-end">
                <router-link v-can="'CitizenStoreService.FindUsers'"
                    :to="{ name: 'Citizens: Info', params: { id: user.getUserId() ?? 0 } }"
                    class="flex-initial text-primary-500 hover:text-primary-400">
                    <EyeIcon class="w-6 h-auto ml-auto mr-2.5" />
                </router-link>
                <button class="flex-initial text-primary-500 hover:text-primary-400" @click="addToClipboard">
                    <ClipboardDocumentIcon class="w-6 h-auto ml-auto mr-2.5" />
                </button>
            </div>
        </td>
    </tr>
</template>
