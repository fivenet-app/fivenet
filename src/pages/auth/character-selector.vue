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
    authOnlyToken: true,
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
    <div class="max-w-full">
        <UContainer :ui="{ constrained: 'max-w-screen' }">
            <UCard class="bg-white/75 backdrop-blur dark:bg-white/5">
                <CharacterSelector />
            </UCard>
        </UContainer>

        <div class="absolute bottom-4 z-10 flex w-full items-center justify-center">
            <UButton icon="i-mdi-logout" :label="$t('common.sign_out')" to="/auth/logout" color="white" />
        </div>
    </div>
</template>
