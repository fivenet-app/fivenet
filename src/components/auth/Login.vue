<script lang="ts" setup>
import { NavigationFailure } from 'vue-router';
import config from '~/config';
import { useAuthStore } from '~/store/auth';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';
import ForgotPasswordForm from './ForgotPasswordForm.vue';
import LoginForm from './LoginForm.vue';
import RegistrationForm from './RegistrationForm.vue';

const authStore = useAuthStore();
const route = useRoute();

const { accessToken } = storeToRefs(authStore);

const forms = ref<{ create: boolean; forgot: boolean }>({
    create: false,
    forgot: false,
});

watch(accessToken, async (): Promise<NavigationFailure | TypedRouteFromName<'auth-character-selector'> | void | undefined> => {
    if (accessToken) {
        return await navigateTo({
            name: 'auth-character-selector',
            query: route.query,
        });
    }
});

const cfg = config;
</script>

<template>
    <div class="max-w-lg w-full sm:w-[32rem] mx-auto">
        <div class="px-4 py-8 rounded-lg bg-base-850 sm:px-10">
            <img class="h-auto mx-auto mb-2 w-36" src="/images/logo.png" alt="FiveNet Logo" />

            <div v-if="forms.create">
                <RegistrationForm @back="forms.create = false" />
            </div>
            <div v-else-if="forms.forgot">
                <ForgotPasswordForm @back="forms.forgot = false" />
            </div>
            <div v-else>
                <LoginForm />
                <div class="mt-6">
                    <button
                        type="button"
                        @click="forms.forgot = true"
                        class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
                    >
                        {{ $t('components.auth.login.forgot_password') }}
                    </button>
                </div>
                <div class="mt-6" v-if="cfg.login.signupEnabled">
                    <button
                        type="button"
                        @click="forms.create = true"
                        class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
                    >
                        {{ $t('components.auth.login.register_account') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
