<script lang="ts" setup>
import { useAuthStore } from '../../store/auth';
import { computed, ref, watch } from 'vue';
import { CreateAccountRequest, LoginRequest, LoginResponse } from '@fivenet/gen/services/auth/auth_pb';
import { XCircleIcon } from '@heroicons/vue/20/solid';
import { RpcError } from 'grpc-web';
import { dispatchNotification } from '../notification';

const { $grpc } = useNuxtApp();
const store = useAuthStore();
const router = useRouter();
const route = useRoute();

const loginError = computed(() => store.$state.loginError);
const accesToken = computed(() => store.$state.accessToken);

watch(accesToken, () => {
    if (accesToken) {
        router.push({ name: 'auth-character-selector', query: route.query });
    }
});

const createAccountForm = ref(false);
const credentials = ref<{ username: string, password: string, }>({ username: '', password: '', });

async function login(): Promise<void> {
    const req = new LoginRequest();
    req.setUsername(credentials.value.username);
    req.setPassword(credentials.value.password);

    // Start login
    store.loginStart();
    store.updateActiveChar(null);
    store.updatePermissions([]);

    try {
        const resp = await $grpc.getUnAuthClient()
            .login(req, null);

        store.loginStop(null);
        store.updateAccessToken(resp.getToken());
    } catch (e) {
        store.loginStop((e as RpcError).message);
        store.updateAccessToken(null);
        return;
    }
}

const accountInfo = ref<{ regToken: string, username: string, password: string, }>({
    regToken: '', username: '', password: '',
});

async function createAccount(): Promise<void> {
    const req = new CreateAccountRequest();
    req.setRegToken(accountInfo.value?.regToken);
    req.setUsername(accountInfo.value?.username);
    req.setPassword(accountInfo.value?.password);

    try {
        await $grpc.getUnAuthClient().
            createAccount(req, null);

        createAccountForm.value = false;
        dispatchNotification({ title: 'Account created successfully!', content: '', type: 'success' });
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
        return;
    }
}
</script>

<template>
    <div class="max-w-xl mx-auto">
        <div class="px-4 py-8 rounded-lg bg-base-850 sm:px-10">
            <img class="h-auto mx-auto mb-2 w-36" src="/images/logo.png" alt="FiveNet Logo" />

            <div v-if="!createAccountForm">
                <h2 class="pb-4 text-3xl text-center text-white">Login</h2>

                <form @submit.prevent="login" class="my-2 space-y-6">
                    <div>
                        <label for="username" class="sr-only">Username</label>
                        <div>
                            <input v-model="credentials.username" id="username" name="username" type="text"
                                autocomplete="username" placeholder="Username" required="true"
                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                        </div>
                    </div>
                    <div>
                        <label for="password" class="sr-only">Password</label>
                        <div>
                            <input v-model="credentials.password" id="password" name="password" type="password"
                                autocomplete="current-password" placeholder="Password" required="true"
                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                        </div>
                    </div>

                    <div>
                        <button type="submit"
                            class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                            Login
                        </button>
                    </div>
                </form>

                <div v-if="loginError" class="p-4 mt-6 rounded-md bg-error-100">
                    <div class="flex">
                        <div class="flex-shrink-0">
                            <XCircleIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                        </div>
                        <div class="ml-3">
                            <h3 class="text-sm font-medium text-error-600">There was an error signing you in, please try
                                again!
                            </h3>
                            <div class="mt-2 text-sm text-error-600">
                                {{ loginError }}
                            </div>
                        </div>
                    </div>
                </div>

                <div class="mt-6">
                    <button type="button" @click="createAccountForm = true"
                        class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                        Register Account using Token
                    </button>
                </div>
            </div>
            <div v-else>
                <h2 class="pb-4 text-3xl text-center text-white">Create Account</h2>

                <form @submit.prevent="createAccount" class="my-2 space-y-6">
                    <div>
                        <label for="regtoken" class="sr-only">Registration Token</label>
                        <div>
                            <input v-model="accountInfo.regToken" id="regtoken" name="regtoken" type="text"
                                inputmode="numeric" aria-describedby="hint" pattern="[0-9]*" autocomplete="regtoken"
                                placeholder="Registration Token" required="true"
                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-lg sm:leading-6" />
                        </div>
                    </div>
                    <div>
                        <label for="username" class="sr-only">Username</label>
                        <div>
                            <input v-model="accountInfo.username" id="username" name="username" type="text"
                                autocomplete="username" placeholder="Username" required="true"
                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                        </div>
                    </div>
                    <div>
                        <label for="password" class="sr-only">Password</label>
                        <div>
                            <input v-model="accountInfo.password" id="password" name="password" type="password"
                                autocomplete="current-password" placeholder="Password" required="true"
                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                        </div>
                    </div>

                    <div>
                        <button type="submit"
                            class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                            Create Account
                        </button>
                    </div>
                </form>

                <div class="mt-6">
                    <button type="button" @click="createAccountForm = false"
                        class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                        Back to Login
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
