<script lang="ts" setup>
import ContentCenterWrapper from '../../components/partials/ContentCenterWrapper.vue';
import CharacterSelector from '../../components/auth/CharacterSelector.vue';

useHead({
    title: 'Character Selector',
});
definePageMeta({
    requiresAuth: true,
    authOnlyToken: true,
});
</script>

<template>
    <ContentCenterWrapper>
        <CharacterSelector />
    </ContentCenterWrapper>
</template>
