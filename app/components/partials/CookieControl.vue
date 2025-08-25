<script lang="ts" setup>
import { emojiBlast } from 'emoji-blast';
import { useCookiesStore } from '~/stores/cookies';

const cookiesStore = useCookiesStore();
const { cookiesState, isConsentModalOpen: open } = storeToRefs(cookiesStore);

const { website } = useAppConfig();
</script>

<template>
    <div>
        <UCard
            v-if="open"
            class="fixed inset-x-0 bottom-8 z-20 mx-auto w-full max-w-lg bg-white/75 backdrop-blur-sm dark:bg-white/5"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl leading-6 font-semibold">
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

                    <UButton class="-my-1" color="neutral" variant="ghost" icon="i-mdi-window-close" @click="open = false" />
                </div>
            </template>

            <div class="flex w-full flex-col gap-2">
                <p>{{ $t('components.CookieControl.subtitle') }}</p>

                <UButtonGroup class="inline-flex w-full flex-1">
                    <UButton
                        v-if="website.links?.privacyPolicy"
                        class="flex-1"
                        variant="link"
                        block
                        :to="website.links.privacyPolicy"
                        external
                    >
                        {{ $t('common.privacy_policy') }}
                    </UButton>

                    <UButton
                        v-if="website.links?.imprint"
                        class="flex-1"
                        variant="link"
                        block
                        :to="website.links.imprint"
                        external
                    >
                        {{ $t('common.imprint') }}
                    </UButton>

                    <UButton class="flex-1" variant="link" block to="/api/clear-site-data">
                        {{ $t('components.CookieControl.clear_data') }}
                    </UButton>
                </UButtonGroup>

                <p class="text-xs">{{ $t('components.CookieControl.description') }}</p>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton class="flex-1" color="neutral" block @click="open = false">
                        {{ $t('common.close', 1) }}
                    </UButton>

                    <UButton
                        class="flex-1"
                        block
                        color="error"
                        :variant="cookiesState === false ? 'soft' : 'solid'"
                        @click="
                            cookiesState = false;
                            open = false;
                        "
                    >
                        {{ $t('common.decline', 1) }}
                    </UButton>

                    <UButton
                        class="flex-1"
                        block
                        color="green"
                        :variant="cookiesState === true ? 'soft' : 'solid'"
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

        <UTooltip class="fixed right-6 bottom-32" :text="$t('components.CookieControl.name')">
            <UButton icon="i-mdi-cookie-cog" size="xl" @click="open = true" />
        </UTooltip>
    </div>
</template>
