<script lang="ts" setup>
import CharacterSelectorCard from './CharacterSelectorCard.vue';
import { GetCharactersRequest } from '@fivenet/gen/services/auth/auth_pb';
import { User } from '@fivenet/gen/resources/users/users_pb';
import { useAuthStore } from '../../store/auth';
import { RpcError } from 'grpc-web';
import DataPendingBlock from '../partials/DataPendingBlock.vue';
import DataErrorBlock from '../partials/DataErrorBlock.vue';

const { $grpc } = useNuxtApp();
const store = useAuthStore();

const { data: chars, pending, refresh, error } = await useLazyAsyncData('chars', () => fetchCharacters());

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

onBeforeMount(async () => {
    await store.updateActiveChar(null);
    await store.updatePermissions([]);
});
</script>

<template>
    <DataPendingBlock v-if="pending" message="Loading your characters..." />
    <DataErrorBlock v-else-if="error" title="Unable to load your characters!" :retry="refresh" />
    <div v-else class="flex flex-row flex-wrap gap-y-2">
        <CharacterSelectorCard v-for="char in chars" :char="char" :key="char.getUserId()"
            class="flex-auto max-w-xl mx-auto" />
    </div>
</template>
