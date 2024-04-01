<script lang="ts" setup>
import CharacterSelector from '~/components/auth/CharacterSelector.vue';
import { useAuthStore } from '~/store/auth';
import { useDocumentEditorStore } from '~/store/documenteditor';

useHead({
    title: 'components.auth.character_selector.title',
});
definePageMeta({
    title: 'components.auth.character_selector.title',
    layout: 'auth',
    requiresAuth: true,
    authOnlyToken: true,
    showQuickButtons: false,
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
    <UCard class="w-full">
        <CharacterSelector />
    </UCard>
</template>
