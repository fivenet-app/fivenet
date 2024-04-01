<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';

const authStore = useAuthStore();

const { accessToken } = storeToRefs(authStore);

defineProps<{
    provider: OAuth2Provider;
}>();
</script>

<template>
    <form method="get" :action="`/api/oauth2/login/${provider.name}`">
        <UButton
            type="submit"
            class="flex w-full justify-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral transition-colors hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
        >
            {{ $t('common.connect') }}
        </UButton>
        <input type="hidden" name="connect-only" value="true" />
        <input type="hidden" name="token" :value="accessToken" />
    </form>
</template>
