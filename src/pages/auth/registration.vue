<script lang="ts" setup>
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import PageFooter from '~/components/partials/PageFooter.vue';
import HeroFull from '~/components/partials/HeroFull.vue';
import RegistrationForm from '~/components/auth/RegistrationForm.vue';
import FormWrapper from '~/components/auth/FormWrapper.vue';

useHead({
    title: 'components.auth.registration_form.title',
});
definePageMeta({
    title: 'components.auth.registration_form.title',
    requiresAuth: false,
    showCookieOptions: true,
});
</script>

<template>
    <div class="flex h-full flex-col justify-between">
        <HeroFull>
            <ContentCenterWrapper>
                <FormWrapper>
                    <template #default>
                        <RegistrationForm />
                    </template>
                </FormWrapper>
            </ContentCenterWrapper>
        </HeroFull>
        <PageFooter />
    </div>
</template>
