<script lang="ts" setup>
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import CharacterSelector from '~/components/auth/CharacterSelector.vue';

useHead({
    title: 'pages.auth.character_selector.title',
});
definePageMeta({
    title: 'pages.auth.character_selector.title',
    requiresAuth: true,
    authOnlyToken: true,
    showQuickButtons: false,
});
</script>

<template>
    <ContentCenterWrapper>
        <CharacterSelector />
    </ContentCenterWrapper>
</template>
