<script lang="ts" setup>
import { useAuthStore } from '../../store/auth';
import { computed, ref, watch } from 'vue';
import { CreateAccountRequest, LoginRequest, LoginResponse } from '@fivenet/gen/services/auth/auth_pb';
import { XCircleIcon } from '@heroicons/vue/20/solid';
import { RpcError } from 'grpc-web';
import { dispatchNotification } from '../notification';
import { NavigationFailure } from 'vue-router';
import { useForm } from 'vee-validate';
import LoginForm from './LoginForm.vue';
import CreateAccountForm from './CreateAccountForm.vue';

const { $grpc } = useNuxtApp();
const store = useAuthStore();
const router = useRouter();
const route = useRoute();

const accesToken = computed(() => store.$state.accessToken);

const createAccountForm = ref(false);

watch(accesToken, async (): Promise<NavigationFailure | void | undefined> => {
    if (accesToken) {
        return await router.push({ name: 'auth-character-selector', query: route.query });
    }
});
</script>

<template>
    <div class="max-w-xl mx-auto">
        <div class="px-4 py-8 rounded-lg bg-base-850 sm:px-10">
            <img class="h-auto mx-auto mb-2 w-36" src="/images/logo.png" alt="FiveNet Logo" />

            <div v-if="!createAccountForm">
                <LoginForm />
                <div class="mt-6">
                    <button type="button" @click="createAccountForm = true"
                        class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                        Register Account using Token
                    </button>
                </div>
            </div>
            <div v-else>
                <CreateAccountForm @back="createAccountForm = false" />
            </div>
        </div>
    </div>
</template>
