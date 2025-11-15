<script lang="ts" setup>
import SocialLogins from '~/components/auth/account/SocialLogins.vue';
import type { GetAccountInfoResponse } from '~~/gen/ts/services/auth/auth';

const account = defineModel<GetAccountInfoResponse>('account', { required: true });

async function removeSocialLogin(provider: string): Promise<void> {
    const idx = account.value?.oauth2Connections.findIndex((v) => v.providerName === provider);
    if (idx !== undefined && idx > -1) {
        account.value?.oauth2Connections.splice(idx, 1);
    }
}
</script>

<template>
    <div v-if="account?.oauth2Providers && account.oauth2Providers.length > 0">
        <UPageCard :title="$t('components.auth.SocialLogins.title')" :description="$t('components.auth.SocialLogins.subtitle')">
            <SocialLogins
                :providers="account.oauth2Providers"
                :connections="account.oauth2Connections"
                @disconnected="removeSocialLogin($event)"
            />
        </UPageCard>
    </div>
</template>
