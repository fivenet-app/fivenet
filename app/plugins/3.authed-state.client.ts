import { useAuthStore } from '~/stores/auth';

export default defineNuxtPlugin({
    name: 'authed-state-sync',
    parallel: true,

    setup() {
        const authStore = useAuthStore();
        const { username } = storeToRefs(authStore);

        // Use `fivenet_authed` cookie for a basic browser-wide logged in/out signal.
        const authedState = useCookie('fivenet_authed');
        useIntervalFn(async () => refreshCookie('fivenet_authed'), 1750);

        async function handleAuthedStateChange(): Promise<void> {
            if (!!authedState.value && username.value === null) {
                await authStore.chooseCharacter(undefined, true);
            } else if (!authedState.value && username.value !== null) {
                await navigateTo('/auth/logout');
            }
        }

        watch(authedState, handleAuthedStateChange, { immediate: true });
    },
});
