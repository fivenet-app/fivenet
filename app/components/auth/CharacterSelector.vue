<script lang="ts" setup>
import CharacterSelectorCard from '~/components/auth/CharacterSelectorCard.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import type { Character } from '~~/gen/ts/resources/accounts/accounts';

const authStore = useAuthStore();

const { chooseCharacter } = authStore;

const { data: chars, pending: loading, refresh, error } = useLazyAsyncData('chars', () => getCharacters());

async function getCharacters(): Promise<Character[]> {
    try {
        const call = getGRPCAuthClient().getCharacters({});
        const { response } = await call;

        return response.chars;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(chars, async () => {
    // If user only has one char, auto select that char
    if (chars.value?.length === 1 && chars.value[0]?.char?.userId !== undefined) {
        await onSubmitThrottle(chars.value[0].char.userId);
    }
});

const charLockActive = computed(() => chars.value?.find((c) => c.available === false));

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (charId: number) => {
    canSubmit.value = false;

    console.log('choose character', charId);
    await chooseCharacter(charId, true).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div>
        <DataPendingBlock
            v-if="loading"
            :message="$t('common.loading', [`${$t('common.your')} ${$t('common.character', 2)}`])"
        />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.not_found', [$t('common.character', 2)])"
            :error="error"
            :retry="refresh"
        />
        <DataNoDataBlock
            v-else-if="!chars || chars.length === 0"
            :type="$t('common.character', 2)"
            :retry="refresh"
            class="w-full min-w-64"
        />

        <UCarousel
            v-else
            v-slot="{ item }"
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
            <CharacterSelectorCard
                :key="item.userId"
                :char="item.char"
                :unavailable="!item.available"
                :can-submit="canSubmit"
                @selected="onSubmitThrottle($event)"
            />
        </UCarousel>

        <UContainer v-if="charLockActive" class="mt-4" :ui="{ constrained: 'max-w-xl' }">
            <UAlert
                :ui="{ wrapper: 'relative overflow-hidden' }"
                icon="i-mdi-information-slab-circle"
                color="primary"
                variant="subtle"
                :title="$t('components.auth.CharacterSelectorCard.char_lock_alert.title')"
                :description="$t('components.auth.CharacterSelectorCard.char_lock_alert.description')"
            />
        </UContainer>
    </div>
</template>
