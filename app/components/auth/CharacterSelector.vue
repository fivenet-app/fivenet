<script lang="ts" setup>
import CharacterSelectorCard from '~/components/auth/CharacterSelectorCard.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/stores/auth';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import type { Character } from '~~/gen/ts/resources/accounts/accounts';

const authAuthClient = await getAuthAuthClient();
const authStore = useAuthStore();

const { chooseCharacter } = authStore;

const { data: chars, status, refresh, error } = useLazyAsyncData('chars', () => getCharacters());

async function getCharacters(): Promise<Character[]> {
    try {
        const call = authAuthClient.getCharacters({});
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

const charLockActive = computed(() => chars.value?.some((c) => c.available === false) ?? false);

const cardsRef = useTemplateRef('cardsRef');

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (charId: number) => {
    canSubmit.value = false;

    await chooseCharacter(charId, true).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div>
        <DataPendingBlock
            v-if="isRequestPending(status)"
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
            class="w-full min-w-64"
            :type="$t('common.character', 2)"
            :retry="refresh"
        />

        <div v-else class="relative overflow-hidden rounded-lg">
            <div ref="cardsRef" class="no-scrollbar relative flex w-full snap-x snap-mandatory overflow-x-auto scroll-smooth">
                <CharacterSelectorCard
                    v-for="char in chars"
                    :key="char.char!.userId"
                    class="basis-full sm:basis-1/3 lg:basis-1/4"
                    :char="char.char!"
                    :unavailable="!char.available"
                    :can-submit="canSubmit"
                    @selected="onSubmitThrottle($event)"
                />
            </div>

            <div class="flex items-center justify-between">
                <UButton
                    class="absolute start-4 top-1/2 -translate-y-1/2 transform rounded-full rtl:[&_span:first-child]:rotate-180"
                    color="gray"
                    icon="i-mdi-arrow-left"
                    aria-label="Prev"
                    @click="cardsRef?.scrollLeft !== null && (cardsRef!.scrollLeft -= 395)"
                />

                <UButton
                    class="absolute end-4 top-1/2 -translate-y-1/2 transform rounded-full rtl:[&_span:last-child]:rotate-180"
                    color="gray"
                    icon="i-mdi-arrow-right"
                    aria-label="Next"
                    @click="cardsRef?.scrollLeft !== null && (cardsRef!.scrollLeft += 395)"
                />
            </div>
        </div>

        <UContainer v-if="charLockActive" class="mt-4" :ui="{ constrained: 'max-w-xl' }">
            <UAlert
                :ui="{ wrapper: 'relative overflow-hidden' }"
                icon="i-mdi-information-outline"
                color="primary"
                variant="subtle"
                :title="$t('components.auth.CharacterSelectorCard.char_lock_alert.title')"
                :description="$t('components.auth.CharacterSelectorCard.char_lock_alert.description')"
            />
        </UContainer>
    </div>
</template>
