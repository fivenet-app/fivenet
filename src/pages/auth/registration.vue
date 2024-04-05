<script lang="ts" setup>
import RegistrationForm from '~/components/auth/RegistrationForm.vue';
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

useHead({
    title: 'components.auth.registration_form.title',
});
definePageMeta({
    title: 'components.auth.registration_form.title',
    layout: 'auth',
    requiresAuth: false,
    showCookieOptions: true,
});
</script>

<template>
    <UCard class="w-full max-w-sm bg-white/75 backdrop-blur dark:bg-white/5">
        <FiveNetLogo class="mx-auto mb-2 h-auto w-20" />

        <RegistrationForm />
    </UCard>
</template>
