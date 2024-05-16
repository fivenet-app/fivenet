<script setup lang="ts">
import { emojiBlast } from 'emoji-blast';
import { useAuthStore } from '~/store/auth';

useHead({
    title: 'pages.notfound.title',
});
definePageMeta({
    title: 'pages.notfound.title',
    layout: 'landing',
    requiresAuth: false,
    redirectIfAuthed: false,
    showCookieOptions: true,
});

const authStore = useAuthStore();
const { username } = storeToRefs(authStore);
</script>

<template>
    <div class="hero flex flex-col">
        <div class="w-full flex-1 bg-black/50">
            <ULandingHero
                :title="$t('pages.notfound.page_not_found')"
                :description="$t('pages.notfound.fun_error')"
                :links="[
                    {
                        label: $t('common.back'),
                        icon: 'i-mdi-arrow-back',
                        size: 'lg',
                        color: 'gray',
                        click: () => useRouter().back(),
                    },
                    username
                        ? {
                              label: $t('common.overview'),
                              trailingIcon: 'i-mdi-home',
                              size: 'lg',
                              to: '/overview',
                          }
                        : { label: $t('common.home'), icon: 'i-mdi-home', size: 'lg', to: '/' },
                ]"
            >
                <template #headline>
                    <UBadge
                        color="gray"
                        variant="solid"
                        size="lg"
                        @click="
                            emojiBlast({
                                emojis: ['😵‍💫', '🔍', '🔎', '👀'],
                            })
                        "
                        >{{ $t('pages.notfound.error') }}</UBadge
                    >
                </template>
            </ULandingHero>
        </div>
    </div>
</template>
