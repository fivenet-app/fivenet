<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { computed, ref, watch } from 'vue';
import { NavigationFailure } from 'vue-router';
import LoginForm from './LoginForm.vue';
import CreateAccountForm from './CreateAccountForm.vue';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';

const store = useAuthStore();
const route = useRoute();

const accesToken = computed(() => store.$state.accessToken);

const createAccountForm = ref(false);

watch(accesToken, async (): Promise<NavigationFailure | TypedRouteFromName<'auth-character-selector'> | void | undefined> => {
    if (accesToken) {
        return await navigateTo({ name: 'auth-character-selector', query: route.query });
    }
});
</script>

<template>
    <div class="max-w-xl mx-auto">
        <div class="px-4 py-8 rounded-lg bg-base-850 sm:px-10">
            <nuxt-img class="h-auto mx-auto mb-2 w-36" src="/images/logo.png" alt="FiveNet Logo" />

            <div v-if="!createAccountForm">
                <LoginForm />
                <div class="mt-6">
                    <button type="button" @click="createAccountForm = true"
                        class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                        {{ $t('components.auth.login.register_account') }}
                    </button>
                </div>
            </div>
            <div v-else>
                <CreateAccountForm @back="createAccountForm = false" />
            </div>
        </div>
    </div>
</template>
