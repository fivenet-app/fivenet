<script lang="ts" setup>
import AccountInfo from '~~/src/components/auth/AccountInfo.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

useHead({
    title: 'Account Info',
});
definePageMeta({
    title: 'Account Info',
    requiresAuth: true,
    authOnlyToken: true,
});
</script>

<template>
    <ContentWrapper>
        <AccountInfo />
    </ContentWrapper>
</template>
