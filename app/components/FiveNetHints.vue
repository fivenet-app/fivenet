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
        }"
    >
        <UCarousel v-slot="{ item: hint }" :items="hints" dots loop :autoplay="{ delay: 7500 }" class="mb-6">
            <div class="mx-auto mb-2 flex items-center gap-2 text-base">
                <span class="grow">{{ $t(`components.hints.${hint.key}.content`) }}</span>

                <div v-if="hint.keyboard || hint.to" class="flex-initial shrink-0">
                    <UKbd v-if="hint.keyboard" size="md">
                        {{ $t(`components.hints.${hint.key}.keyboard`) }}
                    </UKbd>
                    <UButton v-else-if="hint.to" :to="hint.to">
                        {{ $t('components.hints.click_me') }}
                    </UButton>
                </div>
            </div>
        </UCarousel>
    </UPageCard>
</template>
