<script setup lang="ts">
import FiveNetLogo from '../partials/logos/FiveNetLogo.vue';

const { t } = useI18n();

const links = [
    {
        label: t('common.home'),
        to: '/',
    },
    {
        label: t('common.about'),
        to: '/about',
    },
];
</script>

<template>
    <UHeader :links="links">
        <template #logo> <FiveNetLogo class="h-10 w-auto" /> </template>

        <template #right>
            <UButton :label="$t('components.auth.login.title')" color="gray" to="/auth/login" />
            <UButton
                :label="$t('components.auth.registration_form.title')"
                icon="i-mdi-account-plus"
                trailing
                color="black"
                to="/auth/registration"
                class="hidden lg:flex"
            />
        </template>
    </UHeader>
</template>
