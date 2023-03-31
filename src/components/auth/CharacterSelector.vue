<script lang="ts" setup>
import { onBeforeMount } from 'vue';
import CharacterSelectorCard from './CharacterSelectorCard.vue';
import { XCircleIcon } from '@heroicons/vue/20/solid';
import { GetCharactersRequest } from '@fivenet/gen/services/auth/auth_pb';
import { User } from '@fivenet/gen/resources/users/users_pb';
import { useAuthStore } from '../../store/auth';
import { RpcError } from 'grpc-web';
import { ArrowPathIcon } from '@heroicons/vue/24/solid';

const { $grpc } = useNuxtApp();
const store = useAuthStore();

async function fetchCharacters(): Promise<Array<User>> {
    return new Promise(async (res, rej) => {
        try {
            const resp = await $grpc.getAuthClient().
                getCharacters(new GetCharactersRequest(), null);

            return res(resp.getCharsList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const { data: chars, pending, error } = await useLazyAsyncData('chars', () => fetchCharacters());

onBeforeMount(async () => {
    store.updateActiveChar(null);
});
</script>

<template>
    <div v-if="pending">
        <button type="button" disabled
            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2">
            <ArrowPathIcon class="w-12 h-12 mx-auto text-neutral animate-spin" />
            <span class="block mt-2 text-sm font-semibold text-gray-300">
                Loading your characters...
                {{ error }}
            </span>
        </button>
    </div>
    <div v-else-if="error || (chars && chars.length <= 0)" class="rounded-md bg-red-50 p-4 max-w-xs mx-auto">
        <div class="flex">
            <div class="flex-shrink-0">
                <XCircleIcon class="h-5 w-5 text-red-400" aria-hidden="true" />
            </div>
            <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">Unable to load your characters!</h3>
                <div class="mt-2 text-sm text-red-700">
                    <p>
                        Please try again a few minutes.
                    </p>
                </div>
            </div>
        </div>
    </div>
    <div v-else class="flex flex-row flex-wrap gap-y-2">
        <CharacterSelectorCard v-for="char in chars" :char="char" :key="char.getUserId()"
            class="flex-auto max-w-xl mx-auto" />
    </div>
</template>
