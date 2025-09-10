<script lang="ts" setup>
import OAuth2Connections from '~/components/auth/account/OAuth2Connections.vue';
import type { GetAccountInfoResponse } from '~~/gen/ts/services/auth/auth';

const account = defineModel<GetAccountInfoResponse>('account', { required: true });

async function removeOAuth2Connection(provider: string): Promise<void> {
    const idx = account.value?.oauth2Connections.findIndex((v) => v.providerName === provider);
    if (idx !== undefined && idx > -1) {
        account.value?.oauth2Connections.splice(idx, 1);
    }
}
</script>

<template>
    <div v-if="account?.oauth2Providers && account.oauth2Providers.length > 0">
        <UPageCard
            :title="$t('components.auth.OAuth2Connections.title')"
            :description="$t('components.auth.OAuth2Connections.subtitle')"
        >
            <OAuth2Connections
                :providers="account.oauth2Providers"
                :connections="account.oauth2Connections"
                @disconnected="removeOAuth2Connection($event)"
            />
        </UPageCard>
    </div>
</template>
