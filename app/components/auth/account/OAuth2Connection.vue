<script lang="ts" setup>
import OAuth2ConnectButton from '~/components/auth/account/OAuth2ConnectButton.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import type { OAuth2Account, OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';

defineProps<{
    provider: OAuth2Provider;
    account?: OAuth2Account;
}>();

const emit = defineEmits<{
    (e: 'disconnected', provider: string): void;
}>();

async function disconnectOAuth2Connection(provider: OAuth2Provider): Promise<void> {
    try {
        await getGRPCAuthClient().deleteOAuth2Connection({
            provider: provider.name,
        });

        emit('disconnected', provider.name);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const modal = useModal();
</script>

<template>
    <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-5 sm:py-4">
        <dt class="text-sm font-medium">
            <UButton variant="link" :external="true" :to="provider.homepage" target="_blank" class="inline-flex gap-2">
                <img v-if="!provider.icon?.startsWith('i-')" :src="provider.icon" :alt="provider.name" class="size-10" />
                <UIcon
                    v-else
                    :name="provider.icon"
                    class="size-10"
                    :style="provider.name === 'discord' && { color: '#7289da' }"
                />

                {{ provider.label }}
            </UButton>
        </dt>
        <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
            <div v-if="account !== undefined" class="flex items-center justify-between">
                <div class="inline-flex items-center gap-4">
                    <UAvatar size="md" :src="account.avatar" :alt="$t('common.image')" />

                    <span class="text-left" :title="`ID: ${account.externalId}`">
                        {{ account.username }}
                    </span>
                </div>

                <UButton
                    icon="i-mdi-close-circle"
                    color="red"
                    @click="
                        modal.open(ConfirmModal, {
                            confirm: async () => disconnectOAuth2Connection(provider),
                        })
                    "
                >
                    {{ $t('common.disconnect') }}
                </UButton>
            </div>
            <div v-else class="flex flex-row-reverse">
                <template v-if="isNUIAvailable()">
                    <p class="ml-4 text-end text-sm">
                        {{ $t('system.not_supported_on_tablet.title') }}
                    </p>
                </template>
                <template v-else>
                    <OAuth2ConnectButton class="self-end" :provider="provider" />
                </template>
            </div>
        </dd>
    </div>
</template>
