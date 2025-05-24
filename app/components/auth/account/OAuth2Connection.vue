<script lang="ts" setup>
import { NuxtImg } from '#components';
import OAuth2ConnectButton from '~/components/auth/account/OAuth2ConnectButton.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useNotificatorStore } from '~/stores/notificator';
import { useSettingsStore } from '~/stores/settings';
import type { OAuth2Account, OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

defineProps<{
    provider: OAuth2Provider;
    account?: OAuth2Account;
}>();

const emit = defineEmits<{
    (e: 'disconnected', provider: string): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

async function disconnectOAuth2Connection(provider: OAuth2Provider): Promise<void> {
    try {
        await $grpc.auth.auth.deleteOAuth2Connection({
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
                    class="inline-flex flex-1 gap-2"
                    variant="ghost"
                    :external="true"
                    :to="provider.homepage"
                    target="_blank"
                >
                    <NuxtImg
                        v-if="!provider.icon?.startsWith('i-')"
                        class="size-10"
                        :src="provider.icon"
                        :alt="provider.name"
                        placeholder-class="size-10"
                        loading="lazy"
                    />
                    <UIcon
                        v-else
                        class="size-10"
                        :name="provider.icon"
                        :style="provider.name === 'discord' && { color: '#7289da' }"
                    />

                    <div class="flex items-center gap-1.5 text-base font-semibold text-gray-900 dark:text-white">
                        {{ provider.label }}
                    </div>
                </UButton>

                <div v-if="account" class="flex items-center justify-between">
                    <UButton
                        icon="i-mdi-close-circle"
                        color="error"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => disconnectOAuth2Connection(provider),
                            })
                        "
                    >
                        {{ $t('common.disconnect') }}
                    </UButton>
                </div>

                <OAuth2ConnectButton v-if="!account && !nuiEnabled" :provider="provider" />
            </div>
        </template>

        <template v-if="account" #footer>
            <div class="inline-flex items-center gap-4">
                <template v-if="account">
                    <UAvatar :as="NuxtImg" size="md" :src="account.avatar" :alt="$t('common.image')" loading="lazy" />

                    <UTooltip :text="`ID: ${account.externalId}`">
                        <span class="text-left">
                            {{ account.username }}
                        </span>
                    </UTooltip>
                </template>

                <NotSupportedTabletBlock v-else-if="nuiEnabled" class="text-sm" />
            </div>
        </template>
    </UPageCard>
</template>
