<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { alpha_dash, max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import Alert from '~/components/partials/elements/Alert.vue';
import config from '~/config';
import { useAuthStore } from '~/store/auth';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();

const { loginError } = storeToRefs(authStore);
const { loginStart, loginStop, setActiveChar, setPermissions, setAccessToken } = authStore;

async function login(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        // Start login
        loginStart();
        setActiveChar(null);
        setPermissions([]);

        try {
            const call = $grpc.getUnAuthClient().login({
                username: values.username,
                password: values.password,
            });
            const { response } = await call;

            loginStop(null);
            setAccessToken(response.token, toDate(response.expires) as null | Date);

            return res();
        } catch (e) {
            loginStop((e as RpcError).message);
            setAccessToken(null, null);
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const providers = config.login.providers;

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);
defineRule('alpha_dash', alpha_dash);

interface FormData {
    registrationToken: number;
    username: string;
    password: string;
}

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        username: { required: true, min: 3, max: 24, alpha_dash: true },
        password: { required: true, min: 6, max: 70 },
    },
});

const onSubmit = handleSubmit(async (values): Promise<void> => await login(values));
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
                <VeeField
                    name="username"
                    type="text"
                    autocomplete="username"
                    :placeholder="$t('common.username')"
                    :label="$t('common.username')"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                />
                <VeeErrorMessage name="username" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>
        <div>
            <label for="password" class="sr-only">
                {{ $t('common.password') }}
            </label>
            <div>
                <VeeField
                    name="password"
                    type="password"
                    autocomplete="current-password"
                    :placeholder="$t('common.password')"
                    :label="$t('common.password')"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                />
                <VeeErrorMessage name="password" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </div>

        <div>
            <button
                type="submit"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                :disabled="!meta.valid"
                :class="[
                    !meta.valid
                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                ]"
            >
                {{ $t('common.login') }}
            </button>
        </div>
    </form>

    <div class="my-4 space-y-2">
        <div v-for="prov in providers" class="">
            <NuxtLink
                :external="true"
                :to="`/api/oauth2/login/${prov.name}`"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
            >
                {{ prov.label }} {{ $t('common.login') }}
            </NuxtLink>
        </div>
    </div>

    <Alert
        v-if="loginError"
        :title="$t('components.auth.login.login_error')"
        :message="loginError.startsWith('errors.') ? $t(loginError) : loginError"
    />
</template>
