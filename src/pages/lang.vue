<script lang="ts" setup>
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import LanguageSwitcher from '~/components/partials/LanguageSwitcher.vue';

useHead({
    title: 'Change Language',
});
definePageMeta({
    title: 'Change Language',
    requiresAuth: false,
});
</script>

<template>
    <ContentCenterWrapper>
        <LanguageSwitcher />
    </ContentCenterWrapper>
</template>
