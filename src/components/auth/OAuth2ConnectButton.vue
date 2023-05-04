<script lang="ts" setup>
import { OAuth2Provider } from '@fivenet/gen/resources/accounts/oauth2_pb';
import { useAuthStore } from '~/store/auth';

const authStore = useAuthStore();

defineProps({
    provider: {
        type: OAuth2Provider,
        required: true,
    },
});
</script>

<template>
    <form method="get" :action="`/api/oauth2/login/${provider.getName()}`">
        <button type="submit"
            class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
            {{ $t('common.connect') }}
        </button>
        <input type="hidden" name="connect-only" value="true" />
        <input type="hidden" name="token" :value="authStore.$state.accessToken" />
    </form>
</template>
