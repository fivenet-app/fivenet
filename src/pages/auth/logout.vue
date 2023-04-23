<script lang="ts" setup>
import { LogoutRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import { onBeforeMount } from 'vue';
import { useAuthStore } from '~/store/auth';
import HeroFull from '~/components/partials/HeroFull.vue';
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import Footer from '~/components/partials/Footer.vue';
import { useNotificationsStore } from '~/store/notifications';

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
const notifications = useNotificationsStore();

const accessToken = computed(() => store.$state.accessToken);

async function redirect() {
    setTimeout(async () => {
        await router.push({ name: 'index' });
    }, 1500);
}

onBeforeMount(async () => {
    store.clear();

    if (!accessToken.value) {
        redirect();
        return;
    }

    $grpc.getAuthClient()
        .logout(new LogoutRequest(), null)
        .then((resp) => {
            redirect();
        })
        .catch((err: RpcError) => {
            store.loginStop(err.message);
            notifications.dispatchNotification({ title: 'Error during logout!', content: err.message, type: 'error' });
        });
});
</script>

<template>
    <div class="h-full justify-between flex flex-col">
        <HeroFull>
            <ContentCenterWrapper class="max-w-2xl mx-auto text-center">
                <h2 class="text-4xl font-bold tracking-tight text-neutral sm:text-6xl">
                    Signed out
                </h2>
                <p class="mt-6 text-lg leading-8 text-gray-300">
                    You will be redirected to the home page in a moment.
                </p>
            </ContentCenterWrapper>
        </HeroFull>
        <Footer />
    </div>
</template>
