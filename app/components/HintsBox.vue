<script lang="ts" setup>
import type { RoutesNamedLocations } from '@typed-router';

const { discord } = useAppConfig();

const settingsStore = useSettingsStore();
const { eventsShowSnowflakes } = storeToRefs(settingsStore);

type Hint = { key: string; keyboard?: boolean; to?: RoutesNamedLocations; hide?: boolean };

const hints = computed(() =>
    shuffleArray(
        [
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
                hide: discord.botEnabled,
            },
            {
                key: 'toggle_event_effect',
                to: { name: 'user-settings' },
                hide: !eventsShowSnowflakes.value,
            },
        ].flatMap((h) => (!h.hide ? [h] : [])) as Hint[],
    ),
);
</script>

<template>
    <UPageCard
        icon="i-mdi-information-outline"
        :title="$t('components.hints.start_text')"
        :ui="{
            wrapper: 'flex-row gap-2',
            leadingIcon: 'size-6',
            body: 'p-0 sm:p-0',
            root: 'block',
        }"
        v-bind="$attrs"
    >
        <UCarousel
            v-slot="{ item: hint }"
            :items="hints"
            dots
            loop
            :autoplay="{ delay: 7500 }"
            :ui="{ dots: '-bottom-6' }"
            class="mb-4"
        >
            <div class="box-border w-full min-w-0">
                <div class="flex min-w-0 flex-wrap items-center gap-3 overflow-hidden">
                    <p class="min-w-0 flex-1 [overflow-wrap:anywhere] break-words hyphens-auto whitespace-normal">
                        {{ $t(`components.hints.${hint.key}.content`) }}
                    </p>

                    <div v-if="hint.keyboard || hint.to" class="shrink-0">
                        <UKbd v-if="hint.keyboard" size="md" :value="$t(`components.hints.${hint.key}.keyboard`)" />
                        <UButton v-else-if="hint.to" size="sm" :to="hint.to" :label="$t('components.hints.click_me')" />
                    </div>
                </div>
            </div>
        </UCarousel>
    </UPageCard>
</template>
