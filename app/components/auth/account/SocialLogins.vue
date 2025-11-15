<script lang="ts" setup>
import SocialLogin from '~/components/auth/account/SocialLogin.vue';
import type { OAuth2Account, OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';

const props = defineProps<{
    providers: OAuth2Provider[];
    connections: OAuth2Account[];
}>();

defineEmits<{
    (e: 'disconnected', provider: string): void;
}>();

function getProviderConnection(provider: string): undefined | OAuth2Account {
    return props.connections.find((v) => v.providerName === provider);
}
</script>

<template>
    <UPageGrid class="lg:grid-cols-2">
        <SocialLogin
            v-for="provider in providers"
            :key="provider.name"
            :provider="provider"
            :account="getProviderConnection(provider.name)"
            @disconnected="$emit('disconnected', $event)"
        />
    </UPageGrid>
</template>
