<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import { User } from '~~/gen/ts/resources/users/users';
import CharacterSelectorCard from './CharacterSelectorCard.vue';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { accountID } = storeToRefs(authStore);

const { data: chars, pending, refresh, error } = useLazyAsyncData(`chars-${accountID}`, () => fetchCharacters());

async function fetchCharacters(): Promise<User[]> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getAuthClient().getCharacters({});
            const { response } = await call;

            return res(response.chars);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <DataPendingBlock v-if="pending" :message="$t('common.loading', [`${$t('common.your')} ${$t('common.character', 2)}`])" />
    <DataErrorBlock v-else-if="error" :title="$t('common.not_found', [$t('common.character', 2)])" :retry="refresh" />
    <div v-else class="grid grid-flow-row lg:grid-flow-col gap-8 mx-4">
        <CharacterSelectorCard
            v-for="char in chars"
            :char="char"
            :key="char.userId"
            class="flex-auto min-w-[30rem] w-[30rem] max-w-[30rem] mx-auto"
        />
    </div>
</template>
