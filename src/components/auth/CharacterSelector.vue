<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { User } from '~~/gen/ts/resources/users/users';
import CharacterSelectorCard from '~/components/auth/CharacterSelectorCard.vue';

const { $grpc } = useNuxtApp();
const { data: chars, pending, refresh, error } = useLazyAsyncData('chars', () => fetchCharacters());

async function fetchCharacters(): Promise<User[]> {
    try {
        const call = $grpc.getAuthClient().getCharacters({});
        const { response } = await call;

        return response.chars;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <DataPendingBlock v-if="pending" :message="$t('common.loading', [`${$t('common.your')} ${$t('common.character', 2)}`])" />
    <DataErrorBlock v-else-if="error" :title="$t('common.not_found', [$t('common.character', 2)])" :retry="refresh" />
    <div v-else class="mx-4 grid grid-flow-row gap-8 lg:grid-flow-col">
        <CharacterSelectorCard
            v-for="char in chars"
            :key="char.userId"
            :char="char"
            class="mx-auto w-[30rem] min-w-[30rem] max-w-[30rem] flex-auto"
        />
    </div>
</template>
