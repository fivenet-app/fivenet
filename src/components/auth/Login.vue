<script lang="ts" setup>
import { type NavigationFailure } from 'vue-router';
import type { TypedRouteFromName } from '@typed-router';
import { useAuthStore } from '~/store/auth';
import ForgotPasswordForm from '~/components/auth/ForgotPasswordForm.vue';
import LoginForm from '~/components/auth/LoginForm.vue';
import FormWrapper from '~/components/auth/FormWrapper.vue';

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);

const route = useRoute();

const showLogin = ref(true);

watch(accessToken, async (): Promise<NavigationFailure | TypedRouteFromName<'auth-character-selector'> | void | undefined> => {
    if (accessToken.value === null) {
        return;
    }

    return await navigateTo({
        name: 'auth-character-selector',
        query: route.query,
    });
});
</script>

<template>
    <FormWrapper>
        <template #default>
            <component :is="showLogin ? LoginForm : ForgotPasswordForm" @toggle="showLogin = !showLogin" />
        </template>
    </FormWrapper>
</template>
