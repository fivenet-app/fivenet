<script lang="ts" setup>
import { useStore } from '../store/store';
import { computed, ref, watch } from 'vue';
import { LoginRequest } from '@arpanet/gen/services/auth/auth_pb';
import { XCircleIcon } from '@heroicons/vue/20/solid';
import { useRouter } from 'vue-router/auto';

const store = useStore();
const router = useRouter();

const loginError = computed(() => store.state.auth?.loginError);
const accesToken = computed(() => store.state.auth?.accessToken);

watch(accesToken, () => {
    if (accesToken) {
        router.push({ name: 'Character Selector' });
    }
});

const credentials = ref<{ username: string, password: string }>({ username: '', password: '' });

function loginSubmit() {
    const req = new LoginRequest();
    req.setUsername(credentials.value.username);
    req.setPassword(credentials.value.password);
    store.dispatch('auth/doLogin', req);
}
</script>

<template>
    <div class="m-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div class="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
            <form @submit.prevent="loginSubmit" class="space-y-6" action="#">
                <div>
                    <label for="username" class="block text-sm font-medium leading-6 text-gray-900">Username</label>
                    <div class="mt-2">
                        <input v-model="credentials.username" id="username" name="username" type="username"
                            autocomplete="username" required="true"
                            class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" />
                    </div>
                </div>

                <div>
                    <label for="password" class="block text-sm font-medium leading-6 text-gray-900">Password</label>
                    <div class="mt-2">
                        <input v-model="credentials.password" id="password" name="password" type="password"
                            autocomplete="current-password" required="true"
                            class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" />
                    </div>
                </div>

                <div>
                    <button type="submit"
                        class="flex w-full justify-center rounded-md bg-indigo-600 py-2 px-3 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Sign
                        in</button>
                </div>
            </form>

            <div v-if="loginError" class="mt-6 rounded-md bg-red-50 p-4">
                <div class="flex">
                    <div class="flex-shrink-0">
                        <XCircleIcon class="h-5 w-5 text-red-400" aria-hidden="true" />
                    </div>
                    <div class="ml-3">
                        <h3 class="text-sm font-medium text-red-800">There was an error signing you in, please try again!
                        </h3>
                        <div class="mt-2 text-sm text-red-700">
                            {{ loginError }}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
