<script lang="ts" setup>
import type { OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';

defineProps<{
    provider: OAuth2Provider;
}>();
</script>

<template>
    <UButton icon="i-mdi-connection" :label="$t('common.connect')" :to="generateDiscordConnectURL(provider.name)" />
</template>
