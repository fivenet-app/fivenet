<script lang="ts" setup>
import CharacterSelector from '~/components/auth/CharacterSelector.vue';
import { useAuthStore } from '~/stores/auth';

useHead({
    title: 'components.auth.CharacterSelector.title',
});

definePageMeta({
    title: 'components.auth.CharacterSelector.title',
    layout: 'auth',
    requiresAuth: true,
    authTokenOnly: true,
});

const authStore = useAuthStore();

// Clear character-scoped state before CharacterSelector is created. Doing
// this in onBeforeMount lets its auth-scoped async data start with a key that
// immediately becomes stale, causing several cancelled GetCharacters calls.
authStore.activeChar = null;
authStore.permissions = [];
authStore.attributes = [];
authStore.jobProps = null;
</script>

<template>
    <div class="max-w-full overflow-hidden">
        <UContainer class="my-[calc(var(--page-content-bottom-offset)+16*var(--spacing))] max-w-[100vw]">
            <UCard class="bg-white/75 backdrop-blur-sm dark:bg-white/5">
                <CharacterSelector />
            </UCard>
        </UContainer>

        <div class="fixed bottom-4 left-1/2 z-10 flex -translate-x-1/2 items-center justify-center">
            <UFieldGroup>
                <UButton
                    icon="i-mdi-account-cog-outline"
                    :label="$t('components.auth.AccountInfo.title')"
                    to="/auth/account-info"
                    color="neutral"
                />
                <UButton icon="i-mdi-logout" :label="$t('common.sign_out')" to="/auth/logout" color="neutral" />
            </UFieldGroup>
        </div>
    </div>
</template>
