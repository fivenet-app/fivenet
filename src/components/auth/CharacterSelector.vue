<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { User } from '~~/gen/ts/resources/users/users';
import CharacterSelectorCard from '~/components/auth/CharacterSelectorCard.vue';
import { useAuthStore } from '~/store/auth';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();

const { chooseCharacter } = authStore;

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

watch(chars, async () => {
    // If user only has one char, auto select that char
    if (chars.value?.length === 1) {
        await chooseCharacter(chars.value[0].userId);
    }
});
</script>

<template>
    <DataPendingBlock v-if="pending" :message="$t('common.loading', [`${$t('common.your')} ${$t('common.character', 2)}`])" />
    <DataErrorBlock v-else-if="error" :title="$t('common.not_found', [$t('common.character', 2)])" :retry="refresh" />
    <template v-else>
        <div class="md:mx-4 grid grid-flow-row auto-rows-max gap-8 md:grid-flow-col">
            <CharacterSelectorCard
                v-for="char in chars"
                :key="char.userId"
                :char="char"
                class="mx-auto w-[30rem] min-w-[30rem] max-w-[30rem] flex-auto"
            />
        </div>
    </template>
</template>
