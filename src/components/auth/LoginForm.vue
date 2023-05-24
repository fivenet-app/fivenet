<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { LoginRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import { defineRule } from 'vee-validate';
import Alert from '~/components/partials/Alert.vue';
import config from '~/config';
import { required, min, max, alpha_dash } from '@vee-validate/rules';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();

const { loginError } = storeToRefs(authStore);
const { loginStart, loginStop, setActiveChar, setPermissions, setAccessToken } = authStore;

const form = ref<{ username: string; password: string; }>({ username: '', password: '' });

async function login(): Promise<void> {
    return new Promise(async (res, rej) => {
        // Start login
        loginStart();
        setActiveChar(null);
        setPermissions([]);

        const req = new LoginRequest();
        req.setUsername(form.value.username);
        req.setPassword(form.value.password);

        try {
            const resp = await $grpc.getUnAuthClient()
                .login(req, null);

            loginStop(null);
            setAccessToken(resp.getToken(), toDate(resp.getExpires()) as null | Date);

            return res();
        } catch (e) {
            loginStop((e as RpcError).message);
            setAccessToken(null, null);
            return rej(e as RpcError);
        }
    });
}

const providers = config.login.providers;

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);
defineRule('alpha_dash', alpha_dash);
</script>

<template>
    <h2 class="pb-4 text-3xl text-center text-white">
        {{ $t('components.auth.login.title') }}
    </h2>

    <VeeForm @submit="login" class="my-2 space-y-6">
        <div>
            <label for="username" class="sr-only">
                {{ $t('common.username') }}
            </label>
            <div>
                <VeeField name="username" type="text" autocomplete="username" :placeholder="$t('common.username')"
                    :label="$t('common.username')" v-model="form.username"
                    :rules="{ required: true, min: 3, max: 24, alpha_dash: true }"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <VeeErrorMessage name="username" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>
        <div>
            <label for="password" class="sr-only">
                {{ $t('common.password') }}
            </label>
            <div>
                <VeeField name="password" type="password" autocomplete="current-password" :placeholder="$t('common.password')"
                    :label="$t('common.password')" v-model="form.password" :rules="{ required: true, min: 6, max: 70 }"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <VeeErrorMessage name="password" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>

        <div>
            <button type="submit"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                {{ $t('common.login') }}
            </button>
        </div>
    </VeeForm>

    <div class="my-4 space-y-2">
        <div v-for="prov in providers" class="">
            <NuxtLink :external="true" :to="`/api/oauth2/login/${prov.name}`"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                {{ prov.label }} {{ $t('common.login') }}
            </NuxtLink>
        </div>
    </div>

    <Alert v-if="loginError" :title="$t('components.auth.login.login_error')" :message="loginError" />
</template>
