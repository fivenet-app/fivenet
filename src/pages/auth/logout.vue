<script lang="ts" setup>
import { LogoutRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import { onBeforeMount } from 'vue';
import { dispatchNotification } from '~~/src/components/notification';
import { useAuthStore } from '../../store/auth';
import HeroFull from '~/components/partials/HeroFull.vue';
import ContentCenterWrapper from '~~/src/components/partials/ContentCenterWrapper.vue';

useHead({
    title: 'Logout',
});
definePageMeta({
    title: 'Logout',
    requiresAuth: true,
    authOnlyToken: true,
});

const { $grpc } = useNuxtApp();
const store = useAuthStore();
const router = useRouter();

onBeforeMount(async () => {
    store.clear();
    $grpc.getAuthClient()
        .logout(new LogoutRequest(), null)
        .then((resp) => {
            setTimeout(async () => {
                await router.push({ name: 'index' });
            }, 1500);
        })
        .catch((err: RpcError) => {
            store.loginStop(err.message);
            dispatchNotification({ title: 'Error during logout!', content: err.message, type: 'error' });
        });
});
</script>

<template>
    <HeroFull>
        <ContentCenterWrapper>
            <h2 class="text-4xl font-bold tracking-tight text-neutral sm:text-6xl">
                Signed out
            </h2>
            <p class="mt-6 text-lg leading-8 text-gray-300">
                You will be redirected to the home page in a moment.
            </p>
        </ContentCenterWrapper>
    </HeroFull>
</template>
