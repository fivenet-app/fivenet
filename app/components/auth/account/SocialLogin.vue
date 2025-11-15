<script lang="ts" setup>
import { NuxtImg } from '#components';
import SocialLoginConnectButton from '~/components/auth/account/SocialLoginConnectButton.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useSettingsStore } from '~/stores/settings';
import { getAuthAuthClient } from '~~/gen/ts/clients';
import type { OAuth2Account, OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

defineProps<{
    provider: OAuth2Provider;
    account?: OAuth2Account;
}>();

const emit = defineEmits<{
    (e: 'disconnected', provider: string): void;
}>();

const notifications = useNotificationsStore();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const authAuthClient = await getAuthAuthClient();

async function disconnectSocialLogin(provider: OAuth2Provider): Promise<void> {
    try {
        await authAuthClient.deleteSocialLogin({
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

const overlay = useOverlay();
const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <UCard
        :ui="{
            header: 'flex flex-col',
            body: 'flex-1 flex flex-col',
        }"
    >
        <template #header>
            <div class="flex flex-1 gap-2">
                <div class="inline-flex flex-1 gap-2">
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

                    <div class="flex items-center gap-1.5 truncate text-base font-semibold text-highlighted">
                        {{ provider.label }}
                    </div>
                </div>

                <div v-if="account" class="flex items-center justify-between">
                    <UButton
                        icon="i-mdi-close-circle"
                        color="error"
                        @click="
                            confirmModal.open({
                                confirm: async () => disconnectSocialLogin(provider),
                            })
                        "
                    >
                        {{ $t('common.disconnect') }}
                    </UButton>
                </div>

                <SocialLoginConnectButton v-else-if="!nuiEnabled" :provider="provider" />
            </div>
        </template>

        <div v-if="account || nuiEnabled" class="flex flex-1 flex-col items-center justify-center gap-4">
            <template v-if="account">
                <div v-if="account" class="inline-flex items-center gap-2">
                    <UAvatar size="xl" :src="account.avatar" :alt="$t('common.image')" loading="lazy" />

                    <UTooltip :text="`${$t('components.auth.SocialLogins.external_id')}: ${account.externalId}`">
                        <span class="text-left">
                            {{ account.username }}
                        </span>
                    </UTooltip>
                </div>
            </template>

            <NotSupportedTabletBlock v-else-if="nuiEnabled" class="text-sm" />
        </div>

        <template #footer>
            <UButton
                size="xs"
                variant="link"
                color="neutral"
                :label="$t('components.auth.SocialLogins.connection_website')"
                external
                :to="provider.homepage"
                target="_blank"
                trailing-icon="i-mdi-external-link"
            />
        </template>
    </UCard>
</template>
