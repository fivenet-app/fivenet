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
    <div class="flex h-full flex-col items-center justify-center">
        <UContainer>
            <UCard class="w-full">
                <CharacterSelector />
            </UCard>
        </UContainer>

        <UButton
            icon="i-mdi-logout"
            :label="$t('common.sign_out')"
            to="/auth/logout"
            color="white"
            class="absolute bottom-4 z-10"
        />
    </div>
</template>
