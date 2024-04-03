<script lang="ts" setup>
import { OAuth2Account, OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';
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
    <UDashboardPanelContent class="pb-2">
        <UDashboardSection
            :title="$t('components.auth.oauth2_connections.title')"
            :description="$t('components.auth.oauth2_connections.subtitle')"
        >
            <OAuth2Connection
                v-for="provider in providers"
                :key="provider.name"
                :provider="provider"
                :account="getProviderConnection(provider.name)"
                @disconnected="$emit('disconnected', $event)"
            />
        </UDashboardSection>
    </UDashboardPanelContent>
</template>
