<script lang="ts" setup>
import { OAuth2Account, OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';
import GenericContainerPanel from '~/components/partials/elements/GenericContainerPanel.vue';
import OAuth2Connection from '~/components/auth/account/OAuth2Connection.vue';

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
    <GenericContainerPanel>
        <template #title>
            {{ $t('components.auth.oauth2_connections.title') }}
        </template>
        <template #description>
            {{ $t('components.auth.oauth2_connections.subtitle') }}
        </template>
        <template #default>
            <OAuth2Connection
                v-for="provider in providers"
                :key="provider.name"
                :provider="provider"
                :account="getProviderConnection(provider.name)"
                @disconnected="$emit('disconnected', $event)"
            />
        </template>
    </GenericContainerPanel>
</template>
