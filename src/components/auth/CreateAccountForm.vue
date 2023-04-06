<script lang="ts" setup>
import { CreateAccountRequest } from '@fivenet/gen/services/auth/auth_pb';
import { dispatchNotification } from '../notification';
import { RpcError } from 'grpc-web';
import { ErrorMessage, Field, useForm } from 'vee-validate';
import { object, string } from 'yup';
import { toTypedSchema } from '@vee-validate/yup';
import Alert from '../partials/Alert.vue';

const { $grpc } = useNuxtApp();

defineEmits<{
    (e: 'back'): void,
}>();

async function createAccount(regToken: string, username: string, password: string): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new CreateAccountRequest();
        req.setRegToken(regToken);
        req.setUsername(username);
        req.setPassword(password);

        try {
            await $grpc.getUnAuthClient().
                createAccount(req, null);

            dispatchNotification({ title: 'Account created successfully!', content: '', type: 'success' });
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            accountError.value = (e as RpcError).message;
            return rej(e as RpcError);
        }
    });
}

const accountError = ref('');

const { errors, handleSubmit } = useForm({
    validationSchema: toTypedSchema(
        object({
            registrationToken: string().required().length(6),
            username: string().required().min(3).max(24),
            password: string().required().min(6).max(70),
        }),
    ),
});

const onSubmit = handleSubmit(async (values): Promise<void> => await createAccount(values.registrationToken, values.username, values.password));
</script>

<template>
    <h2 class="pb-4 text-3xl text-center text-white">Create Account</h2>

    <form @submit="onSubmit" class="my-2 space-y-6">
        <div>
            <label for="regtoken" class="sr-only">Registration Token</label>
            <div>
                <Field id="regtoken" name="regtoken" type="text" inputmode="numeric" aria-describedby="hint"
                    pattern="[0-9]*" autocomplete="regtoken" placeholder="Registration Token"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-lg sm:leading-6" />
                <ErrorMessage name="registrationToken" as="p" class="mt-2 text-sm text-red-500" />
            </div>
        </div>
        <div>
            <label for="username" class="sr-only">Username</label>
            <div>
                <Field id="username" name="username" type="text" autocomplete="username" placeholder="Username"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <ErrorMessage name="Username" as="p" class="mt-2 text-sm text-red-500" />
            </div>
        </div>
        <div>
            <label for="password" class="sr-only">Password</label>
            <div>
                <Field id="password" name="password" type="password" autocomplete="current-password" placeholder="Password"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                <ErrorMessage name="username" as="p" class="mt-2 text-sm text-red-500" />
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
        <button type="button" @click="$emit('back')"
            class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
            Back to Login
        </button>
    </div>

    <Alert v-if="accountError" title="There was an error signing you in, please try again!" :message="accountError" />
</template>
