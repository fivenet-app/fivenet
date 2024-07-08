<script lang="ts" setup>
import CharacterSelector from '~/components/auth/CharacterSelector.vue';
import { useAuthStore } from '~/store/auth';
import { useDocumentEditorStore } from '~/store/documenteditor';

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
const documentEditorStore = useDocumentEditorStore();

const { setActiveChar, setPermissions, setJobProps } = authStore;

onMounted(async () => {
    setActiveChar(null);
    setPermissions([]);
    setJobProps(undefined);
    documentEditorStore.clear();
});
</script>

<template>
    <div class="max-w-full overflow-hidden">
        <UContainer :ui="{ constrained: 'max-w-screen' }">
            <UCard class="bg-white/75 backdrop-blur dark:bg-white/5">
                <CharacterSelector />
            </UCard>
        </UContainer>

        <div class="absolute bottom-4 left-1/2 z-10 flex -translate-x-1/2 items-center justify-center">
            <UButtonGroup>
                <UButton
                    icon="i-mdi-account-cog-outline"
                    :label="$t('components.auth.AccountInfo.title')"
                    to="/auth/account-info"
                    color="white"
                />
                <UButton icon="i-mdi-logout" :label="$t('common.sign_out')" to="/auth/logout" color="white" />
            </UButtonGroup>
        </div>
    </div>
</template>
