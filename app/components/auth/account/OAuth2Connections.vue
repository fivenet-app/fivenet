<script lang="ts" setup>
import OAuth2Connection from '~/components/auth/account/OAuth2Connection.vue';
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
    <UDashboardPanelContent class="pb-24">
        <UDashboardSection
            :title="$t('components.auth.OAuth2Connections.title')"
            :description="$t('components.auth.OAuth2Connections.subtitle')"
        >
            <UPageGrid :ui="{ wrapper: 'grid-cols-1' }">
                <OAuth2Connection
                    v-for="provider in providers"
                    :key="provider.name"
                    :provider="provider"
                    :account="getProviderConnection(provider.name)"
                    @disconnected="$emit('disconnected', $event)"
                />
            </UPageGrid>
        </UDashboardSection>
    </UDashboardPanelContent>
</template>
