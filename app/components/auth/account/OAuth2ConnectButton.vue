<script lang="ts" setup>
import { OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';

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
    </form>
</template>
