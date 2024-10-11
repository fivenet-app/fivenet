<script lang="ts" setup>
import OAuth2ConnectButton from '~/components/auth/account/OAuth2ConnectButton.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { useNotificatorStore } from '~/store/notificator';
import type { OAuth2Account, OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

defineProps<{
    provider: OAuth2Provider;
    account?: OAuth2Account;
}>();

const emit = defineEmits<{
    (e: 'disconnected', provider: string): void;
}>();

const notifications = useNotificatorStore();

async function disconnectOAuth2Connection(provider: OAuth2Provider): Promise<void> {
    try {
        await getGRPCAuthClient().deleteOAuth2Connection({
            provider: provider.name,
        });

        notifications.add({
            title: { key: 'notifications.auth.oauth2_connect.disconnected.title', parameters: {} },
            description: {
                key: 'notifications.auth.oauth2_connect.disconnected.content',
            },
            type: NotificationType.SUCCESS,
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
    <UPageCard
        :ui="{
            body: {
                padding: 'px-4 py-4 sm:p-4',
            },
            icon: { wrapper: 'mb-1' },
        }"
    >
        <template #title>
            <div class="flex flex-1 gap-2">
                <UButton
                    variant="ghost"
                    :external="true"
                    :to="provider.homepage"
                    target="_blank"
                    class="inline-flex flex-1 gap-2"
                >
                    <img v-if="!provider.icon?.startsWith('i-')" :src="provider.icon" :alt="provider.name" class="size-10" />
                    <UIcon
                        v-else
                        :name="provider.icon"
                        class="size-10"
                        :style="provider.name === 'discord' && { color: '#7289da' }"
                    />

                    <div class="flex items-center gap-1.5 text-base font-semibold text-gray-900 dark:text-white">
                        {{ provider.label }}
                    </div>
                </UButton>

                <div v-if="account !== undefined" class="flex items-center justify-between">
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
                        <OAuth2ConnectButton :provider="provider" />
                    </template>
                </div>
            </div>
        </template>

        <template v-if="account" #footer>
            <div class="inline-flex items-center gap-4">
                <UAvatar size="md" :src="account.avatar" :alt="$t('common.image')" />

                <span class="text-left" :title="`ID: ${account.externalId}`">
                    {{ account.username }}
                </span>
            </div>
        </template>
    </UPageCard>
</template>
