<script lang="ts" setup>
import type { RoutesNamedLocations } from '@typed-router';

type Hint = { key: string; keyboard?: boolean; to?: RoutesNamedLocations };

const hints = shuffleArray([
    {
        key: 'commandpalette',
        keyboard: true,
    },
    {
        key: 'startpage',
        to: { name: 'user-settings', query: { tab: 'settings' }, hash: '#' },
    },
    {
        key: 'documenteditor',
        to: { name: 'user-settings', query: { tab: 'settings' }, hash: '#' },
    },
    {
        key: 'sociallogin_discord',
        to: { name: 'auth-account-info', query: { tab: 'oauth2Connections' }, hash: '#' },
    },
] as Hint[]);
</script>

<template>
    <UPageCard
        icon="i-mdi-information-outline"
        :title="$t('components.hints.start_text')"
        :ui="{
            wrapper: 'flex-row gap-2',
            leadingIcon: 'size-6',
            container: 'p-2 sm:p-2',
        }"
        v-bind="$attrs"
    >
        <UCarousel v-slot="{ item: hint }" :items="hints" dots loop :autoplay="{ delay: 7500 }" class="mb-6">
            <div class="mx-auto flex flex-col items-center gap-2 text-sm sm:flex-row">
                <span class="max-w-xl grow sm:max-w-full">{{ $t(`components.hints.${hint.key}.content`) }}</span>

                <div v-if="hint.keyboard || hint.to" class="flex-initial shrink-0">
                    <UKbd v-if="hint.keyboard" size="md" :value="$t(`components.hints.${hint.key}.keyboard`)" />
                    <UButton v-else-if="hint.to" size="sm" :to="hint.to" :label="$t('components.hints.click_me')" />
                </div>
            </div>
        </UCarousel>
    </UPageCard>
</template>
