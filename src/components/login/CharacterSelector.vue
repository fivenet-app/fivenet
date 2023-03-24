<script lang="ts" setup>
import { ref, onBeforeMount } from 'vue';
import CharacterSelectorCard from './CharacterSelectorCard.vue';
import { XCircleIcon } from '@heroicons/vue/20/solid';
import { GetCharactersRequest } from '@arpanet/gen/services/auth/auth_pb';
import { User } from '@arpanet/gen/resources/users/users_pb';
import { getAuthClient } from '../../grpc/grpc';
import { useStore } from '../../store/store';

const store = useStore();

const chars = ref<Array<User>>([]);

async function fetchCharacters() {
    console.log("FETCH CHARS");
    return getAuthClient().
        getCharacters(new GetCharactersRequest(), null).
        then((resp) => {
            chars.value = resp.getCharsList();
        });
}

onBeforeMount(() => {
    store.dispatch('auth/updateActiveChar', null);

    // Fetch user's characters
    fetchCharacters();
});
</script>

<template>
    <div v-if="chars.length <= 0" class="rounded-md bg-red-50 p-4">
        <div class="flex">
            <div class="flex-shrink-0">
                <XCircleIcon class="h-5 w-5 text-red-400" aria-hidden="true" />
            </div>
            <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">Unable to load your characters!</h3>
                <div class="mt-2 text-sm text-red-700">
                    <p>Please try again a few minutes.</p>
                </div>
            </div>
        </div>
    </div>
    <ul v-else role="list" class="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
        <CharacterSelectorCard v-for="char in chars" :char="char" :key="char.getUserId()" />
    </ul>
</template>
