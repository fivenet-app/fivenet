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

const { auth } = useAppConfig();

onBeforeMount(() => {
    if (!auth.signupEnabled) {
        navigateTo('/auth/login');
    }
});
</script>

<template>
    <div class="w-full max-w-md pt-[calc(var(--page-content-bottom-offset)+16*var(--spacing))] pb-2">
        <UPageCard class="w-full bg-white/75 backdrop-blur-sm dark:bg-white/5">
            <FiveNetLogo class="mx-auto mb-2 h-auto w-20" />

            <RegistrationForm />
        </UPageCard>
    </div>
</template>
