<script lang="ts" setup>
import CharacterSelector from '~/components/auth/CharacterSelector.vue';
import { useAuthStore } from '~/store/auth';
import { useDocumentEditorStore } from '~/store/documenteditor';

useHead({
    title: 'components.auth.character_selector.title',
});
definePageMeta({
    title: 'components.auth.character_selector.title',
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
    <main class="flex size-full">
        <div class="my-auto py-2 max-sm:w-full sm:mx-auto lg:mx-auto">
            <CharacterSelector />
        </div>
    </main>
</template>
