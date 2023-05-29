<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import { User } from '~~/gen/ts/resources/users/users';
import CharacterSelectorCard from './CharacterSelectorCard.vue';

const { $grpc } = useNuxtApp();

const { data: chars, pending, refresh, error } = useLazyAsyncData('chars', () => fetchCharacters());

async function fetchCharacters(): Promise<Array<User>> {
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
    <div v-else class="flex flex-row flex-wrap gap-y-2">
        <CharacterSelectorCard v-for="char in chars" :char="char" :key="char.userId" class="flex-auto max-w-xl mx-auto" />
    </div>
</template>
