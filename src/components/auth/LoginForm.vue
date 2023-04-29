<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { computed } from 'vue';
import { LoginRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import { ErrorMessage, Field, useForm } from 'vee-validate';
import { object, string } from 'yup';
import { toTypedSchema } from '@vee-validate/yup';
import Alert from '~/components/partials/Alert.vue';
import config from '~/config';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();

const loginError = computed(() => authStore.$state.loginError);

async function login(username: string, password: string): Promise<void> {
    return new Promise(async (res, rej) => {
        // Start login
        authStore.loginStart();
        authStore.setActiveChar(null);
        authStore.setPermissions([]);

        const req = new LoginRequest();
        req.setUsername(username);
        req.setPassword(password);

        try {
            const resp = await $grpc.getUnAuthClient()
                .login(req, null);

            authStore.loginStop(null);
            authStore.setAccessToken(resp.getToken(), toDate(resp.getExpires()) as null | Date);

            return res();
        } catch (e) {
            authStore.loginStop((e as RpcError).message);
            authStore.setAccessToken(null, null);
            return rej(e as RpcError);
        }
    });
}

const { handleSubmit } = useForm({
    validationSchema: toTypedSchema(
        object({
            username: string().required().min(3).max(24),
            password: string().required().min(6).max(70),
        }),
    ),
});

const onSubmit = handleSubmit(async (values): Promise<void> => await login(values.username, values.password));

const providers = config.login.providers;
</script>

<template>
    <h2 class="pb-4 text-3xl text-center text-white">
        {{ $t('components.auth.login.title') }}
    </h2>

    <form @submit="onSubmit" class="my-2 space-y-6">
        <div>
            <label for="username" class="sr-only">
                {{ $t('common.username') }}
            </label>
            <div>
                <Field id="username" name="username" type="text" autocomplete="username"
                    :placeholder="$t('common.username')"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <ErrorMessage name="username" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>
        <div>
            <label for="password" class="sr-only">
                {{ $t('common.password') }}
            </label>
            <div>
                <Field id="password" name="password" type="password" autocomplete="current-password"
                    :placeholder="$t('common.password')"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <ErrorMessage name="password" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>

        <div>
            <button type="submit"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                {{ $t('common.login') }}
            </button>
        </div>
    </form>

    <form @submit.prevent="" class="my-4 space-y-2">
        <div v-for="prov in providers" class="">
            <NuxtLink :external="true" :to="`/api/oauth2/login/${prov.name}`"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                {{ prov.label }} {{ $t('common.login') }}
            </NuxtLink>
        </div>
    </form>

    <Alert v-if="loginError" :title="$t('components.auth.login.login_error')" :message="loginError" />
</template>
