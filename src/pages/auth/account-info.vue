<script lang="ts" setup>
import AccountInfo from '~/components/auth/AccountInfo.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

useHead({
    title: 'pages.auth.account_info.title',
});
definePageMeta({
    title: 'pages.auth.account_info.title',
    requiresAuth: true,
    authOnlyToken: true,
    showQuickButtons: false,
});
</script>

<template>
    <ContentWrapper>
        <AccountInfo />
    </ContentWrapper>
</template>
