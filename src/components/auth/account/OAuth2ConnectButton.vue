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
        <UButton type="submit" icon="i-mdi-connection">
            {{ $t('common.connect') }}
        </UButton>
        <input type="hidden" name="connect-only" value="true" />
        <input type="hidden" name="token" :value="accessToken" />
    </form>
</template>
