<script lang="ts" setup>
import { useStore } from '../store/store';
import { computed, ref, watch } from 'vue';
import { LoginRequest } from '@arpanet/gen/services/auth/auth_pb';
import { XCircleIcon } from '@heroicons/vue/20/solid';
import { useRoute, useRouter } from 'vue-router/auto';

const store = useStore();
const router = useRouter();
const route = useRoute();

const loginError = computed(() => store.state.auth?.loginError);
const accesToken = computed(() => store.state.auth?.accessToken);

watch(accesToken, () => {
    if (accesToken) {
        router.push({ name: 'Character Selector', query: route.query });
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
    <div class="max-w-xl mx-auto">
        <div class="px-4 py-8 rounded-lg bg-base-850 sm:px-10">
            <h2 class="pb-4 text-3xl text-center text-white">Login</h2>
            <form @submit.prevent="loginSubmit" class="my-2 space-y-6" action="#">
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
                        class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">Sign
                        in</button>
                </div>
            </form>

            <div v-if="loginError" class="p-4 mt-6 rounded-md bg-error-100">
                <div class="flex">
                    <div class="flex-shrink-0">
                        <XCircleIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                    </div>
                    <div class="ml-3">
                        <h3 class="text-sm font-medium text-error-600">There was an error signing you in, please try again!
                        </h3>
                        <div class="mt-2 text-sm text-error-600">
                            {{ loginError }}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
