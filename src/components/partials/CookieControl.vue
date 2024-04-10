<script lang="ts" setup>
import { useSettingsStore } from '~/store/settings';
import { emojiBlast } from 'emoji-blast';

const settingsStore = useSettingsStore();
const { cookiesState } = storeToRefs(settingsStore);

const open = ref(cookiesState.value === undefined);
</script>

<template>
    <div>
        <UCard
            v-if="open"
            class="absolute inset-x-0 bottom-8 z-20 mx-auto w-full max-w-lg bg-white/75 backdrop-blur dark:bg-white/5"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.CookieControl.title') }}
                        <span
                            class="select-none"
                            @click="
                                emojiBlast({
                                    emojis: ['🍪'],
                                })
                            "
                            >🍪</span
                        >
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="open = false" />
                </div>
            </template>

            <div class="flex w-full flex-col gap-2">
                <p>{{ $t('components.CookieControl.subtitle') }}</p>

                <UButtonGroup class="inline-flex w-full flex-1">
                    <UButton variant="link" block class="flex-1" to="/privacy_policy">
                        {{ $t('common.privacy_policy') }}
                    </UButton>

                    <UButton variant="link" block class="flex-1" to="/imprint">
                        {{ $t('common.imprint') }}
                    </UButton>

                    <UButton variant="link" block class="flex-1" to="/api/clear-site-data">
                        {{ $t('components.CookieControl.clear_data') }}
                    </UButton>
                </UButtonGroup>

                <p class="text-xs">{{ $t('components.CookieControl.description') }}</p>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton color="black" block class="flex-1" @click="open = false">
                        {{ $t('common.close', 1) }}
                    </UButton>

                    <UButton
                        block
                        color="red"
                        class="flex-1"
                        :disabled="cookiesState === false"
                        @click="
                            cookiesState = false;
                            open = false;
                        "
                    >
                        {{ $t('common.decline', 1) }}
                    </UButton>

                    <UButton
                        block
                        color="green"
                        class="flex-1"
                        :disabled="cookiesState === true"
                        @click="
                            cookiesState = true;
                            open = false;
                        "
                    >
                        {{ $t('common.accept', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>

        <UButton
            :ui="{ rounded: 'rounded-full' }"
            icon="i-mdi-cookie-cog"
            size="xl"
            class="absolute bottom-10 right-6"
            @click="open = true"
        />
    </div>
</template>
