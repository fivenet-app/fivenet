<script lang="ts" setup>
import CharacterSelectorCard from './CharacterSelectorCard.vue';
import { GetCharactersRequest } from '@fivenet/gen/services/auth/auth_pb';
import { User } from '@fivenet/gen/resources/users/users_pb';
import { useAuthStore } from '~/store/auth';
import { RpcError } from 'grpc-web';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();

const { data: chars, pending, refresh, error } = useLazyAsyncData('chars', () => fetchCharacters());

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
    await Promise.all([
        authStore.updateActiveChar(null),
        authStore.updatePermissions([])
    ]);
});
</script>

<template>
    <DataPendingBlock v-if="pending"
        :message="$t('common.loading', [`${$t('common.your')} ${$t('common.character', 2)}`])" />
    <DataErrorBlock v-else-if="error" :title="$t('common.not_found', [$t('common.character', 2)])" :retry="refresh" />
    <div v-else class="flex flex-row flex-wrap gap-y-2">
        <CharacterSelectorCard v-for="char in chars" :char="char" :key="char.getUserId()"
            class="flex-auto max-w-xl mx-auto" />
    </div>
</template>
