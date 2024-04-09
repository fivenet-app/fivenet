<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { User } from '~~/gen/ts/resources/users/users';
import CharacterSelectorCard from '~/components/auth/CharacterSelectorCard.vue';
import { useAuthStore } from '~/store/auth';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();

const { chooseCharacter } = authStore;

const { data: chars, pending: loading, refresh, error } = useLazyAsyncData('chars', () => fetchCharacters());

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
    <DataPendingBlock v-if="loading" :message="$t('common.loading', [`${$t('common.your')} ${$t('common.character', 2)}`])" />
    <DataErrorBlock v-else-if="error" :title="$t('common.not_found', [$t('common.character', 2)])" :retry="refresh" />
    <DataNoDataBlock v-else-if="!chars" :title="$t('common.not_found', [$t('common.character', 2)])" />
    <template v-else>
        <UCarousel
            v-slot="{ item }"
            class="backdrop-blur"
            :items="chars"
            arrows
            :prev-button="{
                color: 'gray',
                icon: 'i-mdi-arrow-left',
            }"
            :next-button="{
                color: 'gray',
                icon: 'i-mdi-arrow-right',
            }"
            :ui="{ item: 'basis-full sm:basis-1/4', container: 'rounded-lg' }"
        >
            <CharacterSelectorCard :key="item.userId" :char="item" />
        </UCarousel>
    </template>
</template>
