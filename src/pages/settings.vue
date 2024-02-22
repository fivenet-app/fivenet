<script lang="ts" setup>
import UserSettingsPanel from '~/components/auth/account/UserSettingsPanel.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';

useHead({
    title: 'components.auth.settings_panel.title',
});
definePageMeta({
    title: 'components.auth.settings_panel.title',
    requiresAuth: true,
    authOnlyToken: true,
    showQuickButtons: false,
    showCookieOptions: true,
});
</script>

<template>
    <ContentWrapper>
        <div class="mx-auto max-w-5xl py-2">
            <UserSettingsPanel />
        </div>
    </ContentWrapper>
</template>
